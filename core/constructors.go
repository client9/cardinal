package core

// Constructor functions for core types

// NewBool creates a boolean symbol (True/False) for Mathematica compatibility
// NewObjectExpr creates a new ObjectExpr with the given type name and value
func NewObjectExpr(typeName string, value Expr) ObjectExpr {
	return ObjectExpr{TypeName: typeName, Value: value}
}
