package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"backend/internal/models"

	"github.com/google/uuid"
)

var ErrRiskNotFound = errors.New("risk not found")

type RiskRepository interface {
	Create(ctx context.Context, risk *models.Risk) error
	FindByID(ctx context.Context, id string) (*models.Risk, error)
	List(ctx context.Context, params *models.RiskListParams) (*models.RiskListResponse, error)
	Update(ctx context.Context, risk *models.Risk) error
	Delete(ctx context.Context, id string) error
}

type riskRepository struct {
	db *sql.DB
}

func NewRiskRepository(db *sql.DB) RiskRepository {
	return &riskRepository{db: db}
}

func (r *riskRepository) Create(ctx context.Context, risk *models.Risk) error {
	if risk.ID == "" {
		risk.ID = uuid.New().String()
	}
	now := time.Now()
	risk.CreatedAt = now
	risk.UpdatedAt = now

	query := `
		INSERT INTO risks (id, title, description, owner_id, status, severity, category_id, review_date, created_at, updated_at, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		risk.ID, risk.Title, risk.Description, risk.OwnerID, risk.Status, risk.Severity,
		risk.CategoryID, risk.ReviewDate, risk.CreatedAt, risk.UpdatedAt, risk.CreatedBy, risk.UpdatedBy,
	).Scan(&risk.ID, &risk.CreatedAt, &risk.UpdatedAt)
}

func (r *riskRepository) FindByID(ctx context.Context, id string) (*models.Risk, error) {
	query := `
		SELECT id, title, description, owner_id, status, severity, category_id, review_date, created_at, updated_at, created_by, updated_by
		FROM risks WHERE id = $1
	`
	risk := &models.Risk{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&risk.ID, &risk.Title, &risk.Description, &risk.OwnerID, &risk.Status, &risk.Severity,
		&risk.CategoryID, &risk.ReviewDate, &risk.CreatedAt, &risk.UpdatedAt, &risk.CreatedBy, &risk.UpdatedBy,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRiskNotFound
		}
		return nil, err
	}
	return risk, nil
}

func (r *riskRepository) List(ctx context.Context, params *models.RiskListParams) (*models.RiskListResponse, error) {
	if params == nil {
		params = &models.RiskListParams{Page: 1, Limit: 20}
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}

	// Build WHERE clause
	where := "WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if params.Status != nil {
		where += fmt.Sprintf(" AND status = $%d", argNum)
		args = append(args, *params.Status)
		argNum++
	}
	if params.Severity != nil {
		where += fmt.Sprintf(" AND severity = $%d", argNum)
		args = append(args, *params.Severity)
		argNum++
	}
	if params.CategoryID != nil {
		where += fmt.Sprintf(" AND category_id = $%d", argNum)
		args = append(args, *params.CategoryID)
		argNum++
	}
	if params.OwnerID != nil {
		where += fmt.Sprintf(" AND owner_id = $%d", argNum)
		args = append(args, *params.OwnerID)
		argNum++
	}
	if params.Search != "" {
		where += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argNum, argNum)
		args = append(args, "%"+params.Search+"%", "%"+params.Search+"%")
		argNum++
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM risks %s", where)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Build ORDER BY
	orderBy := "created_at"
	if params.Sort != "" {
		orderBy = params.Sort
	}
	orderDir := "DESC"
	if params.Order == "asc" {
		orderDir = "ASC"
	}

	// Get paginated results
	offset := (params.Page - 1) * params.Limit
	query := fmt.Sprintf(`
		SELECT id, title, description, owner_id, status, severity, category_id, review_date, created_at, updated_at, created_by, updated_by
		FROM risks %s ORDER BY %s %s LIMIT $%d OFFSET $%d
	`, where, orderBy, orderDir, argNum, argNum+1)
	args = append(args, params.Limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var risks []*models.Risk
	for rows.Next() {
		risk := &models.Risk{}
		err := rows.Scan(
			&risk.ID, &risk.Title, &risk.Description, &risk.OwnerID, &risk.Status, &risk.Severity,
			&risk.CategoryID, &risk.ReviewDate, &risk.CreatedAt, &risk.UpdatedAt, &risk.CreatedBy, &risk.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		risks = append(risks, risk)
	}

	return &models.RiskListResponse{
		Data: risks,
		Meta: models.Meta{Page: params.Page, Limit: params.Limit, Total: total},
	}, rows.Err()
}

func (r *riskRepository) Update(ctx context.Context, risk *models.Risk) error {
	now := time.Now()
	risk.UpdatedAt = now

	query := `
		UPDATE risks SET title = $1, description = $2, owner_id = $3, status = $4, severity = $5,
			category_id = $6, review_date = $7, updated_at = $8, updated_by = $9
		WHERE id = $10
		RETURNING updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		risk.Title, risk.Description, risk.OwnerID, risk.Status, risk.Severity,
		risk.CategoryID, risk.ReviewDate, risk.UpdatedAt, risk.UpdatedBy, risk.ID,
	).Scan(&risk.UpdatedAt)
}

func (r *riskRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM risks WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrRiskNotFound
	}
	return nil
}
