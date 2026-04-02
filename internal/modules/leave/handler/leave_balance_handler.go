package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/thecemakin/hr-project/internal/modules/leave/model"
	"github.com/thecemakin/hr-project/internal/modules/leave/service"
)

// Dummy reference to prevent unused import error, required for swag documentation
var _ = model.LeaveBalance{}

type LeaveBalanceHandler struct {
	service service.LeaveService
}

func NewLeaveBalanceHandler(service service.LeaveService) *LeaveBalanceHandler {
	return &LeaveBalanceHandler{service: service}
}

// ListLeaveBalancesByEmployee handles GET /leave-balances/:employeeId
// @Summary List leave balances for an employee
// @Description Get a list of all leave balances (total and used) for a specific employee
// @Tags Leave Balances
// @Produce json
// @Param employeeId path int true "Employee ID"
// @Success 200 {array} model.LeaveBalance
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/leave/leave-balances/{employeeId} [get]
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

// InitializeBalance handles POST /leave-balances/init
// @Summary Initialize leave balance for an employee
// @Description Set the starting leave balance for an employee and leave type
// @Tags Leave Balances
// @Accept json
// @Produce json
// @Param request body object true "Initialization request"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/leave/leave-balances/init [post]
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
