package repository

import (
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"gorm.io/gorm"
)

// AssetFilter defines available filters for asset list
type AssetFilter struct {
	SerialNumber string
	AssetTag     string
	Status       string
	Type         string
	Search       string // For general brand/model search
}

// AssetRepository defines the interface for asset data operations
type AssetRepository interface {
	Create(asset *model.Asset) error
	GetByID(id uint) (*model.Asset, error)
	GetAll(limit, offset int, filter AssetFilter) ([]*model.Asset, error)
	Update(asset *model.Asset) error
	Delete(id uint) error
}

// assetRepository implements the AssetRepository interface
type assetRepository struct {
	db *gorm.DB
}

// NewAssetRepository creates a new instance of assetRepository
func NewAssetRepository(db *gorm.DB) AssetRepository {
	return &assetRepository{
		db: db,
	}
}

// Create adds a new asset to the database
func (r *assetRepository) Create(asset *model.Asset) error {
	return r.db.Create(asset).Error
}

// GetByID retrieves an asset by ID
func (r *assetRepository) GetByID(id uint) (*model.Asset, error) {
	var asset model.Asset
	err := r.db.Preload("Assignments").Preload("Assignments.Employee").First(&asset, id).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}


// GetAll retrieves assets with pagination and optional filtering
func (r *assetRepository) GetAll(limit, offset int, filter AssetFilter) ([]*model.Asset, error) {
	var assets []*model.Asset
	query := r.db.Preload("Assignments").Preload("Assignments.Employee")

	if filter.SerialNumber != "" {
		query = query.Where("serial_number = ?", filter.SerialNumber)
	}
	if filter.AssetTag != "" {
		query = query.Where("asset_tag = ?", filter.AssetTag)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("brand ILIKE ? OR model ILIKE ?", searchTerm, searchTerm)
	}

	err := query.Limit(limit).Offset(offset).Find(&assets).Error
	if err != nil {
		return nil, err
	}
	return assets, nil
}

// Update modifies an existing asset
func (r *assetRepository) Update(asset *model.Asset) error {
	return r.db.Save(asset).Error
}

// Delete removes an asset by ID
func (r *assetRepository) Delete(id uint) error {
	return r.db.Delete(&model.Asset{}, id).Error
}

// GetByStatus retrieves all assets with a specific status
func (r *assetRepository) GetByStatus(status string) ([]*model.Asset, error) {
	var assets []*model.Asset
	err := r.db.Preload("Assignments").Preload("Assignments.Employee").Where("status = ?", status).Find(&assets).Error
	if err != nil {
		return nil, err
	}
	return assets, nil
}

// GetByType retrieves all assets of a specific type
func (r *assetRepository) GetByType(assetType string) ([]*model.Asset, error) {
	var assets []*model.Asset
	err := r.db.Preload("Assignments").Preload("Assignments.Employee").Where("type = ?", assetType).Find(&assets).Error
	if err != nil {
		return nil, err
	}
	return assets, nil
}