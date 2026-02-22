package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"backend/internal/database"
	"backend/internal/middleware"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type mockRiskRepo struct {
	risks map[string]*models.Risk
}

func (m *mockRiskRepo) Create(ctx context.Context, risk *models.Risk) error {
	if risk.ID == "" {
		risk.ID = uuid.New().String()
	}
	risk.CreatedAt = time.Now()
	risk.UpdatedAt = time.Now()
	m.risks[risk.ID] = risk
	return nil
}

func (m *mockRiskRepo) FindByID(ctx context.Context, id string) (*models.Risk, error) {
	if risk, ok := m.risks[id]; ok {
		return risk, nil
	}
	return nil, database.ErrRiskNotFound
}

func (m *mockRiskRepo) List(ctx context.Context, params *models.RiskListParams) (*models.RiskListResponse, error) {
	var risks []*models.Risk
	for _, risk := range m.risks {
		risks = append(risks, risk)
	}
	return &models.RiskListResponse{
		Data: risks,
		Meta: models.Meta{
			Total: len(risks),
			Page:  params.Page,
			Limit: params.Limit,
		},
	}, nil
}

func (m *mockRiskRepo) Update(ctx context.Context, risk *models.Risk) error {
	if _, ok := m.risks[risk.ID]; !ok {
		return nil
	}
	risk.UpdatedAt = time.Now()
	m.risks[risk.ID] = risk
	return nil
}

func (m *mockRiskRepo) Delete(ctx context.Context, id string) error {
	delete(m.risks, id)
	return nil
}

type mockCategoryRepo struct {
	categories map[string]*models.Category
}

func (m *mockCategoryRepo) FindByID(ctx context.Context, id string) (*models.Category, error) {
	if category, ok := m.categories[id]; ok {
		return category, nil
	}
	return nil, nil
}

func (m *mockCategoryRepo) List(ctx context.Context) ([]*models.Category, error) {
	var categories []*models.Category
	for _, category := range m.categories {
		categories = append(categories, category)
	}
	return categories, nil
}

func (m *mockCategoryRepo) Create(ctx context.Context, input *models.CreateCategoryInput) (*models.Category, error) {
	category := &models.Category{
		ID:          uuid.New().String(),
		Name:        input.Name,
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.categories[category.ID] = category
	return category, nil
}

func (m *mockCategoryRepo) Update(ctx context.Context, id string, input *models.UpdateCategoryInput) (*models.Category, error) {
	category, ok := m.categories[id]
	if !ok {
		return nil, database.ErrCategoryNotFound
	}
	if input.Name != nil {
		category.Name = *input.Name
	}
	if input.Description != nil {
		category.Description = *input.Description
	}
	category.UpdatedAt = time.Now()
	m.categories[id] = category
	return category, nil
}

func (m *mockCategoryRepo) Delete(ctx context.Context, id string) error {
	if _, ok := m.categories[id]; !ok {
		return database.ErrCategoryNotFound
	}
	delete(m.categories, id)
	return nil
}

type mockAuditRepo struct {
	logs []*models.AuditLog
}

func (m *mockAuditRepo) ListByEntity(ctx context.Context, entityType, entityID string, limit int) ([]*models.AuditLog, error) {
	var result []*models.AuditLog
	for _, log := range m.logs {
		if log.EntityType == entityType && log.EntityID == entityID {
			result = append(result, log)
		}
	}
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}

func (m *mockAuditRepo) Create(ctx context.Context, entityType, entityID string, action models.AuditAction, changes map[string]any, userID string) error {
	log := &models.AuditLog{
		ID:         uuid.New().String(),
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		Changes:    changes,
		UserID:     userID,
		CreatedAt:  time.Now(),
	}
	m.logs = append(m.logs, log)
	return nil
}

// testAuthMiddleware sets up a user context for testing protected handlers
func testAuthMiddleware(c *fiber.Ctx) error {
	c.Locals(middleware.UserKey, &middleware.UserClaims{
		UserID: "test-user-id",
		Email:  "test@example.com",
		Role:   "member",
	})
	return c.Next()
}

func TestCreateRiskHandler_ValidInput(t *testing.T) {
	app := fiber.New()
	mockRiskRepo := &mockRiskRepo{risks: make(map[string]*models.Risk)}
	mockCategoryRepo := &mockCategoryRepo{categories: make(map[string]*models.Category)}
	mockAuditRepo := &mockAuditRepo{logs: []*models.AuditLog{}}
	handler := NewRiskHandler(mockRiskRepo, mockCategoryRepo, mockAuditRepo)

	app.Post("/risks", testAuthMiddleware, handler.Create)

	// Create test request with valid risk data
	riskInput := map[string]interface{}{
		"title":       "Test Risk",
		"description": "Test description",
		"owner_id":    "user-123",
		"status":      "open",
		"severity":    "high",
	}
	body, _ := json.Marshal(riskInput)
	req := httptest.NewRequest("POST", "/risks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 201 {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}

	// Verify the risk was created
	if len(mockRiskRepo.risks) != 1 {
		t.Errorf("expected 1 risk, got %d", len(mockRiskRepo.risks))
	}

	// Verify audit log was created
	if len(mockAuditRepo.logs) != 1 {
		t.Errorf("expected 1 audit log, got %d", len(mockAuditRepo.logs))
	}
}

func TestCreateRiskHandler_InvalidInput(t *testing.T) {
	app := fiber.New()
	mockRiskRepo := &mockRiskRepo{risks: make(map[string]*models.Risk)}
	mockCategoryRepo := &mockCategoryRepo{categories: make(map[string]*models.Category)}
	mockAuditRepo := &mockAuditRepo{logs: []*models.AuditLog{}}
	handler := NewRiskHandler(mockRiskRepo, mockCategoryRepo, mockAuditRepo)

	app.Post("/risks", testAuthMiddleware, handler.Create)

	// Missing required fields
	body := map[string]interface{}{
		"description": "Test risk description",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/risks", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreateRiskHandler_EmptyCategory(t *testing.T) {
	app := fiber.New()
	mockRiskRepo := &mockRiskRepo{risks: make(map[string]*models.Risk)}
	mockCategoryRepo := &mockCategoryRepo{categories: make(map[string]*models.Category)}
	mockAuditRepo := &mockAuditRepo{logs: []*models.AuditLog{}}
	handler := NewRiskHandler(mockRiskRepo, mockCategoryRepo, mockAuditRepo)

	app.Post("/risks", testAuthMiddleware, handler.Create)

	// Empty category ID should be normalized to nil
	body := map[string]interface{}{
		"title":       "Test Risk",
		"description": "Test description",
		"owner_id":    "user-123",
		"category_id": "",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/risks", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 201 {
		t.Errorf("expected status 201 for empty category, got %d", resp.StatusCode)
	}
}

func TestCreateRiskHandler_InvalidCategory(t *testing.T) {
	app := fiber.New()
	mockRiskRepo := &mockRiskRepo{risks: make(map[string]*models.Risk)}
	mockCategoryRepo := &mockCategoryRepo{categories: make(map[string]*models.Category)}
	mockAuditRepo := &mockAuditRepo{logs: []*models.AuditLog{}}
	handler := NewRiskHandler(mockRiskRepo, mockCategoryRepo, mockAuditRepo)

	app.Post("/risks", testAuthMiddleware, handler.Create)

	// Invalid category ID (non-existent)
	body := map[string]interface{}{
		"title":       "Test Risk",
		"description": "Test description",
		"owner_id":    "user-123",
		"category_id": "invalid-category-id",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/risks", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("expected status 400 for invalid category, got %d", resp.StatusCode)
	}
}

func TestGetRiskHandler_ValidID(t *testing.T) {
	app := fiber.New()
	mockRiskRepo := &mockRiskRepo{risks: make(map[string]*models.Risk)}
	mockCategoryRepo := &mockCategoryRepo{categories: make(map[string]*models.Category)}
	mockAuditRepo := &mockAuditRepo{logs: []*models.AuditLog{}}
	handler := NewRiskHandler(mockRiskRepo, mockCategoryRepo, mockAuditRepo)

	// Create a test risk first
	testRisk := &models.Risk{
		ID:          uuid.New().String(),
		Title:       "Test Risk",
		Description: "Test description",
		Status:      "open",
		Severity:    "high",
	}
	mockRiskRepo.risks[testRisk.ID] = testRisk

	app.Get("/risks/:id", testAuthMiddleware, handler.Get)

	req := httptest.NewRequest("GET", "/risks/"+testRisk.ID, nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	if response["id"] != testRisk.ID {
		t.Errorf("expected ID %s, got %v", testRisk.ID, response["id"])
	}
}

func TestGetRiskHandler_InvalidID(t *testing.T) {
	app := fiber.New()
	mockRiskRepo := &mockRiskRepo{risks: make(map[string]*models.Risk)}
	mockCategoryRepo := &mockCategoryRepo{categories: make(map[string]*models.Category)}
	mockAuditRepo := &mockAuditRepo{logs: []*models.AuditLog{}}
	handler := NewRiskHandler(mockRiskRepo, mockCategoryRepo, mockAuditRepo)

	app.Get("/risks/:id", testAuthMiddleware, handler.Get)

	req := httptest.NewRequest("GET", "/risks/invalid-id", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestListRisksHandler(t *testing.T) {
	app := fiber.New()
	mockRiskRepo := &mockRiskRepo{risks: make(map[string]*models.Risk)}
	mockCategoryRepo := &mockCategoryRepo{categories: make(map[string]*models.Category)}
	mockAuditRepo := &mockAuditRepo{logs: []*models.AuditLog{}}
	handler := NewRiskHandler(mockRiskRepo, mockCategoryRepo, mockAuditRepo)

	// Add some test risks
	for i := 0; i < 3; i++ {
		risk := &models.Risk{
			ID:    uuid.New().String(),
			Title: "Test Risk " + string(rune('A'+i)),
		}
		mockRiskRepo.risks[risk.ID] = risk
	}

	app.Get("/risks", testAuthMiddleware, handler.List)

	req := httptest.NewRequest("GET", "/risks", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var response struct {
		Data []map[string]interface{} `json:"data"`
		Meta map[string]interface{}   `json:"meta"`
	}
	json.NewDecoder(resp.Body).Decode(&response)

	if len(response.Data) != 3 {
		t.Errorf("expected 3 risks, got %d", len(response.Data))
	}
}

func TestUpdateRiskHandler_ValidInput(t *testing.T) {
	app := fiber.New()
	mockRiskRepo := &mockRiskRepo{risks: make(map[string]*models.Risk)}
	mockCategoryRepo := &mockCategoryRepo{categories: make(map[string]*models.Category)}
	mockAuditRepo := &mockAuditRepo{logs: []*models.AuditLog{}}
	handler := NewRiskHandler(mockRiskRepo, mockCategoryRepo, mockAuditRepo)

	// Create a test risk first
	testRisk := &models.Risk{
		ID:          uuid.New().String(),
		Title:       "Original Title",
		Description: "Original description",
		Status:      "open",
		Severity:    "medium",
	}
	mockRiskRepo.risks[testRisk.ID] = testRisk

	app.Put("/risks/:id", testAuthMiddleware, handler.Update)

	body := map[string]interface{}{
		"title": "Updated Title",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/risks/"+testRisk.ID, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	if response["title"] != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got %v", response["title"])
	}
}

func TestDeleteRiskHandler(t *testing.T) {
	app := fiber.New()
	mockRiskRepo := &mockRiskRepo{risks: make(map[string]*models.Risk)}
	mockCategoryRepo := &mockCategoryRepo{categories: make(map[string]*models.Category)}
	mockAuditRepo := &mockAuditRepo{logs: []*models.AuditLog{}}
	handler := NewRiskHandler(mockRiskRepo, mockCategoryRepo, mockAuditRepo)

	// Create a test risk first
	testRisk := &models.Risk{
		ID:    uuid.New().String(),
		Title: "Test Risk to Delete",
	}
	mockRiskRepo.risks[testRisk.ID] = testRisk

	app.Delete("/risks/:id", testAuthMiddleware, handler.Delete)

	req := httptest.NewRequest("DELETE", "/risks/"+testRisk.ID, nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}

	// Verify the risk was deleted
	if _, exists := mockRiskRepo.risks[testRisk.ID]; exists {
		t.Errorf("risk should have been deleted")
	}
}
