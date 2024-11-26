package response

import (
	"time"

	"github.com/google/uuid"
)

type ForumResponse struct {
	ID           uuid.UUID `json:"id" example:"8213280e-2000-403a-b375-cdcda6488450"`
	Name         string    `json:"name" example:"Forum 1"`
	Description  string    `json:"description" example:"some description for forum 1"`
	IsActive     bool      `json:"isActive" example:"true"`
	Type         string    `json:"type" example:"information"`
	RolesAllowed []string  `json:"rolesAllowed" example:"['8213280e-2000-403a-b375-cdcda6488450', 'roleUUID2']"`
	CategoryID   string    `json:"categoryId" example:"8213280e-2000-403a-b375-cdcda6488450"`
	CreatedAt    time.Time `json:"createdAt" example:"2022-08-01T00:00:00Z"`
	UpdatedAt    time.Time `json:"updatedAt" example:"2022-08-01T00:00:00Z"`
}
