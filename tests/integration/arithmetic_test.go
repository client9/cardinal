package integration

import (
	"testing"
)

func TestBasicArithmetic(t *testing.T) {
	tests := []TestCase{
		// Addition
		{
			name:     "Simple addition",
			input:    "Plus(2, 3)",
			expected: "5",
		},
		{
			name:     "Multiple addition",
			input:    "Plus(1, 2, 3, 4)",
			expected: "10",
		},
		{
			name:     "Addition with reals",
			input:    "Plus(1.5, 2.5)",
			expected: "4.0",
		},
		{
			name:     "Mixed integer and real addition",
			input:    "Plus(1, 2.5)",
			expected: "3.5",
		},

		// Subtraction
		{
			name:     "Simple subtraction",
			input:    "Subtract(5, 3)",
			expected: "2",
		},
		{
			name:     "Subtraction with reals",
			input:    "Subtract(5.5, 2.5)",
			expected: "3.0",
		},

		// Multiplication
		{
			name:     "Simple multiplication",
			input:    "Times(3, 4)",
			expected: "12",
		},
		{
			name:     "Multiple multiplication",
			input:    "Times(2, 3, 4)",
			expected: "24",
		},
		{
			name:     "Multiplication with reals",
			input:    "Times(2.5, 4.0)",
			expected: "10.0",
		},

		// Division
		{
			name:     "Simple division",
			input:    "Divide(8, 2)",
			expected: "4",
		},
		{
			name:     "Division with reals",
			input:    "Divide(7.5, 2.5)",
			expected: "3.0",
		},
		{
			name:     "Integer division with remainder",
			input:    "Divide(7, 2)",
			expected: "3",
		},

		// Power
		{
			name:     "Simple power",
			input:    "Power(2, 3)",
			expected: "8",
		},
		{
			name:     "Simple power",
			input:    "Power(2.0, 3.0)",
			expected: "8.0",
		},
		{
			name:     "Simple power",
			input:    "Power(2, 3.0)",
			expected: "8.0",
		},
		{
			name:     "Power with real base",
			input:    "Power(2.5, 2)",
			expected: "6.25",
		},
		{
			name:     "Power of integer zero",
			input:    "Power(5, 0)",
			expected: "1",
		},
		{
			name:     "Power of real zero",
			input:    "Power(5, 0.0)",
			expected: "1.0",
		},
		{
			name:     "Power of one",
			input:    "Power(7, 1)",
			expected: "7",
		},
		{
			name:     "Nested arithmetic",
			input:    "Plus(Times(2, 3), Divide(8, 2))",
			expected: "10",
		},
		{
			name:     "Complex expression",
			input:    "Times(Plus(1, 2), Subtract(5, 2))",
			expected: "9",
		},
		{
			name:     "Power in arithmetic",
			input:    "Plus(Power(2, 3), Times(3, 4))",
			expected: "20",
		},
		{
			name:     "Deeply nested",
			input:    "Plus(1, Times(2, Plus(3, Times(4, 5))))",
			expected: "47",
		},
	}

	runTestCases(t, tests)
}

func TestArithmeticIdentities(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Plus identity (empty)",
			input:    "Plus()",
			expected: "0",
		},
		{
			name:     "Times identity (empty)",
			input:    "Times()",
			expected: "1",
		},
		{
			name:     "Plus with single argument",
			input:    "Plus(42)",
			expected: "42",
		},
		{
			name:     "Times with single argument",
			input:    "Times(42)",
			expected: "42",
		},
	}

	runTestCases(t, tests)
}

func TestArithmeticAttributes(t *testing.T) {
	tests := []TestCase{
		// Test Orderless attribute (commutativity)
		{
			name:     "Plus orderless",
			input:    "Plus(3, 1, 2)",
			expected: "6", // Should work regardless of order
		},
		{
			name:     "Times orderless",
			input:    "Times(4, 1, 3, 2)",
			expected: "24", // Should work regardless of order
		},

		// Test Flat attribute (associativity)
		{
			name:     "Plus flat",
			input:    "Plus(1, Plus(2, 3), 4)",
			expected: "10", // Should flatten to Plus(1, 2, 3, 4)
		},
		{
			name:     "Times flat",
			input:    "Times(2, Times(3, 4))",
			expected: "24", // Should flatten to Times(2, 3, 4)
		},
	}

	runTestCases(t, tests)
}

func TestArithmeticErrors(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		errorType string
	}{
		{
			name:      "Division by zero",
			input:     "Divide(1, 0)",
			errorType: "DivisionByZero",
		},
	}

	runErrorTestCases(t, tests)
}

func TestUnaryMinus(t *testing.T) {
	tests := []TestCase{
		{
			name:     "Minus integer",
			input:    "Minus(5)",
			expected: "-5",
		},
		{
			name:     "Minus real",
			input:    "Minus(3.14)",
			expected: "-3.14",
		},
		{
			name:     "Double negative",
			input:    "Minus(Minus(7))",
			expected: "7",
		},
		{
			name:     "Minus in arithmetic",
			input:    "Plus(Minus(3), 5)",
			expected: "2",
		},
	}

	runTestCases(t, tests)
}
