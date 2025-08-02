package integration

import (
	"testing"
)

func TestMap_BasicFunctionality(t *testing.T) {

	tests := []TestCase{
		{
			name:     "Map with Plus function on single elements",
			input:    `Map(Plus($1, 10) &, [1, 2, 3])`,
			expected: `List(11, 12, 13)`,
		},
		{
			name:     "Map with Times function on single elements",
			input:    `Map(Times($1, 2) &, [2, 3, 4])`,
			expected: `List(4, 6, 8)`,
		},
		{
			name:     "Map with Length function",
			input:    `Map(Length, [[1, 2, 3], [4, 5], [6]])`,
			expected: `List(3, 2, 1)`,
		},
		{
			name:     "Map with Head function",
			input:    `Map(Head, [Plus(a, b), Times(x, y), List(1, 2)])`,
			expected: `List(Plus, Times, List)`,
		},
		{
			name:     "Map with empty list",
			input:    `Map(Plus, [])`,
			expected: `List()`,
		},
		{
			name:     "Map with single element",
			input:    `Map(Length, [[1, 2, 3]])`,
			expected: `List(3)`,
		},
	}
	runTestCases(t, tests)
}

func TestMap_WithAmpersandSyntax(t *testing.T) {

	tests := []TestCase{
		{
			name:     "Map with & syntax - double elements",
			input:    `Map(Plus($1, $1) &, [1, 2, 3, 4])`,
			expected: `List(2, 4, 6, 8)`,
		},
		{
			name:     "Map with & syntax - square elements",
			input:    `Map(Times($1, $1) &, [2, 3, 4])`,
			expected: `List(4, 9, 16)`,
		},
		{
			name:     "Map with & syntax - increment",
			input:    `Map(Plus($1, 1) &, [10, 20, 30])`,
			expected: `List(11, 21, 31)`,
		},
		{
			name:     "Map with & syntax - complex expression",
			input:    `Map(Plus(Times($1, 2), 1) &, [1, 2, 3])`,
			expected: `List(3, 5, 7)`,
		},
		{
			name:     "Map with & syntax - string operations",
			input:    `Map(StringLength($1) &, ["hello", "world", "test"])`,
			expected: `List(5, 5, 4)`,
		},
		{
			name:     "Map with & syntax - nested lists",
			input:    `Map(Length($1) &, [[1, 2], [3, 4, 5], [6]])`,
			expected: `List(2, 3, 1)`,
		},
	}

	runTestCases(t, tests)
}

func TestMap_EdgeCases(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Map with non-list should error",
			input:    `Map(Plus, 42)`,
			errorType: `ArgumentError`,
		},
		{
			name:     "Map with wrong number of arguments returns unevaluated",
			input:    `Map(Plus)`,
			expected: `Map(Plus)`,
		},
		{
			name:     "Map with too many arguments returns unevaluated",
			input:    `Map(Plus, [1, 2], [3, 4])`,
			expected: `Map(Plus, List(1, 2), List(3, 4))`,
		},
		{
			name:     "Map preserves head of input list",
			input:    `Map(Plus($1, 1) &, MyList(1, 2, 3))`,
			expected: `MyList(2, 3, 4)`,
		},
	}
	runTestCases(t, tests)
}

func TestMapApply_Integration(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Map then Apply - sum of squares",
			input:    `Apply(Plus, Map(Times($1, $1) &, [1, 2, 3, 4]))`,
			expected: `30`,
		},
		{
			name:     "Apply then Map - distribute and square",
			input:    `Map(Times($1, $1) &, Apply(List, [2, 3, 4]))`,
			expected: `List(4, 9, 16)`,
		},
		{
			name:     "Nested Map with Apply",
			input:    `Map(Apply(Plus, $1) &, [[1, 2], [3, 4], [5, 6]])`,
			expected: `List(3, 7, 11)`,
		},
		{
			name:     "Apply with Map as function",
			input:    `Apply(Function([x], Map(Plus($1, 1) &, x)), [[1, 2, 3]])`,
			expected: `List(2, 3, 4)`,
		},
	}
	runTestCases(t, tests)
}

func TestMapApply_WithBuiltinFunctions(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Map with IntegerQ predicate",
			input:    `Map(IntegerQ, [1, 2.5, 3, "hello"])`,
			expected: `List(True, False, True, False)`,
		},
		{
			name:     "Map with StringQ predicate",
			input:    `Map(StringQ, [1, "hello", 3.14, "world"])`,
			expected: `List(False, True, False, True)`,
		},
		{
			name:     "Apply with Equal comparison",
			input:    `Apply(Equal, [5, 5])`,
			expected: `True`,
		},
		{
			name:     "Apply with Less comparison",
			input:    `Apply(Less, [3, 7])`,
			expected: `True`,
		},
		{
			name:     "Map with Not function",
			input:    `Map(Not, [True, False, True])`,
			expected: `List(False, True, False)`,
		},
		{
			name:     "Apply with MatchQ",
			input:    `Apply(MatchQ, [42, _Integer])`,
			expected: `True`,
		},
	}
	runTestCases(t, tests)
}
