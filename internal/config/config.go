package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Database  DatabaseConfig  `mapstructure:"database"`
	Defaults  DefaultsConfig  `mapstructure:"defaults"`
	Alerts    AlertsConfig    `mapstructure:"alerts"`
	Recurring RecurringConfig `mapstructure:"recurring"`
	Output    OutputConfig    `mapstructure:"output"`
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	URL                   string `mapstructure:"url"`
	Host                  string `mapstructure:"host"`
	Port                  int    `mapstructure:"port"`
	Database              string `mapstructure:"database"`
	User                  string `mapstructure:"user"`
	Password              string `mapstructure:"password"`
	SSLMode               string `mapstructure:"sslmode"`
	MaxConnections        int    `mapstructure:"max_connections"`
	MaxIdleConnections    int    `mapstructure:"max_idle_connections"`
	ConnectionMaxLifetime string `mapstructure:"connection_max_lifetime"`
}

// DefaultsConfig holds default settings
type DefaultsConfig struct {
	Account    string `mapstructure:"account"`
	Currency   string `mapstructure:"currency"`
	DateFormat string `mapstructure:"date_format"`
	Timezone   string `mapstructure:"timezone"`
}

// AlertsConfig holds alert settings
type AlertsConfig struct {
	Enabled   bool    `mapstructure:"enabled"`
	Threshold float64 `mapstructure:"threshold"`
}

// RecurringConfig holds recurring transaction settings
type RecurringConfig struct {
	AutoGenerate       bool `mapstructure:"auto_generate"`
	GenerateDaysAhead  int  `mapstructure:"generate_days_ahead"`
	ReminderDaysBefore int  `mapstructure:"reminder_days_before"`
}

// OutputConfig holds output formatting settings
type OutputConfig struct {
	DefaultFormat string `mapstructure:"default_format"`
	Color         bool   `mapstructure:"color"`
	Unicode       bool   `mapstructure:"unicode"`
}

var cfg *Config

// Init initializes the configuration
func Init(configFile string) error {
	// Set defaults
	setDefaults()

	// Set config file location
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		// Look for config in standard locations
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		configDir := filepath.Join(home, ".config", "fintrack")
		viper.AddConfigPath(configDir)
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// Read environment variables
	viper.SetEnvPrefix("FINTRACK")
	viper.AutomaticEnv()

	// Read config file (optional - don't error if not found)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal config
	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// Get returns the current configuration
func Get() *Config {
	if cfg == nil {
		cfg = &Config{}
		setDefaults()
	}
	return cfg
}

// setDefaults sets default configuration values
func setDefaults() {
	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.database", "fintrack")
	viper.SetDefault("database.user", "fintrack_user")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_connections", 10)
	viper.SetDefault("database.max_idle_connections", 2)
	viper.SetDefault("database.connection_max_lifetime", "1h")

	// Default settings
	viper.SetDefault("defaults.currency", "USD")
	viper.SetDefault("defaults.date_format", "2006-01-02")
	viper.SetDefault("defaults.timezone", "Local")

	// Alert defaults
	viper.SetDefault("alerts.enabled", true)
	viper.SetDefault("alerts.threshold", 0.80)

	// Recurring defaults
	viper.SetDefault("recurring.auto_generate", false)
	viper.SetDefault("recurring.generate_days_ahead", 3)
	viper.SetDefault("recurring.reminder_days_before", 3)

	// Output defaults
	viper.SetDefault("output.default_format", "table")
	viper.SetDefault("output.color", true)
	viper.SetDefault("output.unicode", true)
}

// GetDatabaseURL returns the database connection URL
func (c *Config) GetDatabaseURL() string {
	if c.Database.URL != "" {
		return c.Database.URL
	}

	// Build URL from components
	password := c.Database.Password
	if password == "" {
		password = os.Getenv("FINTRACK_DB_PASSWORD")
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		password,
		c.Database.Database,
		c.Database.SSLMode,
	)
}
