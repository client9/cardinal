package core

import (
	"testing"
)

func TestSymbolEqual(t *testing.T) {
	tests := []struct {
		name     string
		atom1    Symbol
		atom2    Symbol
		expected bool
	}{
		{
			name:     "Same, Equal",
			atom1:    NewSymbol("Foo"),
			atom2:    NewSymbol("Foo"),
			expected: true,
		},
		{
			name:     "Not Same, Not Equal",
			atom1:    NewSymbol("Foo"),
			atom2:    NewSymbol("Bar"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.atom1.Equal(tt.atom2)
			if result != tt.expected {
				t.Errorf("Symbol Equal expected %v, got %v", tt.expected, result)
			}
			result = tt.atom1 == tt.atom2
			if result != tt.expected {
				t.Errorf("Symbol ==  expected %v, got %v", tt.expected, result)
			}
		})
	}
}
func TestExprSymbolEqual(t *testing.T) {
	tests := []struct {
		name     string
		atom1    Expr
		atom2    Symbol
		expected bool
	}{
		{
			name:     "String/Symbol ",
			atom1:    NewString("Foo"),
			atom2:    NewSymbol("Foo"),
			expected: false,
		},
		{
			name:     "Not Same, Not Equal",
			atom1:    NewString("Foo"),
			atom2:    NewSymbol("Bar"),
			expected: false,
		},
		{
			name:     "Not Same, Not Equal",
			atom1:    NewInteger(123),
			atom2:    NewSymbol("Bar"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.atom1.Equal(tt.atom2)
			if result != tt.expected {
				t.Errorf("Symbol Equal expected %v, got %v", tt.expected, result)
			}
			result = tt.atom1 == tt.atom2
			if result != tt.expected {
				t.Errorf("Symbol ==  expected %v, got %v", tt.expected, result)
			}
		})
	}
}
