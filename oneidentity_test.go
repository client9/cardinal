package sexpr

import (
	"testing"
)

func TestOneIdentityBasics(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic OneIdentity behavior
		{
			name:     "Plus with single argument",
			input:    "Plus(42)",
			expected: "42",
		},
		{
			name:     "Times with single argument",
			input:    "Times(5)",
			expected: "5",
		},
		// Power is not OneIdentity when it doesn't match any builtin pattern
		// Power(x_, y_) pattern requires 2 args, so Power(x) stays unevaluated

		// OneIdentity with nested expressions
		{
			name:     "Plus with nested single argument",
			input:    "Plus(Times(3, 4))",
			expected: "12",
		},
		{
			name:     "Times with nested single argument",
			input:    "Times(Plus(2, 3))",
			expected: "5",
		},

		// OneIdentity with variables
		{
			name:     "Plus with variable",
			input:    "Plus(x)",
			expected: "x",
		},
		{
			name:     "Times with variable",
			input:    "Times(y)",
			expected: "y",
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

func TestOneIdentityWithAttributeInteractions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// OneIdentity with Flat attribute
		{
			name:     "Plus flat and one identity",
			input:    "Plus(Plus(42))",
			expected: "42", // Flat flattens, then OneIdentity applies
		},
		{
			name:     "Times flat and one identity",
			input:    "Times(Times(5))",
			expected: "5",
		},

		// OneIdentity with Orderless attribute
		{
			name:     "Plus orderless doesn't affect single arg",
			input:    "Plus(42)",
			expected: "42", // OneIdentity applies regardless of ordering
		},

		// OneIdentity with multiple attributes
		{
			name:     "Plus flat orderless and one identity",
			input:    "Plus(Plus(Plus(7)))",
			expected: "7", // Multiple levels should all flatten then apply OneIdentity
		},

		// Edge case: OneIdentity with zero arguments (should not apply)
		{
			name:     "Plus with no arguments",
			input:    "Plus()",
			expected: "0", // Identity value, not OneIdentity behavior
		},
		{
			name:     "Times with no arguments",
			input:    "Times()",
			expected: "1", // Identity value, not OneIdentity behavior
		},

		// Edge case: OneIdentity with two or more arguments (should not apply)
		{
			name:     "Plus with two arguments",
			input:    "Plus(1, 2)",
			expected: "3", // Normal evaluation, not OneIdentity
		},
		{
			name:     "Times with multiple arguments",
			input:    "Times(2, 3, 4)",
			expected: "24", // Normal evaluation, not OneIdentity
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

func TestOneIdentityErrorPropagation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// OneIdentity should propagate errors
		{
			name:     "Plus with error argument",
			input:    "Plus(Divide(1, 0))",
			expected: "$Failed(DivisionByZero)",
		},
		{
			name:     "Times with error argument",
			input:    "Times(First(List()))",
			expected: "$Failed(PartError)", // Empty list error
		},

		// OneIdentity with nested error
		{
			name:     "Plus with nested error",
			input:    "Plus(Plus(Divide(1, 0)))",
			expected: "$Failed(DivisionByZero)",
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

func TestOneIdentityWithHoldAttributes(t *testing.T) {
	// Test that OneIdentity doesn't interfere with Hold attributes
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Hold prevents OneIdentity evaluation",
			input:    "Hold(Plus(42))",
			expected: "Hold(Plus(42))", // Hold prevents any evaluation
		},
		{
			name:     "OneIdentity after evaluating held expression",
			input:    "Evaluate(Hold(Plus(42)))",
			expected: "Hold(Plus(42))", // Hold still prevents evaluation inside Evaluate
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

func TestOneIdentityWithUserDefinedFunctions(t *testing.T) {
	// Test that user-defined functions can have OneIdentity behavior
	evaluator := NewEvaluator()

	// Define a function and give it OneIdentity attribute
	evaluateStringSimple(t, evaluator, "f(x_) := x + 10")
	evaluateStringSimple(t, evaluator, "SetAttributes(f, OneIdentity)")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "User function with OneIdentity - single argument",
			input:    "f(5)",
			expected: "5", // OneIdentity: f(x) = x, so f(5) = 5
		},
		{
			name:     "User function OneIdentity with variable",
			input:    "f(x)",
			expected: "x", // OneIdentity: f(x) = x
		},
		{
			name:     "User function OneIdentity with multiple args",
			input:    "f(1, 2)",
			expected: "f(1, 2)", // No pattern matches, stays unevaluated
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

func TestOneIdentityComplexInteractions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// OneIdentity with Listable functions
		{
			name:     "Plus one identity with list",
			input:    "Plus(List(1, 2, 3))",
			expected: "List(1, 2, 3)", // OneIdentity returns the list directly
		},

		// OneIdentity with symbolic expressions
		{
			name:     "Plus one identity with symbol",
			input:    "Plus(Pi)",
			expected: "3.141592653589793", // OneIdentity returns Pi, which evaluates
		},
		{
			name:     "Times one identity with symbol",
			input:    "Times(E)",
			expected: "2.718281828459045", // OneIdentity returns E, which evaluates
		},

		// OneIdentity preserves exact structure
		{
			name:     "Plus one identity with unevaluated expression",
			input:    "Plus(Hold(1 + 2))",
			expected: "Hold(Plus(1, 2))", // OneIdentity returns the Hold expression
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
