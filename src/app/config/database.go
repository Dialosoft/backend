package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Database string
	Port     int
	SSLMode  bool
}

func GetNewDatabaseConfig() DatabaseConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Println("failed to get the port")
		return DatabaseConfig{}
	}
	var SSLMode bool
	if os.Getenv("SSLMODE") == "true" {
		SSLMode = true
	} else {
		SSLMode = false
	}

	return DatabaseConfig{
		Host:     os.Getenv("HOST"),
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
		Database: os.Getenv("DATABASE"),
		Port:     port,
		SSLMode:  SSLMode,
	}
}
