package database

import (
	"context"
	"database/sql"

	"backend/internal/models"
)

type DashboardRepository interface {
	GetSummary(ctx context.Context) (*models.DashboardSummaryResponse, error)
	GetUpcomingReviews(ctx context.Context, days int) (*models.ReviewListResponse, error)
	GetOverdueReviews(ctx context.Context) (*models.ReviewListResponse, error)
}

type dashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetSummary(ctx context.Context) (*models.DashboardSummaryResponse, error) {
	response := &models.DashboardSummaryResponse{
		ByStatus:   make(map[string]int),
		BySeverity: make(map[string]int),
		ByCategory: []models.CategoryCount{},
	}

	// Get total count
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM risks").Scan(&response.TotalRisks)
	if err != nil {
		return nil, err
	}

	// Get counts by status
	rows, err := r.db.QueryContext(ctx, "SELECT status, COUNT(*) FROM risks GROUP BY status")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			rows.Close()
			return nil, err
		}
		response.ByStatus[status] = count
	}
	rows.Close()

	// Get counts by severity
	rows, err = r.db.QueryContext(ctx, "SELECT severity, COUNT(*) FROM risks GROUP BY severity")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var severity string
		var count int
		if err := rows.Scan(&severity, &count); err != nil {
			rows.Close()
			return nil, err
		}
		response.BySeverity[severity] = count
	}
	rows.Close()

	// Get counts by category
	rows, err = r.db.QueryContext(ctx, `
		SELECT c.id, c.name, COUNT(r.id)
		FROM categories c
		LEFT JOIN risks r ON r.category_id = c.id
		GROUP BY c.id, c.name
		ORDER BY COUNT(r.id) DESC
	`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var cc models.CategoryCount
		if err := rows.Scan(&cc.CategoryID, &cc.CategoryName, &cc.Count); err != nil {
			rows.Close()
			return nil, err
		}
		response.ByCategory = append(response.ByCategory, cc)
	}
	rows.Close()

	// Get overdue reviews count
	err = r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM risks WHERE review_date IS NOT NULL AND review_date < NOW()",
	).Scan(&response.OverdueReviews)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *dashboardRepository) GetUpcomingReviews(ctx context.Context, days int) (*models.ReviewListResponse, error) {
	query := `
		SELECT id, title, review_date, severity, status
		FROM risks
		WHERE review_date IS NOT NULL
			AND review_date >= NOW()
			AND review_date <= NOW() + INTERVAL '1 day' * $1
		ORDER BY review_date ASC
	`

	rows, err := r.db.QueryContext(ctx, query, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	response := &models.ReviewListResponse{Risks: []models.ReviewRisk{}}
	for rows.Next() {
		var risk models.ReviewRisk
		if err := rows.Scan(&risk.ID, &risk.Title, &risk.ReviewDate, &risk.Severity, &risk.Status); err != nil {
			return nil, err
		}
		response.Risks = append(response.Risks, risk)
	}

	return response, nil
}

func (r *dashboardRepository) GetOverdueReviews(ctx context.Context) (*models.ReviewListResponse, error) {
	query := `
		SELECT id, title, review_date, severity, status
		FROM risks
		WHERE review_date IS NOT NULL
			AND review_date < NOW()
		ORDER BY review_date ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	response := &models.ReviewListResponse{Risks: []models.ReviewRisk{}}
	for rows.Next() {
		var risk models.ReviewRisk
		if err := rows.Scan(&risk.ID, &risk.Title, &risk.ReviewDate, &risk.Severity, &risk.Status); err != nil {
			return nil, err
		}
		response.Risks = append(response.Risks, risk)
	}

	return response, nil
}
