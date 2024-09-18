package services

import (
	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/mapper"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/pkg/errorsUtils"
	"github.com/Dialosoft/src/pkg/utils/jsonWebToken"
	"github.com/Dialosoft/src/pkg/utils/security"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(user dto.UserDto) (uuid.UUID, string, string, error)
	Login(username, password string) (string, string, error)
}

type authServiceImpl struct {
	userRepository  repository.UserRepository
	roleRepository  repository.RoleRepository
	tokenRepository repository.TokenRepository
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

	userEntity.ID = roleEntity.ID
	userEntity.Role = *roleEntity

	userID, err := service.userRepository.Create(*userEntity)
	if err != nil {
		return uuid.UUID{}, "", "", err
	}

	token, err := jsonWebToken.GenerateJWT(service.jwtKey, userID)
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
		}
		return "", "", err
	}

	accesToken, err := jsonWebToken.GenerateJWT(service.jwtKey, userEntity.ID)
	if err != nil {
		return "", "", nil
	}

	return accesToken, tokenEntity.Token, nil
}

func NewAuthService(userRepository repository.UserRepository, roleRepository repository.RoleRepository, tokenRepository repository.TokenRepository, jwtKey string) AuthService {
	return &authServiceImpl{userRepository: userRepository, roleRepository: roleRepository, tokenRepository: tokenRepository, jwtKey: jwtKey}
}
