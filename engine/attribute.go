package engine

//go:generate go run golang.org/x/tools/cmd/stringer -type=Attribute

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

// Attribute represents a symbol attribute in Mathematica-style
type Attribute int

const (
	HoldAll Attribute = iota
	HoldFirst
	HoldRest
	Flat
	Orderless
	OneIdentity
	Listable
	Constant
	NumericFunction
	Protected
	ReadProtected
	Locked
	Temporary
)

// stringToAttribute provides reverse lookup from string to Attribute
// This map is automatically populated using the stringer-generated String() method
var stringToAttribute map[string]Attribute

// StringToAttribute converts a string name to an attribute
func StringToAttribute(name string) (Attribute, bool) {
	if stringToAttribute == nil {
		initStringToAttributeMap()
	}
	attr, ok := stringToAttribute[name]
	return attr, ok
}

// initStringToAttributeMap initializes the reverse lookup map using stringer output
func initStringToAttributeMap() {
	stringToAttribute = make(map[string]Attribute)
	// Iterate through all possible attribute values
	for i := HoldAll; i <= Temporary; i++ {
		stringToAttribute[i.String()] = i
	}
}

// SymbolTable manages attributes for symbols
type SymbolTable struct {
	attributes map[string][]Attribute
}

// NewSymbolTable creates a new symbol table instance
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		attributes: make(map[string][]Attribute),
	}
}

// Reset clears all attributes from the symbol table (useful for testing)
func (st *SymbolTable) Reset() {
	st.attributes = make(map[string][]Attribute)
}

// SetAttributes sets one or more attributes for a symbol
func (st *SymbolTable) SetAttributes(symbol string, attrs []Attribute) {

	alist := st.attributes[symbol]

	// if doesn't exist, just set directly
	if alist == nil {
		st.attributes[symbol] = attrs
		return
	}

	dirty := false
	for _, attr := range attrs {
		if !slices.Contains(alist, attr) {
			alist = append(alist, attr)
			dirty = true
		}
	}
	if dirty {
		st.attributes[symbol] = alist
	}
}

// ClearAttributes removes one or more attributes from a symbol
func (st *SymbolTable) ClearAttributes(symbol string, attrs []Attribute) {
	alist := st.attributes[symbol]
	if alist == nil {
		return
	}

	for _, attr := range attrs {
		if idx := slices.Index(alist, attr); idx != -1 {
			alist = slices.Delete(alist, idx, idx+1)
		}
	}
	if len(alist) == 0 {
		delete(st.attributes, symbol)
	}
}

// Attributes returns all attributes for a symbol
func (st *SymbolTable) Attributes(symbol string) []Attribute {
	alist := st.attributes[symbol]

	// Sort for consistent output (alphabetically by name)
	sort.Slice(alist, func(i, j int) bool {
		return alist[i].String() < alist[j].String()
	})

	return alist
}

// HasAttribute checks if a symbol has a specific attribute
func (st *SymbolTable) HasAttribute(symbol string, attr Attribute) bool {

	alist := st.attributes[symbol]

	if alist == nil {
		return false
	}
	return slices.Contains(alist, attr)
}

// ClearAllAttributes removes all attributes from a symbol
func (st *SymbolTable) ClearAllAttributes(symbol string) {

	delete(st.attributes, symbol)
}

// AttributesToString returns a formatted string representation of attributes
func AttributesToString(attrs []Attribute) string {
	if len(attrs) == 0 {
		return "{}"
	}

	var names []string
	for _, attr := range attrs {
		names = append(names, attr.String())
	}

	return fmt.Sprintf("{%s}", strings.Join(names, ", "))
}

// AllSymbolsWithAttributes returns all symbols that have attributes
func (st *SymbolTable) AllSymbolsWithAttributes() []string {
	var symbols []string
	for symbol := range st.attributes {
		if len(st.attributes[symbol]) > 0 {
			symbols = append(symbols, symbol)
		}
	}

	sort.Strings(symbols)
	return symbols
}
