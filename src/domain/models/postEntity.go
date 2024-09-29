package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	User      UserEntity     `gorm:"foreignKey:UserID" json:"user"`
	Title     string         `gorm:"type:varchar(255)" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	Views     uint32         `json:"views"`
	Comments  uint32         `json:"comments"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
}
