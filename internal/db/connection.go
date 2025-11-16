package db

import (
	"fmt"
	"time"

	"github.com/fintrack/fintrack/internal/config"
	"github.com/fintrack/fintrack/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// Init initializes the database connection
func Init() error {
	cfg := config.Get()

	dsn := cfg.GetDatabaseURL()
	if dsn == "" {
		return fmt.Errorf("database URL not configured")
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
	db, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database for connection pool configuration
	sqlDB, err := db.DB()
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

	return nil
}

// Get returns the database instance
func Get() *gorm.DB {
	return db
}

// Close closes the database connection
func Close() error {
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

// IsConnected checks if database is connected
func IsConnected() bool {
	if db == nil {
		return false
	}

	sqlDB, err := db.DB()
	if err != nil {
		return false
	}

	return sqlDB.Ping() == nil
}
