package integration

import (
	"testing"
)

func TestSliceSyntaxBasicIndexing(t *testing.T) {
	tests := []TestCase{
		{
			name:     "List single index access",
			input:    "[1,2,3,4,5][3]",
			expected: "3",
		},
		{
			name:     "String single index access",
			input:    `"hello"[2]`,
			expected: "'e'",
		},
		{
			name:     "Variable list indexing",
			input:    "list = [10,20,30]; list[2]",
			expected: "20",
		},
		{
			name:     "Variable string indexing",
			input:    `str = "world"; str[1]`,
			expected: "'w'",
		},
		{
			name:     "First element",
			input:    "[1,2,3,4,5][1]",
			expected: "1",
		},
		{
			name:     "Last element with length",
			input:    "list = [1,2,3,4,5]; list[Length(list)]",
			expected: "5",
		},
	}

	runTestCases(t, tests)
}

func TestSliceSyntaxRanges(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Basic slice range",
			input:    "[1,2,3,4,5][2:4]",
			expected: "List(2, 3, 4)",
		},
		{
			name:     "String slice range",
			input:    `"hello"[2:4]`,
			expected: `"ell"`,
		},
		{
			name:     "Slice from beginning to index",
			input:    "[1,2,3,4,5][:3]",
			expected: "List(1, 2, 3)",
		},
		{
			name:     "Slice from index to end",
			input:    "[1,2,3,4,5][3:]",
			expected: "List(3, 4, 5)",
		},
		{
			name:     "String slice from beginning",
			input:    `"hello"[:3]`,
			expected: `"hel"`,
		},
		{
			name:     "String slice from index to end",
			input:    `"hello"[3:]`,
			expected: `"llo"`,
		},
		{
			name:     "Single element slice",
			input:    "[1,2,3,4,5][3:3]",
			expected: "List(3)",
		},
		{
			name:     "Variable slice range",
			input:    "data = [10,20,30,40,50]; data[2:4]",
			expected: "List(20, 30, 40)",
		},
	}

	runTestCases(t, tests)
}

func TestSliceSyntaxNegativeIndexing(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Negative index single element",
			input:    "[1,2,3,4,5][-1]",
			expected: "5",
		},
		{
			name:     "Negative index from end",
			input:    "[1,2,3,4,5][-2]",
			expected: "4",
		},
		{
			name:     "Negative slice from end",
			input:    "[1,2,3,4,5][-2:]",
			expected: "List(4, 5)",
		},
		{
			name:     "Negative slice range",
			input:    "[1,2,3,4,5][-3:-1]",
			expected: "List(3, 4, 5)",
		},
		{
			name:     "String negative indexing",
			input:    `"hello"[-1]`,
			expected: "'o'",
		},
		{
			name:     "String negative slice",
			input:    `"hello"[-3:]`,
			expected: `"llo"`,
		},
	}

	runTestCases(t, tests)

}

func TestSliceSyntaxSliceableTypes(t *testing.T) {
	tests := []TestCase{
		{
			name:     "List slicing",
			input:    "List(1, 2, 3, 4, 5)[2:4]",
			expected: "List(2, 3, 4)",
		},
		{
			name:     "ByteArray creation and slicing",
			input:    `ByteArray("hello")[2:4]`,
			expected: "ByteArray(101, 108, 108)", // "ell"
		},
		{
			name:     "ByteArray indexing",
			input:    `ByteArray("hello")[1]`,
			expected: "104", // ASCII code for 'h'
		},
	}
	runTestCases(t, tests)
}

/*
func TestSliceSyntaxErrorCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "Non-sliceable type",
			input:       "42[1]",
			expectError: true,
		},
		{
			name:        "Index out of bounds",
			input:       "[1,2,3][10]",
			expectError: true,
		},
		{
			name:        "Negative index too far",
			input:       "[1,2,3][-10]",
			expectError: true,
		},
		{
			name:        "Empty slice bounds",
			input:       "[1,2,3][:]",
			expectError: true, // Should be parse error
		},
		{
			name:        "Invalid slice range",
			input:       "[1,2,3][5:2]", // start > end
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				if tt.expectError {
					return // Expected parse error
				}
				t.Fatalf("Unexpected parse error: %v", err)
			}

			evaluator := NewEvaluator()
			result := evaluator.Evaluate(expr)

			if tt.expectError {
				if !core.IsError(result) {
					t.Errorf("Expected error but got: %v", result)
				}
			} else {
				if core.IsError(result) {
					t.Errorf("Unexpected error: %v", result)
				}
			}
		})
	}
}
*/

func TestSliceSyntaxPrecedence(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Function call then slice",
			input:    "list = [1,2,3,4,5]; First([list, list])[2]",
			expected: "2",
		},
		{
			name:     "Arithmetic in slice index",
			input:    "[1,2,3,4,5][1 + 1]",
			expected: "2",
		},
		{
			name:     "Slice with variable indices",
			input:    "list = [1,2,3,4,5]; start = 2; end = 4; list[start:end]",
			expected: "List(2, 3, 4)",
		},
		{
			name:     "Nested list access",
			input:    "[[1,2],[3,4],[5,6]][2][1]",
			expected: "3",
		},
		{
			name:     "String slice then indexing",
			input:    `"hello world"[1:5][2]`,
			expected: "'e'",
		},
	}
	runTestCases(t, tests)
}
