package controller

import (
	"errors"
	"strings"

	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/services"
	"github.com/Dialosoft/src/pkg/utils/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryController struct {
	CategoryService services.CategoryService
	Layer           string
}

func NewCategoryController(categoryService services.CategoryService, layer string) *CategoryController {
	return &CategoryController{CategoryService: categoryService, Layer: layer}
}

func (ac *CategoryController) GetAllCategories(c fiber.Ctx) error {
	categoriesResponses, err := ac.CategoryService.GetAllCategories()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, ac.Layer)
		}
		return response.ErrInternalServer(c, err, categoriesResponses, ac.Layer)
	}
	return response.Standard(c, "OK", categoriesResponses)
}

// GetAllCategoriesAllowedByRole retrieves all categories accessible to a specific role.
//
// @Summary Get allowed categories by role
// @Description Retrieves a list of categories that the specified role is allowed to access, if is Authenticated can see protected categories by role
// @Tags Categories
// @Accept  json
// @Produce  json
// @Success 200 {array} response.CategoryResponse "List of accessible categories"
// @Failure 403 {object} response.StandardError "FORBIDDEN - Invalid roleID format in token or unauthorized access"
// @Failure 404 {object} response.StandardError "NOT FOUND - No categories found for this role"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Router	/categories/get-all-categories-allowed [get]
func (ac *CategoryController) GetAllCategoriesAllowedByRole(c fiber.Ctx) error {
	roleID := c.Locals("roleID")
	roleIDString, ok := roleID.(string)
	if !ok {
		logger.Error("Invalid roleID format in token", map[string]interface{}{
			"roleID": roleID,
			"route":  c.Path(),
		})
		return response.PersonalizedErr(c, "Error in token: claims", fiber.StatusForbidden)
	}

	categoriesResponses, err := ac.CategoryService.GetAllCategoriesAllowedByRole(roleIDString)
	if err != nil {
		return response.ErrInternalServer(c, err, categoriesResponses, ac.Layer)
	}

	if categoriesResponses == nil {
		return response.ErrNotFound(c, ac.Layer)
	}

	return response.Standard(c, "OK", categoriesResponses)
}

func (ac *CategoryController) GetCategoryByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c, id)
	}

	categoryDto, err := ac.CategoryService.GetCategoryByID(uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, ac.Layer)
		}
		return response.ErrInternalServer(c, err, categoryDto, ac.Layer)
	}

	return response.Standard(c, "OK", categoryDto)
}

func (ac *CategoryController) GetCategoryByName(c fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	categoryDto, err := ac.CategoryService.GetCategoryByName(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, ac.Layer)
		}
		return response.ErrInternalServer(c, err, categoryDto, ac.Layer)
	}

	return response.Standard(c, "OK", categoryDto)
}

// CreateNewCategory creates a new category in the system.
//
// @Summary Create new category
// @Description Creates a new category with a specified name.
// @Tags Categories
// @Accept  json
// @Produce  json
// @Param category body request.NewCategory true "New Category Data"
// @Success 201 {object} string "CREATED - Category created successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid input data"
// @Failure 409 {object} response.StandardError "CONFLICT - Category already exists"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /categories/create-new-category [post]
func (ac *CategoryController) CreateNewCategory(c fiber.Ctx) error {
	var req request.NewCategory
	if err := c.Bind().Body(&req); err != nil {
		body := string(c.Body())
		return response.ErrBadRequest(c, body, err, ac.Layer)
	}

	if req.Name == nil {
		return response.ErrEmptyParametersOrArguments(c)
	}

	categoryUUID, err := ac.CategoryService.CreateCategory(req)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) ||
			strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return response.ErrConflict(c, err, req, ac.Layer)
		}
		return response.ErrInternalServer(c, err, req, ac.Layer)
	}

	return response.StandardCreated(c, "CREATED", fiber.Map{
		"id": categoryUUID.String(),
	})
}

// UpdateCategory updates an existing category based on its ID.
//
// @Summary Update category
// @Description Updates an existing category with new data.
// @Tags Categories
// @Accept  json
// @Produce  json
// @Param id path string true "Category ID"
// @Param category body request.NewCategory true "Updated Category Data"
// @Success 200 {object} response.StandardResponse "UPDATED - Category updated successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid input data"
// @Failure 404 {object} response.StandardError "NOT FOUND - Category not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /categories/protected/update-category/{id} [put]
func (ac *CategoryController) UpdateCategory(c fiber.Ctx) error {
	var req request.NewCategory

	id := c.Params("id")
	if id == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	categoryUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c, id)
	}

	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, ac.Layer)
	}

	err = ac.CategoryService.UpdateCategory(categoryUUID, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, ac.Layer)
		}
		return response.ErrInternalServer(c, err, req, ac.Layer)
	}

	return response.Standard(c, "UPDATED", nil)
}

// DeleteCategory deletes an existing category by its ID.
//
// @Summary Delete category
// @Description Deletes an existing category from the system.
// @Tags Categories
// @Param id path string true "Category ID"
// @Success 200 {object} response.StandardResponse "DELETED - Category deleted successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid category ID"
// @Failure 404 {object} response.StandardError "NOT FOUND - Category not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /categories/protected/delete-category/{id} [delete]
func (ac *CategoryController) DeleteCategory(c fiber.Ctx) error {
	id := c.Params("id")

	categoryUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c, id)
	}

	if err = ac.CategoryService.DeleteCategory(categoryUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, ac.Layer)
		}
		return response.ErrInternalServer(c, err, nil, ac.Layer)
	}

	return response.Standard(c, "DELETED", nil)
}

// RestoreCategory restores a previously deleted category by its ID.
//
// @Summary Restore category
// @Description Restores a previously deleted category.
// @Tags Categories
// @Param id path string true "Category ID"
// @Success 200 {object} response.StandardResponse "RESTORED - Category restored successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid category ID"
// @Failure 404 {object} response.StandardError "NOT FOUND - Category not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /categories/protected/restore-category/{id} [put]
func (ac *CategoryController) RestoreCategory(c fiber.Ctx) error {
	id := c.Params("id")
	categoryUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c, id)
	}

	if err = ac.CategoryService.RestoreCategory(categoryUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, ac.Layer)
		}
		return response.ErrInternalServer(c, err, nil, ac.Layer)
	}

	logger.Info("Category restored successfully", map[string]interface{}{
		"categoryID": id,
		"route":      c.Path(),
		"method":     c.Method(),
	})

	return response.Standard(c, "RESTORED", nil)
}
