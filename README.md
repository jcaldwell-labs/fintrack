# FinTrack

> Terminal-based personal finance tracking and budgeting application

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/jcaldwell-labs/fintrack/pulls)

**Status:** Phase 1 (MVP) - In Development

## Why FinTrack?

*Your finances. Your terminal. Your control.*

Traditional personal finance apps require cloud accounts, sync your sensitive data to third-party servers, and lock you into proprietary formats. FinTrack is different:

- **Privacy-first**: All data stays local in your PostgreSQL database. No cloud accounts, no tracking, no data harvesting
- **Unix philosophy**: Composable commands, text output, JSON for scripting. Pipe, grep, and automate your finances
- **Developer-friendly**: Single binary, cross-platform, scriptable. Integrates with your existing CLI workflow
- **Full control**: Export your data anytime. No vendor lock-in. Open source forever

**Perfect for:**
- Developers who live in the terminal
- Privacy-conscious users who don't want their financial data in the cloud
- Power users who want to script and automate their budgeting
- Anyone migrating from Mint, YNAB, or spreadsheets who wants more control

## Features

### Current (Phase 1)
- ✅ Account management (create, list, update, close)
- ✅ Transaction tracking and categorization
- ✅ Category management with hierarchical support
- ✅ PostgreSQL backend with ACID compliance
- ✅ JSON output for scripting
- ✅ Cross-platform support (Linux, macOS, Windows)

### Coming Soon
- CSV import with bank-specific mappings
- Budget tracking with alerts
- Recurring transaction scheduling
- Cash flow projections
- Financial reports and analytics
- Calendar view
- Reminder system

## Demo

<!-- TODO: Add asciinema recording -->

**Try it yourself:**
```bash
# Create your first account
fintrack account add "Checking" --type checking --balance 5000

# Add a transaction
fintrack tx add -50.00 --account "Checking" --category "Groceries" --payee "Walmart"

# See your accounts
fintrack account list

# Get JSON for scripting
fintrack tx list --json | jq '.[] | select(.category == "Groceries")'
```

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

#### Category Management

```bash
# List all categories
fintrack category list
fintrack cat list

# Add a new category
fintrack category add "Groceries" expense
fintrack cat add "Salary" income --color "#00FF00" --icon "💰"

# Add subcategory
fintrack category add "Coffee" expense --parent "Food & Dining"

# Show category details
fintrack category show 5

# Update category
fintrack category update 5 --name "New Name" --color "#FF5733"

# Delete category (non-system only)
fintrack category delete 5

# Filter by type
fintrack category list --type expense
fintrack cat list -t income
```

**Example output:**
```
ID  NAME              TYPE      PARENT          SYSTEM  COLOR     ICON
1   Salary            income                    Yes               💰
2   Groceries         expense   Food & Dining   No      #FF5733   🛒
3   Transportation    expense                   Yes               🚗
4   Gas/Fuel          expense   Transportation  Yes               ⛽
```

#### Transaction Management

```bash
# Add a transaction
fintrack transaction add -50.00 --account "Checking" --category "Groceries" --payee "Walmart"
fintrack tx add 2500.00 -a "Checking" -c "Salary" --date 2024-01-15

# Add with tags
fintrack tx add -30.00 -a "Checking" -c "Food & Dining" --tags "business,reimbursable"

# List transactions
fintrack transaction list
fintrack tx list

# Filter transactions
fintrack tx list --account "Checking"
fintrack tx list --category "Groceries" --start 2024-01-01 --end 2024-01-31
fintrack tx list --type expense --limit 50

# Show transaction details
fintrack transaction show 42
fintrack tx show 100

# Update transaction
fintrack tx update 42 --amount -75.50
fintrack tx update 42 --category "Entertainment" --payee "Netflix"

# Delete transaction
fintrack tx delete 42
```

**Example output:**
```
ID  DATE        ACCOUNT   CATEGORY      PAYEE      AMOUNT     TYPE
42  2024-01-15  Checking  Groceries     Walmart    -$50.00    expense
43  2024-01-16  Checking  Salary        Employer   $2500.00   income
44  2024-01-17  Checking  Gas/Fuel      Shell      -$45.00    expense

Summary:
  Total Transactions: 3
  Income: $2500.00
  Expenses: $95.00
  Net: $2405.00
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
│   │   ├── category.go        # Category management
│   │   ├── transaction.go     # Transaction management
│   │   └── stubs.go           # Placeholder commands
│   ├── core/                  # Business logic (coming soon)
│   ├── models/                # Data models
│   │   └── models.go
│   ├── db/                    # Database layer
│   │   ├── connection.go
│   │   └── repositories/
│   │       ├── account_repository.go
│   │       ├── category_repository.go
│   │       └── transaction_repository.go
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

## Comparison

| Feature | FinTrack | Mint/YNAB | Ledger-CLI | Spreadsheets |
|---------|----------|-----------|------------|--------------|
| Local-first data | ✅ | ❌ | ✅ | ✅ |
| No cloud account | ✅ | ❌ | ✅ | ❌ (Google/MS) |
| JSON output | ✅ | ❌ | ❌ | ❌ |
| Scriptable CLI | ✅ | ❌ | ✅ | ❌ |
| ACID database | ✅ | ? | ❌ (text) | ❌ |
| Budgeting | 🔄 | ✅ | ❌ | Manual |
| Bank import | 🔄 | ✅ | ✅ | Manual |
| Cross-platform | ✅ | Web | ✅ | Web |
| Open source | ✅ | ❌ | ✅ | ❌ |

## Roadmap

See [.github/planning/ROADMAP.md](.github/planning/ROADMAP.md) for the detailed implementation roadmap.

**Current Phase:** Phase 1 - Core Foundation (MVP)

| Phase | Focus | Status |
|-------|-------|--------|
| Phase 1 | Core Foundation (accounts, transactions, categories) | In Progress |
| Phase 2 | Budgeting & Scheduling | Planned |
| Phase 3 | Projections & Analytics | Planned |
| Phase 4 | Advanced Features (multi-currency, Plaid) | Planned |
| Phase 5 | Polish & Optimization (TUI, shell completion) | Planned |

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

## Documentation

| Document | Description |
|----------|-------------|
| [Quick Reference](docs/FINTRACK_QUICKREF.md) | Command cheat sheet |
| [System Design](docs/FINANCE_TRACKER_PLAN.md) | Complete architecture |
| [Roadmap](.github/planning/ROADMAP.md) | Implementation timeline |
| [Contributing](CONTRIBUTING.md) | How to contribute |
| [Testing](TESTING.md) | Test strategy and coverage |
| [Security](SECURITY.md) | Security guidelines |

## Community

- **Issues:** [Report bugs or request features](https://github.com/jcaldwell-labs/fintrack/issues)
- **Pull Requests:** [Contributions welcome!](https://github.com/jcaldwell-labs/fintrack/pulls)
- **Discussions:** [Ask questions, share ideas](https://github.com/jcaldwell-labs/fintrack/discussions)

## Acknowledgments

Inspired by:
- [ledger-cli](https://www.ledger-cli.org/) - Plain-text accounting
- [hledger](https://hledger.org/) - Accounting tools
- [YNAB](https://www.youneedabudget.com/) - Budget philosophy

---

**Built with Unix philosophy principles: do one thing well, composable, text-based**
