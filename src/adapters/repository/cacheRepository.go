package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

type redisRepositoyryImpl struct {
	client *redis.Client
}

// Get implements RedisRepository.
func (r *redisRepositoyryImpl) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Set implements RedisRepository.
func (r *redisRepositoyryImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Exists implements RedisRepository.
func (r *redisRepositoyryImpl) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	return count > 0, err
}

// Delete implements RedisRepository.
func (r *redisRepositoyryImpl) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func NewRedisRepository(redisConn *redis.Client) RedisRepository {
	return &redisRepositoyryImpl{client: redisConn}
}
