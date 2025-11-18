package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fintrack/fintrack/internal/db"
	"github.com/fintrack/fintrack/internal/db/repositories"
	"github.com/fintrack/fintrack/internal/models"
	"github.com/fintrack/fintrack/internal/output"
	"github.com/spf13/cobra"
)

// NewTransactionCmd creates the transaction command
func NewTransactionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "transaction",
		Aliases: []string{"tx", "t"},
		Short:   "Manage transactions",
		Long: `Manage financial transactions including income, expenses, and transfers.
Transactions are automatically linked to accounts and categories.`,
	}

	cmd.AddCommand(newTransactionAddCmd())
	cmd.AddCommand(newTransactionListCmd())
	cmd.AddCommand(newTransactionShowCmd())
	cmd.AddCommand(newTransactionUpdateCmd())
	cmd.AddCommand(newTransactionDeleteCmd())

	return cmd
}

func newTransactionAddCmd() *cobra.Command {
	var accountName string
	var categoryName string
	var dateStr string
	var payee string
	var description string
	var transactionType string
	var transferAccountName string
	var tags []string

	cmd := &cobra.Command{
		Use:   "add <amount>",
		Short: "Add a new transaction",
		Long: `Add a new transaction to an account.

Amount should be a number (use negative for expenses or positive for income).
Date format: YYYY-MM-DD (defaults to today)

Examples:
  # Add an expense
  fintrack tx add -50.00 --account "Checking" --category "Groceries" --payee "Walmart"

  # Add income
  fintrack tx add 2500.00 --account "Checking" --category "Salary" --date 2024-01-15

  # Add with tags
  fintrack tx add -30.00 --account "Checking" --category "Food & Dining" --tags "business,reimbursable"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse amount
			amount, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}

			// Validate required fields
			if accountName == "" {
				return fmt.Errorf("--account is required")
			}

			// Parse date
			var transactionDate time.Time
			if dateStr != "" {
				transactionDate, err = time.Parse("2006-01-02", dateStr)
				if err != nil {
					return fmt.Errorf("invalid date format (use YYYY-MM-DD): %w", err)
				}
			} else {
				transactionDate = time.Now()
			}

			// Determine transaction type if not specified
			if transactionType == "" {
				if amount >= 0 {
					transactionType = models.TransactionTypeIncome
				} else {
					transactionType = models.TransactionTypeExpense
				}
			}

			// Get account
			accountRepo := repositories.NewAccountRepository(db.Get())
			account, err := getAccountByIDOrName(accountRepo, accountName)
			if err != nil {
				return fmt.Errorf("account not found: %w", err)
			}

			transaction := &models.Transaction{
				AccountID:   account.ID,
				Date:        transactionDate,
				Amount:      amount,
				Payee:       payee,
				Description: description,
				Type:        transactionType,
				Tags:        tags,
			}

			// Get category if specified
			if categoryName != "" {
				categoryRepo := repositories.NewCategoryRepository(db.Get())
				// Try to find category by name (search all types)
				categories, err := categoryRepo.List("")
				if err != nil {
					return fmt.Errorf("failed to search categories: %w", err)
				}

				var foundCategory *models.Category
				for _, cat := range categories {
					if strings.EqualFold(cat.Name, categoryName) {
						foundCategory = cat
						break
					}
				}

				if foundCategory == nil {
					return fmt.Errorf("category '%s' not found", categoryName)
				}

				transaction.CategoryID = &foundCategory.ID
			}

			// Handle transfer
			if transferAccountName != "" {
				transferAccount, err := getAccountByIDOrName(accountRepo, transferAccountName)
				if err != nil {
					return fmt.Errorf("transfer account not found: %w", err)
				}
				transaction.TransferAccountID = &transferAccount.ID
				transaction.Type = models.TransactionTypeTransfer
			}

			// Create transaction
			txRepo := repositories.NewTransactionRepository(db.Get())
			if err := txRepo.Create(transaction); err != nil {
				return output.PrintError(cmd, fmt.Errorf("failed to create transaction: %w", err))
			}

			// Reload with associations
			transaction, err = txRepo.GetByID(transaction.ID)
			if err != nil {
				return fmt.Errorf("failed to reload transaction: %w", err)
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				return output.Print(cmd, transaction)
			}

			fmt.Printf("Transaction created successfully (ID: %d)\n", transaction.ID)
			fmt.Printf("Date: %s\n", transaction.Date.Format("2006-01-02"))
			fmt.Printf("Amount: %s\n", output.FormatCurrency(transaction.Amount, "USD"))
			fmt.Printf("Account: %s\n", transaction.Account.Name)
			if transaction.Category != nil {
				fmt.Printf("Category: %s\n", transaction.Category.Name)
			}
			if transaction.Payee != "" {
				fmt.Printf("Payee: %s\n", transaction.Payee)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&accountName, "account", "a", "", "Account name or ID (required)")
	cmd.Flags().StringVarP(&categoryName, "category", "c", "", "Category name or ID")
	cmd.Flags().StringVarP(&dateStr, "date", "d", "", "Transaction date (YYYY-MM-DD, defaults to today)")
	cmd.Flags().StringVarP(&payee, "payee", "p", "", "Payee name")
	cmd.Flags().StringVar(&description, "description", "", "Transaction description")
	cmd.Flags().StringVarP(&transactionType, "type", "t", "", "Transaction type (income, expense, transfer)")
	cmd.Flags().StringVar(&transferAccountName, "to", "", "Transfer to account (for transfers)")
	cmd.Flags().StringSliceVar(&tags, "tags", []string{}, "Tags (comma-separated)")

	_ = cmd.MarkFlagRequired("account")

	return cmd
}

func newTransactionListCmd() *cobra.Command {
	var accountName string
	var categoryName string
	var startDateStr string
	var endDateStr string
	var transactionType string
	var payee string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List transactions",
		Long: `List transactions with optional filters.

Examples:
  fintrack tx list
  fintrack tx list --account "Checking"
  fintrack tx list --category "Groceries" --start 2024-01-01 --end 2024-01-31
  fintrack tx list --type expense --limit 50`,
		RunE: func(cmd *cobra.Command, args []string) error {
			filter := &repositories.TransactionFilter{
				Type:  transactionType,
				Payee: payee,
				Limit: limit,
			}

			// Parse dates
			if startDateStr != "" {
				startDate, err := time.Parse("2006-01-02", startDateStr)
				if err != nil {
					return fmt.Errorf("invalid start date format (use YYYY-MM-DD): %w", err)
				}
				filter.StartDate = &startDate
			}

			if endDateStr != "" {
				endDate, err := time.Parse("2006-01-02", endDateStr)
				if err != nil {
					return fmt.Errorf("invalid end date format (use YYYY-MM-DD): %w", err)
				}
				filter.EndDate = &endDate
			}

			// Get account if specified
			if accountName != "" {
				accountRepo := repositories.NewAccountRepository(db.Get())
				account, err := getAccountByIDOrName(accountRepo, accountName)
				if err != nil {
					return fmt.Errorf("account not found: %w", err)
				}
				filter.AccountID = &account.ID
			}

			// Get category if specified
			if categoryName != "" {
				categoryRepo := repositories.NewCategoryRepository(db.Get())
				categories, err := categoryRepo.List("")
				if err != nil {
					return fmt.Errorf("failed to search categories: %w", err)
				}

				var foundCategory *models.Category
				for _, cat := range categories {
					if strings.EqualFold(cat.Name, categoryName) {
						foundCategory = cat
						break
					}
				}

				if foundCategory == nil {
					return fmt.Errorf("category '%s' not found", categoryName)
				}
				filter.CategoryID = &foundCategory.ID
			}

			txRepo := repositories.NewTransactionRepository(db.Get())
			transactions, err := txRepo.List(filter)
			if err != nil {
				return output.PrintError(cmd, fmt.Errorf("failed to list transactions: %w", err))
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				// Include summary
				income, expenses, err := txRepo.GetIncomeExpenseTotals(filter.AccountID, filter.StartDate, filter.EndDate)
				if err != nil {
					return output.PrintError(cmd, fmt.Errorf("failed to calculate totals: %w", err))
				}

				data := map[string]interface{}{
					"transactions": transactions,
					"summary": map[string]interface{}{
						"count":    len(transactions),
						"income":   income,
						"expenses": expenses,
						"net":      income - expenses,
					},
				}
				return output.Print(cmd, data)
			}

			if len(transactions) == 0 {
				fmt.Println("No transactions found.")
				return nil
			}

			// Create table
			table := output.NewTable("ID", "Date", "Account", "Category", "Payee", "Amount", "Type")
			for _, tx := range transactions {
				categoryName := ""
				if tx.Category != nil {
					categoryName = tx.Category.Name
				}

				accountName := ""
				if tx.Account != nil {
					accountName = tx.Account.Name
				}

				table.AddRow(
					strconv.FormatUint(uint64(tx.ID), 10),
					tx.Date.Format("2006-01-02"),
					accountName,
					categoryName,
					tx.Payee,
					output.FormatCurrency(tx.Amount, "USD"),
					tx.Type,
				)
			}

			table.Print()

			// Print summary
			income, expenses, err := txRepo.GetIncomeExpenseTotals(filter.AccountID, filter.StartDate, filter.EndDate)
			if err == nil {
				fmt.Printf("\nSummary:\n")
				fmt.Printf("  Total Transactions: %d\n", len(transactions))
				fmt.Printf("  Income: %s\n", output.FormatCurrency(income, "USD"))
				fmt.Printf("  Expenses: %s\n", output.FormatCurrency(expenses, "USD"))
				fmt.Printf("  Net: %s\n", output.FormatCurrency(income-expenses, "USD"))
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&accountName, "account", "a", "", "Filter by account name or ID")
	cmd.Flags().StringVarP(&categoryName, "category", "c", "", "Filter by category name or ID")
	cmd.Flags().StringVar(&startDateStr, "start", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&endDateStr, "end", "", "End date (YYYY-MM-DD)")
	cmd.Flags().StringVarP(&transactionType, "type", "t", "", "Filter by type (income, expense, transfer)")
	cmd.Flags().StringVarP(&payee, "payee", "p", "", "Filter by payee name")
	cmd.Flags().IntVarP(&limit, "limit", "l", 100, "Maximum number of transactions to display")

	return cmd
}

func newTransactionShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <id>",
		Short: "Show transaction details",
		Long: `Display detailed information about a specific transaction.

Examples:
  fintrack tx show 42
  fintrack transaction show 100`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid transaction ID: %w", err)
			}

			txRepo := repositories.NewTransactionRepository(db.Get())
			transaction, err := txRepo.GetByID(uint(id))
			if err != nil {
				return output.PrintError(cmd, err)
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				return output.Print(cmd, transaction)
			}

			fmt.Printf("Transaction ID: %d\n", transaction.ID)
			fmt.Printf("Date: %s\n", transaction.Date.Format("2006-01-02"))
			fmt.Printf("Amount: %s\n", output.FormatCurrency(transaction.Amount, "USD"))
			fmt.Printf("Type: %s\n", transaction.Type)

			if transaction.Account != nil {
				fmt.Printf("Account: %s (ID: %d)\n", transaction.Account.Name, transaction.Account.ID)
			}

			if transaction.Category != nil {
				fmt.Printf("Category: %s (ID: %d)\n", transaction.Category.Name, transaction.Category.ID)
			}

			if transaction.Payee != "" {
				fmt.Printf("Payee: %s\n", transaction.Payee)
			}

			if transaction.Description != "" {
				fmt.Printf("Description: %s\n", transaction.Description)
			}

			if transaction.TransferAccount != nil {
				fmt.Printf("Transfer To: %s (ID: %d)\n", transaction.TransferAccount.Name, transaction.TransferAccount.ID)
			}

			if len(transaction.Tags) > 0 {
				fmt.Printf("Tags: %s\n", strings.Join(transaction.Tags, ", "))
			}

			fmt.Printf("Reconciled: %v\n", transaction.IsReconciled)
			if transaction.ReconciledAt != nil {
				fmt.Printf("Reconciled At: %s\n", transaction.ReconciledAt.Format("2006-01-02 15:04:05"))
			}

			fmt.Printf("Created: %s\n", transaction.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("Updated: %s\n", transaction.UpdatedAt.Format("2006-01-02 15:04:05"))

			return nil
		},
	}

	return cmd
}

func newTransactionUpdateCmd() *cobra.Command {
	var amountValue float64
	var dateStr string
	var categoryName string
	var payee string
	var description string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a transaction",
		Long: `Update an existing transaction's properties.

Examples:
  fintrack tx update 42 --amount -75.50
  fintrack tx update 42 --category "Entertainment" --payee "Netflix"
  fintrack transaction update 100 --date 2024-01-15`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid transaction ID: %w", err)
			}

			txRepo := repositories.NewTransactionRepository(db.Get())
			transaction, err := txRepo.GetByID(uint(id))
			if err != nil {
				return output.PrintError(cmd, err)
			}

			// Apply updates
			updated := false

			if cmd.Flags().Changed("amount") {
				transaction.Amount = amountValue
				updated = true
			}

			if cmd.Flags().Changed("date") {
				date, err := time.Parse("2006-01-02", dateStr)
				if err != nil {
					return fmt.Errorf("invalid date format (use YYYY-MM-DD): %w", err)
				}
				transaction.Date = date
				updated = true
			}

			if cmd.Flags().Changed("category") {
				categoryRepo := repositories.NewCategoryRepository(db.Get())
				categories, err := categoryRepo.List("")
				if err != nil {
					return fmt.Errorf("failed to search categories: %w", err)
				}

				var foundCategory *models.Category
				for _, cat := range categories {
					if strings.EqualFold(cat.Name, categoryName) {
						foundCategory = cat
						break
					}
				}

				if foundCategory == nil {
					return fmt.Errorf("category '%s' not found", categoryName)
				}

				transaction.CategoryID = &foundCategory.ID
				updated = true
			}

			if cmd.Flags().Changed("payee") {
				transaction.Payee = payee
				updated = true
			}

			if cmd.Flags().Changed("description") {
				transaction.Description = description
				updated = true
			}

			if !updated {
				return fmt.Errorf("no updates specified")
			}

			if err := txRepo.Update(transaction); err != nil {
				return output.PrintError(cmd, fmt.Errorf("failed to update transaction: %w", err))
			}

			// Reload with associations
			transaction, err = txRepo.GetByID(transaction.ID)
			if err != nil {
				return fmt.Errorf("failed to reload transaction: %w", err)
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				return output.Print(cmd, transaction)
			}

			fmt.Println("Transaction updated successfully")
			fmt.Printf("Date: %s\n", transaction.Date.Format("2006-01-02"))
			fmt.Printf("Amount: %s\n", output.FormatCurrency(transaction.Amount, "USD"))

			return nil
		},
	}

	cmd.Flags().Float64Var(&amountValue, "amount", 0, "New amount")
	cmd.Flags().StringVar(&dateStr, "date", "", "New date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&categoryName, "category", "", "New category name")
	cmd.Flags().StringVar(&payee, "payee", "", "New payee")
	cmd.Flags().StringVar(&description, "description", "", "New description")

	return cmd
}

func newTransactionDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a transaction",
		Long: `Delete a transaction. This will also update the account balance.

Examples:
  fintrack tx delete 42
  fintrack transaction delete 100`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid transaction ID: %w", err)
			}

			txRepo := repositories.NewTransactionRepository(db.Get())

			if err := txRepo.Delete(uint(id)); err != nil {
				return output.PrintError(cmd, fmt.Errorf("failed to delete transaction: %w", err))
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				return output.Print(cmd, map[string]string{"message": "Transaction deleted successfully"})
			}

			fmt.Println("Transaction deleted successfully")
			return nil
		},
	}

	return cmd
}

// Helper function to get account by ID or name
func getAccountByIDOrName(repo *repositories.AccountRepository, idOrName string) (*models.Account, error) {
	// Try to parse as ID first
	if id, err := strconv.ParseUint(idOrName, 10, 32); err == nil {
		return repo.GetByID(uint(id))
	}

	// Otherwise, try to get by name
	return repo.GetByName(idOrName)
}
