package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Category struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name         string         `gorm:"type:varchar(100);unique;not null" json:"name"`
	Description  string         `gorm:"type:varchar(255)" json:"description"`
	RolesAllowed pq.StringArray `gorm:"type:text[]" json:"roles_allowed"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}
