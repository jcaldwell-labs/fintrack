# Database Migrations

This directory contains SQL migration files for the FinTrack database schema.

## Migration Strategy

FinTrack uses a dual approach for database migrations:

1. **Development**: GORM AutoMigrate for rapid iteration
2. **Production**: Versioned SQL migration files for controlled deployments

## Migration Files

### Current Migrations

- `001_initial_schema.sql` - Initial database schema with all core tables

### File Naming Convention

Migration files follow the format: `{version}_{description}.sql`

- **Version**: 3-digit sequential number (001, 002, 003...)
- **Description**: Snake_case description of the change
- **Examples**:
  - `002_add_tags_to_transactions.sql`
  - `003_create_budget_alerts_table.sql`

## Running Migrations

### Development (GORM AutoMigrate)

The application automatically migrates the database schema on startup when running in development mode.

```go
// In cmd/fintrack/main.go
if err := db.AutoMigrate(); err != nil {
    log.Fatalf("Failed to migrate database: %v", err)
}
```

**Models migrated:**
- Account
- Category
- Transaction
- Budget
- RecurringItem
- Reminder
- CashFlowProjection
- ImportHistory

### Production (SQL Migrations)

For production environments, use migration tools like `golang-migrate` or apply SQL files directly.

#### Using golang-migrate

```bash
# Install
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run all migrations up
migrate -path migrations \
  -database "postgresql://user:pass@localhost:5432/fintrack?sslmode=disable" \
  up

# Run one migration up
migrate -path migrations \
  -database "postgresql://user:pass@localhost:5432/fintrack?sslmode=disable" \
  up 1

# Roll back one migration
migrate -path migrations \
  -database "postgresql://user:pass@localhost:5432/fintrack?sslmode=disable" \
  down 1

# Show current version
migrate -path migrations \
  -database "postgresql://user:pass@localhost:5432/fintrack?sslmode=disable" \
  version
```

#### Manual Migration

```bash
# Run migration manually
psql -d fintrack -f migrations/001_initial_schema.sql

# Or using environment variable
psql $FINTRACK_DB_URL -f migrations/001_initial_schema.sql
```

## Creating New Migrations

### 1. Create SQL Files

Create both UP and DOWN migration files:

```bash
# Up migration (apply changes)
migrations/003_add_user_preferences.sql

# Down migration (revert changes)
migrations/003_add_user_preferences.down.sql
```

### 2. Write Migration SQL

**Up Migration (003_add_user_preferences.sql):**
```sql
-- Add user preferences table
CREATE TABLE user_preferences (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    theme VARCHAR(50) DEFAULT 'light',
    date_format VARCHAR(20) DEFAULT 'YYYY-MM-DD',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_preferences_user_id ON user_preferences(user_id);
```

**Down Migration (003_add_user_preferences.down.sql):**
```sql
-- Revert user preferences table
DROP INDEX IF EXISTS idx_user_preferences_user_id;
DROP TABLE IF EXISTS user_preferences;
```

### 3. Update GORM Models

If adding new tables, create corresponding models in `internal/models/models.go` and update the AutoMigrate call in `internal/db/connection.go`.

### 4. Test Migration

```bash
# Test up migration
psql -d fintrack_test -f migrations/003_add_user_preferences.sql

# Verify schema
psql -d fintrack_test -c "\d user_preferences"

# Test down migration
psql -d fintrack_test -f migrations/003_add_user_preferences.down.sql

# Verify removal
psql -d fintrack_test -c "\d user_preferences"
```

## Migration Best Practices

### DO:
- ✅ Keep migrations small and focused
- ✅ Test migrations on a copy of production data
- ✅ Include both UP and DOWN migrations
- ✅ Add comments explaining complex changes
- ✅ Create indexes for foreign keys
- ✅ Use transactions where appropriate
- ✅ Backup database before running migrations in production

### DON'T:
- ❌ Modify existing migration files after they're deployed
- ❌ Delete or rename old migration files
- ❌ Include data migrations in schema migrations (create separate files)
- ❌ Use database-specific features without documenting alternatives
- ❌ Run migrations without testing first

## Schema Version Tracking

### Using golang-migrate

golang-migrate automatically creates a `schema_migrations` table to track applied migrations:

```sql
SELECT * FROM schema_migrations;
```

### Manual Tracking

If not using a migration tool, create a tracking table:

```sql
CREATE TABLE IF NOT EXISTS schema_version (
    version INTEGER PRIMARY KEY,
    description VARCHAR(255),
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Track applied migration
INSERT INTO schema_version (version, description)
VALUES (1, 'Initial schema');
```

## Rollback Strategy

### Automated Rollback (golang-migrate)

```bash
# Rollback last migration
migrate -path migrations -database $DB_URL down 1

# Rollback to specific version
migrate -path migrations -database $DB_URL goto 2
```

### Manual Rollback

1. Backup current database
2. Apply down migration SQL
3. Verify data integrity
4. Update schema version table

## Migration in CI/CD

### GitHub Actions Example

```yaml
- name: Run Database Migrations
  env:
    DATABASE_URL: ${{ secrets.DATABASE_URL }}
  run: |
    migrate -path migrations -database $DATABASE_URL up
```

### Safety Checks

Before running migrations in production:

1. **Backup**: Create database backup
2. **Test**: Run migrations on staging environment
3. **Review**: Check migration SQL for destructive operations
4. **Monitor**: Watch for errors during migration
5. **Verify**: Confirm data integrity after migration

## Troubleshooting

### Migration Fails Midway

1. Check error message in migration tool output
2. Fix SQL syntax or logic error
3. If needed, manually clean up partial changes
4. Re-run migration

### Schema Drift (Development vs Production)

If GORM AutoMigrate creates different schema than SQL migrations:

1. Generate SQL from GORM models:
```go
// Use GORM's migrator to see what it would create
db.Migrator().CreateTable(&models.Transaction{})
```

2. Update SQL migration to match
3. Test both AutoMigrate and SQL migration produce same schema

### Migration Version Conflicts

If multiple developers create migrations with same version:

1. Rename newer migration to next available version
2. Update references in down migration
3. Commit both changes

## Resources

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [GORM Migration Guide](https://gorm.io/docs/migration.html)
- [PostgreSQL ALTER TABLE](https://www.postgresql.org/docs/current/sql-altertable.html)
- [Database Migration Best Practices](https://www.dbcodestrategies.com/)

## Support

For migration issues or questions:
- Check existing migrations for examples
- Review GORM model definitions
- Consult PostgreSQL documentation
- Open an issue on GitHub
