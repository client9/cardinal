package sexpr

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// REPL represents a Read-Eval-Print Loop for s-expressions
type REPL struct {
	evaluator *Evaluator
	input     io.Reader
	output    io.Writer
	prompt    string
}

// NewREPL creates a new REPL instance
func NewREPL() *REPL {
	evaluator := NewEvaluator()
	// Set up built-in attributes for the evaluator
	setupBuiltinAttributes(evaluator.context.symbolTable)
	
	return &REPL{
		evaluator: evaluator,
		input:     os.Stdin,
		output:    os.Stdout,
		prompt:    "sexpr> ",
	}
}

// NewREPLWithIO creates a new REPL instance with custom input/output
func NewREPLWithIO(input io.Reader, output io.Writer) *REPL {
	evaluator := NewEvaluator()
	// Set up built-in attributes for the evaluator
	setupBuiltinAttributes(evaluator.context.symbolTable)
	
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

// Run starts the REPL loop
func (r *REPL) Run() error {
	scanner := bufio.NewScanner(r.input)
	
	// Print welcome message
	fmt.Fprintf(r.output, "S-Expression REPL v1.0\n")
	fmt.Fprintf(r.output, "Type 'quit' or 'exit' to exit, 'help' for help\n\n")
	
	for {
		// Print prompt
		fmt.Fprint(r.output, r.prompt)
		
		// Read input
		if !scanner.Scan() {
			break
		}
		
		line := strings.TrimSpace(scanner.Text())
		
		// Handle empty input
		if line == "" {
			continue
		}
		
		// Handle special commands
		if r.handleSpecialCommands(line) {
			continue
		}
		
		// Parse and evaluate
		if err := r.processLine(line); err != nil {
			fmt.Fprintf(r.output, "Error: %v\n", err)
		}
	}
	
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %v", err)
	}
	
	return nil
}

// handleSpecialCommands handles special REPL commands
func (r *REPL) handleSpecialCommands(line string) bool {
	switch line {
	case "quit", "exit":
		fmt.Fprintf(r.output, "Goodbye!\n")
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
	expr, err := ParseString(line)
	if err != nil {
		return fmt.Errorf("parse error: %v", err)
	}
	
	// Evaluate the expression
	result := r.evaluator.Evaluate(expr)
	
	// Print the result
	fmt.Fprintf(r.output, "%s\n", result.String())
	
	return nil
}

// printHelp prints help information
func (r *REPL) printHelp() {
	fmt.Fprintf(r.output, `
S-Expression REPL Help
======================

Commands:
  quit, exit     - Exit the REPL
  help           - Show this help message
  clear          - Clear all variable assignments
  attributes     - Show all symbols with their attributes

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
	r.evaluator = NewEvaluator()
	setupBuiltinAttributes(r.evaluator.context.symbolTable)
}

// printAttributes prints all symbols with their attributes
func (r *REPL) printAttributes() {
	fmt.Fprintf(r.output, "\nSymbols with attributes:\n")
	fmt.Fprintf(r.output, "=======================\n")
	
	symbols := r.evaluator.context.symbolTable.AllSymbolsWithAttributes()
	if len(symbols) == 0 {
		fmt.Fprintf(r.output, "No symbols with attributes found.\n")
		return
	}
	
	for _, symbol := range symbols {
		attrs := r.evaluator.context.symbolTable.Attributes(symbol)
		fmt.Fprintf(r.output, "%-15s: %s\n", symbol, AttributesToString(attrs))
	}
	fmt.Fprintf(r.output, "\n")
}

// EvaluateString is a convenience function for evaluating a string expression
func (r *REPL) EvaluateString(input string) (string, error) {
	expr, err := ParseString(input)
	if err != nil {
		return "", fmt.Errorf("parse error: %v", err)
	}
	
	result := r.evaluator.Evaluate(expr)
	return result.String(), nil
}

// GetEvaluator returns the underlying evaluator (for testing purposes)
func (r *REPL) GetEvaluator() *Evaluator {
	return r.evaluator
}

// ExecuteFile executes expressions from a file
func (r *REPL) ExecuteFile(filename string) error {
	// Read file content
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	
	// Split into lines and process each one
	lines := strings.Split(string(content), "\n")
	
	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// Show what we're executing
		fmt.Fprintf(r.output, "In(%d): %s\n", lineNum+1, line)
		
		// Execute the line
		result, err := r.EvaluateString(line)
		if err != nil {
			return fmt.Errorf("error at line %d: %v", lineNum+1, err)
		}
		
		// Show the result
		fmt.Fprintf(r.output, "Out(%d): %s\n", lineNum+1, result)
	}
	
	return nil
}

