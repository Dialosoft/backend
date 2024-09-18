package mapper

import (
	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/domain/models"
)

// UserDtoToUserEntity returns a new UserEntity based on a UserDto, filling in the data.
// It returns an error if any of the required fields (Username, Email, or Password) are empty.
func UserDtoToUserEntity(userDto *dto.UserDto) (*models.UserEntity, error) {
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
	userDto := dto.UserDto{
		Username: userEntity.Username,
		Password: "",
		Email:    userEntity.Email,
	}

	return &userDto, nil
}

func UserRequestToUserDto(userRequest *request.UserRequest) (*dto.UserDto, error) {
	userDto := dto.UserDto{
		Username: userRequest.Username,
		Password: userRequest.Password,
		Email:    userRequest.Email,
	}

	return &userDto, nil
}

func UserUpdateRequestToUserDto(userRequest *request.UpdateUserRequest) (*dto.UserDto, error) {
	userDto := dto.UserDto{
		Username: userRequest.Username,
		Email:    userRequest.Email,
	}

	return &userDto, nil
}
