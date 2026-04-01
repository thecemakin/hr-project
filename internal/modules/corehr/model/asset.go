package model

import (
	"time"

	"gorm.io/gorm"
)

// Asset represents company-owned items such as laptops or phones
type Asset struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	SerialNumber    string `gorm:"uniqueIndex;not null" json:"serial_number"`
	Name            string `gorm:"not null" json:"name"` // e.g., "MacBook Pro 16-inch"
	Description     string `json:"description"`
	AssetTag        string `gorm:"uniqueIndex" json:"asset_tag"` // Internal tracking number
	Type            string `gorm:"not null" json:"type"`         // laptop, phone, tablet, etc.
	Brand           string `json:"brand"`
	Model           string `json:"model"`
	PurchaseDate    *time.Time `json:"purchase_date,omitempty"`
	PurchasePrice   float64    `json:"purchase_price"` // Purchase price in dollars
	Status          string     `gorm:"default:available" json:"status"` // available, assigned, maintenance, retired
	Condition       string     `json:"condition"`      // new, good, fair, poor
	WarrantyExpires *time.Time `json:"warranty_expires,omitempty"`

	// Relationships
	Assignments []AssetAssignment `gorm:"foreignKey:AssetID" json:"assignments,omitempty"`
}

// TableName specifies the table name for Asset
func (Asset) TableName() string {
	return "assets"
}