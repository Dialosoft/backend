package middleware

import (
	"strings"

	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/services"
	"github.com/Dialosoft/src/pkg/utils/jsonWebToken"
	"github.com/Dialosoft/src/pkg/utils/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type SecurityMiddleware struct {
	AuthService  services.AuthService
	CacheService services.CacheService
	JwtKey       string
}

func NewSecurityMiddleware(authService services.AuthService, cacheService services.CacheService, jwtKey string) *SecurityMiddleware {
	return &SecurityMiddleware{AuthService: authService, CacheService: cacheService, JwtKey: jwtKey}
}

// GetAndVerifyAccessToken retrieves the access token from the Authorization header,
// checks its format, and verifies the token. If the token is invalid, expired, or missing,
// appropriate error responses are returned. If valid, the user's ID and role are extracted
// from the token and stored in the request context for further use.
func (sm *SecurityMiddleware) GetAndVerifyAccessToken() fiber.Handler {
	return func(c fiber.Ctx) error {
		accessTokenHeader := c.Get("Authorization")
		if accessTokenHeader == "" {
			logger.Warn("Authorization header missing", map[string]interface{}{
				"route": c.Path(),
			})
			return response.ErrUnauthorizedHeader(c)
		}

		accessTokenParts := strings.Split(accessTokenHeader, " ")
		if len(accessTokenParts) != 2 || accessTokenParts[0] != "Bearer" {
			logger.Warn("Invalid authorization header format", map[string]interface{}{
				"header": accessTokenHeader,
				"route":  c.Path(),
			})
			return response.ErrUnauthorizedInvalidHeader(c)
		}

		accessToken := accessTokenParts[1]

		claimsAccess, err := jsonWebToken.ValidateJWT(accessToken, sm.JwtKey)
		if err != nil {
			if err == jwt.ErrTokenExpired {
				logger.Warn("Access token expired", map[string]interface{}{
					"token": accessToken,
					"route": c.Path(),
				})
				return response.ErrExpiredAccessToken(c)
			}
			logger.Error("Access token validation error", map[string]interface{}{
				"error": err.Error(),
				"route": c.Path(),
			})
			return response.PersonalizedErr(c, "Token is not valid", fiber.StatusUnauthorized)
		}

		userID, ok := claimsAccess["sub"].(string)
		if !ok {
			logger.Error("UserID claim missing in token", map[string]interface{}{
				"token": accessToken,
				"route": c.Path(),
			})
			return response.PersonalizedErr(c, "Error in token: claims", fiber.StatusForbidden)
		}

		roleID, ok := claimsAccess["rid"].(string)
		if !ok {
			logger.Error("RoleID claim missing in token", map[string]interface{}{
				"token": accessToken,
				"route": c.Path(),
			})
			return response.PersonalizedErr(c, "Error in token: claims", fiber.StatusForbidden)
		}

		c.Locals("userID", userID)
		c.Locals("roleID", roleID)

		logger.Info("Access token verified", map[string]interface{}{
			"userID": userID,
			"roleID": roleID,
			"route":  c.Path(),
		})

		return c.Next()
	}
}

// VerifyRefreshToken checks the presence of a refresh token in the X-Refresh-Token header,
// validates it using JWT, and checks if the token has been blacklisted. If the token is invalid
// or blacklisted, an error response is returned. If valid, the request proceeds.
func (sm *SecurityMiddleware) VerifyRefreshToken() fiber.Handler {
	return func(c fiber.Ctx) error {
		refreshToken := c.Get("X-Refresh-Token")
		if refreshToken == "" {
			logger.Warn("Refresh token header missing", map[string]interface{}{
				"route": c.Path(),
			})
			return response.ErrUnauthorizedHeader(c)
		}

		_, err := jsonWebToken.ValidateJWT(refreshToken, sm.JwtKey)
		if err != nil {
			logger.Error("Refresh token validation error", map[string]interface{}{
				"error": err.Error(),
				"token": refreshToken,
				"route": c.Path(),
			})
			return response.ErrUnauthorized(c)
		}

		if sm.CacheService.IsTokenBlacklisted(refreshToken) {
			logger.Warn("Refresh token has been blacklisted", map[string]interface{}{
				"token": refreshToken,
				"route": c.Path(),
			})
			return response.PersonalizedErr(c, "Refresh Token has been invalidated", fiber.StatusUnauthorized)
		}

		logger.Info("Refresh token verified", map[string]interface{}{
			"token": refreshToken,
			"route": c.Path(),
		})

		return c.Next()
	}
}

// RoleRequiredByName ensures that the user has the required role by name to access the route.
// It retrieves the user's role from the context and compares it with the required role.
// If the role doesn't match or is missing, an error is returned.
func (sm *SecurityMiddleware) RoleRequiredByName(roleRequired string) fiber.Handler {
	return func(c fiber.Ctx) error {
		roleID := c.Locals("roleID")
		if roleID == "" {
			logger.Warn("RoleID missing in context", map[string]interface{}{
				"route": c.Path(),
			})
			return response.PersonalizedErr(c, "Missing information", fiber.StatusForbidden)
		}

		roleIDString, ok := roleID.(string)
		if !ok {
			logger.Error("Invalid roleID format in token", map[string]interface{}{
				"roleID": roleID,
				"route":  c.Path(),
			})
			return response.PersonalizedErr(c, "Error in token: claims", fiber.StatusForbidden)
		}

		roleName, err := sm.AuthService.GetRoleInformationByRoleID(roleIDString)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				logger.Warn("Role not found", map[string]interface{}{
					"roleID": roleIDString,
					"route":  c.Path(),
				})
				return response.ErrNotFound(c)
			}
			logger.Error("Error retrieving role information", map[string]interface{}{
				"roleID": roleIDString,
				"error":  err.Error(),
				"route":  c.Path(),
			})
			return response.ErrInternalServer(c)
		}

		if roleName != roleRequired {
			logger.Warn("Insufficient role permissions", map[string]interface{}{
				"requiredRole": roleRequired,
				"roleName":     roleName,
				"route":        c.Path(),
			})
			return response.ErrForbidden(c)
		}

		return c.Next()
	}
}

// RoleRequiredByID ensures that the user has the required role by ID to access the route.
// It retrieves the user's role from the context and compares it with the required role ID.
// If the role doesn't match or is missing, an error is returned.
func (sm *SecurityMiddleware) RoleRequiredByID(roleRequiredID string) fiber.Handler {
	return func(c fiber.Ctx) error {
		roleID := c.Locals("roleID")
		if roleID == "" {
			logger.Warn("RoleID missing in context", map[string]interface{}{
				"route": c.Path(),
			})
			return response.PersonalizedErr(c, "Missing information", fiber.StatusForbidden)
		}

		roleIDString, ok := roleID.(string)
		if !ok {
			logger.Error("Invalid roleID format in token", map[string]interface{}{
				"roleID": roleID,
				"route":  c.Path(),
			})
			return response.PersonalizedErr(c, "Error in token: claims", fiber.StatusForbidden)
		}

		if roleIDString != roleRequiredID {
			logger.Warn("Insufficient role permissions", map[string]interface{}{
				"requiredRoleID": roleRequiredID,
				"roleID":         roleIDString,
				"route":          c.Path(),
			})
			return response.ErrForbidden(c)
		}

		return c.Next()
	}
}

// AuthorizeSelfUserID checks if the user is authorized to access or modify their own resources.
// It compares the user ID from the token (accessToken) with the ID in the request parameters /:id.
// If the IDs do not match, an unauthorized error is returned.
func (sm *SecurityMiddleware) AuthorizeSelfUserID() fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := c.Locals("userID")
		userIDString, ok := userID.(string)

		if !ok || userIDString == "" {
			logger.Error("Invalid userID format in context", map[string]interface{}{
				"route": c.Path(),
			})
			return response.ErrInternalServer(c)
		}

		if c.Params("id") != userIDString {
			logger.Warn("User ID mismatch", map[string]interface{}{
				"providedID": c.Params("id"),
				"userID":     userIDString,
				"route":      c.Path(),
			})
			return response.ErrUnauthorized(c)
		}

		logger.Info("User authorized", map[string]interface{}{
			"userID": userIDString,
			"route":  c.Path(),
		})

		return c.Next()
	}
}

func (sm *SecurityMiddleware) GetRoleFromToken() fiber.Handler {
	return func(c fiber.Ctx) error {
		accessTokenHeader := c.Get("Authorization")
		if accessTokenHeader == "" {
			logger.Warn("Authorization header missing but is not required, continue", map[string]interface{}{
				"route": c.Path(),
			})
			c.Locals("roleID", "")
			return c.Next()
		}

		accessTokenParts := strings.Split(accessTokenHeader, " ")
		if len(accessTokenParts) != 2 || accessTokenParts[0] != "Bearer" {
			logger.Warn("Invalid authorization header format", map[string]interface{}{
				"header": accessTokenHeader,
				"route":  c.Path(),
			})
			return response.ErrUnauthorizedInvalidHeader(c)
		}

		accessToken := accessTokenParts[1]

		claimsAccess, err := jsonWebToken.ValidateJWT(accessToken, sm.JwtKey)
		if err != nil {
			if err == jwt.ErrTokenExpired {
				logger.Warn("Access token expired", map[string]interface{}{
					"token": accessToken,
					"route": c.Path(),
				})
				return response.ErrExpiredAccessToken(c)
			}
			logger.Error("Access token validation error", map[string]interface{}{
				"error": err.Error(),
				"route": c.Path(),
			})
			return response.PersonalizedErr(c, "Token is not valid", fiber.StatusUnauthorized)
		}

		roleID, ok := claimsAccess["rid"].(string)
		if !ok {
			logger.Error("RoleID claim missing in token", map[string]interface{}{
				"token": accessToken,
				"route": c.Path(),
			})
			return response.PersonalizedErr(c, "Error in token: claims", fiber.StatusForbidden)
		}

		c.Locals("roleID", roleID)

		logger.Info("Role obteneid successfully", map[string]interface{}{
			"roleID": roleID,
			"route":  c.Path(),
		})

		return c.Next()
	}
}
