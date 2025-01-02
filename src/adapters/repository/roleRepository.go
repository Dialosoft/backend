package repository

import (
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleRepository defines a set of methods for interacting with roles
// in the system. Each method represents a CRUD action on the RoleEntity.
type RoleRepository interface {
	AbstractRepository[*models.RoleEntity, uuid.UUID]

	// FindByType retrieves a RoleEntity by its type (string).
	// Returns a pointer to the RoleEntity if found, or an error otherwise.
	FindByType(roleType string) (*models.RoleEntity, error)
}

type roleRepositoryImpl struct {
	*abstractRepositoryImpl[*models.RoleEntity, uuid.UUID]
}

// FindByType implements RoleRepository.
func (repo *roleRepositoryImpl) FindByType(roleType string) (*models.RoleEntity, error) {
	return repo.FirstByKey("role_type", roleType)
}

// NewRoleRepository creates a new instance of RoleRepository
func NewRoleRepository(db *gorm.DB) RoleRepository {
	repo := &roleRepositoryImpl{}
	repo.abstractRepositoryImpl = CreateRepository(db, repo)
	return repo
}
