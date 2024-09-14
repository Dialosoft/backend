package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserEntity struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username  string         `json:"username" gorm:"type:varchar(100);unique;not null"`
	Email     string         `json:"email" gorm:"type:varchar(100);unique;not null"`
	Password  string         `json:"password" gorm:"type:varchar(255);not null"`
	Locked    bool           `json:"locked" gorm:"type:boolean;default:false"`
	Disable   bool           `json:"disable" gorm:"type:boolean;default:false"`
	RoleID    uuid.UUID      `json:"roleID" gorm:"type:uuid"`
	Role      RoleEntity     `json:"role" gorm:"foreignKey:RoleID"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (UserEntity) TableName() string {
	return "users"
}
