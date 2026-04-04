package testutils

import (
	"github.com/stretchr/testify/mock"
	corehrmodel "github.com/thecemakin/hr-project/internal/modules/corehr/model"
	corehrrepo "github.com/thecemakin/hr-project/internal/modules/corehr/repository"
	leavemodel "github.com/thecemakin/hr-project/internal/modules/leave/model"
)

// MockLeaveRepository mocks the leave repository
type MockLeaveRepository struct {
	mock.Mock
}

func (m *MockLeaveRepository) CreateLeaveType(lt *leavemodel.LeaveType) error {
	args := m.Called(lt)
	return args.Error(0)
}
func (m *MockLeaveRepository) GetLeaveTypeByID(id uint) (*leavemodel.LeaveType, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*leavemodel.LeaveType), args.Error(1)
}
func (m *MockLeaveRepository) ListLeaveTypes() ([]leavemodel.LeaveType, error) {
	args := m.Called()
	return args.Get(0).([]leavemodel.LeaveType), args.Error(1)
}
func (m *MockLeaveRepository) UpdateLeaveType(lt *leavemodel.LeaveType) error {
	args := m.Called(lt)
	return args.Error(0)
}
func (m *MockLeaveRepository) CreateLeaveBalance(lb *leavemodel.LeaveBalance) error {
	args := m.Called(lb)
	return args.Error(0)
}
func (m *MockLeaveRepository) GetLeaveBalance(employeeID, leaveTypeID uint) (*leavemodel.LeaveBalance, error) {
	args := m.Called(employeeID, leaveTypeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*leavemodel.LeaveBalance), args.Error(1)
}
func (m *MockLeaveRepository) ListLeaveBalancesByEmployee(employeeID uint) ([]leavemodel.LeaveBalance, error) {
	args := m.Called(employeeID)
	return args.Get(0).([]leavemodel.LeaveBalance), args.Error(1)
}
func (m *MockLeaveRepository) UpdateLeaveBalance(lb *leavemodel.LeaveBalance) error {
	args := m.Called(lb)
	return args.Error(0)
}
func (m *MockLeaveRepository) CreateLeaveRequest(req *leavemodel.LeaveRequest) error {
	args := m.Called(req)
	return args.Error(0)
}
func (m *MockLeaveRepository) GetLeaveRequestByID(id uint) (*leavemodel.LeaveRequest, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*leavemodel.LeaveRequest), args.Error(1)
}
func (m *MockLeaveRepository) ListLeaveRequestsByEmployee(employeeID uint) ([]leavemodel.LeaveRequest, error) {
	args := m.Called(employeeID)
	return args.Get(0).([]leavemodel.LeaveRequest), args.Error(1)
}
func (m *MockLeaveRepository) ListPendingLeaveRequestsForManager(managerID uint) ([]leavemodel.LeaveRequest, error) {
	args := m.Called(managerID)
	return args.Get(0).([]leavemodel.LeaveRequest), args.Error(1)
}
func (m *MockLeaveRepository) UpdateLeaveRequest(req *leavemodel.LeaveRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

// MockEmployeeRepository mocks the corehr employee repository
type MockEmployeeRepository struct {
	mock.Mock
}

func (m *MockEmployeeRepository) Create(employee *corehrmodel.Employee) error {
	args := m.Called(employee)
	return args.Error(0)
}
func (m *MockEmployeeRepository) GetByID(id uint) (*corehrmodel.Employee, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*corehrmodel.Employee), args.Error(1)
}
func (m *MockEmployeeRepository) GetByEmail(email string) (*corehrmodel.Employee, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*corehrmodel.Employee), args.Error(1)
}
func (m *MockEmployeeRepository) GetByEmployeeNumber(num string) (*corehrmodel.Employee, error) {
	args := m.Called(num)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*corehrmodel.Employee), args.Error(1)
}
func (m *MockEmployeeRepository) GetAll(l, o int, f corehrrepo.EmployeeFilter) ([]*corehrmodel.Employee, error) {
	args := m.Called(l, o, f)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*corehrmodel.Employee), args.Error(1)
}
func (m *MockEmployeeRepository) Update(employee *corehrmodel.Employee) error {
	args := m.Called(employee)
	return args.Error(0)
}
func (m *MockEmployeeRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockAuditService mocks the audit service
type MockAuditService struct {
	mock.Mock
}

func (m *MockAuditService) Log(actorID uint, action string, resourceType string, resourceID uint, oldValues, newValues any, metadata any) error {
	args := m.Called(actorID, action, resourceType, resourceID, oldValues, newValues, metadata)
	return args.Error(0)
}

// MockNotifier mocks the notification service
type MockNotifier struct {
	mock.Mock
}

func (m *MockNotifier) Send(recipient, subject, message string) error {
	args := m.Called(recipient, subject, message)
	return args.Error(0)
}

// MockAssetRepository mocks the corehr asset repository
type MockAssetRepository struct {
	mock.Mock
}

func (m *MockAssetRepository) Create(asset *corehrmodel.Asset) error {
	args := m.Called(asset)
	return args.Error(0)
}

func (m *MockAssetRepository) GetByID(id uint) (*corehrmodel.Asset, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*corehrmodel.Asset), args.Error(1)
}

func (m *MockAssetRepository) GetAll(l, o int, f corehrrepo.AssetFilter) ([]*corehrmodel.Asset, error) {
	args := m.Called(l, o, f)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*corehrmodel.Asset), args.Error(1)
}

func (m *MockAssetRepository) Update(asset *corehrmodel.Asset) error {
	args := m.Called(asset)
	return args.Error(0)
}

func (m *MockAssetRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockAssetService mocks the corehr asset service
type MockAssetService struct {
	mock.Mock
}

func (m *MockAssetService) CreateAsset(actorID uint, asset *corehrmodel.Asset) error {
	args := m.Called(actorID, asset)
	return args.Error(0)
}

func (m *MockAssetService) GetAssetByID(id uint) (*corehrmodel.Asset, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*corehrmodel.Asset), args.Error(1)
}

func (m *MockAssetService) GetAllAssets(l, o int, f corehrrepo.AssetFilter) ([]*corehrmodel.Asset, error) {
	args := m.Called(l, o, f)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*corehrmodel.Asset), args.Error(1)
}

func (m *MockAssetService) UpdateAsset(actorID uint, asset *corehrmodel.Asset) error {
	args := m.Called(actorID, asset)
	return args.Error(0)
}

func (m *MockAssetService) DeleteAsset(actorID uint, id uint) error {
	args := m.Called(actorID, id)
	return args.Error(0)
}
