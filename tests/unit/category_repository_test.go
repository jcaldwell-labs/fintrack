package unit

import (
	"testing"

	"github.com/fintrack/fintrack/internal/db/repositories"
	"github.com/fintrack/fintrack/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// CategoryRepositoryTestSuite is the test suite for category repository
type CategoryRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *repositories.CategoryRepository
}

// SetupSuite runs once before all tests
func (suite *CategoryRepositoryTestSuite) SetupSuite() {
	// Use in-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)

	suite.db = db
	suite.repo = repositories.NewCategoryRepository(db)

	// Run migrations
	err = db.AutoMigrate(&models.Category{})
	assert.NoError(suite.T(), err)
}

// SetupTest runs before each test
func (suite *CategoryRepositoryTestSuite) SetupTest() {
	// Clean database before each test
	_ = suite.db.Migrator().DropTable(&models.Category{})
	_ = suite.db.AutoMigrate(&models.Category{})
}

// TestCreateCategory tests category creation
func (suite *CategoryRepositoryTestSuite) TestCreateCategory() {
	// Given
	category := &models.Category{
		Name:     "Groceries",
		Type:     models.CategoryTypeExpense,
		Color:    "#FF5733",
		Icon:     "🛒",
		IsSystem: false,
	}

	// When
	err := suite.repo.Create(category)

	// Then
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), category.ID)
	assert.Equal(suite.T(), "Groceries", category.Name)
	assert.Equal(suite.T(), models.CategoryTypeExpense, category.Type)
}

// TestGetCategoryByID tests retrieving category by ID
func (suite *CategoryRepositoryTestSuite) TestGetCategoryByID() {
	// Given
	category := &models.Category{
		Name:     "Salary",
		Type:     models.CategoryTypeIncome,
		IsSystem: true,
	}
	suite.repo.Create(category)

	// When
	retrieved, err := suite.repo.GetByID(category.ID)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), category.ID, retrieved.ID)
	assert.Equal(suite.T(), "Salary", retrieved.Name)
	assert.Equal(suite.T(), models.CategoryTypeIncome, retrieved.Type)
	assert.True(suite.T(), retrieved.IsSystem)
}

// TestGetCategoryByName tests retrieving category by name and type
func (suite *CategoryRepositoryTestSuite) TestGetCategoryByName() {
	// Given
	category := &models.Category{
		Name: "Groceries",
		Type: models.CategoryTypeExpense,
	}
	suite.repo.Create(category)

	// When
	retrieved, err := suite.repo.GetByName("Groceries", models.CategoryTypeExpense)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), category.ID, retrieved.ID)
	assert.Equal(suite.T(), "Groceries", retrieved.Name)
}

// TestGetCategoryByNameNotFound tests getting non-existent category
func (suite *CategoryRepositoryTestSuite) TestGetCategoryByNameNotFound() {
	// When
	_, err := suite.repo.GetByName("NonExistent", models.CategoryTypeExpense)

	// Then
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "not found")
}

// TestListCategories tests listing all categories
func (suite *CategoryRepositoryTestSuite) TestListCategories() {
	// Given
	categories := []models.Category{
		{Name: "Groceries", Type: models.CategoryTypeExpense},
		{Name: "Salary", Type: models.CategoryTypeIncome},
		{Name: "Transportation", Type: models.CategoryTypeExpense},
	}

	for i := range categories {
		suite.repo.Create(&categories[i])
	}

	// When - get all categories
	allCategories, err := suite.repo.List("")

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3, len(allCategories))
}

// TestListCategoriesByType tests listing categories by type
func (suite *CategoryRepositoryTestSuite) TestListCategoriesByType() {
	// Given
	categories := []models.Category{
		{Name: "Groceries", Type: models.CategoryTypeExpense},
		{Name: "Salary", Type: models.CategoryTypeIncome},
		{Name: "Transportation", Type: models.CategoryTypeExpense},
	}

	for i := range categories {
		suite.repo.Create(&categories[i])
	}

	// When - get only expense categories
	expenseCategories, err := suite.repo.ListByType(models.CategoryTypeExpense)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(expenseCategories))

	// When - get only income categories
	incomeCategories, err := suite.repo.ListByType(models.CategoryTypeIncome)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(incomeCategories))
}

// TestListTopLevelCategories tests listing top-level categories
func (suite *CategoryRepositoryTestSuite) TestListTopLevelCategories() {
	// Given - create parent category
	parent := &models.Category{
		Name: "Food & Dining",
		Type: models.CategoryTypeExpense,
	}
	suite.repo.Create(parent)

	// Create subcategory
	subcategory := &models.Category{
		Name:     "Groceries",
		Type:     models.CategoryTypeExpense,
		ParentID: &parent.ID,
	}
	suite.repo.Create(subcategory)

	// When
	topLevel, err := suite.repo.ListTopLevel("")

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(topLevel))
	assert.Equal(suite.T(), "Food & Dining", topLevel[0].Name)
}

// TestListSubcategories tests listing subcategories
func (suite *CategoryRepositoryTestSuite) TestListSubcategories() {
	// Given - create parent category
	parent := &models.Category{
		Name: "Food & Dining",
		Type: models.CategoryTypeExpense,
	}
	suite.repo.Create(parent)

	// Create subcategories
	subcategories := []models.Category{
		{Name: "Groceries", Type: models.CategoryTypeExpense, ParentID: &parent.ID},
		{Name: "Restaurants", Type: models.CategoryTypeExpense, ParentID: &parent.ID},
	}

	for i := range subcategories {
		suite.repo.Create(&subcategories[i])
	}

	// When
	subs, err := suite.repo.ListSubcategories(parent.ID)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(subs))
}

// TestUpdateCategory tests category updates
func (suite *CategoryRepositoryTestSuite) TestUpdateCategory() {
	// Given
	category := &models.Category{
		Name:  "Old Name",
		Type:  models.CategoryTypeExpense,
		Color: "#000000",
	}
	suite.repo.Create(category)

	// When
	category.Name = "New Name"
	category.Color = "#FF5733"
	err := suite.repo.Update(category)

	// Then
	assert.NoError(suite.T(), err)

	var updated models.Category
	suite.db.First(&updated, category.ID)
	assert.Equal(suite.T(), "New Name", updated.Name)
	assert.Equal(suite.T(), "#FF5733", updated.Color)
}

// TestDeleteCategory tests category deletion
func (suite *CategoryRepositoryTestSuite) TestDeleteCategory() {
	// Given - create non-system category
	category := &models.Category{
		Name:     "Test Category",
		Type:     models.CategoryTypeExpense,
		IsSystem: false,
	}
	suite.repo.Create(category)

	// When
	err := suite.repo.Delete(category.ID)

	// Then
	assert.NoError(suite.T(), err)

	// Verify deletion
	_, err = suite.repo.GetByID(category.ID)
	assert.Error(suite.T(), err)
}

// TestDeleteSystemCategory tests that system categories cannot be deleted
func (suite *CategoryRepositoryTestSuite) TestDeleteSystemCategory() {
	// Given - create system category
	category := &models.Category{
		Name:     "System Category",
		Type:     models.CategoryTypeExpense,
		IsSystem: true,
	}
	suite.repo.Create(category)

	// When
	err := suite.repo.Delete(category.ID)

	// Then
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "system category")
}

// TestNameExists tests checking if category name exists
func (suite *CategoryRepositoryTestSuite) TestNameExists() {
	// Given
	category := &models.Category{
		Name: "Groceries",
		Type: models.CategoryTypeExpense,
	}
	suite.repo.Create(category)

	// When - check existing name
	exists, err := suite.repo.NameExists("Groceries", models.CategoryTypeExpense, nil)

	// Then
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)

	// When - check non-existing name
	exists, err = suite.repo.NameExists("NonExistent", models.CategoryTypeExpense, nil)

	// Then
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)

	// When - check with exclusion
	exists, err = suite.repo.NameExists("Groceries", models.CategoryTypeExpense, &category.ID)

	// Then
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists) // Should return false when excluding itself
}

// TestNameExistsWithDifferentType tests that same name can exist for different types
func (suite *CategoryRepositoryTestSuite) TestNameExistsWithDifferentType() {
	// Given - create "Transfer" as expense
	category := &models.Category{
		Name: "Transfer",
		Type: models.CategoryTypeExpense,
	}
	suite.repo.Create(category)

	// When - check if "Transfer" exists as transfer type
	exists, err := suite.repo.NameExists("Transfer", models.CategoryTypeTransfer, nil)

	// Then - should not exist because it's a different type
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

// TestGetSystemCategories tests retrieving system categories
func (suite *CategoryRepositoryTestSuite) TestGetSystemCategories() {
	// Given
	categories := []models.Category{
		{Name: "Salary", Type: models.CategoryTypeIncome, IsSystem: true},
		{Name: "Groceries", Type: models.CategoryTypeExpense, IsSystem: false},
		{Name: "Housing", Type: models.CategoryTypeExpense, IsSystem: true},
	}

	for i := range categories {
		suite.repo.Create(&categories[i])
	}

	// When
	systemCategories, err := suite.repo.GetSystemCategories()

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(systemCategories))
	for _, cat := range systemCategories {
		assert.True(suite.T(), cat.IsSystem)
	}
}

// Run the test suite
func TestCategoryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryRepositoryTestSuite))
}
