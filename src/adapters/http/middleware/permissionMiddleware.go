package middleware

import (
	"github.com/Dialosoft/src/domain/services"
	"github.com/gofiber/fiber/v3"
)

type PermissionMiddleware struct {
	AuthService  services.AuthService
	CacheService services.CacheService
	ForumService services.ForumService
	JwtKey       string
}

func NewPermissionMiddleware(authService services.AuthService, cacheService services.CacheService, jwtKey string) *PermissionMiddleware {
	return &PermissionMiddleware{AuthService: authService, CacheService: cacheService, JwtKey: jwtKey}
}

func (sm *PermissionMiddleware) GetForumProtection() fiber.Handler {
	return func(c fiber.Ctx) error {
		panic("unimplemented")
	}
}
