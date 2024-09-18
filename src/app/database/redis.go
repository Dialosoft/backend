package database

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var redisCtx = context.Background()

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(redisCtx).Result()
	if err != nil {
		log.Fatal("cannot connect to redis")
	}

	return client
}
