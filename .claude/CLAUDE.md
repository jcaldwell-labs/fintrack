# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

FinTrack is a terminal-based personal finance tracking and budgeting application written in Go. It follows Unix philosophy principles: composable commands, text output, scriptable interfaces, and privacy-first local storage.

**Status:** Phase 1 (MVP) - In Development

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
```

### Code Quality
```bash
# Format code
make fmt

# Lint code (requires golangci-lint)
make lint

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
│   ├── account.go  # Account management (CRUD, show, close)
│   └── stubs.go    # Placeholder commands for future features
├── config/         # Viper configuration management
│   └── config.go   # Config struct, env loading, defaults
├── db/             # Database layer
│   ├── connection.go           # GORM connection, pooling, health checks
│   └── repositories/           # Repository pattern for data access
│       └── account_repository.go
├── models/         # GORM models and domain types
│   └── models.go   # Account, Transaction, Budget, etc.
└── output/         # Output formatters (table, JSON)
    └── output.go
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

Example test location: `tests/unit/account_repository_test.go`

### Adding New Commands
1. Create command function in `internal/commands/`
2. Define Cobra command with aliases and flags
3. Implement RunE function with repository logic
4. Add to root command in `cmd/fintrack/main.go`
5. Write tests in `tests/unit/` or `tests/integration/`

### Adding New Repository Methods
1. Add method to repository struct in `internal/db/repositories/`
2. Use GORM query builder for database operations
3. Return domain models from `internal/models`
4. Handle `gorm.ErrRecordNotFound` appropriately
5. Write unit tests

### Database Migrations
Currently using GORM AutoMigrate for development. For production, migrations should be in `migrations/` directory.

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

## Phase 1 Feature Status
- [x] Account CRUD operations
- [x] PostgreSQL connection and GORM setup
- [x] JSON output support
- [ ] Transaction CRUD operations
- [ ] Category management
- [ ] CSV import
- [ ] Basic reporting

See `docs/FINTRACK_ROADMAP.md` for complete implementation timeline.

## Related Documentation
- **Planning:** `docs/FINANCE_TRACKER_PLAN.md` - Complete system design
- **Quick Reference:** `docs/FINTRACK_QUICKREF.md` - Command cheat sheet
- **Roadmap:** `docs/FINTRACK_ROADMAP.md` - Implementation timeline
- **Config Example:** `fintrack_config.example.yaml` - Full configuration reference
- **Security:** `SECURITY.md` - Security best practices, what to commit vs. keep local
