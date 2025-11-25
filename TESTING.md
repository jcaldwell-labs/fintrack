# Testing Strategy for FinTrack

This document outlines the testing approach, current status, and roadmap for quality assurance.

## 🚨 Current Status: Critical Gaps

### Test Coverage: 0%
```
Package                                              Coverage
─────────────────────────────────────────────────────────────
cmd/fintrack                                         0%
internal/commands                                    0%
internal/config                                      0%
internal/db                                          0%
internal/db/repositories                             0%
internal/models                                      0%
internal/output                                      0%
tests/unit (1 file, 2/7 tests failing)               N/A
tests/integration (0 files)                          N/A
```

### Existing Tests
- ✅ `tests/unit/account_repository_test.go` (7 tests)
  - ✅ PASS: TestCreateAccount
  - ✅ PASS: TestGetAccountByID
  - ❌ FAIL: TestListAccounts (expects 2, gets 3)
  - ✅ PASS: TestUpdateAccount
  - ✅ PASS: TestDeleteAccount
  - ✅ PASS: TestAccountTypes
  - ❌ FAIL: TestDuplicateAccountNames (UNIQUE constraint issue)

### Missing Test Coverage
- ❌ No command layer tests (`internal/commands/`)
- ❌ No repository method tests (only basic CRUD)
- ❌ No config loading tests
- ❌ No output formatter tests
- ❌ No database connection tests
- ❌ No integration tests
- ❌ No end-to-end CLI tests
- ❌ No CI/CD automation
- ❌ No regression test suite

## 🎯 Testing Pyramid Strategy

```
                    ┌──────────────┐
                    │     E2E      │  10% - Full CLI workflows
                    │   (Manual)   │
                    └──────────────┘
                 ┌──────────────────┐
                 │  Usage Tests     │  15% - Executable docs
                 │ (Living Docs)    │
                 └──────────────────┘
                   ┌────────────────┐
                   │  Integration   │  20% - DB + Repos + Commands
                   │     Tests      │
                   └────────────────┘
                ┌──────────────────────┐
                │    Unit Tests        │  55% - Business logic, utils
                │  (Fast, Isolated)    │
                └──────────────────────┘
```

### Target Coverage Goals
- **Phase 1 (Current MVP):** 60% overall coverage
  - Unit tests: 70% coverage
  - Integration tests: 40% coverage
  - E2E: Manual smoke tests

- **Phase 2 (Production Ready):** 80% overall coverage
  - Unit tests: 85% coverage
  - Integration tests: 70% coverage
  - E2E: Automated critical paths

## 📋 Testing Types

### 1. Unit Tests (70% of tests)

**Location:** Co-located with source code (same package)
```
internal/
├── commands/
│   ├── account.go
│   └── account_test.go       ← Unit tests here
├── db/repositories/
│   ├── account_repository.go
│   └── account_repository_test.go
└── output/
    ├── output.go
    └── output_test.go
```

**Characteristics:**
- Fast (< 1ms per test)
- Isolated (no external dependencies)
- Use mocks/fakes for DB, config, etc.
- Test single functions/methods

**Example:**
```go
// internal/output/output_test.go
func TestFormatCurrency(t *testing.T) {
    tests := []struct {
        name     string
        amount   float64
        currency string
        want     string
    }{
        {"positive USD", 1234.56, "USD", "$1,234.56"},
        {"negative USD", -1234.56, "USD", "-$1,234.56"},
        {"zero", 0, "USD", "$0.00"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := FormatCurrency(tt.amount, tt.currency)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

### 2. Integration Tests (20% of tests)

**Location:** `tests/integration/`
```
tests/integration/
├── account_workflow_test.go   ← Account CRUD workflow
├── config_loading_test.go     ← Config + DB integration
└── csv_import_test.go         ← Import workflow
```

**Characteristics:**
- Slower (10-100ms per test)
- Uses real PostgreSQL (test database)
- Tests component interactions
- Tests repository + database

**Example:**
```go
// tests/integration/account_workflow_test.go
func TestAccountCRUDWorkflow(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    repo := repositories.NewAccountRepository(db)

    // Create
    account := &models.Account{
        Name: "Test Account",
        Type: models.AccountTypeChecking,
    }
    err := repo.Create(account)
    require.NoError(t, err)

    // Read
    retrieved, err := repo.GetByID(account.ID)
    require.NoError(t, err)
    assert.Equal(t, "Test Account", retrieved.Name)

    // Update
    retrieved.Name = "Updated Name"
    err = repo.Update(retrieved)
    require.NoError(t, err)

    // Delete
    err = repo.Delete(account.ID)
    require.NoError(t, err)
}
```

### 3. End-to-End Tests (10% of tests)

**Location:** `tests/e2e/`
```
tests/e2e/
├── cli_smoke_test.go          ← Basic CLI operations
└── full_workflow_test.sh      ← Complete user workflow
```

**Characteristics:**
- Slowest (100ms - 1s per test)
- Tests actual CLI binary
- Uses real database
- Simulates user interactions

**Example:**
```bash
# tests/e2e/account_workflow.sh
#!/bin/bash
set -e

# Setup
export FINTRACK_DB_URL="postgresql://test:test@localhost:5432/fintrack_test"

# Create account
./bin/fintrack account add "Test Checking" --type checking --balance 1000
assert_success

# List accounts
output=$(./bin/fintrack account list --json)
assert_contains "$output" "Test Checking"

# Update account
./bin/fintrack account update 1 --name "Updated Checking"
assert_success

# Verify update
output=$(./bin/fintrack account show 1)
assert_contains "$output" "Updated Checking"

# Cleanup
./bin/fintrack account close 1
```

### 4. Usage Documentation Tests (15% of tests)

**Location:** `tests/usage/`
```
tests/usage/
├── runner.go                   ← Test harness (parser & executor)
├── usage_test.go               ← Go test runner
└── 01-account-management.md    ← Executable documentation
```

**Characteristics:**
- Medium speed (50-200ms per test)
- Tests actual CLI binary against real database
- Markdown format - human-readable and VCS-friendly
- Self-documenting - serves as user documentation
- Auto-updates with actual results
- Wildcard support for dynamic values (IDs, dates, etc.)

**What Makes Them Special:**
Usage tests are **living documentation** - they're simultaneously:
1. **User documentation** showing real command examples
2. **Automated tests** validating behavior
3. **Regression prevention** catching breaking changes
4. **Onboarding material** for new contributors

**Example:**
```markdown
## Test: Create a checking account
**Purpose:** Verify users can create a basic checking account

### Setup
```bash
# Clean slate
fintrack account delete "Test Checking" 2>/dev/null || true
```

### Execute
```bash
fintrack account create "Test Checking" --type checking --balance 1000.00
```

### Expected Output
```
Account created successfully
ID: <number>
Name: Test Checking
Type: checking
Balance: <money>
```

### Actual Output (auto-updated)
```
Account created successfully
ID: 42
Name: Test Checking
Type: checking
Balance: $1,000.00
```

✅ PASS (last run: 2025-11-23)
```

**Wildcard Patterns:**
- `<any>` - Matches any value
- `<number>` - Matches integers (42, 1, 999)
- `<date>` - Matches YYYY-MM-DD format
- `<uuid>` - Matches UUID format
- `<money>` - Matches currency ($1,234.56)

**Running Usage Tests:**
```bash
# Run all usage tests
make test-usage

# Run and update markdown with results
make test-usage-update

# Run directly with Go
go test -v ./tests/usage/

# Run specific test file
go test -v ./tests/usage/ -run TestUsageDocumentation
```

**When to Add Usage Tests:**
1. Adding new CLI commands
2. Changing output formats or error messages
3. Fixing user-reported issues
4. Documenting critical workflows
5. Regression prevention for stable features

**Benefits:**
- **Catches real-world issues** - Tests actual binary, not mocks
- **Always up-to-date** - Docs update automatically with test runs
- **User-focused** - Tests what users actually see
- **Prevents regressions** - Breaking changes caught immediately
- **Better onboarding** - New devs see working examples

## 🔧 Testing Tools & Frameworks

### Current Stack
- **Test Framework:** Go standard `testing` package
- **Assertions:** `github.com/stretchr/testify/assert`
- **Test Suites:** `github.com/stretchr/testify/suite`
- **Test DB:** SQLite (in-memory) for unit tests
- **Coverage:** `go test -cover`

### Recommended Additions
- **Mocking:** `github.com/stretchr/testify/mock` or `github.com/golang/mock`
- **Database:** Docker PostgreSQL for integration tests
- **E2E:** `github.com/Netflix/go-expect` for CLI testing
- **Coverage Reports:** `gocov` + `gocov-html`
- **Parallel Testing:** `t.Parallel()` for faster test runs

## 📊 Coverage Targets by Package

| Package | Current | Target Phase 1 | Target Phase 2 |
|---------|---------|----------------|----------------|
| `cmd/fintrack` | 0% | 40% | 60% |
| `internal/commands` | 0% | 70% | 85% |
| `internal/db/repositories` | 0% | 80% | 90% |
| `internal/config` | 0% | 70% | 80% |
| `internal/output` | 0% | 85% | 95% |
| `internal/models` | 0% | 60% | 70% |
| **Overall** | **0%** | **60%** | **80%** |

## 🐛 Current Test Failures

### 1. TestListAccounts - Filter Logic Issue
```
Error: Not equal: expected: 2, actual: 3
```

**Root Cause:** Test setup creates 3 accounts, expects filter to return 2 active.
The issue is in the test cleanup - accounts persist across tests.

**Fix:**
```go
func (suite *AccountRepositoryTestSuite) SetupTest() {
    // Clean database before each test
    suite.db.Exec("DELETE FROM accounts")
    suite.db.Exec("DELETE FROM sqlite_sequence WHERE name='accounts'") // Reset ID
}
```

### 2. TestDuplicateAccountNames - Constraint Mismatch
```
Error: UNIQUE constraint failed: accounts.name
```

**Root Cause:** SQLite enforces UNIQUE on `name` column, but test expects duplicates to be allowed (handled at application level).

**Fix:** Update the model to not enforce uniqueness at DB level for SQLite:
```go
// Remove unique constraint for test environment
// Or: Check if error is uniqueness violation and expect it
assert.Error(suite.T(), err)
assert.Contains(suite.T(), err.Error(), "UNIQUE constraint")
```

## 🚀 Testing Roadmap

### Phase 1: Fix Existing Tests (1-2 days)
- [ ] Fix TestListAccounts (cleanup issue)
- [ ] Fix TestDuplicateAccountNames (constraint handling)
- [ ] Move tests to co-located files for better coverage reporting
- [ ] Add repository tests for `NameExists()`, `GetByName()`, `UpdateBalance()`

### Phase 2: Unit Test Coverage (3-5 days)
- [ ] Add `internal/commands/account_test.go`
  - Test command parsing
  - Test flag validation
  - Test error handling
- [ ] Add `internal/config/config_test.go`
  - Test config file loading
  - Test environment variable parsing
  - Test default values
- [ ] Add `internal/output/output_test.go`
  - Test currency formatting
  - Test table generation
  - Test JSON output
- [ ] Add `internal/db/connection_test.go`
  - Test connection pooling
  - Test health checks

### Phase 3: Integration Tests (2-3 days)
- [ ] Create `tests/integration/setup.go` (test DB utilities)
- [ ] Add `account_workflow_test.go` (full CRUD)
- [ ] Add `config_db_integration_test.go` (config → connection)
- [ ] Add `command_integration_test.go` (CLI → DB)

### Phase 4: E2E & Automation (3-4 days)
- [ ] Create E2E test scripts
- [ ] Set up GitHub Actions CI
- [ ] Add pre-commit hooks for tests
- [ ] Add coverage reporting

### Phase 5: Regression Suite (Ongoing)
- [ ] Document critical user workflows
- [ ] Create regression test suite
- [ ] Set up nightly test runs
- [ ] Add performance benchmarks

## 🤖 Test Automation

### GitHub Actions CI Workflow

**File:** `.github/workflows/test.yml`
```yaml
name: Test

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_PASSWORD: test
          POSTGRES_DB: fintrack_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install dependencies
      run: make deps

    - name: Run unit tests
      run: go test -v -race -coverprofile=coverage.out ./...

    - name: Run integration tests
      env:
        FINTRACK_DB_URL: postgresql://postgres:test@localhost:5432/fintrack_test?sslmode=disable
      run: go test -v -tags=integration ./tests/integration/...

    - name: Generate coverage report
      run: go tool cover -html=coverage.out -o coverage.html

    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        files: ./coverage.out

    - name: Check coverage threshold
      run: |
        coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
        if (( $(echo "$coverage < 60.0" | bc -l) )); then
          echo "Coverage $coverage% is below 60% threshold"
          exit 1
        fi
```

### Pre-commit Hook

**File:** `.git/hooks/pre-commit`
```bash
#!/bin/bash
set -e

echo "Running tests..."
make test

echo "Checking coverage..."
coverage=$(go test -cover ./... 2>&1 | grep -oP '\d+\.\d+(?=% of statements)')
if (( $(echo "$coverage < 60.0" | bc -l) )); then
  echo "❌ Coverage $coverage% is below 60% threshold"
  exit 1
fi

echo "✅ All tests passed with $coverage% coverage"
```

### Makefile Enhancements

Add to existing Makefile:
```makefile
# Run tests with race detection
test-race:
	@echo "Running tests with race detector..."
	$(GO) test -race -v ./...

# Run only unit tests
test-unit:
	@echo "Running unit tests..."
	$(GO) test -v -short ./...

# Run only integration tests
test-integration:
	@echo "Running integration tests..."
	$(GO) test -v -tags=integration ./tests/integration/...

# Watch mode - re-run tests on file changes (requires entr)
test-watch:
	find . -name "*.go" | entr -c make test

# Benchmark tests
benchmark:
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./...

# Check test coverage threshold
test-coverage-check:
	@echo "Checking coverage threshold..."
	@coverage=$$($(GO) test -cover ./... 2>&1 | grep -oP '\d+\.\d+(?=% of statements)' | head -1); \
	if (( $$(echo "$$coverage < 60.0" | bc -l) )); then \
		echo "❌ Coverage $$coverage% is below 60% threshold"; \
		exit 1; \
	else \
		echo "✅ Coverage $$coverage% meets threshold"; \
	fi
```

## 📈 Regression Testing Strategy

### What to Test for Regressions

1. **Critical Paths:**
   - Account creation, update, deletion
   - Transaction recording
   - Balance calculations
   - Report generation

2. **Edge Cases:**
   - Negative balances (credit cards)
   - Zero amounts
   - Very large numbers (overflow)
   - Unicode in names
   - SQL injection attempts

3. **Performance:**
   - Large dataset handling (10k+ transactions)
   - Query performance (< 100ms)
   - Memory usage

### Regression Test Suite

**File:** `tests/regression/suite_test.go`
```go
package regression

import (
    "testing"
)

// Critical user workflows that should never break
func TestRegressionSuite(t *testing.T) {
    tests := []struct {
        name     string
        testFunc func(t *testing.T)
    }{
        {"Account CRUD workflow", testAccountCRUD},
        {"Duplicate name prevention", testDuplicateNames},
        {"Negative balance handling", testNegativeBalances},
        {"Large dataset performance", testLargeDataset},
        {"Unicode support", testUnicodeNames},
    }

    for _, tt := range tests {
        t.Run(tt.name, tt.testFunc)
    }
}
```

## 🎯 Success Metrics

### Code Quality Metrics
- **Coverage:** Maintain > 60% (Phase 1), > 80% (Phase 2)
- **Test Count:** > 100 unit tests, > 20 integration tests
- **Test Speed:** Unit tests < 10s total, Integration < 30s
- **Flakiness:** < 1% flaky test rate

### Regression Prevention
- **CI Success Rate:** > 95%
- **Pre-commit Pass Rate:** > 90%
- **Bug Detection:** Catch regressions before production
- **Mean Time to Detection:** < 1 hour (via CI)

## 📚 Testing Best Practices

### 1. Table-Driven Tests
```go
func TestAccountValidation(t *testing.T) {
    tests := []struct {
        name    string
        account *models.Account
        wantErr bool
    }{
        {"valid checking", &models.Account{Name: "Test", Type: "checking"}, false},
        {"empty name", &models.Account{Name: "", Type: "checking"}, true},
        {"invalid type", &models.Account{Name: "Test", Type: "invalid"}, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateAccount(tt.account)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateAccount() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 2. Test Fixtures
```go
// tests/fixtures/accounts.go
func NewTestAccount(overrides ...func(*models.Account)) *models.Account {
    account := &models.Account{
        Name:     "Test Account",
        Type:     models.AccountTypeChecking,
        Currency: "USD",
        IsActive: true,
    }

    for _, override := range overrides {
        override(account)
    }

    return account
}

// Usage
account := NewTestAccount(func(a *models.Account) {
    a.Name = "Custom Name"
    a.InitialBalance = 1000.0
})
```

### 3. Test Helpers
```go
// tests/helpers/db.go
func SetupTestDB(t *testing.T) *gorm.DB {
    t.Helper()

    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)

    err = db.AutoMigrate(&models.Account{})
    require.NoError(t, err)

    return db
}

func CleanupTestDB(t *testing.T, db *gorm.DB) {
    t.Helper()

    sqlDB, err := db.DB()
    require.NoError(t, err)

    err = sqlDB.Close()
    require.NoError(t, err)
}
```

### 4. Parallel Testing
```go
func TestAccountOperations(t *testing.T) {
    t.Parallel() // Run in parallel with other tests

    t.Run("Create", func(t *testing.T) {
        t.Parallel() // Each subtest runs in parallel
        // Test create logic
    })

    t.Run("Update", func(t *testing.T) {
        t.Parallel()
        // Test update logic
    })
}
```

## 🔍 Code Coverage Analysis

### Running Coverage Analysis
```bash
# Generate coverage report
make test-coverage

# View in browser
open coverage.html

# Get coverage summary
go tool cover -func=coverage.out

# Get package-level coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep -v "total"
```

### Coverage Badges
Add to README.md:
```markdown
[![Coverage](https://codecov.io/gh/yourusername/fintrack/branch/main/graph/badge.svg)](https://codecov.io/gh/yourusername/fintrack)
```

## 🚧 Known Issues

1. **Tests in separate package:** Current tests are in `tests/unit/` which doesn't contribute to coverage. Should be co-located.
2. **No CI/CD:** No automated testing on commits/PRs
3. **Manual E2E:** No automated end-to-end testing
4. **Test failures:** 2 of 7 tests failing in existing suite

## 📖 References

- [Go Testing Best Practices](https://golang.org/doc/code.html#Testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Go Test Coverage](https://go.dev/blog/cover)
- [Table Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
