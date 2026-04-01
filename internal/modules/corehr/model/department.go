package model

import (
	"time"

	"gorm.io/gorm"
)

// Department represents organizational groupings
type Department struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Name        string `gorm:"not null;uniqueIndex" json:"name"`
	Description string `json:"description"`
	Code        string `gorm:"uniqueIndex" json:"code"` // e.g., "ENG", "MKT", "FIN"
	HeadID      *uint  `json:"head_id,omitempty"`       // ID of the department head (Employee)

	// Relationships
	Head *Employee `gorm:"foreignKey:HeadID" json:"head,omitempty"`
}

// TableName specifies the table name for Department
func (Department) TableName() string {
	return "departments"
}