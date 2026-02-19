package handlers

import (
	"database/sql"

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
