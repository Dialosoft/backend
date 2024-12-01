//go:build e2e

package e2e

import (
	"context"
	"log"
	"time"

	rd "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"

	"github.com/Dialosoft/src/app/config"
	"github.com/Dialosoft/src/app/database"
)

// E2eTestSuite is a suite of e2e tests for the application.
type E2eTestSuite struct {
	suite.Suite
	testCtx     context.Context          // Test execution context shared by all tests
	db          *gorm.DB                 // Postgres database connection
	pgContainer testcontainers.Container // Reference to the created postgres docker container
	rdContainer testcontainers.Container // Reference to the created Redis container
	rdClient    *rd.Client               // Redis client for communication
}

// SetupSuite is called before any tests in the suite are run.
func (suite *E2eTestSuite) SetupSuite() {
	// Initialize test context
	suite.testCtx = context.Background()

	// Create PostgreSQL container request
	pgReq := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":               "testdb",
			"POSTGRES_USER":             "postgres",
			"POSTGRES_PASSWORD":         "postgres",
			"POSTGRES_HOST_AUTH_METHOD": "trust",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort("5432/tcp"),
		).WithDeadline(60 * time.Second),
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      "./../app/database/init-db.sql",
				ContainerFilePath: "/docker-entrypoint-initdb.d/init-db.sql",
				FileMode:          0o644,
			},
		},
	}

	// Start PostgreSQL container
	pgContainer, err := testcontainers.GenericContainer(
		suite.testCtx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: pgReq,
			Started:          true,
		},
	)
	suite.NoError(err)

	// Get host and port for PostgreSQL connection
	pgHost, err := pgContainer.Host(suite.testCtx)
	suite.NoError(err)
	pgPort, err := pgContainer.MappedPort(suite.testCtx, "5432")
	suite.NoError(err)

    	// Test configuration
	testConfig := config.GeneralConfig{
		Host:     pgHost,
		Port:     pgPort.Int(),
		User:     "postgres",
		Password: "postgres",
		Database: "testdb",
		SSLMode:  "disable",
		JWTKey:   "test-jwt-key",
	}

    	// Use Connecttodatabase to handle migrations and roles
	conn, err := database.ConnectToDatabase(testConfig)
	suite.NoError(err)

	// Use Connecttodatabase to handle migrations and roles
	suite.db = conn.Gorm
	suite.pgContainer = pgContainer

	// Verify PostgreSQL connection
	sqlDB, err := suite.db.DB()
	suite.NoError(err)
	err = sqlDB.Ping()
	suite.NoError(err)

	// Create Redis container request
	redisReq := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("* Ready to accept connections"),
			wait.ForListeningPort("6379/tcp"),
		).WithDeadline(60 * time.Second),
	}

	// Start Redis container
	redisContainer, err := testcontainers.GenericContainer(
		suite.testCtx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: redisReq,
			Started:          true,
		},
	)
	suite.NoError(err)

	redisHost, err := redisContainer.Host(suite.testCtx)
	suite.NoError(err)
	redisPort, err := redisContainer.MappedPort(suite.testCtx, "6379")
	suite.NoError(err)

	testConfig.Redis = &config.RedisConfig{
		Host:     redisHost,
		Port:     redisPort.Int(),
		Password: "",
		DB:       0,
	}

	// Initialize REDIS client using the new configuration
	rdClient := database.NewRedisClient(testConfig)

	suite.rdContainer = redisContainer
	suite.rdClient = rdClient

	// Verify Redis connection
	err = suite.rdClient.Ping(suite.testCtx).Err()
	suite.NoError(err)
}

// TearDownSuite is called after all tests in the suite have been run, regardless of whether they passed or failed.
func (suite *E2eTestSuite) TearDownSuite() {

    if suite.pgContainer != nil {
        err := suite.pgContainer.Terminate(suite.testCtx)
        suite.NoError(err)
        log.Println("PostgreSQL container terminated")
    }

    if suite.rdContainer != nil {
        err := suite.rdContainer.Terminate(suite.testCtx)
        suite.NoError(err)
        log.Println("Redis container terminated")
    }
}
