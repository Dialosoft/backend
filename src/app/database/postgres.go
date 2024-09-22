package database

import (
	"context"
	"errors"
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

	err = db.AutoMigrate(models.UserEntity{}, models.RoleEntity{}, models.TokenEntity{})
	if err != nil {
		return Connection{}, err
	}

	defaultRoles, err := createDefaultRoles(db)
	if err != nil {
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

	for _, role := range roles {
		var existingRole models.RoleEntity
		result := db.Where("role_type = ?", role.RoleType).First(&existingRole)

		if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
			if err := db.Create(&role).Error; err != nil {
				return nil, fmt.Errorf("failed to create role %s: %w", role.RoleType, err)
			}
			roleMap[role.RoleType] = role.ID
		} else {
			roleMap[existingRole.RoleType] = existingRole.ID
		}
	}

	return roleMap, nil
}

func checkOldAndBlockedTokens(db *gorm.DB) {
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	var tokens []models.TokenEntity

	err := db.Where("blocked = ? OR created_at < ?", true, thirtyDaysAgo).Find(&tokens).Error
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if len(tokens) > 0 {
		for _, token := range tokens {
			logger.Info(fmt.Sprintf("Deleting token ID: %d, Blocked: %v, Created At: %s\n", token.ID, token.Blocked, token.CreatedAt))
			if err := db.Delete(&token).Error; err != nil {
				logger.Error(fmt.Sprintf("Error deleting token ID: %d, error: %v\n", token.ID, err))
			} else {
				logger.Info(fmt.Sprintf("Successfully deleted token ID: %d", token.ID))
			}
		}
	}
}
