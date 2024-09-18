package controller

import (
	"errors"
	"strings"

	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/services"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type AuthController struct {
	AuthService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

func (ac *AuthController) Register(c fiber.Ctx) error {
	var req request.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c)
	}
	userDto := dto.UserDto{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	userID, token, refreshToken, err := ac.AuthService.Register(userDto)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) ||
			strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return response.ErrConflict(c)
		}
		return response.ErrInternalServer(c)
	}

	return response.Standard(c, "OK", response.RegisterResponse{
		UserID:       userID.String(),
		Token:        token,
		RefreshToken: refreshToken,
	})
}
