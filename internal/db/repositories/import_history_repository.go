package repositories

import (
	"errors"
	"fmt"

	"github.com/fintrack/fintrack/internal/models"
	"gorm.io/gorm"
)

// ImportHistoryRepository handles import history data operations
type ImportHistoryRepository struct {
	db *gorm.DB
}

// NewImportHistoryRepository creates a new import history repository
func NewImportHistoryRepository(db *gorm.DB) *ImportHistoryRepository {
	return &ImportHistoryRepository{db: db}
}

// Create creates a new import history record
func (r *ImportHistoryRepository) Create(history *models.ImportHistory) error {
	return r.db.Create(history).Error
}

// GetByID retrieves an import history by ID
func (r *ImportHistoryRepository) GetByID(id uint) (*models.ImportHistory, error) {
	var history models.ImportHistory
	err := r.db.Preload("Account").First(&history, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("import history not found")
		}
		return nil, err
	}
	return &history, nil
}

// GetByFileHash retrieves import history by file hash
func (r *ImportHistoryRepository) GetByFileHash(fileHash string) (*models.ImportHistory, error) {
	var history models.ImportHistory
	err := r.db.Where("file_hash = ?", fileHash).First(&history).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error for this case
		}
		return nil, err
	}
	return &history, nil
}

// FileHashExists checks if a file has already been imported
func (r *ImportHistoryRepository) FileHashExists(fileHash string) (bool, error) {
	var count int64
	err := r.db.Model(&models.ImportHistory{}).Where("file_hash = ?", fileHash).Count(&count).Error
	return count > 0, err
}

// List retrieves all import history records
func (r *ImportHistoryRepository) List(limit int) ([]*models.ImportHistory, error) {
	var histories []*models.ImportHistory
	query := r.db.Preload("Account").Order("imported_at desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&histories).Error
	return histories, err
}

// ListByAccount retrieves import history for a specific account
func (r *ImportHistoryRepository) ListByAccount(accountID uint, limit int) ([]*models.ImportHistory, error) {
	var histories []*models.ImportHistory
	query := r.db.Where("account_id = ?", accountID).Order("imported_at desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&histories).Error
	return histories, err
}

// Update updates an import history record
func (r *ImportHistoryRepository) Update(history *models.ImportHistory) error {
	return r.db.Save(history).Error
}

// Delete permanently deletes an import history record
func (r *ImportHistoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.ImportHistory{}, id).Error
}
