package sexpr

import (
	"github.com/client9/sexpr/core"
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
			atom1:    core.NewInteger(42),
			atom2:    core.NewInteger(42),
			expected: true,
		},
		{
			name:     "same floats",
			atom1:    core.NewReal(3.14),
			atom2:    core.NewReal(3.14),
			expected: true,
		},
		{
			name:     "same strings",
			atom1:    core.NewString("hello"),
			atom2:    core.NewString("hello"),
			expected: true,
		},
		{
			name:     "same symbols",
			atom1:    core.NewSymbol("x"),
			atom2:    core.NewSymbol("x"),
			expected: true,
		},

		// Different atoms
		{
			name:     "different integers",
			atom1:    core.NewInteger(42),
			atom2:    core.NewInteger(43),
			expected: false,
		},
		{
			name:     "different floats",
			atom1:    core.NewReal(3.14),
			atom2:    core.NewReal(2.71),
			expected: false,
		},
		{
			name:     "different strings",
			atom1:    core.NewString("hello"),
			atom2:    core.NewString("world"),
			expected: false,
		},
		{
			name:     "different symbols",
			atom1:    core.NewSymbol("x"),
			atom2:    core.NewSymbol("y"),
			expected: false,
		},

		// Different types
		{
			name:     "int vs float",
			atom1:    core.NewInteger(42),
			atom2:    core.NewReal(42.0),
			expected: false,
		},
		{
			name:     "int vs string",
			atom1:    core.NewInteger(42),
			atom2:    core.NewString("42"),
			expected: false,
		},
		{
			name:     "symbol vs string",
			atom1:    core.NewSymbol("hello"),
			atom2:    core.NewString("hello"),
			expected: false,
		},

		// Atom vs non-atom
		{
			name:     "atom vs list",
			atom1:    core.NewInteger(42),
			atom2:    NewList(core.NewSymbol("List"), core.NewInteger(42)),
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
			list1:    NewList(),
			list2:    NewList(),
			expected: true,
		},

		// Same single element lists
		{
			name:     "same single element",
			list1:    NewList(core.NewInteger(42)),
			list2:    NewList(core.NewInteger(42)),
			expected: true,
		},

		// Same multi-element lists
		{
			name:     "same multi-element",
			list1:    NewList(core.NewSymbol("Plus"), core.NewInteger(1), core.NewInteger(2)),
			list2:    NewList(core.NewSymbol("Plus"), core.NewInteger(1), core.NewInteger(2)),
			expected: true,
		},

		// Different lengths
		{
			name:     "different lengths",
			list1:    NewList(core.NewInteger(1)),
			list2:    NewList(core.NewInteger(1), core.NewInteger(2)),
			expected: false,
		},

		// Different elements
		{
			name:     "different elements",
			list1:    NewList(core.NewInteger(1), core.NewInteger(2)),
			list2:    NewList(core.NewInteger(1), core.NewInteger(3)),
			expected: false,
		},

		// Nested lists - same
		{
			name: "nested lists same",
			list1: NewList(
				core.NewSymbol("f"),
				NewList(core.NewSymbol("Plus"), core.NewInteger(1), core.NewInteger(2)),
				core.NewInteger(3),
			),
			list2: NewList(
				core.NewSymbol("f"),
				NewList(core.NewSymbol("Plus"), core.NewInteger(1), core.NewInteger(2)),
				core.NewInteger(3),
			),
			expected: true,
		},

		// Nested lists - different
		{
			name: "nested lists different",
			list1: NewList(
				core.NewSymbol("f"),
				NewList(core.NewSymbol("Plus"), core.NewInteger(1), core.NewInteger(2)),
				core.NewInteger(3),
			),
			list2: NewList(
				core.NewSymbol("f"),
				NewList(core.NewSymbol("Times"), core.NewInteger(1), core.NewInteger(2)),
				core.NewInteger(3),
			),
			expected: false,
		},

		// List vs non-list
		{
			name:     "list vs atom",
			list1:    NewList(core.NewInteger(42)),
			list2:    core.NewInteger(42),
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
	tests := []struct {
		name     string
		error1   *ErrorExpr
		error2   Expr
		expected bool
	}{
		// Same errors
		{
			name:     "same basic error",
			error1:   NewErrorExpr("DivisionByZero", "Division by zero", []Expr{}),
			error2:   NewErrorExpr("DivisionByZero", "Division by zero", []Expr{}),
			expected: true,
		},
		{
			name: "same error with args",
			error1: NewErrorExpr("ArgumentError", "Wrong number of args", []Expr{
				core.NewInteger(1), core.NewInteger(2),
			}),
			error2: NewErrorExpr("ArgumentError", "Wrong number of args", []Expr{
				core.NewInteger(1), core.NewInteger(2),
			}),
			expected: true,
		},

		// Different error types
		{
			name:     "different error types",
			error1:   NewErrorExpr("DivisionByZero", "Division by zero", []Expr{}),
			error2:   NewErrorExpr("ArgumentError", "Division by zero", []Expr{}),
			expected: false,
		},

		// Different messages
		{
			name:     "different messages",
			error1:   NewErrorExpr("DivisionByZero", "Division by zero", []Expr{}),
			error2:   NewErrorExpr("DivisionByZero", "Cannot divide by zero", []Expr{}),
			expected: false,
		},

		// Different arguments
		{
			name: "different arguments",
			error1: NewErrorExpr("ArgumentError", "Wrong args", []Expr{
				core.NewInteger(1),
			}),
			error2: NewErrorExpr("ArgumentError", "Wrong args", []Expr{
				core.NewInteger(2),
			}),
			expected: false,
		},

		// Error vs non-error
		{
			name:     "error vs atom",
			error1:   NewErrorExpr("DivisionByZero", "Division by zero", []Expr{}),
			error2:   core.NewInteger(42),
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

func TestStructuralEqualityIntegration(t *testing.T) {
	// Test that the updated SameQ function works correctly
	evaluator := NewEvaluator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "SameQ identical atoms",
			input:    "SameQ(42, 42)",
			expected: "True",
		},
		{
			name:     "SameQ different atoms",
			input:    "SameQ(42, 43)",
			expected: "False",
		},
		{
			name:     "SameQ identical lists",
			input:    "SameQ(Plus(1, 2), Plus(1, 2))",
			expected: "True",
		},
		{
			name:     "SameQ different lists",
			input:    "SameQ(Plus(1, 2), Plus(1, 3))",
			expected: "False",
		},
		{
			name:     "SameQ different types",
			input:    "SameQ(42, 42.0)",
			expected: "False",
		},
		{
			name:     "SameQ nested expressions",
			input:    "SameQ(f(Plus(a, b), c), f(Plus(a, b), c))",
			expected: "True",
		},
		{
			name:     "SameQ with symbols",
			input:    "SameQ(x, x)",
			expected: "True",
		},
		{
			name:     "SameQ different symbols",
			input:    "SameQ(x, y)",
			expected: "False",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluateStringSimple(t, evaluator, tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestListsEqualFunction(t *testing.T) {
	// Test the updated listsEqual function
	tests := []struct {
		name     string
		list1    List
		list2    List
		expected bool
	}{
		{
			name:     "same lists",
			list1:    NewList(core.NewSymbol("Plus"), core.NewInteger(1), core.NewInteger(2)),
			list2:    NewList(core.NewSymbol("Plus"), core.NewInteger(1), core.NewInteger(2)),
			expected: true,
		},
		{
			name:     "different lists",
			list1:    NewList(core.NewSymbol("Plus"), core.NewInteger(1), core.NewInteger(2)),
			list2:    NewList(core.NewSymbol("Times"), core.NewInteger(1), core.NewInteger(2)),
			expected: false,
		},
		{
			name:     "empty lists",
			list1:    NewList(),
			list2:    NewList(),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := listsEqual(tt.list1, tt.list2)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
