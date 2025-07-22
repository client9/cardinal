package sexpr

import (
	"testing"
)

func TestAndOrMixedBooleanValues(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// And with mixed boolean/non-boolean values
		{
			name:     "And with True and non-boolean",
			input:    "And(True, x)",
			expected: "x", // True && x = x (when x is not boolean)
		},
		{
			name:     "And with False and non-boolean",
			input:    "And(False, x)",
			expected: "False", // False && anything = False
		},
		{
			name:     "And with non-boolean and True",
			input:    "And(x, True)",
			expected: "x", // x && True = x (when x is not boolean)
		},
		{
			name:     "And with non-boolean and False",
			input:    "And(x, False)",
			expected: "False", // x && False = False
		},
		{
			name:     "And with two non-booleans",
			input:    "And(x, y)",
			expected: "And(x, y)", // Cannot evaluate, return symbolic form
		},
		{
			name:     "And with number and boolean",
			input:    "And(42, True)",
			expected: "42", // 42 && True = 42
		},
		{
			name:     "And with string and boolean",
			input:    "And(\"hello\", False)",
			expected: "False", // "hello" && False = False
		},

		// Or with mixed boolean/non-boolean values
		{
			name:     "Or with True and non-boolean",
			input:    "Or(True, x)",
			expected: "True", // True || anything = True
		},
		{
			name:     "Or with False and non-boolean",
			input:    "Or(False, x)",
			expected: "x", // False || x = x (when x is not boolean)
		},
		{
			name:     "Or with non-boolean and True",
			input:    "Or(x, True)",
			expected: "True", // x || True = True
		},
		{
			name:     "Or with non-boolean and False",
			input:    "Or(x, False)",
			expected: "x", // x || False = x
		},
		{
			name:     "Or with two non-booleans",
			input:    "Or(x, y)",
			expected: "Or(x, y)", // Cannot evaluate, return symbolic form
		},
		{
			name:     "Or with number and boolean",
			input:    "Or(42, False)",
			expected: "42", // 42 || False = 42
		},
		{
			name:     "Or with string and boolean",
			input:    "Or(\"hello\", True)",
			expected: "True", // "hello" || True = True
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			result := evaluateStringSimple(t, evaluator, tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestAndOrShortCircuitWithMixed(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Short-circuit behavior with mixed values
		{
			name:     "And short-circuit with False first",
			input:    "And(False, Divide(1, 0))", // Should not evaluate division by zero
			expected: "False",
		},
		{
			name:     "Or short-circuit with True first",
			input:    "Or(True, Divide(1, 0))", // Should not evaluate division by zero
			expected: "True",
		},
		{
			name:     "And with error in second position",
			input:    "And(True, Divide(1, 0))", // Should evaluate and propagate error
			expected: "$Failed(DivisionByZero)",
		},
		{
			name:     "Or with error in second position",
			input:    "Or(False, Divide(1, 0))", // Should evaluate and propagate error
			expected: "$Failed(DivisionByZero)",
		},

		// Multiple arguments with mixed types
		{
			name:     "And with multiple mixed arguments",
			input:    "And(True, x, True, y)",
			expected: "And(x, y)", // True values are eliminated, non-booleans remain
		},
		{
			name:     "Or with multiple mixed arguments",
			input:    "Or(False, x, False, y)",
			expected: "Or(x, y)", // False values are eliminated, non-booleans remain
		},
		{
			name:     "And with False in middle",
			input:    "And(True, x, False, y)",
			expected: "False", // False anywhere makes whole expression False
		},
		{
			name:     "Or with True in middle",
			input:    "Or(False, x, True, y)",
			expected: "True", // True anywhere makes whole expression True
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			result := evaluateStringSimple(t, evaluator, tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestAndOrNestedMixed(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Nested And/Or with mixed values
		{
			name:     "Nested And with mixed",
			input:    "And(Or(True, x), y)",
			expected: "y", // Or(True, x) = True, so And(True, y) = y
		},
		{
			name:     "Nested Or with mixed",
			input:    "Or(And(False, x), y)",
			expected: "y", // And(False, x) = False, so Or(False, y) = y
		},
		{
			name:     "Complex nested expression",
			input:    "And(Or(False, True), Or(x, False))",
			expected: "x", // Or(False, True) = True, Or(x, False) = x, And(True, x) = x
		},
		{
			name:     "Nested with non-evaluable parts",
			input:    "Or(And(x, y), And(False, z))",
			expected: "And(x, y)", // And(False, z) = False, Or(And(x, y), False) = And(x, y)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			result := evaluateStringSimple(t, evaluator, tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestAndOrWithVariableAssignments(t *testing.T) {
	evaluator := NewEvaluator()

	// Set up some variables
	evaluateStringSimple(t, evaluator, "Set(a, True)")
	evaluateStringSimple(t, evaluator, "Set(b, False)")
	evaluateStringSimple(t, evaluator, "Set(c, 42)")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "And with boolean variables",
			input:    "And(a, b)", // True && False
			expected: "False",
		},
		{
			name:     "Or with boolean variables",
			input:    "Or(a, b)", // True || False
			expected: "True",
		},
		{
			name:     "And with mixed variable types",
			input:    "And(a, c)", // True && 42
			expected: "42",
		},
		{
			name:     "Or with mixed variable types",
			input:    "Or(b, c)", // False || 42
			expected: "42",
		},
		{
			name:     "And with undefined variable",
			input:    "And(a, undefined_var)", // True && undefined_var
			expected: "undefined_var",         // Variable remains symbolic
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluateStringSimple(t, evaluator, tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
