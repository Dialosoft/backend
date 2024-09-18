package controller

import (
	"errors"
	"strings"

	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/services"
	"github.com/Dialosoft/src/pkg/errorsUtils"
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

	return response.Standard(c, "Successfully registered", response.RegisterResponse{
		UserID:       userID.String(),
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func (ac *AuthController) Login(c fiber.Ctx) error {
	var req request.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c)
	}

	accesToken, refreshToken, err := ac.AuthService.Login(req.Username, req.Password)
	if err != nil {
		if err == errorsUtils.ErrUnauthorizedAcces || err == gorm.ErrRecordNotFound {
			return response.ErrUnauthorized(c)
		}
		return response.ErrInternalServer(c)
	}

	return response.Standard(c, "Successfully logged in", response.LoginResponse{
		Token:        accesToken,
		RefreshToken: refreshToken,
	})
}

func (ac *AuthController) RefreshToken(c fiber.Ctx) error {
	var req request.RefreshToken
	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c)
	}

	accesToken, err := ac.AuthService.RefreshToken(req.Refresh)
	if err != nil {
		if err == errorsUtils.ErrUnauthorizedAcces || err == gorm.ErrRecordNotFound || err == errorsUtils.ErrRefreshTokenExpiredOrInvalid {
			return response.ErrUnauthorized(c)
		} else if err == errorsUtils.ErrRoleIDInRefreshToken {
			return response.PersonalizedErr(c, err.Error(), fiber.StatusBadRequest)
		} else if err == errorsUtils.ErrInvalidUUID {
			return response.PersonalizedErr(c, err.Error(), fiber.StatusBadRequest)
		}
		return response.ErrInternalServer(c)
	}

	return response.Standard(c, "successfully refreshed", fiber.Map{
		"accesToken": accesToken,
	})
}
