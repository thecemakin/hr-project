package http

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/yokeTH/gofiber-scalar/scalar/v2"
	"github.com/thecemakin/hr-project/docs/openapi"
)

type Server struct {
	App *fiber.App
}

func NewServer() *Server {
	app := fiber.New()

	// A good base middleware stack
	app.Use(recover.New())
	app.Use(logger.New())

	// Base API route
	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Swagger UI
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Scalar UI
	app.Get("/docs*", scalar.New(scalar.Config{
		FileContentString: openapi.SwaggerInfo.ReadDoc(),
		Title:             "HR Project API Documentation",
	}))

	return &Server{
		App: app,
	}
}

func (s *Server) Mount(path string, handler fiber.Router) {
	// Fiber handles mounting differently, this is a simplified approach
	// In Fiber, routes are typically defined directly on the app instance
	log.Printf("Mounting routes at %s - Note: Fiber uses different mounting mechanism", path)
}
