package database

import (
	"context"
	"testing"

	"backend/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	s := New().(*service)
	repo := NewCategoryRepository(s.db)
	ctx := context.Background()

	// 1. Create
	input := &models.CreateCategoryInput{
		Name:        "Integration Category " + uuid.New().String(),
		Description: "Test Desc",
	}
	category, err := repo.Create(ctx, input)
	require.NoError(t, err)
	assert.NotEmpty(t, category.ID)
	assert.Equal(t, input.Name, category.Name)

	// 2. Get
	fetched, err := repo.FindByID(ctx, category.ID)
	require.NoError(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, category.Name, fetched.Name)

	// 3. Update
	newName := "Updated Name " + uuid.New().String()
	updateInput := &models.UpdateCategoryInput{
		Name: &newName,
	}
	updated, err := repo.Update(ctx, category.ID, updateInput)
	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)

	// 4. List
	list, err := repo.List(ctx)
	require.NoError(t, err)
	found := false
	for _, c := range list {
		if c.ID == category.ID {
			found = true
			break
		}
	}
	assert.True(t, found)

	// 5. Delete
	err = repo.Delete(ctx, category.ID)
	require.NoError(t, err)

	fetched, err = repo.FindByID(ctx, category.ID)
	require.NoError(t, err)
	assert.Nil(t, fetched)
}
