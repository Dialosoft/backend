package request

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=4,max=15"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=35"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshToken struct {
	Refresh string `json:"refreshToken"`
}
