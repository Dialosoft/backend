package router

import (
	"github.com/Dialosoft/src/adapters/http/controller"
	"github.com/Dialosoft/src/adapters/http/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CategoryRouter struct {
	CategoryController *controller.CategoryController
}

func NewCategoryRouter(categoryController *controller.CategoryController) *CategoryRouter {
	return &CategoryRouter{CategoryController: categoryController}
}

func (r *CategoryRouter) SetupCategoryRoutes(api fiber.Router, middlewares *middleware.SecurityMiddleware, defaultRoles map[string]uuid.UUID) {
	categoryGroup := api.Group("/categories")

	{
		categoryGroup.Get("/get-all-categories", r.CategoryController.GetAllCategories)
		categoryGroup.Get("/get-category-by-id/:id", r.CategoryController.GetCategoryByID)
		categoryGroup.Get("/get-category-by-name/:name", r.CategoryController.GetCategoryByName)
		categoryGroup.Post("/create-new-category", r.CategoryController.CreateNewCategory,
			middlewares.GetAndVerifyAccesToken(), middlewares.VerifyRefreshToken(), middlewares.RoleRequiredByID(defaultRoles["administrator"].String()),
		)
		categoryGroup.Put("/update-category/:id", r.CategoryController.UpdateCategory, middlewares.VerifyRefreshToken(),
			middlewares.GetAndVerifyAccesToken(), middlewares.VerifyRefreshToken(), middlewares.RoleRequiredByID(defaultRoles["administrator"].String()),
		)
		categoryGroup.Delete("/delete-category/:id", r.CategoryController.DeleteCategory,
			middlewares.GetAndVerifyAccesToken(), middlewares.VerifyRefreshToken(), middlewares.RoleRequiredByID(defaultRoles["administrator"].String()),
		)
		categoryGroup.Put("/restore-category/:id", r.CategoryController.RestoreCategory,
			middlewares.GetAndVerifyAccesToken(), middlewares.VerifyRefreshToken(), middlewares.RoleRequiredByID(defaultRoles["administrator"].String()),
		)
	}
}
