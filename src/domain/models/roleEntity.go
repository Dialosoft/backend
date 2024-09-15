package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleEntity struct {
	ID         uuid.UUID      `json:"id"`
	RoleType   string         `json:"roleType"`
	Permission int            `json:"permission"`
	AdminRole  bool           `json:"adminRole"`
	ModRole    bool           `json:"modRole"`
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (RoleEntity) TableName() string {
	return "roles"
}
