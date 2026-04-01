package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
func (h *AssetAssignmentHandler) CreateAssetAssignment(w http.ResponseWriter, r *http.Request) {
	var body struct {
		AssetID      uint   `json:"asset_id"`
		EmployeeID   uint   `json:"employee_id"`
		AssignedBy   uint   `json:"assigned_by"`
		Notes        string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	assignment := &model.AssetAssignment{
		AssetID:    body.AssetID,
		EmployeeID: body.EmployeeID,
		AssignedBy: body.AssignedBy,
		Notes:      body.Notes,
	}

	if err := h.service.CreateAssetAssignment(assignment); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "Asset assignment created successfully"})
}

// GetAssetAssignmentByID handles GET /asset-assignments/:id
func (h *AssetAssignmentHandler) GetAssetAssignmentByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid assignment ID")
		return
	}

	assignment, err := h.service.GetAssetAssignmentByID(uint(id))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Asset assignment not found")
		return
	}

	respondWithJSON(w, http.StatusOK, assignment)
}

// GetAssetAssignmentsByAssetID handles GET /asset-assignments/asset/:assetId
func (h *AssetAssignmentHandler) GetAssetAssignmentsByAssetID(w http.ResponseWriter, r *http.Request) {
	assetIDParam := chi.URLParam(r, "assetId")
	assetID, err := strconv.ParseUint(assetIDParam, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	assignments, err := h.service.GetAssetAssignmentsByAssetID(uint(assetID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, assignments)
}

// GetAssetAssignmentsByEmployeeID handles GET /asset-assignments/employee/:employeeId
func (h *AssetAssignmentHandler) GetAssetAssignmentsByEmployeeID(w http.ResponseWriter, r *http.Request) {
	employeeIDParam := chi.URLParam(r, "employeeId")
	employeeID, err := strconv.ParseUint(employeeIDParam, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	assignments, err := h.service.GetAssetAssignmentsByEmployeeID(uint(employeeID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, assignments)
}

// GetCurrentAssignmentsByEmployeeID handles GET /asset-assignments/current/employee/:employeeId
func (h *AssetAssignmentHandler) GetCurrentAssignmentsByEmployeeID(w http.ResponseWriter, r *http.Request) {
	employeeIDParam := chi.URLParam(r, "employeeId")
	employeeID, err := strconv.ParseUint(employeeIDParam, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	assignments, err := h.service.GetCurrentAssignmentsByEmployeeID(uint(employeeID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, assignments)
}

// GetActiveAssignmentByAssetID handles GET /asset-assignments/active/asset/:assetId
func (h *AssetAssignmentHandler) GetActiveAssignmentByAssetID(w http.ResponseWriter, r *http.Request) {
	assetIDParam := chi.URLParam(r, "assetId")
	assetID, err := strconv.ParseUint(assetIDParam, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	assignment, err := h.service.GetActiveAssignmentByAssetID(uint(assetID))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No active assignment found for this asset")
		return
	}

	respondWithJSON(w, http.StatusOK, assignment)
}

// GetAllAssetAssignments handles GET /asset-assignments
func (h *AssetAssignmentHandler) GetAllAssetAssignments(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

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
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, assignments)
}

// UpdateAssetAssignment handles PATCH /asset-assignments/:id
func (h *AssetAssignmentHandler) UpdateAssetAssignment(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid assignment ID")
		return
	}

	var body model.AssetAssignment
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	body.ID = uint(id)

	if err := h.service.UpdateAssetAssignment(&body); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Asset assignment updated successfully"})
}

// DeleteAssetAssignment handles DELETE /asset-assignments/:id
func (h *AssetAssignmentHandler) DeleteAssetAssignment(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid assignment ID")
		return
	}

	if err := h.service.DeleteAssetAssignment(uint(id)); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Asset assignment deleted successfully"})
}

// AssignAsset handles POST /assets/assign
func (h *AssetAssignmentHandler) AssignAsset(w http.ResponseWriter, r *http.Request) {
	var body struct {
		AssetID      uint   `json:"asset_id"`
		EmployeeID   uint   `json:"employee_id"`
		AssignedByID uint   `json:"assigned_by"`
		Notes        string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.AssignAsset(body.AssetID, body.EmployeeID, body.AssignedByID, body.Notes); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Asset assigned successfully"})
}

// ReturnAsset handles POST /assets/:assignmentId/return
func (h *AssetAssignmentHandler) ReturnAsset(w http.ResponseWriter, r *http.Request) {
	assignmentIDParam := chi.URLParam(r, "assignmentId")
	assignmentID, err := strconv.ParseUint(assignmentIDParam, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid assignment ID")
		return
	}

	var body struct {
		ReturnedByID uint   `json:"returned_by"`
		Notes        string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.ReturnAsset(uint(assignmentID), body.ReturnedByID, body.Notes); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Asset returned successfully"})
}