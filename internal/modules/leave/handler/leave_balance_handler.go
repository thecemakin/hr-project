package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thecemakin/hr-project/internal/modules/leave/service"
)

type LeaveBalanceHandler struct {
	service service.LeaveService
}

func NewLeaveBalanceHandler(service service.LeaveService) *LeaveBalanceHandler {
	return &LeaveBalanceHandler{service: service}
}

func (h *LeaveBalanceHandler) ListLeaveBalancesByEmployee(w http.ResponseWriter, r *http.Request) {
	empIDStr := chi.URLParam(r, "employeeId")
	empID, err := strconv.ParseUint(empIDStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid employee id", http.StatusBadRequest)
		return
	}

	lbs, err := h.service.ListLeaveBalancesByEmployee(uint(empID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lbs)
}

func (h *LeaveBalanceHandler) InitializeBalance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EmployeeID   uint `json:"employee_id"`
		LeaveTypeID  uint `json:"leave_type_id"`
		StartingDays int  `json:"starting_days"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.InitializeLeaveBalance(req.EmployeeID, req.LeaveTypeID, req.StartingDays); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"success"}`))
}
