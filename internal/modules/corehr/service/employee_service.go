package service

import (
	"errors"
	"fmt"

	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"github.com/thecemakin/hr-project/internal/modules/corehr/repository"
	"github.com/thecemakin/hr-project/internal/platform/audit"
)

// EmployeeService defines the interface for employee business operations
type EmployeeService interface {
	CreateEmployee(actorID uint, employee *model.Employee) error
	GetEmployeeByID(id uint) (*model.Employee, error)
	GetEmployeeByEmail(email string) (*model.Employee, error)
	GetAllEmployees(limit, offset int, filter repository.EmployeeFilter) ([]*model.Employee, error)
	UpdateEmployee(actorID uint, employee *model.Employee) error
	DeleteEmployee(actorID uint, id uint) error
	ValidateManagerRelationship(employeeID, managerID uint) (bool, error)
	GetOrganizationTree() ([]*model.OrganizationNode, error)
}

// employeeService implements the EmployeeService interface
type employeeService struct {
	repo  repository.EmployeeRepository
	audit audit.Service
}

// NewEmployeeService creates a new instance of employeeService
func NewEmployeeService(repo repository.EmployeeRepository, audit audit.Service) EmployeeService {
	return &employeeService{
		repo:  repo,
		audit: audit,
	}
}

// CreateEmployee creates a new employee after validation
func (s *employeeService) CreateEmployee(actorID uint, employee *model.Employee) error {
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
	if err := s.repo.Create(employee); err != nil {
		return err
	}

	// Audit log
	_ = s.audit.Log(actorID, "CREATE", "employees", employee.ID, nil, employee, nil)

	return nil
}

// GetEmployeeByID retrieves an employee by ID
func (s *employeeService) GetEmployeeByID(id uint) (*model.Employee, error) {
	return s.repo.GetByID(id)
}

// GetEmployeeByEmail retrieves an employee by email
func (s *employeeService) GetEmployeeByEmail(email string) (*model.Employee, error) {
	return s.repo.GetByEmail(email)
}

// GetAllEmployees retrieves all employees with pagination and optional filtering
func (s *employeeService) GetAllEmployees(limit, offset int, filter repository.EmployeeFilter) ([]*model.Employee, error) {
	return s.repo.GetAll(limit, offset, filter)
}

// UpdateEmployee updates an existing employee after validation
func (s *employeeService) UpdateEmployee(actorID uint, employee *model.Employee) error {
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
	if err := s.repo.Update(employee); err != nil {
		return err
	}

	// Audit log
	_ = s.audit.Log(actorID, "UPDATE", "employees", employee.ID, existingEmployee, employee, nil)

	return nil
}

// DeleteEmployee removes an employee by ID
func (s *employeeService) DeleteEmployee(actorID uint, id uint) error {
	// Check if employee exists
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("employee with ID %d does not exist", id)
	}

	// TODO: Check if employee has any active assignments or dependencies before deletion
	// For now, just delete the employee
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// Audit log
	_ = s.audit.Log(actorID, "DELETE", "employees", id, existing, nil, nil)

	return nil
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

// GetOrganizationTree builds the company hierarchy tree
func (s *employeeService) GetOrganizationTree() ([]*model.OrganizationNode, error) {
	// 1. Fetch all employees with preloaded Department and Position
	// Using a large limit and empty filter to get everyone for now
	employees, err := s.repo.GetAll(10000, 0, repository.EmployeeFilter{})
	if err != nil {
		return nil, err
	}

	// 2. Create a map of nodes for quick lookup
	nodeMap := make(map[uint]*model.OrganizationNode)
	var roots []*model.OrganizationNode

	for _, emp := range employees {
		node := &model.OrganizationNode{
			ID:        emp.ID,
			FirstName: emp.FirstName,
			LastName:  emp.LastName,
			Status:    emp.Status,
		}
		
		if emp.Position != nil {
			node.JobTitle = emp.Position.Title
		}
		if emp.Department != nil {
			node.Department = emp.Department.Name
		}
		
		nodeMap[emp.ID] = node
	}

	// 3. Link children to parents
	for _, emp := range employees {
		node := nodeMap[emp.ID]
		if emp.ManagerID == nil || *emp.ManagerID == 0 {
			roots = append(roots, node)
		} else {
			if parent, ok := nodeMap[*emp.ManagerID]; ok {
				parent.Subordinates = append(parent.Subordinates, node)
			} else {
				// If manager not found in current set (rare if we fetch all),
				// treat it as a root node for now.
				roots = append(roots, node)
			}
		}
	}

	return roots, nil
}

// GetByEmployeeNumber retrieves an employee by employee number
// This method is needed for validation but wasn't in the interface
// We'll add it to the repository and service