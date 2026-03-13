package models

import "time"

type IncidentStatus string

const (
	IncidentStatusNew          IncidentStatus = "new"
	IncidentStatusAcknowledged IncidentStatus = "acknowledged"
	IncidentStatusInProgress   IncidentStatus = "in_progress"
	IncidentStatusOnHold       IncidentStatus = "on_hold"
	IncidentStatusResolved     IncidentStatus = "resolved"
	IncidentStatusClosed       IncidentStatus = "closed"
)

type IncidentPriority string

const (
	PriorityP1 IncidentPriority = "p1"
	PriorityP2 IncidentPriority = "p2"
	PriorityP3 IncidentPriority = "p3"
	PriorityP4 IncidentPriority = "p4"
)

type IncidentCategory struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateIncidentCategoryInput struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type UpdateIncidentCategoryInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type Incident struct {
	ID              string            `json:"id" db:"id"`
	Title           string            `json:"title" db:"title"`
	Description     string            `json:"description,omitempty" db:"description"`
	CategoryID      *string           `json:"category_id,omitempty" db:"category_id"`
	Category        *IncidentCategory `json:"category,omitempty" db:"-"`
	Priority        IncidentPriority  `json:"priority" db:"priority"`
	Status          IncidentStatus    `json:"status" db:"status"`
	AssigneeID      *string           `json:"assignee_id,omitempty" db:"assignee_id"`
	Assignee        *User             `json:"assignee,omitempty" db:"-"`
	ReporterID      string            `json:"reporter_id" db:"reporter_id"`
	Reporter        *User             `json:"reporter,omitempty" db:"-"`
	ServiceAffected string            `json:"service_affected,omitempty" db:"service_affected"`
	RootCause       string            `json:"root_cause,omitempty" db:"root_cause"`
	ResolutionNotes string            `json:"resolution_notes,omitempty" db:"resolution_notes"`
	OccurredAt      time.Time         `json:"occurred_at" db:"occurred_at"`
	DetectedAt      time.Time         `json:"detected_at" db:"detected_at"`
	ResolvedAt      *time.Time        `json:"resolved_at,omitempty" db:"resolved_at"`
	CreatedAt       time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at" db:"updated_at"`
	CreatedBy       string            `json:"created_by" db:"created_by"`
	UpdatedBy       string            `json:"updated_by" db:"updated_by"`
}

type CreateIncidentInput struct {
	Title           string            `json:"title" validate:"required,min=1,max=255"`
	Description     string            `json:"description"`
	CategoryID      *string           `json:"category_id"`
	Priority        IncidentPriority  `json:"priority"`
	Status          IncidentStatus    `json:"status"`
	AssigneeID      *string           `json:"assignee_id" validate:"omitempty,uuid"`
	ServiceAffected string            `json:"service_affected"`
	OccurredAt      *string           `json:"occurred_at"`
	DetectedAt      *string           `json:"detected_at"`
}

type UpdateIncidentInput struct {
	Title           *string           `json:"title" validate:"omitempty,min=1,max=255"`
	Description     *string           `json:"description"`
	CategoryID      *string           `json:"category_id"`
	Priority        *IncidentPriority `json:"priority"`
	Status          *IncidentStatus   `json:"status"`
	AssigneeID      *string           `json:"assignee_id" validate:"omitempty,uuid"`
	ServiceAffected *string           `json:"service_affected"`
	RootCause       *string           `json:"root_cause"`
	ResolutionNotes *string           `json:"resolution_notes"`
	ResolvedAt      *string           `json:"resolved_at"`
}

type IncidentListParams struct {
	Status     *IncidentStatus
	Priority   *IncidentPriority
	CategoryID *string
	AssigneeID *string
	Search     string
	Sort       string
	Order      string
	Page       int
	Limit      int
}

type IncidentListResponse struct {
	Data []*Incident `json:"data"`
	Meta Meta        `json:"meta"`
}

type IncidentRisk struct {
	ID         string    `json:"id" db:"id"`
	IncidentID string    `json:"incident_id" db:"incident_id"`
	RiskID     string    `json:"risk_id" db:"risk_id"`
	Risk       *Risk     `json:"risk,omitempty" db:"-"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	CreatedBy  string    `json:"created_by" db:"created_by"`
}

type LinkIncidentRiskInput struct {
	RiskID string `json:"risk_id" validate:"required,uuid"`
}
