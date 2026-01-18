//go:build integration

package integration

import (
	"testing"

	"github.com/fintrack/fintrack/internal/db/repositories"
	"github.com/fintrack/fintrack/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryRepository_CRUD(t *testing.T) {
	// Cleanup before test
	cleanupTable(t, "categories")
	defer cleanupTable(t, "categories")

	repo := repositories.NewCategoryRepository()

	// Test Create
	t.Run("Create", func(t *testing.T) {
		category := &models.Category{
			Name: "Groceries",
			Type: "expense",
		}

		err := repo.Create(category)
		require.NoError(t, err)
		assert.NotZero(t, category.ID)
	})

	// Test GetAll
	t.Run("GetAll", func(t *testing.T) {
		categories, err := repo.GetAll()
		require.NoError(t, err)
		assert.Len(t, categories, 1)
		assert.Equal(t, "Groceries", categories[0].Name)
	})

	// Test GetByID
	t.Run("GetByID", func(t *testing.T) {
		categories, _ := repo.GetAll()
		require.Len(t, categories, 1)

		category, err := repo.GetByID(categories[0].ID)
		require.NoError(t, err)
		assert.Equal(t, "Groceries", category.Name)
		assert.Equal(t, "expense", category.Type)
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		categories, _ := repo.GetAll()
		require.Len(t, categories, 1)

		categories[0].Name = "Food & Groceries"
		err := repo.Update(categories[0])
		require.NoError(t, err)

		updated, err := repo.GetByID(categories[0].ID)
		require.NoError(t, err)
		assert.Equal(t, "Food & Groceries", updated.Name)
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		categories, _ := repo.GetAll()
		require.Len(t, categories, 1)

		err := repo.Delete(categories[0].ID)
		require.NoError(t, err)

		remaining, err := repo.GetAll()
		require.NoError(t, err)
		assert.Len(t, remaining, 0)
	})
}

func TestCategoryRepository_GetByID_NotFound(t *testing.T) {
	cleanupTable(t, "categories")
	defer cleanupTable(t, "categories")

	repo := repositories.NewCategoryRepository()
	category, err := repo.GetByID(99999)
	assert.Error(t, err)
	assert.Nil(t, category)
}
