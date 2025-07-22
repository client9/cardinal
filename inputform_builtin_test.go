package sexpr

import (
	"testing"
)

func TestInputFormBuiltin(t *testing.T) {
	eval := NewEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "InputForm of atom",
			input:    `InputForm(42)`,
			expected: `"42"`,
		},
		{
			name:     "InputForm of symbol",
			input:    `InputForm(x)`,
			expected: `"x"`,
		},
		{
			name:     "InputForm of Plus - evaluates first",
			input:    `InputForm(Plus(1, 2))`,
			expected: `"3"`, // Plus(1,2) evaluates to 3 first
		},
		{
			name:     "InputForm of Times - evaluates first",
			input:    `InputForm(Times(3, 4))`,
			expected: `"12"`, // Times(3,4) evaluates to 12 first
		},
		{
			name:     "InputForm of List",
			input:    `InputForm(List(1, 2, 3))`,
			expected: `"[1, 2, 3]"`,
		},
		{
			name:     "InputForm of Set - evaluates first",
			input:    `InputForm(Set(x, 5))`,
			expected: `"5"`, // Set(x, 5) evaluates to 5 (the assigned value)
		},
		{
			name:     "InputForm of precedence example - evaluates first",
			input:    `InputForm(Plus(1, Times(2, 3)))`,
			expected: `"7"`, // Plus(1, Times(2,3)) = Plus(1,6) = 7
		},
		{
			name:     "InputForm with parentheses - evaluates first",
			input:    `InputForm(Times(Plus(1, 2), 3))`,
			expected: `"9"`, // Times(Plus(1,2), 3) = Times(3,3) = 9
		},
		{
			name:     "InputForm evaluates like FullForm",
			input:    `Equal(InputForm(Plus(1, 2)), "3")`,
			expected: `True`, // Both evaluate Plus(1,2) to 3 first
		},
		{
			name:     "FullForm vs InputForm same for evaluated expressions",
			input:    `Equal(FullForm(Plus(1, 2)), InputForm(Plus(1, 2)))`,
			expected: `True`, // Both return "3" for evaluated expressions
		},
		{
			name:     "InputForm of held expression shows infix",
			input:    `InputForm(Hold(Plus(1, 2)))`,
			expected: `"Hold(1 + 2)"`, // Hold prevents evaluation, shows InputForm of structure
		},
		{
			name:     "InputForm vs FullForm difference on held expressions",
			input:    `Unequal(FullForm(Hold(Plus(1, 2))), InputForm(Hold(Plus(1, 2))))`,
			expected: `True`, // FullForm: "Hold(Plus(1, 2))", InputForm: "Hold(1 + 2)"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Errorf("parse error: %v", err)
				return
			}

			result := eval.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.String())
			}
		})
	}
}

func TestInputFormBuiltinErrors(t *testing.T) {
	eval := NewEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string // Expected unevaluated result when pattern doesn't match
	}{
		{
			name:     "InputForm with no arguments",
			input:    `InputForm()`,
			expected: `InputForm()`,
		},
		{
			name:     "InputForm with too many arguments",
			input:    `InputForm(1, 2)`,
			expected: `InputForm(1, 2)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Errorf("parse error: %v", err)
				return
			}

			result := eval.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.String())
			}
		})
	}
}

func TestInputFormBuiltinComparison(t *testing.T) {
	eval := NewEvaluator()

	// Test InputForm vs FullForm on held expressions (to see symbolic structure)
	examples := []string{
		"Hold(Plus(1, 2))",              // InputForm shows "Hold(1 + 2)", FullForm shows "Hold(Plus(1, 2))"
		"Hold(Times(3, 4))",             // InputForm shows "Hold(3 * 4)", FullForm shows "Hold(Times(3, 4))"
		"Hold(List(1, 2, 3))",           // InputForm shows "Hold([1, 2, 3])", FullForm shows "Hold(List(1, 2, 3))"
		"Hold(Set(x, 5))",               // InputForm shows "Hold(x = 5)", FullForm shows "Hold(Set(x, 5))"
		"Hold(Association(Rule(a, b)))", // InputForm shows "Hold({a: b})", FullForm shows "Hold(Association(Rule(a, b)))"
	}

	for _, example := range examples {
		t.Run("compare_"+example, func(t *testing.T) {
			// Evaluate FullForm
			fullFormExpr, err := ParseString("FullForm(" + example + ")")
			if err != nil {
				t.Errorf("parse error for FullForm: %v", err)
				return
			}
			fullFormResult := eval.Evaluate(fullFormExpr)

			// Evaluate InputForm
			inputFormExpr, err := ParseString("InputForm(" + example + ")")
			if err != nil {
				t.Errorf("parse error for InputForm: %v", err)
				return
			}
			inputFormResult := eval.Evaluate(inputFormExpr)

			// They should be different for most expressions (except atoms)
			// Both should be strings
			if !fullFormResult.Equal(inputFormResult) {
				// This is expected for most expressions
				t.Logf("FullForm(%s) = %s", example, fullFormResult.String())
				t.Logf("InputForm(%s) = %s", example, inputFormResult.String())
			}
		})
	}
}
