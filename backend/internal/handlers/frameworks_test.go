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

type mockFrameworkRepo struct {
	frameworks map[string]*models.Framework
}

func (m *mockFrameworkRepo) List(ctx context.Context) ([]*models.Framework, error) {
	var list []*models.Framework
	for _, f := range m.frameworks {
		list = append(list, f)
	}
	return list, nil
}

func (m *mockFrameworkRepo) GetByID(ctx context.Context, id string) (*models.Framework, error) {
	if f, ok := m.frameworks[id]; ok {
		return f, nil
	}
	return nil, database.ErrFrameworkNotFound
}

func (m *mockFrameworkRepo) Create(ctx context.Context, input *models.CreateFrameworkInput) (*models.Framework, error) {
	framework := &models.Framework{
		ID:          uuid.New().String(),
		Name:        input.Name,
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.frameworks[framework.ID] = framework
	return framework, nil
}

func (m *mockFrameworkRepo) Update(ctx context.Context, id string, input *models.UpdateFrameworkInput) (*models.Framework, error) {
	framework, ok := m.frameworks[id]
	if !ok {
		return nil, database.ErrFrameworkNotFound
	}
	if input.Name != nil {
		framework.Name = *input.Name
	}
	if input.Description != nil {
		framework.Description = *input.Description
	}
	framework.UpdatedAt = time.Now()
	m.frameworks[id] = framework
	return framework, nil
}

func (m *mockFrameworkRepo) Delete(ctx context.Context, id string) error {
	if _, ok := m.frameworks[id]; !ok {
		return database.ErrFrameworkNotFound
	}
	delete(m.frameworks, id)
	return nil
}

type mockControlRepo struct {
	controls      map[string]*models.RiskFrameworkControl
	frameworkRepo *mockFrameworkRepo
}

func (m *mockControlRepo) ListByRiskID(ctx context.Context, riskID string) ([]*models.RiskFrameworkControl, error) {
	var list []*models.RiskFrameworkControl
	for _, c := range m.controls {
		if c.RiskID == riskID {
			list = append(list, c)
		}
	}
	return list, nil
}

func (m *mockControlRepo) LinkControl(ctx context.Context, riskID string, input *models.LinkControlInput, createdBy string) (*models.RiskFrameworkControl, error) {
	// Verify framework exists
	framework, err := m.frameworkRepo.GetByID(ctx, input.FrameworkID)
	if err != nil {
		// Real repo would fail foreign key constraint or we check it manually
		// For mock, let's allow it or fail if we want strictness.
		// The handler doesn't check framework existence explicitly, it relies on repo.
		// Let's assume input is valid or just use a default name if not found.
	}
	frameworkName := "Unknown Framework"
	if framework != nil {
		frameworkName = framework.Name
	}

	control := &models.RiskFrameworkControl{
		ID:            uuid.New().String(),
		RiskID:        riskID,
		FrameworkID:   input.FrameworkID,
		FrameworkName: frameworkName,
		ControlRef:    input.ControlRef,
		Notes:         input.Notes,
		CreatedBy:     createdBy,
		CreatedAt:     time.Now(),
	}
	m.controls[control.ID] = control
	return control, nil
}

func (m *mockControlRepo) UnlinkControl(ctx context.Context, id string) error {
	if _, ok := m.controls[id]; !ok {
		return database.ErrFrameworkNotFound // Handler checks for this error
	}
	delete(m.controls, id)
	return nil
}

func TestFrameworkHandler(t *testing.T) {
	app := fiber.New()
	mockFwRepo := &mockFrameworkRepo{frameworks: make(map[string]*models.Framework)}
	mockCtrlRepo := &mockControlRepo{
		controls:      make(map[string]*models.RiskFrameworkControl),
		frameworkRepo: mockFwRepo,
	}
	handler := NewFrameworkHandler(mockFwRepo, mockCtrlRepo)

	// Setup routes
	app.Get("/frameworks", testAuthMiddleware, handler.List)
	app.Post("/frameworks", testAuthMiddleware, handler.Create)
	app.Put("/frameworks/:id", testAuthMiddleware, handler.Update)
	app.Delete("/frameworks/:id", testAuthMiddleware, handler.Delete)

	t.Run("Create Framework", func(t *testing.T) {
		input := models.CreateFrameworkInput{
			Name:        "ISO 27001",
			Description: "Information Security",
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/frameworks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 201 {
			t.Errorf("expected status 201, got %d", resp.StatusCode)
		}

		var created models.Framework
		json.NewDecoder(resp.Body).Decode(&created)
		if created.Name != input.Name {
			t.Errorf("expected name %s, got %s", input.Name, created.Name)
		}
		mockFwRepo.frameworks[created.ID] = &created
	})

	t.Run("List Frameworks", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/frameworks", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		
		var response map[string][]*models.Framework
		json.NewDecoder(resp.Body).Decode(&response)
		if len(response["data"]) != 1 {
			t.Errorf("expected 1 framework, got %d", len(response["data"]))
		}
	})
}

func TestControlHandler(t *testing.T) {
	app := fiber.New()
	mockFwRepo := &mockFrameworkRepo{frameworks: make(map[string]*models.Framework)}
	mockCtrlRepo := &mockControlRepo{
		controls:      make(map[string]*models.RiskFrameworkControl),
		frameworkRepo: mockFwRepo,
	}
	handler := NewControlHandler(mockCtrlRepo)

	// Setup routes
	app.Get("/risks/:riskId/controls", handler.ListControls)
	app.Post("/risks/:riskId/controls", testAuthMiddleware, handler.LinkControl)
	app.Delete("/risks/:riskId/controls/:id", testAuthMiddleware, handler.UnlinkControl)

	// Create a framework first
	fw := &models.Framework{ID: uuid.New().String(), Name: "NIST"}
	mockFwRepo.frameworks[fw.ID] = fw

	riskID := uuid.New().String()

	t.Run("Link Control", func(t *testing.T) {
		input := models.LinkControlInput{
			FrameworkID: fw.ID,
			ControlRef:  "AC-1",
			Notes:       "Access Control Policy",
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/risks/"+riskID+"/controls", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 201 {
			t.Errorf("expected status 201, got %d", resp.StatusCode)
		}

		var linked models.RiskFrameworkControl
		json.NewDecoder(resp.Body).Decode(&linked)
		if linked.ControlRef != input.ControlRef {
			t.Errorf("expected control ref %s, got %s", input.ControlRef, linked.ControlRef)
		}
		if linked.FrameworkName != fw.Name {
			t.Errorf("expected framework name %s, got %s", fw.Name, linked.FrameworkName)
		}
	})

	t.Run("List Controls", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/risks/"+riskID+"/controls", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})
}
