package router

import (
	"github.com/Dialosoft/src/adapters/http/controller"
	"github.com/Dialosoft/src/adapters/http/middleware"
	"github.com/gofiber/fiber/v3"
)

type RoleRouter struct {
	RoleController *controller.RoleController
}

func NewRoleRouter(roleController *controller.RoleController) *RoleRouter {
	return &RoleRouter{RoleController: roleController}
}

func (r *RoleRouter) SetupRoleRouter(api fiber.Router, securityMiddleware *middleware.SecurityMiddleware, permissionMiddleware *middleware.PermissionMiddleware) {

	// public group
	roleGroup := api.Group("/roles")

	// protected group
	roleProtected := roleGroup.Group("/protected",
		securityMiddleware.GetAndVerifyAccessToken(),
		securityMiddleware.VerifyRefreshToken(),
		securityMiddleware.GetRoleFromToken(),
		permissionMiddleware.CanManageRoles() /* permission middleware for Roles */)

	{
		// public
		roleGroup.Get("/get-all-roles", r.RoleController.GetAllRoles)
		roleGroup.Get("/get-role-by-id/:id", r.RoleController.GetRoleByID)
		roleGroup.Get("/get-role-by-type/:type", r.RoleController.GetRoleByType)
	}

	{
		// protected routes by authenticated users and with permission (manage roles permission)
		roleProtected.Get("/get-role-permissions-by-id/:id", r.RoleController.GetRolePermissionsByRoleID)
		roleProtected.Put("/set-role-permissions-by-id/:id", r.RoleController.SetRolePermissionsByRoleID)
		roleProtected.Post("/create-new-role", r.RoleController.CreateNewRole)
		roleProtected.Put("/update-role/:id", r.RoleController.UpdateRole)
		roleProtected.Delete("/delete-role/:id", r.RoleController.DeleteRole)
		roleProtected.Put("/restore-role/:id", r.RoleController.RestoreRole)
	}
}
