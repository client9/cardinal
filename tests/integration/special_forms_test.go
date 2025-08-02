package integration

import (
	"testing"
)

func TestIfConditions(t *testing.T) {
	tests := []TestCase{
		{
			name:     "If with True condition",
			input:    `If(True, "yes", "no")`,
			expected: `"yes"`,
		},
		{
			name:     "If with False condition",
			input:    `If(False, "yes", "no")`,
			expected: `"no"`,
		},
		{
			name:     "If without else clause (True)",
			input:    `If(True, "yes")`,
			expected: `"yes"`,
		},
		{
			name:     "If without else clause (False)",
			input:    `If(False, "yes")`,
			expected: `Null`,
		},
		{
			name:     "If with arithmetic condition",
			input:    `If(Equal(2, Plus(1, 1)), "equal", "not equal")`,
			expected: `"equal"`,
		},
		{
			name:     "Nested If statements",
			input:    `If(True, If(False, "inner true", "inner false"), "outer false")`,
			expected: `"inner false"`,
		},
	}

	runTestCases(t, tests)
}

func TestSetOperations(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Simple Set assignment",
			input:    `Set(x, 42); x`,
			expected: `42`,
		},
		{
			name:     "Set with arithmetic",
			input:    `Set(y, Plus(3, 4)); y`,
			expected: `7`,
		},
		{
			name:     "Set returns assigned value",
			input:    `Set(z, 100)`,
			expected: `100`,
		},
		{
			name:     "Multiple assignments",
			input:    `Set(a, 1); Set(b, 2); Plus(a, b)`,
			expected: `3`,
		},
	}

	runTestCases(t, tests)
}

func TestDoLoop(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Do with simple count",
			input:    `Do(Print("hello"), 0)`,
			expected: `Null`,
		},
		{
			name:     "Do returns Null",
			input:    `Do(Plus(1, 2), 3)`,
			expected: `Null`,
		},
		{
			name:     "Do with variable assignment",
			input:    `Set(counter, 0); Do(Set(counter, Plus(counter, 1)), 3); counter`,
			expected: `3`,
		},
	}

	runTestCases(t, tests)
}

func TestTableGeneration(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Table with simple count",
			input:    `Table(42, 3)`,
			expected: `List(42, 42, 42)`,
		},
		{
			name:     "Table with iterator variable",
			input:    `Table(i, List(i, 3))`,
			expected: `List(1, 2, 3)`,
		},
		{
			name:     "Table with range",
			input:    `Table(Times(i, 2), List(i, 1, 3))`,
			expected: `List(2, 4, 6)`,
		},
		{
			name:     "Table with arithmetic expression",
			input:    `Table(Plus(i, 10), List(i, 1, 3))`,
			expected: `List(11, 12, 13)`,
		},
		{
			name:     "Empty table",
			input:    `Table(x, 0)`,
			expected: `List()`,
		},
	}

	runTestCases(t, tests)
}

func TestHoldForms(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Hold prevents evaluation",
			input:    `Hold(Plus(1, 2))`,
			expected: `Hold(Plus(1, 2))`,
		},
		{
			name:     "Hold with multiple arguments",
			input:    `Hold(Plus(1, 2), Times(3, 4))`,
			expected: `Hold(Plus(1, 2), Times(3, 4))`,
		},
		{
			name:     "Nested Hold",
			input:    `Hold(Hold(Plus(1, 2)))`,
			expected: `Hold(Hold(Plus(1, 2)))`,
		},
	}

	runTestCases(t, tests)
}
