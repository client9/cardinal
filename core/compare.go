package core

// Should be called "Less"

// CanonicalCompare provides a canonical comparison function for expressions
// Used for consistent ordering across mathematical functions and Orderless attribute
// Returns true if expr1 should come before expr2 in canonical ordering
func CanonicalCompare(expr1, expr2 Expr) bool {
	// Mathematical ordering: numbers first, then other expressions
	_, expr1IsNumber := GetNumericValue(expr1)
	_, expr2IsNumber := GetNumericValue(expr2)

	if expr1IsNumber && !expr2IsNumber {
		return true // Numbers come before non-numbers
	}
	if !expr1IsNumber && expr2IsNumber {
		return false // Non-numbers come after numbers
	}

	// Both are numbers or both are non-numbers, use standard ordering
	// First compare by length (complexity)
	cmp := expr1.Length() - expr2.Length()
	if cmp < 0 {
		return true
	}
	if cmp > 0 {
		return false
	}

	// If lengths are equal, compare by string representation for deterministic ordering
	return expr1.String() < expr2.String()
}

