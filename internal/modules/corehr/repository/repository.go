package repository

import (
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"gorm.io/gorm"
)

// Repository defines the data access layer for Core HR module.
type Repository interface {
	CreateEmployee(emp *model.Employee) error
	GetEmployeeByID(id uint) (*model.Employee, error)
	ListEmployees() ([]model.Employee, error)
	UpdateEmployee(emp *model.Employee) error
}

// sqlRepo implements the Repository interface using GORM.
type sqlRepo struct {
	db *gorm.DB
}

// NewSQLRepository creates a new GORM-based repository.
func NewSQLRepository(db *gorm.DB) Repository {
	return &sqlRepo{db: db}
}

func (r *sqlRepo) CreateEmployee(emp *model.Employee) error {
	return r.db.Create(emp).Error
}

func (r *sqlRepo) GetEmployeeByID(id uint) (*model.Employee, error) {
	var emp model.Employee
	if err := r.db.Preload("Department").Preload("Position").Preload("Manager").First(&emp, id).Error; err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *sqlRepo) ListEmployees() ([]model.Employee, error) {
	var emps []model.Employee
	if err := r.db.Preload("Department").Preload("Position").Find(&emps).Error; err != nil {
		return nil, err
	}
	return emps, nil
}

func (r *sqlRepo) UpdateEmployee(emp *model.Employee) error {
	return r.db.Save(emp).Error
}
