package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/fintrack/fintrack/internal/config"
	"github.com/fintrack/fintrack/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db         *gorm.DB
	originalDB *gorm.DB
	initOnce   sync.Once
	initErr    error
	mu         sync.RWMutex
)

// SetTestDB sets a test database instance (for testing only)
func SetTestDB(testDB *gorm.DB) {
	mu.Lock()
	defer mu.Unlock()
	if originalDB == nil {
		originalDB = db
	}
	db = testDB
}

// ResetTestDB resets to the original database instance (for testing only)
func ResetTestDB() {
	mu.Lock()
	defer mu.Unlock()
	if originalDB != nil {
		db = originalDB
		originalDB = nil
	}
	// Reset initOnce so Init() can be called again
	initOnce = sync.Once{}
	initErr = nil
}

// Init initializes the database connection (thread-safe, runs only once)
func Init() error {
	initOnce.Do(func() {
		initErr = initDB()
	})
	return initErr
}

// initDB performs the actual database initialization
func initDB() error {
	cfg := config.Get()

	dsn := cfg.GetDatabaseURL()
	if dsn == "" {
		return fmt.Errorf(`database not configured

Setup required:
  1. Create database: createdb fintrack
  2. Set connection URL (choose one):

     Environment variable:
       export FINTRACK_DB_URL="postgresql://localhost:5432/fintrack?sslmode=disable"

     Config file (~/.config/fintrack/config.yaml):
       database:
         url: "postgresql://localhost:5432/fintrack?sslmode=disable"

For more help, see: https://github.com/jcaldwell-labs/fintrack#configuration`)
	}

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// Connect to database
	var err error
	connDB, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database for connection pool configuration
	sqlDB, err := connDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(cfg.Database.MaxConnections)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConnections)

	// Parse connection max lifetime
	if cfg.Database.ConnectionMaxLifetime != "" {
		duration, err := time.ParseDuration(cfg.Database.ConnectionMaxLifetime)
		if err == nil {
			sqlDB.SetConnMaxLifetime(duration)
		}
	}

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Only assign to global after successful connection
	mu.Lock()
	db = connDB
	mu.Unlock()

	return nil
}

// Get returns the database instance (thread-safe)
func Get() *gorm.DB {
	mu.RLock()
	defer mu.RUnlock()
	return db
}

// Close closes the database connection (thread-safe)
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// AutoMigrate runs database migrations
func AutoMigrate() error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	return db.AutoMigrate(
		&models.Account{},
		&models.Category{},
		&models.Transaction{},
		&models.Budget{},
		&models.RecurringItem{},
		&models.Reminder{},
		&models.CashFlowProjection{},
		&models.ImportHistory{},
	)
}

// IsConnected checks if database is connected (thread-safe)
func IsConnected() bool {
	mu.RLock()
	defer mu.RUnlock()

	if db == nil {
		return false
	}

	sqlDB, err := db.DB()
	if err != nil {
		return false
	}

	return sqlDB.Ping() == nil
}
