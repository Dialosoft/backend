package router

import (
	"github.com/Dialosoft/src/adapters/http/controller"
	"github.com/Dialosoft/src/adapters/http/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type RoleRouter struct {
	RoleController *controller.RoleController
}

func NewRoleRouter(roleController *controller.RoleController) *RoleRouter {
	return &RoleRouter{RoleController: roleController}
}

func (r *RoleRouter) SetupRoleRouter(api fiber.Router, middlewares *middleware.SecurityMiddleware, defaultRoles map[string]uuid.UUID) {
	roleGroup := api.Group("roles")

	{
		roleGroup.Get("/get-all-roles", r.RoleController.GetAllRoles)
		roleGroup.Get("/get-role-by-id/:id", r.RoleController.GetRoleByID)
		roleGroup.Get("/get-role-by-type/:type", r.RoleController.GetRoleByType)
		roleGroup.Post("/create-new-role", r.RoleController.CreateNewRole,
			middlewares.GetAndVerifyAccesToken(), middlewares.VerifyRefreshToken(), middlewares.RoleRequiredByID(defaultRoles["administrator"].String()),
		)
		roleGroup.Put("/update-role/:id", r.RoleController.UpdateRole,
			middlewares.GetAndVerifyAccesToken(), middlewares.VerifyRefreshToken(), middlewares.RoleRequiredByID(defaultRoles["administrator"].String()),
		)
		roleGroup.Delete("/delete-role/:id", r.RoleController.DeleteRole,
			middlewares.GetAndVerifyAccesToken(), middlewares.VerifyRefreshToken(), middlewares.RoleRequiredByID(defaultRoles["administrator"].String()),
		)
		roleGroup.Put("/restore-role/:id", r.RoleController.RestoreRole,
			middlewares.GetAndVerifyAccesToken(), middlewares.VerifyRefreshToken(), middlewares.RoleRequiredByID(defaultRoles["administrator"].String()),
		)
	}
}
