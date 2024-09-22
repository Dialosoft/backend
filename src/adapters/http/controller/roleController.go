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

type RoleController struct {
	RoleService services.RoleService
}

func NewRoleController(roleService services.RoleService) *RoleController {
	return &RoleController{RoleService: roleService}
}

func (rc *RoleController) GetAllRoles(c fiber.Ctx) error {
	rolesDtos, err := rc.RoleService.GetAllRoles()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Error(err.Error())
			return response.ErrNotFound(c)
		}
		logger.Error(err.Error())
		return response.ErrInternalServer(c)
	}

	return c.Status(fiber.StatusOK).JSON(response.Standard(c, "OK", rolesDtos))
}

func (rc *RoleController) GetRoleByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		logger.Error("empty parameters or arguments")
		return response.ErrEmptyParametersOrArguments(c)
	}

	roleUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c)
	}

	roleDto, err := rc.RoleService.GetRoleByID(roleUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return c.Status(fiber.StatusOK).JSON(response.Standard(c, "OK", roleDto))
}

func (rc *RoleController) GetRoleByType(c fiber.Ctx) error {
	roleType := c.Params("type")
	if roleType == "" {
		logger.Error("empty parameters or arguments")
		return response.ErrEmptyParametersOrArguments(c)
	}

	roleDto, err := rc.RoleService.GetRoleByType(roleType)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return c.Status(fiber.StatusOK).JSON(response.Standard(c, "OK", roleDto))
}

func (rc *RoleController) CreateNewRole(c fiber.Ctx) error {
	var req request.NewRole
	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c)
	}

	if req.RoleType == nil {
		return response.ErrEmptyParametersOrArguments(c)
	}

	roleDto := dto.RoleDto{
		RoleType:   *req.RoleType,
		Permission: *req.Permission,
		AdminRole:  *req.AdminRole,
		ModRole:    *req.ModRole,
	}

	roleUUID, err := rc.RoleService.CreateNewRole(roleDto)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) ||
			strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return response.ErrConflict(c)
		}
		return response.ErrInternalServer(c)
	}

	return c.Status(fiber.StatusOK).JSON(response.StandardCreated(c, "CREATED", fiber.Map{
		"id": roleUUID.String(),
	}))
}

func (rc *RoleController) UpdateRole(c fiber.Ctx) error {
	var req request.NewRole

	id := c.Params("id")
	if id == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	roleUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c)
	}

	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c)
	}

	err = rc.RoleService.UpdateRole(roleUUID, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return c.Status(fiber.StatusOK).JSON(response.Standard(c, "UPDATED", nil))
}

func (rc *RoleController) DeleteRole(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	roleUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c)
	}

	if err = rc.RoleService.DeleteRole(roleUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return c.Status(fiber.StatusOK).JSON(response.Standard(c, "DELETED", nil))
}

func (rc *RoleController) RestoreRole(c fiber.Ctx) error {
	id := c.Params("id")
	roleUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c)
	}

	if err = rc.RoleService.RestoreRole(roleUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return c.Status(fiber.StatusOK).JSON(response.Standard(c, "RESTORED", nil))
}
