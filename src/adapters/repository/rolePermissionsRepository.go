package repository

import (
	"fmt"

	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RolePermissionsRepository defines a set of methods for managing role permissions in the system.
type RolePermissionsRepository interface {
	// FindByRoleID retrieves the permissions of a role by its unique identifier (UUID).
	// Returns a pointer to the RolePermissionsEntity if found, or an error otherwise.
	FindByRoleID(roleID uuid.UUID) (*models.RolePermissions, error)

	// Save inserts a new role permissions into the database.
	// Returns the UUID of the newly created role permissions and an error if the operation fails.
	Save(rolePermissions models.RolePermissions) (uuid.UUID, error)
}

type rolePermissionsRepositoryImpl struct {
	db *gorm.DB
}

// FindByRoleID implements RolePermissionsRepository.
func (repo *rolePermissionsRepositoryImpl) FindByRoleID(roleID uuid.UUID) (*models.RolePermissions, error) {
	fmt.Println("entra a FindByRoleID (repository)")
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
