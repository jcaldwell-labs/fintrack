package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Account represents a financial account
type Account struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	Name               string    `gorm:"not null;uniqueIndex:idx_accounts_name_active" json:"name"`
	Type               string    `gorm:"not null;index" json:"type"` // checking, savings, credit, cash, investment, loan
	Currency           string    `gorm:"default:USD" json:"currency"`
	InitialBalance     float64   `gorm:"default:0" json:"initial_balance"`
	CurrentBalance     float64   `gorm:"default:0" json:"current_balance"`
	Institution        string    `json:"institution,omitempty"`
	AccountNumberLast4 string    `gorm:"column:account_number_last4" json:"account_number_last4,omitempty"`
	IsActive           bool      `gorm:"default:true;index" json:"is_active"`
	Notes              string    `json:"notes,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// Category represents a transaction category
type Category struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"not null;uniqueIndex:idx_categories_name_type" json:"name"`
	ParentID  *uint      `gorm:"index" json:"parent_id,omitempty"`
	Parent    *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Type      string     `gorm:"not null;uniqueIndex:idx_categories_name_type" json:"type"` // income, expense, transfer
	Color     string     `json:"color,omitempty"`
	Icon      string     `json:"icon,omitempty"`
	IsSystem  bool       `gorm:"default:false" json:"is_system"`
	CreatedAt time.Time  `json:"created_at"`
}

// StringArray is a custom type for PostgreSQL array
type StringArray []string

// Scan implements the sql.Scanner interface
func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = []string{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, a)
	case string:
		return json.Unmarshal([]byte(v), a)
	default:
		return errors.New("unsupported type for StringArray")
	}
}

// Value implements the driver.Valuer interface
func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return "{}", nil
	}
	return json.Marshal(a)
}

// Transaction represents a financial transaction
type Transaction struct {
	ID                uint        `gorm:"primaryKey" json:"id"`
	AccountID         uint        `gorm:"not null;index" json:"account_id"`
	Account           *Account    `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	Date              time.Time   `gorm:"not null;index:idx_transactions_date,sort:desc" json:"date"`
	Amount            float64     `gorm:"not null" json:"amount"` // Positive for income, negative for expenses
	CategoryID        *uint       `gorm:"index" json:"category_id,omitempty"`
	Category          *Category   `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Payee             string      `gorm:"index" json:"payee,omitempty"`
	Description       string      `json:"description,omitempty"`
	Type              string      `gorm:"not null;index" json:"type"` // income, expense, transfer
	TransferAccountID *uint       `json:"transfer_account_id,omitempty"`
	TransferAccount   *Account    `gorm:"foreignKey:TransferAccountID" json:"transfer_account,omitempty"`
	RecurringID       *uint       `gorm:"index" json:"recurring_id,omitempty"`
	Tags              StringArray `gorm:"type:text[]" json:"tags,omitempty"`
	IsReconciled      bool        `gorm:"default:false;index" json:"is_reconciled"`
	ReconciledAt      *time.Time  `json:"reconciled_at,omitempty"`
	ImportID          *uint       `json:"import_id,omitempty"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

// Budget represents a spending limit
type Budget struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Name            string    `gorm:"not null" json:"name"`
	CategoryID      *uint     `gorm:"index" json:"category_id,omitempty"`
	Category        *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	PeriodType      string    `gorm:"not null" json:"period_type"` // weekly, monthly, quarterly, annual
	PeriodStart     time.Time `gorm:"not null;index:idx_budgets_period" json:"period_start"`
	PeriodEnd       time.Time `gorm:"not null;index:idx_budgets_period" json:"period_end"`
	LimitAmount     float64   `gorm:"not null" json:"limit_amount"`
	RolloverEnabled bool      `gorm:"default:false" json:"rollover_enabled"`
	RolloverAmount  float64   `gorm:"default:0" json:"rollover_amount"`
	AlertThreshold  float64   `gorm:"default:0.80" json:"alert_threshold"`
	IsActive        bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// RecurringItem represents a recurring transaction template
type RecurringItem struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	AccountID           uint      `gorm:"not null;index" json:"account_id"`
	Account             *Account  `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	Name                string    `gorm:"not null" json:"name"`
	Amount              float64   `gorm:"not null" json:"amount"`
	CategoryID          *uint     `json:"category_id,omitempty"`
	Category            *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Description         string    `json:"description,omitempty"`
	Frequency           string    `gorm:"not null" json:"frequency"` // daily, weekly, biweekly, monthly, quarterly, annual
	FrequencyInterval   int       `gorm:"default:1" json:"frequency_interval"`
	DayOfMonth          *int      `json:"day_of_month,omitempty"`
	DayOfWeek           *int      `json:"day_of_week,omitempty"`
	StartDate           time.Time `gorm:"not null" json:"start_date"`
	EndDate             *time.Time `json:"end_date,omitempty"`
	NextDate            time.Time `gorm:"not null;index:idx_recurring_next_date" json:"next_date"`
	LastGeneratedDate   *time.Time `json:"last_generated_date,omitempty"`
	AutoGenerate        bool      `gorm:"default:false" json:"auto_generate"`
	ReminderDaysBefore  int       `gorm:"default:3" json:"reminder_days_before"`
	IsActive            bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// Reminder represents an alert or notification
type Reminder struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Type        string     `gorm:"not null;index" json:"type"` // transaction, budget, bill, low_balance, custom
	RelatedID   *uint      `json:"related_id,omitempty"`
	Title       string     `gorm:"not null" json:"title"`
	Message     string     `json:"message,omitempty"`
	RemindDate  time.Time  `gorm:"not null;index:idx_reminders_date" json:"remind_date"`
	RemindTime  *time.Time `gorm:"type:time;index:idx_reminders_date" json:"remind_time,omitempty"`
	Priority    string     `gorm:"default:normal" json:"priority"` // low, normal, high, urgent
	IsDismissed bool       `gorm:"default:false;index" json:"is_dismissed"`
	DismissedAt *time.Time `json:"dismissed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// CashFlowProjection represents a future cash flow estimate
type CashFlowProjection struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	AccountID        *uint      `gorm:"index:idx_projections_account_date" json:"account_id,omitempty"`
	Account          *Account   `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	ProjectionDate   time.Time  `gorm:"not null;index:idx_projections_account_date,idx_projections_date" json:"projection_date"`
	ProjectedBalance float64    `gorm:"not null" json:"projected_balance"`
	ProjectedIncome  float64    `gorm:"default:0" json:"projected_income"`
	ProjectedExpenses float64   `gorm:"default:0" json:"projected_expenses"`
	ConfidenceLevel  float64    `json:"confidence_level,omitempty"` // 0.0-1.0
	ProjectionType   string     `gorm:"default:moderate;index" json:"projection_type"` // conservative, moderate, optimistic
	GeneratedAt      time.Time  `json:"generated_at"`
}

// ImportHistory tracks file imports
type ImportHistory struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	AccountID       *uint     `json:"account_id,omitempty"`
	Account         *Account  `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	Filename        string    `json:"filename,omitempty"`
	FileHash        string    `gorm:"uniqueIndex" json:"file_hash,omitempty"`
	Format          string    `json:"format,omitempty"`
	ImportedAt      time.Time `gorm:"index:idx_import_date,sort:desc" json:"imported_at"`
	RecordsTotal    int       `json:"records_total"`
	RecordsImported int       `json:"records_imported"`
	RecordsSkipped  int       `json:"records_skipped"`
	RecordsFailed   int       `json:"records_failed"`
	ErrorLog        string    `json:"error_log,omitempty"`
	ImportMetadata  string    `gorm:"type:jsonb" json:"import_metadata,omitempty"`
}

// AccountType constants
const (
	AccountTypeChecking   = "checking"
	AccountTypeSavings    = "savings"
	AccountTypeCredit     = "credit"
	AccountTypeCash       = "cash"
	AccountTypeInvestment = "investment"
	AccountTypeLoan       = "loan"
)

// CategoryType constants
const (
	CategoryTypeIncome    = "income"
	CategoryTypeExpense   = "expense"
	CategoryTypeTransfer  = "transfer"
)

// TransactionType constants
const (
	TransactionTypeIncome    = "income"
	TransactionTypeExpense   = "expense"
	TransactionTypeTransfer  = "transfer"
)

// Frequency constants
const (
	FrequencyDaily     = "daily"
	FrequencyWeekly    = "weekly"
	FrequencyBiweekly  = "biweekly"
	FrequencyMonthly   = "monthly"
	FrequencyQuarterly = "quarterly"
	FrequencyAnnual    = "annual"
)
