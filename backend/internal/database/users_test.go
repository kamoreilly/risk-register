package database

import (
	"testing"
)

func TestUserRepository_Interface(t *testing.T) {
	// This test verifies the interface compiles
	var _ UserRepository = (*userRepository)(nil)
	t.Log("UserRepository interface satisfied")
}
