package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestREPLMultiline(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Simple multiline without comment",
			input: `List(1,2,
3)`,
			expected: "List(1, 2, 3)",
		},
		{
			name: "Multiline with comment",
			input: `List(1,2 # wait
)`,
			expected: "List(1, 2)",
		},
		{
			name: "Complex multiline with comments",
			input: `Plus(1, # first
2,      # second  
3)      # third`,
			expected: "6",
		},
		{
			name: "Multiline with reset command",
			input: `List(1,2 # incomplete
:reset
5 + 3`,
			expected: "8",
		},
		{
			name: "Multiline with double empty line reset",
			input: `List(1,2 # incomplete


7 * 6`,
			expected: "42",
		},
		{
			name: "Nested expressions multiline",
			input: `Times(Plus(1, # inner
2), 
Plus(3, 4))`,
			expected: "21",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := strings.NewReader(test.input)
			output := &bytes.Buffer{}

			repl := NewREPLWithIO(input, output)
			err := repl.Run()
			if err != nil {
				t.Fatalf("REPL error: %v", err)
			}

			result := strings.TrimSpace(output.String())
			// Get the last line which should be the result
			lines := strings.Split(result, "\n")
			lastLine := strings.TrimSpace(lines[len(lines)-1])

			if lastLine != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, lastLine)
				t.Errorf("Full output:\n%s", result)
			}
		})
	}
}

func TestREPLIncompleteExpressionDetection(t *testing.T) {
	tests := []struct {
		name        string
		errorString string
		expected    bool
	}{
		{
			name:        "Missing closing parenthesis",
			errorString: "unexpected EOF, expected ')'",
			expected:    true,
		},
		{
			name:        "Missing closing bracket",
			errorString: "unexpected EOF, expected ']'",
			expected:    true,
		},
		{
			name:        "Missing closing brace",
			errorString: "unexpected EOF, expected '}'",
			expected:    true,
		},
		{
			name:        "General unexpected EOF",
			errorString: "unexpected EOF",
			expected:    true,
		},
		{
			name:        "Invalid symbol error",
			errorString: "undefined symbol: badSymbol",
			expected:    false,
		},
		{
			name:        "Type error",
			errorString: "type mismatch in Plus",
			expected:    false,
		},
	}

	repl := NewREPL()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := repl.isIncompleteExpression(test.errorString)
			if result != test.expected {
				t.Errorf("Expected %v for error %q, got %v", test.expected, test.errorString, result)
			}
		})
	}
}

func TestREPLExecuteString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple expression",
			input:    "Plus(1, 2)",
			expected: "3",
		},
		{
			name: "Multiline with comment",
			input: `List(1,2 #wait
)`,
			expected: "List(1, 2)",
		},
		{
			name: "Multiple expressions, returns last",
			input: `x = 5
y = x + 2
y`,
			expected: "7",
		},
		{
			name: "Complex multiline with comments",
			input: `Plus(1, # first
2,      # second  
3)      # third`,
			expected: "6",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			repl := NewREPLWithIO(nil, output)

			err := repl.ExecuteString(test.input)
			if err != nil {
				t.Fatalf("ExecuteString error: %v", err)
			}

			result := strings.TrimSpace(output.String())
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}
