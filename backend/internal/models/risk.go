package models

import "time"

type RiskStatus string

const (
	StatusOpen       RiskStatus = "open"
	StatusMitigating RiskStatus = "mitigating"
	StatusResolved   RiskStatus = "resolved"
	StatusAccepted   RiskStatus = "accepted"
)

type RiskSeverity string

const (
	SeverityLow      RiskSeverity = "low"
	SeverityMedium   RiskSeverity = "medium"
	SeverityHigh     RiskSeverity = "high"
	SeverityCritical RiskSeverity = "critical"
)

type Risk struct {
	ID          string       `json:"id" db:"id"`
	Title       string       `json:"title" db:"title"`
	Description string       `json:"description,omitempty" db:"description"`
	OwnerID     string       `json:"owner_id" db:"owner_id"`
	Owner       *User        `json:"owner,omitempty" db:"-"`
	Status      RiskStatus   `json:"status" db:"status"`
	Severity    RiskSeverity `json:"severity" db:"severity"`
	CategoryID  *string      `json:"category_id,omitempty" db:"category_id"`
	Category    *Category    `json:"category,omitempty" db:"-"`
	ReviewDate  *time.Time   `json:"review_date,omitempty" db:"review_date"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	CreatedBy   string       `json:"created_by" db:"created_by"`
	UpdatedBy   string       `json:"updated_by" db:"updated_by"`
}

type CreateRiskInput struct {
	Title       string       `json:"title" validate:"required,min=1,max=255"`
	Description string       `json:"description"`
	OwnerID     string       `json:"owner_id" validate:"required,uuid"`
	Status      RiskStatus   `json:"status"`
	Severity    RiskSeverity `json:"severity"`
	CategoryID  *string      `json:"category_id"`
	ReviewDate  *string      `json:"review_date"`
}

type UpdateRiskInput struct {
	Title       *string       `json:"title" validate:"omitempty,min=1,max=255"`
	Description *string       `json:"description"`
	OwnerID     *string       `json:"owner_id" validate:"omitempty,uuid"`
	Status      *RiskStatus   `json:"status"`
	Severity    *RiskSeverity `json:"severity"`
	CategoryID  *string       `json:"category_id"`
	ReviewDate  *string       `json:"review_date"`
}

type RiskListParams struct {
	Status     *RiskStatus
	Severity   *RiskSeverity
	CategoryID *string
	OwnerID    *string
	Search     string
	Sort       string
	Order      string
	Page       int
	Limit      int
}

type RiskListResponse struct {
	Data  []*Risk `json:"data"`
	Meta  Meta    `json:"meta"`
}

type Meta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}
