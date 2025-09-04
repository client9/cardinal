package integration

import (
	"testing"
)

func TestSort(t *testing.T) {
	tests := []TestCase{
		// Basic string sorting (lexicographic with same length)
		{
			name:     "Simple symbols - lexicographic order",
			input:    "Sort(List(z, a, b))",
			expected: "List(a, b, z)",
		},
		{
			name:     "Mixed case symbols",
			input:    "Sort(List(Z, a, B))",
			expected: "List(B, Z, a)", // Capital letters come before lowercase in ASCII
		},

		// Length-based sorting (shorter expressions first)
		{
			name:     "Sort by length - symbols",
			input:    "Sort(List(xx, a, zzz))",
			expected: "List(a, xx, zzz)", // Length 1, 2, 3
		},
		{
			name:     "Sort by length - mixed expressions",
			input:    "Sort(List(Foo(1, 2), x, Bar(1,2,3)))",
			expected: "List(x, Foo(1, 2), Bar(1, 2, 3))",
		},

		// Numbers vs symbols (numbers are shorter in string representation)
		{
			name:     "Numbers vs symbols",
			input:    "Sort(List(x, 1, 10))",
			expected: "List(1, 10, x)", // Numbers sort before symbols lexicographically
		},
		{
			name:     "Mixed types - integers, reals, symbols",
			input:    "Sort(List(a, 2.5, 1))",
			expected: "List(1, 2.5, a)", // Numeric values first, then symbols
		},

		// String literals
		{
			name:     "String literals",
			input:    "Sort(List(\"hello\", \"abc\", x))",
			expected: "List(x, \"abc\", \"hello\")", // Symbol comes first, then quoted strings sorted
		},

		// Complex nested expressions
		{
			name:     "Nested expressions with different heads",
			input:    "Sort(Plus(Foo(a,b), Bar(1,2), x))",
			expected: "Plus(x, Bar(1, 2), Foo(a, b))", // By length: 0, 2, 2, then lexicographic
		},

		// Edge cases
		{
			name:     "Empty list",
			input:    "Sort(List())",
			expected: "List()",
		},
		{
			name:     "Single element",
			input:    "Sort(List(x))",
			expected: "List(x)",
		},
		{
			name:     "Two elements - already sorted",
			input:    "Sort(List(a, b))",
			expected: "List(a, b)",
		},

		// Non-list inputs (should return unchanged)
		{
			name:     "Integer input",
			input:    "Sort(42)",
			expected: "42",
		},
		{
			name:     "Symbol input",
			input:    "Sort(x)",
			expected: "x",
		},

		// Different heads (Sort should work with any head)
		{
			name:     "Custom head",
			input:    "Sort(MyFunction(gamma, alpha, beta))",
			expected: "MyFunction(alpha, beta, gamma)",
		},
	}
	runTestCases(t, tests)
}
