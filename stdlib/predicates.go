package stdlib

import (
	"github.com/client9/sexpr/core"
)

// Type predicate functions - all return bool

// IntegerQExpr checks if an expression is an integer
func IntegerQExpr(expr core.Expr) bool {
	// Check new Integer type first
	if _, ok := expr.(core.Integer); ok {
		return true
	}
	return false
}

// FloatQExpr checks if an expression is a float
func FloatQExpr(expr core.Expr) bool {
	// Check new Real type first
	if _, ok := expr.(core.Real); ok {
		return true
	}
	return false
}

// NumberQExpr checks if an expression is numeric (int or float)
func NumberQExpr(expr core.Expr) bool {
	return core.IsNumeric(expr)
}

// StringQExpr checks if an expression is a string
func StringQExpr(expr core.Expr) bool {
	// Check new String type first
	if _, ok := expr.(core.String); ok {
		return true
	}
	return false
}

// BooleanQExpr checks if an expression is a boolean (True/False symbol)
func BooleanQExpr(expr core.Expr) bool {
	return core.IsBool(expr)
}

// SymbolQExpr checks if an expression is a symbol
func SymbolQExpr(expr core.Expr) bool {
	return core.IsSymbol(expr)
}

// ListQExpr checks if an expression is a list
func ListQExpr(expr core.Expr) bool {
	_, isList := expr.(core.List)
	return isList
}

// AtomQExpr checks if an expression is an atom
func AtomQExpr(expr core.Expr) bool {
	return expr.IsAtom()
}

// TrueQExpr check is an expression is explicity True
func TrueQExpr(expr core.Expr) bool {
	// Check new Symbol type first
	if symbolName, ok := core.ExtractSymbol(expr); ok {
		return symbolName == "True"
	}
	return false
}

// Output format functions - all return string

// FullFormExpr returns the full string representation of an expression
func FullFormExpr(expr core.Expr) string {
	// For now, just return the string representation
	// Pattern conversion logic will be added when patterns are moved to core
	return expr.String()
}

// InputFormExpr returns the user-friendly InputForm representation of an expression
func InputFormExpr(expr core.Expr) string {
	// For now, just return the InputForm representation
	// Pattern conversion logic will be added when patterns are moved to core
	return expr.InputForm()
}

// HeadExpr returns the head/type of an expression
// expr.Type() returns a string, but Head returns a symbol
// Head(Head("foo")) is Symbol("String")
//
func HeadExpr(expr core.Expr) core.Expr {
	return core.NewSymbol(expr.Type())
}
