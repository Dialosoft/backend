package middleware

import (
	"strings"

	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/services"
	"github.com/Dialosoft/src/pkg/utils/jsonWebToken"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	AuthService services.AuthService
	JwtKey      string
}

func NewAuthMiddleware(authService services.AuthService, jwtKey string) *AuthMiddleware {
	return &AuthMiddleware{AuthService: authService, JwtKey: jwtKey}
}

func (am *AuthMiddleware) IsTokenBlacklisted() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.ErrUnauthorizedHeader(c)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.ErrUnauthorizedInvalidHeader(c)
		}

		token := parts[1]

		if am.AuthService.IsTokenBlacklisted(token) {
			return response.PersonalizedErr(c, "Token has been invalidated", fiber.StatusUnauthorized)
		}

		claims, err := jsonWebToken.ValidateJWT(token, am.JwtKey)
		if err != nil {
			return response.ErrUnauthorized(c)
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			return response.PersonalizedErr(c, "Error in token: claims", fiber.StatusForbidden)
		}

		roleID, ok := claims["roleID"].(string)
		if !ok {
			return response.PersonalizedErr(c, "Error in token: claims", fiber.StatusForbidden)
		}

		c.Locals("userID", userID)
		c.Locals("roleID", roleID)

		return c.Next()
	}
}

func (am *AuthMiddleware) RoleRequiredByName(roleRequired string) fiber.Handler {
	return func(c fiber.Ctx) error {
		roleID := c.Locals("roleID")
		if roleID == "" {
			return response.PersonalizedErr(c, "Missing information", fiber.StatusForbidden)
		}

		roleIDString, ok := roleID.(string)
		if !ok {
			return response.PersonalizedErr(c, "Error in token: claims", fiber.StatusForbidden)
		}

		roleName, err := am.AuthService.GetRoleInformationByRoleID(roleIDString)
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

func (am *AuthMiddleware) RoleRequiredByID(roleRequiredID string) fiber.Handler {
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
