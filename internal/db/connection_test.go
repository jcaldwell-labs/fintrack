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

// ============================================================================
// SetTestDB and ResetTestDB Tests (Test Helper Functions)
// ============================================================================

func TestSetTestDB_SetsDatabase(t *testing.T) {
	savedDB := db
	savedOriginalDB := originalDB
	defer func() {
		db = savedDB
		originalDB = savedOriginalDB
	}()

	db = nil
	originalDB = nil

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	SetTestDB(testDB)

	assert.Equal(t, testDB, db)
	assert.Nil(t, originalDB)
}

func TestSetTestDB_PreservesOriginalDB(t *testing.T) {
	savedDB := db
	savedOriginalDB := originalDB
	defer func() {
		db = savedDB
		originalDB = savedOriginalDB
	}()

	initialDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	db = initialDB
	originalDB = nil

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	SetTestDB(testDB)

	assert.Equal(t, testDB, db)
	assert.Equal(t, initialDB, originalDB)
}

func TestSetTestDB_DoesNotOverwriteOriginalDB(t *testing.T) {
	savedDB := db
	savedOriginalDB := originalDB
	defer func() {
		db = savedDB
		originalDB = savedOriginalDB
	}()

	firstDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	secondDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	thirdDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	db = secondDB
	originalDB = firstDB

	SetTestDB(thirdDB)

	assert.Equal(t, thirdDB, db)
	assert.Equal(t, firstDB, originalDB)
}

func TestResetTestDB_RestoresOriginalDB(t *testing.T) {
	savedDB := db
	savedOriginalDB := originalDB
	defer func() {
		db = savedDB
		originalDB = savedOriginalDB
	}()

	originalTestDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	db = testDB
	originalDB = originalTestDB

	ResetTestDB()

	assert.Equal(t, originalTestDB, db)
	assert.Nil(t, originalDB)
}

func TestResetTestDB_WhenNoOriginalDB(t *testing.T) {
	savedDB := db
	savedOriginalDB := originalDB
	defer func() {
		db = savedDB
		originalDB = savedOriginalDB
	}()

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	db = testDB
	originalDB = nil

	ResetTestDB()

	assert.Equal(t, testDB, db)
	assert.Nil(t, originalDB)
}

// ============================================================================
// Init Tests (Connection Initialization with Various Configs)
// ============================================================================

func TestInit_WithInvalidDatabaseURL(t *testing.T) {
	savedDB := db
	defer func() { db = savedDB }()

	_ = config.Init("")
	cfg := config.Get()

	originalURL := cfg.Database.URL
	defer func() { cfg.Database.URL = originalURL }()

	cfg.Database.URL = "invalid://not-a-real-database"

	err := Init()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect to database")
}

func TestInit_WithUnreachableDatabase(t *testing.T) {
	savedDB := db
	defer func() { db = savedDB }()

	_ = config.Init("")
	cfg := config.Get()

	originalURL := cfg.Database.URL
	defer func() { cfg.Database.URL = originalURL }()

	cfg.Database.URL = "postgresql://user:pass@192.0.2.1:5432/testdb?connect_timeout=1"

	err := Init()
	assert.Error(t, err)
}

// ============================================================================
// Pool Configuration Tests
// ============================================================================

func TestPoolConfiguration_MaxConnections(t *testing.T) {
	savedDB := db
	defer func() { db = savedDB }()

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	db = testDB

	sqlDB, err := db.DB()
	require.NoError(t, err)

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(2 * time.Hour)

	stats := sqlDB.Stats()
	assert.NotNil(t, stats)
	assert.Equal(t, 25, stats.MaxOpenConnections)
}

func TestPoolConfiguration_ConnectionLifetime(t *testing.T) {
	savedDB := db
	defer func() { db = savedDB }()

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	db = testDB

	sqlDB, err := db.DB()
	require.NoError(t, err)

	validDurations := []string{"1h", "30m", "2h30m", "1s"}
	for _, dur := range validDurations {
		duration, err := time.ParseDuration(dur)
		assert.NoError(t, err, "Failed to parse duration: %s", dur)
		sqlDB.SetConnMaxLifetime(duration)
	}
}

func TestPoolConfiguration_InvalidLifetimeDuration(t *testing.T) {
	invalidDurations := []string{"invalid", "abc", "1x"}
	for _, dur := range invalidDurations {
		_, err := time.ParseDuration(dur)
		assert.Error(t, err, "Expected error for invalid duration: %s", dur)
	}
}

func TestPoolConfiguration_DefaultValues(t *testing.T) {
	err := config.Init("")
	require.NoError(t, err)

	cfg := config.Get()

	assert.Equal(t, 10, cfg.Database.MaxConnections)
	assert.Equal(t, 2, cfg.Database.MaxIdleConnections)
	assert.Equal(t, "1h", cfg.Database.ConnectionMaxLifetime)
}

// ============================================================================
// Singleton Pattern Tests
// ============================================================================

func TestSingleton_GetReturnsSameInstance(t *testing.T) {
	savedDB := db
	defer func() { db = savedDB }()

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	db = testDB

	instance1 := Get()
	instance2 := Get()
	instance3 := Get()

	assert.Equal(t, instance1, instance2)
	assert.Equal(t, instance2, instance3)
	assert.Same(t, instance1, instance2)
}

func TestSingleton_SetTestDBAndReset(t *testing.T) {
	savedDB := db
	savedOriginalDB := originalDB
	defer func() {
		db = savedDB
		originalDB = savedOriginalDB
	}()

	prodDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	db = prodDB
	originalDB = nil

	assert.Equal(t, prodDB, Get())

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	SetTestDB(testDB)

	assert.Equal(t, testDB, Get())
	assert.NotEqual(t, prodDB, Get())

	ResetTestDB()

	assert.Equal(t, prodDB, Get())
}

func TestSingleton_SetTestDBDoesNotOverwriteOriginalDB(t *testing.T) {
	savedDB := db
	savedOriginalDB := originalDB
	defer func() {
		db = savedDB
		originalDB = savedOriginalDB
	}()

	firstDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	secondDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	thirdDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	db = secondDB
	originalDB = firstDB

	SetTestDB(thirdDB)

	assert.Equal(t, thirdDB, db)
	assert.Equal(t, firstDB, originalDB)
}

// ============================================================================
// Health Check Tests
// ============================================================================

func TestIsConnected_WithHealthyConnection(t *testing.T) {
	savedDB := db
	defer func() { db = savedDB }()

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	db = testDB

	assert.True(t, IsConnected())

	sqlDB, err := db.DB()
	require.NoError(t, err)
	assert.NoError(t, sqlDB.Ping())
}

func TestIsConnected_WithClosedConnection(t *testing.T) {
	savedDB := db
	defer func() { db = savedDB }()

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	db = testDB

	sqlDB, err := db.DB()
	require.NoError(t, err)
	err = sqlDB.Close()
	require.NoError(t, err)

	assert.False(t, IsConnected())
}

func TestIsConnected_MultipleChecks(t *testing.T) {
	savedDB := db
	defer func() { db = savedDB }()

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	db = testDB

	for i := 0; i < 5; i++ {
		assert.True(t, IsConnected(), "Health check %d failed", i)
	}
}

// ============================================================================
// Close Tests (Enhanced)
// ============================================================================

func TestClose_AfterOperations(t *testing.T) {
	savedDB := db
	defer func() { db = savedDB }()

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	db = testDB

	err = AutoMigrate()
	require.NoError(t, err)

	account := &models.Account{
		Name:     "Test Account",
		Type:     "checking",
		Currency: "USD",
		IsActive: true,
	}
	err = db.Create(account).Error
	require.NoError(t, err)

	err = Close()
	assert.NoError(t, err)
}

// ============================================================================
// AutoMigrate Tests (Enhanced)
// ============================================================================

func TestAutoMigrate_CreatesTables(t *testing.T) {
	savedDB := db
	defer func() { db = savedDB }()

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	db = testDB

	err = AutoMigrate()
	require.NoError(t, err)

	expectedTables := []interface{}{
		&models.Account{},
		&models.Category{},
		&models.Transaction{},
		&models.Budget{},
		&models.RecurringItem{},
		&models.Reminder{},
		&models.CashFlowProjection{},
		&models.ImportHistory{},
	}

	for _, model := range expectedTables {
		assert.True(t, testDB.Migrator().HasTable(model), "Table not created for %T", model)
	}
}

func TestAutoMigrate_Idempotent(t *testing.T) {
	savedDB := db
	defer func() { db = savedDB }()

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	db = testDB

	for i := 0; i < 3; i++ {
		err = AutoMigrate()
		assert.NoError(t, err, "AutoMigrate failed on iteration %d", i)
	}

	assert.True(t, testDB.Migrator().HasTable(&models.Account{}))
}

// ============================================================================
// Error Handling Tests
// ============================================================================

func TestErrorHandling_DBOperationAfterClose(t *testing.T) {
	savedDB := db
	defer func() { db = savedDB }()

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	db = testDB

	err = AutoMigrate()
	require.NoError(t, err)

	err = Close()
	require.NoError(t, err)

	assert.False(t, IsConnected())
}

func TestErrorHandling_GetDBInstance(t *testing.T) {
	savedDB := db
	defer func() { db = savedDB }()

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	db = testDB

	sqlDB, err := db.DB()
	assert.NoError(t, err)
	assert.NotNil(t, sqlDB)
}

// ============================================================================
// Connection State Transition Tests
// ============================================================================

func TestConnectionStateTransitions(t *testing.T) {
	savedDB := db
	defer func() { db = savedDB }()

	db = nil
	assert.False(t, IsConnected())
	assert.Nil(t, Get())

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	db = testDB

	assert.True(t, IsConnected())
	assert.NotNil(t, Get())

	err = Close()
	require.NoError(t, err)
	assert.False(t, IsConnected())

	assert.NotNil(t, Get())
}

func TestConnectionStateWithSetTestDB(t *testing.T) {
	savedDB := db
	savedOriginalDB := originalDB
	defer func() {
		db = savedDB
		originalDB = savedOriginalDB
	}()

	prodDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	db = prodDB
	originalDB = nil

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	assert.Equal(t, prodDB, Get())
	assert.True(t, IsConnected())

	SetTestDB(testDB)
	assert.Equal(t, testDB, Get())
	assert.True(t, IsConnected())

	err = Close()
	require.NoError(t, err)
	assert.False(t, IsConnected())

	ResetTestDB()
	assert.Equal(t, prodDB, Get())
	assert.True(t, IsConnected())
}

// ============================================================================
// Configuration-Based Initialization Tests
// ============================================================================

func TestInit_ConfigurationDefaults(t *testing.T) {
	err := config.Init("")
	require.NoError(t, err)

	cfg := config.Get()

	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, "fintrack", cfg.Database.Database)
	assert.Equal(t, "disable", cfg.Database.SSLMode)
	assert.Equal(t, 10, cfg.Database.MaxConnections)
	assert.Equal(t, 2, cfg.Database.MaxIdleConnections)
	assert.Equal(t, "1h", cfg.Database.ConnectionMaxLifetime)
}

func TestInit_GetDatabaseURLBuildsFromComponents(t *testing.T) {
	err := config.Init("")
	require.NoError(t, err)

	cfg := config.Get()

	originalURL := cfg.Database.URL
	originalHost := cfg.Database.Host
	originalPort := cfg.Database.Port
	originalUser := cfg.Database.User
	originalPassword := cfg.Database.Password
	originalDatabase := cfg.Database.Database
	originalSSLMode := cfg.Database.SSLMode

	defer func() {
		cfg.Database.URL = originalURL
		cfg.Database.Host = originalHost
		cfg.Database.Port = originalPort
		cfg.Database.User = originalUser
		cfg.Database.Password = originalPassword
		cfg.Database.Database = originalDatabase
		cfg.Database.SSLMode = originalSSLMode
	}()

	cfg.Database.URL = ""
	cfg.Database.Host = "testhost"
	cfg.Database.Port = 5433
	cfg.Database.User = "testuser"
	cfg.Database.Password = "testpass"
	cfg.Database.Database = "testdb"
	cfg.Database.SSLMode = "require"

	url := cfg.GetDatabaseURL()

	assert.Contains(t, url, "host=testhost")
	assert.Contains(t, url, "port=5433")
	assert.Contains(t, url, "user=testuser")
	assert.Contains(t, url, "password=testpass")
	assert.Contains(t, url, "dbname=testdb")
	assert.Contains(t, url, "sslmode=require")
}

func TestInit_GetDatabaseURLReturnsDirectURL(t *testing.T) {
	err := config.Init("")
	require.NoError(t, err)

	cfg := config.Get()

	originalURL := cfg.Database.URL
	defer func() { cfg.Database.URL = originalURL }()

	expectedURL := "postgresql://direct:url@localhost:5432/db"
	cfg.Database.URL = expectedURL

	url := cfg.GetDatabaseURL()
	assert.Equal(t, expectedURL, url)
}
