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

// FormatCurrency formats a number as currency
func FormatCurrency(amount float64, currency string) string {
	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}
	return fmt.Sprintf("%s$%.2f", sign, amount)
}

// FormatPercentage formats a number as percentage
func FormatPercentage(value float64) string {
	return fmt.Sprintf("%.0f%%", value*100)
}
