package integration

import (
	"testing"
)

func TestFlatten(t *testing.T) {

	runTestCases(t, []TestCase{
		// Basic flattening functionality
		{
			name:     "Simple nested list",
			input:    "Flatten(List(1, 2, List(3, 4)))",
			expected: "List(1, 2, 3, 4)",
		},
		{
			name:     "Deeply nested list",
			input:    "Flatten(List(1, List(2, List(3, 4), 5), 6))",
			expected: "List(1, 2, 3, 4, 5, 6)",
		},
		{
			name:     "Multiple nested lists at same level",
			input:    "Flatten(List(List(1, 2), List(3, 4), List(5, 6)))",
			expected: "List(1, 2, 3, 4, 5, 6)",
		},
		{
			name:     "Empty nested lists",
			input:    "Flatten(List(1, List(), 2, List(), 3))",
			expected: "List(1, 2, 3)",
		},

		// Different head types
		{
			name:     "Zoo with nested Zoo",
			input:    "Flatten(Zoo(a, b, Zoo(c, d)))",
			expected: "Zoo(a, b, c, d)",
		},
		{
			name:     "Plus with nested Plus (using Hold)",
			input:    "Flatten(Hold(Plus(1, Plus(2, 3), 4)))",
			expected: "Hold(Plus(1, 2, 3, 4))",
		},
		{
			name:     "Times with nested Times",
			input:    "Flatten(Times(x, Times(y, z)))",
			expected: "Times(x, y, z)",
		},

		// Mixed head types (should not flatten across different heads)
		{
			name:     "List with Zoo inside - no cross-head flattening",
			input:    "Flatten(List(1, Zoo(2, 3), 4))",
			expected: "List(1, Zoo(2, 3), 4)",
		},
		{
			name:     "Zoo with List inside - no cross-head flattening",
			input:    "Flatten(Zoo(a, List(b, c), d))",
			expected: "Zoo(a, List(b, c), d)",
		},

		// Edge cases
		{
			name:     "Empty list",
			input:    "Flatten(List())",
			expected: "List()",
		},
		{
			name:     "Single element list",
			input:    "Flatten(List(42))",
			expected: "List(42)",
		},
		{
			name:     "Already flat list",
			input:    "Flatten(List(1, 2, 3, 4))",
			expected: "List(1, 2, 3, 4)",
		},
		{
			name:     "Non-list input",
			input:    "Flatten(42)",
			expected: "Flatten(42)",
		},
		{
			name:     "Symbol input",
			input:    "Flatten(x)",
			expected: "Flatten(x)",
		},

		// Complex nested structures
		{
			name:     "Triple nesting",
			input:    "Flatten(List(1, List(2, List(3, List(4, 5)))))",
			expected: "List(1, 2, 3, 4, 5)",
		},
		{
			name:     "Mixed data types",
			input:    "Flatten(List(a, List(1, List(\"hello\", True)), b))",
			expected: "List(a, 1, \"hello\", True, b)",
		},

		// Interaction with other functions
		{
			name:     "Flatten result of Length",
			input:    "Flatten(List(Length(List(1,2,3)), List(4, 5)))",
			expected: "List(3, 4, 5)",
		},
		{
			name:     "Length of flattened list",
			input:    "Length(Flatten(List(List(1, 2), List(3, 4, 5))))",
			expected: "5",
		},
		{
			name:     "First of flattened list",
			input:    "First(Flatten(List(List(a, b), List(c, d))))",
			expected: "a",
		},
		{
			name:     "Last of flattened list",
			input:    "Last(Flatten(List(List(a, b), List(c, d))))",
			expected: "d",
		},

		// Nested with arithmetic
		{
			name:     "Flatten with evaluated arithmetic",
			input:    "Flatten(List(Plus(1, 2), List(Times(2, 3), 4)))",
			expected: "List(3, 6, 4)",
		},

		// Very deep nesting
		{
			name:     "Five levels deep",
			input:    "Flatten(List(1, List(2, List(3, List(4, List(5))))))",
			expected: "List(1, 2, 3, 4, 5)",
		},
		{
			name:     "Flatten with no arguments",
			input:    "Flatten()",
			expected: "Flatten()",
		},
		{
			name:     "Flatten with too many arguments",
			input:    "Flatten(List(1, 2), List(3, 4))",
			expected: "Flatten(List(1, 2), List(3, 4))",
		},
		// Test flattening preserves mathematical structure correctly (using Hold to prevent evaluation)
		{
			name:     "Flatten Plus expression with Hold",
			input:    "Flatten(Hold(Plus(1, Plus(2, 3))))",
			expected: "Hold(Plus(1, 2, 3))",
		},
		{
			name:     "Flatten Times expression",
			input:    "Flatten(Times(a, Times(b, c)))",
			expected: "Times(a, b, c)",
		},
		{
			name:     "Mixed Plus and Times with Hold - no cross-flattening",
			input:    "Flatten(Hold(Plus(1, Times(2, 3))))",
			expected: "Hold(Plus(1, Times(2, 3)))",
		},
		{
			name:     "Nested Plus with Hold",
			input:    "Flatten(Hold(Plus(Plus(1, 2), Plus(3, 4))))",
			expected: "Hold(Plus(1, 2, 3, 4))",
		},
	})

}
