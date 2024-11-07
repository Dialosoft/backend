package config

import (
	"github.com/Dialosoft/src/adapters/http/controller"
	"github.com/Dialosoft/src/adapters/http/middleware"
	"github.com/Dialosoft/src/adapters/http/router"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/services"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Setup for the api
//
// repositories -> services -> controllers -> routers -> Setups for routes
func SetupAPI(db *gorm.DB, redisConn *redis.Client, generalConfig GeneralConfig) *fiber.App {
	app := fiber.New(fiber.Config{})

	api := app.Group("/dialosoft-api/v1")
	validate := validator.New()
	// Repositories
	userRepository := repository.NewUserRepository(db)
	roleRepository := repository.NewRoleRepository(db)
	tokenRepository := repository.NewTokenRepository(db)
	cacheRepository := repository.NewRedisRepository(redisConn)
	forumRepository := repository.NewForumRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)
	postRepository := repository.NewPostRepository(db)
	postLikesRepository := repository.NewPostLikesRepository(db)
	rolePermissionsRepository := repository.NewRolePermissionsRepository(db)

	// Services
	cacheService := services.NewCacheService(cacheRepository)
	userService := services.NewUserService(userRepository, roleRepository)
	authService := services.NewAuthService(userRepository, roleRepository, tokenRepository, cacheService, generalConfig.JWTKey)
	forumService := services.NewForumService(forumRepository, categoryRepository)
	categoryService := services.NewCategoryService(categoryRepository, roleRepository)
	roleService := services.NewRoleRepository(roleRepository, rolePermissionsRepository)
	postService := services.NewPostService(postRepository, postLikesRepository, userRepository)

	// Middlewares
	securityMiddleware := middleware.NewSecurityMiddleware(authService, cacheService, generalConfig.JWTKey, "Middleware/SecurityMiddleware")
	permissionMiddleware := middleware.NewPermissionMiddleware(authService, cacheService, roleService, generalConfig.JWTKey, "Middleware/PermissionMiddleware")

	// Controllers
	userController := controller.NewUserController(userService, "Controller/UserController")
	authController := controller.NewAuthController(authService, validate, "Controller/AuthController")
	forumController := controller.NewForumController(forumService, "Controller/ForumController")
	categoryController := controller.NewCategoryController(categoryService, "Controller/CategoryController")
	roleController := controller.NewRoleController(roleService, "Controller/RoleController")
	postController := controller.NewPostController(postService, "Controller/PostController")
	managementController := controller.NewManagamentController(
		forumService,
		categoryService,
		roleService,
		userService,
		authService,
		cacheService,
		"Controller/ManagementController")

	// Routers
	userRouter := router.NewUserRouter(userController)
	authRouter := router.NewAuthRouter(authController)
	forumRouter := router.NewForumRouter(forumController)
	categoryRouter := router.NewCategoryRouter(categoryController)
	roleRouter := router.NewRoleRouter(roleController)
	managementRouter := router.NewManagementRouter(managementController)
	postRouter := router.NewPostRouter(postController)

	userRouter.SetupUserRoutes(api, securityMiddleware, permissionMiddleware)
	authRouter.SetupAuthRoutes(api, securityMiddleware)
	forumRouter.SetupForumRoutes(api, securityMiddleware, permissionMiddleware)
	categoryRouter.SetupCategoryRoutes(api, securityMiddleware, permissionMiddleware)
	roleRouter.SetupRoleRouter(api, securityMiddleware, permissionMiddleware)
	managementRouter.SetupManagementRoutes(api, securityMiddleware, permissionMiddleware)
	postRouter.SetupPostRoutes(api, securityMiddleware, permissionMiddleware)

	return app
}
