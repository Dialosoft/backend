package request

type ChangeUserRole struct {
	UserID string `json:"userID"`
	RoleID string `json:"roleID"`
}
