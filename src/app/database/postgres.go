package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Dialosoft/src/app/config"
	"github.com/Dialosoft/src/domain/models"
	"github.com/Dialosoft/src/pkg/utils/logger"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connection struct {
	Gorm            *gorm.DB
	DefaultRolesIDs map[string]uuid.UUID
}

func ConnectToDatabase(conf config.GeneralConfig) (Connection, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		conf.Host, conf.User, conf.Password, conf.Database, conf.Port, conf.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return Connection{}, err
	}

	err = db.AutoMigrate(
		models.UserEntity{},
		models.RoleEntity{},
		models.TokenEntity{},
		models.Category{},
		models.Forum{},
		models.Post{},
		models.Comment{},
		models.PostLikes{},
		models.CommentVotes{},
		models.RolePermissions{},
	)
	if err != nil {
		return Connection{}, err
	}

	defaultRoles, err := createDefaultRoles(db)
	if err != nil && err != gorm.ErrRecordNotFound {
		return Connection{}, err
	}

	return Connection{
		Gorm:            db,
		DefaultRolesIDs: defaultRoles,
	}, nil
}

func StartTokenChecker(ctx context.Context, db *gorm.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkOldAndBlockedTokens(db)
		case <-ctx.Done():
			log.Println("Stopping token checker...")
			return
		}
	}
}

func createDefaultRoles(db *gorm.DB) (map[string]uuid.UUID, error) {
	roleMap := make(map[string]uuid.UUID)

	roles := []models.RoleEntity{
		{RoleType: "user", Permission: 1, AdminRole: false, ModRole: false},
		{RoleType: "moderator", Permission: 2, AdminRole: false, ModRole: true},
		{RoleType: "administrator", Permission: 3, AdminRole: true, ModRole: false},
	}

	roleTypes := make([]string, len(roles))
	for i, role := range roles {
		roleTypes[i] = role.RoleType
	}

	var existingRoles []models.RoleEntity
	if err := db.Where("role_type IN ?", roleTypes).Find(&existingRoles).Error; err != nil {
		return nil, err
	}

	existingRoleMap := make(map[string]models.RoleEntity)
	for _, role := range existingRoles {
		existingRoleMap[role.RoleType] = role
		roleMap[role.RoleType] = role.ID
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		for _, role := range roles {
			if existingRole, exists := existingRoleMap[role.RoleType]; !exists {
				if err := tx.Create(&role).Error; err != nil {
					return fmt.Errorf("failed to create role %s: %w", role.RoleType, err)
				}
				roleMap[role.RoleType] = role.ID
			} else {
				role.ID = existingRole.ID
				if err := tx.Save(&role).Error; err != nil {
					return fmt.Errorf("failed to update role %s: %w", role.RoleType, err)
				}
				roleMap[role.RoleType] = existingRole.ID
			}

			rolePermissions := getRolePermissions(role.RoleType, role.ID)
			if err := tx.Where("role_id = ?", role.ID).Assign(rolePermissions).FirstOrCreate(&models.RolePermissions{}).Error; err != nil {
				return fmt.Errorf("failed to create or update role permissions for %s: %w", role.RoleType, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return roleMap, nil
}

func getRolePermissions(roleType string, roleID uuid.UUID) models.RolePermissions {
	switch roleType {
	case "user":
		return models.RolePermissions{
			RoleID:              roleID,
			CanManageCategories: false,
			CanManageForums:     false,
			CanManageRoles:      false,
			CanManageUsers:      false,
		}
	case "moderator":
		return models.RolePermissions{
			RoleID:              roleID,
			CanManageCategories: false,
			CanManageForums:     false,
			CanManageRoles:      false,
			CanManageUsers:      true,
		}
	case "administrator":
		return models.RolePermissions{
			RoleID:              roleID,
			CanManageCategories: true,
			CanManageForums:     true,
			CanManageRoles:      true,
			CanManageUsers:      true,
		}
	default:
		return models.RolePermissions{}
	}
}

func checkOldAndBlockedTokens(db *gorm.DB) {
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	var tokens []models.TokenEntity

	err := db.Where("blocked = ? OR created_at < ?", true, thirtyDaysAgo).Find(&tokens).Error
	if err != nil {
		logger.CaptureError(err, "error in checkOldAndBlockedTokens", nil)
		return
	}

	if len(tokens) > 0 {
		for _, token := range tokens {
			logger.Info(
				fmt.Sprintf("Deleting token ID: %d, Blocked: %v, Created At: %s", token.ID, token.Blocked, token.CreatedAt),
				map[string]interface{}{"tokenID": token.ID, "blocked": token.Blocked, "createdAt": token.CreatedAt},
			)
			if err := db.Delete(&token).Error; err != nil {
				logger.Error(
					"Error deleting token",
					map[string]interface{}{"tokenID": token.ID, "error": err.Error()},
				)
			} else {
				logger.Info(
					fmt.Sprintf("Successfully deleted token ID: %d", token.ID),
					map[string]interface{}{"tokenID": token.ID},
				)
			}
		}
	}
}
