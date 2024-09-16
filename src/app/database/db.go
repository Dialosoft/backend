package database

import (
	"errors"
	"fmt"

	"github.com/Dialosoft/src/app/config"
	"github.com/Dialosoft/src/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase(conf config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		conf.Host, conf.User, conf.Password, conf.Database, conf.Port, conf.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(models.UserEntity{}, models.RoleEntity{})
	if err != nil {
		return nil, err
	}

	err = createDefaultRoles(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createDefaultRoles(db *gorm.DB) error {
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
				return fmt.Errorf("failed to create role %s: %w", role.RoleType, err)
			}
		}
	}

	return nil
}
