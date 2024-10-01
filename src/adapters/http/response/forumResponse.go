package response

import (
	"time"

	"github.com/google/uuid"
)

type ForumResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	IsActive     bool      `json:"isActive"`
	Type         string    `json:"type"`
	RolesAllowed []string  `json:"rolesAllowed"`
	CategoryID   string    `json:"categoryId"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
