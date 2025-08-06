package core

// Constructor functions for core types

// New constructor functions for atomic types
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

// NewObjectExpr creates a new ObjectExpr with the given type name and value
func NewObjectExpr(typeName string, value Expr) ObjectExpr {
	return ObjectExpr{TypeName: typeName, Value: value}
}
