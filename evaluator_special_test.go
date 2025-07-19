package sexpr

import (
	"testing"
)

func TestSpecialFormShortCircuitEvaluation(t *testing.T) {
	// Test And/Or short-circuit evaluation behavior
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// And short-circuit tests
		{"And empty", "And()", "True"},
		{"And early false", "And(False, undefined_symbol)", "False"}, // Should not evaluate undefined_symbol
		{"And all true", "And(True, True, True)", "True"},
		{"And mixed with early false", "And(True, False, undefined_symbol)", "False"}, // Should not evaluate undefined_symbol
		{"And with non-boolean", "And(True, x)", "And(True, x)"}, // Returns symbolic form when encountering non-boolean
		
		// Or short-circuit tests  
		{"Or empty", "Or()", "False"},
		{"Or early true", "Or(True, undefined_symbol)", "True"}, // Should not evaluate undefined_symbol
		{"Or all false", "Or(False, False, False)", "False"},
		{"Or mixed with early true", "Or(False, True, undefined_symbol)", "True"}, // Should not evaluate undefined_symbol
		{"Or with non-boolean", "Or(False, x)", "Or(False, x)"}, // Returns symbolic form when encountering non-boolean
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

func TestSpecialFormNonStandardEvaluation(t *testing.T) {
	// Test that special forms do NOT evaluate all arguments like regular functions
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// If should not evaluate both branches
		{"If true branch", "If(True, 42, undefined_symbol)", "42"},
		{"If false branch", "If(False, undefined_symbol, 24)", "24"},
		
		// SetDelayed should not evaluate RHS immediately
		{"SetDelayed stores unevaluated", "x := Plus(1, 2)", "Null"}, // SetDelayed returns Null
		
		// Hold should prevent evaluation
		{"Hold prevents evaluation", "Hold(Plus(1, 2))", "Hold(Plus(1, 2))"},
		
		// And/Or short-circuit
		{"And short-circuit", "And(False, Plus(undefined, symbol))", "False"}, // Should not evaluate second arg
		{"Or short-circuit", "Or(True, Plus(undefined, symbol))", "True"}, // Should not evaluate second arg
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

// Helper function to evaluate a string and return the result as a string
func evaluateStringSimple(t *testing.T, evaluator *Evaluator, input string) string {
	expr, err := ParseString(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	
	result := evaluator.Evaluate(expr)
	return result.String()
}