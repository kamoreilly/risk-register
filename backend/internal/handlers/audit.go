package handlers

import (
	"backend/internal/database"

	"github.com/gofiber/fiber/v2"
)

type AuditHandler struct {
	auditRepo database.AuditLogRepository
}

func NewAuditHandler(auditRepo database.AuditLogRepository) *AuditHandler {
	return &AuditHandler{auditRepo: auditRepo}
}

func (h *AuditHandler) ListByRisk(c *fiber.Ctx) error {
	riskID := c.Params("riskId")
	if riskID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "risk_id required"})
	}

	logs, err := h.auditRepo.ListByEntity(c.Context(), "risk", riskID, 50)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch audit logs"})
	}

	return c.JSON(fiber.Map{"data": logs})
}
