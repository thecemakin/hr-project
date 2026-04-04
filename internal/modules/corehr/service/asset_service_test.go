package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"github.com/thecemakin/hr-project/internal/platform/testutils"
)

func TestCreateAsset(t *testing.T) {
	t.Run("should fail if required fields are missing", func(t *testing.T) {
		mockRepo := new(testutils.MockAssetRepository)
		mockAudit := new(testutils.MockAuditService)
		svc := NewAssetService(mockRepo, mockAudit)

		asset := &model.Asset{Name: ""} // Missing fields
		err := svc.CreateAsset(1, asset)
		assert.Equal(t, "serial number, name, and type are required", err.Error())
	})

	t.Run("should succeed on valid input", func(t *testing.T) {
		mockRepo := new(testutils.MockAssetRepository)
		mockAudit := new(testutils.MockAuditService)
		svc := NewAssetService(mockRepo, mockAudit)

		asset := &model.Asset{
			SerialNumber: "SN123",
			Name:         "MacBook Pro",
			Type:         "Laptop",
		}

		mockRepo.On("GetAll", 1, 0, mock.Anything).Return([]*model.Asset{}, nil)
		mockRepo.On("Create", asset).Return(nil)
		mockAudit.On("Log", mock.Anything, "CREATE", "assets", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := svc.CreateAsset(1, asset)
		assert.NoError(t, err)
	})
}

func TestUpdateAsset(t *testing.T) {
	t.Run("should fail if asset does not exist", func(t *testing.T) {
		mockRepo := new(testutils.MockAssetRepository)
		mockAudit := new(testutils.MockAuditService)
		svc := NewAssetService(mockRepo, mockAudit)

		asset := &model.Asset{
			ID:           1,
			SerialNumber: "SN123",
			Name:         "MacBook Pro",
			Type:         "Laptop",
		}

		mockRepo.On("GetByID", uint(1)).Return(nil, errors.New("not found"))

		err := svc.UpdateAsset(1, asset)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("should succeed on valid update", func(t *testing.T) {
		mockRepo := new(testutils.MockAssetRepository)
		mockAudit := new(testutils.MockAuditService)
		svc := NewAssetService(mockRepo, mockAudit)

		existing := &model.Asset{ID: 1, SerialNumber: "SN123", Name: "MacBook Pro"}
		updated := &model.Asset{ID: 1, SerialNumber: "SN123", Name: "MacBook Air", Type: "Laptop"}

		mockRepo.On("GetByID", uint(1)).Return(existing, nil)
		mockRepo.On("Update", updated).Return(nil)
		mockAudit.On("Log", mock.Anything, "UPDATE", "assets", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := svc.UpdateAsset(1, updated)
		assert.NoError(t, err)
	})
}
