# Usage Documentation Tests

This directory contains **executable usage documentation** - markdown files that serve as both user-facing documentation and automated regression tests.

## Philosophy

Traditional documentation tends to become outdated because it's not validated. Usage tests solve this by making documentation executable:

- **Write once, verify always** - Examples are tested on every CI run
- **Self-documenting** - Real commands with real output
- **Regression prevention** - Breaking changes are caught immediately
- **Onboarding friendly** - New contributors see working examples

## How It Works

### 1. Write Documentation

Create a markdown file with test cases:

```markdown
## Test: Create a checking account
**Purpose:** Verify users can create a basic checking account

### Setup
```bash
# Clean slate - ensure no existing test accounts
fintrack account delete "Test Checking" 2>/dev/null || true
```

### Execute
```bash
fintrack account create "Test Checking" --type checking --balance 1000.00
```

### Expected Output
```
Account created successfully
ID: <number>
Name: Test Checking
```

### Actual Output (auto-updated)
```
(Results will be inserted here)
```

❌ FAIL (last run: 2025-11-23)
```

### 2. Run Tests

```bash
# Build and run usage tests
make test-usage

# Or run directly
go test -v ./tests/usage/
```

### 3. Review Results

The test harness:
1. Parses markdown files
2. Builds the fintrack binary
3. Executes setup commands
4. Runs the main command
5. Captures output
6. Compares with expected output (with wildcards)
7. Updates markdown with actual output
8. Updates pass/fail status

## Markdown Format

### Required Sections

Each test case must have:

- `## Test: [Name]` - Test case title
- `**Purpose:**` - Brief description
- `### Setup` - Bash commands for test setup
- `### Execute` - The command to test
- `### Expected Output` - Expected output with wildcards
- `### Actual Output (auto-updated)` - Auto-populated results
- Status line: `✅ PASS (last run: YYYY-MM-DD)` or `❌ FAIL (last run: YYYY-MM-DD)`

### Wildcard Patterns

Use wildcards for dynamic values:

| Pattern | Matches | Example |
|---------|---------|---------|
| `<any>` | Any value | `42`, `foo`, `2025-01-01` |
| `<number>` | Integers | `42`, `1`, `999` |
| `<date>` | YYYY-MM-DD | `2025-11-23` |
| `<uuid>` | UUID format | `550e8400-e29b-41d4-a716-446655440000` |
| `<money>` | Currency | `$1,234.56`, `$0.00` |

### Example

```markdown
Expected:
```
Account created successfully
ID: <number>
Name: My Account
Balance: <money>
```

Actual:
```
Account created successfully
ID: 42
Name: My Account
Balance: $1,000.00
```

Result: ✅ PASS - `<number>` matches `42`, `<money>` matches `$1,000.00`
```

## File Naming Convention

Use numbered prefixes for logical ordering:

```
tests/usage/
├── 01-account-management.md
├── 02-transaction-tracking.md
├── 03-budget-management.md
├── 04-reporting.md
└── 99-edge-cases.md
```

## Running Tests

### Locally

```bash
# Run all usage tests
make test-usage

# Run and update markdown files
make test-usage-update

# Run with verbose output
go test -v ./tests/usage/

# Run specific test
go test -v ./tests/usage/ -run "TestUsageDocumentation/01-account-management.md"
```

### In CI/CD

Usage tests run automatically in GitHub Actions:
- On every push to `main` or `develop`
- On every pull request
- Results are uploaded as artifacts

### Test Database

Tests use a dedicated test database to avoid affecting production data:

```bash
# One-time setup
createdb fintrack_test

# Tests require database URL to be set (credentials via environment):
# Format: postgresql://USER:PASS@HOST:PORT/DATABASE
export FINTRACK_DB_URL="postgresql://<redacted>@localhost:5432/fintrack_test"
```

## Writing Good Usage Tests

### DO:

✅ Test common user workflows
✅ Use realistic examples
✅ Include error cases
✅ Test both table and JSON output
✅ Use wildcards for dynamic values
✅ Keep tests focused and simple
✅ Add context in the Purpose section

### DON'T:

❌ Test implementation details
❌ Make tests too complex
❌ Hardcode IDs or timestamps
❌ Rely on state from other tests
❌ Skip setup/cleanup
❌ Test multiple things in one test

## Example Test Structure

```markdown
# Feature Area Usage Tests

Brief overview of what this file tests.

---

## Test: Happy path scenario
**Purpose:** Brief description of what this validates

### Setup
```bash
# Setup commands here
```

### Execute
```bash
fintrack command --flags
```

### Expected Output
```
Output with <wildcards>
```

### Actual Output (auto-updated)
```
(Auto-generated)
```

✅/❌ STATUS (last run: YYYY-MM-DD)

---

## Test: Error case
**Purpose:** Verify error handling

[... similar structure ...]

---

## Notes

Any relevant notes about these tests, gotchas, or maintenance tips.
```

## Maintenance

### Updating Tests

When behavior changes intentionally:
1. Update the "Expected Output" section
2. Run `make test-usage`
3. Verify the new actual output is correct
4. Commit both changes

### Debugging Failures

When tests fail:
1. Check the "Actual Output" section for what actually happened
2. Compare with "Expected Output"
3. Determine if it's a bug or expected behavior change
4. Fix the code or update the expected output

### Adding New Tests

When adding features:
1. Create or update relevant markdown file
2. Add test cases for happy path and error cases
3. Run `make test-usage` to generate initial results
4. Review and commit

## Architecture

### Test Harness Components

- **runner.go** - Core logic for parsing and executing tests
  - `ParseUsageTestFile()` - Parses markdown into test cases
  - `RunTestCase()` - Executes a single test case
  - `matchOutput()` - Compares expected vs actual with wildcards
  - `UpdateMarkdownFile()` - Writes results back to markdown

- **usage_test.go** - Go test integration
  - `TestUsageDocumentation()` - Main test entry point
  - Builds binary, sets up database, runs all tests

### Parser Logic

The parser uses state machine logic to extract test cases:
1. Scans for `## Test:` markers
2. Extracts purpose, setup, execute, expected/actual sections
3. Identifies code blocks by language (`bash`, `expected`, `actual`)
4. Parses status lines for last run timestamp

### Matcher Logic

The output matcher:
1. Splits output into lines
2. Compares line-by-line
3. Converts wildcards to regex patterns
4. Returns pass/fail with detailed diff on failure

## Contributing

When adding or modifying usage tests:

1. **Follow the format** - Use the required sections
2. **Test your tests** - Run `make test-usage` before committing
3. **Use wildcards** - Don't hardcode dynamic values
4. **Add context** - Explain what the test validates
5. **Keep it simple** - One thing per test
6. **Document gotchas** - Add notes if something is tricky

## Questions?

See:
- Main testing docs: `/TESTING.md`
- Project guidance: `/.claude/CLAUDE.md`
- Integration tests: `/tests/integration/`
- Example tests: `/tests/usage/01-account-management.md`
