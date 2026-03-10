package models

import "time"

type FrameworkControl struct {
	ID              string    `json:"id" db:"id"`
	FrameworkID     string    `json:"framework_id" db:"framework_id"`
	FrameworkName   string    `json:"framework_name" db:"framework_name"`
	ControlRef      string    `json:"control_ref" db:"control_ref"`
	Title           string    `json:"title" db:"title"`
	Description     string    `json:"description,omitempty" db:"description"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	LinkedRiskCount int       `json:"linked_risk_count" db:"linked_risk_count"`
}

type CreateFrameworkControlInput struct {
	FrameworkID string `json:"framework_id" validate:"required,uuid"`
	ControlRef  string `json:"control_ref" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}

type UpdateFrameworkControlInput struct {
	ControlRef  *string `json:"control_ref"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type LinkControlInput struct {
	FrameworkControlID string `json:"framework_control_id" validate:"required,uuid"`
	Notes              string `json:"notes"`
}

type RiskFrameworkControl struct {
	ID                 string    `json:"id" db:"id"`
	RiskID             string    `json:"risk_id" db:"risk_id"`
	FrameworkControlID string    `json:"framework_control_id" db:"framework_control_id"`
	FrameworkID        string    `json:"framework_id" db:"framework_id"`
	FrameworkName      string    `json:"framework_name" db:"framework_name"`
	ControlRef         string    `json:"control_ref" db:"control_ref"`
	ControlTitle       string    `json:"control_title" db:"control_title"`
	ControlDescription string    `json:"control_description,omitempty" db:"control_description"`
	Notes              string    `json:"notes,omitempty" db:"notes"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	CreatedBy          string    `json:"created_by" db:"created_by"`
}

type ControlLinkedRisk struct {
	ID           string       `json:"id" db:"id"`
	Title        string       `json:"title" db:"title"`
	Status       RiskStatus   `json:"status" db:"status"`
	Severity     RiskSeverity `json:"severity" db:"severity"`
	CategoryName string       `json:"category_name,omitempty" db:"category_name"`
	OwnerName    string       `json:"owner_name,omitempty" db:"owner_name"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
}
