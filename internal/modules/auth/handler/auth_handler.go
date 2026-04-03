package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecemakin/hr-project/internal/modules/auth/service"
)

type AuthHandler struct {
	svc service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	EmployeeID *uint  `json:"employee_id"`
}

// Login handles user authentication
// @Summary Login user
// @Description Authenticate user and return JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body loginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Token and user data"
// @Failure 401 {object} map[string]string "Error message"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	token, user, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":          user.ID,
			"email":       user.Email,
			"role":        user.Role,
			"employee_id": user.EmployeeID,
		},
	})
}

// Register handles user creation (Admin only)
// @Summary Create user
// @Description Create a new login account (Admin only)
// @Tags Auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param register body registerRequest true "User data"
// @Success 201 {object} map[string]interface{} "Created user"
// @Failure 403 {object} map[string]string "Forbidden"
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.svc.Register(req.Email, req.Password, req.Role, req.EmployeeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"user": fiber.Map{
			"id":          user.ID,
			"email":       user.Email,
			"role":        user.Role,
			"employee_id": user.EmployeeID,
		},
	})
}
