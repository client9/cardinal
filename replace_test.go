package sexpr

import (
	"testing"
)

func TestReplaceFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			result := evaluateStringSimple(t, evaluator, test.input)
			if result != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result)
			}
		})
	}
}
