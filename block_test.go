package sexpr

import (
	"strings"
	"testing"
)

func TestBlockBasic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "Block with single variable assignment",
			input:    `Block(List(Set(x, 5)), x)`,
			expected: `5`,
			hasError: false,
		},
		{
			name:     "Block with variable clearing",
			input:    `Block(List(x), x)`,
			expected: `x`,
			hasError: false,
		},
		{
			name:     "Block with arithmetic",
			input:    `Block(List(Set(x, 3)), Plus(x, 2))`,
			expected: `5`,
			hasError: false,
		},
		{
			name:     "Block with multiple variables",
			input:    `Block(List(Set(x, 1), Set(y, 2)), Plus(x, y))`,
			expected: `3`,
			hasError: false,
		},
		{
			name:     "Block preserves variable isolation",
			input:    `Block(List(Set(x, 10)), Block(List(Set(x, 20)), x))`,
			expected: `20`,
			hasError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			result := evaluateStringSimple(t, evaluator, test.input)

			if test.hasError {
				if !strings.HasPrefix(result, "$Failed") {
					t.Errorf("Expected error, but got: %s", result)
				}
			} else {
				if result != test.expected {
					t.Errorf("Expected %s, got %s", test.expected, result)
				}
			}
		})
	}
}

func TestBlockErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Block with wrong number of arguments",
			input: `Block(List(x))`,
		},
		{
			name:  "Block with non-list first argument",
			input: `Block(x, y)`,
		},
		{
			name:  "Block with invalid variable specification",
			input: `Block(List(123), x)`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			result := evaluateStringSimple(t, evaluator, test.input)

			if !strings.HasPrefix(result, "$Failed") {
				t.Errorf("Expected error, but got: %s", result)
			}
		})
	}
}

func TestBlockVariableScoping(t *testing.T) {
	evaluator := NewEvaluator()

	// Set a global variable
	evaluateStringSimple(t, evaluator, `Set(x, 100)`)

	// Use Block to temporarily change it
	result1 := evaluateStringSimple(t, evaluator, `Block(List(Set(x, 5)), x)`)
	if result1 != "5" {
		t.Errorf("Expected 5 inside Block, got %s", result1)
	}

	// Check that original value is preserved
	result2 := evaluateStringSimple(t, evaluator, `x`)
	if result2 != "100" {
		t.Errorf("Expected 100 after Block, got %s", result2)
	}
}

func TestTableSimple(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "Table with zero repetitions",
			input:    `Table(x, 0)`,
			expected: `List()`,
			hasError: false,
		},
		{
			name:     "Table with one repetition",
			input:    `Table(x, 1)`,
			expected: `List(x)`,
			hasError: false,
		},
		{
			name:     "Table with multiple repetitions",
			input:    `Table(42, 3)`,
			expected: `List(42, 42, 42)`,
			hasError: false,
		},
		{
			name:     "Table with expression",
			input:    `Table(Plus(1, 2), 2)`,
			expected: `List(3, 3)`,
			hasError: false,
		},
		{
			name:     "Table with symbol and arithmetic",
			input:    `Block(List(Set(y, 5)), Table(Times(y, 2), 3))`,
			expected: `List(10, 10, 10)`,
			hasError: false,
		},
		{
			name:     "Table evaluates expression each time",
			input:    `Table(Plus(x, 1), 2)`,
			expected: `List(Plus(1, x), Plus(1, x))`, // Plus is Orderless: Plus(x, 1) → Plus(1, x)
			hasError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			result := evaluateStringSimple(t, evaluator, test.input)

			if test.hasError {
				if !strings.HasPrefix(result, "$Failed") {
					t.Errorf("Expected error, but got: %s", result)
				}
			} else {
				if result != test.expected {
					t.Errorf("Expected %s, got %s", test.expected, result)
				}
			}
		})
	}
}

func TestTableErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Table with negative count",
			input: `Table(x, -1)`,
		},
		{
			name:  "Table with wrong number of arguments",
			input: `Table(x)`,
		},
		{
			name:  "Table with three arguments (not yet implemented)",
			input: `Table(x, y, z)`,
		},
		{
			name:  "Table with non-integer, non-List second argument",
			input: `Table(x, "invalid")`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			result := evaluateStringSimple(t, evaluator, test.input)

			if !strings.HasPrefix(result, "$Failed") {
				t.Errorf("Expected error, but got: %s", result)
			}
		})
	}
}
func TestTableIterator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		// Basic iterator forms
		{
			name:     "Table with List(i, max) - simple case",
			input:    `Table(i, List(i, 3))`,
			expected: `List(1, 2, 3)`,
			hasError: false,
		},
		{
			name:     "Table with List(i, start, end) - simple case",
			input:    `Table(i, List(i, 2, 4))`,
			expected: `List(2, 3, 4)`,
			hasError: false,
		},
		{
			name:     "Table with List(i, start, end, increment) - simple case",
			input:    `Table(i, List(i, 1, 5, 2))`,
			expected: `List(1, 3, 5)`,
			hasError: false,
		},
		{
			name:     "Table with expression using iterator",
			input:    `Table(Times(i, i), List(i, 1, 4))`,
			expected: `List(1, 4, 9, 16)`,
			hasError: false,
		},
		{
			name:     "Table with Plus expression",
			input:    `Table(Plus(i, 10), List(i, 1, 3))`,
			expected: `List(11, 12, 13)`,
			hasError: false,
		},
		{
			name:     "Table with negative increment",
			input:    `Table(i, List(i, 5, 1, Minus(1)))`,
			expected: `List(5, 4, 3, 2, 1)`,
			hasError: false,
		},
		{
			name:     "Table with zero range",
			input:    `Table(i, List(i, 3, 3))`,
			expected: `List(3)`,
			hasError: false,
		},
		{
			name:     "Table with empty range (start > end with positive increment)",
			input:    `Table(i, List(i, 5, 3))`,
			expected: `List()`,
			hasError: false,
		},
		{
			name:     "Table with Real numbers",
			input:    `Table(i, List(i, 1.0, 3.0, 0.5))`,
			expected: `List(1.0, 1.5, 2.0, 2.5, 3.0)`,
			hasError: false,
		},
		{
			name:     "Table with complex expression and scoping",
			input:    `Block(List(Set(x, 10)), Table(Plus(i, x), List(i, 1, 3)))`,
			expected: `List(11, 12, 13)`,
			hasError: false,
		},
		{
			name:     "Table with expressions in iterator spec - original issue",
			input:    `Block(List(Set(n, 2)), Table(x, List(x, n, Times(3, n), n)))`,
			expected: `List(2, 4, 6)`,
			hasError: false,
		},
		{
			name:     "Table with Plus expressions in iterator",
			input:    `Block(List(Set(a, 5)), Table(i, List(i, Plus(a, 1), Plus(a, 3))))`,
			expected: `List(6, 7, 8)`,
			hasError: false,
		},
		{
			name:     "Table with variable increment expression",
			input:    `Block(List(Set(step, 3)), Table(i, List(i, 1, 7, step)))`,
			expected: `List(1, 4, 7)`,
			hasError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			result := evaluateStringSimple(t, evaluator, test.input)

			if test.hasError {
				if !strings.HasPrefix(result, "$Failed") {
					t.Errorf("Expected error, but got: %s", result)
				}
			} else {
				if result != test.expected {
					t.Errorf("Expected %s, got %s", test.expected, result)
				}
			}
		})
	}
}

func TestTableIteratorErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Table with invalid iterator - too few arguments",
			input: `Table(i, List(i))`,
		},
		{
			name:  "Table with invalid iterator - too many arguments",
			input: `Table(i, List(i, 1, 2, 3, 4, 5))`,
		},
		{
			name:  "Table with non-symbol iterator variable",
			input: `Table(i, List(123, 1, 5))`,
		},
		{
			name:  "Table with invalid List structure",
			input: `Table(i, List())`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			result := evaluateStringSimple(t, evaluator, test.input)

			if !strings.HasPrefix(result, "$Failed") {
				t.Errorf("Expected error, but got: %s", result)
			}
		})
	}
}

func TestDoSimple(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Do with simple expression 3 times",
			input:    `Do(Print(x), 3)`,
			expected: "Null",
		},
		{
			name:     "Do zero times",
			input:    `Do(Print("hello"), 0)`,
			expected: "Null",
		},
		{
			name:     "Do once",
			input:    `Do(a = 1, 1); a`,
			expected: "1",
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

func TestDoIterator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Do with iterator variable - basic range",
			input:    `Do(a = i, List(i, 1, 3)); a`,
			expected: "3",
		},
		{
			name:     "Do with iterator variable - step increment",
			input:    `total = 0; Do(total = Plus(total, i), List(i, 1, 10, 2)); total`,
			expected: "25", // 1 + 3 + 5 + 7 + 9 = 25
		},
		{
			name:     "Do with iterator variable - negative increment",
			input:    `last = 0; Do(last = i, List(i, 10, 5, -1)); last`,
			expected: "5",
		},
		{
			name:     "Do with iterator variable - zero increment returns Null",
			input:    `Do(a = i, List(i, 1, 5, 0))`,
			expected: "Null",
		},
		{
			name:     "Do with nested expressions",
			input:    `sum = 0; Do(Block(List(temp), temp = Times(i, 2); sum = Plus(sum, temp)), List(i, 1, 3)); sum`,
			expected: "12", // 2 + 4 + 6 = 12
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

func TestDoErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Do with invalid iterator - too few arguments",
			input: `Do(i, List(i))`,
		},
		{
			name:  "Do with invalid iterator - too many arguments",
			input: `Do(i, List(i, 1, 2, 3, 4, 5))`,
		},
		{
			name:  "Do with non-symbol iterator variable",
			input: `Do(i, List(123, 1, 5))`,
		},
		{
			name:  "Do with invalid List structure",
			input: `Do(i, List())`,
		},
		{
			name:  "Do with symbolic expressions that don't evaluate",
			input: `Do(x, List(x, n, 10*n, n))`,
		},
		{
			name:  "Do with negative count",
			input: `Do(Print("test"), -1)`,
		},
		{
			name:  "Do with wrong number of arguments",
			input: `Do()`,
		},
		{
			name:  "Do with too many arguments",
			input: `Do(x, 1, 2, 3)`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaluator := NewEvaluator()
			result := evaluateStringSimple(t, evaluator, test.input)

			if !strings.HasPrefix(result, "$Failed") {
				t.Errorf("Expected error, but got: %s", result)
			}
		})
	}
}
