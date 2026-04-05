package middleware

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/thecemakin/hr-project/internal/platform/auth"
)

// AuthMiddleware returns a middleware that validates JWT tokens
func AuthMiddleware(tp *auth.TokenProvider, appEnv string, skipAuth bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// Skip authentication in development if configured
			if appEnv == "development" && skipAuth {
				log.Println("[DEV MODE] Bypassing authentication and injecting Admin identity")
				c.Locals("user_id", uint(1))
				c.Locals("email", "admin@hr-project.com")
				c.Locals("role", "admin")
				return c.Next()
			}

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format",
			})
		}

		tokenString := parts[1]
		claims, err := tp.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: " + err.Error(),
			})
		}

		// Store claims in context
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		if claims.EmployeeID != nil {
			c.Locals("employee_id", *claims.EmployeeID)
		}

		return c.Next()
	}
}

// RequireRole returns a middleware that checks if the user has one of the required roles
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")
		if userRole == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Role not found in context",
			})
		}

		roleStr, ok := userRole.(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Invalid role type",
			})
		}

		for _, r := range roles {
			if strings.EqualFold(roleStr, r) {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden: insufficient permissions",
		})
	}
}

// GetUserID helper to extract user id from context
func GetUserID(c *fiber.Ctx) uint {
	id := c.Locals("user_id")
	if id == nil {
		log.Println("GetUserID called on unprotected route")
		return 0
	}
	return id.(uint)
}

// GetUserRole helper to extract role from context
func GetUserRole(c *fiber.Ctx) string {
	role := c.Locals("role")
	if role == nil {
		return ""
	}
	return role.(string)
}
