package response

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostResponse struct {
	ID        uuid.UUID      `json:"id" example:"8213280e-2000-403a-b375-cdcda6488450"`
	User      UserResponse   `json:"user" example:"8213280e-2000-403a-b375-cdcda6488450"`
	Forum     ForumResponse  `json:"forumID" example:"8213280e-2000-403a-b375-cdcda6488450"`
	Title     string         `json:"title" example:"some title"`
	Content   string         `json:"content" example:"some content"`
	Views     uint32         `json:"views" example:"10"`
	Comments  uint32         `json:"comments" example:"5"`
	CreatedAt time.Time      `json:"createdAt" example:"2022-08-01T00:00:00Z"`
	UpdatedAt time.Time      `json:"updatedAt" example:"2022-08-01T00:00:00Z"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" example:"2022-08-01T00:00:00Z"`
}

type SimplePostResponse struct {
	ID        string         `json:"id"`
	UserID    string         `json:"userID"`
	Title     string         `json:"title"`
	CreatedAt string         `json:"createdAt"`
	UpdatedAt string         `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
}
