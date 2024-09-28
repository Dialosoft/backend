package response

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostResponse struct {
	ID        uuid.UUID      `json:"id"`
	User      UserResponse   `json:"user"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Views     uint32         `json:"views"`
	Comments  uint32         `json:"comments"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
}
