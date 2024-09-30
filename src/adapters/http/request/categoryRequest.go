package request

type NewCategory struct {
	Name           *string  `json:"name"`
	Description    *string  `json:"description"`
	RolesAllowedID []string `json:"rolesAllowedID"`
}
