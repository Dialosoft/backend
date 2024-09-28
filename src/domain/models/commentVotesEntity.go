package models

import (
	"time"

	"github.com/google/uuid"
)

type CommentVotes struct {
	CommentID uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
}

func (CommentVotes) TableName() string {
	return "comment_votes"
}
