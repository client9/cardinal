package integration

import (
	"testing"
)

func TestReplaceAllFunction(t *testing.T) {
	tests := []TestCase{
		{
			name:     "ReplaceAll with single Rule - simple replacement",
			input:    `ReplaceAll(x, Rule(x, y))`,
			expected: `y`,
		},
		{
			name:     "ReplaceAll with single Rule - no match",
			input:    `ReplaceAll(x, Rule(z, y))`,
			expected: `x`,
		},
		{
			name:     "ReplaceAll with single Rule - nested expression",
			input:    `ReplaceAll(Plus(x, x), Rule(x, 2))`,
			expected: `4`, // Plus(2, 2) evaluates to 4
		},
		{
			name:     "ReplaceAll with single Rule - deeply nested",
			input:    `ReplaceAll(Plus(x, Times(x, y)), Rule(x, 1))`,
			expected: `Plus(1, y)`, // Times(1, y) simplifies to y, then Plus(1, y)
		},
		{
			name:     "ReplaceAll with single Rule - multiple levels",
			input:    `ReplaceAll(Plus(a, Plus(a, Times(a, b))), Rule(a, z))`,
			expected: `Plus(z, z, Times(b, z))`, // Plus flattens: Plus(z, Plus(z, Times(z, b))) -> Plus(z, z, Times(z, b)) -> Plus(Times(b, z), z, z) (sorted)
		},
		{
			name:     "ReplaceAll with power expressions",
			input:    `ReplaceAll(Power(x, Plus(x, 1)), Rule(x, 2))`,
			expected: `8`, // Power(2, Plus(2, 1)) = Power(2, 3) = 8.0
		},
		{
			name:     "ReplaceAll with colon syntax",
			input:    `ReplaceAll(Plus(a, Times(a, b)), a : 3)`,
			expected: `Plus(3, Times(3, b))`, // Plus is Orderless: Plus(Times(3, b), 3) -> Plus(3, Times(3, b))
		},
		{
			name:     "ReplaceAll stops at first match level",
			input:    `ReplaceAll(Plus(x, y), Rule(Plus(x, y), result))`,
			expected: `result`,
		},
	}

	runTestCases(t, tests)
}

func TestReplaceAllWithRules(t *testing.T) {
	tests := []TestCase{
		{
			name:     "ReplaceAll with List of Rules - both rules apply recursively",
			input:    `ReplaceAll(Plus(x, Times(x, y)), List(x : 1, y : 2))`,
			expected: `3`, // Plus(1, Times(1, 2)) = Plus(1, 2) = 3
		},
		{
			name:     "ReplaceAll with List of Rules - both rules apply recursively variant",
			input:    `ReplaceAll(Plus(y, Times(x, y)), List(x : 1, y : 2))`,
			expected: `4`, // Plus(2, Times(1, 2)) = Plus(2, 2) = 4
		},
		{
			name:     "ReplaceAll with List of Rules - both rules apply to different subexpressions",
			input:    `ReplaceAll(Plus(x, y), List(x : 1, y : 2))`,
			expected: `3`, // Plus(1, 2) = 3
		},
		{
			name:     "ReplaceAll with List of Rules - nested multiple replacements",
			input:    `ReplaceAll(Plus(a, Times(b, Plus(a, b))), List(a : 1, b : 2))`,
			expected: `7`, // Plus(1, Times(2, Plus(1, 2))) = Plus(1, Times(2, 3)) = Plus(1, 6) = 7
		},
		{
			name:     "ReplaceAll with List of Rules - no matches",
			input:    `ReplaceAll(Plus(z, w), List(x : 1, y : 2))`,
			expected: `Plus(w, z)`, // Plus is Orderless, so z,w gets sorted to w,z
		},
		{
			name:     "ReplaceAll with List of Rules - first rule wins at same level",
			input:    `ReplaceAll(x, List(x : first, x : second))`,
			expected: `first`,
		},
		{
			name:     "ReplaceAll with List of Rules - complex nested structure",
			input:    `ReplaceAll(Times(Plus(x, y), Power(x, 2)), List(x : a, y : b))`,
			expected: `Times(Plus(a, b), Power(a, 2))`,
		},
		{
			name:     "ReplaceAll with empty List",
			input:    `ReplaceAll(Plus(x, y), List())`,
			expected: `Plus(x, y)`,
		},
		{
			name:     "ReplaceAll with List containing non-Rules (pattern should not match)",
			input:    `ReplaceAll(x, List(x : a, 42, y : b))`,
			expected: "",
			errorType: "ArgumentError",
		},
		{
			name:     "ReplaceAll with power and arithmetic expressions",
			input:    `ReplaceAll(Plus(Power(x, 2), Times(2, x)), List(x : 3))`,
			expected: `15`, // Plus(Power(3, 2), Times(2, 3)) = Plus(9.0, 6) = 15.0
		},
		{
			name:     "ReplaceAll with recursive function application",
			input:    `ReplaceAll(Plus(f(x), f(y)), List(f(z_) : Times(2, z)))`,
			expected: `Plus(Times(2, x), Times(2, y))`, // Pattern matching works: f(x) -> Times(2, x), f(y) -> Times(2, y)
		},
		{
			name:     "ReplaceAll with deeply nested Lists",
			input:    `ReplaceAll(List(x, List(y, x)), List(x : 1, y : 2))`,
			expected: `List(1, List(2, 1))`,
		},
	}

	runTestCases(t, tests)
}

