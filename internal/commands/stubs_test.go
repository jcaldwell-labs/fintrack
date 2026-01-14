package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTransactionCmd(t *testing.T) {
	cmd := NewTransactionCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "tx", cmd.Use)
	assert.Contains(t, cmd.Aliases, "t")
	assert.Equal(t, "Manage transactions (coming soon)", cmd.Short)
	assert.NotNil(t, cmd.Run)

	// Test that running the command doesn't panic
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	assert.NotPanics(t, func() {
		cmd.Run(cmd, []string{})
	})
	assert.Contains(t, buf.String(), "coming soon")
}

func TestNewBudgetCmd(t *testing.T) {
	cmd := NewBudgetCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "budget", cmd.Use)
	assert.Contains(t, cmd.Aliases, "b")
	assert.Equal(t, "Manage budgets (coming soon)", cmd.Short)
	assert.NotNil(t, cmd.Run)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	assert.NotPanics(t, func() {
		cmd.Run(cmd, []string{})
	})
}

func TestNewScheduleCmd(t *testing.T) {
	cmd := NewScheduleCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "schedule", cmd.Use)
	assert.Contains(t, cmd.Aliases, "s")
	assert.Equal(t, "Manage recurring transactions (coming soon)", cmd.Short)
	assert.NotNil(t, cmd.Run)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	assert.NotPanics(t, func() {
		cmd.Run(cmd, []string{})
	})
}

func TestNewRemindCmd(t *testing.T) {
	cmd := NewRemindCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "remind", cmd.Use)
	assert.Contains(t, cmd.Aliases, "r")
	assert.Equal(t, "Manage reminders (coming soon)", cmd.Short)
	assert.NotNil(t, cmd.Run)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	assert.NotPanics(t, func() {
		cmd.Run(cmd, []string{})
	})
}

func TestNewProjectCmd(t *testing.T) {
	cmd := NewProjectCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "project", cmd.Use)
	assert.Contains(t, cmd.Aliases, "p")
	assert.Equal(t, "Cash flow projection (coming soon)", cmd.Short)
	assert.NotNil(t, cmd.Run)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	assert.NotPanics(t, func() {
		cmd.Run(cmd, []string{})
	})
}

func TestNewReportCmd(t *testing.T) {
	cmd := NewReportCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "report", cmd.Use)
	assert.Contains(t, cmd.Aliases, "rp")
	assert.Equal(t, "Generate reports (coming soon)", cmd.Short)
	assert.NotNil(t, cmd.Run)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	assert.NotPanics(t, func() {
		cmd.Run(cmd, []string{})
	})
}

func TestNewCalendarCmd(t *testing.T) {
	cmd := NewCalendarCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "cal", cmd.Use)
	assert.Contains(t, cmd.Aliases, "c")
	assert.Equal(t, "Calendar view (coming soon)", cmd.Short)
	assert.NotNil(t, cmd.Run)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	assert.NotPanics(t, func() {
		cmd.Run(cmd, []string{})
	})
}

func TestNewImportCmd(t *testing.T) {
	cmd := NewImportCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "import", cmd.Use)
	assert.Equal(t, "Import data from external files", cmd.Short)

	// Check for subcommands
	csvCmd, _, _ := cmd.Find([]string{"csv"})
	assert.NotNil(t, csvCmd)

	histCmd, _, _ := cmd.Find([]string{"history"})
	assert.NotNil(t, histCmd)
}

func TestNewConfigCmd(t *testing.T) {
	cmd := NewConfigCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "config", cmd.Use)
	assert.Equal(t, "Configuration management", cmd.Short)
	assert.True(t, cmd.HasSubCommands())

	// Check that the show subcommand exists
	showCmd, _, err := cmd.Find([]string{"show"})
	assert.NoError(t, err)
	assert.NotNil(t, showCmd)
	assert.Equal(t, "show", showCmd.Use)

	// Test running the show subcommand
	buf := new(bytes.Buffer)
	showCmd.SetOut(buf)
	assert.NotPanics(t, func() {
		showCmd.Run(showCmd, []string{})
	})
}
