# my-context Integration Pattern for FinTrack

> Design document for integrating FinTrack with my-context for automatic context tracking during financial operations.

**Status**: Design Document (not yet implemented)
**Created**: 2025-11-25
**Author**: Phase 3b Session 1

---

## Table of Contents

1. [Overview](#1-overview)
2. [Integration Architecture](#2-integration-architecture)
3. [API Design](#3-api-design)
4. [Example Usage](#4-example-usage)
5. [Implementation Roadmap](#5-implementation-roadmap)
6. [Testing Strategy](#6-testing-strategy)
7. [Future Enhancements](#7-future-enhancements)

---

## 1. Overview

### Why Integrate FinTrack with my-context?

FinTrack is a terminal-based personal finance tracking application. Users often work on financial tasks in sessions: reconciling bank statements, categorizing transactions, setting up budgets, or reviewing spending patterns. These sessions benefit from context tracking for several reasons:

1. **Audit Trail**: Track what financial operations were performed and when
2. **Session Continuity**: Resume complex tasks (like reconciliation) where you left off
3. **Decision Documentation**: Record why certain categorizations or budget adjustments were made
4. **Workflow Optimization**: Analyze patterns in how you manage finances

### What Problems Does This Solve?

| Problem | Current State | With my-context Integration |
|---------|--------------|----------------------------|
| Lost context | "What was I reconciling?" | Session notes track progress |
| No decision history | "Why did I categorize this as X?" | Notes explain reasoning |
| Interrupted workflows | Start over each session | Resume from saved state |
| No operational logging | Manual tracking required | Automatic operation logging |

### Who Benefits?

- **Power users**: Track complex multi-session financial tasks
- **Auditors**: Review history of financial operations
- **Developers**: Debug issues with operational context
- **Teams**: Share context about shared financial management

### Integration Philosophy

Following Unix philosophy principles shared by both projects:

- **Composable**: my-context integration is optional, FinTrack works standalone
- **Non-intrusive**: Operations work the same with or without context tracking
- **Text-based**: All tracked data remains human-readable
- **Scriptable**: Integration works in both interactive and scripted usage

---

## 2. Integration Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        FinTrack CLI                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │ account  │  │category  │  │transaction│  │  budget  │    │
│  │ commands │  │ commands │  │ commands  │  │ commands │    │
│  └────┬─────┘  └────┬─────┘  └────┬──────┘  └────┬─────┘    │
│       │             │             │              │           │
│       └─────────────┴──────┬──────┴──────────────┘           │
│                            │                                 │
│                    ┌───────▼───────┐                        │
│                    │ Context Tracker│ (internal/context/)    │
│                    │   (optional)   │                        │
│                    └───────┬───────┘                        │
└────────────────────────────┼────────────────────────────────┘
                             │
                    ┌────────▼────────┐
                    │   my-context    │
                    │  (external CLI) │
                    └─────────────────┘
```

### Component Responsibilities

#### internal/context/tracker.go
- Wraps my-context CLI operations
- Manages session lifecycle (start, stop, resume)
- Provides consistent note formatting
- Handles graceful degradation when my-context unavailable

#### internal/context/config.go
- Configuration for context tracking behavior
- Enable/disable tracking
- Session naming patterns
- Auto-start preferences

#### internal/context/hooks.go
- Pre/post operation hooks for commands
- Automatic note generation for operations
- File association for affected records

### Where Changes Would Go

```
fintrack/
├── internal/
│   ├── context/              # NEW: Context tracking package
│   │   ├── tracker.go        # Core tracker implementation
│   │   ├── config.go         # Configuration management
│   │   ├── hooks.go          # Command hooks
│   │   └── tracker_test.go   # Unit tests
│   ├── commands/
│   │   ├── account.go        # Add context hooks
│   │   ├── transaction.go    # Add context hooks
│   │   ├── category.go       # Add context hooks
│   │   └── root.go           # Add --context flag
│   └── config/
│       └── config.go         # Add context config section
├── cmd/fintrack/
│   └── main.go               # Initialize context tracker
└── go.mod                    # No new dependencies (uses CLI)
```

### Call Flow

```
User runs: fintrack tx add -50 --account "Checking" --category "Groceries"

1. main.go initializes Tracker (if enabled)
2. transaction.go receives command
3. Pre-hook: tracker.Note("Transaction: add -50.00 to Checking")
4. Execute: Create transaction in database
5. Post-hook: tracker.File("~/.config/fintrack/data/transactions.db")
6. Return result to user
```

### Configuration Flow

```yaml
# ~/.config/fintrack/config.yaml
context:
  enabled: true                    # Enable/disable integration
  auto_start: false                # Auto-start session on first command
  session_prefix: "fintrack"       # Prefix for session names
  track_operations:
    - transaction.add              # Which operations to track
    - transaction.update
    - account.add
    - budget.set
```

---

## 3. API Design

### Core Tracker Interface

```go
// internal/context/tracker.go
package context

import (
    "fmt"
    "os/exec"
    "strings"
    "time"
)

// Tracker manages my-context operations for FinTrack
type Tracker struct {
    enabled     bool
    sessionName string
    config      *Config
}

// Config holds context tracking configuration
type Config struct {
    Enabled       bool     `yaml:"enabled"`
    AutoStart     bool     `yaml:"auto_start"`
    SessionPrefix string   `yaml:"session_prefix"`
    TrackOps      []string `yaml:"track_operations"`
}

// NewTracker creates a new context tracker
func NewTracker(cfg *Config) (*Tracker, error) {
    if cfg == nil {
        cfg = &Config{Enabled: false}
    }

    t := &Tracker{
        enabled: cfg.Enabled,
        config:  cfg,
    }

    // Verify my-context is available
    if t.enabled {
        if err := t.checkMyContext(); err != nil {
            // Graceful degradation: disable if not available
            t.enabled = false
            return t, nil // Not an error, just disabled
        }
    }

    return t, nil
}

// checkMyContext verifies my-context CLI is available
func (t *Tracker) checkMyContext() error {
    cmd := exec.Command("my-context", "show", "--json")
    return cmd.Run()
}

// IsEnabled returns whether context tracking is active
func (t *Tracker) IsEnabled() bool {
    return t.enabled
}
```

### Session Management

```go
// Start begins a new context session
func (t *Tracker) Start(name string) error {
    if !t.enabled {
        return nil
    }

    sessionName := fmt.Sprintf("%s-%s", t.config.SessionPrefix, name)

    cmd := exec.Command("my-context", "start", sessionName)
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to start context: %w", err)
    }

    t.sessionName = sessionName
    return nil
}

// Stop ends the current context session
func (t *Tracker) Stop() error {
    if !t.enabled || t.sessionName == "" {
        return nil
    }

    cmd := exec.Command("my-context", "stop")
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to stop context: %w", err)
    }

    t.sessionName = ""
    return nil
}

// Resume continues an existing session
func (t *Tracker) Resume(name string) error {
    if !t.enabled {
        return nil
    }

    // my-context automatically resumes if session exists
    return t.Start(name)
}
```

### Note Operations

```go
// Note adds a timestamped note to the active context
func (t *Tracker) Note(message string) error {
    if !t.enabled {
        return nil
    }

    cmd := exec.Command("my-context", "note", message)
    return cmd.Run()
}

// NoteOperation formats and logs an operation
func (t *Tracker) NoteOperation(op string, details map[string]interface{}) error {
    if !t.enabled {
        return nil
    }

    // Format: "Operation: transaction.add | amount=-50.00, account=Checking"
    var parts []string
    for k, v := range details {
        parts = append(parts, fmt.Sprintf("%s=%v", k, v))
    }

    message := fmt.Sprintf("Operation: %s | %s", op, strings.Join(parts, ", "))
    return t.Note(message)
}

// NoteDecision records a decision with reasoning
func (t *Tracker) NoteDecision(decision, reason string) error {
    if !t.enabled {
        return nil
    }

    message := fmt.Sprintf("Decision: %s\nReason: %s", decision, reason)
    return t.Note(message)
}
```

### File Association

```go
// File associates a file with the current context
func (t *Tracker) File(path string) error {
    if !t.enabled {
        return nil
    }

    cmd := exec.Command("my-context", "file", path)
    return cmd.Run()
}

// FileDatabase associates the FinTrack database file
func (t *Tracker) FileDatabase() error {
    // Get database path from config
    dbPath := "~/.config/fintrack/fintrack.db" // Example
    return t.File(dbPath)
}
```

### Error Handling

```go
// Error handling approach: graceful degradation

// WrapOperation executes an operation with context tracking
func (t *Tracker) WrapOperation(opName string, fn func() error) error {
    // Pre-hook (ignore errors - don't fail operation due to tracking)
    _ = t.Note(fmt.Sprintf("Starting: %s", opName))

    // Execute actual operation
    err := fn()

    // Post-hook
    if err != nil {
        _ = t.Note(fmt.Sprintf("Failed: %s - %v", opName, err))
    } else {
        _ = t.Note(fmt.Sprintf("Completed: %s", opName))
    }

    return err // Return original error, not tracking errors
}
```

---

## 4. Example Usage

### Example 1: Command-Line Usage with Manual Context

```bash
# Start a fintrack work session
my-context start "fintrack-reconciliation-2025-01"

# Perform fintrack operations
fintrack tx list --account "Checking" --start 2025-01-01
my-context note "Reviewing January checking transactions"

fintrack tx update 42 --category "Groceries"
my-context note "Recategorized tx #42: was Misc, now Groceries (weekly shopping)"

fintrack tx add -150.00 --account "Checking" --category "Utilities" --payee "Electric Co"
my-context note "Added missing utility payment"

# End session
my-context stop

# Export session for records
my-context export "fintrack-reconciliation-2025-01" --to ~/finance-notes/
```

### Example 2: Automatic Context Tracking (Integrated)

```bash
# Enable context tracking in config
cat >> ~/.config/fintrack/config.yaml << EOF
context:
  enabled: true
  auto_start: true
  session_prefix: "fintrack"
EOF

# Now operations are automatically tracked
fintrack tx add -50.00 --account "Checking" --category "Groceries" --payee "Walmart"
# Automatically logs: "Operation: transaction.add | amount=-50.00, account=Checking, category=Groceries"

fintrack budget set "Groceries" 500.00 --month 2025-01
# Automatically logs: "Operation: budget.set | category=Groceries, amount=500.00, month=2025-01"

# View what was tracked
my-context show
```

### Example 3: Programmatic Usage in Go Code

```go
// internal/commands/transaction.go
package commands

import (
    "github.com/fintrack/fintrack/internal/context"
    "github.com/spf13/cobra"
)

var tracker *context.Tracker

func init() {
    // Tracker initialized in main.go
}

func NewTransactionAddCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "add",
        Short: "Add a new transaction",
        RunE: func(cmd *cobra.Command, args []string) error {
            amount, _ := cmd.Flags().GetFloat64("amount")
            account, _ := cmd.Flags().GetString("account")
            category, _ := cmd.Flags().GetString("category")
            payee, _ := cmd.Flags().GetString("payee")

            // Track the operation
            return tracker.WrapOperation("transaction.add", func() error {
                // Log operation details
                tracker.NoteOperation("transaction.add", map[string]interface{}{
                    "amount":   amount,
                    "account":  account,
                    "category": category,
                    "payee":    payee,
                })

                // Actual transaction creation
                tx := &models.Transaction{
                    Amount:   amount,
                    Account:  account,
                    Category: category,
                    Payee:    payee,
                }

                return db.CreateTransaction(tx)
            })
        },
    }

    return cmd
}
```

### Example 4: Session-Based Workflow

```go
// Example of a reconciliation workflow with context tracking
package workflows

import (
    "fmt"
    "github.com/fintrack/fintrack/internal/context"
)

func ReconcileAccount(tracker *context.Tracker, accountID int, statementBalance float64) error {
    // Start a reconciliation session
    sessionName := fmt.Sprintf("reconcile-account-%d-%s", accountID, time.Now().Format("2006-01-02"))
    if err := tracker.Start(sessionName); err != nil {
        return err
    }
    defer tracker.Stop()

    // Track the starting point
    tracker.Note(fmt.Sprintf("Starting reconciliation for account %d", accountID))
    tracker.Note(fmt.Sprintf("Statement balance: $%.2f", statementBalance))

    // Get current balance
    currentBalance, err := db.GetAccountBalance(accountID)
    if err != nil {
        tracker.Note(fmt.Sprintf("Error getting balance: %v", err))
        return err
    }

    tracker.Note(fmt.Sprintf("Current balance: $%.2f", currentBalance))
    tracker.Note(fmt.Sprintf("Difference: $%.2f", statementBalance-currentBalance))

    // If balanced, record success
    if currentBalance == statementBalance {
        tracker.Note("Account reconciled successfully - balances match")
        return nil
    }

    // Otherwise, note the discrepancy for investigation
    tracker.NoteDecision(
        "Account has discrepancy",
        fmt.Sprintf("Statement shows $%.2f but records show $%.2f", statementBalance, currentBalance),
    )

    return nil
}
```

---

## 5. Implementation Roadmap

### Phase 1: Core Integration (PR 1)

**Scope**: Basic tracker implementation

**Files**:
- `internal/context/tracker.go` - Core tracker with Start/Stop/Note
- `internal/context/config.go` - Configuration structure
- `internal/context/tracker_test.go` - Unit tests

**Effort**: 2-3 hours

**Prerequisites**: None

**Acceptance Criteria**:
- Tracker can start/stop my-context sessions
- Notes can be added programmatically
- Graceful degradation when my-context unavailable
- 80%+ test coverage

### Phase 2: Command Integration (PR 2)

**Scope**: Hook tracker into existing commands

**Files**:
- `internal/commands/root.go` - Add --context flag
- `internal/commands/transaction.go` - Add tracking hooks
- `internal/commands/account.go` - Add tracking hooks
- `cmd/fintrack/main.go` - Initialize tracker

**Effort**: 3-4 hours

**Prerequisites**: Phase 1 complete

**Acceptance Criteria**:
- --context flag enables/disables tracking
- Transaction add/update/delete tracked
- Account operations tracked
- No performance degradation

### Phase 3: Configuration (PR 3)

**Scope**: User configuration for tracking behavior

**Files**:
- `internal/config/config.go` - Add context section
- Update example config file
- Update README with configuration docs

**Effort**: 1-2 hours

**Prerequisites**: Phase 2 complete

**Acceptance Criteria**:
- Config file controls tracking behavior
- Can enable/disable per operation type
- Session naming customizable

### Phase 4: Advanced Features (PR 4)

**Scope**: Enhanced tracking capabilities

**Files**:
- `internal/context/hooks.go` - Pre/post operation hooks
- `internal/context/export.go` - Session export helpers
- Additional command integration

**Effort**: 2-3 hours

**Prerequisites**: Phase 3 complete

**Acceptance Criteria**:
- Automatic file association for database
- Export session with fintrack data summary
- Budget and recurring item tracking

### Dependencies

```
Phase 1 ─┬─▶ Phase 2 ─┬─▶ Phase 3 ─┬─▶ Phase 4
         │            │            │
         │            │            └─▶ (Optional enhancements)
         │            │
         │            └─▶ (Core functionality complete)
         │
         └─▶ (Can be released standalone)
```

---

## 6. Testing Strategy

### Unit Tests

```go
// internal/context/tracker_test.go
package context

import (
    "os"
    "os/exec"
    "testing"
)

func TestNewTracker_DisabledByDefault(t *testing.T) {
    tracker, err := NewTracker(nil)
    if err != nil {
        t.Fatalf("NewTracker returned error: %v", err)
    }

    if tracker.IsEnabled() {
        t.Error("Expected tracker to be disabled by default")
    }
}

func TestNewTracker_EnabledWithConfig(t *testing.T) {
    cfg := &Config{
        Enabled:       true,
        SessionPrefix: "test",
    }

    tracker, err := NewTracker(cfg)
    if err != nil {
        t.Fatalf("NewTracker returned error: %v", err)
    }

    // Will be disabled if my-context not installed
    // This is expected behavior (graceful degradation)
}

func TestTracker_NoteWhenDisabled(t *testing.T) {
    tracker := &Tracker{enabled: false}

    err := tracker.Note("test note")
    if err != nil {
        t.Errorf("Note should not error when disabled: %v", err)
    }
}

func TestTracker_NoteOperation_Format(t *testing.T) {
    // Test that NoteOperation formats correctly
    tracker := &Tracker{enabled: false} // Disabled so won't actually call my-context

    details := map[string]interface{}{
        "amount":   -50.00,
        "account":  "Checking",
        "category": "Groceries",
    }

    // This won't actually log, but we can test the format function
    err := tracker.NoteOperation("transaction.add", details)
    if err != nil {
        t.Errorf("NoteOperation should not error: %v", err)
    }
}
```

### Integration Tests

```go
// internal/context/integration_test.go
// +build integration

package context

import (
    "os"
    "os/exec"
    "testing"
)

func TestIntegration_FullWorkflow(t *testing.T) {
    // Skip if my-context not installed
    if _, err := exec.LookPath("my-context"); err != nil {
        t.Skip("my-context not installed, skipping integration test")
    }

    // Use temp directory for test contexts
    tmpDir := t.TempDir()
    os.Setenv("MY_CONTEXT_HOME", tmpDir)
    defer os.Unsetenv("MY_CONTEXT_HOME")

    cfg := &Config{
        Enabled:       true,
        SessionPrefix: "test-fintrack",
    }

    tracker, err := NewTracker(cfg)
    if err != nil {
        t.Fatalf("Failed to create tracker: %v", err)
    }

    // Start session
    if err := tracker.Start("integration-test"); err != nil {
        t.Fatalf("Failed to start session: %v", err)
    }

    // Add notes
    if err := tracker.Note("Test note 1"); err != nil {
        t.Errorf("Failed to add note: %v", err)
    }

    if err := tracker.NoteOperation("test.op", map[string]interface{}{"key": "value"}); err != nil {
        t.Errorf("Failed to add operation note: %v", err)
    }

    // Stop session
    if err := tracker.Stop(); err != nil {
        t.Errorf("Failed to stop session: %v", err)
    }

    // Verify session was created
    cmd := exec.Command("my-context", "list", "--json")
    output, err := cmd.Output()
    if err != nil {
        t.Errorf("Failed to list contexts: %v", err)
    }

    if !strings.Contains(string(output), "test-fintrack-integration-test") {
        t.Error("Session not found in context list")
    }
}
```

### Manual Testing Scenarios

#### Scenario 1: Basic Operation Tracking

```bash
# Setup
export MY_CONTEXT_HOME=/tmp/fintrack-test-context

# Test
fintrack --context tx add -50.00 --account "Test" --category "Test"

# Verify
my-context show
# Should show: "Operation: transaction.add | amount=-50.00, ..."
```

#### Scenario 2: Graceful Degradation

```bash
# Remove my-context temporarily
mv $(which my-context) /tmp/my-context-backup

# Test - should work without errors
fintrack --context tx list

# Restore
mv /tmp/my-context-backup $(which my-context)
```

#### Scenario 3: Session Persistence

```bash
# Start tracking session
fintrack context start "test-session"

# Perform operations
fintrack tx add -50.00 --account "Checking" --category "Food"
fintrack tx add -25.00 --account "Checking" --category "Transport"

# Stop and export
fintrack context stop
my-context export "fintrack-test-session"

# Verify export contains operation history
cat fintrack-test-session.md
```

---

## 7. Future Enhancements

### 7.1 Automatic Session Naming

Intelligent session names based on context:

```go
// Future: Auto-generate meaningful session names
func (t *Tracker) AutoStart() error {
    // Detect what user is doing
    // - If recent account commands: "account-management-{date}"
    // - If budget commands: "budget-planning-{date}"
    // - If many transaction adds: "data-entry-{date}"
    // - Default: "fintrack-{date}"
}
```

### 7.2 Rich Operation Metadata

Track more detailed information:

```go
type OperationMeta struct {
    Operation   string
    Timestamp   time.Time
    User        string
    Duration    time.Duration
    AffectedIDs []int
    Before      interface{} // State before operation
    After       interface{} // State after operation
}
```

### 7.3 Session Templates

Pre-defined session types for common workflows:

```yaml
# ~/.config/fintrack/context-templates.yaml
templates:
  reconciliation:
    prefix: "reconcile"
    auto_track: [transaction.update, transaction.add]
    prompts:
      - "Which account are you reconciling?"
      - "What is the statement balance?"

  budget_planning:
    prefix: "budget"
    auto_track: [budget.set, budget.update]
    prompts:
      - "What month are you planning?"
```

### 7.4 Cross-Session Analytics

Analyze patterns across sessions:

```bash
# Future command
fintrack context analyze --last 30d

Output:
- Sessions: 12
- Total tracked operations: 156
- Most common: transaction.add (89)
- Average session duration: 23 minutes
- Peak activity: Sundays (weekly reconciliation)
```

### 7.5 Integration with Other Tools

Potential integrations:

| Tool | Integration |
|------|-------------|
| Git | Auto-commit context export with finance data |
| Calendar | Link sessions to calendar events |
| Reports | Include context notes in financial reports |
| Backup | Include context data in fintrack backups |

### 7.6 Team Context Sharing

For shared financial management:

```go
// Future: Shared context for couples/families/teams
type SharedTracker struct {
    *Tracker
    SharedDir string // Network/cloud directory for shared contexts
}

func (st *SharedTracker) ShareNote(message string) error {
    // Add note visible to all team members
}
```

### 7.7 Hooks for External Notifications

```yaml
# Future: Trigger external actions on context events
context:
  hooks:
    on_session_end:
      - "notify-send 'FinTrack session ended'"
      - "~/scripts/backup-context.sh"
    on_budget_exceeded:
      - "mail -s 'Budget Alert' user@email.com"
```

---

## Appendix: Reference Links

- [my-context GitHub Repository](https://github.com/jcaldwell-labs/my-context)
- [FinTrack GitHub Repository](https://github.com/jcaldwell-labs/fintrack)
- [Cobra CLI Framework](https://github.com/spf13/cobra)
- [Viper Configuration](https://github.com/spf13/viper)

---

**Document Version**: 1.0
**Last Updated**: 2025-11-25
**Status**: Design Complete - Ready for Implementation
