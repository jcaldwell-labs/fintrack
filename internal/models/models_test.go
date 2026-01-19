package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStringArray_Scan_Nil(t *testing.T) {
	var arr StringArray
	err := arr.Scan(nil)
	assert.NoError(t, err)
	assert.Equal(t, []string{}, []string(arr))
}

func TestStringArray_Scan_ByteSlice(t *testing.T) {
	var arr StringArray
	data := []byte(`["tag1", "tag2", "tag3"]`)
	err := arr.Scan(data)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(arr))
	assert.Equal(t, "tag1", arr[0])
	assert.Equal(t, "tag2", arr[1])
	assert.Equal(t, "tag3", arr[2])
}

func TestStringArray_Scan_String(t *testing.T) {
	var arr StringArray
	data := `["apple", "banana", "cherry"]`
	err := arr.Scan(data)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(arr))
	assert.Equal(t, "apple", arr[0])
	assert.Equal(t, "banana", arr[1])
	assert.Equal(t, "cherry", arr[2])
}

func TestStringArray_Scan_InvalidType(t *testing.T) {
	var arr StringArray
	err := arr.Scan(123)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported type")
}

func TestStringArray_Scan_InvalidJSON(t *testing.T) {
	var arr StringArray
	err := arr.Scan([]byte(`invalid json`))
	assert.Error(t, err)
}

func TestStringArray_Value_Nil(t *testing.T) {
	var arr StringArray
	value, err := arr.Value()
	assert.NoError(t, err)
	assert.Equal(t, "{}", value)
}

func TestStringArray_Value_WithData(t *testing.T) {
	arr := StringArray{"one", "two", "three"}
	value, err := arr.Value()
	assert.NoError(t, err)

	// Value should be JSON-encoded
	var decoded []string
	err = json.Unmarshal(value.([]byte), &decoded)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(decoded))
	assert.Equal(t, "one", decoded[0])
	assert.Equal(t, "two", decoded[1])
	assert.Equal(t, "three", decoded[2])
}

func TestAccount_Structure(t *testing.T) {
	now := time.Now()
	account := Account{
		ID:                  1,
		Name:                "Test Account",
		Type:                AccountTypeChecking,
		Currency:            "USD",
		InitialBalanceCents: 100000, // $1000.00
		CurrentBalanceCents: 150000, // $1500.00
		Institution:         "Test Bank",
		AccountNumberLast4:  "1234",
		IsActive:            true,
		Notes:               "Test notes",
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	assert.Equal(t, uint(1), account.ID)
	assert.Equal(t, "Test Account", account.Name)
	assert.Equal(t, AccountTypeChecking, account.Type)
	assert.Equal(t, "USD", account.Currency)
	assert.Equal(t, int64(100000), account.InitialBalanceCents)
	assert.Equal(t, int64(150000), account.CurrentBalanceCents)
	assert.Equal(t, "Test Bank", account.Institution)
	assert.Equal(t, "1234", account.AccountNumberLast4)
	assert.True(t, account.IsActive)
	assert.Equal(t, "Test notes", account.Notes)
}

func TestAccountType_Constants(t *testing.T) {
	assert.Equal(t, "checking", AccountTypeChecking)
	assert.Equal(t, "savings", AccountTypeSavings)
	assert.Equal(t, "credit", AccountTypeCredit)
	assert.Equal(t, "cash", AccountTypeCash)
	assert.Equal(t, "investment", AccountTypeInvestment)
	assert.Equal(t, "loan", AccountTypeLoan)
}

func TestCategory_Structure(t *testing.T) {
	now := time.Now()
	parentID := uint(10)
	category := Category{
		ID:        1,
		Name:      "Groceries",
		ParentID:  &parentID,
		Type:      CategoryTypeExpense,
		Color:     "#FF5733",
		Icon:      "shopping-cart",
		IsSystem:  false,
		CreatedAt: now,
	}

	assert.Equal(t, uint(1), category.ID)
	assert.Equal(t, "Groceries", category.Name)
	assert.NotNil(t, category.ParentID)
	assert.Equal(t, uint(10), *category.ParentID)
	assert.Equal(t, CategoryTypeExpense, category.Type)
	assert.Equal(t, "#FF5733", category.Color)
	assert.Equal(t, "shopping-cart", category.Icon)
	assert.False(t, category.IsSystem)
}

func TestCategoryType_Constants(t *testing.T) {
	assert.Equal(t, "income", CategoryTypeIncome)
	assert.Equal(t, "expense", CategoryTypeExpense)
	assert.Equal(t, "transfer", CategoryTypeTransfer)
}

func TestTransaction_Structure(t *testing.T) {
	now := time.Now()
	categoryID := uint(5)
	transferAccountID := uint(2)
	recurringID := uint(3)
	importID := uint(4)

	transaction := Transaction{
		ID:                1,
		AccountID:         10,
		Date:              now,
		AmountCents:       -5000, // -$50.00
		CategoryID:        &categoryID,
		Payee:             "Grocery Store",
		Description:       "Weekly groceries",
		Type:              TransactionTypeExpense,
		TransferAccountID: &transferAccountID,
		RecurringID:       &recurringID,
		Tags:              StringArray{"groceries", "food"},
		IsReconciled:      false,
		ImportID:          &importID,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	assert.Equal(t, uint(1), transaction.ID)
	assert.Equal(t, uint(10), transaction.AccountID)
	assert.Equal(t, int64(-5000), transaction.AmountCents)
	assert.NotNil(t, transaction.CategoryID)
	assert.Equal(t, uint(5), *transaction.CategoryID)
	assert.Equal(t, "Grocery Store", transaction.Payee)
	assert.Equal(t, "Weekly groceries", transaction.Description)
	assert.Equal(t, TransactionTypeExpense, transaction.Type)
	assert.Equal(t, 2, len(transaction.Tags))
	assert.Equal(t, "groceries", transaction.Tags[0])
}

func TestTransactionType_Constants(t *testing.T) {
	assert.Equal(t, "income", TransactionTypeIncome)
	assert.Equal(t, "expense", TransactionTypeExpense)
	assert.Equal(t, "transfer", TransactionTypeTransfer)
}

func TestBudget_Structure(t *testing.T) {
	now := time.Now()
	categoryID := uint(5)

	budget := Budget{
		ID:                  1,
		Name:                "Monthly Groceries",
		CategoryID:          &categoryID,
		PeriodType:          "monthly",
		PeriodStart:         now,
		PeriodEnd:           now.AddDate(0, 1, 0),
		LimitAmountCents:    50000, // $500.00
		RolloverEnabled:     true,
		RolloverAmountCents: 5000, // $50.00
		AlertThreshold:      0.80,
		IsActive:            true,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	assert.Equal(t, uint(1), budget.ID)
	assert.Equal(t, "Monthly Groceries", budget.Name)
	assert.Equal(t, int64(50000), budget.LimitAmountCents)
	assert.True(t, budget.RolloverEnabled)
	assert.Equal(t, int64(5000), budget.RolloverAmountCents)
	assert.Equal(t, 0.80, budget.AlertThreshold)
}

func TestRecurringItem_Structure(t *testing.T) {
	now := time.Now()
	categoryID := uint(5)
	dayOfMonth := 15
	dayOfWeek := 1
	endDate := now.AddDate(1, 0, 0)
	lastGenerated := now.AddDate(0, -1, 0)

	recurring := RecurringItem{
		ID:                 1,
		AccountID:          10,
		Name:               "Monthly Rent",
		AmountCents:        -150000, // -$1500.00
		CategoryID:         &categoryID,
		Description:        "Rent payment",
		Frequency:          FrequencyMonthly,
		FrequencyInterval:  1,
		DayOfMonth:         &dayOfMonth,
		DayOfWeek:          &dayOfWeek,
		StartDate:          now,
		EndDate:            &endDate,
		NextDate:           now.AddDate(0, 1, 0),
		LastGeneratedDate:  &lastGenerated,
		AutoGenerate:       true,
		ReminderDaysBefore: 3,
		IsActive:           true,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	assert.Equal(t, uint(1), recurring.ID)
	assert.Equal(t, "Monthly Rent", recurring.Name)
	assert.Equal(t, int64(-150000), recurring.AmountCents)
	assert.Equal(t, FrequencyMonthly, recurring.Frequency)
	assert.True(t, recurring.AutoGenerate)
}

func TestFrequency_Constants(t *testing.T) {
	assert.Equal(t, "daily", FrequencyDaily)
	assert.Equal(t, "weekly", FrequencyWeekly)
	assert.Equal(t, "biweekly", FrequencyBiweekly)
	assert.Equal(t, "monthly", FrequencyMonthly)
	assert.Equal(t, "quarterly", FrequencyQuarterly)
	assert.Equal(t, "annual", FrequencyAnnual)
}

func TestReminder_Structure(t *testing.T) {
	now := time.Now()
	relatedID := uint(10)
	remindTime := now.Add(2 * time.Hour)
	dismissedAt := now.Add(1 * time.Hour)

	reminder := Reminder{
		ID:          1,
		Type:        "budget",
		RelatedID:   &relatedID,
		Title:       "Budget Alert",
		Message:     "You've reached 80% of your budget",
		RemindDate:  now,
		RemindTime:  &remindTime,
		Priority:    "high",
		IsDismissed: true,
		DismissedAt: &dismissedAt,
		CreatedAt:   now,
	}

	assert.Equal(t, uint(1), reminder.ID)
	assert.Equal(t, "budget", reminder.Type)
	assert.Equal(t, "Budget Alert", reminder.Title)
	assert.Equal(t, "high", reminder.Priority)
	assert.True(t, reminder.IsDismissed)
}

func TestCashFlowProjection_Structure(t *testing.T) {
	now := time.Now()
	accountID := uint(10)

	projection := CashFlowProjection{
		ID:                     1,
		AccountID:              &accountID,
		ProjectionDate:         now.AddDate(0, 1, 0),
		ProjectedBalanceCents:  500000, // $5000.00
		ProjectedIncomeCents:   300000, // $3000.00
		ProjectedExpensesCents: 200000, // $2000.00
		ConfidenceLevel:        0.85,
		ProjectionType:         "moderate",
		GeneratedAt:            now,
	}

	assert.Equal(t, uint(1), projection.ID)
	assert.Equal(t, int64(500000), projection.ProjectedBalanceCents)
	assert.Equal(t, int64(300000), projection.ProjectedIncomeCents)
	assert.Equal(t, int64(200000), projection.ProjectedExpensesCents)
	assert.Equal(t, 0.85, projection.ConfidenceLevel)
	assert.Equal(t, "moderate", projection.ProjectionType)
}

func TestImportHistory_Structure(t *testing.T) {
	now := time.Now()
	accountID := uint(10)

	importHistory := ImportHistory{
		ID:              1,
		AccountID:       &accountID,
		Filename:        "transactions.csv",
		FileHash:        "abc123def456",
		Format:          "csv",
		ImportedAt:      now,
		RecordsTotal:    100,
		RecordsImported: 95,
		RecordsSkipped:  3,
		RecordsFailed:   2,
		ErrorLog:        "Some errors occurred",
		ImportMetadata:  `{"source": "bank"}`,
	}

	assert.Equal(t, uint(1), importHistory.ID)
	assert.Equal(t, "transactions.csv", importHistory.Filename)
	assert.Equal(t, "abc123def456", importHistory.FileHash)
	assert.Equal(t, 100, importHistory.RecordsTotal)
	assert.Equal(t, 95, importHistory.RecordsImported)
	assert.Equal(t, 3, importHistory.RecordsSkipped)
	assert.Equal(t, 2, importHistory.RecordsFailed)
}

func TestDollarsToCents(t *testing.T) {
	tests := []struct {
		dollars  float64
		expected int64
	}{
		{0.00, 0},
		{1.00, 100},
		{10.50, 1050},
		{-25.99, -2599},
		{0.01, 1},
		{0.10, 10},
		{1234.56, 123456},
		{-1234.56, -123456},
		// Edge cases with floating point
		{0.1, 10},
		{0.2, 20},
	}

	for _, tt := range tests {
		result := DollarsToCents(tt.dollars)
		assert.Equal(t, tt.expected, result, "DollarsToCents(%f)", tt.dollars)
	}
}

func TestCentsToDollars(t *testing.T) {
	tests := []struct {
		cents    int64
		expected float64
	}{
		{0, 0.00},
		{100, 1.00},
		{1050, 10.50},
		{-2599, -25.99},
		{1, 0.01},
		{10, 0.10},
		{123456, 1234.56},
		{-123456, -1234.56},
	}

	for _, tt := range tests {
		result := CentsToDollars(tt.cents)
		assert.Equal(t, tt.expected, result, "CentsToDollars(%d)", tt.cents)
	}
}

func TestMoneyPrecision(t *testing.T) {
	// Classic floating point precision test: 0.1 + 0.2 should equal 0.3
	// Using cents avoids this problem
	a := DollarsToCents(0.10)
	b := DollarsToCents(0.20)
	sum := a + b
	assert.Equal(t, int64(30), sum, "0.10 + 0.20 in cents should be 30")
	assert.Equal(t, 0.30, CentsToDollars(sum), "30 cents should convert back to 0.30")
}
