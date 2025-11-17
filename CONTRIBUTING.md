# Contributing to FinTrack

Thank you for your interest in contributing to FinTrack! This document provides guidelines and workflows for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Branching Strategy](#branching-strategy)
- [Version & Release Process](#version--release-process)
- [Code Review](#code-review)
- [Testing Requirements](#testing-requirements)

## Code of Conduct

### Our Pledge

We are committed to providing a welcoming and inspiring community for all. Please be respectful and constructive in your interactions.

### Expected Behavior

- Use welcoming and inclusive language
- Be respectful of differing viewpoints and experiences
- Gracefully accept constructive criticism
- Focus on what is best for the community
- Show empathy towards other community members

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git
- Make (optional, for build automation)

### Setting Up Development Environment

```bash
# Clone the repository
git clone https://github.com/yourusername/fintrack.git
cd fintrack

# Install dependencies
make deps

# Set up database
createdb fintrack
psql -d fintrack -f migrations/001_initial_schema.sql

# Configure local settings
cp fintrack_config.example.yaml ~/.config/fintrack/config.yaml
# Edit config.yaml with your database credentials

# Build and test
make build
make test
```

### Running the Application

```bash
# Run directly
make run

# Or with live reload (requires air)
make dev

# Run tests
make test

# Check code coverage
make test-coverage
```

## Development Workflow

### 1. Pick an Issue

- Browse [open issues](https://github.com/yourusername/fintrack/issues)
- Comment on the issue to claim it
- Wait for maintainer approval before starting work

### 2. Create a Feature Branch

```bash
# Update main branch
git checkout main
git pull origin main

# Create feature branch
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
# or
git checkout -b docs/your-documentation-update
```

### 3. Make Your Changes

- Write clean, idiomatic Go code
- Follow existing code style and patterns
- Add tests for new functionality
- Update documentation as needed
- Keep commits focused and atomic

### 4. Test Your Changes

```bash
# Run all tests
make test

# Run tests with race detection
make test-race

# Check code coverage
make test-coverage-check

# Run quality checks
make quality
```

### 5. Commit Your Changes

See [Commit Guidelines](#commit-guidelines) below.

### 6. Push and Create PR

```bash
# Push your branch
git push -u origin feature/your-feature-name

# Create PR on GitHub
# Use the PR template
# Link to related issues
```

## Commit Guidelines

### Commit Message Format

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

#### Type

Must be one of:

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that don't affect code meaning (formatting, missing semicolons, etc.)
- **refactor**: Code change that neither fixes a bug nor adds a feature
- **perf**: Performance improvement
- **test**: Adding or updating tests
- **chore**: Changes to build process, tooling, dependencies, etc.
- **ci**: CI/CD configuration changes
- **revert**: Reverts a previous commit

#### Scope

Optional. The scope should be the name of the package/module affected:

- `commands` - Command layer changes
- `db` - Database layer changes
- `config` - Configuration changes
- `models` - Model changes
- `output` - Output formatting changes
- `cli` - CLI interface changes

#### Subject

- Use imperative, present tense: "add" not "added" nor "adds"
- Don't capitalize first letter
- No period (.) at the end
- Maximum 50 characters

#### Body

- Wrap at 72 characters
- Explain what and why, not how
- Use bullet points for multiple changes

#### Footer

- Reference issues: `Fixes #123`, `Closes #456`
- Breaking changes: `BREAKING CHANGE: description`
- Co-authors: `Co-Authored-By: Name <email>`

### Examples

#### Simple fix
```
fix(commands): handle empty account name in validation

Prevents panic when user provides empty string for account name.
Previously this would crash, now returns proper validation error.

Fixes #123
```

#### New feature
```
feat(commands): add support for multi-currency accounts

- Add currency field to account model
- Support EUR, GBP, JPY in addition to USD
- Auto-format amounts based on currency
- Add currency conversion utilities

This allows users to track accounts in different currencies
and see balances in their preferred display currency.

Closes #456
```

#### Breaking change
```
feat(db)!: change account balance from float to decimal

BREAKING CHANGE: Account.Balance field changed from float64 to
decimal.Decimal for accurate financial calculations. Migration
required for existing databases.

Migration: Run migrations/002_decimal_balance.sql

Closes #789
```

#### Documentation
```
docs(readme): add installation instructions for macOS

- Add Homebrew installation method
- Document macOS-specific setup steps
- Add troubleshooting section for common issues
```

#### Multiple changes
```
chore: update dependencies and tooling

- Update Go to 1.22
- Update testify to v1.9.0
- Update golangci-lint to v1.55
- Add pre-commit hook configuration

All tests passing with updated dependencies.
```

### Commit Attribution

All commits should include attribution to Claude Code if generated with AI assistance:

```
feat(commands): add transaction filtering

Implement advanced filtering for transaction list command.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

## Pull Request Process

### Before Creating a PR

1. ✅ All tests pass (`make test`)
2. ✅ Code coverage meets threshold (`make test-coverage-check`)
3. ✅ Code is formatted (`make fmt`)
4. ✅ No linting errors (`make lint`)
5. ✅ Documentation is updated
6. ✅ CHANGELOG.md is updated (if applicable)

### PR Title

Follow the same format as commit messages:

```
feat(commands): add CSV import functionality
fix(db): resolve connection pool leak
docs(contributing): update PR guidelines
```

### PR Description Template

```markdown
## Description
Brief description of what this PR does.

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Performance improvement
- [ ] Code refactoring
- [ ] Test improvement

## Related Issues
Fixes #(issue number)
Relates to #(issue number)

## Changes Made
- Change 1
- Change 2
- Change 3

## Testing Done
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing completed
- [ ] Tested on multiple platforms (if applicable)

## Screenshots (if applicable)
Before:
[screenshot or CLI output]

After:
[screenshot or CLI output]

## Checklist
- [ ] My code follows the project's code style
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Any dependent changes have been merged and published

## Additional Notes
Any additional information reviewers should know.
```

### PR Size Guidelines

Keep PRs focused and reasonably sized:

- **Small PR** (< 200 lines): Ideal, quick to review
- **Medium PR** (200-500 lines): Acceptable
- **Large PR** (500-1000 lines): Consider splitting
- **Huge PR** (> 1000 lines): Must justify size or split into multiple PRs

### Review Process

1. **Automated Checks**: GitHub Actions runs tests, linting, security scans
2. **Review Assignment**: Maintainer assigned automatically
3. **Code Review**: At least 1 approval required
4. **Address Feedback**: Make requested changes
5. **Final Approval**: Maintainer approves
6. **Merge**: Squash and merge (or rebase, based on preference)

### Review Timeline

- **Initial Response**: Within 2 business days
- **Full Review**: Within 5 business days
- **Urgent Fixes**: Within 1 business day (security, critical bugs)

## Branching Strategy

### Branch Naming Convention

```
<type>/<short-description>
```

**Types:**
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Test additions/improvements
- `chore/` - Maintenance tasks

**Examples:**
```
feature/csv-import
feature/multi-currency-support
fix/balance-calculation-overflow
fix/config-loading-error
docs/api-documentation
docs/getting-started-guide
refactor/repository-pattern
test/integration-account-workflow
chore/update-dependencies
```

### Main Branches

- **`main`** - Production-ready code
  - Always deployable
  - Protected branch (no direct commits)
  - Requires PR approval
  - CI/CD must pass

- **`develop`** - Integration branch (optional for larger teams)
  - Latest development changes
  - Staging/QA deployments
  - Feature branches merge here first

### Branch Protection Rules

**`main` branch:**
- ✅ Require pull request reviews (1 approval minimum)
- ✅ Require status checks to pass
  - All tests passing
  - Code coverage ≥ 60%
  - No linting errors
  - Security scan passing
- ✅ Require branches to be up to date
- ✅ No force pushes
- ✅ No deletions

### Workflow

```
main
  │
  ├─── feature/new-feature ──┐
  │                           │ (PR)
  ├─── fix/bug-fix ──────────┤
  │                           │
  ├─── docs/update-readme ───┤
  │                           │
  └─────────────────────────── (merge)
```

## Version & Release Process

### Semantic Versioning

We follow [Semantic Versioning 2.0.0](https://semver.org/):

```
MAJOR.MINOR.PATCH
```

- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality (backwards compatible)
- **PATCH**: Bug fixes (backwards compatible)

**Examples:**
- `0.1.0` - Initial MVP release
- `0.2.0` - Add transaction tracking feature
- `0.2.1` - Fix transaction date parsing bug
- `1.0.0` - First stable production release
- `2.0.0` - Breaking change: redesigned API

### Version Labeling

#### Pre-release Versions

```
MAJOR.MINOR.PATCH-<pre-release>.<number>
```

- **alpha**: Early development, unstable
- **beta**: Feature complete, testing phase
- **rc**: Release candidate, final testing

**Examples:**
- `0.1.0-alpha.1` - First alpha
- `0.2.0-beta.1` - First beta for 0.2.0
- `1.0.0-rc.1` - First release candidate for 1.0.0

#### Build Metadata

```
MAJOR.MINOR.PATCH+<build>
```

**Example:** `1.0.0+20250117.abc1234`

### Release Process

#### 1. Prepare Release

```bash
# Update main branch
git checkout main
git pull origin main

# Create release branch
git checkout -b release/v0.2.0
```

#### 2. Update Version Files

Update version in:
- `cmd/fintrack/main.go` (version variable)
- `CHANGELOG.md` (add release notes)
- `README.md` (update version badge if applicable)

```go
// cmd/fintrack/main.go
var version = "0.2.0"
```

#### 3. Update CHANGELOG.md

```markdown
## [0.2.0] - 2025-01-17

### Added
- Transaction tracking and categorization
- CSV import with bank-specific mappings
- Budget tracking with alerts

### Changed
- Improved account list performance
- Enhanced error messages

### Fixed
- Balance calculation overflow issue
- Config loading from environment variables

### Security
- Updated dependencies to address CVE-2024-XXXXX
```

#### 4. Create Release Commit

```bash
git add cmd/fintrack/main.go CHANGELOG.md README.md
git commit -m "chore(release): prepare v0.2.0 release

Update version to 0.2.0 and add changelog entries.
"
```

#### 5. Create PR for Release

```bash
git push -u origin release/v0.2.0
# Create PR: release/v0.2.0 → main
# Title: "Release v0.2.0"
```

#### 6. Tag the Release

After PR is merged:

```bash
git checkout main
git pull origin main

# Create annotated tag
git tag -a v0.2.0 -m "Release v0.2.0

New Features:
- Transaction tracking
- CSV import
- Budget alerts

See CHANGELOG.md for full details."

# Push tag
git push origin v0.2.0
```

#### 7. Create GitHub Release

Use GitHub's release interface:

1. Go to **Releases** → **Draft a new release**
2. Choose tag: `v0.2.0`
3. Release title: `v0.2.0 - Transaction Tracking & Budgets`
4. Description: Copy from CHANGELOG.md
5. Attach binaries (optional):
   - `fintrack-linux-amd64`
   - `fintrack-darwin-amd64`
   - `fintrack-darwin-arm64`
   - `fintrack-windows-amd64.exe`
6. Check "Set as the latest release"
7. Publish release

#### 8. Build & Publish Binaries (Optional)

```bash
# Build for all platforms
make build-all

# Upload to GitHub release
# (Can be automated via GitHub Actions)
```

### Automated Release Workflow

Create `.github/workflows/release.yml` for automated releases:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Build binaries
        run: make build-all

      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: |
            bin/fintrack-*
          generate_release_notes: true
```

## Code Review

### What Reviewers Look For

#### Functionality
- ✅ Does the code do what it's supposed to do?
- ✅ Are edge cases handled?
- ✅ Are error cases handled properly?
- ✅ Is the logic correct?

#### Code Quality
- ✅ Is the code readable and maintainable?
- ✅ Are variable/function names descriptive?
- ✅ Is the code well-structured?
- ✅ Are there code smells or anti-patterns?
- ✅ Is there unnecessary complexity?

#### Testing
- ✅ Are there sufficient tests?
- ✅ Do tests cover happy paths and edge cases?
- ✅ Are tests readable and maintainable?
- ✅ Does coverage meet threshold?

#### Documentation
- ✅ Are complex parts commented?
- ✅ Is the README updated if needed?
- ✅ Are API docs updated?
- ✅ Is CHANGELOG.md updated?

#### Security
- ✅ No SQL injection vulnerabilities
- ✅ No command injection vulnerabilities
- ✅ No hardcoded secrets
- ✅ Proper input validation
- ✅ No security warnings from linters

#### Performance
- ✅ No obvious performance issues
- ✅ Efficient algorithms used
- ✅ No memory leaks
- ✅ Database queries optimized

### Review Guidelines

**For Reviewers:**
- Be constructive and respectful
- Explain the "why" behind suggestions
- Differentiate between "must fix" and "nice to have"
- Approve when satisfied, even if minor nits remain
- Use GitHub's suggestion feature for quick fixes

**For Authors:**
- Respond to all comments
- Ask for clarification if needed
- Don't take feedback personally
- Mark resolved conversations
- Thank reviewers for their time

### Review Checklist

```markdown
## Code Review Checklist

### Functionality
- [ ] Code does what it's supposed to do
- [ ] Edge cases handled
- [ ] Error handling is appropriate
- [ ] No obvious bugs

### Code Quality
- [ ] Code is readable and maintainable
- [ ] Follows Go best practices
- [ ] No code duplication
- [ ] Appropriate abstractions
- [ ] Consistent with existing codebase style

### Testing
- [ ] Tests are comprehensive
- [ ] Tests are readable
- [ ] Coverage meets threshold (60%+)
- [ ] All tests pass

### Documentation
- [ ] Code comments where needed
- [ ] README updated if needed
- [ ] API docs updated
- [ ] CHANGELOG.md updated

### Security
- [ ] No hardcoded credentials
- [ ] Input validation present
- [ ] No SQL injection risks
- [ ] No security warnings

### Performance
- [ ] No obvious performance issues
- [ ] Efficient algorithms
- [ ] Database queries optimized

### Additional
- [ ] Commit messages follow guidelines
- [ ] PR description is clear
- [ ] No unnecessary dependencies added
```

## Testing Requirements

### Minimum Requirements

- ✅ All existing tests pass
- ✅ Code coverage ≥ 60% overall
- ✅ New code has ≥ 70% coverage
- ✅ Critical paths have ≥ 90% coverage

### Test Types Required

**Unit Tests:**
- All new functions/methods
- All bug fixes
- All edge cases

**Integration Tests:**
- New features touching multiple components
- Database interactions
- API integrations

**E2E Tests (for major features):**
- Complete user workflows
- Critical business logic

### Test Naming Convention

```go
func TestFunctionName_Scenario(t *testing.T) {
    // Test implementation
}
```

**Examples:**
```go
func TestCreateAccount_ValidInput(t *testing.T)
func TestCreateAccount_DuplicateName(t *testing.T)
func TestCreateAccount_EmptyName_ReturnsError(t *testing.T)
func TestUpdateBalance_NegativeAmount(t *testing.T)
```

### Running Tests

```bash
# All tests
make test

# With race detection
make test-race

# With coverage
make test-coverage

# Watch mode
make test-watch

# Quality checks (all)
make quality
```

## Getting Help

- **Questions**: Open a [Discussion](https://github.com/yourusername/fintrack/discussions)
- **Bugs**: Open an [Issue](https://github.com/yourusername/fintrack/issues)
- **Security**: See [SECURITY.md](SECURITY.md)
- **Chat**: Join our [Discord/Slack] (if available)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to FinTrack! 🎉
