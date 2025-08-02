package integration

import (
	"testing"
)

func TestArithmeticSimplification(t *testing.T) {
	runTestCases(t, []TestCase{
		// Basic Plus simplification
		{
			name:     "Plus with all integers",
			input:    "Plus(1, 2, 3, 4)",
			expected: "10",
		},
		{
			name:     "Plus with integers and symbols",
			input:    "Plus(1, 2, 3, x)",
			expected: "Plus(6, x)",
		},
		{
			name:     "Plus with mixed positions",
			input:    "Plus(x, y, z, 10, 2)",
			expected: "Plus(12, x, y, z)",
		},
		{
			name:     "Plus with mixed types - integer and real",
			input:    "Plus(x, y, z, 2.0, 1)",
			expected: "Plus(3.0, x, y, z)",
		},
		{
			name:     "Plus with all reals",
			input:    "Plus(1.5, 2.5, 3.0)",
			expected: "7.0",
		},
		{
			name:     "Plus with non-numeric types",
			input:    "Plus(True, 1, 2)",
			expected: "Plus(3, True)",
		},
		{
			name:     "Plus with string",
			input:    "Plus(\"foo\", 1, 2)",
			expected: "Plus(3, \"foo\")",
		},
		{
			name:     "Plus with nested expression",
			input:    "Plus(Times(x, y), 1, 2)",
			expected: "Plus(3, Times(x, y))",
		},

		// Edge cases for Plus
		{
			name:     "Plus with zero",
			input:    "Plus(0, x, y)",
			expected: "Plus(x, y)", // Zero should be omitted when there are other terms
		},
		{
			name:     "Plus with only zero and symbols",
			input:    "Plus(0, x)",
			expected: "x", // OneIdentity behavior: Plus with one non-zero element returns that element
		},
		{
			name:     "Plus empty",
			input:    "Plus()",
			expected: "0",
		},
		{
			name:     "Plus single element",
			input:    "Plus(42)",
			expected: "42", // OneIdentity behavior
		},
		{
			name:     "Plus only symbols",
			input:    "Plus(x, y, z)",
			expected: "Plus(x, y, z)",
		},

		// Basic Times simplification
		{
			name:     "Times with all integers",
			input:    "Times(2, 3, 4)",
			expected: "24",
		},
		{
			name:     "Times with integers and symbols",
			input:    "Times(2, 3, 4, x)",
			expected: "Times(24, x)",
		},
		{
			name:     "Times with mixed types",
			input:    "Times(2, 3.0, x)",
			expected: "Times(6.0, x)",
		},
		{
			name:     "Times with zero",
			input:    "Times(0, x, y)",
			expected: "0", // Zero short-circuit
		},
		{
			name:     "Times with one",
			input:    "Times(1, x, y)",
			expected: "Times(x, y)", // One should be omitted
		},
		{
			name:     "Times empty",
			input:    "Times()",
			expected: "1",
		},
		{
			name:     "Times single element",
			input:    "Times(42)",
			expected: "42", // OneIdentity behavior
		},

		// Complex mixed cases
		{
			name:     "Nested arithmetic",
			input:    "Plus(Times(2, 3), 4, x)",
			expected: "Plus(10, x)", // Times(2,3) → 6, then Plus(6,4,x) → Plus(10,x)
		},
		{
			name:     "Multiple nested expressions",
			input:    "Times(Plus(1, 2), 4, x)",
			expected: "Times(12, x)", // Plus(1,2) → 3, then Times(3,4,x) → Times(12,x)
		},

		// Type promotion edge cases
		{
			name:     "Plus integer overflow behavior",
			input:    "Plus(9223372036854775807, 1)", // max int64 + 1, wraps to min int64
			expected: "-9223372036854775808",
		},
		{
			name:     "Mixed precision",
			input:    "Plus(1, 2.5, 3)",
			expected: "6.5",
		},
	})
}

func TestArithmeticSimplificationWithReplace(t *testing.T) {
	runTestCases(t, []TestCase{
		// Test that replacement rules work correctly with the new arithmetic
		{
			name:     "Replace on simplified Plus",
			input:    "Replace(Plus(1, 2, x), Rule(Plus(n_Integer, x_), transformed(n, x)))",
			expected: "transformed(3, x)",
		},
		{
			name:     "Replace on simplified Times",
			input:    "Replace(Times(2, 3, x), Rule(Times(n_Integer, x_), scaled(n, x)))",
			expected: "scaled(6, x)",
		},
		{
			name:     "ReplaceAll with nested arithmetic",
			input:    "ReplaceAll(Plus(Times(2, 3), 4, x), Rule(Times(a_Integer, b_Integer), Power(a, b)))",
			expected: "Plus(10, x)", // Times(2,3) already simplified to 6, so rule doesn't match
		},
		{
			name:     "Pattern matching on mixed types",
			input:    "MatchQ(Plus(1, 2.0, x), Plus(n_Real, x_))",
			expected: "True", // 1+2.0 becomes 3.0 (Real), so pattern matches
		},
	})

}
