package handler

import (
	"github.com/go-chi/chi/v5"
)

// SetupRoutes configures all the routes for the Core HR module
func SetupRoutes(
	employeeHandler *EmployeeHandler,
	departmentHandler *DepartmentHandler,
	positionHandler *PositionHandler,
	assetHandler *AssetHandler,
	assetAssignmentHandler *AssetAssignmentHandler,
) *chi.Mux {
	r := chi.NewRouter()

	// Employees
	r.Route("/employees", func(r chi.Router) {
		r.Post("/", employeeHandler.CreateEmployee)
		r.Get("/", employeeHandler.GetAllEmployees)
		r.Get("/{id}", employeeHandler.GetEmployeeByID)
		r.Put("/{id}", employeeHandler.UpdateEmployee)
		r.Delete("/{id}", employeeHandler.DeleteEmployee)
		
		r.Get("/email/{email}", employeeHandler.GetEmployeeByEmail)
		r.Get("/manager/{managerId}", employeeHandler.GetEmployeesByManagerID)
		r.Get("/department/{departmentId}", employeeHandler.GetEmployeesByDepartmentID)
		r.Get("/status/{status}", employeeHandler.GetEmployeesByStatus)
	})

	// Departments
	r.Route("/departments", func(r chi.Router) {
		r.Post("/", departmentHandler.CreateDepartment)
		r.Get("/", departmentHandler.GetAllDepartments)
		r.Get("/{id}", departmentHandler.GetDepartmentByID)
		r.Put("/{id}", departmentHandler.UpdateDepartment)
		r.Delete("/{id}", departmentHandler.DeleteDepartment)
		
		r.Get("/name/{name}", departmentHandler.GetDepartmentByName)
	})

	// Positions
	r.Route("/positions", func(r chi.Router) {
		r.Post("/", positionHandler.CreatePosition)
		r.Get("/", positionHandler.GetAllPositions)
		r.Get("/{id}", positionHandler.GetPositionByID)
		r.Put("/{id}", positionHandler.UpdatePosition)
		r.Delete("/{id}", positionHandler.DeletePosition)
		
		r.Get("/title/{title}", positionHandler.GetPositionByTitle)
	})

	// Assets
	r.Route("/assets", func(r chi.Router) {
		r.Post("/", assetHandler.CreateAsset)
		r.Get("/", assetHandler.GetAllAssets)
		r.Get("/{id}", assetHandler.GetAssetByID)
		r.Put("/{id}", assetHandler.UpdateAsset)
		r.Delete("/{id}", assetHandler.DeleteAsset)
		
		r.Get("/serial/{serialNumber}", assetHandler.GetAssetBySerialNumber)
		r.Get("/tag/{assetTag}", assetHandler.GetAssetByAssetTag)
		r.Get("/status/{status}", assetHandler.GetAssetsByStatus)
		r.Get("/type/{type}", assetHandler.GetAssetsByType)
		
		// Assignments from asset root
		r.Post("/assign", assetAssignmentHandler.AssignAsset)
		r.Post("/{assignmentId}/return", assetAssignmentHandler.ReturnAsset)
	})

	// Asset Assignments
	r.Route("/asset-assignments", func(r chi.Router) {
		r.Post("/", assetAssignmentHandler.CreateAssetAssignment)
		r.Get("/", assetAssignmentHandler.GetAllAssetAssignments)
		r.Get("/{id}", assetAssignmentHandler.GetAssetAssignmentByID)
		r.Put("/{id}", assetAssignmentHandler.UpdateAssetAssignment)
		r.Delete("/{id}", assetAssignmentHandler.DeleteAssetAssignment)
		
		r.Get("/asset/{assetId}", assetAssignmentHandler.GetAssetAssignmentsByAssetID)
		r.Get("/employee/{employeeId}", assetAssignmentHandler.GetAssetAssignmentsByEmployeeID)
		r.Get("/current/employee/{employeeId}", assetAssignmentHandler.GetCurrentAssignmentsByEmployeeID)
		r.Get("/active/asset/{assetId}", assetAssignmentHandler.GetActiveAssignmentByAssetID)
	})

	return r
}
