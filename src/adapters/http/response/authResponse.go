package response

type RegisterResponse struct {
	UserID       string `json:"userID"`
	AccesToken   string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type LoginResponse struct {
	AccesToken   string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
