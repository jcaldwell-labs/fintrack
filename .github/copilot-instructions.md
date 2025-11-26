# GitHub Copilot Instructions

This file provides guidance to GitHub Copilot when working with code in this repository.

## Project Overview

FinTrack is a terminal-based personal finance tracking and budgeting application written in Go. It follows Unix philosophy principles: composable commands, text output, scriptable interfaces, and privacy-first local storage.

**Status**: Phase 1 (MVP) - In Development
**Language**: Go 1.21+
**Dependencies**: Cobra (CLI), Viper (config), GORM (ORM), PostgreSQL

Key features:
- Account management (CRUD operations)
- Category management (hierarchical)
- Dual output: table (human-readable) and JSON (machine-readable)
- Test-Driven Development (TDD) workflow
- Usage documentation tests (executable markdown)

## Build System

```bash
# Build
make build          # Build the application to bin/fintrack
make build-all      # Cross-platform builds (Linux, macOS, Windows)
make install        # Install to /usr/local/bin

# Testing
make test           # Run all tests
make test-race      # Run tests with race detector
make test-unit      # Run only unit tests (fast)
make test-integration # Run integration tests
make test-usage     # Run usage documentation tests
make test-coverage  # Run tests with coverage report
make test-coverage-check # Verify 60% coverage threshold

# Code Quality
make fmt            # Format code with go fmt
make lint           # Lint with golangci-lint
make quality        # Run all quality checks (fmt, lint, test-race, coverage)

# Development
make run            # Build and run
make dev            # Live reload (requires air)
make test-watch     # Watch mode for tests (requires entr)

# Dependencies
make deps           # Download and tidy dependencies
make verify         # Verify dependencies

make help           # Show all targets
```

## Architecture

### Directory Structure

```
fintrack/
├── cmd/fintrack/              # CLI entry point
│   └── main.go                # Root command setup, global flags
├── internal/
│   ├── commands/              # Cobra command implementations
│   │   ├── account.go         # Account CRUD operations
│   │   ├── account_test.go    # Co-located tests
│   │   ├── category.go        # Category management
│   │   └── stubs.go           # Placeholder commands
│   ├── config/                # Viper configuration
│   │   ├── config.go          # Config struct, defaults
│   │   └── config_test.go     # 93.9% coverage
│   ├── db/                    # Database layer
│   │   ├── connection.go      # GORM connection, pooling
│   │   └── repositories/      # Repository pattern
│   │       ├── account_repository.go
│   │       └── account_repository_test.go  # 94.9% coverage
│   ├── models/                # GORM models
│   │   ├── models.go          # Account, Transaction, Budget, etc.
│   │   └── models_test.go     # 100.0% coverage
│   └── output/                # Output formatters
│       ├── output.go          # Table and JSON output
│       └── output_test.go     # 91.8% coverage
└── tests/
    ├── integration/           # Integration tests
    └── usage/                 # Executable documentation tests
        ├── runner.go
        └── 01-account-management.md
```

### Key Patterns

**Repository Pattern:**
- All data access via `internal/db/repositories/`
- Each repository takes `*gorm.DB` in constructor
- Returns domain models from `internal/models`

**Command Pattern (Cobra):**
- Root command in `cmd/fintrack/main.go`
- Feature commands in `internal/commands/`
- Each command creates its own repository instances

**Configuration Hierarchy:**
1. Config file: `~/.config/fintrack/config.yaml`
2. Environment variables: `FINTRACK_*`
3. Command-line flags

**Output Formatting:**
- Table format (default): Human-readable
- JSON format (`--json`): Machine-readable
- Use `output.GetFormat(cmd)` and `output.Print()`

## Code Style and Conventions

- **Go conventions**: gofmt, effective go guidelines
- **Testing**: Table-driven tests, testify assertions
- **Coverage targets**: models (100%), repos (90%+), config (90%+)
- **Error handling**: Wrap errors with context, use `output.PrintError()`
- **Database**: GORM with PostgreSQL, repository pattern

**Test Writing Pattern:**
```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "foo", "bar", false},
        {"invalid input", "", "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if tt.wantErr {
                require.Error(t, err)
            } else {
                require.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }
}
```

## Before Committing (Required Steps)

Run these commands before every commit:

1. **Format**: `make fmt` - Format code
2. **Lint**: `make lint` - Check for issues
3. **Test**: `make test-race` - Run tests with race detection
4. **Coverage**: `make test-coverage-check` - Verify 60% threshold

```bash
# Quick pre-commit check (all in one)
make quality
```

## Common Development Tasks

### Adding New Commands
1. Write tests first in co-located `*_test.go` file (TDD)
2. Create command function in `internal/commands/`
3. Define Cobra command with aliases and flags
4. Implement RunE function with repository logic
5. Add to root command in `cmd/fintrack/main.go`
6. Run tests: `make test`

### Adding Repository Methods
1. Write tests first in co-located `*_test.go` file (TDD)
2. Add method to repository struct
3. Use GORM query builder
4. Return domain models from `internal/models`
5. Target 90%+ coverage

### Usage Documentation Tests
Add executable examples to `tests/usage/*.md`:
```markdown
## Test: Create account
### Execute
```bash
fintrack account create "Test" --type checking
```
### Expected Output
```
Account created successfully
ID: <number>
```
```

**Wildcards**: `<any>`, `<number>`, `<date>`, `<uuid>`, `<money>`

## Pull Request Standards

When creating PRs, follow these rules:

1. **Always link the issue**: Use `Fixes #N` or `Closes #N`
2. **Fill in all sections**: Never leave placeholder text

**Required PR format:**
```markdown
## Summary
[2-3 sentences describing what and why]

Fixes #[issue-number]

## Changes
- [Actual change 1]
- [Actual change 2]

## Testing
- [x] All tests pass (`make test`)
- [x] Race detection clean (`make test-race`)
- [x] Coverage threshold met (`make test-coverage-check`)

## Type
- [x] New feature | Bug fix | Refactor | Docs | CI
```

## Database Configuration

**Option 1: Direct URL**
```yaml
database:
  url: "postgresql://user:pass@localhost:5432/fintrack?sslmode=disable"
```

**Option 2: Components**
```yaml
database:
  host: localhost
  port: 5432
  database: fintrack
  user: fintrack_user
  password: secret
```

**Option 3: Environment**
```bash
export FINTRACK_DB_URL="postgresql://localhost:5432/fintrack"
# Or individual: FINTRACK_DB_HOST, FINTRACK_DB_PORT, etc.
```

## Debugging Tips

### Running Specific Tests
```bash
go test -v -run TestAccountCreate ./internal/commands/
go test -v -cover ./internal/db/repositories/
```

### Common Test Failures
- **"UNIQUE constraint failed"**: Database not cleaned between tests
- **"record not found"**: Test data not set up properly
- **Flaky tests**: Usually race conditions - run `make test-race`

## Additional Documentation

- **docs/FINTRACK_ROADMAP.md** - Implementation timeline
- **docs/FINTRACK_QUICKREF.md** - Command cheat sheet
- **TESTING.md** - Testing strategy and coverage targets
- **CONTRIBUTING.md** - Contribution guidelines
