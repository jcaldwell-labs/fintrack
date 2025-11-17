package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Simple smoke tests to ensure commands can be created
// These tests only cover command structure, not execution

func TestAccountStatus_Function(t *testing.T) {
	assert.Equal(t, "Active", accountStatus(true))
	assert.Equal(t, "Closed", accountStatus(false))
}
