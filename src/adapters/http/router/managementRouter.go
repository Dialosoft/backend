package router

import (
	"github.com/Dialosoft/src/adapters/http/controller"
	"github.com/Dialosoft/src/adapters/http/middleware"
	"github.com/gofiber/fiber/v3"
)

type ManagementRouter struct {
	ManagementController *controller.ManagementController
}

func NewManagementRouter(managementController *controller.ManagementController) *ManagementRouter {
	return &ManagementRouter{ManagementController: managementController}
}

func (r *ManagementRouter) SetupManagementRoutes(api fiber.Router, securityMiddleware *middleware.SecurityMiddleware, permissionMiddleware *middleware.PermissionMiddleware) {
	managementGroup := api.Group("/management",
		securityMiddleware.GetAndVerifyAccessToken(),
		securityMiddleware.VerifyRefreshToken(),
		securityMiddleware.GetRoleFromToken(),
		permissionMiddleware.CanManageUsers(), /* permission middleware for Users */
		permissionMiddleware.CanManageRoles(), /* permission middleware for Roles */
	)

	/*
		Here cannot add public routes, because the user must be authenticated to access them
		and the permission middleware will check the role permissions !!!
	*/

	{
		managementGroup.Post("/change-user-role", r.ManagementController.ChangeUserRole)
	}
}
