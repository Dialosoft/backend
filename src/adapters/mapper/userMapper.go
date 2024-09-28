package mapper

import (
	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/models"
)

// UserDtoToUserEntity returns a new UserEntity based on a UserDto, filling in the data.
// It returns an error if any of the required fields (Username, Email, or Password) are empty.
func UserDtoToUserEntity(userDto *dto.UserDto) *models.UserEntity {
	userEntity := models.UserEntity{
		Username: userDto.Username,
		Email:    userDto.Email,
		Password: userDto.Password,
	}

	return &userEntity
}

// UserEntityToUserDto converts a UserEntity to a UserDto.
// Returns an error if the UserEntity has missing required fields (Username or Email).
// The Password field is intentionally left blank in the resulting UserDto.
func UserEntityToUserDto(userEntity *models.UserEntity) *dto.UserDto {
	userDto := dto.UserDto{
		Username: userEntity.Username,
		Password: "",
		Email:    userEntity.Email,
	}

	return &userDto
}

func UserRequestToUserDto(userRequest *request.UserRequest) *dto.UserDto {
	userDto := dto.UserDto{
		Username: userRequest.Username,
		Password: userRequest.Password,
		Email:    userRequest.Email,
	}

	return &userDto
}

func UserUpdateRequestToUserDto(userRequest *request.UpdateUserRequest) *dto.UserDto {
	userDto := dto.UserDto{
		Username: userRequest.Username,
		Email:    userRequest.Email,
	}

	return &userDto
}

func UserEntityToUserResponse(userEntity *models.UserEntity) response.UserResponse {
	return response.UserResponse{
		ID:          userEntity.ID,
		Username:    userEntity.Username,
		Email:       userEntity.Email,
		Name:        userEntity.Name,
		Description: userEntity.Description,
		Banned:      userEntity.Banned,
		Role:        RoleEntityToRoleResponse(&userEntity.Role),
		CreatedAt:   userEntity.CreatedAt,
		UpdatedAt:   userEntity.UpdatedAt,
		DeletedAt:   userEntity.DeletedAt,
	}
}

func UserResponseToUserEntity(userResponse *response.UserResponse) *models.UserEntity {
	return &models.UserEntity{
		ID:          userResponse.ID,
		Username:    userResponse.Username,
		Email:       userResponse.Email,
		Name:        userResponse.Name,
		Description: userResponse.Description,
		Banned:      userResponse.Banned,
		RoleID:      userResponse.Role.ID,
		Role:        *RoleResponseToRoleEntity(&userResponse.Role),
		CreatedAt:   userResponse.CreatedAt,
		UpdatedAt:   userResponse.UpdatedAt,
		DeletedAt:   userResponse.DeletedAt,
	}
}
