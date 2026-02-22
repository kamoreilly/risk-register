package handlers

import (
	"backend/internal/database"
	"backend/internal/middleware"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type FrameworkHandler struct {
	frameworkRepo database.FrameworkRepository
	controlRepo   database.RiskFrameworkControlRepository
}

func NewFrameworkHandler(frameworkRepo database.FrameworkRepository, controlRepo database.RiskFrameworkControlRepository) *FrameworkHandler {
	return &FrameworkHandler{
		frameworkRepo: frameworkRepo,
		controlRepo:   controlRepo,
	}
}

// List returns all frameworks
func (h *FrameworkHandler) List(c *fiber.Ctx) error {
	frameworks, err := h.frameworkRepo.List(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch frameworks"})
	}
	return c.JSON(fiber.Map{"data": frameworks})
}

// Create creates a new framework (admin only)
func (h *FrameworkHandler) Create(c *fiber.Ctx) error {
	var input models.CreateFrameworkInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Validate required fields
	if input.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name is required"})
	}

	framework, err := h.frameworkRepo.Create(c.Context(), &input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create framework"})
	}
	return c.Status(201).JSON(framework)
}

// Update updates a framework (admin only)
func (h *FrameworkHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "framework id required"})
	}

	var input models.UpdateFrameworkInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if input.Name == nil && input.Description == nil {
		return c.Status(400).JSON(fiber.Map{"error": "at least one field must be provided"})
	}
	if input.Name != nil && *input.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name cannot be empty"})
	}

	framework, err := h.frameworkRepo.Update(c.Context(), id, &input)
	if err != nil {
		if err == database.ErrFrameworkNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "framework not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to update framework"})
	}
	return c.JSON(framework)
}

// Delete deletes a framework (admin only)
func (h *FrameworkHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "framework id required"})
	}

	if err := h.frameworkRepo.Delete(c.Context(), id); err != nil {
		if err == database.ErrFrameworkNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "framework not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete framework"})
	}
	return c.SendStatus(204)
}

type ControlHandler struct {
	controlRepo database.RiskFrameworkControlRepository
}

func NewControlHandler(controlRepo database.RiskFrameworkControlRepository) *ControlHandler {
	return &ControlHandler{controlRepo: controlRepo}
}

// ListControls returns all controls linked to a risk
func (h *ControlHandler) ListControls(c *fiber.Ctx) error {
	riskID := c.Params("riskId")
	if riskID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "risk_id required"})
	}

	controls, err := h.controlRepo.ListByRiskID(c.Context(), riskID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch controls"})
	}
	return c.JSON(fiber.Map{"data": controls})
}

// LinkControl links a framework control to a risk
func (h *ControlHandler) LinkControl(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	riskID := c.Params("riskId")
	if riskID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "risk_id required"})
	}

	var input models.LinkControlInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Validate required fields
	if input.FrameworkID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "framework_id is required"})
	}
	if input.ControlRef == "" {
		return c.Status(400).JSON(fiber.Map{"error": "control_ref is required"})
	}

	control, err := h.controlRepo.LinkControl(c.Context(), riskID, &input, user.UserID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to link control"})
	}
	return c.Status(201).JSON(control)
}

// UnlinkControl removes a control from a risk
func (h *ControlHandler) UnlinkControl(c *fiber.Ctx) error {
	controlID := c.Params("id")
	if controlID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "control_id required"})
	}

	if err := h.controlRepo.UnlinkControl(c.Context(), controlID); err != nil {
		if err == database.ErrFrameworkNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "control not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to unlink control"})
	}
	return c.SendStatus(204)
}
