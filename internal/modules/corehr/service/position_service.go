package service

import (
	"errors"
	"fmt"

	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"github.com/thecemakin/hr-project/internal/modules/corehr/repository"
)

// PositionService defines the interface for position business operations
type PositionService interface {
	CreatePosition(position *model.Position) error
	GetPositionByID(id uint) (*model.Position, error)
	GetAllPositions(limit, offset int, filter repository.PositionFilter) ([]*model.Position, error)
	UpdatePosition(position *model.Position) error
	DeletePosition(id uint) error
}

// positionService implements the PositionService interface
type positionService struct {
	repo repository.PositionRepository
}

// NewPositionService creates a new instance of positionService
func NewPositionService(repo repository.PositionRepository) PositionService {
	return &positionService{
		repo: repo,
	}
}

// CreatePosition creates a new position after validation
func (s *positionService) CreatePosition(position *model.Position) error {
	// Validate required fields
	if position.Title == "" {
		return errors.New("position title is required")
	}

	// Check if position title already exists
	titleFilter := repository.PositionFilter{Title: position.Title}
	existingWithTitle, _ := s.repo.GetAll(1, 0, titleFilter)
	if len(existingWithTitle) > 0 {
		return fmt.Errorf("position with title %s already exists", position.Title)
	}

	// Validate department if provided
	if position.DepartmentID != nil {
		// In a real implementation, we would validate that the department exists
		// For now, we'll just proceed with creation
	}

	// Create the position
	return s.repo.Create(position)
}

// GetPositionByID retrieves a position by ID
func (s *positionService) GetPositionByID(id uint) (*model.Position, error) {
	return s.repo.GetByID(id)
}

// GetAllPositions retrieves all positions with pagination and optional filtering
func (s *positionService) GetAllPositions(limit, offset int, filter repository.PositionFilter) ([]*model.Position, error) {
	return s.repo.GetAll(limit, offset, filter)
}

// UpdatePosition updates an existing position after validation
func (s *positionService) UpdatePosition(position *model.Position) error {
	// Validate required fields
	if position.Title == "" {
		return errors.New("position title is required")
	}

	// Check if position exists
	existingPos, err := s.repo.GetByID(position.ID)
	if err != nil {
		return fmt.Errorf("position with ID %d does not exist", position.ID)
	}

	// Check if title is being changed and if it already exists for another position
	if existingPos.Title != position.Title {
		titleFilter := repository.PositionFilter{Title: position.Title}
		existingWithTitle, _ := s.repo.GetAll(1, 0, titleFilter)
		if len(existingWithTitle) > 0 && existingWithTitle[0].ID != position.ID {
			return fmt.Errorf("position with title %s already exists", position.Title)
		}
	}

	// Update the position
	return s.repo.Update(position)
}

// DeletePosition removes a position by ID
func (s *positionService) DeletePosition(id uint) error {
	// Check if position exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("position with ID %d does not exist", id)
	}

	// TODO: Check if position has any employees before deletion
	// For now, just delete the position
	return s.repo.Delete(id)
}