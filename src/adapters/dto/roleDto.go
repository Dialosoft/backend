package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleDto struct {
	ID         uuid.UUID      `json:"id"`
	RoleType   string         `json:"username"`
	Permission int            `json:"permission"`
	AdminRole  bool           `json:"adminRole"`
	ModRole    bool           `json:"modRole"`
	Email      string         `json:"email"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt"`
}
