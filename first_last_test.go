package sexpr

import (
	"testing"
)

func TestEvaluateFirst(t *testing.T) {
	tests := []struct {
		name      string
		input     Expr
		expected  string
		hasError  bool
		errorType string
	}{
		{
			name: "First of simple function",
			input: List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewIntAtom(1),
				NewIntAtom(2),
				NewIntAtom(3),
			}},
			expected: "1",
			hasError: false,
		},
		{
			name: "First of two-element list",
			input: List{Elements: []Expr{
				NewSymbolAtom("List"),
				NewStringAtom("hello"),
			}},
			expected: "\"hello\"",
			hasError: false,
		},
		{
			name: "First of nested expression",
			input: List{Elements: []Expr{
				NewSymbolAtom("Equal"),
				List{Elements: []Expr{
					NewSymbolAtom("Plus"),
					NewIntAtom(1),
					NewIntAtom(2),
				}},
				NewIntAtom(3),
			}},
			expected: "Plus(1, 2)",
			hasError: false,
		},
		{
			name: "First of list with mixed types",
			input: List{Elements: []Expr{
				NewSymbolAtom("Mixed"),
				NewIntAtom(42),
				NewFloatAtom(3.14),
				NewBoolAtom(true),
				NewStringAtom("test"),
			}},
			expected: "42",
			hasError: false,
		},
		{
			name: "First of empty list - should error",
			input: List{Elements: []Expr{}},
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
		{
			name: "First of single element list (head only) - should error",
			input: List{Elements: []Expr{
				NewSymbolAtom("OnlyHead"),
			}},
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
		{
			name: "First of integer atom - should error",
			input: NewIntAtom(42),
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
		{
			name: "First of string atom - should error",
			input: NewStringAtom("hello"),
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
		{
			name: "First of symbol atom - should error",
			input: NewSymbolAtom("x"),
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
		{
			name: "First of boolean atom - should error",
			input: NewBoolAtom(true),
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateFirst([]Expr{tt.input})
			
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

func TestEvaluateLast(t *testing.T) {
	tests := []struct {
		name      string
		input     Expr
		expected  string
		hasError  bool
		errorType string
	}{
		{
			name: "Last of simple function",
			input: List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewIntAtom(1),
				NewIntAtom(2),
				NewIntAtom(3),
			}},
			expected: "3",
			hasError: false,
		},
		{
			name: "Last of two-element list",
			input: List{Elements: []Expr{
				NewSymbolAtom("List"),
				NewStringAtom("hello"),
			}},
			expected: "\"hello\"",
			hasError: false,
		},
		{
			name: "Last of nested expression",
			input: List{Elements: []Expr{
				NewSymbolAtom("Equal"),
				List{Elements: []Expr{
					NewSymbolAtom("Plus"),
					NewIntAtom(1),
					NewIntAtom(2),
				}},
				NewIntAtom(3),
			}},
			expected: "3",
			hasError: false,
		},
		{
			name: "Last of list with mixed types",
			input: List{Elements: []Expr{
				NewSymbolAtom("Mixed"),
				NewIntAtom(42),
				NewFloatAtom(3.14),
				NewBoolAtom(true),
				NewStringAtom("test"),
			}},
			expected: "\"test\"",
			hasError: false,
		},
		{
			name: "Last of function with single argument",
			input: List{Elements: []Expr{
				NewSymbolAtom("Head"),
				NewSymbolAtom("x"),
			}},
			expected: "x",
			hasError: false,
		},
		{
			name: "Last of empty list - should error",
			input: List{Elements: []Expr{}},
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
		{
			name: "Last of single element list (head only) - should error",
			input: List{Elements: []Expr{
				NewSymbolAtom("OnlyHead"),
			}},
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
		{
			name: "Last of integer atom - should error",
			input: NewIntAtom(42),
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
		{
			name: "Last of string atom - should error",
			input: NewStringAtom("hello"),
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
		{
			name: "Last of symbol atom - should error",
			input: NewSymbolAtom("x"),
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
		{
			name: "Last of boolean atom - should error",
			input: NewBoolAtom(false),
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateLast([]Expr{tt.input})
			
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

// Test argument validation for First and Last
func TestFirstLast_ArgumentValidation(t *testing.T) {
	functions := []struct {
		name string
		fn   func([]Expr) Expr
	}{
		{"First", EvaluateFirst},
		{"Last", EvaluateLast},
	}
	
	for _, fn := range functions {
		t.Run(fn.name+"_no_args", func(t *testing.T) {
			result := fn.fn([]Expr{})
			if !IsError(result) {
				t.Errorf("expected error for no arguments, got %s", result.String())
			}
			
			errorExpr := result.(*ErrorExpr)
			if errorExpr.ErrorType != "ArgumentError" {
				t.Errorf("expected ArgumentError, got %s", errorExpr.ErrorType)
			}
		})
		
		t.Run(fn.name+"_too_many_args", func(t *testing.T) {
			result := fn.fn([]Expr{NewIntAtom(1), NewIntAtom(2)})
			if !IsError(result) {
				t.Errorf("expected error for too many arguments, got %s", result.String())
			}
			
			errorExpr := result.(*ErrorExpr)
			if errorExpr.ErrorType != "ArgumentError" {
				t.Errorf("expected ArgumentError, got %s", errorExpr.ErrorType)
			}
		})
	}
}

// Integration tests with the evaluator
func TestFirstLast_Integration(t *testing.T) {
	eval := setupTestEvaluator()
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// First tests
		{
			name:     "First of simple list",
			input:    "First([1, 2, 3])",
			expected: "1",
		},
		{
			name:     "First of function call",
			input:    "First(Plus(a, b, c))",
			expected: "a",
		},
		{
			name:     "First of nested structure",
			input:    "First(Equal(Plus(1, 2), 3))",
			expected: "$Failed(PartError)", // Equal(3, 3) evaluates to True, First(True) errors
		},
		{
			name:     "First error on atom",
			input:    "First(42)",
			expected: "$Failed(PartError)",
		},
		{
			name:     "First error on empty list",
			input:    "First([])",
			expected: "$Failed(PartError)",
		},
		
		// Last tests
		{
			name:     "Last of simple list",
			input:    "Last([1, 2, 3])",
			expected: "3",
		},
		{
			name:     "Last of function call",
			input:    "Last(Plus(a, b, c))",
			expected: "c",
		},
		{
			name:     "Last of nested structure",
			input:    "Last(Equal(Plus(1, 2), 3))",
			expected: "$Failed(PartError)", // Equal(3, 3) evaluates to True, Last(True) errors
		},
		{
			name:     "Last error on atom",
			input:    "Last(42)",
			expected: "$Failed(PartError)",
		},
		{
			name:     "Last error on empty list",
			input:    "Last([])",
			expected: "$Failed(PartError)",
		},
		
		// Combined tests
		{
			name:     "First of Last result",
			input:    "First([Last([1, 2, 3]), 4, 5])",
			expected: "3",
		},
		{
			name:     "Last of First result",
			input:    "Last([1, 2, First([10, 11, 12])])",
			expected: "10",
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

// Test First and Last with complex expressions
func TestFirstLast_ComplexExpressions(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		function func([]Expr) Expr
		expected string
	}{
		{
			name: "First of deeply nested expression",
			expr: List{Elements: []Expr{
				NewSymbolAtom("Outer"),
				List{Elements: []Expr{
					NewSymbolAtom("Inner"),
					List{Elements: []Expr{
						NewSymbolAtom("Deep"),
						NewIntAtom(1),
						NewIntAtom(2),
					}},
					NewIntAtom(3),
				}},
				NewIntAtom(4),
			}},
			function: EvaluateFirst,
			expected: "Inner(Deep(1, 2), 3)",
		},
		{
			name: "Last of deeply nested expression",
			expr: List{Elements: []Expr{
				NewSymbolAtom("Outer"),
				List{Elements: []Expr{
					NewSymbolAtom("Inner"),
					List{Elements: []Expr{
						NewSymbolAtom("Deep"),
						NewIntAtom(1),
						NewIntAtom(2),
					}},
					NewIntAtom(3),
				}},
				NewIntAtom(4),
			}},
			function: EvaluateLast,
			expected: "4",
		},
		{
			name: "First of list with string elements",
			expr: List{Elements: []Expr{
				NewSymbolAtom("StringList"),
				NewStringAtom("hello"),
				NewStringAtom("world"),
				NewStringAtom("test"),
			}},
			function: EvaluateFirst,
			expected: "\"hello\"",
		},
		{
			name: "Last of list with mixed atom types",
			expr: List{Elements: []Expr{
				NewSymbolAtom("Mixed"),
				NewIntAtom(42),
				NewFloatAtom(3.14159),
				NewBoolAtom(true),
				NewSymbolAtom("x"),
			}},
			function: EvaluateLast,
			expected: "x",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function([]Expr{tt.expr})
			
			if IsError(result) {
				t.Errorf("unexpected error: %s", result.String())
				return
			}
			
			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}
