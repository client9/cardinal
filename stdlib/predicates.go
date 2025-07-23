package stdlib

import (
	"github.com/client9/sexpr/core"
)

// Type predicate functions - all return bool

// IntegerQExpr checks if an expression is an integer
func IntegerQExpr(expr core.Expr) bool {
	if atom, ok := expr.(core.Atom); ok {
		return atom.AtomType == core.IntAtom
	}
	return false
}

// FloatQExpr checks if an expression is a float
func FloatQExpr(expr core.Expr) bool {
	if atom, ok := expr.(core.Atom); ok {
		return atom.AtomType == core.FloatAtom
	}
	return false
}

// NumberQExpr checks if an expression is numeric (int or float)
func NumberQExpr(expr core.Expr) bool {
	return core.IsNumeric(expr)
}

// StringQExpr checks if an expression is a string
func StringQExpr(expr core.Expr) bool {
	if atom, ok := expr.(core.Atom); ok {
		return atom.AtomType == core.StringAtom
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
	_, isAtom := expr.(core.Atom)
	return isAtom
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
func HeadExpr(expr core.Expr) core.Expr {
	switch ex := expr.(type) {
	case core.Atom:
		switch ex.AtomType {
		case core.IntAtom:
			return core.NewSymbolAtom("Integer")
		case core.FloatAtom:
			return core.NewSymbolAtom("Real")
		case core.StringAtom:
			return core.NewSymbolAtom("String")
		case core.SymbolAtom:
			return core.NewSymbolAtom("Symbol")
		default:
			return core.NewSymbolAtom("Unknown")
		}
	case core.List:
		if len(ex.Elements) == 0 {
			return core.NewSymbolAtom("List")
		} else {
			// For non-empty lists, the head is the first element
			// This matches Mathematica semantics where f[x,y] has head f
			return ex.Elements[0]
		}
	case core.ObjectExpr:
		return core.NewSymbolAtom(ex.TypeName)
	default:
		return core.NewSymbolAtom("Unknown")
	}
	// Note: ErrorExpr is not handled here - wrapper will propagate errors
}
