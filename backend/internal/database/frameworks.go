package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"backend/internal/models"

	"github.com/google/uuid"
)

var ErrFrameworkNotFound = errors.New("framework not found")

type FrameworkRepository interface {
	List(ctx context.Context) ([]*models.Framework, error)
	GetByID(ctx context.Context, id string) (*models.Framework, error)
	Create(ctx context.Context, input *models.CreateFrameworkInput) (*models.Framework, error)
}

type RiskFrameworkControlRepository interface {
	ListByRiskID(ctx context.Context, riskID string) ([]*models.RiskFrameworkControl, error)
	LinkControl(ctx context.Context, riskID string, input *models.LinkControlInput, createdBy string) (*models.RiskFrameworkControl, error)
	UnlinkControl(ctx context.Context, id string) error
}

type frameworkRepository struct {
	db *sql.DB
}

type riskFrameworkControlRepository struct {
	db *sql.DB
}

func NewFrameworkRepository(db *sql.DB) FrameworkRepository {
	return &frameworkRepository{db: db}
}

func NewRiskFrameworkControlRepository(db *sql.DB) RiskFrameworkControlRepository {
	return &riskFrameworkControlRepository{db: db}
}

func (r *frameworkRepository) List(ctx context.Context) ([]*models.Framework, error) {
	query := `
		SELECT id, name, description, created_at
		FROM frameworks ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var frameworks []*models.Framework
	for rows.Next() {
		framework := &models.Framework{}
		err := rows.Scan(
			&framework.ID,
			&framework.Name,
			&framework.Description,
			&framework.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		frameworks = append(frameworks, framework)
	}

	return frameworks, rows.Err()
}

func (r *frameworkRepository) GetByID(ctx context.Context, id string) (*models.Framework, error) {
	query := `
		SELECT id, name, description, created_at
		FROM frameworks WHERE id = $1
	`

	framework := &models.Framework{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&framework.ID,
		&framework.Name,
		&framework.Description,
		&framework.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFrameworkNotFound
		}
		return nil, err
	}

	return framework, nil
}

func (r *frameworkRepository) Create(ctx context.Context, input *models.CreateFrameworkInput) (*models.Framework, error) {
	framework := &models.Framework{
		ID:          uuid.New().String(),
		Name:        input.Name,
		Description: input.Description,
		CreatedAt:   time.Now(),
	}

	query := `
		INSERT INTO frameworks (id, name, description, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(ctx, query,
		framework.ID,
		framework.Name,
		framework.Description,
		framework.CreatedAt,
	).Scan(&framework.ID, &framework.CreatedAt)

	if err != nil {
		return nil, err
	}

	return framework, nil
}

func (r *riskFrameworkControlRepository) ListByRiskID(ctx context.Context, riskID string) ([]*models.RiskFrameworkControl, error) {
	query := `
		SELECT rfc.id, rfc.risk_id, rfc.framework_id, f.name as framework_name, rfc.control_ref, rfc.notes, rfc.created_at, rfc.created_by
		FROM risk_framework_controls rfc
		JOIN frameworks f ON rfc.framework_id = f.id
		WHERE rfc.risk_id = $1
		ORDER BY f.name ASC, rfc.control_ref ASC
	`

	rows, err := r.db.QueryContext(ctx, query, riskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var controls []*models.RiskFrameworkControl
	for rows.Next() {
		control := &models.RiskFrameworkControl{}
		err := rows.Scan(
			&control.ID,
			&control.RiskID,
			&control.FrameworkID,
			&control.FrameworkName,
			&control.ControlRef,
			&control.Notes,
			&control.CreatedAt,
			&control.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		controls = append(controls, control)
	}

	return controls, rows.Err()
}

func (r *riskFrameworkControlRepository) LinkControl(ctx context.Context, riskID string, input *models.LinkControlInput, createdBy string) (*models.RiskFrameworkControl, error) {
	control := &models.RiskFrameworkControl{
		ID:          uuid.New().String(),
		RiskID:      riskID,
		FrameworkID: input.FrameworkID,
		ControlRef:  input.ControlRef,
		Notes:       input.Notes,
		CreatedAt:   time.Now(),
		CreatedBy:   createdBy,
	}

	query := `
		INSERT INTO risk_framework_controls (id, risk_id, framework_id, control_ref, notes, created_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(ctx, query,
		control.ID,
		control.RiskID,
		control.FrameworkID,
		control.ControlRef,
		control.Notes,
		control.CreatedAt,
		control.CreatedBy,
	).Scan(&control.ID, &control.CreatedAt)

	if err != nil {
		return nil, err
	}

	// Fetch the framework name for the response
	frameworkQuery := `SELECT name FROM frameworks WHERE id = $1`
	err = r.db.QueryRowContext(ctx, frameworkQuery, control.FrameworkID).Scan(&control.FrameworkName)
	if err != nil {
		return nil, err
	}

	return control, nil
}

func (r *riskFrameworkControlRepository) UnlinkControl(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM risk_framework_controls WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrFrameworkNotFound
	}

	return nil
}
