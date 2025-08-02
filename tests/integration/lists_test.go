package integration

import (
	"testing"
)

func TestRotateLeft_Integration(t *testing.T) {
	tests := []TestCase{
		{
			name:     "RotateLeft basic list",
			input:    "RotateLeft([1, 2, 3, 4], 1)",
			expected: "List(2, 3, 4, 1)",
		},
		{
			name:     "RotateLeft by multiple positions",
			input:    "RotateLeft([1, 2, 3, 4, 5], 2)",
			expected: "List(3, 4, 5, 1, 2)",
		},
		{
			name:     "RotateLeft by zero",
			input:    "RotateLeft([1, 2, 3], 0)",
			expected: "List(1, 2, 3)",
		},
		{
			name:     "RotateLeft negative (equivalent to right)",
			input:    "RotateLeft([1, 2, 3, 4, 5], -1)",
			expected: "List(5, 1, 2, 3, 4)",
		},
		{
			name:     "RotateLeft with nested expression",
			input:    "RotateLeft([Plus(1, 2), Times(3, 4), 5], 1)",
			expected: "List(12, 5, 3)",
		},
		{
			name:     "RotateLeft with string elements",
			input:    "RotateLeft([\"a\", \"b\", \"c\"], 1)",
			expected: "List(\"b\", \"c\", \"a\")",
		},
	}

	runTestCases(t, tests)
}

func TestRotateRight(t *testing.T) {
	tests := []TestCase{
		{
			name:     "RotateRight basic list",
			input:    "RotateRight([1, 2, 3, 4], 1)",
			expected: "List(4, 1, 2, 3)",
		},
		{
			name:     "RotateRight by multiple positions",
			input:    "RotateRight([1, 2, 3, 4, 5], 2)",
			expected: "List(4, 5, 1, 2, 3)",
		},
		{
			name:     "RotateRight by zero",
			input:    "RotateRight([1, 2, 3], 0)",
			expected: "List(1, 2, 3)",
		},
		{
			name:     "RotateRight negative (equivalent to left)",
			input:    "RotateRight([1, 2, 3, 4, 5], -1)",
			expected: "List(2, 3, 4, 5, 1)",
		},
		{
			name:     "RotateRight with nested expression",
			input:    "RotateRight([Plus(1, 2), Times(3, 4), 5], 1)",
			expected: "List(5, 3, 12)",
		},
		{
			name:     "RotateRight chained operations",
			input:    "RotateRight(RotateLeft([1, 2, 3, 4], 1), 1)",
			expected: "List(1, 2, 3, 4)", // Should return to original
		},
	}

	runTestCases(t, tests)
}

func TestTakeDrop_Integration(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Take first n elements",
			input:    "Take([1, 2, 3, 4, 5], 3)",
			expected: "List(1, 2, 3)",
		},
		{
			name:     "Take with range specification",
			input:    "Take([1, 2, 3, 4, 5], [2, 4])",
			expected: "List(2, 3, 4)",
		},
		{
			name:     "Drop first n elements",
			input:    "Drop([1, 2, 3, 4, 5], 2)",
			expected: "List(3, 4, 5)",
		},
		{
			name:     "Drop with range specification - NOT IMPLEMENTED",
			input:    "Drop([1, 2, 3, 4, 5], [2, 3])",
			expected: "$Failed(NotImplemented)",
			skip:     true,
		},
		{
			name:     "Take negative count",
			input:    "Take([1, 2, 3, 4, 5], -2)",
			expected: "List(4, 5)",
		},
		{
			name:     "Drop negative count",
			input:    "Drop([1, 2, 3, 4, 5], -2)",
			expected: "List(1, 2, 3)",
		},
	}

	runTestCases(t, tests)
}

func TestListAccess_Integration(t *testing.T) {
	tests := []TestCase{
		{
			name:     "First element of list",
			input:    "First([1, 2, 3])",
			expected: "1",
		},
		{
			name:     "Last element of list",
			input:    "Last([1, 2, 3])",
			expected: "3",
		},
		{
			name:     "Rest of list (all but first)",
			input:    "Rest([1, 2, 3, 4])",
			expected: "List(2, 3, 4)",
		},
		{
			name:     "Most of list (all but last)",
			input:    "Most([1, 2, 3, 4])",
			expected: "List(1, 2, 3)",
		},
		{
			name:     "Part with index",
			input:    "Part([1, 2, 3, 4], 2)",
			expected: "2",
		},
		{
			name:     "Length of list",
			input:    "Length([1, 2, 3, 4, 5])",
			expected: "5",
		},
	}

	runTestCases(t, tests)
}

func TestListManipulation_Integration(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Append to list",
			input:    "Append([1, 2, 3], 4)",
			expected: "List(1, 2, 3, 4)",
		},
		{
			name:     "Reverse list",
			input:    "Reverse([1, 2, 3, 4])",
			expected: "Reverse(List(1, 2, 3, 4))",
		},
		{
			name:     "Nested list operations",
			input:    "First(Rest([1, 2, 3, 4]))",
			expected: "2",
		},
		{
			name:     "Complex list expression",
			input:    "Length(Take(Drop([1, 2, 3, 4, 5], 1), 3))",
			expected: "3",
		},
	}

	runTestCases(t, tests)
}
