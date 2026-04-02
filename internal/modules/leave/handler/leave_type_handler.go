package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thecemakin/hr-project/internal/modules/leave/model"
	"github.com/thecemakin/hr-project/internal/modules/leave/service"
)

type LeaveTypeHandler struct {
	service service.LeaveService
}

func NewLeaveTypeHandler(service service.LeaveService) *LeaveTypeHandler {
	return &LeaveTypeHandler{service: service}
}

func (h *LeaveTypeHandler) CreateLeaveType(w http.ResponseWriter, r *http.Request) {
	var lt model.LeaveType
	if err := json.NewDecoder(r.Body).Decode(&lt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateLeaveType(&lt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(lt)
}

func (h *LeaveTypeHandler) GetLeaveType(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	lt, err := h.service.GetLeaveTypeByID(uint(id))
	if err != nil {
		http.Error(w, "leave type not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lt)
}

func (h *LeaveTypeHandler) ListLeaveTypes(w http.ResponseWriter, r *http.Request) {
	lts, err := h.service.ListLeaveTypes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lts)
}
