package response

type RegisterResponse struct {
	UserID       string `json:"userID"`
	AccesToken   string `json:"accesToken"`
	RefreshToken string `json:"refreshToken"`
}

type LoginResponse struct {
	AccesToken   string `json:"accesToken"`
	RefreshToken string `json:"refreshToken"`
}
