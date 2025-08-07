package stdlib

import (
	"testing"

	"github.com/client9/sexpr/core"
)

func TestSort(t *testing.T) {
	// Helper function to create test lists with mixed elements
	createMixedList := func(head string, elements ...core.Expr) core.List {
		exprs := make([]core.Expr, len(elements)+1)
		exprs[0] = core.NewSymbol(head)
		copy(exprs[1:], elements)
		return core.NewListFromExprs(exprs...)
	}

	tests := []struct {
		name     string
		input    core.Expr
		expected string
	}{
		// Basic string sorting (lexicographic with same length)
		{
			name:     "Simple symbols - lexicographic order",
			input:    createMixedList("List", core.NewSymbol("z"), core.NewSymbol("a"), core.NewSymbol("b")),
			expected: "List(a, b, z)",
		},
		{
			name:     "Mixed case symbols",
			input:    createMixedList("List", core.NewSymbol("Z"), core.NewSymbol("a"), core.NewSymbol("B")),
			expected: "List(B, Z, a)", // Capital letters come before lowercase in ASCII
		},

		// Length-based sorting (shorter expressions first)
		{
			name:     "Sort by length - symbols",
			input:    createMixedList("List", core.NewSymbol("xx"), core.NewSymbol("a"), core.NewSymbol("zzz")),
			expected: "List(a, xx, zzz)", // Length 1, 2, 3
		},
		{
			name: "Sort by length - mixed expressions",
			input: createMixedList("List",
				createMixedList("Plus", core.NewInteger(1), core.NewInteger(2)), // Length 2 (Plus has 2 args)
				core.NewSymbol("x"), // Length 0 (atom)
				createMixedList("Times", core.NewInteger(1), core.NewInteger(2), core.NewInteger(3)), // Length 3
			),
			expected: "List(x, Plus(1, 2), Times(1, 2, 3))",
		},

		// Numbers vs symbols (numbers are shorter in string representation)
		{
			name:     "Numbers vs symbols",
			input:    createMixedList("List", core.NewSymbol("x"), core.NewInteger(1), core.NewInteger(10)),
			expected: "List(1, 10, x)", // Numbers sort before symbols lexicographically
		},
		{
			name:     "Mixed types - integers, reals, symbols",
			input:    createMixedList("List", core.NewSymbol("a"), core.NewReal(2.5), core.NewInteger(1)),
			expected: "List(1, 2.5, a)", // Numeric values first, then symbols
		},

		// String literals
		{
			name:     "String literals",
			input:    createMixedList("List", core.NewString("hello"), core.NewString("abc"), core.NewSymbol("x")),
			expected: "List(x, \"abc\", \"hello\")", // Symbol comes first, then quoted strings sorted
		},

		// Complex nested expressions
		{
			name: "Nested expressions with different heads",
			input: createMixedList("Plus",
				createMixedList("Times", core.NewSymbol("a"), core.NewSymbol("b")),
				createMixedList("Plus", core.NewInteger(1), core.NewInteger(2)),
				core.NewSymbol("x"),
			),
			expected: "Plus(x, Plus(1, 2), Times(a, b))", // By length: 0, 2, 2, then lexicographic
		},

		// Edge cases
		{
			name:     "Empty list",
			input:    createMixedList("List"),
			expected: "List()",
		},
		{
			name:     "Single element",
			input:    createMixedList("List", core.NewSymbol("x")),
			expected: "List(x)",
		},
		{
			name:     "Two elements - already sorted",
			input:    createMixedList("List", core.NewSymbol("a"), core.NewSymbol("b")),
			expected: "List(a, b)",
		},
		{
			name:     "Two elements - reverse order",
			input:    createMixedList("List", core.NewSymbol("b"), core.NewSymbol("a")),
			expected: "List(a, b)",
		},

		// Non-list inputs (should return unchanged)
		{
			name:     "Integer input",
			input:    core.NewInteger(42),
			expected: "42",
		},
		{
			name:     "Symbol input",
			input:    core.NewSymbol("x"),
			expected: "x",
		},

		// Different heads (Sort should work with any head)
		{
			name:     "Plus expression",
			input:    createMixedList("Plus", core.NewSymbol("z"), core.NewSymbol("a"), core.NewInteger(1)),
			expected: "Plus(1, a, z)",
		},
		{
			name:     "Times expression",
			input:    createMixedList("Times", core.NewSymbol("c"), core.NewSymbol("a"), core.NewSymbol("b")),
			expected: "Times(a, b, c)",
		},
		{
			name:     "Custom head",
			input:    createMixedList("MyFunction", core.NewSymbol("gamma"), core.NewSymbol("alpha"), core.NewSymbol("beta")),
			expected: "MyFunction(alpha, beta, gamma)",
		},

		// Test the exact same ordering as Orderless attribute
		{
			name: "Orderless-style complex sorting",
			input: createMixedList("Plus",
				core.NewSymbol("zzz"), // Length 0, string "zzz"
				core.NewSymbol("a"),   // Length 0, string "a"
				core.NewInteger(10),   // Length 0, string "10"
				core.NewInteger(2),    // Length 0, string "2"
				createMixedList("Times", core.NewSymbol("x"), core.NewSymbol("y")), // Length 2
				createMixedList("Power", core.NewSymbol("x"), core.NewInteger(2)),  // Length 2
			),
			expected: "Plus(10, 2, a, zzz, Power(x, 2), Times(x, y))", // Length 0 items first (sorted), then length 2 items (sorted)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sort(tt.input)
			if result.String() != tt.expected {
				t.Errorf("Sort(%s) = %s, expected %s",
					tt.input.String(), result.String(), tt.expected)
			}
		})
	}
}
