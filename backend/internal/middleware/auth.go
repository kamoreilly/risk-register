package middleware

import (
	"strings"

	"backend/internal/auth"

	"github.com/gofiber/fiber/v2"
)

type contextKey string

const UserKey contextKey = "user"

type UserClaims struct {
	UserID string
	Email  string
	Role   string
}

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "missing authorization header",
		})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid authorization format",
		})
	}

	claims, err := auth.ValidateToken(parts[1])
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid or expired token",
		})
	}

	c.Locals(UserKey, &UserClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   string(claims.Role),
	})

	return c.Next()
}

func GetUserFromContext(c *fiber.Ctx) *UserClaims {
	user, ok := c.Locals(UserKey).(*UserClaims)
	if !ok {
		return nil
	}
	return user
}

func RequireAdmin(c *fiber.Ctx) error {
	user := GetUserFromContext(c)
	if user == nil || user.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "admin access required",
		})
	}
	return c.Next()
}
