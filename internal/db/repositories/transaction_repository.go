package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/fintrack/fintrack/internal/models"
	"gorm.io/gorm"
)

// TransactionFilter contains filter options for listing transactions
type TransactionFilter struct {
	AccountID   *uint
	CategoryID  *uint
	Type        string
	DateFrom    *time.Time
	DateTo      *time.Time
	Payee       string
	IsReconciled *bool
	Limit       int
	Offset      int
}

// TransactionRepository handles transaction data operations
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create creates a new transaction and updates the account balance
func (r *TransactionRepository) Create(tx *models.Transaction) error {
	return r.db.Transaction(func(dbTx *gorm.DB) error {
		// Create the transaction
		if err := dbTx.Create(tx).Error; err != nil {
			return err
		}

		// Update account balance
		return r.updateAccountBalance(dbTx, tx.AccountID, tx.Amount)
	})
}

// GetByID retrieves a transaction by ID with related entities
func (r *TransactionRepository) GetByID(id uint) (*models.Transaction, error) {
	var tx models.Transaction
	err := r.db.Preload("Account").Preload("Category").First(&tx, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, err
	}
	return &tx, nil
}

// List retrieves transactions with optional filters
func (r *TransactionRepository) List(filter TransactionFilter) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	query := r.db.Preload("Account").Preload("Category")

	// Apply filters
	if filter.AccountID != nil {
		query = query.Where("account_id = ?", *filter.AccountID)
	}
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.DateFrom != nil {
		query = query.Where("date >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("date <= ?", *filter.DateTo)
	}
	if filter.Payee != "" {
		query = query.Where("payee LIKE ?", "%"+filter.Payee+"%")
	}
	if filter.IsReconciled != nil {
		query = query.Where("is_reconciled = ?", *filter.IsReconciled)
	}

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	err := query.Order("date desc, id desc").Find(&transactions).Error
	return transactions, err
}

// ListByAccount retrieves all transactions for a specific account
func (r *TransactionRepository) ListByAccount(accountID uint, limit int) ([]*models.Transaction, error) {
	return r.List(TransactionFilter{AccountID: &accountID, Limit: limit})
}

// ListByCategory retrieves all transactions for a specific category
func (r *TransactionRepository) ListByCategory(categoryID uint, limit int) ([]*models.Transaction, error) {
	return r.List(TransactionFilter{CategoryID: &categoryID, Limit: limit})
}

// ListByDateRange retrieves transactions within a date range
func (r *TransactionRepository) ListByDateRange(from, to time.Time, limit int) ([]*models.Transaction, error) {
	return r.List(TransactionFilter{DateFrom: &from, DateTo: &to, Limit: limit})
}

// Update updates a transaction and adjusts account balance if amount changed
func (r *TransactionRepository) Update(tx *models.Transaction) error {
	return r.db.Transaction(func(dbTx *gorm.DB) error {
		// Get the original transaction to calculate balance adjustment
		var original models.Transaction
		if err := dbTx.First(&original, tx.ID).Error; err != nil {
			return err
		}

		// Calculate the balance difference
		balanceDiff := tx.Amount - original.Amount

		// Update the transaction
		if err := dbTx.Save(tx).Error; err != nil {
			return err
		}

		// Adjust account balance if amount changed
		if balanceDiff != 0 {
			// If account changed, revert old account and update new
			if original.AccountID != tx.AccountID {
				// Revert original account
				if err := r.updateAccountBalance(dbTx, original.AccountID, -original.Amount); err != nil {
					return err
				}
				// Update new account
				return r.updateAccountBalance(dbTx, tx.AccountID, tx.Amount)
			}
			// Same account, just apply the difference
			return r.updateAccountBalance(dbTx, tx.AccountID, balanceDiff)
		}

		return nil
	})
}

// Delete deletes a transaction and adjusts the account balance
func (r *TransactionRepository) Delete(id uint) error {
	return r.db.Transaction(func(dbTx *gorm.DB) error {
		// Get the transaction first
		var tx models.Transaction
		if err := dbTx.First(&tx, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("transaction not found")
			}
			return err
		}

		// Delete the transaction
		if err := dbTx.Delete(&models.Transaction{}, id).Error; err != nil {
			return err
		}

		// Revert the account balance
		return r.updateAccountBalance(dbTx, tx.AccountID, -tx.Amount)
	})
}

// GetTotalByAccount calculates total amount for an account (optionally filtered by type)
func (r *TransactionRepository) GetTotalByAccount(accountID uint, txType string) (float64, error) {
	var total float64
	query := r.db.Model(&models.Transaction{}).Where("account_id = ?", accountID)

	if txType != "" {
		query = query.Where("type = ?", txType)
	}

	err := query.Select("COALESCE(SUM(amount), 0)").Scan(&total).Error
	return total, err
}

// GetTotalByCategory calculates total amount for a category within a date range
func (r *TransactionRepository) GetTotalByCategory(categoryID uint, from, to *time.Time) (float64, error) {
	var total float64
	query := r.db.Model(&models.Transaction{}).Where("category_id = ?", categoryID)

	if from != nil {
		query = query.Where("date >= ?", *from)
	}
	if to != nil {
		query = query.Where("date <= ?", *to)
	}

	err := query.Select("COALESCE(SUM(amount), 0)").Scan(&total).Error
	return total, err
}

// Reconcile marks a transaction as reconciled
func (r *TransactionRepository) Reconcile(id uint) error {
	now := time.Now()
	return r.db.Model(&models.Transaction{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_reconciled": true,
			"reconciled_at": now,
		}).Error
}

// Unreconcile marks a transaction as not reconciled
func (r *TransactionRepository) Unreconcile(id uint) error {
	return r.db.Model(&models.Transaction{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_reconciled": false,
			"reconciled_at": nil,
		}).Error
}

// Count returns the total number of transactions matching the filter
func (r *TransactionRepository) Count(filter TransactionFilter) (int64, error) {
	var count int64
	query := r.db.Model(&models.Transaction{})

	if filter.AccountID != nil {
		query = query.Where("account_id = ?", *filter.AccountID)
	}
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.DateFrom != nil {
		query = query.Where("date >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("date <= ?", *filter.DateTo)
	}

	err := query.Count(&count).Error
	return count, err
}

// updateAccountBalance updates an account's current balance
func (r *TransactionRepository) updateAccountBalance(db *gorm.DB, accountID uint, amount float64) error {
	return db.Model(&models.Account{}).
		Where("id = ?", accountID).
		Update("current_balance", gorm.Expr("current_balance + ?", amount)).Error
}
