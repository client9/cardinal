package sexpr

import (
	"strings"
	"testing"

	"github.com/client9/sexpr/core"
)

func TestFunction_RegularSyntax(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Single parameter function",
			input:    "Function(x, x + 1)(5)",
			expected: "6",
		},
		{
			name:     "Multiple parameter function",
			input:    "Function([x, y], x + y)(3, 4)",
			expected: "7",
		},
		{
			name:     "Function creation without application",
			input:    "Function(x, x * 2)",
			expected: "Function(x, Times(x, 2))",
		},
		{
			name:     "Multiple parameter function creation",
			input:    "Function([a, b, c], a * b + c)",
			expected: "Function([a, b, c], Plus(Times(a, b), c))",
		},
		{
			name:     "Zero parameter function",
			input:    "Function([], x + 1)()",
			expected: "Plus(1, x)",
		},
		{
			name:     "Zero parameter function creation",
			input:    "Function([], x + 1)",
			expected: "Function([], Plus(x, 1))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			result := evaluator.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestFunction_SlotBasedSyntax(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic slot function",
			input:    "Function($1 + $2)(10, 20)",
			expected: "30",
		},
		{
			name:     "Bare $ as first slot",
			input:    "Function($ * 2)(5)",
			expected: "10",
		},
		{
			name:     "$ and $1 both refer to first argument",
			input:    "Function($ + $1)(7)",
			expected: "14",
		},
		{
			name:     "Missing intermediate slots",
			input:    "Function($1 + $3)(100, 200, 300)",
			expected: "400",
		},
		{
			name:     "High-numbered slots",
			input:    "Function($10 + $11)(1,2,3,4,5,6,7,8,9,10,11)",
			expected: "21",
		},
		{
			name:     "Single slot function",
			input:    "Function($1 * $1)(6)",
			expected: "36",
		},
		{
			name:     "Slot function creation without application",
			input:    "Function($1 + $2)",
			expected: "Function([slot1, slot2], Plus($1, $2))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			result := evaluator.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestFunction_ConstantFunctions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Constant numeric function",
			input:    "Function(42)()",
			expected: "42",
		},
		{
			name:     "Constant string function",
			input:    "Function(\"hello\")()",
			expected: "\"hello\"",
		},
		{
			name:     "Constant expression function",
			input:    "Function(2 + 3)()",
			expected: "5",
		},
		{
			name:     "Constant function creation",
			input:    "Function(42)",
			expected: "Function([], 42)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			result := evaluator.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestFunction_ErrorCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "Wrong number of arguments to regular function",
			input:       "Function(x, x + 1)()",
			expectError: true,
		},
		{
			name:        "Wrong number of arguments to slot function",
			input:       "Function($1 + $2)(5)",
			expectError: true,
		},
		{
			name:        "Too many arguments to Function",
			input:       "Function(x, y, z)",
			expectError: true,
		},
		{
			name:        "Non-symbol parameter in regular function",
			input:       "Function(42, x + 1)",
			expectError: true,
		},
		{
			name:        "Zero parameter function should work",
			input:       "Function([], x + 1)()",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			result := evaluator.Evaluate(expr)

			if tt.expectError {
				resultStr := result.String()
				if !strings.HasPrefix(resultStr, "$Failed") {
					t.Errorf("Expected error (starting with $Failed), but got result: %s", resultStr)
				}
			} else {
				resultStr := result.String()
				if strings.HasPrefix(resultStr, "$Failed") {
					t.Errorf("Unexpected error: %s", resultStr)
				}
			}
		})
	}
}

func TestFunction_Scoping(t *testing.T) {
	tests := []struct {
		name     string
		setup    []string // setup statements to run first
		input    string
		expected string
	}{
		{
			name:     "Function parameter shadows global variable",
			setup:    []string{"x = 100"},
			input:    "Function(x, x + 1)(5)",
			expected: "6",
		},
		{
			name:     "Function can access global variables",
			setup:    []string{"y = 10"},
			input:    "Function(x, x + y)(5)",
			expected: "15",
		},
		{
			name:     "Slot function can access global variables",
			setup:    []string{"z = 20"},
			input:    "Function($1 + z)(5)",
			expected: "25",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()

			// Run setup statements
			for _, setupStmt := range tt.setup {
				setupExpr, err := ParseString(setupStmt)
				if err != nil {
					t.Fatalf("Setup parse error: %v", err)
				}
				result := evaluator.Evaluate(setupExpr)
				if core.IsError(result) {
					t.Fatalf("Setup error: %s", result.String())
				}
			}

			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			result := evaluator.Evaluate(expr)
			if err != nil {
				t.Fatalf("Evaluation error: %v", err)
			}
			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestFunction_NestedFunctions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Nested slot-based function creation",
			input:    "Function($1 + Function($2 * $3))",
			expected: "Function(slot1, Plus($1, Function([slot1, slot2, slot3], Times($2, $3))))",
		},
		{
			name:     "Nested slot-based function application",
			input:    "Function($1 + Function($2 * $3))(1)",
			expected: "Plus(1, Function([slot1, slot2, slot3], Times($2, $3)))",
		},
		{
			name:     "Deeply nested slot function creation",
			input:    "Function($1 + Function($2 + Function($3)))",
			expected: "Function(slot1, Plus($1, Function([slot1, slot2], Plus($2, Function([slot1, slot2, slot3], $3)))))",
		},
		{
			name:     "Multiple nested slots in same level",
			input:    "Function(Function($1 + $2) + Function($3 * $4))",
			expected: "Function([], Plus(Function([slot1, slot2], Plus($1, $2)), Function([slot1, slot2, slot3, slot4], Times($3, $4))))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			result := evaluator.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestFunction_AmpersandSyntax(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic & syntax tests
		{
			name:     "Simple single slot with &",
			input:    "$1 &",
			expected: "Function(slot1, $1)",
		},
		{
			name:     "Two slots with &",
			input:    "$1 + $2 &",
			expected: "Function([slot1, slot2], Plus($1, $2))",
		},
		{
			name:     "Multiple slots with complex expression",
			input:    "$1 * $2 + $3 &",
			expected: "Function([slot1, slot2, slot3], Plus(Times($1, $2), $3))",
		},
		{
			name:     "Parenthesized expression with &",
			input:    "($1 + $2) * $3 &",
			expected: "Function([slot1, slot2, slot3], Times(Plus($1, $2), $3))",
		},
		{
			name:     "Mixed slot and constant with &",
			input:    "$1 + 10 &",
			expected: "Function(slot1, Plus($1, 10))",
		},
		{
			name:     "Complex arithmetic with &",
			input:    "$1 * 2 + $2 / 3 &",
			expected: "Function([slot1, slot2], Plus(Times($1, 2), Divide($2, 3)))",
		},

		// Function application tests
		{
			name:     "Simple & function application",
			input:    "($1 * 2 &)(5)",
			expected: "10",
		},
		{
			name:     "Two parameter & function application",
			input:    "($1 + $2 &)(10, 20)",
			expected: "30",
		},
		{
			name:     "Complex & function application",
			input:    "($1 * $2 + $3 &)(2, 3, 4)",
			expected: "10",
		},
		{
			name:     "& function with mixed types",
			input:    "(Append($1, \" world\") &)(\"hello\")",
			expected: "\"hello world\"",
		},

		// Precedence tests
		{
			name:     "& has lower precedence than arithmetic",
			input:    "$1 + $2 * $3 &",
			expected: "Function([slot1, slot2, slot3], Plus($1, Times($2, $3)))",
		},
		{
			name:     "& binds to entire arithmetic expression",
			input:    "$1 + $2 - $3 &",
			expected: "Function([slot1, slot2, slot3], Subtract(Plus($1, $2), $3))",
		},
		{
			name:     "& with power operator precedence",
			input:    "$1 ^ $2 + $3 &",
			expected: "Function([slot1, slot2, slot3], Plus(Power($1, $2), $3))",
		},

		// Nested and edge cases
		{
			name:     "Constant expression with &",
			input:    "42 &",
			expected: "Function([], 42)",
		},
		{
			name:     "String expression with &",
			input:    "\"hello\" &",
			expected: "Function([], \"hello\")",
		},
		{
			name:     "& with function calls",
			input:    "Plus($1, $2) &",
			expected: "Function([slot1, slot2], Plus($1, $2))",
		},
		{
			name:     "& with nested expressions",
			input:    "Times($1, Plus($2, $3)) &",
			expected: "Function([slot1, slot2, slot3], Times($1, Plus($2, $3)))",
		},

		// High numbered slots
		{
			name:     "& with high numbered slots",
			input:    "$10 + $5 &",
			expected: "Function([slot1, slot2, slot3, slot4, slot5, slot6, slot7, slot8, slot9, slot10], Plus($10, $5))",
		},

		// Bare $ (first slot)
		{
			name:     "Bare $ with &",
			input:    "$ * 2 &",
			expected: "Function(slot1, Times($, 2))",
		},
		{
			name:     "$ and $1 mixed with &",
			input:    "$ + $1 &",
			expected: "Function(slot1, Plus($, $1))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			result := evaluator.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestFunction_AmpersandSyntaxParsing(t *testing.T) {
	// Test that & parses correctly with different precedence contexts
	tests := []struct {
		name  string
		input string
		// We'll test the parsed AST structure to ensure precedence is correct
		shouldParse bool
	}{
		{
			name:        "& with addition",
			input:       "$1 + $2 &",
			shouldParse: true,
		},
		{
			name:        "& with multiplication",
			input:       "$1 * $2 &",
			shouldParse: true,
		},
		{
			name:        "& with power",
			input:       "$1 ^ $2 &",
			shouldParse: true,
		},
		{
			name:        "& with parentheses",
			input:       "($1 + $2) &",
			shouldParse: true,
		},
		{
			name:        "& with comparison",
			input:       "$1 < $2 &",
			shouldParse: true,
		},
		{
			name:        "& with logical operations",
			input:       "$1 && $2 &",
			shouldParse: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseString(tt.input)
			if tt.shouldParse && err != nil {
				t.Errorf("Expected %s to parse successfully, but got error: %v", tt.input, err)
			}
			if !tt.shouldParse && err == nil {
				t.Errorf("Expected %s to fail parsing, but it succeeded", tt.input)
			}
		})
	}
}
