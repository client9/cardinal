package core

import (
	"fmt"
	"strings"
)

// FunctionExpr represents a pure/lambda function created by Function(args, body)
type FunctionExpr struct {
	Parameters []string // Parameter names (e.g., ["x", "y"])
	Body       Expr     // Function body (held unevaluated)
}

// String returns the string representation of the function
func (f FunctionExpr) String() string {
	if len(f.Parameters) == 0 {
		// Constant function: Function([], body)
		return fmt.Sprintf("Function([], %s)", f.Body.String())
	}

	if len(f.Parameters) == 1 {
		// Single parameter: Function(x, body)
		return fmt.Sprintf("Function(%s, %s)", f.Parameters[0], f.Body.String())
	}

	// Multiple parameters: Function([x, y], body)
	paramList := "[" + strings.Join(f.Parameters, ", ") + "]"
	return fmt.Sprintf("Function(%s, %s)", paramList, f.Body.String())
}

// InputForm returns the input form representation
func (f FunctionExpr) InputForm() string {
	return f.String() // Same as String for now
}

// Head returns the head of the expression
func (f FunctionExpr) Head() string {
	return "Function"
}

// Length returns the length (number of parameters + 1 for body)
func (f FunctionExpr) Length() int64 {
	return int64(len(f.Parameters)) + 1
}

// Equal checks equality with another expression
func (f FunctionExpr) Equal(rhs Expr) bool {
	if other, ok := rhs.(FunctionExpr); ok {
		// Check if parameters match
		if len(f.Parameters) != len(other.Parameters) {
			return false
		}
		for i, param := range f.Parameters {
			if param != other.Parameters[i] {
				return false
			}
		}
		// Check if bodies are equal
		return f.Body.Equal(other.Body)
	}
	return false
}

// IsAtom returns false since functions are composite
func (f FunctionExpr) IsAtom() bool {
	return false
}

// NewFunction creates a new FunctionExpr
func NewFunction(parameters []string, body Expr) FunctionExpr {
	return FunctionExpr{
		Parameters: parameters,
		Body:       body,
	}
}
