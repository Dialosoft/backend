package database

import (
	"fmt"

	"github.com/Dialosoft/src/app/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase(conf config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%t",
		conf.Host, conf.User, conf.Password, conf.Database, conf.Port, conf.SSLMode)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
