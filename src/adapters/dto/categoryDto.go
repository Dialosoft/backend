package dto

import (
	"time"

	"github.com/google/uuid"
)

type CategoryDto struct {
	ID           uuid.UUID `json:"id,omitempty"`
	Name         string    `json:"name" validate:"required"`
	Description  string    `json:"description,omitempty"`
	RolesAllowed []string  `json:"rolesAllowed"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}
