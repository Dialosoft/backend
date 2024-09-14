package repository

import (
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
)

type RoleRepository interface {
	GetAllRoles() ([]models.RoleEntity, error)
	GetRoleByID() (models.RoleEntity, error)
	GetRoleByUsername() (models.RoleEntity, error)
	Create() (uuid.UUID, error)
	Update() (uuid.UUID, error)
	Delete() (uuid.UUID, error)
	Restore() (uuid.UUID, error)
}
