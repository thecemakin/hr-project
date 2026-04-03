package repository

import (
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"gorm.io/gorm"
)

// EmployeeFilter defines the available filters for employee list
type EmployeeFilter struct {
	Email        string
	Status       string
	DepartmentID *uint
	ManagerID    *uint
	Search       string // For general name search
}

// EmployeeRepository defines the interface for employee data operations
type EmployeeRepository interface {
	Create(employee *model.Employee) error
	GetByID(id uint) (*model.Employee, error)
	GetByEmail(email string) (*model.Employee, error)
	GetByEmployeeNumber(employeeNumber string) (*model.Employee, error)
	GetAll(limit, offset int, filter EmployeeFilter) ([]*model.Employee, error)
	Update(employee *model.Employee) error
	Delete(id uint) error
}

// employeeRepository implements the EmployeeRepository interface
type employeeRepository struct {
	db *gorm.DB
}

// NewEmployeeRepository creates a new instance of employeeRepository
func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{
		db: db,
	}
}

// Create adds a new employee to the database
func (r *employeeRepository) Create(employee *model.Employee) error {
	return r.db.Create(employee).Error
}

// GetByID retrieves an employee by ID
func (r *employeeRepository) GetByID(id uint) (*model.Employee, error) {
	var employee model.Employee
	err := r.db.Preload("Department").Preload("Position").Preload("Manager").First(&employee, id).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// GetByEmail retrieves an employee by email
func (r *employeeRepository) GetByEmail(email string) (*model.Employee, error) {
	var employee model.Employee
	err := r.db.Preload("Department").Preload("Position").Preload("Manager").Where("email = ?", email).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// GetByEmployeeNumber retrieves an employee by employee number
func (r *employeeRepository) GetByEmployeeNumber(employeeNumber string) (*model.Employee, error) {
	var employee model.Employee
	err := r.db.Preload("Department").Preload("Position").Preload("Manager").Where("employee_number = ?", employeeNumber).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// GetAll retrieves employees with pagination and optional filtering
func (r *employeeRepository) GetAll(limit, offset int, filter EmployeeFilter) ([]*model.Employee, error) {
	var employees []*model.Employee
	query := r.db.Preload("Department").Preload("Position").Preload("Manager")

	if filter.Email != "" {
		query = query.Where("email = ?", filter.Email)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.DepartmentID != nil {
		query = query.Where("department_id = ?", *filter.DepartmentID)
	}
	if filter.ManagerID != nil {
		query = query.Where("manager_id = ?", *filter.ManagerID)
	}
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ?", searchTerm, searchTerm)
	}

	err := query.Limit(limit).Offset(offset).Find(&employees).Error
	if err != nil {
		return nil, err
	}
	return employees, nil
}

// Update modifies an existing employee
func (r *employeeRepository) Update(employee *model.Employee) error {
	return r.db.Save(employee).Error
}

// Delete removes an employee by ID
func (r *employeeRepository) Delete(id uint) error {
	return r.db.Delete(&model.Employee{}, id).Error
}