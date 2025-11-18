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
	// Given - create active accounts
	checking := &models.Account{Name: "Checking", Type: models.AccountTypeChecking, IsActive: true}
	_ = suite.db.Create(checking).Error

	savings := &models.Account{Name: "Savings", Type: models.AccountTypeSavings, IsActive: true}
	_ = suite.db.Create(savings).Error

	// Create inactive account by first creating it active, then updating
	credit := &models.Account{Name: "Credit", Type: models.AccountTypeCredit, IsActive: true}
	_ = suite.db.Create(credit).Error
	_ = suite.db.Model(credit).Update("is_active", false).Error

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

// TestDuplicateAccountNames tests the unique constraint on (name, is_active)
func (suite *AccountRepositoryTestSuite) TestDuplicateAccountNames() {
	// Given - create an active account
	account1 := &models.Account{
		Name:     "Duplicate Account",
		Type:     models.AccountTypeChecking,
		Currency: "USD",
		IsActive: true,
	}
	err := suite.db.Create(account1).Error
	assert.NoError(suite.T(), err)

	// When - try to create another ACTIVE account with same name
	// Then - should fail due to unique constraint on (name, is_active)
	account2 := &models.Account{
		Name:     "Duplicate Account",
		Type:     models.AccountTypeSavings,
		Currency: "USD",
		IsActive: true,
	}
	err = suite.db.Create(account2).Error
	assert.Error(suite.T(), err) // Should fail - duplicate active account name

	// When - deactivate the first account and create another with same name
	// Then - should succeed because the constraint is on (name, is_active)
	_ = suite.db.Model(account1).Update("is_active", false).Error
	account3 := &models.Account{
		Name:     "Duplicate Account",
		Type:     models.AccountTypeSavings,
		IsActive: true,
	}
	err = suite.db.Create(account3).Error
	assert.NoError(suite.T(), err) // Should succeed - only one active account with this name
}

// Run the test suite
func TestAccountRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AccountRepositoryTestSuite))
}
