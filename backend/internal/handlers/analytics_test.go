package handlers

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type mockAnalyticsRepo struct {
	response *models.AnalyticsResponse
	err      error
}

func (m *mockAnalyticsRepo) GetAnalytics(ctx context.Context, granularity models.AnalyticsGranularity) (*models.AnalyticsResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

func TestAnalyticsHandler_Get(t *testing.T) {
	mockResp := &models.AnalyticsResponse{
		TotalRisks: 10,
		BySeverity: map[string]int{"critical": 2, "high": 3, "medium": 3, "low": 2},
		ByStatus:   map[string]int{"open": 6, "mitigating": 2, "resolved": 2},
		ByCategory: []models.CategoryCount{
			{CategoryID: "1", CategoryName: "Security", Count: 5},
			{CategoryID: "2", CategoryName: "Compliance", Count: 3},
		},
		CreatedOverTime: []models.TimeDataPoint{
			{Period: "2024-01", Count: 2},
			{Period: "2024-02", Count: 3},
		},
		StatusOverTime: []models.StatusTimeDataPoint{
			{Period: "2024-01", Open: 2, Closed: 1},
			{Period: "2024-02", Open: 3, Closed: 2},
		},
	}

	app := fiber.New()
	handler := NewAnalyticsHandler(&mockAnalyticsRepo{response: mockResp})
	app.Get("/analytics", handler.Get)

	tests := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{"default granularity", "/analytics", 200},
		{"monthly granularity", "/analytics?granularity=monthly", 200},
		{"weekly granularity", "/analytics?granularity=weekly", 200},
		{"invalid granularity defaults to monthly", "/analytics?granularity=invalid", 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			var response models.AnalyticsResponse
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.TotalRisks != 10 {
				t.Errorf("expected 10 total risks, got %d", response.TotalRisks)
			}
		})
	}
}
