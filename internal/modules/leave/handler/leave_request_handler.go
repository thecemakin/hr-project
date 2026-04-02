package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/thecemakin/hr-project/internal/modules/leave/model"
	"github.com/thecemakin/hr-project/internal/modules/leave/service"
)

type LeaveRequestHandler struct {
	service service.LeaveService
}

func NewLeaveRequestHandler(service service.LeaveService) *LeaveRequestHandler {
	return &LeaveRequestHandler{service: service}
}

func (h *LeaveRequestHandler) SubmitRequest(c *fiber.Ctx) error {
	var req model.LeaveRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.SubmitLeaveRequest(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(req)
}

func (h *LeaveRequestHandler) ListOwnRequests(c *fiber.Ctx) error {
	// In a real app, empID would come from JWT / auth context.
	// For testing, we might pass it via a query param or header until Auth is built.
	tempIDStr := c.Query("employeeId")
	tempID, err := strconv.ParseUint(tempIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing or invalid employeeId query param"})
	}

	reqs, err := h.service.ListOwnLeaveRequests(uint(tempID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(reqs)
}

func (h *LeaveRequestHandler) ListPendingApprovals(c *fiber.Ctx) error {
	// From JWT context representing the logged-in manager
	managerIDStr := c.Query("managerId")
	managerID, err := strconv.ParseUint(managerIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing or invalid managerId query param"})
	}

	reqs, err := h.service.ListPendingApprovals(uint(managerID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(reqs)
}

func (h *LeaveRequestHandler) ApproveRequest(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	// From JWT context
	managerIDStr := c.Query("managerId")
	managerID, err := strconv.ParseUint(managerIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing or invalid managerId query param"})
	}

	var payload struct {
		Note string `json:"note"`
	}
	// Note is optional
	_ = c.BodyParser(&payload)

	if err := h.service.ApproveLeaveRequest(uint(id), uint(managerID), payload.Note); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "approved"})
}

func (h *LeaveRequestHandler) RejectRequest(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	// From JWT context
	managerIDStr := c.Query("managerId")
	managerID, err := strconv.ParseUint(managerIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing or invalid managerId query param"})
	}

	var payload struct {
		Note string `json:"note"`
	}
	// Note is optional
	_ = c.BodyParser(&payload)

	if err := h.service.RejectLeaveRequest(uint(id), uint(managerID), payload.Note); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "rejected"})
}
