package core

import (
	"github.com/client9/sexpr/core/atom"
)

type Symbol struct {
	atom atom.Atom
	name string
}

var symbolTrue Symbol
var symbolFalse Symbol
var symbolNull Symbol
var symbolList Symbol

func init() {
	atomTrue := atom.Lookup("True")
	atomFalse := atom.Lookup("False")
	atomNull := atom.Lookup("Null")
	atomList := atom.Lookup("List")
	symbolTrue = Symbol{atom: atomTrue, name: atomTrue.String()}
	symbolFalse = Symbol{atom: atomFalse, name: atomFalse.String()}
	symbolNull = Symbol{atom: atomNull, name: atomNull.String()}
	symbolList = Symbol{atom: atomList, name: atomList.String()}
}

// SymbolFor makes a symbol with the given atom "SymbolFor(atom.Xyz)"
func SymbolFor(a atom.Atom) Symbol {
	// TODO -- why even bother with the name/string?
	return Symbol{
		atom: a,
		name: a.String(),
	}
}

func NewSymbol(s string) Symbol {
	a := atom.Lookup(s)
	if a == 0 {
		return Symbol{atom: 0, name: s}
	}
	return Symbol{atom: a, name: a.String()}
}

// NewSymbolNull creates the Null symbol
func NewSymbolNull() Symbol { return symbolNull }

// Symbol type implementation
func (s Symbol) String() string {
	return s.name
}

func (s Symbol) InputForm() string {
	return s.String()
}

func (s Symbol) Head() string {
	return "Symbol"
}

func (s Symbol) Length() int64 {
	return 0
}

func (s Symbol) Equal(rhs Expr) bool {
	other, ok := rhs.(Symbol)
	if !ok || other.atom != s.atom {
		return false
	}
	if other.atom == 0 {
		return other.name == s.name
	}
	return true
}

func (s Symbol) IsAtom() bool {
	return true
}
