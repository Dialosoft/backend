package response

type RegisterResponse struct {
	UserID       string
	Token        string
	RefreshToken string
}

type LoginResponse struct {
	Token        string
	RefreshToken string
}
