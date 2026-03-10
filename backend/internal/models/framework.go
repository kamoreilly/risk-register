package models

import "time"

type Framework struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateFrameworkInput struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type UpdateFrameworkInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
