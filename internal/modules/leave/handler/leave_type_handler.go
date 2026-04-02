package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/thecemakin/hr-project/internal/modules/leave/model"
	"github.com/thecemakin/hr-project/internal/modules/leave/service"
)

type LeaveTypeHandler struct {
	service service.LeaveService
}

func NewLeaveTypeHandler(service service.LeaveService) *LeaveTypeHandler {
	return &LeaveTypeHandler{service: service}
}

// CreateLeaveType handles POST /leave-types
// @Summary Create a new leave type
// @Description Create a new category for leave requests
// @Tags Leave Types
// @Accept json
// @Produce json
// @Param leaveType body model.LeaveType true "Leave Type object"
// @Success 201 {object} model.LeaveType
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/leave/leave-types [post]
func (h *LeaveTypeHandler) CreateLeaveType(c *fiber.Ctx) error {
	var lt model.LeaveType
	if err := c.BodyParser(&lt); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.CreateLeaveType(&lt); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(lt)
}

// GetLeaveType handles GET /leave-types/:id
// @Summary Get leave type by ID
// @Description Get detailed information about a leave type
// @Tags Leave Types
// @Produce json
// @Param id path int true "Leave Type ID"
// @Success 200 {object} model.LeaveType
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/leave/leave-types/{id} [get]
func (h *LeaveTypeHandler) GetLeaveType(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	lt, err := h.service.GetLeaveTypeByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "leave type not found"})
	}

	return c.JSON(lt)
}

// ListLeaveTypes handles GET /leave-types
// @Summary List all leave types
// @Description Get a list of all available leave types
// @Tags Leave Types
// @Produce json
// @Success 200 {array} model.LeaveType
// @Failure 500 {object} map[string]string
// @Router /api/v1/leave/leave-types [get]
func (h *LeaveTypeHandler) ListLeaveTypes(c *fiber.Ctx) error {
	lts, err := h.service.ListLeaveTypes()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(lts)
}
