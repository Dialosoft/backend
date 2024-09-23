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
		return response.ErrBadRequest(c)
	}

	if req.RoleID == "" || req.UserID == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		logger.Error(err.Error())
		return response.ErrUUIDParse(c)
	}

	newUserRequest := request.NewUser{
		RoleID: &req.RoleID,
	}

	if err := mc.UserService.UpdateUser(userUUID, newUserRequest); err != nil {
		logger.Error(err.Error())
		return response.ErrInternalServer(c)
	}

	refreshToken, err := mc.CacheService.GetRefreshTokenByID(userUUID)
	if err != nil {
		logger.Error(err.Error())
		return response.ErrInternalServer(c)
	}

	if err := mc.CacheService.InvalidateRefreshToken(refreshToken); err != nil {
		logger.Error(err.Error())
		return response.ErrInternalServer(c)
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
