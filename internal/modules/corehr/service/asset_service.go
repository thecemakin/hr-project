package service

import (
	"errors"
	"fmt"

	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"github.com/thecemakin/hr-project/internal/modules/corehr/repository"
)

// AssetService defines the interface for asset business operations
type AssetService interface {
	CreateAsset(asset *model.Asset) error
	GetAssetByID(id uint) (*model.Asset, error)
	GetAllAssets(limit, offset int, filter repository.AssetFilter) ([]*model.Asset, error)
	UpdateAsset(asset *model.Asset) error
	DeleteAsset(id uint) error
}

// assetService implements the AssetService interface
type assetService struct {
	repo repository.AssetRepository
}

// NewAssetService creates a new instance of assetService
func NewAssetService(repo repository.AssetRepository) AssetService {
	return &assetService{
		repo: repo,
	}
}

// CreateAsset creates a new asset after validation
func (s *assetService) CreateAsset(asset *model.Asset) error {
	// Validate required fields
	if asset.SerialNumber == "" || asset.Name == "" || asset.Type == "" {
		return errors.New("serial number, name, and type are required")
	}

	// Check if serial number already exists
	snFilter := repository.AssetFilter{SerialNumber: asset.SerialNumber}
	existingSn, _ := s.repo.GetAll(1, 0, snFilter)
	if len(existingSn) > 0 {
		return fmt.Errorf("asset with serial number %s already exists", asset.SerialNumber)
	}

	// Check if asset tag is provided and if it already exists
	if asset.AssetTag != "" {
		atFilter := repository.AssetFilter{AssetTag: asset.AssetTag}
		existingTag, _ := s.repo.GetAll(1, 0, atFilter)
		if len(existingTag) > 0 {
			return fmt.Errorf("asset with tag %s already exists", asset.AssetTag)
		}
	}

	// Create the asset
	return s.repo.Create(asset)
}

// GetAssetByID retrieves an asset by ID
func (s *assetService) GetAssetByID(id uint) (*model.Asset, error) {
	return s.repo.GetByID(id)
}


// GetAllAssets retrieves all assets with pagination and optional filtering
func (s *assetService) GetAllAssets(limit, offset int, filter repository.AssetFilter) ([]*model.Asset, error) {
	return s.repo.GetAll(limit, offset, filter)
}

// UpdateAsset updates an existing asset after validation
func (s *assetService) UpdateAsset(asset *model.Asset) error {
	// Validate required fields
	if asset.SerialNumber == "" || asset.Name == "" || asset.Type == "" {
		return errors.New("serial number, name, and type are required")
	}

	// Check if asset exists
	existingAsset, err := s.repo.GetByID(asset.ID)
	if err != nil {
		return fmt.Errorf("asset with ID %d does not exist", asset.ID)
	}

	// Check if serial number is being changed and if it already exists for another asset
	if existingAsset.SerialNumber != asset.SerialNumber {
		snFilter := repository.AssetFilter{SerialNumber: asset.SerialNumber}
		existingSn, _ := s.repo.GetAll(1, 0, snFilter)
		if len(existingSn) > 0 && existingSn[0].ID != asset.ID {
			return fmt.Errorf("asset with serial number %s already exists", asset.SerialNumber)
		}
	}

	// Check if asset tag is being changed and if it already exists for another asset
	if existingAsset.AssetTag != asset.AssetTag && asset.AssetTag != "" {
		atFilter := repository.AssetFilter{AssetTag: asset.AssetTag}
		existingTag, _ := s.repo.GetAll(1, 0, atFilter)
		if len(existingTag) > 0 && existingTag[0].ID != asset.ID {
			return fmt.Errorf("asset with tag %s already exists", asset.AssetTag)
		}
	}

	// Update the asset
	return s.repo.Update(asset)
}

// DeleteAsset removes an asset by ID
func (s *assetService) DeleteAsset(id uint) error {
	// Check if asset exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("asset with ID %d does not exist", id)
	}

	// TODO: Check if asset has any active assignments before deletion
	// For now, just delete the asset
	return s.repo.Delete(id)
}