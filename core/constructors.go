package core

// Constructor functions for core types

// New constructor functions for atomic types
func NewString(s string) String  { return String(s) }
func NewInteger(i int64) Integer { return Integer(i) }
func NewReal(f float64) Real     { return Real(f) }
func NewSymbol(s string) Symbol  { return Symbol(s) }

// NewSymbolNull creates the Null symbol
func NewSymbolNull() Symbol { return Symbol("Null") }

// NewBool creates a boolean symbol (True/False) for Mathematica compatibility
func NewBool(value bool) Symbol {
	if value {
		return Symbol("True")
	} else {
		return Symbol("False")
	}
}

// NewErrorExpr creates a new error expression
func NewErrorExpr(errorType, message string, args []Expr) *ErrorExpr {
	return &ErrorExpr{
		ErrorType:  errorType,
		Message:    message,
		Args:       args,
		StackTrace: []StackFrame{},
	}
}

// NewErrorExprWithStack creates a new error expression with stack trace
func NewErrorExprWithStack(errorType, message string, args []Expr, stack []StackFrame) *ErrorExpr {
	return &ErrorExpr{
		ErrorType:  errorType,
		Message:    message,
		Args:       args,
		StackTrace: stack,
	}
}

func NewList(elements ...Expr) List {
	return List{Elements: elements}
}

// NewObjectExpr creates a new ObjectExpr with the given type name and value
func NewObjectExpr(typeName string, value Expr) ObjectExpr {
	return ObjectExpr{TypeName: typeName, Value: value}
}
