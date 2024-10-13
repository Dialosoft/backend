package router

import (
	"github.com/Dialosoft/src/adapters/http/controller"
	"github.com/Dialosoft/src/adapters/http/middleware"
	"github.com/gofiber/fiber/v3"
)

type CategoryRouter struct {
	CategoryController *controller.CategoryController
}

func NewCategoryRouter(categoryController *controller.CategoryController) *CategoryRouter {
	return &CategoryRouter{CategoryController: categoryController}
}

func (r *CategoryRouter) SetupCategoryRoutes(api fiber.Router, securityMiddleware *middleware.SecurityMiddleware, permissionMiddleware *middleware.PermissionMiddleware) {
	categoryGroup := api.Group("/categories")
	protectedGroup := categoryGroup.Group("/protected",
		securityMiddleware.GetAndVerifyAccessToken(),
		securityMiddleware.VerifyRefreshToken(),
		securityMiddleware.GetRoleFromToken(),
		permissionMiddleware.CanManageCategories() /* permission middleware for categories */)

	{
		// public

		categoryGroup.Get("/get-all-categories-allowed", r.CategoryController.GetAllCategoriesAllowedByRole, securityMiddleware.GetRoleFromToken())
	}

	{
		// protected routes by authenticated users and with permission

		// categoryGroup.Get("/get-all-categories", r.CategoryController.GetAllCategories)
		// categoryGroup.Get("/get-category-by-id/:id", r.CategoryController.GetCategoryByID)
		// categoryGroup.Get("/get-category-by-name/:name", r.CategoryController.GetCategoryByName)

		protectedGroup.Post("/create-new-category", r.CategoryController.CreateNewCategory)
		protectedGroup.Put("/update-category/:id", r.CategoryController.UpdateCategory)
		protectedGroup.Delete("/delete-category/:id", r.CategoryController.DeleteCategory)
		protectedGroup.Put("/restore-category/:id", r.CategoryController.RestoreCategory)
	}
}
