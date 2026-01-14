package services

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fintrack/fintrack/internal/db/repositories"
	"github.com/fintrack/fintrack/internal/models"
	"gorm.io/gorm"
)

type CSVColumnMapping struct {
	DateColumn        int
	AmountColumn      int
	DescriptionColumn int
	PayeeColumn       int
	CategoryColumn    int
	DateFormat        string
	HasHeader         bool
	AmountNegative    bool
}

func DefaultColumnMapping() CSVColumnMapping {
	return CSVColumnMapping{
		DateColumn:        0,
		AmountColumn:      1,
		DescriptionColumn: 2,
		PayeeColumn:       -1,
		CategoryColumn:    -1,
		DateFormat:        "2006-01-02",
		HasHeader:         true,
		AmountNegative:    true,
	}
}

type ImportResult struct {
	TotalRecords    int
	ImportedRecords int
	SkippedRecords  int
	FailedRecords   int
	Transactions    []*models.Transaction
	Errors          []ImportError
	FileHash        string
}

type ImportError struct {
	Line    int
	Message string
	Data    string
}

type CSVImporter struct {
	db          *gorm.DB
	txRepo      *repositories.TransactionRepository
	historyRepo *repositories.ImportHistoryRepository
	accountRepo *repositories.AccountRepository
}

func NewCSVImporter(db *gorm.DB) *CSVImporter {
	return &CSVImporter{
		db:          db,
		txRepo:      repositories.NewTransactionRepository(db),
		historyRepo: repositories.NewImportHistoryRepository(db),
		accountRepo: repositories.NewAccountRepository(db),
	}
}

type ImportOptions struct {
	AccountID      uint
	Mapping        CSVColumnMapping
	DryRun         bool
	SkipDuplicates bool
	BatchSize      int
}

func (i *CSVImporter) Import(filePath string, opts ImportOptions) (*ImportResult, error) {
	account, err := i.accountRepo.GetByID(opts.AccountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash, err := calculateFileHash(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate file hash: %w", err)
	}

	if !opts.DryRun {
		exists, err := i.historyRepo.FileHashExists(hash)
		if err != nil {
			return nil, fmt.Errorf("failed to check import history: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("file has already been imported (hash: %s)", hash[:16])
		}
	}

	file.Seek(0, 0)

	result, err := i.parseCSV(file, account, opts)
	if err != nil {
		return nil, err
	}
	result.FileHash = hash

	if opts.DryRun {
		return result, nil
	}

	err = i.db.Transaction(func(tx *gorm.DB) error {
		history := &models.ImportHistory{
			AccountID:       &opts.AccountID,
			Filename:        filePath,
			FileHash:        hash,
			Format:          "csv",
			ImportedAt:      time.Now(),
			RecordsTotal:    result.TotalRecords,
			RecordsImported: result.ImportedRecords,
			RecordsSkipped:  result.SkippedRecords,
			RecordsFailed:   result.FailedRecords,
		}

		historyRepo := repositories.NewImportHistoryRepository(tx)
		if err := historyRepo.Create(history); err != nil {
			return fmt.Errorf("failed to create import history: %w", err)
		}

		for _, txn := range result.Transactions {
			txn.ImportID = &history.ID
		}

		if len(result.Transactions) > 0 {
			txnRepo := repositories.NewTransactionRepository(tx)
			batchSize := opts.BatchSize
			if batchSize <= 0 {
				batchSize = 100
			}
			if err := txnRepo.CreateBatch(result.Transactions, batchSize); err != nil {
				return fmt.Errorf("failed to create transactions: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (i *CSVImporter) parseCSV(reader io.Reader, account *models.Account, opts ImportOptions) (*ImportResult, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true

	result := &ImportResult{
		Transactions: make([]*models.Transaction, 0),
		Errors:       make([]ImportError, 0),
	}

	lineNum := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Errors = append(result.Errors, ImportError{
				Line:    lineNum + 1,
				Message: fmt.Sprintf("CSV parse error: %v", err),
			})
			result.FailedRecords++
			lineNum++
			continue
		}

		lineNum++

		if lineNum == 1 && opts.Mapping.HasHeader {
			continue
		}

		result.TotalRecords++

		txn, err := i.parseRecord(record, account, opts, lineNum)
		if err != nil {
			result.Errors = append(result.Errors, ImportError{
				Line:    lineNum,
				Message: err.Error(),
				Data:    strings.Join(record, ","),
			})
			result.FailedRecords++
			continue
		}

		if opts.SkipDuplicates {
			dup, err := i.txRepo.FindDuplicate(account.ID, repositories.DuplicateCheck{
				Date:        txn.Date,
				Amount:      txn.Amount,
				Description: txn.Description,
			})
			if err != nil {
				result.Errors = append(result.Errors, ImportError{
					Line:    lineNum,
					Message: fmt.Sprintf("duplicate check failed: %v", err),
				})
				result.FailedRecords++
				continue
			}
			if dup != nil {
				result.SkippedRecords++
				continue
			}
		}

		result.Transactions = append(result.Transactions, txn)
		result.ImportedRecords++
	}

	return result, nil
}

func (i *CSVImporter) parseRecord(record []string, account *models.Account, opts ImportOptions, lineNum int) (*models.Transaction, error) {
	mapping := opts.Mapping

	maxCol := maxInt(mapping.DateColumn, mapping.AmountColumn, mapping.DescriptionColumn)
	if len(record) <= maxCol {
		return nil, fmt.Errorf("not enough columns (expected at least %d, got %d)", maxCol+1, len(record))
	}

	dateStr := strings.TrimSpace(record[mapping.DateColumn])
	date, err := parseDate(dateStr, mapping.DateFormat)
	if err != nil {
		return nil, fmt.Errorf("invalid date: %v", err)
	}

	amountStr := strings.TrimSpace(record[mapping.AmountColumn])
	amount, err := parseAmount(amountStr)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %v", err)
	}

	description := strings.TrimSpace(record[mapping.DescriptionColumn])
	if description == "" {
		description = "Imported transaction"
	}

	txType := models.TransactionTypeExpense
	if amount > 0 {
		txType = models.TransactionTypeIncome
	}

	if !mapping.AmountNegative {
		if amount > 0 {
			txType = models.TransactionTypeExpense
			amount = -amount
		} else {
			txType = models.TransactionTypeIncome
			amount = -amount
		}
	}

	payee := ""
	if mapping.PayeeColumn >= 0 && len(record) > mapping.PayeeColumn {
		payee = strings.TrimSpace(record[mapping.PayeeColumn])
	}

	txn := &models.Transaction{
		AccountID:   account.ID,
		Date:        date,
		Amount:      amount,
		Description: description,
		Payee:       payee,
		Type:        txType,
	}

	return txn, nil
}

func calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func parseDate(dateStr string, format string) (time.Time, error) {
	if t, err := time.Parse(format, dateStr); err == nil {
		return t, nil
	}
	formats := []string{"2006-01-02", "01/02/2006", "02/01/2006", "1/2/2006", "2/1/2006", "01-02-2006", "02-01-2006", "2006/01/02", "Jan 2, 2006", "January 2, 2006", "2 Jan 2006"}
	for _, f := range formats {
		if t, err := time.Parse(f, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date")
}

func parseAmount(amountStr string) (float64, error) {
	amountStr = strings.TrimSpace(amountStr)
	amountStr = strings.ReplaceAll(amountStr, "$", "")
	amountStr = strings.ReplaceAll(amountStr, ",", "")
	amountStr = strings.ReplaceAll(amountStr, " ", "")
	if strings.HasPrefix(amountStr, "(") && strings.HasSuffix(amountStr, ")") {
		amountStr = "-" + amountStr[1:len(amountStr)-1]
	}
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0, err
	}
	return math.Round(amount*100) / 100, nil
}

func maxInt(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}
	m := nums[0]
	for _, n := range nums[1:] {
		if n > m {
			m = n
		}
	}
	return m
}
