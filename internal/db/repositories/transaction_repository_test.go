package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestNewTransactionRepository tests repository creation
func TestNewTransactionRepository(t *testing.T) {
	var db *gorm.DB // nil db for structure test only
	repo := NewTransactionRepository(db)
	assert.NotNil(t, repo)
}

// TestTransactionFilter tests filter struct initialization
func TestTransactionFilter_Defaults(t *testing.T) {
	filter := TransactionFilter{}
	assert.Nil(t, filter.AccountID)
	assert.Nil(t, filter.CategoryID)
	assert.Equal(t, "", filter.Type)
	assert.Nil(t, filter.DateFrom)
	assert.Nil(t, filter.DateTo)
	assert.Equal(t, "", filter.Payee)
	assert.Nil(t, filter.IsReconciled)
	assert.Equal(t, 0, filter.Limit)
	assert.Equal(t, 0, filter.Offset)
}

// TestTransactionFilter_WithValues tests filter with values set
func TestTransactionFilter_WithValues(t *testing.T) {
	accountID := uint(1)
	categoryID := uint(2)
	isReconciled := true

	filter := TransactionFilter{
		AccountID:    &accountID,
		CategoryID:   &categoryID,
		Type:         "expense",
		Payee:        "Test Payee",
		IsReconciled: &isReconciled,
		Limit:        50,
		Offset:       10,
	}

	assert.Equal(t, uint(1), *filter.AccountID)
	assert.Equal(t, uint(2), *filter.CategoryID)
	assert.Equal(t, "expense", filter.Type)
	assert.Equal(t, "Test Payee", filter.Payee)
	assert.True(t, *filter.IsReconciled)
	assert.Equal(t, 50, filter.Limit)
	assert.Equal(t, 10, filter.Offset)
}
