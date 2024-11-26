package request

type NewCategory struct {
	Name           *string  `json:"name" example:"Category 2"`
	Description    *string  `json:"description" example:"some description for category 2"`
	RolesAllowedID []string `json:"rolesAllowedID" example:"['8213280e-2000-403a-b375-cdcda6488450', 'roleUUID2']"`
}
