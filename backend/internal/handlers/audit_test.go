package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func TestAuditHandler(t *testing.T) {
	app := fiber.New()
	mockRepo := &mockAuditRepo{logs: []*models.AuditLog{}}
	handler := NewAuditHandler(mockRepo)

	// Add some logs
	riskID := uuid.New().String()
	log1 := &models.AuditLog{
		ID:         uuid.New().String(),
		EntityType: "risk",
		EntityID:   riskID,
		Action:     models.AuditActionCreated,
		UserID:     "user1",
		CreatedAt:  time.Now(),
	}
	mockRepo.logs = append(mockRepo.logs, log1)

	// Add unrelated log
	log2 := &models.AuditLog{
		ID:         uuid.New().String(),
		EntityType: "risk",
		EntityID:   "other-risk",
		Action:     models.AuditActionCreated,
		UserID:     "user1",
		CreatedAt:  time.Now(),
	}
	mockRepo.logs = append(mockRepo.logs, log2)

	app.Get("/risks/:riskId/audit", handler.ListByRisk)

	t.Run("List Logs", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/risks/"+riskID+"/audit", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var response map[string][]*models.AuditLog
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		
		if len(response["data"]) != 1 {
			t.Errorf("expected 1 log, got %d", len(response["data"]))
		}
		if response["data"][0].ID != log1.ID {
			t.Errorf("expected log ID %s, got %s", log1.ID, response["data"][0].ID)
		}
	})
}
