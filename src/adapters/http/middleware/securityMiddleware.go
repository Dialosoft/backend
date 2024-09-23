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

func (sm *SecurityMiddleware) GetAndVerifyAccesToken() fiber.Handler {
	return func(c fiber.Ctx) error {
		accessTokenHeader := c.Get("Authorization")
		if accessTokenHeader == "" {
			logger.Error("refresh token empty")
			return response.ErrUnauthorizedHeader(c)
		}

		accessTokenParts := strings.Split(accessTokenHeader, " ")
		if len(accessTokenParts) != 2 || accessTokenParts[0] != "Bearer" {
			logger.Error("bad access token")
			return response.ErrUnauthorizedInvalidHeader(c)
		}

		accesToken := accessTokenParts[1]

		claimsAcess, err := jsonWebToken.ValidateJWT(accesToken, sm.JwtKey)
		if err != nil {
			if err == jwt.ErrTokenExpired {
				logger.Error(err.Error())
				return response.ErrExpiredAccessToken(c)
			}
			logger.Error(err.Error())
			return response.PersonalizedErr(c, "token is not valid", fiber.StatusUnauthorized)
		}

		userID, ok := claimsAcess["sub"].(string)
		if !ok {
			logger.Error("Error in token: claims")
			return response.PersonalizedErr(c, "Error in token: claims", fiber.StatusForbidden)
		}

		roleID, ok := claimsAcess["rid"].(string)
		if !ok {
			logger.Error("Error in token: claims")
			return response.PersonalizedErr(c, "Error in token: claims", fiber.StatusForbidden)
		}

		c.Locals("userID", userID)
		c.Locals("roleID", roleID)

		return c.Next()
	}
}

func (sm *SecurityMiddleware) VerifyRefreshToken() fiber.Handler {
	return func(c fiber.Ctx) error {
		refreshToken := c.Get("X-Refresh-Token")
		if refreshToken == "" {
			logger.Info("refresh token empty")
			return response.ErrUnauthorizedHeader(c)
		}

		_, err := jsonWebToken.ValidateJWT(refreshToken, sm.JwtKey)
		if err != nil {
			logger.Error(err.Error())
			return response.ErrUnauthorized(c)
		}

		if sm.CacheService.IsTokenBlacklisted(refreshToken) {
			return response.PersonalizedErr(c, "Refresh Token has been invalidated", fiber.StatusUnauthorized)
		}

		return c.Next()
	}
}

func (sm *SecurityMiddleware) RoleRequiredByName(roleRequired string) fiber.Handler {
	return func(c fiber.Ctx) error {
		roleID := c.Locals("roleID")
		if roleID == "" {
			return response.PersonalizedErr(c, "Missing information", fiber.StatusForbidden)
		}

		roleIDString, ok := roleID.(string)
		if !ok {
			return response.PersonalizedErr(c, "Error in token: claims", fiber.StatusForbidden)
		}

		roleName, err := sm.AuthService.GetRoleInformationByRoleID(roleIDString)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return response.ErrNotFound(c)
			}
		}

		if roleName != roleRequired {
			return response.ErrForbidden(c)
		}

		return c.Next()
	}
}

func (sm *SecurityMiddleware) RoleRequiredByID(roleRequiredID string) fiber.Handler {
	return func(c fiber.Ctx) error {
		roleID := c.Locals("roleID")
		if roleID == "" {
			return response.PersonalizedErr(c, "Missing information", fiber.StatusForbidden)
		}

		roleIDString, ok := roleID.(string)
		if !ok {
			return response.PersonalizedErr(c, "Error in token: claims", fiber.StatusForbidden)
		}

		if roleIDString != roleRequiredID {
			return response.ErrForbidden(c)
		}

		return c.Next()
	}
}

func (sm *SecurityMiddleware) AuthorizeSelfUserID() fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := c.Locals("userID")
		userIDString, ok := userID.(string)

		if !ok || userIDString == "" {
			return response.ErrInternalServer(c)
		}

		if c.Params("id") != userIDString {
			return response.ErrUnauthorized(c)
		}

		return c.Next()
	}
}
