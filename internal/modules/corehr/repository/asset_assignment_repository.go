package repository

import (
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"gorm.io/gorm"
)

// AssetAssignmentFilter defines available filters for asset assignments
type AssetAssignmentFilter struct {
	AssetID    uint
	EmployeeID uint
	Status     string // e.g., 'assigned', 'returned'
}

// AssetAssignmentRepository defines the interface for asset assignment data operations
type AssetAssignmentRepository interface {
	Create(assignment *model.AssetAssignment) error
	GetByID(id uint) (*model.AssetAssignment, error)
	GetAll(limit, offset int, filter AssetAssignmentFilter) ([]*model.AssetAssignment, error)
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


// GetAll retrieves asset assignments with pagination and optional filtering
func (r *assetAssignmentRepository) GetAll(limit, offset int, filter AssetAssignmentFilter) ([]*model.AssetAssignment, error) {
	var assignments []*model.AssetAssignment
	query := r.db.Preload("Asset").Preload("Employee").Preload("Assigner").Preload("Returner")

	if filter.AssetID != 0 {
		query = query.Where("asset_id = ?", filter.AssetID)
	}
	if filter.EmployeeID != 0 {
		query = query.Where("employee_id = ?", filter.EmployeeID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	err := query.Limit(limit).Offset(offset).Find(&assignments).Error
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