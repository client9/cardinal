package core

// needsSequenceSplicing determines if a substitution should be spliced (for sequence patterns)
func needsSequenceSplicing(originalElem, newElem Expr, bindings PatternBindings) bool {
	// Check if original element is a symbol that was bound to a List
	if elemSym, ok := originalElem.(Symbol); ok {
		varName := string(elemSym)
		if boundValue, exists := bindings[varName]; exists {
			// Check if the bound value is a List (indicating sequence pattern)
			return boundValue.Head() == "List"
		}
	}
	return false
}

// SubstituteBindings replaces pattern variables in an expression with their bound values
func SubstituteBindings(expr Expr, bindings PatternBindings) Expr {
	switch e := expr.(type) {
	case Symbol:
		// Check if this symbol is a bound variable
		if value, exists := bindings[string(e)]; exists {
			return value
		}
		return e

	case List:
		// Recursively substitute in list elements with sequence splicing
		newElements := make([]Expr, 0, e.Length()+1)
		changed := false

		for i, elem := range e.AsSlice() {
			newElem := SubstituteBindings(elem, bindings)

			// Check if this is a sequence variable substitution that needs splicing
			if i > 0 && needsSequenceSplicing(elem, newElem, bindings) {
				// This is a sequence variable - splice its elements
				if elemSym, ok := elem.(Symbol); ok {
					varName := string(elemSym)
					if boundValue, exists := bindings[varName]; exists {
						if boundList, ok := boundValue.(List); ok {
							// Check if it's an empty sequence (just "List" head)
							if boundList.Length() == 0 {
								// Empty sequence - skip adding anything
								changed = true
								continue
							} else {
								// Non-empty sequence - skip the "List" head and add the actual elements
								newElements = append(newElements, boundList.Tail()...)
								changed = true
								continue
							}
						}
					}
				}
			}

			// Regular substitution (not a sequence)
			newElements = append(newElements, newElem)
			if !newElem.Equal(elem) {
				changed = true
			}
		}

		if changed {
			return NewListFromExprs(newElements...)
		}
		return e

	default:
		// For atomic types, no substitution needed
		return e
	}
}
