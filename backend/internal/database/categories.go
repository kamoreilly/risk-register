package database

import (
	"context"
	"database/sql"
	"errors"

	"backend/internal/models"
)

var ErrCategoryNotFound = errors.New("category not found")

type CategoryRepository interface {
	List(ctx context.Context) ([]*models.Category, error)
	FindByID(ctx context.Context, id string) (*models.Category, error)
	Create(ctx context.Context, input *models.CreateCategoryInput) (*models.Category, error)
	Update(ctx context.Context, id string, input *models.UpdateCategoryInput) (*models.Category, error)
	Delete(ctx context.Context, id string) error
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) List(ctx context.Context) ([]*models.Category, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM categories ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		c := &models.Category{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *categoryRepository) FindByID(ctx context.Context, id string) (*models.Category, error) {
	c := &models.Category{}
	query := `SELECT id, name, description, created_at, updated_at FROM categories WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return c, nil
}

func (r *categoryRepository) Create(ctx context.Context, input *models.CreateCategoryInput) (*models.Category, error) {
	c := &models.Category{}
	query := `
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, input.Name, input.Description).Scan(
		&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *categoryRepository) Update(ctx context.Context, id string, input *models.UpdateCategoryInput) (*models.Category, error) {
	if input.Name == nil && input.Description == nil {
		return nil, errors.New("at least one field must be updated")
	}

	c := &models.Category{}
	query := `
		UPDATE categories
		SET name = COALESCE($1, name), description = COALESCE($2, description), updated_at = NOW()
		WHERE id = $3
		RETURNING id, name, description, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, input.Name, input.Description, id).Scan(
		&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return c, nil
}

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM categories WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrCategoryNotFound
	}
	return nil
}
