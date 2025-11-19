package repositories

import (
	"errors"
	"fmt"

	"github.com/fintrack/fintrack/internal/models"
	"gorm.io/gorm"
)

// CategoryRepository handles category data operations
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create creates a new category
func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// GetByID retrieves a category by ID
func (r *CategoryRepository) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.Preload("Parent").First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category not found")
		}
		return nil, err
	}
	return &category, nil
}

// GetByName retrieves a category by name and type
func (r *CategoryRepository) GetByName(name string, categoryType string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("name = ? AND type = ?", name, categoryType).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category not found")
		}
		return nil, err
	}
	return &category, nil
}

// List retrieves all categories with optional type filter
func (r *CategoryRepository) List(categoryType string) ([]*models.Category, error) {
	var categories []*models.Category
	query := r.db.Preload("Parent")

	if categoryType != "" {
		query = query.Where("type = ?", categoryType)
	}

	err := query.Order("type, name").Find(&categories).Error
	return categories, err
}

// ListByType retrieves all categories of a specific type
func (r *CategoryRepository) ListByType(categoryType string) ([]*models.Category, error) {
	return r.List(categoryType)
}

// ListTopLevel retrieves all top-level categories (no parent)
func (r *CategoryRepository) ListTopLevel(categoryType string) ([]*models.Category, error) {
	var categories []*models.Category
	query := r.db.Where("parent_id IS NULL")

	if categoryType != "" {
		query = query.Where("type = ?", categoryType)
	}

	err := query.Order("type, name").Find(&categories).Error
	return categories, err
}

// ListSubcategories retrieves all subcategories of a parent category
func (r *CategoryRepository) ListSubcategories(parentID uint) ([]*models.Category, error) {
	var categories []*models.Category
	err := r.db.Where("parent_id = ?", parentID).Order("name").Find(&categories).Error
	return categories, err
}

// Update updates a category
func (r *CategoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

// Delete deletes a category (only if not a system category)
func (r *CategoryRepository) Delete(id uint) error {
	// Check if it's a system category
	var category models.Category
	if err := r.db.First(&category, id).Error; err != nil {
		return err
	}

	if category.IsSystem {
		return fmt.Errorf("cannot delete system category")
	}

	return r.db.Delete(&models.Category{}, id).Error
}

// NameExists checks if a category name already exists for the given type
func (r *CategoryRepository) NameExists(name string, categoryType string, excludeID *uint) (bool, error) {
	query := r.db.Model(&models.Category{}).Where("name = ? AND type = ?", name, categoryType)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

// GetSystemCategories retrieves all system categories
func (r *CategoryRepository) GetSystemCategories() ([]*models.Category, error) {
	var categories []*models.Category
	err := r.db.Where("is_system = ?", true).Order("type, name").Find(&categories).Error
	return categories, err
}
