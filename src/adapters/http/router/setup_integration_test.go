//go:build integration
package router

import (
	"context"
	"fmt"
	"time"

	rd "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Dialosoft/src/domain/models"
)


type IntegrationTestSuite struct {
    suite.Suite
    testCtx            context.Context               // Test execution context shared by all tests
    db                 *gorm.DB                      // Postgres database connection
    pgContainer        testcontainers.Container      // Reference to the created postgres docker container
    pgConnectionString string                        // Connection string for postgres database
    rdContainer        testcontainers.Container      // Reference to the created Redis container
    rdConnectionString string                        // Connection string for Redis server
    rdClient           *rd.Client                    // Redis client for communication
}

func (suite *IntegrationTestSuite) SetupSuite() {
    // Initialize test context
    suite.testCtx = context.Background()

    // Create PostgreSQL container request
    pgReq := testcontainers.ContainerRequest{
        Image:        "postgres:15.3-alpine",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_DB":       "database",
            "POSTGRES_USER":     "postgres",
            "POSTGRES_PASSWORD": "postgres",
        },
        WaitingFor: wait.ForLog("database system is ready to accept connections").
            WithOccurrence(2).
            WithStartupTimeout(5 * time.Second),
    }

    // Start PostgreSQL container
    pgContainer, err := testcontainers.GenericContainer(
        suite.testCtx,
        testcontainers.GenericContainerRequest{
            ContainerRequest: pgReq,
            Started:         true,
        },
    )
    suite.NoError(err)

    // Get host and port for PostgreSQL connection
    pgHost, err := pgContainer.Host(suite.testCtx)
    suite.NoError(err)
    pgPort, err := pgContainer.MappedPort(suite.testCtx, "5432")
    suite.NoError(err)

    // Build PostgreSQL connection string
    connStr := fmt.Sprintf("host=%s port=%s user=postgres password=postgres dbname=database sslmode=disable",
        pgHost, pgPort.Port())

    // Initialize GORM with PostgreSQL connection
    db, err := gorm.Open(pg.Open(connStr), &gorm.Config{})
    suite.NoError(err)

    // Store container and connection information
    suite.pgContainer = pgContainer
    suite.pgConnectionString = connStr
    suite.db = db

    // Verify PostgreSQL connection
    sqlDB, err := suite.db.DB()
    suite.NoError(err)
    err = sqlDB.Ping()
    suite.NoError(err)

    // Create Redis container request
    redisReq := testcontainers.ContainerRequest{
        Image:        "redis:6",
        ExposedPorts: []string{"6379/tcp"},
        WaitingFor:   wait.ForLog("Ready to accept connections"),
    }

    // Start Redis container
    redisContainer, err := testcontainers.GenericContainer(
        suite.testCtx,
        testcontainers.GenericContainerRequest{
            ContainerRequest: redisReq,
            Started:         true,
        },
    )
    suite.NoError(err)

    // Get host and port for Redis connection
    redisHost, err := redisContainer.Host(suite.testCtx)
    suite.NoError(err)
    redisPort, err := redisContainer.MappedPort(suite.testCtx, "6379")
    suite.NoError(err)

    // Build Redis connection string
    rdConnStr := fmt.Sprintf("redis://%s:%s", redisHost, redisPort.Port())

    // Parse Redis connection options
    rdConnOptions, err := rd.ParseURL(rdConnStr)
    suite.NoError(err)

    // Initialize Redis client
    rdClient := rd.NewClient(rdConnOptions)

    // Store Redis container and client information
    suite.rdContainer = redisContainer
    suite.rdConnectionString = rdConnStr
    suite.rdClient = rdClient

    // Verify Redis connection
    err = suite.rdClient.Ping(suite.testCtx).Err()
    suite.NoError(err)
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	if suite.pgContainer != nil {
		err := suite.pgContainer.Terminate(suite.ctx)
		suite.NoError(err)
		log.Println("PostgreSQL container terminated")
	}

	if suite.rdContainer != nil {
		err := suite.rdContainer.Terminate(suite.ctx)
		suite.NoError(err)
		log.Println("Redis container terminated")
	}
}



