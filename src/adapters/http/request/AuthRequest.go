package request

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=4,max=15" example:"busta"`
	Email    string `json:"email" validate:"required,email" example:"busta@email.com"`
	Password string `json:"password" validate:"required,min=6,max=35" example:"12345678"`
}

type LoginRequest struct {
	Username string `json:"username" example:"matdevcoder"`
	Password string `json:"password" example:"12345678"`
}

type RefreshToken struct {
	Refresh string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
