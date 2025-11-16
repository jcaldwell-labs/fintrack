package main

import (
	"fmt"
	"os"

	"github.com/fintrack/fintrack/internal/commands"
	"github.com/fintrack/fintrack/internal/config"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0-dev"
	cfgFile string
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "fintrack",
		Short: "Terminal-based personal finance tracking and budgeting",
		Long: `FinTrack is a command-line tool for managing personal finances.
Track transactions, set budgets, schedule recurring expenses,
and project cash flow - all from your terminal.

Privacy-first: All data stored locally in PostgreSQL.
Unix philosophy: Composable commands, text output, scriptable.`,
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Initialize configuration
			if err := config.Init(cfgFile); err != nil {
				return fmt.Errorf("failed to initialize config: %w", err)
			}
			return nil
		},
	}

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: $HOME/.config/fintrack/config.yaml)")
	rootCmd.PersistentFlags().Bool("json", false, "output in JSON format")

	// Add subcommands
	rootCmd.AddCommand(commands.NewAccountCmd())
	rootCmd.AddCommand(commands.NewTransactionCmd())
	rootCmd.AddCommand(commands.NewBudgetCmd())
	rootCmd.AddCommand(commands.NewScheduleCmd())
	rootCmd.AddCommand(commands.NewRemindCmd())
	rootCmd.AddCommand(commands.NewProjectCmd())
	rootCmd.AddCommand(commands.NewReportCmd())
	rootCmd.AddCommand(commands.NewCalendarCmd())
	rootCmd.AddCommand(commands.NewImportCmd())
	rootCmd.AddCommand(commands.NewConfigCmd())

	return rootCmd
}
