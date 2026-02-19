package handlers

import (
	"context"

	"backend/internal/auth"
	"backend/internal/database"
	"backend/internal/middleware"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	users database.UserRepository
}

func NewAuthHandler(users database.UserRepository) *AuthHandler {
	return &AuthHandler{users: users}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var input models.RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate input
	if input.Email == "" || input.Password == "" || input.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "email, password, and name are required",
		})
	}

	if len(input.Password) < 8 {
		return c.Status(400).JSON(fiber.Map{
			"error": "password must be at least 8 characters",
		})
	}

	// Check if user exists
	existing, _ := h.users.FindByEmail(context.Background(), input.Email)
	if existing != nil {
		return c.Status(409).JSON(fiber.Map{
			"error": "email already registered",
		})
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to hash password",
		})
	}

	// Create user
	user := &models.User{
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Name:         input.Name,
		Role:         models.RoleMember,
	}

	if err := h.users.Create(context.Background(), user); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to create user",
		})
	}

	// Generate token
	token, err := auth.GenerateToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to generate token",
		})
	}

	return c.Status(201).JSON(models.AuthResponse{
		User:  user,
		Token: token,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input models.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Find user
	user, err := h.users.FindByEmail(context.Background(), input.Email)
	if err != nil || user == nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}

	// Check password
	if !auth.CheckPassword(input.Password, user.PasswordHash) {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}

	// Generate token
	token, err := auth.GenerateToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to generate token",
		})
	}

	return c.JSON(models.AuthResponse{
		User:  user,
		Token: token,
	})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "not authenticated",
		})
	}

	fullUser, err := h.users.FindByID(context.Background(), user.UserID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return c.JSON(fullUser)
}
