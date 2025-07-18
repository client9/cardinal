package sexpr

import (
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "simple symbol",
			input:    "x",
			expected: "x",
			hasError: false,
		},
		{
			name:     "integer",
			input:    "42",
			expected: "42",
			hasError: false,
		},
		{
			name:     "float",
			input:    "3.14",
			expected: "3.14",
			hasError: false,
		},
		{
			name:     "string",
			input:    `"hello world"`,
			expected: `"hello world"`,
			hasError: false,
		},
		{
			name:     "boolean true",
			input:    "True",
			expected: "True",
			hasError: false,
		},
		{
			name:     "boolean false",
			input:    "False",
			expected: "False",
			hasError: false,
		},
		{
			name:     "empty list",
			input:    "{}",
			expected: "List[]",
			hasError: false,
		},
		{
			name:     "simple function call",
			input:    "Plus[1, 2]",
			expected: "Plus[1, 2]",
			hasError: false,
		},
		{
			name:     "function with no args",
			input:    "Random[]",
			expected: "Random[]",
			hasError: false,
		},
		{
			name:     "function with mixed types",
			input:    `List[1, 2.5, "hello", True, x]`,
			expected: `List[1, 2.5, "hello", True, x]`,
			hasError: false,
		},
		{
			name:     "nested function",
			input:    "Plus[1, Times[2, 3]]",
			expected: "Plus[1, Times[2, 3]]",
			hasError: false,
		},
		{
			name:     "deeply nested",
			input:    "Plus[Times[2, Power[x, 2]], Minus[5, 1]]",
			expected: "Plus[Times[2, Power[x, 2]], Minus[5, 1]]",
			hasError: false,
		},
		{
			name:     "function with string parameter",
			input:    `Function["x", Power[x, 2]]`,
			expected: `Function["x", Power[x, 2]]`,
			hasError: false,
		},
		{
			name:     "escaped string",
			input:    `Print["hello\\nworld"]`,
			expected: `Print["hello\nworld"]`,
			hasError: false,
		},
		{
			name:     "simple module expression",
			input:    `Module[x, Plus[x, Times[2, x]]]`,
			expected: `Module[x, Plus[x, Times[2, x]]]`,
			hasError: false,
		},
		{
			name:     "missing closing bracket",
			input:    "Plus[1, 2",
			expected: "",
			hasError: true,
		},
		{
			name:     "missing opening bracket",
			input:    "Plus 1, 2]",
			expected: "Plus", // Parser will parse "Plus" as a symbol, then encounter issues with rest
			hasError: false,
		},
		{
			name:     "invalid token",
			input:    "Plus[1 @ 2]",
			expected: "",
			hasError: true,
		},
		{
			name:     "unclosed brace",
			input:    "{",
			expected: "",
			hasError: true,
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
			hasError: true,
		},
		{
			name:     "multiple expressions",
			input:    "Plus[1, 2] Times[3, 4]",
			expected: "Plus[1, 2]", // Parser should parse first expression
			hasError: false,
		},
		{
			name:     "simple addition",
			input:    "1 + 2",
			expected: "Plus[1, 2]",
			hasError: false,
		},
		{
			name:     "simple subtraction",
			input:    "5 - 3",
			expected: "Subtract[5, 3]",
			hasError: false,
		},
		{
			name:     "simple multiplication",
			input:    "4 * 6",
			expected: "Times[4, 6]",
			hasError: false,
		},
		{
			name:     "simple division",
			input:    "8 / 2",
			expected: "Divide[8, 2]",
			hasError: false,
		},
		{
			name:     "precedence multiplication over addition",
			input:    "1 + 2 * 3",
			expected: "Plus[1, Times[2, 3]]",
			hasError: false,
		},
		{
			name:     "precedence division over subtraction",
			input:    "10 - 6 / 2",
			expected: "Subtract[10, Divide[6, 2]]",
			hasError: false,
		},
		{
			name:     "left associativity same precedence",
			input:    "1 + 2 + 3",
			expected: "Plus[Plus[1, 2], 3]",
			hasError: false,
		},
		{
			name:     "left associativity multiplication",
			input:    "2 * 3 * 4",
			expected: "Times[Times[2, 3], 4]",
			hasError: false,
		},
		{
			name:     "complex expression",
			input:    "1 + 2 * 3 - 4 / 2",
			expected: "Subtract[Plus[1, Times[2, 3]], Divide[4, 2]]",
			hasError: false,
		},
		{
			name:     "unary minus",
			input:    "-5",
			expected: "Minus[5]",
			hasError: false,
		},
		{
			name:     "unary plus",
			input:    "+5",
			expected: "5",
			hasError: false,
		},
		{
			name:     "unary minus with expression",
			input:    "-(2 + 3)",
			expected: "Minus[Plus[2, 3]]",
			hasError: false,
		},
		{
			name:     "simple assignment",
			input:    "x = 5",
			expected: "Set[x, 5]",
			hasError: false,
		},
		{
			name:     "delayed assignment",
			input:    "f := g[x]",
			expected: "SetDelayed[f, g[x]]",
			hasError: false,
		},
		{
			name:     "unset assignment",
			input:    "x =.",
			expected: "Unset[x]",
			hasError: false,
		},
		{
			name:     "assignment with arithmetic",
			input:    "y = 2 + 3",
			expected: "Set[y, Plus[2, 3]]",
			hasError: false,
		},
		{
			name:     "assignment precedence",
			input:    "x = y + z",
			expected: "Set[x, Plus[y, z]]",
			hasError: false,
		},
		{
			name:     "equality operator",
			input:    "x == y",
			expected: "Equal[x, y]",
			hasError: false,
		},
		{
			name:     "inequality operator",
			input:    "x != y",
			expected: "Unequal[x, y]",
			hasError: false,
		},
		{
			name:     "less than operator",
			input:    "x < y",
			expected: "Less[x, y]",
			hasError: false,
		},
		{
			name:     "greater than operator",
			input:    "x > y",
			expected: "Greater[x, y]",
			hasError: false,
		},
		{
			name:     "less equal operator",
			input:    "x <= y",
			expected: "LessEqual[x, y]",
			hasError: false,
		},
		{
			name:     "greater equal operator",
			input:    "x >= y",
			expected: "GreaterEqual[x, y]",
			hasError: false,
		},
		{
			name:     "logical and operator",
			input:    "x && y",
			expected: "And[x, y]",
			hasError: false,
		},
		{
			name:     "logical or operator",
			input:    "x || y",
			expected: "Or[x, y]",
			hasError: false,
		},
		{
			name:     "sameq operator",
			input:    "x === y",
			expected: "SameQ[x, y]",
			hasError: false,
		},
		{
			name:     "unsameq operator",
			input:    "x =!= y",
			expected: "UnsameQ[x, y]",
			hasError: false,
		},
		{
			name:     "comparison precedence",
			input:    "x + y == z * w",
			expected: "Equal[Plus[x, y], Times[z, w]]",
			hasError: false,
		},
		{
			name:     "logical precedence",
			input:    "x == y && z != w",
			expected: "And[Equal[x, y], Unequal[z, w]]",
			hasError: false,
		},
		{
			name:     "complex precedence",
			input:    "x + y < z && a || b",
			expected: "Or[And[Less[Plus[x, y], z], a], b]",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			
			if tt.hasError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			
			if expr == nil {
				t.Errorf("expected expression but got nil")
				return
			}
			
			result := expr.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestParser_ParseAtoms(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType string
		expectedVal  interface{}
	}{
		{
			name:         "symbol atom",
			input:        "mySymbol",
			expectedType: "symbol",
			expectedVal:  "mySymbol",
		},
		{
			name:         "integer atom",
			input:        "123",
			expectedType: "int",
			expectedVal:  123,
		},
		{
			name:         "float atom",
			input:        "45.67",
			expectedType: "float64",
			expectedVal:  45.67,
		},
		{
			name:         "string atom",
			input:        `"test string"`,
			expectedType: "string",
			expectedVal:  "test string",
		},
		{
			name:         "boolean true atom",
			input:        "True",
			expectedType: "bool",
			expectedVal:  true,
		},
		{
			name:         "boolean false atom",
			input:        "False",
			expectedType: "bool",
			expectedVal:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			
			if expr.Type() != tt.expectedType {
				t.Errorf("expected type %q, got %q", tt.expectedType, expr.Type())
			}
			
			atom, ok := expr.(*Atom)
			if !ok {
				t.Errorf("expected Atom, got %T", expr)
				return
			}
			
			if atom.Value != tt.expectedVal {
				t.Errorf("expected value %v, got %v", tt.expectedVal, atom.Value)
			}
		})
	}
}

func TestParser_ParseLists(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedHead    string
		expectedArgCount int
	}{
		{
			name:            "simple list",
			input:           "Plus[1, 2, 3]",
			expectedHead:    "Plus",
			expectedArgCount: 3,
		},
		{
			name:            "empty list",
			input:           "{}",
			expectedHead:    "List",
			expectedArgCount: 0,
		},
		{
			name:            "no args function",
			input:           "Random[]",
			expectedHead:    "Random",
			expectedArgCount: 0,
		},
		{
			name:            "nested list",
			input:           "Plus[1, Times[2, 3]]",
			expectedHead:    "Plus",
			expectedArgCount: 2,
		},
		{
			name:            "single arg",
			input:           "Not[True]",
			expectedHead:    "Not",
			expectedArgCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			
			list, ok := expr.(*List)
			if !ok {
				t.Errorf("expected List, got %T", expr)
				return
			}
			
			
			if len(list.Elements) == 0 {
				t.Errorf("expected non-empty list")
				return
			}
			
			head, ok := list.Elements[0].(*Atom)
			if !ok {
				t.Errorf("expected head to be Atom, got %T", list.Elements[0])
				return
			}
			
			if head.Value != tt.expectedHead {
				t.Errorf("expected head %q, got %q", tt.expectedHead, head.Value)
			}
			
			argCount := len(list.Elements) - 1 // Subtract head
			if argCount != tt.expectedArgCount {
				t.Errorf("expected %d args, got %d", tt.expectedArgCount, argCount)
			}
		})
	}
}

func TestParser_ErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name:          "unclosed bracket",
			input:         "Plus[1, 2",
			expectedError: "unexpected EOF, expected ']'",
		},
		{
			name:          "invalid token in list",
			input:         "Plus[1 @ 2]",
			expectedError: "expected ',' or ']', got ILLEGAL(@)",
		},
		{
			name:          "unclosed brace",
			input:         "{",
			expectedError: "unexpected EOF, expected '}'",
		},
		{
			name:          "invalid integer",
			input:         "999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999",
			expectedError: "invalid integer",
		},
		{
			name:          "unexpected token at start",
			input:         "[1, 2]",
			expectedError: "unexpected token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseString(tt.input)
			if err == nil {
				t.Errorf("expected error but got none")
				return
			}
			
			if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("expected error to contain %q, got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestParser_StringEscaping(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "newline escape",
			input:    `Print["hello\\nworld"]`,
			expected: `Print["hello\nworld"]`,
		},
		{
			name:     "tab escape",
			input:    `Print["hello\\tworld"]`,
			expected: `Print["hello\tworld"]`,
		},
		{
			name:     "quote escape",
			input:    `Print["say \\\"hello\\\""]`,
			expected: `Print["say \"hello\""]`,
		},
		{
			name:     "backslash escape",
			input:    `Print["path\\\\file"]`,
			expected: `Print["path\\file"]`,
		},
		{
			name:     "carriage return escape",
			input:    `Print["line1\\rline2"]`,
			expected: `Print["line1\rline2"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			
			result := expr.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestParseString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		hasError bool
	}{
		{
			name:     "valid expression",
			input:    "Plus[1, 2]",
			hasError: false,
		},
		{
			name:     "invalid expression",
			input:    "Plus[1,",
			hasError: true,
		},
		{
			name:     "empty string",
			input:    "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseString(tt.input)
			
			if tt.hasError && err == nil {
				t.Errorf("expected error but got none")
			}
			
			if !tt.hasError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}