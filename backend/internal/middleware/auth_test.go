package middleware

import (
	"net/http/httptest"
	"testing"

	"backend/internal/auth"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	app := fiber.New()
	app.Use(AuthMiddleware)
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	user := &models.User{
		ID:    "test-id",
		Email: "test@example.com",
		Role:  models.RoleMember,
	}
	token, err := auth.GenerateToken(user)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	app := fiber.New()
	app.Use(AuthMiddleware)
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/protected", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 401 {
		t.Errorf("expected status 401, got %d", resp.StatusCode)
	}
}
