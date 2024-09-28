package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID      `json:"userID"`
	User      UserEntity     `gorm:"foreignKey:UserID" json:"user"`
	PostID    uuid.UUID      `gorm:"type:uuid;index" json:"postId"`
	CommentID *uuid.UUID     `gorm:"type:uuid;index" json:"commentId"`
	Content   string         `gorm:"type:text" json:"content"`
	IsBest    bool           `json:"isBest"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}
