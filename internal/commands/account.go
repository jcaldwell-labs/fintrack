package commands

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fintrack/fintrack/internal/db"
	"github.com/fintrack/fintrack/internal/db/repositories"
	"github.com/fintrack/fintrack/internal/models"
	"github.com/fintrack/fintrack/internal/output"
	"github.com/spf13/cobra"
)

// NewAccountCmd creates the account command
func NewAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account",
		Aliases: []string{"a"},
		Short:   "Manage financial accounts",
		Long: `Manage financial accounts including bank accounts, credit cards, and cash.

Examples:
  fintrack account list
  fintrack account add "Chase Checking" --type checking --balance 5000
  fintrack account show 1
  fintrack account update 1 --name "Chase Premier Checking"`,
	}

	cmd.AddCommand(newAccountListCmd())
	cmd.AddCommand(newAccountAddCmd())
	cmd.AddCommand(newAccountShowCmd())
	cmd.AddCommand(newAccountUpdateCmd())
	cmd.AddCommand(newAccountCloseCmd())

	return cmd
}

func newAccountListCmd() *cobra.Command {
	var activeOnly bool

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo := repositories.NewAccountRepository(db.Get())
			accounts, err := repo.List(activeOnly)
			if err != nil {
				return output.PrintError(cmd, err)
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				return output.Print(cmd, accounts)
			}

			// Table format
			table := output.NewTable("ID", "NAME", "TYPE", "BALANCE", "LAST ACTIVITY")
			for _, acc := range accounts {
				lastActivity := acc.UpdatedAt.Format("2006-01-02")
				table.AddRow(
					fmt.Sprintf("%d", acc.ID),
					acc.Name,
					acc.Type,
					output.FormatCurrency(acc.CurrentBalance, acc.Currency),
					lastActivity,
				)
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().BoolVar(&activeOnly, "active", true, "Show only active accounts")

	return cmd
}

func newAccountAddCmd() *cobra.Command {
	var (
		accountType   string
		balance       float64
		currency      string
		institution   string
		notes         string
	)

	cmd := &cobra.Command{
		Use:     "add NAME",
		Aliases: []string{"create"},
		Short:   "Add a new account",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			// Validate account type
			validTypes := map[string]bool{
				models.AccountTypeChecking:   true,
				models.AccountTypeSavings:    true,
				models.AccountTypeCredit:     true,
				models.AccountTypeCash:       true,
				models.AccountTypeInvestment: true,
				models.AccountTypeLoan:       true,
			}
			if !validTypes[accountType] {
				return fmt.Errorf("invalid account type: %s (valid types: checking, savings, credit, cash, investment, loan)", accountType)
			}

			account := &models.Account{
				Name:           name,
				Type:           accountType,
				Currency:       currency,
				InitialBalance: balance,
				CurrentBalance: balance,
				Institution:    institution,
				Notes:          notes,
				IsActive:       true,
			}

			repo := repositories.NewAccountRepository(db.Get())

			// Check if name already exists
			exists, err := repo.NameExists(name, nil)
			if err != nil {
				return output.PrintError(cmd, err)
			}
			if exists {
				return output.PrintError(cmd, fmt.Errorf("account with name '%s' already exists", name))
			}

			if err := repo.Create(account); err != nil {
				return output.PrintError(cmd, err)
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				return output.Print(cmd, account)
			}

			fmt.Printf("✓ Created account #%d\n", account.ID)
			fmt.Printf("Name: %s\n", account.Name)
			fmt.Printf("Type: %s\n", account.Type)
			fmt.Printf("Balance: %s\n", output.FormatCurrency(account.CurrentBalance, account.Currency))

			return nil
		},
	}

	cmd.Flags().StringVarP(&accountType, "type", "t", "checking", "Account type (checking, savings, credit, cash, investment, loan)")
	cmd.Flags().Float64VarP(&balance, "balance", "b", 0, "Initial balance")
	cmd.Flags().StringVarP(&currency, "currency", "c", "USD", "Currency code")
	cmd.Flags().StringVar(&institution, "institution", "", "Financial institution name")
	cmd.Flags().StringVar(&notes, "notes", "", "Additional notes")

	cmd.MarkFlagRequired("type")

	return cmd
}

func newAccountShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show ID",
		Aliases: []string{"get"},
		Short:   "Show account details",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseAccountID(args[0])
			if err != nil {
				return output.PrintError(cmd, err)
			}

			repo := repositories.NewAccountRepository(db.Get())
			account, err := repo.GetByID(id)
			if err != nil {
				return output.PrintError(cmd, err)
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				return output.Print(cmd, account)
			}

			// Table format
			fmt.Printf("Account #%d\n", account.ID)
			fmt.Printf("Name: %s\n", account.Name)
			fmt.Printf("Type: %s\n", account.Type)
			fmt.Printf("Currency: %s\n", account.Currency)
			fmt.Printf("Current Balance: %s\n", output.FormatCurrency(account.CurrentBalance, account.Currency))
			fmt.Printf("Initial Balance: %s\n", output.FormatCurrency(account.InitialBalance, account.Currency))
			if account.Institution != "" {
				fmt.Printf("Institution: %s\n", account.Institution)
			}
			if account.AccountNumberLast4 != "" {
				fmt.Printf("Account Number: ****%s\n", account.AccountNumberLast4)
			}
			fmt.Printf("Status: %s\n", accountStatus(account.IsActive))
			fmt.Printf("Created: %s\n", account.CreatedAt.Format("2006-01-02 15:04:05"))
			if account.Notes != "" {
				fmt.Printf("Notes: %s\n", account.Notes)
			}

			return nil
		},
	}

	return cmd
}

func newAccountUpdateCmd() *cobra.Command {
	var (
		name        string
		institution string
		notes       string
	)

	cmd := &cobra.Command{
		Use:   "update ID",
		Short: "Update account details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseAccountID(args[0])
			if err != nil {
				return output.PrintError(cmd, err)
			}

			repo := repositories.NewAccountRepository(db.Get())
			account, err := repo.GetByID(id)
			if err != nil {
				return output.PrintError(cmd, err)
			}

			// Update fields if provided
			if name != "" {
				// Check if new name conflicts
				exists, err := repo.NameExists(name, &id)
				if err != nil {
					return output.PrintError(cmd, err)
				}
				if exists {
					return output.PrintError(cmd, fmt.Errorf("account with name '%s' already exists", name))
				}
				account.Name = name
			}
			if institution != "" {
				account.Institution = institution
			}
			if notes != "" {
				account.Notes = notes
			}

			if err := repo.Update(account); err != nil {
				return output.PrintError(cmd, err)
			}

			return output.PrintSuccess(cmd, fmt.Sprintf("Account #%d updated successfully", id))
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "New account name")
	cmd.Flags().StringVar(&institution, "institution", "", "Financial institution")
	cmd.Flags().StringVar(&notes, "notes", "", "Additional notes")

	return cmd
}

func newAccountCloseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close ID",
		Short: "Close an account (soft delete)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseAccountID(args[0])
			if err != nil {
				return output.PrintError(cmd, err)
			}

			repo := repositories.NewAccountRepository(db.Get())
			if err := repo.Delete(id); err != nil {
				return output.PrintError(cmd, err)
			}

			return output.PrintSuccess(cmd, fmt.Sprintf("Account #%d closed successfully", id))
		},
	}

	return cmd
}

// Helper functions

func parseAccountID(idStr string) (uint, error) {
	// Try to parse as ID
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err == nil {
		return uint(id), nil
	}

	// If not a number, try to look up by name
	repo := repositories.NewAccountRepository(db.Get())
	account, err := repo.GetByName(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid account ID or name: %s", idStr)
	}

	return account.ID, nil
}

func accountStatus(isActive bool) string {
	if isActive {
		return "Active"
	}
	return "Closed"
}
