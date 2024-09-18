package services

import (
	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/mapper"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/google/uuid"
)

type UserService interface {
	GetAllUsers() ([]*dto.UserDto, error)
	GetUserByID(userID uuid.UUID) (*dto.UserDto, error)
	GetUserByUsername(username string) (*dto.UserDto, error)
	CreateNewUser(newUser dto.UserDto) (uuid.UUID, error)
	UpdateUser(userID uuid.UUID, updatedUser dto.UserDto) error
	DeleteUser(userID uuid.UUID) error
	RestoreUser(userID uuid.UUID) error
}

type userServiceImpl struct {
	repository     repository.UserRepository
	roleRepository repository.RoleRepository
}

// GetAllUsers implements UserService.
func (service *userServiceImpl) GetAllUsers() ([]*dto.UserDto, error) {
	var usersDtos []*dto.UserDto

	usersEntities, err := service.repository.FindAllUsers()
	if err != nil {
		return nil, err
	}

	for _, v := range usersEntities {
		userDto, err := mapper.UserEntityToUserDto(v)
		if err != nil {
			return nil, err
		} else {
			usersDtos = append(usersDtos, userDto)
		}
	}

	return usersDtos, nil
}

// GetUserByID implements UserService.
func (service *userServiceImpl) GetUserByID(userID uuid.UUID) (*dto.UserDto, error) {
	userEntity, err := service.repository.FindByID(userID)
	if err != nil {
		return nil, err
	}
	userDto, err := mapper.UserEntityToUserDto(userEntity)
	if err != nil {
		return nil, err
	}

	return userDto, nil
}

// GetUserByUsername implements UserService.
func (service *userServiceImpl) GetUserByUsername(username string) (*dto.UserDto, error) {
	userEntity, err := service.repository.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	userDto, err := mapper.UserEntityToUserDto(userEntity)
	if err != nil {
		return nil, err
	}

	return userDto, nil
}

// CreateNewUser implements UserService.
func (service *userServiceImpl) CreateNewUser(newUser dto.UserDto) (uuid.UUID, error) {

	roleEntity, err := service.roleRepository.FindByType("user")
	if err != nil {
		return uuid.UUID{}, err
	}

	userEntity, err := mapper.UserDtoToUserEntity(&newUser)
	if err != nil {
		return uuid.UUID{}, err
	}

	userEntity.ID = roleEntity.ID
	userEntity.Role = *roleEntity

	id, err := service.repository.Create(*userEntity)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

// UpdateUser implements UserService.
func (service *userServiceImpl) UpdateUser(userID uuid.UUID, updatedUser dto.UserDto) error {
	userDto, err := mapper.UserDtoToUserEntity(&updatedUser)
	if err != nil {
		return err
	}
	if err = service.repository.Update(userID, *userDto); err != nil {
		return err
	}

	return nil
}

// DeleteUser implements UserService.
func (service *userServiceImpl) DeleteUser(userID uuid.UUID) error {
	return service.repository.Delete(userID)
}

// RestoreUser implements UserService.
func (service *userServiceImpl) RestoreUser(userID uuid.UUID) error {
	return service.repository.Restore(userID)
}

func NewUserService(userRepository repository.UserRepository, roleRepository repository.RoleRepository) UserService {
	return &userServiceImpl{repository: userRepository, roleRepository: roleRepository}
}
