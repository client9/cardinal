package integration

import (
	"testing"
)

func TestEvaluateSpecialForm(t *testing.T) {
	tests := []TestCase{
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
			expected: "Evaluate(3, 12)",
		},
		// No arguments
		{
			name:     "evaluate empty",
			input:    "Evaluate()",
			expected: "Evaluate()",
		},
		{
			name:     "evaluate single argument",
			input:    "Evaluate(42)",
			expected: "42",
		},
		{
			name:     "evaluate symbol",
			input:    "x := 2.0; Evaluate(x)",
			expected: "2.0",
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
			name:      "evaluate error propagation",
			input:     "Evaluate(Divide(1, 0))",
			expected:  "",
			errorType: "DivisionByZero",
		},
		{
			name:      "evaluate multiple with error",
			input:     "Evaluate(Plus(1, 2), Divide(1, 0), Times(2, 3))",
			expected:  "",
			errorType: "DivisionByZero", // Should stop at first error
		},

		// Complex expressions
		{
			name:     "evaluate conditional",
			input:    "Evaluate(If(Greater(5, 3), Plus(1, 2), Times(2, 3)))",
			expected: "3",
		},
		{
			name:     "evaluate variable binding",
			input:    "Evaluate(Set(y, 10)); Times(y, 2))",
			expected: "20",
		},

		// Nested Evaluate
		{
			name:     "nested evaluate",
			input:    "Evaluate(Plus(Evaluate(Times(2, 3)), 4))",
			expected: "10",
		},
	}
	runTestCases(t, tests)

}

func TestEvaluateWithAttributes(t *testing.T) {
	// Test that Evaluate doesn't interfere with normal attribute processing
	tests := []TestCase{
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

	runTestCases(t, tests)
}

/*
func TestEvaluateStackBehavior(t *testing.T) {
	// Test that Evaluate properly manages the evaluation stack
	evaluator := engine.NewEvaluator()

	// Create a deeply nested expression to test stack management
	expr, err := cardinal.ParseString("Evaluate(Plus(1, Times(2, Plus(3, 4))))")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result := evaluator.Evaluate(expr)
	expected := "15" // 1 + (2 * (3 + 4)) = 1 + (2 * 7) = 15

	if result.String() != expected {
		t.Errorf("Expected %s, got %s", expected, result.String())
	}

	// TODO -- this is an interface -- need real object
	// Verify stack is clean after evaluation
	stackDepth := evaluator.GetContext().stack.Depth()
	if stackDepth != 0 {
		t.Errorf("Expected stack depth 0 after evaluation, got %d", stackDepth)
	}
}
*/
