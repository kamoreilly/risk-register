package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type AIHandler struct{}

func NewAIHandler() *AIHandler {
	return &AIHandler{}
}

// SummarizeRequest is the input for summarize endpoint
type SummarizeRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Status      string `json:"status"`
}

// SummarizeResponse is the output
type SummarizeResponse struct {
	Summary string `json:"summary"`
}

// Summarize returns a stubbed summary of a risk
func (h *AIHandler) Summarize(c *fiber.Ctx) error {
	var req SummarizeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Stub: return a placeholder summary
	summary := "This risk requires attention. "
	if req.Severity == "critical" || req.Severity == "high" {
		summary += "Given its high severity, immediate action is recommended. "
	}
	summary += "Consider implementing appropriate mitigations to reduce the risk to an acceptable level."

	return c.JSON(SummarizeResponse{Summary: summary})
}

// DraftMitigationRequest is the input for draft-mitigation endpoint
type DraftMitigationRequest struct {
	RiskTitle       string `json:"risk_title"`
	RiskDescription string `json:"risk_description"`
	Severity        string `json:"severity"`
}

// DraftMitigationResponse is the output
type DraftMitigationResponse struct {
	Draft string `json:"draft"`
}

// DraftMitigation returns a stubbed mitigation draft
func (h *AIHandler) DraftMitigation(c *fiber.Ctx) error {
	var req DraftMitigationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Stub: return a placeholder mitigation draft
	draft := "Recommended mitigation actions:\n"
	draft += "1. Assess the current state and identify gaps\n"
	draft += "2. Implement technical controls to address the vulnerability\n"
	draft += "3. Establish monitoring and alerting\n"
	draft += "4. Document the mitigation plan and assign ownership\n"
	draft += "5. Schedule regular reviews to ensure effectiveness"

	return c.JSON(DraftMitigationResponse{Draft: draft})
}
