package sexpr

import (
	"testing"

	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/stdlib"
)

func TestRotateLeft(t *testing.T) {
	tests := []struct {
		name     string
		input    core.Expr
		n        int64
		expected string
		isError  bool
	}{
		// Basic List rotation
		{
			name:     "List rotate left by 1",
			input:    core.List{Elements: []core.Expr{core.NewSymbol("List"), core.NewInteger(1), core.NewInteger(2), core.NewInteger(3), core.NewInteger(4), core.NewInteger(5)}},
			n:        1,
			expected: "List(2, 3, 4, 5, 1)",
		},
		{
			name:     "List rotate left by 2",
			input:    core.List{Elements: []core.Expr{core.NewSymbol("List"), core.NewInteger(1), core.NewInteger(2), core.NewInteger(3), core.NewInteger(4), core.NewInteger(5)}},
			n:        2,
			expected: "List(3, 4, 5, 1, 2)",
		},
		{
			name:     "List rotate left by 3",
			input:    core.List{Elements: []core.Expr{core.NewSymbol("List"), core.NewInteger(1), core.NewInteger(2), core.NewInteger(3)}},
			n:        3,
			expected: "List(1, 2, 3)",
		},

		// Zero rotation
		{
			name:     "List rotate left by 0",
			input:    core.List{Elements: []core.Expr{core.NewSymbol("List"), core.NewInteger(1), core.NewInteger(2), core.NewInteger(3)}},
			n:        0,
			expected: "List(1, 2, 3)",
		},

		// Negative rotation (should behave like RotateRight)
		{
			name:     "List rotate left by -1 (equivalent to rotate right by 1)",
			input:    core.List{Elements: []core.Expr{core.NewSymbol("List"), core.NewInteger(1), core.NewInteger(2), core.NewInteger(3), core.NewInteger(4), core.NewInteger(5)}},
			n:        -1,
			expected: "List(5, 1, 2, 3, 4)",
		},
		{
			name:     "List rotate left by -3",
			input:    core.List{Elements: []core.Expr{core.NewSymbol("List"), core.NewInteger(1), core.NewInteger(2), core.NewInteger(3), core.NewInteger(4), core.NewInteger(5)}},
			n:        -3,
			expected: "List(3, 4, 5, 1, 2)",
		},

		// Rotation larger than length
		{
			name:     "List rotate left by length + 1",
			input:    core.List{Elements: []core.Expr{core.NewSymbol("List"), core.NewInteger(1), core.NewInteger(2), core.NewInteger(3)}},
			n:        4, // length is 3, so this is 4 % 3 = 1
			expected: "List(2, 3, 1)",
		},
		{
			name:     "List rotate left by 2 * length + 2",
			input:    core.List{Elements: []core.Expr{core.NewSymbol("List"), core.NewInteger(1), core.NewInteger(2), core.NewInteger(3)}},
			n:        8, // length is 3, so this is 8 % 3 = 2
			expected: "List(3, 1, 2)",
		},

		// String rotation
		{
			name:     "String rotate left by 1",
			input:    core.NewString("hello"),
			n:        1,
			expected: "\"elloh\"",
		},
		{
			name:     "String rotate left by 0",
			input:    core.NewString("hello"),
			n:        0,
			expected: "\"hello\"",
		},
		{
			name:     "String rotate left by -1",
			input:    core.NewString("hello"),
			n:        -1,
			expected: "\"ohell\"",
		},
		{
			name:     "String rotate left by length + 2",
			input:    core.NewString("abc"),
			n:        5, // length is 3, so this is 5 % 3 = 2
			expected: "\"cab\"",
		},

		// ByteArray rotation
		{
			name:     "ByteArray rotate left by 1",
			input:    core.NewByteArray([]byte{1, 2, 3, 4, 5}),
			n:        1,
			expected: "ByteArray(2, 3, 4, 5, 1)",
		},
		{
			name:     "ByteArray rotate left by 0",
			input:    core.NewByteArray([]byte{1, 2, 3}),
			n:        0,
			expected: "ByteArray(1, 2, 3)",
		},
		{
			name:     "ByteArray rotate left by -2",
			input:    core.NewByteArray([]byte{1, 2, 3, 4, 5}),
			n:        -2,
			expected: "ByteArray(4, 5, 1, 2, 3)",
		},
		{
			name:     "ByteArray rotate left by length + 1",
			input:    core.NewByteArray([]byte{1, 2, 3}),
			n:        4, // length is 3, so this is 4 % 3 = 1
			expected: "ByteArray(2, 3, 1)",
		},

		// Empty sequences
		{
			name:     "Empty list rotate left",
			input:    core.List{Elements: []core.Expr{core.NewSymbol("List")}},
			n:        5,
			expected: "List()",
		},
		{
			name:     "Empty string rotate left",
			input:    core.NewString(""),
			n:        3,
			expected: "\"\"",
		},
		{
			name:     "Empty ByteArray rotate left",
			input:    core.NewByteArray([]byte{}),
			n:        2,
			expected: "ByteArray()",
		},

		// Single element sequences
		{
			name:     "Single element list rotate left",
			input:    core.List{Elements: []core.Expr{core.NewSymbol("List"), core.NewInteger(42)}},
			n:        10,
			expected: "List(42)",
		},
		{
			name:     "Single character string rotate left",
			input:    core.NewString("x"),
			n:        -5,
			expected: "\"x\"",
		},

		// Error cases
		{
			name:    "Non-sliceable type (integer)",
			input:   core.NewInteger(42),
			n:       1,
			isError: true,
		},
		{
			name:    "Non-sliceable type (symbol)",
			input:   core.NewSymbol("test"),
			n:       2,
			isError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stdlib.RotateLeft(tt.input, tt.n)

			if tt.isError {
				if !core.IsError(result) {
					t.Errorf("Expected error but got: %v", result)
				}
				return
			}

			if core.IsError(result) {
				t.Fatalf("Unexpected error: %v", result)
			}

			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestRotateRight(t *testing.T) {
	tests := []struct {
		name     string
		input    core.Expr
		n        int64
		expected string
		isError  bool
	}{
		// Basic List rotation
		{
			name:     "List rotate right by 1",
			input:    core.List{Elements: []core.Expr{core.NewSymbol("List"), core.NewInteger(1), core.NewInteger(2), core.NewInteger(3), core.NewInteger(4), core.NewInteger(5)}},
			n:        1,
			expected: "List(5, 1, 2, 3, 4)",
		},
		{
			name:     "List rotate right by 2",
			input:    core.List{Elements: []core.Expr{core.NewSymbol("List"), core.NewInteger(1), core.NewInteger(2), core.NewInteger(3), core.NewInteger(4), core.NewInteger(5)}},
			n:        2,
			expected: "List(4, 5, 1, 2, 3)",
		},
		{
			name:     "List rotate right by 3",
			input:    core.List{Elements: []core.Expr{core.NewSymbol("List"), core.NewInteger(1), core.NewInteger(2), core.NewInteger(3)}},
			n:        3,
			expected: "List(1, 2, 3)",
		},

		// Zero rotation
		{
			name:     "List rotate right by 0",
			input:    core.List{Elements: []core.Expr{core.NewSymbol("List"), core.NewInteger(1), core.NewInteger(2), core.NewInteger(3)}},
			n:        0,
			expected: "List(1, 2, 3)",
		},

		// Negative rotation (should behave like RotateLeft)
		{
			name:     "List rotate right by -1 (equivalent to rotate left by 1)",
			input:    core.List{Elements: []core.Expr{core.NewSymbol("List"), core.NewInteger(1), core.NewInteger(2), core.NewInteger(3), core.NewInteger(4), core.NewInteger(5)}},
			n:        -1,
			expected: "List(2, 3, 4, 5, 1)",
		},
		{
			name:     "List rotate right by -3",
			input:    core.List{Elements: []core.Expr{core.NewSymbol("List"), core.NewInteger(1), core.NewInteger(2), core.NewInteger(3), core.NewInteger(4), core.NewInteger(5)}},
			n:        -3,
			expected: "List(4, 5, 1, 2, 3)",
		},

		// Rotation larger than length
		{
			name:     "List rotate right by length + 1",
			input:    core.List{Elements: []core.Expr{core.NewSymbol("List"), core.NewInteger(1), core.NewInteger(2), core.NewInteger(3)}},
			n:        4, // length is 3, so this is 4 % 3 = 1
			expected: "List(3, 1, 2)",
		},
		{
			name:     "List rotate right by 2 * length + 2",
			input:    core.List{Elements: []core.Expr{core.NewSymbol("List"), core.NewInteger(1), core.NewInteger(2), core.NewInteger(3)}},
			n:        8, // length is 3, so this is 8 % 3 = 2
			expected: "List(2, 3, 1)",
		},

		// String rotation
		{
			name:     "String rotate right by 1",
			input:    core.NewString("hello"),
			n:        1,
			expected: "\"ohell\"",
		},
		{
			name:     "String rotate right by 0",
			input:    core.NewString("hello"),
			n:        0,
			expected: "\"hello\"",
		},
		{
			name:     "String rotate right by -1",
			input:    core.NewString("hello"),
			n:        -1,
			expected: "\"elloh\"",
		},
		{
			name:     "String rotate right by length + 2",
			input:    core.NewString("abc"),
			n:        5, // length is 3, so this is 5 % 3 = 2
			expected: "\"bca\"",
		},

		// ByteArray rotation
		{
			name:     "ByteArray rotate right by 1",
			input:    core.NewByteArray([]byte{1, 2, 3, 4, 5}),
			n:        1,
			expected: "ByteArray(5, 1, 2, 3, 4)",
		},
		{
			name:     "ByteArray rotate right by 0",
			input:    core.NewByteArray([]byte{1, 2, 3}),
			n:        0,
			expected: "ByteArray(1, 2, 3)",
		},
		{
			name:     "ByteArray rotate right by -2",
			input:    core.NewByteArray([]byte{1, 2, 3, 4, 5}),
			n:        -2,
			expected: "ByteArray(3, 4, 5, 1, 2)",
		},
		{
			name:     "ByteArray rotate right by length + 1",
			input:    core.NewByteArray([]byte{1, 2, 3}),
			n:        4, // length is 3, so this is 4 % 3 = 1
			expected: "ByteArray(3, 1, 2)",
		},

		// Empty sequences
		{
			name:     "Empty list rotate right",
			input:    core.List{Elements: []core.Expr{core.NewSymbol("List")}},
			n:        5,
			expected: "List()",
		},
		{
			name:     "Empty string rotate right",
			input:    core.NewString(""),
			n:        3,
			expected: "\"\"",
		},
		{
			name:     "Empty ByteArray rotate right",
			input:    core.NewByteArray([]byte{}),
			n:        2,
			expected: "ByteArray()",
		},

		// Single element sequences
		{
			name:     "Single element list rotate right",
			input:    core.List{Elements: []core.Expr{core.NewSymbol("List"), core.NewInteger(42)}},
			n:        10,
			expected: "List(42)",
		},
		{
			name:     "Single character string rotate right",
			input:    core.NewString("x"),
			n:        -5,
			expected: "\"x\"",
		},

		// Error cases
		{
			name:    "Non-sliceable type (integer)",
			input:   core.NewInteger(42),
			n:       1,
			isError: true,
		},
		{
			name:    "Non-sliceable type (symbol)",
			input:   core.NewSymbol("test"),
			n:       2,
			isError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stdlib.RotateRight(tt.input, tt.n)

			if tt.isError {
				if !core.IsError(result) {
					t.Errorf("Expected error but got: %v", result)
				}
				return
			}

			if core.IsError(result) {
				t.Fatalf("Unexpected error: %v", result)
			}

			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

// Test integration through the evaluator (using the generated wrappers)
func TestRotateIntegration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		isError  bool
	}{
		// RotateLeft integration tests
		{
			name:     "RotateLeft List through evaluator",
			input:    "RotateLeft([1,2,3,4,5], 2)",
			expected: "List(3, 4, 5, 1, 2)",
		},
		{
			name:     "RotateLeft String through evaluator",
			input:    `RotateLeft("hello", 1)`,
			expected: `"elloh"`,
		},
		{
			name:     "RotateLeft with zero",
			input:    "RotateLeft([1,2,3], 0)",
			expected: "List(1, 2, 3)",
		},
		{
			name:     "RotateLeft with negative number",
			input:    "RotateLeft([1,2,3,4,5], -2)",
			expected: "List(4, 5, 1, 2, 3)",
		},
		{
			name:     "RotateLeft with large number",
			input:    "RotateLeft([1,2,3], 7)", // 7 % 3 = 1
			expected: "List(2, 3, 1)",
		},

		// RotateRight integration tests
		{
			name:     "RotateRight List through evaluator",
			input:    "RotateRight([1,2,3,4,5], 2)",
			expected: "List(4, 5, 1, 2, 3)",
		},
		{
			name:     "RotateRight String through evaluator",
			input:    `RotateRight("hello", 1)`,
			expected: `"ohell"`,
		},
		{
			name:     "RotateRight with zero",
			input:    "RotateRight([1,2,3], 0)",
			expected: "List(1, 2, 3)",
		},
		{
			name:     "RotateRight with negative number",
			input:    "RotateRight([1,2,3,4,5], -2)",
			expected: "List(3, 4, 5, 1, 2)",
		},
		{
			name:     "RotateRight with large number",
			input:    "RotateRight([1,2,3], 8)", // 8 % 3 = 2
			expected: "List(2, 3, 1)",
		},

		// Error cases
		{
			name:    "RotateLeft with non-sliceable type",
			input:   "RotateLeft(42, 1)",
			isError: true,
		},
		{
			name:     "RotateRight with non-integer rotation",
			input:    "RotateRight([1,2,3], 1.5)",
			expected: "RotateRight(List(1, 2, 3), 1.5)", // Pattern doesn't match, returns unevaluated
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			evaluator := NewEvaluator()
			result := evaluator.Evaluate(expr)

			if tt.isError {
				if !core.IsError(result) {
					t.Errorf("Expected error but got: %v", result)
				}
				return
			}

			if core.IsError(result) {
				t.Fatalf("Evaluation error: %v", result)
			}

			if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}
