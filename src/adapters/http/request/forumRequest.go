package request

type NewForum struct {
	Name         *string  `json:"name" example:"Forum 2"`
	Description  *string  `json:"description" example:"some description for forum 2"`
	Type         *string  `json:"type" example:"information"`
	IsActive     *bool    `json:"isActive" example:"true"`
	RolesAllowed []string `json:"rolesAllowed" example:"['8213280e-2000-403a-b375-cdcda6488450', 'roleUUID2']"`
	CategoryID   *string  `json:"categoryID" validate:"required" example:"8213280e-2000-403a-b375-cdcda6488450"`
}
