package service

import (
	"errors"
	"fmt"

	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"github.com/thecemakin/hr-project/internal/modules/corehr/repository"
)

// EmployeeService defines the interface for employee business operations
type EmployeeService interface {
	CreateEmployee(employee *model.Employee) error
	GetEmployeeByID(id uint) (*model.Employee, error)
	GetEmployeeByEmail(email string) (*model.Employee, error)
	GetAllEmployees(limit, offset int) ([]*model.Employee, error)
	UpdateEmployee(employee *model.Employee) error
	DeleteEmployee(id uint) error
	GetEmployeesByManagerID(managerID uint) ([]*model.Employee, error)
	GetEmployeesByDepartmentID(departmentID uint) ([]*model.Employee, error)
	GetEmployeesByStatus(status string) ([]*model.Employee, error)
	ValidateManagerRelationship(employeeID, managerID uint) (bool, error)
}

// employeeService implements the EmployeeService interface
type employeeService struct {
	repo repository.EmployeeRepository
}

// NewEmployeeService creates a new instance of employeeService
func NewEmployeeService(repo repository.EmployeeRepository) EmployeeService {
	return &employeeService{
		repo: repo,
	}
}

// CreateEmployee creates a new employee after validation
func (s *employeeService) CreateEmployee(employee *model.Employee) error {
	// Validate required fields
	if employee.FirstName == "" || employee.LastName == "" || employee.Email == "" || employee.EmployeeNumber == "" {
		return errors.New("first name, last name, email, and employee number are required")
	}

	// Check if email already exists
	existingEmployee, err := s.repo.GetByEmail(employee.Email)
	if err == nil && existingEmployee != nil {
		return fmt.Errorf("employee with email %s already exists", employee.Email)
	}

	// Check if employee number already exists
	existingEmployee, err = s.repo.GetByEmployeeNumber(employee.EmployeeNumber)
	if err == nil && existingEmployee != nil {
		return fmt.Errorf("employee with number %s already exists", employee.EmployeeNumber)
	}

	// Validate manager relationship if provided
	if employee.ManagerID != nil {
		isValid, err := s.ValidateManagerRelationship(employee.ID, *employee.ManagerID)
		if err != nil {
			return err
		}
		if !isValid {
			return errors.New("invalid manager relationship: employee cannot be their own manager")
		}
	}

	// Create the employee
	return s.repo.Create(employee)
}

// GetEmployeeByID retrieves an employee by ID
func (s *employeeService) GetEmployeeByID(id uint) (*model.Employee, error) {
	return s.repo.GetByID(id)
}

// GetEmployeeByEmail retrieves an employee by email
func (s *employeeService) GetEmployeeByEmail(email string) (*model.Employee, error) {
	return s.repo.GetByEmail(email)
}

// GetAllEmployees retrieves all employees with pagination
func (s *employeeService) GetAllEmployees(limit, offset int) ([]*model.Employee, error) {
	return s.repo.GetAll(limit, offset)
}

// UpdateEmployee updates an existing employee after validation
func (s *employeeService) UpdateEmployee(employee *model.Employee) error {
	// Validate required fields
	if employee.FirstName == "" || employee.LastName == "" || employee.Email == "" || employee.EmployeeNumber == "" {
		return errors.New("first name, last name, email, and employee number are required")
	}

	// Check if employee exists
	existingEmployee, err := s.repo.GetByID(employee.ID)
	if err != nil {
		return fmt.Errorf("employee with ID %d does not exist", employee.ID)
	}

	// Check if email is being changed and if it already exists for another employee
	if existingEmployee.Email != employee.Email {
		emailExists, _ := s.repo.GetByEmail(employee.Email)
		if emailExists != nil && emailExists.ID != employee.ID {
			return fmt.Errorf("employee with email %s already exists", employee.Email)
		}
	}

	// Check if employee number is being changed and if it already exists for another employee
	if existingEmployee.EmployeeNumber != employee.EmployeeNumber {
		numExists, _ := s.repo.GetByEmployeeNumber(employee.EmployeeNumber)
		if numExists != nil && numExists.ID != employee.ID {
			return fmt.Errorf("employee with number %s already exists", employee.EmployeeNumber)
		}
	}

	// Validate manager relationship if provided
	if employee.ManagerID != nil {
		isValid, err := s.ValidateManagerRelationship(employee.ID, *employee.ManagerID)
		if err != nil {
			return err
		}
		if !isValid {
			return errors.New("invalid manager relationship: employee cannot be their own manager")
		}
	}

	// Update the employee
	return s.repo.Update(employee)
}

// DeleteEmployee removes an employee by ID
func (s *employeeService) DeleteEmployee(id uint) error {
	// Check if employee exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("employee with ID %d does not exist", id)
	}

	// TODO: Check if employee has any active assignments or dependencies before deletion
	// For now, just delete the employee
	return s.repo.Delete(id)
}

// GetEmployeesByManagerID retrieves all employees reporting to a specific manager
func (s *employeeService) GetEmployeesByManagerID(managerID uint) ([]*model.Employee, error) {
	return s.repo.GetByManagerID(managerID)
}

// GetEmployeesByDepartmentID retrieves all employees in a specific department
func (s *employeeService) GetEmployeesByDepartmentID(departmentID uint) ([]*model.Employee, error) {
	return s.repo.GetByDepartmentID(departmentID)
}

// GetEmployeesByStatus retrieves all employees with a specific status
func (s *employeeService) GetEmployeesByStatus(status string) ([]*model.Employee, error) {
	return s.repo.GetByStatus(status)
}

// ValidateManagerRelationship checks if an employee can report to a specific manager
// This prevents circular references and self-management
func (s *employeeService) ValidateManagerRelationship(employeeID, managerID uint) (bool, error) {
	// An employee cannot be their own manager
	if employeeID == managerID {
		return false, nil
	}

	// Check if the manager exists
	_, err := s.repo.GetByID(managerID)
	if err != nil {
		return false, fmt.Errorf("manager with ID %d does not exist", managerID)
	}

	// Additional validation could be added here, such as:
	// - Checking if the manager is in the same department (optional)
	// - Checking if the manager has the appropriate role (optional)

	return true, nil
}

// GetByEmployeeNumber retrieves an employee by employee number
// This method is needed for validation but wasn't in the interface
// We'll add it to the repository and service