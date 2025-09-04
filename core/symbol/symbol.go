package symbol

import (
	"unique"
)

// SymbolExpr is core symbol atom.
// It can't be named "Symbol" since it conflicts with the Symbol.Symbol literal.
//
// basically a pointer to string
type SymbolExpr unique.Handle[string]

func NewSymbol(s string) SymbolExpr {
	return SymbolExpr(unique.Make(s))
}

func (s SymbolExpr) String() string {
	return unique.Handle[string](s).Value()
}

func (s SymbolExpr) InputForm() string {
	return s.String()
}

func (s SymbolExpr) Head() Expr {
	return Symbol
}

func (s SymbolExpr) Length() int64 {
	return 0
}

func (s SymbolExpr) Equal(rhs Expr) bool {
	return s == rhs

	/*
	   	if other, ok := rhs.(Symbol); ok && s == other  {
	   		return true
	   	}

	   return false
	*/
}

func (s SymbolExpr) IsAtom() bool {
	return true
}
