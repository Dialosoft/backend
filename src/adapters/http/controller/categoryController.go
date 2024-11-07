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
