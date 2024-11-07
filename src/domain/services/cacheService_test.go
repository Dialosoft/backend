package services

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/models"
	"github.com/Dialosoft/src/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// RedisRepository Mock
type MockRedisRepository struct {
	mock.Mock
}

func (m *MockRedisRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisRepository) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisRepository) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisRepository) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

// Test InvalidateRefreshToken
func TestCacheService_InvalidateRefreshToken(t *testing.T) {
	mockRepo := new(MockRedisRepository)
	service := services.NewCacheService(mockRepo)

	t.Run("success", func(t *testing.T) {
		token := "testToken"
		cacheKey := "blacklist:" + token

		mockRepo.On("Set", mock.Anything, cacheKey, "true", time.Hour*720).Return(nil)

		err := service.InvalidateRefreshToken(token)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error setting token in cache", func(t *testing.T) {
		token := "testToken"
		cacheKey := "blacklist:" + token

		mockRepo.On("Set", mock.Anything, cacheKey, "true", time.Hour*720).Return(errors.New("cache error"))

		err := service.InvalidateRefreshToken(token)

		assert.EqualError(t, err, "cache error")
	})
}

// Test IsTokenBlacklisted
func TestCacheService_IsTokenBlacklisted(t *testing.T) {
	mockRepo := new(MockRedisRepository)
	service := services.NewCacheService(mockRepo)

	t.Run("token is blacklisted", func(t *testing.T) {
		token := "testToken"
		cacheKey := "blacklist:" + token

		mockRepo.On("Exists", mock.Anything, cacheKey).Return(true, nil)

		isBlacklisted := service.IsTokenBlacklisted(token)

		assert.True(t, isBlacklisted)
	})

	t.Run("token is not blacklisted", func(t *testing.T) {
		token := "testToken"
		cacheKey := "blacklist:" + token

		mockRepo.On("Exists", mock.Anything, cacheKey).Return(false, nil)

		isBlacklisted := service.IsTokenBlacklisted(token)

		assert.False(t, isBlacklisted)
	})

	t.Run("error checking blacklist status", func(t *testing.T) {
		token := "testToken"
		cacheKey := "blacklist:" + token

		mockRepo.On("Exists", mock.Anything, cacheKey).Return(false, errors.New("cache error"))

		isBlacklisted := service.IsTokenBlacklisted(token)

		assert.False(t, isBlacklisted)
	})
}

// Test SetUserInfoByID
func TestCacheService_SetUserInfoByID(t *testing.T) {
	mockRepo := new(MockRedisRepository)
	service := services.NewCacheService(mockRepo)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		userEntity := &models.UserEntity{ID: userID, Username: "testuser"}
		cacheKey := "user:" + userID.String()

		userEntity.Password = "" // Simular la eliminación de la contraseña antes de almacenar en caché
		userJSON, _ := json.Marshal(userEntity)

		mockRepo.On("Set", mock.Anything, cacheKey, string(userJSON), time.Hour*24).Return(nil)

		err := service.SetUserInfoByID(userID, userEntity)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error marshalling", func(t *testing.T) {
		userID := uuid.New()
		userEntity := &models.UserEntity{ID: userID}

		mockRepo.On("Set", mock.Anything).Return(nil).Maybe() // No se espera que se llame

		err := service.SetUserInfoByID(userID, userEntity)

		assert.NoError(t, err)
	})
}

// Test GetUserInfoByID
func TestCacheService_GetUserInfoByID(t *testing.T) {
	mockRepo := new(MockRedisRepository)
	service := services.NewCacheService(mockRepo)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		cacheKey := "user:" + userID.String()
		userEntity := &models.UserEntity{ID: userID, Username: "testuser"}

		userJSON, _ := json.Marshal(userEntity)

		mockRepo.On("Get", mock.Anything, cacheKey).Return(string(userJSON), nil)

		result, err := service.GetUserInfoByID(userID)

		assert.NoError(t, err)
		assert.Equal(t, userEntity.Username, result.Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error retrieving from cache", func(t *testing.T) {
		userID := uuid.New()
		cacheKey := "user:" + userID.String()

		mockRepo.On("Get", mock.Anything, cacheKey).Return("", errors.New("cache miss"))

		result, err := service.GetUserInfoByID(userID)

		assert.Nil(t, result)
		assert.EqualError(t, err, "cache miss")
	})

	t.Run("error unmarshalling data", func(t *testing.T) {
		userID := uuid.New()
		cacheKey := "user:" + userID.String()

		mockRepo.On("Get", mock.Anything, cacheKey).Return("invalid json", nil)

		result, err := service.GetUserInfoByID(userID)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

// Test SetRefreshTokenByID
func TestCacheService_SetRefreshTokenByID(t *testing.T) {
	mockRepo := new(MockRedisRepository)
	service := services.NewCacheService(mockRepo)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		token := "refreshToken"
		cacheKey := "refreshToken:" + userID.String()

		mockRepo.On("Set", mock.Anything, cacheKey, token, time.Hour*120).Return(nil)

		err := service.SetRefreshTokenByID(userID, token)

		assert.NoError(t, err)
	})

	t.Run("error setting refresh token in cache", func(t *testing.T) {
		userID := uuid.New()
		token := "refreshToken"
		cacheKey := "refreshToken:" + userID.String()

		mockRepo.On("Set", mock.Anything, cacheKey, token, time.Hour*120).Return(errors.New("cache error"))

		err := service.SetRefreshTokenByID(userID, token)

		assert.EqualError(t, err, "cache error")
	})
}

// Test GetRefreshTokenByID
func TestCacheService_GetRefreshTokenByID(t *testing.T) {
	mockRepo := new(MockRedisRepository)
	service := services.NewCacheService(mockRepo)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		cacheKey := "refreshToken:" + userID.String()

		mockRepo.On("Get", mock.AnythingOfType("*context.emptyCtx"), cacheKey).Return("refresh-token-value", nil)

		tokenValue, err := service.GetRefreshTokenByID(userID)

		assert.NoError(t, err)
		assert.Equal(t, "refresh-token-value", tokenValue)
	})

	t.Run("error retrieving refresh token from cache", func(t *testing.T) {
		userID := uuid.New()
		cacheKey := "refreshToken:" + userID.String()

		mockRepo.On("Get", mock.AnythingOfType("*context.emptyCtx"), cacheKey).Return("", errors.New("cache miss"))

		tokenValue, err := service.GetRefreshTokenByID(userID)

		assert.Empty(t, tokenValue)
		assert.EqualError(t, err, "cache miss")
	})
}

// Test DeleteRefreshTokenByID
func TestCacheService_DeleteRefreshTokenByID(t *testing.T) {
	mockRepo := new(MockRedisRepository)
	service := services.NewCacheService(mockRepo)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		cacheKey := "refreshToken:" + userID.String()

		mockRepo.On("Delete", mock.AnythingOfType("*context.emptyCtx"), cacheKey).Return(nil)

		err := service.DeleteRefreshTokenByID(userID)

		assert.NoError(t, err)
	})

	t.Run("error deleting refresh token from cache", func(t *testing.T) {
		userID := uuid.New()
		cacheKey := "refreshToken:" + userID.String()

		mockRepo.On("Delete", mock.AnythingOfType("*context.emptyCtx"), cacheKey).Return(errors.New("cache error"))

		err := service.DeleteRefreshTokenByID(userID)

		assert.EqualError(t, err, "cache error")
	})
}
