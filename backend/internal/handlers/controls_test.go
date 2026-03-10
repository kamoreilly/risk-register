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

func TestFrameworkControlHandler(t *testing.T) {
	app := fiber.New()
	repo := &mockFrameworkControlRepo{controls: make(map[string]*models.FrameworkControl)}
	handler := NewFrameworkControlHandler(repo)

	app.Get("/controls", testAuthMiddleware, handler.List)
	app.Post("/controls", testAuthMiddleware, handler.Create)
	app.Put("/controls/:id", testAuthMiddleware, handler.Update)
	app.Delete("/controls/:id", testAuthMiddleware, handler.Delete)

	frameworkID := uuid.New().String()

	t.Run("Create Control", func(t *testing.T) {
		input := models.CreateFrameworkControlInput{
			FrameworkID: frameworkID,
			ControlRef:  "A.12.1.1",
			Title:       "Operating procedures",
			Description: "Ensure documented procedures exist.",
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/controls", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 201 {
			t.Fatalf("expected status 201, got %d", resp.StatusCode)
		}

		var created models.FrameworkControl
		json.NewDecoder(resp.Body).Decode(&created)
		if created.ControlRef != input.ControlRef {
			t.Fatalf("expected control ref %s, got %s", input.ControlRef, created.ControlRef)
		}
	})

	t.Run("List Controls", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/controls", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var response map[string][]*models.FrameworkControl
		json.NewDecoder(resp.Body).Decode(&response)
		if len(response["data"]) != 1 {
			t.Fatalf("expected 1 control, got %d", len(response["data"]))
		}
	})
}
