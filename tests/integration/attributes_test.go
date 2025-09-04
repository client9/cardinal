package integration

import (
	"testing"
)

// Test SetAttributes
func TestSetAttributes(t *testing.T) {
	tests := []TestCase{
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
			name:      "Error: invalid symbol",
			input:     "SetAttributes(42, Protected)",
			expected:  "SetAttributes(42, Protected)",
			errorType: "",
		},
		{
			name:      "Error: invalid attribute",
			input:     "SetAttributes(myFunc, InvalidAttribute)",
			expected:  "",
			errorType: "UnknownAttribute",
		},
		{
			name:     "Wrong number of arguments returns unevaluated",
			input:    "SetAttributes(myFunc)",
			expected: "SetAttributes(myFunc)",
		},
		{
			name:     "Get attributes of symbol with no attributes",
			input:    "Attributes(newSymbol)",
			expected: "List()",
		},
		{
			name:     "Get attributes of symbol with single attribute",
			input:    "SetAttributes(testFunc, Protected); Attributes(testFunc)",
			expected: "List(Protected)",
		},
		{
			name:     "Get attributes of symbol with multiple attributes",
			input:    "SetAttributes(testFunc, List(Protected, HoldFirst, Constant)); Attributes(testFunc)",
			expected: "List(Constant, HoldFirst, Protected)", // Alphabetically sorted
		},
		{
			name:     "Get attributes of builtin function",
			input:    "Attributes(Plus)",
			expected: "List(Flat, Listable, NumericFunction, OneIdentity, Orderless, Protected)",
		},
		{
			name:      "Attributes Error: invalid symbol",
			input:     "Attributes(42)",
			expected:  "Attributes(42)",
			errorType: "",
		},
		{
			name:     "Wrong number of arguments returns unevaluated",
			input:    "Attributes()",
			expected: "Attributes()",
		},
		{
			name:     "Clear single attribute",
			input:    "SetAttributes(testFunc, List(Protected, HoldFirst)); ClearAttributes(testFunc, Protected)",
			expected: "Null",
		},
		{
			name:      "Error: invalid symbol",
			input:     "ClearAttributes(42, Protected)",
			expected:  "ClearAttributes(42, Protected)",
			errorType: "",
		},
		{
			name:     "No arguments returns unevaluated",
			input:    "ClearAttributes()",
			expected: "ClearAttributes()",
		},
		{
			name:     "Clear all attributes with list",
			input:    "SetAttributes(x, [Flat, Listable, NumericFunction, OneIdentity, Orderless]); ClearAttributes(x, Attributes(x)); Length(Attributes(x))",
			expected: "0",
		},
	}
	runTestCases(t, tests)
}

// TestProtectedAttributeEnforcement tests that Protected symbols cannot be reassigned
func TestAttributeProtection(t *testing.T) {
	tests := []TestCase{
		{
			name:      "Protected symbol cannot be assigned with Set",
			input:     "x = 5; SetAttributes(x, Protected); x = 10",
			expected:  "",
			errorType: "Protected",
		},
		{
			name:      "Protected symbol cannot be assigned with SetDelayed",
			input:     "y = 3; SetAttributes(y, Protected); y := 42",
			expected:  "",
			errorType: "Protected",
		},
		{
			name:      "Built-in Protected symbol cannot be reassigned",
			input:     "Plus = 42",
			expected:  "",
			errorType: "Protected",
		},
		{
			name:     "Symbol can be reassigned after clearing Protected",
			input:    "w = 1; SetAttributes(w, Protected); ClearAttributes(w, Protected); w = 2",
			expected: "2",
		},
	}
	runTestCases(t, tests)
}
