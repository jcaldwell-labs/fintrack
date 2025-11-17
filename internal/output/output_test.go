package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetFormat_JSON(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("json", true, "JSON output")
	cmd.Flags().Set("json", "true")

	format := GetFormat(cmd)
	assert.Equal(t, FormatJSON, format)
}

func TestGetFormat_Table(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("json", false, "JSON output")

	format := GetFormat(cmd)
	assert.Equal(t, FormatTable, format)
}

func TestPrintJSON(t *testing.T) {
	var buf bytes.Buffer
	data := Response{
		Status: "success",
		Data:   map[string]string{"key": "value"},
	}

	err := PrintJSON(&buf, data)
	assert.NoError(t, err)

	var result Response
	err = json.Unmarshal(buf.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, "success", result.Status)
}

func TestPrint_JSON(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("json", true, "JSON output")
	cmd.Flags().Set("json", "true")

	data := map[string]string{"test": "data"}
	err := Print(cmd, data)
	assert.NoError(t, err)
}

func TestPrint_Table(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("json", false, "JSON output")

	data := map[string]string{"test": "data"}
	err := Print(cmd, data)
	assert.NoError(t, err)
}

func TestPrintError_JSON(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("json", true, "JSON output")
	cmd.Flags().Set("json", "true")

	err := PrintError(cmd, assert.AnError)
	assert.NoError(t, err)
}

func TestPrintSuccess_JSON(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("json", true, "JSON output")
	cmd.Flags().Set("json", "true")

	err := PrintSuccess(cmd, "success message")
	assert.NoError(t, err)
}

func TestNewTable(t *testing.T) {
	table := NewTable("Header1", "Header2", "Header3")
	assert.NotNil(t, table)
	assert.Equal(t, 3, len(table.Headers))
	assert.Equal(t, "Header1", table.Headers[0])
	assert.Equal(t, "Header2", table.Headers[1])
	assert.Equal(t, "Header3", table.Headers[2])
}

func TestTable_AddRow(t *testing.T) {
	table := NewTable("Col1", "Col2")
	table.AddRow("value1", "value2")
	table.AddRow("value3", "value4")

	assert.Equal(t, 2, len(table.Rows))
	assert.Equal(t, "value1", table.Rows[0][0])
	assert.Equal(t, "value4", table.Rows[1][1])
}

func TestTable_Print_Empty(t *testing.T) {
	var buf bytes.Buffer
	table := &Table{
		Headers: []string{},
		Rows:    [][]string{},
		writer:  &buf,
	}

	table.Print()
	assert.Equal(t, "", buf.String())
}

func TestTable_Print_WithData(t *testing.T) {
	var buf bytes.Buffer
	table := &Table{
		Headers: []string{"Name", "Type", "Balance"},
		Rows: [][]string{
			{"Checking", "checking", "$100.00"},
			{"Savings", "savings", "$500.00"},
		},
		writer: &buf,
	}

	table.Print()
	output := buf.String()

	assert.Contains(t, output, "Name")
	assert.Contains(t, output, "Type")
	assert.Contains(t, output, "Balance")
	assert.Contains(t, output, "Checking")
	assert.Contains(t, output, "Savings")
	assert.Contains(t, output, "---")
}

func TestTable_Print_AlignColumns(t *testing.T) {
	var buf bytes.Buffer
	table := &Table{
		Headers: []string{"Short", "Very Long Header"},
		Rows: [][]string{
			{"A", "B"},
			{"Long Value", "C"},
		},
		writer: &buf,
	}

	table.Print()
	output := buf.String()

	// Check that the table is formatted with proper spacing
	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.GreaterOrEqual(t, len(lines), 4) // Header, separator, 2 rows
}

func TestFormatCurrency_Positive(t *testing.T) {
	result := FormatCurrency(123.45, "USD")
	assert.Equal(t, "$123.45", result)
}

func TestFormatCurrency_Negative(t *testing.T) {
	result := FormatCurrency(-67.89, "USD")
	assert.Equal(t, "-$67.89", result)
}

func TestFormatCurrency_Zero(t *testing.T) {
	result := FormatCurrency(0, "USD")
	assert.Equal(t, "$0.00", result)
}

func TestFormatCurrency_Rounds(t *testing.T) {
	result := FormatCurrency(99.999, "USD")
	assert.Equal(t, "$100.00", result)
}

func TestFormatPercentage_Whole(t *testing.T) {
	result := FormatPercentage(0.75)
	assert.Equal(t, "75%", result)
}

func TestFormatPercentage_Zero(t *testing.T) {
	result := FormatPercentage(0)
	assert.Equal(t, "0%", result)
}

func TestFormatPercentage_GreaterThanOne(t *testing.T) {
	result := FormatPercentage(1.5)
	assert.Equal(t, "150%", result)
}

func TestResponse_JSON(t *testing.T) {
	resp := Response{
		Status: "success",
		Data:   "test data",
	}

	data, err := json.Marshal(resp)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "success")
	assert.Contains(t, string(data), "test data")
}

func TestResponse_Error(t *testing.T) {
	resp := Response{
		Status: "error",
		Error:  "test error",
	}

	data, err := json.Marshal(resp)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "error")
	assert.Contains(t, string(data), "test error")
}
