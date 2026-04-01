package model

import (
	"time"

	"gorm.io/gorm"
)

// Employee represents the core person record in the system
type Employee struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Personal Information
	FirstName        string         `gorm:"not null" json:"first_name"`
	LastName         string         `gorm:"not null" json:"last_name"`
	Email            string         `gorm:"uniqueIndex;not null" json:"email"`
	Phone            string         `json:"phone"`
	DateOfBirth      *time.Time     `json:"date_of_birth,omitempty"`
	Gender           string         `json:"gender"` // male, female, other

	// Address Information
	AddressLine1     string         `json:"address_line1"`
	AddressLine2     string         `json:"address_line2"`
	City             string         `json:"city"`
	State            string         `json:"state"`
	PostalCode       string         `json:"postal_code"`
	Country          string         `json:"country"`

	// Emergency Contact
	EmergencyContactName     string `json:"emergency_contact_name"`
	EmergencyContactPhone    string `json:"emergency_contact_phone"`
	EmergencyContactRelation string `json:"emergency_contact_relation"`

	// Employment Information
	EmployeeNumber   string         `gorm:"uniqueIndex;not null" json:"employee_number"`
	HireDate         *time.Time     `json:"hire_date,omitempty"`
	TerminationDate  *time.Time     `json:"termination_date,omitempty"`
	Status           string         `gorm:"default:active" json:"status"` // active, inactive, terminated
	DepartmentID     *uint          `json:"department_id,omitempty"`
	PositionID       *uint          `json:"position_id,omitempty"`
	ManagerID        *uint          `json:"manager_id,omitempty"` // Self-referencing foreign key

	// Bank Account Details
	BankAccountNumber string `json:"bank_account_number"`
	BankName         string `json:"bank_name"`
	BankRoutingNumber string `json:"bank_routing_number"`

	// Relationships
	Department *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Position   *Position   `gorm:"foreignKey:PositionID" json:"position,omitempty"`
	Manager    *Employee   `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
	Subordinates []Employee `gorm:"foreignKey:ManagerID" json:"subordinates,omitempty"`
}

// TableName specifies the table name for Employee
func (Employee) TableName() string {
	return "employees"
}