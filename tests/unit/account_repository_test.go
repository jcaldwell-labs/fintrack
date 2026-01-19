package unit

import (
	"strings"
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
	// Clean database before each test by dropping and recreating table
	_ = suite.db.Migrator().DropTable(&models.Account{})
	_ = suite.db.AutoMigrate(&models.Account{})
}

// TestCreateAccount tests account creation
func (suite *AccountRepositoryTestSuite) TestCreateAccount() {
	// Given
	account := &models.Account{
		Name:                "Test Checking",
		Type:                models.AccountTypeChecking,
		Currency:            "USD",
		InitialBalanceCents: 100000, // $1000.00
		IsActive:            true,
	}

	// When
	err := suite.db.Create(account).Error

	// Then
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), account.ID)
	assert.Equal(suite.T(), "Test Checking", account.Name)
	assert.Equal(suite.T(), int64(100000), account.InitialBalanceCents)
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
	// Given - create test accounts (all active initially due to GORM default:true)
	accounts := []models.Account{
		{Name: "Checking", Type: models.AccountTypeChecking, Currency: "USD"},
		{Name: "Savings", Type: models.AccountTypeSavings, Currency: "USD"},
		{Name: "Credit", Type: models.AccountTypeCredit, Currency: "USD"},
	}

	for i := range accounts {
		suite.db.Create(&accounts[i])
	}

	// Now set Credit account to inactive using raw SQL to bypass GORM defaults
	suite.db.Exec("UPDATE accounts SET is_active = ? WHERE name = ?", false, "Credit")

	// When - get all accounts
	var allAccounts []models.Account
	err := suite.db.Find(&allAccounts).Error

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3, len(allAccounts))

	// When - get only active accounts
	var activeAccounts []models.Account
	err = suite.db.Where("is_active = ?", true).Find(&activeAccounts).Error

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(activeAccounts))
	for _, acc := range activeAccounts {
		assert.True(suite.T(), acc.IsActive)
	}

	// When - get only inactive accounts
	var inactiveAccounts []models.Account
	err = suite.db.Where("is_active = ?", false).Find(&inactiveAccounts).Error

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(inactiveAccounts))
	if len(inactiveAccounts) > 0 {
		assert.Equal(suite.T(), "Credit", inactiveAccounts[0].Name)
	}
}

// TestUpdateAccount tests account updates
func (suite *AccountRepositoryTestSuite) TestUpdateAccount() {
	// Given
	account := &models.Account{
		Name:                "Old Name",
		Type:                models.AccountTypeChecking,
		CurrentBalanceCents: 10000, // $100.00
		IsActive:            true,
	}
	suite.db.Create(account)

	// When
	account.Name = "New Name"
	account.CurrentBalanceCents = 20000 // $200.00
	err := suite.db.Save(account).Error

	// Then
	assert.NoError(suite.T(), err)

	var updated models.Account
	suite.db.First(&updated, account.ID)
	assert.Equal(suite.T(), "New Name", updated.Name)
	assert.Equal(suite.T(), int64(20000), updated.CurrentBalanceCents)
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
			Name:     accountType,
			Type:     accountType,
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

// TestDuplicateAccountNames tests that duplicate account names are rejected by database constraint
func (suite *AccountRepositoryTestSuite) TestDuplicateAccountNames() {
	// Given - create first account
	account1 := &models.Account{
		Name:     "Duplicate Account",
		Type:     models.AccountTypeChecking,
		Currency: "USD",
		IsActive: true,
	}
	err := suite.db.Create(account1).Error
	assert.NoError(suite.T(), err)

	// When - try to create account with same name
	account2 := &models.Account{
		Name:     "Duplicate Account",
		Type:     models.AccountTypeSavings,
		Currency: "USD",
		IsActive: true,
	}
	err = suite.db.Create(account2).Error

	// Then - should fail with unique constraint violation
	assert.Error(suite.T(), err)
	assert.True(suite.T(),
		strings.Contains(strings.ToLower(err.Error()), "duplicate") ||
			strings.Contains(strings.ToLower(err.Error()), "unique") ||
			strings.Contains(strings.ToLower(err.Error()), "constraint"),
		"Expected unique constraint error, got: %v", err)
}

// Run the test suite
func TestAccountRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AccountRepositoryTestSuite))
}
