package database

import (
	"errors"
	"fmt"

	"github.com/Dialosoft/src/app/config"
	"github.com/Dialosoft/src/domain/models"
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
			// Si el rol no existe, lo creamos y agregamos su ID al mapa
			if err := db.Create(&role).Error; err != nil {
				return nil, fmt.Errorf("failed to create role %s: %w", role.RoleType, err)
			}
			roleMap[role.RoleType] = role.ID
		} else {
			// Si el rol ya existe, agregamos su ID al mapa
			roleMap[existingRole.RoleType] = existingRole.ID
		}
	}

	return roleMap, nil
}
