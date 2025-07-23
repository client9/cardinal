package sexpr

import (
	"fmt"
	"sort"
	"strings"
	"sync"
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

// AttributeNames maps attribute constants to their string representations
var AttributeNames = map[Attribute]string{
	HoldAll:         "HoldAll",
	HoldFirst:       "HoldFirst",
	HoldRest:        "HoldRest",
	Flat:            "Flat",
	Orderless:       "Orderless",
	OneIdentity:     "OneIdentity",
	Listable:        "Listable",
	Constant:        "Constant",
	NumericFunction: "NumericFunction",
	Protected:       "Protected",
	ReadProtected:   "ReadProtected",
	Locked:          "Locked",
	Temporary:       "Temporary",
}

// StringToAttribute converts a string name to an attribute
func StringToAttribute(name string) (Attribute, bool) {
	for attr, attrName := range AttributeNames {
		if attrName == name {
			return attr, true
		}
	}
	return 0, false
}

// String returns the string representation of an attribute
func (a Attribute) String() string {
	if name, ok := AttributeNames[a]; ok {
		return name
	}
	return fmt.Sprintf("Attribute(%d)", int(a))
}

// SymbolTable manages attributes for symbols
type SymbolTable struct {
	attributes map[string]map[Attribute]bool
	mu         sync.RWMutex
}

// NewSymbolTable creates a new symbol table instance
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		attributes: make(map[string]map[Attribute]bool),
	}
}

// SetAttributes sets one or more attributes for a symbol
func (st *SymbolTable) SetAttributes(symbol string, attrs []Attribute) {
	st.mu.Lock()
	defer st.mu.Unlock()

	if st.attributes[symbol] == nil {
		st.attributes[symbol] = make(map[Attribute]bool)
	}

	for _, attr := range attrs {
		st.attributes[symbol][attr] = true
	}
}

// ClearAttributes removes one or more attributes from a symbol
func (st *SymbolTable) ClearAttributes(symbol string, attrs []Attribute) {
	st.mu.Lock()
	defer st.mu.Unlock()

	if st.attributes[symbol] == nil {
		return
	}

	for _, attr := range attrs {
		delete(st.attributes[symbol], attr)
	}

	// Clean up empty attribute maps
	if len(st.attributes[symbol]) == 0 {
		delete(st.attributes, symbol)
	}
}

// Attributes returns all attributes for a symbol
func (st *SymbolTable) Attributes(symbol string) []Attribute {
	st.mu.RLock()
	defer st.mu.RUnlock()

	if st.attributes[symbol] == nil {
		return nil
	}

	var attrs []Attribute
	for attr := range st.attributes[symbol] {
		attrs = append(attrs, attr)
	}

	// Sort for consistent output (alphabetically by name)
	sort.Slice(attrs, func(i, j int) bool {
		return attrs[i].String() < attrs[j].String()
	})

	return attrs
}

// HasAttribute checks if a symbol has a specific attribute
func (st *SymbolTable) HasAttribute(symbol string, attr Attribute) bool {
	st.mu.RLock()
	defer st.mu.RUnlock()

	if st.attributes[symbol] == nil {
		return false
	}

	return st.attributes[symbol][attr]
}

// ClearAllAttributes removes all attributes from a symbol
func (st *SymbolTable) ClearAllAttributes(symbol string) {
	st.mu.Lock()
	defer st.mu.Unlock()

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
	st.mu.RLock()
	defer st.mu.RUnlock()

	var symbols []string
	for symbol := range st.attributes {
		if len(st.attributes[symbol]) > 0 {
			symbols = append(symbols, symbol)
		}
	}

	sort.Strings(symbols)
	return symbols
}

// Reset clears all attributes from the symbol table (useful for testing)
func (st *SymbolTable) Reset() {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.attributes = make(map[string]map[Attribute]bool)
}

// parseAttribute parses an attribute name string into an Attribute enum value
func parseAttribute(attrName string) (Attribute, error) {
	switch attrName {
	case "HoldAll":
		return HoldAll, nil
	case "HoldFirst":
		return HoldFirst, nil
	case "HoldRest":
		return HoldRest, nil
	case "Flat":
		return Flat, nil
	case "Orderless":
		return Orderless, nil
	case "OneIdentity":
		return OneIdentity, nil
	case "Listable":
		return Listable, nil
	case "Constant":
		return Constant, nil
	case "NumericFunction":
		return NumericFunction, nil
	case "Protected":
		return Protected, nil
	case "ReadProtected":
		return ReadProtected, nil
	case "Locked":
		return Locked, nil
	case "Temporary":
		return Temporary, nil
	default:
		return HoldAll, fmt.Errorf("unknown attribute: %s", attrName)
	}
}
