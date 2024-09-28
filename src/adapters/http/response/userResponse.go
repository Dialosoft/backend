package response

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserResponse struct {
	ID          uuid.UUID      `json:"id"`
	Username    string         `json:"username"`
	Email       string         `json:"email"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Banned      bool           `json:"banned"`
	Role        RoleResponse   `json:"role"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt"`
}
