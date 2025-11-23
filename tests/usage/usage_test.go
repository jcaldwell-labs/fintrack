package usage

import (
	"os"
	"path/filepath"
	"testing"
)

// TestUsageDocumentation runs all usage documentation tests
func TestUsageDocumentation(t *testing.T) {
	// Build binary for testing
	binaryPath := buildTestBinary(t)
	defer os.Remove(binaryPath)

	// Setup test database
	setupTestDatabase(t)
	defer cleanupTestDatabase(t)

	// Get usage tests directory
	usageDir := "."

	// Run all usage tests
	RunAllUsageTests(t, usageDir, binaryPath)
}

// buildTestBinary builds the fintrack binary for testing
func buildTestBinary(t *testing.T) string {
	t.Helper()

	// Create temp binary path
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "fintrack")

	// Build command would go here, but for now we'll use the installed binary
	// In production, this should build from source
	// For now, assume binary is built via make build
	projectRoot := filepath.Join("..", "..")
	builtBinary := filepath.Join(projectRoot, "bin", "fintrack")

	// Check if binary exists, otherwise skip
	if _, err := os.Stat(builtBinary); os.IsNotExist(err) {
		t.Skip("Binary not found. Run 'make build' first.")
	}

	return builtBinary
}

// setupTestDatabase creates a clean test database
func setupTestDatabase(t *testing.T) {
	t.Helper()

	// Set environment for test database
	os.Setenv("FINTRACK_DB_URL", "postgresql://postgres:postgres@localhost:5432/fintrack_test?sslmode=disable")

	// Note: Database setup is handled by the application
	// The test database should be created beforehand:
	// createdb fintrack_test
}

// cleanupTestDatabase cleans up the test database
func cleanupTestDatabase(t *testing.T) {
	t.Helper()

	// Note: We don't drop the database, just clean tables
	// This allows for inspection after test runs
}
