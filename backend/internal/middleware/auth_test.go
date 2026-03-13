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

func TestRequireResponder_AdminRole(t *testing.T) {
	app := fiber.New()
	app.Use(AuthMiddleware)
	app.Use(RequireResponder)
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	user := &models.User{
		ID:    "test-id",
		Email: "admin@example.com",
		Role:  models.RoleAdmin,
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
		t.Errorf("expected status 200 for admin, got %d", resp.StatusCode)
	}
}

func TestRequireResponder_ResponderRole(t *testing.T) {
	app := fiber.New()
	app.Use(AuthMiddleware)
	app.Use(RequireResponder)
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	user := &models.User{
		ID:    "test-id",
		Email: "responder@example.com",
		Role:  models.RoleResponder,
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
		t.Errorf("expected status 200 for responder, got %d", resp.StatusCode)
	}
}

func TestRequireResponder_MemberRoleForbidden(t *testing.T) {
	app := fiber.New()
	app.Use(AuthMiddleware)
	app.Use(RequireResponder)
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	user := &models.User{
		ID:    "test-id",
		Email: "member@example.com",
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

	if resp.StatusCode != 403 {
		t.Errorf("expected status 403 for member, got %d", resp.StatusCode)
	}
}
