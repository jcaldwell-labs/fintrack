# Personal Finance Tracker - System Design Plan

## Executive Summary

A terminal-based personal financial tracking and budgeting application following Unix philosophy principles. PostgreSQL-backed, command-oriented CLI with support for transaction tracking, budgeting, recurring items, reminders, and cash flow projection.

**Project Name:** `fintrack` (working name)

---

## Design Principles

### Unix Philosophy
1. **Do One Thing Well**: Each command focuses on a single responsibility
2. **Text I/O**: Human-readable output, machine-readable with `--json` flag
3. **Composable**: Commands can be chained via pipes and shell scripts
4. **Minimal Surface Area**: ~10-12 core commands with single-letter aliases
5. **Configuration as Code**: Database connection, preferences in plain text config

### User Experience
- **Fast**: Sub-100ms response for most commands
- **Intuitive**: Natural language-style commands (`fintrack tx add -100 "Groceries"`)
- **Offline-First**: Local PostgreSQL, no cloud dependencies
- **Portable**: CSV export/import for data migration
- **Privacy**: All data stays local

---

## Architecture Overview

### Technology Stack
```
┌─────────────────────────────────────┐
│   CLI Interface (Cobra)             │
│   - Command parsing                 │
│   - Flag validation                 │
│   - Output formatting               │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   Business Logic Layer              │
│   - Account management              │
│   - Transaction processing          │
│   - Budget calculations             │
│   - Projection engine               │
│   - Reminder scheduler              │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   Data Access Layer (GORM/sqlx)     │
│   - Repository pattern              │
│   - Transaction management          │
│   - Query optimization              │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   PostgreSQL Database               │
│   - Transactional integrity         │
│   - Constraints & indexes           │
│   - Materialized views for reports  │
└─────────────────────────────────────┘
```

### Directory Structure
```
fintrack/
├── cmd/fintrack/           # CLI entry point
│   └── main.go
├── internal/
│   ├── commands/           # Command implementations
│   │   ├── account.go
│   │   ├── transaction.go
│   │   ├── budget.go
│   │   ├── import.go
│   │   ├── report.go
│   │   ├── schedule.go
│   │   ├── remind.go
│   │   ├── project.go
│   │   └── calendar.go
│   ├── core/              # Business logic
│   │   ├── accounts.go
│   │   ├── transactions.go
│   │   ├── budgets.go
│   │   ├── recurring.go
│   │   ├── reminders.go
│   │   └── projections.go
│   ├── models/            # Data structures
│   │   └── models.go
│   ├── db/               # Database layer
│   │   ├── connection.go
│   │   ├── migrations/
│   │   └── repositories/
│   ├── importers/        # CSV and connector logic
│   │   ├── csv.go
│   │   └── plaid.go      # Future: Plaid connector
│   ├── output/           # Output formatters
│   │   ├── table.go
│   │   ├── json.go
│   │   └── calendar.go
│   └── config/           # Configuration management
│       └── config.go
├── tests/
│   ├── integration/
│   └── unit/
├── migrations/           # SQL migration files
│   ├── 001_initial_schema.sql
│   ├── 002_add_recurring.sql
│   └── ...
└── config.example.yaml   # Example configuration
```

---

## Database Schema

### Core Tables

#### accounts
```sql
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,  -- checking, savings, credit, cash, investment
    currency VARCHAR(3) DEFAULT 'USD',
    initial_balance DECIMAL(15,2) DEFAULT 0,
    current_balance DECIMAL(15,2) DEFAULT 0,
    institution VARCHAR(255),
    account_number_last4 VARCHAR(4),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_accounts_active ON accounts(is_active);
CREATE INDEX idx_accounts_type ON accounts(type);
```

#### categories
```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id INTEGER REFERENCES categories(id),
    type VARCHAR(20) NOT NULL,  -- income, expense, transfer
    color VARCHAR(7),  -- Hex color for visualization
    icon VARCHAR(50),  -- Icon identifier for TUI
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_categories_parent ON categories(parent_id);
CREATE INDEX idx_categories_type ON categories(type);

-- Hierarchical categories: Groceries > Food & Dining
```

#### transactions
```sql
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    amount DECIMAL(15,2) NOT NULL,  -- Negative for expenses, positive for income
    category_id INTEGER REFERENCES categories(id),
    payee VARCHAR(255),
    description TEXT,
    type VARCHAR(20) NOT NULL,  -- income, expense, transfer
    transfer_account_id INTEGER REFERENCES accounts(id),  -- For transfers
    recurring_id INTEGER REFERENCES recurring_items(id),  -- Link to recurring template
    tags TEXT[],  -- PostgreSQL array for tags
    is_reconciled BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_transactions_account ON transactions(account_id);
CREATE INDEX idx_transactions_date ON transactions(date DESC);
CREATE INDEX idx_transactions_category ON transactions(category_id);
CREATE INDEX idx_transactions_recurring ON transactions(recurring_id);
CREATE INDEX idx_transactions_tags ON transactions USING GIN(tags);
```

#### budgets
```sql
CREATE TABLE budgets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category_id INTEGER REFERENCES categories(id),
    period_type VARCHAR(20) NOT NULL,  -- monthly, quarterly, annual
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    limit_amount DECIMAL(15,2) NOT NULL,
    rollover_enabled BOOLEAN DEFAULT false,
    alert_threshold DECIMAL(5,2) DEFAULT 0.80,  -- Alert at 80% of limit
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_budgets_period ON budgets(period_start, period_end);
CREATE INDEX idx_budgets_category ON budgets(category_id);
```

#### recurring_items
```sql
CREATE TABLE recurring_items (
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL REFERENCES accounts(id),
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    category_id INTEGER REFERENCES categories(id),
    description TEXT,
    frequency VARCHAR(20) NOT NULL,  -- daily, weekly, biweekly, monthly, quarterly, annual
    frequency_interval INTEGER DEFAULT 1,  -- Every N periods
    start_date DATE NOT NULL,
    end_date DATE,
    next_date DATE NOT NULL,
    last_generated_date DATE,
    auto_generate BOOLEAN DEFAULT false,  -- Auto-create transactions
    reminder_days_before INTEGER DEFAULT 3,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_recurring_next_date ON recurring_items(next_date) WHERE is_active = true;
CREATE INDEX idx_recurring_account ON recurring_items(account_id);
```

#### reminders
```sql
CREATE TABLE reminders (
    id SERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL,  -- transaction, budget, bill, custom
    related_id INTEGER,  -- ID of related entity (recurring_item, budget, etc.)
    title VARCHAR(255) NOT NULL,
    message TEXT,
    remind_date DATE NOT NULL,
    remind_time TIME,
    is_dismissed BOOLEAN DEFAULT false,
    dismissed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_reminders_date ON reminders(remind_date, remind_time) WHERE is_dismissed = false;
CREATE INDEX idx_reminders_type ON reminders(type);
```

#### cash_flow_projections
```sql
CREATE TABLE cash_flow_projections (
    id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts(id),  -- NULL for all accounts
    projection_date DATE NOT NULL,
    projected_balance DECIMAL(15,2) NOT NULL,
    confidence_level DECIMAL(5,2),  -- 0.0 to 1.0
    projection_type VARCHAR(20),  -- conservative, moderate, optimistic
    generated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_projections_account_date ON cash_flow_projections(account_id, projection_date);
```

#### import_history
```sql
CREATE TABLE import_history (
    id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts(id),
    filename VARCHAR(255),
    imported_at TIMESTAMP DEFAULT NOW(),
    records_imported INTEGER,
    records_skipped INTEGER,
    records_failed INTEGER,
    import_hash VARCHAR(64),  -- SHA256 of file to prevent duplicates
    error_log TEXT
);
```

### Materialized Views (for performance)

#### monthly_spending_by_category
```sql
CREATE MATERIALIZED VIEW monthly_spending_by_category AS
SELECT
    DATE_TRUNC('month', date) as month,
    category_id,
    c.name as category_name,
    SUM(amount) as total_amount,
    COUNT(*) as transaction_count
FROM transactions t
JOIN categories c ON t.category_id = c.id
WHERE type = 'expense'
GROUP BY DATE_TRUNC('month', date), category_id, c.name;

CREATE INDEX idx_monthly_spending_month ON monthly_spending_by_category(month);
```

---

## Command Interface Design

### Core Commands

#### 1. `fintrack account` (alias: `a`)
Manage financial accounts

```bash
# List all accounts
fintrack account list
fintrack a ls

# Add new account
fintrack account add "Chase Checking" --type checking --balance 5000
fintrack a add "Amex Gold" -t credit -b -1200

# Show account details with recent transactions
fintrack account show "Chase Checking"
fintrack a show 1  # By ID

# Update account
fintrack account update 1 --name "Chase Premier Checking"

# Close/archive account
fintrack account close 1

# Balance reconciliation
fintrack account reconcile 1 --balance 5432.10

# Output formats
fintrack account list --json
fintrack account list --format table  # Default
```

**Output Example:**
```
ID  NAME                 TYPE       BALANCE     LAST ACTIVITY
1   Chase Checking       checking   $5,234.10   2025-11-15
2   Amex Gold           credit     -$1,234.00   2025-11-14
3   Ally Savings        savings    $15,000.00   2025-11-10
```

#### 2. `fintrack tx` (alias: `t`)
Transaction operations

```bash
# Add transaction
fintrack tx add -100 "Whole Foods" --account 1 --category groceries --date 2025-11-15
fintrack t add -100 "Whole Foods" -a "Chase Checking" -c groceries -d today

# Quick entry (uses default account)
fintrack tx -50 "Coffee"

# Add income
fintrack tx add 3500 "Salary" --category salary --account 1

# Transfer between accounts
fintrack tx transfer 500 --from 1 --to 3 --description "Savings"

# List transactions
fintrack tx list --account 1 --limit 20
fintrack tx list --from 2025-11-01 --to 2025-11-30
fintrack tx list --category groceries --month 2025-11

# Search transactions
fintrack tx search "whole foods"
fintrack tx search --payee "Amazon" --min 50 --max 200

# Edit transaction
fintrack tx edit 123 --amount -105.50 --category "dining"

# Delete transaction
fintrack tx delete 123

# Tag transactions
fintrack tx tag 123 "business-expense" "reimbursable"

# Import from CSV
fintrack import csv transactions.csv --account 1 --format chase
fintrack import csv --account 2 --format generic --map config.yaml
```

**Output Example:**
```
DATE        AMOUNT      CATEGORY     PAYEE           DESCRIPTION
2025-11-15  -$100.00    Groceries    Whole Foods     Weekly shopping
2025-11-14  -$45.23     Dining       Chipotle        Lunch
2025-11-13  -$3.50      Coffee       Starbucks       Morning coffee
2025-11-10  $3,500.00   Salary       Acme Corp       Biweekly paycheck

Total: -$148.73 | Count: 4
```

#### 3. `fintrack budget` (alias: `b`)
Budget management

```bash
# Create budget
fintrack budget create "Groceries" --category groceries --limit 600 --period monthly
fintrack b create "Dining Out" -c dining -l 300 -p monthly --alert 0.75

# List budgets
fintrack budget list
fintrack budget list --active

# Show budget status
fintrack budget show groceries
fintrack budget status --month 2025-11

# Update budget
fintrack budget update groceries --limit 700

# Budget performance
fintrack budget report --month 2025-11
```

**Output Example:**
```
CATEGORY     BUDGET      SPENT       REMAINING   % USED  STATUS
Groceries    $600.00     $423.45     $176.55     71%     ●●●●●○
Dining       $300.00     $267.89     $32.11      89%     ●●●●●● ⚠
Transport    $200.00     $145.00     $55.00      73%     ●●●●●○
Utilities    $150.00     $148.23     $1.77       99%     ●●●●●● ⚠

Total: $1,250.00 | Spent: $984.57 | Remaining: $265.43
```

#### 4. `fintrack schedule` (alias: `s`)
Manage recurring transactions

```bash
# Create recurring item
fintrack schedule add "Rent" -1500 --category rent --frequency monthly --day 1
fintrack s add "Netflix" -15.99 -c subscriptions -f monthly --account 2

# List recurring items
fintrack schedule list
fintrack s ls --active
fintrack s ls --upcoming 30  # Next 30 days

# Show upcoming recurring transactions
fintrack schedule upcoming --days 14

# Generate transactions from recurring items
fintrack schedule generate --dry-run
fintrack schedule generate --confirm

# Pause/resume recurring item
fintrack schedule pause 3
fintrack schedule resume 3

# Delete recurring item
fintrack schedule delete 3
```

**Output Example:**
```
ID  NAME         AMOUNT      FREQUENCY  NEXT DATE   AUTO  ACCOUNT
1   Rent         -$1,500.00  monthly    2025-12-01  yes   Checking
2   Netflix      -$15.99     monthly    2025-11-20  yes   Amex Gold
3   Salary       $3,500.00   biweekly   2025-11-22  no    Checking
4   Electric     -$120.00    monthly    2025-11-25  no    Checking

Next 7 days: 2 items | Total: -$135.99
```

#### 5. `fintrack remind` (alias: `r`)
Reminder and alert system

```bash
# Create custom reminder
fintrack remind add "Pay credit card" --date 2025-11-20 --time 09:00

# List upcoming reminders
fintrack remind list
fintrack remind list --today
fintrack remind list --week

# Dismiss reminder
fintrack remind dismiss 5

# Show overdue reminders
fintrack remind overdue

# Configure budget alerts
fintrack remind budget-alert groceries --threshold 0.80
```

**Output Example:**
```
DATE        TIME   TYPE     MESSAGE                           STATUS
2025-11-16  09:00  bill     Pay Amex Gold credit card         pending
2025-11-18  --     budget   Groceries budget at 85%           pending
2025-11-20  10:00  custom   Review investment portfolio       pending
2025-11-15  --     bill     Electric bill due                 overdue ⚠

Total: 4 reminders | Overdue: 1
```

#### 6. `fintrack project` (alias: `p`)
Cash flow projection

```bash
# Generate projection
fintrack project --days 90
fintrack project --months 6 --type conservative

# Show projection by account
fintrack project --account "Chase Checking" --days 30

# Projection summary
fintrack project summary --weeks 4

# Export projection data
fintrack project --days 90 --json > projection.json
```

**Output Example:**
```
CASH FLOW PROJECTION (90 days, moderate scenario)

DATE        INCOME      EXPENSES    NET FLOW    BALANCE     TREND
2025-11-16  $0.00       -$45.00     -$45.00     $5,189.10   ━
2025-11-20  $0.00       -$135.99    -$135.99    $5,053.11   ━
2025-11-22  $3,500.00   -$120.00    $3,380.00   $8,433.11   ━━━━
2025-12-01  $0.00       -$1,500.00  -$1,500.00  $6,933.11   ━━
2025-12-06  $3,500.00   -$580.00    $2,920.00   $9,853.11   ━━━━
...

Summary:
  Total Income:    $14,000.00
  Total Expenses:  -$8,234.50
  Net Change:      +$5,765.50
  Ending Balance:  $10,999.60

Confidence: 82% (based on 3 months of transaction history)
```

#### 7. `fintrack report` (alias: `rp`)
Generate financial reports

```bash
# Income statement
fintrack report income --month 2025-11
fintrack report income --year 2025

# Spending by category
fintrack report spending --month 2025-11
fintrack report spending --category --chart

# Net worth over time
fintrack report networth --months 12

# Custom date range
fintrack report summary --from 2025-01-01 --to 2025-11-15

# Export report
fintrack report income --year 2025 --format csv > income_2025.csv
fintrack report spending --month 2025-11 --format json
```

**Output Example:**
```
INCOME STATEMENT - November 2025

INCOME
  Salary                    $7,000.00
  Freelance                 $1,200.00
  Interest                     $12.45
                          ───────────
  Total Income              $8,212.45

EXPENSES
  Housing
    Rent                   -$1,500.00
    Utilities                -$145.00
  Food & Dining
    Groceries                -$523.45
    Dining Out               -$287.00
  Transportation             -$245.00
  Entertainment              -$156.00
  Other                      -$234.89
                          ───────────
  Total Expenses           -$3,091.34

NET INCOME                  $5,121.11
```

#### 8. `fintrack cal` (alias: `c`)
Calendar view of financial events

```bash
# Show calendar with transactions
fintrack cal
fintrack cal --month 2025-12

# Show specific types
fintrack cal --type bills
fintrack cal --type income
fintrack cal --category groceries

# Upcoming events
fintrack cal next --days 7
```

**Output Example:**
```
            November 2025
Su  Mo  Tu  We  Th  Fr  Sa
                        1  💰 Rent -$1,500
 2   3   4   5   6   7   8  💰 Salary +$3,500
 9  10  11  12  13  14  15  🛒 Groceries -$100
16  17  18  19  20  21  22  💰 Salary +$3,500, 📺 Netflix -$16
23  24  25  26  27  28  29  ⚡ Electric -$120
30

Upcoming (next 7 days):
  Nov 20: Netflix subscription -$15.99
  Nov 22: Salary deposit +$3,500.00
  Nov 25: Electric bill (est.) -$120.00
```

#### 9. `fintrack config`
Configuration management

```bash
# Show current configuration
fintrack config show

# Set database connection
fintrack config set db.url "postgresql://user:pass@localhost:5432/fintrack"

# Set default account
fintrack config set defaults.account "Chase Checking"

# Set currency
fintrack config set defaults.currency USD

# Initialize database
fintrack config init-db

# Run migrations
fintrack config migrate
```

#### 10. `fintrack stats`
Quick statistics and insights

```bash
# Dashboard summary
fintrack stats

# Spending trends
fintrack stats trends --months 6

# Top categories
fintrack stats top-spending --month 2025-11

# Savings rate
fintrack stats savings-rate --year 2025
```

---

## Data Flow Examples

### Adding a Transaction

```
User: fintrack tx add -100 "Whole Foods" --category groceries

↓
commands.TransactionCmd parses flags
  - amount: -100
  - payee: "Whole Foods"
  - category: "groceries"
  - account: default (from config)
  - date: today

↓
core.AddTransaction():
  1. Validate account exists
  2. Resolve category name → ID
  3. Begin database transaction
  4. Insert transaction record
  5. Update account.current_balance
  6. Check budget status for category
  7. Create reminder if budget threshold exceeded
  8. Commit transaction

↓
output.PrintTransaction() displays:
  ✓ Added transaction #123
  Date: 2025-11-16
  Amount: -$100.00
  Category: Groceries
  Account: Chase Checking (new balance: $5,134.10)

  Budget Alert: Groceries budget is at 84% ($502/$600)
```

### Generating Cash Flow Projection

```
User: fintrack project --days 30

↓
commands.ProjectCmd initiates projection

↓
core.GenerateProjection():
  1. Fetch all active accounts and balances
  2. Query recurring_items for next 30 days
  3. Calculate historical spending averages by category
  4. For each future date:
     a. Add scheduled income/expenses
     b. Estimate variable expenses based on averages
     c. Calculate projected balance
     d. Assign confidence score
  5. Store projections in cash_flow_projections table

↓
output.PrintProjection() displays table:
  - Daily breakdown
  - Balance trends
  - Summary statistics
  - Confidence metrics
```

### Processing Recurring Transactions

```
Scheduled Job (cron/systemd timer):
fintrack schedule generate

↓
core.ProcessRecurringItems():
  1. Query recurring_items WHERE next_date <= TODAY
  2. For each item:
     a. Create transaction record
     b. Link to recurring_id
     c. Update recurring_item.next_date based on frequency
     d. Update recurring_item.last_generated_date
  3. Create reminders for items due in 3 days

↓
Output:
  Generated 3 transactions:
  - Rent: -$1,500.00
  - Netflix: -$15.99
  - Electric (estimated): -$120.00

  Created 2 reminders for upcoming bills
```

---

## CSV Import Design

### Generic CSV Format

```csv
date,amount,payee,category,description
2025-11-15,-100.00,Whole Foods,Groceries,Weekly shopping
2025-11-14,-45.23,Chipotle,Dining,Lunch
2025-11-13,3500.00,Acme Corp,Salary,Biweekly paycheck
```

### Bank-Specific Mappings

**Chase Format:**
```yaml
# chase_mapping.yaml
columns:
  date: "Posting Date"
  amount: "Amount"
  payee: "Description"
  category: "Type"  # Map to categories
date_format: "MM/DD/YYYY"
amount_sign: inverse  # Chase uses positive for debits
skip_rows: 1
```

**Import Command:**
```bash
fintrack import csv chase_transactions.csv \
  --account "Chase Checking" \
  --format chase \
  --map chase_mapping.yaml \
  --dry-run

# Review dry-run output
fintrack import csv chase_transactions.csv \
  --account "Chase Checking" \
  --format chase \
  --confirm
```

### Duplicate Detection

- Generate SHA256 hash of: `date + amount + payee`
- Check `import_history.import_hash` before inserting
- Skip duplicates, report in import summary

---

## Configuration File

**Location:** `~/.config/fintrack/config.yaml`

```yaml
# Database connection
database:
  url: "postgresql://localhost:5432/fintrack"
  # Or separate components:
  host: "localhost"
  port: 5432
  database: "fintrack"
  user: "fintrack_user"
  password: "${FINTRACK_DB_PASSWORD}"  # From environment
  sslmode: "disable"

# Default settings
defaults:
  account: "Chase Checking"  # Default for quick entries
  currency: "USD"
  date_format: "2006-01-02"  # Go time format

# Budget alerts
alerts:
  enabled: true
  threshold: 0.80  # Alert at 80% of budget
  methods:
    - terminal  # Print on next command
    - desktop   # Desktop notification (future)

# Recurring transaction automation
recurring:
  auto_generate: false  # Require manual confirmation
  generate_days_ahead: 3  # Generate transactions 3 days before due

# Projection settings
projection:
  default_days: 90
  confidence_threshold: 0.70  # Minimum confidence to show
  scenario: "moderate"  # conservative, moderate, optimistic

# Import settings
import:
  duplicate_window_days: 30  # Check for duplicates in last 30 days
  auto_categorize: true  # Use ML to suggest categories (future)

# Output preferences
output:
  default_format: "table"  # table, json, csv
  color: true
  unicode: true  # Use Unicode symbols
  timezone: "America/New_York"
```

**Environment Variables:**
- `FINTRACK_DB_URL`: Override database URL
- `FINTRACK_CONFIG`: Custom config file path
- `FINTRACK_DB_PASSWORD`: Database password

---

## Reminder & Alert System

### Reminder Types

1. **Bill Reminders**: Auto-created from recurring items
2. **Budget Alerts**: Triggered when spending exceeds threshold
3. **Custom Reminders**: User-defined one-time reminders
4. **Low Balance Warnings**: When account balance drops below threshold

### Delivery Methods (Phase 1)

- **Terminal**: Show on next `fintrack` command invocation
- **Check Command**: `fintrack remind list --today`

### Future Enhancements

- Desktop notifications (via `notify-send` on Linux, `osascript` on macOS)
- Email alerts (via SMTP configuration)
- Webhook integration (POST to custom URL)

---

## Calendar Integration

### Display Modes

1. **Month View**: Traditional calendar with financial events
2. **Week View**: Detailed weekly breakdown
3. **List View**: Upcoming events in chronological order

### Event Types

- 💰 Income (salary, deposits)
- 💳 Expenses (categorized by type)
- 🔄 Recurring items
- 📅 Bill due dates
- ⚠️ Budget alerts
- 📊 End-of-period summaries

### Calendar Data Sources

- Completed transactions
- Scheduled recurring items
- Budget period boundaries
- Custom reminders

---

## Cash Flow Projection Engine

### Projection Algorithm

```
For each future date D in projection range:
  1. Starting balance = previous_balance

  2. Add scheduled income/expenses:
     - Query recurring_items WHERE next_date = D
     - Add confirmed amounts

  3. Estimate variable expenses:
     - Calculate category averages from last 90 days
     - Prorate by days in period
     - Apply confidence factor

  4. Calculate ending balance:
     - ending_balance = starting_balance + income - expenses

  5. Assign confidence level:
     - High (>90%): Scheduled recurring items
     - Medium (70-90%): Based on consistent spending patterns
     - Low (<70%): Variable/unpredictable categories

  6. Store projection record
```

### Scenario Types

- **Conservative**: Use 75th percentile of historical spending
- **Moderate**: Use median (50th percentile)
- **Optimistic**: Use 25th percentile

### Confidence Scoring

```
Confidence = (
  0.4 * data_completeness_factor +
  0.3 * spending_consistency_factor +
  0.3 * time_decay_factor
)

Where:
  data_completeness = transaction_count / expected_count
  spending_consistency = 1 - (std_dev / mean)
  time_decay = exp(-days_ahead / 90)  # Exponential decay
```

---

## Testing Strategy

### Unit Tests

- Transaction validation logic
- Budget calculation accuracy
- Recurring item frequency calculations
- Date/time parsing and formatting
- Category hierarchy resolution

### Integration Tests

- Database CRUD operations
- Transaction rollback scenarios
- CSV import with various formats
- Projection generation end-to-end
- Reminder creation and dismissal

### Cross-Platform Testing

- PostgreSQL compatibility (versions 12+)
- Path handling (Windows vs Unix)
- Date/time zone handling
- Binary compilation (Linux, macOS, Windows)

### Test Data

- Sample transaction dataset (6 months)
- Multiple account types
- Hierarchical categories
- Edge cases: transfers, refunds, split transactions

---

## Development Phases

### Phase 1: Core Foundation (MVP)
**Goal:** Basic transaction tracking and account management

- [ ] Database schema and migrations
- [ ] Account CRUD operations
- [ ] Transaction CRUD operations
- [ ] Category management
- [ ] CSV import (generic format)
- [ ] Basic reporting (income statement, spending by category)
- [ ] Configuration management

**Deliverable:** Users can track transactions, categorize spending, import CSV files

### Phase 2: Budgeting & Scheduling
**Goal:** Budget tracking and recurring transactions

- [ ] Budget creation and monitoring
- [ ] Budget alerts and thresholds
- [ ] Recurring transaction templates
- [ ] Manual recurring transaction generation
- [ ] Reminder system (terminal-based)
- [ ] Calendar view

**Deliverable:** Users can set budgets, schedule recurring expenses, view financial calendar

### Phase 3: Projections & Analytics
**Goal:** Cash flow forecasting and insights

- [ ] Cash flow projection engine
- [ ] Confidence scoring algorithm
- [ ] Scenario modeling (conservative/moderate/optimistic)
- [ ] Spending trend analysis
- [ ] Savings rate calculation
- [ ] Net worth tracking over time

**Deliverable:** Users can forecast cash flow, analyze trends, understand financial health

### Phase 4: Advanced Features
**Goal:** Automation and integrations

- [ ] Auto-generate recurring transactions
- [ ] Bank-specific CSV mappings (Chase, BoA, etc.)
- [ ] Desktop notifications
- [ ] Multi-currency support
- [ ] Split transactions
- [ ] Reconciliation workflow
- [ ] Plaid connector (optional, requires API key)

**Deliverable:** Reduced manual entry, automated workflows, bank integrations

### Phase 5: Polish & Optimization
**Goal:** Performance and user experience

- [ ] Query optimization and indexing
- [ ] Materialized view refresh strategies
- [ ] Interactive TUI mode (bubble tea framework)
- [ ] Colorized output themes
- [ ] Shell completion (bash, zsh, fish)
- [ ] Comprehensive documentation
- [ ] Performance benchmarks

**Deliverable:** Production-ready, performant, delightful UX

---

## Technical Decisions

### Why PostgreSQL?

- **ACID compliance**: Critical for financial data integrity
- **Rich data types**: Arrays (tags), JSONB (metadata), date/time handling
- **Materialized views**: Pre-computed reports for performance
- **Constraint enforcement**: Foreign keys, check constraints
- **Wide deployment**: Easy to self-host, cloud options available

### Why Go?

- **Single binary**: No runtime dependencies, easy distribution
- **Cross-platform**: Compiles for Linux, macOS, Windows, ARM
- **Performance**: Fast startup, low memory footprint
- **Concurrency**: Goroutines for async operations (future)
- **Strong typing**: Compile-time safety for financial calculations
- **Rich ecosystem**: Cobra (CLI), GORM/sqlx (DB), viper (config)

### Why Cobra Framework?

- Industry standard for Go CLIs
- Built-in flag parsing and validation
- Automatic help generation
- Subcommand structure aligns with Unix philosophy
- Used by kubectl, Hugo, GitHub CLI

### Text Output vs TUI

**Phase 1-3**: Plain text output (tables, lists)
- Pipeable, scriptable, Unix-friendly
- Works over SSH, in scripts, cron jobs
- Minimal dependencies

**Phase 5**: Optional TUI mode with Bubble Tea
- Interactive account/transaction browsing
- Real-time budget gauges
- Calendar navigation
- Activated with `fintrack tui` or `--interactive` flag

---

## Security Considerations

### Database Credentials

- Never hardcode passwords
- Use environment variables or OS keychain
- Support PostgreSQL connection strings with credential masking

### Data Privacy

- All data stored locally (no cloud sync in Phase 1)
- No telemetry or analytics
- Clear data export/deletion workflows

### Input Validation

- Sanitize all user inputs (SQL injection prevention)
- Validate amounts (prevent arithmetic overflow)
- Date range validation
- File path sanitization for CSV imports

### CSV Import Safety

- Limit file size (prevent resource exhaustion)
- Validate CSV structure before processing
- Transaction-based imports (rollback on error)
- Duplicate detection (prevent accidental re-imports)

---

## Performance Targets

- **Command response**: <100ms for read operations
- **Transaction insertion**: <50ms including balance update
- **CSV import**: 1000 transactions/second
- **Projection generation**: <500ms for 90-day projection
- **Database queries**: All <100ms with proper indexing

### Optimization Strategies

1. **Indexes**: On frequently queried columns (date, account_id, category_id)
2. **Materialized views**: Pre-aggregate monthly spending
3. **Connection pooling**: Reuse database connections
4. **Prepared statements**: Avoid query re-compilation
5. **Batch operations**: Bulk CSV inserts

---

## Migration Path (for existing users)

### From Mint

1. Export transactions CSV from Mint
2. Map Mint categories → fintrack categories
3. Import with custom mapping: `fintrack import csv mint_export.csv --format mint`

### From YNAB (You Need A Budget)

1. Export budget and transactions
2. Convert YNAB categories to fintrack categories
3. Import transactions: `fintrack import csv ynab_export.csv --format ynab`

### From GnuCash

1. Export to CSV from GnuCash
2. Map accounts and categories
3. Import: `fintrack import csv gnucash_export.csv --format gnucash`

---

## Future Enhancements (Beyond Phase 5)

### API Server Mode

```bash
fintrack serve --port 8080
```

- REST API for mobile/web frontends
- JWT authentication
- Read/write endpoints
- WebSocket for real-time updates

### Mobile Companion App

- Quick expense entry
- Receipt photo capture
- Push notifications for reminders
- Sync via API server

### Machine Learning Features

- Auto-categorization based on payee/description
- Anomaly detection (unusual spending)
- Budget recommendations based on historical patterns
- Income/expense forecasting with ML models

### Multi-User Support

- Shared accounts/budgets (household)
- User roles and permissions
- Separate views per user

### Investment Tracking

- Portfolio accounts
- Stock/crypto prices (via APIs)
- Capital gains calculations
- Asset allocation views

### Bill Pay Integration

- Direct payment from CLI (via bank APIs)
- Payment confirmations
- Transaction auto-reconciliation

---

## Open Questions & Design Decisions Needed

1. **Default behavior for recurring transactions:**
   - Auto-generate or require manual confirmation?
   - *Recommendation:* Manual confirmation (safer), with `--auto` flag for power users

2. **Budget rollover:**
   - Allow unused budget to roll to next period?
   - *Recommendation:* Optional per-budget setting

3. **Multi-currency:**
   - Support in Phase 1 or defer to Phase 4?
   - *Recommendation:* Defer to Phase 4, focus on USD first

4. **Transaction splits:**
   - Single transaction → multiple categories (e.g., Costco: groceries + household)
   - *Recommendation:* Phase 4, complex schema change

5. **Reconciliation workflow:**
   - Mark transactions as reconciled, track unreconciled balance
   - *Recommendation:* Phase 4, add `is_reconciled` flag and reconciliation report

6. **Backup/restore:**
   - Built-in backup command or rely on `pg_dump`?
   - *Recommendation:* Document `pg_dump` usage, add `fintrack backup` wrapper in Phase 5

7. **Soft deletes:**
   - Mark records as deleted vs hard delete?
   - *Recommendation:* Hard delete for transactions, soft delete for accounts (set `is_active=false`)

8. **Timezone handling:**
   - Store timestamps in UTC or user's local timezone?
   - *Recommendation:* Store in UTC, display in user's timezone (from config)

---

## Success Metrics

### User Experience
- [ ] Can add transaction in <10 seconds
- [ ] Can import bank statement in <30 seconds
- [ ] Projection generates in <1 second
- [ ] Zero data loss in transaction processing

### Technical
- [ ] 95% test coverage for core business logic
- [ ] All queries <100ms with 10,000+ transactions
- [ ] Binary size <15MB
- [ ] Memory usage <50MB for typical workloads

### Feature Completeness
- [ ] All 10 core commands implemented
- [ ] JSON output for all commands
- [ ] CSV import for 5+ major banks
- [ ] 90-day projection with 80%+ confidence

---

## Documentation Plan

1. **README.md**: Quick start, installation, basic usage
2. **USAGE.md**: Comprehensive command reference
3. **IMPORT_GUIDE.md**: CSV import for various banks
4. **DATABASE.md**: Schema documentation, migration guide
5. **DEVELOPMENT.md**: Contributing guide, architecture deep-dive
6. **FAQ.md**: Common questions and troubleshooting

---

## Conclusion

This design provides a solid foundation for a Unix-philosophy personal finance CLI tool. Key strengths:

- **Simple**: Core commands map to natural workflows
- **Portable**: PostgreSQL + single binary, no vendor lock-in
- **Composable**: Text output enables piping and scripting
- **Extensible**: Clear phases for adding advanced features
- **Privacy-focused**: Local-first, no cloud dependencies

Next steps:
1. Validate design with potential users
2. Create database schema SQL files
3. Scaffold Go project structure
4. Implement Phase 1 MVP
5. Gather feedback and iterate

---

**Prepared by:** Claude (AI Assistant)
**Date:** 2025-11-16
**Version:** 1.0 - Initial Planning Document
