package sexpr

import (
	"testing"
)

func TestErrorExpr_Basic(t *testing.T) {
	// Test ErrorExpr construction and methods
	err := NewErrorExpr("TestError", "This is a test error", []Expr{NewIntAtom(1), NewIntAtom(2)})
	
	if err.ErrorType != "TestError" {
		t.Errorf("expected ErrorType 'TestError', got '%s'", err.ErrorType)
	}
	
	if err.Message != "This is a test error" {
		t.Errorf("expected Message 'This is a test error', got '%s'", err.Message)
	}
	
	expected := "$Failed[TestError]"
	if err.String() != expected {
		t.Errorf("expected String() '%s', got '%s'", expected, err.String())
	}
	
	if err.Type() != "error" {
		t.Errorf("expected Type() 'error', got '%s'", err.Type())
	}
}

func TestIsError(t *testing.T) {
	// Test IsError function
	errorExpr := NewErrorExpr("TestError", "test", nil)
	intExpr := NewIntAtom(42)
	listExpr := &List{Elements: []Expr{NewSymbolAtom("Plus"), NewIntAtom(1), NewIntAtom(2)}}
	
	if !IsError(errorExpr) {
		t.Error("IsError should return true for ErrorExpr")
	}
	
	if IsError(intExpr) {
		t.Error("IsError should return false for IntAtom")
	}
	
	if IsError(listExpr) {
		t.Error("IsError should return false for List")
	}
}

func TestArithmeticErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		function    string
		args        []Expr
		expectedErr string
	}{
		{
			name:        "Subtract wrong args",
			function:    "Subtract",
			args:        []Expr{NewIntAtom(1)},
			expectedErr: "ArgumentError",
		},
		{
			name:        "Subtract too many args",
			function:    "Subtract",
			args:        []Expr{NewIntAtom(1), NewIntAtom(2), NewIntAtom(3)},
			expectedErr: "ArgumentError",
		},
		{
			name:        "Divide wrong args",
			function:    "Divide",
			args:        []Expr{NewIntAtom(1)},
			expectedErr: "ArgumentError",
		},
		{
			name:        "Divide by zero",
			function:    "Divide",
			args:        []Expr{NewIntAtom(10), NewIntAtom(0)},
			expectedErr: "DivisionByZero",
		},
		{
			name:        "Power wrong args",
			function:    "Power",
			args:        []Expr{NewIntAtom(2)},
			expectedErr: "ArgumentError",
		},
		{
			name:        "Power invalid result",
			function:    "Power",
			args:        []Expr{NewFloatAtom(-1), NewFloatAtom(0.5)}, // sqrt of negative
			expectedErr: "MathematicalError",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Expr
			
			switch tt.function {
			case "Subtract":
				result = EvaluateSubtract(tt.args)
			case "Divide":
				result = EvaluateDivide(tt.args)
			case "Power":
				result = EvaluatePower(tt.args)
			default:
				t.Fatalf("Unknown function: %s", tt.function)
			}
			
			if !IsError(result) {
				t.Errorf("expected error, got: %s", result.String())
				return
			}
			
			errorExpr := result.(*ErrorExpr)
			if errorExpr.ErrorType != tt.expectedErr {
				t.Errorf("expected error type '%s', got '%s'", tt.expectedErr, errorExpr.ErrorType)
			}
		})
	}
}

func TestComparisonErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		function    string
		args        []Expr
		expectedErr string
	}{
		{
			name:        "Equal wrong args",
			function:    "Equal",
			args:        []Expr{NewIntAtom(1)},
			expectedErr: "ArgumentError",
		},
		{
			name:        "Less wrong args",
			function:    "Less",
			args:        []Expr{NewIntAtom(1), NewIntAtom(2), NewIntAtom(3)},
			expectedErr: "ArgumentError",
		},
		{
			name:        "Greater no args",
			function:    "Greater",
			args:        []Expr{},
			expectedErr: "ArgumentError",
		},
		{
			name:        "SameQ wrong args",
			function:    "SameQ",
			args:        []Expr{NewIntAtom(1)},
			expectedErr: "ArgumentError",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Expr
			
			switch tt.function {
			case "Equal":
				result = EvaluateEqual(tt.args)
			case "Less":
				result = EvaluateLess(tt.args)
			case "Greater":
				result = EvaluateGreater(tt.args)
			case "SameQ":
				result = EvaluateSameQ(tt.args)
			default:
				t.Fatalf("Unknown function: %s", tt.function)
			}
			
			if !IsError(result) {
				t.Errorf("expected error, got: %s", result.String())
				return
			}
			
			errorExpr := result.(*ErrorExpr)
			if errorExpr.ErrorType != tt.expectedErr {
				t.Errorf("expected error type '%s', got '%s'", tt.expectedErr, errorExpr.ErrorType)
			}
		})
	}
}

func TestLogicalErrorHandling(t *testing.T) {
	// Test Not with wrong number of arguments
	result := EvaluateNot([]Expr{NewBoolAtom(true), NewBoolAtom(false)})
	
	if !IsError(result) {
		t.Errorf("expected error for Not with 2 arguments, got: %s", result.String())
		return
	}
	
	errorExpr := result.(*ErrorExpr)
	if errorExpr.ErrorType != "ArgumentError" {
		t.Errorf("expected ArgumentError, got: %s", errorExpr.ErrorType)
	}
}

func TestErrorPropagation(t *testing.T) {
	eval := setupTestEvaluator()
	
	// Create an expression that will cause division by zero in nested evaluation
	expr, err := ParseString("And[True, Equal[1, Divide[1, 0]]]")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	
	result := eval.Evaluate(expr)
	
	// The error should propagate up through the And expression
	if !IsError(result) {
		t.Errorf("expected error to propagate, got: %s", result.String())
		return
	}
	
	errorExpr := result.(*ErrorExpr)
	if errorExpr.ErrorType != "DivisionByZero" {
		t.Errorf("expected DivisionByZero error, got: %s", errorExpr.ErrorType)
	}
}

func TestValidOperationsStillWork(t *testing.T) {
	// Ensure that valid operations still work correctly
	tests := []struct {
		name     string
		function string
		args     []Expr
		expected string
	}{
		{
			name:     "Valid Subtract",
			function: "Subtract",
			args:     []Expr{NewIntAtom(5), NewIntAtom(3)},
			expected: "2",
		},
		{
			name:     "Valid Divide",
			function: "Divide",
			args:     []Expr{NewIntAtom(10), NewIntAtom(2)},
			expected: "5",
		},
		{
			name:     "Valid Power",
			function: "Power",
			args:     []Expr{NewIntAtom(2), NewIntAtom(3)},
			expected: "8",
		},
		{
			name:     "Valid Equal",
			function: "Equal",
			args:     []Expr{NewIntAtom(5), NewIntAtom(5)},
			expected: "True",
		},
		{
			name:     "Valid Not",
			function: "Not",
			args:     []Expr{NewBoolAtom(false)},
			expected: "True",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Expr
			
			switch tt.function {
			case "Subtract":
				result = EvaluateSubtract(tt.args)
			case "Divide":
				result = EvaluateDivide(tt.args)
			case "Power":
				result = EvaluatePower(tt.args)
			case "Equal":
				result = EvaluateEqual(tt.args)
			case "Not":
				result = EvaluateNot(tt.args)
			default:
				t.Fatalf("Unknown function: %s", tt.function)
			}
			
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