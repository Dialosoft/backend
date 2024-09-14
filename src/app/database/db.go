package database

import (
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

	db.AutoMigrate(models.UserEntity{}, models.RoleEntity{})

	return db, nil
}
