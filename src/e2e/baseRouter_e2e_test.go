//go:build e2e

package e2e

import (
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/Dialosoft/src/app/config"
)

type TestServer struct {
	App         *fiber.App
	DB          *gorm.DB
	RedisClient *redis.Client
	Config      *config.GeneralConfig
}

func NewTestServer(db *gorm.DB, redisClient *redis.Client) *TestServer {
	// Create test configuration
	testConfig := config.GeneralConfig{
		Host:     "localhost",
		User:     "postgresql",
		Password: "postgresql",
		Database: "testdb",
		Port:     5432,
		SSLMode:  "disable",
		JWTKey:   "test-jwt-key",
	}

	// Initialize the Fiber app using the existing SetupAPI function
	app := config.SetupAPI(db, redisClient, testConfig)

	return &TestServer{
		App:         app,
		DB:          db,
		RedisClient: redisClient,
		Config:      &testConfig,
	}
}

