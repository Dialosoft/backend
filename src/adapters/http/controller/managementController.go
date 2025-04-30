package controller

import (
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/services"
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
	Layer           string
}

// ChangeUserRole changes the role of a specified user.
//
// @Summary Change user role
// @Description Updates the role of a user and invalidates their current refresh token to ensure role changes are applied.
// @Tags Management
// @Accept json
// @Produce json
// @Param request body request.ChangeUserRole true "Change User Role Data"
// @Success 200 {object} response.StandardResponse "UPDATED - User role updated successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid input data or parameters"
// @Failure 404 {object} response.StandardError "NOT FOUND - User or role not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /management/change-user-role [post]
func (mc *ManagementController) ChangeUserRole(c fiber.Ctx) error {
	var req request.ChangeUserRole
	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, mc.Layer)
	}

	if req.RoleID == "" || req.UserID == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		return response.ErrUUIDParse(c, req.UserID)
	}

	stateUserRequest := request.StateUser{
		RoleID: &req.RoleID,
	}

	// Update the user's role
	if err := mc.UserService.UpdateUser(userUUID, stateUserRequest); err != nil {
		return response.ErrInternalServer(c, err, nil, mc.Layer)
	}

	// Retrieve the refresh token
	refreshToken, err := mc.CacheService.GetRefreshTokenByID(userUUID)
	if err != nil {
		return response.ErrInternalServer(c, err, nil, mc.Layer)
	}

	// Invalidate the refresh token
	if err := mc.CacheService.InvalidateRefreshToken(refreshToken); err != nil {
		return response.ErrInternalServer(c, err, nil, mc.Layer)
	}

	return response.Standard(c, "UPDATED", nil)
}

func NewManagamentController(
	forumService services.ForumService,
	categoryService services.CategoryService,
	roleService services.RoleService,
	userService services.UserService,
	AuthService services.AuthService,
	CacheService services.CacheService,
	Layer string,
) *ManagementController {

	return &ManagementController{
		ForumService:    forumService,
		CategoryService: categoryService,
		RoleService:     roleService,
		UserService:     userService,
		AuthService:     AuthService,
		CacheService:    CacheService,
		Layer:           Layer,
	}
}
