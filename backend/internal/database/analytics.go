package database

import (
	"context"
	"database/sql"
	"fmt"
	"sort"

	"backend/internal/models"
)

type AnalyticsRepository interface {
	GetAnalytics(ctx context.Context, granularity models.AnalyticsGranularity) (*models.AnalyticsResponse, error)
}

type analyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) AnalyticsRepository {
	return &analyticsRepository{db: db}
}

func (r *analyticsRepository) GetAnalytics(ctx context.Context, granularity models.AnalyticsGranularity) (*models.AnalyticsResponse, error) {
	response := &models.AnalyticsResponse{
		BySeverity:      make(map[string]int),
		ByStatus:        make(map[string]int),
		ByCategory:      []models.CategoryCount{},
		CreatedOverTime: []models.TimeDataPoint{},
		StatusOverTime:  []models.StatusTimeDataPoint{},
	}

	// Get total count
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM risks").Scan(&response.TotalRisks); err != nil {
		return nil, fmt.Errorf("failed to get total risks: %w", err)
	}

	// Get counts by severity
	if err := r.populateByField(ctx, "severity", &response.BySeverity); err != nil {
		return nil, err
	}

	// Get counts by status
	if err := r.populateByField(ctx, "status", &response.ByStatus); err != nil {
		return nil, err
	}

	// Get counts by category
	if err := r.populateByCategory(ctx, &response.ByCategory); err != nil {
		return nil, err
	}

	// Get created over time
	if err := r.populateCreatedOverTime(ctx, granularity, &response.CreatedOverTime); err != nil {
		return nil, err
	}

	// Get status over time
	if err := r.populateStatusOverTime(ctx, granularity, &response.StatusOverTime); err != nil {
		return nil, err
	}

	return response, nil
}

func (r *analyticsRepository) populateByField(ctx context.Context, field string, target *map[string]int) error {
	// Allowlist for valid field names to prevent SQL injection
	var query string
	switch field {
	case "severity":
		query = "SELECT severity, COUNT(*) FROM risks GROUP BY severity"
	case "status":
		query = "SELECT status, COUNT(*) FROM risks GROUP BY status"
	default:
		return fmt.Errorf("invalid field for grouping: %s", field)
	}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to get counts by %s: %w", field, err)
	}
	defer rows.Close()

	for rows.Next() {
		var value string
		var count int
		if err := rows.Scan(&value, &count); err != nil {
			return fmt.Errorf("failed to scan %s row: %w", field, err)
		}
		(*target)[value] = count
	}
	return rows.Err()
}

func (r *analyticsRepository) populateByCategory(ctx context.Context, target *[]models.CategoryCount) error {
	query := `
		SELECT c.id, c.name, COUNT(r.id)
		FROM categories c
		LEFT JOIN risks r ON r.category_id = c.id
		GROUP BY c.id, c.name
		ORDER BY COUNT(r.id) DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to get counts by category: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cc models.CategoryCount
		if err := rows.Scan(&cc.CategoryID, &cc.CategoryName, &cc.Count); err != nil {
			return fmt.Errorf("failed to scan category row: %w", err)
		}
		*target = append(*target, cc)
	}
	return rows.Err()
}

func (r *analyticsRepository) populateCreatedOverTime(ctx context.Context, granularity models.AnalyticsGranularity, target *[]models.TimeDataPoint) error {
	var dateFormat string
	if granularity == models.GranularityWeekly {
		dateFormat = "YYYY-\"W\"WW"
	} else {
		dateFormat = "YYYY-MM"
	}

	query := fmt.Sprintf(`
		SELECT TO_CHAR(created_at, '%s') as period, COUNT(*) as count
		FROM risks
		WHERE created_at >= NOW() - INTERVAL '12 months'
		GROUP BY period
		ORDER BY period ASC
	`, dateFormat)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to get created over time: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var dp models.TimeDataPoint
		if err := rows.Scan(&dp.Period, &dp.Count); err != nil {
			return fmt.Errorf("failed to scan created over time row: %w", err)
		}
		*target = append(*target, dp)
	}
	return rows.Err()
}

func (r *analyticsRepository) populateStatusOverTime(ctx context.Context, granularity models.AnalyticsGranularity, target *[]models.StatusTimeDataPoint) error {
	var dateFormat string
	if granularity == models.GranularityWeekly {
		dateFormat = "YYYY-\"W\"WW"
	} else {
		dateFormat = "YYYY-MM"
	}

	// Get opened (created) per period
	openedQuery := fmt.Sprintf(`
		SELECT TO_CHAR(created_at, '%s') as period, COUNT(*) as count
		FROM risks
		WHERE created_at >= NOW() - INTERVAL '12 months'
		GROUP BY period
	`, dateFormat)

	openedCounts := make(map[string]int)
	rows, err := r.db.QueryContext(ctx, openedQuery)
	if err != nil {
		return fmt.Errorf("failed to get opened counts: %w", err)
	}
	for rows.Next() {
		var period string
		var count int
		if err := rows.Scan(&period, &count); err != nil {
			rows.Close()
			return err
		}
		openedCounts[period] = count
	}
	rows.Close()

	// Get closed (status changed to resolved/accepted) per period
	closedQuery := fmt.Sprintf(`
		SELECT TO_CHAR(updated_at, '%s') as period, COUNT(*) as count
		FROM risks
		WHERE status IN ('resolved', 'accepted')
		  AND updated_at >= NOW() - INTERVAL '12 months'
		GROUP BY period
	`, dateFormat)

	closedCounts := make(map[string]int)
	rows, err = r.db.QueryContext(ctx, closedQuery)
	if err != nil {
		return fmt.Errorf("failed to get closed counts: %w", err)
	}
	for rows.Next() {
		var period string
		var count int
		if err := rows.Scan(&period, &count); err != nil {
			rows.Close()
			return err
		}
		closedCounts[period] = count
	}
	rows.Close()

	// Merge all periods
	allPeriods := make(map[string]bool)
	for p := range openedCounts {
		allPeriods[p] = true
	}
	for p := range closedCounts {
		allPeriods[p] = true
	}

	// Build sorted result
	for period := range allPeriods {
		*target = append(*target, models.StatusTimeDataPoint{
			Period: period,
			Open:   openedCounts[period],
			Closed: closedCounts[period],
		})
	}

	// Sort by period
	sort.Slice(*target, func(i, j int) bool {
		return (*target)[i].Period < (*target)[j].Period
	})

	return nil
}
