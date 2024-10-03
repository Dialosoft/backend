package middleware

import (
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/models"
	"github.com/Dialosoft/src/domain/services"
	"github.com/Dialosoft/src/pkg/utils/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionMiddleware struct {
	AuthService  services.AuthService
	CacheService services.CacheService
	ForumService services.ForumService
	RoleService  services.RoleService
	JwtKey       string
}

func NewPermissionMiddleware(authService services.AuthService, cacheService services.CacheService, roleService services.RoleService, jwtKey string) *PermissionMiddleware {
	return &PermissionMiddleware{AuthService: authService, CacheService: cacheService, RoleService: roleService, JwtKey: jwtKey}
}

func (sm *PermissionMiddleware) CanManageCategories() fiber.Handler {
	return func(c fiber.Ctx) error {
		rolePermission, err := sm.processBeforeCheckPermissionHelper(c)
		if err != nil {
			return err
		}

		if rolePermission.CanManageCategories {
			return c.Next()
		}

		logger.Warn("Insufficient role permissions", map[string]interface{}{
			"route": c.Path(),
		})
		return response.ErrForbidden(c)
	}
}

func (sm *PermissionMiddleware) CanManageForums() fiber.Handler {
	return func(c fiber.Ctx) error {
		rolePermission, err := sm.processBeforeCheckPermissionHelper(c)
		if err != nil {
			return err
		}

		if rolePermission.CanManageForums {
			return c.Next()
		}

		logger.Warn("Insufficient role permissions", map[string]interface{}{
			"route": c.Path(),
		})
		return response.ErrForbidden(c)
	}
}

func (sm *PermissionMiddleware) CanManageRoles() fiber.Handler {
	return func(c fiber.Ctx) error {
		rolePermission, err := sm.processBeforeCheckPermissionHelper(c)
		if err != nil {
			return err
		}

		if rolePermission.CanManageRoles {
			return c.Next()
		}

		logger.Warn("Insufficient role permissions", map[string]interface{}{
			"route": c.Path(),
		})
		return response.ErrForbidden(c)
	}
}

func (sm *PermissionMiddleware) CanManageUsers() fiber.Handler {
	return func(c fiber.Ctx) error {
		rolePermission, err := sm.processBeforeCheckPermissionHelper(c)
		if err != nil {
			return err
		}

		if rolePermission.CanManageUsers {
			return c.Next()
		}

		logger.Warn("Insufficient role permissions", map[string]interface{}{
			"route": c.Path(),
		})
		return response.ErrForbidden(c)
	}
}

func (sm *PermissionMiddleware) processBeforeCheckPermissionHelper(c fiber.Ctx) (*models.RolePermissions, error) {
	roleID := c.Locals("roleID")
	if roleID == "" || roleID == nil {
		logger.Warn("RoleID missing in context", map[string]interface{}{
			"route": c.Path(),
		})
		return nil, response.ErrEmptyParametersOrArguments(c)
	}

	roleIDString, ok := roleID.(string)
	if !ok {
		logger.Error("Invalid roleID format in token", map[string]interface{}{
			"roleID": roleID,
			"route":  c.Path(),
		})
		return nil, response.ErrInternalServer(c)
	}

	roleUUID, err := uuid.Parse(roleIDString)
	if err != nil {
		logger.Error("Invalid UUID format", map[string]interface{}{
			"provided-id": roleIDString,
			"route":       c.Path(),
		})
		return nil, response.ErrUUIDParse(c)
	}

	rolePermission, err := sm.RoleService.GetRolePermissionsByRoleID(roleUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("Role not found", map[string]interface{}{
				"roleID": roleIDString,
				"route":  c.Path(),
			})
			return nil, response.ErrNotFound(c)
		}
		logger.Error("Error retrieving role information", map[string]interface{}{
			"roleID": roleIDString,
			"route":  c.Path(),
		})
		return nil, response.ErrInternalServer(c)
	}

	return rolePermission, nil
}
