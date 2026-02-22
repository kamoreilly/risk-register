package handlers

import (
	"errors"

	"backend/internal/database"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categories database.CategoryRepository
}

func NewCategoryHandler(categories database.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{categories: categories}
}

func (h *CategoryHandler) List(c *fiber.Ctx) error {
	categories, err := h.categories.List(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch categories"})
	}
	return c.JSON(categories)
}

func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	var input models.CreateCategoryInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if input.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name is required"})
	}

	category, err := h.categories.Create(c.Context(), &input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create category"})
	}
	return c.Status(201).JSON(category)
}

func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "category id required"})
	}

	var input models.UpdateCategoryInput
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
		if errors.Is(err, database.ErrCategoryNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "category not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to update category"})
	}
	return c.JSON(category)
}

func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "category id required"})
	}

	if err := h.categories.Delete(c.Context(), id); err != nil {
		if errors.Is(err, database.ErrCategoryNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "category not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete category"})
	}
	return c.SendStatus(204)
}
