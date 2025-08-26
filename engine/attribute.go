package engine

//go:generate stringer -type=Attribute

import (
	"sort"

	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/atom"
)

// Attribute represents a symbol attribute in Mathematica-style
type Attribute int

// in alphabetical order
const (
	Constant Attribute = 1 << iota
	Flat
	HoldAll
	HoldFirst
	HoldRest
	Listable
	Locked
	NumericFunction
	OneIdentity
	Orderless
	Protected
	ReadProtected
	Temporary
	AttributeLast
)

var symbolToAttribute map[atom.Atom]Attribute

// stringToAttribute provides reverse lookup from string to Attribute
// This map is automatically populated using the stringer-generated String() method
var atomToAttribute = map[atom.Atom]Attribute{
	atom.HoldAll:         HoldAll,
	atom.HoldFirst:       HoldFirst,
	atom.HoldRest:        HoldRest,
	atom.Flat:            Flat,
	atom.OneIdentity:     OneIdentity,
	atom.Orderless:       Orderless,
	atom.Listable:        Listable,
	atom.Constant:        Constant,
	atom.NumericFunction: NumericFunction,
	atom.Protected:       Protected,
	atom.ReadProtected:   ReadProtected,
	atom.Locked:          Locked,
	atom.Temporary:       Temporary,
}

func SymbolToAttribute(e core.Expr) Attribute {
	if s, ok := e.(core.Symbol); ok {
		return atomToAttribute[s.Atom()]
	}
	return 0
}

func attributeToSymbol(a Attribute) core.Expr {
	for k, v := range atomToAttribute {
		if v == a {
			return core.SymbolFor(k)
		}
	}
	return nil
}

func AttributeToSymbols(a Attribute) []core.Expr {
	out := []core.Expr{}
	if a == 0 {
		return out
	}
	for i := Attribute(1); i < AttributeLast; i = i << 1 {
		if a&i == i {
			out = append(out, attributeToSymbol(Attribute(i)))
		}
	}
	return out
}

// SymbolTable manages attributes for symbols
type SymbolTable struct {
	attributes map[string]Attribute
}

// NewSymbolTable creates a new symbol table instance
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		attributes: make(map[string]Attribute),
	}
}

// Reset clears all attributes from the symbol table (useful for testing)
func (st *SymbolTable) Reset() {
	st.attributes = make(map[string]Attribute)
}

// SetAttributes sets one or more attributes for a symbol
func (st *SymbolTable) SetAttributes(symbol string, attrs Attribute) {
	alist := st.attributes[symbol]
	st.attributes[symbol] = alist | attrs
}

// ClearAttributes removes one or more attributes from a symbol
func (st *SymbolTable) ClearAttributes(symbol string, attrs Attribute) {
	alist := st.attributes[symbol]
	if alist == 0 {
		return
	}
	alist &^= attrs
	if alist == 0 {
		delete(st.attributes, symbol)
		return
	}
	st.attributes[symbol] = alist
}

// Attributes returns all attributes for a symbol
func (st *SymbolTable) Attributes(symbol string) Attribute {
	return st.attributes[symbol]
}

// HasAttribute checks if a symbol has a specific attribute
func (st *SymbolTable) HasAttribute(symbol string, attr Attribute) bool {
	alist := st.attributes[symbol]
	return alist&attr == attr
}

// ClearAllAttributes removes all attributes from a symbol
func (st *SymbolTable) ClearAllAttributes(symbol string) {
	delete(st.attributes, symbol)
}

// AllSymbolsWithAttributes returns all symbols that have attributes
func (st *SymbolTable) AllSymbolsWithAttributes() []string {
	var symbols []string
	for symbol := range st.attributes {
		if st.attributes[symbol] != 0 {
			symbols = append(symbols, symbol)
		}
	}

	sort.Strings(symbols)
	return symbols
}
