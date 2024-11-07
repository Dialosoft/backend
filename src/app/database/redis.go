package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var redisCtx = context.Background()

func NewRedisClient() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST") // Get the Redis host from environment variables
	redisPort := os.Getenv("REDIS_PORT") // Get the Redis port from environment variables

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(redisCtx).Result()
	if err != nil {
		log.Fatal("cannot connect to redis")
	}

	return client
}
