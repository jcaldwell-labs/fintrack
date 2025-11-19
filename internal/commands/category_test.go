package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCategoryCmd(t *testing.T) {
	cmd := NewCategoryCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "category", cmd.Use)
	assert.Contains(t, cmd.Aliases, "cat")
	assert.Contains(t, cmd.Aliases, "c")
	assert.Equal(t, "Manage transaction categories", cmd.Short)
	assert.True(t, cmd.HasSubCommands())

	// Verify all subcommands exist
	subcommands := map[string]string{
		"add":    "add <name> <type>",
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

func TestCategoryAddCmd_Structure(t *testing.T) {
	cmd := NewCategoryCmd()
	addCmd, _, err := cmd.Find([]string{"add"})
	assert.NoError(t, err)
	assert.NotNil(t, addCmd)
	assert.Equal(t, "add <name> <type>", addCmd.Use)
	assert.Equal(t, "Add a new category", addCmd.Short)

	// Verify flags exist
	flags := []string{"parent", "color", "icon"}
	for _, flagName := range flags {
		flag := addCmd.Flags().Lookup(flagName)
		assert.NotNil(t, flag, "Flag --%s should exist", flagName)
	}
}

func TestCategoryListCmd_Structure(t *testing.T) {
	cmd := NewCategoryCmd()
	listCmd, _, err := cmd.Find([]string{"list"})
	assert.NoError(t, err)
	assert.NotNil(t, listCmd)
	assert.Equal(t, "list", listCmd.Use)
	assert.Equal(t, "List all categories", listCmd.Short)

	// Verify filter flags exist
	flags := []string{"type", "top-level"}
	for _, flagName := range flags {
		flag := listCmd.Flags().Lookup(flagName)
		assert.NotNil(t, flag, "Flag --%s should exist", flagName)
	}
}

func TestCategoryShowCmd_Structure(t *testing.T) {
	cmd := NewCategoryCmd()
	showCmd, _, err := cmd.Find([]string{"show"})
	assert.NoError(t, err)
	assert.NotNil(t, showCmd)
	assert.Equal(t, "show <id>", showCmd.Use)
	assert.Equal(t, "Show category details", showCmd.Short)
}

func TestCategoryUpdateCmd_Structure(t *testing.T) {
	cmd := NewCategoryCmd()
	updateCmd, _, err := cmd.Find([]string{"update"})
	assert.NoError(t, err)
	assert.NotNil(t, updateCmd)
	assert.Equal(t, "update <id>", updateCmd.Use)
	assert.Equal(t, "Update a category", updateCmd.Short)

	// Verify update flags exist
	flags := []string{"name", "color", "icon"}
	for _, flagName := range flags {
		flag := updateCmd.Flags().Lookup(flagName)
		assert.NotNil(t, flag, "Flag --%s should exist", flagName)
	}
}

func TestCategoryDeleteCmd_Structure(t *testing.T) {
	cmd := NewCategoryCmd()
	deleteCmd, _, err := cmd.Find([]string{"delete"})
	assert.NoError(t, err)
	assert.NotNil(t, deleteCmd)
	assert.Equal(t, "delete <id>", deleteCmd.Use)
	assert.Equal(t, "Delete a category", deleteCmd.Short)
}

func TestCategoryCmd_Aliases(t *testing.T) {
	// Test that aliases work
	catCmd := NewCategoryCmd()

	// Test "cat" and "c" aliases
	assert.Contains(t, catCmd.Aliases, "cat")
	assert.Contains(t, catCmd.Aliases, "c")

	// Verify command structure is the same regardless of alias used
	assert.True(t, catCmd.HasSubCommands())
	assert.Equal(t, 5, len(catCmd.Commands()))
}

func TestCategoryCmd_Subcommands(t *testing.T) {
	cmd := NewCategoryCmd()
	subcommands := cmd.Commands()

	assert.Len(t, subcommands, 5, "Should have exactly 5 subcommands")

	// Verify each subcommand has RunE defined (not just Run)
	for _, subcmd := range subcommands {
		assert.NotNil(t, subcmd.RunE, "Subcommand %s should have RunE defined", subcmd.Use)
	}
}

func TestCategoryAddCmd_ArgsValidation(t *testing.T) {
	cmd := NewCategoryCmd()
	addCmd, _, err := cmd.Find([]string{"add"})
	assert.NoError(t, err)

	// Verify it requires exactly 2 arguments (name and type)
	assert.NotNil(t, addCmd.Args)
}

func TestCategoryUpdateCmd_ArgsValidation(t *testing.T) {
	cmd := NewCategoryCmd()
	updateCmd, _, err := cmd.Find([]string{"update"})
	assert.NoError(t, err)

	// Verify it requires exactly 1 argument (id)
	assert.NotNil(t, updateCmd.Args)
}

func TestCategoryShowCmd_ArgsValidation(t *testing.T) {
	cmd := NewCategoryCmd()
	showCmd, _, err := cmd.Find([]string{"show"})
	assert.NoError(t, err)

	// Verify it requires exactly 1 argument (id)
	assert.NotNil(t, showCmd.Args)
}

func TestCategoryDeleteCmd_ArgsValidation(t *testing.T) {
	cmd := NewCategoryCmd()
	deleteCmd, _, err := cmd.Find([]string{"delete"})
	assert.NoError(t, err)

	// Verify it requires exactly 1 argument (id)
	assert.NotNil(t, deleteCmd.Args)
}

func TestCategoryListCmd_NoArgsRequired(t *testing.T) {
	cmd := NewCategoryCmd()
	listCmd, _, err := cmd.Find([]string{"list"})
	assert.NoError(t, err)

	// List command should accept 0 args (all optional flags)
	assert.Nil(t, listCmd.Args)
}

func TestCategoryCmd_FlagDefaults(t *testing.T) {
	cmd := NewCategoryCmd()
	addCmd, _, err := cmd.Find([]string{"add"})
	assert.NoError(t, err)

	// Verify default values for flags
	parentFlag := addCmd.Flags().Lookup("parent")
	assert.NotNil(t, parentFlag)
	assert.Equal(t, "", parentFlag.DefValue)

	colorFlag := addCmd.Flags().Lookup("color")
	assert.NotNil(t, colorFlag)
	assert.Equal(t, "", colorFlag.DefValue)

	iconFlag := addCmd.Flags().Lookup("icon")
	assert.NotNil(t, iconFlag)
	assert.Equal(t, "", iconFlag.DefValue)
}

func TestCategoryListCmd_FlagDefaults(t *testing.T) {
	cmd := NewCategoryCmd()
	listCmd, _, err := cmd.Find([]string{"list"})
	assert.NoError(t, err)

	// Verify filter flag defaults
	typeFlag := listCmd.Flags().Lookup("type")
	assert.NotNil(t, typeFlag)
	assert.Equal(t, "", typeFlag.DefValue)

	topLevelFlag := listCmd.Flags().Lookup("top-level")
	assert.NotNil(t, topLevelFlag)
	assert.Equal(t, "false", topLevelFlag.DefValue)
}

func TestCategoryUpdateCmd_FlagDefaults(t *testing.T) {
	cmd := NewCategoryCmd()
	updateCmd, _, err := cmd.Find([]string{"update"})
	assert.NoError(t, err)

	// Verify all update flags have empty defaults (all optional)
	flags := []string{"name", "color", "icon"}
	for _, flagName := range flags {
		flag := updateCmd.Flags().Lookup(flagName)
		assert.NotNil(t, flag, "Flag --%s should exist", flagName)
		assert.Equal(t, "", flag.DefValue, "Flag --%s should have empty default", flagName)
	}
}
