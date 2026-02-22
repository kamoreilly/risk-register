package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAIHandler(t *testing.T) {
	app := fiber.New()
	handler := NewAIHandler()

	app.Post("/ai/summarize", handler.Summarize)
	app.Post("/ai/draft-mitigation", handler.DraftMitigation)

	t.Run("Summarize", func(t *testing.T) {
		reqBody := map[string]string{
			"title":       "High Risk",
			"description": "Description",
			"severity":    "high",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/ai/summarize", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var response SummarizeResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if response.Summary == "" {
			t.Error("expected summary to be non-empty")
		}
	})

	t.Run("Draft Mitigation", func(t *testing.T) {
		reqBody := map[string]string{
			"risk_title":       "Risk",
			"risk_description": "Desc",
			"severity":         "medium",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/ai/draft-mitigation", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var response DraftMitigationResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if response.Draft == "" {
			t.Error("expected draft to be non-empty")
		}
	})
}
