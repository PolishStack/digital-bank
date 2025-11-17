package http

import "github.com/gofiber/fiber/v2"

func (h *Handler) RegisterRoutes(app *fiber.App) {
	g := app.Group("/api/v1/auth")
	g.Post("/register", h.Register)
	g.Post("/login", h.Login)
}
