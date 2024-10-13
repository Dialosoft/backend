package controller

import (
	"errors"
	"strings"

	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/services"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleController struct {
	RoleService services.RoleService
	Layer       string
}

func NewRoleController(roleService services.RoleService, Layer string) *RoleController {
	return &RoleController{RoleService: roleService, Layer: Layer}
}

func (rc *RoleController) GetAllRoles(c fiber.Ctx) error {
	rolesDtos, err := rc.RoleService.GetAllRoles()
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(rolesDtos) == 0 {
			return response.ErrNotFound(c, rc.Layer)
		}
		return response.ErrInternalServer(c, err, rolesDtos, rc.Layer)
	}

	return response.Standard(c, "OK", rolesDtos)
}

func (rc *RoleController) GetRoleByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	roleUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c, id)
	}

	roleDto, err := rc.RoleService.GetRoleByID(roleUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, rc.Layer)
		}
		return response.ErrInternalServer(c, err, roleDto, rc.Layer)
	}

	return response.Standard(c, "OK", roleDto)
}

func (rc *RoleController) GetRoleByType(c fiber.Ctx) error {
	roleType := c.Params("type")
	if roleType == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	roleDto, err := rc.RoleService.GetRoleByType(roleType)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, rc.Layer)
		}
		return response.ErrInternalServer(c, err, roleDto, rc.Layer)
	}

	return response.Standard(c, "OK", roleDto)
}

func (rc *RoleController) CreateNewRole(c fiber.Ctx) error {
	var req request.NewRole
	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, rc.Layer)
	}

	if *req.RoleType == "" {
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
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return response.ErrConflict(c, err, roleDto, rc.Layer)
		}
		return response.ErrInternalServer(c, err, roleDto, rc.Layer)
	}

	return response.StandardCreated(c, "CREATED", fiber.Map{
		"id": roleUUID.String(),
	})
}

func (rc *RoleController) UpdateRole(c fiber.Ctx) error {
	var req request.NewRole

	id := c.Params("id")
	if id == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	roleUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c, id)
	}

	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, rc.Layer)
	}

	err = rc.RoleService.UpdateRole(roleUUID, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, rc.Layer)
		}
		return response.ErrInternalServer(c, err, nil, rc.Layer)
	}

	return response.Standard(c, "UPDATED", nil)
}

func (rc *RoleController) SetRolePermissionsByRoleID(c fiber.Ctx) error {
	var req request.NewRolePermissions
	id := c.Params("id")
	if id == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	roleUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c, id)
	}

	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, rc.Layer)
	}

	err = rc.RoleService.SetRolePermissionsByRoleID(roleUUID, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, rc.Layer)
		}
		return response.ErrInternalServer(c, err, nil, rc.Layer)
	}

	return response.Standard(c, "UPDATED", nil)
}

func (rc *RoleController) GetRolePermissionsByRoleID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	roleUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c, id)
	}

	rolePermission, err := rc.RoleService.GetRolePermissionsByRoleID(roleUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, rc.Layer)
		}
		return response.ErrInternalServer(c, err, rolePermission, rc.Layer)
	}

	return response.Standard(c, "OK", rolePermission)
}

func (rc *RoleController) DeleteRole(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	roleUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c, id)
	}

	if err = rc.RoleService.DeleteRole(roleUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, rc.Layer)
		}
		return response.ErrInternalServer(c, err, nil, rc.Layer)
	}

	return response.Standard(c, "DELETED", nil)
}

func (rc *RoleController) RestoreRole(c fiber.Ctx) error {
	id := c.Params("id")
	roleUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c, id)
	}

	if err = rc.RoleService.RestoreRole(roleUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, rc.Layer)
		}
		return response.ErrInternalServer(c, err, nil, rc.Layer)
	}

	return response.Standard(c, "RESTORED", nil)
}
