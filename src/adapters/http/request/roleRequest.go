package request

type NewRole struct {
	RoleType   *string `json:"roleType" example:"moderator"`
	Permission *int    `json:"permission" example:"2"`
	AdminRole  *bool   `json:"adminRole" example:"true"`
	ModRole    *bool   `json:"modRole" example:"false"`
}
