# FinTrack Quick Reference

## Command Cheat Sheet

### Account Management
```bash
# List accounts
fintrack account list
fintrack a ls

# Add account
fintrack account add "Account Name" --type checking --balance 5000
fintrack a add "Credit Card" -t credit -b -1200

# Show account
fintrack account show "Account Name"
fintrack a show 1

# Update account
fintrack account update 1 --name "New Name" --balance 5432.10

# Close account
fintrack account close 1
```

### Transactions
```bash
# Add expense
fintrack tx add -100 "Payee" --account 1 --category groceries
fintrack t add -50.25 "Coffee Shop" -a checking -c dining

# Add income
fintrack tx add 3500 "Employer" --category salary

# Quick entry (uses default account)
fintrack tx -75 "Quick expense"

# Transfer
fintrack tx transfer 500 --from checking --to savings

# List transactions
fintrack tx list --limit 20
fintrack tx list --from 2025-11-01 --to 2025-11-30
fintrack tx list --account checking --month 2025-11

# Search
fintrack tx search "amazon"
fintrack tx search --payee "Whole Foods" --min 50

# Edit/delete
fintrack tx edit 123 --amount -105.50
fintrack tx delete 123

# Tag
fintrack tx tag 123 "business" "reimbursable"
```

### Budgets
```bash
# Create budget
fintrack budget create "Groceries" --category groceries --limit 600 --period monthly
fintrack b create "Dining" -c dining -l 300 -p monthly

# List budgets
fintrack budget list
fintrack b ls --active

# Show budget status
fintrack budget show groceries
fintrack budget status --month 2025-11

# Update
fintrack budget update groceries --limit 700

# Report
fintrack budget report --month 2025-11
```

### Recurring Transactions
```bash
# Add recurring item
fintrack schedule add "Rent" -1500 --category rent --frequency monthly --day 1
fintrack s add "Netflix" -15.99 -c subscriptions -f monthly

# List
fintrack schedule list
fintrack s ls --upcoming 30

# Generate transactions
fintrack schedule generate --dry-run
fintrack schedule generate --confirm

# Pause/resume
fintrack schedule pause 3
fintrack schedule resume 3
```

### Reminders
```bash
# List reminders
fintrack remind list
fintrack r ls --today
fintrack r ls --week

# Add custom reminder
fintrack remind add "Pay credit card" --date 2025-11-20 --time 09:00

# Dismiss
fintrack remind dismiss 5

# Show overdue
fintrack remind overdue
```

### Cash Flow Projection
```bash
# Generate projection
fintrack project --days 90
fintrack p --months 6 --type conservative

# By account
fintrack project --account checking --days 30

# Summary
fintrack project summary --weeks 4
```

### Reports
```bash
# Income statement
fintrack report income --month 2025-11
fintrack rp income --year 2025

# Spending analysis
fintrack report spending --month 2025-11
fintrack report spending --category

# Net worth
fintrack report networth --months 12

# Custom range
fintrack report summary --from 2025-01-01 --to 2025-11-15

# Export
fintrack report income --year 2025 --format csv > income.csv
```

### Calendar
```bash
# Show calendar
fintrack cal
fintrack c --month 2025-12

# Specific types
fintrack cal --type bills
fintrack cal --type income

# Upcoming
fintrack cal next --days 7
```

### Import
```bash
# Import CSV
fintrack import csv file.csv --account checking --format generic
fintrack import csv chase.csv --account checking --format chase --dry-run
fintrack import csv transactions.csv -a 1 -f generic --confirm
```

### Configuration
```bash
# Show config
fintrack config show

# Set values
fintrack config set db.url "postgresql://localhost:5432/fintrack"
fintrack config set defaults.account "Checking"
fintrack config set defaults.currency USD

# Database setup
fintrack config init-db
fintrack config migrate
```

### Statistics
```bash
# Dashboard
fintrack stats

# Trends
fintrack stats trends --months 6

# Top spending
fintrack stats top-spending --month 2025-11

# Savings rate
fintrack stats savings-rate --year 2025
```

---

## Common Workflows

### Daily Use
```bash
# Check today's reminders
fintrack remind list --today

# Quick expense entry
fintrack tx -45 "Grocery store" -c groceries

# See current account balances
fintrack account list
```

### Monthly Review
```bash
# Budget status
fintrack budget report --month 2025-11

# Income statement
fintrack report income --month 2025-11

# Spending breakdown
fintrack report spending --month 2025-11 --category

# Update projections
fintrack project --days 90
```

### Weekly Tasks
```bash
# Check upcoming bills
fintrack cal next --days 7

# Review unreconciled transactions
fintrack tx list --unreconciled

# Generate recurring transactions
fintrack schedule generate --confirm
```

### Import Bank Statement
```bash
# 1. Download CSV from bank
# 2. Dry-run import to check
fintrack import csv statement.csv --account checking --format chase --dry-run

# 3. Review output
# 4. Confirm import
fintrack import csv statement.csv --account checking --format chase --confirm

# 5. Reconcile account
fintrack account reconcile checking --balance 5432.10
```

---

## JSON Output (for scripting)

All commands support `--json` flag for machine-readable output:

```bash
# Get account data as JSON
fintrack account list --json

# Parse with jq
fintrack tx list --month 2025-11 --json | jq '.data[] | select(.amount < -100)'

# Export to file
fintrack report income --year 2025 --json > income_2025.json
```

---

## Environment Variables

```bash
# Database connection
export FINTRACK_DB_URL="postgresql://user:pass@localhost:5432/fintrack"

# Custom config location
export FINTRACK_CONFIG="$HOME/.config/fintrack/config.yaml"

# Database password (for security)
export FINTRACK_DB_PASSWORD="secret"
```

---

## Configuration File Example

`~/.config/fintrack/config.yaml`:

```yaml
database:
  url: "postgresql://localhost:5432/fintrack"
  # Or use password from environment:
  password: "${FINTRACK_DB_PASSWORD}"

defaults:
  account: "Chase Checking"
  currency: "USD"

alerts:
  enabled: true
  threshold: 0.80

recurring:
  auto_generate: false
  generate_days_ahead: 3

projection:
  default_days: 90
  scenario: "moderate"

output:
  default_format: "table"
  color: true
  unicode: true
```

---

## Tips & Tricks

### 1. Quick Transaction Entry
Set a default account to skip the `--account` flag:
```bash
fintrack config set defaults.account "Checking"
fintrack tx -50 "Coffee"  # Uses default account
```

### 2. Shell Aliases
```bash
alias ft="fintrack"
alias fta="fintrack account"
alias ftt="fintrack tx"
alias ftb="fintrack budget"

# Now use:
ft tx -100 "Groceries"
ftb status
```

### 3. Budget Alerts
Set aggressive thresholds for important categories:
```bash
fintrack budget create "Dining Out" -c dining -l 200 --alert 0.70
```

### 4. Search and Tag
Tag transactions for later analysis:
```bash
fintrack tx search "amazon" --month 2025-11 --json | \
  jq -r '.data[].id' | \
  xargs -I {} fintrack tx tag {} "online-shopping"
```

### 5. Monthly Automation
Create a script for month-end tasks:
```bash
#!/bin/bash
# monthly_review.sh

MONTH=$(date +%Y-%m)

echo "=== Budget Report ==="
fintrack budget report --month $MONTH

echo -e "\n=== Income Statement ==="
fintrack report income --month $MONTH

echo -e "\n=== Spending Analysis ==="
fintrack report spending --month $MONTH

echo -e "\n=== Next Month Projection ==="
fintrack project --days 30
```

### 6. Backup Automation
```bash
#!/bin/bash
# backup_fintrack.sh

BACKUP_DIR="$HOME/backups/fintrack"
DATE=$(date +%Y%m%d)

mkdir -p $BACKUP_DIR
pg_dump fintrack | gzip > "$BACKUP_DIR/fintrack_$DATE.sql.gz"

# Keep last 30 days
find $BACKUP_DIR -name "fintrack_*.sql.gz" -mtime +30 -delete
```

---

## CSV Import Formats

### Generic Format
```csv
date,amount,payee,category,description
2025-11-15,-100.00,Whole Foods,Groceries,Weekly shopping
2025-11-14,3500.00,Acme Corp,Salary,Paycheck
```

### Chase Format Mapping
```yaml
# chase_mapping.yaml
columns:
  date: "Posting Date"
  amount: "Amount"
  payee: "Description"
  category: "Type"
date_format: "MM/DD/YYYY"
amount_sign: inverse  # Chase uses positive for debits
skip_rows: 1
```

Import:
```bash
fintrack import csv chase.csv --account checking --format chase --map chase_mapping.yaml
```

---

## Troubleshooting

### Connection Issues
```bash
# Test database connection
psql -U fintrack_user -d fintrack -c "SELECT version();"

# Verify config
fintrack config show
```

### Transaction Not Updating Balance
```bash
# Manually recalculate balances (should not be needed with triggers)
fintrack account reconcile 1 --balance $(psql -U fintrack_user -d fintrack -tAc \
  "SELECT SUM(amount) FROM transactions WHERE account_id = 1")
```

### Performance Issues
```bash
# Refresh materialized views
psql -U fintrack_user -d fintrack -c \
  "REFRESH MATERIALIZED VIEW CONCURRENTLY monthly_spending_by_category;"
```

---

## Keyboard Shortcuts (TUI Mode - Future)

When using `fintrack tui` (interactive mode):

- `j/k`: Navigate down/up
- `h/l`: Navigate left/right (tabs)
- `Enter`: Select/edit
- `d`: Delete
- `a`: Add new
- `?`: Help
- `q`: Quit

---

## Advanced Queries (PostgreSQL)

Direct database access for power users:

```sql
-- Find all unreconciled transactions
SELECT * FROM transactions WHERE is_reconciled = false ORDER BY date;

-- Monthly spending trend
SELECT
  DATE_TRUNC('month', date) as month,
  SUM(amount) as total_spent
FROM transactions
WHERE type = 'expense'
GROUP BY month
ORDER BY month DESC
LIMIT 12;

-- Top 10 payees by total spending
SELECT
  payee,
  COUNT(*) as count,
  SUM(amount) as total
FROM transactions
WHERE type = 'expense'
GROUP BY payee
ORDER BY total ASC
LIMIT 10;

-- Category hierarchy with totals
WITH RECURSIVE category_tree AS (
  SELECT id, name, parent_id, 0 as level
  FROM categories
  WHERE parent_id IS NULL

  UNION ALL

  SELECT c.id, c.name, c.parent_id, ct.level + 1
  FROM categories c
  JOIN category_tree ct ON c.parent_id = ct.id
)
SELECT
  ct.level,
  ct.name,
  COALESCE(SUM(t.amount), 0) as total
FROM category_tree ct
LEFT JOIN transactions t ON t.category_id = ct.id
GROUP BY ct.id, ct.level, ct.name
ORDER BY ct.level, ct.name;
```

---

**Version:** 1.0
**Last Updated:** 2025-11-16
**For:** FinTrack Personal Finance CLI
