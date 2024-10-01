package request

type NewForum struct {
	Name         *string  `json:"name"`
	Description  *string  `json:"description"`
	Type         *string  `json:"type"`
	IsActive     *bool    `json:"isActive"`
	RolesAllowed []string `json:"rolesAllowed"`
	CategoryID   *string  `json:"categoryID" validate:"required"`
}
