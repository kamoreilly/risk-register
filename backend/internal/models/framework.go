package models

import "time"

type Framework struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type RiskFrameworkControl struct {
	ID            string    `json:"id" db:"id"`
	RiskID        string    `json:"risk_id" db:"risk_id"`
	FrameworkID   string    `json:"framework_id" db:"framework_id"`
	FrameworkName string    `json:"framework_name" db:"framework_name"` // joined from frameworks
	ControlRef    string    `json:"control_ref" db:"control_ref"`
	Notes         string    `json:"notes,omitempty" db:"notes"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	CreatedBy     string    `json:"created_by" db:"created_by"`
}

type CreateFrameworkInput struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type LinkControlInput struct {
	FrameworkID string `json:"framework_id" validate:"required,uuid"`
	ControlRef  string `json:"control_ref" validate:"required"`
	Notes       string `json:"notes"`
}
