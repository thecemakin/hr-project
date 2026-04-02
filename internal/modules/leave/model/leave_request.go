package model

import (
	"time"

	corehrmodel "github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"gorm.io/gorm"
)

// LeaveRequest represents a formal request for leave by an employee
type LeaveRequest struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	EmployeeID         uint           `gorm:"not null" json:"employee_id"`
	LeaveTypeID        uint           `gorm:"not null" json:"leave_type_id"`
	StartDate          time.Time      `gorm:"type:date;not null" json:"start_date"`
	EndDate            time.Time      `gorm:"type:date;not null" json:"end_date"`
	TotalDaysRequested int            `gorm:"not null" json:"total_days_requested"`
	
	Status             string         `gorm:"default:pending;not null" json:"status"` // pending, approved, rejected
	Reason             string         `json:"reason"`
	
	ReviewerID         *uint          `json:"reviewer_id,omitempty"`
	ReviewerNote       string         `json:"reviewer_note"`

	// Relationships
	Employee  *corehrmodel.Employee `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
	LeaveType *LeaveType            `gorm:"foreignKey:LeaveTypeID" json:"leave_type,omitempty"`
	Reviewer  *corehrmodel.Employee `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
}

func (LeaveRequest) TableName() string {
	return "leave_requests"
}
