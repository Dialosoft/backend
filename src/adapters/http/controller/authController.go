package controller

import (
	"errors"
	"strings"

	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/services"
	"github.com/Dialosoft/src/pkg/errorsUtils"
	"github.com/Dialosoft/src/pkg/utils/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

type AuthController struct {
	AuthService services.AuthService
	Layer       string
}

func NewAuthController(authService services.AuthService, validator *validator.Validate, layer string) *AuthController {
	return &AuthController{AuthService: authService, Validator: validator, Layer: layer}
}

func (ac *AuthController) Register(c fiber.Ctx) error {
	var req request.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequestParse(c, err, req, ac.Layer)
	}
	if err := ac.Validator.Struct(req); err != nil {
		logger.CaptureError(err, "Failed validator for RegisterRequest in Controller", map[string]interface{}{
			"request-tried": req,
			"route":         c.Path(),
			"method":        c.Method(),
		})
		return response.PersonalizedErr(c, "validate information error", fiber.StatusBadRequest)
	}
	userDto := dto.UserDto{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	userID, token, refreshToken, err := ac.AuthService.Register(userDto)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return response.ErrConflict(c, err, userDto, ac.Layer)
		}
		return response.ErrInternalServer(c, err, userDto, ac.Layer)
	}

	return response.Standard(c, "Successfully registered", response.RegisterResponse{
		UserID:       userID.String(),
		AccessToken:  token,
		RefreshToken: refreshToken,
	})
}

func (ac *AuthController) Login(c fiber.Ctx) error {
	var req request.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequestParse(c, err, req, ac.Layer)
	}

	accessToken, refreshToken, err := ac.AuthService.Login(req.Username, req.Password)
	if err != nil {
		if err == errorsUtils.ErrUnauthorizedAcces || err == gorm.ErrRecordNotFound {
			return response.ErrUnauthorized(c, req, err, ac.Layer)
		}
		return response.ErrInternalServer(c, err, req, ac.Layer)
	}

	return response.Standard(c, "Successfully logged in", response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (ac *AuthController) RefreshToken(c fiber.Ctx) error {
	var req request.RefreshToken
	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequestParse(c, err, req, ac.Layer)
	}

	accessToken, err := ac.AuthService.RefreshToken(req.Refresh)
	if err != nil {
		if err == errorsUtils.ErrUnauthorizedAcces || err == gorm.ErrRecordNotFound ||
			err == errorsUtils.ErrRefreshTokenExpiredOrInvalid || err == errorsUtils.ErrNotFound {
			return response.ErrUnauthorized(c, req, errorsUtils.ErrRefreshTokenExpiredOrInvalid, ac.Layer)
		}
		return response.ErrInternalServer(c, err, req, ac.Layer)
	}

	return response.Standard(c, "successfully refreshed", fiber.Map{
		"accessToken": accessToken,
	})
}
