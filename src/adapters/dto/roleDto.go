package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleDto struct {
	ID         uuid.UUID      `json:"id"`
	RoleType   string         `json:"roleType"`
	Permission int            `json:"permission"`
	AdminRole  bool           `json:"adminRole"`
	ModRole    bool           `json:"modRole"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt"`
}
