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

// Transaction types
const (
	TxTypeIncome   = "income"
	TxTypeExpense  = "expense"
	TxTypeTransfer = "transfer"
)

// NewTransactionCmd creates the transaction command
func NewTransactionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "transaction",
		Aliases: []string{"tx", "t"},
		Short:   "Manage transactions",
		Long: `Manage financial transactions including income, expenses, and transfers.

Examples:
  fintrack tx list
  fintrack tx add --account 1 --amount -50.00 --payee "Grocery Store" --category 5
  fintrack tx show 1
  fintrack tx update 1 --payee "Updated Payee"
  fintrack tx delete 1`,
	}

	cmd.AddCommand(newTransactionListCmd())
	cmd.AddCommand(newTransactionAddCmd())
	cmd.AddCommand(newTransactionShowCmd())
	cmd.AddCommand(newTransactionUpdateCmd())
	cmd.AddCommand(newTransactionDeleteCmd())

	return cmd
}

func newTransactionListCmd() *cobra.Command {
	var (
		accountID  uint
		categoryID uint
		txType     string
		dateFrom   string
		dateTo     string
		payee      string
		limit      int
	)

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List transactions",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo := repositories.NewTransactionRepository(db.Get())

			filter := repositories.TransactionFilter{
				Limit: limit,
			}

			if accountID > 0 {
				filter.AccountID = &accountID
			}
			if categoryID > 0 {
				filter.CategoryID = &categoryID
			}
			if txType != "" {
				filter.Type = txType
			}
			if payee != "" {
				filter.Payee = payee
			}
			if dateFrom != "" {
				t, err := time.Parse("2006-01-02", dateFrom)
				if err != nil {
					return output.PrintError(cmd, fmt.Errorf("invalid from date format (use YYYY-MM-DD): %v", err))
				}
				filter.DateFrom = &t
			}
			if dateTo != "" {
				t, err := time.Parse("2006-01-02", dateTo)
				if err != nil {
					return output.PrintError(cmd, fmt.Errorf("invalid to date format (use YYYY-MM-DD): %v", err))
				}
				filter.DateTo = &t
			}

			transactions, err := repo.List(filter)
			if err != nil {
				return output.PrintError(cmd, err)
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				return output.Print(cmd, transactions)
			}

			// Table format
			table := output.NewTable("ID", "DATE", "AMOUNT", "TYPE", "PAYEE", "CATEGORY", "ACCOUNT")

			// Track totals for summary
			var incomeCents, expenseCents int64
			txCount := len(transactions)

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
					fmt.Sprintf("%d", tx.ID),
					tx.Date.Format("2006-01-02"),
					formatAmountCents(tx.AmountCents),
					tx.Type,
					tx.Payee,
					categoryName,
					accountName,
				)

				// Track income vs expenses
				if tx.AmountCents > 0 {
					incomeCents += tx.AmountCents
				} else {
					expenseCents += tx.AmountCents
				}
			}
			table.Print()

			// Print summary
			if txCount > 0 {
				netCents := incomeCents + expenseCents
				fmt.Printf("\nSummary: %d transactions | Income: %s | Expenses: %s | Net: %s\n",
					txCount,
					output.FormatCurrencyCents(incomeCents, "USD"),
					output.FormatCurrencyCents(-expenseCents, "USD"),
					output.FormatCurrencyCents(netCents, "USD"),
				)
			}

			return nil
		},
	}

	cmd.Flags().UintVar(&accountID, "account", 0, "Filter by account ID")
	cmd.Flags().UintVar(&categoryID, "category", 0, "Filter by category ID")
	cmd.Flags().StringVar(&txType, "type", "", "Filter by type (income, expense, transfer)")
	cmd.Flags().StringVar(&dateFrom, "from", "", "Filter from date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&dateTo, "to", "", "Filter to date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&payee, "payee", "", "Filter by payee (partial match)")
	cmd.Flags().IntVar(&limit, "limit", 50, "Maximum number of transactions to show")

	return cmd
}

func newTransactionAddCmd() *cobra.Command {
	var (
		accountID   uint
		amount      float64
		categoryID  uint
		payee       string
		description string
		txType      string
		date        string
		tags        string
	)

	cmd := &cobra.Command{
		Use:     "add",
		Aliases: []string{"create", "new"},
		Short:   "Add a new transaction",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate required fields
			if accountID == 0 {
				return output.PrintError(cmd, fmt.Errorf("--account is required"))
			}
			if amount == 0 {
				return output.PrintError(cmd, fmt.Errorf("--amount is required and cannot be zero"))
			}

			// Validate transaction type
			if txType == "" {
				// Auto-determine type based on amount
				if amount > 0 {
					txType = TxTypeIncome
				} else {
					txType = TxTypeExpense
				}
			}
			validTypes := map[string]bool{
				TxTypeIncome:   true,
				TxTypeExpense:  true,
				TxTypeTransfer: true,
			}
			if !validTypes[txType] {
				return output.PrintError(cmd, fmt.Errorf("invalid transaction type: %s (valid: income, expense, transfer)", txType))
			}

			// Parse date
			txDate := time.Now()
			if date != "" {
				var err error
				txDate, err = time.Parse("2006-01-02", date)
				if err != nil {
					return output.PrintError(cmd, fmt.Errorf("invalid date format (use YYYY-MM-DD): %v", err))
				}
			}

			// Parse tags
			var tagList []string
			if tags != "" {
				tagList = strings.Split(tags, ",")
				for i := range tagList {
					tagList[i] = strings.TrimSpace(tagList[i])
				}
			}

			// Convert dollars to cents for storage
			amountCents := models.DollarsToCents(amount)

			tx := &models.Transaction{
				AccountID:   accountID,
				AmountCents: amountCents,
				Payee:       payee,
				Description: description,
				Type:        txType,
				Date:        txDate,
				Tags:        tagList,
			}

			if categoryID > 0 {
				tx.CategoryID = &categoryID
			}

			repo := repositories.NewTransactionRepository(db.Get())
			if err := repo.Create(tx); err != nil {
				return output.PrintError(cmd, err)
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				return output.Print(cmd, tx)
			}

			fmt.Printf("✓ Created transaction #%d\n", tx.ID)
			fmt.Printf("Date: %s\n", tx.Date.Format("2006-01-02"))
			fmt.Printf("Amount: %s\n", formatAmountCents(tx.AmountCents))
			fmt.Printf("Type: %s\n", tx.Type)
			if tx.Payee != "" {
				fmt.Printf("Payee: %s\n", tx.Payee)
			}

			return nil
		},
	}

	cmd.Flags().UintVar(&accountID, "account", 0, "Account ID (required)")
	cmd.Flags().Float64VarP(&amount, "amount", "a", 0, "Transaction amount in dollars (negative for expense, positive for income)")
	cmd.Flags().UintVar(&categoryID, "category", 0, "Category ID")
	cmd.Flags().StringVarP(&payee, "payee", "p", "", "Payee name")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Description")
	cmd.Flags().StringVarP(&txType, "type", "t", "", "Transaction type (income, expense, transfer)")
	cmd.Flags().StringVar(&date, "date", "", "Transaction date (YYYY-MM-DD, default: today)")
	cmd.Flags().StringVar(&tags, "tags", "", "Comma-separated tags")

	mustMarkRequired(cmd, "account")
	mustMarkRequired(cmd, "amount")

	return cmd
}

func newTransactionShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show ID",
		Aliases: []string{"get"},
		Short:   "Show transaction details",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return output.PrintError(cmd, fmt.Errorf("invalid transaction ID: %s", args[0]))
			}

			repo := repositories.NewTransactionRepository(db.Get())
			tx, err := repo.GetByID(uint(id))
			if err != nil {
				return output.PrintError(cmd, err)
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				return output.Print(cmd, tx)
			}

			// Detailed output
			fmt.Printf("Transaction #%d\n", tx.ID)
			fmt.Printf("Date: %s\n", tx.Date.Format("2006-01-02"))
			fmt.Printf("Amount: %s\n", formatAmountCents(tx.AmountCents))
			fmt.Printf("Type: %s\n", tx.Type)
			if tx.Account != nil {
				fmt.Printf("Account: %s (#%d)\n", tx.Account.Name, tx.AccountID)
			}
			if tx.Payee != "" {
				fmt.Printf("Payee: %s\n", tx.Payee)
			}
			if tx.Category != nil {
				fmt.Printf("Category: %s\n", tx.Category.Name)
			}
			if tx.Description != "" {
				fmt.Printf("Description: %s\n", tx.Description)
			}
			if len(tx.Tags) > 0 {
				fmt.Printf("Tags: %s\n", strings.Join(tx.Tags, ", "))
			}
			fmt.Printf("Reconciled: %v\n", tx.IsReconciled)
			fmt.Printf("Created: %s\n", tx.CreatedAt.Format("2006-01-02 15:04:05"))

			return nil
		},
	}

	return cmd
}

func newTransactionUpdateCmd() *cobra.Command {
	var (
		amount      float64
		categoryID  uint
		payee       string
		description string
		date        string
		tags        string
		reconcile   bool
	)

	cmd := &cobra.Command{
		Use:   "update ID",
		Short: "Update a transaction",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return output.PrintError(cmd, fmt.Errorf("invalid transaction ID: %s", args[0]))
			}

			repo := repositories.NewTransactionRepository(db.Get())
			tx, err := repo.GetByID(uint(id))
			if err != nil {
				return output.PrintError(cmd, err)
			}

			// Update fields if provided
			if cmd.Flags().Changed("amount") {
				tx.AmountCents = models.DollarsToCents(amount)
			}
			if cmd.Flags().Changed("category") {
				tx.CategoryID = &categoryID
			}
			if cmd.Flags().Changed("payee") {
				tx.Payee = payee
			}
			if cmd.Flags().Changed("description") {
				tx.Description = description
			}
			if cmd.Flags().Changed("date") {
				t, err := time.Parse("2006-01-02", date)
				if err != nil {
					return output.PrintError(cmd, fmt.Errorf("invalid date format (use YYYY-MM-DD): %v", err))
				}
				tx.Date = t
			}
			if cmd.Flags().Changed("tags") {
				var tagList []string
				if tags != "" {
					tagList = strings.Split(tags, ",")
					for i := range tagList {
						tagList[i] = strings.TrimSpace(tagList[i])
					}
				}
				tx.Tags = tagList
			}
			if cmd.Flags().Changed("reconcile") {
				if reconcile {
					if err := repo.Reconcile(uint(id)); err != nil {
						return output.PrintError(cmd, err)
					}
				} else {
					if err := repo.Unreconcile(uint(id)); err != nil {
						return output.PrintError(cmd, err)
					}
				}
			}

			if err := repo.Update(tx); err != nil {
				return output.PrintError(cmd, err)
			}

			return output.PrintSuccess(cmd, fmt.Sprintf("Transaction #%d updated successfully", id))
		},
	}

	cmd.Flags().Float64VarP(&amount, "amount", "a", 0, "New amount in dollars")
	cmd.Flags().UintVar(&categoryID, "category", 0, "New category ID")
	cmd.Flags().StringVarP(&payee, "payee", "p", "", "New payee")
	cmd.Flags().StringVarP(&description, "description", "d", "", "New description")
	cmd.Flags().StringVar(&date, "date", "", "New date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&tags, "tags", "", "New tags (comma-separated)")
	cmd.Flags().BoolVar(&reconcile, "reconcile", false, "Mark as reconciled")

	return cmd
}

func newTransactionDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete ID",
		Aliases: []string{"rm", "remove"},
		Short:   "Delete a transaction",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return output.PrintError(cmd, fmt.Errorf("invalid transaction ID: %s", args[0]))
			}

			repo := repositories.NewTransactionRepository(db.Get())
			if err := repo.Delete(uint(id)); err != nil {
				return output.PrintError(cmd, err)
			}

			return output.PrintSuccess(cmd, fmt.Sprintf("Transaction #%d deleted successfully", id))
		},
	}

	return cmd
}

// Helper function to format amount with sign (for cents)
func formatAmountCents(cents int64) string {
	dollars := float64(cents) / 100
	if cents >= 0 {
		return fmt.Sprintf("+%.2f", dollars)
	}
	return fmt.Sprintf("%.2f", dollars)
}
