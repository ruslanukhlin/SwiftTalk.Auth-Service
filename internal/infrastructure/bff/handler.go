package bff

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/token"
	"github.com/ruslanukhlin/SwiftTalk.auth-service/pkg/config"
)

var (
	ErrInvalidRefreshToken = errors.New("refresh token не валидный")
	ErrInvalidAccessToken = errors.New("access token не валидный")
)

type Handler struct {
	authService *AuthService
	jwtService token.TokenRepository
	config *config.Config
}

func NewHandler(authService *AuthService, jwtService token.TokenRepository, config *config.Config) *Handler {
	return &Handler{
		authService: authService,
		jwtService: jwtService,
		config: config,
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

	c.Cookie(&fiber.Cookie{
		Name: "access_token",
		Value: tokens.AccessToken,
		Expires: time.Now().Add(h.config.JWT.ExpiresAfter),
	})
	c.Cookie(&fiber.Cookie{
		Name: "refresh_token",
		Value: tokens.RefreshToken,
		Expires: time.Now().Add(h.config.JWT.RefreshExpiresAfter),
	})

	return c.JSON(tokens)
}

func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": ErrInvalidRefreshToken.Error(),
		})
	}

	tokens, err := h.authService.RefreshToken(c.Context(), &RefreshTokenPayload{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return handleGRPCError(c, err)
	}

	c.Cookie(&fiber.Cookie{
		Name: "access_token",
		Value: tokens.AccessToken,
		Expires: time.Now().Add(h.config.JWT.ExpiresAfter),
	})
	c.Cookie(&fiber.Cookie{
		Name: "refresh_token",
		Value: tokens.RefreshToken,
		Expires: time.Now().Add(h.config.JWT.RefreshExpiresAfter),
	})

	return c.JSON(tokens)
}

func (h *Handler) VerifyToken(c *fiber.Ctx) error {
	accessToken := c.Cookies("access_token")
	if accessToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": ErrInvalidAccessToken.Error(),
		})
	}

	response, err := h.authService.VerifyToken(c.Context(), &VerifyTokenPayload{
		AccessToken: accessToken,
	})
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