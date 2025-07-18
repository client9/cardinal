package sexpr

import (
	"strings"
	"testing"
)

func TestEvaluationStack_Basic(t *testing.T) {
	stack := NewEvaluationStack(10)
	
	// Test initial state
	if stack.Depth() != 0 {
		t.Errorf("expected initial depth 0, got %d", stack.Depth())
	}
	
	// Test push
	err := stack.Push("test", "test expression")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	
	if stack.Depth() != 1 {
		t.Errorf("expected depth 1 after push, got %d", stack.Depth())
	}
	
	// Test pop
	stack.Pop()
	if stack.Depth() != 0 {
		t.Errorf("expected depth 0 after pop, got %d", stack.Depth())
	}
}

func TestEvaluationStack_RecursionLimit(t *testing.T) {
	stack := NewEvaluationStack(3)
	
	// Push up to limit
	for i := 0; i < 3; i++ {
		err := stack.Push("test", "test expression")
		if err != nil {
			t.Errorf("unexpected error at depth %d: %v", i, err)
		}
	}
	
	// Try to exceed limit
	err := stack.Push("test", "test expression")
	if err == nil {
		t.Error("expected error when exceeding recursion limit")
	}
	
	if !strings.Contains(err.Error(), "maximum recursion depth exceeded") {
		t.Errorf("expected recursion error message, got: %v", err)
	}
}

func TestEvaluationStack_GetFrames(t *testing.T) {
	stack := NewEvaluationStack(10)
	
	// Push multiple frames
	stack.Push("func1", "expr1")
	stack.Push("func2", "expr2")
	stack.Push("func3", "expr3")
	
	frames := stack.GetFrames()
	
	if len(frames) != 3 {
		t.Errorf("expected 3 frames, got %d", len(frames))
	}
	
	// Check frame contents
	expected := []struct {
		function   string
		expression string
	}{
		{"func1", "expr1"},
		{"func2", "expr2"},
		{"func3", "expr3"},
	}
	
	for i, frame := range frames {
		if frame.Function != expected[i].function {
			t.Errorf("frame %d: expected function %s, got %s", i, expected[i].function, frame.Function)
		}
		if frame.Expression != expected[i].expression {
			t.Errorf("frame %d: expected expression %s, got %s", i, expected[i].expression, frame.Expression)
		}
	}
}

func TestRecursionPrevention_SimpleCase(t *testing.T) {
	eval := NewEvaluator()
	
	// Create a simple recursive definition that would cause infinite recursion
	// Set f[x_] := f[x + 1] (this would recurse infinitely)
	// We'll simulate this by directly testing the recursion limit
	
	// Set a very low recursion limit for testing
	eval.context.stack = NewEvaluationStack(5)
	
	// Create a deeply nested expression that will hit the limit
	// Plus[Plus[Plus[Plus[Plus[Plus[1, 2], 3], 4], 5], 6], 7]
	expr := NewList(
		NewSymbolAtom("Plus"),
		NewList(
			NewSymbolAtom("Plus"),
			NewList(
				NewSymbolAtom("Plus"),
				NewList(
					NewSymbolAtom("Plus"),
					NewList(
						NewSymbolAtom("Plus"),
						NewList(
							NewSymbolAtom("Plus"),
							NewIntAtom(1),
							NewIntAtom(2),
						),
						NewIntAtom(3),
					),
					NewIntAtom(4),
				),
				NewIntAtom(5),
			),
			NewIntAtom(6),
		),
		NewIntAtom(7),
	)
	
	result := eval.Evaluate(expr)
	
	// Should get a recursion error
	if !IsError(result) {
		t.Error("expected recursion error, got successful result")
		return
	}
	
	errorExpr := result.(*ErrorExpr)
	if errorExpr.ErrorType != "RecursionError" {
		t.Errorf("expected RecursionError, got %s", errorExpr.ErrorType)
	}
	
	if !strings.Contains(errorExpr.Message, "maximum recursion depth exceeded") {
		t.Errorf("expected recursion error message, got: %s", errorExpr.Message)
	}
}

func TestStackTrace_ErrorPropagation(t *testing.T) {
	eval := NewEvaluator()
	
	// Test that errors include stack traces
	// Divide[1, 0] should give an error with stack trace
	expr := NewList(
		NewSymbolAtom("Divide"),
		NewIntAtom(1),
		NewIntAtom(0),
	)
	
	result := eval.Evaluate(expr)
	
	if !IsError(result) {
		t.Error("expected error for division by zero")
		return
	}
	
	errorExpr := result.(*ErrorExpr)
	
	// Check that stack trace is present
	if len(errorExpr.StackTrace) == 0 {
		t.Error("expected stack trace in error, but got none")
	}
	
	// Check stack trace content
	stackTrace := errorExpr.GetStackTrace()
	if !strings.Contains(stackTrace, "Divide") {
		t.Errorf("expected Divide in stack trace, got: %s", stackTrace)
	}
}

func TestStackTrace_NestedErrors(t *testing.T) {
	eval := NewEvaluator()
	
	// Test nested function calls with errors
	// Plus[1, Divide[2, 0]] should show both Plus and Divide in stack trace
	expr := NewList(
		NewSymbolAtom("Plus"),
		NewIntAtom(1),
		NewList(
			NewSymbolAtom("Divide"),
			NewIntAtom(2),
			NewIntAtom(0),
		),
	)
	
	result := eval.Evaluate(expr)
	
	if !IsError(result) {
		t.Error("expected error for nested division by zero")
		return
	}
	
	errorExpr := result.(*ErrorExpr)
	stackTrace := errorExpr.GetStackTrace()
	
	// Should contain both functions in the stack trace
	if !strings.Contains(stackTrace, "Plus") {
		t.Errorf("expected Plus in stack trace, got: %s", stackTrace)
	}
	if !strings.Contains(stackTrace, "Divide") {
		t.Errorf("expected Divide in stack trace, got: %s", stackTrace)
	}
}

func TestStackTrace_StringFunctions(t *testing.T) {
	eval := NewEvaluator()
	
	// Test that string functions also get proper stack traces
	// StringLength[42] should error with stack trace
	expr := NewList(
		NewSymbolAtom("StringLength"),
		NewIntAtom(42),
	)
	
	result := eval.Evaluate(expr)
	
	if !IsError(result) {
		t.Error("expected error for StringLength on non-string")
		return
	}
	
	errorExpr := result.(*ErrorExpr)
	stackTrace := errorExpr.GetStackTrace()
	
	if !strings.Contains(stackTrace, "StringLength") {
		t.Errorf("expected StringLength in stack trace, got: %s", stackTrace)
	}
}

func TestErrorExpr_DetailedMessage(t *testing.T) {
	// Test the GetDetailedMessage method
	frames := []StackFrame{
		{Function: "Plus", Expression: "Plus[1, 2]", Location: ""},
		{Function: "Divide", Expression: "Divide[1, 0]", Location: ""},
	}
	
	errorExpr := NewErrorExprWithStack("DivisionByZero", "Cannot divide by zero", []Expr{NewIntAtom(1), NewIntAtom(0)}, frames)
	
	detailed := errorExpr.GetDetailedMessage()
	
	// Check that all components are present
	if !strings.Contains(detailed, "DivisionByZero") {
		t.Error("detailed message should contain error type")
	}
	if !strings.Contains(detailed, "Cannot divide by zero") {
		t.Error("detailed message should contain error message")
	}
	if !strings.Contains(detailed, "Plus") {
		t.Error("detailed message should contain stack trace")
	}
	if !strings.Contains(detailed, "Divide") {
		t.Error("detailed message should contain stack trace")
	}
}