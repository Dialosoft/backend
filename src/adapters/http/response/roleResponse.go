package response

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleResponse struct {
	ID         uuid.UUID      `json:"id" example:"8213280e-2000-403a-b375-cdcda6488450"`
	RoleType   string         `json:"username" example:"moderator"`
	Permission int            `json:"permission" example:"2"`
	AdminRole  bool           `json:"adminRole" example:"true"`
	ModRole    bool           `json:"modRole" example:"false"`
	Email      string         `json:"email" example:"matdevcoder@email.com"`
	CreatedAt  time.Time      `json:"createdAt" example:"2022-08-01T00:00:00Z"`
	UpdatedAt  time.Time      `json:"updatedAt" example:"2022-08-01T00:00:00Z"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt" example:"2022-08-01T00:00:00Z"`
}
