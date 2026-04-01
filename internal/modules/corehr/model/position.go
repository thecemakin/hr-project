package model

import (
	"time"

	"gorm.io/gorm"
)

// Position represents job titles or roles within the organization
type Position struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Title       string `gorm:"not null" json:"title"`
	Description string `json:"description"`
	Code        string `gorm:"uniqueIndex" json:"code"` // e.g., "SENG", "MGR", "VP"
	Level       int    `json:"level"`                   // Hierarchical level (1-10 scale)
	SalaryMin   int64  `json:"salary_min"`              // Minimum salary in cents
	SalaryMax   int64  `json:"salary_max"`              // Maximum salary in cents
	IsActive    bool   `gorm:"default:true" json:"is_active"`

	// Relationships
	DepartmentID *uint       `json:"department_id,omitempty"`
	Department   *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
}

// TableName specifies the table name for Position
func (Position) TableName() string {
	return "positions"
}