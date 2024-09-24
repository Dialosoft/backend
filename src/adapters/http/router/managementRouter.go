package router

import (
	"github.com/Dialosoft/src/adapters/http/controller"
	"github.com/Dialosoft/src/adapters/http/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ManagementRouter struct {
	ManagementController *controller.ManagementController
}

func NewManagementRouter(managementController *controller.ManagementController) *ManagementRouter {
	return &ManagementRouter{ManagementController: managementController}
}

func (r *ManagementRouter) SetupManagementRoutes(api fiber.Router, middlewares *middleware.SecurityMiddleware, defaultRoles map[string]uuid.UUID) {
	managementGroup := api.Group("/management")

	{
		managementGroup.Post("/change-user-role", r.ManagementController.ChangeUserRole,
			middlewares.GetAndVerifyAccessToken(),
			middlewares.VerifyRefreshToken(),
			middlewares.RoleRequiredByID(defaultRoles["administrator"].String()),
		)
		managementGroup.Get("/test", func(c fiber.Ctx) error {
			return c.SendString("pudiste!")
		}, middlewares.GetAndVerifyAccessToken(), middlewares.VerifyRefreshToken())
	}
}
