package repository

import (
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"gorm.io/gorm"
)

// DepartmentRepository defines the interface for department data operations
type DepartmentRepository interface {
	Create(department *model.Department) error
	GetByID(id uint) (*model.Department, error)
	GetByName(name string) (*model.Department, error)
	GetAll(limit, offset int) ([]*model.Department, error)
	Update(department *model.Department) error
	Delete(id uint) error
}

// departmentRepository implements the DepartmentRepository interface
type departmentRepository struct {
	db *gorm.DB
}

// NewDepartmentRepository creates a new instance of departmentRepository
func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{
		db: db,
	}
}

// Create adds a new department to the database
func (r *departmentRepository) Create(department *model.Department) error {
	return r.db.Create(department).Error
}

// GetByID retrieves a department by ID
func (r *departmentRepository) GetByID(id uint) (*model.Department, error) {
	var department model.Department
	err := r.db.Preload("Head").First(&department, id).Error
	if err != nil {
		return nil, err
	}
	return &department, nil
}

// GetByName retrieves a department by name
func (r *departmentRepository) GetByName(name string) (*model.Department, error) {
	var department model.Department
	err := r.db.Preload("Head").Where("name = ?", name).First(&department).Error
	if err != nil {
		return nil, err
	}
	return &department, nil
}

// GetAll retrieves all departments with pagination
func (r *departmentRepository) GetAll(limit, offset int) ([]*model.Department, error) {
	var departments []*model.Department
	err := r.db.Preload("Head").Limit(limit).Offset(offset).Find(&departments).Error
	if err != nil {
		return nil, err
	}
	return departments, nil
}

// Update modifies an existing department
func (r *departmentRepository) Update(department *model.Department) error {
	return r.db.Save(department).Error
}

// Delete removes a department by ID
func (r *departmentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Department{}, id).Error
}