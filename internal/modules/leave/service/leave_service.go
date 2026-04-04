package service

import (
	"errors"
	"time"

	corehrrepo "github.com/thecemakin/hr-project/internal/modules/corehr/repository"
	"github.com/thecemakin/hr-project/internal/modules/leave/model"
	"github.com/thecemakin/hr-project/internal/modules/leave/repository"
	"github.com/thecemakin/hr-project/internal/platform/audit"
	"github.com/thecemakin/hr-project/internal/platform/notification"
)

// LeaveService defines business logic for Leave module
type LeaveService interface {
	// Leave Types
	CreateLeaveType(lt *model.LeaveType) error
	GetLeaveTypeByID(id uint) (*model.LeaveType, error)
	ListLeaveTypes() ([]model.LeaveType, error)

	// Leave Balances
	GetLeaveBalance(employeeID, leaveTypeID uint) (*model.LeaveBalance, error)
	ListLeaveBalancesByEmployee(employeeID uint) ([]model.LeaveBalance, error)
	InitializeLeaveBalance(employeeID, leaveTypeID uint, startingDays int) error

	// Leave Requests
	SubmitLeaveRequest(req *model.LeaveRequest) error
	GetLeaveRequestByID(id uint) (*model.LeaveRequest, error)
	ListOwnLeaveRequests(employeeID uint) ([]model.LeaveRequest, error)
	ListPendingApprovals(managerID uint) ([]model.LeaveRequest, error)
	ApproveLeaveRequest(requestID, managerID uint, note string) error
	RejectLeaveRequest(requestID, managerID uint, note string) error
}

type leaveService struct {
	repo       repository.Repository
	corehrRepo corehrrepo.EmployeeRepository
	audit      audit.Service
	notifier   notification.Notifier
}

// NewLeaveService creates a new LeaveService instance
func NewLeaveService(repo repository.Repository, corehrRepo corehrrepo.EmployeeRepository, audit audit.Service, notifier notification.Notifier) LeaveService {
	return &leaveService{
		repo:       repo,
		corehrRepo: corehrRepo,
		audit:      audit,
		notifier:   notifier,
	}
}

// --- Leave Types ---

func (s *leaveService) CreateLeaveType(lt *model.LeaveType) error {
	return s.repo.CreateLeaveType(lt)
}

func (s *leaveService) GetLeaveTypeByID(id uint) (*model.LeaveType, error) {
	return s.repo.GetLeaveTypeByID(id)
}

func (s *leaveService) ListLeaveTypes() ([]model.LeaveType, error) {
	return s.repo.ListLeaveTypes()
}

// --- Leave Balances ---

func (s *leaveService) GetLeaveBalance(employeeID, leaveTypeID uint) (*model.LeaveBalance, error) {
	return s.repo.GetLeaveBalance(employeeID, leaveTypeID)
}

func (s *leaveService) ListLeaveBalancesByEmployee(employeeID uint) ([]model.LeaveBalance, error) {
	return s.repo.ListLeaveBalancesByEmployee(employeeID)
}

func (s *leaveService) InitializeLeaveBalance(employeeID, leaveTypeID uint, startingDays int) error {
	balance := &model.LeaveBalance{
		EmployeeID:  employeeID,
		LeaveTypeID: leaveTypeID,
		TotalDays:   startingDays,
		UsedDays:    0,
	}
	return s.repo.CreateLeaveBalance(balance)
}

// --- Leave Requests ---

func (s *leaveService) SubmitLeaveRequest(req *model.LeaveRequest) error {
	// Basic validation
	if req.StartDate.After(req.EndDate) {
		return errors.New("start date must be before or equal to end date")
	}

	// Calculate requested days
	// Note: In a real system, you'd exclude weekends/holidays depending on policy.
	// We'll trust the requested days for now or implement a basic calculation here.
	if req.TotalDaysRequested <= 0 {
		return errors.New("total days requested must be greater than zero")
	}

	// Check leave balance
	balance, err := s.repo.GetLeaveBalance(req.EmployeeID, req.LeaveTypeID)
	if err != nil {
		return errors.New("could not find leave balance for this type")
	}

	availableDays := balance.TotalDays - balance.UsedDays
	if req.TotalDaysRequested > availableDays {
		return errors.New("insufficient leave balance")
	}

	req.Status = "pending"
	if err := s.repo.CreateLeaveRequest(req); err != nil {
		return err
	}

	// Notify (Simulation)
	emp, _ := s.corehrRepo.GetByID(req.EmployeeID)
	if emp != nil {
		_ = s.notifier.Send(emp.Email, "Leave Request Submitted", "Your leave request has been submitted and is pending approval.")
	}

	// Audit log
	_ = s.audit.Log(req.EmployeeID, "CREATE", "leave_requests", req.ID, nil, req, nil)

	return nil
}

func (s *leaveService) GetLeaveRequestByID(id uint) (*model.LeaveRequest, error) {
	return s.repo.GetLeaveRequestByID(id)
}

func (s *leaveService) ListOwnLeaveRequests(employeeID uint) ([]model.LeaveRequest, error) {
	return s.repo.ListLeaveRequestsByEmployee(employeeID)
}

func (s *leaveService) ListPendingApprovals(managerID uint) ([]model.LeaveRequest, error) {
	return s.repo.ListPendingLeaveRequestsForManager(managerID)
}

func (s *leaveService) ApproveLeaveRequest(requestID, managerID uint, note string) error {
	req, err := s.repo.GetLeaveRequestByID(requestID)
	if err != nil {
		return err
	}

	if req.Status != "pending" {
		return errors.New("only pending requests can be approved")
	}

	// Verify manager is indeed the employee's manager
	// Our List queries handle this, but for deep linking/direct API hits, we should verify
	emp, err := s.corehrRepo.GetByID(req.EmployeeID)
	if err != nil {
		return errors.New("employee not found")
	}

	if emp.ManagerID == nil || *emp.ManagerID != managerID {
		return errors.New("unauthorized: you are not the manager of this employee")
	}

	// Deduct balance
	balance, err := s.repo.GetLeaveBalance(req.EmployeeID, req.LeaveTypeID)
	if err != nil {
		return errors.New("could not find leave balance")
	}

	availableDays := balance.TotalDays - balance.UsedDays
	if req.TotalDaysRequested > availableDays {
		// Just in case things changed while it was pending
		return errors.New("insufficient leave balance to approve")
	}

	balance.UsedDays += req.TotalDaysRequested
	if err := s.repo.UpdateLeaveBalance(balance); err != nil {
		return err
	}

	// Mark as approved
	req.Status = "approved"
	req.ReviewerID = &managerID
	req.ReviewerNote = note
	req.UpdatedAt = time.Now()

	if err := s.repo.UpdateLeaveRequest(req); err != nil {
		return err
	}

	// Notify
	_ = s.notifier.Send(emp.Email, "Leave Request Approved", "Your leave request has been approved.")

	// Audit log
	_ = s.audit.Log(managerID, "APPROVE", "leave_requests", req.ID, "pending", "approved", map[string]string{"note": note})

	return nil
}

func (s *leaveService) RejectLeaveRequest(requestID, managerID uint, note string) error {
	req, err := s.repo.GetLeaveRequestByID(requestID)
	if err != nil {
		return err
	}

	if req.Status != "pending" {
		return errors.New("only pending requests can be rejected")
	}

	// Verify manager
	emp, err := s.corehrRepo.GetByID(req.EmployeeID)
	if err != nil {
		return errors.New("employee not found")
	}

	if emp.ManagerID == nil || *emp.ManagerID != managerID {
		return errors.New("unauthorized: you are not the manager of this employee")
	}

	// Mark as rejected
	req.Status = "rejected"
	req.ReviewerID = &managerID
	req.ReviewerNote = note
	req.UpdatedAt = time.Now()

	if err := s.repo.UpdateLeaveRequest(req); err != nil {
		return err
	}

	// Notify
	_ = s.notifier.Send(emp.Email, "Leave Request Approved", "Your leave request has been approved.")

	// Audit log
	_ = s.audit.Log(managerID, "REJECT", "leave_requests", req.ID, "pending", "rejected", map[string]string{"note": note})

	return nil
}
