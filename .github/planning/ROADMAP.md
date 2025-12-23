# FinTrack Implementation Roadmap

## Overview

This document outlines the implementation roadmap for FinTrack, a terminal-based personal finance tracking and budgeting application. The project is divided into 5 major phases, each delivering incrementally valuable features while following Unix philosophy and test-driven development principles.

**Target Users:** Developers, power users, privacy-conscious individuals who prefer local-first, CLI-based tools

**Core Values:**
- Privacy (local-first, no cloud dependencies)
- Simplicity (Unix philosophy, composable commands)
- Transparency (PostgreSQL, plain text output, open source)
- Portability (single binary, cross-platform, data export)

---

## Project Milestones

### Phase 1: Core Foundation (MVP)
**Timeline:** 4-6 weeks
**Goal:** Basic transaction tracking and account management

#### Features
- [x] Project setup and structure
- [ ] Database schema and migrations
- [ ] Account CRUD operations
- [ ] Transaction CRUD operations
- [ ] Category management (hierarchical)
- [ ] CSV import (generic format)
- [ ] Basic reporting (income statement, spending by category)
- [ ] Configuration management
- [ ] CLI framework (Cobra)

#### Deliverables
1. Working binary (`fintrack`)
2. Core commands: `account`, `tx`, `import`, `report`, `config`
3. PostgreSQL database with triggers
4. Test suite (>80% coverage)
5. Basic documentation

#### Success Criteria
- [ ] Can create accounts and add transactions
- [ ] Can import 1000+ transactions from CSV in <5 seconds
- [ ] Account balances update automatically
- [ ] Reports show accurate income/expense summaries
- [ ] All commands support `--json` output
- [ ] Works on Linux, macOS, Windows

#### Technical Tasks
1. **Project Scaffold**
   - Initialize Go module
   - Set up Cobra CLI framework
   - Create directory structure
   - Configure build scripts

2. **Database Layer**
   - Implement schema (fintrack_schema.sql)
   - Create migration system
   - Set up GORM/sqlx connections
   - Implement repository pattern
   - Add balance update triggers

3. **Core Commands**
   - `fintrack account`: CRUD operations
   - `fintrack tx`: Transaction management
   - `fintrack import csv`: CSV import with duplicate detection
   - `fintrack report`: Income and spending reports
   - `fintrack config`: Configuration management

4. **Testing**
   - Unit tests for business logic
   - Integration tests for database operations
   - CSV import test cases
   - Cross-platform testing (Linux, macOS, Windows)

5. **Documentation**
   - Installation guide
   - Quick start tutorial
   - Command reference
   - CSV import guide

---

### Phase 2: Budgeting & Scheduling
**Timeline:** 3-4 weeks
**Goal:** Budget tracking and recurring transactions

#### Features
- [ ] Budget creation and monitoring
- [ ] Budget alerts and thresholds
- [ ] Recurring transaction templates
- [ ] Manual recurring transaction generation
- [ ] Reminder system (terminal-based)
- [ ] Calendar view (monthly/weekly)

#### Deliverables
1. Commands: `budget`, `schedule`, `remind`, `cal`
2. Budget alert system
3. Recurring transaction engine
4. Calendar visualization (ASCII)

#### Success Criteria
- [ ] Can create monthly budgets per category
- [ ] Alerts show when 80% of budget spent
- [ ] Can schedule recurring bills/income
- [ ] Calendar shows upcoming financial events
- [ ] Reminders appear on command invocation

#### Technical Tasks
1. **Budget System**
   - Budget CRUD operations
   - Spending calculation against budgets
   - Alert threshold checking
   - Rollover logic (optional)

2. **Recurring Transactions**
   - Frequency calculation (daily, weekly, monthly, etc.)
   - Next occurrence algorithm
   - Generate transactions command
   - Link generated transactions to templates

3. **Reminder Engine**
   - Reminder creation from recurring items
   - Budget threshold reminders
   - Custom reminders
   - Dismissal logic

4. **Calendar View**
   - ASCII calendar rendering
   - Event overlay (transactions, bills, reminders)
   - Color coding by type
   - Upcoming events list

---

### Phase 3: Projections & Analytics
**Timeline:** 4-5 weeks
**Goal:** Cash flow forecasting and insights

#### Features
- [ ] Cash flow projection engine
- [ ] Confidence scoring algorithm
- [ ] Scenario modeling (conservative/moderate/optimistic)
- [ ] Spending trend analysis
- [ ] Savings rate calculation
- [ ] Net worth tracking over time
- [ ] Category spending trends

#### Deliverables
1. Command: `project`, enhanced `report`
2. Projection algorithm with confidence scores
3. Trend analysis reports
4. Statistics dashboard

#### Success Criteria
- [ ] 90-day projection generates in <500ms
- [ ] Projections achieve 80%+ accuracy against actuals
- [ ] Confidence scores reflect data quality
- [ ] Trend reports identify spending patterns
- [ ] Savings rate calculated correctly

#### Technical Tasks
1. **Projection Engine**
   - Historical spending analysis
   - Category average calculations
   - Recurring item integration
   - Confidence scoring formula
   - Scenario implementation (percentile-based)

2. **Analytics**
   - Monthly/quarterly/annual aggregations
   - Category trend analysis
   - Savings rate calculation
   - Net worth tracking
   - Spending velocity metrics

3. **Optimization**
   - Materialized view refreshes
   - Query performance tuning
   - Index optimization
   - Caching strategies

---

### Phase 4: Advanced Features
**Timeline:** 5-6 weeks
**Goal:** Automation and integrations

#### Features
- [ ] Auto-generate recurring transactions (cron job)
- [ ] Bank-specific CSV mappings (Chase, BoA, Amex, etc.)
- [ ] Desktop notifications (Linux/macOS)
- [ ] Multi-currency support
- [ ] Split transactions (one transaction, multiple categories)
- [ ] Reconciliation workflow
- [ ] Transaction rules (auto-categorization)
- [ ] Data backup/restore commands
- [ ] Plaid connector (optional)

#### Deliverables
1. Automation scripts (cron/systemd)
2. Bank CSV mappings for top 10 banks
3. Desktop notification integration
4. Multi-currency conversion
5. Advanced transaction features
6. Plaid integration (opt-in)

#### Success Criteria
- [ ] Recurring transactions generate automatically
- [ ] Import works for 10+ major banks
- [ ] Desktop notifications on Linux and macOS
- [ ] Multi-currency accounts calculate correctly
- [ ] Split transactions allocate to multiple categories
- [ ] Plaid syncs transactions automatically (if enabled)

#### Technical Tasks
1. **Automation**
   - Systemd service for recurring transaction generation
   - Cron job examples
   - Reminder notification daemon
   - Materialized view refresh scheduler

2. **Import Enhancements**
   - Bank-specific CSV parsers
   - Format auto-detection
   - Intelligent duplicate detection
   - Category mapping rules

3. **Multi-Currency**
   - Exchange rate API integration
   - Currency conversion logic
   - Multi-currency reporting
   - Historical rate storage

4. **Split Transactions**
   - Schema update (transaction_splits table)
   - UI for split entry
   - Report aggregation updates

5. **Plaid Integration** (Optional)
   - Plaid API client
   - Account linking flow
   - Transaction sync
   - Balance reconciliation

---

### Phase 5: Polish & Optimization
**Timeline:** 3-4 weeks
**Goal:** Production-ready, performant, delightful UX

#### Features
- [ ] Query optimization and indexing review
- [ ] Materialized view refresh strategies
- [ ] Interactive TUI mode (Bubble Tea framework)
- [ ] Colorized output themes
- [ ] Shell completion (bash, zsh, fish)
- [ ] Comprehensive documentation
- [ ] Performance benchmarks
- [ ] Security audit
- [ ] Professional packaging (homebrew, apt, etc.)

#### Deliverables
1. TUI mode (`fintrack tui`)
2. Shell completion scripts
3. Performance benchmarks
4. Full documentation site
5. Installation packages
6. 1.0 release

#### Success Criteria
- [ ] All queries <100ms with 100k+ transactions
- [ ] TUI mode is responsive and intuitive
- [ ] Shell completion works in bash/zsh/fish
- [ ] Documentation is comprehensive
- [ ] Installation is one command on major platforms
- [ ] Security vulnerabilities addressed

#### Technical Tasks
1. **TUI Development**
   - Bubble Tea framework integration
   - Interactive account/transaction browser
   - Budget gauge visualization
   - Calendar navigation
   - Search and filter UI

2. **Performance**
   - Query profiling and optimization
   - Index review
   - Connection pooling tuning
   - Memory profiling
   - Load testing (10k, 100k, 1M transactions)

3. **Shell Integration**
   - Bash completion
   - Zsh completion
   - Fish completion
   - Man page generation

4. **Documentation**
   - User guide (comprehensive)
   - Developer documentation
   - API reference (if applicable)
   - Video tutorials (optional)
   - Example workflows

5. **Packaging**
   - Homebrew formula
   - Debian package (.deb)
   - RPM package
   - AUR package (Arch Linux)
   - Chocolatey package (Windows)
   - Installation scripts

6. **Security**
   - Input validation audit
   - SQL injection testing
   - Dependency vulnerability scan
   - Credential storage review
   - Data encryption (optional)

---

## Development Workflow

### Test-Driven Development (TDD)

**All features MUST follow TDD:**

1. **Red**: Write failing test
2. **Green**: Implement minimum code to pass
3. **Refactor**: Improve code quality

**Example workflow:**
```bash
# 1. Write test
cat > tests/unit/account_test.go <<EOF
func TestCreateAccount(t *testing.T) {
    account := CreateAccount("Test Account", "checking", 1000)
    assert.Equal(t, "Test Account", account.Name)
    assert.Equal(t, 1000.0, account.Balance)
}
EOF

# 2. Run test (should fail)
go test ./tests/unit/account_test.go
# FAIL: undefined: CreateAccount

# 3. Implement feature
cat > internal/core/accounts.go <<EOF
func CreateAccount(name, accountType string, balance float64) *Account {
    return &Account{Name: name, Type: accountType, Balance: balance}
}
EOF

# 4. Run test (should pass)
go test ./tests/unit/account_test.go
# PASS

# 5. Refactor and ensure tests still pass
```

### Git Workflow

**Branch strategy:**
- `main`: Stable releases only
- `develop`: Integration branch
- `feature/account-management`: Feature branches
- `fix/balance-calculation-bug`: Bug fix branches

**Commit conventions:**
```
feat(account): Add account creation command
fix(transaction): Correct balance update trigger
docs(readme): Update installation instructions
test(budget): Add budget alert test cases
refactor(core): Extract transaction logic to service
```

### Code Review Checklist

- [ ] All tests pass (`go test ./...`)
- [ ] Test coverage >80% for new code
- [ ] No hardcoded credentials or secrets
- [ ] Error handling is comprehensive
- [ ] Input validation is present
- [ ] Documentation is updated
- [ ] CHANGELOG.md is updated
- [ ] Performance is acceptable (<100ms for reads)

---

## Technology Decisions

### Core Stack

| Component | Technology | Rationale |
|-----------|-----------|-----------|
| Language | Go 1.21+ | Single binary, cross-platform, performance |
| Database | PostgreSQL 12+ | ACID compliance, rich features, self-hostable |
| CLI Framework | Cobra | Industry standard, auto-help, subcommands |
| Config | Viper | YAML/ENV support, hot reload |
| Testing | Testify | Assertions, mocking, suite support |
| ORM | GORM or sqlx | GORM for ease, sqlx for performance |
| TUI | Bubble Tea | Modern, composable, delightful UX |

### Database Choice: PostgreSQL

**Why PostgreSQL over SQLite?**
- Better concurrency (multiple processes)
- Advanced features (arrays, JSONB, triggers)
- Materialized views for performance
- Suitable for future API server mode

**Why PostgreSQL over MySQL?**
- Better standards compliance
- Superior JSON support
- More permissive license (PostgreSQL License)

### Go vs Python/Ruby/Node

**Why Go?**
- Single binary distribution (no runtime)
- Fast startup (<10ms)
- Strong typing (financial safety)
- Excellent concurrency (future features)
- Native cross-compilation

---

## Testing Strategy

### Test Pyramid

```
         /\
        /  \       E2E Tests (5%)
       /    \      - Full workflow tests
      /------\     - Cross-platform binaries
     /        \
    /  INTEG.  \   Integration Tests (25%)
   /------------\  - Database operations
  /              \ - CSV import end-to-end
 /   UNIT TESTS  \ Unit Tests (70%)
/________________\ - Business logic
                    - Calculations
                    - Validation
```

### Test Coverage Goals

- **Unit Tests**: 80%+ coverage
- **Integration Tests**: All CRUD operations
- **E2E Tests**: Critical workflows (add transaction, import CSV, generate projection)

### Cross-Platform Testing

**Platforms:**
- Linux (Ubuntu 22.04, Arch)
- macOS (Intel, Apple Silicon)
- Windows (10, 11 via WSL and native)

**Test matrix:**
```yaml
platform:
  - linux-amd64
  - linux-arm64
  - darwin-amd64
  - darwin-arm64
  - windows-amd64

postgres_version:
  - 12
  - 13
  - 14
  - 15
```

---

## Documentation Plan

### User Documentation

1. **README.md** - Quick overview and installation
2. **INSTALLATION.md** - Detailed installation for all platforms
3. **QUICKSTART.md** - 5-minute tutorial
4. **USAGE.md** - Comprehensive command reference (from FINTRACK_QUICKREF.md)
5. **IMPORT_GUIDE.md** - CSV import for major banks
6. **FAQ.md** - Common questions and troubleshooting

### Developer Documentation

1. **DEVELOPMENT.md** - Contributing guide, architecture
2. **DATABASE.md** - Schema documentation (from fintrack_schema.sql)
3. **TESTING.md** - Testing strategy and guidelines
4. **ROADMAP.md** - This file
5. **CHANGELOG.md** - Version history

### API Documentation (Future)

- OpenAPI/Swagger spec for REST API
- Authentication guide
- Client library examples (Python, JavaScript)

---

## Performance Targets

### Response Time

| Operation | Target | Measured |
|-----------|--------|----------|
| Add transaction | <50ms | TBD |
| List 100 transactions | <100ms | TBD |
| Generate 90-day projection | <500ms | TBD |
| Import 1000 CSV rows | <2s | TBD |
| Monthly report | <200ms | TBD |
| Database query (indexed) | <50ms | TBD |

### Scalability

| Metric | Target |
|--------|--------|
| Transactions | 1M+ without degradation |
| Accounts | 100+ |
| Categories | 500+ |
| Concurrent users | 10+ (API mode) |

### Resource Usage

| Resource | Target |
|----------|--------|
| Binary size | <15MB |
| Memory usage | <50MB (typical) |
| Database size | <100MB per year of data |
| Startup time | <100ms |

---

## Security Considerations

### Data Protection

1. **Database Credentials**
   - Use environment variables
   - Support OS keychain integration (future)
   - Never log passwords

2. **Input Validation**
   - Sanitize all user inputs
   - Use parameterized queries (prevent SQL injection)
   - Validate file paths (prevent directory traversal)
   - Limit CSV file sizes

3. **Data Encryption** (Future)
   - Encrypt sensitive fields (account numbers)
   - Use AES-256-GCM
   - Store encryption key in OS keychain

### Audit Trail

- All transactions are immutable (no updates, only reversals)
- Track import history (prevent duplicate imports)
- Log all delete operations

---

## Migration Path for Existing Users

### From Mint

1. Download Mint export CSV
2. Create category mapping file
3. Import: `fintrack import csv mint.csv --format mint`

### From YNAB

1. Export budget and transactions from YNAB
2. Convert categories to FinTrack format
3. Import: `fintrack import csv ynab.csv --format ynab`

### From GnuCash

1. Export to CSV from GnuCash
2. Map accounts to FinTrack accounts
3. Import: `fintrack import csv gnucash.csv --format gnucash`

### From Excel/Sheets

1. Format data to generic CSV:
   ```csv
   date,amount,payee,category,description
   2025-11-15,-100.00,Store,Groceries,Weekly shopping
   ```
2. Import: `fintrack import csv data.csv --format generic`

---

## Future Enhancements (Post-1.0)

### API Server Mode

- REST API for mobile/web clients
- WebSocket for real-time updates
- JWT authentication
- Rate limiting

### Mobile App

- React Native or Flutter
- Quick expense entry
- Receipt OCR
- Push notifications

### Machine Learning

- Auto-categorization based on payee/description
- Anomaly detection (unusual spending)
- Budget recommendations
- Income/expense forecasting

### Investment Tracking

- Portfolio accounts
- Stock/crypto price integration
- Capital gains calculations
- Asset allocation views

### Multi-User / Shared Accounts

- Household budgets
- User roles (admin, viewer)
- Separate views per user
- Shared vs personal accounts

---

## Release Strategy

### Version Numbering

Follow [Semantic Versioning](https://semver.org/):
- `MAJOR.MINOR.PATCH`
- `1.0.0` = Initial stable release
- `1.1.0` = New features (backward compatible)
- `1.0.1` = Bug fixes
- `2.0.0` = Breaking changes

### Release Checklist

- [ ] All tests pass
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Version bumped in code
- [ ] Git tag created
- [ ] Binaries built for all platforms
- [ ] Release notes written
- [ ] GitHub release created
- [ ] Homebrew formula updated
- [ ] Package repositories updated

### Support Policy

- **Latest major version**: Full support (features + bug fixes)
- **Previous major version**: Security fixes only (6 months)
- **Older versions**: No support

---

## Success Metrics

### User Adoption

- GitHub stars: 1000+ (first year)
- Active users: 5000+ (first year)
- Community contributions: 10+ contributors

### Technical Quality

- Test coverage: >80%
- Bug reports: <5 open critical bugs
- Performance: All targets met
- Documentation completeness: 100%

### Community Health

- Issue response time: <48 hours
- PR review time: <1 week
- Monthly releases (during active development)

---

## Risks and Mitigations

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| PostgreSQL complexity scares users | High | Medium | Provide Docker Compose setup, one-command install |
| Performance issues with large datasets | Medium | Low | Early optimization, materialized views, benchmarks |
| Lack of user adoption | High | Medium | Focus on developer community, Reddit/HN launch |
| Security vulnerability | High | Low | Regular dependency updates, security audit |
| Feature creep | Medium | High | Stick to roadmap, defer nice-to-haves to post-1.0 |
| Cross-platform issues | Medium | Medium | CI/CD testing on all platforms, early user testing |

---

## Contributing

### How to Contribute

1. **Check open issues** for "good first issue" labels
2. **Comment on issue** to claim it
3. **Fork repository** and create feature branch
4. **Write tests first** (TDD)
5. **Implement feature** with tests passing
6. **Submit PR** with description and test evidence
7. **Respond to review** feedback

### Contribution Areas

- **Code**: Implement features, fix bugs
- **Documentation**: Improve guides, fix typos
- **Testing**: Add test cases, report bugs
- **Design**: UI/UX improvements for TUI mode
- **Localization**: Translate to other languages (future)

---

## Resources

### Documentation

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Cobra CLI Framework](https://cobra.dev/)
- [Bubble Tea TUI](https://github.com/charmbracelet/bubbletea)
- [Go Best Practices](https://go.dev/doc/effective_go)

### Related Projects

- [ledger-cli](https://www.ledger-cli.org/) - Inspiration for plain-text accounting
- [hledger](https://hledger.org/) - Haskell-based accounting tool
- [YNAB](https://www.youneedabudget.com/) - Budgeting philosophy
- [Firefly III](https://www.firefly-iii.org/) - Web-based personal finance

---

## Conclusion

FinTrack aims to be the terminal-based personal finance tool for developers and power users who value:
- **Privacy** (local-first, no cloud)
- **Simplicity** (Unix philosophy, composable)
- **Transparency** (open source, plain text)
- **Control** (full data ownership, scriptable)

By following this roadmap and adhering to test-driven development, we can deliver a production-ready 1.0 release in approximately **6 months** with a small team (2-3 developers).

---

**Document Version:** 1.0
**Last Updated:** 2025-11-16
**Next Review:** Start of each phase
