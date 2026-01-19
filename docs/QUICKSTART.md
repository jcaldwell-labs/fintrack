# FinTrack Quick Start Guide

Get up and running with FinTrack in 5 minutes.

## Step 1: Install (1 minute)

### Option A: From Source (Recommended)

```bash
# Clone the repository
git clone https://github.com/fintrack/fintrack.git
cd fintrack

# Build
make build

# (Optional) Install to PATH
sudo make install
```

### Option B: Go Install

```bash
go install github.com/fintrack/fintrack@latest
```

## Step 2: Database Setup (2 minutes)

FinTrack uses PostgreSQL for reliable, local data storage.

### Install PostgreSQL (if needed)

```bash
# Ubuntu/Debian
sudo apt install postgresql

# macOS (Homebrew)
brew install postgresql@15
brew services start postgresql@15

# Verify it's running
pg_isready
```

### Create Database

```bash
createdb fintrack
```

### Configure Connection

**Option A: Environment Variable (Recommended)**

```bash
# Add to your shell profile
export FINTRACK_DB_URL="postgresql://localhost:5432/fintrack?sslmode=disable"
```

**Option B: Config File**

Create ~/.config/fintrack/config.yaml:

```yaml
database:
  url: "postgresql://localhost:5432/fintrack?sslmode=disable"
```

### Verify Connection

```bash
fintrack account list
```

If you see an empty table (not an error), you're ready to go.

## Step 3: Create Your First Account (30 seconds)

```bash
# Create a checking account
fintrack account add "Main Checking" --type checking --balance 2500

# Create a savings account
fintrack account add "Emergency Fund" --type savings --balance 10000

# Create a credit card
fintrack account add "Chase Visa" --type credit --balance -450

# View your accounts
fintrack account list
```

**Output:**

```
ID  NAME            TYPE      BALANCE      CURRENCY  ACTIVE
1   Main Checking   checking  $2,500.00    USD       Yes
2   Emergency Fund  savings   $10,000.00   USD       Yes
3   Chase Visa      credit    -$450.00     USD       Yes
```

## Step 4: Set Up Categories (30 seconds)

```bash
# Add expense categories
fintrack category add "Groceries" expense
fintrack category add "Dining" expense
fintrack category add "Transportation" expense
fintrack category add "Utilities" expense

# Add income category
fintrack category add "Salary" income

# View categories
fintrack category list
```

## Step 5: Track Transactions (1 minute)

```bash
# Add an expense (negative amount)
fintrack tx add --account 1 --amount -75.50 --payee "Whole Foods" --category 1

# Add income (positive amount)
fintrack tx add --account 1 --amount 3500 --payee "Employer" --category 5 --type income

# Add a bill payment
fintrack tx add --account 1 --amount -125.00 --payee "Electric Company" --category 4

# View recent transactions
fintrack tx list

# View transactions for a specific account
fintrack tx list --account 1
```

## Command Reference

### Accounts

| Command | Description |
|---------|-------------|
| fintrack account list | List all accounts |
| fintrack account add NAME --type TYPE --balance BAL | Create account |
| fintrack account show ID | Show account details |
| fintrack account update ID --name NAME | Update account |
| fintrack account close ID | Close account |

**Account Types:** checking, savings, credit, cash, investment, loan

### Categories

| Command | Description |
|---------|-------------|
| fintrack category list | List all categories |
| fintrack category add NAME TYPE | Create category (TYPE: income/expense) |
| fintrack category show ID | Show category details |
| fintrack category update ID --name NAME | Update category |
| fintrack category delete ID | Delete category |

### Transactions

| Command | Description |
|---------|-------------|
| fintrack tx list | List recent transactions |
| fintrack tx add --account ID --amount AMT | Add transaction |
| fintrack tx show ID | Show transaction details |
| fintrack tx update ID --payee NAME | Update transaction |
| fintrack tx delete ID | Delete transaction |

**Transaction Filters:**

```bash
fintrack tx list --account 1          # By account
fintrack tx list --category 2         # By category
fintrack tx list --type expense       # By type
fintrack tx list --from 2026-01-01    # From date
fintrack tx list --to 2026-01-31      # To date
fintrack tx list --payee "Whole"      # By payee (partial match)
fintrack tx list --limit 100          # Limit results
```

### JSON Output

Add --json to any list command for machine-readable output:

```bash
fintrack account list --json
fintrack tx list --json
fintrack category list --json
```

## Aliases

FinTrack supports short aliases for faster typing:

| Full Command | Alias |
|--------------|-------|
| fintrack account | fintrack a |
| fintrack category | fintrack cat |
| fintrack transaction | fintrack tx or fintrack t |
| fintrack account list | fintrack a ls |
| fintrack transaction list | fintrack tx ls |

## Troubleshooting

### "database not configured"

Set the database URL:

```bash
export FINTRACK_DB_URL="postgresql://localhost:5432/fintrack?sslmode=disable"
```

### "connection refused"

PostgreSQL isn't running:

```bash
# Linux
sudo systemctl start postgresql

# macOS
brew services start postgresql@15
```

### "relation does not exist"

Tables are created automatically on first run. If issues persist:

```bash
# Reset database
dropdb fintrack
createdb fintrack
fintrack account list  # Triggers auto-migration
```

## Next Steps

- Read the full [README](../README.md) for detailed documentation
- Check [DEMO.md](DEMO.md) for presentation tips
- Explore fintrack --help for all commands
