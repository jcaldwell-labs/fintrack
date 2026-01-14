package services

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDefaultColumnMapping(t *testing.T) {
	mapping := DefaultColumnMapping()
	assert.Equal(t, 0, mapping.DateColumn)
	assert.Equal(t, 1, mapping.AmountColumn)
	assert.Equal(t, 2, mapping.DescriptionColumn)
	assert.True(t, mapping.HasHeader)
}

func TestParseAmount(t *testing.T) {
	result, err := parseAmount("100.00")
	assert.NoError(t, err)
	assert.Equal(t, 100.0, result)

	result, err = parseAmount("-50.25")
	assert.NoError(t, err)
	assert.Equal(t, -50.25, result)

	result, err = parseAmount("1234.56")
	assert.NoError(t, err)
	assert.Equal(t, 1234.56, result)

	_, err = parseAmount("invalid")
	assert.Error(t, err)
}

func TestParseDate(t *testing.T) {
	_, err := parseDate("2024-01-15", "2006-01-02")
	assert.NoError(t, err)

	_, err = parseDate("01/15/2024", "01/02/2006")
	assert.NoError(t, err)

	_, err = parseDate("not-a-date", "2006-01-02")
	assert.Error(t, err)
}

func TestMaxInt(t *testing.T) {
	assert.Equal(t, 0, maxInt())
	assert.Equal(t, 5, maxInt(5))
	assert.Equal(t, 10, maxInt(1, 5, 10, 3))
}
