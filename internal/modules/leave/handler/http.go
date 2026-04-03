package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecemakin/hr-project/internal/modules/leave/service"
	"github.com/thecemakin/hr-project/internal/platform/auth"
	"github.com/thecemakin/hr-project/internal/platform/http/middleware"
)

// SetupRoutesFiber configures all the routes for the Leave module using Fiber
func SetupRoutesFiber(app *fiber.App, svc service.LeaveService, tp *auth.TokenProvider) {
	leaveTypeHandler := NewLeaveTypeHandler(svc)
	leaveBalanceHandler := NewLeaveBalanceHandler(svc)
	leaveRequestHandler := NewLeaveRequestHandler(svc)

	// Create group for Leave routes and apply authentication
	v1 := app.Group("/api/v1/leave", middleware.AuthMiddleware(tp))

	// Leave Types
	v1.Get("/leave-types", leaveTypeHandler.ListLeaveTypes)
	v1.Post("/leave-types", leaveTypeHandler.CreateLeaveType)
	v1.Get("/leave-types/:id", leaveTypeHandler.GetLeaveType)

	// Leave Balances
	v1.Get("/leave-balances/:employeeId", leaveBalanceHandler.ListLeaveBalancesByEmployee)
	v1.Post("/leave-balances/init", leaveBalanceHandler.InitializeBalance)

	// Leave Requests
	v1.Post("/leave-requests", leaveRequestHandler.SubmitRequest)
	v1.Get("/leave-requests/me", leaveRequestHandler.ListOwnRequests)
	v1.Get("/leave-requests/pending-approvals", leaveRequestHandler.ListPendingApprovals)
	v1.Post("/leave-requests/:id/approve", leaveRequestHandler.ApproveRequest)
	v1.Post("/leave-requests/:id/reject", leaveRequestHandler.RejectRequest)
}

// RegisterRoutes is kept for backward compatibility but should not be used with Fiber
func RegisterRoutes(r interface{}, svc service.LeaveService) {
	// This function is kept for backward compatibility but should not be used with Fiber
	// The Fiber routes are registered via SetupRoutesFiber
	// For now, leaving empty as we're migrating to Fiber
}
