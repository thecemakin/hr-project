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

func (h *LeaveTypeHandler) ListLeaveTypes(c *fiber.Ctx) error {
	lts, err := h.service.ListLeaveTypes()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(lts)
}
