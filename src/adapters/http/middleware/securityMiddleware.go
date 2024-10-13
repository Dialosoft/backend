package middleware

import (
	"strings"

	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/services"
	"github.com/Dialosoft/src/pkg/errorsUtils"
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
	Layer        string
}

func NewSecurityMiddleware(authService services.AuthService, cacheService services.CacheService, jwtKey string, Layer string) *SecurityMiddleware {
	return &SecurityMiddleware{AuthService: authService, CacheService: cacheService, JwtKey: jwtKey, Layer: Layer}
}

// GetAndVerifyAccessToken retrieves the access token from the Authorization header,
// checks its format, and verifies the token. If the token is invalid, expired, or missing,
// appropriate error responses are returned. If valid, the user's ID and role are extracted
// from the token and stored in the request context for further use.
func (sm *SecurityMiddleware) GetAndVerifyAccessToken() fiber.Handler {
	return func(c fiber.Ctx) error {
		accessTokenHeader := c.Get("Authorization")
		if accessTokenHeader == "" {
			return response.ErrUnauthorizedHeader(c, sm.Layer)
		}

		accessTokenParts := strings.Split(accessTokenHeader, " ")
		if len(accessTokenParts) != 2 || accessTokenParts[0] != "Bearer" {
			return response.ErrUnauthorizedInvalidHeader(c, sm.Layer)
		}

		accessToken := accessTokenParts[1]

		claimsAccess, err := jsonWebToken.ValidateJWT(accessToken, sm.JwtKey)
		if err != nil {
			if err == jwt.ErrTokenExpired {
				return response.ErrExpiredAccessToken(c, sm.Layer)
			}
			return response.ErrInvalidToken(c, sm.Layer)
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
			return response.ErrUnauthorizedHeader(c, sm.Layer)
		}

		_, err := jsonWebToken.ValidateJWT(refreshToken, sm.JwtKey)
		if err != nil {
			return response.ErrUnauthorized(c, refreshToken, err, sm.Layer)
		}

		if sm.CacheService.IsTokenBlacklisted(refreshToken) {
			return response.ErrTokenIsBlacklisted(c, sm.Layer)
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
				return response.ErrNotFound(c, sm.Layer)
			}
			return response.ErrInternalServer(c, err, nil, sm.Layer)
		}

		if roleName != roleRequired {
			return response.ErrForbidden(c, sm.Layer)
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
			return response.ErrForbidden(c, sm.Layer)
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
			return response.ErrInternalServer(c, errorsUtils.ErrInternalServer, nil, sm.Layer)
		}

		if c.Params("id") != userIDString {
			return response.ErrUnauthorized(c, userIDString, errorsUtils.ErrUnauthorizedAcces, sm.Layer)
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
			return response.ErrUnauthorizedInvalidHeader(c, sm.Layer)
		}

		accessToken := accessTokenParts[1]

		claimsAccess, err := jsonWebToken.ValidateJWT(accessToken, sm.JwtKey)
		if err != nil {
			if err == jwt.ErrTokenExpired {
				return response.ErrExpiredAccessToken(c, sm.Layer)
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
