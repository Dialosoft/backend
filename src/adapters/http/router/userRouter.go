package router

import (
	"github.com/Dialosoft/src/adapters/http/controller"
	"github.com/Dialosoft/src/adapters/http/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

type UserRouter struct {
	UserController *controller.UserController
}

func NewUserRouter(userController *controller.UserController) *UserRouter {
	return &UserRouter{UserController: userController}
}

func (r *UserRouter) SetupUserRoutes(api fiber.Router, securityMiddleware *middleware.SecurityMiddleware, permissionMiddleware *middleware.PermissionMiddleware) {

	// public group
	userGroup := api.Group("/users")

	// protected routes by authenticated users
	userProtected := userGroup.Group("/protected",
		securityMiddleware.VerifyRefreshToken(),
		securityMiddleware.GetAndVerifyAccessToken(),
		securityMiddleware.GetRoleFromToken(),
		permissionMiddleware.CanManageUsers(),
	)

	{

		// public
		userGroup.Get("/get-all-users", r.UserController.GetAllUsers)
		userGroup.Get("/get-user-by-id/:id", r.UserController.GetUserByID)
		userGroup.Get("/get-user-by-name/:username", r.UserController.GetUserByUsername)
		userGroup.Get("/avatars/*", static.New("./images/avatars"))
	}

	{

		// protected routes by authenticated users and with permission (manage users permission)
		userProtected.Delete("/delete-user/:id", r.UserController.DeleteUser)
		userProtected.Patch("/restore-user/:id", r.UserController.RestoreUser, r.UserController.DeleteUser)
		userProtected.Put("/update-user/:id", r.UserController.UpdateUser)
		userProtected.Put("/change-user-avatar/:id", r.UserController.ChangeUserAvatar)
	}
}
