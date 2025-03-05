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
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type AuthController struct {
	AuthService services.AuthService
	Layer       string
	Validator   *validator.Validate
}

func NewAuthController(authService services.AuthService, validator *validator.Validate, layer string) *AuthController {
	return &AuthController{AuthService: authService, Validator: validator, Layer: layer}
}

// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param register body request.RegisterRequest true "Register Request"
// @Success 200 {object} response.RegisterResponse
// @Failure 400 {object} response.StandardError "BAD REQUEST"
// @Failure 409 {object} response.StandardError "CONFLICT"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR"
// @Router /auth/register [post]
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
		return response.RegisterValidatiorErr(c, err)
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

	return response.StandardCreated(c, "Successfully registered", response.RegisterResponse{
		UserID:       userID.String(),
		AccessToken:  token,
		RefreshToken: refreshToken,
	})
}

// @Summary Log in a user
// @Description Log in a user with username and password
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param login body request.LoginRequest true "Login Request"
// @Success 200 {object} response.LoginResponse
// @Failure 400 {object} response.StandardError "BAD REQUEST"
// @Failure 400 {object} response.ValidatorError "Credential validation failed"
// @Failure 401 {object} response.StandardError "UNAUTHORIZED"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR"
// @Router /auth/login [post]
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

// @Summary Refresh an access token
// @Description Refresh the access token using a refresh token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param refreshToken body request.RefreshToken true "Refresh Token Request"
// @Success 200 {object} response.RefreshTokenResponse
// @Failure 400 {object} response.StandardError "BAD REQUEST"
// @Failure 401 {object} response.StandardError "UNAUTHORIZED"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR"
// @Router /auth/refresh [post]
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

	return response.Standard(c, "successfully refreshed", response.RefreshTokenResponse{
		AccessToken: accessToken,
	})
}
