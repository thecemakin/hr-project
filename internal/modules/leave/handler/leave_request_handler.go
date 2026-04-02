package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thecemakin/hr-project/internal/modules/leave/model"
	"github.com/thecemakin/hr-project/internal/modules/leave/service"
)

type LeaveRequestHandler struct {
	service service.LeaveService
}

func NewLeaveRequestHandler(service service.LeaveService) *LeaveRequestHandler {
	return &LeaveRequestHandler{service: service}
}

func (h *LeaveRequestHandler) SubmitRequest(w http.ResponseWriter, r *http.Request) {
	var req model.LeaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.SubmitLeaveRequest(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func (h *LeaveRequestHandler) ListOwnRequests(w http.ResponseWriter, r *http.Request) {
	// In a real app, empID would come from JWT / auth context.
	// For testing, we might pass it via a query param or header until Auth is built.
	empIDStr := r.URL.Query().Get("employeeId")
	empID, err := strconv.ParseUint(empIDStr, 10, 32)
	if err != nil {
		http.Error(w, "missing or invalid employeeId query param", http.StatusUnauthorized)
		return
	}

	reqs, err := h.service.ListOwnLeaveRequests(uint(empID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reqs)
}

func (h *LeaveRequestHandler) ListPendingApprovals(w http.ResponseWriter, r *http.Request) {
	// From JWT context representing the logged-in manager
	managerIDStr := r.URL.Query().Get("managerId")
	managerID, err := strconv.ParseUint(managerIDStr, 10, 32)
	if err != nil {
		http.Error(w, "missing or invalid managerId query param", http.StatusUnauthorized)
		return
	}

	reqs, err := h.service.ListPendingApprovals(uint(managerID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reqs)
}

func (h *LeaveRequestHandler) ApproveRequest(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// From JWT context
	managerIDStr := r.URL.Query().Get("managerId")
	managerID, err := strconv.ParseUint(managerIDStr, 10, 32)
	if err != nil {
		http.Error(w, "missing or invalid managerId query param", http.StatusUnauthorized)
		return
	}

	var payload struct {
		Note string `json:"note"`
	}
	// Note is optional
	_ = json.NewDecoder(r.Body).Decode(&payload)

	if err := h.service.ApproveLeaveRequest(uint(id), uint(managerID), payload.Note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"approved"}`))
}

func (h *LeaveRequestHandler) RejectRequest(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// From JWT context
	managerIDStr := r.URL.Query().Get("managerId")
	managerID, err := strconv.ParseUint(managerIDStr, 10, 32)
	if err != nil {
		http.Error(w, "missing or invalid managerId query param", http.StatusUnauthorized)
		return
	}

	var payload struct {
		Note string `json:"note"`
	}
	// Note is optional
	_ = json.NewDecoder(r.Body).Decode(&payload)

	if err := h.service.RejectLeaveRequest(uint(id), uint(managerID), payload.Note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"rejected"}`))
}
