package repositories

import (
	"errors"
	"fmt"

	"github.com/fintrack/fintrack/internal/models"
	"gorm.io/gorm"
)

// AccountRepository handles account data operations
type AccountRepository struct {
	db *gorm.DB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// Create creates a new account
func (r *AccountRepository) Create(account *models.Account) error {
	// Set initial balance as current balance if not set
	if account.CurrentBalance == 0 && account.InitialBalance != 0 {
		account.CurrentBalance = account.InitialBalance
	}

	return r.db.Create(account).Error
}

// GetByID retrieves an account by ID
func (r *AccountRepository) GetByID(id uint) (*models.Account, error) {
	var account models.Account
	err := r.db.First(&account, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("account not found")
		}
		return nil, err
	}
	return &account, nil
}

// GetByName retrieves an account by name
func (r *AccountRepository) GetByName(name string) (*models.Account, error) {
	var account models.Account
	err := r.db.Where("name = ? AND is_active = ?", name, true).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("account not found")
		}
		return nil, err
	}
	return &account, nil
}

// List retrieves all accounts with optional filters
func (r *AccountRepository) List(activeOnly bool) ([]*models.Account, error) {
	var accounts []*models.Account
	query := r.db

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	err := query.Order("created_at desc").Find(&accounts).Error
	return accounts, err
}

// Update updates an account
func (r *AccountRepository) Update(account *models.Account) error {
	return r.db.Save(account).Error
}

// Delete soft-deletes an account (sets is_active = false)
func (r *AccountRepository) Delete(id uint) error {
	return r.db.Model(&models.Account{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

// HardDelete permanently deletes an account
func (r *AccountRepository) HardDelete(id uint) error {
	return r.db.Delete(&models.Account{}, id).Error
}

// UpdateBalance updates the account balance
func (r *AccountRepository) UpdateBalance(id uint, newBalance float64) error {
	return r.db.Model(&models.Account{}).
		Where("id = ?", id).
		Update("current_balance", newBalance).Error
}

// GetBalance retrieves the current balance for an account
func (r *AccountRepository) GetBalance(id uint) (float64, error) {
	var account models.Account
	err := r.db.Select("current_balance").First(&account, id).Error
	if err != nil {
		return 0, err
	}
	return account.CurrentBalance, nil
}

// NameExists checks if an account name already exists
func (r *AccountRepository) NameExists(name string, excludeID *uint) (bool, error) {
	query := r.db.Model(&models.Account{}).Where("name = ? AND is_active = ?", name, true)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}
