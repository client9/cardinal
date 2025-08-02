package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

)

func TestREPL_Basic(t *testing.T) {
	// Test basic REPL functionality
	input := strings.NewReader("1 + 2\n3 * 4\nquit\n")
	output := &bytes.Buffer{}

	repl := NewREPLWithIO(input, output)
	repl.SetPrompt("test> ")

	// This should not return an error (quit command exits)
	// We'll capture panic from os.Exit in a real scenario
	defer func() {
		if r := recover(); r != nil {
			// Expected behavior when quit is called
		}
	}()

	// Just test that we can create and configure the REPL
	if repl == nil {
		t.Fatal("Failed to create REPL")
	}

	// Test prompt setting
	if repl.prompt != "test> " {
		t.Errorf("Expected prompt 'test> ', got '%s'", repl.prompt)
	}
}

func TestREPL_EvaluateString(t *testing.T) {
	repl := NewREPL()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple arithmetic",
			input:    "1 + 2",
			expected: "3",
		},
		{
			name:     "function call",
			input:    "Plus(1, 2, 3)",
			expected: "6",
		},
		{
			name:     "comparison",
			input:    "Greater(5, 3)",
			expected: "True",
		},
		{
			name:     "constant",
			input:    "True",
			expected: "True",
		},
		{
			name:     "variable assignment",
			input:    "x = 5",
			expected: "5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repl.EvaluateString(tt.input)
			if err != nil {
				t.Fatalf("EvaluateString error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestREPL_EvaluateString_WithContext(t *testing.T) {
	repl := NewREPL()

	// Set a variable
	result, err := repl.EvaluateString("x = 10")
	if err != nil {
		t.Fatalf("Failed to set variable: %v", err)
	}
	if result != "10" {
		t.Errorf("Expected '10', got '%s'", result)
	}

	// Use the variable
	result, err = repl.EvaluateString("x + 5")
	if err != nil {
		t.Fatalf("Failed to use variable: %v", err)
	}
	if result != "15" {
		t.Errorf("Expected '15', got '%s'", result)
	}
}

func TestREPL_EvaluateString_Errors(t *testing.T) {
	repl := NewREPL()

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "unclosed bracket",
			input: "Plus(1, 2",
		},
		{
			name:  "invalid token",
			input: "1 + @",
		},
		{
			name:  "invalid expression",
			input: "Plus(1 2)", // Missing comma
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repl.EvaluateString(tt.input)
			if err == nil {
				t.Errorf("Expected error for input %q, but got none", tt.input)
			}
		})
	}
}

func TestREPL_ProcessLine(t *testing.T) {
	output := &bytes.Buffer{}
	repl := NewREPLWithIO(strings.NewReader(""), output)

	err := repl.processLine("1 + 2")
	if err != nil {
		t.Fatalf("processLine error: %v", err)
	}

	result := strings.TrimSpace(output.String())
	if result != "3" {
		t.Errorf("Expected '3', got '%s'", result)
	}
}

func TestREPL_SpecialCommands(t *testing.T) {
	output := &bytes.Buffer{}
	repl := NewREPLWithIO(strings.NewReader(""), output)

	// Test help command
	if !repl.handleSpecialCommands("help") {
		t.Error("help command should return true")
	}

	// Test clear command
	if !repl.handleSpecialCommands("clear") {
		t.Error("clear command should return true")
	}

	// Test attributes command
	if !repl.handleSpecialCommands("attributes") {
		t.Error("attributes command should return true")
	}

	// Test non-special command
	if repl.handleSpecialCommands("1 + 2") {
		t.Error("regular expression should return false")
	}
}


func TestREPL_ClearContext(t *testing.T) {
	repl := NewREPL()

	// Set a variable
	repl.EvaluateString("x = 10")

	// Verify variable exists
	result, _ := repl.EvaluateString("x")
	if result != "10" {
		t.Errorf("Expected '10', got '%s'", result)
	}

	// Clear context
	repl.clearContext()

	// Verify variable is cleared (should return the symbol itself)
	result, _ = repl.EvaluateString("x")
	if result != "x" {
		t.Errorf("Expected 'x' (undefined), got '%s'", result)
	}
}

// Example demonstrates using the REPL programmatically
func Example_repl() {
	repl := NewREPL()

	// Evaluate some expressions
	result, _ := repl.EvaluateString("1 + 2")
	fmt.Println("1 + 2 =", result)

	result, _ = repl.EvaluateString("x = 5")
	fmt.Println("x = 5 returns", result)

	result, _ = repl.EvaluateString("x * 2")
	fmt.Println("x * 2 =", result)

	result, _ = repl.EvaluateString("Greater(10, 5)")
	fmt.Println("Greater[10, 5] =", result)

	// Output:
	// 1 + 2 = 3
	// x = 5 returns 5
	// x * 2 = 10
	// Greater[10, 5] = True
}
