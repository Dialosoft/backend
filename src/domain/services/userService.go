package services

import (
	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/mapper"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
)

// UserService defines a set of methods for handling business logic related to users.
// It provides operations like retrieving, creating, updating, and deleting users in the system.
type UserService interface {

	// GetAllUsers retrieves all users as data transfer objects (DTOs).
	// Returns a slice of pointers to UserDto and an error if something goes wrong.
	GetAllUsers() ([]*dto.UserDto, error)

	// GetUserByID retrieves a user by their unique identifier (UUID) as a DTO.
	// Returns a pointer to UserDto if found, or an error otherwise.
	GetUserByID(userID uuid.UUID) (*dto.UserDto, error)

	// GetUserByUsername retrieves a user by their username (string) as a DTO.
	// Returns a pointer to UserDto if found, or an error otherwise.
	GetUserByUsername(username string) (*dto.UserDto, error)

	// CreateNewUser creates a new user based on the provided UserDto.
	// Returns the UUID of the created user and an error if the creation fails.
	CreateNewUser(newUser dto.UserDto) (uuid.UUID, error)

	// UpdateUser modifies an existing user identified by their UUID based on the provided UserDto.
	// Returns an error if the update fails.
	UpdateUser(userID uuid.UUID, req request.NewUser) error

	// DeleteUser marks a user as deleted by their UUID.
	// Returns an error if the deletion fails.
	DeleteUser(userID uuid.UUID) error

	// RestoreUser restores a previously deleted user by their UUID.
	// Returns an error if the restoration fails.
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
		userDto := mapper.UserEntityToUserDto(v)
		usersDtos = append(usersDtos, userDto)
	}

	return usersDtos, nil
}

// GetUserByID implements UserService.
func (service *userServiceImpl) GetUserByID(userID uuid.UUID) (*dto.UserDto, error) {
	userEntity, err := service.repository.FindByID(userID)
	if err != nil {
		return nil, err
	}
	userDto := mapper.UserEntityToUserDto(userEntity)

	return userDto, nil
}

// GetUserByUsername implements UserService.
func (service *userServiceImpl) GetUserByUsername(username string) (*dto.UserDto, error) {
	userEntity, err := service.repository.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	userDto := mapper.UserEntityToUserDto(userEntity)

	return userDto, nil
}

// CreateNewUser implements UserService.
func (service *userServiceImpl) CreateNewUser(newUser dto.UserDto) (uuid.UUID, error) {

	roleEntity, err := service.roleRepository.FindByType("user")
	if err != nil {
		return uuid.UUID{}, err
	}

	userEntity := mapper.UserDtoToUserEntity(&newUser)

	userEntity.ID = roleEntity.ID
	userEntity.Role = *roleEntity

	id, err := service.repository.Create(*userEntity)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

// UpdateUser implements UserService.
func (service *userServiceImpl) UpdateUser(userID uuid.UUID, req request.NewUser) error {
	var roleEntity *models.RoleEntity

	userEntity, err := service.repository.FindByID(userID)
	if err != nil {
		return err
	}

	if req.Username != nil && *req.Username != "" {
		userEntity.Username = *req.Username
	}

	if req.Locked != nil {
		userEntity.Banned = *req.Locked
	}

	if req.RoleID != nil {
		roleUUID, err := uuid.Parse(*req.RoleID)
		if err != nil {
			return err
		}

		roleEntity, err = service.roleRepository.FindByID(roleUUID)
		if err != nil {
			return err
		}

		userEntity.RoleID = roleEntity.ID
		userEntity.Role = *roleEntity
	}

	if err := service.repository.Update(userID, *userEntity); err != nil {
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
