package core

// Constructor functions for core types

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

// NewStringAtom creates a new string atom
func NewStringAtom(value string) Atom {
	return Atom{AtomType: StringAtom, Value: value}
}

// NewIntAtom creates a new integer atom
func NewIntAtom(value int) Atom {
	return Atom{AtomType: IntAtom, Value: value}
}

// NewFloatAtom creates a new float atom
func NewFloatAtom(value float64) Atom {
	return Atom{AtomType: FloatAtom, Value: value}
}

// NewBoolAtom creates a boolean atom (returns True/False symbols)
func NewBoolAtom(value bool) Atom {
	// Return True/False symbols instead of BoolAtom for Mathematica compatibility
	// This makes our system fully symbolic like Mathematica
	if value {
		return NewSymbolAtom("True")
	} else {
		return NewSymbolAtom("False")
	}
}

// NewSymbolAtom creates a new symbol atom
func NewSymbolAtom(value string) Atom {
	return Atom{AtomType: SymbolAtom, Value: value}
}

// NewList creates a new list with the given elements
func NewList(elements ...Expr) List {
	return List{Elements: elements}
}

// NewObjectExpr creates a new ObjectExpr with the given type name and value
func NewObjectExpr(typeName string, value Expr) ObjectExpr {
	return ObjectExpr{TypeName: typeName, Value: value}
}