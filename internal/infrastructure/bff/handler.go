package bff

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ruslanukhlin/SwiftTalk.Auth-service/internal/domain/token"
	"github.com/ruslanukhlin/SwiftTalk.Auth-service/pkg/config"
)

var (
	ErrInvalidRefreshToken = errors.New("refresh token не валидный")
	ErrInvalidAccessToken  = errors.New("access token не валидный")
)

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

type Handler struct {
	authService *AuthService
	jwtService  token.TokenRepository
	config      *config.Config
}

func NewHandler(authService *AuthService, jwtService token.TokenRepository) *Handler {
	cfg := config.LoadConfigFromEnv()
	return &Handler{
		authService: authService,
		jwtService:  jwtService,
		config:      cfg,
	}
}

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя и возвращает токены доступа
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body RegisterPayload true "Данные для регистрации"
// @Success 200 {object} TokenResponse "Успешная регистрация"
// @Failure 400 {object} ErrorResponse "Ошибка в параметрах запроса"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /register [post]
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

// Login godoc
// @Summary Вход в систему
// @Description Аутентифицирует пользователя и возвращает токены доступа
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body LoginPayload true "Данные для входа"
// @Success 200 {object} TokenResponse "Успешный вход"
// @Failure 400 {object} ErrorResponse "Ошибка в параметрах запроса"
// @Failure 401 {string} string Jwt is expired
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /login [post]
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
		Name:    "access_token",
		Value:   tokens.AccessToken,
		Expires: time.Now().Add(h.config.JWT.ExpiresAfter),
	})
	c.Cookie(&fiber.Cookie{
		Name:    "refresh_token",
		Value:   tokens.RefreshToken,
		Expires: time.Now().Add(h.config.JWT.RefreshExpiresAfter),
	})

	return c.JSON(tokens)
}

// RefreshToken godoc
// @Summary Обновление токена
// @Description Обновляет access и refresh токены
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} TokenResponse "Успешное обновление токенов"
// @Failure 401 {object} ErrorResponse "Невалидный refresh token"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /refresh [post]
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
		Name:    "access_token",
		Value:   tokens.AccessToken,
		Expires: time.Now().Add(h.config.JWT.ExpiresAfter),
	})
	c.Cookie(&fiber.Cookie{
		Name:    "refresh_token",
		Value:   tokens.RefreshToken,
		Expires: time.Now().Add(h.config.JWT.RefreshExpiresAfter),
	})

	return c.JSON(tokens)
}

// VerifyToken godoc
// @Summary Проверка токена
// @Description Проверяет валидность access token и возвращает информацию о пользователе
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} VerifyTokenResponse "Успешная проверка токена"
// @Failure 401 {object} ErrorResponse "Невалидный access token"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /verify [get]
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
