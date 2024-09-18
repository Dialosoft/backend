package models

import (
	"time"

	"github.com/google/uuid"
)

type TokenEntity struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Token     string    `json:"token" gorm:"type:text;not null"`
	RoleID    uuid.UUID `json:"roleID" gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `json:"userID" gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
}

func (TokenEntity) TableName() string {
	return "tokens"
}
