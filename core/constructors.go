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

// NewList creates a new List with a Symbol head and arguments
// This reflects the s-expression semantics where lists are function calls
func NewList(head string, args ...Expr) List {
	elements := make([]Expr, len(args)+1)
	elements[0] = NewSymbol(head)
	copy(elements[1:], args)
	return List{Elements: elements}
}

// NewListFromExprs creates a List directly from expressions (for special cases)
// Use NewList instead when possible, as it enforces the Symbol-head convention
func NewListFromExprs(elements ...Expr) List {
	return List{Elements: elements}
}

// NewObjectExpr creates a new ObjectExpr with the given type name and value
func NewObjectExpr(typeName string, value Expr) ObjectExpr {
	return ObjectExpr{TypeName: typeName, Value: value}
}
