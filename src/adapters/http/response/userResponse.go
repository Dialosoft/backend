package response

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserResponse struct {
	ID        uuid.UUID      `json:"id"`
	Username  string         `json:"username"`
	Password  string         `json:"password"`
	Email     string         `json:"email"`
	Locked    bool           `json:"locked"`
	Disable   bool           `json:"disable"`
	Role      RoleResponse   `json:"role"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
}
