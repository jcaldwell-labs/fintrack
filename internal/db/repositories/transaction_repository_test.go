package repositories

import (
	"testing"
	"time"

	"github.com/fintrack/fintrack/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TransactionRepositoryTestSuite struct {
	suite.Suite
	db      *gorm.DB
	repo    *TransactionRepository
	account *models.Account
}

func (suite *TransactionRepositoryTestSuite) SetupTest() {
	// Create in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		suite.T().Fatalf("failed to connect database: %v", err)
	}

	// Auto-migrate schemas
	err = db.AutoMigrate(&models.Account{}, &models.Category{}, &models.Transaction{})
	if err != nil {
		suite.T().Fatalf("failed to migrate database: %v", err)
	}

	suite.db = db
	suite.repo = NewTransactionRepository(db)

	// Create a test account
	suite.account = &models.Account{
		Name:           "Test Account",
		Type:           models.AccountTypeChecking,
		Currency:       "USD",
		InitialBalance: 1000.0,
		CurrentBalance: 1000.0,
		IsActive:       true,
	}
	_ = db.Create(suite.account).Error
}

func (suite *TransactionRepositoryTestSuite) TearDownTest() {
	if suite.db != nil {
		sqlDB, _ := suite.db.DB()
		_ = sqlDB.Close()
	}
}

func TestTransactionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionRepositoryTestSuite))
}

func (suite *TransactionRepositoryTestSuite) TestCreate() {
	tx := &models.Transaction{
		AccountID:   suite.account.ID,
		Date:        time.Now(),
		Amount:      -50.0,
		Type:        models.TransactionTypeExpense,
		Description: "Test transaction",
	}

	err := suite.repo.Create(tx)
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), tx.ID)
}

func (suite *TransactionRepositoryTestSuite) TestCreate_SetsDateIfNotProvided() {
	tx := &models.Transaction{
		AccountID:   suite.account.ID,
		Amount:      100.0,
		Type:        models.TransactionTypeIncome,
		Description: "Test income",
	}

	before := time.Now()
	err := suite.repo.Create(tx)
	after := time.Now()

	assert.NoError(suite.T(), err)
	assert.False(suite.T(), tx.Date.IsZero())
	assert.True(suite.T(), tx.Date.After(before) || tx.Date.Equal(before))
	assert.True(suite.T(), tx.Date.Before(after) || tx.Date.Equal(after))
}

func (suite *TransactionRepositoryTestSuite) TestCreate_NilTransaction() {
	err := suite.repo.Create(nil)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "cannot be nil")
}

func (suite *TransactionRepositoryTestSuite) TestCreate_MissingAccountID() {
	tx := &models.Transaction{
		Amount: 100.0,
		Type:   models.TransactionTypeIncome,
	}

	err := suite.repo.Create(tx)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "account_id is required")
}

func (suite *TransactionRepositoryTestSuite) TestCreate_InvalidType() {
	tx := &models.Transaction{
		AccountID: suite.account.ID,
		Amount:    100.0,
		Type:      "invalid",
	}

	err := suite.repo.Create(tx)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid transaction type")
}

func (suite *TransactionRepositoryTestSuite) TestGetByID() {
	tx := &models.Transaction{
		AccountID:   suite.account.ID,
		Date:        time.Now(),
		Amount:      -25.0,
		Type:        models.TransactionTypeExpense,
		Description: "Grocery shopping",
	}
	_ = suite.repo.Create(tx)

	retrieved, err := suite.repo.GetByID(tx.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), tx.ID, retrieved.ID)
	assert.Equal(suite.T(), tx.Amount, retrieved.Amount)
	assert.NotNil(suite.T(), retrieved.Account)
	assert.Equal(suite.T(), suite.account.Name, retrieved.Account.Name)
}

func (suite *TransactionRepositoryTestSuite) TestGetByID_NotFound() {
	retrieved, err := suite.repo.GetByID(9999)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), retrieved)
	assert.Contains(suite.T(), err.Error(), "not found")
}

func (suite *TransactionRepositoryTestSuite) TestList() {
	transactions := []*models.Transaction{
		{AccountID: suite.account.ID, Date: time.Now().AddDate(0, 0, -2), Amount: -50.0, Type: models.TransactionTypeExpense},
		{AccountID: suite.account.ID, Date: time.Now().AddDate(0, 0, -1), Amount: 100.0, Type: models.TransactionTypeIncome},
		{AccountID: suite.account.ID, Date: time.Now(), Amount: -25.0, Type: models.TransactionTypeExpense},
	}
	for _, tx := range transactions {
		_ = suite.repo.Create(tx)
	}

	result, err := suite.repo.List(nil, nil, nil, nil, 0)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3, len(result))
	// Should be ordered by date DESC
	assert.True(suite.T(), result[0].Date.After(result[1].Date) || result[0].Date.Equal(result[1].Date))
}

func (suite *TransactionRepositoryTestSuite) TestList_FilterByAccount() {
	// Create another account
	account2 := &models.Account{
		Name:     "Account 2",
		Type:     models.AccountTypeSavings,
		Currency: "USD",
		IsActive: true,
	}
	_ = suite.db.Create(account2).Error

	// Create transactions for both accounts
	_ = suite.repo.Create(&models.Transaction{AccountID: suite.account.ID, Amount: 100.0, Type: models.TransactionTypeIncome, Date: time.Now()})
	_ = suite.repo.Create(&models.Transaction{AccountID: account2.ID, Amount: 50.0, Type: models.TransactionTypeIncome, Date: time.Now()})

	result, err := suite.repo.List(&suite.account.ID, nil, nil, nil, 0)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(result))
	assert.Equal(suite.T(), suite.account.ID, result[0].AccountID)
}

func (suite *TransactionRepositoryTestSuite) TestList_FilterByDateRange() {
	now := time.Now()
	startDate := now.AddDate(0, 0, -5)
	endDate := now.AddDate(0, 0, -1)

	// Create transactions inside and outside the range
	_ = suite.repo.Create(&models.Transaction{AccountID: suite.account.ID, Amount: 50.0, Type: models.TransactionTypeIncome, Date: now.AddDate(0, 0, -10)})
	_ = suite.repo.Create(&models.Transaction{AccountID: suite.account.ID, Amount: 100.0, Type: models.TransactionTypeIncome, Date: now.AddDate(0, 0, -3)})
	_ = suite.repo.Create(&models.Transaction{AccountID: suite.account.ID, Amount: 25.0, Type: models.TransactionTypeIncome, Date: now})

	result, err := suite.repo.List(nil, &startDate, &endDate, nil, 0)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(result))
	assert.Equal(suite.T(), 100.0, result[0].Amount)
}

func (suite *TransactionRepositoryTestSuite) TestList_FilterByType() {
	_ = suite.repo.Create(&models.Transaction{AccountID: suite.account.ID, Amount: 100.0, Type: models.TransactionTypeIncome, Date: time.Now()})
	_ = suite.repo.Create(&models.Transaction{AccountID: suite.account.ID, Amount: -50.0, Type: models.TransactionTypeExpense, Date: time.Now()})

	txType := models.TransactionTypeExpense
	result, err := suite.repo.List(nil, nil, nil, &txType, 0)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(result))
	assert.Equal(suite.T(), models.TransactionTypeExpense, result[0].Type)
}

func (suite *TransactionRepositoryTestSuite) TestList_WithLimit() {
	for i := 0; i < 5; i++ {
		_ = suite.repo.Create(&models.Transaction{AccountID: suite.account.ID, Amount: float64(i * 10), Type: models.TransactionTypeIncome, Date: time.Now()})
	}

	result, err := suite.repo.List(nil, nil, nil, nil, 3)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3, len(result))
}

func (suite *TransactionRepositoryTestSuite) TestUpdate() {
	tx := &models.Transaction{
		AccountID:   suite.account.ID,
		Date:        time.Now(),
		Amount:      100.0,
		Type:        models.TransactionTypeIncome,
		Description: "Original",
	}
	_ = suite.repo.Create(tx)

	tx.Amount = 150.0
	tx.Description = "Updated"
	err := suite.repo.Update(tx)
	assert.NoError(suite.T(), err)

	updated, _ := suite.repo.GetByID(tx.ID)
	assert.Equal(suite.T(), 150.0, updated.Amount)
	assert.Equal(suite.T(), "Updated", updated.Description)
}

func (suite *TransactionRepositoryTestSuite) TestUpdate_NilTransaction() {
	err := suite.repo.Update(nil)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "cannot be nil")
}

func (suite *TransactionRepositoryTestSuite) TestUpdate_MissingID() {
	tx := &models.Transaction{
		AccountID: suite.account.ID,
		Amount:    100.0,
		Type:      models.TransactionTypeIncome,
	}

	err := suite.repo.Update(tx)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "ID is required")
}

func (suite *TransactionRepositoryTestSuite) TestUpdate_InvalidType() {
	tx := &models.Transaction{
		AccountID: suite.account.ID,
		Amount:    100.0,
		Type:      models.TransactionTypeIncome,
		Date:      time.Now(),
	}
	_ = suite.repo.Create(tx)

	tx.Type = "invalid"
	err := suite.repo.Update(tx)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid transaction type")
}

func (suite *TransactionRepositoryTestSuite) TestDelete() {
	tx := &models.Transaction{
		AccountID: suite.account.ID,
		Amount:    100.0,
		Type:      models.TransactionTypeIncome,
		Date:      time.Now(),
	}
	_ = suite.repo.Create(tx)

	err := suite.repo.Delete(tx.ID)
	assert.NoError(suite.T(), err)

	retrieved, err := suite.repo.GetByID(tx.ID)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), retrieved)
}

func (suite *TransactionRepositoryTestSuite) TestDelete_MissingID() {
	err := suite.repo.Delete(0)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "ID is required")
}

func (suite *TransactionRepositoryTestSuite) TestDelete_NotFound() {
	err := suite.repo.Delete(9999)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "not found")
}

func (suite *TransactionRepositoryTestSuite) TestGetAccountTotal() {
	_ = suite.repo.Create(&models.Transaction{AccountID: suite.account.ID, Amount: 100.0, Type: models.TransactionTypeIncome, Date: time.Now()})
	_ = suite.repo.Create(&models.Transaction{AccountID: suite.account.ID, Amount: -50.0, Type: models.TransactionTypeExpense, Date: time.Now()})
	_ = suite.repo.Create(&models.Transaction{AccountID: suite.account.ID, Amount: 25.0, Type: models.TransactionTypeIncome, Date: time.Now()})

	total, err := suite.repo.GetAccountTotal(suite.account.ID, nil, nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 75.0, total) // 100 - 50 + 25
}

func (suite *TransactionRepositoryTestSuite) TestGetAccountTotal_WithDateRange() {
	now := time.Now()
	startDate := now.AddDate(0, 0, -5)
	endDate := now.AddDate(0, 0, -1)

	_ = suite.repo.Create(&models.Transaction{AccountID: suite.account.ID, Amount: 50.0, Type: models.TransactionTypeIncome, Date: now.AddDate(0, 0, -10)})
	_ = suite.repo.Create(&models.Transaction{AccountID: suite.account.ID, Amount: 100.0, Type: models.TransactionTypeIncome, Date: now.AddDate(0, 0, -3)})
	_ = suite.repo.Create(&models.Transaction{AccountID: suite.account.ID, Amount: 25.0, Type: models.TransactionTypeIncome, Date: now})

	total, err := suite.repo.GetAccountTotal(suite.account.ID, &startDate, &endDate)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 100.0, total)
}

func (suite *TransactionRepositoryTestSuite) TestGetCategoryTotal() {
	// Create a category
	category := &models.Category{
		Name: "Food",
		Type: models.CategoryTypeExpense,
	}
	_ = suite.db.Create(category).Error

	_ = suite.repo.Create(&models.Transaction{AccountID: suite.account.ID, CategoryID: &category.ID, Amount: -50.0, Type: models.TransactionTypeExpense, Date: time.Now()})
	_ = suite.repo.Create(&models.Transaction{AccountID: suite.account.ID, CategoryID: &category.ID, Amount: -25.0, Type: models.TransactionTypeExpense, Date: time.Now()})

	total, err := suite.repo.GetCategoryTotal(category.ID, nil, nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), -75.0, total)
}

func (suite *TransactionRepositoryTestSuite) TestReconcile() {
	tx := &models.Transaction{
		AccountID: suite.account.ID,
		Amount:    100.0,
		Type:      models.TransactionTypeIncome,
		Date:      time.Now(),
	}
	_ = suite.repo.Create(tx)

	err := suite.repo.Reconcile(tx.ID)
	assert.NoError(suite.T(), err)

	retrieved, _ := suite.repo.GetByID(tx.ID)
	assert.True(suite.T(), retrieved.IsReconciled)
	assert.NotNil(suite.T(), retrieved.ReconciledAt)
}

func (suite *TransactionRepositoryTestSuite) TestUnreconcile() {
	tx := &models.Transaction{
		AccountID:    suite.account.ID,
		Amount:       100.0,
		Type:         models.TransactionTypeIncome,
		Date:         time.Now(),
		IsReconciled: true,
	}
	now := time.Now()
	tx.ReconciledAt = &now
	_ = suite.repo.Create(tx)

	err := suite.repo.Unreconcile(tx.ID)
	assert.NoError(suite.T(), err)

	retrieved, _ := suite.repo.GetByID(tx.ID)
	assert.False(suite.T(), retrieved.IsReconciled)
	assert.Nil(suite.T(), retrieved.ReconciledAt)
}
