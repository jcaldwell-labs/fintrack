# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

FinTrack is a terminal-based personal finance tracking and budgeting application written in Go. It follows Unix philosophy principles: composable commands, text output, scriptable interfaces, and privacy-first local storage.

**Status:** Phase 1 (MVP) - In Development

**Test Coverage:** ~45-50% overall (targeting 60% for Phase 1, 80% for production)

## Common Commands

### Build and Run
```bash
# Build the application
make build

# Run without installing
make run

# Install to /usr/local/bin
make install

# Run with live reload (requires air)
make dev
```

### Testing
```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# View coverage report in browser
open coverage.html

# Run tests with race detection
make test-race

# Run only unit tests (fast)
make test-unit

# Run only integration tests
make test-integration

# Check coverage threshold (60% target)
make test-coverage-check

# Watch mode - re-run tests on file changes (requires entr)
make test-watch

# Run benchmarks
make benchmark
```

### Code Quality
```bash
# Format code
make fmt

# Lint code (requires golangci-lint)
make lint

# Run all quality checks (fmt, lint, test-race, coverage)
make quality

# Download and tidy dependencies
make deps

# Verify dependencies
make verify
```

### Cross-Platform Builds
```bash
# Build for all platforms
make build-all

# Platform-specific builds
make build-linux
make build-darwin
make build-windows
```

### Database Setup
```bash
# Create database
createdb fintrack

# Run migrations (if available)
psql -d fintrack -f migrations/schema.sql

# Or use GORM auto-migration through the application
```

## Architecture

### Technology Stack
- **Language:** Go 1.21+
- **CLI Framework:** Cobra (command structure and parsing)
- **Config:** Viper (YAML/ENV support with automatic environment variable binding)
- **Database:** PostgreSQL 12+ with GORM ORM
- **Testing:** Testify

### Project Structure
```
internal/
├── commands/       # Cobra command implementations
│   ├── account.go          # Account management (CRUD, show, close)
│   ├── account_test.go     # Command tests (co-located)
│   ├── stubs.go            # Placeholder commands for future features
│   └── stubs_test.go       # Stub command tests
├── config/         # Viper configuration management
│   ├── config.go           # Config struct, env loading, defaults
│   └── config_test.go      # Config tests (93.9% coverage)
├── db/             # Database layer
│   ├── connection.go       # GORM connection, pooling, health checks
│   ├── connection_test.go  # Connection tests (36.8% coverage)
│   └── repositories/       # Repository pattern for data access
│       ├── account_repository.go
│       └── account_repository_test.go  # Repository tests (94.9% coverage)
├── models/         # GORM models and domain types
│   ├── models.go           # Account, Transaction, Budget, etc.
│   └── models_test.go      # Model tests (100.0% coverage)
└── output/         # Output formatters (table, JSON)
    ├── output.go
    └── output_test.go      # Output formatter tests (91.8% coverage)

tests/
├── integration/    # Integration tests (future)
└── unit/          # Legacy unit tests (being migrated to co-located)
```

### Key Architecture Patterns

**Repository Pattern:**
All data access goes through repository structs in `internal/db/repositories/`. Each repository:
- Takes a `*gorm.DB` in constructor
- Encapsulates all queries for a specific model
- Returns domain models from `internal/models`
- Example: `AccountRepository.GetByID()`, `AccountRepository.Create()`

**Command Pattern (Cobra):**
Each feature is a Cobra command with subcommands:
- Root command in `cmd/fintrack/main.go`
- Feature commands in `internal/commands/`
- Each command creates its own repository instances
- Commands use `output` package for consistent formatting

**Configuration Hierarchy:**
1. Config file: `~/.config/fintrack/config.yaml`
2. Environment variables: `FINTRACK_*` (auto-bound via Viper)
3. Command-line flags (defined in commands)

Example: Database URL can be set via:
- `database.url` in YAML
- `FINTRACK_DB_URL` environment variable
- Individual components: `FINTRACK_DB_HOST`, `FINTRACK_DB_PORT`, etc.

**Database Connection:**
- Connection initialized in `db.Init()` called from root command's `PersistentPreRunE`
- Singleton pattern: `db.Get()` returns shared `*gorm.DB` instance
- Connection pooling configured from config
- Health checks available via `db.IsConnected()`

**Output Formatting:**
All commands support dual output:
- **Table format:** Human-readable table output (default)
- **JSON format:** Machine-readable via `--json` flag

Use `output.GetFormat(cmd)` to check mode and `output.Print()` for consistent JSON output.

### Domain Models

All models in `internal/models/models.go` follow GORM conventions:
- Struct tags for PostgreSQL mapping
- Time fields use `time.Time`
- Soft deletes via `is_active` boolean (not GORM's DeletedAt)
- Currency stored as string codes (e.g., "USD")
- Amounts as `float64` (consider decimal library for production)

**Account Model:**
- Unique name constraint (per active accounts)
- Types: checking, savings, credit, cash, investment, loan
- Tracks both initial and current balance
- Supports institution and last 4 digits of account number

**Transaction Model:**
- Linked to accounts and categories
- Positive amounts = income, negative = expenses
- Supports transfers between accounts
- Tagging system via PostgreSQL arrays
- Reconciliation tracking

## Configuration

### Database Configuration
Database connection can be configured via:

**Option 1: Direct URL**
```yaml
database:
  url: "postgresql://user:pass@localhost:5432/fintrack?sslmode=disable"
```

**Option 2: Component-based**
```yaml
database:
  host: localhost
  port: 5432
  database: fintrack
  user: fintrack_user
  password: secret  # Or use FINTRACK_DB_PASSWORD env var
  sslmode: disable
  max_connections: 10
  max_idle_connections: 2
  connection_max_lifetime: "1h"
```

**Option 3: Environment variables**
```bash
export FINTRACK_DB_URL="postgresql://localhost:5432/fintrack?sslmode=disable"
# Or individual components
export FINTRACK_DB_HOST="localhost"
export FINTRACK_DB_PORT="5432"
export FINTRACK_DB_PASSWORD="secret"
```

### Application Defaults
```yaml
defaults:
  currency: "USD"
  date_format: "2006-01-02"  # Go time format
  timezone: "Local"

output:
  default_format: "table"  # or "json"
  color: true
  unicode: true
```

## Development Workflow

### Test-Driven Development (TDD)
This project follows TDD:
1. **Red:** Write failing test first
2. **Green:** Implement minimum code to pass
3. **Refactor:** Improve code quality

**Testing Best Practices:**
- Tests are co-located with source code (e.g., `account.go` → `account_test.go`)
- Use table-driven tests for multiple test cases
- Use testify/assert and testify/require for assertions
- Aim for high coverage: models (100%), repositories (90%+), output (90%+), config (90%+)
- Use race detection (`make test-race`) to catch concurrency issues
- Write benchmarks for performance-critical code

**Current Coverage Status:**
- `internal/models`: 100.0% ✅
- `internal/db/repositories`: 94.9% ✅
- `internal/config`: 93.9% ✅
- `internal/output`: 91.8% ✅
- `internal/commands`: 15.9% ⚠️ (needs improvement)
- `internal/db`: 36.8% ⚠️ (needs improvement)
- `cmd/fintrack`: 0.0% ❌ (integration tests planned)

### Adding New Commands
1. **Write tests first** in co-located `*_test.go` file (TDD)
2. Create command function in `internal/commands/`
3. Define Cobra command with aliases and flags
4. Implement RunE function with repository logic
5. Add to root command in `cmd/fintrack/main.go`
6. Run tests and ensure they pass: `make test`
7. Check coverage: `make test-coverage-check`

### Adding New Repository Methods
1. **Write tests first** in co-located `*_test.go` file (TDD)
2. Add method to repository struct in `internal/db/repositories/`
3. Use GORM query builder for database operations
4. Return domain models from `internal/models`
5. Handle `gorm.ErrRecordNotFound` appropriately
6. Run tests: `make test`
7. Target 90%+ coverage for repository code

### Database Migrations
Currently using GORM AutoMigrate for development. For production, migrations should be in `migrations/` directory.

### CI/CD & Automation

**GitHub Actions Workflows:**
- `.github/workflows/test.yml` - Test & Coverage pipeline
  - Runs on push to `main`/`develop` and all PRs
  - Tests with Go 1.21 and 1.22 (matrix)
  - PostgreSQL 15 service container for integration tests
  - Race detection enabled
  - Coverage threshold check (currently 45%, targeting 60%)
  - Codecov integration for coverage tracking
  - Security scanning with Gosec

- `.github/workflows/pr-checks.yml` - PR validation
  - Code formatting (go fmt)
  - Linting (golangci-lint)
  - Go vet checks

- `.github/workflows/release.yml` - Release automation
  - Cross-platform binary builds
  - Automated releases

**Coverage Thresholds:**
- Current minimum: 45% (enforced in CI)
- Phase 1 target: 60%
- Long-term target: 80%

**Pre-commit Best Practices:**
- Always run `make test` before committing
- Run `make fmt` to format code
- Run `make lint` to catch issues
- Use `make quality` to run all checks at once

## Important Notes

### Account Lookups
Commands can accept account ID or name:
```go
id, err := parseAccountID(args[0])  // Tries numeric ID first, then name lookup
```

### Error Handling
Use `output.PrintError(cmd, err)` for consistent error formatting across table/JSON modes.

### Currency Formatting
Use `output.FormatCurrency(amount, currency)` for consistent currency display.

### Unique Constraints
Account names must be unique among active accounts. The `NameExists()` repository method supports exclusion for updates.

## Common Pitfalls & Best Practices

### Testing Gotchas
1. **Test Database Cleanup:** Always clean test database between tests to avoid state leakage
   ```go
   func setupTestDB(t *testing.T) *gorm.DB {
       db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
       db.AutoMigrate(&models.Account{})  // Fresh schema per test
       return db
   }
   ```

2. **Race Conditions:** Always run `make test-race` to catch concurrency issues

3. **Coverage Blind Spots:** Co-located tests automatically count toward coverage. If adding tests to `tests/unit/`, they won't contribute to package coverage metrics.

4. **GORM ErrRecordNotFound:** Don't treat as error for GET operations. Return nil, nil instead.
   ```go
   if errors.Is(err, gorm.ErrRecordNotFound) {
       return nil, nil  // Not found is not an error
   }
   ```

### Code Quality
1. **Always format before committing:** `make fmt`
2. **Check linting:** `make lint` (requires golangci-lint)
3. **Verify coverage threshold:** `make test-coverage-check` (minimum 45%, targeting 60%)
4. **Run full quality check:** `make quality` (runs fmt, lint, test-race, coverage)

### Error Handling
- Use `output.PrintError(cmd, err)` for consistent error output across table/JSON modes
- Wrap errors with context: `fmt.Errorf("failed to create account: %w", err)`
- Always handle database errors gracefully

### Performance
- Use GORM's `Preload()` for eager loading to avoid N+1 queries
- Keep repository methods focused and single-purpose
- Benchmark performance-critical paths: `make benchmark`

## Test Writing Patterns

### Table-Driven Tests (Preferred)
```go
func TestAccountValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   *models.Account
        wantErr bool
        errMsg  string
    }{
        {
            name: "valid checking account",
            input: &models.Account{Name: "Test", Type: "checking"},
            wantErr: false,
        },
        {
            name: "empty name",
            input: &models.Account{Name: "", Type: "checking"},
            wantErr: true,
            errMsg: "name is required",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateAccount(tt.input)
            if tt.wantErr {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

### Test Fixtures
Use builder pattern for test data:
```go
func newTestAccount(opts ...func(*models.Account)) *models.Account {
    acc := &models.Account{
        Name: "Test Account",
        Type: models.AccountTypeChecking,
        Currency: "USD",
        IsActive: true,
    }
    for _, opt := range opts {
        opt(acc)
    }
    return acc
}

// Usage:
account := newTestAccount(func(a *models.Account) {
    a.Name = "Custom Name"
    a.InitialBalance = 1000.0
})
```

### Mocking Best Practices
- Mock external dependencies (database, APIs) in unit tests
- Use real database (SQLite in-memory) for repository tests
- Keep mocks simple and focused on behavior being tested

## Debugging Tips

### Running Specific Tests
```bash
# Run single test
go test -v -run TestAccountCreate ./internal/commands/

# Run tests in single package
go test -v ./internal/db/repositories/

# Run with verbose output
go test -v -race ./...

# Run with coverage for specific package
go test -v -cover ./internal/commands/
```

### Debugging Test Failures
1. **Check test output:** `go test -v` shows detailed output
2. **Enable database logging:** Set `GORM_LOG_LEVEL=info` for query debugging
3. **Check race conditions:** `make test-race` catches concurrency bugs
4. **Isolate the test:** Run failing test alone to check for test pollution
5. **Review coverage gaps:** `make test-coverage` shows uncovered code paths

### Common Test Failures
- **"UNIQUE constraint failed":** Database not cleaned between tests
- **"record not found":** Test data not set up properly
- **"nil pointer dereference":** Missing nil checks or uninitialized mocks
- **Flaky tests:** Usually indicates race conditions or timing issues

### CI/CD Debugging
- Check GitHub Actions logs for detailed error messages
- Local reproduction: Run same commands as CI (`make test-race`, `make test-coverage-check`)
- Coverage threshold failures: Run `make test-coverage` to see which packages need more tests
- Linting failures: Run `make lint` locally before pushing

## Phase 1 Feature Status
- [x] Account CRUD operations
- [x] PostgreSQL connection and GORM setup
- [x] JSON output support
- [x] Comprehensive test coverage (~45-50% overall)
- [x] CI/CD pipeline with GitHub Actions
- [x] Co-located unit tests
- [x] Security scanning (Gosec)
- [ ] Transaction CRUD operations
- [ ] Category management
- [ ] CSV import
- [ ] Basic reporting
- [ ] Integration test suite

See `docs/FINTRACK_ROADMAP.md` for complete implementation timeline.

## Related Documentation
- **Planning:** `docs/FINANCE_TRACKER_PLAN.md` - Complete system design
- **Quick Reference:** `docs/FINTRACK_QUICKREF.md` - Command cheat sheet
- **Roadmap:** `docs/FINTRACK_ROADMAP.md` - Implementation timeline
- **Testing:** `TESTING.md` - Comprehensive testing strategy, coverage targets, best practices
- **Contributing:** `CONTRIBUTING.md` - Contribution guidelines, development workflow
- **Security:** `SECURITY.md` - Security best practices, what to commit vs. keep local
- **Config Example:** `fintrack_config.example.yaml` - Full configuration reference
- **Changelog:** `CHANGELOG.md` - Version history and changes
