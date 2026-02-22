package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func TestListCategoriesHandler(t *testing.T) {
	app := fiber.New()
	mockRepo := &mockCategoryRepo{categories: make(map[string]*models.Category)}
	handler := NewCategoryHandler(mockRepo)

	// Add test data
	for i := 0; i < 3; i++ {
		cat := &models.Category{
			ID:   uuid.New().String(),
			Name: "Category " + string(rune('A'+i)),
		}
		mockRepo.categories[cat.ID] = cat
	}

	app.Get("/categories", handler.List)

	req := httptest.NewRequest("GET", "/categories", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var response []*models.Category
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 3 {
		t.Errorf("expected 3 categories, got %d", len(response))
	}
}

func TestCreateCategoryHandler(t *testing.T) {
	app := fiber.New()
	mockRepo := &mockCategoryRepo{categories: make(map[string]*models.Category)}
	handler := NewCategoryHandler(mockRepo)

	app.Post("/categories", handler.Create)

	t.Run("Valid Input", func(t *testing.T) {
		input := models.CreateCategoryInput{
			Name:        "New Category",
			Description: "Description",
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 201 {
			t.Errorf("expected status 201, got %d", resp.StatusCode)
		}

		var created models.Category
		json.NewDecoder(resp.Body).Decode(&created)
		if created.Name != input.Name {
			t.Errorf("expected name %s, got %s", input.Name, created.Name)
		}
	})

	t.Run("Missing Name", func(t *testing.T) {
		input := models.CreateCategoryInput{
			Description: "Description only",
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/categories", bytes.NewReader(body))
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

func TestUpdateCategoryHandler(t *testing.T) {
	app := fiber.New()
	mockRepo := &mockCategoryRepo{categories: make(map[string]*models.Category)}
	handler := NewCategoryHandler(mockRepo)

	cat := &models.Category{
		ID:   uuid.New().String(),
		Name: "Old Name",
	}
	mockRepo.categories[cat.ID] = cat

	app.Put("/categories/:id", handler.Update)

	t.Run("Valid Update", func(t *testing.T) {
		newName := "New Name"
		input := models.UpdateCategoryInput{
			Name: &newName,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("PUT", "/categories/"+cat.ID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var updated models.Category
		json.NewDecoder(resp.Body).Decode(&updated)
		if updated.Name != newName {
			t.Errorf("expected name %s, got %s", newName, updated.Name)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		newName := "New Name"
		input := models.UpdateCategoryInput{
			Name: &newName,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("PUT", "/categories/invalid-id", bytes.NewReader(body))
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

func TestDeleteCategoryHandler(t *testing.T) {
	app := fiber.New()
	mockRepo := &mockCategoryRepo{categories: make(map[string]*models.Category)}
	handler := NewCategoryHandler(mockRepo)

	cat := &models.Category{
		ID:   uuid.New().String(),
		Name: "To Delete",
	}
	mockRepo.categories[cat.ID] = cat

	app.Delete("/categories/:id", handler.Delete)

	t.Run("Valid Delete", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/categories/"+cat.ID, nil)

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
		req := httptest.NewRequest("DELETE", "/categories/invalid-id", nil)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 404 {
			t.Errorf("expected status 404, got %d", resp.StatusCode)
		}
	})
}
