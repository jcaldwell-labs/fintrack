# Changelog

All notable changes to FinTrack will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Testing infrastructure with GitHub Actions CI/CD
- Comprehensive testing documentation (TESTING.md)
- Security documentation (SECURITY.md)
- Contributing guidelines (CONTRIBUTING.md)
- PR and issue templates

### Changed
- Enhanced Makefile with additional testing targets
- Updated .claude/CLAUDE.md with security reference

### Fixed
- Database initialization in main.go
- Removed unused imports

## [0.1.0] - 2025-11-16

### Added
- Initial project setup
- Account management (create, list, update, close)
- PostgreSQL database backend
- JSON output support
- Cross-platform support (Linux, macOS, Windows)
- Configuration via YAML file or environment variables
- Command-line interface using Cobra
- GORM for database interactions

### Database
- Complete schema with 8 core tables
- Automatic triggers for balance updates
- Materialized views for reporting
- ACID compliance

### Documentation
- README.md with quick start guide
- CLAUDE.md for development context
- Example configuration file
- Quick reference guide
- Complete roadmap

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
