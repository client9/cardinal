package integration

import (
	"testing"
)

func TestSetDelayed(t *testing.T) {
	tests := []TestCase{
		{
			name:     "SetDelayed simple variable",
			input:    `SetDelayed(x, Plus(1, 2)); x`,
			expected: `3`,
		},
		{
			name:     "SetDelayed returns Null",
			input:    `SetDelayed(y, 42)`,
			expected: `Null`,
		},
		{
			name:     "SetDelayed with function definition",
			input:    `SetDelayed(f(x_), Times(x, 2)); f(5)`,
			expected: `10`,
		},
		{
			name:  "SetDelayed invokes function",
			input: "x := RandomReal(); x != x",
			expected: "True",
		},
	}

	runTestCases(t, tests)
}

