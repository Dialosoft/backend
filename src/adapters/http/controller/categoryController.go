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
			logger.Error(err.Error())
			return response.ErrNotFound(c)
		}
		logger.Error(err.Error())
		return response.ErrInternalServer(c)
	}

	return response.Standard(c, "OK", categoriesDtos)
}

func (ac *CategoryController) GetCategoryByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		logger.Error("empty parameters or arguments")
		return response.ErrEmptyParametersOrArguments(c)
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c)
	}

	categoryDto, err := ac.CategoryService.GetCategoryByID(uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return response.Standard(c, "OK", categoryDto)
}

func (ac *CategoryController) GetCategoryByName(c fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		logger.Error("empty parameters or arguments")
		return response.ErrEmptyParametersOrArguments(c)
	}

	categoryDto, err := ac.CategoryService.GetCategoryByName(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return response.Standard(c, "OK", categoryDto)
}

func (ac *CategoryController) CreateNewCategory(c fiber.Ctx) error {
	var req request.NewCategory
	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c)
	}

	if req.Name == nil || req.Description == nil {
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
			return response.ErrConflict(c)
		}
		return response.ErrInternalServer(c)
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
		return response.ErrUUIDParse(c)
	}

	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c)
	}

	err = ac.CategoryService.UpdateCategory(categoryUUID, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return response.Standard(c, "UPDATED", nil)
}

func (ac *CategoryController) DeleteCategory(c fiber.Ctx) error {
	id := c.Params("id")

	categoryUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c)
	}

	if err = ac.CategoryService.DeleteCategory(categoryUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return response.Standard(c, "DELETED", nil)
}

func (ac *CategoryController) RestoreCategory(c fiber.Ctx) error {
	id := c.Params("id")
	categoryUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c)
	}

	if err = ac.CategoryService.RestoreCategory(categoryUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return response.Standard(c, "RESTORED", nil)
}
