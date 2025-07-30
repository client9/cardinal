package sexpr

import (
	"testing"
)

func TestRuleDelayed_Basic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := evaluator.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestRuleDelayed_Replace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := evaluator.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestRuleDelayed_LexicalScoping(t *testing.T) {
	tests := []struct {
		name        string
		setup       string
		input       string
		expected    string
		description string
	}{
		{
			name:        "Local bindings don't pollute global scope",
			setup:       "x = 100",
			input:       "Replace(5, RuleDelayed(y_, x + y))",
			expected:    "105",
			description: "RuleDelayed evaluates x + y with global x=100 and local y=5",
		},
		{
			name:        "Pattern variables are locally scoped",
			setup:       "y = 999",
			input:       "Replace(42, RuleDelayed(y_, y * 2))",
			expected:    "84",
			description: "Pattern variable y shadows global y in RuleDelayed",
		},
		{
			name:        "Complex expression with local evaluation",
			setup:       "f(x_) := x^2",
			input:       "Replace(3, RuleDelayed(x_, f(x) + x))",
			expected:    "12",
			description: "Function f is evaluated with local x=3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()

			// Setup global variables
			if tt.setup != "" {
				setupExpr, err := ParseString(tt.setup)
				if err != nil {
					t.Fatalf("Setup parse error: %v", err)
				}
				evaluator.Evaluate(setupExpr)
			}

			// Test the main expression
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := evaluator.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s (%s)", tt.expected, result.String(), tt.description)
			}
		})
	}
}

func TestRuleDelayed_ReplaceAll(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
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
			input:    "ReplaceAll(Plus(1, 2, 3), RuleDelayed(x_, x * 10))",
			expected: "60",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := evaluator.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestRuleDelayed_CompareWithRule(t *testing.T) {
	tests := []struct {
		name         string
		setup        string
		ruleInput    string
		delayedInput string
		description  string
	}{
		{
			name:         "Pattern variable scoping difference",
			setup:        "y = 999",
			ruleInput:    "Replace(5, Rule(y_, y + 1))",
			delayedInput: "Replace(5, RuleDelayed(y_, y + 1))",
			description:  "Rule substitutes pattern var, RuleDelayed uses local binding",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator1 := NewEvaluator()
			evaluator2 := NewEvaluator()

			// Setup global variables in both evaluators
			if tt.setup != "" {
				setupExpr, err := ParseString(tt.setup)
				if err != nil {
					t.Fatalf("Setup parse error: %v", err)
				}
				evaluator1.Evaluate(setupExpr)
				evaluator2.Evaluate(setupExpr)
			}

			// Test Rule
			ruleExpr, err := ParseString(tt.ruleInput)
			if err != nil {
				t.Fatalf("Rule parse error: %v", err)
			}
			ruleResult := evaluator1.Evaluate(ruleExpr)

			// Test RuleDelayed
			delayedExpr, err := ParseString(tt.delayedInput)
			if err != nil {
				t.Fatalf("RuleDelayed parse error: %v", err)
			}
			delayedResult := evaluator2.Evaluate(delayedExpr)

			t.Logf("Rule result: %s", ruleResult.String())
			t.Logf("RuleDelayed result: %s", delayedResult.String())
			t.Logf("Description: %s", tt.description)

			// They should be different to demonstrate the scoping difference
			if ruleResult.Equal(delayedResult) {
				t.Errorf("Expected Rule and RuleDelayed to produce different results")
			}
		})
	}
}
