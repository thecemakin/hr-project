package model

import (
	"time"

	"gorm.io/gorm"
)

// LeaveType represents a kind of leave an employee can take
type LeaveType struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	DefaultDays int            `gorm:"default:0" json:"default_days"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
}

func (LeaveType) TableName() string {
	return "leave_types"
}
