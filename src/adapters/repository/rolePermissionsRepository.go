package repository

import (
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RolePermissionsRepository interface {
	FindByRoleID(roleID uuid.UUID) (*models.RolePermissions, error)
	Save(rolePermissions models.RolePermissions) (uuid.UUID, error)
}

type rolePermissionsRepositoryImpl struct {
	db *gorm.DB
}

// FindByRoleID implements RolePermissionsRepository.
func (repo *rolePermissionsRepositoryImpl) FindByRoleID(roleID uuid.UUID) (*models.RolePermissions, error) {
	var rolePermission models.RolePermissions
	if err := repo.db.First(&rolePermission, "role_id = ?", roleID.String()).Error; err != nil {
		return nil, err
	}

	return &rolePermission, nil
}

// Save implements RolePermissionsRepository.
func (repo *rolePermissionsRepositoryImpl) Save(rolePermissions models.RolePermissions) (uuid.UUID, error) {
	result := repo.db.Save(&rolePermissions)
	if result.Error != nil {
		return uuid.UUID{}, result.Error
	}
	return rolePermissions.RoleID, nil
}

func NewRolePermissionsRepository(db *gorm.DB) RolePermissionsRepository {
	return &rolePermissionsRepositoryImpl{db: db}
}
