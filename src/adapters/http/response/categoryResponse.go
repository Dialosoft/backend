package response

import (
	"time"

	"github.com/google/uuid"
)

type CategoryResponse struct {
	ID           uuid.UUID `json:"id" example:"8213280e-2000-403a-b375-cdcda6488450"`
	Name         string    `json:"name" example:"Category 1"`
	Description  string    `json:"description" example:"some description for category 1"`
	RolesAllowed []string  `json:"rolesAllowed" example:"['8213280e-2000-403a-b375-cdcda6488450', 'roleUUID2']"`
	CreatedAt    time.Time `json:"created_at" example:"2022-08-01T00:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2022-08-01T00:00:00Z"`
}
