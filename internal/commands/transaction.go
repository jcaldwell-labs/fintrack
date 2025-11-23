package commands

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/fintrack/fintrack/internal/config"
	"github.com/fintrack/fintrack/internal/db"
	"github.com/fintrack/fintrack/internal/db/repositories"
	"github.com/fintrack/fintrack/internal/models"
	"github.com/fintrack/fintrack/internal/output"
	"github.com/spf13/cobra"
)

// NewTransactionCmd creates the transaction command
func NewTransactionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tx",
		Aliases: []string{"t", "transaction"},
		Short:   "Manage transactions",
		Long:    `Add, list, update, and delete financial transactions`,
	}

	cmd.AddCommand(newTxAddCmd())
	cmd.AddCommand(newTxListCmd())
	cmd.AddCommand(newTxShowCmd())
	cmd.AddCommand(newTxUpdateCmd())
	cmd.AddCommand(newTxDeleteCmd())
	cmd.AddCommand(newTxReconcileCmd())

	return cmd
}

func newTxAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add ACCOUNT_ID AMOUNT",
		Short: "Add a new transaction",
		Long:  `Add a new transaction to an account. Use negative amounts for expenses, positive for income.`,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse account ID
			accountID, err := parseAccountID(args[0])
			if err != nil {
				return output.PrintError(cmd, fmt.Errorf("invalid account: %w", err))
			}

			// Parse amount
			amount, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return output.PrintError(cmd, fmt.Errorf("invalid amount: %w", err))
			}

			// Get flags
			dateStr, _ := cmd.Flags().GetString("date")
			txType, _ := cmd.Flags().GetString("type")
			description, _ := cmd.Flags().GetString("description")
			payee, _ := cmd.Flags().GetString("payee")
			categoryID, _ := cmd.Flags().GetUint("category")

			// Parse date
			var txDate time.Time
			if dateStr != "" {
				cfg := config.Get()
				txDate, err = time.Parse(cfg.Defaults.DateFormat, dateStr)
				if err != nil {
					return output.PrintError(cmd, fmt.Errorf("invalid date format (expected %s): %w", cfg.Defaults.DateFormat, err))
				}
			} else {
				txDate = time.Now()
			}

			// Determine transaction type
			if txType == "" {
				if amount > 0 {
					txType = models.TransactionTypeIncome
				} else {
					txType = models.TransactionTypeExpense
				}
			}

			// Validate type
			if txType != models.TransactionTypeIncome &&
				txType != models.TransactionTypeExpense &&
				txType != models.TransactionTypeTransfer {
				return output.PrintError(cmd, errors.New("type must be income, expense, or transfer"))
			}

			// Create transaction
			tx := &models.Transaction{
				AccountID:   accountID,
				Date:        txDate,
				Amount:      amount,
				Type:        txType,
				Description: description,
				Payee:       payee,
			}

			if categoryID > 0 {
				tx.CategoryID = &categoryID
			}

			// Save to database
			repo := repositories.NewTransactionRepository(db.Get())
			if err := repo.Create(tx); err != nil {
				return output.PrintError(cmd, err)
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				return output.PrintSuccess(cmd, fmt.Sprintf("Transaction %d created", tx.ID))
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "✓ Transaction added successfully (ID: %d)\n", tx.ID)
			return nil
		},
	}

	cmd.Flags().String("date", "", "Transaction date (YYYY-MM-DD)")
	cmd.Flags().StringP("type", "t", "", "Transaction type (income, expense, transfer)")
	cmd.Flags().StringP("description", "d", "", "Transaction description")
	cmd.Flags().StringP("payee", "p", "", "Payee/merchant name")
	cmd.Flags().UintP("category", "c", 0, "Category ID")

	return cmd
}

func newTxListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [ACCOUNT_ID]",
		Aliases: []string{"ls"},
		Short:   "List transactions",
		Long:    `List all transactions or filter by account, date range, or type`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var accountID *uint
			if len(args) > 0 {
				id, err := parseAccountID(args[0])
				if err != nil {
					return output.PrintError(cmd, fmt.Errorf("invalid account: %w", err))
				}
				accountID = &id
			}

			// Get flags
			startDateStr, _ := cmd.Flags().GetString("start-date")
			endDateStr, _ := cmd.Flags().GetString("end-date")
			txType, _ := cmd.Flags().GetString("type")
			limit, _ := cmd.Flags().GetInt("limit")

			var startDate, endDate *time.Time
			cfg := config.Get()

			if startDateStr != "" {
				d, err := time.Parse(cfg.Defaults.DateFormat, startDateStr)
				if err != nil {
					return output.PrintError(cmd, fmt.Errorf("invalid start date: %w", err))
				}
				startDate = &d
			}

			if endDateStr != "" {
				d, err := time.Parse(cfg.Defaults.DateFormat, endDateStr)
				if err != nil {
					return output.PrintError(cmd, fmt.Errorf("invalid end date: %w", err))
				}
				endDate = &d
			}

			var typeFilter *string
			if txType != "" {
				typeFilter = &txType
			}

			// Fetch transactions
			repo := repositories.NewTransactionRepository(db.Get())
			transactions, err := repo.List(accountID, startDate, endDate, typeFilter, limit)
			if err != nil {
				return output.PrintError(cmd, err)
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				return output.Print(cmd, transactions)
			}

			// Table output
			table := output.NewTable("ID", "Date", "Account", "Amount", "Type", "Payee", "Description")
			for _, tx := range transactions {
				accountName := ""
				if tx.Account != nil {
					accountName = tx.Account.Name
				}

				table.AddRow(
					fmt.Sprintf("%d", tx.ID),
					tx.Date.Format(cfg.Defaults.DateFormat),
					accountName,
					output.FormatCurrency(tx.Amount, cfg.Defaults.Currency),
					tx.Type,
					tx.Payee,
					tx.Description,
				)
			}
			table.Print()
			return nil
		},
	}

	cmd.Flags().String("start-date", "", "Start date filter (YYYY-MM-DD)")
	cmd.Flags().String("end-date", "", "End date filter (YYYY-MM-DD)")
	cmd.Flags().StringP("type", "t", "", "Type filter (income, expense, transfer)")
	cmd.Flags().IntP("limit", "l", 50, "Limit number of results")

	return cmd
}

func newTxShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show TRANSACTION_ID",
		Short: "Show transaction details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return output.PrintError(cmd, fmt.Errorf("invalid transaction ID: %w", err))
			}

			repo := repositories.NewTransactionRepository(db.Get())
			tx, err := repo.GetByID(uint(id))
			if err != nil {
				return output.PrintError(cmd, err)
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				return output.Print(cmd, tx)
			}

			cfg := config.Get()
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Transaction ID:  %d\n", tx.ID)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Date:            %s\n", tx.Date.Format(cfg.Defaults.DateFormat))
			if tx.Account != nil {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Account:         %s\n", tx.Account.Name)
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Amount:          %s\n", output.FormatCurrency(tx.Amount, cfg.Defaults.Currency))
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Type:            %s\n", tx.Type)
			if tx.Payee != "" {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Payee:           %s\n", tx.Payee)
			}
			if tx.Description != "" {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Description:     %s\n", tx.Description)
			}
			if tx.Category != nil {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Category:        %s\n", tx.Category.Name)
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Reconciled:      %v\n", tx.IsReconciled)
			if tx.ReconciledAt != nil {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Reconciled At:   %s\n", tx.ReconciledAt.Format(cfg.Defaults.DateFormat))
			}

			return nil
		},
	}
}

func newTxUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update TRANSACTION_ID",
		Short: "Update a transaction",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return output.PrintError(cmd, fmt.Errorf("invalid transaction ID: %w", err))
			}

			repo := repositories.NewTransactionRepository(db.Get())
			tx, err := repo.GetByID(uint(id))
			if err != nil {
				return output.PrintError(cmd, err)
			}

			// Update fields if flags are provided
			if cmd.Flags().Changed("amount") {
				amount, _ := cmd.Flags().GetFloat64("amount")
				tx.Amount = amount
			}

			if cmd.Flags().Changed("date") {
				dateStr, _ := cmd.Flags().GetString("date")
				cfg := config.Get()
				date, err := time.Parse(cfg.Defaults.DateFormat, dateStr)
				if err != nil {
					return output.PrintError(cmd, fmt.Errorf("invalid date format: %w", err))
				}
				tx.Date = date
			}

			if cmd.Flags().Changed("type") {
				txType, _ := cmd.Flags().GetString("type")
				tx.Type = txType
			}

			if cmd.Flags().Changed("description") {
				description, _ := cmd.Flags().GetString("description")
				tx.Description = description
			}

			if cmd.Flags().Changed("payee") {
				payee, _ := cmd.Flags().GetString("payee")
				tx.Payee = payee
			}

			if cmd.Flags().Changed("category") {
				categoryID, _ := cmd.Flags().GetUint("category")
				if categoryID > 0 {
					tx.CategoryID = &categoryID
				} else {
					tx.CategoryID = nil
				}
			}

			if err := repo.Update(tx); err != nil {
				return output.PrintError(cmd, err)
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				return output.PrintSuccess(cmd, "Transaction updated")
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "✓ Transaction updated successfully")
			return nil
		},
	}

	cmd.Flags().Float64P("amount", "a", 0, "New amount")
	cmd.Flags().String("date", "", "New date (YYYY-MM-DD)")
	cmd.Flags().StringP("type", "t", "", "New type (income, expense, transfer)")
	cmd.Flags().StringP("description", "d", "", "New description")
	cmd.Flags().StringP("payee", "p", "", "New payee")
	cmd.Flags().UintP("category", "c", 0, "New category ID (0 to remove)")

	return cmd
}

func newTxDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete TRANSACTION_ID",
		Aliases: []string{"del", "rm"},
		Short:   "Delete a transaction",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return output.PrintError(cmd, fmt.Errorf("invalid transaction ID: %w", err))
			}

			repo := repositories.NewTransactionRepository(db.Get())
			if err := repo.Delete(uint(id)); err != nil {
				return output.PrintError(cmd, err)
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				return output.PrintSuccess(cmd, "Transaction deleted")
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "✓ Transaction deleted successfully")
			return nil
		},
	}
}

func newTxReconcileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reconcile TRANSACTION_ID",
		Short: "Mark transaction as reconciled",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return output.PrintError(cmd, fmt.Errorf("invalid transaction ID: %w", err))
			}

			unreconcile, _ := cmd.Flags().GetBool("unreconcile")

			repo := repositories.NewTransactionRepository(db.Get())
			if unreconcile {
				err = repo.Unreconcile(uint(id))
			} else {
				err = repo.Reconcile(uint(id))
			}

			if err != nil {
				return output.PrintError(cmd, err)
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				msg := "Transaction reconciled"
				if unreconcile {
					msg = "Transaction unreconciled"
				}
				return output.PrintSuccess(cmd, msg)
			}

			if unreconcile {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "✓ Transaction marked as unreconciled")
			} else {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "✓ Transaction marked as reconciled")
			}
			return nil
		},
	}

	cmd.Flags().BoolP("unreconcile", "u", false, "Unreconcile the transaction")
	return cmd
}
