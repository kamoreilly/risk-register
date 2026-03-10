package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"backend/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrFrameworkNotFound = errors.New("framework not found")
var ErrFrameworkControlNotFound = errors.New("framework control not found")
var ErrFrameworkControlInUse = errors.New("framework control is linked to risks")

type FrameworkRepository interface {
	List(ctx context.Context) ([]*models.Framework, error)
	GetByID(ctx context.Context, id string) (*models.Framework, error)
	Create(ctx context.Context, input *models.CreateFrameworkInput) (*models.Framework, error)
	Update(ctx context.Context, id string, input *models.UpdateFrameworkInput) (*models.Framework, error)
	Delete(ctx context.Context, id string) error
}

type FrameworkControlRepository interface {
	List(ctx context.Context, frameworkID, search string) ([]*models.FrameworkControl, error)
	GetByID(ctx context.Context, id string) (*models.FrameworkControl, error)
	ListLinkedRisks(ctx context.Context, id string) ([]*models.ControlLinkedRisk, error)
	Create(ctx context.Context, input *models.CreateFrameworkControlInput) (*models.FrameworkControl, error)
	Update(ctx context.Context, id string, input *models.UpdateFrameworkControlInput) (*models.FrameworkControl, error)
	Delete(ctx context.Context, id string) error
}

type RiskFrameworkControlRepository interface {
	ListByRiskID(ctx context.Context, riskID string) ([]*models.RiskFrameworkControl, error)
	LinkControl(ctx context.Context, riskID string, input *models.LinkControlInput, createdBy string) (*models.RiskFrameworkControl, error)
	UnlinkControl(ctx context.Context, id string) error
}

type frameworkRepository struct {
	db *sql.DB
}

type frameworkControlRepository struct {
	db *sql.DB
}

type riskFrameworkControlRepository struct {
	db *sql.DB
}

func NewFrameworkRepository(db *sql.DB) FrameworkRepository {
	return &frameworkRepository{db: db}
}

func NewFrameworkControlRepository(db *sql.DB) FrameworkControlRepository {
	return &frameworkControlRepository{db: db}
}

func NewRiskFrameworkControlRepository(db *sql.DB) RiskFrameworkControlRepository {
	return &riskFrameworkControlRepository{db: db}
}

func (r *frameworkRepository) List(ctx context.Context) ([]*models.Framework, error) {
	query := `
		SELECT id, name, COALESCE(description, ''), created_at, updated_at
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
			&framework.UpdatedAt,
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
		SELECT id, name, COALESCE(description, ''), created_at, updated_at
		FROM frameworks WHERE id = $1
	`

	framework := &models.Framework{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&framework.ID,
		&framework.Name,
		&framework.Description,
		&framework.CreatedAt,
		&framework.UpdatedAt,
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
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		framework.ID,
		framework.Name,
		framework.Description,
		framework.CreatedAt,
	).Scan(&framework.ID, &framework.CreatedAt, &framework.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return framework, nil
}

func (r *frameworkRepository) Update(ctx context.Context, id string, input *models.UpdateFrameworkInput) (*models.Framework, error) {
	if input.Name == nil && input.Description == nil {
		return nil, errors.New("at least one field must be updated")
	}

	framework := &models.Framework{}
	query := `
		UPDATE frameworks
		SET name = COALESCE($1, name), description = COALESCE($2, description), updated_at = NOW()
		WHERE id = $3
		RETURNING id, name, description, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, input.Name, input.Description, id).Scan(
		&framework.ID, &framework.Name, &framework.Description, &framework.CreatedAt, &framework.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFrameworkNotFound
		}
		return nil, err
	}
	return framework, nil
}

func (r *frameworkRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM frameworks WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
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

func (r *frameworkControlRepository) List(ctx context.Context, frameworkID, search string) ([]*models.FrameworkControl, error) {
	// Validate UUID format if frameworkID is provided
	if frameworkID != "" {
		if _, err := uuid.Parse(frameworkID); err != nil {
			return nil, errors.New("invalid framework ID format")
		}
	}

	query := `
		SELECT fc.id, fc.framework_id, f.name, fc.control_ref, fc.title, COALESCE(fc.description, ''),
			fc.created_at, fc.updated_at, COUNT(rfc.id) AS linked_risk_count
		FROM framework_controls fc
		JOIN frameworks f ON f.id = fc.framework_id
		LEFT JOIN risk_framework_controls rfc ON rfc.framework_control_id = fc.id
		WHERE (NULLIF($1, '') IS NULL OR fc.framework_id = NULLIF($1, '')::uuid)
		  AND (
			$2 = '' OR
			fc.control_ref ILIKE '%' || $2 || '%' OR
			fc.title ILIKE '%' || $2 || '%' OR
			COALESCE(fc.description, '') ILIKE '%' || $2 || '%'
		  )
		GROUP BY fc.id, fc.framework_id, f.name, fc.control_ref, fc.title, fc.description, fc.created_at, fc.updated_at
		ORDER BY f.name ASC, fc.control_ref ASC
	`

	rows, err := r.db.QueryContext(ctx, query, frameworkID, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var controls []*models.FrameworkControl
	for rows.Next() {
		control := &models.FrameworkControl{}
		if err := rows.Scan(
			&control.ID,
			&control.FrameworkID,
			&control.FrameworkName,
			&control.ControlRef,
			&control.Title,
			&control.Description,
			&control.CreatedAt,
			&control.UpdatedAt,
			&control.LinkedRiskCount,
		); err != nil {
			return nil, err
		}
		controls = append(controls, control)
	}

	return controls, rows.Err()
}

func (r *frameworkControlRepository) GetByID(ctx context.Context, id string) (*models.FrameworkControl, error) {
	query := `
		SELECT fc.id, fc.framework_id, f.name, fc.control_ref, fc.title, COALESCE(fc.description, ''),
			fc.created_at, fc.updated_at,
			(SELECT COUNT(*) FROM risk_framework_controls rfc WHERE rfc.framework_control_id = fc.id) AS linked_risk_count
		FROM framework_controls fc
		JOIN frameworks f ON f.id = fc.framework_id
		WHERE fc.id = $1
	`

	control := &models.FrameworkControl{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&control.ID,
		&control.FrameworkID,
		&control.FrameworkName,
		&control.ControlRef,
		&control.Title,
		&control.Description,
		&control.CreatedAt,
		&control.UpdatedAt,
		&control.LinkedRiskCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFrameworkControlNotFound
		}
		return nil, err
	}

	return control, nil
}

func (r *frameworkControlRepository) ListLinkedRisks(ctx context.Context, id string) ([]*models.ControlLinkedRisk, error) {
	if _, err := r.GetByID(ctx, id); err != nil {
		return nil, err
	}

	query := `
		SELECT r.id, r.title, r.status, r.severity,
			COALESCE(c.name, '') AS category_name,
			COALESCE(u.name, '') AS owner_name,
			r.updated_at
		FROM risk_framework_controls rfc
		JOIN risks r ON r.id = rfc.risk_id
		LEFT JOIN categories c ON c.id = r.category_id
		LEFT JOIN users u ON u.id = r.owner_id
		WHERE rfc.framework_control_id = $1
		ORDER BY r.updated_at DESC, r.title ASC
	`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var risks []*models.ControlLinkedRisk
	for rows.Next() {
		risk := &models.ControlLinkedRisk{}
		if err := rows.Scan(
			&risk.ID,
			&risk.Title,
			&risk.Status,
			&risk.Severity,
			&risk.CategoryName,
			&risk.OwnerName,
			&risk.UpdatedAt,
		); err != nil {
			return nil, err
		}
		risks = append(risks, risk)
	}

	return risks, rows.Err()
}

func (r *frameworkControlRepository) Create(ctx context.Context, input *models.CreateFrameworkControlInput) (*models.FrameworkControl, error) {
	control := &models.FrameworkControl{
		ID:          uuid.New().String(),
		FrameworkID: input.FrameworkID,
		ControlRef:  input.ControlRef,
		Title:       input.Title,
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `
		INSERT INTO framework_controls (id, framework_id, control_ref, title, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	if _, err := r.db.ExecContext(ctx, query,
		control.ID,
		control.FrameworkID,
		control.ControlRef,
		control.Title,
		control.Description,
		control.CreatedAt,
		control.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, control.ID)
}

func (r *frameworkControlRepository) Update(ctx context.Context, id string, input *models.UpdateFrameworkControlInput) (*models.FrameworkControl, error) {
	if input.ControlRef == nil && input.Title == nil && input.Description == nil {
		return nil, errors.New("at least one field must be updated")
	}

	query := `
		UPDATE framework_controls
		SET control_ref = COALESCE($1, control_ref),
			title = COALESCE($2, title),
			description = COALESCE($3, description),
			updated_at = NOW()
		WHERE id = $4
		RETURNING id
	`

	var updatedID string
	err := r.db.QueryRowContext(ctx, query, input.ControlRef, input.Title, input.Description, id).Scan(&updatedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFrameworkControlNotFound
		}
		return nil, err
	}

	return r.GetByID(ctx, updatedID)
}

func (r *frameworkControlRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM framework_controls WHERE id = $1`, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return ErrFrameworkControlInUse
		}
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrFrameworkControlNotFound
	}

	return nil
}

func (r *riskFrameworkControlRepository) ListByRiskID(ctx context.Context, riskID string) ([]*models.RiskFrameworkControl, error) {
	query := `
		SELECT rfc.id, rfc.risk_id, rfc.framework_control_id,
			fc.framework_id, f.name as framework_name, fc.control_ref, fc.title, COALESCE(fc.description, ''),
			rfc.notes, rfc.created_at, rfc.created_by
		FROM risk_framework_controls rfc
		JOIN framework_controls fc ON fc.id = rfc.framework_control_id
		JOIN frameworks f ON fc.framework_id = f.id
		WHERE rfc.risk_id = $1
		ORDER BY f.name ASC, fc.control_ref ASC
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
			&control.FrameworkControlID,
			&control.FrameworkID,
			&control.FrameworkName,
			&control.ControlRef,
			&control.ControlTitle,
			&control.ControlDescription,
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
		ID:                 uuid.New().String(),
		RiskID:             riskID,
		FrameworkControlID: input.FrameworkControlID,
		Notes:              input.Notes,
		CreatedAt:          time.Now(),
		CreatedBy:          createdBy,
	}

	query := `
		INSERT INTO risk_framework_controls (id, risk_id, framework_control_id, notes, created_at, created_by)
		SELECT $1, $2, fc.id, $3, $4, $5
		FROM framework_controls fc
		WHERE fc.id = $6
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(ctx, query,
		control.ID,
		control.RiskID,
		control.Notes,
		control.CreatedAt,
		control.CreatedBy,
		control.FrameworkControlID,
	).Scan(&control.ID, &control.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFrameworkControlNotFound
		}
		return nil, err
	}

	return r.getByID(ctx, control.ID)
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
		return ErrFrameworkControlNotFound
	}

	return nil
}

func (r *riskFrameworkControlRepository) getByID(ctx context.Context, id string) (*models.RiskFrameworkControl, error) {
	query := `
		SELECT rfc.id, rfc.risk_id, rfc.framework_control_id,
			fc.framework_id, f.name as framework_name, fc.control_ref, fc.title, COALESCE(fc.description, ''),
			rfc.notes, rfc.created_at, rfc.created_by
		FROM risk_framework_controls rfc
		JOIN framework_controls fc ON fc.id = rfc.framework_control_id
		JOIN frameworks f ON fc.framework_id = f.id
		WHERE rfc.id = $1
	`

	control := &models.RiskFrameworkControl{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&control.ID,
		&control.RiskID,
		&control.FrameworkControlID,
		&control.FrameworkID,
		&control.FrameworkName,
		&control.ControlRef,
		&control.ControlTitle,
		&control.ControlDescription,
		&control.Notes,
		&control.CreatedAt,
		&control.CreatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFrameworkControlNotFound
		}
		return nil, err
	}

	return control, nil
}
