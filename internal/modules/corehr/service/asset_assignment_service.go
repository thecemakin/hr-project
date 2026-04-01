package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"github.com/thecemakin/hr-project/internal/modules/corehr/repository"
)

// AssetAssignmentService defines the interface for asset assignment business operations
type AssetAssignmentService interface {
	CreateAssetAssignment(assignment *model.AssetAssignment) error
	GetAssetAssignmentByID(id uint) (*model.AssetAssignment, error)
	GetAssetAssignmentsByAssetID(assetID uint) ([]*model.AssetAssignment, error)
	GetAssetAssignmentsByEmployeeID(employeeID uint) ([]*model.AssetAssignment, error)
	GetCurrentAssignmentsByEmployeeID(employeeID uint) ([]*model.AssetAssignment, error)
	GetActiveAssignmentByAssetID(assetID uint) (*model.AssetAssignment, error)
	GetAllAssetAssignments(limit, offset int) ([]*model.AssetAssignment, error)
	UpdateAssetAssignment(assignment *model.AssetAssignment) error
	DeleteAssetAssignment(id uint) error
	AssignAsset(assetID, employeeID, assignedByID uint, notes string) error
	ReturnAsset(assignmentID, returnedByID uint, notes string) error
}

// assetAssignmentService implements the AssetAssignmentService interface
type assetAssignmentService struct {
	repo repository.AssetAssignmentRepository
	assetRepo repository.AssetRepository
	employeeRepo repository.EmployeeRepository
}

// NewAssetAssignmentService creates a new instance of assetAssignmentService
func NewAssetAssignmentService(
	assignmentRepo repository.AssetAssignmentRepository,
	assetRepo repository.AssetRepository,
	employeeRepo repository.EmployeeRepository,
) AssetAssignmentService {
	return &assetAssignmentService{
		repo: assignmentRepo,
		assetRepo: assetRepo,
		employeeRepo: employeeRepo,
	}
}

// CreateAssetAssignment creates a new asset assignment after validation
func (s *assetAssignmentService) CreateAssetAssignment(assignment *model.AssetAssignment) error {
	// Validate required fields
	if assignment.AssetID == 0 || assignment.EmployeeID == 0 || assignment.AssignedBy == 0 {
		return errors.New("asset ID, employee ID, and assigned by are required")
	}

	// Validate that asset exists
	_, err := s.assetRepo.GetByID(assignment.AssetID)
	if err != nil {
		return fmt.Errorf("asset with ID %d does not exist", assignment.AssetID)
	}

	// Validate that employee exists
	_, err = s.employeeRepo.GetByID(assignment.EmployeeID)
	if err != nil {
		return fmt.Errorf("employee with ID %d does not exist", assignment.EmployeeID)
	}

	// Validate that assigner exists
	_, err = s.employeeRepo.GetByID(assignment.AssignedBy)
	if err != nil {
		return fmt.Errorf("assigner with ID %d does not exist", assignment.AssignedBy)
	}

	// Check if asset is available for assignment (not already assigned)
	activeAssignment, err := s.repo.GetActiveAssignmentsByAssetID(assignment.AssetID)
	if err == nil && activeAssignment != nil {
		return fmt.Errorf("asset with ID %d is already assigned to employee with ID %d", assignment.AssetID, activeAssignment.EmployeeID)
	}

	// Set status to assigned if not already set
	if assignment.Status == "" {
		assignment.Status = "assigned"
	}

	// Set assigned date if not already set
	if assignment.AssignedDate.IsZero() {
		assignment.AssignedDate = time.Now()
	}

	// Create the assignment
	return s.repo.Create(assignment)
}

// GetAssetAssignmentByID retrieves an asset assignment by ID
func (s *assetAssignmentService) GetAssetAssignmentByID(id uint) (*model.AssetAssignment, error) {
	return s.repo.GetByID(id)
}

// GetAssetAssignmentsByAssetID retrieves all assignments for a specific asset
func (s *assetAssignmentService) GetAssetAssignmentsByAssetID(assetID uint) ([]*model.AssetAssignment, error) {
	return s.repo.GetByAssetID(assetID)
}

// GetAssetAssignmentsByEmployeeID retrieves all assignments for a specific employee
func (s *assetAssignmentService) GetAssetAssignmentsByEmployeeID(employeeID uint) ([]*model.AssetAssignment, error) {
	return s.repo.GetByEmployeeID(employeeID)
}

// GetCurrentAssignmentsByEmployeeID retrieves all current (non-returned) assignments for a specific employee
func (s *assetAssignmentService) GetCurrentAssignmentsByEmployeeID(employeeID uint) ([]*model.AssetAssignment, error) {
	return s.repo.GetCurrentAssignmentsByEmployeeID(employeeID)
}

// GetActiveAssignmentByAssetID retrieves the currently active assignment for a specific asset
func (s *assetAssignmentService) GetActiveAssignmentByAssetID(assetID uint) (*model.AssetAssignment, error) {
	return s.repo.GetActiveAssignmentsByAssetID(assetID)
}

// GetAllAssetAssignments retrieves all asset assignments with pagination
func (s *assetAssignmentService) GetAllAssetAssignments(limit, offset int) ([]*model.AssetAssignment, error) {
	return s.repo.GetAll(limit, offset)
}

// UpdateAssetAssignment updates an existing asset assignment after validation
func (s *assetAssignmentService) UpdateAssetAssignment(assignment *model.AssetAssignment) error {
	// Validate required fields
	if assignment.AssetID == 0 || assignment.EmployeeID == 0 || assignment.AssignedBy == 0 {
		return errors.New("asset ID, employee ID, and assigned by are required")
	}

	// Check if assignment exists
	_, err := s.repo.GetByID(assignment.ID)
	if err != nil {
		return fmt.Errorf("asset assignment with ID %d does not exist", assignment.ID)
	}

	// Update the assignment
	return s.repo.Update(assignment)
}

// DeleteAssetAssignment removes an asset assignment by ID
func (s *assetAssignmentService) DeleteAssetAssignment(id uint) error {
	// Check if assignment exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("asset assignment with ID %d does not exist", id)
	}

	// TODO: Check if assignment is active before deletion
	// For now, just delete the assignment
	return s.repo.Delete(id)
}

// AssignAsset assigns an asset to an employee
func (s *assetAssignmentService) AssignAsset(assetID, employeeID, assignedByID uint, notes string) error {
	// Validate that asset exists
	_, err := s.assetRepo.GetByID(assetID)
	if err != nil {
		return fmt.Errorf("asset with ID %d does not exist", assetID)
	}

	// Validate that employee exists
	_, err = s.employeeRepo.GetByID(employeeID)
	if err != nil {
		return fmt.Errorf("employee with ID %d does not exist", employeeID)
	}

	// Validate that assigner exists
	_, err = s.employeeRepo.GetByID(assignedByID)
	if err != nil {
		return fmt.Errorf("assigner with ID %d does not exist", assignedByID)
	}

	// Check if asset is available for assignment (not already assigned)
	activeAssignment, err := s.repo.GetActiveAssignmentsByAssetID(assetID)
	if err == nil && activeAssignment != nil {
		return fmt.Errorf("asset with ID %d is already assigned to employee with ID %d", assetID, activeAssignment.EmployeeID)
	}

	// Create the assignment
	assignment := &model.AssetAssignment{
		AssetID:      assetID,
		EmployeeID:   employeeID,
		AssignedBy:   assignedByID,
		AssignedDate: time.Now(),
		Status:       "assigned",
		Notes:        notes,
	}

	return s.repo.Create(assignment)
}

// ReturnAsset returns an assigned asset
func (s *assetAssignmentService) ReturnAsset(assignmentID, returnedByID uint, notes string) error {
	// Get the assignment
	assignment, err := s.repo.GetByID(assignmentID)
	if err != nil {
		return fmt.Errorf("asset assignment with ID %d does not exist", assignmentID)
	}

	// Validate that returner exists
	_, err = s.employeeRepo.GetByID(returnedByID)
	if err != nil {
		return fmt.Errorf("returner with ID %d does not exist", returnedByID)
	}

	// Check if assignment is already returned
	if assignment.Status == "returned" {
		return fmt.Errorf("asset assignment with ID %d is already returned", assignmentID)
	}

	// Update the assignment
	now := time.Now()
	assignment.ReturnedDate = &now
	assignment.ReturnedBy = &returnedByID
	assignment.Status = "returned"
	
	// Append notes if provided
	if notes != "" {
		if assignment.Notes != "" {
			assignment.Notes += "; " + notes
		} else {
			assignment.Notes = notes
		}
	}

	return s.repo.Update(assignment)
}