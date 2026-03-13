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
