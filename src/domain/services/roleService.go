package services

import (
	"fmt"

	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/mapper"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/google/uuid"
)

// RoleService defines a set of methods for handling business logic related to roles.
// It interacts with role data transfer objects (RoleDto) for operations like retrieving,
// creating, updating, and deleting roles in the system.
type RoleService interface {

	// GetAllRoles retrieves all roles as data transfer objects (DTOs).
	// Returns a slice of pointers to RoleDto and an error if something goes wrong.
	GetAllRoles() ([]*dto.RoleDto, error)

	// GetRoleByID retrieves a role by its unique identifier (UUID) as a DTO.
	// Returns a pointer to RoleDto if found, or an error otherwise.
	GetRoleByID(roleID uuid.UUID) (*dto.RoleDto, error)

	// GetRoleByType retrieves a role by its type (string) as a DTO.
	// Returns a pointer to RoleDto if found, or an error otherwise.
	GetRoleByType(roleType string) (*dto.RoleDto, error)

	// CreateNewRole creates a new role based on the provided RoleDto.
	// Returns the UUID of the created role and an error if the creation fails.
	CreateNewRole(newRole dto.RoleDto) (uuid.UUID, error)

	// UpdateRole modifies an existing role identified by its UUID based on the provided request.
	// Returns an error if the update fails.
	UpdateRole(roleID uuid.UUID, req request.NewRole) error

	// DeleteRole marks a role as deleted by its UUID.
	// Returns an error if the deletion fails.
	DeleteRole(roleID uuid.UUID) error

	// RestoreRole restores a previously deleted role by its UUID.
	// Returns an error if the restoration fails.
	RestoreRole(roleID uuid.UUID) error

	// GetDefaultRoles retrieves a map of default role types and their corresponding UUIDs.
	// Returns the map and an error if something goes wrong.
	GetDefaultRoles() (map[string]uuid.UUID, error)
}

type roleServiceImpl struct {
	repository repository.RoleRepository
}

// GetDefaultRoles implements RoleService.
func (service *roleServiceImpl) GetDefaultRoles() (map[string]uuid.UUID, error) {
	roles := map[string]uuid.UUID{}

	roleTypes := []string{"user", "moderator", "administrator"}

	for _, roleType := range roleTypes {
		role, err := service.GetRoleByType(roleType)
		if err != nil {
			return nil, fmt.Errorf("failed to get default role %s: %w", roleType, err)
		}
		roles[roleType] = role.ID
	}

	return roles, nil
}

// GetAllRoles implements RoleService.
func (service *roleServiceImpl) GetAllRoles() ([]*dto.RoleDto, error) {
	var rolesDtos []*dto.RoleDto

	rolesEntities, err := service.repository.FindAllRoles()
	if err != nil {
		return nil, err
	}

	for _, v := range rolesEntities {
		roleDto := mapper.RoleEntityToRoleDto(v)
		rolesDtos = append(rolesDtos, roleDto)
	}

	return rolesDtos, nil
}

// GetRoleByID implements RoleService.
func (service *roleServiceImpl) GetRoleByID(roleID uuid.UUID) (*dto.RoleDto, error) {
	roleEntity, err := service.repository.FindByID(roleID)
	if err != nil {
		return nil, err
	}
	roleDto := mapper.RoleEntityToRoleDto(roleEntity)

	return roleDto, nil
}

// GetRoleByType implements RoleService.
func (service *roleServiceImpl) GetRoleByType(roleType string) (*dto.RoleDto, error) {
	roleEntity, err := service.repository.FindByType(roleType)
	if err != nil {
		return nil, err
	}
	roleDto := mapper.RoleEntityToRoleDto(roleEntity)

	return roleDto, nil
}

// CreateNewRole implements RoleService.
func (service *roleServiceImpl) CreateNewRole(newRole dto.RoleDto) (uuid.UUID, error) {
	roleEntity := mapper.RoleDtoToRoleEntity(&newRole)

	roleUUID, err := service.repository.Create(*roleEntity)
	if err != nil {
		return uuid.UUID{}, err
	}

	return roleUUID, nil
}

// UpdateRole implements RoleService.
func (service *roleServiceImpl) UpdateRole(roleID uuid.UUID, req request.NewRole) error {
	existingRole, err := service.repository.FindByID(roleID)
	if err != nil {
		return err
	}

	if req.RoleType != nil {
		existingRole.RoleType = *req.RoleType
	}
	if req.Permission != nil {
		existingRole.Permission = *req.Permission
	}
	if req.AdminRole != nil {
		existingRole.AdminRole = *req.AdminRole
	}
	if req.ModRole != nil {
		existingRole.ModRole = *req.ModRole
	}

	return service.repository.Update(roleID, *existingRole)
}

// DeleteRole implements RoleService.
func (service *roleServiceImpl) DeleteRole(roleID uuid.UUID) error {
	return service.repository.Delete(roleID)
}

// RestoreRole implements RoleService.
func (service *roleServiceImpl) RestoreRole(roleID uuid.UUID) error {
	return service.repository.Restore(roleID)
}

func NewRoleRepository(roleRepository repository.RoleRepository) RoleService {
	return &roleServiceImpl{repository: roleRepository}
}
