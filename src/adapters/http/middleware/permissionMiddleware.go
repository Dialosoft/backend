package middleware

import (
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/models"
	"github.com/Dialosoft/src/domain/services"
	"github.com/Dialosoft/src/pkg/errorsUtils"
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
	Layer        string
}

func NewPermissionMiddleware(authService services.AuthService, cacheService services.CacheService, roleService services.RoleService, jwtKey string, Layer string) *PermissionMiddleware {
	return &PermissionMiddleware{AuthService: authService, CacheService: cacheService, RoleService: roleService, JwtKey: jwtKey, Layer: Layer}
}

// CanManageCategories checks if the user has the permission to manage categories.
func (sm *PermissionMiddleware) CanManageCategories() fiber.Handler {
	return func(c fiber.Ctx) error {
		rolePermission, err := sm.processBeforeCheckPermissionHelper(c)
		if err != nil {
			return err
		}

		if rolePermission.CanManageCategories {
			return c.Next()
		}

		return response.ErrForbidden(c, sm.Layer)
	}
}

// CanManageForums checks if the user has the permission to manage forums.
func (sm *PermissionMiddleware) CanManageForums() fiber.Handler {
	return func(c fiber.Ctx) error {
		rolePermission, err := sm.processBeforeCheckPermissionHelper(c)
		if err != nil {
			return err
		}

		if rolePermission.CanManageForums {
			return c.Next()
		}

		return response.ErrForbidden(c, sm.Layer)
	}
}

// CanManageRoles checks if the user has the permission to manage roles.
func (sm *PermissionMiddleware) CanManageRoles() fiber.Handler {
	return func(c fiber.Ctx) error {
		rolePermission, err := sm.processBeforeCheckPermissionHelper(c)
		if err != nil {
			return err
		}

		if rolePermission.CanManageRoles {
			return c.Next()
		}

		return response.ErrForbidden(c, sm.Layer)
	}
}

// CanManageUsers checks if the user has the permission to manage users.
func (sm *PermissionMiddleware) CanManageUsers() fiber.Handler {
	return func(c fiber.Ctx) error {
		rolePermission, err := sm.processBeforeCheckPermissionHelper(c)
		if err != nil {
			return err
		}

		if rolePermission.CanManageUsers {
			return c.Next()
		}

		return response.ErrForbidden(c, sm.Layer)
	}
}

// CanManagePosts checks if the user has the permission to manage posts.
func (sm *PermissionMiddleware) CanManagePosts() fiber.Handler {
	return func(c fiber.Ctx) error {
		rolePermission, err := sm.processBeforeCheckPermissionHelper(c)
		if err != nil {
			return err
		}

		if rolePermission.CanManagePosts {
			return c.Next()
		}

		logger.Warn("Insufficient role permissions", map[string]interface{}{
			"route": c.Path(),
		})
		return response.ErrForbidden(c, sm.Layer)
	}
}

// processBeforeCheckPermissionHelper is fuction helper to checks if the user has the permission to access the resource.
// It retrieves the role permissions for the user based on their role ID and checks if the user has the required permission.
// If the user has the required permission, it returns the role permissions.
// If the user does not have the required permission, it returns an error.
func (sm *PermissionMiddleware) processBeforeCheckPermissionHelper(c fiber.Ctx) (*models.RolePermissions, error) {
	roleID := c.Locals("roleID")
	if roleID == "" || roleID == nil {
		return nil, response.ErrEmptyParametersOrArguments(c)
	}

	roleIDString, ok := roleID.(string)
	if !ok {
		return nil, response.ErrInternalServer(c, errorsUtils.ErrInternalServer, nil, sm.Layer)
	}

	roleUUID, err := uuid.Parse(roleIDString)
	if err != nil {
		return nil, response.ErrUUIDParse(c, roleIDString)
	}

	rolePermission, err := sm.RoleService.GetRolePermissionsByRoleID(roleUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, response.ErrNotFound(c, sm.Layer)
		}
		logger.Error("Error retrieving role information", map[string]interface{}{
			"roleID": roleIDString,
			"route":  c.Path(),
		})
		return nil, response.ErrInternalServer(c, err, nil, sm.Layer)
	}

	return rolePermission, nil
}
