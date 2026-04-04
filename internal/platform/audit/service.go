package audit

import (
	"encoding/json"
	"log"

	"gorm.io/gorm"
)

type auditService struct {
	db *gorm.DB
}

// NewAuditService creates a new instance of the audit service
func NewAuditService(db *gorm.DB) Service {
	return &auditService{
		db: db,
	}
}

func (s *auditService) Log(actorID uint, action string, resourceType string, resourceID uint, oldValues, newValues any, metadata any) error {
	var oldV, newV, meta []byte
	var err error

	if oldValues != nil {
		oldV, err = json.Marshal(oldValues)
		if err != nil {
			log.Printf("[AUDIT ERROR] marshaling oldValues: %v", err)
		}
	}

	if newValues != nil {
		newV, err = json.Marshal(newValues)
		if err != nil {
			log.Printf("[AUDIT ERROR] marshaling newValues: %v", err)
		}
	}

	if metadata != nil {
		meta, err = json.Marshal(metadata)
		if err != nil {
			log.Printf("[AUDIT ERROR] marshaling metadata: %v", err)
		}
	}

	auditLog := &AuditLog{
		ActorID:      actorID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		OldValues:    string(oldV),
		NewValues:    string(newV),
		Metadata:     string(meta),
	}

	// We use a separate Goroutine or ensure this doesn't block the main flow
	// For v1, we just Save it. In a high-traffic app, we might use a queue.
	if err := s.db.Create(auditLog).Error; err != nil {
		log.Printf("[AUDIT ERROR] failed to save audit log: %v", err)
		return err
	}

	return nil
}
