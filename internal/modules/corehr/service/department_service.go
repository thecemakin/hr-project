package service

import (
	"errors"
	"fmt"

	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"github.com/thecemakin/hr-project/internal/modules/corehr/repository"
)

// DepartmentService defines the interface for department business operations
type DepartmentService interface {
	CreateDepartment(department *model.Department) error
	GetDepartmentByID(id uint) (*model.Department, error)
	GetDepartmentByName(name string) (*model.Department, error)
	GetAllDepartments(limit, offset int) ([]*model.Department, error)
	UpdateDepartment(department *model.Department) error
	DeleteDepartment(id uint) error
}

// departmentService implements the DepartmentService interface
type departmentService struct {
	repo repository.DepartmentRepository
}

// NewDepartmentService creates a new instance of departmentService
func NewDepartmentService(repo repository.DepartmentRepository) DepartmentService {
	return &departmentService{
		repo: repo,
	}
}

// CreateDepartment creates a new department after validation
func (s *departmentService) CreateDepartment(department *model.Department) error {
	// Validate required fields
	if department.Name == "" {
		return errors.New("department name is required")
	}

	// Check if department name already exists
	existingDept, err := s.repo.GetByName(department.Name)
	if err == nil && existingDept != nil {
		return fmt.Errorf("department with name %s already exists", department.Name)
	}

	// Validate head if provided
	if department.HeadID != nil {
		// In a real implementation, we would validate that the head exists
		// For now, we'll just proceed with creation
	}

	// Create the department
	return s.repo.Create(department)
}

// GetDepartmentByID retrieves a department by ID
func (s *departmentService) GetDepartmentByID(id uint) (*model.Department, error) {
	return s.repo.GetByID(id)
}

// GetDepartmentByName retrieves a department by name
func (s *departmentService) GetDepartmentByName(name string) (*model.Department, error) {
	return s.repo.GetByName(name)
}

// GetAllDepartments retrieves all departments with pagination
func (s *departmentService) GetAllDepartments(limit, offset int) ([]*model.Department, error) {
	return s.repo.GetAll(limit, offset)
}

// UpdateDepartment updates an existing department after validation
func (s *departmentService) UpdateDepartment(department *model.Department) error {
	// Validate required fields
	if department.Name == "" {
		return errors.New("department name is required")
	}

	// Check if department exists
	existingDept, err := s.repo.GetByID(department.ID)
	if err != nil {
		return fmt.Errorf("department with ID %d does not exist", department.ID)
	}

	// Check if name is being changed and if it already exists for another department
	if existingDept.Name != department.Name {
		nameExists, _ := s.repo.GetByName(department.Name)
		if nameExists != nil && nameExists.ID != department.ID {
			return fmt.Errorf("department with name %s already exists", department.Name)
		}
	}

	// Update the department
	return s.repo.Update(department)
}

// DeleteDepartment removes a department by ID
func (s *departmentService) DeleteDepartment(id uint) error {
	// Check if department exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("department with ID %d does not exist", id)
	}

	// TODO: Check if department has any employees before deletion
	// For now, just delete the department
	return s.repo.Delete(id)
}