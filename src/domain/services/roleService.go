package services

import (
	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/mapper"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/google/uuid"
)

type RoleService interface {
	GetAllRoles() ([]*dto.RoleDto, error)
	GetRoleByID(roleID uuid.UUID) (*dto.RoleDto, error)
	GetRoleByType(roleType string) (*dto.RoleDto, error)
	CreateNewRole(newRole dto.RoleDto) (uuid.UUID, error)
	UpdateRole(roleID uuid.UUID, updatedRole dto.RoleDto) error
	DeleteRole(roleID uuid.UUID) error
	RestoreRole(roleID uuid.UUID) error
}

type roleServiceImpl struct {
	repository repository.RoleRepository
}

// GetAllRoles implements RoleService.
func (service *roleServiceImpl) GetAllRoles() ([]*dto.RoleDto, error) {
	var rolesDtos []*dto.RoleDto

	rolesEntities, err := service.repository.FindAllRoles()
	if err != nil {
		return nil, err
	}

	for _, v := range rolesEntities {
		roleDto, err := mapper.RoleEntityToRoleDto(v)
		if err != nil {
			return nil, err
		} else {
			rolesDtos = append(rolesDtos, roleDto)
		}
	}

	return rolesDtos, nil
}

// GetRoleByID implements RoleService.
func (service *roleServiceImpl) GetRoleByID(roleID uuid.UUID) (*dto.RoleDto, error) {
	roleEntity, err := service.repository.FindByID(roleID)
	if err != nil {
		return nil, err
	}
	roleDto, err := mapper.RoleEntityToRoleDto(roleEntity)
	if err != nil {
		return nil, err
	}

	return roleDto, nil
}

// GetRoleByType implements RoleService.
func (service *roleServiceImpl) GetRoleByType(roleType string) (*dto.RoleDto, error) {
	roleEntity, err := service.repository.FindByType(roleType)
	if err != nil {
		return nil, err
	}
	roleDto, err := mapper.RoleEntityToRoleDto(roleEntity)
	if err != nil {
		return nil, err
	}

	return roleDto, nil
}

// CreateNewRole implements RoleService.
func (service *roleServiceImpl) CreateNewRole(newRole dto.RoleDto) (uuid.UUID, error) {
	roleEntity, err := mapper.RoleDtoToRoleEntity(&newRole)
	if err != nil {
		return uuid.UUID{}, err
	}

	id, err := service.repository.Create(*roleEntity)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

// UpdateRole implements RoleService.
func (service *roleServiceImpl) UpdateRole(roleID uuid.UUID, updatedRole dto.RoleDto) error {
	roleEntity, err := mapper.RoleDtoToRoleEntity(&updatedRole)
	if err != nil {
		return err
	}

	return service.repository.Update(roleID, *roleEntity)
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
