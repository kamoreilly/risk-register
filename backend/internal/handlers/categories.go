package handlers

import (
	"backend/internal/database"

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
