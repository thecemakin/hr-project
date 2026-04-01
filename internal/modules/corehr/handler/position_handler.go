package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"github.com/thecemakin/hr-project/internal/modules/corehr/service"
)

// PositionHandler handles HTTP requests for position operations
type PositionHandler struct {
	service service.PositionService
}

// NewPositionHandler creates a new instance of PositionHandler
func NewPositionHandler(service service.PositionService) *PositionHandler {
	return &PositionHandler{
		service: service,
	}
}

// CreatePosition handles POST /positions
func (h *PositionHandler) CreatePosition(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Code        string `json:"code"`
		Level       int    `json:"level"`
		SalaryMin   int64  `json:"salary_min"`
		SalaryMax   int64  `json:"salary_max"`
		IsActive    bool   `json:"is_active"`
		DepartmentID *uint `json:"department_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	position := &model.Position{
		Title:       body.Title,
		Description: body.Description,
		Code:        body.Code,
		Level:       body.Level,
		SalaryMin:   body.SalaryMin,
		SalaryMax:   body.SalaryMax,
		IsActive:    body.IsActive,
		DepartmentID: body.DepartmentID,
	}

	if err := h.service.CreatePosition(position); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "Position created successfully"})
}

// GetPositionByID handles GET /positions/:id
func (h *PositionHandler) GetPositionByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid position ID")
		return
	}

	position, err := h.service.GetPositionByID(uint(id))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Position not found")
		return
	}

	respondWithJSON(w, http.StatusOK, position)
}

// GetPositionByTitle handles GET /positions/title/:title
func (h *PositionHandler) GetPositionByTitle(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	position, err := h.service.GetPositionByTitle(title)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Position not found")
		return
	}

	respondWithJSON(w, http.StatusOK, position)
}

// GetAllPositions handles GET /positions
func (h *PositionHandler) GetAllPositions(w http.ResponseWriter, r *http.Request) {
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

	positions, err := h.service.GetAllPositions(limit, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, positions)
}

// UpdatePosition handles PATCH /positions/:id
func (h *PositionHandler) UpdatePosition(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid position ID")
		return
	}

	var body model.Position
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	body.ID = uint(id)

	if err := h.service.UpdatePosition(&body); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Position updated successfully"})
}

// DeletePosition handles DELETE /positions/:id
func (h *PositionHandler) DeletePosition(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid position ID")
		return
	}

	if err := h.service.DeletePosition(uint(id)); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Position deleted successfully"})
}