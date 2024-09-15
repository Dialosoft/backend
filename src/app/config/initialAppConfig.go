package config

import (
	"github.com/Dialosoft/src/adapters/http/controller"
	"github.com/Dialosoft/src/adapters/http/router"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/services"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// Setup for the api
//
// repositories -> services -> controllers -> routers
func SetupAPI(db *gorm.DB, fiberConfigs fiber.Config) *fiber.App {

	app := fiber.New(fiberConfigs)

	// Repositories
	userRepository := repository.NewUserRepository(db)
	roleRepository := repository.NewRoleRepository(db)

	// Services
	userService := services.NewUserService(userRepository)
	_ = services.NewRoleRepository(roleRepository)

	// Controllers
	userController := controller.NewUserController(userService)

	// Routers
	userRouter := router.NewUserRouter(userController)

	userRouter.SetupUserRoutes(app)

	return app
}
