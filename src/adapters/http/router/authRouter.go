package router

import (
	"github.com/Dialosoft/src/adapters/http/controller"
	"github.com/Dialosoft/src/adapters/http/middleware"
	"github.com/gofiber/fiber/v3"
)

type AuthRouter struct {
	AuthController *controller.AuthController
}

func NewAuthRouter(authController *controller.AuthController) *AuthRouter {
	return &AuthRouter{AuthController: authController}
}

func (r *AuthRouter) SetupAuthRoutes(api fiber.Router, middlewares *middleware.SecurityMiddleware) {
	authGroup := api.Group("/auth")
	{
		authGroup.Post("/register", r.AuthController.Register)
		authGroup.Post("/login", r.AuthController.Login)
		authGroup.Post("/refresh-token", r.AuthController.RefreshToken, middlewares.VerifyRefreshToken())
	}
}
