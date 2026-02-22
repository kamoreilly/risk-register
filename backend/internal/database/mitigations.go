package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"backend/internal/models"

	"github.com/google/uuid"
)

var ErrMitigationNotFound = errors.New("mitigation not found")

// parseDate attempts to parse a date string in multiple formats
func parseDate(dateStr string) (time.Time, error) {
	// Try RFC3339 format first (e.g., "2006-01-02T15:04:05Z07:00")
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return t, nil
	}
	// Try date-only format (e.g., "2006-01-02")
	return time.Parse("2006-01-02", dateStr)
}

type MitigationRepository interface {
	Create(ctx context.Context, input *models.CreateMitigationInput, createdBy string) (*models.Mitigation, error)
	FindByID(ctx context.Context, id string) (*models.Mitigation, error)
	ListByRiskID(ctx context.Context, riskID string) ([]*models.Mitigation, error)
	Update(ctx context.Context, id string, input *models.UpdateMitigationInput, updatedBy string) (*models.Mitigation, error)
	Delete(ctx context.Context, id string) error
}

type mitigationRepository struct {
	db *sql.DB
}

func NewMitigationRepository(db *sql.DB) MitigationRepository {
	return &mitigationRepository{db: db}
}

func (r *mitigationRepository) Create(ctx context.Context, input *models.CreateMitigationInput, createdBy string) (*models.Mitigation, error) {
	mitigation := &models.Mitigation{
		ID:          uuid.New().String(),
		RiskID:      input.RiskID,
		Description: input.Description,
		Owner:       input.Owner,
		Status:      input.Status,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}

	// Parse due_date if provided
	if input.DueDate != nil && *input.DueDate != "" {
		dueDate, err := parseDate(*input.DueDate)
		if err != nil {
			return nil, err
		}
		mitigation.DueDate = &dueDate
	}

	now := time.Now()
	mitigation.CreatedAt = now
	mitigation.UpdatedAt = now

	query := `
		INSERT INTO mitigations (id, risk_id, description, owner, status, due_date, created_at, updated_at, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		mitigation.ID,
		mitigation.RiskID,
		mitigation.Description,
		mitigation.Owner,
		mitigation.Status,
		mitigation.DueDate,
		mitigation.CreatedAt,
		mitigation.UpdatedAt,
		mitigation.CreatedBy,
		mitigation.UpdatedBy,
	).Scan(&mitigation.ID, &mitigation.CreatedAt, &mitigation.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return mitigation, nil
}

func (r *mitigationRepository) FindByID(ctx context.Context, id string) (*models.Mitigation, error) {
	query := `
		SELECT id, risk_id, description, owner, status, due_date, created_at, updated_at, created_by, updated_by
		FROM mitigations WHERE id = $1
	`

	mitigation := &models.Mitigation{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&mitigation.ID,
		&mitigation.RiskID,
		&mitigation.Description,
		&mitigation.Owner,
		&mitigation.Status,
		&mitigation.DueDate,
		&mitigation.CreatedAt,
		&mitigation.UpdatedAt,
		&mitigation.CreatedBy,
		&mitigation.UpdatedBy,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMitigationNotFound
		}
		return nil, err
	}

	return mitigation, nil
}

func (r *mitigationRepository) ListByRiskID(ctx context.Context, riskID string) ([]*models.Mitigation, error) {
	query := `
		SELECT id, risk_id, description, owner, status, due_date, created_at, updated_at, created_by, updated_by
		FROM mitigations WHERE risk_id = $1 ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, riskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mitigations []*models.Mitigation
	for rows.Next() {
		mitigation := &models.Mitigation{}
		err := rows.Scan(
			&mitigation.ID,
			&mitigation.RiskID,
			&mitigation.Description,
			&mitigation.Owner,
			&mitigation.Status,
			&mitigation.DueDate,
			&mitigation.CreatedAt,
			&mitigation.UpdatedAt,
			&mitigation.CreatedBy,
			&mitigation.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		mitigations = append(mitigations, mitigation)
	}

	return mitigations, rows.Err()
}

func (r *mitigationRepository) Update(ctx context.Context, id string, input *models.UpdateMitigationInput, updatedBy string) (*models.Mitigation, error) {
	// First, get the existing mitigation
	mitigation, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if input.Description != nil {
		mitigation.Description = *input.Description
	}
	if input.Owner != nil {
		mitigation.Owner = *input.Owner
	}
	if input.Status != nil {
		mitigation.Status = *input.Status
	}
	if input.DueDate != nil {
		if *input.DueDate == "" {
			mitigation.DueDate = nil
		} else {
			dueDate, err := parseDate(*input.DueDate)
			if err != nil {
				return nil, err
			}
			mitigation.DueDate = &dueDate
		}
	}

	mitigation.UpdatedBy = updatedBy
	now := time.Now()
	mitigation.UpdatedAt = now

	query := `
		UPDATE mitigations SET description = $1, owner = $2, status = $3, due_date = $4, updated_at = $5, updated_by = $6
		WHERE id = $7
		RETURNING updated_at
	`

	err = r.db.QueryRowContext(ctx, query,
		mitigation.Description,
		mitigation.Owner,
		mitigation.Status,
		mitigation.DueDate,
		mitigation.UpdatedAt,
		mitigation.UpdatedBy,
		mitigation.ID,
	).Scan(&mitigation.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMitigationNotFound
		}
		return nil, err
	}

	return mitigation, nil
}

func (r *mitigationRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM mitigations WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrMitigationNotFound
	}

	return nil
}
