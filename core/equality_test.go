package core

import (
	"github.com/client9/cardinal/core/symbol"
	"testing"
)

func TestAtomEqual(t *testing.T) {
	tests := []struct {
		name     string
		atom1    Expr
		atom2    Expr
		expected bool
	}{
		// Same atoms
		{
			name:     "same integers",
			atom1:    NewInteger(42),
			atom2:    NewInteger(42),
			expected: true,
		},
		{
			name:     "same floats",
			atom1:    NewReal(3.14),
			atom2:    NewReal(3.14),
			expected: true,
		},
		{
			name:     "same strings",
			atom1:    NewString("hello"),
			atom2:    NewString("hello"),
			expected: true,
		},
		{
			name:     "same symbols",
			atom1:    NewSymbol("x"),
			atom2:    NewSymbol("x"),
			expected: true,
		},

		// Different atoms
		{
			name:     "different integers",
			atom1:    NewInteger(42),
			atom2:    NewInteger(43),
			expected: false,
		},
		{
			name:     "different floats",
			atom1:    NewReal(3.14),
			atom2:    NewReal(2.71),
			expected: false,
		},
		{
			name:     "different strings",
			atom1:    NewString("hello"),
			atom2:    NewString("world"),
			expected: false,
		},
		{
			name:     "different symbols",
			atom1:    NewSymbol("x"),
			atom2:    NewSymbol("y"),
			expected: false,
		},

		// Different types
		{
			name:     "int vs float",
			atom1:    NewInteger(42),
			atom2:    NewReal(42.0),
			expected: false,
		},
		{
			name:     "int vs string",
			atom1:    NewInteger(42),
			atom2:    NewString("42"),
			expected: false,
		},
		{
			name:     "symbol vs string",
			atom1:    NewSymbol("hello"),
			atom2:    NewString("hello"),
			expected: false,
		},

		// Atom vs non-atom
		{
			name:     "atom vs list",
			atom1:    NewInteger(42),
			atom2:    NewList(symbol.List, NewInteger(42)),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.atom1.Equal(tt.atom2)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestListEqual(t *testing.T) {
	tests := []struct {
		name     string
		list1    List
		list2    Expr
		expected bool
	}{
		// Empty lists
		{
			name:     "empty lists",
			list1:    NewList(symbol.List),
			list2:    NewList(symbol.List),
			expected: true,
		},

		// Empty lists
		{
			name:     "empty lists different heads",
			list1:    NewList(symbol.List),
			list2:    NewList(symbol.Plus),
			expected: false,
		},

		// Same multi-element lists
		{
			name:     "same multi-element",
			list1:    NewList(symbol.Plus, NewInteger(1), NewInteger(2)),
			list2:    NewList(symbol.Plus, NewInteger(1), NewInteger(2)),
			expected: true,
		},

		// Different lengths
		{
			name:     "different lengths",
			list1:    NewList(symbol.Plus, NewInteger(1)),
			list2:    NewList(symbol.Plus, NewInteger(1), NewInteger(2)),
			expected: false,
		},

		// Different elements
		{
			name:     "different elements",
			list1:    NewList(symbol.List, NewInteger(1), NewInteger(2)),
			list2:    NewList(symbol.List, NewInteger(1), NewInteger(3)),
			expected: false,
		},

		// Nested lists - same
		{
			name: "nested lists same",
			list1: NewList(NewSymbol("f"),
				NewList(symbol.Plus, NewInteger(1), NewInteger(2)),
				NewInteger(3),
			),
			list2: NewList(NewSymbol("f"),
				NewList(symbol.Plus, NewInteger(1), NewInteger(2)),
				NewInteger(3),
			),
			expected: true,
		},

		// Nested lists - different
		{
			name: "nested lists different",
			list1: NewList(NewSymbol("f"),
				NewList(symbol.Plus, NewInteger(1), NewInteger(2)),
				NewInteger(3),
			),
			list2: NewList(NewSymbol("f"),
				NewList(symbol.Times, NewInteger(1), NewInteger(2)),
				NewInteger(3),
			),
			expected: false,
		},

		// List vs non-list
		{
			name:     "list vs atom",
			list1:    NewList(symbol.List, NewInteger(42)),
			list2:    NewInteger(42),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.list1.Equal(tt.list2)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestErrorEqual(t *testing.T) {
	t.Skip()
	tests := []struct {
		name     string
		error1   ErrorExpr
		error2   Expr
		expected bool
	}{
		// Same errors
		{
			name:     "same basic error",
			error1:   NewError("DivisionByZero", "Division by zero"),
			error2:   NewError("DivisionByZero", "Division by zero"),
			expected: true,
		},
		{
			name:     "same error with args",
			error1:   NewError("ArgumentError", "Wrong number of args"),
			error2:   NewError("ArgumentError", "Wrong number of args"),
			expected: true,
		},

		// Different error types
		{
			name:     "different error types",
			error1:   NewError("DivisionByZero", "Division by zero"),
			error2:   NewError("ArgumentError", "Division by zero"),
			expected: false,
		},

		// Different messages
		{
			name:     "different messages",
			error1:   NewError("DivisionByZero", "Division by zero"),
			error2:   NewError("DivisionByZero", "Cannot divide by zero"),
			expected: false,
		},

		// Different arguments
		{
			name:     "different arguments",
			error1:   NewError("ArgumentError", "Wrong args"),
			error2:   NewError("ArgumentError", "Wrong args"),
			expected: false,
		},

		// Error vs non-error
		{
			name:     "error vs atom",
			error1:   NewError("DivisionByZero", "Division by zero"),
			error2:   NewInteger(42),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.error1.Equal(tt.error2)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestListEqualMethod(t *testing.T) {
	// Test the List.Equal method directly
	tests := []struct {
		name     string
		list1    List
		list2    List
		expected bool
	}{
		{
			name:     "same lists",
			list1:    NewList(symbol.Plus, NewInteger(1), NewInteger(2)),
			list2:    NewList(symbol.Plus, NewInteger(1), NewInteger(2)),
			expected: true,
		},
		{
			name:     "different lists",
			list1:    NewList(symbol.Plus, NewInteger(1), NewInteger(2)),
			list2:    NewList(symbol.Times, NewInteger(1), NewInteger(2)),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.list1.Equal(tt.list2)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
