package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/thecemakin/hr-project/internal/modules/leave/service"
)

type LeaveBalanceHandler struct {
	service service.LeaveService
}

func NewLeaveBalanceHandler(service service.LeaveService) *LeaveBalanceHandler {
	return &LeaveBalanceHandler{service: service}
}

func (h *LeaveBalanceHandler) ListLeaveBalancesByEmployee(c *fiber.Ctx) error {
	tempIDStr := c.Params("employeeId")
	tempID, err := strconv.ParseUint(tempIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid employee id"})
	}

	lbs, err := h.service.ListLeaveBalancesByEmployee(uint(tempID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(lbs)
}

func (h *LeaveBalanceHandler) InitializeBalance(c *fiber.Ctx) error {
	var req struct {
		EmployeeID   uint `json:"employee_id"`
		LeaveTypeID  uint `json:"leave_type_id"`
		StartingDays int  `json:"starting_days"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.InitializeLeaveBalance(req.EmployeeID, req.LeaveTypeID, req.StartingDays); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success"})
}
