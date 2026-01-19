# FinTrack Demo Script

> 2-minute demo for showcasing FinTrack's core value proposition

## The Pitch (30 seconds)

> "FinTrack is personal finance tracking for developers who live in the terminal. Your data stays local in PostgreSQL - no cloud accounts, no subscriptions, no data harvesting. And here's the killer feature: JSON output. Pipe to jq, script your budgets in bash, build dashboards with your existing tools."

## Prerequisites

Before running the demo, ensure:

```bash
# PostgreSQL is running
pg_isready

# Database exists
createdb fintrack_demo 2>/dev/null || true

# Environment configured
export FINTRACK_DB_URL="postgresql://localhost:5432/fintrack_demo?sslmode=disable"

# Binary is built
make build
```

## Demo Flow (90 seconds)

### 1. Create an Account (15 seconds)

```bash
# Create a checking account with initial balance
./bin/fintrack account add "Checking" --type checking --balance 5000

# Show it worked
./bin/fintrack account list
```

**Expected Output:**

```
ID  NAME      TYPE      BALANCE     CURRENCY  ACTIVE
1   Checking  checking  $5,000.00   USD       Yes
```

**Talking Point:** "Notice the balance is stored locally in PostgreSQL - your financial data never leaves your machine."

### 2. Add Some Transactions (30 seconds)

```bash
# Add a category first
./bin/fintrack category add "Groceries" expense

# Add expense transactions
./bin/fintrack tx add --account 1 --amount -85.50 --payee "Whole Foods" --category 1
./bin/fintrack tx add --account 1 --amount -42.30 --payee "Trader Joe's" --category 1
./bin/fintrack tx add --account 1 --amount -15.99 --payee "Coffee Shop" --category 1

# Add income
./bin/fintrack tx add --account 1 --amount 2500.00 --payee "Employer" --type income

# View transactions
./bin/fintrack tx list
```

**Expected Output:**

```
ID  DATE        AMOUNT     TYPE     PAYEE          CATEGORY   ACCOUNT
1   2026-01-19  -85.50     expense  Whole Foods    Groceries  Checking
2   2026-01-19  -42.30     expense  Trader Joe's   Groceries  Checking
3   2026-01-19  -15.99     expense  Coffee Shop    Groceries  Checking
4   2026-01-19  +2500.00   income   Employer                  Checking

Summary: 4 transactions | Income: $2,500.00 | Expenses: $143.79 | Net: $2,356.21
```

**Talking Point:** "FinTrack auto-detects transaction types from the amount sign. Negative is expense, positive is income."

### 3. The "Aha" Moment - JSON Output (45 seconds)

```bash
# Get JSON output
./bin/fintrack tx list --json
```

**Now the magic - pipe to jq:**

```bash
# Find all expenses
./bin/fintrack tx list --json | jq '.[] | select(.amount_cents < 0) | .payee'

# Calculate total spent at grocery stores
./bin/fintrack tx list --json | jq '[.[] | select(.payee | test("Foods|Joe"))] | map(.amount_cents) | add / 100'

# Get account balances as JSON for dashboards
./bin/fintrack account list --json | jq '.[] | {name, balance: .current_balance_cents / 100}'
```

**Talking Point:** "This is what sets FinTrack apart. Every command outputs clean JSON. Build shell scripts, create cron jobs, pipe to monitoring systems, or feed your own dashboards. Your data, your tools, your automation."

## Closing (15 seconds)

> "FinTrack is open source, MIT licensed, and designed for developers who want control over their financial data. No vendor lock-in, no cloud dependency, just clean Unix-philosophy tooling."

**Call to Action:**

```bash
# Install
go install github.com/fintrack/fintrack@latest

# Or clone and build
git clone https://github.com/fintrack/fintrack && cd fintrack && make install
```

## Reset Demo Environment

```bash
# Drop and recreate demo database
dropdb fintrack_demo
createdb fintrack_demo
```

## Common Demo Variations

### Quick Version (60 seconds)

Skip the category creation and use just:

```bash
./bin/fintrack account add "Checking" -t checking -b 5000
./bin/fintrack tx add --account 1 --amount -50 --payee "Store"
./bin/fintrack tx list --json | jq '.'
```

### Technical Audience

Add these commands to show power user features:

```bash
# Filter by date range
./bin/fintrack tx list --from 2026-01-01 --to 2026-01-31

# Filter by payee (partial match)
./bin/fintrack tx list --payee "Foods"

# Combine filters
./bin/fintrack tx list --type expense --limit 10 --json
```

### Import Focus

If demonstrating CSV import capabilities:

```bash
# Dry run first (experimental feature)
./bin/fintrack import csv bank_export.csv --account 1 --dry-run

# Then import
./bin/fintrack import csv bank_export.csv --account 1 --skip-duplicates

# View import history
./bin/fintrack import history
```

## Troubleshooting

| Issue                     | Solution                                   |
| ------------------------- | ------------------------------------------ |
| "database not configured" | Set `FINTRACK_DB_URL` environment variable |
| "connection refused"      | Start PostgreSQL: `pg_ctl start`           |
| "relation does not exist" | Tables auto-migrate on first run           |
| JSON output malformed     | Ensure `--json` flag is at the end         |
