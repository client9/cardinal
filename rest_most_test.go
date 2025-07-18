package sexpr

import (
	"testing"
)

func TestEvaluateRest(t *testing.T) {
	tests := []struct {
		name      string
		input     Expr
		expected  string
		hasError  bool
		errorType string
	}{
		{
			name: "Rest of simple function with three elements",
			input: &List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewIntAtom(1),
				NewIntAtom(2),
				NewIntAtom(3),
			}},
			expected: "Plus[2, 3]",
			hasError: false,
		},
		{
			name: "Rest of function with four elements",
			input: &List{Elements: []Expr{
				NewSymbolAtom("Times"),
				NewIntAtom(1),
				NewIntAtom(2),
				NewIntAtom(3),
				NewIntAtom(4),
			}},
			expected: "Times[2, 3, 4]",
			hasError: false,
		},
		{
			name: "Rest of two-element list (head + one element)",
			input: &List{Elements: []Expr{
				NewSymbolAtom("List"),
				NewStringAtom("hello"),
			}},
			expected: "List[]", // Only head remains
			hasError: false,
		},
		{
			name: "Rest of nested expression",
			input: &List{Elements: []Expr{
				NewSymbolAtom("Equal"),
				&List{Elements: []Expr{
					NewSymbolAtom("Plus"),
					NewIntAtom(1),
					NewIntAtom(2),
				}},
				NewIntAtom(3),
				NewIntAtom(4),
			}},
			expected: "Equal[3, 4]",
			hasError: false,
		},
		{
			name: "Rest of list with mixed types",
			input: &List{Elements: []Expr{
				NewSymbolAtom("Mixed"),
				NewIntAtom(42),
				NewFloatAtom(3.14),
				NewBoolAtom(true),
				NewStringAtom("test"),
			}},
			expected: "Mixed[3.14, True, \"test\"]",
			hasError: false,
		},
		{
			name: "Rest of empty list - should error",
			input: &List{Elements: []Expr{}},
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
		{
			name: "Rest of single element list (head only) - should error",
			input: &List{Elements: []Expr{
				NewSymbolAtom("OnlyHead"),
			}},
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
		{
			name: "Rest of integer atom - should error",
			input: NewIntAtom(42),
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
		{
			name: "Rest of string atom - should error",
			input: NewStringAtom("hello"),
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
		{
			name: "Rest of symbol atom - should error",
			input: NewSymbolAtom("x"),
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
		{
			name: "Rest of boolean atom - should error",
			input: NewBoolAtom(true),
			expected: "",
			hasError: true,
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

func TestEvaluateMost(t *testing.T) {
	tests := []struct {
		name      string
		input     Expr
		expected  string
		hasError  bool
		errorType string
	}{
		{
			name: "Most of simple function with three elements",
			input: &List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewIntAtom(1),
				NewIntAtom(2),
				NewIntAtom(3),
			}},
			expected: "Plus[1, 2]",
			hasError: false,
		},
		{
			name: "Most of function with four elements",
			input: &List{Elements: []Expr{
				NewSymbolAtom("Times"),
				NewIntAtom(1),
				NewIntAtom(2),
				NewIntAtom(3),
				NewIntAtom(4),
			}},
			expected: "Times[1, 2, 3]",
			hasError: false,
		},
		{
			name: "Most of two-element list (head + one element)",
			input: &List{Elements: []Expr{
				NewSymbolAtom("List"),
				NewStringAtom("hello"),
			}},
			expected: "List[]", // Only head remains
			hasError: false,
		},
		{
			name: "Most of nested expression",
			input: &List{Elements: []Expr{
				NewSymbolAtom("Equal"),
				&List{Elements: []Expr{
					NewSymbolAtom("Plus"),
					NewIntAtom(1),
					NewIntAtom(2),
				}},
				NewIntAtom(3),
				NewIntAtom(4),
			}},
			expected: "Equal[Plus[1, 2], 3]",
			hasError: false,
		},
		{
			name: "Most of list with mixed types",
			input: &List{Elements: []Expr{
				NewSymbolAtom("Mixed"),
				NewIntAtom(42),
				NewFloatAtom(3.14),
				NewBoolAtom(true),
				NewStringAtom("test"),
			}},
			expected: "Mixed[42, 3.14, True]",
			hasError: false,
		},
		{
			name: "Most of function with single argument",
			input: &List{Elements: []Expr{
				NewSymbolAtom("Head"),
				NewSymbolAtom("x"),
			}},
			expected: "Head[]",
			hasError: false,
		},
		{
			name: "Most of empty list - should error",
			input: &List{Elements: []Expr{}},
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
		{
			name: "Most of single element list (head only) - should error",
			input: &List{Elements: []Expr{
				NewSymbolAtom("OnlyHead"),
			}},
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
		{
			name: "Most of integer atom - should error",
			input: NewIntAtom(42),
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
		{
			name: "Most of string atom - should error",
			input: NewStringAtom("hello"),
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
		{
			name: "Most of symbol atom - should error",
			input: NewSymbolAtom("x"),
			expected: "",
			hasError: true,
			errorType: "PartError",
		},
		{
			name: "Most of boolean atom - should error",
			input: NewBoolAtom(false),
			expected: "",
			hasError: true,
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

// Test argument validation for Rest and Most
func TestRestMost_ArgumentValidation(t *testing.T) {
	functions := []struct {
		name string
		fn   func([]Expr) Expr
	}{
		{"Rest", EvaluateRest},
		{"Most", EvaluateMost},
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
			input:    "Rest[{1, 2, 3, 4}]",
			expected: "List[2, 3, 4]",
		},
		{
			name:     "Rest of function call",
			input:    "Rest[Plus[a, b, c, d]]",
			expected: "Plus[b, c, d]",
		},
		{
			name:     "Rest of two-element list",
			input:    "Rest[{x, y}]",
			expected: "List[y]",
		},
		{
			name:     "Rest of single-element list",
			input:    "Rest[{x}]",
			expected: "List[]",
		},
		{
			name:     "Rest error on atom",
			input:    "Rest[42]",
			expected: "$Failed[PartError]",
		},
		{
			name:     "Rest error on empty list",
			input:    "Rest[{}]",
			expected: "$Failed[PartError]",
		},
		
		// Most tests
		{
			name:     "Most of simple list",
			input:    "Most[{1, 2, 3, 4}]",
			expected: "List[1, 2, 3]",
		},
		{
			name:     "Most of function call",
			input:    "Most[Plus[a, b, c, d]]",
			expected: "Plus[a, b, c]",
		},
		{
			name:     "Most of two-element list",
			input:    "Most[{x, y}]",
			expected: "List[x]",
		},
		{
			name:     "Most of single-element list",
			input:    "Most[{x}]",
			expected: "List[]",
		},
		{
			name:     "Most error on atom",
			input:    "Most[42]",
			expected: "$Failed[PartError]",
		},
		{
			name:     "Most error on empty list",
			input:    "Most[{}]",
			expected: "$Failed[PartError]",
		},
		
		// Combined tests
		{
			name:     "Rest of Most result",
			input:    "Rest[Most[{1, 2, 3, 4}]]",
			expected: "List[2, 3]",
		},
		{
			name:     "Most of Rest result",
			input:    "Most[Rest[{1, 2, 3, 4}]]",
			expected: "List[2, 3]",
		},
		{
			name:     "First of Rest result",
			input:    "First[Rest[{a, b, c, d}]]",
			expected: "b",
		},
		{
			name:     "Last of Most result",
			input:    "Last[Most[{a, b, c, d}]]",
			expected: "c",
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

// Test Rest and Most with complex expressions
func TestRestMost_ComplexExpressions(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		function func([]Expr) Expr
		expected string
	}{
		{
			name: "Rest of deeply nested expression",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("Outer"),
				&List{Elements: []Expr{
					NewSymbolAtom("Inner"),
					NewIntAtom(1),
					NewIntAtom(2),
				}},
				NewIntAtom(3),
				NewIntAtom(4),
			}},
			function: EvaluateRest,
			expected: "Outer[3, 4]",
		},
		{
			name: "Most of deeply nested expression",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("Outer"),
				&List{Elements: []Expr{
					NewSymbolAtom("Inner"),
					NewIntAtom(1),
					NewIntAtom(2),
				}},
				NewIntAtom(3),
				NewIntAtom(4),
			}},
			function: EvaluateMost,
			expected: "Outer[Inner[1, 2], 3]",
		},
		{
			name: "Rest of list with string elements",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("StringList"),
				NewStringAtom("hello"),
				NewStringAtom("world"),
				NewStringAtom("test"),
			}},
			function: EvaluateRest,
			expected: "StringList[\"world\", \"test\"]",
		},
		{
			name: "Most of list with mixed atom types",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("Mixed"),
				NewIntAtom(42),
				NewFloatAtom(3.14159),
				NewBoolAtom(true),
				NewSymbolAtom("x"),
			}},
			function: EvaluateMost,
			expected: "Mixed[42, 3.14159, True]",
		},
		{
			name: "Rest of function with many arguments",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("ManyArgs"),
				NewIntAtom(1),
				NewIntAtom(2),
				NewIntAtom(3),
				NewIntAtom(4),
				NewIntAtom(5),
				NewIntAtom(6),
			}},
			function: EvaluateRest,
			expected: "ManyArgs[2, 3, 4, 5, 6]",
		},
		{
			name: "Most of function with many arguments",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("ManyArgs"),
				NewIntAtom(1),
				NewIntAtom(2),
				NewIntAtom(3),
				NewIntAtom(4),
				NewIntAtom(5),
				NewIntAtom(6),
			}},
			function: EvaluateMost,
			expected: "ManyArgs[1, 2, 3, 4, 5]",
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

// Test edge cases and special scenarios
func TestRestMost_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		restFunc func([]Expr) Expr
		mostFunc func([]Expr) Expr
		expectedRest string
		expectedMost string
	}{
		{
			name: "Two-element list",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("F"),
				NewIntAtom(42),
			}},
			restFunc: EvaluateRest,
			mostFunc: EvaluateMost,
			expectedRest: "F[]",
			expectedMost: "F[]",
		},
		{
			name: "Three-element list",
			expr: &List{Elements: []Expr{
				NewSymbolAtom("G"),
				NewIntAtom(1),
				NewIntAtom(2),
			}},
			restFunc: EvaluateRest,
			mostFunc: EvaluateMost,
			expectedRest: "G[2]",
			expectedMost: "G[1]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_Rest", func(t *testing.T) {
			result := tt.restFunc([]Expr{tt.expr})
			
			if IsError(result) {
				t.Errorf("unexpected error: %s", result.String())
				return
			}
			
			if result.String() != tt.expectedRest {
				t.Errorf("Rest: expected %s, got %s", tt.expectedRest, result.String())
			}
		})
		
		t.Run(tt.name+"_Most", func(t *testing.T) {
			result := tt.mostFunc([]Expr{tt.expr})
			
			if IsError(result) {
				t.Errorf("unexpected error: %s", result.String())
				return
			}
			
			if result.String() != tt.expectedMost {
				t.Errorf("Most: expected %s, got %s", tt.expectedMost, result.String())
			}
		})
	}
}