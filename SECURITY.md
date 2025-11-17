# Security Best Practices for FinTrack

This document outlines what data should be committed to git and what should remain local to protect sensitive information.

## ✅ Safe to Commit

These files/changes contain NO sensitive data and should be committed:

### Code Changes
- **Source code** (`*.go` files) - Application logic, no credentials
- **Database migrations** (`migrations/*.sql`) - Schema only, no actual data
- **Build configuration** (`Makefile`, `go.mod`, `go.sum`)
- **Documentation** (`.claude/CLAUDE.md`, `README.md`, `docs/*`)
- **Example configs** (`fintrack_config.example.yaml`) - Templates without real credentials

### Current Pending Changes
```bash
# These changes are SAFE to commit:
modified:   cmd/fintrack/main.go          # Adds db.Init() call
modified:   internal/commands/account.go  # Removes unused import
modified:   go.mod                        # Dependency updates
new file:   .claude/CLAUDE.md             # Documentation
new file:   go.sum                        # Dependency checksums
```

## 🔒 NEVER Commit (PII/Sensitive Data)

These contain personal or sensitive information and must stay local:

### User Configuration
- `~/.config/fintrack/config.yaml` - Contains database credentials
  - **Location:** Outside repository (in user's home directory)
  - **Already protected:** ✅ Not tracked by git
  - **Contains:** Database passwords, connection strings

### Database Data
- PostgreSQL database contents - Thomas's actual accounts and transactions
  - **Location:** Docker container `postgres-dev` or separate PostgreSQL server
  - **Already protected:** ✅ Outside repository
  - **Contains:** All financial data (balances, transactions, etc.)

### Build Artifacts
- `bin/*` - Compiled binaries
  - **Already protected:** ✅ Listed in `.gitignore`
- `coverage.out`, `coverage.html` - Test coverage reports
  - **Already protected:** ✅ Listed in `.gitignore`

### Local Development Files
- `.air.toml` - Local dev server config
- `*.log` - Application logs (may contain sensitive data)
- `tmp/`, `temp/` - Temporary files

## 🛡️ Current .gitignore Protection

The repository is already configured to ignore:

```gitignore
# Sensitive configuration
config.yaml              # Main config file
*.local.yaml            # Local overrides

# Database files (if using SQLite locally)
*.db
*.sqlite
*.sqlite3

# Build artifacts
bin/
dist/
*.exe

# Logs (may contain sensitive data)
*.log

# IDE files (may contain local paths)
.idea/
.vscode/
```

## 🔐 Security Recommendations

### 1. Configuration Management

**Current Setup (Recommended):**
```bash
# User config with credentials
~/.config/fintrack/config.yaml          # ✅ Outside repo, not tracked

# Repository example (no credentials)
fintrack_config.example.yaml            # ✅ Committed as template
```

**Best Practice:**
```yaml
# ~/.config/fintrack/config.yaml (NEVER commit)
database:
  url: "postgresql://thomas:thomas123@localhost:5432/fintrack"  # Real credentials

# fintrack_config.example.yaml (Safe to commit)
database:
  url: "postgresql://USERNAME:PASSWORD@localhost:5432/fintrack"  # Placeholders
```

### 2. Environment Variables (Alternative)

For production/CI/CD, use environment variables instead:

```bash
# Set in shell, not committed to git
export FINTRACK_DB_URL="postgresql://user:pass@host:5432/fintrack"
export FINTRACK_DB_PASSWORD="secret"

# Run application (reads from environment)
fintrack account list
```

### 3. Database Credentials Rotation

If you accidentally commit credentials:

1. **Immediately rotate them:**
   ```bash
   docker exec postgres-dev psql -U devuser -d fintrack -c "ALTER USER thomas WITH PASSWORD 'new-secure-password';"
   ```

2. **Remove from git history:**
   ```bash
   # If credentials were committed, use git filter-branch or BFG Repo-Cleaner
   # This is complex - prevention is better!
   ```

3. **Update your local config:**
   ```bash
   # Edit ~/.config/fintrack/config.yaml with new password
   ```

### 4. Pre-Commit Checklist

Before committing, verify:

- [ ] No `config.yaml` files with real credentials
- [ ] No database dumps or exports
- [ ] No log files with transaction data
- [ ] No hardcoded passwords in code
- [ ] Check `git diff` for sensitive data
- [ ] Run `git status` to see untracked files

### 5. Code Review for Sensitive Data

When reviewing code changes, look for:

```go
// ❌ BAD - Hardcoded credentials
db.Connect("postgresql://thomas:thomas123@localhost/fintrack")

// ✅ GOOD - From config
db.Connect(config.Get().GetDatabaseURL())
```

## 📝 What Gets Committed vs. Local Only

### Repository (Git)
```
fintrack/
├── .claude/
│   └── CLAUDE.md                    ✅ Commit (documentation)
├── cmd/                             ✅ Commit (code)
├── internal/                        ✅ Commit (code)
├── migrations/                      ✅ Commit (schema only, no data)
├── docs/                            ✅ Commit (documentation)
├── tests/                           ✅ Commit (test code)
├── fintrack_config.example.yaml    ✅ Commit (template without credentials)
├── go.mod, go.sum                  ✅ Commit (dependencies)
├── Makefile                         ✅ Commit (build scripts)
└── README.md                        ✅ Commit (documentation)
```

### Local Only (Not in Git)
```
~/.config/fintrack/
└── config.yaml                      🔒 LOCAL ONLY (has real credentials)

Docker Container:
└── fintrack database                🔒 LOCAL ONLY (Thomas's financial data)

Build Directory:
fintrack/bin/
└── fintrack                         🔒 LOCAL ONLY (binary, ignored by git)
```

## 🚨 Emergency: Credential Leak Response

If you accidentally commit sensitive data:

1. **Immediate Actions:**
   - Don't push to remote (if not pushed yet)
   - Rotate all credentials immediately
   - Contact your team/security officer

2. **Clean Git History:**
   ```bash
   # If only in latest commit and not pushed
   git reset --soft HEAD~1
   git restore --staged path/to/sensitive/file

   # If pushed, you may need to force-push (coordinate with team)
   # Use tools like BFG Repo-Cleaner for deep history cleaning
   ```

3. **Verify Protection:**
   ```bash
   # Check what would be committed
   git add -n .

   # Check for sensitive patterns
   git grep -i "password"
   git grep -i "secret"
   git grep -i "thomas123"
   ```

## 📊 Security Checklist Summary

| Data Type | Location | Git Status | Contains PII? |
|-----------|----------|------------|---------------|
| Source code | `internal/`, `cmd/` | ✅ Committed | No |
| CLAUDE.md | `.claude/CLAUDE.md` | ✅ Committed | No |
| Config template | `fintrack_config.example.yaml` | ✅ Committed | No (placeholders) |
| User config | `~/.config/fintrack/config.yaml` | 🔒 Local only | Yes (credentials) |
| Database data | Docker/PostgreSQL | 🔒 Local only | Yes (financial data) |
| Build binaries | `bin/` | 🔒 Ignored | No |
| Logs | `*.log` | 🔒 Ignored | Maybe |

## 🎯 Current Status

**Safe to commit right now:**
```bash
git add cmd/fintrack/main.go
git add internal/commands/account.go
git add go.mod go.sum
git add .claude/CLAUDE.md
git commit -m "Fix: Initialize database connection in main.go

- Add db.Init() call to PersistentPreRunE
- Remove unused time import from account.go
- Add CLAUDE.md documentation for future development
- Update dependencies (go.mod, go.sum)
"
```

**Remains local (protected):**
- `~/.config/fintrack/config.yaml` - Contains Thomas's DB credentials
- Docker container `postgres-dev` - Contains Thomas's account data ($3,250.75 checking, $15,000 savings, etc.)
- `bin/fintrack` - Compiled binary (ignored by git)
