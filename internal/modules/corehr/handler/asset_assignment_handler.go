package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"github.com/thecemakin/hr-project/internal/modules/corehr/service"
)

// AssetAssignmentHandler handles HTTP requests for asset assignment operations
type AssetAssignmentHandler struct {
	service service.AssetAssignmentService
}

// NewAssetAssignmentHandler creates a new instance of AssetAssignmentHandler
func NewAssetAssignmentHandler(service service.AssetAssignmentService) *AssetAssignmentHandler {
	return &AssetAssignmentHandler{
		service: service,
	}
}

// CreateAssetAssignment handles POST /asset-assignments
func (h *AssetAssignmentHandler) CreateAssetAssignment(c *fiber.Ctx) error {
	var body struct {
		AssetID      uint   `json:"asset_id"`
		EmployeeID   uint   `json:"employee_id"`
		AssignedBy   uint   `json:"assigned_by"`
		Notes        string `json:"notes"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	assignment := &model.AssetAssignment{
		AssetID:    body.AssetID,
		EmployeeID: body.EmployeeID,
		AssignedBy: body.AssignedBy,
		Notes:      body.Notes,
	}

	if err := h.service.CreateAssetAssignment(assignment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Asset assignment created successfully"})
}

// GetAssetAssignmentByID handles GET /asset-assignments/:id
func (h *AssetAssignmentHandler) GetAssetAssignmentByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid assignment ID"})
	}

	assignment, err := h.service.GetAssetAssignmentByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Asset assignment not found"})
	}

	return c.Status(fiber.StatusOK).JSON(assignment)
}

// GetAssetAssignmentsByAssetID handles GET /asset-assignments/asset/:assetId
func (h *AssetAssignmentHandler) GetAssetAssignmentsByAssetID(c *fiber.Ctx) error {
	assetIDParam := c.Params("assetId")
	assetID, err := strconv.ParseUint(assetIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid asset ID"})
	}

	assignments, err := h.service.GetAssetAssignmentsByAssetID(uint(assetID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(assignments)
}

// GetAssetAssignmentsByEmployeeID handles GET /asset-assignments/employee/:employeeId
func (h *AssetAssignmentHandler) GetAssetAssignmentsByEmployeeID(c *fiber.Ctx) error {
	employeeIDParam := c.Params("employeeId")
	employeeID, err := strconv.ParseUint(employeeIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid employee ID"})
	}

	assignments, err := h.service.GetAssetAssignmentsByEmployeeID(uint(employeeID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(assignments)
}

// GetCurrentAssignmentsByEmployeeID handles GET /asset-assignments/current/employee/:employeeId
func (h *AssetAssignmentHandler) GetCurrentAssignmentsByEmployeeID(c *fiber.Ctx) error {
	employeeIDParam := c.Params("employeeId")
	employeeID, err := strconv.ParseUint(employeeIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid employee ID"})
	}

	assignments, err := h.service.GetCurrentAssignmentsByEmployeeID(uint(employeeID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(assignments)
}

// GetActiveAssignmentByAssetID handles GET /asset-assignments/active/asset/:assetId
func (h *AssetAssignmentHandler) GetActiveAssignmentByAssetID(c *fiber.Ctx) error {
	assetIDParam := c.Params("assetId")
	assetID, err := strconv.ParseUint(assetIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid asset ID"})
	}

	assignment, err := h.service.GetActiveAssignmentByAssetID(uint(assetID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No active assignment found for this asset"})
	}

	return c.Status(fiber.StatusOK).JSON(assignment)
}

// GetAllAssetAssignments handles GET /asset-assignments
func (h *AssetAssignmentHandler) GetAllAssetAssignments(c *fiber.Ctx) error {
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

	assignments, err := h.service.GetAllAssetAssignments(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(assignments)
}

// UpdateAssetAssignment handles PATCH /asset-assignments/:id
func (h *AssetAssignmentHandler) UpdateAssetAssignment(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid assignment ID"})
	}

	var body model.AssetAssignment
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	body.ID = uint(id)

	if err := h.service.UpdateAssetAssignment(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Asset assignment updated successfully"})
}

// DeleteAssetAssignment handles DELETE /asset-assignments/:id
func (h *AssetAssignmentHandler) DeleteAssetAssignment(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid assignment ID"})
	}

	if err := h.service.DeleteAssetAssignment(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Asset assignment deleted successfully"})
}

// AssignAsset handles POST /assets/assign
func (h *AssetAssignmentHandler) AssignAsset(c *fiber.Ctx) error {
	var body struct {
		AssetID      uint   `json:"asset_id"`
		EmployeeID   uint   `json:"employee_id"`
		AssignedByID uint   `json:"assigned_by"`
		Notes        string `json:"notes"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.AssignAsset(body.AssetID, body.EmployeeID, body.AssignedByID, body.Notes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Asset assigned successfully"})
}

// ReturnAsset handles POST /assets/:assignmentId/return
func (h *AssetAssignmentHandler) ReturnAsset(c *fiber.Ctx) error {
	assignmentIDParam := c.Params("assignmentId")
	assignmentID, err := strconv.ParseUint(assignmentIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid assignment ID"})
	}

	var body struct {
		ReturnedByID uint   `json:"returned_by"`
		Notes        string `json:"notes"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.ReturnAsset(uint(assignmentID), body.ReturnedByID, body.Notes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Asset returned successfully"})
}