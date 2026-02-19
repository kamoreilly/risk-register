package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"backend/internal/auth"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type mockUserRepo struct {
	users map[string]*models.User
}

func (m *mockUserRepo) Create(ctx context.Context, user *models.User) error {
	if _, exists := m.users[user.Email]; exists {
		return nil // simulate duplicate
	}
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	if user, ok := m.users[email]; ok {
		return user, nil
	}
	return nil, nil
}

func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*models.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, nil
}

func TestRegisterHandler_ValidInput(t *testing.T) {
	app := fiber.New()
	mockRepo := &mockUserRepo{users: make(map[string]*models.User)}
	handler := NewAuthHandler(mockRepo)

	app.Post("/register", handler.Register)

	body := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
		"name":     "Test User",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 201 {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}
}

func TestLoginHandler_ValidCredentials(t *testing.T) {
	app := fiber.New()
	mockRepo := &mockUserRepo{users: make(map[string]*models.User)}
	handler := NewAuthHandler(mockRepo)

	// First create a user with hashed password
	hashedPassword, _ := auth.HashPassword("password123")
	mockRepo.users["test@example.com"] = &models.User{
		ID:           "test-id",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		Name:         "Test User",
		Role:         models.RoleMember,
	}

	app.Post("/login", handler.Login)

	body := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}
