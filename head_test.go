package sexpr

import (
	"testing"
)

func TestEvaluateHead_Atoms(t *testing.T) {
	tests := []struct {
		name     string
		arg      Expr
		expected string
	}{
		{
			name:     "Integer",
			arg:      NewIntAtom(42),
			expected: "Integer",
		},
		{
			name:     "Real",
			arg:      NewFloatAtom(3.14),
			expected: "Real",
		},
		{
			name:     "String",
			arg:      NewStringAtom("hello"),
			expected: "String",
		},
		{
			name:     "Boolean true",
			arg:      NewBoolAtom(true),
			expected: "Boolean",
		},
		{
			name:     "Boolean false",
			arg:      NewBoolAtom(false),
			expected: "Boolean",
		},
		{
			name:     "Symbol",
			arg:      NewSymbolAtom("x"),
			expected: "Symbol",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateHead([]Expr{tt.arg})
			
			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
			
			// Verify it returns a Symbol
			if !isSymbol(result) {
				t.Errorf("Head should return a Symbol, got %T", result)
			}
		})
	}
}

func TestEvaluateHead_Lists(t *testing.T) {
	tests := []struct {
		name     string
		arg      Expr
		expected string
	}{
		{
			name:     "Empty list",
			arg:      List{Elements: []Expr{}},
			expected: "List",
		},
		{
			name: "Function call",
			arg: List{Elements: []Expr{
				NewSymbolAtom("Plus"),
				NewIntAtom(1),
				NewIntAtom(2),
			}},
			expected: "Plus",
		},
		{
			name: "Nested expression", 
			arg: List{Elements: []Expr{
				NewSymbolAtom("Times"),
				NewIntAtom(2),
				List{Elements: []Expr{
					NewSymbolAtom("Plus"),
					NewIntAtom(1),
					NewIntAtom(3),
				}},
			}},
			expected: "Times",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateHead([]Expr{tt.arg})
			
			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestEvaluateHead_Errors(t *testing.T) {
	// Test error expressions
	errorExpr := NewErrorExpr("TestError", "test message", nil)
	result := EvaluateHead([]Expr{errorExpr})
	
	expected := "Error"
	if result.String() != expected {
		t.Errorf("expected %s, got %s", expected, result.String())
	}
}

func TestEvaluateHead_ArgumentValidation(t *testing.T) {
	tests := []struct {
		name string
		args []Expr
	}{
		{
			name: "No arguments",
			args: []Expr{},
		},
		{
			name: "Too many arguments",
			args: []Expr{NewIntAtom(1), NewIntAtom(2)},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateHead(tt.args)
			
			if !IsError(result) {
				t.Errorf("expected error for %s, got %s", tt.name, result.String())
				return
			}
			
			errorExpr := result.(*ErrorExpr)
			if errorExpr.ErrorType != "ArgumentError" {
				t.Errorf("expected ArgumentError, got %s", errorExpr.ErrorType)
			}
		})
	}
}

func TestHead_Integration(t *testing.T) {
	eval := setupTestEvaluator()
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Head of integer",
			input:    "Head(42)",
			expected: "Integer",
		},
		{
			name:     "Head of float",
			input:    "Head(3.14)",
			expected: "Real",
		},
		{
			name:     "Head of string",
			input:    "Head(\"test\")",
			expected: "String",
		},
		{
			name:     "Head of boolean",
			input:    "Head(True)",
			expected: "Boolean",
		},
		{
			name:     "Head of symbol",
			input:    "Head(x)",
			expected: "Symbol",
		},
		{
			name:     "Head of empty list",
			input:    "Head([])",
			expected: "List",
		},
		{
			name:     "Head of unevaluated function",
			input:    "Head(Plus(x, y))",
			expected: "Plus",
		},
		{
			name:     "Head of evaluated function",
			input:    "Head(Plus(1, 2))",
			expected: "Integer", // Plus(1,2) evaluates to 3, head of 3 is Integer
		},
		{
			name:     "Head of held expression",
			input:    "Head(Hold(Plus(1, 2)))",
			expected: "Hold",
		},
		{
			name:     "Head of comparison",
			input:    "Head(Equal(1, 2))",
			expected: "Boolean", // Equal(1,2) evaluates to False, head of False is Boolean
		},
		{
			name:     "Head of symbolic comparison",
			input:    "Head(Equal(x, y))",
			expected: "Boolean", // Equal(x,y) evaluates to False (string comparison), so head is Boolean
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

func TestHead_ErrorPropagation(t *testing.T) {
	eval := setupTestEvaluator()
	
	// Test that Head propagates errors
	expr, err := ParseString("Head(Divide(1, 0))")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	
	result := eval.Evaluate(expr)
	
	// Should get the DivisionByZero error, not the head of the error
	if !IsError(result) {
		t.Errorf("expected error propagation, got %s", result.String())
		return
	}
	
	errorExpr := result.(*ErrorExpr)
	if errorExpr.ErrorType != "DivisionByZero" {
		t.Errorf("expected DivisionByZero error, got %s", errorExpr.ErrorType)
	}
}

func TestHead_CompareWithCurrentType(t *testing.T) {
	// Compare Head output with current Type() method
	testExprs := []Expr{
		NewIntAtom(42),
		NewFloatAtom(3.14),
		NewStringAtom("hello"),
		NewBoolAtom(true),
		NewSymbolAtom("x"),
		List{Elements: []Expr{}},
		List{Elements: []Expr{NewSymbolAtom("Plus"), NewIntAtom(1)}},
		NewErrorExpr("TestError", "test", nil),
	}
	
	for _, expr := range testExprs {
		head := EvaluateHead([]Expr{expr})
		currentType := expr.Type()
		
		t.Logf("Expression: %s", expr.String())
		t.Logf("  Head(): %s", head.String())
		t.Logf("  Type(): %s", currentType)
		t.Logf("")
	}
}
