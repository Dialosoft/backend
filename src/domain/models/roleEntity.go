package models

import "github.com/google/uuid"

type RoleEntity struct {
	ID         uuid.UUID `json:"id"`
	RoleType   string    `json:"roleType"`
	Permission int       `json:"permission"`
	AdminRole  bool      `json:"adminRole"`
	ModRole    bool      `json:"modRole"`
}

func (RoleEntity) TableName() string {
	return "roles"
}
