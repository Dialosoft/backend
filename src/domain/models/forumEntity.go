package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Forum struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name         string         `gorm:"type:varchar(100);unique;not null" json:"name"`
	Description  string         `gorm:"type:varchar(255)" json:"description"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	Type         string         `gorm:"type:varchar(100);not null" json:"type"`
	RolesAllowed pq.StringArray `gorm:"type:text[]" json:"roles_allowed"`
	CategoryID   string         `gorm:"not null" json:"category_id"`
	Category     Category       `gorm:"foreignKey:CategoryID;references:ID" json:"category"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeleteAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Forum) TableName() string {
	return "forums"
}
