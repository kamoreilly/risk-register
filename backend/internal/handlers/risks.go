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
}

func NewRiskHandler(risks database.RiskRepository, categories database.CategoryRepository) *RiskHandler {
	return &RiskHandler{risks: risks, categories: categories}
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

	// Apply updates
	if input.Title != nil {
		risk.Title = *input.Title
	}
	if input.Description != nil {
		risk.Description = *input.Description
	}
	if input.OwnerID != nil {
		risk.OwnerID = *input.OwnerID
	}
	if input.Status != nil {
		risk.Status = *input.Status
	}
	if input.Severity != nil {
		risk.Severity = *input.Severity
	}
	if input.CategoryID != nil {
		risk.CategoryID = input.CategoryID
	}
	if input.ReviewDate != nil {
		t, err := time.Parse("2006-01-02", *input.ReviewDate)
		if err == nil {
			risk.ReviewDate = &t
		}
	}

	if err := h.risks.Update(c.Context(), risk); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update risk"})
	}

	return c.JSON(risk)
}

func (h *RiskHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.risks.Delete(c.Context(), id); err != nil {
		if err == database.ErrRiskNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "risk not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete risk"})
	}

	return c.SendStatus(204)
}
