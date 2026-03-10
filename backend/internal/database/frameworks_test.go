package database

import (
	"context"
	"testing"

	"backend/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFrameworkRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	s := New().(*service)
	repo := NewFrameworkRepository(s.db)
	ctx := context.Background()

	// 1. Create
	input := &models.CreateFrameworkInput{
		Name:        "Test Framework " + uuid.New().String(),
		Description: "Test Desc",
	}
	framework, err := repo.Create(ctx, input)
	require.NoError(t, err)
	assert.NotEmpty(t, framework.ID)

	// 2. Get
	fetched, err := repo.GetByID(ctx, framework.ID)
	require.NoError(t, err)
	assert.Equal(t, framework.Name, fetched.Name)

	// 3. Update
	newName := "Updated Framework " + uuid.New().String()
	updateInput := &models.UpdateFrameworkInput{
		Name: &newName,
	}
	updated, err := repo.Update(ctx, framework.ID, updateInput)
	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)

	// 4. List
	list, err := repo.List(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)

	// 5. Delete
	err = repo.Delete(ctx, framework.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, framework.ID)
	assert.Error(t, err)
	assert.Equal(t, ErrFrameworkNotFound, err)
}

func TestFrameworkControlRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	s := New().(*service)
	frameworkRepo := NewFrameworkRepository(s.db)
	controlRepo := NewFrameworkControlRepository(s.db)
	ctx := context.Background()

	framework, err := frameworkRepo.Create(ctx, &models.CreateFrameworkInput{
		Name:        "Framework Controls " + uuid.New().String(),
		Description: "Test framework for controls",
	})
	require.NoError(t, err)

	created, err := controlRepo.Create(ctx, &models.CreateFrameworkControlInput{
		FrameworkID: framework.ID,
		ControlRef:  "AC-" + uuid.New().String()[:8],
		Title:       "Access control policy",
		Description: "Ensure documented access policy exists",
	})
	require.NoError(t, err)
	assert.Equal(t, framework.ID, created.FrameworkID)
	assert.Zero(t, created.LinkedRiskCount)

	fetched, err := controlRepo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.Title, fetched.Title)

	list, err := controlRepo.List(ctx, framework.ID, "access")
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, created.ID, list[0].ID)

	updatedTitle := "Updated access control policy"
	updatedDescription := "Updated description"
	updated, err := controlRepo.Update(ctx, created.ID, &models.UpdateFrameworkControlInput{
		Title:       &updatedTitle,
		Description: &updatedDescription,
	})
	require.NoError(t, err)
	assert.Equal(t, updatedTitle, updated.Title)
	assert.Equal(t, updatedDescription, updated.Description)

	linkedRisks, err := controlRepo.ListLinkedRisks(ctx, created.ID)
	require.NoError(t, err)
	assert.Len(t, linkedRisks, 0)

	err = controlRepo.Delete(ctx, created.ID)
	require.NoError(t, err)

	_, err = controlRepo.GetByID(ctx, created.ID)
	assert.ErrorIs(t, err, ErrFrameworkControlNotFound)
}

func TestRiskFrameworkControlRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	s := New().(*service)
	frameworkRepo := NewFrameworkRepository(s.db)
	controlRepo := NewFrameworkControlRepository(s.db)
	riskControlRepo := NewRiskFrameworkControlRepository(s.db)
	riskRepo := NewRiskRepository(s.db)
	userRepo := NewUserRepository(s.db)
	categoryRepo := NewCategoryRepository(s.db)
	ctx := context.Background()

	framework, err := frameworkRepo.Create(ctx, &models.CreateFrameworkInput{
		Name:        "Linked Controls " + uuid.New().String(),
		Description: "Framework for linkage tests",
	})
	require.NoError(t, err)

	definition, err := controlRepo.Create(ctx, &models.CreateFrameworkControlInput{
		FrameworkID: framework.ID,
		ControlRef:  "CC-" + uuid.New().String()[:8],
		Title:       "Control link definition",
	})
	require.NoError(t, err)

	user := &models.User{
		ID:           uuid.New().String(),
		Email:        "framework-control-" + uuid.New().String() + "@example.com",
		PasswordHash: "hash",
		Name:         "Framework Control Tester",
		Role:         models.RoleMember,
	}
	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

	category, err := categoryRepo.Create(ctx, &models.CreateCategoryInput{
		Name:        "Control Category " + uuid.New().String(),
		Description: "Category for linked risks",
	})
	require.NoError(t, err)

	risk := &models.Risk{
		Title:       "Linked Risk",
		Description: "Risk linked to a framework control",
		OwnerID:     user.ID,
		Status:      models.StatusOpen,
		Severity:    models.SeverityHigh,
		CategoryID:  &category.ID,
		CreatedBy:   user.ID,
		UpdatedBy:   user.ID,
	}
	err = riskRepo.Create(ctx, risk)
	require.NoError(t, err)

	linked, err := riskControlRepo.LinkControl(ctx, risk.ID, &models.LinkControlInput{
		FrameworkControlID: definition.ID,
		Notes:              "Mapped during integration test",
	}, user.ID)
	require.NoError(t, err)
	assert.Equal(t, definition.ID, linked.FrameworkControlID)
	assert.Equal(t, framework.Name, linked.FrameworkName)
	assert.Equal(t, definition.Title, linked.ControlTitle)

	riskControls, err := riskControlRepo.ListByRiskID(ctx, risk.ID)
	require.NoError(t, err)
	require.Len(t, riskControls, 1)
	assert.Equal(t, linked.ID, riskControls[0].ID)

	linkedRisks, err := controlRepo.ListLinkedRisks(ctx, definition.ID)
	require.NoError(t, err)
	require.Len(t, linkedRisks, 1)
	assert.Equal(t, risk.ID, linkedRisks[0].ID)
	assert.Equal(t, risk.Title, linkedRisks[0].Title)

	err = controlRepo.Delete(ctx, definition.ID)
	assert.ErrorIs(t, err, ErrFrameworkControlInUse)

	err = riskControlRepo.UnlinkControl(ctx, linked.ID)
	require.NoError(t, err)

	err = controlRepo.Delete(ctx, definition.ID)
	require.NoError(t, err)
}
