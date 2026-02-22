package handlers

import (
	"strconv"

	"backend/internal/database"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	repo database.DashboardRepository
}

func NewDashboardHandler(repo database.DashboardRepository) *DashboardHandler {
	return &DashboardHandler{repo: repo}
}

// Summary returns aggregated risk statistics for the dashboard
func (h *DashboardHandler) Summary(c *fiber.Ctx) error {
	response, err := h.repo.GetSummary(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch dashboard summary"})
	}

	return c.JSON(response)
}

// UpcomingReviews returns risks with review_date in the next N days
func (h *DashboardHandler) UpcomingReviews(c *fiber.Ctx) error {
	// Parse days query parameter (default: 30)
	days := 30
	if daysParam := c.Query("days"); daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}

	response, err := h.repo.GetUpcomingReviews(c.Context(), days)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch upcoming reviews"})
	}

	return c.JSON(response)
}

// OverdueReviews returns risks where review_date is in the past
func (h *DashboardHandler) OverdueReviews(c *fiber.Ctx) error {
	response, err := h.repo.GetOverdueReviews(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch overdue reviews"})
	}

	return c.JSON(response)
}
