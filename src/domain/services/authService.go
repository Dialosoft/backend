package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/mapper"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/models"
	"github.com/Dialosoft/src/pkg/errorsUtils"
	"github.com/Dialosoft/src/pkg/utils/jsonWebToken"
	"github.com/Dialosoft/src/pkg/utils/security"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(user dto.UserDto) (uuid.UUID, string, string, error)
	Login(username, password string) (string, string, error)
	RefreshToken(token string) (string, error)
	InvalidateRefreshToken(token string) error
	IsTokenBlacklisted(token string) bool
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

	refreshToken, tokenEntity, err := jsonWebToken.GenerateRefreshToken(service.jwtKey, userID)
	if err != nil {
		return uuid.UUID{}, "", "", err
	}

	service.tokenRepository.Save(tokenEntity)

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

	tokenEntity, err := service.tokenRepository.FindTokenByUserID(userEntity.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			_, newTokenEntity, err := jsonWebToken.GenerateRefreshToken(service.jwtKey, userEntity.ID)
			if err != nil {
				return "", "", err
			}

			err = service.tokenRepository.Save(newTokenEntity)
			if err != nil {
				return "", "", err
			}
			tokenEntity = &newTokenEntity
		} else {
			return "", "", err
		}
	}

	accesToken, err := jsonWebToken.GenerateAccessJWT(service.jwtKey, userEntity.ID, userEntity.RoleID)
	if err != nil {
		return "", "", nil
	}

	return accesToken, tokenEntity.Token, nil
}

// RefreshToken implements AuthService.
func (service *authServiceImpl) RefreshToken(refreshToken string) (string, error) {
	claims, err := jsonWebToken.ValidateJWT(refreshToken, service.jwtKey)
	if err != nil {
		return "", errorsUtils.ErrRefreshTokenExpiredOrInvalid
	}

	userID, err := claims.GetSubject()
	if err != nil {
		return "", err
	}

	var roleID string
	if roleIDClaim, ok := claims["roleID"].(string); ok {
		roleID = roleIDClaim
	}

	if roleID == "" {
		return "", errorsUtils.ErrRoleIDInRefreshToken
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return "", errorsUtils.ErrInvalidUUID
	}

	roleUUID, err := uuid.Parse(roleID)
	if err != nil {
		return "", errorsUtils.ErrInvalidUUID
	}

	accesToken, err := jsonWebToken.GenerateAccessJWT(service.jwtKey, userUUID, roleUUID)
	if err != nil {
		return "", err
	}

	return accesToken, nil
}

// InvalidateToken implements AuthService.
func (service *authServiceImpl) InvalidateRefreshToken(token string) error {
	expiration := time.Hour * 720
	return service.cacheRepository.Set(context.Background(), fmt.Sprintf("blacklist:%s", token), "true", expiration)
}

func (service *authServiceImpl) IsTokenBlacklisted(token string) bool {
	is, err := service.cacheRepository.Exists(context.Background(), fmt.Sprintf("blacklist:%s", token))
	if err != nil {
		log.Println("ERROR:", err)
		return false
	}
	return is
}

func (service *authServiceImpl) GetRoleInformationByRoleID(roleID string) (string, error) {
	var roleModel *models.RoleEntity
	roleUUID, err := uuid.Parse(roleID)
	if err != nil {
		return "", err
	}

	key, err := service.cacheRepository.Get(context.Background(), roleID)
	if err != nil || key == "" {
		roleModel, err = service.roleRepository.FindByID(roleUUID)
		if err != nil {
			return "", err
		}

		err = service.cacheRepository.Set(context.Background(), roleModel.ID.String(), "true", time.Hour*48)
		if err != nil {
			return "", err
		}

		return roleModel.RoleType, nil
	}

	return key, nil
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
