package sexpr

import (
	"testing"
)

func TestEvaluateRest_Basic(t *testing.T) {
	tests := []struct {
		name      string
		input     Expr
		expected  string
		hasError  bool
		errorType string
	}{
		{
			name: "Rest of simple function",
			input: List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewIntAtom(1),
				NewIntAtom(2),
				NewIntAtom(3),
			}},
			expected: "Plus(2, 3)",
			hasError: false,
		},
		{
			name: "Rest of two-element list",
			input: List{Elements: []Expr{
				NewSymbolAtom("List"),
				NewStringAtom("hello"),
			}},
			expected: "List()",
			hasError: false,
		},
		{
			name:      "Rest of empty list - should error",
			input:     List{Elements: []Expr{}},
			expected:  "",
			hasError:  true,
			errorType: "PartError",
		},
		{
			name: "Rest of single element list (head only) - should error",
			input: List{Elements: []Expr{
				NewSymbolAtom("OnlyHead"),
			}},
			expected:  "",
			hasError:  true,
			errorType: "PartError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateRest([]Expr{tt.input})

			if tt.hasError {
				if !IsError(result) {
					t.Errorf("expected error for %s, got %s", tt.name, result.String())
					return
				}

				errorExpr := result.(*ErrorExpr)
				if errorExpr.ErrorType != tt.errorType {
					t.Errorf("expected error type %s, got %s", tt.errorType, errorExpr.ErrorType)
				}
			} else {
				if IsError(result) {
					t.Errorf("unexpected error: %s", result.String())
					return
				}

				if result.String() != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, result.String())
				}
			}
		})
	}
}

func TestEvaluateMost_Basic(t *testing.T) {
	tests := []struct {
		name      string
		input     Expr
		expected  string
		hasError  bool
		errorType string
	}{
		{
			name: "Most of simple function",
			input: List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewIntAtom(1),
				NewIntAtom(2),
				NewIntAtom(3),
			}},
			expected: "Plus(1, 2)",
			hasError: false,
		},
		{
			name: "Most of two-element list",
			input: List{Elements: []Expr{
				NewSymbolAtom("List"),
				NewStringAtom("hello"),
			}},
			expected: "List()",
			hasError: false,
		},
		{
			name:      "Most of empty list - should error",
			input:     List{Elements: []Expr{}},
			expected:  "",
			hasError:  true,
			errorType: "PartError",
		},
		{
			name: "Most of single element list (head only) - should error",
			input: List{Elements: []Expr{
				NewSymbolAtom("OnlyHead"),
			}},
			expected:  "",
			hasError:  true,
			errorType: "PartError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateMost([]Expr{tt.input})

			if tt.hasError {
				if !IsError(result) {
					t.Errorf("expected error for %s, got %s", tt.name, result.String())
					return
				}

				errorExpr := result.(*ErrorExpr)
				if errorExpr.ErrorType != tt.errorType {
					t.Errorf("expected error type %s, got %s", tt.errorType, errorExpr.ErrorType)
				}
			} else {
				if IsError(result) {
					t.Errorf("unexpected error: %s", result.String())
					return
				}

				if result.String() != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, result.String())
				}
			}
		})
	}
}

func TestRestMost_Integration(t *testing.T) {
	eval := setupTestEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Rest tests
		{
			name:     "Rest of simple list",
			input:    "Rest([1, 2, 3, 4])",
			expected: "List(2, 3, 4)",
		},
		{
			name:     "Rest of function call",
			input:    "Rest(Plus(a, b, c, d))",
			expected: "Plus(b, c, d)",
		},
		{
			name:     "Rest of two-element list",
			input:    "Rest([x, y])",
			expected: "List(y)",
		},
		{
			name:     "Rest of single-element list",
			input:    "Rest([x])",
			expected: "List()",
		},
		{
			name:     "Rest error on atom",
			input:    "Rest(42)",
			expected: "Rest(42)", // Pattern doesn't match, returns unchanged
		},
		{
			name:     "Rest error on empty list",
			input:    "Rest([])",
			expected: "$Failed(PartError)",
		},

		// Most tests
		{
			name:     "Most of simple list",
			input:    "Most([1, 2, 3, 4])",
			expected: "List(1, 2, 3)",
		},
		{
			name:     "Most of function call",
			input:    "Most(Plus(a, b, c, d))",
			expected: "Plus(a, b, c)",
		},
		{
			name:     "Most of two-element list",
			input:    "Most([x, y])",
			expected: "List(x)",
		},
		{
			name:     "Most of single-element list",
			input:    "Most([x])",
			expected: "List()",
		},
		{
			name:     "Most error on atom",
			input:    "Most(42)",
			expected: "Most(42)", // Pattern doesn't match, returns unchanged
		},
		{
			name:     "Most error on empty list",
			input:    "Most([])",
			expected: "$Failed(PartError)",
		},

		// Combined tests
		{
			name:     "Rest of Most result",
			input:    "Rest(Most([1, 2, 3, 4]))",
			expected: "List(2, 3)",
		},
		{
			name:     "Most of Rest result",
			input:    "Most(Rest([1, 2, 3, 4]))",
			expected: "List(2, 3)",
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
