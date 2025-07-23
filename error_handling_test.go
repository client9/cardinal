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

	expected := "$Failed(TestError)"
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
	listExpr := List{Elements: []Expr{NewSymbolAtom("Plus"), NewIntAtom(1), NewIntAtom(2)}}

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

func TestErrorPropagation(t *testing.T) {
	eval := setupTestEvaluator()

	// Create an expression that will cause division by zero in nested evaluation
	expr, err := ParseString("And(True, Equal(1, Divide(1, 0)))")
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
