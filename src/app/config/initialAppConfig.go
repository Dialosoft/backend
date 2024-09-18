package config

import (
	"github.com/Dialosoft/src/adapters/http/controller"
	"github.com/Dialosoft/src/adapters/http/router"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/services"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Setup for the api
//
// repositories -> services -> controllers -> routers
func SetupAPI(db *gorm.DB, redisConn *redis.Client, generalConfig GeneralConfig) *fiber.App {

	app := fiber.New(fiber.Config{})

	api := app.Group("/dialosoft-api/v1")

	// Repositories
	userRepository := repository.NewUserRepository(db)
	roleRepository := repository.NewRoleRepository(db)
	tokenRepository := repository.NewTokenRepository(db)
	cacheRepository := repository.NewRedisRepository(redisConn)

	// Services
	userService := services.NewUserService(userRepository, roleRepository)
	authService := services.NewAuthService(userRepository, roleRepository, tokenRepository, cacheRepository, generalConfig.JWTKey)

	// Controllers
	userController := controller.NewUserController(userService)
	authController := controller.NewAuthController(authService)

	// Routers
	userRouter := router.NewUserRouter(userController)
	authRouter := router.NewAuthRouter(authController)

	userRouter.SetupUserRoutes(api)
	authRouter.SetupAuthRoutes(api)

	return app
}
