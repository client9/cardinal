package sexpr

import (
	"sort"
	"sync"
	"testing"
)

func TestAttribute_String(t *testing.T) {
	tests := []struct {
		name      string
		attribute Attribute
		expected  string
	}{
		{"HoldAll", HoldAll, "HoldAll"},
		{"HoldFirst", HoldFirst, "HoldFirst"},
		{"HoldRest", HoldRest, "HoldRest"},
		{"Flat", Flat, "Flat"},
		{"Orderless", Orderless, "Orderless"},
		{"OneIdentity", OneIdentity, "OneIdentity"},
		{"Listable", Listable, "Listable"},
		{"Constant", Constant, "Constant"},
		{"NumericFunction", NumericFunction, "NumericFunction"},
		{"Protected", Protected, "Protected"},
		{"ReadProtected", ReadProtected, "ReadProtected"},
		{"Locked", Locked, "Locked"},
		{"Temporary", Temporary, "Temporary"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.attribute.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSetAttributes(t *testing.T) {
	tests := []struct {
		name       string
		symbol     string
		attributes []Attribute
		expected   []Attribute
	}{
		{
			name:       "single attribute",
			symbol:     "Plus",
			attributes: []Attribute{Flat},
			expected:   []Attribute{Flat},
		},
		{
			name:       "multiple attributes",
			symbol:     "Times",
			attributes: []Attribute{Flat, Orderless, OneIdentity},
			expected:   []Attribute{Flat, OneIdentity, Orderless}, // sorted
		},
		{
			name:       "adding to existing attributes",
			symbol:     "Plus",
			attributes: []Attribute{Orderless},
			expected:   []Attribute{Flat, Orderless}, // Plus already had Flat, now adding Orderless
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := NewSymbolTable()
			
			// For the "adding to existing" test, set up initial state
			if tt.name == "adding to existing attributes" {
				st.SetAttributes(tt.symbol, []Attribute{Flat})
			}
			
			st.SetAttributes(tt.symbol, tt.attributes)
			
			result := st.Attributes(tt.symbol)
			if !attributeSlicesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestClearAttributes(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*SymbolTable)
		symbol     string
		clearAttrs []Attribute
		expected   []Attribute
	}{
		{
			name: "clear single attribute",
			setup: func(st *SymbolTable) {
				st.SetAttributes("Plus", []Attribute{Flat, Orderless})
			},
			symbol:     "Plus",
			clearAttrs: []Attribute{Flat},
			expected:   []Attribute{Orderless},
		},
		{
			name: "clear multiple attributes",
			setup: func(st *SymbolTable) {
				st.SetAttributes("Times", []Attribute{Flat, Orderless, OneIdentity})
			},
			symbol:     "Times",
			clearAttrs: []Attribute{Flat, OneIdentity},
			expected:   []Attribute{Orderless},
		},
		{
			name: "clear non-existent attribute",
			setup: func(st *SymbolTable) {
				st.SetAttributes("Power", []Attribute{OneIdentity})
			},
			symbol:     "Power",
			clearAttrs: []Attribute{Flat},
			expected:   []Attribute{OneIdentity},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := NewSymbolTable()
			tt.setup(st)
			
			st.ClearAttributes(tt.symbol, tt.clearAttrs)
			
			result := st.Attributes(tt.symbol)
			if !attributeSlicesEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestHasAttribute(t *testing.T) {
	st := NewSymbolTable()
	st.SetAttributes("TestSymbol", []Attribute{Flat, Orderless})

	tests := []struct {
		name      string
		symbol    string
		attribute Attribute
		expected  bool
	}{
		{
			name:      "has attribute",
			symbol:    "TestSymbol",
			attribute: Flat,
			expected:  true,
		},
		{
			name:      "does not have attribute",
			symbol:    "TestSymbol",
			attribute: OneIdentity,
			expected:  false,
		},
		{
			name:      "non-existent symbol",
			symbol:    "NonExistent",
			attribute: Flat,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := st.HasAttribute(tt.symbol, tt.attribute)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestClearAllAttributes(t *testing.T) {
	st := NewSymbolTable()
	st.SetAttributes("TestSymbol", []Attribute{Flat, Orderless, OneIdentity})
	
	// Verify attributes are set
	if len(st.Attributes("TestSymbol")) == 0 {
		t.Fatal("Expected attributes to be set")
	}
	
	// Clear all attributes
	st.ClearAllAttributes("TestSymbol")
	
	// Verify all attributes are cleared
	result := st.Attributes("TestSymbol")
	if len(result) != 0 {
		t.Errorf("expected no attributes, got %v", result)
	}
}

func TestAttributesToString(t *testing.T) {
	tests := []struct {
		name       string
		attributes []Attribute
		expected   string
	}{
		{
			name:       "empty attributes",
			attributes: []Attribute{},
			expected:   "{}",
		},
		{
			name:       "single attribute",
			attributes: []Attribute{Flat},
			expected:   "{Flat}",
		},
		{
			name:       "multiple attributes",
			attributes: []Attribute{Flat, Orderless, OneIdentity},
			expected:   "{Flat, Orderless, OneIdentity}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AttributesToString(tt.attributes)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestAllSymbolsWithAttributes(t *testing.T) {
	st := NewSymbolTable()
	st.SetAttributes("Plus", []Attribute{Flat})
	st.SetAttributes("Times", []Attribute{Orderless})
	st.SetAttributes("Power", []Attribute{OneIdentity})
	
	symbols := st.AllSymbolsWithAttributes()
	
	expected := []string{"Plus", "Power", "Times"} // sorted
	if !stringSlicesEqual(symbols, expected) {
		t.Errorf("expected %v, got %v", expected, symbols)
	}
}

func TestSymbolTable_ThreadSafety(t *testing.T) {
	st := NewSymbolTable()
	var wg sync.WaitGroup
	
	// Test concurrent access
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			symbol := "TestSymbol"
			if n%2 == 0 {
				st.SetAttributes(symbol, []Attribute{Flat})
			} else {
				_ = st.HasAttribute(symbol, Flat)
			}
		}(i)
	}
	
	wg.Wait()
	// If we get here without data races, the test passed
}

func TestSymbolTable_Reset(t *testing.T) {
	st := NewSymbolTable()
	st.SetAttributes("TestSymbol", []Attribute{Flat, Orderless})
	
	// Verify attributes are set
	if len(st.AllSymbolsWithAttributes()) == 0 {
		t.Fatal("Expected symbols with attributes")
	}
	
	// Reset the symbol table
	st.Reset()
	
	// Verify all symbols are cleared
	symbols := st.AllSymbolsWithAttributes()
	if len(symbols) != 0 {
		t.Errorf("expected no symbols after reset, got %v", symbols)
	}
}

func TestMathematicaStyleUsage(t *testing.T) {
	st := NewSymbolTable()
	
	// Set up typical Mathematica built-ins
	st.SetAttributes("Plus", []Attribute{Flat, Orderless, OneIdentity})
	st.SetAttributes("Times", []Attribute{Flat, Orderless, OneIdentity})
	st.SetAttributes("Hold", []Attribute{HoldAll})
	st.SetAttributes("Sin", []Attribute{Listable, NumericFunction})
	st.SetAttributes("Pi", []Attribute{Constant, Protected})
	
	// Test typical queries
	if !st.HasAttribute("Plus", Flat) {
		t.Error("Plus should have Flat attribute")
	}
	
	// Test clearing protected symbols (should still work at the table level)
	st.ClearAllAttributes("Pi")
	if st.HasAttribute("Pi", Flat) {
		t.Error("Pi should not have any attributes after clearing")
	}
	
	// Test attribute combinations
	attrs := st.Attributes("Sin")
	expectedAttrs := []Attribute{Listable, NumericFunction}
	if !attributeSlicesEqual(attrs, expectedAttrs) {
		t.Errorf("expected %v, got %v", expectedAttrs, attrs)
	}
}

// Helper functions

func attributeSlicesEqual(a, b []Attribute) bool {
	if len(a) != len(b) {
		return false
	}
	
	// Sort both slices for comparison
	aCopy := make([]Attribute, len(a))
	bCopy := make([]Attribute, len(b))
	copy(aCopy, a)
	copy(bCopy, b)
	
	sort.Slice(aCopy, func(i, j int) bool {
		return aCopy[i] < aCopy[j]
	})
	sort.Slice(bCopy, func(i, j int) bool {
		return bCopy[i] < bCopy[j]
	})
	
	for i := range aCopy {
		if aCopy[i] != bCopy[i] {
			return false
		}
	}
	
	return true
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	
	return true
}