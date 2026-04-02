package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/thecemakin/hr-project/internal/modules/leave/service"
)

func RegisterRoutes(r chi.Router, svc service.LeaveService) {
	leaveTypeHandler := NewLeaveTypeHandler(svc)
	leaveBalanceHandler := NewLeaveBalanceHandler(svc)
	leaveRequestHandler := NewLeaveRequestHandler(svc)

	// Leave Types
	r.Route("/leave-types", func(r chi.Router) {
		r.Get("/", leaveTypeHandler.ListLeaveTypes)
		r.Post("/", leaveTypeHandler.CreateLeaveType)
		r.Get("/{id}", leaveTypeHandler.GetLeaveType)
	})

	// Leave Balances
	r.Route("/leave-balances", func(r chi.Router) {
		r.Get("/{employeeId}", leaveBalanceHandler.ListLeaveBalancesByEmployee)
		r.Post("/init", leaveBalanceHandler.InitializeBalance)
	})

	// Leave Requests
	r.Route("/leave-requests", func(r chi.Router) {
		r.Post("/", leaveRequestHandler.SubmitRequest)
		r.Get("/me", leaveRequestHandler.ListOwnRequests)
		r.Get("/pending-approvals", leaveRequestHandler.ListPendingApprovals)
		r.Post("/{id}/approve", leaveRequestHandler.ApproveRequest)
		r.Post("/{id}/reject", leaveRequestHandler.RejectRequest)
	})
}
