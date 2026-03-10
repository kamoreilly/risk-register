package handlers

import (
	"errors"

	"backend/internal/database"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type FrameworkHandler struct {
	frameworkRepo database.FrameworkRepository
}

func NewFrameworkHandler(frameworkRepo database.FrameworkRepository) *FrameworkHandler {
	return &FrameworkHandler{
		frameworkRepo: frameworkRepo,
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
		if errors.Is(err, database.ErrFrameworkNotFound) {
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
		if errors.Is(err, database.ErrFrameworkNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "framework not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete framework"})
	}
	return c.SendStatus(204)
}

func mapFrameworkControlError(c *fiber.Ctx, err error, defaultMessage string) error {
	if errors.Is(err, database.ErrFrameworkControlNotFound) {
		return c.Status(404).JSON(fiber.Map{"error": "control not found"})
	}
	if errors.Is(err, database.ErrFrameworkControlInUse) {
		return c.Status(409).JSON(fiber.Map{"error": "control is linked to one or more risks"})
	}
	return c.Status(500).JSON(fiber.Map{"error": defaultMessage})
}
