-- FinTrack Database Schema
-- PostgreSQL 12+
-- Personal Finance Tracking Application
-- Version: 1.0
-- Created: 2025-11-16

-- Drop existing tables (for clean setup)
DROP TABLE IF EXISTS cash_flow_projections CASCADE;
DROP TABLE IF EXISTS import_history CASCADE;
DROP TABLE IF EXISTS reminders CASCADE;
DROP TABLE IF EXISTS recurring_items CASCADE;
DROP TABLE IF EXISTS budgets CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS accounts CASCADE;

-- ============================================================================
-- ACCOUNTS
-- ============================================================================
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('checking', 'savings', 'credit', 'cash', 'investment', 'loan')),
    currency VARCHAR(3) DEFAULT 'USD',
    initial_balance DECIMAL(15,2) DEFAULT 0,
    current_balance DECIMAL(15,2) DEFAULT 0,
    institution VARCHAR(255),
    account_number_last4 VARCHAR(4),
    is_active BOOLEAN DEFAULT true,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_accounts_active ON accounts(is_active);
CREATE INDEX idx_accounts_type ON accounts(type);
CREATE UNIQUE INDEX idx_accounts_name_active ON accounts(name) WHERE is_active = true;

COMMENT ON TABLE accounts IS 'Financial accounts (bank, credit card, cash, investment)';
COMMENT ON COLUMN accounts.type IS 'Account type: checking, savings, credit, cash, investment, loan';
COMMENT ON COLUMN accounts.current_balance IS 'Calculated balance based on transactions';

-- ============================================================================
-- CATEGORIES
-- ============================================================================
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('income', 'expense', 'transfer')),
    color VARCHAR(7),  -- Hex color code (e.g., #FF5733)
    icon VARCHAR(50),  -- Icon identifier for UI
    is_system BOOLEAN DEFAULT false,  -- System categories cannot be deleted
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_categories_parent ON categories(parent_id);
CREATE INDEX idx_categories_type ON categories(type);
CREATE UNIQUE INDEX idx_categories_name_type ON categories(name, type);

COMMENT ON TABLE categories IS 'Hierarchical transaction categories';
COMMENT ON COLUMN categories.parent_id IS 'Parent category for hierarchical organization (e.g., Groceries -> Food & Dining)';
COMMENT ON COLUMN categories.is_system IS 'System categories created by app, cannot be deleted';

-- ============================================================================
-- TRANSACTIONS
-- ============================================================================
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    payee VARCHAR(255),
    description TEXT,
    type VARCHAR(20) NOT NULL CHECK (type IN ('income', 'expense', 'transfer')),
    transfer_account_id INTEGER REFERENCES accounts(id) ON DELETE SET NULL,
    recurring_id INTEGER REFERENCES recurring_items(id) ON DELETE SET NULL,
    tags TEXT[],  -- PostgreSQL array for flexible tagging
    is_reconciled BOOLEAN DEFAULT false,
    reconciled_at TIMESTAMP,
    import_id INTEGER REFERENCES import_history(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_transactions_account ON transactions(account_id);
CREATE INDEX idx_transactions_date ON transactions(date DESC);
CREATE INDEX idx_transactions_category ON transactions(category_id);
CREATE INDEX idx_transactions_recurring ON transactions(recurring_id);
CREATE INDEX idx_transactions_tags ON transactions USING GIN(tags);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_payee ON transactions(payee);
CREATE INDEX idx_transactions_reconciled ON transactions(is_reconciled) WHERE is_reconciled = false;

COMMENT ON TABLE transactions IS 'Individual financial transactions';
COMMENT ON COLUMN transactions.amount IS 'Positive for income, negative for expenses';
COMMENT ON COLUMN transactions.transfer_account_id IS 'For transfer transactions, the destination account';
COMMENT ON COLUMN transactions.tags IS 'Array of tags for flexible categorization (e.g., ["business", "reimbursable"])';

-- ============================================================================
-- BUDGETS
-- ============================================================================
CREATE TABLE budgets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category_id INTEGER REFERENCES categories(id) ON DELETE CASCADE,
    period_type VARCHAR(20) NOT NULL CHECK (period_type IN ('weekly', 'monthly', 'quarterly', 'annual')),
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    limit_amount DECIMAL(15,2) NOT NULL CHECK (limit_amount > 0),
    rollover_enabled BOOLEAN DEFAULT false,
    rollover_amount DECIMAL(15,2) DEFAULT 0,
    alert_threshold DECIMAL(5,2) DEFAULT 0.80 CHECK (alert_threshold BETWEEN 0 AND 1),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_budgets_period ON budgets(period_start, period_end);
CREATE INDEX idx_budgets_category ON budgets(category_id);
CREATE INDEX idx_budgets_active ON budgets(is_active) WHERE is_active = true;

COMMENT ON TABLE budgets IS 'Spending limits by category and time period';
COMMENT ON COLUMN budgets.alert_threshold IS 'Alert when spending reaches this percentage (0.0-1.0)';
COMMENT ON COLUMN budgets.rollover_amount IS 'Amount rolled over from previous period';

-- ============================================================================
-- RECURRING ITEMS
-- ============================================================================
CREATE TABLE recurring_items (
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    description TEXT,
    frequency VARCHAR(20) NOT NULL CHECK (frequency IN ('daily', 'weekly', 'biweekly', 'monthly', 'quarterly', 'annual')),
    frequency_interval INTEGER DEFAULT 1 CHECK (frequency_interval > 0),
    day_of_month INTEGER CHECK (day_of_month BETWEEN 1 AND 31),
    day_of_week INTEGER CHECK (day_of_week BETWEEN 0 AND 6),  -- 0=Sunday
    start_date DATE NOT NULL,
    end_date DATE,
    next_date DATE NOT NULL,
    last_generated_date DATE,
    auto_generate BOOLEAN DEFAULT false,
    reminder_days_before INTEGER DEFAULT 3 CHECK (reminder_days_before >= 0),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_recurring_next_date ON recurring_items(next_date) WHERE is_active = true;
CREATE INDEX idx_recurring_account ON recurring_items(account_id);
CREATE INDEX idx_recurring_active ON recurring_items(is_active);

COMMENT ON TABLE recurring_items IS 'Templates for recurring income and expenses';
COMMENT ON COLUMN recurring_items.frequency_interval IS 'Repeat every N periods (e.g., every 2 weeks)';
COMMENT ON COLUMN recurring_items.auto_generate IS 'Automatically create transactions without confirmation';
COMMENT ON COLUMN recurring_items.reminder_days_before IS 'Create reminder N days before due date';

-- ============================================================================
-- REMINDERS
-- ============================================================================
CREATE TABLE reminders (
    id SERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL CHECK (type IN ('transaction', 'budget', 'bill', 'low_balance', 'custom')),
    related_id INTEGER,  -- ID of related entity (recurring_item, budget, etc.)
    title VARCHAR(255) NOT NULL,
    message TEXT,
    remind_date DATE NOT NULL,
    remind_time TIME,
    priority VARCHAR(20) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    is_dismissed BOOLEAN DEFAULT false,
    dismissed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_reminders_date ON reminders(remind_date, remind_time) WHERE is_dismissed = false;
CREATE INDEX idx_reminders_type ON reminders(type);
CREATE INDEX idx_reminders_dismissed ON reminders(is_dismissed);

COMMENT ON TABLE reminders IS 'Alerts and notifications for bills, budgets, and custom events';
COMMENT ON COLUMN reminders.related_id IS 'ID of related entity based on type';

-- ============================================================================
-- CASH FLOW PROJECTIONS
-- ============================================================================
CREATE TABLE cash_flow_projections (
    id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts(id) ON DELETE CASCADE,  -- NULL for all accounts
    projection_date DATE NOT NULL,
    projected_balance DECIMAL(15,2) NOT NULL,
    projected_income DECIMAL(15,2) DEFAULT 0,
    projected_expenses DECIMAL(15,2) DEFAULT 0,
    confidence_level DECIMAL(5,2) CHECK (confidence_level BETWEEN 0 AND 1),
    projection_type VARCHAR(20) DEFAULT 'moderate' CHECK (projection_type IN ('conservative', 'moderate', 'optimistic')),
    generated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_projections_account_date ON cash_flow_projections(account_id, projection_date);
CREATE INDEX idx_projections_date ON cash_flow_projections(projection_date);
CREATE INDEX idx_projections_type ON cash_flow_projections(projection_type);

COMMENT ON TABLE cash_flow_projections IS 'Future cash flow estimates and balance projections';
COMMENT ON COLUMN cash_flow_projections.confidence_level IS 'Confidence score 0.0-1.0 based on data quality';

-- ============================================================================
-- IMPORT HISTORY
-- ============================================================================
CREATE TABLE import_history (
    id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts(id) ON DELETE SET NULL,
    filename VARCHAR(255),
    file_hash VARCHAR(64) UNIQUE,  -- SHA256 to prevent duplicate imports
    format VARCHAR(50),  -- chase, generic, mint, ynab, etc.
    imported_at TIMESTAMP DEFAULT NOW(),
    records_total INTEGER,
    records_imported INTEGER,
    records_skipped INTEGER,
    records_failed INTEGER,
    error_log TEXT,
    import_metadata JSONB  -- Store format-specific details
);

CREATE INDEX idx_import_account ON import_history(account_id);
CREATE INDEX idx_import_date ON import_history(imported_at DESC);
CREATE INDEX idx_import_hash ON import_history(file_hash);

COMMENT ON TABLE import_history IS 'Track CSV/file imports to prevent duplicates';
COMMENT ON COLUMN import_history.file_hash IS 'SHA256 hash of file to detect duplicate imports';

-- ============================================================================
-- MATERIALIZED VIEWS (for performance)
-- ============================================================================

-- Monthly spending by category
CREATE MATERIALIZED VIEW monthly_spending_by_category AS
SELECT
    DATE_TRUNC('month', t.date)::DATE as month,
    t.category_id,
    c.name as category_name,
    c.type as category_type,
    SUM(t.amount) as total_amount,
    AVG(t.amount) as avg_amount,
    COUNT(*) as transaction_count
FROM transactions t
JOIN categories c ON t.category_id = c.id
WHERE t.type = 'expense'
GROUP BY DATE_TRUNC('month', t.date)::DATE, t.category_id, c.name, c.type;

CREATE UNIQUE INDEX idx_monthly_spending_month_category ON monthly_spending_by_category(month, category_id);
CREATE INDEX idx_monthly_spending_month ON monthly_spending_by_category(month DESC);

-- Account balance summary
CREATE MATERIALIZED VIEW account_balance_summary AS
SELECT
    a.id as account_id,
    a.name as account_name,
    a.type as account_type,
    a.current_balance,
    COUNT(t.id) as transaction_count,
    MAX(t.date) as last_transaction_date,
    SUM(CASE WHEN t.date >= CURRENT_DATE - INTERVAL '30 days' THEN ABS(t.amount) ELSE 0 END) as last_30_days_activity
FROM accounts a
LEFT JOIN transactions t ON a.id = t.account_id
WHERE a.is_active = true
GROUP BY a.id, a.name, a.type, a.current_balance;

CREATE UNIQUE INDEX idx_account_summary_id ON account_balance_summary(account_id);

COMMENT ON MATERIALIZED VIEW monthly_spending_by_category IS 'Pre-aggregated monthly spending for reports';
COMMENT ON MATERIALIZED VIEW account_balance_summary IS 'Account overview with transaction stats';

-- ============================================================================
-- FUNCTIONS AND TRIGGERS
-- ============================================================================

-- Function to update account balance after transaction
CREATE OR REPLACE FUNCTION update_account_balance()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- Add transaction amount to account balance
        UPDATE accounts
        SET current_balance = current_balance + NEW.amount,
            updated_at = NOW()
        WHERE id = NEW.account_id;

        -- If transfer, update destination account
        IF NEW.type = 'transfer' AND NEW.transfer_account_id IS NOT NULL THEN
            UPDATE accounts
            SET current_balance = current_balance - NEW.amount,
                updated_at = NOW()
            WHERE id = NEW.transfer_account_id;
        END IF;

    ELSIF TG_OP = 'UPDATE' THEN
        -- Reverse old amount
        UPDATE accounts
        SET current_balance = current_balance - OLD.amount,
            updated_at = NOW()
        WHERE id = OLD.account_id;

        -- Apply new amount
        UPDATE accounts
        SET current_balance = current_balance + NEW.amount,
            updated_at = NOW()
        WHERE id = NEW.account_id;

        -- Handle transfer account changes
        IF OLD.type = 'transfer' AND OLD.transfer_account_id IS NOT NULL THEN
            UPDATE accounts
            SET current_balance = current_balance + OLD.amount,
                updated_at = NOW()
            WHERE id = OLD.transfer_account_id;
        END IF;

        IF NEW.type = 'transfer' AND NEW.transfer_account_id IS NOT NULL THEN
            UPDATE accounts
            SET current_balance = current_balance - NEW.amount,
                updated_at = NOW()
            WHERE id = NEW.transfer_account_id;
        END IF;

    ELSIF TG_OP = 'DELETE' THEN
        -- Reverse transaction amount
        UPDATE accounts
        SET current_balance = current_balance - OLD.amount,
            updated_at = NOW()
        WHERE id = OLD.account_id;

        -- Reverse transfer if applicable
        IF OLD.type = 'transfer' AND OLD.transfer_account_id IS NOT NULL THEN
            UPDATE accounts
            SET current_balance = current_balance + OLD.amount,
                updated_at = NOW()
            WHERE id = OLD.transfer_account_id;
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to automatically update account balance
CREATE TRIGGER trg_update_account_balance
    AFTER INSERT OR UPDATE OR DELETE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION update_account_balance();

-- Function to update timestamp on row modification
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply updated_at trigger to relevant tables
CREATE TRIGGER trg_accounts_updated_at
    BEFORE UPDATE ON accounts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_transactions_updated_at
    BEFORE UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_budgets_updated_at
    BEFORE UPDATE ON budgets
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_recurring_items_updated_at
    BEFORE UPDATE ON recurring_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- SEED DATA (Default Categories)
-- ============================================================================

-- Income categories
INSERT INTO categories (name, type, icon, is_system) VALUES
    ('Salary', 'income', '💰', true),
    ('Freelance', 'income', '💼', true),
    ('Investment Income', 'income', '📈', true),
    ('Gifts', 'income', '🎁', true),
    ('Refunds', 'income', '↩️', true),
    ('Other Income', 'income', '💵', true);

-- Expense categories (top-level)
INSERT INTO categories (name, type, icon, is_system) VALUES
    ('Housing', 'expense', '🏠', true),
    ('Transportation', 'expense', '🚗', true),
    ('Food & Dining', 'expense', '🍽️', true),
    ('Utilities', 'expense', '⚡', true),
    ('Healthcare', 'expense', '🏥', true),
    ('Entertainment', 'expense', '🎬', true),
    ('Shopping', 'expense', '🛍️', true),
    ('Personal Care', 'expense', '💇', true),
    ('Education', 'expense', '📚', true),
    ('Subscriptions', 'expense', '📺', true),
    ('Insurance', 'expense', '🛡️', true),
    ('Savings & Investments', 'expense', '🏦', true),
    ('Taxes', 'expense', '🧾', true),
    ('Gifts & Donations', 'expense', '🎁', true),
    ('Other Expenses', 'expense', '📦', true);

-- Subcategories (examples)
INSERT INTO categories (name, parent_id, type, icon, is_system)
SELECT 'Rent/Mortgage', id, 'expense', '🏡', true FROM categories WHERE name = 'Housing' AND type = 'expense';

INSERT INTO categories (name, parent_id, type, icon, is_system)
SELECT 'Groceries', id, 'expense', '🛒', true FROM categories WHERE name = 'Food & Dining' AND type = 'expense';

INSERT INTO categories (name, parent_id, type, icon, is_system)
SELECT 'Restaurants', id, 'expense', '🍔', true FROM categories WHERE name = 'Food & Dining' AND type = 'expense';

INSERT INTO categories (name, parent_id, type, icon, is_system)
SELECT 'Gas/Fuel', id, 'expense', '⛽', true FROM categories WHERE name = 'Transportation' AND type = 'expense';

INSERT INTO categories (name, parent_id, type, icon, is_system)
SELECT 'Electric', id, 'expense', '💡', true FROM categories WHERE name = 'Utilities' AND type = 'expense';

INSERT INTO categories (name, parent_id, type, icon, is_system)
SELECT 'Internet', id, 'expense', '🌐', true FROM categories WHERE name = 'Utilities' AND type = 'expense';

-- Transfer category
INSERT INTO categories (name, type, icon, is_system) VALUES
    ('Transfer', 'transfer', '🔄', true);

-- ============================================================================
-- HELPER QUERIES (for reference)
-- ============================================================================

-- View all transactions with category names
-- SELECT t.id, t.date, t.amount, c.name as category, t.payee, t.description
-- FROM transactions t
-- LEFT JOIN categories c ON t.category_id = c.id
-- ORDER BY t.date DESC;

-- View account balances
-- SELECT a.name, a.type, a.current_balance
-- FROM accounts a
-- WHERE a.is_active = true
-- ORDER BY a.current_balance DESC;

-- View budget status for current month
-- SELECT b.name, b.limit_amount,
--        COALESCE(SUM(t.amount), 0) as spent,
--        b.limit_amount - COALESCE(SUM(t.amount), 0) as remaining
-- FROM budgets b
-- LEFT JOIN transactions t ON t.category_id = b.category_id
--     AND t.date BETWEEN b.period_start AND b.period_end
-- WHERE b.is_active = true
--   AND CURRENT_DATE BETWEEN b.period_start AND b.period_end
-- GROUP BY b.id, b.name, b.limit_amount;

-- Refresh materialized views (run periodically)
-- REFRESH MATERIALIZED VIEW CONCURRENTLY monthly_spending_by_category;
-- REFRESH MATERIALIZED VIEW CONCURRENTLY account_balance_summary;

-- ============================================================================
-- SCHEMA VERSION
-- ============================================================================
CREATE TABLE schema_version (
    version VARCHAR(20) PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO schema_version (version) VALUES ('1.0.0');

COMMENT ON TABLE schema_version IS 'Track database schema version for migrations';

-- End of schema
