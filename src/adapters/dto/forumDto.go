package dto

import (
	"time"

	"github.com/google/uuid"
)

type ForumDto struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"isActive"`
	Type        string    `json:"type" validate:"required"`
	PostCount   uint32    `json:"postCount"`
	CategoryID  string    `json:"categoryId" validate:"required"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}
