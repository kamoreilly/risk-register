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
var ErrIncidentRiskNotFound = errors.New("incident risk link not found")
var ErrIncidentRiskAlreadyExists = errors.New("incident is already linked to this risk")

type IncidentRepository interface {
	Create(ctx context.Context, incident *models.Incident) error
	FindByID(ctx context.Context, id string) (*models.Incident, error)
	List(ctx context.Context, params *models.IncidentListParams) (*models.IncidentListResponse, error)
	Update(ctx context.Context, incident *models.Incident) error
	Delete(ctx context.Context, id string) error
}

type IncidentRiskRepository interface {
	ListByIncident(ctx context.Context, incidentID string) ([]*models.IncidentRisk, error)
	LinkRisk(ctx context.Context, incidentID, riskID, createdBy string) (*models.IncidentRisk, error)
	UnlinkRisk(ctx context.Context, incidentID, riskID string) error
}

type incidentRepository struct {
	db *sql.DB
}

type incidentRiskRepository struct {
	db *sql.DB
}

func NewIncidentRepository(db *sql.DB) IncidentRepository {
	return &incidentRepository{db: db}
}

func NewIncidentRiskRepository(db *sql.DB) IncidentRiskRepository {
	return &incidentRiskRepository{db: db}
}

func (r *incidentRepository) Create(ctx context.Context, incident *models.Incident) error {
	if incident.ID == "" {
		incident.ID = uuid.New().String()
	}
	now := time.Now()
	incident.CreatedAt = now
	incident.UpdatedAt = now

	// Set default status and priority if not provided
	if incident.Status == "" {
		incident.Status = models.IncidentStatusNew
	}
	if incident.Priority == "" {
		incident.Priority = models.PriorityP3
	}

	// Set default timestamps
	if incident.OccurredAt.IsZero() {
		incident.OccurredAt = now
	}
	if incident.DetectedAt.IsZero() {
		incident.DetectedAt = now
	}

	query := `
		INSERT INTO incidents (id, title, description, category_id, priority, status, assignee_id, reporter_id,
			service_affected, root_cause, resolution_notes, occurred_at, detected_at, resolved_at,
			created_at, updated_at, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		incident.ID, incident.Title, incident.Description, incident.CategoryID, incident.Priority, incident.Status,
		incident.AssigneeID, incident.ReporterID, incident.ServiceAffected, incident.RootCause, incident.ResolutionNotes,
		incident.OccurredAt, incident.DetectedAt, incident.ResolvedAt,
		incident.CreatedAt, incident.UpdatedAt, incident.CreatedBy, incident.UpdatedBy,
	).Scan(&incident.ID, &incident.CreatedAt, &incident.UpdatedAt)
}

func (r *incidentRepository) FindByID(ctx context.Context, id string) (*models.Incident, error) {
	query := `
		SELECT i.id, i.title, i.description, i.category_id, i.priority, i.status, i.assignee_id, i.reporter_id,
			i.service_affected, i.root_cause, i.resolution_notes, i.occurred_at, i.detected_at, i.resolved_at,
			i.created_at, i.updated_at, i.created_by, i.updated_by,
			c.id, c.name, c.description
		FROM incidents i
		LEFT JOIN incident_categories c ON i.category_id = c.id
		WHERE i.id = $1
	`
	incident := &models.Incident{}
	var catID, catName, catDesc sql.NullString
	var assigneeID, resolvedAt sql.NullString
	var description, serviceAffected, rootCause, resolutionNotes sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&incident.ID, &incident.Title, &description, &incident.CategoryID, &incident.Priority, &incident.Status,
		&assigneeID, &incident.ReporterID, &serviceAffected, &rootCause, &resolutionNotes,
		&incident.OccurredAt, &incident.DetectedAt, &resolvedAt,
		&incident.CreatedAt, &incident.UpdatedAt, &incident.CreatedBy, &incident.UpdatedBy,
		&catID, &catName, &catDesc,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrIncidentNotFound
		}
		return nil, err
	}

	// Set nullable fields
	if description.Valid {
		incident.Description = description.String
	}
	if serviceAffected.Valid {
		incident.ServiceAffected = serviceAffected.String
	}
	if rootCause.Valid {
		incident.RootCause = rootCause.String
	}
	if resolutionNotes.Valid {
		incident.ResolutionNotes = resolutionNotes.String
	}
	if assigneeID.Valid {
		incident.AssigneeID = &assigneeID.String
	}
	if resolvedAt.Valid {
		t, err := time.Parse(time.RFC3339, resolvedAt.String)
		if err == nil {
			incident.ResolvedAt = &t
		}
	}

	if catID.Valid {
		incident.Category = &models.IncidentCategory{
			ID:          catID.String,
			Name:        catName.String,
			Description: catDesc.String,
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

	// Build WHERE clause
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

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM incidents i %s", where)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Build ORDER BY
	orderBy := "i.created_at"
	if params.Sort != "" {
		switch params.Sort {
		case "title", "status", "priority", "category_id", "occurred_at", "detected_at", "resolved_at", "updated_at":
			orderBy = "i." + params.Sort
		case "created_at":
			orderBy = "i.created_at"
		default:
			orderBy = "i.created_at"
		}
	}
	orderDir := "DESC"
	if params.Order == "asc" {
		orderDir = "ASC"
	}

	// Get paginated results
	offset := (params.Page - 1) * params.Limit
	query := fmt.Sprintf(`
		SELECT i.id, i.title, i.description, i.category_id, i.priority, i.status, i.assignee_id, i.reporter_id,
			i.service_affected, i.root_cause, i.resolution_notes, i.occurred_at, i.detected_at, i.resolved_at,
			i.created_at, i.updated_at, i.created_by, i.updated_by,
			c.id, c.name, c.description
		FROM incidents i
		LEFT JOIN incident_categories c ON i.category_id = c.id
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
		var assigneeID, resolvedAt sql.NullString
		var description, serviceAffected, rootCause, resolutionNotes sql.NullString

		err := rows.Scan(
			&incident.ID, &incident.Title, &description, &incident.CategoryID, &incident.Priority, &incident.Status,
			&assigneeID, &incident.ReporterID, &serviceAffected, &rootCause, &resolutionNotes,
			&incident.OccurredAt, &incident.DetectedAt, &resolvedAt,
			&incident.CreatedAt, &incident.UpdatedAt, &incident.CreatedBy, &incident.UpdatedBy,
			&catID, &catName, &catDesc,
		)
		if err != nil {
			return nil, err
		}

		// Set nullable fields
		if description.Valid {
			incident.Description = description.String
		}
		if serviceAffected.Valid {
			incident.ServiceAffected = serviceAffected.String
		}
		if rootCause.Valid {
			incident.RootCause = rootCause.String
		}
		if resolutionNotes.Valid {
			incident.ResolutionNotes = resolutionNotes.String
		}
		if assigneeID.Valid {
			incident.AssigneeID = &assigneeID.String
		}
		if resolvedAt.Valid {
			t, err := time.Parse(time.RFC3339, resolvedAt.String)
			if err == nil {
				incident.ResolvedAt = &t
			}
		}

		if catID.Valid {
			incident.Category = &models.IncidentCategory{
				ID:          catID.String,
				Name:        catName.String,
				Description: catDesc.String,
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
		UPDATE incidents SET title = $1, description = $2, category_id = $3, priority = $4, status = $5,
			assignee_id = $6, service_affected = $7, root_cause = $8, resolution_notes = $9,
			resolved_at = $10, updated_at = $11, updated_by = $12
		WHERE id = $13
		RETURNING updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		incident.Title, incident.Description, incident.CategoryID, incident.Priority, incident.Status,
		incident.AssigneeID, incident.ServiceAffected, incident.RootCause, incident.ResolutionNotes,
		incident.ResolvedAt, incident.UpdatedAt, incident.UpdatedBy, incident.ID,
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

// IncidentRiskRepository methods

func (r *incidentRiskRepository) ListByIncident(ctx context.Context, incidentID string) ([]*models.IncidentRisk, error) {
	query := `
		SELECT ir.id, ir.incident_id, ir.risk_id, ir.created_at, ir.created_by,
			r.id, r.title, r.description, r.status, r.severity
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
		link := &models.IncidentRisk{
			Risk: &models.Risk{},
		}
		var riskDesc sql.NullString

		err := rows.Scan(
			&link.ID, &link.IncidentID, &link.RiskID, &link.CreatedAt, &link.CreatedBy,
			&link.Risk.ID, &link.Risk.Title, &riskDesc, &link.Risk.Status, &link.Risk.Severity,
		)
		if err != nil {
			return nil, err
		}
		if riskDesc.Valid {
			link.Risk.Description = riskDesc.String
		}
		links = append(links, link)
	}

	return links, rows.Err()
}

func (r *incidentRiskRepository) LinkRisk(ctx context.Context, incidentID, riskID, createdBy string) (*models.IncidentRisk, error) {
	// Check if link already exists
	var exists bool
	err := r.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM incident_risks WHERE incident_id = $1 AND risk_id = $2)",
		incidentID, riskID,
	).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrIncidentRiskAlreadyExists
	}

	link := &models.IncidentRisk{
		ID:         uuid.New().String(),
		IncidentID: incidentID,
		RiskID:     riskID,
		CreatedBy:  createdBy,
	}

	query := `
		INSERT INTO incident_risks (id, incident_id, risk_id, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`
	err = r.db.QueryRowContext(ctx, query, link.ID, link.IncidentID, link.RiskID, link.CreatedBy).Scan(&link.CreatedAt)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (r *incidentRiskRepository) UnlinkRisk(ctx context.Context, incidentID, riskID string) error {
	result, err := r.db.ExecContext(ctx,
		"DELETE FROM incident_risks WHERE incident_id = $1 AND risk_id = $2",
		incidentID, riskID,
	)
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
