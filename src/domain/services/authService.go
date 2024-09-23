package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/mapper"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/models"
	"github.com/Dialosoft/src/pkg/errorsUtils"
	"github.com/Dialosoft/src/pkg/utils/jsonWebToken"
	"github.com/Dialosoft/src/pkg/utils/logger"
	"github.com/Dialosoft/src/pkg/utils/security"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthService defines a set of methods for managing authentication processes,
// such as user registration, login, token refresh, and token invalidation.
type AuthService interface {

	// Register registers a new user based on the provided UserDto.
	// Returns the UUID of the created user, an access token, a refresh token, and an error if the operation fails.
	Register(user dto.UserDto) (uuid.UUID, string, string, error)

	// Login authenticates a user with the provided username and password.
	// Returns an access token, a refresh token, and an error if authentication fails.
	Login(username, password string) (string, string, error)

	// RefreshToken refreshes the access token using the provided refresh token.
	// Returns a new access token and an error if the operation fails.
	RefreshToken(token string) (string, error)

	// InvalidateRefreshToken invalidates a refresh token, preventing it from being used again.
	// Returns an error if the invalidation fails.
	InvalidateRefreshToken(token string) error

	// IsTokenBlacklisted checks if the provided token has been blacklisted.
	// Returns true if the token is blacklisted, false otherwise.
	IsTokenBlacklisted(token string) bool

	// GetRoleInformationByRoleID retrieves role information based on the provided role ID.
	// Returns the role information as a string and an error if the retrieval fails.
	GetRoleInformationByRoleID(roleID string) (string, error)
}

type authServiceImpl struct {
	userRepository  repository.UserRepository
	roleRepository  repository.RoleRepository
	tokenRepository repository.TokenRepository
	cacheRepository repository.RedisRepository
	jwtKey          string
}

// Register implements AuthService.
func (service *authServiceImpl) Register(user dto.UserDto) (uuid.UUID, string, string, error) {
	roleEntity, err := service.roleRepository.FindByType("user")
	if err != nil {
		return uuid.UUID{}, "", "", err
	}

	passwordHashed, err := security.HashPassword(user.Password)
	if err != nil {
		return uuid.UUID{}, "", "", err
	}
	user.Password = passwordHashed

	userEntity := mapper.UserDtoToUserEntity(&user)

	userEntity.RoleID = roleEntity.ID
	userEntity.Role = *roleEntity

	userID, err := service.userRepository.Create(*userEntity)
	if err != nil {
		return uuid.UUID{}, "", "", err
	}

	token, err := jsonWebToken.GenerateAccessJWT(service.jwtKey, userID, userEntity.RoleID)
	if err != nil {
		return uuid.UUID{}, "", "", err
	}

	refreshToken, err := service.getOrSaveRefreshToken(*userEntity)
	if err != nil {
		return uuid.UUID{}, "", "", err
	}

	return userID, token, refreshToken, nil
}

// Login implements AuthService.
func (service *authServiceImpl) Login(username string, password string) (string, string, error) {

	userEntity, err := service.userRepository.FindByUsername(username)
	if err != nil {
		return "", "", err
	}

	if !security.CheckPasswordHash(password, userEntity.Password) {
		return "", "", errorsUtils.ErrUnauthorizedAcces
	}

	refreshToken, err := service.getOrSaveRefreshToken(*userEntity)
	if err != nil {
		return "", "", err
	}

	accesToken, err := jsonWebToken.GenerateAccessJWT(service.jwtKey, userEntity.ID, userEntity.RoleID)
	if err != nil {
		return "", "", nil
	}

	return accesToken, refreshToken, nil
}

// RefreshToken implements AuthService.
func (service *authServiceImpl) RefreshToken(refreshToken string) (string, error) {
	var userEntity *models.UserEntity
	claims, err := jsonWebToken.ValidateJWT(refreshToken, service.jwtKey)
	if err != nil {
		return "", errorsUtils.ErrRefreshTokenExpiredOrInvalid
	}

	if service.IsTokenBlacklisted(refreshToken) {
		return "", errorsUtils.ErrUnauthorizedAcces
	}

	userID, err := claims.GetSubject()
	if err != nil {
		return "", err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return "", errorsUtils.ErrInvalidUUID
	}

	cacheKey := fmt.Sprintf("user:%s", userUUID.String())
	userData, err := service.cacheRepository.Get(context.Background(), cacheKey)
	if err != nil {
		userEntity, err = service.userRepository.FindByID(userUUID)
		if err != nil {
			return "", errorsUtils.ErrNotFound
		}
		userEntity.Password = "" // delete the password hash, no needed
		json, err := json.Marshal(userEntity)
		if err != nil {
			return "", errorsUtils.ErrInternalServer
		}

		err = service.cacheRepository.Set(context.Background(), cacheKey, string(json), time.Hour*24)
		if err != nil {
			return "", errorsUtils.ErrInternalServer
		}
	} else {
		err = json.Unmarshal([]byte(userData), &userEntity)
		if err != nil {
			return "", errorsUtils.ErrInternalServer
		}
	}

	if userEntity.Locked {
		return "", errorsUtils.ErrUnauthorizedAcces
	}
	if userEntity.DeletedAt.Valid {
		return "", gorm.ErrRecordNotFound
	}
	if userEntity.Disable {
		return "", errorsUtils.ErrUnauthorizedAcces
	}

	accessToken, err := jsonWebToken.GenerateAccessJWT(service.jwtKey, userUUID, userEntity.RoleID)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// InvalidateToken implements AuthService.
func (service *authServiceImpl) InvalidateRefreshToken(token string) error {
	expiration := time.Hour * 720
	cacheKey := fmt.Sprintf("blacklist:%s", token)
	err := service.cacheRepository.Set(context.Background(), cacheKey, "true", expiration)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (service *authServiceImpl) IsTokenBlacklisted(token string) bool {
	cacheKey := fmt.Sprintf("blacklist:%s", token)
	is, err := service.cacheRepository.Exists(context.Background(), cacheKey)
	if err != nil {
		logger.Error(err.Error())
		return false
	}
	return is
}

func (service *authServiceImpl) GetRoleInformationByRoleID(roleID string) (string, error) {
	var roleModel *models.RoleEntity
	roleUUID, err := uuid.Parse(roleID)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	key, err := service.cacheRepository.Get(context.Background(), roleID)
	if err != nil || key == "" {
		roleModel, err = service.roleRepository.FindByID(roleUUID)
		if err != nil {
			logger.Error(err.Error())
			return "", err
		}

		err = service.cacheRepository.Set(context.Background(), roleModel.ID.String(), "true", time.Hour*48)
		if err != nil {
			logger.Error(err.Error())
			return "", err
		}

		return roleModel.RoleType, nil
	}

	return key, nil
}

func (service *authServiceImpl) getOrSaveRefreshToken(userEntity models.UserEntity) (string, error) {
	var refreshToken string
	cacheKey := fmt.Sprintf("refreshToken:%s", userEntity.ID)

	refreshToken, err := service.cacheRepository.Get(context.Background(), cacheKey)
	if err != nil {
		logger.Error(err.Error())
	}

	if refreshToken == "" {
		tokenEntity, err := service.tokenRepository.FindTokenByUserID(userEntity.ID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				_, newTokenEntity, err := jsonWebToken.GenerateRefreshToken(service.jwtKey, userEntity.ID)
				if err != nil {
					logger.Error(err.Error())
					return "", err
				}

				err = service.tokenRepository.Save(newTokenEntity)
				if err != nil {
					logger.Error(err.Error())
					return "", err
				}

				err = service.cacheRepository.Set(context.Background(), cacheKey, newTokenEntity.Token, time.Hour*120) // 5 days
				if err != nil {
					logger.Error(err.Error())
					return "", err
				}

				refreshToken = newTokenEntity.Token

			} else {
				logger.Error(err.Error())
				return "", err
			}
		} else {
			err = service.cacheRepository.Set(context.Background(), cacheKey, tokenEntity.Token, time.Hour*120) // 5 days
			if err != nil {
				logger.Error(err.Error())
				return "", err
			}
			refreshToken = tokenEntity.Token
		}
	} else {
		claims, validatErr := jsonWebToken.ValidateJWT(refreshToken, service.jwtKey)
		if err != nil {
			logger.Error(err.Error())
		}

		if service.IsTokenBlacklisted(refreshToken) || validatErr != nil {
			_, newTokenEntity, err := jsonWebToken.GenerateRefreshToken(service.jwtKey, userEntity.ID)
			if err != nil {
				logger.Error(err.Error())
				return "", err
			}

			jti, ok := claims["jti"].(string)
			if !ok {
				return "", errorsUtils.ErrInternalServer
			}

			jtiUUID, err := uuid.Parse(jti)
			if err != nil {
				logger.Error(err.Error())
				return "", errorsUtils.ErrInvalidUUID
			}

			{
				//delete

				err = service.tokenRepository.Delete(jtiUUID)
				if err != nil {
					logger.Error(err.Error())
					return "", err
				}

				err = service.cacheRepository.Delete(context.Background(), cacheKey)
				if err != nil {
					logger.Error(err.Error())
				}
			}

			err = service.tokenRepository.Save(newTokenEntity)
			if err != nil {
				logger.Error(err.Error())
				return "", err
			}

			err = service.cacheRepository.Set(context.Background(), cacheKey, newTokenEntity.Token, time.Hour*120) // 5 days
			if err != nil {
				logger.Error(err.Error())
				return "", err
			}

			refreshToken = newTokenEntity.Token
		}
	}

	return refreshToken, nil
}

func NewAuthService(userRepository repository.UserRepository,
	roleRepository repository.RoleRepository,
	tokenRepository repository.TokenRepository,
	cacheRepository repository.RedisRepository,
	jwtKey string) AuthService {
	return &authServiceImpl{
		userRepository:  userRepository,
		roleRepository:  roleRepository,
		tokenRepository: tokenRepository,
		cacheRepository: cacheRepository,
		jwtKey:          jwtKey}
}
