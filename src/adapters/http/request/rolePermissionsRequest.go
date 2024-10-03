package request

type NewRolePermissions struct {
	CanManageCategories *bool `json:"canManageCategories"`
	CanManageForums     *bool `json:"canManageForums"`
	CanManageRoles      *bool `json:"canManageRoles"`
	CanManageUsers      *bool `json:"canManageUsers"`
}
