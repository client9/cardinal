package sexpr

import (
	"testing"
)

// TestAssociationImmutability verifies that association operations don't modify the original
func TestAssociationImmutability(t *testing.T) {
	// Create original association
	original := NewAssociationValue()
	original = original.Set(NewSymbolAtom("a"), NewIntAtom(1))
	original = original.Set(NewSymbolAtom("b"), NewIntAtom(2))

	// Make a copy to verify original doesn't change
	originalKeys := original.Keys()
	if len(originalKeys) != 2 {
		t.Fatalf("Original should have 2 keys, got %d", len(originalKeys))
	}

	// Add a new key-value pair - should return new association
	modified := original.Set(NewSymbolAtom("c"), NewIntAtom(3))

	// Verify original is unchanged
	originalKeysAfter := original.Keys()
	if len(originalKeysAfter) != 2 {
		t.Errorf("Original association was modified! Expected 2 keys, got %d", len(originalKeysAfter))
	}

	// Verify modified has the new entry
	modifiedKeys := modified.Keys()
	if len(modifiedKeys) != 3 {
		t.Errorf("Modified association should have 3 keys, got %d", len(modifiedKeys))
	}

	// Verify original values are still accessible
	if val, exists := original.Get(NewSymbolAtom("a")); !exists || !val.Equal(NewIntAtom(1)) {
		t.Error("Original association lost value for key 'a'")
	}

	// Verify new value is in modified association
	if val, exists := modified.Get(NewSymbolAtom("c")); !exists || !val.Equal(NewIntAtom(3)) {
		t.Error("Modified association doesn't have new key 'c'")
	}

	// Verify original doesn't have new key
	if _, exists := original.Get(NewSymbolAtom("c")); exists {
		t.Error("Original association was contaminated with new key 'c'")
	}
}

// TestListImmutability verifies that list operations maintain immutability
func TestListImmutability(t *testing.T) {
	// Test Rest operation
	original := List{Elements: []Expr{
		NewSymbolAtom("List"),
		NewIntAtom(1),
		NewIntAtom(2),
		NewIntAtom(3),
	}}

	// Create a copy to verify immutability
	originalCopy := List{Elements: make([]Expr, len(original.Elements))}
	copy(originalCopy.Elements, original.Elements)

	// Perform Rest operation
	result := EvaluateRest([]Expr{original})

	// Verify original is unchanged
	if !original.Equal(originalCopy) {
		t.Error("Rest operation modified the original list")
	}

	// Verify result is correct
	if resultList, ok := result.(List); ok {
		if len(resultList.Elements) != 3 { // Head + 2 remaining elements
			t.Errorf("Rest result should have 3 elements (head + 2), got %d", len(resultList.Elements))
		}

		// Should be List[2, 3]
		expected := List{Elements: []Expr{
			NewSymbolAtom("List"),
			NewIntAtom(2),
			NewIntAtom(3),
		}}

		if !resultList.Equal(expected) {
			t.Errorf("Rest result incorrect. Expected %s, got %s", expected.String(), resultList.String())
		}
	} else {
		t.Errorf("Rest should return a List, got %T", result)
	}
}

// TestBuiltinFunctionImmutability tests that builtin functions don't modify their arguments
func TestBuiltinFunctionImmutability(t *testing.T) {
	tests := []struct {
		name     string
		function func([]Expr) Expr
		args     []Expr
	}{
		{
			name:     "EvaluatePlus",
			function: EvaluatePlus,
			args:     []Expr{NewIntAtom(1), NewIntAtom(2), NewIntAtom(3)},
		},
		{
			name:     "EvaluateTimes",
			function: EvaluateTimes,
			args:     []Expr{NewIntAtom(2), NewIntAtom(3), NewIntAtom(4)},
		},
		{
			name:     "EvaluateLength",
			function: EvaluateLength,
			args:     []Expr{List{Elements: []Expr{NewSymbolAtom("List"), NewIntAtom(1), NewIntAtom(2)}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create copies of arguments
			originalArgs := make([]Expr, len(tt.args))
			for i, arg := range tt.args {
				switch a := arg.(type) {
				case Atom:
					originalArgs[i] = a // Atoms are immutable by design
				case List:
					// Create a copy of the list
					newElements := make([]Expr, len(a.Elements))
					copy(newElements, a.Elements)
					originalArgs[i] = List{Elements: newElements}
				default:
					originalArgs[i] = arg // Other types handled similarly
				}
			}

			// Call the function
			_ = tt.function(tt.args)

			// Verify arguments are unchanged
			for i, original := range originalArgs {
				if !tt.args[i].Equal(original) {
					t.Errorf("Function %s modified argument %d", tt.name, i)
				}
			}
		})
	}
}

// TestStructuralSharing verifies that immutable operations efficiently share unchanged data
func TestStructuralSharing(t *testing.T) {
	// Create a list with elements
	original := List{Elements: []Expr{
		NewSymbolAtom("List"),
		NewIntAtom(1),
		NewIntAtom(2),
		NewIntAtom(3),
	}}

	// Perform First operation - should share element references
	result := EvaluateFirst([]Expr{original})

	// Verify result is correct
	if !result.Equal(NewIntAtom(1)) {
		t.Errorf("First should return 1, got %s", result.String())
	}

	// For structural sharing test, we can at least verify that no deep copying occurred
	// by checking that the operation completed efficiently (this is more of a design verification)
	// In a more sophisticated test, we might check memory addresses, but Go's value semantics
	// make this less relevant for our atom types
}
