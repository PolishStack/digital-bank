package jwt

import (
    "context"
    "strings"

    "github.com/gofiber/fiber/v2"
)

// headerName: "Authorization", ctxKey: where to store claims in Locals
func FiberAuthMiddleware(v Verifier, headerName string, ctxKey string) fiber.Handler {
    if headerName == "" {
        headerName = "Authorization"
    }
    if ctxKey == "" {
        ctxKey = "claims"
    }
    return func(c *fiber.Ctx) error {
        h := c.Get(headerName)
        if h == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing_authorization"})
        }
        parts := strings.SplitN(h, " ", 2)
        if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid_authorization"})
        }
        tokenStr := parts[1]
        claims, err := v.Verify(context.Background(), tokenStr)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
        }
        c.Locals(ctxKey, claims)
        return c.Next()
    }
}
