package repository

import (
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"gorm.io/gorm"
)

// PositionFilter defines the available filters for position list
type PositionFilter struct {
	Title string
}

// PositionRepository defines the interface for position data operations
type PositionRepository interface {
	Create(position *model.Position) error
	GetByID(id uint) (*model.Position, error)
	GetAll(limit, offset int, filter PositionFilter) ([]*model.Position, error)
	Update(position *model.Position) error
	Delete(id uint) error
}

// positionRepository implements the PositionRepository interface
type positionRepository struct {
	db *gorm.DB
}

// NewPositionRepository creates a new instance of positionRepository
func NewPositionRepository(db *gorm.DB) PositionRepository {
	return &positionRepository{
		db: db,
	}
}

// Create adds a new position to the database
func (r *positionRepository) Create(position *model.Position) error {
	return r.db.Create(position).Error
}

// GetByID retrieves a position by ID
func (r *positionRepository) GetByID(id uint) (*model.Position, error) {
	var position model.Position
	err := r.db.Preload("Department").First(&position, id).Error
	if err != nil {
		return nil, err
	}
	return &position, nil
}


// GetAll retrieves positions with pagination and optional filtering
func (r *positionRepository) GetAll(limit, offset int, filter PositionFilter) ([]*model.Position, error) {
	var positions []*model.Position
	query := r.db.Preload("Department")

	if filter.Title != "" {
		query = query.Where("title ILIKE ?", "%"+filter.Title+"%")
	}

	err := query.Limit(limit).Offset(offset).Find(&positions).Error
	if err != nil {
		return nil, err
	}
	return positions, nil
}

// Update modifies an existing position
func (r *positionRepository) Update(position *model.Position) error {
	return r.db.Save(position).Error
}

// Delete removes a position by ID
func (r *positionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Position{}, id).Error
}