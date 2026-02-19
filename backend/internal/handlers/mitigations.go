package handlers

import (
	"backend/internal/database"
	"backend/internal/middleware"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type MitigationHandler struct {
	mitigationRepo database.MitigationRepository
}

func NewMitigationHandler(mitigationRepo database.MitigationRepository) *MitigationHandler {
	return &MitigationHandler{mitigationRepo: mitigationRepo}
}

// List returns all mitigations for a specific risk
func (h *MitigationHandler) List(c *fiber.Ctx) error {
	riskID := c.Params("riskId")
	if riskID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "risk_id is required"})
	}

	mitigations, err := h.mitigationRepo.ListByRiskID(c.Context(), riskID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch mitigations"})
	}

	return c.JSON(mitigations)
}

// Create creates a new mitigation for a risk
func (h *MitigationHandler) Create(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	riskID := c.Params("riskId")
	if riskID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "risk_id is required"})
	}

	var input models.CreateMitigationInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Set risk_id from URL param
	input.RiskID = riskID

	// Validate required fields
	if input.Description == "" {
		return c.Status(400).JSON(fiber.Map{"error": "description is required"})
	}
	if input.Owner == "" {
		return c.Status(400).JSON(fiber.Map{"error": "owner is required"})
	}

	// Set default status to "planned"
	if input.Status == "" {
		input.Status = models.MitigationStatusPlanned
	}

	mitigation, err := h.mitigationRepo.Create(c.Context(), &input, user.UserID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create mitigation"})
	}

	return c.Status(201).JSON(mitigation)
}

// Update updates an existing mitigation
func (h *MitigationHandler) Update(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	riskID := c.Params("riskId")
	if riskID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "risk_id is required"})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "mitigation id is required"})
	}

	var input models.UpdateMitigationInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	mitigation, err := h.mitigationRepo.Update(c.Context(), id, &input, user.UserID)
	if err != nil {
		if err == database.ErrMitigationNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "mitigation not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to update mitigation"})
	}

	return c.JSON(mitigation)
}

// Delete removes a mitigation
func (h *MitigationHandler) Delete(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	riskID := c.Params("riskId")
	if riskID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "risk_id is required"})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "mitigation id is required"})
	}

	if err := h.mitigationRepo.Delete(c.Context(), id); err != nil {
		if err == database.ErrMitigationNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "mitigation not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete mitigation"})
	}

	return c.SendStatus(204)
}
