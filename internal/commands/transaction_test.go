package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewTransactionCmd tests the transaction command structure
func TestNewTransactionCmd(t *testing.T) {
	cmd := NewTransactionCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "transaction", cmd.Use)
	assert.Contains(t, cmd.Aliases, "tx")
	assert.Contains(t, cmd.Aliases, "t")
	assert.Equal(t, "Manage transactions", cmd.Short)
	assert.True(t, cmd.HasSubCommands())
}

func TestTransactionCmd_Subcommands(t *testing.T) {
	cmd := NewTransactionCmd()

	// Check that expected subcommands exist
	subcommands := []string{"list", "add", "show", "update", "delete"}
	for _, sub := range subcommands {
		found := false
		for _, c := range cmd.Commands() {
			if c.Name() == sub {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected subcommand '%s' not found", sub)
	}
}

func TestTransactionListCmd_Structure(t *testing.T) {
	cmd := NewTransactionCmd()
	listCmd, _, err := cmd.Find([]string{"list"})
	assert.NoError(t, err)
	assert.NotNil(t, listCmd)
	assert.Equal(t, "list", listCmd.Use)
	assert.Contains(t, listCmd.Aliases, "ls")
}

func TestTransactionAddCmd_Structure(t *testing.T) {
	cmd := NewTransactionCmd()
	addCmd, _, err := cmd.Find([]string{"add"})
	assert.NoError(t, err)
	assert.NotNil(t, addCmd)
	assert.Equal(t, "add", addCmd.Use)
	assert.Contains(t, addCmd.Aliases, "create")
	assert.Contains(t, addCmd.Aliases, "new")
}

func TestTransactionShowCmd_Structure(t *testing.T) {
	cmd := NewTransactionCmd()
	showCmd, _, err := cmd.Find([]string{"show"})
	assert.NoError(t, err)
	assert.NotNil(t, showCmd)
	assert.Equal(t, "show ID", showCmd.Use)
	assert.Contains(t, showCmd.Aliases, "get")
}

func TestTransactionUpdateCmd_Structure(t *testing.T) {
	cmd := NewTransactionCmd()
	updateCmd, _, err := cmd.Find([]string{"update"})
	assert.NoError(t, err)
	assert.NotNil(t, updateCmd)
	assert.Equal(t, "update ID", updateCmd.Use)
}

func TestTransactionDeleteCmd_Structure(t *testing.T) {
	cmd := NewTransactionCmd()
	deleteCmd, _, err := cmd.Find([]string{"delete"})
	assert.NoError(t, err)
	assert.NotNil(t, deleteCmd)
	assert.Equal(t, "delete ID", deleteCmd.Use)
	assert.Contains(t, deleteCmd.Aliases, "rm")
	assert.Contains(t, deleteCmd.Aliases, "remove")
}

func TestFormatAmountCents(t *testing.T) {
	tests := []struct {
		cents    int64
		expected string
	}{
		{10000, "+100.00"},   // $100.00
		{-5025, "-50.25"},    // -$50.25
		{0, "+0.00"},         // $0.00
		{-1, "-0.01"},        // -$0.01
		{123456, "+1234.56"}, // $1234.56
	}

	for _, tt := range tests {
		result := formatAmountCents(tt.cents)
		assert.Equal(t, tt.expected, result)
	}
}

func TestTransactionListCmd_Flags(t *testing.T) {
	cmd := NewTransactionCmd()
	listCmd, _, _ := cmd.Find([]string{"list"})

	// Check that expected flags exist
	flags := []string{"account", "category", "type", "from", "to", "payee", "limit"}
	for _, flag := range flags {
		f := listCmd.Flags().Lookup(flag)
		assert.NotNil(t, f, "Expected flag '%s' not found", flag)
	}
}

func TestTransactionAddCmd_Flags(t *testing.T) {
	cmd := NewTransactionCmd()
	addCmd, _, _ := cmd.Find([]string{"add"})

	// Check that expected flags exist
	flags := []string{"account", "amount", "category", "payee", "description", "type", "date", "tags"}
	for _, flag := range flags {
		f := addCmd.Flags().Lookup(flag)
		assert.NotNil(t, f, "Expected flag '%s' not found", flag)
	}
}

func TestTransactionUpdateCmd_Flags(t *testing.T) {
	cmd := NewTransactionCmd()
	updateCmd, _, _ := cmd.Find([]string{"update"})

	// Check that expected flags exist
	flags := []string{"amount", "category", "payee", "description", "date", "tags", "reconcile"}
	for _, flag := range flags {
		f := updateCmd.Flags().Lookup(flag)
		assert.NotNil(t, f, "Expected flag '%s' not found", flag)
	}
}
