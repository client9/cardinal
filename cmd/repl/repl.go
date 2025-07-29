package main

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"io"
	"os"
	"strings"

	"github.com/client9/sexpr"
)

// REPL represents a Read-Eval-Print Loop for s-expressions
type REPL struct {
	evaluator *sexpr.Evaluator
	input     io.Reader
	output    io.Writer
	prompt    string
}

// NewREPL creates a new REPL instance
func NewREPL() *REPL {
	evaluator := sexpr.NewEvaluator()
	// Set up built-in attributes for the evaluator
	sexpr.SetupBuiltinAttributes(evaluator.GetContext().GetSymbolTable())

	return &REPL{
		evaluator: evaluator,
		input:     os.Stdin,
		output:    os.Stdout,
		prompt:    "sexpr> ",
	}
}

// NewREPLWithIO creates a new REPL instance with custom input/output
func NewREPLWithIO(input io.Reader, output io.Writer) *REPL {
	evaluator := sexpr.NewEvaluator()
	// Set up built-in attributes for the evaluator
	sexpr.SetupBuiltinAttributes(evaluator.GetContext().GetSymbolTable())

	return &REPL{
		evaluator: evaluator,
		input:     input,
		output:    output,
		prompt:    "sexpr> ",
	}
}

// SetPrompt sets the REPL prompt
func (r *REPL) SetPrompt(prompt string) {
	r.prompt = prompt
}

// isInteractive returns true if the REPL is running in interactive mode
func (r *REPL) isInteractive() bool {
	// Check if input is stdin and if stdin is a terminal
	if r.input == os.Stdin {
		return term.IsTerminal(int(os.Stdin.Fd()))
	}
	return false
}

// Run starts the REPL loop
func (r *REPL) Run() error {
	scanner := bufio.NewScanner(r.input)

	// Print welcome message only in interactive mode (when stdin is a terminal)
	if r.isInteractive() {
		_, _ = fmt.Fprintf(r.output, "S-Expression REPL v1.0\n")
		_, _ = fmt.Fprintf(r.output, "Type 'quit' or 'exit' to exit, 'help' for help\n")
	}

	var currentExpr strings.Builder
	var emptyLineCount int

	for {
		// Print prompt only in interactive mode
		if r.isInteractive() {
			if currentExpr.Len() == 0 {
				_, _ = fmt.Fprint(r.output, r.prompt)
			} else {
				_, _ = fmt.Fprint(r.output, "   ... ")
			}
		}

		// Read input
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())

		// Handle empty input
		if line == "" {
			if currentExpr.Len() == 0 {
				continue
			}
			// Count consecutive empty lines in multi-line mode
			emptyLineCount++
			if emptyLineCount >= 2 {
				// Two empty lines in a row - abandon current expression
				if r.isInteractive() {
					_, _ = fmt.Fprintf(r.output, "Expression abandoned.\n")
				}
				currentExpr.Reset()
				emptyLineCount = 0
				continue
			}
			// Single empty line in multi-line - continue building
		} else {
			emptyLineCount = 0 // Reset empty line counter

			// Check for special reset command even when building expression
			if line == ":reset" || line == ":clear" {
				if currentExpr.Len() > 0 {
					if r.isInteractive() {
						_, _ = fmt.Fprintf(r.output, "Expression abandoned.\n")
					}
					currentExpr.Reset()
				}
				continue
			}

			// Handle special commands only if we're not building an expression
			if currentExpr.Len() == 0 && r.handleSpecialCommands(line) {
				continue
			}

			// Add line to current expression
			if currentExpr.Len() > 0 {
				// Use newline instead of space to preserve comment boundaries
				currentExpr.WriteString("\n")
			}
			currentExpr.WriteString(line)
		}

		// Try to parse and evaluate the current expression
		if currentExpr.Len() > 0 {
			expr := currentExpr.String()
			if r.tryProcessExpression(expr) {
				// Successfully processed, reset for next expression
				currentExpr.Reset()
			}
			// If parsing failed, continue accumulating lines
		}
	}

	// Handle any incomplete expression at the end
	if currentExpr.Len() > 0 {
		expr := currentExpr.String()
		if err := r.processLine(expr); err != nil {
			_, _ = fmt.Fprintf(r.output, "Error: %v\n", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %v", err)
	}

	return nil
}

// tryProcessExpression attempts to parse and evaluate an expression
// Returns true if successful, false if the expression is incomplete
func (r *REPL) tryProcessExpression(expr string) bool {
	// Try to parse the expression
	_, err := sexpr.ParseString(expr)
	if err != nil {
		errStr := err.Error()
		// Check if this looks like an incomplete expression
		if r.isIncompleteExpression(errStr) {
			return false
		}

		// This looks like a real error, not just incomplete
		_, _ = fmt.Fprintf(r.output, "Parse error: %v\n", err)
		_, _ = fmt.Fprintf(r.output, "(Type an empty line to reset if stuck)\n")
		return true // Reset the expression
	}

	// Parse succeeded, now evaluate it
	if err := r.processLine(expr); err != nil {
		_, _ = fmt.Fprintf(r.output, "Error: %v\n", err)
	}

	return true
}

// isIncompleteExpression tries to determine if a parse error indicates
// an incomplete expression (should wait for more input) vs a real error
func (r *REPL) isIncompleteExpression(errStr string) bool {
	// Common patterns that indicate incomplete expressions
	incompletePatterns := []string{
		"unexpected EOF, expected ')'",
		"unexpected EOF, expected ']'",
		"unexpected EOF, expected '}'",
		"unexpected EOF, expected ','",
		"unexpected EOF",
	}

	for _, pattern := range incompletePatterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}

	return false
}

// handleSpecialCommands handles special REPL commands
func (r *REPL) handleSpecialCommands(line string) bool {
	switch line {
	case "quit", "exit":
		if r.isInteractive() {
			_, _ = fmt.Fprintf(r.output, "Goodbye!\n")
		}
		os.Exit(0)
		return true
	case "help":
		r.printHelp()
		return true
	case "clear":
		r.clearContext()
		return true
	case "attributes":
		r.printAttributes()
		return true
	default:
		return false
	}
}

// processLine parses and evaluates a single line of input
func (r *REPL) processLine(line string) error {
	// Parse the expression
	expr, err := sexpr.ParseString(line)
	if err != nil {
		return fmt.Errorf("parse error: %v", err)
	}

	// Evaluate the expression
	result := r.evaluator.Evaluate(expr)

	// Print the result
	_, _ = fmt.Fprintf(r.output, "%s\n", result.String())

	return nil
}

// printHelp prints help information
func (r *REPL) printHelp() {
	_, _ = fmt.Fprintf(r.output, `
S-Expression REPL Help
======================

Commands:
  quit, exit     - Exit the REPL
  help           - Show this help message
  clear          - Clear all variable assignments
  attributes     - Show all symbols with their attributes
  :reset, :clear - Abandon current multi-line expression
  
Multi-line input:
  - Incomplete expressions (missing ) ] }) continue on next line
  - Type two empty lines or :reset to abandon incomplete expression

Examples:
  1 + 2 * 3                    # Arithmetic with infix notation
  Plus(1, 2, 3)                # Function call syntax
  x = 5                        # Variable assignment
  y := 2 * x                   # Delayed assignment
  If(x > 3, "big", "small")    # Conditional expression
  And(True, False)             # Logical operations
  Greater(5, 3)                # Comparison operations
  Hold(1 + 2)                  # Prevent evaluation
  Pi                           # Mathematical constants
  SameQ(3, 3)                  # Identity comparison

Operators:
  +, -, *, /     - Arithmetic operators
  ==, !=         - Equality/inequality
  <, >, <=, >=   - Comparison operators
  &&, ||         - Logical operators
  ===, =!=       - Identity operators
  =              - Assignment
  :=             - Delayed assignment

`)
}

// clearContext clears all variable assignments
func (r *REPL) clearContext() {
	r.evaluator = sexpr.NewEvaluator()
	sexpr.SetupBuiltinAttributes(r.evaluator.GetContext().GetSymbolTable())
}

// printAttributes prints all symbols with their attributes
func (r *REPL) printAttributes() {
	_, _ = fmt.Fprintf(r.output, "\nSymbols with attributes:\n")
	_, _ = fmt.Fprintf(r.output, "=======================\n")

	symbols := r.evaluator.GetContext().GetSymbolTable().AllSymbolsWithAttributes()
	if len(symbols) == 0 {
		_, _ = fmt.Fprintf(r.output, "No symbols with attributes found.\n")
		return
	}

	for _, symbol := range symbols {
		attrs := r.evaluator.GetContext().GetSymbolTable().Attributes(symbol)
		_, _ = fmt.Fprintf(r.output, "%-15s: %s\n", symbol, sexpr.AttributesToString(attrs))
	}
	_, _ = fmt.Fprintf(r.output, "\n")
}

// exprInfo represents a parsed expression with its location information
type exprInfo struct {
	text      string // The complete expression text
	startLine int    // Line number where the expression starts (1-based)
}

// parseFileContent parses file content into complete expressions, handling multi-line expressions
func (r *REPL) parseFileContent(content string) ([]exprInfo, error) {
	var expressions []exprInfo
	var currentExpr strings.Builder
	var startLine int

	lines := strings.Split(content, "\n")
	lineNum := 0

	for lineNum < len(lines) {
		lineNum++
		line := strings.TrimSpace(lines[lineNum-1])

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Start a new expression
		if currentExpr.Len() == 0 {
			startLine = lineNum
		}

		// Add the current line to the expression
		if currentExpr.Len() > 0 {
			// Use newline instead of space to preserve comment boundaries
			currentExpr.WriteString("\n")
		}
		currentExpr.WriteString(line)

		// Try to parse the current accumulated expression
		_, err := sexpr.ParseString(currentExpr.String())
		if err == nil {
			// Successfully parsed - we have a complete expression
			expressions = append(expressions, exprInfo{
				text:      currentExpr.String(),
				startLine: startLine,
			})
			currentExpr.Reset()
		} else {
			// Parse failed - this might be a multi-line expression
			// Continue accumulating lines until we get a complete expression
			// or reach end of file
		}
	}

	// Check if we have an incomplete expression at the end
	if currentExpr.Len() > 0 {
		// Try to parse one more time
		_, err := sexpr.ParseString(currentExpr.String())
		if err != nil {
			return nil, fmt.Errorf("incomplete expression starting at line %d: %v", startLine, err)
		}
		expressions = append(expressions, exprInfo{
			text:      currentExpr.String(),
			startLine: startLine,
		})
	}

	return expressions, nil
}

// EvaluateString is a convenience function for evaluating a string expression
func (r *REPL) EvaluateString(input string) (string, error) {
	expr, err := sexpr.ParseString(input)
	if err != nil {
		return "", fmt.Errorf("parse error: %v", err)
	}

	result := r.evaluator.Evaluate(expr)
	return result.String(), nil
}

// GetEvaluator returns the underlying evaluator (for testing purposes)
func (r *REPL) GetEvaluator() *sexpr.Evaluator {
	return r.evaluator
}

// ExecuteString executes expressions from a string, handling multi-line expressions
func (r *REPL) ExecuteString(content string) error {
	// Parse expressions from string content, handling multi-line expressions
	expressions, err := r.parseFileContent(content)
	if err != nil {
		return err
	}

	// Execute each complete expression, showing only final result for -c flag
	var lastResult string
	for _, exprInfo := range expressions {
		// Execute the expression
		result, err := r.EvaluateString(exprInfo.text)
		if err != nil {
			return fmt.Errorf("error in expression (line %d): %v", exprInfo.startLine, err)
		}
		lastResult = result
	}

	// For -c flag, just show the final result
	if lastResult != "" {
		_, _ = fmt.Fprintf(r.output, "%s\n", lastResult)
	}

	return nil
}

// ExecuteFile executes expressions from a file
func (r *REPL) ExecuteFile(filename string) error {
	// Read file content
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Parse expressions from file content, handling multi-line expressions
	expressions, err := r.parseFileContent(string(content))
	if err != nil {
		return err
	}

	// Execute each complete expression
	for i, exprInfo := range expressions {
		// Show what we're executing
		_, _ = fmt.Fprintf(r.output, "In(%d): %s\n", i+1, exprInfo.text)

		// Execute the expression
		result, err := r.EvaluateString(exprInfo.text)
		if err != nil {
			return fmt.Errorf("error at expression %d (line %d): %v", i+1, exprInfo.startLine, err)
		}

		// Show the result
		_, _ = fmt.Fprintf(r.output, "Out(%d): %s\n", i+1, result)
	}

	return nil
}
