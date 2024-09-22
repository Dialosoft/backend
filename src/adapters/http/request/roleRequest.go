package request

type NewRole struct {
	RoleType   *string `json:"roleType"`
	Permission *int    `json:"permission"`
	AdminRole  *bool   `json:"adminRole"`
	ModRole    *bool   `json:"modRole"`
}
