# Incident Management Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add incident management and reporting capabilities to the risk register application.

**Architecture:** Parallel entity model - incidents as a first-class entity alongside risks with many-to-many risk linkage. Uses existing patterns for models, database layer, handlers, and frontend components.

**Tech Stack:** Go (Fiber), PostgreSQL, React (TanStack Start), TanStack Query, shadcn/ui, Tailwind CSS

---

## Phase 1: Database Schema

### Task 1: Add responder role to user_role enum

**Files:**
- Create: `backend/internal/migrations/migrations/010_incident_role.up.sql`
- Create: `backend/internal/migrations/migrations/010_incident_role.down.sql`

**Step 1: Create up migration**

```sql
-- backend/internal/migrations/migrations/010_incident_role.up.sql
ALTER TYPE user_role ADD VALUE IF NOT EXISTS 'responder';
```

**Step 2: Create down migration**

```sql
-- backend/internal/migrations/migrations/010_incident_role.down.sql
-- Note: PostgreSQL doesn't support removing enum values directly
-- This is a no-op migration; role will remain but won't be used
-- To fully remove, would need to recreate the enum type
```

**Step 3: Commit**

```bash
git add backend/internal/migrations/migrations/010_incident_role.*
git commit -m "feat(db): add responder role to user_role enum"
```

---

### Task 2: Create incident categories table

**Files:**
- Create: `backend/internal/migrations/migrations/011_incident_categories.up.sql`
- Create: `backend/internal/migrations/migrations/011_incident_categories.down.sql`

**Step 1: Create up migration**

```sql
-- backend/internal/migrations/migrations/011_incident_categories.up.sql
CREATE TABLE incident_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_incident_categories_name ON incident_categories(name);

-- Insert default categories
INSERT INTO incident_categories (name, description) VALUES
    ('Outage', 'System or service unavailability'),
    ('Breach', 'Security breach or data compromise'),
    ('Error', 'Application error or bug'),
    ('External', 'External event or third-party issue'),
    ('Performance', 'Performance degradation or slowdown');
```

**Step 2: Create down migration**

```sql
-- backend/internal/migrations/migrations/011_incident_categories.down.sql
DROP TABLE IF EXISTS incident_categories;
```

**Step 3: Commit**

```bash
git add backend/internal/migrations/migrations/011_incident_categories.*
git commit -m "feat(db): add incident_categories table with seed data"
```

---

### Task 3: Create incidents and incident_risks tables

**Files:**
- Create: `backend/internal/migrations/migrations/012_incidents.up.sql`
- Create: `backend/internal/migrations/migrations/012_incidents.down.sql`

**Step 1: Create up migration**

```sql
-- backend/internal/migrations/migrations/012_incidents.up.sql
-- Create enums
CREATE TYPE incident_status AS ENUM ('new', 'acknowledged', 'in_progress', 'on_hold', 'resolved', 'closed');
CREATE TYPE incident_priority AS ENUM ('p1', 'p2', 'p3', 'p4');

-- Create incidents table
CREATE TABLE incidents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category_id UUID REFERENCES incident_categories(id) ON DELETE SET NULL,
    priority incident_priority NOT NULL DEFAULT 'p3',
    status incident_status NOT NULL DEFAULT 'new',
    assignee_id UUID REFERENCES users(id) ON DELETE SET NULL,
    reporter_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    service_affected VARCHAR(255),
    root_cause TEXT,
    resolution_notes TEXT,
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    detected_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    updated_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT
);

-- Create incident_risks junction table
CREATE TABLE incident_risks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    risk_id UUID NOT NULL REFERENCES risks(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(incident_id, risk_id)
);

-- Create indexes
CREATE INDEX idx_incidents_status ON incidents(status);
CREATE INDEX idx_incidents_priority ON incidents(priority);
CREATE INDEX idx_incidents_category ON incidents(category_id);
CREATE INDEX idx_incidents_assignee ON incidents(assignee_id);
CREATE INDEX idx_incidents_reporter ON incidents(reporter_id);
CREATE INDEX idx_incidents_occurred_at ON incidents(occurred_at);
CREATE INDEX idx_incident_risks_incident ON incident_risks(incident_id);
CREATE INDEX idx_incident_risks_risk ON incident_risks(risk_id);
```

**Step 2: Create down migration**

```sql
-- backend/internal/migrations/migrations/012_incidents.down.sql
DROP TABLE IF EXISTS incident_risks;
DROP TABLE IF EXISTS incidents;
DROP TYPE IF EXISTS incident_status;
DROP TYPE IF EXISTS incident_priority;
```

**Step 3: Commit**

```bash
git add backend/internal/migrations/migrations/012_incidents.*
git commit -m "feat(db): add incidents and incident_risks tables"
```

---

## Phase 2: Backend Models

### Task 4: Create incident models

**Files:**
- Create: `backend/internal/models/incident.go`

**Step 1: Create incident model file**

```go
// backend/internal/models/incident.go
package models

import "time"

type IncidentStatus string

const (
	IncidentStatusNew           IncidentStatus = "new"
	IncidentStatusAcknowledged  IncidentStatus = "acknowledged"
	IncidentStatusInProgress    IncidentStatus = "in_progress"
	IncidentStatusOnHold        IncidentStatus = "on_hold"
	IncidentStatusResolved      IncidentStatus = "resolved"
	IncidentStatusClosed        IncidentStatus = "closed"
)

type IncidentPriority string

const (
	PriorityP1 IncidentPriority = "p1"
	PriorityP2 IncidentPriority = "p2"
	PriorityP3 IncidentPriority = "p3"
	PriorityP4 IncidentPriority = "p4"
)

type IncidentCategory struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateIncidentCategoryInput struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type UpdateIncidentCategoryInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type Incident struct {
	ID               string           `json:"id" db:"id"`
	Title            string           `json:"title" db:"title"`
	Description      string           `json:"description,omitempty" db:"description"`
	CategoryID       *string          `json:"category_id,omitempty" db:"category_id"`
	Category         *IncidentCategory `json:"category,omitempty" db:"-"`
	Priority         IncidentPriority `json:"priority" db:"priority"`
	Status           IncidentStatus   `json:"status" db:"status"`
	AssigneeID       *string          `json:"assignee_id,omitempty" db:"assignee_id"`
	Assignee         *User            `json:"assignee,omitempty" db:"-"`
	ReporterID       string           `json:"reporter_id" db:"reporter_id"`
	Reporter         *User            `json:"reporter,omitempty" db:"-"`
	ServiceAffected  string           `json:"service_affected,omitempty" db:"service_affected"`
	RootCause        string           `json:"root_cause,omitempty" db:"root_cause"`
	ResolutionNotes  string           `json:"resolution_notes,omitempty" db:"resolution_notes"`
	OccurredAt       time.Time        `json:"occurred_at" db:"occurred_at"`
	DetectedAt       time.Time        `json:"detected_at" db:"detected_at"`
	ResolvedAt       *time.Time       `json:"resolved_at,omitempty" db:"resolved_at"`
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at" db:"updated_at"`
	CreatedBy        string           `json:"created_by" db:"created_by"`
	UpdatedBy        string           `json:"updated_by" db:"updated_by"`
}

type CreateIncidentInput struct {
	Title           string           `json:"title" validate:"required,min=1,max=255"`
	Description     string           `json:"description"`
	CategoryID      *string          `json:"category_id"`
	Priority        IncidentPriority `json:"priority"`
	Status          IncidentStatus   `json:"status"`
	AssigneeID      *string          `json:"assignee_id" validate:"omitempty,uuid"`
	ServiceAffected string           `json:"service_affected"`
	OccurredAt      *string          `json:"occurred_at"`
	DetectedAt      *string          `json:"detected_at"`
}

type UpdateIncidentInput struct {
	Title           *string           `json:"title" validate:"omitempty,min=1,max=255"`
	Description     *string           `json:"description"`
	CategoryID      *string           `json:"category_id"`
	Priority        *IncidentPriority `json:"priority"`
	Status          *IncidentStatus   `json:"status"`
	AssigneeID      *string           `json:"assignee_id" validate:"omitempty,uuid"`
	ServiceAffected *string           `json:"service_affected"`
	RootCause       *string           `json:"root_cause"`
	ResolutionNotes *string           `json:"resolution_notes"`
	ResolvedAt      *string           `json:"resolved_at"`
}

type IncidentListParams struct {
	Status     *IncidentStatus
	Priority   *IncidentPriority
	CategoryID *string
	AssigneeID *string
	Search     string
	Sort       string
	Order      string
	Page       int
	Limit      int
}

type IncidentListResponse struct {
	Data  []*Incident `json:"data"`
	Meta  Meta        `json:"meta"`
}

type IncidentRisk struct {
	ID         string    `json:"id" db:"id"`
	IncidentID string    `json:"incident_id" db:"incident_id"`
	RiskID     string    `json:"risk_id" db:"risk_id"`
	Risk       *Risk     `json:"risk,omitempty" db:"-"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	CreatedBy  string    `json:"created_by" db:"created_by"`
}

type LinkIncidentRiskInput struct {
	RiskID string `json:"risk_id" validate:"required,uuid"`
}
```

**Step 2: Commit**

```bash
git add backend/internal/models/incident.go
git commit -m "feat(models): add incident-related models and types"
```

---

### Task 5: Add responder role constant to user model

**Files:**
- Modify: `backend/internal/models/user.go`

**Step 1: Add responder role**

Add to `backend/internal/models/user.go` after the existing role constants:

```go
const (
	RoleAdmin     UserRole = "admin"
	RoleResponder UserRole = "responder"
	RoleMember    UserRole = "member"
)
```

**Step 2: Commit**

```bash
git add backend/internal/models/user.go
git commit -m "feat(models): add responder role constant"
```

---

## Phase 3: Backend Database Layer

### Task 6: Create incident categories repository

**Files:**
- Create: `backend/internal/database/incident_categories.go`

**Step 1: Create repository file**

```go
// backend/internal/database/incident_categories.go
package database

import (
	"context"
	"database/sql"
	"errors"

	"backend/internal/models"
)

var ErrIncidentCategoryNotFound = errors.New("incident category not found")

type IncidentCategoryRepository interface {
	List(ctx context.Context) ([]*models.IncidentCategory, error)
	FindByID(ctx context.Context, id string) (*models.IncidentCategory, error)
	Create(ctx context.Context, input *models.CreateIncidentCategoryInput) (*models.IncidentCategory, error)
	Update(ctx context.Context, id string, input *models.UpdateIncidentCategoryInput) (*models.IncidentCategory, error)
	Delete(ctx context.Context, id string) error
}

type incidentCategoryRepository struct {
	db *sql.DB
}

func NewIncidentCategoryRepository(db *sql.DB) IncidentCategoryRepository {
	return &incidentCategoryRepository{db: db}
}

func (r *incidentCategoryRepository) List(ctx context.Context) ([]*models.IncidentCategory, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM incident_categories ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.IncidentCategory
	for rows.Next() {
		c := &models.IncidentCategory{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *incidentCategoryRepository) FindByID(ctx context.Context, id string) (*models.IncidentCategory, error) {
	c := &models.IncidentCategory{}
	query := `SELECT id, name, description, created_at, updated_at FROM incident_categories WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return c, nil
}

func (r *incidentCategoryRepository) Create(ctx context.Context, input *models.CreateIncidentCategoryInput) (*models.IncidentCategory, error) {
	c := &models.IncidentCategory{}
	query := `
		INSERT INTO incident_categories (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, input.Name, input.Description).Scan(
		&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *incidentCategoryRepository) Update(ctx context.Context, id string, input *models.UpdateIncidentCategoryInput) (*models.IncidentCategory, error) {
	if input.Name == nil && input.Description == nil {
		return nil, errors.New("at least one field must be updated")
	}

	c := &models.IncidentCategory{}
	query := `
		UPDATE incident_categories
		SET name = COALESCE($1, name), description = COALESCE($2, description), updated_at = NOW()
		WHERE id = $3
		RETURNING id, name, description, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, input.Name, input.Description, id).Scan(
		&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrIncidentCategoryNotFound
		}
		return nil, err
	}
	return c, nil
}

func (r *incidentCategoryRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM incident_categories WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrIncidentCategoryNotFound
	}
	return nil
}
```

**Step 2: Commit**

```bash
git add backend/internal/database/incident_categories.go
git commit -m "feat(db): add incident categories repository"
```

---

### Task 7: Create incidents repository

**Files:**
- Create: `backend/internal/database/incidents.go`

**Step 1: Create repository file**

```go
// backend/internal/database/incidents.go
package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"backend/internal/models"

	"github.com/google/uuid"
)

var ErrIncidentNotFound = errors.New("incident not found")

type IncidentRepository interface {
	Create(ctx context.Context, incident *models.Incident) error
	FindByID(ctx context.Context, id string) (*models.Incident, error)
	List(ctx context.Context, params *models.IncidentListParams) (*models.IncidentListResponse, error)
	Update(ctx context.Context, incident *models.Incident) error
	Delete(ctx context.Context, id string) error
}

type incidentRepository struct {
	db *sql.DB
}

func NewIncidentRepository(db *sql.DB) IncidentRepository {
	return &incidentRepository{db: db}
}

func (r *incidentRepository) Create(ctx context.Context, incident *models.Incident) error {
	if incident.ID == "" {
		incident.ID = uuid.New().String()
	}
	now := time.Now()
	incident.CreatedAt = now
	incident.UpdatedAt = now

	query := `
		INSERT INTO incidents (id, title, description, category_id, priority, status,
			assignee_id, reporter_id, service_affected, root_cause, resolution_notes,
			occurred_at, detected_at, resolved_at, created_at, updated_at, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		incident.ID, incident.Title, incident.Description, incident.CategoryID, incident.Priority,
		incident.Status, incident.AssigneeID, incident.ReporterID, incident.ServiceAffected,
		incident.RootCause, incident.ResolutionNotes, incident.OccurredAt, incident.DetectedAt,
		incident.ResolvedAt, incident.CreatedAt, incident.UpdatedAt, incident.CreatedBy, incident.UpdatedBy,
	).Scan(&incident.ID, &incident.CreatedAt, &incident.UpdatedAt)
}

func (r *incidentRepository) FindByID(ctx context.Context, id string) (*models.Incident, error) {
	query := `
		SELECT i.id, i.title, i.description, i.category_id, i.priority, i.status,
			i.assignee_id, i.reporter_id, i.service_affected, i.root_cause, i.resolution_notes,
			i.occurred_at, i.detected_at, i.resolved_at, i.created_at, i.updated_at, i.created_by, i.updated_by,
			c.id, c.name, c.description,
			assignee.id, assignee.name, assignee.email,
			reporter.id, reporter.name, reporter.email
		FROM incidents i
		LEFT JOIN incident_categories c ON i.category_id = c.id
		LEFT JOIN users assignee ON i.assignee_id = assignee.id
		LEFT JOIN users reporter ON i.reporter_id = reporter.id
		WHERE i.id = $1
	`
	incident := &models.Incident{}
	var catID, catName, catDesc sql.NullString
	var assigneeID, assigneeName, assigneeEmail sql.NullString
	var reporterID, reporterName, reporterEmail sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&incident.ID, &incident.Title, &incident.Description, &incident.CategoryID,
		&incident.Priority, &incident.Status, &incident.AssigneeID, &incident.ReporterID,
		&incident.ServiceAffected, &incident.RootCause, &incident.ResolutionNotes,
		&incident.OccurredAt, &incident.DetectedAt, &incident.ResolvedAt,
		&incident.CreatedAt, &incident.UpdatedAt, &incident.CreatedBy, &incident.UpdatedBy,
		&catID, &catName, &catDesc,
		&assigneeID, &assigneeName, &assigneeEmail,
		&reporterID, &reporterName, &reporterEmail,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrIncidentNotFound
		}
		return nil, err
	}

	if catID.Valid {
		incident.Category = &models.IncidentCategory{
			ID:          catID.String,
			Name:        catName.String,
			Description: catDesc.String,
		}
	}
	if assigneeID.Valid {
		incident.Assignee = &models.User{
			ID:    assigneeID.String,
			Name:  assigneeName.String,
			Email: assigneeEmail.String,
		}
	}
	if reporterID.Valid {
		incident.Reporter = &models.User{
			ID:    reporterID.String,
			Name:  reporterName.String,
			Email: reporterEmail.String,
		}
	}

	return incident, nil
}

func (r *incidentRepository) List(ctx context.Context, params *models.IncidentListParams) (*models.IncidentListResponse, error) {
	if params == nil {
		params = &models.IncidentListParams{Page: 1, Limit: 20}
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}

	where := "WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if params.Status != nil {
		where += fmt.Sprintf(" AND i.status = $%d", argNum)
		args = append(args, *params.Status)
		argNum++
	}
	if params.Priority != nil {
		where += fmt.Sprintf(" AND i.priority = $%d", argNum)
		args = append(args, *params.Priority)
		argNum++
	}
	if params.CategoryID != nil {
		where += fmt.Sprintf(" AND i.category_id = $%d", argNum)
		args = append(args, *params.CategoryID)
		argNum++
	}
	if params.AssigneeID != nil {
		where += fmt.Sprintf(" AND i.assignee_id = $%d", argNum)
		args = append(args, *params.AssigneeID)
		argNum++
	}
	if params.Search != "" {
		where += fmt.Sprintf(" AND (i.title ILIKE $%d OR i.description ILIKE $%d OR i.service_affected ILIKE $%d)", argNum, argNum, argNum)
		args = append(args, "%"+params.Search+"%")
		argNum++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM incidents i %s", where)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	orderBy := "i.created_at"
	if params.Sort != "" {
		switch params.Sort {
		case "title", "status", "priority", "occurred_at", "updated_at":
			orderBy = "i." + params.Sort
		case "created_at":
			orderBy = "i.created_at"
		}
	}
	orderDir := "DESC"
	if params.Order == "asc" {
		orderDir = "ASC"
	}

	offset := (params.Page - 1) * params.Limit
	query := fmt.Sprintf(`
		SELECT i.id, i.title, i.description, i.category_id, i.priority, i.status,
			i.assignee_id, i.reporter_id, i.service_affected, i.root_cause, i.resolution_notes,
			i.occurred_at, i.detected_at, i.resolved_at, i.created_at, i.updated_at, i.created_by, i.updated_by,
			c.id, c.name, c.description,
			assignee.id, assignee.name, assignee.email,
			reporter.id, reporter.name, reporter.email
		FROM incidents i
		LEFT JOIN incident_categories c ON i.category_id = c.id
		LEFT JOIN users assignee ON i.assignee_id = assignee.id
		LEFT JOIN users reporter ON i.reporter_id = reporter.id
		%s ORDER BY %s %s LIMIT $%d OFFSET $%d
	`, where, orderBy, orderDir, argNum, argNum+1)
	args = append(args, params.Limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []*models.Incident
	for rows.Next() {
		incident := &models.Incident{}
		var catID, catName, catDesc sql.NullString
		var assigneeID, assigneeName, assigneeEmail sql.NullString
		var reporterID, reporterName, reporterEmail sql.NullString
		err := rows.Scan(
			&incident.ID, &incident.Title, &incident.Description, &incident.CategoryID,
			&incident.Priority, &incident.Status, &incident.AssigneeID, &incident.ReporterID,
			&incident.ServiceAffected, &incident.RootCause, &incident.ResolutionNotes,
			&incident.OccurredAt, &incident.DetectedAt, &incident.ResolvedAt,
			&incident.CreatedAt, &incident.UpdatedAt, &incident.CreatedBy, &incident.UpdatedBy,
			&catID, &catName, &catDesc,
			&assigneeID, &assigneeName, &assigneeEmail,
			&reporterID, &reporterName, &reporterEmail,
		)
		if err != nil {
			return nil, err
		}
		if catID.Valid {
			incident.Category = &models.IncidentCategory{
				ID:          catID.String,
				Name:        catName.String,
				Description: catDesc.String,
			}
		}
		if assigneeID.Valid {
			incident.Assignee = &models.User{
				ID:    assigneeID.String,
				Name:  assigneeName.String,
				Email: assigneeEmail.String,
			}
		}
		if reporterID.Valid {
			incident.Reporter = &models.User{
				ID:    reporterID.String,
				Name:  reporterName.String,
				Email: reporterEmail.String,
			}
		}
		incidents = append(incidents, incident)
	}

	return &models.IncidentListResponse{
		Data: incidents,
		Meta: models.Meta{Page: params.Page, Limit: params.Limit, Total: total},
	}, rows.Err()
}

func (r *incidentRepository) Update(ctx context.Context, incident *models.Incident) error {
	now := time.Now()
	incident.UpdatedAt = now

	query := `
		UPDATE incidents SET title = $1, description = $2, category_id = $3, priority = $4,
			status = $5, assignee_id = $6, service_affected = $7, root_cause = $8,
			resolution_notes = $9, resolved_at = $10, updated_at = $11, updated_by = $12
		WHERE id = $13
		RETURNING updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		incident.Title, incident.Description, incident.CategoryID, incident.Priority,
		incident.Status, incident.AssigneeID, incident.ServiceAffected, incident.RootCause,
		incident.ResolutionNotes, incident.ResolvedAt, incident.UpdatedAt, incident.UpdatedBy, incident.ID,
	).Scan(&incident.UpdatedAt)
}

func (r *incidentRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM incidents WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrIncidentNotFound
	}
	return nil
}
```

**Step 2: Commit**

```bash
git add backend/internal/database/incidents.go
git commit -m "feat(db): add incidents repository"
```

---

### Task 8: Create incident_risks repository

**Files:**
- Create: `backend/internal/database/incident_risks.go`

**Step 1: Create repository file**

```go
// backend/internal/database/incident_risks.go
package database

import (
	"context"
	"database/sql"
	"errors"

	"backend/internal/models"

	"github.com/google/uuid"
)

var ErrIncidentRiskNotFound = errors.New("incident risk link not found")

type IncidentRiskRepository interface {
	ListByIncident(ctx context.Context, incidentID string) ([]*models.IncidentRisk, error)
	ListByRisk(ctx context.Context, riskID string) ([]*models.IncidentRisk, error)
	Create(ctx context.Context, incidentRisk *models.IncidentRisk) error
	Delete(ctx context.Context, incidentID, riskID string) error
}

type incidentRiskRepository struct {
	db *sql.DB
}

func NewIncidentRiskRepository(db *sql.DB) IncidentRiskRepository {
	return &incidentRiskRepository{db: db}
}

func (r *incidentRiskRepository) ListByIncident(ctx context.Context, incidentID string) ([]*models.IncidentRisk, error) {
	query := `
		SELECT ir.id, ir.incident_id, ir.risk_id, ir.created_at, ir.created_by,
			r.id, r.title, r.status, r.severity
		FROM incident_risks ir
		JOIN risks r ON ir.risk_id = r.id
		WHERE ir.incident_id = $1
		ORDER BY ir.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, incidentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*models.IncidentRisk
	for rows.Next() {
		link := &models.IncidentRisk{}
		link.Risk = &models.Risk{}
		err := rows.Scan(
			&link.ID, &link.IncidentID, &link.RiskID, &link.CreatedAt, &link.CreatedBy,
			&link.Risk.ID, &link.Risk.Title, &link.Risk.Status, &link.Risk.Severity,
		)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, rows.Err()
}

func (r *incidentRiskRepository) ListByRisk(ctx context.Context, riskID string) ([]*models.IncidentRisk, error) {
	query := `
		SELECT ir.id, ir.incident_id, ir.risk_id, ir.created_at, ir.created_by,
			i.id, i.title, i.status, i.priority, i.occurred_at
		FROM incident_risks ir
		JOIN incidents i ON ir.incident_id = i.id
		WHERE ir.risk_id = $1
		ORDER BY i.occurred_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, riskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*models.IncidentRisk
	for rows.Next() {
		link := &models.IncidentRisk{}
		link.Risk = &models.Risk{}
		err := rows.Scan(
			&link.ID, &link.IncidentID, &link.RiskID, &link.CreatedAt, &link.CreatedBy,
			&link.ID, &link.Title, &link.Status, &link.Priority, &link.OccurredAt,
		)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, rows.Err()
}

func (r *incidentRiskRepository) Create(ctx context.Context, incidentRisk *models.IncidentRisk) error {
	if incidentRisk.ID == "" {
		incidentRisk.ID = uuid.New().String()
	}

	query := `
		INSERT INTO incident_risks (id, incident_id, risk_id, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	return r.db.QueryRowContext(ctx, query,
		incidentRisk.ID, incidentRisk.IncidentID, incidentRisk.RiskID, incidentRisk.CreatedBy,
	).Scan(&incidentRisk.ID, &incidentRisk.CreatedAt)
}

func (r *incidentRiskRepository) Delete(ctx context.Context, incidentID, riskID string) error {
	result, err := r.db.ExecContext(ctx,
		"DELETE FROM incident_risks WHERE incident_id = $1 AND risk_id = $2",
		incidentID, riskID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrIncidentRiskNotFound
	}
	return nil
}
```

**Step 2: Commit**

```bash
git add backend/internal/database/incident_risks.go
git commit -m "feat(db): add incident_risks repository for risk linkage"
```

---

## Phase 4: Backend Handlers

### Task 9: Create RequireResponder middleware

**Files:**
- Modify: `backend/internal/middleware/auth.go`

**Step 1: Add RequireResponder function**

Add to `backend/internal/middleware/auth.go`:

```go
func RequireResponder(c *fiber.Ctx) error {
	user := GetUserFromContext(c)
	if user == nil {
		return c.Status(403).JSON(fiber.Map{
			"error": "access required",
		})
	}
	if user.Role != "admin" && user.Role != "responder" {
		return c.Status(403).JSON(fiber.Map{
			"error": "responder or admin access required",
		})
	}
	return c.Next()
}
```

**Step 2: Commit**

```bash
git add backend/internal/middleware/auth.go
git commit -m "feat(middleware): add RequireResponder middleware"
```

---

### Task 10: Create incident category handler

**Files:**
- Create: `backend/internal/handlers/incident_categories.go`

**Step 1: Create handler file**

```go
// backend/internal/handlers/incident_categories.go
package handlers

import (
	"errors"

	"backend/internal/database"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type IncidentCategoryHandler struct {
	categories database.IncidentCategoryRepository
}

func NewIncidentCategoryHandler(categories database.IncidentCategoryRepository) *IncidentCategoryHandler {
	return &IncidentCategoryHandler{categories: categories}
}

func (h *IncidentCategoryHandler) List(c *fiber.Ctx) error {
	categories, err := h.categories.List(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch incident categories"})
	}
	return c.JSON(categories)
}

func (h *IncidentCategoryHandler) Create(c *fiber.Ctx) error {
	var input models.CreateIncidentCategoryInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if input.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name is required"})
	}

	category, err := h.categories.Create(c.Context(), &input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create incident category"})
	}
	return c.Status(201).JSON(category)
}

func (h *IncidentCategoryHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "category id required"})
	}

	var input models.UpdateIncidentCategoryInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if input.Name == nil && input.Description == nil {
		return c.Status(400).JSON(fiber.Map{"error": "at least one field must be provided"})
	}
	if input.Name != nil && *input.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name cannot be empty"})
	}

	category, err := h.categories.Update(c.Context(), id, &input)
	if err != nil {
		if errors.Is(err, database.ErrIncidentCategoryNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "incident category not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to update incident category"})
	}
	return c.JSON(category)
}

func (h *IncidentCategoryHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "category id required"})
	}

	if err := h.categories.Delete(c.Context(), id); err != nil {
		if errors.Is(err, database.ErrIncidentCategoryNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "incident category not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete incident category"})
	}
	return c.SendStatus(204)
}
```

**Step 2: Commit**

```bash
git add backend/internal/handlers/incident_categories.go
git commit -m "feat(handlers): add incident category handler"
```

---

### Task 11: Create incident handler

**Files:**
- Create: `backend/internal/handlers/incidents.go`

**Step 1: Create handler file**

```go
// backend/internal/handlers/incidents.go
package handlers

import (
	"backend/internal/database"
	"backend/internal/middleware"
	"backend/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

type IncidentHandler struct {
	incidents  database.IncidentRepository
	categories database.IncidentCategoryRepository
	audit      database.AuditLogRepository
}

func NewIncidentHandler(incidents database.IncidentRepository, categories database.IncidentCategoryRepository, audit database.AuditLogRepository) *IncidentHandler {
	return &IncidentHandler{incidents: incidents, categories: categories, audit: audit}
}

func (h *IncidentHandler) List(c *fiber.Ctx) error {
	params := &models.IncidentListParams{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 20),
		Search: c.Query("search"),
		Sort:   c.Query("sort", "created_at"),
		Order:  c.Query("order", "desc"),
	}

	if status := c.Query("status"); status != "" {
		s := models.IncidentStatus(status)
		params.Status = &s
	}
	if priority := c.Query("priority"); priority != "" {
		p := models.IncidentPriority(priority)
		params.Priority = &p
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		params.CategoryID = &categoryID
	}
	if assigneeID := c.Query("assignee_id"); assigneeID != "" {
		params.AssigneeID = &assigneeID
	}

	response, err := h.incidents.List(c.Context(), params)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch incidents"})
	}
	return c.JSON(response)
}

func (h *IncidentHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	incident, err := h.incidents.FindByID(c.Context(), id)
	if err != nil {
		if err == database.ErrIncidentNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "incident not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch incident"})
	}
	return c.JSON(incident)
}

func (h *IncidentHandler) Create(c *fiber.Ctx) error {
	var input models.CreateIncidentInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if input.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "title is required"})
	}

	user := middleware.GetUserFromContext(c)

	// Set defaults
	if input.Status == "" {
		input.Status = models.IncidentStatusNew
	}
	if input.Priority == "" {
		input.Priority = models.PriorityP3
	}

	now := time.Now()
	incident := &models.Incident{
		Title:           input.Title,
		Description:     input.Description,
		CategoryID:      input.CategoryID,
		Priority:        input.Priority,
		Status:          input.Status,
		AssigneeID:      input.AssigneeID,
		ReporterID:      user.UserID,
		ServiceAffected: input.ServiceAffected,
		OccurredAt:      now,
		DetectedAt:      now,
		CreatedBy:       user.UserID,
		UpdatedBy:       user.UserID,
	}

	if input.OccurredAt != nil {
		t, err := time.Parse(time.RFC3339, *input.OccurredAt)
		if err == nil {
			incident.OccurredAt = t
		}
	}
	if input.DetectedAt != nil {
		t, err := time.Parse(time.RFC3339, *input.DetectedAt)
		if err == nil {
			incident.DetectedAt = t
		}
	}

	if err := h.incidents.Create(c.Context(), incident); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create incident"})
	}

	// Log audit event
	changes := map[string]any{
		"title":    incident.Title,
		"priority": incident.Priority,
		"status":   incident.Status,
	}
	h.audit.Create(c.Context(), "incident", incident.ID, models.AuditActionCreated, changes, user.UserID)

	return c.Status(201).JSON(incident)
}

func (h *IncidentHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var input models.UpdateIncidentInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	incident, err := h.incidents.FindByID(c.Context(), id)
	if err != nil {
		if err == database.ErrIncidentNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "incident not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch incident"})
	}

	user := middleware.GetUserFromContext(c)
	incident.UpdatedBy = user.UserID

	changes := make(map[string]any)

	if input.Title != nil {
		changes["title"] = map[string]any{"from": incident.Title, "to": *input.Title}
		incident.Title = *input.Title
	}
	if input.Description != nil {
		changes["description"] = map[string]any{"from": incident.Description, "to": *input.Description}
		incident.Description = *input.Description
	}
	if input.Status != nil {
		changes["status"] = map[string]any{"from": incident.Status, "to": *input.Status}
		incident.Status = *input.Status
	}
	if input.Priority != nil {
		changes["priority"] = map[string]any{"from": incident.Priority, "to": *input.Priority}
		incident.Priority = *input.Priority
	}
	if input.AssigneeID != nil {
		changes["assignee_id"] = map[string]any{"from": incident.AssigneeID, "to": *input.AssigneeID}
		incident.AssigneeID = input.AssigneeID
	}
	if input.CategoryID != nil {
		incident.CategoryID = input.CategoryID
	}
	if input.ServiceAffected != nil {
		incident.ServiceAffected = *input.ServiceAffected
	}
	if input.RootCause != nil {
		incident.RootCause = *input.RootCause
	}
	if input.ResolutionNotes != nil {
		incident.ResolutionNotes = *input.ResolutionNotes
	}
	if input.ResolvedAt != nil {
		t, err := time.Parse(time.RFC3339, *input.ResolvedAt)
		if err == nil {
			incident.ResolvedAt = &t
		}
	}

	// Auto-set resolved_at when status changes to resolved
	if input.Status != nil && *input.Status == models.IncidentStatusResolved && incident.ResolvedAt == nil {
		now := time.Now()
		incident.ResolvedAt = &now
	}

	if err := h.incidents.Update(c.Context(), incident); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update incident"})
	}

	if len(changes) > 0 {
		h.audit.Create(c.Context(), "incident", incident.ID, models.AuditActionUpdated, changes, user.UserID)
	}

	return c.JSON(incident)
}

func (h *IncidentHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	user := middleware.GetUserFromContext(c)

	h.audit.Create(c.Context(), "incident", id, models.AuditActionDeleted, nil, user.UserID)

	if err := h.incidents.Delete(c.Context(), id); err != nil {
		if err == database.ErrIncidentNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "incident not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete incident"})
	}

	return c.SendStatus(204)
}
```

**Step 2: Commit**

```bash
git add backend/internal/handlers/incidents.go
git commit -m "feat(handlers): add incident handler with CRUD operations"
```

---

### Task 12: Create incident-risk link handler

**Files:**
- Create: `backend/internal/handlers/incident_risks.go`

**Step 1: Create handler file**

```go
// backend/internal/handlers/incident_risks.go
package handlers

import (
	"backend/internal/database"
	"backend/internal/middleware"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type IncidentRiskHandler struct {
	incidentRisks database.IncidentRiskRepository
	risks         database.RiskRepository
}

func NewIncidentRiskHandler(incidentRisks database.IncidentRiskRepository, risks database.RiskRepository) *IncidentRiskHandler {
	return &IncidentRiskHandler{incidentRisks: incidentRisks, risks: risks}
}

func (h *IncidentRiskHandler) ListByIncident(c *fiber.Ctx) error {
	incidentID := c.Params("incidentId")
	links, err := h.incidentRisks.ListByIncident(c.Context(), incidentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch incident risks"})
	}
	return c.JSON(links)
}

func (h *IncidentRiskHandler) ListByRisk(c *fiber.Ctx) error {
	riskID := c.Params("riskId")
	links, err := h.incidentRisks.ListByRisk(c.Context(), riskID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch risk incidents"})
	}
	return c.JSON(links)
}

func (h *IncidentRiskHandler) Link(c *fiber.Ctx) error {
	incidentID := c.Params("incidentId")

	var input models.LinkIncidentRiskInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if input.RiskID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "risk_id is required"})
	}

	// Verify risk exists
	risk, err := h.risks.FindByID(c.Context(), input.RiskID)
	if err != nil {
		if err == database.ErrRiskNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "risk not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to verify risk"})
	}

	user := middleware.GetUserFromContext(c)

	link := &models.IncidentRisk{
		IncidentID: incidentID,
		RiskID:     input.RiskID,
		CreatedBy:  user.UserID,
	}

	if err := h.incidentRisks.Create(c.Context(), link); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to link risk to incident"})
	}

	link.Risk = risk
	return c.Status(201).JSON(link)
}

func (h *IncidentRiskHandler) Unlink(c *fiber.Ctx) error {
	incidentID := c.Params("incidentId")
	riskID := c.Params("riskId")

	if err := h.incidentRisks.Delete(c.Context(), incidentID, riskID); err != nil {
		if err == database.ErrIncidentRiskNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "incident risk link not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to unlink risk from incident"})
	}

	return c.SendStatus(204)
}
```

**Step 2: Commit**

```bash
git add backend/internal/handlers/incident_risks.go
git commit -m "feat(handlers): add incident-risk link handler"
```

---

### Task 13: Update server to include incident routes

**Files:**
- Modify: `backend/internal/server/server.go`
- Modify: `backend/internal/server/routes.go`

**Step 1: Update server.go to add incident repositories and handlers**

In `backend/internal/server/server.go`, add to the FiberServer struct:

```go
	incidentCategories      database.IncidentCategoryRepository
	incidents               database.IncidentRepository
	incidentRisks           database.IncidentRiskRepository
	incidentCategoryHandler *handlers.IncidentCategoryHandler
	incidentHandler         *handlers.IncidentHandler
	incidentRiskHandler     *handlers.IncidentRiskHandler
```

And add to the New() function after existing repositories:

```go
	incidentCategories := database.NewIncidentCategoryRepository(rawDB)
	incidents := database.NewIncidentRepository(rawDB)
	incidentRisks := database.NewIncidentRiskRepository(rawDB)
```

And add to the server initialization:

```go
		incidentCategories:      incidentCategories,
		incidents:               incidents,
		incidentRisks:           incidentRisks,
		incidentCategoryHandler: handlers.NewIncidentCategoryHandler(incidentCategories),
		incidentHandler:         handlers.NewIncidentHandler(incidents, incidentCategories, audit),
		incidentRiskHandler:     handlers.NewIncidentRiskHandler(incidentRisks, risks),
```

**Step 2: Update routes.go to add incident routes**

Add to `backend/internal/server/routes.go` before the AI routes section:

```go
	// Incident category routes (admin only)
	incidentCategories := protected.Group("/incident-categories")
	incidentCategories.Get("/", middleware.RequireAdmin, s.incidentCategoryHandler.List)
	incidentCategories.Post("/", middleware.RequireAdmin, s.incidentCategoryHandler.Create)
	incidentCategories.Put("/:id", middleware.RequireAdmin, s.incidentCategoryHandler.Update)
	incidentCategories.Delete("/:id", middleware.RequireAdmin, s.incidentCategoryHandler.Delete)

	// Incident routes
	incidents := protected.Group("/incidents")
	incidents.Get("/", s.incidentHandler.List)
	incidents.Post("/", middleware.RequireResponder, s.incidentHandler.Create)
	incidents.Get("/:id", s.incidentHandler.Get)
	incidents.Put("/:id", middleware.RequireResponder, s.incidentHandler.Update)
	incidents.Delete("/:id", middleware.RequireAdmin, s.incidentHandler.Delete)

	// Incident-risk linking (nested under incidents)
	incidents.Get("/:incidentId/risks", s.incidentRiskHandler.ListByIncident)
	incidents.Post("/:incidentId/risks", middleware.RequireResponder, s.incidentRiskHandler.Link)
	incidents.Delete("/:incidentId/risks/:riskId", middleware.RequireResponder, s.incidentRiskHandler.Unlink)

	// Reverse lookup: incidents for a risk
	risks.Get("/:riskId/incidents", s.incidentRiskHandler.ListByRisk)

	// Audit log for incidents
	incidents.Get("/:incidentId/audit", s.auditHandler.ListByIncident)
```

**Step 3: Update audit handler to support incidents**

Add to `backend/internal/handlers/audit.go`:

```go
func (h *AuditHandler) ListByIncident(c *fiber.Ctx) error {
	incidentID := c.Params("incidentId")
	logs, err := h.audit.ListByEntity(c.Context(), "incident", incidentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch audit logs"})
	}
	return c.JSON(logs)
}
```

**Step 4: Commit**

```bash
git add backend/internal/server/server.go backend/internal/server/routes.go backend/internal/handlers/audit.go
git commit -m "feat(server): register incident routes and handlers"
```

---

## Phase 5: Frontend Types

### Task 14: Create incident types

**Files:**
- Create: `frontend/apps/web/src/types/incident.ts`

**Step 1: Create types file**

```typescript
// frontend/apps/web/src/types/incident.ts
export type IncidentStatus = 'new' | 'acknowledged' | 'in_progress' | 'on_hold' | 'resolved' | 'closed';
export type IncidentPriority = 'p1' | 'p2' | 'p3' | 'p4';

export interface IncidentCategory {
  id: string;
  name: string;
  description?: string;
  created_at: string;
}

export interface Incident {
  id: string;
  title: string;
  description?: string;
  category_id?: string;
  category?: IncidentCategory;
  priority: IncidentPriority;
  status: IncidentStatus;
  assignee_id?: string;
  assignee?: {
    id: string;
    name: string;
    email: string;
  };
  reporter_id: string;
  reporter?: {
    id: string;
    name: string;
    email: string;
  };
  service_affected?: string;
  root_cause?: string;
  resolution_notes?: string;
  occurred_at: string;
  detected_at: string;
  resolved_at?: string;
  created_at: string;
  updated_at: string;
  created_by: string;
  updated_by: string;
}

export interface CreateIncidentInput {
  title: string;
  description?: string;
  category_id?: string;
  priority?: IncidentPriority;
  status?: IncidentStatus;
  assignee_id?: string;
  service_affected?: string;
  occurred_at?: string;
  detected_at?: string;
}

export interface UpdateIncidentInput {
  title?: string;
  description?: string;
  category_id?: string;
  priority?: IncidentPriority;
  status?: IncidentStatus;
  assignee_id?: string;
  service_affected?: string;
  root_cause?: string;
  resolution_notes?: string;
  resolved_at?: string;
}

export interface IncidentListParams {
  status?: IncidentStatus;
  priority?: IncidentPriority;
  category_id?: string;
  assignee_id?: string;
  search?: string;
  sort?: string;
  order?: 'asc' | 'desc';
  page?: number;
  limit?: number;
}

export interface IncidentListResponse {
  data: Incident[];
  meta: {
    page: number;
    limit: number;
    total: number;
  };
}

export interface IncidentRisk {
  id: string;
  incident_id: string;
  risk_id: string;
  risk?: {
    id: string;
    title: string;
    status: string;
    severity: string;
  };
  created_at: string;
  created_by: string;
}

export interface LinkIncidentRiskInput {
  risk_id: string;
}
```

**Step 2: Commit**

```bash
git add frontend/apps/web/src/types/incident.ts
git commit -m "feat(frontend): add incident TypeScript types"
```

---

## Phase 6: Frontend Hooks

### Task 15: Create incident hooks

**Files:**
- Create: `frontend/apps/web/src/hooks/useIncidents.ts`

**Step 1: Create hooks file**

```typescript
// frontend/apps/web/src/hooks/useIncidents.ts
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type {
  IncidentCategory,
  CreateIncidentInput,
  Incident,
  IncidentListParams,
  IncidentListResponse,
  IncidentStatus,
  IncidentRisk,
  LinkIncidentRiskInput,
  UpdateIncidentInput,
} from '@/types/incident';

const INCIDENTS_KEY = ['incidents'];
const INCIDENT_CATEGORIES_KEY = ['incident-categories'];
const INCIDENT_RISKS_KEY = ['incident-risks'];

// Incident Categories
export function useIncidentCategories() {
  return useQuery({
    queryKey: INCIDENT_CATEGORIES_KEY,
    queryFn: () => api.get<IncidentCategory[]>('/api/v1/incident-categories'),
    staleTime: 5 * 60 * 1000,
  });
}

export function useCreateIncidentCategory() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: { name: string; description?: string }) =>
      api.post<IncidentCategory>('/api/v1/incident-categories', input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: INCIDENT_CATEGORIES_KEY });
    },
  });
}

// Incidents List
export function useIncidents(params?: IncidentListParams) {
  const queryString = buildQueryString(params);
  return useQuery({
    queryKey: [...INCIDENTS_KEY, params],
    queryFn: () => api.get<IncidentListResponse>(`/api/v1/incidents${queryString}`),
  });
}

// Single Incident
export function useIncident(id: string) {
  return useQuery({
    queryKey: [...INCIDENTS_KEY, id],
    queryFn: () => api.get<Incident>(`/api/v1/incidents/${id}`),
    enabled: !!id,
  });
}

// Create Incident
export function useCreateIncident() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: CreateIncidentInput) =>
      api.post<Incident>('/api/v1/incidents', input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: INCIDENTS_KEY });
    },
  });
}

// Update Incident
export function useUpdateIncident(id: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: UpdateIncidentInput) =>
      api.put<Incident>(`/api/v1/incidents/${id}`, input),
    onSuccess: (updatedIncident) => {
      queryClient.setQueryData([...INCIDENTS_KEY, id], updatedIncident);
      queryClient.invalidateQueries({ queryKey: INCIDENTS_KEY });
    },
  });
}

// Delete Incident
export function useDeleteIncident() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => api.delete(`/api/v1/incidents/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: INCIDENTS_KEY });
    },
  });
}

// Incident-Risk Links
export function useIncidentRisks(incidentId: string) {
  return useQuery({
    queryKey: [...INCIDENT_RISKS_KEY, incidentId],
    queryFn: () => api.get<IncidentRisk[]>(`/api/v1/incidents/${incidentId}/risks`),
    enabled: !!incidentId,
  });
}

export function useLinkIncidentRisk(incidentId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: LinkIncidentRiskInput) =>
      api.post<IncidentRisk>(`/api/v1/incidents/${incidentId}/risks`, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [...INCIDENT_RISKS_KEY, incidentId] });
    },
  });
}

export function useUnlinkIncidentRisk(incidentId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (riskId: string) =>
      api.delete(`/api/v1/incidents/${incidentId}/risks/${riskId}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [...INCIDENT_RISKS_KEY, incidentId] });
    },
  });
}

// Helper to build query string
function buildQueryString(params?: IncidentListParams): string {
  if (!params) return '';

  const searchParams = new URLSearchParams();

  if (params.status) searchParams.set('status', params.status);
  if (params.priority) searchParams.set('priority', params.priority);
  if (params.category_id) searchParams.set('category_id', params.category_id);
  if (params.assignee_id) searchParams.set('assignee_id', params.assignee_id);
  if (params.search) searchParams.set('search', params.search);
  if (params.sort) searchParams.set('sort', params.sort);
  if (params.order) searchParams.set('order', params.order);
  if (params.page) searchParams.set('page', String(params.page));
  if (params.limit) searchParams.set('limit', String(params.limit));

  const queryString = searchParams.toString();
  return queryString ? `?${queryString}` : '';
}
```

**Step 2: Commit**

```bash
git add frontend/apps/web/src/hooks/useIncidents.ts
git commit -m "feat(frontend): add incident hooks for data fetching"
```

---

## Phase 7: Frontend Pages

### Task 16: Create incidents list page

**Files:**
- Create: `frontend/apps/web/src/routes/app/incidents/index.tsx`

**Step 1: Create page file**

Create the incidents list page following the pattern from `routes/app/risks/index.tsx` with:
- Filter bar for status, priority, category, search
- Data table with columns: Priority, Title, Status, Category, Assignee, Occurred, Age
- Pagination
- Link to create new incident

**Step 2: Commit**

```bash
git add frontend/apps/web/src/routes/app/incidents/index.tsx
git commit -m "feat(frontend): add incidents list page with filters"
```

---

### Task 17: Create incident detail page

**Files:**
- Create: `frontend/apps/web/src/routes/app/incidents/$id.tsx`

**Step 1: Create page file**

Create the incident detail page following the pattern from `routes/app/risks/$id.tsx` with:
- Header with priority badge, status badge
- Two-column layout: description/service/cause/resolution on left, metadata on right
- Audit history timeline
- Linked risks section with add/remove capability

**Step 2: Commit**

```bash
git add frontend/apps/web/src/routes/app/incidents/$id.tsx
git commit -m "feat(frontend): add incident detail page with editing"
```

---

### Task 18: Create new incident page

**Files:**
- Create: `frontend/apps/web/src/routes/app/incidents/new.tsx`

**Step 1: Create page file**

Create the new incident form page with:
- Title input (required)
- Description textarea
- Priority select (P1-P4)
- Category select
- Service affected input
- Occurred/Detected date pickers
- Form validation and submission

**Step 2: Commit**

```bash
git add frontend/apps/web/src/routes/app/incidents/new.tsx
git commit -m "feat(frontend): add new incident form page"
```

---

## Phase 8: Navigation Updates

### Task 19: Add incidents to navigation

**Files:**
- Modify: `frontend/apps/web/src/components/AppSidebar.tsx` (or equivalent navigation component)

**Step 1: Add Incidents link**

Add "Incidents" link to the main navigation sidebar, between "Risks" and "Categories".

**Step 2: Commit**

```bash
git add frontend/apps/web/src/components/AppSidebar.tsx
git commit -m "feat(frontend): add incidents to navigation"
```

---

## Phase 9: Testing & Verification

### Task 20: Run backend tests

**Step 1: Run tests**

```bash
cd backend && make test
```

Expected: All existing tests pass.

**Step 2: Run integration tests**

```bash
cd backend && make itest
```

Expected: All integration tests pass.

---

### Task 21: Run frontend type check

**Step 1: Run type check**

```bash
cd frontend && bun run check-types
```

Expected: No type errors.

---

### Task 22: Manual verification

**Step 1: Start services**

```bash
make docker-run && make dev
```

**Step 2: Verify in browser**

1. Navigate to http://localhost:3001/app/incidents
2. Create a new incident
3. View incident detail
4. Edit incident
5. Link a risk to the incident
6. Check audit history

---

## Summary

This plan creates:
- **3 database migrations**: role enum, incident_categories, incidents + incident_risks
- **3 model files**: incident.go (with all types)
- **3 repository files**: incident_categories.go, incidents.go, incident_risks.go
- **3 handler files**: incident_categories.go, incidents.go, incident_risks.go
- **1 middleware update**: RequireResponder
- **2 frontend type files**: incident.ts
- **1 frontend hooks file**: useIncidents.ts
- **3 frontend pages**: list, detail, new
- **1 navigation update**: sidebar

Total: ~20 commits following TDD principles where applicable.
