package controller

import (
	"errors"
	"strings"

	"github.com/Dialosoft/src/adapters/dto"
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
}

func NewCategoryController(categoryService services.CategoryService) *CategoryController {
	return &CategoryController{CategoryService: categoryService}
}

func (ac *CategoryController) GetAllCategories(c fiber.Ctx) error {
	categoriesDtos, err := ac.CategoryService.GetAllCategories()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("No categories found", map[string]interface{}{
				"route":  c.Path(),
				"method": c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error retrieving categories", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Categories retrieved successfully", map[string]interface{}{
		"route":           c.Path(),
		"method":          c.Method(),
		"categoriesCount": len(categoriesDtos),
	})

	return response.Standard(c, "OK", categoriesDtos)
}

func (ac *CategoryController) GetCategoryByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		logger.Error("Empty parameters or arguments", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrEmptyParametersOrArguments(c)
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		logger.Error("Invalid UUID format", map[string]interface{}{
			"provided-id": id,
			"route":       c.Path(),
			"method":      c.Method(),
		})
		return response.ErrUUIDParse(c)
	}

	categoryDto, err := ac.CategoryService.GetCategoryByID(uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("Category not found", map[string]interface{}{
				"categoryID": id,
				"route":      c.Path(),
				"method":     c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error retrieving category by ID", map[string]interface{}{
			"categoryID": id,
			"route":      c.Path(),
			"method":     c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Category retrieved successfully", map[string]interface{}{
		"categoryID": id,
		"route":      c.Path(),
		"method":     c.Method(),
	})

	return response.Standard(c, "OK", categoryDto)
}

func (ac *CategoryController) GetCategoryByName(c fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		logger.Error("Empty parameters or arguments", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrEmptyParametersOrArguments(c)
	}

	categoryDto, err := ac.CategoryService.GetCategoryByName(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("Category not found", map[string]interface{}{
				"categoryName": name,
				"route":        c.Path(),
				"method":       c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error retrieving category by name", map[string]interface{}{
			"categoryName": name,
			"route":        c.Path(),
			"method":       c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Category retrieved successfully", map[string]interface{}{
		"categoryName": name,
		"route":        c.Path(),
		"method":       c.Method(),
	})

	return response.Standard(c, "OK", categoryDto)
}

func (ac *CategoryController) CreateNewCategory(c fiber.Ctx) error {
	var req request.NewCategory
	if err := c.Bind().Body(&req); err != nil {
		logger.CaptureError(err, "Failed to parse CreateNewCategory request", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrBadRequest(c)
	}

	if req.Name == nil || req.Description == nil {
		logger.Error("Missing parameters in CreateNewCategory request", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrEmptyParametersOrArguments(c)
	}

	categoryDto := dto.CategoryDto{
		Name:        *req.Name,
		Description: *req.Description,
	}

	categoryUUID, err := ac.CategoryService.CreateCategory(categoryDto)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) ||
			strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			logger.Warn("Category creation conflict", map[string]interface{}{
				"categoryDto": categoryDto,
				"route":       c.Path(),
				"method":      c.Method(),
			})
			return response.ErrConflict(c)
		}
		logger.CaptureError(err, "Error creating new category", map[string]interface{}{
			"categoryDto": categoryDto,
			"route":       c.Path(),
			"method":      c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Category created successfully", map[string]interface{}{
		"categoryID": categoryUUID.String(),
		"route":      c.Path(),
		"method":     c.Method(),
	})

	return response.StandardCreated(c, "CREATED", fiber.Map{
		"id": categoryUUID.String(),
	})
}

func (ac *CategoryController) UpdateCategory(c fiber.Ctx) error {
	var req request.NewCategory

	id := c.Params("id")
	if id == "" {
		logger.Error("Empty parameters or arguments", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrEmptyParametersOrArguments(c)
	}

	categoryUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error("Invalid UUID format", map[string]interface{}{
			"provided-id": id,
			"route":       c.Path(),
			"method":      c.Method(),
		})
		return response.ErrUUIDParse(c)
	}

	if err := c.Bind().Body(&req); err != nil {
		logger.CaptureError(err, "Failed to parse UpdateCategory request", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrBadRequest(c)
	}

	err = ac.CategoryService.UpdateCategory(categoryUUID, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("Category not found for update", map[string]interface{}{
				"categoryID": id,
				"route":      c.Path(),
				"method":     c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error updating category", map[string]interface{}{
			"categoryID": id,
			"route":      c.Path(),
			"method":     c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Category updated successfully", map[string]interface{}{
		"categoryID": id,
		"route":      c.Path(),
		"method":     c.Method(),
	})

	return response.Standard(c, "UPDATED", nil)
}

func (ac *CategoryController) DeleteCategory(c fiber.Ctx) error {
	id := c.Params("id")

	categoryUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error("Invalid UUID format", map[string]interface{}{
			"provided-id": id,
			"route":       c.Path(),
			"method":      c.Method(),
		})
		return response.ErrUUIDParse(c)
	}

	if err = ac.CategoryService.DeleteCategory(categoryUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("Category not found for deletion", map[string]interface{}{
				"categoryID": id,
				"route":      c.Path(),
				"method":     c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error deleting category", map[string]interface{}{
			"categoryID": id,
			"route":      c.Path(),
			"method":     c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Category deleted successfully", map[string]interface{}{
		"categoryID": id,
		"route":      c.Path(),
		"method":     c.Method(),
	})

	return response.Standard(c, "DELETED", nil)
}

func (ac *CategoryController) RestoreCategory(c fiber.Ctx) error {
	id := c.Params("id")
	categoryUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error("Invalid UUID format", map[string]interface{}{
			"provided-id": id,
			"route":       c.Path(),
			"method":      c.Method(),
		})
		return response.ErrUUIDParse(c)
	}

	if err = ac.CategoryService.RestoreCategory(categoryUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("Category not found for restoration", map[string]interface{}{
				"categoryID": id,
				"route":      c.Path(),
				"method":     c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error restoring category", map[string]interface{}{
			"categoryID": id,
			"route":      c.Path(),
			"method":     c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Category restored successfully", map[string]interface{}{
		"categoryID": id,
		"route":      c.Path(),
		"method":     c.Method(),
	})

	return response.Standard(c, "RESTORED", nil)
}
