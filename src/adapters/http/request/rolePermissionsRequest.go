package request

type NewRolePermissions struct {
	CanManageCategories *bool `json:"canManageCategories" example:"true"`
	CanManageForums     *bool `json:"canManageForums" example:"false"`
	CanManageRoles      *bool `json:"canManageRoles" example:"true"`
	CanManageUsers      *bool `json:"canManageUsers" example:"false"`
}
