package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type UserEntity struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username  string         `json:"username" gorm:"type:varchar(100);unique;not null"`
	Email     string         `json:"email" gorm:"type:varchar(100);unique;not null"`
	Password  string         `json:"password" gorm:"type:varchar(255);not null"`
	Locked    bool           `json:"locked" gorm:"type:boolean;default:false"`
	Disable   bool           `json:"disable" gorm:"type:boolean;default:false"`
	Role      RoleEntity     `json:"role" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (UserEntity) TableName() string {
	return "users"
}
