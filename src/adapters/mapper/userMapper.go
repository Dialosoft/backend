package mapper

import (
	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/domain/models"
	"github.com/Dialosoft/src/pkg/errorsUtils"
)

// UserDtoToUserEntity returns a new UserEntity based on a UserDto, filling in the data.
// It returns an error if any of the required fields (Username, Email, or Password) are empty.
func UserDtoToUserEntity(userDto *dto.UserDto) (*models.UserEntity, error) {
	if userDto.Username == "" ||
		userDto.Email == "" ||
		userDto.Password == "" {
		return nil, errorsUtils.ErrParameterCannotBeNull
	}

	userEntity := models.UserEntity{
		Username: userDto.Username,
		Email:    userDto.Email,
		Password: userDto.Password,
	}

	return &userEntity, nil
}

// UserEntityToUserDto converts a UserEntity to a UserDto.
// Returns an error if the UserEntity has missing required fields (Username or Email).
// The Password field is intentionally left blank in the resulting UserDto.
func UserEntityToUserDto(userEntity *models.UserEntity) (*dto.UserDto, error) {
	if userEntity.Username == "" ||
		userEntity.Email == "" {
		return nil, errorsUtils.ErrParameterCannotBeNull
	}

	userDto := dto.UserDto{
		Username: userEntity.Username,
		Password: "",
		Email:    userEntity.Email,
	}

	return &userDto, nil
}
