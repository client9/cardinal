package sexpr

import (
	"testing"

	"github.com/client9/sexpr/core"
)

// TestPatternSpecificityComparison tests relative pattern specificity
// This test is robust against changes to the specificity algorithm as it only
// tests relative ordering, not absolute values.
func TestPatternSpecificityComparison(t *testing.T) {
	tests := []struct {
		name         string
		moreSpecific string // Pattern that should be MORE specific
		lessSpecific string // Pattern that should be LESS specific
	}{
		// Literal vs Pattern Variable Tests
		{"Literal vs General", "42", "x_"},
		{"Literal vs Typed", "42", "x_Integer"},
		{"Literal vs Sequence", "42", "x__"},
		{"Literal vs Null Sequence", "42", "x___"},
		{"Symbol literal vs General", "Plus", "x_"},
		{"Symbol literal vs Symbol type", "Plus", "x_Symbol"},

		// Pattern Type Hierarchy Tests
		{"Single vs Sequence", "x_", "x__"},
		{"Single vs Null Sequence", "x_", "x___"},
		{"Sequence vs Null Sequence", "x__", "x___"},

		// Type Constraint Tests
		{"Typed vs General", "x_Integer", "x_"},
		{"Typed vs General sequence", "x__Integer", "x__"},
		{"Typed vs General null sequence", "x___Integer", "x___"},

		// Built-in vs User Types
		{"Builtin type vs General", "x_Integer", "x_"},
		{"User type vs General", "x_Color", "x_"},
		// Note: Builtin vs User type ordering depends on specificity constants
		// but both should be more specific than general patterns

		// Complex Pattern Tests
		// Note: Current algorithm prioritizes more arguments, but this could be debated
		{"Single arg vs Zero arg", "Plus(x_)", "Plus()"},
		{"Multi arg vs Single arg", "Plus(x_, y_)", "Plus(x_)"},
		{"Typed arg vs General arg", "Plus(x_Integer)", "Plus(x_)"},
		{"Mixed types", "Plus(42, x_Integer)", "Plus(x_, y_)"},

		// Sequence in Functions
		{"Function with single vs sequence", "Plus(x_)", "Plus(x__)"},
		{"Function with sequence vs null sequence", "Plus(x__)", "Plus(x___)"},

		// Compound Specificity
		// Current algorithm: more arguments = higher specificity
		{"More args higher specificity", "Plus(x_, y_, z_)", "Plus(x_, y_)"},
		{"Typed compound more specific", "Plus(x_Integer, y_Integer)", "Plus(x_, y_)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse both patterns
			moreSpecific, err1 := ParseString(tt.moreSpecific)
			if err1 != nil {
				t.Fatalf("Failed to parse more specific pattern %q: %v", tt.moreSpecific, err1)
			}

			lessSpecific, err2 := ParseString(tt.lessSpecific)
			if err2 != nil {
				t.Fatalf("Failed to parse less specific pattern %q: %v", tt.lessSpecific, err2)
			}

			// Calculate specificities
			moreSpec := core.GetPatternSpecificity(moreSpecific)
			lessSpec := core.GetPatternSpecificity(lessSpecific)

			// More specific should have HIGHER numeric value
			if moreSpec <= lessSpec {
				t.Errorf("Pattern %q (specificity=%d) should be MORE specific than %q (specificity=%d)",
					tt.moreSpecific, moreSpec, tt.lessSpecific, lessSpec)
			}
		})
	}
}

// TestPatternSpecificityInFunctionRegistry tests that the function registry
// correctly orders patterns by specificity
func TestPatternSpecificityInFunctionRegistry(t *testing.T) {
	testCases := []struct {
		definitions []string
		testInput   string
		expected    string
		description string
	}{
		{
			definitions: []string{
				`test(x_) := "general"`,
				`test(42) := "literal"`,
				`test(x_Integer) := "integer"`,
			},
			testInput:   "test(42)",
			expected:    `"literal"`,
			description: "Literal should win over typed and general",
		},
		{
			definitions: []string{
				`test(x_) := "general"`,
				`test(42) := "literal"`,
				`test(x_Integer) := "integer"`,
			},
			testInput:   "test(99)",
			expected:    `"integer"`,
			description: "Typed should win over general for integers",
		},
		{
			definitions: []string{
				`test(x_) := "general"`,
				`test(42) := "literal"`,
				`test(x_Integer) := "integer"`,
			},
			testInput:   `test("hello")`,
			expected:    `"general"`,
			description: "General should match non-integers",
		},
		{
			definitions: []string{
				`seq(x___) := "null sequence"`,
				`seq(x__) := "sequence"`,
				`seq(x_) := "single"`,
			},
			testInput:   "seq(1)",
			expected:    `"single"`,
			description: "Single pattern most specific for one arg",
		},
		{
			definitions: []string{
				`seq(x___) := "null sequence"`,
				`seq(x__) := "sequence"`,
				`seq(x_) := "single"`,
			},
			testInput:   "seq(1, 2)",
			expected:    `"sequence"`,
			description: "Sequence pattern for multiple args",
		},
		{
			definitions: []string{
				`seq(x___) := "null sequence"`,
				`seq(x__) := "sequence"`,
				`seq(x_) := "single"`,
			},
			testInput:   "seq()",
			expected:    `"null sequence"`,
			description: "Null sequence pattern for zero args",
		},
		{
			definitions: []string{
				`func(x_, y_) := "two general"`,
				`func(x_Integer, y_Integer) := "two integers"`,
				`func(42, y_Integer) := "literal and integer"`,
			},
			testInput:   "func(42, 99)",
			expected:    `"literal and integer"`,
			description: "Mixed literal and typed should be most specific",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Create fresh evaluator for each test
			evaluator := NewEvaluator()

			// Define all patterns
			for _, def := range tc.definitions {
				result := evaluateStringHelper(t, evaluator, def)
				// Definitions should return the symbol being defined
				if !contains(result, "test") && !contains(result, "seq") && !contains(result, "func") {
					t.Logf("Definition result: %s", result)
				}
			}

			// Test the input
			result := evaluateStringHelper(t, evaluator, tc.testInput)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s for input %s", tc.expected, result, tc.testInput)
			}
		})
	}
}

// TestPatternSpecificityEdgeCases tests edge cases in pattern specificity
func TestPatternSpecificityEdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		moreSpecific string
		lessSpecific string
	}{
		// Empty vs non-empty patterns (current algorithm favors more args)
		{"One arg vs zero args", "f(x_)", "f()"},

		// Complex nested patterns
		{"Nested literal", "f(Plus(42))", "f(Plus(x_))"},
		{"Nested typed", "f(Plus(x_Integer))", "f(Plus(x_))"},

		// Association patterns
		{"Typed association", "f(x_Association)", "f(x_)"},

		// Multiple type constraints
		{"Multiple specific types", "f(x_Integer, y_String)", "f(x_, y_)"},
		{"Mixed specificity", "f(x_Integer, y_)", "f(x_, y_)"},

		// Rule patterns
		{"Typed rule", "f(x_Rule)", "f(x_)"},

		// Complex sequence patterns
		{"Typed sequence", "f(x__Integer)", "f(x__)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moreSpecific, err1 := ParseString(tt.moreSpecific)
			if err1 != nil {
				t.Fatalf("Failed to parse more specific pattern %q: %v", tt.moreSpecific, err1)
			}

			lessSpecific, err2 := ParseString(tt.lessSpecific)
			if err2 != nil {
				t.Fatalf("Failed to parse less specific pattern %q: %v", tt.lessSpecific, err2)
			}

			moreSpec := core.GetPatternSpecificity(moreSpecific)
			lessSpec := core.GetPatternSpecificity(lessSpecific)

			if moreSpec <= lessSpec {
				t.Errorf("Pattern %q (specificity=%d) should be MORE specific than %q (specificity=%d)",
					tt.moreSpecific, moreSpec, tt.lessSpecific, lessSpec)
			}
		})
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
