package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/fintrack/fintrack/internal/models"
	"gorm.io/gorm"
)

// TransactionRepository handles transaction data operations
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create creates a new transaction
func (r *TransactionRepository) Create(tx *models.Transaction) error {
	return r.db.Create(tx).Error
}

// CreateBatch creates multiple transactions in a single operation
func (r *TransactionRepository) CreateBatch(txs []*models.Transaction, batchSize int) error {
	if batchSize <= 0 {
		batchSize = 100
	}
	return r.db.CreateInBatches(txs, batchSize).Error
}

// GetByID retrieves a transaction by ID
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
func (r *TransactionRepository) List(accountID *uint, limit int, offset int) ([]*models.Transaction, error) {
	var txs []*models.Transaction
	query := r.db.Preload("Account").Preload("Category").Order("date desc")

	if accountID != nil {
		query = query.Where("account_id = ?", *accountID)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&txs).Error
	return txs, err
}

// Update updates a transaction
func (r *TransactionRepository) Update(tx *models.Transaction) error {
	return r.db.Save(tx).Error
}

// Delete permanently deletes a transaction
func (r *TransactionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Transaction{}, id).Error
}

// DuplicateCheck holds the fields used for duplicate detection
type DuplicateCheck struct {
	Date        time.Time
	Amount      float64
	Description string
}

// FindDuplicate checks if a transaction with the same date, amount, and description exists
func (r *TransactionRepository) FindDuplicate(accountID uint, check DuplicateCheck) (*models.Transaction, error) {
	var tx models.Transaction

	// Normalize date to day-level precision
	dateStart := time.Date(check.Date.Year(), check.Date.Month(), check.Date.Day(), 0, 0, 0, 0, check.Date.Location())
	dateEnd := dateStart.Add(24 * time.Hour)

	err := r.db.Where(
		"account_id = ? AND date >= ? AND date < ? AND amount = ? AND description = ?",
		accountID, dateStart, dateEnd, check.Amount, check.Description,
	).First(&tx).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No duplicate found
		}
		return nil, err
	}
	return &tx, nil
}

// FindDuplicates checks multiple transactions for duplicates in batch
func (r *TransactionRepository) FindDuplicates(accountID uint, checks []DuplicateCheck) (map[int]bool, error) {
	result := make(map[int]bool)

	for i, check := range checks {
		dup, err := r.FindDuplicate(accountID, check)
		if err != nil {
			return nil, err
		}
		result[i] = dup != nil
	}

	return result, nil
}

// CountByImportID counts transactions linked to a specific import
func (r *TransactionRepository) CountByImportID(importID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Transaction{}).Where("import_id = ?", importID).Count(&count).Error
	return count, err
}

// GetSummaryByAccount returns transaction summary for an account
func (r *TransactionRepository) GetSummaryByAccount(accountID uint) (int64, float64, error) {
	var count int64
	var sum float64

	err := r.db.Model(&models.Transaction{}).
		Where("account_id = ?", accountID).
		Count(&count).Error
	if err != nil {
		return 0, 0, err
	}

	err = r.db.Model(&models.Transaction{}).
		Where("account_id = ?", accountID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&sum).Error

	return count, sum, err
}
