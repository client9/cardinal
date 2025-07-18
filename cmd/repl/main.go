package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/client9/sexpr"
)

func main() {
	// Define command line flags
	var (
		prompt = flag.String("prompt", "sexpr> ", "REPL prompt string")
		help   = flag.Bool("help", false, "Show help message")
		file   = flag.String("file", "", "Execute expressions from file instead of interactive mode")
	)
	
	flag.Parse()
	
	// Show help if requested
	if *help {
		showHelp()
		return
	}
	
	// Create REPL instance
	repl := sexpr.NewREPL()
	repl.SetPrompt(*prompt)
	
	// If file is specified, execute it
	if *file != "" {
		if err := executeFile(repl, *file); err != nil {
			fmt.Fprintf(os.Stderr, "Error executing file: %v\n", err)
			os.Exit(1)
		}
		return
	}
	
	// Start interactive REPL
	if err := repl.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "REPL error: %v\n", err)
		os.Exit(1)
	}
}

// showHelp displays help information
func showHelp() {
	fmt.Println(`S-Expression REPL - Interactive calculator and symbolic computation system

Usage:
  repl [flags)

Flags:
  -prompt string    Set the REPL prompt (default "sexpr> ")
  -file string      Execute expressions from file instead of interactive mode
  -help            Show this help message

Examples:
  repl                           # Start interactive REPL
  repl -prompt "calc> "          # Start with custom prompt
  repl -file examples.sexpr      # Execute file and exit

Interactive Commands:
  quit, exit       - Exit the REPL
  help             - Show help message
  clear            - Clear all variable assignments
  attributes       - Show all symbols with their attributes

Expression Examples:
  1 + 2 * 3                    # Arithmetic: 7
  Plus(1, 2, 3)                # Function call: 6
  x = 5                        # Variable assignment: 5
  y := 2 * x                   # Delayed assignment: Null
  If(x > 3, "big", "small")    # Conditional: "big"
  And(True, False)             # Logical: False
  Greater(5, 3)                # Comparison: True
  Hold(1 + 2)                  # Prevent evaluation: Hold(Plus(1, 2))
  Pi                           # Mathematical constant: 3.141592653589793

For more information, type 'help' in the REPL.`)
}

// executeFile executes expressions from a file
func executeFile(repl *sexpr.REPL, filename string) error {
	// Read file content
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	
	// Split into lines and process each one
	lines := splitLines(string(content))
	
	for lineNum, line := range lines {
		line = trimLine(line)
		
		// Skip empty lines and comments
		if line == "" || line[0] == '#' {
			continue
		}
		
		// Show what we're executing
		fmt.Printf("In[%d): %s\n", lineNum+1, line)
		
		// Execute the line
		result, err := repl.EvaluateString(line)
		if err != nil {
			return fmt.Errorf("error at line %d: %v", lineNum+1, err)
		}
		
		// Show the result
		fmt.Printf("Out[%d): %s\n", lineNum+1, result)
	}
	
	return nil
}

// splitLines splits content into lines
func splitLines(content string) []string {
	lines := []string{}
	current := ""
	
	for _, char := range content {
		if char == '\n' {
			lines = append(lines, current)
			current = ""
		} else if char != '\r' {
			current += string(char)
		}
	}
	
	if current != "" {
		lines = append(lines, current)
	}
	
	return lines
}

// trimLine removes leading and trailing whitespace
func trimLine(line string) string {
	// Simple trim implementation
	start := 0
	end := len(line)
	
	// Find first non-whitespace character
	for start < end && isWhitespace(line[start]) {
		start++
	}
	
	// Find last non-whitespace character
	for end > start && isWhitespace(line[end-1]) {
		end--
	}
	
	return line[start:end]
}

// isWhitespace checks if a character is whitespace
func isWhitespace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}