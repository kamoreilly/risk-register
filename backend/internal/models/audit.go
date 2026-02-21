package models

import "time"

type AuditAction string

const (
	AuditActionCreated AuditAction = "created"
	AuditActionUpdated AuditAction = "updated"
	AuditActionDeleted AuditAction = "deleted"
)

type AuditLog struct {
	ID         string         `json:"id" db:"id"`
	EntityType string         `json:"entity_type" db:"entity_type"`
	EntityID   string         `json:"entity_id" db:"entity_id"`
	Action     AuditAction    `json:"action" db:"action"`
	Changes    map[string]any `json:"changes,omitempty" db:"changes"`
	UserID     string         `json:"user_id" db:"user_id"`
	UserName   string         `json:"user_name,omitempty" db:"user_name"` // joined from users
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
}
