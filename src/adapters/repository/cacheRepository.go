package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {

	// Set sets a value in Redis with the given key and expiration time.
	// Accepts an interface{} as the value to store, allowing flexibility in the data type.
	// Returns an error if the operation fails.
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Get retrieves the value associated with the given key from Redis.
	// Returns the value as a string and an error if the key does not exist or there is a Redis error.
	Get(ctx context.Context, key string) (string, error)

	// Delete removes the given key from Redis.
	// Returns an error if the key does not exist or if there is a Redis error during deletion.
	Delete(ctx context.Context, key string) error

	// Exists checks if the given key exists in Redis.
	// Returns true if the key exists, false otherwise, along with an error if the operation fails.
	Exists(ctx context.Context, key string) (bool, error)
}

type redisRepositoyryImpl struct {
	client *redis.Client
}

func (r *redisRepositoyryImpl) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *redisRepositoyryImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisRepositoyryImpl) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	return count > 0, err
}

func (r *redisRepositoyryImpl) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func NewRedisRepository(redisConn *redis.Client) RedisRepository {
	return &redisRepositoyryImpl{client: redisConn}
}
