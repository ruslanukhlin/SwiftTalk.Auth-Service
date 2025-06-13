package bff

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, h *Handler) {	
	app.Post("/register", h.Register)
	app.Post("/login", h.Login)
	app.Post("/refresh", h.RefreshToken)
	app.Post("/verify", h.VerifyToken)
	
	app.Get("/.well-known/jwks.json", h.GetJWKS)
} 