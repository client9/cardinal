package sexpr

import (
	"testing"
)

func TestPrefixNotOperator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		// Basic NOT operator tests
		{
			name:     "NOT with True",
			input:    "!True",
			expected: "False",
			hasError: false,
		},
		{
			name:     "NOT with False",
			input:    "!False",
			expected: "True",
			hasError: false,
		},
		{
			name:     "NOT with boolean expression",
			input:    "!(1 == 1)",
			expected: "False",
			hasError: false,
		},
		{
			name:     "NOT with false boolean expression",
			input:    "!(1 == 2)",
			expected: "True",
			hasError: false,
		},
		{
			name:     "NOT with variable",
			input:    "x = True; !x",
			expected: "False",
			hasError: false,
		},
		{
			name:     "NOT with complex expression",
			input:    "!(True && False)",
			expected: "True",
			hasError: false,
		},
		{
			name:     "Double NOT",
			input:    "!!True",
			expected: "True",
			hasError: false,
		},
		{
			name:     "Double NOT with False",
			input:    "!!False",
			expected: "False",
			hasError: false,
		},
		{
			name:     "NOT with comparison",
			input:    "!(5 > 3)",
			expected: "False",
			hasError: false,
		},
		{
			name:     "NOT with inequality",
			input:    "!(2 != 2)",
			expected: "True",
			hasError: false,
		},
		{
			name:     "NOT with grouped expression",
			input:    "!((1 < 2) && (3 > 4))",
			expected: "True",
			hasError: false,
		},

		// Precedence tests
		{
			name:     "NOT precedence with AND",
			input:    "!True && False",
			expected: "False",
			hasError: false,
		},
		{
			name:     "NOT precedence with OR",
			input:    "!False || True",
			expected: "True",
			hasError: false,
		},
		{
			name:     "NOT with arithmetic expression",
			input:    "!(1 + 2 == 3)",
			expected: "False",
			hasError: false,
		},

		// Mixed with other operators
		{
			name:     "NOT with assignment",
			input:    "result = !True; result",
			expected: "False",
			hasError: false,
		},
		{
			name:     "NOT in compound statement",
			input:    "a = True; b = !a; b",
			expected: "False",
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

func TestPrefixNotOperatorParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string // Expected AST structure
		hasError bool
	}{
		{
			name:     "Simple NOT parsing",
			input:    "!True",
			expected: "Not(True)",
			hasError: false,
		},
		{
			name:     "NOT with parentheses",
			input:    "!(True)",
			expected: "Not(True)",
			hasError: false,
		},
		{
			name:     "NOT with expression",
			input:    "!(1 == 2)",
			expected: "Not(Equal(1, 2))",
			hasError: false,
		},
		{
			name:     "Double NOT",
			input:    "!!True",
			expected: "Not(Not(True))",
			hasError: false,
		},
		{
			name:     "NOT with variable",
			input:    "!x",
			expected: "Not(x)",
			hasError: false,
		},
		{
			name:     "NOT with function call",
			input:    "!f(x)",
			expected: "Not(f(x))",
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

func TestPrefixNotOperatorWithExistingOperators(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		// Verify existing unary operators still work
		{
			name:     "Unary minus still works",
			input:    "-5",
			expected: "-5",
			hasError: false,
		},
		{
			name:     "Unary plus still works",
			input:    "+5",
			expected: "5",
			hasError: false,
		},
		{
			name:     "Unary minus with expression",
			input:    "-(2 + 3)",
			expected: "-5",
			hasError: false,
		},

		// Mix NOT with other unary operators
		{
			name:     "NOT with negative number",
			input:    "!(-1 == -1)",
			expected: "False",
			hasError: false,
		},
		{
			name:     "NOT with positive number",
			input:    "!(+1 == 1)",
			expected: "False",
			hasError: false,
		},

		// Verify != still works (should not be confused with ! and =)
		{
			name:     "Inequality operator still works",
			input:    "1 != 2",
			expected: "True",
			hasError: false,
		},
		{
			name:     "Inequality vs NOT",
			input:    "!(1 == 2)",
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

func TestNotOperatorPrecedence(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		desc     string
	}{
		{
			name:     "NOT has higher precedence than AND",
			input:    "!True && False",
			expected: "False",
			desc:     "Should parse as (!True) && False, not !(True && False)",
		},
		{
			name:     "NOT has higher precedence than OR",
			input:    "!False || False",
			expected: "True",
			desc:     "Should parse as (!False) || False, not !(False || False)",
		},
		{
			name:     "NOT with parentheses overrides precedence",
			input:    "!(True && False)",
			expected: "True",
			desc:     "Parentheses force evaluation of AND first",
		},
		{
			name:     "Multiple NOTs are right-associative",
			input:    "!!!True",
			expected: "False",
			desc:     "Should parse as !(!(!True))",
		},
		{
			name:     "NOT with comparison",
			input:    "!1 == 2",
			expected: "False",
			desc:     "Should parse as (!1) == 2, not !(1 == 2)",
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

func TestLexerNotToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []TokenType
	}{
		{
			name:     "Single NOT token",
			input:    "!",
			expected: []TokenType{NOT, EOF},
		},
		{
			name:     "NOT followed by symbol",
			input:    "!True",
			expected: []TokenType{NOT, SYMBOL, EOF},
		},
		{
			name:     "NOT vs UNEQUAL",
			input:    "! !=",
			expected: []TokenType{NOT, UNEQUAL, EOF},
		},
		{
			name:     "Multiple NOTs",
			input:    "!!!",
			expected: []TokenType{NOT, NOT, NOT, EOF},
		},
		{
			name:     "NOT with parentheses",
			input:    "!()",
			expected: []TokenType{NOT, LPAREN, RPAREN, EOF},
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
