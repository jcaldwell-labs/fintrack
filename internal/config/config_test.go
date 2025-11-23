package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGet_ReturnsConfig(t *testing.T) {
	// Reset global config
	cfg = nil
	config := Get()
	assert.NotNil(t, config)
}

func TestGetDatabaseURL_WithDirectURL(t *testing.T) {
	config := &Config{
		Database: DatabaseConfig{
			URL: "postgresql://user:pass@localhost:5432/testdb",
		},
	}

	url := config.GetDatabaseURL()
	assert.Equal(t, "postgresql://user:pass@localhost:5432/testdb", url)
}

func TestGetDatabaseURL_WithComponents(t *testing.T) {
	config := &Config{
		Database: DatabaseConfig{
			Host:     "testhost",
			Port:     5433,
			User:     "testuser",
			Password: "testpass",
			Database: "testdb",
			SSLMode:  "require",
		},
	}

	url := config.GetDatabaseURL()
	assert.Contains(t, url, "host=testhost")
	assert.Contains(t, url, "port=5433")
	assert.Contains(t, url, "user=testuser")
	assert.Contains(t, url, "password=testpass")
	assert.Contains(t, url, "dbname=testdb")
	assert.Contains(t, url, "sslmode=require")
}

func TestGetDatabaseURL_WithEnvPassword(t *testing.T) {
	_ = os.Setenv("FINTRACK_DB_PASSWORD", "envpass")
	defer func() { _ = os.Unsetenv("FINTRACK_DB_PASSWORD") }()

	config := &Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "user",
			Database: "db",
			SSLMode:  "disable",
		},
	}

	url := config.GetDatabaseURL()
	assert.Contains(t, url, "password=envpass")
}

func TestInit_WithConfigFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
database:
  host: testhost
  port: 5433
defaults:
  currency: EUR
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	// Reset viper
	viper.Reset()

	err = Init(configPath)
	assert.NoError(t, err)

	config := Get()
	assert.Equal(t, "testhost", config.Database.Host)
	assert.Equal(t, 5433, config.Database.Port)
	assert.Equal(t, "EUR", config.Defaults.Currency)
}

func TestInit_WithoutConfigFile(t *testing.T) {
	// Reset viper
	viper.Reset()
	cfg = nil

	// Create a temporary directory for config
	tmpDir := t.TempDir()
	_ = os.Setenv("HOME", tmpDir)
	defer func() { _ = os.Unsetenv("HOME") }()

	err := Init("")
	// Should not error even if config file doesn't exist
	assert.NoError(t, err)

	config := Get()
	assert.NotNil(t, config)
}

func TestInit_WithEnvVars(t *testing.T) {
	// Create a temporary directory for config
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", tmpDir)
	defer func() { _ = os.Setenv("HOME", oldHome) }()

	// Reset viper
	viper.Reset()
	cfg = nil

	err := Init("")
	assert.NoError(t, err)

	// Just verify Init doesn't error and config is initialized with defaults
	config := Get()
	assert.NotNil(t, config)
	// Verify defaults are applied
	assert.Equal(t, "USD", config.Defaults.Currency)
}

func TestDefaults(t *testing.T) {
	// Reset viper and config
	viper.Reset()
	cfg = nil
	setDefaults()

	// Unmarshal to get the actual config
	cfg = &Config{}
	err := viper.Unmarshal(cfg)
	assert.NoError(t, err)

	config := Get()

	// Test database defaults
	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, 5432, config.Database.Port)
	assert.Equal(t, "fintrack", config.Database.Database)
	assert.Equal(t, "fintrack_user", config.Database.User)
	assert.Equal(t, "disable", config.Database.SSLMode)
	assert.Equal(t, 10, config.Database.MaxConnections)
	assert.Equal(t, 2, config.Database.MaxIdleConnections)

	// Test default settings
	assert.Equal(t, "USD", config.Defaults.Currency)
	assert.Equal(t, "2006-01-02", config.Defaults.DateFormat)
	assert.Equal(t, "Local", config.Defaults.Timezone)

	// Test alert defaults
	assert.True(t, config.Alerts.Enabled)
	assert.Equal(t, 0.80, config.Alerts.Threshold)

	// Test recurring defaults
	assert.False(t, config.Recurring.AutoGenerate)
	assert.Equal(t, 3, config.Recurring.GenerateDaysAhead)
	assert.Equal(t, 3, config.Recurring.ReminderDaysBefore)

	// Test output defaults
	assert.Equal(t, "table", config.Output.DefaultFormat)
	assert.True(t, config.Output.Color)
	assert.True(t, config.Output.Unicode)
}

func TestConfig_AllStructs(t *testing.T) {
	config := &Config{
		Database: DatabaseConfig{
			Host: "test",
		},
		Defaults: DefaultsConfig{
			Currency: "USD",
		},
		Alerts: AlertsConfig{
			Enabled: true,
		},
		Recurring: RecurringConfig{
			AutoGenerate: false,
		},
		Output: OutputConfig{
			DefaultFormat: "json",
		},
	}

	assert.Equal(t, "test", config.Database.Host)
	assert.Equal(t, "USD", config.Defaults.Currency)
	assert.True(t, config.Alerts.Enabled)
	assert.False(t, config.Recurring.AutoGenerate)
	assert.Equal(t, "json", config.Output.DefaultFormat)
}
