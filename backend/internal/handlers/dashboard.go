package handlers

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	db *sql.DB
}

func NewDashboardHandler(db *sql.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// CategoryCount represents the count of risks per category
type CategoryCount struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	Count        int    `json:"count"`
}

// DashboardSummaryResponse represents the dashboard summary data
type DashboardSummaryResponse struct {
	TotalRisks     int             `json:"total_risks"`
	ByStatus       map[string]int  `json:"by_status"`
	BySeverity     map[string]int  `json:"by_severity"`
	ByCategory     []CategoryCount `json:"by_category"`
	OverdueReviews int             `json:"overdue_reviews"`
}

// ReviewRisk represents a risk with review date information
type ReviewRisk struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	ReviewDate time.Time `json:"review_date"`
	Severity   string    `json:"severity"`
	Status     string    `json:"status"`
}

// ReviewListResponse represents the response for review endpoints
type ReviewListResponse struct {
	Risks []ReviewRisk `json:"risks"`
}

// Summary returns aggregated risk statistics for the dashboard
func (h *DashboardHandler) Summary(c *fiber.Ctx) error {
	ctx := c.Context()

	response := DashboardSummaryResponse{
		ByStatus:   make(map[string]int),
		BySeverity: make(map[string]int),
		ByCategory: []CategoryCount{},
	}

	// Get total count
	err := h.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM risks").Scan(&response.TotalRisks)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch total risks"})
	}

	// Get counts by status
	rows, err := h.db.QueryContext(ctx, "SELECT status, COUNT(*) FROM risks GROUP BY status")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch risks by status"})
	}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			rows.Close()
			return c.Status(500).JSON(fiber.Map{"error": "failed to scan status row"})
		}
		response.ByStatus[status] = count
	}
	rows.Close()

	// Get counts by severity
	rows, err = h.db.QueryContext(ctx, "SELECT severity, COUNT(*) FROM risks GROUP BY severity")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch risks by severity"})
	}
	for rows.Next() {
		var severity string
		var count int
		if err := rows.Scan(&severity, &count); err != nil {
			rows.Close()
			return c.Status(500).JSON(fiber.Map{"error": "failed to scan severity row"})
		}
		response.BySeverity[severity] = count
	}
	rows.Close()

	// Get counts by category
	rows, err = h.db.QueryContext(ctx, `
		SELECT c.id, c.name, COUNT(r.id)
		FROM categories c
		LEFT JOIN risks r ON r.category_id = c.id
		GROUP BY c.id, c.name
		ORDER BY COUNT(r.id) DESC
	`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch risks by category"})
	}
	for rows.Next() {
		var cc CategoryCount
		if err := rows.Scan(&cc.CategoryID, &cc.CategoryName, &cc.Count); err != nil {
			rows.Close()
			return c.Status(500).JSON(fiber.Map{"error": "failed to scan category row"})
		}
		response.ByCategory = append(response.ByCategory, cc)
	}
	rows.Close()

	// Get overdue reviews count
	err = h.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM risks WHERE review_date IS NOT NULL AND review_date < NOW()",
	).Scan(&response.OverdueReviews)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch overdue reviews"})
	}

	return c.JSON(response)
}

// UpcomingReviews returns risks with review_date in the next N days
func (h *DashboardHandler) UpcomingReviews(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse days query parameter (default: 30)
	days := 30
	if daysParam := c.Query("days"); daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}

	query := `
		SELECT id, title, review_date, severity, status
		FROM risks
		WHERE review_date IS NOT NULL
			AND review_date >= NOW()
			AND review_date <= NOW() + INTERVAL '1 day' * $1
		ORDER BY review_date ASC
	`

	rows, err := h.db.QueryContext(ctx, query, days)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch upcoming reviews"})
	}
	defer rows.Close()

	response := ReviewListResponse{Risks: []ReviewRisk{}}
	for rows.Next() {
		var risk ReviewRisk
		if err := rows.Scan(&risk.ID, &risk.Title, &risk.ReviewDate, &risk.Severity, &risk.Status); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to scan review risk row"})
		}
		response.Risks = append(response.Risks, risk)
	}

	return c.JSON(response)
}

// OverdueReviews returns risks where review_date is in the past
func (h *DashboardHandler) OverdueReviews(c *fiber.Ctx) error {
	ctx := c.Context()

	query := `
		SELECT id, title, review_date, severity, status
		FROM risks
		WHERE review_date IS NOT NULL
			AND review_date < NOW()
		ORDER BY review_date ASC
	`

	rows, err := h.db.QueryContext(ctx, query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch overdue reviews"})
	}
	defer rows.Close()

	response := ReviewListResponse{Risks: []ReviewRisk{}}
	for rows.Next() {
		var risk ReviewRisk
		if err := rows.Scan(&risk.ID, &risk.Title, &risk.ReviewDate, &risk.Severity, &risk.Status); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to scan overdue review risk row"})
		}
		response.Risks = append(response.Risks, risk)
	}

	return c.JSON(response)
}
