package core

// Sliceable interface for types that support indexed access and slicing operations
// Uses 1-based indexing to match Mathematica conventions
type Sliceable interface {
	// ElementAt returns the nth element (1-indexed)
	// Returns an error Expr if index is out of bounds
	ElementAt(n int64) Expr

	// Slice returns a new Expr containing elements from start to stop (inclusive, 1-indexed)
	// Returns an error Expr if indices are out of bounds
	Slice(start, stop int64) Expr

	// Join joins this sliceable with another sliceable of the same type
	// Returns an error Expr if the types are incompatible
	Join(other Sliceable) Expr

	// SetElementAt returns a new Expr with the nth element replaced (1-indexed)
	// Returns an error Expr if index is out of bounds or value is incompatible
	SetElementAt(n int64, value Expr) Expr

	// SetSlice returns a new Expr with elements from start to stop replaced by values (1-indexed)
	// values can be a single Expr, List, or other Sliceable
	// Returns an error Expr if indices are out of bounds or values are incompatible
	SetSlice(start, stop int64, values Expr) Expr
}

// IsSliceable checks if an expression implements the Sliceable interface and actually supports slicing
func IsSliceable(expr Expr) bool {
	switch e := expr.(type) {
	case List:
		return true
	case Atom:
		return e.AtomType == StringAtom
	case ByteArray:
		return true
	default:
		return false
	}
}

// AsSliceable safely casts an Expr to Sliceable, returning nil if not sliceable
func AsSliceable(expr Expr) Sliceable {
	if IsSliceable(expr) {
		if sliceable, ok := expr.(Sliceable); ok {
			return sliceable
		}
	}
	return nil
}
