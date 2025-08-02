package integration

import (
	"testing"
)

func TestRuleDelayed_Basic(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Basic RuleDelayed creation",
			input:    "RuleDelayed(x, x + 1)",
			expected: "RuleDelayed(x, Plus(x, 1))",
		},
		{
			name:     "RuleDelayed with pattern variable (=> syntax)",
			input:    "x_ => x + 1",
			expected: "RuleDelayed(Pattern(x, Blank()), Plus(x, 1))",
		},
		{
			name:     "RuleDelayed with pattern variable (function form)",
			input:    "RuleDelayed(x_, x + 1)",
			expected: "RuleDelayed(Pattern(x, Blank()), Plus(x, 1))",
		},
		{
			name:     "RuleDelayed RHS held unevaluated",
			input:    "RuleDelayed(x_, 1 + 2)",
			expected: "RuleDelayed(Pattern(x, Blank()), Plus(1, 2))",
		},
	}

	runTestCases(t, tests)
}

func TestRuleDelayed_Replace(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Replace with RuleDelayed - simple case (=> syntax)",
			input:    "Replace(5, x_ => x + 1)",
			expected: "6",
		},
		{
			name:     "Replace with RuleDelayed - simple case (function form)",
			input:    "Replace(5, RuleDelayed(x_, x + 1))",
			expected: "6",
		},
		{
			name:     "Replace with RuleDelayed - no match",
			input:    "Replace(5, RuleDelayed(6, x + 1))",
			expected: "5",
		},
		{
			name:     "Replace with RuleDelayed - exact match",
			input:    "Replace(5, RuleDelayed(5, 42))",
			expected: "42",
		},
		{
			name:     "Replace with RuleDelayed - list pattern",
			input:    "Replace([1, 2, 3], RuleDelayed(List(x_, y_, z_), x + y + z))",
			expected: "6",
		},
	}

	runTestCases(t, tests)
}

func TestRuleDelayed_LexicalScoping(t *testing.T) {
	tests := []TestCase{
		{
			name:        "Local bindings don't pollute global scope",
			input:       "x = 100; Replace(5, RuleDelayed(y_, x + y))",
			expected:    "105",
		},
		{
			name:        "Pattern variables are locally scoped",
			input:       "y = 999; Replace(42, RuleDelayed(y_, y * 2))",
			expected:    "84",
		},
		{
			name:        "Complex expression with local evaluation",
			input:       "f(x_) := x^2; Replace(3, RuleDelayed(x_, f(x) + x))",
			expected:    "12",
		},
	}

	runTestCases(t, tests)
}

func TestRuleDelayed_ReplaceAll(t *testing.T) {
	tests := []TestCase{
		{
			name:     "ReplaceAll with RuleDelayed - all occurrences",
			input:    "ReplaceAll(Plus(x, x, y), RuleDelayed(x, 2))",
			expected: "Plus(4, y)",
		},
		{
			name:     "ReplaceAll with RuleDelayed - nested expressions",
			input:    "ReplaceAll(Plus(f(x), g(x)), RuleDelayed(x, 42))",
			expected: "Plus(f(42), g(42))",
		},
		{
			name:     "ReplaceAll with RuleDelayed - pattern variables",
			//input:    "ReplaceAll(6, RuleDelayed(x_, x * 10))",
			input:    "ReplaceAll(Plus(1, 2, 3), RuleDelayed(x_, x * 10))",
			expected: "60",
			
		},
	}

	runTestCases(t, tests)
}
