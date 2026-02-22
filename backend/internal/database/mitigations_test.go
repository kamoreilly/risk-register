package database

import (
	"context"
	"testing"

	"backend/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMitigationRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	s := New().(*service)
	riskRepo := NewRiskRepository(s.db)
	mitigationRepo := NewMitigationRepository(s.db)
	userRepo := NewUserRepository(s.db)
	catRepo := NewCategoryRepository(s.db)

	ctx := context.Background()

	// 1. Setup Data
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        "test-mitigation-" + uuid.New().String() + "@example.com",
		PasswordHash: "hash",
		Name:         "Mitigation Tester",
		Role:         models.RoleMember,
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	categoryInput := &models.CreateCategoryInput{
		Name:        "Mitigation Category " + uuid.New().String(),
		Description: "Test Desc",
	}
	category, err := catRepo.Create(ctx, categoryInput)
	require.NoError(t, err)

	risk := &models.Risk{
		Title:       "Test Risk for Mitigation",
		Description: "Description",
		OwnerID:     user.ID,
		Status:      models.StatusOpen,
		Severity:    models.SeverityHigh,
		CategoryID:  &category.ID,
		CreatedBy:   user.ID,
		UpdatedBy:   user.ID,
	}
	err = riskRepo.Create(ctx, risk)
	require.NoError(t, err)

	// 2. Create Mitigation
	mitigationInput := &models.CreateMitigationInput{
		RiskID:      risk.ID,
		Description: "Mitigation Plan",
		Owner:       "Mitigation Owner",
		Status:      models.MitigationStatusPlanned,
	}
	mitigation, err := mitigationRepo.Create(ctx, mitigationInput, user.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, mitigation.ID)
	assert.Equal(t, mitigationInput.Description, mitigation.Description)

	// 3. Get Mitigation
	fetched, err := mitigationRepo.FindByID(ctx, mitigation.ID)
	require.NoError(t, err)
	assert.Equal(t, mitigation.Description, fetched.Description)
	assert.Equal(t, mitigation.RiskID, fetched.RiskID)

	// 4. Update Mitigation
	newDesc := "Updated Mitigation Plan"
	updateInput := &models.UpdateMitigationInput{
		Description: &newDesc,
	}
	updated, err := mitigationRepo.Update(ctx, mitigation.ID, updateInput, user.ID)
	require.NoError(t, err)
	assert.Equal(t, newDesc, updated.Description)

	// 5. List Mitigations
	list, err := mitigationRepo.ListByRiskID(ctx, risk.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)

	// 6. Delete Mitigation
	err = mitigationRepo.Delete(ctx, mitigation.ID)
	require.NoError(t, err)

	fetched, err = mitigationRepo.FindByID(ctx, mitigation.ID)
	assert.Error(t, err)
	assert.Equal(t, ErrMitigationNotFound, err)
}
