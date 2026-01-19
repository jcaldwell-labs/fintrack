# Account Management Usage Tests

This file contains executable usage tests for account management features. These tests serve as both documentation and automated validation.

**Format Notes:**

- `<any>` - Matches any value (for IDs, timestamps, etc.)
- `<number>` - Matches any numeric value
- `<date>` - Matches date format YYYY-MM-DD
- `<money>` - Matches currency format $X,XXX.XX

---

## Test: Create a checking account

**Purpose:** Verify users can create a basic checking account with initial balance

### Setup

```bash
# Clean slate - ensure no existing test accounts
fintrack account delete "Test Checking" 2>/dev/null || true
```

### Execute

```bash
fintrack account add "Test Checking" --type checking --balance 1000.00
```

### Expected Output

```
Account created successfully
ID: <number>
Name: Test Checking
Type: checking
Balance: <money>
```

### Actual Output (auto-updated)

```

```

❌ FAIL (last run: 2026-01-19)

---

## Test: List all accounts

**Purpose:** Verify account listing displays all active accounts in table format

### Setup

```bash
# Ensure at least one account exists
fintrack account add "Test List Account" --type savings --balance 500.00 2>/dev/null || true
```

### Execute

```bash
fintrack account list
```

### Expected Output

```
ID    Name                  Type       Balance      Currency  Status
<number>    Test List Account      savings    <money>      USD       active
```

### Actual Output (auto-updated)

```

```

❌ FAIL (last run: 2026-01-19)

---

## Test: Show specific account details

**Purpose:** Verify detailed account information retrieval by ID

### Setup

```bash
# Create a test account and capture its ID
fintrack account add "Test Show Account" --type checking --balance 2500.50 2>/dev/null || true
```

### Execute

```bash
fintrack account show 1
```

### Expected Output

```
Account Details
─────────────────────────────────
ID:              <number>
Name:            Test Show Account
Type:            checking
Balance:         <money>
Currency:        USD
Institution:
Account Number:
Status:          active
Created:         <date>
```

### Actual Output (auto-updated)

```

```

❌ FAIL (last run: 2026-01-19)

---

## Test: Update account name

**Purpose:** Verify account name can be updated

### Setup

```bash
# Create account to update
fintrack account add "Original Name" --type checking --balance 1000.00
```

### Execute

```bash
fintrack account update 1 --name "Updated Name"
```

### Expected Output

```
Account updated successfully
ID: <number>
Name: Updated Name
Type: checking
Balance: <money>
```

### Actual Output (auto-updated)

```

```

❌ FAIL (last run: 2026-01-19)

---

## Test: Close account (soft delete)

**Purpose:** Verify account can be closed (soft delete) while preserving history

### Setup

```bash
# Create account to close
fintrack account add "Account to Close" --type savings --balance 100.00
```

### Execute

```bash
fintrack account close 1
```

### Expected Output

```
Account closed successfully
ID: <number>
Name: Account to Close
Status: inactive
```

### Actual Output (auto-updated)

```

```

❌ FAIL (last run: 2026-01-19)

---

## Test: JSON output format

**Purpose:** Verify JSON output format for scripting and automation

### Setup

```bash
# Ensure test account exists
fintrack account add "JSON Test Account" --type checking --balance 999.99 2>/dev/null || true
```

### Execute

```bash
fintrack account list --json
```

### Expected Output

```
[
  {
    "id": <number>,
    "name": "JSON Test Account",
    "type": "checking",
    "balance": 999.99,
    "currency": "USD",
    "is_active": true
  }
]
```

### Actual Output (auto-updated)

```

```

❌ FAIL (last run: 2026-01-19)

---

## Test: Account type validation

**Purpose:** Verify invalid account types are rejected with helpful error message

### Setup

```bash
# No setup needed
```

### Execute

```bash
fintrack account add "Invalid Type Account" --type invalid_type --balance 100.00 2>&1
```

### Expected Output

```
Error: invalid account type "invalid_type". Valid types: checking, savings, credit, cash, investment, loan
```

### Actual Output (auto-updated)

```

```

❌ FAIL (last run: 2026-01-19)

---

## Test: Duplicate account name prevention

**Purpose:** Verify duplicate account names are prevented

### Setup

```bash
# Create first account
fintrack account add "Duplicate Name" --type checking --balance 100.00
```

### Execute

```bash
fintrack account add "Duplicate Name" --type savings --balance 200.00 2>&1
```

### Expected Output

```
Error: account with name "Duplicate Name" already exists
```

### Actual Output (auto-updated)

```

```

❌ FAIL (last run: 2026-01-19)

---

## Notes

### Running These Tests

```bash
# Run all usage tests
make test-usage

# Run with verbose output
go test -v ./tests/usage/

# Update documentation with actual results
make test-usage-update
```

### Test Database

These tests use a separate `fintrack_test` database to avoid affecting production data:

```bash
# Create test database (one-time setup)
createdb fintrack_test

# Clean test data between runs
psql -d fintrack_test -c "TRUNCATE accounts CASCADE;"
```

### Wildcard Patterns

- `<any>` - Matches any value
- `<number>` - Matches integers (e.g., 42, 1, 999)
- `<date>` - Matches YYYY-MM-DD format
- `<uuid>` - Matches UUID format
- `<money>` - Matches currency like $1,234.56

### Test Philosophy

These usage tests are designed to:

1. **Document real-world usage** - Show exactly how users interact with the CLI
2. **Validate output formatting** - Ensure consistent, user-friendly output
3. **Catch regressions** - Detect when changes break existing behavior
4. **Serve as examples** - Provide copy-paste examples for documentation

### Maintenance

- Tests automatically update actual output after each run
- PASS/FAIL status is tracked with timestamps
- Failed tests show diffs to help debug issues
- Markdown format is human-readable and VCS-friendly
