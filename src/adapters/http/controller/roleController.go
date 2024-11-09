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

// GetAllRoles retrieves all roles in the system.
//
// @Summary Get all roles
// @Description Retrieves a list of all roles available in the system.
// @Tags Roles
// @Accept json
// @Produce json
// @Success 200 {array} response.RoleResponse "List of roles"
// @Failure 404 {object} response.StandardError "NOT FOUND - No roles found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Router /roles/get-all-roles [get]
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

// GetRoleByID retrieves a role by its unique ID.
//
// @Summary Get role by ID
// @Description Retrieves role details by its unique identifier.
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} response.RoleResponse "Role data"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid ID format"
// @Failure 404 {object} response.StandardError "NOT FOUND - Role not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Router /roles/get-role-by-id/{id} [get]
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

// GetRoleByType retrieves a role by its type.
//
// @Summary Get role by type
// @Description Retrieves role details based on its type.
// @Tags Roles
// @Accept json
// @Produce json
// @Param type path string true "Role type"
// @Success 200 {object} response.RoleResponse "Role data"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid role type"
// @Failure 404 {object} response.StandardError "NOT FOUND - Role not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Router /roles/get-role-by-type/{type} [get]
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

// CreateNewRole creates a new role in the system.
//
// @Summary Create new role
// @Description Creates a new role with the specified data.
// @Tags Roles
// @Accept json
// @Produce json
// @Param role body request.NewRole true "New Role Data"
// @Success 201 {object} string "CREATED - New role created successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid input data"
// @Failure 409 {object} response.StandardError "CONFLICT - Role already exists"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /roles/protected/create-new-role [post]
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

// UpdateRole updates an existing role's details.
//
// @Summary Update role
// @Description Updates an existing role with new data by its ID.
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param role body request.NewRole true "Updated Role Data"
// @Success 200 {object} response.StandardResponse "UPDATED - Role updated successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid input data"
// @Failure 404 {object} response.StandardError "NOT FOUND - Role not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /roles/protected/update-role/{id} [put]
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

// SetRolePermissionsByRoleID updates the permissions of a role by its ID.
//
// @Summary Set role permissions by ID
// @Description Updates the permissions for a role by its unique identifier.
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param permissions body request.NewRolePermissions true "Updated Permissions Data"
// @Success 200 {object} response.StandardResponse "UPDATED - Role permissions updated successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid input data"
// @Failure 404 {object} response.StandardError "NOT FOUND - Role not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /roles/protected/set-role-permissions-by-id/{id} [put]
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

// GetRolePermissionsByRoleID retrieves the permissions associated with a role.
//
// @Summary Get role permissions by ID
// @Description Retrieves permissions for a specific role by its ID.
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} response.RolePermissionsResponse "Role permissions"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid ID format"
// @Failure 404 {object} response.StandardError "NOT FOUND - Role or permissions not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /roles/protected/get-role-permissions-by-id/{id} [get]
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

// DeleteRole deletes a role by its ID.
//
// @Summary Delete role
// @Description Deletes an existing role by its unique identifier.
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} response.StandardResponse "DELETED - Role deleted successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid role ID"
// @Failure 404 {object} response.StandardError "NOT FOUND - Role not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /roles/protected/delete-role/{id} [delete]
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

// RestoreRole restores a previously deleted role by its ID.
//
// @Summary Restore role
// @Description Restores a previously deleted role by its unique identifier.
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} response.StandardResponse "RESTORED - Role restored successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid role ID"
// @Failure 404 {object} response.StandardError "NOT FOUND - Role not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /roles/protected/restore-role/{id} [put]
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
