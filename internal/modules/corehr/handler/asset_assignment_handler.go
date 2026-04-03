package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"github.com/thecemakin/hr-project/internal/modules/corehr/repository"
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
// @Security ApiKeyAuth
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
// @Security ApiKeyAuth
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

func (h *AssetAssignmentHandler) GetAssetAssignmentsByAssetID(c *fiber.Ctx) error {
	assetID, _ := strconv.ParseUint(c.Params("assetId"), 10, 32)
	filter := repository.AssetAssignmentFilter{AssetID: uint(assetID)}
	assignments, err := h.service.GetAllAssetAssignments(100, 0, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(assignments)
}

func (h *AssetAssignmentHandler) GetAssetAssignmentsByEmployeeID(c *fiber.Ctx) error {
	employeeID, _ := strconv.ParseUint(c.Params("employeeId"), 10, 32)
	filter := repository.AssetAssignmentFilter{EmployeeID: uint(employeeID)}
	assignments, err := h.service.GetAllAssetAssignments(100, 0, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(assignments)
}

func (h *AssetAssignmentHandler) GetCurrentAssignmentsByEmployeeID(c *fiber.Ctx) error {
	employeeID, _ := strconv.ParseUint(c.Params("employeeId"), 10, 32)
	filter := repository.AssetAssignmentFilter{EmployeeID: uint(employeeID), Status: "assigned"}
	assignments, err := h.service.GetAllAssetAssignments(100, 0, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(assignments)
}

func (h *AssetAssignmentHandler) GetActiveAssignmentByAssetID(c *fiber.Ctx) error {
	assetID, _ := strconv.ParseUint(c.Params("assetId"), 10, 32)
	assignment, err := h.service.GetActiveAssignmentByAssetID(uint(assetID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(assignment)
}

// GetAllAssetAssignments handles GET /asset-assignments
// @Summary Get asset assignments with filtering
// @Description Get a list of all asset assignments with pagination and optional filters (asset_id, employee_id, status)
// @Tags Asset Assignments
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Param asset_id query int false "Filter by asset ID"
// @Param employee_id query int false "Filter by employee ID"
// @Param status query string false "Filter by status"
// @Security ApiKeyAuth
// @Success 200 {array} model.AssetAssignment
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/asset-assignments [get]
func (h *AssetAssignmentHandler) GetAllAssetAssignments(c *fiber.Ctx) error {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")
	assetIDStr := c.Query("asset_id")
	employeeIDStr := c.Query("employee_id")
	status := c.Query("status")

	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset, _ := strconv.Atoi(offsetStr)
	if offset < 0 {
		offset = 0
	}

	assetID, _ := strconv.ParseUint(assetIDStr, 10, 32)
	employeeID, _ := strconv.ParseUint(employeeIDStr, 10, 32)

	filter := repository.AssetAssignmentFilter{
		AssetID:    uint(assetID),
		EmployeeID: uint(employeeID),
		Status:     status,
	}

	assignments, err := h.service.GetAllAssetAssignments(limit, offset, filter)
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
// @Security ApiKeyAuth
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
// @Security ApiKeyAuth
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
// @Security ApiKeyAuth
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
// @Security ApiKeyAuth
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