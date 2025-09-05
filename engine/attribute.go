package engine

//go:generate stringer -type=Attribute

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/symbol"
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
	NHoldFirst
	NHoldRest
	NumericFunction
	OneIdentity
	Orderless
	Protected
	ReadProtected
	Temporary
	AttributeLast
)

// stringToAttribute provides reverse lookup from string to Attribute
// This map is automatically populated using the stringer-generated String() method
var symbolToAttribute = map[core.Symbol]Attribute{
	symbol.HoldAll:         HoldAll,
	symbol.HoldFirst:       HoldFirst,
	symbol.HoldRest:        HoldRest,
	symbol.Flat:            Flat,
	symbol.OneIdentity:     OneIdentity,
	symbol.Orderless:       Orderless,
	symbol.Listable:        Listable,
	symbol.Constant:        Constant,
	symbol.NumericFunction: NumericFunction,
	symbol.Protected:       Protected,
	symbol.ReadProtected:   ReadProtected,
	symbol.Locked:          Locked,
	symbol.Temporary:       Temporary,
	symbol.NHoldFirst:      NHoldFirst,
	symbol.NHoldRest:       NHoldRest,
}

func SymbolToAttribute(e core.Expr) Attribute {
	if s, ok := e.(core.Symbol); ok {
		return symbolToAttribute[s]
	}
	return 0
}

// scan values to return key
func attributeToSymbol(a Attribute) core.Expr {
	for k, v := range symbolToAttribute {
		if v == a {
			return k
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
	attributes map[core.Symbol]Attribute
}

// NewSymbolTable creates a new symbol table instance
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		attributes: make(map[core.Symbol]Attribute),
	}
}

// Reset clears all attributes from the symbol table (useful for testing)
func (st *SymbolTable) Reset() {
	st.attributes = make(map[core.Symbol]Attribute)
}

// SetAttributes sets one or more attributes for a symbol
func (st *SymbolTable) SetAttributes(symbol core.Symbol, attrs Attribute) {
	alist := st.attributes[symbol]
	st.attributes[symbol] = alist | attrs
}

// ClearAttributes removes one or more attributes from a symbol
func (st *SymbolTable) ClearAttributes(symbol core.Symbol, attrs Attribute) {
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
func (st *SymbolTable) Attributes(symbol core.Symbol) Attribute {
	return st.attributes[symbol]
}

// HasAttribute checks if a symbol has a specific attribute
func (st *SymbolTable) HasAttribute(symbol core.Symbol, attr Attribute) bool {
	alist := st.attributes[symbol]
	return alist&attr == attr
}

// ClearAllAttributes removes all attributes from a symbol
func (st *SymbolTable) ClearAllAttributes(symbol core.Symbol) {
	delete(st.attributes, symbol)
}

// AllSymbolsWithAttributes returns all symbols that have attributes
func (st *SymbolTable) AllSymbolsWithAttributes() []core.Symbol {
	// TODO SORT
	var symbols []core.Symbol
	for sym := range st.attributes {
		if st.attributes[sym] != 0 {
			symbols = append(symbols, sym)
		}
	}

	//sort.Strings(symbols)
	return symbols
}
