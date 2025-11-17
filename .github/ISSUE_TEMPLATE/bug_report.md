---
name: Bug Report
about: Report a bug to help us improve
title: '[BUG] '
labels: bug
assignees: ''
---

## Bug Description

<!-- A clear and concise description of what the bug is -->

## Steps to Reproduce

1. Run command '...'
2. With configuration '...'
3. See error

## Expected Behavior

<!-- A clear and concise description of what you expected to happen -->

## Actual Behavior

<!-- What actually happened -->

```
<!-- Paste error message or unexpected output here -->
```

## Environment

**FinTrack Version:**
```bash
$ fintrack --version
```

**OS:**
- [ ] Linux (specify distro and version)
- [ ] macOS (specify version)
- [ ] Windows (specify version)
- [ ] WSL (specify version)

**Go Version:**
```bash
$ go version
```

**Database:**
- Database: PostgreSQL
- Version:
- Connection method: (local/Docker/remote)

**Configuration:**
<!-- Paste relevant config (REMOVE ANY CREDENTIALS!) -->
```yaml
# ~/.config/fintrack/config.yaml
database:
  # ... (redact credentials)
```

## Additional Context

<!-- Add any other context about the problem here -->

### Logs

<!-- If applicable, add relevant log output -->
```
<!-- Paste logs here -->
```

### Screenshots

<!-- If applicable, add screenshots to help explain your problem -->

## Possible Solution

<!-- If you have ideas on how to fix this, please share -->

## Checklist

- [ ] I have checked existing issues for duplicates
- [ ] I have provided all required information
- [ ] I have removed any sensitive data from logs/config
- [ ] I have tested with the latest version
