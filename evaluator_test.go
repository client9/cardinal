package sexpr

import (
	"github.com/client9/sexpr/core"
	"testing"
)

// setupTestEvaluator creates an evaluator with built-in attributes for testing
func setupTestEvaluator() *Evaluator {
	eval := NewEvaluator()
	// Built-in attributes are already set up in NewEvaluator()
	return eval
}

func TestEvaluator_ArithmeticOperations(t *testing.T) {
	eval := setupTestEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple addition",
			input:    "Plus(1, 2)",
			expected: "3",
		},
		{
			name:     "simple multiplication",
			input:    "Times(2, 3)",
			expected: "6",
		},
		{
			name:     "subtraction",
			input:    "Subtract(5, 3)",
			expected: "2",
		},
		{
			name:     "division",
			input:    "Divide(10, 2)",
			expected: "5",
		},
		{
			name:     "power",
			input:    "Power(2, 3)",
			expected: "8.0",
		},
		{
			name:     "mixed types",
			input:    "Plus(1, 2.5)",
			expected: "3.5",
		},
		{
			name:     "multiple arguments",
			input:    "Plus(1, 2, 3, 4)",
			expected: "10",
		},
		{
			name:     "nested operations",
			input:    "Plus(1, Times(2, 3))",
			expected: "7",
		},
		{
			name:     "symbolic arithmetic",
			input:    "Plus(x, y)",
			expected: "Plus(x, y)",
		},
		{
			name:     "division by zero",
			input:    "Divide(1, 0)",
			expected: "$Failed(DivisionByZero)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := eval.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.String())
			}
		})
	}
}

func TestEvaluator_ComparisonOperations(t *testing.T) {
	eval := setupTestEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "equal numbers",
			input:    "Equal(3, 3)",
			expected: "True",
		},
		{
			name:     "unequal numbers",
			input:    "Equal(3, 4)",
			expected: "False",
		},
		{
			name:     "less than",
			input:    "Less(2, 3)",
			expected: "True",
		},
		{
			name:     "greater than",
			input:    "Greater(5, 3)",
			expected: "True",
		},
		{
			name:     "less equal",
			input:    "LessEqual(3, 3)",
			expected: "True",
		},
		{
			name:     "greater equal",
			input:    "GreaterEqual(4, 3)",
			expected: "True",
		},
		{
			name:     "unequal",
			input:    "Unequal(3, 4)",
			expected: "True",
		},
		{
			name:     "symbolic comparison",
			input:    "Equal(x, y)",
			expected: "False",
		},
		{
			name:     "same symbols",
			input:    "Equal(x, x)",
			expected: "True",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := eval.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.String())
			}
		})
	}
}

func TestEvaluator_LogicalOperations(t *testing.T) {
	eval := setupTestEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "and true true",
			input:    "And(True, True)",
			expected: "True",
		},
		{
			name:     "and true false",
			input:    "And(True, False)",
			expected: "False",
		},
		{
			name:     "or false true",
			input:    "Or(False, True)",
			expected: "True",
		},
		{
			name:     "or false false",
			input:    "Or(False, False)",
			expected: "False",
		},
		{
			name:     "not true",
			input:    "Not(True)",
			expected: "False",
		},
		{
			name:     "not false",
			input:    "Not(False)",
			expected: "True",
		},
		{
			name:     "multiple and",
			input:    "And(True, True, True)",
			expected: "True",
		},
		{
			name:     "multiple or",
			input:    "Or(False, False, True)",
			expected: "True",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := eval.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.String())
			}
		})
	}
}

func TestEvaluator_AttributeTransformations(t *testing.T) {
	eval := setupTestEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "flat transformation",
			input:    "Plus(1, Plus(2, 3))",
			expected: "6", // Should flatten and evaluate
		},
		{
			name:     "orderless transformation",
			input:    "Plus(3, 1, 2)",
			expected: "6", // Should reorder and evaluate
		},
		{
			name:     "one identity",
			input:    "Plus(5)",
			expected: "5", // Plus(x) -> x
		},
		{
			name:     "times flat",
			input:    "Times(2, Times(3, 4))",
			expected: "24", // Should flatten and evaluate
		},
		{
			name:     "symbolic flat",
			input:    "Plus(x, Plus(y, z))",
			expected: "Plus(x, y, z)", // Should flatten symbolically
		},
		{
			name:     "symbolic orderless",
			input:    "Plus(z, a, b)",
			expected: "Plus(a, b, z)", // Should reorder symbolically
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := eval.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.String())
			}
		})
	}
}

func TestEvaluator_AssignmentOperations(t *testing.T) {
	eval := setupTestEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
		setup    func()
	}{
		{
			name:     "simple assignment",
			input:    "Set(x, 5)",
			expected: "5",
		},
		{
			name:     "use assigned variable",
			input:    "Plus(x, 3)",
			expected: "8",
			setup: func() {
				expr, _ := ParseString("Set(x, 5)")
				eval.Evaluate(expr)
			},
		},
		{
			name:     "delayed assignment",
			input:    "SetDelayed(y, Plus(1, 2))",
			expected: "Null",
		},
		{
			name:     "use delayed variable",
			input:    "y",
			expected: "Plus(1, 2)",
			setup: func() {
				expr, _ := ParseString("SetDelayed(y, Plus(1, 2))")
				eval.Evaluate(expr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := eval.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.String())
			}
		})
	}
}

func TestEvaluator_ControlStructures(t *testing.T) {
	eval := setupTestEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "if true",
			input:    "If(True, 1, 2)",
			expected: "1",
		},
		{
			name:     "if false",
			input:    "If(False, 1, 2)",
			expected: "2",
		},
		{
			name:     "if no else",
			input:    "If(False, 1)",
			expected: "Null",
		},
		{
			name:     "hold expression",
			input:    "Hold(Plus(1, 2))",
			expected: "Hold(Plus(1, 2))",
		},
		{
			name:     "evaluate expression",
			input:    "Evaluate(Plus(1, 2))",
			expected: "3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := eval.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.String())
			}
		})
	}
}

func TestEvaluator_BuiltinConstants(t *testing.T) {
	eval := setupTestEvaluator()

	tests := []struct {
		name      string
		input     string
		checkFunc func(result core.Expr) bool
	}{
		{
			name:  "pi constant",
			input: "Pi",
			checkFunc: func(result core.Expr) bool {
				if real, ok := result.(core.Real); ok {
					val := float64(real)
					return val > 3.14 && val < 3.15 // Approximate check
				}
				return false
			},
		},
		{
			name:  "e constant",
			input: "E",
			checkFunc: func(result core.Expr) bool {
				if real, ok := result.(core.Real); ok {
					val := float64(real)
					return val > 2.71 && val < 2.72 // Approximate check
				}
				return false
			},
		},
		{
			name:  "true constant",
			input: "True",
			checkFunc: func(result core.Expr) bool {
				// True is now a symbol, not a BoolAtom (Mathematica compatibility)
				if symbolName, ok := core.ExtractSymbol(result); ok {
					return symbolName == "True"
				}
				return false
			},
		},
		{
			name:  "false constant",
			input: "False",
			checkFunc: func(result core.Expr) bool {
				// False is now a symbol, not a BoolAtom (Mathematica compatibility)
				if symbolName, ok := core.ExtractSymbol(result); ok {
					return symbolName == "False"
				}
				return false
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := eval.Evaluate(expr)
			if !tt.checkFunc(result) {
				t.Errorf("constant %s did not evaluate correctly: %s", tt.input, result.String())
			}
		})
	}
}

func TestEvaluator_SameQUnsameQ(t *testing.T) {
	eval := setupTestEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "same numbers",
			input:    "SameQ(3, 3)",
			expected: "True",
		},
		{
			name:     "different numbers",
			input:    "SameQ(3, 4)",
			expected: "False",
		},
		{
			name:     "same symbols",
			input:    "SameQ(x, x)",
			expected: "True",
		},
		{
			name:     "different symbols",
			input:    "SameQ(x, y)",
			expected: "False",
		},
		{
			name:     "unsame numbers",
			input:    "UnsameQ(3, 4)",
			expected: "True",
		},
		{
			name:     "unsame same",
			input:    "UnsameQ(3, 3)",
			expected: "False",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := eval.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.String())
			}
		})
	}
}

func TestEvaluator_ComplexExpressions(t *testing.T) {
	eval := setupTestEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "nested arithmetic",
			input:    "Plus(Times(2, 3), Power(2, 2))",
			expected: "10.0",
		},
		{
			name:     "comparison with arithmetic",
			input:    "Greater(Plus(2, 3), 4)",
			expected: "True",
		},
		{
			name:     "logical with comparison",
			input:    "And(Greater(5, 3), Less(2, 4))",
			expected: "True",
		},
		{
			name:     "conditional with arithmetic",
			input:    "If(Greater(5, 3), Plus(1, 2), Times(2, 3))",
			expected: "3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := eval.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.String())
			}
		})
	}
}

func TestEvaluator_InfixNotation(t *testing.T) {
	eval := setupTestEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "infix addition",
			input:    "1 + 2",
			expected: "3",
		},
		{
			name:     "infix multiplication",
			input:    "2 * 3",
			expected: "6",
		},
		{
			name:     "infix subtraction",
			input:    "5 - 3",
			expected: "2",
		},
		{
			name:     "infix division",
			input:    "10 / 2",
			expected: "5",
		},
		{
			name:     "complex infix",
			input:    "1 + 2 * 3",
			expected: "7",
		},
		{
			name:     "infix comparison",
			input:    "3 > 2",
			expected: "True",
		},
		{
			name:     "infix equality",
			input:    "3 == 3",
			expected: "True",
		},
		{
			name:     "infix logical",
			input:    "True && False",
			expected: "False",
		},
		{
			name:     "infix sameq",
			input:    "3 === 3",
			expected: "True",
		},
		{
			name:     "infix unsameq",
			input:    "3 =!= 4",
			expected: "True",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := eval.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.String())
			}
		})
	}
}
