package model

import (
	"time"

	corehrmodel "github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"gorm.io/gorm"
)

// LeaveBalance tracks how many days an employee has and used
type LeaveBalance struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	EmployeeID  uint           `gorm:"uniqueIndex:idx_emp_leave_type;not null" json:"employee_id"`
	LeaveTypeID uint           `gorm:"uniqueIndex:idx_emp_leave_type;not null" json:"leave_type_id"`
	TotalDays   int            `gorm:"default:0;not null" json:"total_days"`
	UsedDays    int            `gorm:"default:0;not null" json:"used_days"`

	// Relationships
	Employee  *corehrmodel.Employee `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
	LeaveType *LeaveType            `gorm:"foreignKey:LeaveTypeID" json:"leave_type,omitempty"`
}

func (LeaveBalance) TableName() string {
	return "leave_balances"
}
