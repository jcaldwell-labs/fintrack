package repositories

import (
	"testing"

	"github.com/fintrack/fintrack/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AccountRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *AccountRepository
}

func (suite *AccountRepositoryTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)

	suite.db = db
	suite.repo = NewAccountRepository(db)

	err = db.AutoMigrate(&models.Account{})
	assert.NoError(suite.T(), err)
}

func (suite *AccountRepositoryTestSuite) SetupTest() {
	suite.db.Exec("DELETE FROM accounts")
}

func (suite *AccountRepositoryTestSuite) TestCreate() {
	account := &models.Account{
		Name:           "Test Account",
		Type:           models.AccountTypeChecking,
		Currency:       "USD",
		InitialBalance: 1000.0,
		IsActive:       true,
	}

	err := suite.repo.Create(account)
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), account.ID)
	assert.Equal(suite.T(), 1000.0, account.CurrentBalance)
}

func (suite *AccountRepositoryTestSuite) TestCreate_SetsCurrentBalanceFromInitial() {
	account := &models.Account{
		Name:           "Balance Test",
		Type:           models.AccountTypeSavings,
		Currency:       "USD",
		InitialBalance: 500.0,
		IsActive:       true,
	}

	err := suite.repo.Create(account)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 500.0, account.CurrentBalance)
}

func (suite *AccountRepositoryTestSuite) TestGetByID() {
	account := &models.Account{
		Name:     "Get By ID Test",
		Type:     models.AccountTypeChecking,
		Currency: "USD",
		IsActive: true,
	}
	_ = suite.repo.Create(account)

	retrieved, err := suite.repo.GetByID(account.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), account.ID, retrieved.ID)
	assert.Equal(suite.T(), "Get By ID Test", retrieved.Name)
}

func (suite *AccountRepositoryTestSuite) TestGetByID_NotFound() {
	retrieved, err := suite.repo.GetByID(9999)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), retrieved)
	assert.Contains(suite.T(), err.Error(), "not found")
}

func (suite *AccountRepositoryTestSuite) TestGetByName() {
	account := &models.Account{
		Name:     "Unique Name",
		Type:     models.AccountTypeSavings,
		Currency: "USD",
		IsActive: true,
	}
	_ = suite.repo.Create(account)

	retrieved, err := suite.repo.GetByName("Unique Name")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), account.ID, retrieved.ID)
	assert.Equal(suite.T(), "Unique Name", retrieved.Name)
}

func (suite *AccountRepositoryTestSuite) TestGetByName_NotFound() {
	retrieved, err := suite.repo.GetByName("Nonexistent")
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), retrieved)
	assert.Contains(suite.T(), err.Error(), "not found")
}

func (suite *AccountRepositoryTestSuite) TestGetByName_OnlyReturnsActive() {
	account := &models.Account{
		Name:     "Inactive Account",
		Type:     models.AccountTypeChecking,
		Currency: "USD",
		IsActive: false,
	}
	// Use raw SQL to bypass GORM defaults
	suite.db.Exec("INSERT INTO accounts (name, type, currency, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))",
		account.Name, account.Type, account.Currency, false)

	retrieved, err := suite.repo.GetByName("Inactive Account")
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), retrieved)
}

func (suite *AccountRepositoryTestSuite) TestList_All() {
	accounts := []*models.Account{
		{Name: "Account 1", Type: models.AccountTypeChecking, Currency: "USD", IsActive: true},
		{Name: "Account 2", Type: models.AccountTypeSavings, Currency: "USD", IsActive: true},
		{Name: "Account 3", Type: models.AccountTypeCredit, Currency: "USD", IsActive: false},
	}
	for _, acc := range accounts {
		_ = suite.repo.Create(acc)
	}

	result, err := suite.repo.List(false)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3, len(result))
}

func (suite *AccountRepositoryTestSuite) TestList_ActiveOnly() {
	// Create active accounts
	activeAccounts := []*models.Account{
		{Name: "Active 1", Type: models.AccountTypeChecking, Currency: "USD", IsActive: true},
		{Name: "Active 2", Type: models.AccountTypeSavings, Currency: "USD", IsActive: true},
	}
	for _, acc := range activeAccounts {
		_ = suite.repo.Create(acc)
	}

	// Create inactive account using raw SQL to bypass GORM defaults
	suite.db.Exec("INSERT INTO accounts (name, type, currency, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))",
		"Inactive", models.AccountTypeCredit, "USD", false)

	result, err := suite.repo.List(true)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(result))
}

func (suite *AccountRepositoryTestSuite) TestUpdate() {
	account := &models.Account{
		Name:           "Original Name",
		Type:           models.AccountTypeChecking,
		Currency:       "USD",
		CurrentBalance: 100.0,
		IsActive:       true,
	}
	_ = suite.repo.Create(account)

	account.Name = "Updated Name"
	account.CurrentBalance = 200.0
	err := suite.repo.Update(account)
	assert.NoError(suite.T(), err)

	updated, _ := suite.repo.GetByID(account.ID)
	assert.Equal(suite.T(), "Updated Name", updated.Name)
	assert.Equal(suite.T(), 200.0, updated.CurrentBalance)
}

func (suite *AccountRepositoryTestSuite) TestDelete() {
	account := &models.Account{
		Name:     "To Delete",
		Type:     models.AccountTypeChecking,
		Currency: "USD",
		IsActive: true,
	}
	_ = suite.repo.Create(account)

	err := suite.repo.Delete(account.ID)
	assert.NoError(suite.T(), err)

	var deleted models.Account
	suite.db.First(&deleted, account.ID)
	assert.False(suite.T(), deleted.IsActive)
}

func (suite *AccountRepositoryTestSuite) TestHardDelete() {
	account := &models.Account{
		Name:     "To Hard Delete",
		Type:     models.AccountTypeChecking,
		Currency: "USD",
		IsActive: true,
	}
	_ = suite.repo.Create(account)

	err := suite.repo.HardDelete(account.ID)
	assert.NoError(suite.T(), err)

	retrieved, err := suite.repo.GetByID(account.ID)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), retrieved)
}

func (suite *AccountRepositoryTestSuite) TestUpdateBalance() {
	account := &models.Account{
		Name:           "Balance Update Test",
		Type:           models.AccountTypeChecking,
		Currency:       "USD",
		CurrentBalance: 100.0,
		IsActive:       true,
	}
	_ = suite.repo.Create(account)

	err := suite.repo.UpdateBalance(account.ID, 250.0)
	assert.NoError(suite.T(), err)

	updated, _ := suite.repo.GetByID(account.ID)
	assert.Equal(suite.T(), 250.0, updated.CurrentBalance)
}

func (suite *AccountRepositoryTestSuite) TestGetBalance() {
	account := &models.Account{
		Name:           "Get Balance Test",
		Type:           models.AccountTypeSavings,
		Currency:       "USD",
		CurrentBalance: 350.0,
		IsActive:       true,
	}
	_ = suite.repo.Create(account)

	balance, err := suite.repo.GetBalance(account.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 350.0, balance)
}

func (suite *AccountRepositoryTestSuite) TestGetBalance_NotFound() {
	balance, err := suite.repo.GetBalance(9999)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 0.0, balance)
}

func (suite *AccountRepositoryTestSuite) TestNameExists_True() {
	account := &models.Account{
		Name:     "Existing Name",
		Type:     models.AccountTypeChecking,
		Currency: "USD",
		IsActive: true,
	}
	_ = suite.repo.Create(account)

	exists, err := suite.repo.NameExists("Existing Name", nil)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)
}

func (suite *AccountRepositoryTestSuite) TestNameExists_False() {
	exists, err := suite.repo.NameExists("Nonexistent Name", nil)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

func (suite *AccountRepositoryTestSuite) TestNameExists_WithExclude() {
	account := &models.Account{
		Name:     "Test Name",
		Type:     models.AccountTypeChecking,
		Currency: "USD",
		IsActive: true,
	}
	_ = suite.repo.Create(account)

	// Should return false when excluding the ID of the account with that name
	exists, err := suite.repo.NameExists("Test Name", &account.ID)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

func (suite *AccountRepositoryTestSuite) TestNameExists_OnlyActiveAccounts() {
	// Use raw SQL to bypass GORM defaults
	suite.db.Exec("INSERT INTO accounts (name, type, currency, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))",
		"Inactive Name", models.AccountTypeChecking, "USD", false)

	exists, err := suite.repo.NameExists("Inactive Name", nil)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

func TestAccountRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AccountRepositoryTestSuite))
}
