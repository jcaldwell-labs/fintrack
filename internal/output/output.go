package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Format represents the output format
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
)

// Response represents a standard JSON response structure
type Response struct {
	Status string      `json:"status"` // success, error
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// GetFormat determines the output format from command flags
func GetFormat(cmd *cobra.Command) Format {
	jsonFlag, _ := cmd.Flags().GetBool("json")
	if jsonFlag {
		return FormatJSON
	}
	return FormatTable
}

// Print outputs data in the specified format
func Print(cmd *cobra.Command, data interface{}) error {
	format := GetFormat(cmd)

	if format == FormatJSON {
		return PrintJSON(os.Stdout, Response{
			Status: "success",
			Data:   data,
		})
	}

	// Default to table format - delegate to specific printers
	return nil
}

// PrintJSON outputs data as JSON
func PrintJSON(w io.Writer, data interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// PrintError outputs an error message
func PrintError(cmd *cobra.Command, err error) error {
	format := GetFormat(cmd)

	if format == FormatJSON {
		return PrintJSON(os.Stdout, Response{
			Status: "error",
			Error:  err.Error(),
		})
	}

	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	return nil
}

// PrintSuccess outputs a success message
func PrintSuccess(cmd *cobra.Command, message string) error {
	format := GetFormat(cmd)

	if format == FormatJSON {
		return PrintJSON(os.Stdout, Response{
			Status: "success",
			Data:   map[string]string{"message": message},
		})
	}

	fmt.Println(message)
	return nil
}

// Table represents a simple text table
type Table struct {
	Headers []string
	Rows    [][]string
	writer  io.Writer
}

// NewTable creates a new table
func NewTable(headers ...string) *Table {
	return &Table{
		Headers: headers,
		Rows:    [][]string{},
		writer:  os.Stdout,
	}
}

// AddRow adds a row to the table
func (t *Table) AddRow(values ...string) {
	t.Rows = append(t.Rows, values)
}

// Print prints the table
func (t *Table) Print() {
	if len(t.Rows) == 0 && len(t.Headers) == 0 {
		return
	}

	// Calculate column widths
	widths := make([]int, len(t.Headers))
	for i, header := range t.Headers {
		widths[i] = len(header)
	}

	for _, row := range t.Rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Print header
	for i, header := range t.Headers {
		fmt.Fprintf(t.writer, "%-*s  ", widths[i], header)
	}
	fmt.Fprintln(t.writer)

	// Print separator
	for i := range t.Headers {
		fmt.Fprint(t.writer, strings.Repeat("-", widths[i]), "  ")
	}
	fmt.Fprintln(t.writer)

	// Print rows
	for _, row := range t.Rows {
		for i, cell := range row {
			if i < len(widths) {
				fmt.Fprintf(t.writer, "%-*s  ", widths[i], cell)
			}
		}
		fmt.Fprintln(t.writer)
	}
}

// FormatCurrencyCents formats cents (int64) as currency with thousand separators
func FormatCurrencyCents(cents int64, currency string) string {
	sign := ""
	if cents < 0 {
		sign = "-"
		cents = -cents
	}
	dollars := cents / 100
	remainingCents := cents % 100
	return fmt.Sprintf("%s$%s.%02d", sign, formatWithCommas(dollars), remainingCents)
}

// FormatCurrency formats a dollar amount as currency (deprecated: use FormatCurrencyCents)
func FormatCurrency(amount float64, currency string) string {
	// Convert to cents and use the cents formatter for consistency
	cents := int64(amount*100 + 0.5*sign(amount))
	return FormatCurrencyCents(cents, currency)
}

// sign returns 1 for positive, -1 for negative, 0 for zero
func sign(x float64) float64 {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	}
	return 0
}

// formatWithCommas adds thousand separators to a number
func formatWithCommas(n int64) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}

	// Build the string from right to left
	s := fmt.Sprintf("%d", n)
	result := make([]byte, len(s)+(len(s)-1)/3)

	j := len(result) - 1
	for i := len(s) - 1; i >= 0; i-- {
		result[j] = s[i]
		j--
		// Add comma after every 3 digits (except at the start)
		if (len(s)-i)%3 == 0 && i > 0 {
			result[j] = ','
			j--
		}
	}

	return string(result)
}

// FormatPercentage formats a number as percentage
func FormatPercentage(value float64) string {
	return fmt.Sprintf("%.0f%%", value*100)
}
