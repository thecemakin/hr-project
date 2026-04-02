package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"github.com/thecemakin/hr-project/internal/modules/corehr/service"
)

// PositionHandler handles HTTP requests for position operations
type PositionHandler struct {
	service service.PositionService
}

// NewPositionHandler creates a new instance of PositionHandler
func NewPositionHandler(service service.PositionService) *PositionHandler {
	return &PositionHandler{
		service: service,
	}
}

// CreatePosition handles POST /positions
func (h *PositionHandler) CreatePosition(c *fiber.Ctx) error {
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Code        string `json:"code"`
		Level       int    `json:"level"`
		SalaryMin   int64  `json:"salary_min"`
		SalaryMax   int64  `json:"salary_max"`
		IsActive    bool   `json:"is_active"`
		DepartmentID *uint `json:"department_id"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	position := &model.Position{
		Title:       body.Title,
		Description: body.Description,
		Code:        body.Code,
		Level:       body.Level,
		SalaryMin:   body.SalaryMin,
		SalaryMax:   body.SalaryMax,
		IsActive:    body.IsActive,
		DepartmentID: body.DepartmentID,
	}

	if err := h.service.CreatePosition(position); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Position created successfully"})
}

// GetPositionByID handles GET /positions/:id
func (h *PositionHandler) GetPositionByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid position ID"})
	}

	position, err := h.service.GetPositionByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Position not found"})
	}

	return c.Status(fiber.StatusOK).JSON(position)
}

// GetPositionByTitle handles GET /positions/title/:title
func (h *PositionHandler) GetPositionByTitle(c *fiber.Ctx) error {
	title := c.Params("title")
	position, err := h.service.GetPositionByTitle(title)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Position not found"})
	}

	return c.Status(fiber.StatusOK).JSON(position)
}

// GetAllPositions handles GET /positions
func (h *PositionHandler) GetAllPositions(c *fiber.Ctx) error {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	positions, err := h.service.GetAllPositions(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(positions)
}

// UpdatePosition handles PATCH /positions/:id
func (h *PositionHandler) UpdatePosition(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid position ID"})
	}

	var body model.Position
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	body.ID = uint(id)

	if err := h.service.UpdatePosition(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Position updated successfully"})
}

// DeletePosition handles DELETE /positions/:id
func (h *PositionHandler) DeletePosition(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid position ID"})
	}

	if err := h.service.DeletePosition(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Position deleted successfully"})
}