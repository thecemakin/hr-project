package model

import (
	"time"

	"gorm.io/gorm"
)

// AssetAssignment represents the relationship between an asset and an employee
type AssetAssignment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	AssetID      uint      `gorm:"not null" json:"asset_id"`
	EmployeeID   uint      `gorm:"not null" json:"employee_id"`
	AssignedBy   uint      `json:"assigned_by"`   // ID of the person who assigned the asset
	AssignedDate time.Time `gorm:"not null" json:"assigned_date"`
	ReturnedDate *time.Time `json:"returned_date,omitempty"`
	ReturnedBy   *uint      `json:"returned_by,omitempty"` // ID of the person who returned the asset
	Notes        string     `json:"notes"`

	// Status indicates current state of the assignment
	Status string `gorm:"default:assigned;not null" json:"status"` // assigned, returned

	// Relationships
	Asset    *Asset    `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
	Employee *Employee `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
	Assigner *Employee `gorm:"foreignKey:AssignedBy" json:"assigner,omitempty"`
	Returner *Employee `gorm:"foreignKey:ReturnedBy" json:"returner,omitempty"`
}

// TableName specifies the table name for AssetAssignment
func (AssetAssignment) TableName() string {
	return "asset_assignments"
}