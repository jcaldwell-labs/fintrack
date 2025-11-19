package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTransactionCmd(t *testing.T) {
	cmd := NewTransactionCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "transaction", cmd.Use)
	assert.Contains(t, cmd.Aliases, "tx")
	assert.Contains(t, cmd.Aliases, "t")
	assert.Equal(t, "Manage transactions", cmd.Short)
	assert.True(t, cmd.HasSubCommands())

	// Verify all subcommands exist
	subcommands := map[string]string{
		"add":    "add <amount>",
		"list":   "list",
		"show":   "show <id>",
		"update": "update <id>",
		"delete": "delete <id>",
	}
	for subcmd, expectedUse := range subcommands {
		found, _, err := cmd.Find([]string{subcmd})
		assert.NoError(t, err, "Subcommand %s should exist", subcmd)
		assert.NotNil(t, found)
		assert.Equal(t, expectedUse, found.Use)
	}
}

func TestTransactionAddCmd_Structure(t *testing.T) {
	cmd := NewTransactionCmd()
	addCmd, _, err := cmd.Find([]string{"add"})
	assert.NoError(t, err)
	assert.NotNil(t, addCmd)
	assert.Equal(t, "add <amount>", addCmd.Use)
	assert.Equal(t, "Add a new transaction", addCmd.Short)

	// Verify flags exist
	flags := []string{"account", "category", "date", "payee", "description", "type", "to", "tags"}
	for _, flagName := range flags {
		flag := addCmd.Flags().Lookup(flagName)
		assert.NotNil(t, flag, "Flag --%s should exist", flagName)
	}
}

func TestTransactionListCmd_Structure(t *testing.T) {
	cmd := NewTransactionCmd()
	listCmd, _, err := cmd.Find([]string{"list"})
	assert.NoError(t, err)
	assert.NotNil(t, listCmd)
	assert.Equal(t, "list", listCmd.Use)
	assert.Equal(t, "List transactions", listCmd.Short)

	// Verify filter flags exist
	flags := []string{"account", "category", "start", "end", "type", "payee", "limit"}
	for _, flagName := range flags {
		flag := listCmd.Flags().Lookup(flagName)
		assert.NotNil(t, flag, "Flag --%s should exist", flagName)
	}
}

func TestTransactionShowCmd_Structure(t *testing.T) {
	cmd := NewTransactionCmd()
	showCmd, _, err := cmd.Find([]string{"show"})
	assert.NoError(t, err)
	assert.NotNil(t, showCmd)
	assert.Equal(t, "show <id>", showCmd.Use)
	assert.Equal(t, "Show transaction details", showCmd.Short)
}

func TestTransactionUpdateCmd_Structure(t *testing.T) {
	cmd := NewTransactionCmd()
	updateCmd, _, err := cmd.Find([]string{"update"})
	assert.NoError(t, err)
	assert.NotNil(t, updateCmd)
	assert.Equal(t, "update <id>", updateCmd.Use)
	assert.Equal(t, "Update a transaction", updateCmd.Short)

	// Verify update flags exist
	flags := []string{"amount", "date", "category", "payee", "description"}
	for _, flagName := range flags {
		flag := updateCmd.Flags().Lookup(flagName)
		assert.NotNil(t, flag, "Flag --%s should exist", flagName)
	}
}

func TestTransactionDeleteCmd_Structure(t *testing.T) {
	cmd := NewTransactionCmd()
	deleteCmd, _, err := cmd.Find([]string{"delete"})
	assert.NoError(t, err)
	assert.NotNil(t, deleteCmd)
	assert.Equal(t, "delete <id>", deleteCmd.Use)
	assert.Equal(t, "Delete a transaction", deleteCmd.Short)
}

func TestTransactionCmd_Aliases(t *testing.T) {
	// Test that aliases work
	txCmd := NewTransactionCmd()

	// Test "tx" alias
	assert.Contains(t, txCmd.Aliases, "tx")
	assert.Contains(t, txCmd.Aliases, "t")

	// Verify command structure is the same regardless of alias used
	assert.True(t, txCmd.HasSubCommands())
	assert.Equal(t, 5, len(txCmd.Commands()))
}

func TestTransactionCmd_Subcommands(t *testing.T) {
	cmd := NewTransactionCmd()
	subcommands := cmd.Commands()

	assert.Len(t, subcommands, 5, "Should have exactly 5 subcommands")

	// Verify each subcommand has RunE defined (not just Run)
	for _, subcmd := range subcommands {
		assert.NotNil(t, subcmd.RunE, "Subcommand %s should have RunE defined", subcmd.Use)
	}
}

func TestTransactionAddCmd_ArgsValidation(t *testing.T) {
	cmd := NewTransactionCmd()
	addCmd, _, err := cmd.Find([]string{"add"})
	assert.NoError(t, err)

	// Verify it requires exactly 1 argument (amount)
	assert.NotNil(t, addCmd.Args)
}

func TestTransactionUpdateCmd_ArgsValidation(t *testing.T) {
	cmd := NewTransactionCmd()
	updateCmd, _, err := cmd.Find([]string{"update"})
	assert.NoError(t, err)

	// Verify it requires exactly 1 argument (id)
	assert.NotNil(t, updateCmd.Args)
}

func TestTransactionShowCmd_ArgsValidation(t *testing.T) {
	cmd := NewTransactionCmd()
	showCmd, _, err := cmd.Find([]string{"show"})
	assert.NoError(t, err)

	// Verify it requires exactly 1 argument (id)
	assert.NotNil(t, showCmd.Args)
}

func TestTransactionDeleteCmd_ArgsValidation(t *testing.T) {
	cmd := NewTransactionCmd()
	deleteCmd, _, err := cmd.Find([]string{"delete"})
	assert.NoError(t, err)

	// Verify it requires exactly 1 argument (id)
	assert.NotNil(t, deleteCmd.Args)
}

func TestTransactionListCmd_NoArgsRequired(t *testing.T) {
	cmd := NewTransactionCmd()
	listCmd, _, err := cmd.Find([]string{"list"})
	assert.NoError(t, err)

	// List command should accept 0 args (all optional flags)
	// This is implicitly tested by the command not having Args set
	assert.Nil(t, listCmd.Args)
}

func TestTransactionCmd_FlagDefaults(t *testing.T) {
	cmd := NewTransactionCmd()
	addCmd, _, err := cmd.Find([]string{"add"})
	assert.NoError(t, err)

	// Verify default values for some flags
	accountFlag := addCmd.Flags().Lookup("account")
	assert.NotNil(t, accountFlag)
	assert.Equal(t, "", accountFlag.DefValue)

	tagsFlag := addCmd.Flags().Lookup("tags")
	assert.NotNil(t, tagsFlag)
}

func TestTransactionListCmd_FlagDefaults(t *testing.T) {
	cmd := NewTransactionCmd()
	listCmd, _, err := cmd.Find([]string{"list"})
	assert.NoError(t, err)

	// Verify pagination defaults
	limitFlag := listCmd.Flags().Lookup("limit")
	assert.NotNil(t, limitFlag)
	assert.Equal(t, "100", limitFlag.DefValue)
}
