package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"backend/internal/database"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type mockMitigationRepo struct {
	mitigations map[string]*models.Mitigation
}

func (m *mockMitigationRepo) Create(ctx context.Context, input *models.CreateMitigationInput, createdBy string) (*models.Mitigation, error) {
	mitigation := &models.Mitigation{
		ID:          uuid.New().String(),
		RiskID:      input.RiskID,
		Description: input.Description,
		Owner:       input.Owner,
		Status:      input.Status,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if input.DueDate != nil && *input.DueDate != "" {
		t, _ := time.Parse(time.RFC3339, *input.DueDate)
		mitigation.DueDate = &t
	}
	m.mitigations[mitigation.ID] = mitigation
	return mitigation, nil
}

func (m *mockMitigationRepo) FindByID(ctx context.Context, id string) (*models.Mitigation, error) {
	if val, ok := m.mitigations[id]; ok {
		return val, nil
	}
	return nil, database.ErrMitigationNotFound
}

func (m *mockMitigationRepo) ListByRiskID(ctx context.Context, riskID string) ([]*models.Mitigation, error) {
	var list []*models.Mitigation
	for _, val := range m.mitigations {
		if val.RiskID == riskID {
			list = append(list, val)
		}
	}
	return list, nil
}

func (m *mockMitigationRepo) Update(ctx context.Context, id string, input *models.UpdateMitigationInput, updatedBy string) (*models.Mitigation, error) {
	mitigation, ok := m.mitigations[id]
	if !ok {
		return nil, database.ErrMitigationNotFound
	}
	if input.Description != nil {
		mitigation.Description = *input.Description
	}
	if input.Owner != nil {
		mitigation.Owner = *input.Owner
	}
	if input.Status != nil {
		mitigation.Status = *input.Status
	}
	if input.DueDate != nil {
		if *input.DueDate == "" {
			mitigation.DueDate = nil
		} else {
			t, _ := time.Parse(time.RFC3339, *input.DueDate)
			mitigation.DueDate = &t
		}
	}
	mitigation.UpdatedBy = updatedBy
	mitigation.UpdatedAt = time.Now()
	m.mitigations[id] = mitigation
	return mitigation, nil
}

func (m *mockMitigationRepo) Delete(ctx context.Context, id string) error {
	if _, ok := m.mitigations[id]; !ok {
		return database.ErrMitigationNotFound
	}
	delete(m.mitigations, id)
	return nil
}

func TestListMitigationsHandler(t *testing.T) {
	app := fiber.New()
	mockRepo := &mockMitigationRepo{mitigations: make(map[string]*models.Mitigation)}
	handler := NewMitigationHandler(mockRepo)

	riskID := uuid.New().String()
	// Add test data
	for i := 0; i < 3; i++ {
		mit := &models.Mitigation{
			ID:          uuid.New().String(),
			RiskID:      riskID,
			Description: "Mitigation " + string(rune('A'+i)),
		}
		mockRepo.mitigations[mit.ID] = mit
	}

	app.Get("/risks/:riskId/mitigations", handler.List)

	req := httptest.NewRequest("GET", "/risks/"+riskID+"/mitigations", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var response []*models.Mitigation
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 3 {
		t.Errorf("expected 3 mitigations, got %d", len(response))
	}
}

func TestCreateMitigationHandler(t *testing.T) {
	app := fiber.New()
	mockRepo := &mockMitigationRepo{mitigations: make(map[string]*models.Mitigation)}
	handler := NewMitigationHandler(mockRepo)

	// Use testAuthMiddleware from risks_test.go
	app.Post("/risks/:riskId/mitigations", testAuthMiddleware, handler.Create)

	riskID := uuid.New().String()

	t.Run("Valid Input", func(t *testing.T) {
		input := models.CreateMitigationInput{
			Description: "New Mitigation",
			Owner:       "Owner",
			Status:      "planned",
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/risks/"+riskID+"/mitigations", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 201 {
			t.Errorf("expected status 201, got %d", resp.StatusCode)
		}

		var created models.Mitigation
		json.NewDecoder(resp.Body).Decode(&created)
		if created.Description != input.Description {
			t.Errorf("expected description %s, got %s", input.Description, created.Description)
		}
		if created.RiskID != riskID {
			t.Errorf("expected risk_id %s, got %s", riskID, created.RiskID)
		}
	})

	t.Run("Missing Description", func(t *testing.T) {
		input := models.CreateMitigationInput{
			Owner: "Owner",
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/risks/"+riskID+"/mitigations", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 400 {
			t.Errorf("expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestUpdateMitigationHandler(t *testing.T) {
	app := fiber.New()
	mockRepo := &mockMitigationRepo{mitigations: make(map[string]*models.Mitigation)}
	handler := NewMitigationHandler(mockRepo)

	riskID := uuid.New().String()
	mit := &models.Mitigation{
		ID:          uuid.New().String(),
		RiskID:      riskID,
		Description: "Old Description",
		Owner:       "Old Owner",
	}
	mockRepo.mitigations[mit.ID] = mit

	app.Put("/risks/:riskId/mitigations/:id", testAuthMiddleware, handler.Update)

	t.Run("Valid Update", func(t *testing.T) {
		newDesc := "New Description"
		input := models.UpdateMitigationInput{
			Description: &newDesc,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("PUT", "/risks/"+riskID+"/mitigations/"+mit.ID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var updated models.Mitigation
		json.NewDecoder(resp.Body).Decode(&updated)
		if updated.Description != newDesc {
			t.Errorf("expected description %s, got %s", newDesc, updated.Description)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		newDesc := "New Description"
		input := models.UpdateMitigationInput{
			Description: &newDesc,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("PUT", "/risks/"+riskID+"/mitigations/invalid-id", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 404 {
			t.Errorf("expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestDeleteMitigationHandler(t *testing.T) {
	app := fiber.New()
	mockRepo := &mockMitigationRepo{mitigations: make(map[string]*models.Mitigation)}
	handler := NewMitigationHandler(mockRepo)

	riskID := uuid.New().String()
	mit := &models.Mitigation{
		ID:          uuid.New().String(),
		RiskID:      riskID,
		Description: "To Delete",
	}
	mockRepo.mitigations[mit.ID] = mit

	app.Delete("/risks/:riskId/mitigations/:id", testAuthMiddleware, handler.Delete)

	t.Run("Valid Delete", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/risks/"+riskID+"/mitigations/"+mit.ID, nil)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 204 {
			t.Errorf("expected status 204, got %d", resp.StatusCode)
		}

		if _, exists := mockRepo.mitigations[mit.ID]; exists {
			t.Error("mitigation should have been deleted")
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/risks/"+riskID+"/mitigations/invalid-id", nil)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 404 {
			t.Errorf("expected status 404, got %d", resp.StatusCode)
		}
	})
}
