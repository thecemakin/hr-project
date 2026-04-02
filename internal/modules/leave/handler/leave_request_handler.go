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

// SubmitRequest handles POST /leave-requests
// @Summary Submit a new leave request
// @Description Create a new leave request for approval
// @Tags Leave Requests
// @Accept json
// @Produce json
// @Param request body model.LeaveRequest true "Leave Request object"
// @Success 201 {object} model.LeaveRequest
// @Failure 400 {object} map[string]string
// @Router /api/v1/leave/leave-requests [post]
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

// ListOwnRequests handles GET /leave-requests/me
// @Summary List own leave requests
// @Description Get a list of leave requests submitted by the current employee
// @Tags Leave Requests
// @Produce json
// @Param employeeId query int true "Employee ID (simulation of auth context)"
// @Success 200 {array} model.LeaveRequest
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/leave/leave-requests/me [get]
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

// ListPendingApprovals handles GET /leave-requests/pending-approvals
// @Summary List pending leave requests for approval
// @Description Get a list of leave requests pending approval for a specific manager
// @Tags Leave Requests
// @Produce json
// @Param managerId query int true "Manager ID (simulation of auth context)"
// @Success 200 {array} model.LeaveRequest
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/leave/leave-requests/pending-approvals [get]
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

// ApproveRequest handles POST /leave-requests/:id/approve
// @Summary Approve a leave request
// @Description Approve a pending leave request
// @Tags Leave Requests
// @Accept json
// @Produce json
// @Param id path int true "Leave Request ID"
// @Param managerId query int true "Manager ID (simulation of auth context)"
// @Param body body object false "Approval note"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/leave/leave-requests/{id}/approve [post]
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

// RejectRequest handles POST /leave-requests/:id/reject
// @Summary Reject a leave request
// @Description Reject a pending leave request
// @Tags Leave Requests
// @Accept json
// @Produce json
// @Param id path int true "Leave Request ID"
// @Param managerId query int true "Manager ID (simulation of auth context)"
// @Param body body object false "Rejection note"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/leave/leave-requests/{id}/reject [post]
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
