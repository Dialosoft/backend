package request

type NewRolePermissions struct {
	CanCreateCategory *bool `json:"canCreateCategory"`
	CanCreateForum    *bool `json:"canCreateForum"`
	CanCreateNewRoles *bool `json:"canCreateNewRoles"`
	CanManageRoles    *bool `json:"canManageRoles"`
	CanManageUsers    *bool `json:"canManageUsers"`
}
