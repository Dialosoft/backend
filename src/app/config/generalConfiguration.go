package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type GeneralConfig struct {
	Host         string
	User         string
	Password     string
	Database     string
	Port         int
	SSLMode      string
	SMTPHost     string
	SMTPPort     string
	MailUsername string
	MailPassword string
	FromAddress  string
	JWTKey       string
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

	jwtKey := os.Getenv("JWTKEY")
	if jwtKey == "" {
		log.Fatal("JWT key is missing")
	}

	return GeneralConfig{
		Host:         os.Getenv("HOST"),
		User:         os.Getenv("USER"),
		Password:     os.Getenv("PASSWORD"),
		Database:     os.Getenv("DATABASE"),
		SMTPHost:     os.Getenv("SMTPHOST"),
		SMTPPort:     os.Getenv("SMTPPORT"),
		MailUsername: os.Getenv("MAILUSERNAME"),
		MailPassword: os.Getenv("MAILPASSWORD"),
		FromAddress:  os.Getenv("FROMADDRESS"),
		Port:         port,
		SSLMode:      SSLMode,
		JWTKey:       jwtKey,
	}
}
