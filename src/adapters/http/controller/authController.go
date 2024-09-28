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
		logger.CaptureError(err, "Failed to parse RegisterRequest in Controller", map[string]interface{}{
			"request-tried": req,
			"route":         c.Path(),
			"method":        c.Method(),
		})
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
			logger.Warn("Duplicate key error during registration", map[string]interface{}{
				"error":   err.Error(),
				"userDto": userDto,
				"route":   c.Path(),
				"method":  c.Method(),
			})
			return response.ErrConflict(c)
		}
		logger.CaptureError(err, "Error during user registration", map[string]interface{}{
			"userDto": userDto,
			"route":   c.Path(),
			"method":  c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("User registered successfully", map[string]interface{}{
		"userID": userID.String(),
		"route":  c.Path(),
		"method": c.Method(),
	})

	return response.Standard(c, "Successfully registered", response.RegisterResponse{
		UserID:       userID.String(),
		AccessToken:  token,
		RefreshToken: refreshToken,
	})
}

func (ac *AuthController) Login(c fiber.Ctx) error {
	var req request.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		logger.CaptureError(err, "Failed to parse LoginRequest in Controller", map[string]interface{}{
			"request-tried": req,
			"route":         c.Path(),
			"method":        c.Method(),
		})
		return response.ErrBadRequest(c)
	}

	accessToken, refreshToken, err := ac.AuthService.Login(req.Username, req.Password)
	if err != nil {
		if err == errorsUtils.ErrUnauthorizedAcces || err == gorm.ErrRecordNotFound {
			logger.Warn("Unauthorized login attempt", map[string]interface{}{
				"username": req.Username,
				"error":    err.Error(),
				"route":    c.Path(),
				"method":   c.Method(),
			})
			return response.ErrUnauthorized(c)
		}
		logger.CaptureError(err, "Error during user login", map[string]interface{}{
			"username": req.Username,
			"route":    c.Path(),
			"method":   c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("User logged in successfully", map[string]interface{}{
		"username": req.Username,
		"route":    c.Path(),
		"method":   c.Method(),
	})

	return response.Standard(c, "Successfully logged in", response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (ac *AuthController) RefreshToken(c fiber.Ctx) error {
	var req request.RefreshToken
	if err := c.Bind().Body(&req); err != nil {
		logger.CaptureError(err, "Failed to parse RefreshToken request in Controller", map[string]interface{}{
			"request-tried": req,
			"route":         c.Path(),
			"method":        c.Method(),
		})
		return response.ErrBadRequest(c)
	}

	accessToken, err := ac.AuthService.RefreshToken(req.Refresh)
	if err != nil {
		if err == errorsUtils.ErrUnauthorizedAcces || err == gorm.ErrRecordNotFound ||
			err == errorsUtils.ErrRefreshTokenExpiredOrInvalid || err == errorsUtils.ErrNotFound {
			logger.Warn("Invalid or expired refresh token", map[string]interface{}{
				"refreshToken": req.Refresh,
				"error":        err.Error(),
				"route":        c.Path(),
				"method":       c.Method(),
			})
			return response.ErrUnauthorized(c)
		}
		logger.CaptureError(err, "Error during token refresh", map[string]interface{}{
			"refreshToken": req.Refresh,
			"route":        c.Path(),
			"method":       c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Refresh token successfully generated", map[string]interface{}{
		"route":  c.Path(),
		"method": c.Method(),
	})

	return response.Standard(c, "successfully refreshed", fiber.Map{
		"accessToken": accessToken,
	})
}
