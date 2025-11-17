# Workflow Quick Reference

Quick reference guide for common development workflows.

## 🚀 Quick Start Development

```bash
# Clone and setup
git clone https://github.com/yourusername/fintrack.git
cd fintrack
make deps
make build

# Create database
createdb fintrack
psql -d fintrack -f migrations/001_initial_schema.sql

# Configure
cp fintrack_config.example.yaml ~/.config/fintrack/config.yaml
# Edit config with your database credentials

# Test
make test
```

## 📝 Commit Messages

```bash
# Format
<type>(<scope>): <subject>

# Examples
git commit -m "feat(commands): add CSV import"
git commit -m "fix(db): resolve connection leak"
git commit -m "docs(readme): update install instructions"
git commit -m "test(repo): add integration tests"
```

**Types:** `feat` `fix` `docs` `style` `refactor` `perf` `test` `chore` `ci` `revert`

**Scopes:** `commands` `db` `config` `models` `output` `cli`

## 🌿 Branch Workflow

```bash
# Start new feature
git checkout main
git pull origin main
git checkout -b feature/my-feature

# Make changes, test, commit

# Push and create PR
git push -u origin feature/my-feature
# Create PR on GitHub

# After merge, cleanup
git checkout main
git pull origin main
git branch -d feature/my-feature
```

**Branch naming:** `feature/` `fix/` `docs/` `refactor/` `test/` `chore/`

## ✅ Pre-PR Checklist

```bash
# Run all checks
make quality

# Or individually
make fmt               # Format code
make lint              # Lint check
make test              # Run tests
make test-race         # Race detection
make test-coverage-check  # Coverage threshold
```

**Before creating PR:**
- [ ] All tests pass
- [ ] Coverage ≥ 60%
- [ ] Code formatted
- [ ] No lint errors
- [ ] Documentation updated
- [ ] CHANGELOG.md updated

## 🎯 Testing Commands

```bash
make test                   # Run all tests
make test-race              # With race detector
make test-unit              # Unit tests only
make test-integration       # Integration tests only
make test-coverage          # Generate coverage report
make test-coverage-check    # Check 60% threshold
make test-watch             # Watch mode (requires entr)
make benchmark              # Run benchmarks
```

## 📦 Release Process

```bash
# 1. Create release branch
git checkout main
git pull origin main
git checkout -b release/v0.2.0

# 2. Update version
# Edit cmd/fintrack/main.go: version = "0.2.0"
# Edit CHANGELOG.md: Add release notes
# Edit README.md: Update version badge

# 3. Commit and PR
git add cmd/fintrack/main.go CHANGELOG.md README.md
git commit -m "chore(release): prepare v0.2.0"
git push -u origin release/v0.2.0
# Create PR on GitHub

# 4. After merge, tag
git checkout main
git pull origin main
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0

# GitHub Actions automatically:
# - Builds binaries for all platforms
# - Creates checksums
# - Creates GitHub release
# - Uploads artifacts
```

## 🏷️ Version Numbers

**Format:** `MAJOR.MINOR.PATCH`

- **MAJOR** - Breaking changes (1.0.0 → 2.0.0)
- **MINOR** - New features (1.0.0 → 1.1.0)
- **PATCH** - Bug fixes (1.0.0 → 1.0.1)

**Pre-releases:**
- `0.1.0-alpha.1` - Early development
- `0.2.0-beta.1` - Feature complete
- `1.0.0-rc.1` - Release candidate

## 🔍 Code Review

**What to check:**
- ✅ Functionality correct
- ✅ Edge cases handled
- ✅ Tests sufficient
- ✅ Code readable
- ✅ Documentation updated
- ✅ No security issues

**Response times:**
- Initial response: 2 business days
- Full review: 5 business days
- Urgent fixes: 1 business day

## 🐛 Reporting Issues

**Bug reports:** Use `.github/ISSUE_TEMPLATE/bug_report.md`
**Features:** Use `.github/ISSUE_TEMPLATE/feature_request.md`

Include:
- Clear description
- Steps to reproduce
- Expected vs actual behavior
- Environment details
- Logs (with credentials removed)

## 📚 Documentation Files

- **README.md** - Project overview and quick start
- **CONTRIBUTING.md** - Full contributor guide
- **TESTING.md** - Testing strategy
- **SECURITY.md** - Security best practices
- **CHANGELOG.md** - Version history
- **.claude/CLAUDE.md** - Development context

## 🔒 Security

**Never commit:**
- ❌ Passwords or credentials
- ❌ Database dumps with real data
- ❌ API keys or tokens
- ❌ ~/.config/fintrack/config.yaml

**Always:**
- ✅ Use config files outside repo
- ✅ Use environment variables
- ✅ Review git diff before commit
- ✅ Check git status for untracked files

## 🤝 Getting Help

- **Questions:** [GitHub Discussions](https://github.com/yourusername/fintrack/discussions)
- **Bugs:** [Issues](https://github.com/yourusername/fintrack/issues)
- **Security:** See [SECURITY.md](../SECURITY.md)

## ⚡ Common Commands

```bash
# Development
make dev                    # Live reload
make run                    # Build and run
make build                  # Build binary
make install                # Install to /usr/local/bin

# Quality
make fmt                    # Format code
make lint                   # Lint code
make quality                # All checks

# Dependencies
make deps                   # Download dependencies
make verify                 # Verify dependencies

# Cross-platform
make build-all              # Build for all platforms
make build-linux            # Linux only
make build-darwin           # macOS only
make build-windows          # Windows only

# Cleanup
make clean                  # Remove build artifacts
```

## 📋 PR Template Checklist

When creating a PR:
- [ ] Description clear
- [ ] Type of change selected
- [ ] Related issues linked
- [ ] Changes listed
- [ ] Tests added/updated
- [ ] Manual testing done
- [ ] Screenshots/output included
- [ ] Security considered
- [ ] Performance impact noted
- [ ] Breaking changes documented
- [ ] CHANGELOG updated

---

For full details, see [CONTRIBUTING.md](../CONTRIBUTING.md)
