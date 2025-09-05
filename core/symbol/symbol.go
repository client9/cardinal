package symbol

import (
	"strconv"
	"unicode"
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
	v := unique.Handle[string](s).Value()
	if isSymbolLiteral(v) {
		return v
	}
	return "Symbol(" + strconv.Quote(v) + ")"
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

func isSymbolLiteral(s string) bool {
	for i, r := range s {
		if i == 0 {
			if !IsSymbolRuneFirst(r) {
				return false
			}
		} else if !IsSymbolRuneRest(r) {
			return false
		}
	}
	return true
}

func IsSymbolRuneFirst(r rune) bool {
	// symbols can't start with a number
	if r >= '0' && r <= '9' {
		return false
	}
	return IsSymbolRuneRest(r)
}

func IsSymbolRuneRest(r rune) bool {
	return unicode.IsPrint(r) && !unicode.IsSpace(r) && r != '_'
}
