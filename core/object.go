package core

// ObjectExpr wraps a user-defined Expr implementation with a type name
// This allows users to register custom Go types that implement Expr
// and use them with pattern matching (e.g., x_Uint64)
type ObjectExpr struct {
	TypeName string // e.g., "Uint64", "BigInt", "Matrix"
	Value    Expr   // User-defined type that implements Expr interface
}

func (o ObjectExpr) String() string {
	return o.Value.String() // Delegate to the wrapped Expr
}

func (o ObjectExpr) InputForm() string {
	// Delegate to the wrapped Expr's InputForm if it has one,
	// otherwise fall back to String()
	return o.Value.InputForm()
}

func (o ObjectExpr) Length() int64 {
	return o.Value.Length() // Delegate to wrapper Expr
}

func (o ObjectExpr) Head() string {
	return o.TypeName // Return the registered type name
}

func (o ObjectExpr) IsAtom() bool {
	return false // ObjectExpr is a wrapper, so it's not atomic
}

func (o ObjectExpr) Equal(rhs Expr) bool {
	rhsObj, ok := rhs.(ObjectExpr)
	if !ok || o.TypeName != rhsObj.TypeName {
		return false
	}
	return o.Value.Equal(rhsObj.Value) // Delegate to wrapped Expr
}
