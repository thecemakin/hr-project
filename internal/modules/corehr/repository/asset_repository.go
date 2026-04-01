package repository

import (
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"gorm.io/gorm"
)

// AssetRepository defines the interface for asset data operations
type AssetRepository interface {
	Create(asset *model.Asset) error
	GetByID(id uint) (*model.Asset, error)
	GetBySerialNumber(serialNumber string) (*model.Asset, error)
	GetByAssetTag(assetTag string) (*model.Asset, error)
	GetAll(limit, offset int) ([]*model.Asset, error)
	Update(asset *model.Asset) error
	Delete(id uint) error
	GetByStatus(status string) ([]*model.Asset, error)
	GetByType(assetType string) ([]*model.Asset, error)
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

// GetBySerialNumber retrieves an asset by serial number
func (r *assetRepository) GetBySerialNumber(serialNumber string) (*model.Asset, error) {
	var asset model.Asset
	err := r.db.Preload("Assignments").Preload("Assignments.Employee").Where("serial_number = ?", serialNumber).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// GetByAssetTag retrieves an asset by asset tag
func (r *assetRepository) GetByAssetTag(assetTag string) (*model.Asset, error) {
	var asset model.Asset
	err := r.db.Preload("Assignments").Preload("Assignments.Employee").Where("asset_tag = ?", assetTag).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// GetAll retrieves all assets with pagination
func (r *assetRepository) GetAll(limit, offset int) ([]*model.Asset, error) {
	var assets []*model.Asset
	err := r.db.Preload("Assignments").Preload("Assignments.Employee").Limit(limit).Offset(offset).Find(&assets).Error
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