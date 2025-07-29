package sexpr

import (
	"github.com/client9/sexpr/core"
	"testing"
)

func TestErrorExpr_Basic(t *testing.T) {
	// Test ErrorExpr construction and methods
	err := core.NewErrorExpr("TestError", "This is a test error", []core.Expr{core.NewInteger(1), core.NewInteger(2)})

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

	if err.Head() != "error" {
		t.Errorf("expected Type() 'error', got '%s'", err.Head())
	}
}

func TestIsError(t *testing.T) {
	// Test core.IsError function
	errorExpr := core.NewErrorExpr("TestError", "test", nil)
	intExpr := core.NewInteger(42)
	listExpr := core.NewList("Plus", core.NewInteger(1), core.NewInteger(2))

	if !core.IsError(errorExpr) {
		t.Error("core.IsError should return true for ErrorExpr")
	}

	if core.IsError(intExpr) {
		t.Error("core.IsError should return false for IntAtom")
	}

	if core.IsError(listExpr) {
		t.Error("core.IsError should return false for List")
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
	if !core.IsError(result) {
		t.Errorf("expected error to propagate, got: %s", result.String())
		return
	}

	errorExpr := result.(*core.ErrorExpr)
	if errorExpr.ErrorType != "DivisionByZero" {
		t.Errorf("expected DivisionByZero error, got: %s", errorExpr.ErrorType)
	}
}
