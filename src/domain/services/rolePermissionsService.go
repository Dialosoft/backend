package services

import (
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
)

type RolePermissionsService interface {
	GetRolePermissionsByRoleID(roleID uuid.UUID) (*models.RolePermissions, error)
	SetRolePermissionsByRoleID(roleID uuid.UUID, rolePermissionsRequest request.NewRolePermissions) (uuid.UUID, error)
}

type rolePermissionsServiceImpl struct {
	repository repository.RolePermissionsRepository
}

// GetRolePermissionsByRoleID implements RolePermissionsService.
func (service *rolePermissionsServiceImpl) GetRolePermissionsByRoleID(roleID uuid.UUID) (*models.RolePermissions, error) {
	rolePermission, err := service.repository.FindByRoleID(roleID)
	if err != nil {
		return nil, err
	}
	return rolePermission, nil
}

// SetRolePermissionsByRoleID implements RolePermissionsService.
func (service *rolePermissionsServiceImpl) SetRolePermissionsByRoleID(roleID uuid.UUID, rolePermissionsRequest request.NewRolePermissions) (uuid.UUID, error) {
	rolePermission, err := service.repository.FindByRoleID(roleID)
	if err != nil {
		return uuid.UUID{}, err
	}

	if rolePermissionsRequest.CanCreateCategory != nil {
		rolePermission.CanCreateCategory = *rolePermissionsRequest.CanCreateCategory
	}
	if rolePermissionsRequest.CanCreateForum != nil {
		rolePermission.CanCreateForum = *rolePermissionsRequest.CanCreateForum
	}
	if rolePermissionsRequest.CanCreateNewRoles != nil {
		rolePermission.CanCreateNewRoles = *rolePermissionsRequest.CanCreateNewRoles
	}
	if rolePermissionsRequest.CanManageRoles != nil {
		rolePermission.CanManageRoles = *rolePermissionsRequest.CanManageRoles
	}
	if rolePermissionsRequest.CanManageUsers != nil {
		rolePermission.CanManageUsers = *rolePermissionsRequest.CanManageUsers
	}

	roleUUID, err := service.repository.Save(*rolePermission)
	if err != nil {
		return uuid.UUID{}, err
	}

	return roleUUID, nil
}

func NewRolePermissionsService(repository repository.RolePermissionsRepository) RolePermissionsService {
	return &rolePermissionsServiceImpl{repository: repository}
}
