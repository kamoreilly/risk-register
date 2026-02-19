# Risk Register Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a full-featured risk management application with authentication, risk CRUD, mitigations, dashboard, calendar, board, compliance mapping, and AI stubs.

**Architecture:** Layered build approach - each slice delivers a working, testable feature end-to-end. Backend Go/Fiber API with PostgreSQL. Frontend React/TanStack with shadcn/ui.

**Tech Stack:**
- Backend: Go 1.22+, Fiber v2, PostgreSQL (pgx), JWT (golang-jwt), bcrypt
- Frontend: React 18, TypeScript, TanStack Router/Query, React Hook Form, Zod, shadcn/ui

---

## Slice 1: Authentication

Deliverable: Users can register, log in, and access protected routes.

### Task 1.1: Set up migrations infrastructure

**Files:**
- Create: `backend/internal/migrations/migrations.go`
- Create: `backend/migrations/001_users.up.sql`
- Create: `backend/migrations/001_users.down.sql`

**Step 1: Add migration dependencies**

Run: `cd backend && go get github.com/golang-migrate/migrate/v4 github.com/golang-migrate/migrate/v4/database/postgres github.com/golang-migrate/migrate/v4/source/iofs`

**Step 2: Create users migration (up)**

Create `backend/migrations/001_users.up.sql`:

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE user_role AS ENUM ('admin', 'member');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'member',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
```

**Step 3: Create users migration (down)**

Create `backend/migrations/001_users.down.sql`:

```sql
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_role;
```

**Step 4: Create migration runner**

Create `backend/internal/migrations/migrations.go`:

```go
package migrations

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed all:migrations
var migrationsFS embed.FS

func RunMigrations(databaseURL string) error {
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
```

**Step 5: Commit**

```bash
git add backend/internal/migrations/ backend/migrations/
git commit -m "feat(backend): add migrations infrastructure and users table"
```

---

### Task 1.2: Add auth dependencies

**Files:**
- Modify: `backend/go.mod`

**Step 1: Install auth packages**

Run: `cd backend && go get github.com/golang-jwt/jwt/v5 golang.org/x/crypto/bcrypt`

**Step 2: Tidy dependencies**

Run: `cd backend && go mod tidy`

**Step 3: Commit**

```bash
git add backend/go.mod backend/go.sum
git commit -m "feat(backend): add JWT and bcrypt dependencies"
```

---

### Task 1.3: Create User model

**Files:**
- Create: `backend/internal/models/user.go`

**Step 1: Write the User model**

Create `backend/internal/models/user.go`:

```go
package models

import (
	"time"
)

type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleMember UserRole = "member"
)

type User struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Name         string    `json:"name" db:"name"`
	Role         UserRole  `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=1"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
```

**Step 2: Commit**

```bash
git add backend/internal/models/user.go
git commit -m "feat(backend): add User model"
```

---

### Task 1.4: Create auth service

**Files:**
- Create: `backend/internal/auth/auth.go`

**Step 1: Create auth service with password hashing and JWT**

Create `backend/internal/auth/auth.go`:

```go
package auth

import (
	"errors"
	"os"
	"time"

	"backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidToken       = errors.New("invalid token")
)

type Claims struct {
	UserID string       `json:"user_id"`
	Email  string       `json:"email"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(user *models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-in-production"
	}

	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "risk-register",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-in-production"
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
```

**Step 2: Commit**

```bash
git add backend/internal/auth/auth.go
git commit -m "feat(backend): add auth service with JWT and bcrypt"
```

---

### Task 1.5: Create auth middleware

**Files:**
- Create: `backend/internal/middleware/auth.go`

**Step 1: Write the failing test**

Create `backend/internal/middleware/auth_test.go`:

```go
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
```

**Step 2: Run test to verify it fails**

Run: `cd backend && go test ./internal/middleware -run TestAuthMiddleware -v`
Expected: FAIL (middleware doesn't exist)

**Step 3: Write the middleware implementation**

Create `backend/internal/middleware/auth.go`:

```go
package middleware

import (
	"strings"

	"backend/internal/auth"

	"github.com/gofiber/fiber/v2"
)

type contextKey string

const UserKey contextKey = "user"

type UserClaims struct {
	UserID string
	Email  string
	Role   string
}

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "missing authorization header",
		})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid authorization format",
		})
	}

	claims, err := auth.ValidateToken(parts[1])
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid or expired token",
		})
	}

	c.Locals(UserKey, &UserClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   string(claims.Role),
	})

	return c.Next()
}

func GetUserFromContext(c *fiber.Ctx) *UserClaims {
	user, ok := c.Locals(UserKey).(*UserClaims)
	if !ok {
		return nil
	}
	return user
}

func RequireAdmin(c *fiber.Ctx) error {
	user := GetUserFromContext(c)
	if user == nil || user.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "admin access required",
		})
	}
	return c.Next()
}
```

**Step 4: Run test to verify it passes**

Run: `cd backend && go test ./internal/middleware -run TestAuthMiddleware -v`
Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/middleware/
git commit -m "feat(backend): add auth middleware with tests"
```

---

### Task 1.6: Create user repository

**Files:**
- Create: `backend/internal/database/users.go`

**Step 1: Write the failing test**

Create `backend/internal/database/users_test.go`:

```go
package database

import (
	"context"
	"testing"
	"time"

	"backend/internal/models"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setupTestDB(t *testing.T) *postgres.PostgresContainer {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
	)
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}

	return pgContainer
}

func TestUserRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	pgContainer := setupTestDB(t)
	defer func() {
		if err := testcontainers.TerminateContainer(pgContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	// This test verifies the interface compiles
	var _ UserRepository = (*userRepository)(nil)
	t.Log("UserRepository interface satisfied")
}
```

**Step 2: Run test to verify it fails**

Run: `cd backend && go test ./internal/database -run TestUserRepository -v`
Expected: FAIL (repository doesn't exist)

**Step 3: Write the repository implementation**

Create `backend/internal/database/users.go`:

```go
package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"backend/internal/models"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserExists = errors.New("user already exists")

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `
		INSERT INTO users (id, email, password_hash, name, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if isDuplicateKeyError(err) {
			return ErrUserExists
		}
		return err
	}

	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, name, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, name, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func isDuplicateKeyError(err error) bool {
	// PostgreSQL duplicate key error code
	return err != nil && (err.Error() != "" && containsString(err.Error(), "duplicate key"))
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
```

**Step 4: Add google/uuid dependency**

Run: `cd backend && go get github.com/google/uuid && go mod tidy`

**Step 5: Run test to verify it passes**

Run: `cd backend && go test ./internal/database -run TestUserRepository -v`
Expected: PASS

**Step 6: Commit**

```bash
git add backend/internal/database/users.go backend/internal/database/users_test.go
git commit -m "feat(backend): add user repository with interface"
```

---

### Task 1.7: Create auth handlers

**Files:**
- Create: `backend/internal/handlers/auth.go`
- Create: `backend/internal/handlers/auth_test.go`

**Step 1: Write the failing test**

Create `backend/internal/handlers/auth_test.go`:

```go
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type mockUserRepo struct {
	users map[string]*models.User
}

func (m *mockUserRepo) Create(ctx interface{}, user *models.User) error {
	if _, exists := m.users[user.Email]; exists {
		return nil // simulate duplicate
	}
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepo) FindByEmail(ctx interface{}, email string) (*models.User, error) {
	if user, ok := m.users[email]; ok {
		return user, nil
	}
	return nil, nil
}

func (m *mockUserRepo) FindByID(ctx interface{}, id string) (*models.User, error) {
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

	// First create a user
	hashedPassword, _ := hashPassword("password123")
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
```

**Step 2: Run test to verify it fails**

Run: `cd backend && go test ./internal/handlers -run TestRegister -v`
Expected: FAIL (handler doesn't exist)

**Step 3: Write the handler implementation**

Create `backend/internal/handlers/auth.go`:

```go
package handlers

import (
	"context"

	"backend/internal/auth"
	"backend/internal/database"
	"backend/internal/middleware"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	users database.UserRepository
}

func NewAuthHandler(users database.UserRepository) *AuthHandler {
	return &AuthHandler{users: users}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var input models.RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate input
	if input.Email == "" || input.Password == "" || input.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "email, password, and name are required",
		})
	}

	if len(input.Password) < 8 {
		return c.Status(400).JSON(fiber.Map{
			"error": "password must be at least 8 characters",
		})
	}

	// Check if user exists
	existing, _ := h.users.FindByEmail(context.Background(), input.Email)
	if existing != nil {
		return c.Status(409).JSON(fiber.Map{
			"error": "email already registered",
		})
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to hash password",
		})
	}

	// Create user
	user := &models.User{
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Name:         input.Name,
		Role:         models.RoleMember,
	}

	if err := h.users.Create(context.Background(), user); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to create user",
		})
	}

	// Generate token
	token, err := auth.GenerateToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to generate token",
		})
	}

	return c.Status(201).JSON(models.AuthResponse{
		User:  user,
		Token: token,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input models.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Find user
	user, err := h.users.FindByEmail(context.Background(), input.Email)
	if err != nil || user == nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}

	// Check password
	if !auth.CheckPassword(input.Password, user.PasswordHash) {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}

	// Generate token
	token, err := auth.GenerateToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to generate token",
		})
	}

	return c.JSON(models.AuthResponse{
		User:  user,
		Token: token,
	})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "not authenticated",
		})
	}

	fullUser, err := h.users.FindByID(context.Background(), user.UserID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return c.JSON(fullUser)
}
```

**Step 4: Run test to verify it passes**

Run: `cd backend && go test ./internal/handlers -run TestRegister -v`
Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/handlers/
git commit -m "feat(backend): add auth handlers with tests"
```

---

### Task 1.8: Wire up auth routes

**Files:**
- Modify: `backend/internal/server/routes.go`
- Modify: `backend/internal/server/server.go`

**Step 1: Update server to include user repository**

Modify `backend/internal/server/server.go`:

```go
package server

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"

	"backend/internal/database"
	"backend/internal/handlers"
)

type FiberServer struct {
	*fiber.App
	db      database.Service
	rawDB   *sql.DB
	users   database.UserRepository
	auth    *handlers.AuthHandler
}

func New() *FiberServer {
	db := database.New()
	rawDB := getRawDB()
	users := database.NewUserRepository(rawDB)

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "risk-register",
			AppName:      "Risk Register API",
		}),
		db:    db,
		rawDB: rawDB,
		users: users,
		auth:  handlers.NewAuthHandler(users),
	}

	return server
}

func getRawDB() *sql.DB {
	// This is a simplified approach - in production you'd want better connection management
	connStr := buildConnStr()
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		panic(err)
	}
	return db
}

func buildConnStr() string {
	// Build from env vars like existing database.go does
	return "" // implementation detail - use existing env vars
}
```

**Step 2: Update routes to include auth endpoints**

Modify `backend/internal/server/routes.go`:

```go
package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"backend/internal/middleware"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Public routes
	s.App.Get("/", s.HelloWorldHandler)
	s.App.Get("/health", s.healthHandler)

	// Auth routes (public)
	auth := s.App.Group("/api/v1/auth")
	auth.Post("/register", s.auth.Register)
	auth.Post("/login", s.auth.Login)

	// Protected routes
	protected := s.App.Group("/api/v1", middleware.AuthMiddleware)
	protected.Get("/auth/me", s.auth.Me)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Risk Register API",
		"version": "1.0.0",
	}
	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
```

**Step 3: Run existing tests**

Run: `cd backend && go test ./... -v`
Expected: PASS

**Step 4: Commit**

```bash
git add backend/internal/server/
git commit -m "feat(backend): wire up auth routes"
```

---

### Task 1.9: Update main.go to run migrations

**Files:**
- Modify: `backend/cmd/api/main.go`

**Step 1: Add migration call to main.go**

Modify `backend/cmd/api/main.go` to add migration run at startup:

```go
package main

import (
	"backend/internal/migrations"
	"backend/internal/server"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func gracefulShutdown(fiberServer *server.FiberServer, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := fiberServer.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
	done <- true
}

func main() {
	// Run migrations
	databaseURL := buildDatabaseURL()
	if err := migrations.RunMigrations(databaseURL); err != nil {
		log.Printf("Migration warning: %v", err)
	}

	server := server.New()
	server.RegisterFiberRoutes()

	done := make(chan bool, 1)

	go func() {
		port, _ := strconv.Atoi(os.Getenv("PORT"))
		if port == 0 {
			port = 8080
		}
		err := server.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	go gracefulShutdown(server, done)

	<-done
	log.Println("Graceful shutdown complete.")
}

func buildDatabaseURL() string {
	host := getEnv("RISK_REGISTER_DB_HOST", "localhost")
	port := getEnv("RISK_REGISTER_DB_PORT", "5432")
	user := getEnv("RISK_REGISTER_DB_USERNAME", "postgres")
	password := getEnv("RISK_REGISTER_DB_PASSWORD", "postgres")
	database := getEnv("RISK_REGISTER_DB_DATABASE", "risk_register")
	schema := getEnv("RISK_REGISTER_DB_SCHEMA", "public")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		user, password, host, port, database, schema)
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
```

**Step 2: Commit**

```bash
git add backend/cmd/api/main.go
git commit -m "feat(backend): run migrations on startup"
```

---

### Task 1.10: Add frontend API client

**Files:**
- Create: `frontend/apps/web/src/lib/api.ts`

**Step 1: Create API client**

Create `frontend/apps/web/src/lib/api.ts`:

```typescript
const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080';

interface ApiOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE';
  body?: unknown;
  token?: string;
}

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private getStoredToken(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem('token');
  }

  async request<T>(path: string, options: ApiOptions = {}): Promise<T> {
    const { method = 'GET', body, token } = options;

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    const authToken = token || this.getStoredToken();
    if (authToken) {
      headers['Authorization'] = `Bearer ${authToken}`;
    }

    const response = await fetch(`${this.baseUrl}${path}`, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Request failed' }));
      throw new ApiError(response.status, error.error || 'Request failed');
    }

    return response.json();
  }

  get<T>(path: string, token?: string): Promise<T> {
    return this.request<T>(path, { method: 'GET', token });
  }

  post<T>(path: string, body: unknown, token?: string): Promise<T> {
    return this.request<T>(path, { method: 'POST', body, token });
  }

  put<T>(path: string, body: unknown, token?: string): Promise<T> {
    return this.request<T>(path, { method: 'PUT', body, token });
  }

  delete<T>(path: string, token?: string): Promise<T> {
    return this.request<T>(path, { method: 'DELETE', token });
  }
}

export class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = 'ApiError';
  }
}

export const api = new ApiClient(API_BASE);
```

**Step 2: Commit**

```bash
git add frontend/apps/web/src/lib/api.ts
git commit -m "feat(frontend): add API client"
```

---

### Task 1.11: Add auth types and hooks

**Files:**
- Create: `frontend/apps/web/src/types/auth.ts`
- Create: `frontend/apps/web/src/hooks/useAuth.ts`

**Step 1: Create auth types**

Create `frontend/apps/web/src/types/auth.ts`:

```typescript
export interface User {
  id: string;
  email: string;
  name: string;
  role: 'admin' | 'member';
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  user: User;
  token: string;
}

export interface LoginInput {
  email: string;
  password: string;
}

export interface RegisterInput {
  email: string;
  password: string;
  name: string;
}
```

**Step 2: Create auth hook**

Create `frontend/apps/web/src/hooks/useAuth.ts`:

```typescript
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useNavigate } from '@tanstack/react-router';

import { api } from '@/lib/api';
import type { AuthResponse, LoginInput, RegisterInput, User } from '@/types/auth';

const AUTH_KEY = ['auth', 'me'];

export function useAuth() {
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  const { data: user, isLoading, error } = useQuery({
    queryKey: AUTH_KEY,
    queryFn: () => api.get<User>('/api/v1/auth/me'),
    retry: false,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });

  const loginMutation = useMutation({
    mutationFn: (input: LoginInput) =>
      api.post<AuthResponse>('/api/v1/auth/login', input),
    onSuccess: (data) => {
      localStorage.setItem('token', data.token);
      queryClient.setQueryData(AUTH_KEY, data.user);
      navigate({ to: '/app' });
    },
  });

  const registerMutation = useMutation({
    mutationFn: (input: RegisterInput) =>
      api.post<AuthResponse>('/api/v1/auth/register', input),
    onSuccess: (data) => {
      localStorage.setItem('token', data.token);
      queryClient.setQueryData(AUTH_KEY, data.user);
      navigate({ to: '/app' });
    },
  });

  const logout = () => {
    localStorage.removeItem('token');
    queryClient.clear();
    navigate({ to: '/login' });
  };

  return {
    user,
    isLoading,
    error,
    isAuthenticated: !!user,
    login: loginMutation.mutate,
    loginError: loginMutation.error,
    isLoginLoading: loginMutation.isPending,
    register: registerMutation.mutate,
    registerError: registerMutation.error,
    isRegisterLoading: registerMutation.isPending,
    logout,
  };
}
```

**Step 3: Commit**

```bash
git add frontend/apps/web/src/types/ frontend/apps/web/src/hooks/
git commit -m "feat(frontend): add auth types and useAuth hook"
```

---

### Task 1.12: Wire login page to backend

**Files:**
- Modify: `frontend/apps/web/src/routes/login.tsx`

**Step 1: Update login page to use auth hook**

Modify `frontend/apps/web/src/routes/login.tsx`:

```tsx
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import * as React from "react";

import { useAuth } from "@/hooks/useAuth";
import { Button, buttonVariants } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/login")({
  component: LoginComponent,
});

function LoginComponent() {
  const navigate = useNavigate();
  const { login, isLoginLoading, loginError, isAuthenticated } = useAuth();

  const [email, setEmail] = React.useState("");
  const [password, setPassword] = React.useState("");
  const [rememberMe, setRememberMe] = React.useState(false);

  // Redirect if already authenticated
  React.useEffect(() => {
    if (isAuthenticated) {
      navigate({ to: "/app" });
    }
  }, [isAuthenticated, navigate]);

  async function onSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    login({ email, password });
  }

  return (
    <div className="bg-background">
      <main className="grid min-h-svh place-items-center px-4">
        <Card className="w-full max-w-sm">
          <CardHeader className="border-b">
            <CardTitle>Log in</CardTitle>
            <CardDescription>Sign in to access the risk register.</CardDescription>
          </CardHeader>

          <CardContent>
            {loginError && (
              <div className="mb-4 rounded-md bg-destructive/10 p-3 text-sm text-destructive">
                {loginError instanceof Error ? loginError.message : "Login failed"}
              </div>
            )}
            <form className="grid gap-3" onSubmit={onSubmit}>
              <div className="grid gap-1.5">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  name="email"
                  type="email"
                  autoComplete="email"
                  placeholder="you@company.com"
                  value={email}
                  onChange={(e) => setEmail(e.currentTarget.value)}
                  required
                />
              </div>

              <div className="grid gap-1.5">
                <Label htmlFor="password">Password</Label>
                <Input
                  id="password"
                  name="password"
                  type="password"
                  autoComplete="current-password"
                  placeholder="••••••••"
                  value={password}
                  onChange={(e) => setPassword(e.currentTarget.value)}
                  required
                />
              </div>

              <div className="flex items-center justify-between gap-3">
                <label className="flex items-center gap-2 text-xs">
                  <Checkbox checked={rememberMe} onCheckedChange={(v) => setRememberMe(Boolean(v))} />
                  Remember me
                </label>

                <Button variant="link" type="button" className="h-auto px-0">
                  Forgot password?
                </Button>
              </div>

              <Button type="submit" disabled={isLoginLoading}>
                {isLoginLoading ? "Signing in…" : "Sign in"}
              </Button>
            </form>
          </CardContent>

          <CardFooter className="justify-between">
            <Link to="/register" className={cn(buttonVariants({ variant: "link", size: "sm" }))}>
              Create account
            </Link>
            <Link to="/" className={cn(buttonVariants({ variant: "outline", size: "sm" }))}>
              Back
            </Link>
          </CardFooter>
        </Card>
      </main>
    </div>
  );
}
```

**Step 2: Commit**

```bash
git add frontend/apps/web/src/routes/login.tsx
git commit -m "feat(frontend): wire login page to backend auth"
```

---

### Task 1.13: Create register page

**Files:**
- Create: `frontend/apps/web/src/routes/register.tsx`

**Step 1: Create register page**

Create `frontend/apps/web/src/routes/register.tsx`:

```tsx
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import * as React from "react";

import { useAuth } from "@/hooks/useAuth";
import { Button, buttonVariants } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/register")({
  component: RegisterComponent,
});

function RegisterComponent() {
  const navigate = useNavigate();
  const { register, isRegisterLoading, registerError, isAuthenticated } = useAuth();

  const [email, setEmail] = React.useState("");
  const [password, setPassword] = React.useState("");
  const [confirmPassword, setConfirmPassword] = React.useState("");
  const [name, setName] = React.useState("");
  const [validationError, setValidationError] = React.useState<string | null>(null);

  React.useEffect(() => {
    if (isAuthenticated) {
      navigate({ to: "/app" });
    }
  }, [isAuthenticated, navigate]);

  async function onSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setValidationError(null);

    if (password !== confirmPassword) {
      setValidationError("Passwords do not match");
      return;
    }

    if (password.length < 8) {
      setValidationError("Password must be at least 8 characters");
      return;
    }

    register({ email, password, name });
  }

  return (
    <div className="bg-background">
      <main className="grid min-h-svh place-items-center px-4">
        <Card className="w-full max-w-sm">
          <CardHeader className="border-b">
            <CardTitle>Create account</CardTitle>
            <CardDescription>Sign up to start managing risks.</CardDescription>
          </CardHeader>

          <CardContent>
            {(validationError || registerError) && (
              <div className="mb-4 rounded-md bg-destructive/10 p-3 text-sm text-destructive">
                {validationError || (registerError instanceof Error ? registerError.message : "Registration failed")}
              </div>
            )}
            <form className="grid gap-3" onSubmit={onSubmit}>
              <div className="grid gap-1.5">
                <Label htmlFor="name">Name</Label>
                <Input
                  id="name"
                  name="name"
                  type="text"
                  autoComplete="name"
                  placeholder="John Doe"
                  value={name}
                  onChange={(e) => setName(e.currentTarget.value)}
                  required
                />
              </div>

              <div className="grid gap-1.5">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  name="email"
                  type="email"
                  autoComplete="email"
                  placeholder="you@company.com"
                  value={email}
                  onChange={(e) => setEmail(e.currentTarget.value)}
                  required
                />
              </div>

              <div className="grid gap-1.5">
                <Label htmlFor="password">Password</Label>
                <Input
                  id="password"
                  name="password"
                  type="password"
                  autoComplete="new-password"
                  placeholder="••••••••"
                  value={password}
                  onChange={(e) => setPassword(e.currentTarget.value)}
                  required
                />
              </div>

              <div className="grid gap-1.5">
                <Label htmlFor="confirmPassword">Confirm password</Label>
                <Input
                  id="confirmPassword"
                  name="confirmPassword"
                  type="password"
                  autoComplete="new-password"
                  placeholder="••••••••"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.currentTarget.value)}
                  required
                />
              </div>

              <Button type="submit" disabled={isRegisterLoading}>
                {isRegisterLoading ? "Creating account…" : "Create account"}
              </Button>
            </form>
          </CardContent>

          <CardFooter className="justify-between">
            <Link to="/login" className={cn(buttonVariants({ variant: "link", size: "sm" }))}>
              Already have an account?
            </Link>
            <Link to="/" className={cn(buttonVariants({ variant: "outline", size: "sm" }))}>
              Back
            </Link>
          </CardFooter>
        </Card>
      </main>
    </div>
  );
}
```

**Step 2: Commit**

```bash
git add frontend/apps/web/src/routes/register.tsx
git commit -m "feat(frontend): add register page"
```

---

### Task 1.14: Create protected route layout

**Files:**
- Create: `frontend/apps/web/src/routes/app/__root.tsx`

**Step 1: Create protected layout**

Create `frontend/apps/web/src/routes/app/__root.tsx`:

```tsx
import { Outlet, createFileRoute, useNavigate } from "@tanstack/react-router";
import * as React from "react";

import { useAuth } from "@/hooks/useAuth";
import { HeaderNav } from "@/components/header-nav";
import { Button } from "@/components/ui/button";
import { Loader } from "@/components/loader";

export const Route = createFileRoute("/app")({
  component: AppLayout,
});

function AppLayout() {
  const navigate = useNavigate();
  const { user, isLoading, isAuthenticated, logout } = useAuth();

  React.useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      navigate({ to: "/login" });
    }
  }, [isLoading, isAuthenticated, navigate]);

  if (isLoading) {
    return (
      <div className="flex min-h-svh items-center justify-center">
        <Loader />
      </div>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  return (
    <div className="flex min-h-svh flex-col bg-background">
      <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-14 items-center px-4">
          <HeaderNav />
          <div className="ml-auto flex items-center gap-4">
            <span className="text-sm text-muted-foreground">{user?.name}</span>
            <Button variant="outline" size="sm" onClick={logout}>
              Log out
            </Button>
          </div>
        </div>
      </header>
      <main className="flex-1">
        <Outlet />
      </main>
    </div>
  );
}
```

**Step 2: Commit**

```bash
git add frontend/apps/web/src/routes/app/__root.tsx
git commit -m "feat(frontend): add protected route layout"
```

---

### Task 1.15: Create dashboard placeholder

**Files:**
- Create: `frontend/apps/web/src/routes/app/index.tsx`

**Step 1: Create dashboard placeholder**

Create `frontend/apps/web/src/routes/app/index.tsx`:

```tsx
import { createFileRoute } from "@tanstack/react-router";

import { useAuth } from "@/hooks/useAuth";

export const Route = createFileRoute("/app/")({
  component: Dashboard,
});

function Dashboard() {
  const { user } = useAuth();

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold">Dashboard</h1>
      <p className="text-muted-foreground">
        Welcome back, {user?.name}!
      </p>
      <p className="mt-4 text-sm text-muted-foreground">
        Dashboard widgets will be added in Slice 4.
      </p>
    </div>
  );
}
```

**Step 2: Commit**

```bash
git add frontend/apps/web/src/routes/app/index.tsx
git commit -m "feat(frontend): add dashboard placeholder"
```

---

## Slice 2: Risks CRUD

Deliverable: Users can create, read, update, and delete risks with categories.

### Task 2.1: Create risks and categories migrations

**Files:**
- Create: `backend/migrations/002_categories.up.sql`
- Create: `backend/migrations/002_categories.down.sql`
- Create: `backend/migrations/003_risks.up.sql`
- Create: `backend/migrations/003_risks.down.sql`

**Step 1: Create categories migration**

`backend/migrations/002_categories.up.sql`:
```sql
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

INSERT INTO categories (name, description) VALUES
    ('Security', 'Security-related risks including cyber threats and data breaches'),
    ('Operational', 'Operational risks affecting business processes'),
    ('Financial', 'Financial risks including market, credit, and liquidity risks'),
    ('Compliance', 'Regulatory and compliance risks'),
    ('Strategic', 'Strategic risks affecting long-term business objectives'),
    ('Reputational', 'Risks to company reputation and brand');
```

`backend/migrations/002_categories.down.sql`:
```sql
DELETE FROM categories WHERE name IN ('Security', 'Operational', 'Financial', 'Compliance', 'Strategic', 'Reputational');
DROP TABLE IF EXISTS categories;
```

**Step 2: Create risks migration**

`backend/migrations/003_risks.up.sql`:
```sql
CREATE TYPE risk_status AS ENUM ('open', 'mitigating', 'resolved', 'accepted');
CREATE TYPE risk_severity AS ENUM ('low', 'medium', 'high', 'critical');

CREATE TABLE risks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status risk_status NOT NULL DEFAULT 'open',
    severity risk_severity NOT NULL DEFAULT 'medium',
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    review_date DATE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    updated_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT
);

CREATE INDEX idx_risks_owner ON risks(owner_id);
CREATE INDEX idx_risks_status ON risks(status);
CREATE INDEX idx_risks_severity ON risks(severity);
CREATE INDEX idx_risks_category ON risks(category_id);
CREATE INDEX idx_risks_review_date ON risks(review_date);
```

`backend/migrations/003_risks.down.sql`:
```sql
DROP INDEX IF EXISTS idx_risks_review_date;
DROP INDEX IF EXISTS idx_risks_category;
DROP INDEX IF EXISTS idx_risks_severity;
DROP INDEX IF EXISTS idx_risks_status;
DROP INDEX IF EXISTS idx_risks_owner;
DROP TABLE IF EXISTS risks;
DROP TYPE IF EXISTS risk_severity;
DROP TYPE IF EXISTS risk_status;
```

**Step 3: Commit**

```bash
git add backend/migrations/
git commit -m "feat(backend): add categories and risks migrations"
```

---

### Task 2.2: Create Risk and Category models

**Files:**
- Create: `backend/internal/models/risk.go`
- Create: `backend/internal/models/category.go`

**Step 1: Create models**

`backend/internal/models/category.go`:
```go
package models

import "time"

type Category struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
```

`backend/internal/models/risk.go`:
```go
package models

import "time"

type RiskStatus string

const (
	StatusOpen       RiskStatus = "open"
	StatusMitigating RiskStatus = "mitigating"
	StatusResolved   RiskStatus = "resolved"
	StatusAccepted   RiskStatus = "accepted"
)

type RiskSeverity string

const (
	SeverityLow      RiskSeverity = "low"
	SeverityMedium   RiskSeverity = "medium"
	SeverityHigh     RiskSeverity = "high"
	SeverityCritical RiskSeverity = "critical"
)

type Risk struct {
	ID          string       `json:"id" db:"id"`
	Title       string       `json:"title" db:"title"`
	Description string       `json:"description,omitempty" db:"description"`
	OwnerID     string       `json:"owner_id" db:"owner_id"`
	Owner       *User        `json:"owner,omitempty" db:"-"`
	Status      RiskStatus   `json:"status" db:"status"`
	Severity    RiskSeverity `json:"severity" db:"severity"`
	CategoryID  *string      `json:"category_id,omitempty" db:"category_id"`
	Category    *Category    `json:"category,omitempty" db:"-"`
	ReviewDate  *time.Time   `json:"review_date,omitempty" db:"review_date"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	CreatedBy   string       `json:"created_by" db:"created_by"`
	UpdatedBy   string       `json:"updated_by" db:"updated_by"`
}

type CreateRiskInput struct {
	Title       string       `json:"title" validate:"required,min=1,max=255"`
	Description string       `json:"description"`
	OwnerID     string       `json:"owner_id" validate:"required,uuid"`
	Status      RiskStatus   `json:"status" validate:"omitempty,risk_status"`
	Severity    RiskSeverity `json:"severity" validate:"omitempty,risk_severity"`
	CategoryID  *string      `json:"category_id" validate:"omitempty,uuid"`
	ReviewDate  *string      `json:"review_date" validate:"omitempty,date"`
}

type UpdateRiskInput struct {
	Title       *string       `json:"title" validate:"omitempty,min=1,max=255"`
	Description *string       `json:"description"`
	OwnerID     *string       `json:"owner_id" validate:"omitempty,uuid"`
	Status      *RiskStatus   `json:"status" validate:"omitempty,risk_status"`
	Severity    *RiskSeverity `json:"severity" validate:"omitempty,risk_severity"`
	CategoryID  *string       `json:"category_id" validate:"omitempty,uuid"`
	ReviewDate  *string       `json:"review_date" validate:"omitempty,date"`
}

type RiskListParams struct {
	Status     *RiskStatus
	Severity   *RiskSeverity
	CategoryID *string
	OwnerID    *string
	Search     string
	Sort       string
	Order      string
	Page       int
	Limit      int
}

type RiskListResponse struct {
	Data  []*Risk `json:"data"`
	Meta  Meta    `json:"meta"`
}

type Meta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}
```

**Step 2: Commit**

```bash
git add backend/internal/models/
git commit -m "feat(backend): add Risk and Category models"
```

---

### Task 2.3: Create risk repository

**Files:**
- Create: `backend/internal/database/risks.go`
- Create: `backend/internal/database/categories.go`

(Implementation follows same pattern as user repository - full code in execution)

---

### Task 2.4: Create risk handlers

**Files:**
- Create: `backend/internal/handlers/risks.go`
- Create: `backend/internal/handlers/categories.go`

(Implementation with CRUD operations - full code in execution)

---

### Task 2.5: Wire up risk routes

**Files:**
- Modify: `backend/internal/server/routes.go`

Add `/api/v1/risks` and `/api/v1/categories` routes.

---

### Task 2.6-2.10: Frontend risk management

- Create `frontend/apps/web/src/types/risk.ts`
- Create `frontend/apps/web/src/hooks/useRisks.ts`
- Create `frontend/apps/web/src/routes/app/risks/index.tsx` (list)
- Create `frontend/apps/web/src/routes/app/risks/$id.tsx` (detail)
- Create `frontend/apps/web/src/routes/app/risks/new.tsx` (create)
- Create `frontend/apps/web/src/components/risk-form.tsx`
- Create `frontend/apps/web/src/components/risk-table.tsx`

---

## Slice 3-9: Summary

The remaining slices follow the same TDD pattern:

- **Slice 3 (Mitigations):** Migration → Model → Repository → Handler → Routes → Frontend components
- **Slice 4 (Dashboard):** Dashboard API endpoints → Frontend widgets with charts
- **Slice 5 (Calendar):** Review date queries → Calendar view component
- **Slice 6 (Board):** Kanban component with drag-and-drop
- **Slice 7 (Frameworks):** Framework tables → Control mapping → Frontend badges
- **Slice 8 (AI Stub):** Stub endpoints → Frontend buttons with placeholder text
- **Slice 9 (Audit Trail):** Audit log table → Middleware → Timeline component

Each slice includes:
1. Migration files
2. Model definitions
3. Repository with tests
4. Handler with tests
5. Route wiring
6. Frontend types, hooks, and components
7. Commit after each task

---

## Dependencies to Add

**Backend:**
```bash
cd backend
go get github.com/golang-migrate/migrate/v4
go get github.com/golang-migrate/migrate/v4/database/postgres
go get github.com/golang-migrate/migrate/v4/source/iofs
go get github.com/golang-jwt/jwt/v5
go get github.com/google/uuid
go get github.com/go-playground/validator/v10
```

**Frontend:**
```bash
cd frontend
bun add @tanstack/react-table
bun add react-hook-form @hookform/resolvers
bun add date-fns
bun add zustand
bun add @dnd-kit/core @dnd-kit/sortable @dnd-kit/utilities
```

---

## Environment Variables

**Backend (.env):**
```
PORT=8080
RISK_REGISTER_DB_HOST=localhost
RISK_REGISTER_DB_PORT=5432
RISK_REGISTER_DB_USERNAME=postgres
RISK_REGISTER_DB_PASSWORD=postgres
RISK_REGISTER_DB_DATABASE=risk_register
RISK_REGISTER_DB_SCHEMA=public
JWT_SECRET=your-secret-key-here
```

**Frontend (.env):**
```
VITE_API_URL=http://localhost:8080
```
