package request

type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Username string `jso:"username"`
	Email    string `json:"email"`
}

type NewUser struct {
	Username *string `json:"username"`
	Locked   *bool   `json:"locked"`
	Disable  *bool   `json:"disable"`
	RoleID   *string `json:"userID"`
}
