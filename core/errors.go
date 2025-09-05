package core

import (
	"fmt"
	"slices"

	"github.com/client9/cardinal/core/symbol"
)

// ErrorExpr represents an error that occurred during evaluation
type ErrorExpr struct {
	ErrorType string // "DivisionByZero", "ArgumentError", etc.
	Message   string // Detailed error message
	Arg       Expr
	Err       *ErrorExpr
}

func NewError(etype string, message string) ErrorExpr {
	return ErrorExpr{
		ErrorType: etype,
		Message:   message,
	}
}

func (e ErrorExpr) SetCaller(arg Expr) ErrorExpr {
	e.Arg = arg
	return e
}

// AsError returns the ErrorExpr or nil.
//
// For error stack traces do something like:
//
//	if err, ok := AsError(expr); ok {
//	    return WrapError(err, arg);
//	}
func AsError(arg Expr) (ErrorExpr, bool) {
	err, ok := arg.(ErrorExpr)
	return err, ok
}

func (e ErrorExpr) Wrap(arg Expr) Expr {
	return ErrorExpr{
		ErrorType: "Stack",
		Message:   "tracing the stack",
		Arg:       arg,
		Err:       &e,
	}
}

func (e ErrorExpr) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorType, e.Message)
}

func (e ErrorExpr) Unwrap() error {
	return e.Err
}

// StackTrace unwrapps the error chain, return deepest error (the originator) first.
// Output always has at least one element
func (e ErrorExpr) StackTrace() []ErrorExpr {
	out := []ErrorExpr{e}
	for e.Err != nil {
		out = append(out, *e.Err)
		e = *e.Err
	}
	slices.Reverse(out)
	return out
}

// Length of an Error is 0 (zero).
func (e ErrorExpr) Length() int64 {
	return 0
}

func (e ErrorExpr) String() string {
	return fmt.Sprintf("$Failed(%s)", e.ErrorType)
}

func (e ErrorExpr) InputForm() string {
	// For errors, InputForm is the same as String()
	return e.String()
}

func (e ErrorExpr) Head() Expr {
	return symbol.Error
}

func (e ErrorExpr) IsAtom() bool {
	return true // ErrorExpr is atomic in symbolic computation - a complete value that doesn't need re-evaluation
}

func (e ErrorExpr) Equal(rhs Expr) bool {
	rhsError, ok := rhs.(ErrorExpr)
	if !ok {
		return false
	}

	// Compare error type and message
	if e.ErrorType != rhsError.ErrorType {
		return false
	}
	if e.Arg == nil && rhsError.Arg == nil {
		return true
	}

	if e.Arg == nil || rhsError.Arg == nil {
		return false
	}
	return e.Arg.Equal(rhsError.Arg)
}
