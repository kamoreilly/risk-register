package models

import "time"

type MitigationStatus string

const (
	MitigationStatusPlanned    MitigationStatus = "planned"
	MitigationStatusInProgress MitigationStatus = "in_progress"
	MitigationStatusCompleted  MitigationStatus = "completed"
	MitigationStatusCancelled  MitigationStatus = "cancelled"
)

type Mitigation struct {
	ID          string           `json:"id" db:"id"`
	RiskID      string           `json:"risk_id" db:"risk_id"`
	Description string           `json:"description,omitempty" db:"description"`
	Owner       string           `json:"owner,omitempty" db:"owner"`
	Status      MitigationStatus `json:"status" db:"status"`
	DueDate     *time.Time       `json:"due_date,omitempty" db:"due_date"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
	CreatedBy   string           `json:"created_by" db:"created_by"`
	UpdatedBy   string           `json:"updated_by" db:"updated_by"`
}

type CreateMitigationInput struct {
	RiskID      string           `json:"risk_id" validate:"required,uuid"`
	Description string           `json:"description" validate:"required"`
	Owner       string           `json:"owner" validate:"required,min=1,max=255"`
	Status      MitigationStatus `json:"status"`
	DueDate     *string          `json:"due_date"`
}

type UpdateMitigationInput struct {
	Description *string           `json:"description" validate:"omitempty"`
	Owner       *string           `json:"owner" validate:"omitempty,min=1,max=255"`
	Status      *MitigationStatus `json:"status"`
	DueDate     *string           `json:"due_date"`
}
