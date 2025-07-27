package sexpr

import (
	"testing"
)

func TestCaretOperatorBasic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		// Basic power operations
		{
			name:     "Integer power",
			input:    "2 ^ 3",
			expected: "8.0",
			hasError: false,
		},
		{
			name:     "Real base, integer exponent",
			input:    "2.5 ^ 2",
			expected: "6.25",
			hasError: false,
		},
		{
			name:     "Integer to real power",
			input:    "4 ^ 0.5",
			expected: "2.0",
			hasError: false,
		},
		{
			name:     "Real to real power",
			input:    "2.0 ^ 3.0",
			expected: "8.0",
			hasError: false,
		},
		{
			name:     "Power of zero",
			input:    "0 ^ 5",
			expected: "0.0",
			hasError: false,
		},
		{
			name:     "Power of one",
			input:    "1 ^ 100",
			expected: "1.0",
			hasError: false,
		},
		{
			name:     "Zero exponent",
			input:    "5 ^ 0",
			expected: "1.0",
			hasError: false,
		},
		{
			name:     "Negative base, even exponent",
			input:    "(-2) ^ 2",
			expected: "4.0",
			hasError: false,
		},
		{
			name:     "Negative base, odd exponent",
			input:    "(-2) ^ 3",
			expected: "-8.0",
			hasError: false,
		},
		{
			name:     "Negative exponent",
			input:    "2 ^ (-2)",
			expected: "0.25",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()

			expr, err := ParseString(tt.input)
			if err != nil {
				if !tt.hasError {
					t.Fatalf("Parse error: %v", err)
				}
				return
			}

			result := evaluator.Evaluate(expr)

			if tt.hasError {
				if !IsError(result) {
					t.Errorf("Expected error, but got: %s", result.String())
				}
			} else {
				if IsError(result) {
					t.Errorf("Unexpected error: %s", result.String())
				} else if result.String() != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, result.String())
				}
			}
		})
	}
}

func TestCaretOperatorParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string // Expected AST structure
		hasError bool
	}{
		{
			name:     "Simple power parsing",
			input:    "a ^ b",
			expected: "Power(a, b)",
			hasError: false,
		},
		{
			name:     "Power with parentheses",
			input:    "(x + 1) ^ 2",
			expected: "Power(Plus(x, 1), 2)",
			hasError: false,
		},
		{
			name:     "Power with variables",
			input:    "x ^ y",
			expected: "Power(x, y)",
			hasError: false,
		},
		{
			name:     "Function call as base",
			input:    "f(x) ^ 2",
			expected: "Power(f(x), 2)",
			hasError: false,
		},
		{
			name:     "Function call as exponent",
			input:    "2 ^ g(x)",
			expected: "Power(2, g(x))",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				if !tt.hasError {
					t.Fatalf("Parse error: %v", err)
				}
				return
			}

			if tt.hasError {
				t.Errorf("Expected parse error, but parsing succeeded")
				return
			}

			result := expr.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestCaretOperatorRightAssociativity(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		desc     string
	}{
		{
			name:     "Triple power right associative",
			input:    "2 ^ 3 ^ 2",
			expected: "512.0",
			desc:     "Should parse as 2^(3^2) = 2^9 = 512, not (2^3)^2 = 8^2 = 64",
		},
		{
			name:     "Four powers right associative",
			input:    "2 ^ 2 ^ 2 ^ 2",
			expected: "65536.0",
			desc:     "Should parse as 2^(2^(2^2)) = 2^(2^4) = 2^16 = 65536",
		},
		{
			name:     "AST structure for triple power",
			input:    "a ^ b ^ c",
			expected: "Power(a, Power(b, c))",
			desc:     "AST should show right associativity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()

			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			if tt.name == "AST structure for triple power" {
				// Test AST structure directly
				result := expr.String()
				if result != tt.expected {
					t.Errorf("Expected %s, got %s (%s)", tt.expected, result, tt.desc)
				}
			} else {
				// Test evaluation result
				result := evaluator.Evaluate(expr)
				if IsError(result) {
					t.Errorf("Unexpected error: %s", result.String())
				} else if result.String() != tt.expected {
					t.Errorf("Expected %s, got %s (%s)", tt.expected, result.String(), tt.desc)
				}
			}
		})
	}
}

func TestCaretOperatorPrecedence(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		desc     string
	}{
		{
			name:     "Power has higher precedence than multiplication",
			input:    "2 * 3 ^ 2",
			expected: "18.0",
			desc:     "Should parse as 2 * (3^2) = 2 * 9 = 18, not (2*3)^2 = 6^2 = 36",
		},
		{
			name:     "Power has higher precedence than division",
			input:    "8 / 2 ^ 2",
			expected: "2.0",
			desc:     "Should parse as 8 / (2^2) = 8 / 4 = 2, not (8/2)^2 = 4^2 = 16",
		},
		{
			name:     "Power has higher precedence than addition",
			input:    "1 + 2 ^ 3",
			expected: "9.0",
			desc:     "Should parse as 1 + (2^3) = 1 + 8 = 9, not (1+2)^3 = 3^3 = 27",
		},
		{
			name:     "Power has higher precedence than subtraction",
			input:    "10 - 2 ^ 2",
			expected: "6.0",
			desc:     "Should parse as 10 - (2^2) = 10 - 4 = 6, not (10-2)^2 = 8^2 = 64",
		},
		{
			name:     "Multiple operators with power",
			input:    "2 + 3 * 4 ^ 2",
			expected: "50.0",
			desc:     "Should parse as 2 + (3 * (4^2)) = 2 + (3 * 16) = 2 + 48 = 50",
		},
		{
			name:     "Parentheses override precedence",
			input:    "(2 + 3) ^ 2",
			expected: "25.0",
			desc:     "Parentheses force addition first",
		},
		{
			name:     "Unary minus with power",
			input:    "-2 ^ 2",
			expected: "-4.0",
			desc:     "Should parse as -(2^2) = -4 with Mathematica precedence",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()

			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := evaluator.Evaluate(expr)

			if IsError(result) {
				t.Errorf("Unexpected error: %s", result.String())
			} else if result.String() != tt.expected {
				t.Errorf("Expected %s, got %s (%s)", tt.expected, result.String(), tt.desc)
			}
		})
	}
}

func TestCaretOperatorWithExistingOperators(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		// Verify existing operators still work
		{
			name:     "Multiplication still works",
			input:    "3 * 4",
			expected: "12",
			hasError: false,
		},
		{
			name:     "Addition still works",
			input:    "5 + 7",
			expected: "12",
			hasError: false,
		},
		{
			name:     "Complex expression with power",
			input:    "2 ^ 3 + 3 ^ 2",
			expected: "17.0",
			hasError: false,
		},
		{
			name:     "Power in compound statement",
			input:    "x = 2; y = x ^ 3; y",
			expected: "8.0",
			hasError: false,
		},
		{
			name:     "Power with comparison",
			input:    "2 ^ 3 == 8",
			expected: "True",
			hasError: false,
		},
		{
			name:     "Power with boolean operators",
			input:    "2 ^ 2 == 4 && 3 ^ 2 == 9",
			expected: "True",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()

			expr, err := ParseString(tt.input)
			if err != nil {
				if !tt.hasError {
					t.Fatalf("Parse error: %v", err)
				}
				return
			}

			result := evaluator.Evaluate(expr)

			if tt.hasError {
				if !IsError(result) {
					t.Errorf("Expected error, but got: %s", result.String())
				}
			} else {
				if IsError(result) {
					t.Errorf("Unexpected error: %s", result.String())
				} else if result.String() != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, result.String())
				}
			}
		})
	}
}

func TestLexerCaretToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []TokenType
	}{
		{
			name:     "Single caret token",
			input:    "^",
			expected: []TokenType{CARET, EOF},
		},
		{
			name:     "Caret between numbers",
			input:    "2^3",
			expected: []TokenType{INTEGER, CARET, INTEGER, EOF},
		},
		{
			name:     "Caret with spaces",
			input:    "2 ^ 3",
			expected: []TokenType{INTEGER, CARET, INTEGER, EOF},
		},
		{
			name:     "Multiple carets",
			input:    "2^3^4",
			expected: []TokenType{INTEGER, CARET, INTEGER, CARET, INTEGER, EOF},
		},
		{
			name:     "Caret with parentheses",
			input:    "2^(3+4)",
			expected: []TokenType{INTEGER, CARET, LPAREN, INTEGER, PLUS, INTEGER, RPAREN, EOF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			var tokens []TokenType

			for {
				tok := lexer.NextToken()
				tokens = append(tokens, tok.Type)
				if tok.Type == EOF {
					break
				}
			}

			if len(tokens) != len(tt.expected) {
				t.Errorf("Expected %d tokens, got %d", len(tt.expected), len(tokens))
				return
			}

			for i, expected := range tt.expected {
				if tokens[i] != expected {
					t.Errorf("Token %d: expected %v, got %v", i, expected, tokens[i])
				}
			}
		})
	}
}
