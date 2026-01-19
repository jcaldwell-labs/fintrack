package commands

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fintrack/fintrack/internal/db"
	"github.com/fintrack/fintrack/internal/db/repositories"
	"github.com/fintrack/fintrack/internal/output"
	"github.com/fintrack/fintrack/internal/services"
	"github.com/spf13/cobra"
)

// NewImportCmd creates the import command
func NewImportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import data from external files (experimental)",
		Long: `Import transactions from CSV files.

⚠️  EXPERIMENTAL: This feature is under active development.
    Bank-specific mappings and edge cases may not be fully supported.`,
	}

	cmd.AddCommand(newImportCSVCmd())
	cmd.AddCommand(newImportHistoryCmd())

	return cmd
}

func newImportCSVCmd() *cobra.Command {
	var (
		accountID      string
		dateCol        int
		amountCol      int
		descCol        int
		payeeCol       int
		dateFormat     string
		noHeader       bool
		dryRun         bool
		skipDuplicates bool
		batchSize      int
	)

	cmd := &cobra.Command{
		Use:   "csv FILE",
		Short: "Import transactions from CSV file",
		Long:  "Import transactions from CSV files.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]

			// Resolve account ID
			accID, err := resolveAccountID(accountID)
			if err != nil {
				return output.PrintError(cmd, err)
			}

			// Build column mapping
			mapping := services.DefaultColumnMapping()
			mapping.DateColumn = dateCol
			mapping.AmountColumn = amountCol
			mapping.DescriptionColumn = descCol
			if payeeCol >= 0 {
				mapping.PayeeColumn = payeeCol
			}
			if dateFormat != "" {
				mapping.DateFormat = dateFormat
			}
			mapping.HasHeader = !noHeader

			// Build import options
			opts := services.ImportOptions{
				AccountID:      accID,
				Mapping:        mapping,
				DryRun:         dryRun,
				SkipDuplicates: skipDuplicates,
				BatchSize:      batchSize,
			}

			// Run import
			importer := services.NewCSVImporter(db.Get())
			result, err := importer.Import(filePath, opts)
			if err != nil {
				return output.PrintError(cmd, err)
			}

			// Output results
			if output.GetFormat(cmd) == output.FormatJSON {
				return output.Print(cmd, result)
			}

			// Table format output
			printImportSummary(cmd, filePath, result, dryRun)

			return nil
		},
	}

	cmd.Flags().StringVarP(&accountID, "account", "a", "", "Account ID or name (required)")
	cmd.Flags().IntVar(&dateCol, "date-col", 0, "Column index for date (0-based)")
	cmd.Flags().IntVar(&amountCol, "amount-col", 1, "Column index for amount")
	cmd.Flags().IntVar(&descCol, "desc-col", 2, "Column index for description")
	cmd.Flags().IntVar(&payeeCol, "payee-col", -1, "Column index for payee (optional)")
	cmd.Flags().StringVar(&dateFormat, "date-format", "", "Date format (Go time format, e.g., 2006-01-02)")
	cmd.Flags().BoolVar(&noHeader, "no-header", false, "CSV has no header row")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview import without saving")
	cmd.Flags().BoolVar(&skipDuplicates, "skip-duplicates", false, "Skip duplicate transactions")
	cmd.Flags().IntVar(&batchSize, "batch-size", 100, "Batch size for database inserts")

	_ = cmd.MarkFlagRequired("account")

	return cmd
}

func newImportHistoryCmd() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:     "history",
		Aliases: []string{"hist"},
		Short:   "Show import history",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo := repositories.NewImportHistoryRepository(db.Get())
			histories, err := repo.List(limit)
			if err != nil {
				return output.PrintError(cmd, err)
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				return output.Print(cmd, histories)
			}

			if len(histories) == 0 {
				fmt.Println("No import history found.")
				return nil
			}

			table := output.NewTable("ID", "FILE", "ACCOUNT", "IMPORTED", "SKIPPED", "FAILED", "DATE")
			for _, h := range histories {
				accountName := "N/A"
				if h.Account != nil {
					accountName = h.Account.Name
				}
				table.AddRow(
					fmt.Sprintf("%d", h.ID),
					filepath.Base(h.Filename),
					accountName,
					fmt.Sprintf("%d", h.RecordsImported),
					fmt.Sprintf("%d", h.RecordsSkipped),
					fmt.Sprintf("%d", h.RecordsFailed),
					h.ImportedAt.Format("2006-01-02 15:04"),
				)
			}
			table.Print()

			return nil
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "n", 20, "Maximum number of history records")

	return cmd
}

func resolveAccountID(idOrName string) (uint, error) {
	if idOrName == "" {
		return 0, fmt.Errorf("account is required")
	}

	// Try to parse as ID first
	if id, err := strconv.ParseUint(idOrName, 10, 32); err == nil {
		return uint(id), nil
	}

	// Otherwise, look up by name
	repo := repositories.NewAccountRepository(db.Get())
	account, err := repo.GetByName(idOrName)
	if err != nil {
		return 0, fmt.Errorf("account not found: %s", idOrName)
	}

	return account.ID, nil
}

func printImportSummary(cmd *cobra.Command, filePath string, result *services.ImportResult, dryRun bool) {
	mode := "Import"
	if dryRun {
		mode = "Dry Run"
	}

	fmt.Printf("\n%s Summary\n", mode)
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("File: %s\n", filepath.Base(filePath))
	fmt.Printf("Total records: %d\n", result.TotalRecords)
	fmt.Printf("Imported: %d\n", result.ImportedRecords)
	fmt.Printf("Skipped: %d\n", result.SkippedRecords)
	fmt.Printf("Failed: %d\n", result.FailedRecords)

	if len(result.Errors) > 0 {
		fmt.Println("\nErrors:")
		maxErrors := 10
		for i, e := range result.Errors {
			if i >= maxErrors {
				fmt.Printf("  ... and %d more errors\n", len(result.Errors)-maxErrors)
				break
			}
			fmt.Printf("  Line %d: %s\n", e.Line, e.Message)
		}
	}

	if dryRun {
		fmt.Println("\nThis was a dry run. No data was saved.")
		fmt.Println("Run without --dry-run to import transactions.")
	}
}
