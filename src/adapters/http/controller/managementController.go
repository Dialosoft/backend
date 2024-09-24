package controller

import (
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/services"
	"github.com/Dialosoft/src/pkg/utils/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ManagementController struct {
	ForumService    services.ForumService
	CategoryService services.CategoryService
	RoleService     services.RoleService
	UserService     services.UserService
	AuthService     services.AuthService
	CacheService    services.CacheService
}

func (mc *ManagementController) ChangeUserRole(c fiber.Ctx) error {
	var req request.ChangeUserRole
	if err := c.Bind().Body(&req); err != nil {
		logger.CaptureError(err, "Failed to bind ChangeUserRole request", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrBadRequest(c)
	}

	if req.RoleID == "" || req.UserID == "" {
		logger.Error("RoleID or UserID is missing", map[string]interface{}{
			"roleID": req.RoleID,
			"userID": req.UserID,
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrEmptyParametersOrArguments(c)
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		logger.Error("Invalid UUID format for UserID", map[string]interface{}{
			"userID": req.UserID,
			"route":  c.Path(),
			"method": c.Method(),
			"error":  err.Error(),
		})
		return response.ErrUUIDParse(c)
	}

	newUserRequest := request.NewUser{
		RoleID: &req.RoleID,
	}

	// Update the user's role
	if err := mc.UserService.UpdateUser(userUUID, newUserRequest); err != nil {
		logger.CaptureError(err, "Failed to update user role", map[string]interface{}{
			"userUUID": userUUID,
			"roleID":   req.RoleID,
			"route":    c.Path(),
			"method":   c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	// Retrieve the refresh token
	refreshToken, err := mc.CacheService.GetRefreshTokenByID(userUUID)
	if err != nil {
		logger.CaptureError(err, "Failed to retrieve refresh token for user", map[string]interface{}{
			"userUUID": userUUID,
			"route":    c.Path(),
			"method":   c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	// Invalidate the refresh token
	if err := mc.CacheService.InvalidateRefreshToken(refreshToken); err != nil {
		logger.CaptureError(err, "Failed to invalidate refresh token", map[string]interface{}{
			"userUUID":     userUUID,
			"refreshToken": refreshToken,
			"route":        c.Path(),
			"method":       c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("User role updated and refresh token invalidated successfully", map[string]interface{}{
		"userUUID": userUUID,
		"roleID":   req.RoleID,
		"route":    c.Path(),
		"method":   c.Method(),
	})

	return response.Standard(c, "UPDATED", nil)
}

func NewManagamentController(
	forumService services.ForumService,
	categoryService services.CategoryService,
	roleService services.RoleService,
	userService services.UserService,
	AuthService services.AuthService,
	CacheService services.CacheService,
) *ManagementController {

	return &ManagementController{
		ForumService:    forumService,
		CategoryService: categoryService,
		RoleService:     roleService,
		UserService:     userService,
		AuthService:     AuthService,
		CacheService:    CacheService,
	}
}
