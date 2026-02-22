package handlers

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type mockDashboardRepo struct {
	summary  *models.DashboardSummaryResponse
	upcoming *models.ReviewListResponse
	overdue  *models.ReviewListResponse
}

func (m *mockDashboardRepo) GetSummary(ctx context.Context) (*models.DashboardSummaryResponse, error) {
	return m.summary, nil
}

func (m *mockDashboardRepo) GetUpcomingReviews(ctx context.Context, days int) (*models.ReviewListResponse, error) {
	return m.upcoming, nil
}

func (m *mockDashboardRepo) GetOverdueReviews(ctx context.Context) (*models.ReviewListResponse, error) {
	return m.overdue, nil
}

func TestDashboardHandler(t *testing.T) {
	app := fiber.New()
	mockRepo := &mockDashboardRepo{
		summary: &models.DashboardSummaryResponse{
			TotalRisks: 10,
			ByStatus:   map[string]int{"open": 5, "closed": 5},
		},
		upcoming: &models.ReviewListResponse{
			Risks: []models.ReviewRisk{{ID: "1", Title: "Upcoming"}},
		},
		overdue: &models.ReviewListResponse{
			Risks: []models.ReviewRisk{{ID: "2", Title: "Overdue"}},
		},
	}
	handler := NewDashboardHandler(mockRepo)

	// Use testAuthMiddleware if protected routes, but unit tests can skip if we don't apply middleware in test app setup
	// The actual routes in server.go are protected.
	// But here we just test the handler logic.
	app.Get("/dashboard/summary", handler.Summary)
	app.Get("/dashboard/reviews/upcoming", handler.UpcomingReviews)
	app.Get("/dashboard/reviews/overdue", handler.OverdueReviews)

	t.Run("Summary", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/dashboard/summary", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var response models.DashboardSummaryResponse
		json.NewDecoder(resp.Body).Decode(&response)
		if response.TotalRisks != 10 {
			t.Errorf("expected 10 total risks, got %d", response.TotalRisks)
		}
	})

	t.Run("Upcoming Reviews", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/dashboard/reviews/upcoming?days=30", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var response models.ReviewListResponse
		json.NewDecoder(resp.Body).Decode(&response)
		if len(response.Risks) != 1 {
			t.Errorf("expected 1 upcoming review, got %d", len(response.Risks))
		}
	})

	t.Run("Overdue Reviews", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/dashboard/reviews/overdue", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var response models.ReviewListResponse
		json.NewDecoder(resp.Body).Decode(&response)
		if len(response.Risks) != 1 {
			t.Errorf("expected 1 overdue review, got %d", len(response.Risks))
		}
	})
}
