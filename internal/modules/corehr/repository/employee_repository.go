package repository

import (
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"gorm.io/gorm"
)

// EmployeeRepository defines the interface for employee data operations
type EmployeeRepository interface {
	Create(employee *model.Employee) error
	GetByID(id uint) (*model.Employee, error)
	GetByEmail(email string) (*model.Employee, error)
	GetByEmployeeNumber(employeeNumber string) (*model.Employee, error)
	GetAll(limit, offset int) ([]*model.Employee, error)
	Update(employee *model.Employee) error
	Delete(id uint) error
	GetByManagerID(managerID uint) ([]*model.Employee, error)
	GetByDepartmentID(departmentID uint) ([]*model.Employee, error)
	GetByStatus(status string) ([]*model.Employee, error)
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

// GetAll retrieves all employees with pagination
func (r *employeeRepository) GetAll(limit, offset int) ([]*model.Employee, error) {
	var employees []*model.Employee
	err := r.db.Preload("Department").Preload("Position").Preload("Manager").Limit(limit).Offset(offset).Find(&employees).Error
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

// GetByManagerID retrieves all employees reporting to a specific manager
func (r *employeeRepository) GetByManagerID(managerID uint) ([]*model.Employee, error) {
	var employees []*model.Employee
	err := r.db.Preload("Department").Preload("Position").Preload("Manager").Where("manager_id = ?", managerID).Find(&employees).Error
	if err != nil {
		return nil, err
	}
	return employees, nil
}

// GetByDepartmentID retrieves all employees in a specific department
func (r *employeeRepository) GetByDepartmentID(departmentID uint) ([]*model.Employee, error) {
	var employees []*model.Employee
	err := r.db.Preload("Department").Preload("Position").Preload("Manager").Where("department_id = ?", departmentID).Find(&employees).Error
	if err != nil {
		return nil, err
	}
	return employees, nil
}

// GetByStatus retrieves all employees with a specific status
func (r *employeeRepository) GetByStatus(status string) ([]*model.Employee, error) {
	var employees []*model.Employee
	err := r.db.Preload("Department").Preload("Position").Preload("Manager").Where("status = ?", status).Find(&employees).Error
	if err != nil {
		return nil, err
	}
	return employees, nil
}