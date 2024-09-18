package router

import (
	"github.com/Dialosoft/src/adapters/http/controller"
	"github.com/gofiber/fiber/v3"
)

type AuthRouter struct {
	AuthController *controller.AuthController
}

func NewAuthRouter(authRouter *controller.AuthController) *AuthRouter {
	return &AuthRouter{AuthController: authRouter}
}

func (r *AuthRouter) SetupAuthRoutes(api fiber.Router) {
	authGroup := api.Group("/auth")
	{
		authGroup.Post("/register", r.AuthController.Register)
		authGroup.Post("/login", r.AuthController.Login)
	}
}
