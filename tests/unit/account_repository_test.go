package unit

import (
	"testing"

	"github.com/fintrack/fintrack/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// AccountRepositoryTestSuite is the test suite for account repository
type AccountRepositoryTestSuite struct {
	suite.Suite
	db *gorm.DB
}

// SetupSuite runs once before all tests
func (suite *AccountRepositoryTestSuite) SetupSuite() {
	// Use in-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)

	suite.db = db

	// Run migrations
	err = db.AutoMigrate(&models.Account{})
	assert.NoError(suite.T(), err)
}

// SetupTest runs before each test
func (suite *AccountRepositoryTestSuite) SetupTest() {
	// Clean database before each test
	suite.db.Exec("DELETE FROM accounts")
	// Reset auto-increment counter for SQLite
	suite.db.Exec("DELETE FROM sqlite_sequence WHERE name='accounts'")
}

// TestCreateAccount tests account creation
func (suite *AccountRepositoryTestSuite) TestCreateAccount() {
	// Given
	account := &models.Account{
		Name:           "Test Checking",
		Type:           models.AccountTypeChecking,
		Currency:       "USD",
		InitialBalance: 1000.0,
		IsActive:       true,
	}

	// When
	err := suite.db.Create(account).Error

	// Then
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), account.ID)
	assert.Equal(suite.T(), "Test Checking", account.Name)
	assert.Equal(suite.T(), 1000.0, account.InitialBalance)
}

// TestGetAccountByID tests retrieving account by ID
func (suite *AccountRepositoryTestSuite) TestGetAccountByID() {
	// Given
	account := &models.Account{
		Name:     "Test Savings",
		Type:     models.AccountTypeSavings,
		Currency: "USD",
		IsActive: true,
	}
	suite.db.Create(account)

	// When
	var retrieved models.Account
	err := suite.db.First(&retrieved, account.ID).Error

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), account.ID, retrieved.ID)
	assert.Equal(suite.T(), "Test Savings", retrieved.Name)
}

// TestListAccounts tests listing all accounts
func (suite *AccountRepositoryTestSuite) TestListAccounts() {
	// Create active accounts
	activeAccounts := []models.Account{
		{Name: "ListTest_Checking", Type: models.AccountTypeChecking, Currency: "USD", IsActive: true},
		{Name: "ListTest_Savings", Type: models.AccountTypeSavings, Currency: "USD", IsActive: true},
	}
	for i := range activeAccounts {
		err := suite.db.Create(&activeAccounts[i]).Error
		assert.NoError(suite.T(), err)
	}

	// Create inactive account using raw SQL to bypass GORM defaults
	suite.db.Exec("INSERT INTO accounts (name, type, currency, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))",
		"ListTest_Credit", models.AccountTypeCredit, "USD", false)

	// When - get all accounts with our test names
	var allAccounts []models.Account
	err := suite.db.Where("name LIKE ?", "ListTest_%").Find(&allAccounts).Error

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3, len(allAccounts))

	// When - get only active accounts with our test names
	var activeOnly []models.Account
	err = suite.db.Where("is_active = ? AND name LIKE ?", true, "ListTest_%").Find(&activeOnly).Error

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(activeOnly))
}

// TestUpdateAccount tests account updates
func (suite *AccountRepositoryTestSuite) TestUpdateAccount() {
	// Given
	account := &models.Account{
		Name:           "Old Name",
		Type:           models.AccountTypeChecking,
		CurrentBalance: 100.0,
		IsActive:       true,
	}
	suite.db.Create(account)

	// When
	account.Name = "New Name"
	account.CurrentBalance = 200.0
	err := suite.db.Save(account).Error

	// Then
	assert.NoError(suite.T(), err)

	var updated models.Account
	suite.db.First(&updated, account.ID)
	assert.Equal(suite.T(), "New Name", updated.Name)
	assert.Equal(suite.T(), 200.0, updated.CurrentBalance)
}

// TestDeleteAccount tests account soft deletion
func (suite *AccountRepositoryTestSuite) TestDeleteAccount() {
	// Given
	account := &models.Account{
		Name:     "To Delete",
		Type:     models.AccountTypeChecking,
		IsActive: true,
	}
	suite.db.Create(account)

	// When - soft delete (set is_active = false)
	err := suite.db.Model(&models.Account{}).
		Where("id = ?", account.ID).
		Update("is_active", false).Error

	// Then
	assert.NoError(suite.T(), err)

	var deleted models.Account
	suite.db.First(&deleted, account.ID)
	assert.False(suite.T(), deleted.IsActive)
}

// TestAccountTypes tests different account types
func (suite *AccountRepositoryTestSuite) TestAccountTypes() {
	accountTypes := []string{
		models.AccountTypeChecking,
		models.AccountTypeSavings,
		models.AccountTypeCredit,
		models.AccountTypeCash,
		models.AccountTypeInvestment,
		models.AccountTypeLoan,
	}

	for i, accountType := range accountTypes {
		account := &models.Account{
			Name:     accountType + "_account", // Make names unique
			Type:     accountType,
			Currency: "USD",
			IsActive: true,
		}
		err := suite.db.Create(account).Error
		assert.NoError(suite.T(), err, "Failed to create account type: %s", accountType)
		assert.NotZero(suite.T(), account.ID)

		// Verify it was created
		var retrieved models.Account
		suite.db.First(&retrieved, account.ID)
		assert.Equal(suite.T(), accountType, retrieved.Type, "Account type %d mismatch", i)
	}
}

// TestDuplicateAccountNames tests that duplicate active account names are prevented by DB constraint
func (suite *AccountRepositoryTestSuite) TestDuplicateAccountNames() {
	// Given
	account1 := &models.Account{
		Name:     "Duplicate Name",
		Type:     models.AccountTypeChecking,
		IsActive: true,
	}
	suite.db.Create(account1)

	// When - try to create another account with same name (active)
	account2 := &models.Account{
		Name:     "Duplicate Name",
		Type:     models.AccountTypeSavings,
		IsActive: true,
	}

	err := suite.db.Create(account2).Error

	// Then - should fail with unique constraint error
	assert.Error(suite.T(), err)
}

// Run the test suite
func TestAccountRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AccountRepositoryTestSuite))
}
