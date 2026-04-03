package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecemakin/hr-project/internal/platform/auth"
	"github.com/thecemakin/hr-project/internal/platform/http/middleware"
)

func SetupRoutesFiber(app *fiber.App, handler *AuthHandler, tp *auth.TokenProvider) {
	v1 := app.Group("/api/v1/auth")

	// Public routes
	v1.Post("/login", handler.Login)

	// Protected routes (Admin only for registration)
	v1.Post("/register",
		middleware.AuthMiddleware(tp),
		middleware.RequireRole("admin"),
		handler.Register,
	)
}
