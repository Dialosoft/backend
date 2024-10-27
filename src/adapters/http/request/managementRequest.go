package request

type ChangeUserRole struct {
	UserID string `json:"userID" example:"8213280e-2000-403a-b375-cdcda6488450"`
	RoleID string `json:"roleID" example:"8213280e-2000-403a-b375-cdcda6488452"`
}
