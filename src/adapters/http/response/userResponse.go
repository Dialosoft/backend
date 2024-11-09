package response

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserResponse struct {
	ID          uuid.UUID      `json:"id" example:"8213280e-2000-403a-b375-cdcda6488450"`
	Username    string         `json:"username" example:"matdevcoder"`
	Email       string         `json:"email" example:"matdevcoder@email.com"`
	Name        string         `json:"name" example:"Mat"`
	Description string         `json:"description" example:"some description"`
	Banned      bool           `json:"banned" example:"false"`
	Role        RoleResponse   `json:"role" example:"8213280e-2000-403a-b375-cdcda6488450"`
	CreatedAt   time.Time      `json:"createdAt" example:"2022-08-01T00:00:00Z"`
	UpdatedAt   time.Time      `json:"updatedAt" example:"2022-08-01T00:00:00Z"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt" example:"2022-08-01T00:00:00Z"`
}
