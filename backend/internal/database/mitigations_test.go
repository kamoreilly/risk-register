package database

import (
	"testing"
)

func TestMitigationRepository_Interface(t *testing.T) {
	// This test verifies the interface compiles
	var _ MitigationRepository = (*mitigationRepository)(nil)
	t.Log("MitigationRepository interface satisfied")
}
