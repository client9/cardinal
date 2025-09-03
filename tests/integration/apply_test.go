package integration

import (
	"testing"
)

func TestApply_BasicFunctionality(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Apply with Plus function",
			input:    `Apply(Plus, [1, 2, 3, 4])`,
			expected: `10`,
		},
		{
			name:     "Apply with Times function",
			input:    `Apply(Times, [2, 3, 4])`,
			expected: `24`,
		},
		{
			name:     "Apply with List function",
			input:    `Apply(List, [1, 2, 3])`,
			expected: `List(1, 2, 3)`,
		},
		{
			name:     "Apply with single argument",
			input:    `Apply(Length, [[1, 2, 3, 4]])`,
			expected: `4`,
		},
		{
			name:     "Apply with empty list",
			input:    `Apply(Plus, [])`,
			expected: `0`,
		},
		{
			name:     "Apply with two arguments",
			input:    `Apply(Power, [2, 8])`,
			expected: `256`,
		},
	}

	runTestCases(t, tests)
}

func TestApply_WithAmpersandSyntax(t *testing.T) {
	t.Skip()
	tests := []TestCase{
		{
			name:     "Apply with & syntax - two arguments",
			input:    `Apply(Plus($1, $2) &, [10, 20])`,
			expected: `30`,
		},
		{
			name:     "Apply with & syntax - three arguments",
			input:    `Apply($1 + $2 + $3 &, [10, 20, 30])`,
			expected: `60`,
		},
		{
			name:     "Apply with & syntax - multiplication",
			input:    `Apply(Times($1, $2, $3) &, [2, 3, 4])`,
			expected: `24`,
		},
		{
			name:     "Apply with & syntax - complex expression",
			input:    `Apply(Plus(Times($1, $2), $3) &, [3, 4, 5])`,
			expected: `17`,
		},
		{
			name:     "Apply with & syntax - single argument",
			input:    `Apply(Times($1, $1) &, [7])`,
			expected: `49`,
		},
		{
			name:     "Apply with & syntax - higher slots",
			input:    `Apply($1 + $2 + $3 + $4 &, [1, 10, 100, 1000])`,
			expected: `1111`,
		},
	}
	runTestCases(t, tests)
}

func TestApply_WithRegularFunction(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Apply with Function syntax - two parameters",
			input:    `Apply(Function([x, y], x + y), [15, 25])`,
			expected: `40`,
		},
		{
			name:     "Apply with Function syntax - three parameters",
			input:    `Apply(Function([a, b, c], a * b + c), [2, 5, 3])`,
			expected: `13`,
		},
		{
			name:     "Apply with Function syntax - single parameter",
			input:    `Apply(Function([x], x * x), [6])`,
			expected: `36`,
		},
	}
	runTestCases(t, tests)
}

func TestApply_EdgeCases(t *testing.T) {

	tests := []TestCase{
		{
			name:      "Apply with non-list should error",
			input:     `Apply(Plus, 42)`,
			expected:  "",
			errorType: "ArgumentError",
		},
		{
			name:     "Apply with wrong number of arguments returns unevaluated",
			input:    `Apply(Plus)`,
			expected: `Apply(Plus)`,
		},
		{
			name:     "Apply with too many arguments returns unevaluated",
			input:    `Apply(Plus, [1, 2], [3, 4])`,
			expected: `Apply(Plus, List(1, 2), List(3, 4))`,
		},
		{
			name:     "Apply ignores list head",
			input:    `Apply(Plus, MyList(5, 10, 15))`,
			expected: `30`,
		},
	}
	runTestCases(t, tests)
}
