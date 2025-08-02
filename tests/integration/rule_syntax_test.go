package integration

import (
	"testing"
)

func TestRuleShorthandSyntax(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Simple Rule shorthand",
			input:    "a:b",
			expected: "Rule(a, b)",
		},
		{
			name:     "Rule shorthand with numbers",
			input:    "1:2",
			expected: "Rule(1, 2)",
		},
		{
			name:     "Rule shorthand with expressions",
			input:    "Plus(a, b):Times(x, y)",
			expected: "Rule(Plus(a, b), Times(x, y))",
		},
		{
			name:     "Complex pattern Rule shorthand",
			input:    "(x_^a_ * x_^b_) : x^(a+b)",
			expected: "Rule(Times(Power(Pattern(x, Blank()), Pattern(a, Blank())), Power(Pattern(x, Blank()), Pattern(b, Blank()))), Power(x, Plus(a, b)))",
		},
	}

	runTestCases(t, tests)
}

func TestRuleShorthandWithAssignment(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Simple Rule assignment",
			input:    "r1 = a:b; FullForm(r1)",
			expected: "\"Rule(a, b)\"",
		},
		{
			name:     "Rule assignment with numbers",
			input:    "r2 = 1:2; FullForm(r2)",
			expected: "\"Rule(1, 2)\"",
		},
		{
			name:     "Rule assignment with expressions",
			input:    "r3 = Plus(a, b):Times(x, y); FullForm(r3)",
			expected: "\"Rule(Plus(a, b), Times(x, y))\"",
		},
		{
			name:     "Complex pattern Rule assignment - BUG CASE",
			input:    "r4 = (x_^a_ * x_^b_) : x^(a+b); FullForm(r4)",
			expected: "\"Rule(Times(Power(Pattern(x, Blank()), Pattern(a, Blank())), Power(Pattern(x, Blank()), Pattern(b, Blank()))), Power(x, Plus(a, b)))\"",
		},
	}

	runTestCases(t, tests)
}

func TestRuleShorthandVsExplicit(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Simple rule comparison",
			input:    "SameQ( a:b, Rule(a, b))",
			expected: "True",
		},
		{
			name:     "Complex pattern rule comparison",
			input:    "SameQ( (x_^a_ * x_^b_) : x^(a+b), Rule(x_^a_ * x_^b_, x^(a+b)))",
			expected: "True",
		},
		{
			name:     "Assignment comparison - BUG CASE",
			input:    "SameQ((x_^a_ * x_^b_) : x^(a+b), Rule(x_^a_ * x_^b_, x^(a+b)))",
			expected: "True",
		},
	}

	runTestCases(t, tests)
}
