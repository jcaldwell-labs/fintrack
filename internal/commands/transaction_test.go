package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Simple smoke tests to ensure transaction commands can be created
// These tests only cover command structure, not database execution

func TestNewTransactionCmd(t *testing.T) {
	cmd := NewTransactionCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "tx", cmd.Use)
	assert.Contains(t, cmd.Aliases, "t")
	assert.Contains(t, cmd.Aliases, "transaction")
	assert.Equal(t, "Manage transactions", cmd.Short)
	assert.True(t, cmd.HasSubCommands())
}

func TestNewTxAddCmd(t *testing.T) {
	cmd := newTxAddCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "add ACCOUNT_ID AMOUNT", cmd.Use)
	assert.Equal(t, "Add a new transaction", cmd.Short)

	// Check flags exist
	assert.NotNil(t, cmd.Flags().Lookup("date"))
	assert.NotNil(t, cmd.Flags().Lookup("type"))
	assert.NotNil(t, cmd.Flags().Lookup("description"))
	assert.NotNil(t, cmd.Flags().Lookup("payee"))
	assert.NotNil(t, cmd.Flags().Lookup("category"))
}

func TestNewTxListCmd(t *testing.T) {
	cmd := newTxListCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "list [ACCOUNT_ID]", cmd.Use)
	assert.Contains(t, cmd.Aliases, "ls")
	assert.Equal(t, "List transactions", cmd.Short)

	// Check flags exist
	assert.NotNil(t, cmd.Flags().Lookup("start-date"))
	assert.NotNil(t, cmd.Flags().Lookup("end-date"))
	assert.NotNil(t, cmd.Flags().Lookup("type"))
	assert.NotNil(t, cmd.Flags().Lookup("limit"))

	// Check limit default value
	limit, err := cmd.Flags().GetInt("limit")
	assert.NoError(t, err)
	assert.Equal(t, 50, limit)
}

func TestNewTxShowCmd(t *testing.T) {
	cmd := newTxShowCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "show TRANSACTION_ID", cmd.Use)
	assert.Equal(t, "Show transaction details", cmd.Short)
}

func TestNewTxUpdateCmd(t *testing.T) {
	cmd := newTxUpdateCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "update TRANSACTION_ID", cmd.Use)
	assert.Equal(t, "Update a transaction", cmd.Short)

	// Check flags exist
	assert.NotNil(t, cmd.Flags().Lookup("amount"))
	assert.NotNil(t, cmd.Flags().Lookup("date"))
	assert.NotNil(t, cmd.Flags().Lookup("type"))
	assert.NotNil(t, cmd.Flags().Lookup("description"))
	assert.NotNil(t, cmd.Flags().Lookup("payee"))
	assert.NotNil(t, cmd.Flags().Lookup("category"))
}

func TestNewTxDeleteCmd(t *testing.T) {
	cmd := newTxDeleteCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "delete TRANSACTION_ID", cmd.Use)
	assert.Contains(t, cmd.Aliases, "del")
	assert.Contains(t, cmd.Aliases, "rm")
	assert.Equal(t, "Delete a transaction", cmd.Short)
}

func TestNewTxReconcileCmd(t *testing.T) {
	cmd := newTxReconcileCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "reconcile TRANSACTION_ID", cmd.Use)
	assert.Equal(t, "Mark transaction as reconciled", cmd.Short)

	// Check unreconcile flag exists
	assert.NotNil(t, cmd.Flags().Lookup("unreconcile"))

	// Check default value is false
	unreconcile, err := cmd.Flags().GetBool("unreconcile")
	assert.NoError(t, err)
	assert.False(t, unreconcile)
}

func TestTransactionCmd_HasAllSubcommands(t *testing.T) {
	cmd := NewTransactionCmd()

	// Test that all expected subcommands exist
	subcommands := []string{"add", "list", "show", "update", "delete", "reconcile"}
	for _, subcmd := range subcommands {
		found, _, err := cmd.Find([]string{subcmd})
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, subcmd, found.Use[:len(subcmd)])
	}
}

func TestTxAddCmd_ArgsValidation(t *testing.T) {
	cmd := newTxAddCmd()

	// Test that command expects minimum 2 arguments
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Not testing execution with database, just verifying args requirement exists
	assert.NotNil(t, cmd.Args)
}

func TestTxShowCmd_ArgsValidation(t *testing.T) {
	cmd := newTxShowCmd()

	// Test that command expects exactly 1 argument
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Not testing execution with database, just verifying args requirement exists
	assert.NotNil(t, cmd.Args)
}

func TestTxUpdateCmd_ArgsValidation(t *testing.T) {
	cmd := newTxUpdateCmd()

	// Test that command expects exactly 1 argument
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Not testing execution with database, just verifying args requirement exists
	assert.NotNil(t, cmd.Args)
}

func TestTxDeleteCmd_ArgsValidation(t *testing.T) {
	cmd := newTxDeleteCmd()

	// Test that command expects exactly 1 argument
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Not testing execution with database, just verifying args requirement exists
	assert.NotNil(t, cmd.Args)
}

func TestTxReconcileCmd_ArgsValidation(t *testing.T) {
	cmd := newTxReconcileCmd()

	// Test that command expects exactly 1 argument
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Not testing execution with database, just verifying args requirement exists
	assert.NotNil(t, cmd.Args)
}
