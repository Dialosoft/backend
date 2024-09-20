package router

import (
	"log"

	"github.com/Dialosoft/src/adapters/http/controller"
	"github.com/Dialosoft/src/adapters/http/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UserRouter struct {
	UserController *controller.UserController
}

func NewUserRouter(userController *controller.UserController) *UserRouter {
	return &UserRouter{UserController: userController}
}

func (r *UserRouter) SetupUserRoutes(api fiber.Router, middleware *middleware.AuthMiddleware, defaultRoles map[string]uuid.UUID) {
	userGroup := api.Group("/users")

	adminRoleID := defaultRoles["administrator"]
	if adminRoleID.String() == "" {
		log.Panicf("map adminRole is clean!")
	}

	{
		userGroup.Get("/get-all-users", r.UserController.GetAllUsers)
		userGroup.Get("/get-user-by-id/:id", r.UserController.GetUserByID)
		userGroup.Put("/update-user/:id", r.UserController.UpdateUser,
			middleware.IsTokenBlacklisted(),
			middleware.AuthorizeSelfUserID(),
		)
		userGroup.Delete("/delete-user/:id", r.UserController.DeleteUser,
			middleware.IsTokenBlacklisted(),
			middleware.RoleRequiredByID(adminRoleID.String()),
		)
		userGroup.Patch("/restore-user/:id", r.UserController.RestoreUser,
			middleware.IsTokenBlacklisted(),
			middleware.RoleRequiredByID(adminRoleID.String()))
	}
}
