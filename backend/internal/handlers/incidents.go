package handlers

import (
	"errors"
	"time"

	"backend/internal/database"
	"backend/internal/middleware"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type IncidentHandler struct {
	incidents         database.IncidentRepository
	incidentCategories database.IncidentCategoryRepository
	incidentRisks     database.IncidentRiskRepository
	audit             database.AuditLogRepository
}

func NewIncidentHandler(
	incidents database.IncidentRepository,
	incidentCategories database.IncidentCategoryRepository,
	incidentRisks database.IncidentRiskRepository,
	audit database.AuditLogRepository,
) *IncidentHandler {
	return &IncidentHandler{
		incidents:          incidents,
		incidentCategories: incidentCategories,
		incidentRisks:      incidentRisks,
		audit:              audit,
	}
}

// normalizeCategoryID handles empty string CategoryID by setting it to nil
func (h *IncidentHandler) normalizeCategoryID(categoryID *string) *string {
	if categoryID != nil && *categoryID == "" {
		return nil
	}
	return categoryID
}

// validateCategoryInput validates that the category exists if provided
func (h *IncidentHandler) validateCategoryInput(c *fiber.Ctx, categoryID *string) error {
	if categoryID == nil {
		return nil
	}

	cat, err := h.incidentCategories.FindByID(c.Context(), *categoryID)
	if err != nil || cat == nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid category_id"})
	}
	return nil
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

	// Normalize category ID first
	categoryID := h.normalizeCategoryID(input.CategoryID)

	// Validate category input
	if err := h.validateCategoryInput(c, categoryID); err != nil {
		return err
	}

	// Parse timestamps
	var occurredAt, detectedAt time.Time
	if input.OccurredAt != nil {
		t, err := time.Parse(time.RFC3339, *input.OccurredAt)
		if err == nil {
			occurredAt = t
		} else {
			occurredAt = time.Now()
		}
	} else {
		occurredAt = time.Now()
	}

	if input.DetectedAt != nil {
		t, err := time.Parse(time.RFC3339, *input.DetectedAt)
		if err == nil {
			detectedAt = t
		} else {
			detectedAt = time.Now()
		}
	} else {
		detectedAt = time.Now()
	}

	incident := &models.Incident{
		Title:           input.Title,
		Description:     input.Description,
		CategoryID:      categoryID,
		Priority:        input.Priority,
		Status:          input.Status,
		AssigneeID:      input.AssigneeID,
		ReporterID:      user.UserID,
		ServiceAffected: input.ServiceAffected,
		OccurredAt:      occurredAt,
		DetectedAt:      detectedAt,
		CreatedBy:       user.UserID,
		UpdatedBy:       user.UserID,
	}

	if err := h.incidents.Create(c.Context(), incident); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create incident"})
	}

	// Log audit event
	changes := map[string]any{
		"title":           incident.Title,
		"description":     incident.Description,
		"priority":        incident.Priority,
		"status":          incident.Status,
		"service_affected": incident.ServiceAffected,
	}
	if incident.CategoryID != nil {
		changes["category_id"] = *incident.CategoryID
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

	// Get existing incident
	incident, err := h.incidents.FindByID(c.Context(), id)
	if err != nil {
		if err == database.ErrIncidentNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "incident not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch incident"})
	}

	user := middleware.GetUserFromContext(c)
	incident.UpdatedBy = user.UserID

	// Track changes for audit log
	changes := make(map[string]any)

	// Apply updates and track changes
	if input.Title != nil {
		changes["title"] = map[string]any{"from": incident.Title, "to": *input.Title}
		incident.Title = *input.Title
	}
	if input.Description != nil {
		changes["description"] = map[string]any{"from": incident.Description, "to": *input.Description}
		incident.Description = *input.Description
	}
	if input.Priority != nil {
		changes["priority"] = map[string]any{"from": incident.Priority, "to": *input.Priority}
		incident.Priority = *input.Priority
	}
	if input.Status != nil {
		changes["status"] = map[string]any{"from": incident.Status, "to": *input.Status}
		incident.Status = *input.Status

		// Auto-set resolved_at when status changes to resolved or closed
		if (*input.Status == models.IncidentStatusResolved || *input.Status == models.IncidentStatusClosed) && incident.ResolvedAt == nil {
			now := time.Now()
			incident.ResolvedAt = &now
		}
	}
	if input.CategoryID != nil {
		var oldCategoryID string
		if incident.CategoryID != nil {
			oldCategoryID = *incident.CategoryID
		}

		normalizedCategoryID := h.normalizeCategoryID(input.CategoryID)

		if normalizedCategoryID == nil {
			incident.CategoryID = nil
			incident.Category = nil
			changes["category_id"] = map[string]any{"from": oldCategoryID, "to": nil}
		} else {
			if err := h.validateCategoryInput(c, normalizedCategoryID); err != nil {
				return err
			}
			incident.CategoryID = normalizedCategoryID
			cat, _ := h.incidentCategories.FindByID(c.Context(), *normalizedCategoryID)
			incident.Category = cat
			changes["category_id"] = map[string]any{"from": oldCategoryID, "to": *normalizedCategoryID}
		}
	}
	if input.AssigneeID != nil {
		changes["assignee_id"] = map[string]any{"from": incident.AssigneeID, "to": input.AssigneeID}
		incident.AssigneeID = input.AssigneeID
	}
	if input.ServiceAffected != nil {
		changes["service_affected"] = map[string]any{"from": incident.ServiceAffected, "to": *input.ServiceAffected}
		incident.ServiceAffected = *input.ServiceAffected
	}
	if input.RootCause != nil {
		changes["root_cause"] = map[string]any{"from": incident.RootCause, "to": *input.RootCause}
		incident.RootCause = *input.RootCause
	}
	if input.ResolutionNotes != nil {
		changes["resolution_notes"] = map[string]any{"from": incident.ResolutionNotes, "to": *input.ResolutionNotes}
		incident.ResolutionNotes = *input.ResolutionNotes
	}
	if input.ResolvedAt != nil {
		t, err := time.Parse(time.RFC3339, *input.ResolvedAt)
		if err == nil {
			incident.ResolvedAt = &t
			changes["resolved_at"] = map[string]any{"to": *input.ResolvedAt}
		}
	}

	if err := h.incidents.Update(c.Context(), incident); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update incident"})
	}

	// Log audit event if there were changes
	if len(changes) > 0 {
		h.audit.Create(c.Context(), "incident", incident.ID, models.AuditActionUpdated, changes, user.UserID)
	}

	return c.JSON(incident)
}

func (h *IncidentHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	user := middleware.GetUserFromContext(c)

	// Log audit event before deletion
	h.audit.Create(c.Context(), "incident", id, models.AuditActionDeleted, nil, user.UserID)

	if err := h.incidents.Delete(c.Context(), id); err != nil {
		if err == database.ErrIncidentNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "incident not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete incident"})
	}

	return c.SendStatus(204)
}

// IncidentRiskHandler handles incident-risk link operations
type IncidentRiskHandler struct {
	incidentRisks database.IncidentRiskRepository
	audit         database.AuditLogRepository
}

func NewIncidentRiskHandler(incidentRisks database.IncidentRiskRepository, audit database.AuditLogRepository) *IncidentRiskHandler {
	return &IncidentRiskHandler{incidentRisks: incidentRisks, audit: audit}
}

func (h *IncidentRiskHandler) ListRisks(c *fiber.Ctx) error {
	incidentID := c.Params("incidentId")
	links, err := h.incidentRisks.ListByIncident(c.Context(), incidentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch linked risks"})
	}
	return c.JSON(links)
}

func (h *IncidentRiskHandler) LinkRisk(c *fiber.Ctx) error {
	incidentID := c.Params("incidentId")

	var input models.LinkIncidentRiskInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if input.RiskID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "risk_id is required"})
	}

	user := middleware.GetUserFromContext(c)

	link, err := h.incidentRisks.LinkRisk(c.Context(), incidentID, input.RiskID, user.UserID)
	if err != nil {
		if errors.Is(err, database.ErrIncidentRiskAlreadyExists) {
			return c.Status(409).JSON(fiber.Map{"error": "incident is already linked to this risk"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to link risk"})
	}

	// Log audit event
	h.audit.Create(c.Context(), "incident", incidentID, models.AuditActionUpdated, map[string]any{
		"action": "link_risk",
		"risk_id": input.RiskID,
	}, user.UserID)

	return c.Status(201).JSON(link)
}

func (h *IncidentRiskHandler) UnlinkRisk(c *fiber.Ctx) error {
	incidentID := c.Params("incidentId")
	riskID := c.Params("riskId")

	user := middleware.GetUserFromContext(c)

	if err := h.incidentRisks.UnlinkRisk(c.Context(), incidentID, riskID); err != nil {
		if errors.Is(err, database.ErrIncidentRiskNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "incident risk link not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to unlink risk"})
	}

	// Log audit event
	h.audit.Create(c.Context(), "incident", incidentID, models.AuditActionUpdated, map[string]any{
		"action": "unlink_risk",
		"risk_id": riskID,
	}, user.UserID)

	return c.SendStatus(204)
}
