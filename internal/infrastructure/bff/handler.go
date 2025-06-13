package bff

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/token"
)

type Handler struct {
	authService *AuthService
	jwtService token.TokenRepository
}

func NewHandler(authService *AuthService, jwtService token.TokenRepository) *Handler {
	return &Handler{
		authService: authService,
		jwtService: jwtService,
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

func (h *Handler) GetJWKS(c *fiber.Ctx) error {
    jwks, err := h.jwtService.GetJWKS()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to get JWKS",
        })
    }

    c.Set("Content-Type", "application/json")
    return c.Send(jwks)
} 