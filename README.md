# FinTrack

> Terminal-based personal finance tracking and budgeting application

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)

FinTrack is a command-line tool for managing personal finances following Unix philosophy principles. Track transactions, set budgets, schedule recurring expenses, and project cash flow - all from your terminal with complete privacy.

**Status:** 🚧 Phase 1 (MVP) - In Development

## Features

### Current (Phase 1)
- ✅ Account management (create, list, update, close)
- ✅ PostgreSQL backend with ACID compliance
- ✅ JSON output for scripting
- ✅ Cross-platform support (Linux, macOS, Windows)

### Coming Soon
- 🔄 Transaction tracking and categorization
- 🔄 CSV import with bank-specific mappings
- 🔄 Budget tracking with alerts
- 🔄 Recurring transaction scheduling
- 🔄 Cash flow projections
- 🔄 Financial reports and analytics
- 🔄 Calendar view
- 🔄 Reminder system

## Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Make (optional, for build automation)

### Installation

1. **Clone the repository:**
```bash
git clone https://github.com/fintrack/fintrack.git
cd fintrack
```

2. **Install dependencies:**
```bash
make deps
```

3. **Set up PostgreSQL database:**
```bash
# Create database
createdb fintrack

# Run migrations (from planning docs)
psql -d fintrack -f ../fintrack_schema.sql
```

4. **Configure database connection:**

Create `~/.config/fintrack/config.yaml`:
```yaml
database:
  url: "postgresql://localhost:5432/fintrack?sslmode=disable"
  # Or use environment variable:
  # export FINTRACK_DB_URL="postgresql://localhost:5432/fintrack"
```

5. **Build and install:**
```bash
make build
make install
```

### Usage

#### Account Management

```bash
# List all accounts
fintrack account list
fintrack a ls

# Add a new account
fintrack account add "Chase Checking" --type checking --balance 5000
fintrack a add "Amex Gold" -t credit -b -1200

# Show account details
fintrack account show 1
fintrack a show "Chase Checking"

# Update account
fintrack account update 1 --name "Chase Premier Checking"

# Close account
fintrack account close 1

# JSON output (for scripting)
fintrack account list --json
```

**Example output:**
```
ID  NAME             TYPE       BALANCE      LAST ACTIVITY
1   Chase Checking   checking   $5,234.10    2025-11-16
2   Amex Gold        credit     -$1,234.00   2025-11-15
3   Ally Savings     savings    $15,000.00   2025-11-10
```

## Architecture

### Technology Stack

- **Language:** Go 1.21+
- **Database:** PostgreSQL 12+
- **CLI Framework:** Cobra
- **Config:** Viper (YAML/ENV support)
- **ORM:** GORM
- **Testing:** Testify

### Project Structure

```
fintrack/
├── cmd/fintrack/              # CLI entry point
│   └── main.go
├── internal/
│   ├── commands/              # Command implementations
│   │   ├── account.go         # Account management
│   │   └── stubs.go           # Placeholder commands
│   ├── core/                  # Business logic (coming soon)
│   ├── models/                # Data models
│   │   └── models.go
│   ├── db/                    # Database layer
│   │   ├── connection.go
│   │   └── repositories/
│   │       └── account_repository.go
│   ├── output/                # Output formatters
│   │   └── output.go
│   └── config/                # Configuration
│       └── config.go
├── tests/
│   ├── integration/           # Integration tests
│   └── unit/                  # Unit tests
├── Makefile                   # Build automation
└── README.md
```

### Design Principles

1. **Unix Philosophy:** Do one thing well, composable commands, text I/O
2. **Privacy-First:** Local storage only, no cloud dependencies
3. **Cross-Platform:** Single binary for Linux, macOS, Windows
4. **Developer-Friendly:** JSON output, scriptable, pipeable
5. **Test-Driven:** All features developed with TDD

## Configuration

Configuration can be provided via:

1. **Config file:** `~/.config/fintrack/config.yaml`
2. **Environment variables:** `FINTRACK_*`
3. **Command-line flags**

### Example Configuration

```yaml
database:
  url: "postgresql://localhost:5432/fintrack?sslmode=disable"
  max_connections: 10

defaults:
  currency: "USD"
  date_format: "2006-01-02"

alerts:
  enabled: true
  threshold: 0.80

output:
  default_format: "table"
  color: true
  unicode: true
```

### Environment Variables

```bash
# Database connection
export FINTRACK_DB_URL="postgresql://localhost:5432/fintrack"

# Or individual components
export FINTRACK_DB_PASSWORD="secret"
export FINTRACK_DB_HOST="localhost"
export FINTRACK_DB_PORT="5432"
```

## Development

### Building

```bash
# Build binary
make build

# Build for all platforms
make build-all

# Run without installing
make run
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# View coverage report
open coverage.html
```

### Code Quality

```bash
# Format code
make fmt

# Lint code (requires golangci-lint)
make lint

# Verify dependencies
make verify
```

## Database Schema

The application uses a comprehensive PostgreSQL schema with:

- **8 core tables:** accounts, transactions, categories, budgets, recurring_items, reminders, cash_flow_projections, import_history
- **Automatic triggers:** Balance updates, timestamp tracking
- **Materialized views:** Performance-optimized reporting
- **ACID compliance:** Financial data integrity

See `../fintrack_schema.sql` for the complete schema.

## Roadmap

### Phase 1: Core Foundation (Current)
- [x] Project setup and structure
- [x] Database connection and models
- [x] Account CRUD operations
- [ ] Transaction CRUD operations
- [ ] Category management
- [ ] CSV import (generic format)
- [ ] Basic reporting

**Timeline:** 4-6 weeks

### Phase 2: Budgeting & Scheduling
- [ ] Budget tracking with alerts
- [ ] Recurring transaction templates
- [ ] Reminder system
- [ ] Calendar view

**Timeline:** 3-4 weeks

### Phase 3: Projections & Analytics
- [ ] Cash flow projection engine
- [ ] Trend analysis
- [ ] Savings rate calculation

**Timeline:** 4-5 weeks

### Phase 4: Advanced Features
- [ ] Bank-specific CSV mappings
- [ ] Multi-currency support
- [ ] Split transactions
- [ ] Plaid integration (optional)

**Timeline:** 5-6 weeks

### Phase 5: Polish & Optimization
- [ ] Interactive TUI mode
- [ ] Shell completion
- [ ] Performance optimization
- [ ] Professional packaging

**Timeline:** 3-4 weeks

See `../FINTRACK_ROADMAP.md` for detailed implementation plan.

## Contributing

Contributions are welcome! Please follow these guidelines:

1. **Fork the repository**
2. **Create a feature branch:** `git checkout -b feature/my-feature`
3. **Write tests first** (TDD approach)
4. **Implement the feature**
5. **Run tests:** `make test`
6. **Format code:** `make fmt`
7. **Submit a pull request**

### Development Workflow

This project follows test-driven development (TDD):

1. ❌ **Red:** Write failing test
2. ✅ **Green:** Implement minimum code to pass
3. ♻️ **Refactor:** Improve code quality

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Related Documentation

- **Planning:** `../FINANCE_TRACKER_PLAN.md` - Complete system design
- **Schema:** `../fintrack_schema.sql` - Database schema
- **Quick Reference:** `../FINTRACK_QUICKREF.md` - Command cheat sheet
- **Roadmap:** `../FINTRACK_ROADMAP.md` - Implementation timeline
- **Config Example:** `../fintrack_config.example.yaml` - Full configuration

## Support

- **Issues:** [GitHub Issues](https://github.com/fintrack/fintrack/issues)
- **Documentation:** See planning docs in parent directory

## Acknowledgments

Inspired by:
- [ledger-cli](https://www.ledger-cli.org/) - Plain-text accounting
- [hledger](https://hledger.org/) - Accounting tools
- [YNAB](https://www.youneedabudget.com/) - Budget philosophy

---

**Built with ❤️ following Unix philosophy principles**
