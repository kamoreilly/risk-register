package handlers

import (
	"errors"

	"backend/internal/database"
	"backend/internal/middleware"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
)

type FrameworkControlHandler struct {
	controls database.FrameworkControlRepository
}

func NewFrameworkControlHandler(controls database.FrameworkControlRepository) *FrameworkControlHandler {
	return &FrameworkControlHandler{controls: controls}
}

func (h *FrameworkControlHandler) List(c *fiber.Ctx) error {
	frameworkID := c.Query("framework_id")
	search := c.Query("search")

	controls, err := h.controls.List(c.Context(), frameworkID, search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch controls"})
	}

	return c.JSON(fiber.Map{"data": controls})
}

func (h *FrameworkControlHandler) Create(c *fiber.Ctx) error {
	var input models.CreateFrameworkControlInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if input.FrameworkID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "framework_id is required"})
	}
	if input.ControlRef == "" {
		return c.Status(400).JSON(fiber.Map{"error": "control_ref is required"})
	}
	if input.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "title is required"})
	}

	control, err := h.controls.Create(c.Context(), &input)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return c.Status(409).JSON(fiber.Map{"error": "Control reference already exists for this framework"})
			}
			if pgErr.Code == "23503" {
				return c.Status(400).JSON(fiber.Map{"error": "Framework not found"})
			}
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to create control"})
	}

	return c.Status(201).JSON(control)
}

func (h *FrameworkControlHandler) ListLinkedRisks(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "control id required"})
	}

	risks, err := h.controls.ListLinkedRisks(c.Context(), id)
	if err != nil {
		return mapFrameworkControlError(c, err, "failed to fetch linked risks")
	}

	return c.JSON(fiber.Map{"data": risks})
}

func (h *FrameworkControlHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "control id required"})
	}

	var input models.UpdateFrameworkControlInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if input.ControlRef == nil && input.Title == nil && input.Description == nil {
		return c.Status(400).JSON(fiber.Map{"error": "at least one field must be provided"})
	}
	if input.ControlRef != nil && *input.ControlRef == "" {
		return c.Status(400).JSON(fiber.Map{"error": "control_ref cannot be empty"})
	}
	if input.Title != nil && *input.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "title cannot be empty"})
	}

	control, err := h.controls.Update(c.Context(), id, &input)
	if err != nil {
		if errors.Is(err, database.ErrFrameworkControlNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "control not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to update control"})
	}

	return c.JSON(control)
}

func (h *FrameworkControlHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "control id required"})
	}

	if err := h.controls.Delete(c.Context(), id); err != nil {
		return mapFrameworkControlError(c, err, "failed to delete control")
	}

	return c.SendStatus(204)
}

type ControlHandler struct {
	controlRepo database.RiskFrameworkControlRepository
}

func NewControlHandler(controlRepo database.RiskFrameworkControlRepository) *ControlHandler {
	return &ControlHandler{controlRepo: controlRepo}
}

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

	if input.FrameworkControlID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "framework_control_id is required"})
	}

	control, err := h.controlRepo.LinkControl(c.Context(), riskID, &input, user.UserID)
	if err != nil {
		return mapFrameworkControlError(c, err, "failed to link control")
	}
	return c.Status(201).JSON(control)
}

func (h *ControlHandler) UnlinkControl(c *fiber.Ctx) error {
	controlID := c.Params("id")
	if controlID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "control_id required"})
	}

	if err := h.controlRepo.UnlinkControl(c.Context(), controlID); err != nil {
		if errors.Is(err, database.ErrFrameworkControlNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "control not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to unlink control"})
	}
	return c.SendStatus(204)
}
