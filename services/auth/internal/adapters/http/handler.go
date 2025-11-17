package http

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"bitka/auth-service/internal/entities/dto"
	"bitka/auth-service/internal/usecases"
)

type Handler struct {
	authUC usecases.AuthUsecase
}

func NewHandler(a usecases.AuthUsecase) *Handler {
	return &Handler{authUC: a}
}

// See if it need logger & custom error or not

func (h *Handler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}
	// validate (use validator)
	ctx := context.Background()
	user, err := h.authUC.Register(ctx, req.Email, req.Password)
	if err != nil {
		// map error to status code (use internal/errors)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	// TODO: map to response DTO
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": user.ID, "email": user.Email})
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}
	ctx := context.Background()
	access, refresh, err := h.authUC.Login(ctx, req.Email, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}
	return c.JSON(fiber.Map{"access_token": access, "refresh_token": refresh})
}
