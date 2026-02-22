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
