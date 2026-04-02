package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"github.com/thecemakin/hr-project/internal/modules/corehr/service"
)

// DepartmentHandler handles HTTP requests for department operations
type DepartmentHandler struct {
	service service.DepartmentService
}

// NewDepartmentHandler creates a new instance of DepartmentHandler
func NewDepartmentHandler(service service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{
		service: service,
	}
}

// CreateDepartment handles POST /departments
func (h *DepartmentHandler) CreateDepartment(c *fiber.Ctx) error {
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Code        string `json:"code"`
		HeadID      *uint  `json:"head_id"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	department := &model.Department{
		Name:        body.Name,
		Description: body.Description,
		Code:        body.Code,
		HeadID:      body.HeadID,
	}

	if err := h.service.CreateDepartment(department); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Department created successfully"})
}

// GetDepartmentByID handles GET /departments/:id
func (h *DepartmentHandler) GetDepartmentByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid department ID"})
	}

	department, err := h.service.GetDepartmentByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Department not found"})
	}

	return c.Status(fiber.StatusOK).JSON(department)
}

// GetDepartmentByName handles GET /departments/name/:name
func (h *DepartmentHandler) GetDepartmentByName(c *fiber.Ctx) error {
	name := c.Params("name")
	department, err := h.service.GetDepartmentByName(name)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Department not found"})
	}

	return c.Status(fiber.StatusOK).JSON(department)
}

// GetAllDepartments handles GET /departments
func (h *DepartmentHandler) GetAllDepartments(c *fiber.Ctx) error {
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

	departments, err := h.service.GetAllDepartments(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(departments)
}

// UpdateDepartment handles PATCH /departments/:id
func (h *DepartmentHandler) UpdateDepartment(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid department ID"})
	}

	var body model.Department
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	body.ID = uint(id)

	if err := h.service.UpdateDepartment(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Department updated successfully"})
}

// DeleteDepartment handles DELETE /departments/:id
func (h *DepartmentHandler) DeleteDepartment(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid department ID"})
	}

	if err := h.service.DeleteDepartment(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Department deleted successfully"})
}