package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecemakin/hr-project/internal/platform/auth"
	"github.com/thecemakin/hr-project/internal/platform/http/middleware"
)

// SetupRoutesFiber configures all the routes for the Core HR module using Fiber
func SetupRoutesFiber(
	app *fiber.App,
	employeeHandler *EmployeeHandler,
	departmentHandler *DepartmentHandler,
	positionHandler *PositionHandler,
	assetHandler *AssetHandler,
	assetAssignmentHandler *AssetAssignmentHandler,
	organizationHandler *OrganizationHandler,
	tp *auth.TokenProvider,
	appEnv string,
	skipAuth bool,
) {
	// Create group for Core HR routes and apply authentication
	v1 := app.Group("/api/v1/corehr", middleware.AuthMiddleware(tp, appEnv, skipAuth))

	// Organization
	v1.Get("/organization/tree", organizationHandler.GetTree)

	// Employees
	v1.Post("/employees", middleware.RequireRole("admin", "hr"), employeeHandler.CreateEmployee)
	v1.Get("/employees", employeeHandler.GetAllEmployees)
	v1.Get("/employees/:id", employeeHandler.GetEmployeeByID)
	v1.Put("/employees/:id", middleware.RequireRole("admin", "hr"), employeeHandler.UpdateEmployee)
	v1.Delete("/employees/:id", middleware.RequireRole("admin"), employeeHandler.DeleteEmployee)

	// Departments
	v1.Post("/departments", departmentHandler.CreateDepartment)
	v1.Get("/departments", departmentHandler.GetAllDepartments)
	v1.Get("/departments/:id", departmentHandler.GetDepartmentByID)
	v1.Put("/departments/:id", departmentHandler.UpdateDepartment)
	v1.Delete("/departments/:id", departmentHandler.DeleteDepartment)

	// Positions
	v1.Post("/positions", positionHandler.CreatePosition)
	v1.Get("/positions", positionHandler.GetAllPositions)
	v1.Get("/positions/:id", positionHandler.GetPositionByID)
	v1.Put("/positions/:id", positionHandler.UpdatePosition)
	v1.Delete("/positions/:id", positionHandler.DeletePosition)

	// Assets
	v1.Post("/assets", assetHandler.CreateAsset)
	v1.Get("/assets", assetHandler.GetAllAssets)
	v1.Get("/assets/:id", assetHandler.GetAssetByID)
	v1.Put("/assets/:id", assetHandler.UpdateAsset)
	v1.Delete("/assets/:id", assetHandler.DeleteAsset)
	
	// Assignments from asset root
	v1.Post("/assets/assign", assetAssignmentHandler.AssignAsset)
	v1.Post("/assets/:assignmentId/return", assetAssignmentHandler.ReturnAsset)

	// Asset Assignments
	v1.Post("/asset-assignments", assetAssignmentHandler.CreateAssetAssignment)
	v1.Get("/asset-assignments", assetAssignmentHandler.GetAllAssetAssignments)
	v1.Get("/asset-assignments/:id", assetAssignmentHandler.GetAssetAssignmentByID)
	v1.Put("/asset-assignments/:id", assetAssignmentHandler.UpdateAssetAssignment)
	v1.Delete("/asset-assignments/:id", assetAssignmentHandler.DeleteAssetAssignment)
}

// SetupRoutes configures all the routes for the Core HR module (for backward compatibility)
func SetupRoutes(
	employeeHandler *EmployeeHandler,
	departmentHandler *DepartmentHandler,
	positionHandler *PositionHandler,
	assetHandler *AssetHandler,
	assetAssignmentHandler *AssetAssignmentHandler,
) interface{} {
	// This function is kept for backward compatibility but should not be used with Fiber
	return nil
}
