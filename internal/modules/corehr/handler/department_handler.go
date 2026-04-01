package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"github.com/thecemakin/hr-project/internal/modules/corehr/service"
)

// DepartmentHandler handles HTTP requests for department operations
type DepartmentHandler struct {
	service service.DepartmentService
}

// NewDepartmentHandler creates a new instance of DepartmentHandler
func NewDepartmentHandler(service service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{
		service: service,
	}
}

// CreateDepartment handles POST /departments
func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Code        string `json:"code"`
		HeadID      *uint  `json:"head_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	department := &model.Department{
		Name:        body.Name,
		Description: body.Description,
		Code:        body.Code,
		HeadID:      body.HeadID,
	}

	if err := h.service.CreateDepartment(department); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "Department created successfully"})
}

// GetDepartmentByID handles GET /departments/:id
func (h *DepartmentHandler) GetDepartmentByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid department ID")
		return
	}

	department, err := h.service.GetDepartmentByID(uint(id))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Department not found")
		return
	}

	respondWithJSON(w, http.StatusOK, department)
}

// GetDepartmentByName handles GET /departments/name/:name
func (h *DepartmentHandler) GetDepartmentByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	department, err := h.service.GetDepartmentByName(name)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Department not found")
		return
	}

	respondWithJSON(w, http.StatusOK, department)
}

// GetAllDepartments handles GET /departments
func (h *DepartmentHandler) GetAllDepartments(w http.ResponseWriter, r *http.Request) {
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

	departments, err := h.service.GetAllDepartments(limit, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, departments)
}

// UpdateDepartment handles PATCH /departments/:id
func (h *DepartmentHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid department ID")
		return
	}

	var body model.Department
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	body.ID = uint(id)

	if err := h.service.UpdateDepartment(&body); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Department updated successfully"})
}

// DeleteDepartment handles DELETE /departments/:id
func (h *DepartmentHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid department ID")
		return
	}

	if err := h.service.DeleteDepartment(uint(id)); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Department deleted successfully"})
}