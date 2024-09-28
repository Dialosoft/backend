package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserDto struct {
	ID          uuid.UUID      `json:"id"`
	Username    string         `json:"username"`
	Password    string         `json:"password"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Email       string         `json:"email"`
	Banned      bool           `json:"locked"`
	Role        RoleDto        `json:"role"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt"`
}
