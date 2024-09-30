package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleEntity struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RoleType   string         `json:"roleType" gorm:"type:varchar(250);unique;not null"`
	Permission int            `json:"permission"`
	AdminRole  bool           `json:"adminRole"`
	ModRole    bool           `json:"modRole"`
	UserRole   bool           `json:"userRole"`
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (RoleEntity) TableName() string {
	return "roles"
}
