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

	// protected routes by authenticated users
	userProtectedForAuthenticatedUsersGroup := api.Group("/users/users-only",
		middleware.VerifyRefreshToken(),
		middleware.GetAndVerifyAccessToken())

	// protectd routes by self user
	userProtectedForSelfUser := userProtectedForAuthenticatedUsersGroup.Group("/users/self-user",
		middleware.AuthorizeSelfUserID())

	{
		userGroup.Get("/get-all-users", r.UserController.GetAllUsers)
		userGroup.Get("/get-user-by-id/:id", r.UserController.GetUserByID)
		userProtectedForSelfUser.Put("/update-user/:id", r.UserController.UpdateUser,
			middleware.VerifyRefreshToken(),
			middleware.GetAndVerifyAccessToken(),
			middleware.AuthorizeSelfUserID(),
		)
		userProtectedForAuthenticatedUsersGroup.Delete("/delete-user/:id", r.UserController.DeleteUser,
			middleware.RoleRequiredByID(defaultRoles["administrator"].String()),
		)
		userProtectedForAuthenticatedUsersGroup.Patch("/restore-user/:id", r.UserController.RestoreUser, r.UserController.DeleteUser,
			middleware.RoleRequiredByID(defaultRoles["administrator"].String()),
		)
		userProtectedForSelfUser.Put("/change-user-avatar/:id", r.UserController.ChangeUserAvatar)
		userGroup.Get("/avatars/*", static.New("./images/avatars"))
	}
}
