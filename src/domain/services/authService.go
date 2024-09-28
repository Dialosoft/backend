package services

import (
	"fmt"

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
	RefreshToken(refreshToken string) (string, error)

	// GetRoleInformationByRoleID retrieves role information based on the provided role ID.
	// Returns the role information as a string and an error if the retrieval fails.
	GetRoleInformationByRoleID(roleID string) (string, error)
}

type authServiceImpl struct {
	userRepository  repository.UserRepository
	roleRepository  repository.RoleRepository
	tokenRepository repository.TokenRepository
	cacheService    CacheService
	jwtKey          string
}

// GetRoleInformationByRoleID implements AuthService.
func (service *authServiceImpl) GetRoleInformationByRoleID(roleID string) (string, error) {
	panic("unimplemented")
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

	fmt.Println(userEntity)

	refreshToken, err := service.getOrSaveRefreshToken(userID)
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

	refreshToken, err := service.getOrSaveRefreshToken(userEntity.ID)
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

	if service.cacheService.IsTokenBlacklisted(refreshToken) {
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

	userEntity, err = service.cacheService.GetUserInfoByID(userUUID)
	if err != nil {
		userEntity, err = service.userRepository.FindByID(userUUID)
		if err != nil {
			return "", err
		}
		err = service.cacheService.SetUserInfoByID(userUUID, userEntity)
		if err != nil {
			return "", err
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

func (service *authServiceImpl) getOrSaveRefreshToken(userID uuid.UUID) (string, error) {
	var refreshToken string

	refreshToken, err := service.cacheService.GetRefreshTokenByID(userID)
	if err != nil {
		logger.CaptureError(err, "Error in cache: GetRefreshTokenByID", map[string]interface{}{
			"userID": userID,
		})
	}

	if refreshToken == "" {
		tokenEntity, err := service.tokenRepository.FindTokenByUserID(userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				_, newTokenEntity, err := jsonWebToken.GenerateRefreshToken(service.jwtKey, userID)
				if err != nil {
					return "", err
				}

				err = service.tokenRepository.Save(newTokenEntity)
				if err != nil {
					return "", err
				}

				err = service.cacheService.SetRefreshTokenByID(userID, newTokenEntity.Token)
				if err != nil {
					return "", err
				}

				refreshToken = newTokenEntity.Token

			} else {
				return "", err
			}
		} else {
			err = service.cacheService.SetRefreshTokenByID(userID, tokenEntity.Token) // 5 days
			if err != nil {
				return "", err
			}
			refreshToken = tokenEntity.Token
		}
	} else {
		claims, validatErr := jsonWebToken.ValidateJWT(refreshToken, service.jwtKey)
		if err != nil {
		}

		if service.cacheService.IsTokenBlacklisted(refreshToken) || validatErr != nil {
			_, newTokenEntity, err := jsonWebToken.GenerateRefreshToken(service.jwtKey, userID)
			if err != nil {
				return "", err
			}

			jti, ok := claims["jti"].(string)
			if !ok {
				return "", errorsUtils.ErrInternalServer
			}

			jtiUUID, err := uuid.Parse(jti)
			if err != nil {
				return "", errorsUtils.ErrInvalidUUID
			}

			{
				//delete

				err = service.tokenRepository.Delete(jtiUUID)
				if err != nil {
					return "", err
				}

				err = service.cacheService.DeleteRefreshTokenByID(userID)
				if err != nil {
					return "", err
				}
			}

			err = service.tokenRepository.Save(newTokenEntity)
			if err != nil {
				return "", err
			}

			err = service.cacheService.SetRefreshTokenByID(userID, newTokenEntity.Token) // 5 days
			if err != nil {
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
	cacheService CacheService,
	jwtKey string) AuthService {
	return &authServiceImpl{
		userRepository:  userRepository,
		roleRepository:  roleRepository,
		tokenRepository: tokenRepository,
		cacheService:    cacheService,
		jwtKey:          jwtKey}
}
