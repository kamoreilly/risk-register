package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"backend/internal/database"
	"backend/internal/middleware"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	pgContainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// =============================================================================
// Integration Tests for Incident Management
// =============================================================================

// TestMain sets up testcontainers for integration tests
func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		os.Exit(m.Run())
	}

	teardown, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	if err := runMigrations(); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	code := m.Run()

	if err := teardown(context.Background()); err != nil {
		log.Printf("warning: could not teardown postgres container: %v", err)
	}

	os.Exit(code)
}

func mustStartPostgresContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	var (
		dbName = "database"
		dbPwd  = "password"
		dbUser = "user"
	)

	dbContainer, err := pgContainer.Run(
		context.Background(),
		"postgres:latest",
		pgContainer.WithDatabase(dbName),
		pgContainer.WithUsername(dbUser),
		pgContainer.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	// Set environment variables for database package
	os.Setenv("RISK_REGISTER_DB_DATABASE", dbName)
	os.Setenv("RISK_REGISTER_DB_PASSWORD", dbPwd)
	os.Setenv("RISK_REGISTER_DB_USERNAME", dbUser)

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}
	os.Setenv("RISK_REGISTER_DB_HOST", dbHost)

	dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}
	os.Setenv("RISK_REGISTER_DB_PORT", dbPort.Port())

	return dbContainer.Terminate, err
}

func runMigrations() error {
	s := database.New()
	db := database.GetDB(s)
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	_, filename, _, _ := runtime.Caller(0)
	migrationsPath := filepath.Join(filepath.Dir(filename), "../../database/migrations/migrations")

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver)
	if err != nil {
		return err
	}
	return m.Up()
}

// integrationTestAuthMiddleware sets up a user context for integration testing
func integrationTestAuthMiddleware(userID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals(middleware.UserKey, &middleware.UserClaims{
			UserID: userID,
			Email:  "test@example.com",
			Role:   "member",
		})
		return c.Next()
	}
}

func TestIncidentHandler_Integration_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := database.GetDB(database.New())
	incidentRepo := database.NewIncidentRepository(db)
	incidentCategoryRepo := database.NewIncidentCategoryRepository(db)
	incidentRiskRepo := database.NewIncidentRiskRepository(db)
	auditRepo := database.NewAuditLogRepository(db)
	userRepo := database.NewUserRepository(db)

	ctx := context.Background()

	// Setup: Create a test user
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        "incident-test-" + uuid.New().String() + "@example.com",
		PasswordHash: "hash",
		Name:         "Incident Tester",
		Role:         models.RoleMember,
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	// Setup: Create an incident category
	category := &models.CreateIncidentCategoryInput{
		Name:        "Security",
		Description: "Security incidents",
	}
	createdCategory, err := incidentCategoryRepo.Create(ctx, category)
	require.NoError(t, err)

	app := fiber.New()
	handler := NewIncidentHandler(incidentRepo, incidentCategoryRepo, incidentRiskRepo, auditRepo)
	app.Post("/incidents", integrationTestAuthMiddleware(user.ID), handler.Create)

	tests := []struct {
		name           string
		input          models.CreateIncidentInput
		expectedStatus int
		checkResponse  func(t *testing.T, incident *models.Incident)
	}{
		{
			name: "create incident with all fields",
			input: models.CreateIncidentInput{
				Title:           "Production Database Outage",
				Description:     "Primary database became unresponsive",
				Priority:        models.PriorityP1,
				Status:          models.IncidentStatusNew,
				ServiceAffected: "PostgreSQL Primary",
			},
			expectedStatus: 201,
			checkResponse: func(t *testing.T, incident *models.Incident) {
				assert.NotEmpty(t, incident.ID)
				assert.Equal(t, "Production Database Outage", incident.Title)
				assert.Equal(t, models.PriorityP1, incident.Priority)
				assert.Equal(t, models.IncidentStatusNew, incident.Status)
				assert.Equal(t, user.ID, incident.ReporterID)
				assert.Equal(t, user.ID, incident.CreatedBy)
				assert.NotZero(t, incident.CreatedAt)
				assert.NotZero(t, incident.OccurredAt)
				assert.NotZero(t, incident.DetectedAt)
			},
		},
		{
			name: "create incident with defaults",
			input: models.CreateIncidentInput{
				Title: "Minimal Incident",
			},
			expectedStatus: 201,
			checkResponse: func(t *testing.T, incident *models.Incident) {
				assert.Equal(t, models.PriorityP3, incident.Priority)
				assert.Equal(t, models.IncidentStatusNew, incident.Status)
			},
		},
		{
			name: "create incident with category",
			input: models.CreateIncidentInput{
				Title:      "Security Breach",
				CategoryID: &createdCategory.ID,
			},
			expectedStatus: 201,
			checkResponse: func(t *testing.T, incident *models.Incident) {
				assert.NotNil(t, incident.CategoryID)
				assert.Equal(t, createdCategory.ID, *incident.CategoryID)
			},
		},
		{
			name: "create incident missing title",
			input: models.CreateIncidentInput{
				Description: "No title provided",
			},
			expectedStatus: 400,
			checkResponse:  nil,
		},
		{
			name: "create incident with invalid category",
			input: models.CreateIncidentInput{
				Title:      "Bad Category",
				CategoryID: ptrIncident("non-existent-category-id"),
			},
			expectedStatus: 400,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/incidents", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.checkResponse != nil && resp.StatusCode == 201 {
				var incident models.Incident
				err := json.NewDecoder(resp.Body).Decode(&incident)
				require.NoError(t, err)
				tt.checkResponse(t, &incident)
			}
		})
	}
}

func TestIncidentHandler_Integration_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := database.GetDB(database.New())
	incidentRepo := database.NewIncidentRepository(db)
	incidentCategoryRepo := database.NewIncidentCategoryRepository(db)
	incidentRiskRepo := database.NewIncidentRiskRepository(db)
	auditRepo := database.NewAuditLogRepository(db)
	userRepo := database.NewUserRepository(db)

	ctx := context.Background()

	// Setup: Create a test user
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        "incident-get-" + uuid.New().String() + "@example.com",
		PasswordHash: "hash",
		Name:         "Incident Get Tester",
		Role:         models.RoleMember,
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	// Setup: Create an incident
	incident := &models.Incident{
		Title:       "Test Incident for Get",
		Description: "Description for get test",
		Status:      models.IncidentStatusNew,
		Priority:    models.PriorityP2,
		ReporterID:  user.ID,
		CreatedBy:   user.ID,
		UpdatedBy:   user.ID,
	}
	err = incidentRepo.Create(ctx, incident)
	require.NoError(t, err)

	app := fiber.New()
	handler := NewIncidentHandler(incidentRepo, incidentCategoryRepo, incidentRiskRepo, auditRepo)
	app.Get("/incidents/:id", handler.Get)

	t.Run("get existing incident", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/incidents/"+incident.ID, nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var fetched models.Incident
		err = json.NewDecoder(resp.Body).Decode(&fetched)
		require.NoError(t, err)
		assert.Equal(t, incident.ID, fetched.ID)
		assert.Equal(t, incident.Title, fetched.Title)
		assert.Equal(t, incident.Priority, fetched.Priority)
	})

	t.Run("get non-existent incident", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/incidents/non-existent-id", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
	})
}

func TestIncidentHandler_Integration_List(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := database.GetDB(database.New())
	incidentRepo := database.NewIncidentRepository(db)
	incidentCategoryRepo := database.NewIncidentCategoryRepository(db)
	incidentRiskRepo := database.NewIncidentRiskRepository(db)
	auditRepo := database.NewAuditLogRepository(db)
	userRepo := database.NewUserRepository(db)

	ctx := context.Background()

	// Setup: Create a test user
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        "incident-list-" + uuid.New().String() + "@example.com",
		PasswordHash: "hash",
		Name:         "Incident List Tester",
		Role:         models.RoleMember,
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	// Setup: Create a category
	category := &models.CreateIncidentCategoryInput{
		Name:        "Infrastructure",
		Description: "Infra incidents",
	}
	createdCategory, err := incidentCategoryRepo.Create(ctx, category)
	require.NoError(t, err)

	// Setup: Create multiple incidents with different attributes
	incidents := []*models.Incident{
		{
			Title:       "Critical Outage",
			Description: "System down",
			Status:      models.IncidentStatusNew,
			Priority:    models.PriorityP1,
			ReporterID:  user.ID,
			CategoryID:  &createdCategory.ID,
			CreatedBy:   user.ID,
			UpdatedBy:   user.ID,
		},
		{
			Title:       "Medium Issue",
			Description: "Performance degradation",
			Status:      models.IncidentStatusInProgress,
			Priority:    models.PriorityP2,
			ReporterID:  user.ID,
			CreatedBy:   user.ID,
			UpdatedBy:   user.ID,
		},
		{
			Title:       "Low Priority Bug",
			Description: "Minor UI issue",
			Status:      models.IncidentStatusResolved,
			Priority:    models.PriorityP4,
			ReporterID:  user.ID,
			CreatedBy:   user.ID,
			UpdatedBy:   user.ID,
		},
	}

	for _, inc := range incidents {
		err := incidentRepo.Create(ctx, inc)
		require.NoError(t, err)
	}

	app := fiber.New()
	handler := NewIncidentHandler(incidentRepo, incidentCategoryRepo, incidentRiskRepo, auditRepo)
	app.Get("/incidents", handler.List)

	tests := []struct {
		name          string
		queryParams   string
		expectedCount int
		checkResults  func(t *testing.T, response *models.IncidentListResponse)
	}{
		{
			name:          "list all incidents",
			queryParams:   "",
			expectedCount: 3,
			checkResults:  nil,
		},
		{
			name:        "filter by status",
			queryParams: "?status=new",
			expectedCount: 1,
			checkResults: func(t *testing.T, response *models.IncidentListResponse) {
				for _, inc := range response.Data {
					assert.Equal(t, models.IncidentStatusNew, inc.Status)
				}
			},
		},
		{
			name:        "filter by priority",
			queryParams: "?priority=p1",
			expectedCount: 1,
			checkResults: func(t *testing.T, response *models.IncidentListResponse) {
				for _, inc := range response.Data {
					assert.Equal(t, models.PriorityP1, inc.Priority)
				}
			},
		},
		{
			name:        "filter by category",
			queryParams: "?category_id=" + createdCategory.ID,
			expectedCount: 1,
			checkResults: func(t *testing.T, response *models.IncidentListResponse) {
				for _, inc := range response.Data {
					assert.NotNil(t, inc.CategoryID)
					assert.Equal(t, createdCategory.ID, *inc.CategoryID)
				}
			},
		},
		{
			name:        "search by title",
			queryParams: "?search=Outage",
			expectedCount: 1,
			checkResults: func(t *testing.T, response *models.IncidentListResponse) {
				assert.Equal(t, "Critical Outage", response.Data[0].Title)
			},
		},
		{
			name:        "pagination",
			queryParams: "?page=1&limit=2",
			expectedCount: 2,
			checkResults: func(t *testing.T, response *models.IncidentListResponse) {
				assert.Equal(t, 1, response.Meta.Page)
				assert.Equal(t, 2, response.Meta.Limit)
				assert.Equal(t, 3, response.Meta.Total)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/incidents"+tt.queryParams, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)

			var response models.IncidentListResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(response.Data), tt.expectedCount)

			if tt.checkResults != nil {
				tt.checkResults(t, &response)
			}
		})
	}
}

func TestIncidentHandler_Integration_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := database.GetDB(database.New())
	incidentRepo := database.NewIncidentRepository(db)
	incidentCategoryRepo := database.NewIncidentCategoryRepository(db)
	incidentRiskRepo := database.NewIncidentRiskRepository(db)
	auditRepo := database.NewAuditLogRepository(db)
	userRepo := database.NewUserRepository(db)

	ctx := context.Background()

	// Setup: Create a test user
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        "incident-update-" + uuid.New().String() + "@example.com",
		PasswordHash: "hash",
		Name:         "Incident Update Tester",
		Role:         models.RoleMember,
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	// Setup: Create a category
	category := &models.CreateIncidentCategoryInput{
		Name:        "Network",
		Description: "Network incidents",
	}
	createdCategory, err := incidentCategoryRepo.Create(ctx, category)
	require.NoError(t, err)

	app := fiber.New()
	handler := NewIncidentHandler(incidentRepo, incidentCategoryRepo, incidentRiskRepo, auditRepo)
	app.Put("/incidents/:id", integrationTestAuthMiddleware(user.ID), handler.Update)

	t.Run("update title and description", func(t *testing.T) {
		// Create an incident
		incident := &models.Incident{
			Title:       "Original Title",
			Description: "Original Description",
			Status:      models.IncidentStatusNew,
			Priority:    models.PriorityP3,
			ReporterID:  user.ID,
			CreatedBy:   user.ID,
			UpdatedBy:   user.ID,
		}
		err := incidentRepo.Create(ctx, incident)
		require.NoError(t, err)

		newTitle := "Updated Title"
		newDesc := "Updated Description"
		input := models.UpdateIncidentInput{
			Title:       &newTitle,
			Description: &newDesc,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("PUT", "/incidents/"+incident.ID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var updated models.Incident
		err = json.NewDecoder(resp.Body).Decode(&updated)
		require.NoError(t, err)
		assert.Equal(t, newTitle, updated.Title)
		assert.Equal(t, newDesc, updated.Description)
	})

	t.Run("update status to resolved sets resolved_at", func(t *testing.T) {
		incident := &models.Incident{
			Title:       "To Resolve",
			Status:      models.IncidentStatusInProgress,
			Priority:    models.PriorityP2,
			ReporterID:  user.ID,
			CreatedBy:   user.ID,
			UpdatedBy:   user.ID,
		}
		err := incidentRepo.Create(ctx, incident)
		require.NoError(t, err)

		newStatus := models.IncidentStatusResolved
		input := models.UpdateIncidentInput{
			Status: &newStatus,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("PUT", "/incidents/"+incident.ID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var updated models.Incident
		err = json.NewDecoder(resp.Body).Decode(&updated)
		require.NoError(t, err)
		assert.Equal(t, models.IncidentStatusResolved, updated.Status)
		assert.NotNil(t, updated.ResolvedAt)
	})

	t.Run("update status to closed sets resolved_at", func(t *testing.T) {
		incident := &models.Incident{
			Title:       "To Close",
			Status:      models.IncidentStatusInProgress,
			Priority:    models.PriorityP2,
			ReporterID:  user.ID,
			CreatedBy:   user.ID,
			UpdatedBy:   user.ID,
		}
		err := incidentRepo.Create(ctx, incident)
		require.NoError(t, err)

		newStatus := models.IncidentStatusClosed
		input := models.UpdateIncidentInput{
			Status: &newStatus,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("PUT", "/incidents/"+incident.ID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var updated models.Incident
		err = json.NewDecoder(resp.Body).Decode(&updated)
		require.NoError(t, err)
		assert.Equal(t, models.IncidentStatusClosed, updated.Status)
		assert.NotNil(t, updated.ResolvedAt)
	})

	t.Run("update with category", func(t *testing.T) {
		incident := &models.Incident{
			Title:       "Categorized Incident",
			Status:      models.IncidentStatusNew,
			Priority:    models.PriorityP3,
			ReporterID:  user.ID,
			CreatedBy:   user.ID,
			UpdatedBy:   user.ID,
		}
		err := incidentRepo.Create(ctx, incident)
		require.NoError(t, err)

		input := models.UpdateIncidentInput{
			CategoryID: &createdCategory.ID,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("PUT", "/incidents/"+incident.ID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var updated models.Incident
		err = json.NewDecoder(resp.Body).Decode(&updated)
		require.NoError(t, err)
		assert.NotNil(t, updated.CategoryID)
		assert.Equal(t, createdCategory.ID, *updated.CategoryID)
	})

	t.Run("update with invalid category returns error", func(t *testing.T) {
		incident := &models.Incident{
			Title:       "Bad Category Update",
			Status:      models.IncidentStatusNew,
			Priority:    models.PriorityP3,
			ReporterID:  user.ID,
			CreatedBy:   user.ID,
			UpdatedBy:   user.ID,
		}
		err := incidentRepo.Create(ctx, incident)
		require.NoError(t, err)

		input := models.UpdateIncidentInput{
			CategoryID: ptrIncident("non-existent-category"),
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("PUT", "/incidents/"+incident.ID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("update non-existent incident returns 404", func(t *testing.T) {
		newTitle := "Updated"
		input := models.UpdateIncidentInput{
			Title: &newTitle,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("PUT", "/incidents/non-existent-id", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
	})
}

func TestIncidentHandler_Integration_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := database.GetDB(database.New())
	incidentRepo := database.NewIncidentRepository(db)
	incidentCategoryRepo := database.NewIncidentCategoryRepository(db)
	incidentRiskRepo := database.NewIncidentRiskRepository(db)
	auditRepo := database.NewAuditLogRepository(db)
	userRepo := database.NewUserRepository(db)

	ctx := context.Background()

	// Setup: Create a test user
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        "incident-delete-" + uuid.New().String() + "@example.com",
		PasswordHash: "hash",
		Name:         "Incident Delete Tester",
		Role:         models.RoleMember,
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	app := fiber.New()
	handler := NewIncidentHandler(incidentRepo, incidentCategoryRepo, incidentRiskRepo, auditRepo)
	app.Delete("/incidents/:id", integrationTestAuthMiddleware(user.ID), handler.Delete)

	t.Run("delete existing incident", func(t *testing.T) {
		incident := &models.Incident{
			Title:       "To Delete",
			Status:      models.IncidentStatusNew,
			Priority:    models.PriorityP3,
			ReporterID:  user.ID,
			CreatedBy:   user.ID,
			UpdatedBy:   user.ID,
		}
		err := incidentRepo.Create(ctx, incident)
		require.NoError(t, err)

		req := httptest.NewRequest("DELETE", "/incidents/"+incident.ID, nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 204, resp.StatusCode)

		// Verify deletion
		_, err = incidentRepo.FindByID(ctx, incident.ID)
		assert.Error(t, err)
		assert.Equal(t, database.ErrIncidentNotFound, err)
	})

	t.Run("delete non-existent incident returns 404", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/incidents/non-existent-id", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
	})
}

func TestIncidentRiskHandler_Integration_LinkUnlink(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := database.GetDB(database.New())
	incidentRepo := database.NewIncidentRepository(db)
	incidentRiskRepo := database.NewIncidentRiskRepository(db)
	auditRepo := database.NewAuditLogRepository(db)
	userRepo := database.NewUserRepository(db)
	riskRepo := database.NewRiskRepository(db)
	catRepo := database.NewCategoryRepository(db)

	ctx := context.Background()

	// Setup: Create a test user
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        "incident-risk-" + uuid.New().String() + "@example.com",
		PasswordHash: "hash",
		Name:         "Incident Risk Tester",
		Role:         models.RoleMember,
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	// Setup: Create a category for risk
	category := &models.CreateCategoryInput{
		Name:        "Risk Category",
		Description: "Test category",
	}
	createdCategory, err := catRepo.Create(ctx, category)
	require.NoError(t, err)

	// Setup: Create a risk
	risk := &models.Risk{
		Title:       "Linked Risk",
		Description: "Risk to be linked to incident",
		OwnerID:     user.ID,
		Status:      models.StatusOpen,
		Severity:    models.SeverityHigh,
		CategoryID:  &createdCategory.ID,
		CreatedBy:   user.ID,
		UpdatedBy:   user.ID,
	}
	err = riskRepo.Create(ctx, risk)
	require.NoError(t, err)

	// Setup: Create an incident
	incident := &models.Incident{
		Title:       "Incident with Risk",
		Status:      models.IncidentStatusNew,
		Priority:    models.PriorityP2,
		ReporterID:  user.ID,
		CreatedBy:   user.ID,
		UpdatedBy:   user.ID,
	}
	err = incidentRepo.Create(ctx, incident)
	require.NoError(t, err)

	app := fiber.New()
	handler := NewIncidentRiskHandler(incidentRiskRepo, auditRepo)

	app.Get("/incidents/:incidentId/risks", handler.ListRisks)
	app.Post("/incidents/:incidentId/risks", integrationTestAuthMiddleware(user.ID), handler.LinkRisk)
	app.Delete("/incidents/:incidentId/risks/:riskId", integrationTestAuthMiddleware(user.ID), handler.UnlinkRisk)

	t.Run("link risk to incident", func(t *testing.T) {
		input := models.LinkIncidentRiskInput{
			RiskID: risk.ID,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/incidents/"+incident.ID+"/risks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)

		var link models.IncidentRisk
		err = json.NewDecoder(resp.Body).Decode(&link)
		require.NoError(t, err)
		assert.Equal(t, incident.ID, link.IncidentID)
		assert.Equal(t, risk.ID, link.RiskID)
		assert.Equal(t, user.ID, link.CreatedBy)
	})

	t.Run("list linked risks", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/incidents/"+incident.ID+"/risks", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var links []*models.IncidentRisk
		err = json.NewDecoder(resp.Body).Decode(&links)
		require.NoError(t, err)
		assert.Len(t, links, 1)
		assert.Equal(t, risk.ID, links[0].RiskID)
		assert.NotNil(t, links[0].Risk)
		assert.Equal(t, risk.Title, links[0].Risk.Title)
	})

	t.Run("duplicate link returns 409", func(t *testing.T) {
		input := models.LinkIncidentRiskInput{
			RiskID: risk.ID,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/incidents/"+incident.ID+"/risks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 409, resp.StatusCode)
	})

	t.Run("link with missing risk_id returns 400", func(t *testing.T) {
		input := models.LinkIncidentRiskInput{
			RiskID: "",
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/incidents/"+incident.ID+"/risks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("unlink risk from incident", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/incidents/"+incident.ID+"/risks/"+risk.ID, nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 204, resp.StatusCode)

		// Verify unlinking
		links, err := incidentRiskRepo.ListByIncident(ctx, incident.ID)
		require.NoError(t, err)
		assert.Len(t, links, 0)
	})

	t.Run("unlink non-existent link returns 404", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/incidents/"+incident.ID+"/risks/non-existent-risk", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
	})

	t.Run("list empty risks", func(t *testing.T) {
		// Create a new incident with no linked risks
		newIncident := &models.Incident{
			Title:       "Incident No Risks",
			Status:      models.IncidentStatusNew,
			Priority:    models.PriorityP3,
			ReporterID:  user.ID,
			CreatedBy:   user.ID,
			UpdatedBy:   user.ID,
		}
		err := incidentRepo.Create(ctx, newIncident)
		require.NoError(t, err)

		req := httptest.NewRequest("GET", "/incidents/"+newIncident.ID+"/risks", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var links []*models.IncidentRisk
		err = json.NewDecoder(resp.Body).Decode(&links)
		require.NoError(t, err)
		assert.Len(t, links, 0)
	})
}

func TestIncidentHandler_Integration_FullWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// This test simulates a complete incident lifecycle:
	// Create -> List -> Update (multiple times) -> Link Risk -> Unlink Risk -> Delete

	db := database.GetDB(database.New())
	incidentRepo := database.NewIncidentRepository(db)
	incidentCategoryRepo := database.NewIncidentCategoryRepository(db)
	incidentRiskRepo := database.NewIncidentRiskRepository(db)
	auditRepo := database.NewAuditLogRepository(db)
	userRepo := database.NewUserRepository(db)
	riskRepo := database.NewRiskRepository(db)
	catRepo := database.NewCategoryRepository(db)

	ctx := context.Background()

	// Setup: Create a test user
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        "incident-workflow-" + uuid.New().String() + "@example.com",
		PasswordHash: "hash",
		Name:         "Workflow Tester",
		Role:         models.RoleMember,
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	// Setup: Create a category for risk
	category := &models.CreateCategoryInput{
		Name:        "Workflow Category",
		Description: "Test category",
	}
	createdCategory, err := catRepo.Create(ctx, category)
	require.NoError(t, err)

	// Setup: Create a risk to link
	risk := &models.Risk{
		Title:       "Workflow Risk",
		Description: "Risk for workflow test",
		OwnerID:     user.ID,
		Status:      models.StatusOpen,
		Severity:    models.SeverityHigh,
		CategoryID:  &createdCategory.ID,
		CreatedBy:   user.ID,
		UpdatedBy:   user.ID,
	}
	err = riskRepo.Create(ctx, risk)
	require.NoError(t, err)

	app := fiber.New()
	incidentHandler := NewIncidentHandler(incidentRepo, incidentCategoryRepo, incidentRiskRepo, auditRepo)
	incidentRiskHandler := NewIncidentRiskHandler(incidentRiskRepo, auditRepo)

	app.Post("/incidents", integrationTestAuthMiddleware(user.ID), incidentHandler.Create)
	app.Get("/incidents/:id", incidentHandler.Get)
	app.Get("/incidents", incidentHandler.List)
	app.Put("/incidents/:id", integrationTestAuthMiddleware(user.ID), incidentHandler.Update)
	app.Delete("/incidents/:id", integrationTestAuthMiddleware(user.ID), incidentHandler.Delete)
	app.Post("/incidents/:incidentId/risks", integrationTestAuthMiddleware(user.ID), incidentRiskHandler.LinkRisk)
	app.Delete("/incidents/:incidentId/risks/:riskId", integrationTestAuthMiddleware(user.ID), incidentRiskHandler.UnlinkRisk)
	app.Get("/incidents/:incidentId/risks", incidentRiskHandler.ListRisks)

	// Step 1: Create incident
	t.Run("step 1 - create incident", func(t *testing.T) {
		input := models.CreateIncidentInput{
			Title:           "Production Outage",
			Description:     "API servers are down",
			Priority:        models.PriorityP1,
			Status:          models.IncidentStatusNew,
			ServiceAffected: "API Gateway",
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/incidents", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)
	})

	// Step 2: Get the incident
	var incidentID string
	t.Run("step 2 - get incident", func(t *testing.T) {
		// First list to get the ID
		req := httptest.NewRequest("GET", "/incidents?search=Production+Outage", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		var listResp models.IncidentListResponse
		err = json.NewDecoder(resp.Body).Decode(&listResp)
		require.NoError(t, err)
		require.Len(t, listResp.Data, 1)
		incidentID = listResp.Data[0].ID

		// Now get by ID
		req = httptest.NewRequest("GET", "/incidents/"+incidentID, nil)
		resp, err = app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	// Step 3: Update to in_progress
	t.Run("step 3 - update to in_progress", func(t *testing.T) {
		newStatus := models.IncidentStatusInProgress
		input := models.UpdateIncidentInput{
			Status: &newStatus,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("PUT", "/incidents/"+incidentID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	// Step 4: Link risk
	t.Run("step 4 - link risk", func(t *testing.T) {
		input := models.LinkIncidentRiskInput{
			RiskID: risk.ID,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/incidents/"+incidentID+"/risks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)
	})

	// Step 5: Verify linked risk
	t.Run("step 5 - verify linked risk", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/incidents/"+incidentID+"/risks", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		var links []*models.IncidentRisk
		err = json.NewDecoder(resp.Body).Decode(&links)
		require.NoError(t, err)
		assert.Len(t, links, 1)
		assert.Equal(t, risk.ID, links[0].RiskID)
	})

	// Step 6: Update to resolved
	t.Run("step 6 - update to resolved", func(t *testing.T) {
		newStatus := models.IncidentStatusResolved
		rootCause := "Database connection pool exhausted"
		resolutionNotes := "Increased connection pool size and added monitoring"
		input := models.UpdateIncidentInput{
			Status:          &newStatus,
			RootCause:       &rootCause,
			ResolutionNotes: &resolutionNotes,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("PUT", "/incidents/"+incidentID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var updated models.Incident
		err = json.NewDecoder(resp.Body).Decode(&updated)
		require.NoError(t, err)
		assert.Equal(t, models.IncidentStatusResolved, updated.Status)
		assert.NotNil(t, updated.ResolvedAt)
	})

	// Step 7: Unlink risk
	t.Run("step 7 - unlink risk", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/incidents/"+incidentID+"/risks/"+risk.ID, nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 204, resp.StatusCode)
	})

	// Step 8: Update to closed
	t.Run("step 8 - update to closed", func(t *testing.T) {
		newStatus := models.IncidentStatusClosed
		input := models.UpdateIncidentInput{
			Status: &newStatus,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("PUT", "/incidents/"+incidentID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	// Step 9: Delete incident
	t.Run("step 9 - delete incident", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/incidents/"+incidentID, nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 204, resp.StatusCode)

		// Verify deletion
		req = httptest.NewRequest("GET", "/incidents/"+incidentID, nil)
		resp, err = app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
	})
}

func TestIncidentHandler_Integration_Timestamps(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := database.GetDB(database.New())
	incidentRepo := database.NewIncidentRepository(db)
	incidentCategoryRepo := database.NewIncidentCategoryRepository(db)
	incidentRiskRepo := database.NewIncidentRiskRepository(db)
	auditRepo := database.NewAuditLogRepository(db)
	userRepo := database.NewUserRepository(db)

	ctx := context.Background()

	// Setup: Create a test user
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        "incident-timestamp-" + uuid.New().String() + "@example.com",
		PasswordHash: "hash",
		Name:         "Timestamp Tester",
		Role:         models.RoleMember,
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	app := fiber.New()
	handler := NewIncidentHandler(incidentRepo, incidentCategoryRepo, incidentRiskRepo, auditRepo)
	app.Post("/incidents", integrationTestAuthMiddleware(user.ID), handler.Create)

	t.Run("custom occurred_at and detected_at", func(t *testing.T) {
		occurredAt := time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
		detectedAt := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)

		input := models.CreateIncidentInput{
			Title:       "Timestamp Test",
			OccurredAt:  &occurredAt,
			DetectedAt:  &detectedAt,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/incidents", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)

		var incident models.Incident
		err = json.NewDecoder(resp.Body).Decode(&incident)
		require.NoError(t, err)

		// Verify timestamps are approximately correct (within 1 minute due to parsing)
		expectedOccurred, _ := time.Parse(time.RFC3339, occurredAt)
		expectedDetected, _ := time.Parse(time.RFC3339, detectedAt)

		assert.WithinDuration(t, expectedOccurred, incident.OccurredAt, time.Minute)
		assert.WithinDuration(t, expectedDetected, incident.DetectedAt, time.Minute)
	})

	t.Run("default timestamps when not provided", func(t *testing.T) {
		beforeCreate := time.Now()

		input := models.CreateIncidentInput{
			Title: "Default Timestamp Test",
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/incidents", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)

		var incident models.Incident
		err = json.NewDecoder(resp.Body).Decode(&incident)
		require.NoError(t, err)

		// Verify timestamps are set to approximately now
		assert.WithinDuration(t, beforeCreate, incident.OccurredAt, 5*time.Second)
		assert.WithinDuration(t, beforeCreate, incident.DetectedAt, 5*time.Second)
	})
}

func TestIncidentHandler_Integration_EmptyCategoryNormalization(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := database.GetDB(database.New())
	incidentRepo := database.NewIncidentRepository(db)
	incidentCategoryRepo := database.NewIncidentCategoryRepository(db)
	incidentRiskRepo := database.NewIncidentRiskRepository(db)
	auditRepo := database.NewAuditLogRepository(db)
	userRepo := database.NewUserRepository(db)

	ctx := context.Background()

	// Setup: Create a test user
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        "incident-emptycat-" + uuid.New().String() + "@example.com",
		PasswordHash: "hash",
		Name:         "Empty Cat Tester",
		Role:         models.RoleMember,
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	app := fiber.New()
	handler := NewIncidentHandler(incidentRepo, incidentCategoryRepo, incidentRiskRepo, auditRepo)
	app.Post("/incidents", integrationTestAuthMiddleware(user.ID), handler.Create)
	app.Put("/incidents/:id", integrationTestAuthMiddleware(user.ID), handler.Update)

	t.Run("empty category_id normalized to nil on create", func(t *testing.T) {
		emptyStr := ""
		input := models.CreateIncidentInput{
			Title:      "Empty Category Create",
			CategoryID: &emptyStr,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/incidents", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)

		var incident models.Incident
		err = json.NewDecoder(resp.Body).Decode(&incident)
		require.NoError(t, err)
		assert.Nil(t, incident.CategoryID)
	})

	t.Run("empty category_id normalized to nil on update", func(t *testing.T) {
		// Create incident first
		incident := &models.Incident{
			Title:       "To Update Category",
			Status:      models.IncidentStatusNew,
			Priority:    models.PriorityP3,
			ReporterID:  user.ID,
			CreatedBy:   user.ID,
			UpdatedBy:   user.ID,
		}
		err := incidentRepo.Create(ctx, incident)
		require.NoError(t, err)

		emptyStr := ""
		input := models.UpdateIncidentInput{
			CategoryID: &emptyStr,
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("PUT", "/incidents/"+incident.ID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var updated models.Incident
		err = json.NewDecoder(resp.Body).Decode(&updated)
		require.NoError(t, err)
		assert.Nil(t, updated.CategoryID)
	})
}
