package sexpr

import (
	"testing"
)

func TestEvaluateSpecialForm(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic evaluation
		{
			name:     "evaluate arithmetic",
			input:    "Evaluate(Plus(1, 2))",
			expected: "3",
		},
		{
			name:     "evaluate nested expression",
			input:    "Evaluate(Times(Plus(1, 2), 4))",
			expected: "12",
		},

		// Multiple arguments
		{
			name:     "evaluate multiple arguments",
			input:    "Evaluate(Plus(1, 2), Times(3, 4))",
			expected: "12", // Should return last result
		},
		{
			name:     "evaluate sequence of assignments",
			input:    "Evaluate(Set(x, 5), Plus(x, 3))",
			expected: "8",
		},

		// Edge cases
		{
			name:     "evaluate empty",
			input:    "Evaluate()",
			expected: "Null",
		},
		{
			name:     "evaluate single argument",
			input:    "Evaluate(42)",
			expected: "42",
		},
		{
			name:     "evaluate symbol",
			input:    "Evaluate(Pi)",
			expected: "3.141592653589793",
		},

		// Interaction with Hold
		{
			name:     "evaluate held expression",
			input:    "Evaluate(Hold(Plus(1, 2)))",
			expected: "Hold(Plus(1, 2))", // Hold prevents evaluation even inside Evaluate
		},
		{
			name:     "hold prevents evaluate from working",
			input:    "Hold(Evaluate(Plus(1, 2)))",
			expected: "Hold(Evaluate(Plus(1, 2)))", // Hold prevents evaluation entirely
		},

		// Error propagation
		{
			name:     "evaluate error propagation",
			input:    "Evaluate(Divide(1, 0))",
			expected: "$Failed(DivisionByZero)",
		},
		{
			name:     "evaluate multiple with error",
			input:    "Evaluate(Plus(1, 2), Divide(1, 0), Times(2, 3))",
			expected: "$Failed(DivisionByZero)", // Should stop at first error
		},

		// Complex expressions
		{
			name:     "evaluate conditional",
			input:    "Evaluate(If(Greater(5, 3), Plus(1, 2), Times(2, 3)))",
			expected: "3",
		},
		{
			name:     "evaluate variable binding",
			input:    "Evaluate(Set(y, 10), Times(y, 2))",
			expected: "20",
		},

		// Nested Evaluate
		{
			name:     "nested evaluate",
			input:    "Evaluate(Plus(Evaluate(Times(2, 3)), 4))",
			expected: "10",
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

func TestEvaluateWithAttributes(t *testing.T) {
	// Test that Evaluate doesn't interfere with normal attribute processing
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "evaluate with orderless",
			input:    "Evaluate(Plus(3, 1, 2))",
			expected: "6", // Should still apply Orderless and evaluate correctly
		},
		{
			name:     "evaluate with flat",
			input:    "Evaluate(Plus(1, Plus(2, 3)))",
			expected: "6", // Should apply Flat and evaluate
		},
		{
			name:     "evaluate with one identity",
			input:    "Evaluate(Plus(42))",
			expected: "42", // Should apply OneIdentity
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

func TestEvaluateStackBehavior(t *testing.T) {
	// Test that Evaluate properly manages the evaluation stack
	evaluator := NewEvaluator()

	// Create a deeply nested expression to test stack management
	expr, err := ParseString("Evaluate(Plus(1, Times(2, Plus(3, 4))))")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result := evaluator.Evaluate(expr)
	expected := "15" // 1 + (2 * (3 + 4)) = 1 + (2 * 7) = 15

	if result.String() != expected {
		t.Errorf("Expected %s, got %s", expected, result.String())
	}

	// Verify stack is clean after evaluation
	stackDepth := evaluator.GetContext().stack.Depth()
	if stackDepth != 0 {
		t.Errorf("Expected stack depth 0 after evaluation, got %d", stackDepth)
	}
}
