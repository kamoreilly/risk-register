package database

import (
	"context"
	"testing"

	"backend/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRiskRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	s := New().(*service)
	riskRepo := NewRiskRepository(s.db)
	userRepo := NewUserRepository(s.db)
	catRepo := NewCategoryRepository(s.db)

	ctx := context.Background()

	// 1. Setup Data
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        "test-risk-" + uuid.New().String() + "@example.com",
		PasswordHash: "hash",
		Name:         "Risk Tester",
		Role:         models.RoleMember,
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	categoryInput := &models.CreateCategoryInput{
		Name:        "Risk Category " + uuid.New().String(),
		Description: "Test Desc",
	}
	category, err := catRepo.Create(ctx, categoryInput)
	require.NoError(t, err)

	// 2. Create Risk
	risk := &models.Risk{
		Title:       "Test Risk",
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
	assert.NotEmpty(t, risk.ID)

	// 3. Get Risk
	fetchedRisk, err := riskRepo.FindByID(ctx, risk.ID)
	require.NoError(t, err)
	assert.Equal(t, risk.Title, fetchedRisk.Title)
	assert.Equal(t, risk.OwnerID, fetchedRisk.OwnerID)
	assert.Equal(t, risk.CategoryID, fetchedRisk.CategoryID)

	// 4. Update Risk
	newTitle := "Updated Risk Title"
	risk.Title = newTitle
	err = riskRepo.Update(ctx, risk)
	require.NoError(t, err)

	fetchedRisk, err = riskRepo.FindByID(ctx, risk.ID)
	require.NoError(t, err)
	assert.Equal(t, newTitle, fetchedRisk.Title)

	// 5. List Risks
	params := &models.RiskListParams{
		Page:  1,
		Limit: 10,
	}
	listResp, err := riskRepo.List(ctx, params)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(listResp.Data), 1)

	// 6. Delete Risk
	err = riskRepo.Delete(ctx, risk.ID)
	require.NoError(t, err)

	_, err = riskRepo.FindByID(ctx, risk.ID)
	assert.Error(t, err)
	assert.Equal(t, ErrRiskNotFound, err)
}
