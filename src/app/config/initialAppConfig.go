package config

import (
	"github.com/Dialosoft/src/adapters/http/controller"
	"github.com/Dialosoft/src/adapters/http/middleware"
	"github.com/Dialosoft/src/adapters/http/router"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/services"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Setup for the api
//
// repositories -> services -> controllers -> routers
func SetupAPI(db *gorm.DB, redisConn *redis.Client, generalConfig GeneralConfig, defaultRoles map[string]uuid.UUID) *fiber.App {

	app := fiber.New(fiber.Config{})

	api := app.Group("/dialosoft-api/v1")

	// Repositories
	userRepository := repository.NewUserRepository(db)
	roleRepository := repository.NewRoleRepository(db)
	tokenRepository := repository.NewTokenRepository(db)
	cacheRepository := repository.NewRedisRepository(redisConn)
	forumRepository := repository.NewForumRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)
	postRepository := repository.NewPostRepository(db)
	postLikesRepository := repository.NewPostLikesRepository(db)

	// Services
	cacheService := services.NewCacheService(cacheRepository)
	userService := services.NewUserService(userRepository, roleRepository)
	authService := services.NewAuthService(userRepository, roleRepository, tokenRepository, cacheService, generalConfig.JWTKey)
	forumService := services.NewForumService(forumRepository)
	categoryService := services.NewCategoryService(categoryRepository, roleRepository)
	roleService := services.NewRoleRepository(roleRepository)
	postService := services.NewPostService(postRepository, postLikesRepository, userRepository)

	// Middlewares
	securityMiddleware := middleware.NewSecurityMiddleware(authService, cacheService, generalConfig.JWTKey)

	// Controllers
	userController := controller.NewUserController(userService)
	authController := controller.NewAuthController(authService)
	forumController := controller.NewForumController(forumService)
	categoryController := controller.NewCategoryController(categoryService)
	roleController := controller.NewRoleController(roleService)
	postController := controller.NewPostController(postService)
	managementController := controller.NewManagamentController(
		forumService,
		categoryService,
		roleService,
		userService,
		authService,
		cacheService)

	// Routers
	userRouter := router.NewUserRouter(userController)
	authRouter := router.NewAuthRouter(authController)
	forumRouter := router.NewForumRouter(forumController)
	categoryRouter := router.NewCategoryRouter(categoryController)
	roleRouter := router.NewRoleRouter(roleController)
	managementRouter := router.NewManagementRouter(managementController)
	postRouter := router.NewPostRouter(postController)

	userRouter.SetupUserRoutes(api, securityMiddleware, defaultRoles)
	authRouter.SetupAuthRoutes(api, securityMiddleware)
	forumRouter.SetupForumRoutes(api, securityMiddleware, defaultRoles)
	categoryRouter.SetupCategoryRoutes(api, securityMiddleware, defaultRoles)
	roleRouter.SetupRoleRouter(api, securityMiddleware, defaultRoles)
	managementRouter.SetupManagementRoutes(api, securityMiddleware, defaultRoles)
	postRouter.SetupPostRoutes(api, securityMiddleware, defaultRoles)

	return app
}
