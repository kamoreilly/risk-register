package database

import (
	"context"
	"database/sql"
	"errors"

	"backend/internal/models"
)

var ErrIncidentCategoryNotFound = errors.New("incident category not found")

type IncidentCategoryRepository interface {
	List(ctx context.Context) ([]*models.IncidentCategory, error)
	FindByID(ctx context.Context, id string) (*models.IncidentCategory, error)
	Create(ctx context.Context, input *models.CreateIncidentCategoryInput) (*models.IncidentCategory, error)
	Update(ctx context.Context, id string, input *models.UpdateIncidentCategoryInput) (*models.IncidentCategory, error)
	Delete(ctx context.Context, id string) error
}

type incidentCategoryRepository struct {
	db *sql.DB
}

func NewIncidentCategoryRepository(db *sql.DB) IncidentCategoryRepository {
	return &incidentCategoryRepository{db: db}
}

func (r *incidentCategoryRepository) List(ctx context.Context) ([]*models.IncidentCategory, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM incident_categories ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.IncidentCategory
	for rows.Next() {
		c := &models.IncidentCategory{}
		var desc sql.NullString
		if err := rows.Scan(&c.ID, &c.Name, &desc, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		if desc.Valid {
			c.Description = desc.String
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *incidentCategoryRepository) FindByID(ctx context.Context, id string) (*models.IncidentCategory, error) {
	c := &models.IncidentCategory{}
	var desc sql.NullString
	query := `SELECT id, name, description, created_at, updated_at FROM incident_categories WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.Name, &desc, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if desc.Valid {
		c.Description = desc.String
	}
	return c, nil
}

func (r *incidentCategoryRepository) Create(ctx context.Context, input *models.CreateIncidentCategoryInput) (*models.IncidentCategory, error) {
	c := &models.IncidentCategory{}
	query := `
		INSERT INTO incident_categories (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at
	`
	var desc sql.NullString
	err := r.db.QueryRowContext(ctx, query, input.Name, sql.NullString{String: input.Description, Valid: input.Description != ""}).Scan(
		&c.ID, &c.Name, &desc, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if desc.Valid {
		c.Description = desc.String
	}
	return c, nil
}

func (r *incidentCategoryRepository) Update(ctx context.Context, id string, input *models.UpdateIncidentCategoryInput) (*models.IncidentCategory, error) {
	if input.Name == nil && input.Description == nil {
		return nil, errors.New("at least one field must be updated")
	}

	c := &models.IncidentCategory{}
	query := `
		UPDATE incident_categories
		SET name = COALESCE($1, name), description = COALESCE($2, description), updated_at = NOW()
		WHERE id = $3
		RETURNING id, name, description, created_at, updated_at
	`
	var desc sql.NullString
	err := r.db.QueryRowContext(ctx, query, input.Name, input.Description, id).Scan(
		&c.ID, &c.Name, &desc, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrIncidentCategoryNotFound
		}
		return nil, err
	}
	if desc.Valid {
		c.Description = desc.String
	}
	return c, nil
}

func (r *incidentCategoryRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM incident_categories WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrIncidentCategoryNotFound
	}
	return nil
}
