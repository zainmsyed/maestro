package models

import "time"

type DateAuditLog struct {
	ID         int64
	EntityType string
	EntityID   string
	ChangedBy  string
	OldDate    *time.Time
	NewDate    *time.Time
	DeltaDays  int
	Reason     *string
	ChangedAt  time.Time
}
