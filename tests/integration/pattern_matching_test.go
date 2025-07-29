package integration

import (
	"testing"
)

func TestMatchQSequencePatterns(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic sequence patterns with BlankNullSequence (z___)
		{
			name:     "BlankNullSequence with zero elements",
			input:    "MatchQ(Zoo(1,2), Zoo(x_Integer, y_Integer, z___))",
			expected: "True",
		},
		{
			name:     "BlankNullSequence with one element",
			input:    "MatchQ(Zoo(1,2,a), Zoo(x_Integer, y_Integer, z___))",
			expected: "True",
		},
		{
			name:     "BlankNullSequence with multiple elements",
			input:    "MatchQ(Zoo(1,2,a,b), Zoo(x_Integer, y_Integer, z___))",
			expected: "True",
		},
		{
			name:     "BlankNullSequence with many elements",
			input:    "MatchQ(Zoo(1,2,a,b,c,d), Zoo(x_Integer, y_Integer, z___))",
			expected: "True",
		},

		// Type constraints on sequence patterns
		{
			name:     "BlankNullSequence with type constraint - success",
			input:    "MatchQ(Zoo(1,2,3,4), Zoo(x_Integer, y_Integer, z___Integer))",
			expected: "True",
		},
		{
			name:     "BlankNullSequence with type constraint - failure",
			input:    "MatchQ(Zoo(1,2,a,b), Zoo(x_Integer, y_Integer, z___Integer))",
			expected: "False",
		},
		{
			name:     "BlankNullSequence with mixed types - partial match",
			input:    "MatchQ(Zoo(1,2,3,a), Zoo(x_Integer, y_Integer, z___Integer))",
			expected: "False",
		},

		// BlankSequence patterns (z__) - must match at least one element
		{
			name:     "BlankSequence with zero elements - should fail",
			input:    "MatchQ(Zoo(1,2), Zoo(x_Integer, y_Integer, z__))",
			expected: "False",
		},
		{
			name:     "BlankSequence with one element",
			input:    "MatchQ(Zoo(1,2,a), Zoo(x_Integer, y_Integer, z__))",
			expected: "True",
		},
		{
			name:     "BlankSequence with multiple elements",
			input:    "MatchQ(Zoo(1,2,a,b), Zoo(x_Integer, y_Integer, z__))",
			expected: "True",
		},

		// Multiple sequence patterns
		{
			name:     "Multiple BlankNullSequence patterns",
			input:    "MatchQ(Zoo(1,a,b,2), Zoo(x_Integer, y___, z_Integer))",
			expected: "True",
		},
		{
			name:     "Complex pattern with multiple sequences",
			input:    "MatchQ(Zoo(1,a,b,c,2,d,e), Zoo(x_Integer, y___, z_Integer, w___))",
			expected: "True",
		},

		// Edge cases
		{
			name:     "Pattern longer than expression",
			input:    "MatchQ(Zoo(1), Zoo(x_Integer, y_Integer, z___))",
			expected: "False",
		},
		{
			name:     "Empty Zoo with sequence pattern",
			input:    "MatchQ(Zoo(), Zoo(z___))",
			expected: "True",
		},
	}

	runTestCases(t, tests)
}

func TestMatchQBasicPatterns(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic wildcard patterns
		{
			name:     "Simple wildcard match",
			input:    "MatchQ(Zoo(1,2), Zoo(x_,y_))",
			expected: "True",
		},
		{
			name:     "Type-constrained wildcard match",
			input:    "MatchQ(Zoo(1,2), Zoo(x_Integer, y_Integer))",
			expected: "True",
		},
		{
			name:     "Type-constrained wildcard failure",
			input:    "MatchQ(Zoo(1,a), Zoo(x_Integer, y_Integer))",
			expected: "False",
		},

		// Literal matches
		{
			name:     "Exact literal match",
			input:    "MatchQ(Zoo(1,2), Zoo(1,2))",
			expected: "True",
		},
		{
			name:     "Literal mismatch",
			input:    "MatchQ(Zoo(1,2), Zoo(1,3))",
			expected: "False",
		},

		// Mixed patterns
		{
			name:     "Mixed literal and wildcard",
			input:    "MatchQ(Zoo(1,a), Zoo(1,x_))",
			expected: "True",
		},
		{
			name:     "Mixed with type constraint",
			input:    "MatchQ(Zoo(1,42), Zoo(1,x_Integer))",
			expected: "True",
		},
	}

	runTestCases(t, tests)
}

func TestMatchQRegressionCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// These were the original failing cases
		{
			name:     "Original failing case - sequence with one extra",
			input:    "MatchQ(Zoo(1,2,a), Zoo(x_Integer, y_Integer, z___))",
			expected: "True",
		},
		{
			name:     "Original failing case - sequence with two extras",
			input:    "MatchQ(Zoo(1,2,a,b), Zoo(x_Integer, y_Integer, z___))",
			expected: "True",
		},

		// Plus pattern (should work with Hold)
		{
			name:     "Plus pattern with Hold",
			input:    "MatchQ(Hold(Plus(1,3)), Hold(Plus(x_,y_)))",
			expected: "True",
		},
	}

	runTestCases(t, tests)
}
