package audit

import (
	"time"
)

// AuditLog represents a record in the audit_logs table
type AuditLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	
	ActorID   uint           `json:"actor_id"`
	Action    string         `gorm:"size:50" json:"action"`
	ResourceType string      `gorm:"size:50" json:"resource_type"`
	ResourceID   uint        `json:"resource_id"`
	
	OldValues    string      `gorm:"type:jsonb" json:"old_values"`
	NewValues    string      `gorm:"type:jsonb" json:"new_values"`
	
	IPAddress    string      `gorm:"size:45" json:"ip_address"`
	UserAgent    string      `json:"user_agent"`
	Metadata     string      `gorm:"type:jsonb" json:"metadata"`
}

// Service defines the interface for creating audit logs
type Service interface {
	Log(actorID uint, action string, resourceType string, resourceID uint, oldValues, newValues any, metadata any) error
}
