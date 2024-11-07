package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"myapp/internal/adapters/repository"
)

// Mock Redis Client
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return redis.NewStringResult(args.String(0), args.Error(1))
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return redis.NewStatusResult(args.Error(0))
}

func (m *MockRedisClient) Exists(ctx context.Context, key string) *redis.IntCmd {
	args := m.Called(ctx, key)
	return redis.NewIntResult(int64(args.Int(0)), args.Error(1))
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return redis.NewIntResult(int64(args.Int(0)), args.Error(1))
}

func TestRedisRepository_Set(t *testing.T) {
	mockClient := new(MockRedisClient)
	repo := repository.NewRedisRepository(mockClient)

	ctx := context.Background()
	key := "testKey"
	value := "testValue"
	expiration := time.Hour

	mockClient.On("Set", ctx, key, value, expiration).Return(nil)

	err := repo.Set(ctx, key, value, expiration)
	assert.NoError(t, err)

	mockClient.AssertExpectations(t)
}

func TestRedisRepository_Get(t *testing.T) {
	mockClient := new(MockRedisClient)
	repo := repository.NewRedisRepository(mockClient)

	ctx := context.Background()
	key := "testKey"

	mockClient.On("Get", ctx, key).Return("testValue", nil)

	result, err := repo.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, "testValue", result)

	mockClient.AssertExpectations(t)
}

func TestRedisRepository_Delete(t *testing.T) {
	mockClient := new(MockRedisClient)
	repo := repository.NewRedisRepository(mockClient)

	ctx := context.Background()
	key := "testKey"

	mockClient.On("Del", ctx, key).Return(nil)

	err := repo.Delete(ctx, key)
	assert.NoError(t, err)

	mockClient.AssertExpectations(t)
}

func TestRedisRepository_Exists(t *testing.T) {
	mockClient := new(MockRedisClient)
	repo := repository.NewRedisRepository(mockClient)

	ctx := context.Background()
	key := "testKey"

	mockClient.On("Exists", ctx, key).Return(1, nil)

	exists, err := repo.Exists(ctx, key)
	assert.NoError(t, err)
	assert.True(t, exists)

	mockClient.AssertExpectations(t)
}