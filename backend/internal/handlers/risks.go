package handlers

import (
	"backend/internal/database"
	"backend/internal/middleware"
	"backend/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RiskHandler struct {
	risks      database.RiskRepository
	categories database.CategoryRepository
	audit      database.AuditLogRepository
}

func NewRiskHandler(risks database.RiskRepository, categories database.CategoryRepository, audit database.AuditLogRepository) *RiskHandler {
	return &RiskHandler{risks: risks, categories: categories, audit: audit}
}

func (h *RiskHandler) List(c *fiber.Ctx) error {
	params := &models.RiskListParams{
		Page:  c.QueryInt("page", 1),
		Limit: c.QueryInt("limit", 20),
		Search: c.Query("search"),
		Sort:   c.Query("sort", "created_at"),
		Order:  c.Query("order", "desc"),
	}

	if status := c.Query("status"); status != "" {
		s := models.RiskStatus(status)
		params.Status = &s
	}
	if severity := c.Query("severity"); severity != "" {
		s := models.RiskSeverity(severity)
		params.Severity = &s
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		params.CategoryID = &categoryID
	}
	if ownerID := c.Query("owner_id"); ownerID != "" {
		params.OwnerID = &ownerID
	}

	response, err := h.risks.List(c.Context(), params)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch risks"})
	}
	return c.JSON(response)
}

func (h *RiskHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	risk, err := h.risks.FindByID(c.Context(), id)
	if err != nil {
		if err == database.ErrRiskNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "risk not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch risk"})
	}
	return c.JSON(risk)
}

func (h *RiskHandler) Create(c *fiber.Ctx) error {
	var input models.CreateRiskInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if input.Title == "" || input.OwnerID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "title and owner_id are required"})
	}

	user := middleware.GetUserFromContext(c)

	// Set defaults
	if input.Status == "" {
		input.Status = models.StatusOpen
	}
	if input.Severity == "" {
		input.Severity = models.SeverityMedium
	}

	risk := &models.Risk{
		Title:       input.Title,
		Description: input.Description,
		OwnerID:     input.OwnerID,
		Status:      input.Status,
		Severity:    input.Severity,
		CategoryID:  input.CategoryID,
		CreatedBy:   user.UserID,
		UpdatedBy:   user.UserID,
	}

	if input.ReviewDate != nil {
		t, err := time.Parse("2006-01-02", *input.ReviewDate)
		if err == nil {
			risk.ReviewDate = &t
		}
	}

	if err := h.risks.Create(c.Context(), risk); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create risk"})
	}

	// Log audit event
	changes := map[string]any{
		"title":       risk.Title,
		"description": risk.Description,
		"owner_id":    risk.OwnerID,
		"status":      risk.Status,
		"severity":    risk.Severity,
	}
	if risk.CategoryID != nil {
		changes["category_id"] = *risk.CategoryID
	}
	h.audit.Create(c.Context(), "risk", risk.ID, models.AuditActionCreated, changes, user.UserID)

	return c.Status(201).JSON(risk)
}

func (h *RiskHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var input models.UpdateRiskInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Get existing risk
	risk, err := h.risks.FindByID(c.Context(), id)
	if err != nil {
		if err == database.ErrRiskNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "risk not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch risk"})
	}

	user := middleware.GetUserFromContext(c)
	risk.UpdatedBy = user.UserID

	// Track changes for audit log
	changes := make(map[string]any)

	// Apply updates and track changes
	if input.Title != nil {
		changes["title"] = map[string]any{"from": risk.Title, "to": *input.Title}
		risk.Title = *input.Title
	}
	if input.Description != nil {
		changes["description"] = map[string]any{"from": risk.Description, "to": *input.Description}
		risk.Description = *input.Description
	}
	if input.OwnerID != nil {
		changes["owner_id"] = map[string]any{"from": risk.OwnerID, "to": *input.OwnerID}
		risk.OwnerID = *input.OwnerID
	}
	if input.Status != nil {
		changes["status"] = map[string]any{"from": risk.Status, "to": *input.Status}
		risk.Status = *input.Status
	}
	if input.Severity != nil {
		changes["severity"] = map[string]any{"from": risk.Severity, "to": *input.Severity}
		risk.Severity = *input.Severity
	}
	if input.CategoryID != nil {
		var oldCategoryID string
		if risk.CategoryID != nil {
			oldCategoryID = *risk.CategoryID
		}
		changes["category_id"] = map[string]any{"from": oldCategoryID, "to": *input.CategoryID}
		risk.CategoryID = input.CategoryID
	}
	if input.ReviewDate != nil {
		var oldReviewDate string
		if risk.ReviewDate != nil {
			oldReviewDate = risk.ReviewDate.Format("2006-01-02")
		}
		t, err := time.Parse("2006-01-02", *input.ReviewDate)
		if err == nil {
			changes["review_date"] = map[string]any{"from": oldReviewDate, "to": *input.ReviewDate}
			risk.ReviewDate = &t
		}
	}

	if err := h.risks.Update(c.Context(), risk); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update risk"})
	}

	// Log audit event if there were changes
	if len(changes) > 0 {
		h.audit.Create(c.Context(), "risk", risk.ID, models.AuditActionUpdated, changes, user.UserID)
	}

	return c.JSON(risk)
}

func (h *RiskHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	user := middleware.GetUserFromContext(c)

	// Log audit event before deletion
	h.audit.Create(c.Context(), "risk", id, models.AuditActionDeleted, nil, user.UserID)

	if err := h.risks.Delete(c.Context(), id); err != nil {
		if err == database.ErrRiskNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "risk not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete risk"})
	}

	return c.SendStatus(204)
}
