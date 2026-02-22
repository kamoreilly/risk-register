package server

import (
	"backend/internal/auth"
	"backend/internal/database"
	"backend/internal/models"
	"context"
	"errors"
	"log"
)

func (s *FiberServer) SeedDevUsers() {
	ctx := context.Background()

	users := []struct {
		Email    string
		Password string
		Name     string
		Role     models.UserRole
	}{
		{
			Email:    "admin@example.com",
			Password: "password123",
			Name:     "Admin User",
			Role:     models.RoleAdmin,
		},
		{
			Email:    "member@example.com",
			Password: "password123",
			Name:     "Member User",
			Role:     models.RoleMember,
		},
	}

	for _, u := range users {
		_, err := s.users.FindByEmail(ctx, u.Email)
		if err == nil {
			// User exists
			continue
		}

		if !errors.Is(err, database.ErrUserNotFound) {
			log.Printf("Error checking user %s: %v", u.Email, err)
			continue
		}

		// User not found, create it
		hash, err := auth.HashPassword(u.Password)
		if err != nil {
			log.Printf("Error hashing password for %s: %v", u.Email, err)
			continue
		}

		user := &models.User{
			Email:        u.Email,
			PasswordHash: hash,
			Name:         u.Name,
			Role:         u.Role,
		}

		if err := s.users.Create(ctx, user); err != nil {
			log.Printf("Error creating user %s: %v", u.Email, err)
		} else {
			log.Printf("Seeded dev user: %s (%s)", u.Email, u.Role)
		}
	}
}
