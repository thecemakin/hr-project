package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	corehrmodel "github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"github.com/thecemakin/hr-project/internal/modules/leave/model"
	"github.com/thecemakin/hr-project/internal/platform/testutils"
)

func TestSubmitLeaveRequest(t *testing.T) {
	t.Run("should fail if start date is after end date", func(t *testing.T) {
		mockRepo := new(testutils.MockLeaveRepository)
		mockEmpRepo := new(testutils.MockEmployeeRepository)
		mockAudit := new(testutils.MockAuditService)
		mockNotifier := new(testutils.MockNotifier)
		svc := NewLeaveService(mockRepo, mockEmpRepo, mockAudit, mockNotifier)

		req := &model.LeaveRequest{
			StartDate: time.Now().AddDate(0, 0, 5),
			EndDate:   time.Now().AddDate(0, 0, 2),
		}
		err := svc.SubmitLeaveRequest(req)
		assert.Equal(t, "start date must be before or equal to end date", err.Error())
	})

	t.Run("should fail if insufficient balance", func(t *testing.T) {
		mockRepo := new(testutils.MockLeaveRepository)
		mockEmpRepo := new(testutils.MockEmployeeRepository)
		mockAudit := new(testutils.MockAuditService)
		mockNotifier := new(testutils.MockNotifier)
		svc := NewLeaveService(mockRepo, mockEmpRepo, mockAudit, mockNotifier)

		req := &model.LeaveRequest{
			EmployeeID:          1,
			LeaveTypeID:         1,
			StartDate:           time.Now(),
			EndDate:             time.Now().AddDate(0, 0, 5),
			TotalDaysRequested: 10,
		}
		balance := &model.LeaveBalance{TotalDays: 15, UsedDays: 10} // 5 available

		mockRepo.On("GetLeaveBalance", uint(1), uint(1)).Return(balance, nil)

		err := svc.SubmitLeaveRequest(req)
		assert.Equal(t, "insufficient leave balance", err.Error())
	})

	t.Run("should succeed on valid request", func(t *testing.T) {
		mockRepo := new(testutils.MockLeaveRepository)
		mockEmpRepo := new(testutils.MockEmployeeRepository)
		mockAudit := new(testutils.MockAuditService)
		mockNotifier := new(testutils.MockNotifier)
		svc := NewLeaveService(mockRepo, mockEmpRepo, mockAudit, mockNotifier)

		req := &model.LeaveRequest{
			EmployeeID:          1,
			LeaveTypeID:         1,
			StartDate:           time.Now(),
			EndDate:             time.Now().AddDate(0, 0, 2),
			TotalDaysRequested: 3,
		}
		balance := &model.LeaveBalance{TotalDays: 15, UsedDays: 5} // 10 available

		mockRepo.On("GetLeaveBalance", uint(1), uint(1)).Return(balance, nil)
		mockRepo.On("CreateLeaveRequest", req).Return(nil)
		mockEmpRepo.On("GetByID", uint(1)).Return(&corehrmodel.Employee{Email: "test@example.com"}, nil)
		mockNotifier.On("Send", "test@example.com", mock.Anything, mock.Anything).Return(nil)
		mockAudit.On("Log", mock.Anything, "CREATE", "leave_requests", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := svc.SubmitLeaveRequest(req)
		assert.NoError(t, err)
		assert.Equal(t, "pending", req.Status)
	})
}

func TestApproveLeaveRequest(t *testing.T) {
	t.Run("should fail if manager is not authorized", func(t *testing.T) {
		mockRepo := new(testutils.MockLeaveRepository)
		mockEmpRepo := new(testutils.MockEmployeeRepository)
		mockAudit := new(testutils.MockAuditService)
		mockNotifier := new(testutils.MockNotifier)
		svc := NewLeaveService(mockRepo, mockEmpRepo, mockAudit, mockNotifier)

		req := &model.LeaveRequest{
			ID:                 1,
			EmployeeID:         10,
			Status:             "pending",
			TotalDaysRequested: 2,
		}
		managerID := uint(50)
		otherManagerID := uint(99)
		emp := &corehrmodel.Employee{ID: 10, ManagerID: &otherManagerID}

		mockRepo.On("GetLeaveRequestByID", uint(1)).Return(req, nil)
		mockEmpRepo.On("GetByID", uint(10)).Return(emp, nil)

		err := svc.ApproveLeaveRequest(1, managerID, "Approved")
		assert.Equal(t, "unauthorized: you are not the manager of this employee", err.Error())
	})

	t.Run("should succeed and deduct balance", func(t *testing.T) {
		mockRepo := new(testutils.MockLeaveRepository)
		mockEmpRepo := new(testutils.MockEmployeeRepository)
		mockAudit := new(testutils.MockAuditService)
		mockNotifier := new(testutils.MockNotifier)
		svc := NewLeaveService(mockRepo, mockEmpRepo, mockAudit, mockNotifier)

		req := &model.LeaveRequest{
			ID:                 1,
			EmployeeID:         10,
			LeaveTypeID:        1,
			Status:             "pending",
			TotalDaysRequested: 2,
		}
		managerID := uint(50)
		emp := &corehrmodel.Employee{ID: 10, ManagerID: &managerID, Email: "test@example.com"}
		balance := &model.LeaveBalance{EmployeeID: 10, LeaveTypeID: 1, TotalDays: 10, UsedDays: 2}

		mockRepo.On("GetLeaveRequestByID", uint(1)).Return(req, nil)
		mockEmpRepo.On("GetByID", uint(10)).Return(emp, nil)
		mockRepo.On("GetLeaveBalance", uint(10), uint(1)).Return(balance, nil)
		mockRepo.On("UpdateLeaveBalance", mock.Anything).Return(nil)
		mockRepo.On("UpdateLeaveRequest", mock.Anything).Return(nil)
		mockNotifier.On("Send", "test@example.com", mock.Anything, mock.Anything).Return(nil)
		mockAudit.On("Log", mock.Anything, "APPROVE", "leave_requests", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := svc.ApproveLeaveRequest(1, managerID, "Approved")
		assert.NoError(t, err)
		assert.Equal(t, 4, balance.UsedDays)
		assert.Equal(t, "approved", req.Status)
	})
}
