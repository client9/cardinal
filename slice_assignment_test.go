package sexpr

import (
	"github.com/client9/sexpr/core"
	"testing"
)

func TestSliceAssignment(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		// List element assignment
		{
			name:     "List element assignment",
			input:    "[1,2,3][1] = 9",
			expected: "List(9, 2, 3)",
			hasError: false,
		},
		{
			name:     "List element assignment middle",
			input:    "[1,2,3][2] = 99",
			expected: "List(1, 99, 3)",
			hasError: false,
		},
		{
			name:     "List element assignment last",
			input:    "[1,2,3][3] = 999",
			expected: "List(1, 2, 999)",
			hasError: false,
		},
		{
			name:     "List negative indexing",
			input:    "[1,2,3][-1] = 88",
			expected: "List(1, 2, 88)",
			hasError: false,
		},

		// List range assignment
		{
			name:     "List range assignment same size",
			input:    "[1,2,3,4,5][2:3] = [7,8]",
			expected: "List(1, 7, 8, 4, 5)",
			hasError: false,
		},
		{
			name:     "List range assignment expansion",
			input:    "[1,2,3][2:2] = [7,8,9]",
			expected: "List(1, 7, 8, 9, 3)",
			hasError: false,
		},
		{
			name:     "List range assignment contraction",
			input:    "[1,2,3,4,5][2:4] = [99]",
			expected: "List(1, 99, 5)",
			hasError: false,
		},
		{
			name:     "List range assignment beginning",
			input:    "[1,2,3][1:2] = [7,8]",
			expected: "List(7, 8, 3)",
			hasError: false,
		},
		{
			name:     "List range assignment end",
			input:    "[1,2,3][2:3] = [7,8]",
			expected: "List(1, 7, 8)",
			hasError: false,
		},

		// String character assignment
		{
			name:     "String character assignment",
			input:    "\"abc\"[2] = \"x\"",
			expected: "\"axc\"",
			hasError: false,
		},
		{
			name:     "String character assignment first",
			input:    "\"abc\"[1] = \"z\"",
			expected: "\"zbc\"",
			hasError: false,
		},
		{
			name:     "String character assignment last",
			input:    "\"abc\"[3] = \"z\"",
			expected: "\"abz\"",
			hasError: false,
		},

		// String slice assignment
		{
			name:     "String slice assignment same size",
			input:    "\"hello\"[2:3] = \"xy\"",
			expected: "\"hxylo\"",
			hasError: false,
		},
		{
			name:     "String slice assignment expansion",
			input:    "\"abc\"[2:2] = \"xyz\"",
			expected: "\"axyzc\"",
			hasError: false,
		},
		{
			name:     "String slice assignment contraction",
			input:    "\"hello\"[2:4] = \"x\"",
			expected: "\"hxo\"",
			hasError: false,
		},
		{
			name:     "String slice assignment beginning",
			input:    "\"hello\"[1:2] = \"XY\"",
			expected: "\"XYllo\"",
			hasError: false,
		},

		// Complex expressions
		{
			name:     "List with expression index",
			input:    "[1,2,3][1+1] = 99",
			expected: "List(1, 99, 3)",
			hasError: false,
		},
		{
			name:     "List with expression value",
			input:    "[1,2,3][2] = 10*2",
			expected: "List(1, 20, 3)",
			hasError: false,
		},
		{
			name:     "Nested list assignment",
			input:    "[[1,2],[3,4]][1] = [9,9]",
			expected: "List(List(9, 9), List(3, 4))",
			hasError: false,
		},

		// Error cases
		{
			name:     "Index out of bounds",
			input:    "[1,2,3][5] = 9",
			expected: "",
			hasError: true,
		},
		{
			name:     "Negative index out of bounds",
			input:    "[1,2,3][-5] = 9",
			expected: "",
			hasError: true,
		},
		{
			name:     "String multi-char assignment to single position",
			input:    "\"abc\"[2] = \"xyz\"",
			expected: "",
			hasError: true,
		},
		{
			name:     "String non-string value assignment",
			input:    "\"abc\"[2] = 123",
			expected: "",
			hasError: true,
		},
		{
			name:     "Assignment to non-sliceable",
			input:    "123[1] = 9",
			expected: "",
			hasError: true,
		},
		{
			name:     "Invalid slice range",
			input:    "[1,2,3][3:1] = [9]",
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			expr, err := ParseString(tt.input)

			if err != nil {
				if !tt.hasError {
					t.Errorf("Unexpected parse error: %v", err)
				}
				return
			}

			result := evaluator.Evaluate(expr)

			if tt.hasError {
				if !IsError(result) {
					t.Errorf("Expected error but got result: %v", result)
				}
				return
			}

			if IsError(result) {
				t.Errorf("Unexpected evaluation error: %v", result)
				return
			}

			resultStr := result.String()
			if resultStr != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, resultStr)
			}
		})
	}
}

func TestSliceAssignmentImmutability(t *testing.T) {
	// Test that slice assignment returns new objects and doesn't modify originals
	evaluator := NewEvaluator()

	// Set up original list
	evaluator.GetContext().Set("original", NewList(core.NewSymbol("List"), core.NewInteger(1), core.NewInteger(2), core.NewInteger(3)))

	// Perform assignment
	expr, err := ParseString("original[1] = 99")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result := evaluator.Evaluate(expr)
	if IsError(result) {
		t.Fatalf("Evaluation error: %v", result)
	}

	// Check that result is modified
	if result.String() != "List(99, 2, 3)" {
		t.Errorf("Expected modified result List(99, 2, 3), got %v", result)
	}

	// Check that original is unchanged
	original, _ := evaluator.GetContext().Get("original")
	if original.String() != "List(1, 2, 3)" {
		t.Errorf("Original should be unchanged, but got %v", original)
	}
}

func TestByteArrayAssignment(t *testing.T) {
	tests := []struct {
		name     string
		setup    string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "ByteArray element assignment",
			setup:    "arr = ByteArray(\"hello\")",
			input:    "arr[1] = 72", // 'H' = 72
			expected: "ByteArray(72, 101, 108, 108, 111)",
			hasError: false,
		},
		{
			name:     "ByteArray slice assignment",
			setup:    "arr = ByteArray(\"hello\")",
			input:    "arr[1:2] = [65, 66]", // 'AB' = [65, 66]
			expected: "ByteArray(65, 66, 108, 108, 111)",
			hasError: false,
		},
		{
			name:     "ByteArray invalid byte value",
			setup:    "arr = ByteArray(\"hello\")",
			input:    "arr[1] = 256", // Invalid byte value
			expected: "",
			hasError: true,
		},
		{
			name:     "ByteArray non-integer assignment",
			setup:    "arr = ByteArray(\"hello\")",
			input:    "arr[1] = \"x\"", // String instead of integer
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()

			// Setup
			if tt.setup != "" {
				setupExpr, err := ParseString(tt.setup)
				if err != nil {
					t.Fatalf("Setup parse error: %v", err)
				}
				setupResult := evaluator.Evaluate(setupExpr)
				if IsError(setupResult) {
					t.Fatalf("Setup evaluation error: %v", setupResult)
				}
			}

			// Test assignment
			expr, err := ParseString(tt.input)
			if err != nil {
				if !tt.hasError {
					t.Errorf("Unexpected parse error: %v", err)
				}
				return
			}

			result := evaluator.Evaluate(expr)

			if tt.hasError {
				if !IsError(result) {
					t.Errorf("Expected error but got result: %v", result)
				}
				return
			}

			if IsError(result) {
				t.Errorf("Unexpected evaluation error: %v", result)
				return
			}

			resultStr := result.String()
			if resultStr != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, resultStr)
			}
		})
	}
}

func TestSliceAssignmentEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		// Empty assignments
		{
			name:     "Empty list assignment",
			input:    "[1,2,3][2:3] = []",
			expected: "List(1)",
			hasError: false,
		},
		{
			name:     "Empty string assignment",
			input:    "\"hello\"[2:4] = \"\"",
			expected: "\"ho\"",
			hasError: false,
		},

		// Single element to multiple
		{
			name:     "Single element to list",
			input:    "[1,2,3][2] = [7,8,9]",
			expected: "List(1, List(7, 8, 9), 3)",
			hasError: false,
		},

		// Assignment at boundaries
		{
			name:     "Assignment at start",
			input:    "[1,2,3][1:1] = [99]",
			expected: "List(99, 2, 3)",
			hasError: false,
		},
		{
			name:     "Assignment at end",
			input:    "[1,2,3][3:3] = [99]",
			expected: "List(1, 2, 99)",
			hasError: false,
		},

		// Unicode string handling
		{
			name:     "Unicode string assignment",
			input:    "\"café\"[2] = \"x\"",
			expected: "\"cxfé\"",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			expr, err := ParseString(tt.input)

			if err != nil {
				if !tt.hasError {
					t.Errorf("Unexpected parse error: %v", err)
				}
				return
			}

			result := evaluator.Evaluate(expr)

			if tt.hasError {
				if !IsError(result) {
					t.Errorf("Expected error but got result: %v", result)
				}
				return
			}

			if IsError(result) {
				t.Errorf("Unexpected evaluation error: %v", result)
				return
			}

			resultStr := result.String()
			if resultStr != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, resultStr)
			}
		})
	}
}
