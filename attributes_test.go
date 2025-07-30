package sexpr

import (
	"sort"
	"sync"
	"testing"

	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

func TestAttribute_String(t *testing.T) {
	tests := []struct {
		name      string
		attribute engine.Attribute
		expected  string
	}{
		{"HoldAll", engine.HoldAll, "HoldAll"},
		{"HoldFirst", engine.HoldFirst, "HoldFirst"},
		{"HoldRest", engine.HoldRest, "HoldRest"},
		{"Flat", engine.Flat, "Flat"},
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

// Test SetAttributes builtin function
func TestSetAttributesBuiltin(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		shouldError bool
	}{
		{
			name:     "Set single attribute",
			input:    "SetAttributes(myFunc, Protected)",
			expected: "Null",
		},
		{
			name:     "Set multiple attributes with list",
			input:    "SetAttributes(myFunc, List(Protected, HoldFirst))",
			expected: "Null",
		},
		{
			name:        "Error: invalid symbol",
			input:       "SetAttributes(42, Protected)",
			shouldError: true,
		},
		{
			name:        "Error: invalid attribute",
			input:       "SetAttributes(myFunc, InvalidAttribute)",
			shouldError: true,
		},
		{
			name:     "Wrong number of arguments returns unevaluated",
			input:    "SetAttributes(myFunc)",
			expected: "SetAttributes(myFunc)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()

			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := evaluator.Evaluate(expr)

			if tt.shouldError {
				if !core.IsError(result) {
					t.Errorf("Expected error, but got: %s", result.String())
				}
			} else {
				if core.IsError(result) {
					t.Errorf("Unexpected error: %s", result.String())
				} else if result.String() != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, result.String())
				}
			}
		})
	}
}

// Test ClearAttributes builtin function
func TestClearAttributesBuiltin(t *testing.T) {
	tests := []struct {
		name        string
		setup       []string
		input       string
		expected    string
		shouldError bool
	}{
		{
			name: "Clear single attribute",
			setup: []string{
				"SetAttributes(testFunc, List(Protected, HoldFirst))",
			},
			input:    "ClearAttributes(testFunc, Protected)",
			expected: "Null",
		},
		{
			name:        "Error: invalid symbol",
			input:       "ClearAttributes(42, Protected)",
			shouldError: true,
		},
		{
			name:     "No arguments returns unevaluated",
			input:    "ClearAttributes()",
			expected: "ClearAttributes()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()

			// Run setup commands
			for _, setupCmd := range tt.setup {
				setupExpr, err := ParseString(setupCmd)
				if err != nil {
					t.Fatalf("Setup parse error: %v", err)
				}
				evaluator.Evaluate(setupExpr)
			}

			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := evaluator.Evaluate(expr)

			if tt.shouldError {
				if !core.IsError(result) {
					t.Errorf("Expected error, but got: %s", result.String())
				}
			} else {
				if core.IsError(result) {
					t.Errorf("Unexpected error: %s", result.String())
				} else if result.String() != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, result.String())
				}
			}
		})
	}
}

// Test Attributes builtin function
func TestAttributesBuiltin(t *testing.T) {
	tests := []struct {
		name        string
		setup       []string
		input       string
		expected    string
		shouldError bool
	}{
		{
			name:     "Get attributes of symbol with no attributes",
			input:    "Attributes(newSymbol)",
			expected: "List()",
		},
		{
			name: "Get attributes of symbol with single attribute",
			setup: []string{
				"SetAttributes(testFunc, Protected)",
			},
			input:    "Attributes(testFunc)",
			expected: "List(Protected)",
		},
		{
			name: "Get attributes of symbol with multiple attributes",
			setup: []string{
				"SetAttributes(testFunc, List(Protected, HoldFirst, Constant))",
			},
			input:    "Attributes(testFunc)",
			expected: "List(Constant, HoldFirst, Protected)", // Alphabetically sorted
		},
		{
			name:     "Get attributes of builtin function",
			input:    "Attributes(Plus)",
			expected: "List(Flat, Listable, NumericFunction, OneIdentity, Orderless, Protected)",
		},
		{
			name:        "Error: invalid symbol",
			input:       "Attributes(42)",
			shouldError: true,
		},
		{
			name:     "Wrong number of arguments returns unevaluated",
			input:    "Attributes()",
			expected: "Attributes()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()

			// Run setup commands
			for _, setupCmd := range tt.setup {
				setupExpr, err := ParseString(setupCmd)
				if err != nil {
					t.Fatalf("Setup parse error: %v", err)
				}
				evaluator.Evaluate(setupExpr)
			}

			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := evaluator.Evaluate(expr)

			if tt.shouldError {
				if !core.IsError(result) {
					t.Errorf("Expected error, but got: %s", result.String())
				}
			} else {
				if core.IsError(result) {
					t.Errorf("Unexpected error: %s", result.String())
				} else if result.String() != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, result.String())
				}
			}
		})
	}
}

// Test integration of all attribute functions
func TestAttributeFunctionsIntegration(t *testing.T) {
	evaluator := NewEvaluator()

	// Test setting attributes
	expr1, _ := ParseString("SetAttributes(myFunc, List(Protected, HoldFirst))")
	result1 := evaluator.Evaluate(expr1)
	if result1.String() != "Null" {
		t.Errorf("SetAttributes failed: %s", result1.String())
	}

	// Test getting attributes
	expr2, _ := ParseString("Attributes(myFunc)")
	result2 := evaluator.Evaluate(expr2)
	expected2 := "List(HoldFirst, Protected)" // Should be sorted
	if result2.String() != expected2 {
		t.Errorf("Attributes: expected %s, got %s", expected2, result2.String())
	}

	// Test clearing specific attribute
	expr3, _ := ParseString("ClearAttributes(myFunc, Protected)")
	result3 := evaluator.Evaluate(expr3)
	if result3.String() != "Null" {
		t.Errorf("ClearAttributes failed: %s", result3.String())
	}

	// Verify attribute was cleared
	expr4, _ := ParseString("Attributes(myFunc)")
	result4 := evaluator.Evaluate(expr4)
	expected4 := "List(HoldFirst)"
	if result4.String() != expected4 {
		t.Errorf("After clearing: expected %s, got %s", expected4, result4.String())
	}

	// Test clearing all attributes
	expr5, _ := ParseString("ClearAttributes(myFunc, Attributes(myFunc))")
	result5 := evaluator.Evaluate(expr5)
	if result5.String() != "Null" {
		t.Errorf("ClearAttributes all failed: %s", result5.String())
	}

	// Verify all attributes were cleared
	expr6, _ := ParseString("Attributes(myFunc)")
	result6 := evaluator.Evaluate(expr6)
	expected6 := "List()"
	if result6.String() != expected6 {
		t.Errorf("After clearing all: expected %s, got %s", expected6, result6.String())
	}
}

// TestProtectedAttributeEnforcement tests that Protected symbols cannot be reassigned
func TestProtectedAttributeEnforcement(t *testing.T) {
	tests := []struct {
		name        string
		setup       []string
		input       string
		expectError bool
	}{
		{
			name: "Protected symbol cannot be assigned with Set",
			setup: []string{
				"x = 5",
				"SetAttributes(x, Protected)",
			},
			input:       "x = 10",
			expectError: true,
		},
		{
			name: "Protected symbol cannot be assigned with SetDelayed",
			setup: []string{
				"y = 3",
				"SetAttributes(y, Protected)",
			},
			input:       "y := 42",
			expectError: true,
		},
		{
			name:        "Built-in Protected symbol cannot be reassigned",
			input:       "Plus = 42",
			expectError: true,
		},
		{
			name:        "Unprotected symbol can be reassigned",
			input:       "z = 5; z = 10; z",
			expectError: false,
		},
		{
			name: "Symbol can be reassigned after clearing Protected",
			setup: []string{
				"w = 1",
				"SetAttributes(w, Protected)",
				"ClearAttributes(w, Protected)",
			},
			input:       "w = 2",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewEvaluator()

			// Run setup commands
			for _, setupCmd := range tt.setup {
				expr, err := ParseString(setupCmd)
				if err != nil {
					t.Fatalf("Parse error in setup: %v", err)
				}
				result := evaluator.Evaluate(expr)
				if core.IsError(result) {
					t.Fatalf("Setup error: %s", result.String())
				}
			}

			// Parse and evaluate the test input
			expr, err := ParseString(tt.input)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := evaluator.Evaluate(expr)

			if tt.expectError {
				if !core.IsError(result) {
					t.Errorf("Expected error, but got: %s", result.String())
				} else {
					// Check that it's specifically a ProtectionError
					if errorExpr, ok := result.(*core.ErrorExpr); ok {
						if errorExpr.ErrorType != "ProtectionError" {
							t.Errorf("Expected ProtectionError, got %s: %s", errorExpr.ErrorType, result.String())
						}
					}
				}
			} else {
				if core.IsError(result) {
					t.Errorf("Unexpected error: %s", result.String())
				}
			}
		})
	}
}
