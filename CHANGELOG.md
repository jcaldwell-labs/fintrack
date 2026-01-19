# Changelog

All notable changes to FinTrack will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

_No unreleased changes_

## [0.1.0] - 2026-01-19 (Debut Release)

This is the debut release of FinTrack - a privacy-first, terminal-based personal finance tracker for developers.

### Added

- **Account management** - Create, list, update, show, and close accounts
- **Transaction tracking** - Add, list, update, delete transactions with category assignment
- **Category management** - Hierarchical categories for income and expenses
- **CSV import** (experimental) - Import transactions from bank CSV exports
- **JSON output** - All list commands support `--json` for scripting
- **Transaction summary** - List command shows income/expense/net totals
- **Currency formatting** - Amounts display with thousand separators ($1,234.56)
- **Helpful first-run experience** - Clear setup instructions when database not configured
- Testing infrastructure with GitHub Actions CI/CD
- Comprehensive testing documentation (TESTING.md)
- Security documentation (SECURITY.md)
- Contributing guidelines (CONTRIBUTING.md)
- PR and issue templates
- DEMO.md - 2-minute demo script
- QUICKSTART.md - 5-minute getting started guide

### Changed

- **Money storage** - All amounts stored as int64 cents (not float64) for precision
- **Database initialization** - Thread-safe with sync.Once pattern
- Enhanced Makefile with additional testing targets
- Coverage threshold set to 40% (realistic baseline)
- README updated with honest "What Works Today" feature list
- CSV import marked as experimental in help text

### Fixed

- **Money precision bug** - float64 → int64 cents conversion (0.1 + 0.2 now equals 0.3)
- **Database race condition** - Added sync.Once for thread-safe initialization
- **CSV amount negation** - Correct income/expense detection during import
- Usage documentation tests - Updated command syntax
- Database initialization in main.go
- Removed unused imports

### Removed

- Stub commands hidden from help output (budget, schedule, remind, project, report, calendar, config)
  - These will be re-enabled when implemented

### Documentation

- README.md with honest feature list and "What Works Today" section
- DEMO.md for presentations and showcases
- QUICKSTART.md for 5-minute onboarding
- CLAUDE.md for development context
- Example configuration file
- Quick reference guide
- Complete roadmap

### Database

- PostgreSQL backend with GORM ORM
- Complete schema with accounts, transactions, categories, import history
- Automatic triggers for balance updates
- ACID compliance for financial data integrity

---

## Release Types

### Major Version (1.0.0 → 2.0.0)

Breaking changes that require migration or code changes by users.

### Minor Version (1.0.0 → 1.1.0)

New features that are backwards compatible.

### Patch Version (1.0.0 → 1.0.1)

Bug fixes that are backwards compatible.

## Change Categories

### Added

New features or functionality added to the project.

### Changed

Changes to existing functionality.

### Deprecated

Features that will be removed in future releases.

### Removed

Features that have been removed.

### Fixed

Bug fixes.

### Security

Security improvements or vulnerability fixes.

---

[Unreleased]: https://github.com/yourusername/fintrack/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/yourusername/fintrack/releases/tag/v0.1.0
