package commands

import (
	"github.com/spf13/cobra"
)

// Stub implementations for commands not yet implemented

func NewTransactionCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "tx",
		Aliases: []string{"t"},
		Short:   "Manage transactions (coming soon)",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Transaction management coming soon!")
		},
	}
}

func NewBudgetCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "budget",
		Aliases: []string{"b"},
		Short:   "Manage budgets (coming soon)",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Budget management coming soon!")
		},
	}
}

func NewScheduleCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "schedule",
		Aliases: []string{"s"},
		Short:   "Manage recurring transactions (coming soon)",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Recurring transaction scheduling coming soon!")
		},
	}
}

func NewRemindCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "remind",
		Aliases: []string{"r"},
		Short:   "Manage reminders (coming soon)",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Reminder management coming soon!")
		},
	}
}

func NewProjectCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "project",
		Aliases: []string{"p"},
		Short:   "Cash flow projection (coming soon)",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Cash flow projection coming soon!")
		},
	}
}

func NewReportCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "report",
		Aliases: []string{"rp"},
		Short:   "Generate reports (coming soon)",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Report generation coming soon!")
		},
	}
}

func NewCalendarCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "cal",
		Aliases: []string{"c"},
		Short:   "Calendar view (coming soon)",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Calendar view coming soon!")
		},
	}
}

func NewImportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "import",
		Short: "Import data from CSV (coming soon)",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("CSV import coming soon!")
		},
	}
}

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration management",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Configuration management coming soon!")
		},
	})

	return cmd
}
