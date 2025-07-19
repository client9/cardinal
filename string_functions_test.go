package sexpr

import (
	"testing"
)

func TestEvaluateStringQ(t *testing.T) {
	tests := []struct {
		name     string
		arg      Expr
		expected bool
	}{
		{
			name:     "String atom",
			arg:      NewStringAtom("hello"),
			expected: true,
		},
		{
			name:     "Empty string",
			arg:      NewStringAtom(""),
			expected: true,
		},
		{
			name:     "Long string",
			arg:      NewStringAtom("This is a longer string with spaces and punctuation!"),
			expected: true,
		},
		{
			name:     "String with escape sequences",
			arg:      NewStringAtom("Hello\nWorld\t!"),
			expected: true,
		},
		{
			name:     "Integer atom",
			arg:      NewIntAtom(42),
			expected: false,
		},
		{
			name:     "Float atom",
			arg:      NewFloatAtom(3.14),
			expected: false,
		},
		{
			name:     "Boolean atom",
			arg:      NewBoolAtom(true),
			expected: false,
		},
		{
			name:     "Symbol atom",
			arg:      NewSymbolAtom("x"),
			expected: false,
		},
		{
			name: "List",
			arg: List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewIntAtom(1),
				NewIntAtom(2),
			}},
			expected: false,
		},
		{
			name:     "Empty list",
			arg:      List{Elements: []Expr{}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateStringQ([]Expr{tt.arg})

			if !isBool(result) {
				t.Errorf("expected boolean result, got %T", result)
				return
			}

			val, _ := getBoolValue(result)
			if val != tt.expected {
				t.Errorf("expected %t, got %t", tt.expected, val)
			}
		})
	}
}

func TestEvaluateStringLength(t *testing.T) {
	tests := []struct {
		name     string
		arg      Expr
		expected int
		hasError bool
	}{
		{
			name:     "Empty string",
			arg:      NewStringAtom(""),
			expected: 0,
			hasError: false,
		},
		{
			name:     "Single character",
			arg:      NewStringAtom("a"),
			expected: 1,
			hasError: false,
		},
		{
			name:     "Simple string",
			arg:      NewStringAtom("hello"),
			expected: 5,
			hasError: false,
		},
		{
			name:     "String with spaces",
			arg:      NewStringAtom("hello world"),
			expected: 11,
			hasError: false,
		},
		{
			name:     "String with special characters",
			arg:      NewStringAtom("Hello, World! 123"),
			expected: 17,
			hasError: false,
		},
		{
			name:     "String with escape sequences",
			arg:      NewStringAtom("Line1\nLine2\tTabbed"),
			expected: 18, // \n and \t count as single characters
			hasError: false,
		},
		{
			name:     "Unicode string",
			arg:      NewStringAtom("Hello 世界"),
			expected: 8, // Unicode characters count correctly
			hasError: false,
		},
		{
			name:     "Integer atom - should error",
			arg:      NewIntAtom(42),
			expected: 0,
			hasError: true,
		},
		{
			name:     "Float atom - should error",
			arg:      NewFloatAtom(3.14),
			expected: 0,
			hasError: true,
		},
		{
			name:     "Boolean atom - should error",
			arg:      NewBoolAtom(true),
			expected: 0,
			hasError: true,
		},
		{
			name:     "Symbol atom - should error",
			arg:      NewSymbolAtom("x"),
			expected: 0,
			hasError: true,
		},
		{
			name: "List - should error",
			arg: List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewIntAtom(1),
			}},
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateStringLength([]Expr{tt.arg})

			if tt.hasError {
				if !IsError(result) {
					t.Errorf("expected error for %s, got %s", tt.name, result.String())
				}
				return
			}

			if IsError(result) {
				t.Errorf("unexpected error: %s", result.String())
				return
			}

			if !isNumeric(result) {
				t.Errorf("expected numeric result, got %T", result)
				return
			}

			val, _ := getNumericValue(result)
			if int(val) != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, int(val))
			}
		})
	}
}

func TestEvaluateFullForm(t *testing.T) {
	tests := []struct {
		name     string
		arg      Expr
		expected string
	}{
		{
			name:     "Integer atom",
			arg:      NewIntAtom(42),
			expected: "42",
		},
		{
			name:     "Float atom",
			arg:      NewFloatAtom(3.14),
			expected: "3.14",
		},
		{
			name:     "String atom",
			arg:      NewStringAtom("hello"),
			expected: "\"hello\"",
		},
		{
			name:     "Boolean true",
			arg:      NewBoolAtom(true),
			expected: "True",
		},
		{
			name:     "Boolean false",
			arg:      NewBoolAtom(false),
			expected: "False",
		},
		{
			name:     "Symbol atom",
			arg:      NewSymbolAtom("x"),
			expected: "x",
		},
		{
			name: "Simple function call",
			arg: List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewIntAtom(1),
				NewIntAtom(2),
			}},
			expected: "Plus(1, 2)",
		},
		{
			name:     "Empty list",
			arg:      List{Elements: []Expr{}},
			expected: "List()",
		},
		{
			name: "List literal",
			arg: List{Elements: []Expr{
				NewSymbolAtom("List"),
				NewIntAtom(1),
				NewIntAtom(2),
				NewIntAtom(3),
			}},
			expected: "List(1, 2, 3)",
		},
		{
			name: "Nested expression",
			arg: List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewIntAtom(1),
				List{Elements: []Expr{
					NewSymbolAtom("Times"),
					NewIntAtom(2),
					NewIntAtom(3),
				}},
			}},
			expected: "Plus(1, Times(2, 3))",
		},
		{
			name: "Complex expression",
			arg: List{Elements: []Expr{
				NewSymbolAtom("Equal"),
				List{Elements: []Expr{
					NewSymbolAtom("Plus"),
					NewSymbolAtom("x"),
					NewIntAtom(1),
				}},
				NewIntAtom(5),
			}},
			expected: "Equal(Plus(x, 1), 5)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateFullForm([]Expr{tt.arg})

			if IsError(result) {
				t.Errorf("unexpected error: %s", result.String())
				return
			}

			// FullForm should return a string atom
			if atom, ok := result.(Atom); ok && atom.AtomType == StringAtom {
				val := atom.Value.(string)
				if val != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, val)
				}
			} else {
				t.Errorf("expected string atom, got %T", result)
			}
		})
	}
}

// Test argument validation for all string functions
func TestStringFunctions_ArgumentValidation(t *testing.T) {
	functions := []struct {
		name string
		fn   func([]Expr) Expr
	}{
		{"StringQ", EvaluateStringQ},
		{"StringLength", EvaluateStringLength},
		{"FullForm", EvaluateFullForm},
	}

	for _, fn := range functions {
		t.Run(fn.name+"_no_args", func(t *testing.T) {
			result := fn.fn([]Expr{})
			if !IsError(result) {
				t.Errorf("expected error for no arguments, got %s", result.String())
			}
		})

		t.Run(fn.name+"_too_many_args", func(t *testing.T) {
			result := fn.fn([]Expr{NewIntAtom(1), NewIntAtom(2)})
			if !IsError(result) {
				t.Errorf("expected error for too many arguments, got %s", result.String())
			}
		})
	}
}

// Integration tests with the evaluator
func TestStringFunctions_Integration(t *testing.T) {
	eval := setupTestEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// StringQ tests
		{
			name:     "StringQ on string",
			input:    "StringQ(\"hello\")",
			expected: "True",
		},
		{
			name:     "StringQ on integer",
			input:    "StringQ(42)",
			expected: "False",
		},
		{
			name:     "StringQ on symbol",
			input:    "StringQ(x)",
			expected: "False",
		},
		{
			name:     "StringQ on list",
			input:    "StringQ([1, 2, 3])",
			expected: "False",
		},

		// StringLength tests
		{
			name:     "StringLength empty string",
			input:    "StringLength(\"\")",
			expected: "0",
		},
		{
			name:     "StringLength simple string",
			input:    "StringLength(\"hello\")",
			expected: "5",
		},
		{
			name:     "StringLength string with spaces",
			input:    "StringLength(\"hello world\")",
			expected: "11",
		},
		{
			name:     "StringLength on non-string",
			input:    "StringLength(42)",
			expected: "$Failed(ArgumentError)",
		},

		// FullForm tests
		{
			name:     "FullForm integer",
			input:    "FullForm(42)",
			expected: "\"42\"",
		},
		{
			name:     "FullForm string",
			input:    "FullForm(\"hello\")",
			expected: "\"\"hello\"\"",
		},
		{
			name:     "FullForm symbol",
			input:    "FullForm(x)",
			expected: "\"x\"",
		},
		{
			name:     "FullForm evaluated expression",
			input:    "FullForm(Plus(1, 2))",
			expected: "\"3\"", // Plus(1,2) evaluates to 3, then FullForm(3) = "3"
		},
		{
			name:     "FullForm unevaluated expression",
			input:    "FullForm(Hold(Plus(1, 2)))",
			expected: "\"Hold(Plus(1, 2))\"",
		},
		{
			name:     "FullForm list literal",
			input:    "FullForm([1, 2, 3])",
			expected: "\"List(1, 2, 3)\"",
		},

		// Combined tests
		{
			name:     "StringLength of FullForm",
			input:    "StringLength(FullForm(42))",
			expected: "2", // "42" has 2 characters
		},
		{
			name:     "StringQ of FullForm",
			input:    "StringQ(FullForm(42))",
			expected: "True",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := eval.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

// Test FullForm with complex expressions
func TestFullForm_ComplexExpressions(t *testing.T) {
	eval := setupTestEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "FullForm of unevaluated expression",
			input:    "FullForm(Hold(1 + 2))",
			expected: "\"Hold(Plus(1, 2))\"",
		},
		{
			name:     "FullForm of nested list",
			input:    "FullForm([[1, 2], [3, 4]])",
			expected: "\"List(List(1, 2), List(3, 4))\"",
		},
		{
			name:     "FullForm of evaluated comparison",
			input:    "FullForm(Equal(x, y))",
			expected: "\"False\"", // Equal(x,y) evaluates to False first
		},
		{
			name:     "FullForm of held comparison",
			input:    "FullForm(Hold(Equal(x, y)))",
			expected: "\"Hold(Equal(x, y))\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := eval.Evaluate(expr)
			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}
