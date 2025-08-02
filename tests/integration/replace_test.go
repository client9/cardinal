package integration

import (
	"testing"
)

func TestReplaceFunction(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Replace with exact match using Rule",
			input:    `Replace(x^2, Rule(x^2, a + b))`,
			expected: `Plus(a, b)`,
		},
		{
			name:     "Replace with exact match using colon syntax",
			input:    `Replace(x^2, x^2 : a + b)`,
			expected: `Plus(a, b)`,
		},
		{
			name:     "Replace with Power function syntax",
			input:    `Replace(Power(x, 2), Rule(Power(x, 2), Plus(a, b)))`,
			expected: `Plus(a, b)`,
		},
		{
			name:     "Replace with no match",
			input:    `Replace(x^3, Rule(x^2, a + b))`,
			expected: `Power(x, 3)`,
		},
		{
			name:     "Replace simple expression",
			input:    `Replace(3, Rule(3, Times(x, y)))`,
			expected: `Times(x, y)`,
		},
	}
	
	runTestCases(t, tests)
}

func TestReplaceWithRules(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Replace with List of Rules - first rule matches",
			input:    `Replace(x, List(x : a, y : b))`,
			expected: `a`,
		},
		{
			name:     "Replace with List of Rules - second rule matches",
			input:    `Replace(y, List(x : a, y : b))`,
			expected: `b`,
		},
		{
			name:     "Replace with List of Rules - no matches",
			input:    `Replace(z, List(x : a, y : b))`,
			expected: `z`,
		},
		{
			name:     "Replace with List of Rules - first rule wins",
			input:    `Replace(x, List(x : first, x : second))`,
			expected: `first`,
		},
		{
			name:     "Replace with List of Rules - complex expressions",
			input:    `Replace(Plus(1, 2), List(Plus(1, 2) : Times(a, b), Plus(2, 3) : Times(c, d)))`,
			expected: `Times(a, b)`,
		},
		{
			name:     "Replace with List of Rules - power expressions",
			input:    `Replace(x^2, List(x^3 : cube, x^2 : square, x : linear))`,
			expected: `square`,
		},
		{
			name:     "Replace with List of Rules using Rule function",
			input:    `Replace(a, List(Rule(a, first), Rule(b, second)))`,
			expected: `first`,
		},
		{
			name:     "Replace with empty List",
			input:    `Replace(x, List())`,
			expected: `x`,
		},
		{
			name:     "Replace with List containing non-Rules (pattern should not match)",
			input:    `Replace(x, List(x : a, 42, y : b))`,
			expected: "",
			errorType: "ArgumentError",
		},
		{
			name:     "Replace with nested expressions",
			input:    `Replace(Times(x, y), List(Times(x, y) : result1, Plus(x, y) : result2))`,
			expected: `result1`,
		},
	}

	runTestCases(t, tests)
}
