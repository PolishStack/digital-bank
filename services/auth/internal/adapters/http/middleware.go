package http

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"bitka/jwt"
	"github.com/rs/zerolog"
)

const CtxUserKey = "user_claims"

func JWTMiddleware(jwtSvc jwt.Service, logger *zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Next() // no auth provided; allow public endpoints (or return 401 if desired)
		}
		if !strings.HasPrefix(auth, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}
		token := strings.TrimPrefix(auth, "Bearer ")

		claims, err := jwtSvc.ValidateAccessToken(context.Background(), token)
		if err != nil {
			logger.Debug().Err(err).Msg("invalid access token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}

		// attach claims to context (fiber Locals)
		c.Locals(CtxUserKey, claims)
		return c.Next()
	}
}

// Helper to get claims inside handlers
func GetClaims(c *fiber.Ctx) (*jwt.Claims, bool) {
	v := c.Locals(CtxUserKey)
	if v == nil {
		return nil, false
	}
	claims, ok := v.(*jwt.Claims)
	return claims, ok
}
