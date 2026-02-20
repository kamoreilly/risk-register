package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"backend/internal/middleware"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Public routes
	s.App.Get("/", s.HelloWorldHandler)
	s.App.Get("/health", s.healthHandler)

	// Auth routes (public)
	auth := s.App.Group("/api/v1/auth")
	auth.Post("/register", s.auth.Register)
	auth.Post("/login", s.auth.Login)

	// Protected routes
	protected := s.App.Group("/api/v1", middleware.AuthMiddleware)
	protected.Get("/auth/me", s.auth.Me)

	// Dashboard routes
	dashboard := protected.Group("/dashboard")
	dashboard.Get("/summary", s.dashboardHandler.Summary)
	dashboard.Get("/reviews/upcoming", s.dashboardHandler.UpcomingReviews)
	dashboard.Get("/reviews/overdue", s.dashboardHandler.OverdueReviews)

	// Category routes (public read)
	categories := protected.Group("/categories")
	categories.Get("/", s.categoryHandler.List)

	// Risk routes
	risks := protected.Group("/risks")
	risks.Get("/", s.riskHandler.List)
	risks.Post("/", s.riskHandler.Create)
	risks.Get("/:id", s.riskHandler.Get)
	risks.Put("/:id", s.riskHandler.Update)
	risks.Delete("/:id", s.riskHandler.Delete)

	// Nested mitigation routes under a specific risk
	risks.Get("/:riskId/mitigations", s.mitigationHandler.List)
	risks.Post("/:riskId/mitigations", s.mitigationHandler.Create)
	risks.Put("/:riskId/mitigations/:id", s.mitigationHandler.Update)
	risks.Delete("/:riskId/mitigations/:id", s.mitigationHandler.Delete)

	// Framework routes (public list, admin create)
	protected.Get("/frameworks", s.frameworkHandler.List)
	protected.Post("/frameworks", middleware.RequireAdmin, s.frameworkHandler.Create)

	// Nested control routes under a specific risk
	risks.Get("/:riskId/controls", s.controlHandler.ListControls)
	risks.Post("/:riskId/controls", s.controlHandler.LinkControl)
	risks.Delete("/:riskId/controls/:id", s.controlHandler.UnlinkControl)

	// AI routes (stubbed)
	ai := protected.Group("/ai")
	ai.Post("/summarize", s.aiHandler.Summarize)
	ai.Post("/draft-mitigation", s.aiHandler.DraftMitigation)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Risk Register API",
		"version": "1.0.0",
	}
	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
