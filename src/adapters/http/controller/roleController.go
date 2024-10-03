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
		if err == gorm.ErrRecordNotFound || len(rolesDtos) == 0 {
			logger.Warn("No roles found", map[string]interface{}{
				"route":  c.Path(),
				"method": c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error retrieving all roles", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Roles retrieved successfully", map[string]interface{}{
		"route":  c.Path(),
		"method": c.Method(),
		"count":  len(rolesDtos),
	})

	return response.Standard(c, "OK", rolesDtos)
}

func (rc *RoleController) GetRoleByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		logger.Error("Empty parameters or arguments", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrEmptyParametersOrArguments(c)
	}

	roleUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error("Invalid UUID format", map[string]interface{}{
			"provided-id": id,
			"route":       c.Path(),
			"method":      c.Method(),
		})
		return response.ErrUUIDParse(c)
	}

	roleDto, err := rc.RoleService.GetRoleByID(roleUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("Role not found", map[string]interface{}{
				"roleID": id,
				"route":  c.Path(),
				"method": c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error retrieving role by ID", map[string]interface{}{
			"roleID": id,
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Role retrieved successfully", map[string]interface{}{
		"roleID": id,
		"route":  c.Path(),
		"method": c.Method(),
	})

	return response.Standard(c, "OK", roleDto)
}

func (rc *RoleController) GetRoleByType(c fiber.Ctx) error {
	roleType := c.Params("type")
	if roleType == "" {
		logger.Error("Empty parameters or arguments", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrEmptyParametersOrArguments(c)
	}

	roleDto, err := rc.RoleService.GetRoleByType(roleType)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("Role not found by type", map[string]interface{}{
				"roleType": roleType,
				"route":    c.Path(),
				"method":   c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error retrieving role by type", map[string]interface{}{
			"roleType": roleType,
			"route":    c.Path(),
			"method":   c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Role retrieved successfully by type", map[string]interface{}{
		"roleType": roleType,
		"route":    c.Path(),
		"method":   c.Method(),
	})

	return response.Standard(c, "OK", roleDto)
}

func (rc *RoleController) CreateNewRole(c fiber.Ctx) error {
	var req request.NewRole
	if err := c.Bind().Body(&req); err != nil {
		logger.CaptureError(err, "Failed to bind request for creating new role", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrBadRequest(c)
	}

	if *req.RoleType == "" {
		logger.Error("RoleType is missing", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
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
			logger.Warn("Duplicate role detected", map[string]interface{}{
				"roleDto": roleDto,
				"route":   c.Path(),
				"method":  c.Method(),
			})
			return response.ErrConflict(c)
		}
		logger.CaptureError(err, "Error creating new role", map[string]interface{}{
			"roleDto": roleDto,
			"route":   c.Path(),
			"method":  c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Role created successfully", map[string]interface{}{
		"roleID": roleUUID.String(),
		"route":  c.Path(),
		"method": c.Method(),
	})

	return response.StandardCreated(c, "CREATED", fiber.Map{
		"id": roleUUID.String(),
	})
}

func (rc *RoleController) UpdateRole(c fiber.Ctx) error {
	var req request.NewRole

	id := c.Params("id")
	if id == "" {
		logger.Error("Empty parameters or arguments", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrEmptyParametersOrArguments(c)
	}

	roleUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error("Invalid UUID format", map[string]interface{}{
			"provided-id": id,
			"route":       c.Path(),
			"method":      c.Method(),
		})
		return response.ErrUUIDParse(c)
	}

	if err := c.Bind().Body(&req); err != nil {
		logger.CaptureError(err, "Failed to bind request for updating role", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrBadRequest(c)
	}

	err = rc.RoleService.UpdateRole(roleUUID, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("Role not found for update", map[string]interface{}{
				"roleID": id,
				"route":  c.Path(),
				"method": c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error updating role", map[string]interface{}{
			"roleID": id,
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Role updated successfully", map[string]interface{}{
		"roleID": id,
		"route":  c.Path(),
		"method": c.Method(),
	})

	return response.Standard(c, "UPDATED", nil)
}

func (rc *RoleController) SetRolePermissionsByRoleID(c fiber.Ctx) error {
	var req request.NewRolePermissions
	id := c.Params("id")
	if id == "" {
		logger.Error("Empty parameters or arguments", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrEmptyParametersOrArguments(c)
	}

	roleUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error("Invalid UUID format", map[string]interface{}{
			"provided-id": id,
			"route":       c.Path(),
			"method":      c.Method(),
		})
		return response.ErrUUIDParse(c)
	}

	if err := c.Bind().Body(&req); err != nil {
		logger.CaptureError(err, "Failed to bind request for updating role", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrBadRequest(c)
	}

	err = rc.RoleService.SetRolePermissionsByRoleID(roleUUID, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("Role not found for update", map[string]interface{}{
				"roleID": id,
				"route":  c.Path(),
				"method": c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error updating role", map[string]interface{}{
			"roleID": id,
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Role updated successfully", map[string]interface{}{
		"roleID": id,
		"route":  c.Path(),
		"method": c.Method(),
	})

	return response.Standard(c, "UPDATED", nil)
}

func (rc *RoleController) GetRolePermissionsByRoleID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		logger.Error("Empty parameters or arguments", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrEmptyParametersOrArguments(c)
	}

	roleUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error("Invalid UUID format", map[string]interface{}{
			"provided-id": id,
			"route":       c.Path(),
			"method":      c.Method(),
		})
		return response.ErrUUIDParse(c)
	}

	rolePermission, err := rc.RoleService.GetRolePermissionsByRoleID(roleUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("Role not found", map[string]interface{}{
				"roleID": id,
				"route":  c.Path(),
				"method": c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error retrieving role by ID", map[string]interface{}{
			"roleID": id,
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Role retrieved successfully", map[string]interface{}{
		"roleID": id,
		"route":  c.Path(),
		"method": c.Method(),
	})

	return response.Standard(c, "OK", rolePermission)
}

func (rc *RoleController) DeleteRole(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		logger.Error("Empty parameters or arguments", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrEmptyParametersOrArguments(c)
	}

	roleUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error("Invalid UUID format", map[string]interface{}{
			"provided-id": id,
			"route":       c.Path(),
			"method":      c.Method(),
		})
		return response.ErrUUIDParse(c)
	}

	if err = rc.RoleService.DeleteRole(roleUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("Role not found for deletion", map[string]interface{}{
				"roleID": id,
				"route":  c.Path(),
				"method": c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error deleting role", map[string]interface{}{
			"roleID": id,
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Role deleted successfully", map[string]interface{}{
		"roleID": id,
		"route":  c.Path(),
		"method": c.Method(),
	})

	return response.Standard(c, "DELETED", nil)
}

func (rc *RoleController) RestoreRole(c fiber.Ctx) error {
	id := c.Params("id")
	roleUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error("Invalid UUID format", map[string]interface{}{
			"provided-id": id,
			"route":       c.Path(),
			"method":      c.Method(),
		})
		return response.ErrUUIDParse(c)
	}

	if err = rc.RoleService.RestoreRole(roleUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("Role not found for restoration", map[string]interface{}{
				"roleID": id,
				"route":  c.Path(),
				"method": c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error restoring role", map[string]interface{}{
			"roleID": id,
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Role restored successfully", map[string]interface{}{
		"roleID": id,
		"route":  c.Path(),
		"method": c.Method(),
	})

	return response.Standard(c, "RESTORED", nil)
}
