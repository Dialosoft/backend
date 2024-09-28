package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
)

// CacheService defines the methods for managing user-related data and tokens in a cache.
type CacheService interface {
	// InvalidateRefreshToken invalidates the given refresh token by marking it as unusable.
	InvalidateRefreshToken(token string) error

	// IsTokenBlacklisted checks if the given token has been blacklisted and is no longer valid.
	IsTokenBlacklisted(token string) bool

	// SetUserInfoByID stores user information in the cache associated with the given user ID.
	SetUserInfoByID(userID uuid.UUID, userEntity *models.UserEntity) error

	// GetUserInfoByID retrieves user information from the cache based on the provided user ID.
	// Returns the user entity or an error if not found.
	GetUserInfoByID(userID uuid.UUID) (*models.UserEntity, error)

	// SetRefreshTokenByID stores a refresh token in the cache, associated with the given user ID.
	SetRefreshTokenByID(userID uuid.UUID, token string) error

	// GetRefreshTokenByID retrieves a refresh token from the cache based on the provided user ID.
	// Returns the token or an error if not found.
	GetRefreshTokenByID(userID uuid.UUID) (string, error)

	// DeleteRefreshTokenByID removes the refresh token from the cache associated with the given user ID.
	DeleteRefreshTokenByID(userID uuid.UUID) error
}

type cacheServiceImpl struct {
	cacheRepository repository.RedisRepository
}

// GetUserInfoByID implements CacheService.
func (service *cacheServiceImpl) GetUserInfoByID(userID uuid.UUID) (*models.UserEntity, error) {
	var userEntity *models.UserEntity

	cacheKey := fmt.Sprintf("user:%s", userID.String())
	userString, err := service.cacheRepository.Get(context.Background(), cacheKey)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(userString), &userEntity)
	if err != nil {
		return nil, err
	}

	return userEntity, nil
}

// SetUserInfoByID implements CacheService.
func (service *cacheServiceImpl) SetUserInfoByID(userID uuid.UUID, userEntity *models.UserEntity) error {
	cacheKey := fmt.Sprintf("user:%s", userID.String())
	userEntity.Password = "" // delete the password hash, no needed

	json, err := json.Marshal(userEntity)
	if err != nil {
		return err
	}

	return service.cacheRepository.Set(context.Background(), cacheKey, string(json), time.Hour*24)
}

// InvalidateRefreshToken implements CacheService.
func (service *cacheServiceImpl) InvalidateRefreshToken(token string) error {
	expiration := time.Hour * 720
	cacheKey := fmt.Sprintf("blacklist:%s", token)
	err := service.cacheRepository.Set(context.Background(), cacheKey, "true", expiration)
	if err != nil {
		return err
	}

	return nil
}

// IsTokenBlacklisted implements CacheService.
func (service *cacheServiceImpl) IsTokenBlacklisted(token string) bool {
	cacheKey := fmt.Sprintf("blacklist:%s", token)
	is, err := service.cacheRepository.Exists(context.Background(), cacheKey)
	if err != nil {
		return false
	}
	return is
}

// SetRefreshTokenByID implements CacheService.
func (service *cacheServiceImpl) SetRefreshTokenByID(userID uuid.UUID, token string) error {
	cacheKey := fmt.Sprintf("refreshToken:%s", userID.String())
	return service.cacheRepository.Set(context.Background(), cacheKey, token, time.Hour*120) // 5 days
}

// GetRefreshTokenByID implements CacheService.
func (service *cacheServiceImpl) GetRefreshTokenByID(userID uuid.UUID) (string, error) {
	cacheKey := fmt.Sprintf("refreshToken:%s", userID.String())
	refreshToken, err := service.cacheRepository.Get(context.Background(), cacheKey)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (service *cacheServiceImpl) DeleteRefreshTokenByID(userID uuid.UUID) error {
	cacheKey := fmt.Sprintf("refreshToken:%s", userID.String())
	return service.cacheRepository.Delete(context.Background(), cacheKey)
}

func NewCacheService(cacheRepository repository.RedisRepository) CacheService {
	return &cacheServiceImpl{cacheRepository: cacheRepository}
}
