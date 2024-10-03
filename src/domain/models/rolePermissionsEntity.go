package models

import "github.com/google/uuid"

type RolePermissions struct {
	RoleID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	CanManageCategories bool      `gorm:"type:bool"`
	CanManageForums     bool      `gorm:"type:bool"`
	CanManageRoles      bool      `gorm:"type:bool"`
	CanManageUsers      bool      `gorm:"type:bool"`
}

func (RolePermissions) TableName() string {
	return "role_permissions"
}
