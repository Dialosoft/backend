package models

import "github.com/google/uuid"

type RoleEntity struct {
	Id         uuid.UUID `json:"id"`
	RoleType   string    `json:"roleType"`
	Permission int       `json:"permission"`
	AdminRole  bool      `json:"adminRole"`
	ModRole    bool      `json:"modRole"`
}
