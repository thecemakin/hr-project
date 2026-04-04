package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"github.com/thecemakin/hr-project/internal/modules/corehr/repository"
	"github.com/thecemakin/hr-project/internal/platform/testutils"
)

func TestCreateEmployee(t *testing.T) {
	t.Run("should fail if required fields are missing", func(t *testing.T) {
		mockRepo := new(testutils.MockEmployeeRepository)
		mockAudit := new(testutils.MockAuditService)
		svc := NewEmployeeService(mockRepo, mockAudit)

		emp := &model.Employee{FirstName: ""} // Missing fields
		err := svc.CreateEmployee(1, emp)
		assert.Equal(t, "first name, last name, email, and employee number are required", err.Error())
	})

	t.Run("should fail if email already exists", func(t *testing.T) {
		mockRepo := new(testutils.MockEmployeeRepository)
		mockAudit := new(testutils.MockAuditService)
		svc := NewEmployeeService(mockRepo, mockAudit)

		emp := &model.Employee{
			FirstName:      "John",
			LastName:       "Doe",
			Email:          "john@example.com",
			EmployeeNumber: "EMP001",
		}

		mockRepo.On("GetByEmail", emp.Email).Return(emp, nil)

		err := svc.CreateEmployee(1, emp)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("should fail if self-managed", func(t *testing.T) {
		mockRepo := new(testutils.MockEmployeeRepository)
		mockAudit := new(testutils.MockAuditService)
		svc := NewEmployeeService(mockRepo, mockAudit)

		id := uint(1)
		emp := &model.Employee{
			ID:             id,
			FirstName:      "John",
			LastName:       "Doe",
			Email:          "john@example.com",
			EmployeeNumber: "EMP001",
			ManagerID:      &id,
		}

		mockRepo.On("GetByEmail", emp.Email).Return(nil, errors.New("not found"))
		mockRepo.On("GetByEmployeeNumber", emp.EmployeeNumber).Return(nil, errors.New("not found"))

		err := svc.CreateEmployee(1, emp)
		assert.Equal(t, "invalid manager relationship: employee cannot be their own manager", err.Error())
	})

	t.Run("should succeed on valid input", func(t *testing.T) {
		mockRepo := new(testutils.MockEmployeeRepository)
		mockAudit := new(testutils.MockAuditService)
		svc := NewEmployeeService(mockRepo, mockAudit)

		emp := &model.Employee{
			FirstName:      "John",
			LastName:       "Doe",
			Email:          "john@example.com",
			EmployeeNumber: "EMP001",
		}

		mockRepo.On("GetByEmail", emp.Email).Return(nil, errors.New("not found"))
		mockRepo.On("GetByEmployeeNumber", emp.EmployeeNumber).Return(nil, errors.New("not found"))
		mockRepo.On("Create", emp).Return(nil)
		mockAudit.On("Log", mock.Anything, "CREATE", "employees", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := svc.CreateEmployee(1, emp)
		assert.NoError(t, err)
	})
}

func TestValidateManagerRelationship(t *testing.T) {
	t.Run("should fail if employee is their own manager", func(t *testing.T) {
		mockRepo := new(testutils.MockEmployeeRepository)
		mockAudit := new(testutils.MockAuditService)
		svc := NewEmployeeService(mockRepo, mockAudit)

		isValid, err := svc.ValidateManagerRelationship(1, 1)
		assert.NoError(t, err)
		assert.False(t, isValid)
	})

	t.Run("should fail if manager does not exist", func(t *testing.T) {
		mockRepo := new(testutils.MockEmployeeRepository)
		mockAudit := new(testutils.MockAuditService)
		svc := NewEmployeeService(mockRepo, mockAudit)

		mockRepo.On("GetByID", uint(2)).Return(nil, errors.New("not found"))

		isValid, err := svc.ValidateManagerRelationship(1, 2)
		assert.Error(t, err)
		assert.False(t, isValid)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("should succeed if manager exists", func(t *testing.T) {
		mockRepo := new(testutils.MockEmployeeRepository)
		mockAudit := new(testutils.MockAuditService)
		svc := NewEmployeeService(mockRepo, mockAudit)

		mockRepo.On("GetByID", uint(2)).Return(&model.Employee{ID: 2}, nil)

		isValid, err := svc.ValidateManagerRelationship(1, 2)
		assert.NoError(t, err)
		assert.True(t, isValid)
	})
}

func TestGetOrganizationTree(t *testing.T) {
	t.Run("should build a valid tree from a flat list", func(t *testing.T) {
		mockRepo := new(testutils.MockEmployeeRepository)
		mockAudit := new(testutils.MockAuditService)
		svc := NewEmployeeService(mockRepo, mockAudit)

		mgrID := uint(1)
		employees := []*model.Employee{
			{ID: 1, FirstName: "Manager", LastName: "One", ManagerID: nil},
			{ID: 2, FirstName: "Employee", LastName: "One", ManagerID: &mgrID},
		}

		mockRepo.On("GetAll", 10000, 0, repository.EmployeeFilter{}).Return(employees, nil)

		tree, err := svc.GetOrganizationTree()
		assert.NoError(t, err)
		assert.Len(t, tree, 1) // Only one root
		assert.Equal(t, uint(1), tree[0].ID)
		assert.Len(t, tree[0].Subordinates, 1)
		assert.Equal(t, uint(2), tree[0].Subordinates[0].ID)
	})
}
