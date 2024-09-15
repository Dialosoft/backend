package router

import (
	"github.com/Dialosoft/src/adapters/http/controller"
	"github.com/gofiber/fiber/v3"
)

type UserRouter struct {
	UserController *controller.UserController
}

func NewUserRouter(userController *controller.UserController) *UserRouter {
	return &UserRouter{UserController: userController}
}

func (r *UserRouter) SetupUserRoutes(api *fiber.App) {
	userGroup := api.Group("/users")

	{
		userGroup.Get("", r.UserController.GetAllUsers)
	}
}
