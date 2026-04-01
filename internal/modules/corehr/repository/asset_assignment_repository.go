package repository

import (
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"gorm.io/gorm"
)

// AssetAssignmentRepository defines the interface for asset assignment data operations
type AssetAssignmentRepository interface {
	Create(assignment *model.AssetAssignment) error
	GetByID(id uint) (*model.AssetAssignment, error)
	GetByAssetID(assetID uint) ([]*model.AssetAssignment, error)
	GetByEmployeeID(employeeID uint) ([]*model.AssetAssignment, error)
	GetCurrentAssignmentsByEmployeeID(employeeID uint) ([]*model.AssetAssignment, error)
	GetActiveAssignmentsByAssetID(assetID uint) (*model.AssetAssignment, error)
	GetAll(limit, offset int) ([]*model.AssetAssignment, error)
	Update(assignment *model.AssetAssignment) error
	Delete(id uint) error
}

// assetAssignmentRepository implements the AssetAssignmentRepository interface
type assetAssignmentRepository struct {
	db *gorm.DB
}

// NewAssetAssignmentRepository creates a new instance of assetAssignmentRepository
func NewAssetAssignmentRepository(db *gorm.DB) AssetAssignmentRepository {
	return &assetAssignmentRepository{
		db: db,
	}
}

// Create adds a new asset assignment to the database
func (r *assetAssignmentRepository) Create(assignment *model.AssetAssignment) error {
	return r.db.Create(assignment).Error
}

// GetByID retrieves an asset assignment by ID
func (r *assetAssignmentRepository) GetByID(id uint) (*model.AssetAssignment, error) {
	var assignment model.AssetAssignment
	err := r.db.Preload("Asset").Preload("Employee").Preload("Assigner").Preload("Returner").First(&assignment, id).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

// GetByAssetID retrieves all assignments for a specific asset
func (r *assetAssignmentRepository) GetByAssetID(assetID uint) ([]*model.AssetAssignment, error) {
	var assignments []*model.AssetAssignment
	err := r.db.Preload("Asset").Preload("Employee").Preload("Assigner").Preload("Returner").Where("asset_id = ?", assetID).Find(&assignments).Error
	if err != nil {
		return nil, err
	}
	return assignments, nil
}

// GetByEmployeeID retrieves all assignments for a specific employee
func (r *assetAssignmentRepository) GetByEmployeeID(employeeID uint) ([]*model.AssetAssignment, error) {
	var assignments []*model.AssetAssignment
	err := r.db.Preload("Asset").Preload("Employee").Preload("Assigner").Preload("Returner").Where("employee_id = ?", employeeID).Find(&assignments).Error
	if err != nil {
		return nil, err
	}
	return assignments, nil
}

// GetCurrentAssignmentsByEmployeeID retrieves all current (non-returned) assignments for a specific employee
func (r *assetAssignmentRepository) GetCurrentAssignmentsByEmployeeID(employeeID uint) ([]*model.AssetAssignment, error) {
	var assignments []*model.AssetAssignment
	err := r.db.Preload("Asset").Preload("Employee").Preload("Assigner").Preload("Returner").Where("employee_id = ? AND status = 'assigned'", employeeID).Find(&assignments).Error
	if err != nil {
		return nil, err
	}
	return assignments, nil
}

// GetActiveAssignmentsByAssetID retrieves the currently active assignment for a specific asset
func (r *assetAssignmentRepository) GetActiveAssignmentsByAssetID(assetID uint) (*model.AssetAssignment, error) {
	var assignment model.AssetAssignment
	err := r.db.Preload("Asset").Preload("Employee").Preload("Assigner").Preload("Returner").Where("asset_id = ? AND status = 'assigned'", assetID).First(&assignment).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

// GetAll retrieves all asset assignments with pagination
func (r *assetAssignmentRepository) GetAll(limit, offset int) ([]*model.AssetAssignment, error) {
	var assignments []*model.AssetAssignment
	err := r.db.Preload("Asset").Preload("Employee").Preload("Assigner").Preload("Returner").Limit(limit).Offset(offset).Find(&assignments).Error
	if err != nil {
		return nil, err
	}
	return assignments, nil
}

// Update modifies an existing asset assignment
func (r *assetAssignmentRepository) Update(assignment *model.AssetAssignment) error {
	return r.db.Save(assignment).Error
}

// Delete removes an asset assignment by ID
func (r *assetAssignmentRepository) Delete(id uint) error {
	return r.db.Delete(&model.AssetAssignment{}, id).Error
}