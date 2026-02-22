package handlers

import (
	"backend/internal/database"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type AnalyticsHandler struct {
	repo database.AnalyticsRepository
}

func NewAnalyticsHandler(repo database.AnalyticsRepository) *AnalyticsHandler {
	return &AnalyticsHandler{repo: repo}
}

// Get returns all analytics data
func (h *AnalyticsHandler) Get(c *fiber.Ctx) error {
	granularity := models.AnalyticsGranularity(c.Query("granularity", "monthly"))
	if granularity != models.GranularityMonthly && granularity != models.GranularityWeekly {
		granularity = models.GranularityMonthly
	}

	response, err := h.repo.GetAnalytics(c.Context(), granularity)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch analytics"})
	}

	return c.JSON(response)
}
