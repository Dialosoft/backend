package repository

import (
	"fmt"

	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RolePermissionsRepository defines a set of methods for managing role permissions in the system.
type RolePermissionsRepository interface {
	AbstractRepository[*models.RolePermissions, uuid.UUID]
	
	// FindByRoleID retrieves the permissions of a role by its unique identifier (UUID).
	// Returns a pointer to the RolePermissionsEntity if found, or an error otherwise.
	FindByRoleID(roleID uuid.UUID) (*models.RolePermissions, error)
}

type rolePermissionsRepositoryImpl struct {
	*abstractRepositoryImpl[*models.RolePermissions, uuid.UUID]
}

// NewRolePermissionsRepository creates a new instance of RolePermissionsRepository
func NewRolePermissionsRepository(db *gorm.DB) RolePermissionsRepository {
	repo := &rolePermissionsRepositoryImpl{}
	repo.abstractRepositoryImpl = CreateRepository(db, repo)
	return repo
}

// FindByRoleID implements RolePermissionsRepository.
func (repo *rolePermissionsRepositoryImpl) FindByRoleID(roleID uuid.UUID) (*models.RolePermissions, error) {
	fmt.Println("get into to FindByRoleID (repository)")
	return repo.FirstByKey("role_id", roleID.String())
}

func (repo *rolePermissionsRepositoryImpl) GetKeyIdName() string {
    return "role_id"
}
