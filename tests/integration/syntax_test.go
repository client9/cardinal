package integration

import (
	"testing"
)

func TestCompoundStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple two expressions",
			input:    "1 + 2; 3 + 4",
			expected: "7",
		},
		{
			name:     "Sequential assignments",
			input:    "a = 5; b = 10; a + b",
			expected: "15",
		},
		{
			name:     "Chained variable dependencies",
			input:    "x = 2; y = x * 3; z = y + 1; z",
			expected: "7",
		},
		{
			name:     "Multiple simple values",
			input:    "1; 2; 3; 4",
			expected: "4",
		},
		{
			name:     "Assignment with function call",
			input:    "list = [1, 2, 3]; len = Length(list); len",
			expected: "3",
		},
		{
			name:     "Empty compound statement",
			input:    ";",
			expected: "Null",
		},
		{
			name:     "Trailing semicolon",
			input:    "42;",
			expected: "Null",
		},
		{
			name:     "Leading semicolon",
			input:    ";42",
			expected: "42",
		},
	}

	runTestCases(t, tests)
}

func TestCaretOperatorBasic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic power operations
		{
			name:     "Integer power",
			input:    "2 ^ 3",
			expected: "8",
		},
		{
			name:     "Real base, integer exponent",
			input:    "2.5 ^ 2",
			expected: "6.25",
		},
		{
			name:     "Zero exponent",
			input:    "5 ^ 0",
			expected: "1",
		},
		{
			name:     "One exponent",
			input:    "7 ^ 1",
			expected: "7",
		},
		{
			name:     "Base one",
			input:    "1 ^ 100",
			expected: "1",
		},
		{
			name:     "Base zero",
			input:    "0 ^ 5",
			expected: "0",
		},
		{
			name:     "Negative base positive exponent",
			input:    "(-2) ^ 3",
			expected: "-8",
		},
	}

	runTestCases(t, tests)
}

func TestCaretOperatorPrecedence(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Right associative",
			input:    "2 ^ 3 ^ 2",
			expected: "512", // 2^(3^2) = 2^9 = 512
		},
		{
			name:     "Higher than multiplication",
			input:    "2 * 3 ^ 2",
			expected: "18", // 2 * (3^2) = 2 * 9 = 18
		},
		{
			name:     "Higher than addition",
			input:    "1 + 2 ^ 3",
			expected: "9", // 1 + (2^3) = 1 + 8 = 9
		},
		{
			name:     "With parentheses override",
			input:    "(2 + 3) ^ 2",
			expected: "25", // (5)^2 = 25
		},
		{
			name:     "Complex expression",
			input:    "2 + 3 * 4 ^ 2 - 1",
			expected: "49", // 2 + 3 * 16 - 1 = 2 + 48 - 1 = 49
		},
	}

	runTestCases(t, tests)
}

func TestUnaryOperators(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Unary minus integer",
			input:    "-5",
			expected: "-5",
		},
		{
			name:     "Unary plus integer",
			input:    "+5",
			expected: "5",
		},
		{
			name:     "Unary minus with expression",
			input:    "-(2 + 3)",
			expected: "-5",
		},
		{
			name:     "Double negative",
			input:    "-(-5)",
			expected: "5",
		},
		{
			name:     "Unary with multiplication",
			input:    "-2 * 3",
			expected: "-6",
		},
	}

	runTestCases(t, tests)
}

func TestSlicingSyntax(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "List indexing",
			input:    "[1, 2, 3, 4][2]",
			expected: "2",
		},
		{
			name:     "List slicing range",
			input:    "[1, 2, 3, 4, 5][2:4]",
			expected: "List(2, 3, 4)",
		},
		{
			name:     "String indexing",
			input:    "\"hello\"[1]",
			expected: "\"h\"",
		},
		{
			name:     "String slicing",
			input:    "\"hello\"[2:4]",
			expected: "\"ell\"",
		},
		{
			name:     "Negative indexing",
			input:    "[1, 2, 3][-1]",
			expected: "3",
		},
	}

	runTestCases(t, tests)
}
