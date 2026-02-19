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
