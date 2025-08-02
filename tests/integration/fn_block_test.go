package integration

import (
	"testing"
)

func TestBlockBasic(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Block with single variable assignment",
			input:    `Block(List(Set(x, 5)), x)`,
			expected: `5`,
		},
		{
			name:     "Block with variable clearing",
			input:    `Block(List(x), x)`,
			expected: `x`,
		},
		{
			name:     "Block with arithmetic",
			input:    `Block(List(Set(x, 3)), Plus(x, 2))`,
			expected: `5`,
		},
		{
			name:     "Block with multiple variables",
			input:    `Block(List(Set(x, 1), Set(y, 2)), Plus(x, y))`,
			expected: `3`,
		},
		{
			name:     "Block preserves variable isolation",
			input:    `Block(List(Set(localVar, 42)), localVar)`,
			expected: `42`,
		},
	}

	runTestCases(t, tests)
}
