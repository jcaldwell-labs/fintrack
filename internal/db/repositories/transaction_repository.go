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

// TransactionFilter defines filters for querying transactions
type TransactionFilter struct {
	AccountID  *uint
	CategoryID *uint
	Type       string
	StartDate  *time.Time
	EndDate    *time.Time
	Payee      string
	MinAmount  *float64
	MaxAmount  *float64
	Tags       []string
	Reconciled *bool
	Limit      int
	Offset     int
}

// Create creates a new transaction
func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

// GetByID retrieves a transaction by ID
func (r *TransactionRepository) GetByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Account").
		Preload("Category").
		Preload("TransferAccount").
		First(&transaction, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, err
	}
	return &transaction, nil
}

// List retrieves transactions with optional filters
func (r *TransactionRepository) List(filter *TransactionFilter) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	query := r.db.Preload("Account").Preload("Category").Preload("TransferAccount")

	query = r.applyFilters(query, filter)

	// Default ordering by date descending
	query = query.Order("date DESC, id DESC")

	// Apply pagination
	if filter != nil {
		if filter.Limit > 0 {
			query = query.Limit(filter.Limit)
		}
		if filter.Offset > 0 {
			query = query.Offset(filter.Offset)
		}
	}

	err := query.Find(&transactions).Error
	return transactions, err
}

// Count returns the total number of transactions matching the filter
func (r *TransactionRepository) Count(filter *TransactionFilter) (int64, error) {
	var count int64
	query := r.db.Model(&models.Transaction{})

	query = r.applyFilters(query, filter)

	err := query.Count(&count).Error
	return count, err
}

// applyFilters applies filter conditions to the query
func (r *TransactionRepository) applyFilters(query *gorm.DB, filter *TransactionFilter) *gorm.DB {
	if filter == nil {
		return query
	}

	if filter.AccountID != nil {
		query = query.Where("account_id = ?", *filter.AccountID)
	}

	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	if filter.StartDate != nil {
		query = query.Where("date >= ?", *filter.StartDate)
	}

	if filter.EndDate != nil {
		query = query.Where("date <= ?", *filter.EndDate)
	}

	if filter.Payee != "" {
		// Use LIKE for compatibility (works with both PostgreSQL and SQLite)
		query = query.Where("payee LIKE ?", "%"+filter.Payee+"%")
	}

	if filter.MinAmount != nil {
		query = query.Where("amount >= ?", *filter.MinAmount)
	}

	if filter.MaxAmount != nil {
		query = query.Where("amount <= ?", *filter.MaxAmount)
	}

	if len(filter.Tags) > 0 {
		// NOTE: PostgreSQL array contains operator (@>) is not compatible with SQLite
		// For production use with PostgreSQL, this will work correctly
		// For SQLite testing, tag filtering is skipped to maintain compatibility
		// TODO: Implement dialect-aware filtering or use JSON functions for SQLite
		if r.db.Dialector.Name() == "postgres" {
			query = query.Where("tags @> ?", filter.Tags)
		}
		// Skip tag filtering for SQLite and other dialects
	}

	if filter.Reconciled != nil {
		query = query.Where("is_reconciled = ?", *filter.Reconciled)
	}

	return query
}

// Update updates a transaction
func (r *TransactionRepository) Update(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}

// Delete deletes a transaction
func (r *TransactionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Transaction{}, id).Error
}

// GetTotalByAccount calculates total transaction amount for an account
func (r *TransactionRepository) GetTotalByAccount(accountID uint, startDate, endDate *time.Time) (float64, error) {
	var total float64
	query := r.db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("account_id = ?", accountID)

	if startDate != nil {
		query = query.Where("date >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("date <= ?", *endDate)
	}

	err := query.Scan(&total).Error
	return total, err
}

// GetTotalByCategory calculates total transaction amount for a category
func (r *TransactionRepository) GetTotalByCategory(categoryID uint, startDate, endDate *time.Time) (float64, error) {
	var total float64
	query := r.db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("category_id = ?", categoryID)

	if startDate != nil {
		query = query.Where("date >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("date <= ?", *endDate)
	}

	err := query.Scan(&total).Error
	return total, err
}

// GetIncomeExpenseTotals calculates total income and expenses
func (r *TransactionRepository) GetIncomeExpenseTotals(accountID *uint, startDate, endDate *time.Time) (income, expenses float64, err error) {
	// Build base query filters
	buildQuery := func() *gorm.DB {
		query := r.db.Model(&models.Transaction{})
		if accountID != nil {
			query = query.Where("account_id = ?", *accountID)
		}
		if startDate != nil {
			query = query.Where("date >= ?", *startDate)
		}
		if endDate != nil {
			query = query.Where("date <= ?", *endDate)
		}
		return query
	}

	// Get income (positive amounts)
	err = buildQuery().Where("amount > 0").Select("COALESCE(SUM(amount), 0)").Scan(&income).Error
	if err != nil {
		return 0, 0, err
	}

	// Get expenses (negative amounts, convert to positive)
	err = buildQuery().Where("amount < 0").Select("COALESCE(SUM(ABS(amount)), 0)").Scan(&expenses).Error
	if err != nil {
		return 0, 0, err
	}

	return income, expenses, nil
}

// ListByDateRange retrieves transactions within a date range
func (r *TransactionRepository) ListByDateRange(startDate, endDate time.Time, accountID *uint) ([]*models.Transaction, error) {
	filter := &TransactionFilter{
		StartDate: &startDate,
		EndDate:   &endDate,
		AccountID: accountID,
	}
	return r.List(filter)
}

// MarkReconciled marks a transaction as reconciled
func (r *TransactionRepository) MarkReconciled(id uint) error {
	now := time.Now()
	return r.db.Model(&models.Transaction{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_reconciled": true,
			"reconciled_at": now,
		}).Error
}

// UnmarkReconciled marks a transaction as not reconciled
func (r *TransactionRepository) UnmarkReconciled(id uint) error {
	return r.db.Model(&models.Transaction{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_reconciled": false,
			"reconciled_at": nil,
		}).Error
}

// GetRecentTransactions retrieves the most recent transactions
func (r *TransactionRepository) GetRecentTransactions(limit int, accountID *uint) ([]*models.Transaction, error) {
	filter := &TransactionFilter{
		Limit:     limit,
		AccountID: accountID,
	}
	return r.List(filter)
}

// SearchByPayee searches transactions by payee name
func (r *TransactionRepository) SearchByPayee(payee string) ([]*models.Transaction, error) {
	filter := &TransactionFilter{
		Payee: payee,
	}
	return r.List(filter)
}

// GetTransactionsByTags retrieves transactions with specific tags
func (r *TransactionRepository) GetTransactionsByTags(tags []string) ([]*models.Transaction, error) {
	filter := &TransactionFilter{
		Tags: tags,
	}
	return r.List(filter)
}
