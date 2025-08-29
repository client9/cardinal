package core

import (
	"unique"
)

// basically a pointer to string
type Symbol unique.Handle[string]

func NewSymbol(s string) Symbol {
	return Symbol(unique.Make(s))
}

func (s Symbol) String() string {
	return unique.Handle[string](s).Value()
}

func (s Symbol) InputForm() string {
	return s.String()
}

func (s Symbol) HeadExpr() Symbol {
	return symbolSymbol
}

func (s Symbol) Length() int64 {
	return 0
}

func (s Symbol) Equal(rhs Expr) bool {
	return s == rhs

	/*
	   	if other, ok := rhs.(Symbol); ok && s == other  {
	   		return true
	   	}

	   return false
	*/
}

func (s Symbol) IsAtom() bool {
	return true
}
