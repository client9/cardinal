package sexpr

import (
	"testing"
)

func TestInputForm_Basic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "atom string",
			input:    `"hello"`,
			expected: `"hello"`,
		},
		{
			name:     "atom integer",
			input:    "42",
			expected: "42",
		},
		{
			name:     "atom symbol",
			input:    "x",
			expected: "x",
		},
		{
			name:     "list literal",
			input:    "List(1, 2, 3)",
			expected: "[1, 2, 3]",
		},
		{
			name:     "empty list",
			input:    "List()",
			expected: "[]",
		},
		{
			name:     "addition infix",
			input:    "Plus(1, 2)",
			expected: "1 + 2",
		},
		{
			name:     "multiplication infix",
			input:    "Times(3, 4)",
			expected: "3 * 4",
		},
		{
			name:     "assignment",
			input:    "Set(x, 5)",
			expected: "x = 5",
		},
		{
			name:     "delayed assignment",
			input:    "SetDelayed(f, Plus(x, 1))",
			expected: "f := x + 1",
		},
		{
			name:     "rule",
			input:    "Rule(a, b)",
			expected: "a: b",
		},
		{
			name:     "equality",
			input:    "Equal(x, y)",
			expected: "x == y",
		},
		{
			name:     "comparison",
			input:    "Greater(a, b)",
			expected: "a > b",
		},
		{
			name:     "logical and",
			input:    "And(True, False)",
			expected: "True && False",
		},
		{
			name:     "logical or",
			input:    "Or(x, y)",
			expected: "x || y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Errorf("unexpected parse error: %v", err)
				return
			}

			result := expr.InputForm()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestInputForm_Precedence(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "addition and multiplication",
			input:    "Plus(1, Times(2, 3))",
			expected: "1 + 2 * 3",
		},
		{
			name:     "multiplication and addition with parentheses",
			input:    "Times(Plus(1, 2), 3)",
			expected: "(1 + 2) * 3",
		},
		{
			name:     "nested operations",
			input:    "Plus(Times(a, b), Times(c, d))",
			expected: "a * b + c * d",
		},
		{
			name:     "comparison in logical expression",
			input:    "And(Greater(x, 0), Less(y, 10))",
			expected: "x > 0 && y < 10",
		},
		{
			name:     "assignment with arithmetic",
			input:    "Set(x, Plus(1, Times(2, 3)))",
			expected: "x = 1 + 2 * 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Errorf("unexpected parse error: %v", err)
				return
			}

			result := expr.InputForm()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestInputForm_Association(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty association",
			input:    "Association()",
			expected: "{}",
		},
		{
			name:     "simple association",
			input:    "Association(Rule(name, \"Bob\"))",
			expected: "{name: \"Bob\"}",
		},
		{
			name:     "multiple rules",
			input:    "Association(Rule(name, \"Bob\"), Rule(age, 30))",
			expected: "{name: \"Bob\", age: 30}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Errorf("unexpected parse error: %v", err)
				return
			}

			result := expr.InputForm()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestInputForm_MultipleArguments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "multiple addition",
			input:    "Plus(1, 2, 3, 4)",
			expected: "1 + 2 + 3 + 4",
		},
		{
			name:     "multiple multiplication",
			input:    "Times(a, b, c)",
			expected: "a * b * c",
		},
		{
			name:     "multiple logical and",
			input:    "And(True, False, True)",
			expected: "True && False && True",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Errorf("unexpected parse error: %v", err)
				return
			}

			result := expr.InputForm()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestInputForm_FallbackToFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "custom function",
			input:    "MyFunction(x, y)",
			expected: "MyFunction(x, y)",
		},
		{
			name:     "Sin function",
			input:    "Sin(x)",
			expected: "Sin(x)",
		},
		{
			name:     "Head function",
			input:    "Head(List(1, 2, 3))",
			expected: "Head([1, 2, 3])",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Errorf("unexpected parse error: %v", err)
				return
			}

			result := expr.InputForm()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
