package middleware

import (
	"strings"

	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/services"
	"github.com/Dialosoft/src/pkg/utils/jsonWebToken"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
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
		refreshToken := c.Get("X-Refresh-Token")
		if refreshToken == "" {
			return response.ErrUnauthorizedHeader(c)
		}

		accesTokenHeader := c.Get("Authorization")
		if accesTokenHeader == "" {
			return response.ErrUnauthorizedHeader(c)
		}

		accessTokenParts := strings.Split(accesTokenHeader, " ")
		if len(accessTokenParts) != 2 || accessTokenParts[0] != "Bearer" {
			return response.ErrUnauthorizedInvalidHeader(c)
		}

		accesToken := accessTokenParts[1]

		_, err := jsonWebToken.ValidateJWT(refreshToken, am.JwtKey)
		if err != nil {
			return response.ErrUnauthorized(c)
		}

		if am.AuthService.IsTokenBlacklisted(refreshToken) {
			return response.PersonalizedErr(c, "Token has been invalidated", fiber.StatusUnauthorized)
		}

		claimsAcess, err := jsonWebToken.ValidateJWT(accesToken, am.JwtKey)
		if err != nil {
			if err == jwt.ErrTokenExpired {
				return response.ErrExpiredAccessToken(c)
			} else {
				return response.PersonalizedErr(c, "token is not valid", fiber.StatusUnauthorized)
			}
		}

		userID, ok := claimsAcess["sub"].(string)
		if !ok {
			return response.PersonalizedErr(c, "Error in token: claims", fiber.StatusForbidden)
		}

		roleID, ok := claimsAcess["rid"].(string)
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

func (am *AuthMiddleware) AuthorizeSelfUserID() fiber.Handler {
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
