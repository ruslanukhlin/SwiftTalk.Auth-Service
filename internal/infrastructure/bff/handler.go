package bff

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	authService *AuthService
}

func NewHandler(authService *AuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}

func (h *Handler) Register(c *fiber.Ctx) error {
	payload := new(RegisterPayload)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат данных",
		})
	}

	tokens, err := h.authService.Register(c.Context(), payload)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.JSON(tokens)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	payload := new(LoginPayload)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат данных",
		})
	}

	tokens, err := h.authService.Login(c.Context(), payload)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.JSON(tokens)
}

func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	payload := new(RefreshTokenPayload)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат данных",
		})
	}

	tokens, err := h.authService.RefreshToken(c.Context(), payload)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.JSON(tokens)
}

func (h *Handler) VerifyToken(c *fiber.Ctx) error {
	payload := new(VerifyTokenPayload)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат данных",
		})
	}

	response, err := h.authService.VerifyToken(c.Context(), payload)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.JSON(response)
} 