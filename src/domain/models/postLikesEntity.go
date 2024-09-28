package models

import (
	"time"

	"github.com/google/uuid"
)

type PostLikes struct {
	PostID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
}

func (PostLikes) TableName() string {
	return "posts_likes"
}
