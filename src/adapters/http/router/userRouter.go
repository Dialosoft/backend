package router

import (
	"github.com/Dialosoft/src/adapters/http/controller"
	"github.com/Dialosoft/src/adapters/http/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/google/uuid"
)

type UserRouter struct {
	UserController *controller.UserController
}

func NewUserRouter(userController *controller.UserController) *UserRouter {
	return &UserRouter{UserController: userController}
}

func (r *UserRouter) SetupUserRoutes(api fiber.Router, middleware *middleware.SecurityMiddleware, defaultRoles map[string]uuid.UUID) {

	// free routes
	userGroup := api.Group("/users")

	{
		userGroup.Get("/get-all-users", r.UserController.GetAllUsers)
		userGroup.Get("/get-user-by-id/:id", r.UserController.GetUserByID)
		userGroup.Get("/get-user-by-name/:username", r.UserController.GetUserByUsername)
	}

	// protected routes by authenticated users
	userProtectedForAuthenticatedUsersGroup := userGroup.Group("/auth-only",
		middleware.VerifyRefreshToken(),
		middleware.GetAndVerifyAccessToken())

	{
		userProtectedForAuthenticatedUsersGroup.Delete("/delete-user/:id", r.UserController.DeleteUser,
			middleware.RoleRequiredByID(defaultRoles["administrator"].String()))
		userProtectedForAuthenticatedUsersGroup.Patch("/restore-user/:id", r.UserController.RestoreUser, r.UserController.DeleteUser,
			middleware.RoleRequiredByID(defaultRoles["administrator"].String()))
	}

	// protectd routes by self user
	userProtectedForSelfUser := userGroup.Group("/self-user",
		middleware.AuthorizeSelfUserID())

	{
		userProtectedForSelfUser.Put("/update-user/:id", r.UserController.UpdateUser,
			middleware.VerifyRefreshToken(),
			middleware.GetAndVerifyAccessToken(),
			middleware.AuthorizeSelfUserID())
		userProtectedForSelfUser.Put("/change-user-avatar/:id", r.UserController.ChangeUserAvatar)
		userGroup.Get("/avatars/*", static.New("./images/avatars"))
	}
}
