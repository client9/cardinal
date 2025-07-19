package sexpr

import (
	"testing"
)

func TestLexer_NextToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "symbols",
			input: "Plus Minus x y",
			expected: []Token{
				{Type: SYMBOL, Value: "Plus"},
				{Type: SYMBOL, Value: "Minus"},
				{Type: SYMBOL, Value: "x"},
				{Type: SYMBOL, Value: "y"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "integers",
			input: "42 0 123",
			expected: []Token{
				{Type: INTEGER, Value: "42"},
				{Type: INTEGER, Value: "0"},
				{Type: INTEGER, Value: "123"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "floats",
			input: "3.14 0.5 123.456",
			expected: []Token{
				{Type: FLOAT, Value: "3.14"},
				{Type: FLOAT, Value: "0.5"},
				{Type: FLOAT, Value: "123.456"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "strings",
			input: `"hello" "world" "test string"`,
			expected: []Token{
				{Type: STRING, Value: "hello"},
				{Type: STRING, Value: "world"},
				{Type: STRING, Value: "test string"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "escaped strings",
			input: `"hello\nworld" "tab\there" "quote\"test"`,
			expected: []Token{
				{Type: STRING, Value: "hello\\nworld"},
				{Type: STRING, Value: "tab\\there"},
				{Type: STRING, Value: "quote\\\"test"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "booleans",
			input: "True False",
			expected: []Token{
				{Type: BOOLEAN, Value: "True"},
				{Type: BOOLEAN, Value: "False"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "brackets and braces",
			input: "[]{}",
			expected: []Token{
				{Type: LBRACKET, Value: "["},
				{Type: RBRACKET, Value: "]"},
				{Type: LBRACE, Value: "{"},
				{Type: RBRACE, Value: "}"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "comma",
			input: "a, b, c",
			expected: []Token{
				{Type: SYMBOL, Value: "a"},
				{Type: COMMA, Value: ","},
				{Type: SYMBOL, Value: "b"},
				{Type: COMMA, Value: ","},
				{Type: SYMBOL, Value: "c"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "mathematical operators",
			input: "a + b - c * d / e",
			expected: []Token{
				{Type: SYMBOL, Value: "a"},
				{Type: PLUS, Value: "+"},
				{Type: SYMBOL, Value: "b"},
				{Type: MINUS, Value: "-"},
				{Type: SYMBOL, Value: "c"},
				{Type: MULTIPLY, Value: "*"},
				{Type: SYMBOL, Value: "d"},
				{Type: DIVIDE, Value: "/"},
				{Type: SYMBOL, Value: "e"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "set operator",
			input: "x = 5",
			expected: []Token{
				{Type: SYMBOL, Value: "x"},
				{Type: SET, Value: "="},
				{Type: INTEGER, Value: "5"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "setdelayed operator",
			input: "f := g[x]",
			expected: []Token{
				{Type: SYMBOL, Value: "f"},
				{Type: SETDELAYED, Value: ":="},
				{Type: SYMBOL, Value: "g"},
				{Type: LBRACKET, Value: "["},
				{Type: SYMBOL, Value: "x"},
				{Type: RBRACKET, Value: "]"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "unset operator",
			input: "x =.",
			expected: []Token{
				{Type: SYMBOL, Value: "x"},
				{Type: UNSET, Value: "=."},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "comparison operators",
			input: "x == y != z < a > b <= c >= d",
			expected: []Token{
				{Type: SYMBOL, Value: "x"},
				{Type: EQUAL, Value: "=="},
				{Type: SYMBOL, Value: "y"},
				{Type: UNEQUAL, Value: "!="},
				{Type: SYMBOL, Value: "z"},
				{Type: LESS, Value: "<"},
				{Type: SYMBOL, Value: "a"},
				{Type: GREATER, Value: ">"},
				{Type: SYMBOL, Value: "b"},
				{Type: LESSEQUAL, Value: "<="},
				{Type: SYMBOL, Value: "c"},
				{Type: GREATEREQUAL, Value: ">="},
				{Type: SYMBOL, Value: "d"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "sameq operator",
			input: "x === y",
			expected: []Token{
				{Type: SYMBOL, Value: "x"},
				{Type: SAMEQ, Value: "==="},
				{Type: SYMBOL, Value: "y"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "unsameq operator",
			input: "x =!= y",
			expected: []Token{
				{Type: SYMBOL, Value: "x"},
				{Type: UNSAMEQ, Value: "=!="},
				{Type: SYMBOL, Value: "y"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "logical operators",
			input: "x && y || z",
			expected: []Token{
				{Type: SYMBOL, Value: "x"},
				{Type: AND, Value: "&&"},
				{Type: SYMBOL, Value: "y"},
				{Type: OR, Value: "||"},
				{Type: SYMBOL, Value: "z"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "function call",
			input: "Plus[1, 2, 3]",
			expected: []Token{
				{Type: SYMBOL, Value: "Plus"},
				{Type: LBRACKET, Value: "["},
				{Type: INTEGER, Value: "1"},
				{Type: COMMA, Value: ","},
				{Type: INTEGER, Value: "2"},
				{Type: COMMA, Value: ","},
				{Type: INTEGER, Value: "3"},
				{Type: RBRACKET, Value: "]"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "nested function",
			input: `Function["x", Power[x, 2]]`,
			expected: []Token{
				{Type: SYMBOL, Value: "Function"},
				{Type: LBRACKET, Value: "["},
				{Type: STRING, Value: "x"},
				{Type: COMMA, Value: ","},
				{Type: SYMBOL, Value: "Power"},
				{Type: LBRACKET, Value: "["},
				{Type: SYMBOL, Value: "x"},
				{Type: COMMA, Value: ","},
				{Type: INTEGER, Value: "2"},
				{Type: RBRACKET, Value: "]"},
				{Type: RBRACKET, Value: "]"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "empty list",
			input: "{}",
			expected: []Token{
				{Type: LBRACE, Value: "{"},
				{Type: RBRACE, Value: "}"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "whitespace handling",
			input: " \t\n\r Plus \t\n [ \r 1 \n ] ",
			expected: []Token{
				{Type: SYMBOL, Value: "Plus"},
				{Type: LBRACKET, Value: "["},
				{Type: INTEGER, Value: "1"},
				{Type: RBRACKET, Value: "]"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "illegal characters",
			input: "Plus @ 123",
			expected: []Token{
				{Type: SYMBOL, Value: "Plus"},
				{Type: ILLEGAL, Value: "@"},
				{Type: INTEGER, Value: "123"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "unclosed string",
			input: `"unclosed string`,
			expected: []Token{
				{Type: STRING, Value: "unclosed string"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "empty input",
			input: "",
			expected: []Token{
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "complex expression",
			input: `List[1, 2.5, "hello", True, Plus[x, y]]`,
			expected: []Token{
				{Type: SYMBOL, Value: "List"},
				{Type: LBRACKET, Value: "["},
				{Type: INTEGER, Value: "1"},
				{Type: COMMA, Value: ","},
				{Type: FLOAT, Value: "2.5"},
				{Type: COMMA, Value: ","},
				{Type: STRING, Value: "hello"},
				{Type: COMMA, Value: ","},
				{Type: BOOLEAN, Value: "True"},
				{Type: COMMA, Value: ","},
				{Type: SYMBOL, Value: "Plus"},
				{Type: LBRACKET, Value: "["},
				{Type: SYMBOL, Value: "x"},
				{Type: COMMA, Value: ","},
				{Type: SYMBOL, Value: "y"},
				{Type: RBRACKET, Value: "]"},
				{Type: RBRACKET, Value: "]"},
				{Type: EOF, Value: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)

			for i, expectedToken := range tt.expected {
				token := lexer.NextToken()

				if token.Type != expectedToken.Type {
					t.Errorf("test[%d] - token type wrong. expected=%v, got=%v", i, expectedToken.Type, token.Type)
				}

				if token.Value != expectedToken.Value {
					t.Errorf("test[%d] - token value wrong. expected=%q, got=%q", i, expectedToken.Value, token.Value)
				}
			}
		})
	}
}

func TestLexer_Tokenize(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedCount int
	}{
		{
			name:          "simple expression",
			input:         "Plus[1, 2]",
			expectedCount: 7, // Plus, [, 1, ,, 2, ], EOF
		},
		{
			name:          "empty input",
			input:         "",
			expectedCount: 1, // EOF
		},
		{
			name:          "complex expression",
			input:         `Function["x", Power[x, 2]]`,
			expectedCount: 12, // Function, [, "x", ,, Power, [, x, ,, 2, ], ], EOF
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens := lexer.Tokenize()

			if len(tokens) != tt.expectedCount {
				t.Errorf("expected %d tokens, got %d", tt.expectedCount, len(tokens))
			}

			// Last token should always be EOF
			if len(tokens) > 0 && tokens[len(tokens)-1].Type != EOF {
				t.Errorf("expected last token to be EOF, got %v", tokens[len(tokens)-1].Type)
			}
		})
	}
}

func TestToken_String(t *testing.T) {
	tests := []struct {
		name     string
		token    Token
		expected string
	}{
		{
			name:     "symbol token",
			token:    Token{Type: SYMBOL, Value: "Plus"},
			expected: "SYMBOL(Plus)",
		},
		{
			name:     "integer token",
			token:    Token{Type: INTEGER, Value: "42"},
			expected: "INTEGER(42)",
		},
		{
			name:     "float token",
			token:    Token{Type: FLOAT, Value: "3.14"},
			expected: "FLOAT(3.14)",
		},
		{
			name:     "string token",
			token:    Token{Type: STRING, Value: "hello"},
			expected: "STRING(hello)",
		},
		{
			name:     "boolean token",
			token:    Token{Type: BOOLEAN, Value: "True"},
			expected: "BOOLEAN(True)",
		},
		{
			name:     "left bracket",
			token:    Token{Type: LBRACKET, Value: "["},
			expected: "LBRACKET",
		},
		{
			name:     "right bracket",
			token:    Token{Type: RBRACKET, Value: "]"},
			expected: "RBRACKET",
		},
		{
			name:     "left brace",
			token:    Token{Type: LBRACE, Value: "{"},
			expected: "LBRACE",
		},
		{
			name:     "right brace",
			token:    Token{Type: RBRACE, Value: "}"},
			expected: "RBRACE",
		},
		{
			name:     "comma",
			token:    Token{Type: COMMA, Value: ","},
			expected: "COMMA",
		},
		{
			name:     "EOF",
			token:    Token{Type: EOF, Value: ""},
			expected: "EOF",
		},
		{
			name:     "illegal token",
			token:    Token{Type: ILLEGAL, Value: "@"},
			expected: "ILLEGAL(@)",
		},
		{
			name:     "plus token",
			token:    Token{Type: PLUS, Value: "+"},
			expected: "PLUS",
		},
		{
			name:     "minus token",
			token:    Token{Type: MINUS, Value: "-"},
			expected: "MINUS",
		},
		{
			name:     "multiply token",
			token:    Token{Type: MULTIPLY, Value: "*"},
			expected: "MULTIPLY",
		},
		{
			name:     "divide token",
			token:    Token{Type: DIVIDE, Value: "/"},
			expected: "DIVIDE",
		},
		{
			name:     "set token",
			token:    Token{Type: SET, Value: "="},
			expected: "SET",
		},
		{
			name:     "setdelayed token",
			token:    Token{Type: SETDELAYED, Value: ":="},
			expected: "SETDELAYED",
		},
		{
			name:     "unset token",
			token:    Token{Type: UNSET, Value: "=."},
			expected: "UNSET",
		},
		{
			name:     "equal token",
			token:    Token{Type: EQUAL, Value: "=="},
			expected: "EQUAL",
		},
		{
			name:     "unequal token",
			token:    Token{Type: UNEQUAL, Value: "!="},
			expected: "UNEQUAL",
		},
		{
			name:     "less token",
			token:    Token{Type: LESS, Value: "<"},
			expected: "LESS",
		},
		{
			name:     "greater token",
			token:    Token{Type: GREATER, Value: ">"},
			expected: "GREATER",
		},
		{
			name:     "lessequal token",
			token:    Token{Type: LESSEQUAL, Value: "<="},
			expected: "LESSEQUAL",
		},
		{
			name:     "greaterequal token",
			token:    Token{Type: GREATEREQUAL, Value: ">="},
			expected: "GREATEREQUAL",
		},
		{
			name:     "and token",
			token:    Token{Type: AND, Value: "&&"},
			expected: "AND",
		},
		{
			name:     "or token",
			token:    Token{Type: OR, Value: "||"},
			expected: "OR",
		},
		{
			name:     "sameq token",
			token:    Token{Type: SAMEQ, Value: "==="},
			expected: "SAMEQ",
		},
		{
			name:     "unsameq token",
			token:    Token{Type: UNSAMEQ, Value: "=!="},
			expected: "UNSAMEQ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.token.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
