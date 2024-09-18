package response

type RegisterResponse struct {
	UserID       string `json:"userID"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
