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

func (r *UserRouter) SetupUserRoutes(api fiber.Router) {
	userGroup := api.Group("/users")

	{
		userGroup.Get("/get-all-users", r.UserController.GetAllUsers)
		userGroup.Get("/get-user-by-id/:id", r.UserController.GetUserByID)
		userGroup.Post("/create-user", r.UserController.CreateNewUser)
		userGroup.Put("/update-user/:id", r.UserController.UpdateUser)
		userGroup.Delete("/delete-user/:id", r.UserController.DeleteUser)
		userGroup.Patch("/restore-user/:id", r.UserController.RestoreUser)
	}
}
