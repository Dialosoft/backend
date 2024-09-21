package dto

import (
	"time"

	"github.com/google/uuid"
)

type ForumDto struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	Type        string    `json:"type" validate:"required"`
	PostCount   uint32    `json:"post_count"`
	CategoryID  string    `json:"category_id" validate:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
