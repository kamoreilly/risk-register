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

// mockIncidentCategoryRepo implements database.IncidentCategoryRepository for testing
type mockIncidentCategoryRepo struct {
	categories map[string]*models.IncidentCategory
}

func newMockIncidentCategoryRepo() *mockIncidentCategoryRepo {
	return &mockIncidentCategoryRepo{categories: make(map[string]*models.IncidentCategory)}
}

func (m *mockIncidentCategoryRepo) List(ctx context.Context) ([]*models.IncidentCategory, error) {
	var categories []*models.IncidentCategory
	for _, category := range m.categories {
		categories = append(categories, category)
	}
	return categories, nil
}

func (m *mockIncidentCategoryRepo) FindByID(ctx context.Context, id string) (*models.IncidentCategory, error) {
	if category, ok := m.categories[id]; ok {
		return category, nil
	}
	return nil, nil
}

func (m *mockIncidentCategoryRepo) Create(ctx context.Context, input *models.CreateIncidentCategoryInput) (*models.IncidentCategory, error) {
	category := &models.IncidentCategory{
		ID:          uuid.New().String(),
		Name:        input.Name,
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.categories[category.ID] = category
	return category, nil
}

func (m *mockIncidentCategoryRepo) Update(ctx context.Context, id string, input *models.UpdateIncidentCategoryInput) (*models.IncidentCategory, error) {
	category, ok := m.categories[id]
	if !ok {
		return nil, database.ErrIncidentCategoryNotFound
	}

	if input.Name != nil {
		category.Name = *input.Name
	}
	if input.Description != nil {
		category.Description = *input.Description
	}
	category.UpdatedAt = time.Now()

	return category, nil
}

func (m *mockIncidentCategoryRepo) Delete(ctx context.Context, id string) error {
	if _, ok := m.categories[id]; !ok {
		return database.ErrIncidentCategoryNotFound
	}
	delete(m.categories, id)
	return nil
}

func TestListIncidentCategoriesHandler(t *testing.T) {
	app := fiber.New()
	mockRepo := newMockIncidentCategoryRepo()
	handler := NewIncidentCategoryHandler(mockRepo)

	// Add test data
	for i := 0; i < 3; i++ {
		cat := &models.IncidentCategory{
			ID:   uuid.New().String(),
			Name: "Category " + string(rune('A'+i)),
		}
		mockRepo.categories[cat.ID] = cat
	}

	app.Get("/incident-categories", handler.List)

	req := httptest.NewRequest("GET", "/incident-categories", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var response []*models.IncidentCategory
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 3 {
		t.Errorf("expected 3 categories, got %d", len(response))
	}
}

func TestCreateIncidentCategoryHandler(t *testing.T) {
	app := fiber.New()
	mockRepo := newMockIncidentCategoryRepo()
	handler := NewIncidentCategoryHandler(mockRepo)

	app.Post("/incident-categories", handler.Create)

	t.Run("Valid Input", func(t *testing.T) {
		input := models.CreateIncidentCategoryInput{
			Name:        "Security Incident",
			Description: "Security related incidents",
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/incident-categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 201 {
			t.Errorf("expected status 201, got %d", resp.StatusCode)
		}

		var created models.IncidentCategory
		json.NewDecoder(resp.Body).Decode(&created)
		if created.Name != input.Name {
			t.Errorf("expected name %s, got %s", input.Name, created.Name)
		}
	})

	t.Run("Missing Name", func(t *testing.T) {
		input := models.CreateIncidentCategoryInput{
			Description: "Description only",
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/incident-categories", bytes.NewReader(body))
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

func TestUpdateIncidentCategoryHandler(t *testing.T) {
	app := fiber.New()
	mockRepo := newMockIncidentCategoryRepo()
	handler := NewIncidentCategoryHandler(mockRepo)

	cat := &models.IncidentCategory{
		ID:   uuid.New().String(),
		Name: "Old Name",
	}
	mockRepo.categories[cat.ID] = cat

	app.Put("/incident-categories/:id", handler.Update)

	t.Run("Valid Update", func(t *testing.T) {
		newName := "New Name"
		input := models.UpdateIncidentCategoryInput{
			Name: &newName,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("PUT", "/incident-categories/"+cat.ID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var updated models.IncidentCategory
		json.NewDecoder(resp.Body).Decode(&updated)
		if updated.Name != newName {
			t.Errorf("expected name %s, got %s", newName, updated.Name)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		newName := "New Name"
		input := models.UpdateIncidentCategoryInput{
			Name: &newName,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("PUT", "/incident-categories/invalid-id", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 404 {
			t.Errorf("expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("Empty Name", func(t *testing.T) {
		emptyName := ""
		input := models.UpdateIncidentCategoryInput{
			Name: &emptyName,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("PUT", "/incident-categories/"+cat.ID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 400 {
			t.Errorf("expected status 400 for empty name, got %d", resp.StatusCode)
		}
	})

	t.Run("No Fields Provided", func(t *testing.T) {
		input := models.UpdateIncidentCategoryInput{}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("PUT", "/incident-categories/"+cat.ID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 400 {
			t.Errorf("expected status 400 for no fields, got %d", resp.StatusCode)
		}
	})
}

func TestDeleteIncidentCategoryHandler(t *testing.T) {
	app := fiber.New()
	mockRepo := newMockIncidentCategoryRepo()
	handler := NewIncidentCategoryHandler(mockRepo)

	cat := &models.IncidentCategory{
		ID:   uuid.New().String(),
		Name: "To Delete",
	}
	mockRepo.categories[cat.ID] = cat

	app.Delete("/incident-categories/:id", handler.Delete)

	t.Run("Valid Delete", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/incident-categories/"+cat.ID, nil)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 204 {
			t.Errorf("expected status 204, got %d", resp.StatusCode)
		}

		if _, exists := mockRepo.categories[cat.ID]; exists {
			t.Error("category should have been deleted")
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/incident-categories/invalid-id", nil)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 404 {
			t.Errorf("expected status 404, got %d", resp.StatusCode)
		}
	})
}
