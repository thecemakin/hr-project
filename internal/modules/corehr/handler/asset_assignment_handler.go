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
// @Summary Create a new asset assignment
// @Description Create a new record for an asset assigned to an employee
// @Tags Asset Assignments
// @Accept json
// @Produce json
// @Param assignment body model.AssetAssignment true "Asset Assignment object"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/asset-assignments [post]
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
// @Summary Get asset assignment by ID
// @Description Get detailed information about an asset assignment
// @Tags Asset Assignments
// @Produce json
// @Param id path int true "Assignment ID"
// @Success 200 {object} model.AssetAssignment
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/corehr/asset-assignments/{id} [get]
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
// @Summary Get asset assignments by asset ID
// @Description Get all assignment history for a specific asset
// @Tags Asset Assignments
// @Produce json
// @Param assetId path int true "Asset ID"
// @Success 200 {array} model.AssetAssignment
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/asset-assignments/asset/{assetId} [get]
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
// @Summary Get asset assignments by employee ID
// @Description Get all assignment history for a specific employee
// @Tags Asset Assignments
// @Produce json
// @Param employeeId path int true "Employee ID"
// @Success 200 {array} model.AssetAssignment
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/asset-assignments/employee/{employeeId} [get]
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
// @Summary Get current active assignments by employee ID
// @Description Get only the currently active asset assignments for a specific employee
// @Tags Asset Assignments
// @Produce json
// @Param employeeId path int true "Employee ID"
// @Success 200 {array} model.AssetAssignment
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/asset-assignments/current/employee/{employeeId} [get]
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
// @Summary Get active assignment by asset ID
// @Description Get the current active assignment for a specific asset
// @Tags Asset Assignments
// @Produce json
// @Param assetId path int true "Asset ID"
// @Success 200 {object} model.AssetAssignment
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/corehr/asset-assignments/active/asset/{assetId} [get]
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
// @Summary Get all asset assignments
// @Description Get a list of all asset assignments with pagination
// @Tags Asset Assignments
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} model.AssetAssignment
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/asset-assignments [get]
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

// UpdateAssetAssignment handles PUT /asset-assignments/:id
// @Summary Update an asset assignment
// @Description Update an existing asset assignment's details
// @Tags Asset Assignments
// @Accept json
// @Produce json
// @Param id path int true "Assignment ID"
// @Param assignment body model.AssetAssignment true "Asset Assignment update object"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/asset-assignments/{id} [put]
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
// @Summary Delete an asset assignment
// @Description Delete an asset assignment by its ID
// @Tags Asset Assignments
// @Produce json
// @Param id path int true "Assignment ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/asset-assignments/{id} [delete]
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
// @Summary Assign an asset to an employee (convenience)
// @Description Directly assign an asset to an employee
// @Tags Assets
// @Accept json
// @Produce json
// @Param assignment body object true "Assignment request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/assets/assign [post]
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
// @Summary Return an assigned asset
// @Description Mark an assigned asset as returned
// @Tags Assets
// @Accept json
// @Produce json
// @Param assignmentId path int true "Assignment ID"
// @Param return body object true "Return request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/assets/{assignmentId}/return [post]
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