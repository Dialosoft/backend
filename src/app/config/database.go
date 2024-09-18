package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type GeneralConfig struct {
	Host     string
	User     string
	Password string
	Database string
	Port     int
	SSLMode  string
	JWTKey   string
}

func GetGeneralConfig() GeneralConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Println("failed to get the port")
		return GeneralConfig{}
	}
	var SSLMode string
	if os.Getenv("SSLMODE") == "enable" {
		SSLMode = "enable"
	} else {
		SSLMode = "disable"
	}

	return GeneralConfig{
		Host:     os.Getenv("HOST"),
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
		Database: os.Getenv("DATABASE"),
		Port:     port,
		SSLMode:  SSLMode,
		JWTKey:   os.Getenv("JWTKEY"),
	}
}
