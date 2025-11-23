package repositories

import (
	"errors"
	"time"

	"github.com/fintrack/fintrack/internal/models"
	"gorm.io/gorm"
)

// TransactionRepository handles transaction data access
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create creates a new transaction
func (r *TransactionRepository) Create(tx *models.Transaction) error {
	if tx == nil {
		return errors.New("transaction cannot be nil")
	}

	// Validate account exists
	if tx.AccountID == 0 {
		return errors.New("account_id is required")
	}

	// Set date to now if not provided
	if tx.Date.IsZero() {
		tx.Date = time.Now()
	}

	// Validate type
	if tx.Type != models.TransactionTypeIncome &&
		tx.Type != models.TransactionTypeExpense &&
		tx.Type != models.TransactionTypeTransfer {
		return errors.New("invalid transaction type")
	}

	return r.db.Create(tx).Error
}

// GetByID retrieves a transaction by ID
func (r *TransactionRepository) GetByID(id uint) (*models.Transaction, error) {
	var tx models.Transaction
	err := r.db.Preload("Account").Preload("Category").Preload("TransferAccount").
		First(&tx, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return &tx, nil
}

// List retrieves transactions with optional filters
func (r *TransactionRepository) List(accountID *uint, startDate, endDate *time.Time, txType *string, limit int) ([]*models.Transaction, error) {
	var transactions []*models.Transaction

	query := r.db.Preload("Account").Preload("Category").Preload("TransferAccount").
		Order("date DESC, id DESC")

	if accountID != nil {
		query = query.Where("account_id = ?", *accountID)
	}

	if startDate != nil {
		query = query.Where("date >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("date <= ?", *endDate)
	}

	if txType != nil {
		query = query.Where("type = ?", *txType)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&transactions).Error
	return transactions, err
}

// Update updates a transaction
func (r *TransactionRepository) Update(tx *models.Transaction) error {
	if tx == nil {
		return errors.New("transaction cannot be nil")
	}

	if tx.ID == 0 {
		return errors.New("transaction ID is required")
	}

	// Validate type
	if tx.Type != models.TransactionTypeIncome &&
		tx.Type != models.TransactionTypeExpense &&
		tx.Type != models.TransactionTypeTransfer {
		return errors.New("invalid transaction type")
	}

	return r.db.Save(tx).Error
}

// Delete deletes a transaction
func (r *TransactionRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("transaction ID is required")
	}

	result := r.db.Delete(&models.Transaction{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("transaction not found")
	}

	return nil
}

// GetAccountTotal returns the sum of transactions for an account
func (r *TransactionRepository) GetAccountTotal(accountID uint, startDate, endDate *time.Time) (float64, error) {
	var total float64

	query := r.db.Model(&models.Transaction{}).
		Where("account_id = ?", accountID).
		Select("COALESCE(SUM(amount), 0)")

	if startDate != nil {
		query = query.Where("date >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("date <= ?", *endDate)
	}

	err := query.Scan(&total).Error
	return total, err
}

// GetCategoryTotal returns the sum of transactions for a category
func (r *TransactionRepository) GetCategoryTotal(categoryID uint, startDate, endDate *time.Time) (float64, error) {
	var total float64

	query := r.db.Model(&models.Transaction{}).
		Where("category_id = ?", categoryID).
		Select("COALESCE(SUM(amount), 0)")

	if startDate != nil {
		query = query.Where("date >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("date <= ?", *endDate)
	}

	err := query.Scan(&total).Error
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
