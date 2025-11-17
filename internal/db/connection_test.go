package db

import (
	"testing"

	"github.com/fintrack/fintrack/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGet_WhenNotInitialized(t *testing.T) {
	// Save original db
	originalDB := db
	defer func() { db = originalDB }()

	db = nil
	result := Get()
	assert.Nil(t, result)
}

func TestGet_WhenInitialized(t *testing.T) {
	// Save original db
	originalDB := db
	defer func() { db = originalDB }()

	// Create a test database
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	db = testDB
	result := Get()
	assert.NotNil(t, result)
	assert.Equal(t, testDB, result)
}

func TestClose_WhenNotInitialized(t *testing.T) {
	// Save original db
	originalDB := db
	defer func() { db = originalDB }()

	db = nil
	err := Close()
	assert.NoError(t, err)
}

func TestClose_WhenInitialized(t *testing.T) {
	// Save original db
	originalDB := db
	defer func() { db = originalDB }()

	// Create a test database
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	db = testDB
	err = Close()
	assert.NoError(t, err)
}

func TestAutoMigrate_WhenNotInitialized(t *testing.T) {
	// Save original db
	originalDB := db
	defer func() { db = originalDB }()

	db = nil
	err := AutoMigrate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database not initialized")
}

func TestAutoMigrate_WhenInitialized(t *testing.T) {
	// Save original db
	originalDB := db
	defer func() { db = originalDB }()

	// Create a test database
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	db = testDB
	err = AutoMigrate()
	assert.NoError(t, err)

	// Verify migrations worked by checking if tables exist
	assert.True(t, testDB.Migrator().HasTable(&models.Account{}))
	assert.True(t, testDB.Migrator().HasTable(&models.Category{}))
	assert.True(t, testDB.Migrator().HasTable(&models.Transaction{}))
	assert.True(t, testDB.Migrator().HasTable(&models.Budget{}))
	assert.True(t, testDB.Migrator().HasTable(&models.RecurringItem{}))
	assert.True(t, testDB.Migrator().HasTable(&models.Reminder{}))
	assert.True(t, testDB.Migrator().HasTable(&models.CashFlowProjection{}))
	assert.True(t, testDB.Migrator().HasTable(&models.ImportHistory{}))
}

func TestIsConnected_WhenNotInitialized(t *testing.T) {
	// Save original db
	originalDB := db
	defer func() { db = originalDB }()

	db = nil
	connected := IsConnected()
	assert.False(t, connected)
}

func TestIsConnected_WhenInitialized(t *testing.T) {
	// Save original db
	originalDB := db
	defer func() { db = originalDB }()

	// Create a test database
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	db = testDB
	connected := IsConnected()
	assert.True(t, connected)
}

func TestIsConnected_AfterClose(t *testing.T) {
	// Save original db
	originalDB := db
	defer func() { db = originalDB }()

	// Create a test database
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	db = testDB

	// Close the connection
	err = Close()
	assert.NoError(t, err)

	// Check if connected (should be false after closing)
	connected := IsConnected()
	assert.False(t, connected)
}
