package repository

import (
	"github.com/thecemakin/hr-project/internal/modules/leave/model"
	"gorm.io/gorm"
)

// Repository defines the data access layer for Leave module.
type Repository interface {
	// Leave Types
	CreateLeaveType(lt *model.LeaveType) error
	GetLeaveTypeByID(id uint) (*model.LeaveType, error)
	ListLeaveTypes() ([]model.LeaveType, error)
	UpdateLeaveType(lt *model.LeaveType) error

	// Leave Balances
	CreateLeaveBalance(lb *model.LeaveBalance) error
	GetLeaveBalance(employeeID, leaveTypeID uint) (*model.LeaveBalance, error)
	ListLeaveBalancesByEmployee(employeeID uint) ([]model.LeaveBalance, error)
	UpdateLeaveBalance(lb *model.LeaveBalance) error

	// Leave Requests
	CreateLeaveRequest(req *model.LeaveRequest) error
	GetLeaveRequestByID(id uint) (*model.LeaveRequest, error)
	ListLeaveRequestsByEmployee(employeeID uint) ([]model.LeaveRequest, error)
	ListPendingLeaveRequestsForManager(managerID uint) ([]model.LeaveRequest, error)
	UpdateLeaveRequest(req *model.LeaveRequest) error
}

type sqlRepo struct {
	db *gorm.DB
}

// NewSQLRepository creates a new GORM-based repository.
func NewSQLRepository(db *gorm.DB) Repository {
	return &sqlRepo{db: db}
}

// --- Leave Types -------------------------------------------------------------

func (r *sqlRepo) CreateLeaveType(lt *model.LeaveType) error {
	return r.db.Create(lt).Error
}

func (r *sqlRepo) GetLeaveTypeByID(id uint) (*model.LeaveType, error) {
	var lt model.LeaveType
	if err := r.db.First(&lt, id).Error; err != nil {
		return nil, err
	}
	return &lt, nil
}

func (r *sqlRepo) ListLeaveTypes() ([]model.LeaveType, error) {
	var lts []model.LeaveType
	if err := r.db.Find(&lts).Error; err != nil {
		return nil, err
	}
	return lts, nil
}

func (r *sqlRepo) UpdateLeaveType(lt *model.LeaveType) error {
	return r.db.Save(lt).Error
}

// --- Leave Balances ----------------------------------------------------------

func (r *sqlRepo) CreateLeaveBalance(lb *model.LeaveBalance) error {
	return r.db.Create(lb).Error
}

func (r *sqlRepo) GetLeaveBalance(employeeID, leaveTypeID uint) (*model.LeaveBalance, error) {
	var lb model.LeaveBalance
	if err := r.db.Where("employee_id = ? AND leave_type_id = ?", employeeID, leaveTypeID).First(&lb).Error; err != nil {
		return nil, err
	}
	return &lb, nil
}

func (r *sqlRepo) ListLeaveBalancesByEmployee(employeeID uint) ([]model.LeaveBalance, error) {
	var lbs []model.LeaveBalance
	if err := r.db.Preload("LeaveType").Where("employee_id = ?", employeeID).Find(&lbs).Error; err != nil {
		return nil, err
	}
	return lbs, nil
}

func (r *sqlRepo) UpdateLeaveBalance(lb *model.LeaveBalance) error {
	return r.db.Save(lb).Error
}

// --- Leave Requests ----------------------------------------------------------

func (r *sqlRepo) CreateLeaveRequest(req *model.LeaveRequest) error {
	return r.db.Create(req).Error
}

func (r *sqlRepo) GetLeaveRequestByID(id uint) (*model.LeaveRequest, error) {
	var req model.LeaveRequest
	if err := r.db.Preload("Employee").Preload("LeaveType").Preload("Reviewer").First(&req, id).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *sqlRepo) ListLeaveRequestsByEmployee(employeeID uint) ([]model.LeaveRequest, error) {
	var reqs []model.LeaveRequest
	if err := r.db.Preload("LeaveType").Where("employee_id = ?", employeeID).Order("created_at desc").Find(&reqs).Error; err != nil {
		return nil, err
	}
	return reqs, nil
}

func (r *sqlRepo) ListPendingLeaveRequestsForManager(managerID uint) ([]model.LeaveRequest, error) {
	var reqs []model.LeaveRequest
	// Joins employees table to find requests where the employee's manager is managerID AND status is pending
	err := r.db.Joins("JOIN employees ON employees.id = leave_requests.employee_id").
		Where("employees.manager_id = ? AND leave_requests.status = ?", managerID, "pending").
		Preload("Employee").Preload("LeaveType").
		Order("leave_requests.created_at asc").
		Find(&reqs).Error

	if err != nil {
		return nil, err
	}
	return reqs, nil
}

func (r *sqlRepo) UpdateLeaveRequest(req *model.LeaveRequest) error {
	return r.db.Save(req).Error
}
