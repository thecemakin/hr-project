package handler

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutesFiber configures all the routes for the Core HR module using Fiber
func SetupRoutesFiber(
	app *fiber.App,
	employeeHandler *EmployeeHandler,
	departmentHandler *DepartmentHandler,
	positionHandler *PositionHandler,
	assetHandler *AssetHandler,
	assetAssignmentHandler *AssetAssignmentHandler,
) {
	// Create group for Core HR routes
	v1 := app.Group("/api/v1/corehr")

	// Employees
	v1.Post("/employees", employeeHandler.CreateEmployee)
	v1.Get("/employees", employeeHandler.GetAllEmployees)
	v1.Get("/employees/:id", employeeHandler.GetEmployeeByID)
	v1.Put("/employees/:id", employeeHandler.UpdateEmployee)
	v1.Delete("/employees/:id", employeeHandler.DeleteEmployee)
	
	v1.Get("/employees/email/:email", employeeHandler.GetEmployeeByEmail)
	v1.Get("/employees/manager/:managerId", employeeHandler.GetEmployeesByManagerID)
	v1.Get("/employees/department/:departmentId", employeeHandler.GetEmployeesByDepartmentID)
	v1.Get("/employees/status/:status", employeeHandler.GetEmployeesByStatus)

	// Departments
	v1.Post("/departments", departmentHandler.CreateDepartment)
	v1.Get("/departments", departmentHandler.GetAllDepartments)
	v1.Get("/departments/:id", departmentHandler.GetDepartmentByID)
	v1.Put("/departments/:id", departmentHandler.UpdateDepartment)
	v1.Delete("/departments/:id", departmentHandler.DeleteDepartment)
	
	v1.Get("/departments/name/:name", departmentHandler.GetDepartmentByName)

	// Positions
	v1.Post("/positions", positionHandler.CreatePosition)
	v1.Get("/positions", positionHandler.GetAllPositions)
	v1.Get("/positions/:id", positionHandler.GetPositionByID)
	v1.Put("/positions/:id", positionHandler.UpdatePosition)
	v1.Delete("/positions/:id", positionHandler.DeletePosition)
	
	v1.Get("/positions/title/:title", positionHandler.GetPositionByTitle)

	// Assets
	v1.Post("/assets", assetHandler.CreateAsset)
	v1.Get("/assets", assetHandler.GetAllAssets)
	v1.Get("/assets/:id", assetHandler.GetAssetByID)
	v1.Put("/assets/:id", assetHandler.UpdateAsset)
	v1.Delete("/assets/:id", assetHandler.DeleteAsset)
	
	v1.Get("/assets/serial/:serialNumber", assetHandler.GetAssetBySerialNumber)
	v1.Get("/assets/tag/:assetTag", assetHandler.GetAssetByAssetTag)
	v1.Get("/assets/status/:status", assetHandler.GetAssetsByStatus)
	v1.Get("/assets/type/:type", assetHandler.GetAssetsByType)
	
	// Assignments from asset root
	v1.Post("/assets/assign", assetAssignmentHandler.AssignAsset)
	v1.Post("/assets/:assignmentId/return", assetAssignmentHandler.ReturnAsset)

	// Asset Assignments
	v1.Post("/asset-assignments", assetAssignmentHandler.CreateAssetAssignment)
	v1.Get("/asset-assignments", assetAssignmentHandler.GetAllAssetAssignments)
	v1.Get("/asset-assignments/:id", assetAssignmentHandler.GetAssetAssignmentByID)
	v1.Put("/asset-assignments/:id", assetAssignmentHandler.UpdateAssetAssignment)
	v1.Delete("/asset-assignments/:id", assetAssignmentHandler.DeleteAssetAssignment)
	
	v1.Get("/asset-assignments/asset/:assetId", assetAssignmentHandler.GetAssetAssignmentsByAssetID)
	v1.Get("/asset-assignments/employee/:employeeId", assetAssignmentHandler.GetAssetAssignmentsByEmployeeID)
	v1.Get("/asset-assignments/current/employee/:employeeId", assetAssignmentHandler.GetCurrentAssignmentsByEmployeeID)
	v1.Get("/asset-assignments/active/asset/:assetId", assetAssignmentHandler.GetActiveAssignmentByAssetID)
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
