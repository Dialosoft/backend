package database

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"

	"github.com/Dialosoft/src/app/config"
)


var redisCtx = context.Background()

func NewRedisClient(config config.GeneralConfig) *redis.Client {

    client := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
        Password: config.Redis.Password,
        DB:       config.Redis.DB,
    })

	_, err := client.Ping(redisCtx).Result()
	if err != nil {
		log.Fatal("cannot connect to redis")
	}

	return client
}