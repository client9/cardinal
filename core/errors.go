package core

import "fmt"

// StackFrame represents a single frame in the evaluation stack
type StackFrame struct {
	Function   string // Function name being evaluated
	Expression string // String representation of the expression
	Location   string // Optional location information
}

// ErrorExpr represents an error that occurred during evaluation
type ErrorExpr struct {
	ErrorType  string       // "DivisionByZero", "ArgumentError", etc.
	Message    string       // Detailed error message
	Args       []Expr       // Arguments that caused the error
	StackTrace []StackFrame // Stack trace of evaluation frames
}

// Possible we expose the arguments above differently and use ObjectExpr
// for now return 0
func (e *ErrorExpr) Length() int64 {
	return 0
}

func (e *ErrorExpr) String() string {
	return fmt.Sprintf("$Failed(%s)", e.ErrorType)
}

func (e *ErrorExpr) InputForm() string {
	// For errors, InputForm is the same as String()
	return e.String()
}

func (e *ErrorExpr) Type() string {
	return "error"
}

func (e *ErrorExpr) IsAtom() bool {
	return false
}

func (e *ErrorExpr) Equal(rhs Expr) bool {
	rhsError, ok := rhs.(*ErrorExpr)
	if !ok {
		return false
	}

	// Compare error type and message
	if e.ErrorType != rhsError.ErrorType || e.Message != rhsError.Message {
		return false
	}

	// Compare argument lists
	if len(e.Args) != len(rhsError.Args) {
		return false
	}

	for i, arg := range e.Args {
		if !arg.Equal(rhsError.Args[i]) {
			return false
		}
	}

	// Note: We don't compare stack traces as they are context-dependent
	return true
}
