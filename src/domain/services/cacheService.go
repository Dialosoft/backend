package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/models"
	"github.com/Dialosoft/src/pkg/utils/logger"
	"github.com/google/uuid"
)

type CacheService interface {
	InvalidateRefreshToken(token string) error
	IsTokenBlacklisted(token string) bool
	SetUserInfoByID(userID uuid.UUID, userEntity *models.UserEntity) error
	GetUserInfoByID(userID uuid.UUID) (*models.UserEntity, error)
	SetRefreshTokenByID(userID uuid.UUID, token string) error
	GetRefreshTokenByID(userID uuid.UUID) (string, error)
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
		logger.Error(err.Error())
		return err
	}

	return nil
}

// IsTokenBlacklisted implements CacheService.
func (service *cacheServiceImpl) IsTokenBlacklisted(token string) bool {
	cacheKey := fmt.Sprintf("blacklist:%s", token)
	is, err := service.cacheRepository.Exists(context.Background(), cacheKey)
	if err != nil {
		logger.Error(err.Error())
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
