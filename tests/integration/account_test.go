//go:build integration

package integration

import (
	"testing"

	"github.com/fintrack/fintrack/internal/db/repositories"
	"github.com/fintrack/fintrack/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountRepository_CRUD(t *testing.T) {
	// Cleanup before test
	cleanupTable(t, "accounts")
	defer cleanupTable(t, "accounts")

	repo := repositories.NewAccountRepository()

	// Test Create
	t.Run("Create", func(t *testing.T) {
		account := &models.Account{
			Name:    "Test Checking",
			Type:    "checking",
			Balance: 1000.00,
		}

		err := repo.Create(account)
		require.NoError(t, err)
		assert.NotZero(t, account.ID)
	})

	// Test GetAll
	t.Run("GetAll", func(t *testing.T) {
		accounts, err := repo.GetAll()
		require.NoError(t, err)
		assert.Len(t, accounts, 1)
		assert.Equal(t, "Test Checking", accounts[0].Name)
	})

	// Test GetByID
	t.Run("GetByID", func(t *testing.T) {
		accounts, _ := repo.GetAll()
		require.Len(t, accounts, 1)

		account, err := repo.GetByID(accounts[0].ID)
		require.NoError(t, err)
		assert.Equal(t, "Test Checking", account.Name)
		assert.Equal(t, 1000.00, account.Balance)
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		accounts, _ := repo.GetAll()
		require.Len(t, accounts, 1)

		accounts[0].Balance = 1500.00
		err := repo.Update(accounts[0])
		require.NoError(t, err)

		updated, err := repo.GetByID(accounts[0].ID)
		require.NoError(t, err)
		assert.Equal(t, 1500.00, updated.Balance)
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		accounts, _ := repo.GetAll()
		require.Len(t, accounts, 1)

		err := repo.Delete(accounts[0].ID)
		require.NoError(t, err)

		remaining, err := repo.GetAll()
		require.NoError(t, err)
		assert.Len(t, remaining, 0)
	})
}

func TestAccountRepository_GetByID_NotFound(t *testing.T) {
	cleanupTable(t, "accounts")
	defer cleanupTable(t, "accounts")

	repo := repositories.NewAccountRepository()
	account, err := repo.GetByID(99999)
	assert.Error(t, err)
	assert.Nil(t, account)
}
