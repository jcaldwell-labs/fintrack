package usage

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestCase represents a single executable test case from markdown
type TestCase struct {
	Name           string
	Description    string
	Setup          []string
	Execute        string
	ExpectedOutput string
	ActualOutput   string
	Status         string
	LastRun        time.Time
	FilePath       string
	LineNumber     int
}

// UsageTestFile represents a parsed markdown file with usage tests
type UsageTestFile struct {
	FilePath  string
	TestCases []TestCase
}

// ParseUsageTestFile parses a markdown file and extracts test cases
func ParseUsageTestFile(filePath string) (*UsageTestFile, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	file := &UsageTestFile{
		FilePath:  filePath,
		TestCases: []TestCase{},
	}

	lines := strings.Split(string(content), "\n")
	var currentTest *TestCase
	var inCodeBlock bool
	var codeBlockType string // "setup", "execute", "expected", "actual"
	var codeBlockContent strings.Builder

	for i, line := range lines {
		// Detect test case start (## Test: ...)
		if strings.HasPrefix(line, "## Test:") {
			// Save previous test if exists
			if currentTest != nil {
				file.TestCases = append(file.TestCases, *currentTest)
			}

			// Start new test
			currentTest = &TestCase{
				Name:       strings.TrimSpace(strings.TrimPrefix(line, "## Test:")),
				FilePath:   filePath,
				LineNumber: i + 1,
			}
			continue
		}

		if currentTest == nil {
			continue
		}

		// Extract description from **Purpose:** line
		if strings.HasPrefix(line, "**Purpose:**") {
			currentTest.Description = strings.TrimSpace(strings.TrimPrefix(line, "**Purpose:**"))
			continue
		}

		// Detect code block boundaries
		if strings.HasPrefix(line, "```") {
			if !inCodeBlock {
				// Starting a code block
				inCodeBlock = true
				if strings.Contains(line, "bash") {
					// Determine type based on preceding header
					codeBlockType = "bash"
				}
				codeBlockContent.Reset()
			} else {
				// Ending a code block
				inCodeBlock = false
				content := strings.TrimSpace(codeBlockContent.String())

				// Determine what type of block this was based on context
				if codeBlockType == "bash" {
					if len(currentTest.Execute) == 0 && len(currentTest.Setup) > 0 {
						// This is the execute block
						currentTest.Execute = content
					} else if len(currentTest.Execute) == 0 {
						// First bash block after setup section
						// Check if there are multiple commands (setup)
						commands := strings.Split(content, "\n")
						if len(commands) > 1 || strings.Contains(content, "# Clean") || strings.Contains(content, "||") {
							currentTest.Setup = commands
						} else {
							currentTest.Execute = content
						}
					}
				} else if strings.HasPrefix(codeBlockType, "expected") {
					currentTest.ExpectedOutput = content
				} else if strings.HasPrefix(codeBlockType, "actual") {
					currentTest.ActualOutput = content
				}

				codeBlockType = ""
			}
			continue
		}

		if inCodeBlock {
			codeBlockContent.WriteString(line + "\n")
			continue
		}

		// Detect section headers to determine code block type
		if strings.HasPrefix(line, "### Setup") {
			codeBlockType = "setup"
		} else if strings.HasPrefix(line, "### Execute") {
			codeBlockType = "execute"
		} else if strings.HasPrefix(line, "### Expected Output") {
			codeBlockType = "expected"
		} else if strings.HasPrefix(line, "### Actual Output") {
			codeBlockType = "actual"
		}

		// Extract status from status line
		if strings.HasPrefix(line, "✅ PASS") || strings.HasPrefix(line, "❌ FAIL") {
			if strings.HasPrefix(line, "✅ PASS") {
				currentTest.Status = "PASS"
			} else {
				currentTest.Status = "FAIL"
			}

			// Extract timestamp
			re := regexp.MustCompile(`\(last run: ([\d-]+)\)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				if t, err := time.Parse("2006-01-02", matches[1]); err == nil {
					currentTest.LastRun = t
				}
			}
		}
	}

	// Don't forget the last test
	if currentTest != nil {
		file.TestCases = append(file.TestCases, *currentTest)
	}

	return file, nil
}

// RunTestCase executes a single test case
func RunTestCase(t *testing.T, tc *TestCase, binaryPath string) bool {
	t.Helper()

	// Run setup commands
	for _, cmd := range tc.Setup {
		if strings.TrimSpace(cmd) == "" || strings.HasPrefix(cmd, "#") {
			continue
		}
		runCommand(t, cmd, binaryPath, false) // Don't fail on setup errors
	}

	// Run the main command
	output := runCommand(t, tc.Execute, binaryPath, true)
	tc.ActualOutput = output
	tc.LastRun = time.Now()

	// Compare output
	if matchOutput(tc.ExpectedOutput, output) {
		tc.Status = "PASS"
		return true
	}

	tc.Status = "FAIL"
	t.Errorf("Output mismatch for test '%s':\nExpected:\n%s\n\nActual:\n%s",
		tc.Name, tc.ExpectedOutput, output)
	return false
}

// runCommand executes a shell command and returns output
func runCommand(t *testing.T, cmdStr string, binaryPath string, failOnError bool) string {
	t.Helper()

	// Replace 'fintrack' with actual binary path
	cmdStr = strings.ReplaceAll(cmdStr, "fintrack ", binaryPath+" ")

	// Split command for exec
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return ""
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Set test database environment
	cmd.Env = append(os.Environ(),
		"FINTRACK_DB_URL=postgresql://postgres:postgres@localhost:5432/fintrack_test?sslmode=disable",
	)

	err := cmd.Run()
	output := stdout.String()

	if err != nil && failOnError {
		t.Logf("Command failed: %s\nStderr: %s", cmdStr, stderr.String())
	}

	return strings.TrimSpace(output)
}

// matchOutput compares expected and actual output with wildcard support
func matchOutput(expected, actual string) bool {
	expected = strings.TrimSpace(expected)
	actual = strings.TrimSpace(actual)

	if expected == actual {
		return true
	}

	// Split into lines for comparison
	expectedLines := strings.Split(expected, "\n")
	actualLines := strings.Split(actual, "\n")

	if len(expectedLines) != len(actualLines) {
		return false
	}

	for i := range expectedLines {
		if !matchLine(expectedLines[i], actualLines[i]) {
			return false
		}
	}

	return true
}

// matchLine compares a single line with wildcard support
// Supports: <any>, <number>, <date>, <uuid>
func matchLine(expected, actual string) bool {
	expected = strings.TrimSpace(expected)
	actual = strings.TrimSpace(actual)

	// Exact match
	if expected == actual {
		return true
	}

	// Build regex pattern from expected with wildcards
	pattern := regexp.QuoteMeta(expected)

	// Replace wildcards with regex patterns
	pattern = strings.ReplaceAll(pattern, `\<any\>`, `.*`)
	pattern = strings.ReplaceAll(pattern, `\<number\>`, `\d+`)
	pattern = strings.ReplaceAll(pattern, `\<date\>`, `\d{4}-\d{2}-\d{2}`)
	pattern = strings.ReplaceAll(pattern, `\<uuid\>`, `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	pattern = strings.ReplaceAll(pattern, `\<money\>`, `\$[\d,]+\.\d{2}`)

	// Match full line
	pattern = "^" + pattern + "$"

	matched, err := regexp.MatchString(pattern, actual)
	return err == nil && matched
}

// UpdateMarkdownFile updates the markdown file with actual results
func UpdateMarkdownFile(file *UsageTestFile) error {
	content, err := os.ReadFile(file.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	var result strings.Builder
	var inActualBlock bool
	var inStatusLine bool
	testIndex := 0

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Check if we're at a test boundary
		if strings.HasPrefix(line, "## Test:") && testIndex < len(file.TestCases) {
			result.WriteString(line + "\n")
			continue
		}

		// Handle Actual Output section
		if strings.HasPrefix(line, "### Actual Output") {
			result.WriteString(line + "\n")
			// Skip old content until next section or code block end
			i++
			if i < len(lines) && strings.HasPrefix(lines[i], "```") {
				result.WriteString("```\n")
				// Skip until closing ```
				i++
				for i < len(lines) && !strings.HasPrefix(lines[i], "```") {
					i++
				}
				// Write new actual output
				if testIndex < len(file.TestCases) {
					result.WriteString(file.TestCases[testIndex].ActualOutput + "\n")
				}
				result.WriteString("```\n")
				continue
			}
		}

		// Handle status line
		if (strings.HasPrefix(line, "✅ PASS") || strings.HasPrefix(line, "❌ FAIL")) && testIndex < len(file.TestCases) {
			tc := file.TestCases[testIndex]
			statusEmoji := "✅"
			if tc.Status == "FAIL" {
				statusEmoji = "❌"
			}
			result.WriteString(fmt.Sprintf("%s %s (last run: %s)\n",
				statusEmoji, tc.Status, tc.LastRun.Format("2006-01-02")))
			testIndex++
			continue
		}

		result.WriteString(line + "\n")
	}

	return os.WriteFile(file.FilePath, []byte(result.String()), 0644)
}

// RunAllUsageTests discovers and runs all usage tests
func RunAllUsageTests(t *testing.T, usageDir string, binaryPath string) {
	files, err := filepath.Glob(filepath.Join(usageDir, "*.md"))
	assert.NoError(t, err, "Failed to discover usage test files")

	for _, filePath := range files {
		t.Run(filepath.Base(filePath), func(t *testing.T) {
			file, err := ParseUsageTestFile(filePath)
			assert.NoError(t, err, "Failed to parse usage test file")

			for i := range file.TestCases {
				tc := &file.TestCases[i]
				t.Run(tc.Name, func(t *testing.T) {
					RunTestCase(t, tc, binaryPath)
				})
			}

			// Update markdown with results
			err = UpdateMarkdownFile(file)
			assert.NoError(t, err, "Failed to update markdown file")
		})
	}
}
