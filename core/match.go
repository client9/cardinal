package core

// Pure pattern matching (no variable binding)

// PatternMatcher provides pure pattern matching without side effects
type PatternMatcher struct{}

// NewPatternMatcher creates a new PatternMatcher
func NewPatternMatcher() *PatternMatcher {
	return &PatternMatcher{}
}

// TestMatch tests if an expression matches a pattern (pure function, no binding)
func (pm *PatternMatcher) TestMatch(pattern, expr Expr) bool {
	return pm.testMatchInternal(pattern, expr)
}

// testMatchInternal implements the core matching logic
func (pm *PatternMatcher) testMatchInternal(pattern, expr Expr) bool {
	// Handle symbolic patterns
	if isPattern, _, blankExpr := IsSymbolicPattern(pattern); isPattern {
		// For pure matching, just test the blank part (no variable binding)
		return pm.testMatchBlank(blankExpr, expr)
	}

	// Handle direct symbolic blanks
	if isBlank, _, _ := IsSymbolicBlank(pattern); isBlank {
		return pm.testMatchBlank(pattern, expr)
	}

	// pattern and expr are both lists
	if patList, ok := pattern.(List); ok {
		if exprList, ok := expr.(List); ok {
			return pm.testMatchList(patList, exprList)
		}
		return false
	}

	return pattern.Equal(expr)
}

// testMatchBlank tests if an expression matches a blank pattern
func (pm *PatternMatcher) testMatchBlank(blankExpr, expr Expr) bool {
	isBlank, blankType, typeExpr := IsSymbolicBlank(blankExpr)
	if !isBlank {
		return false
	}

	// Check type constraint if present
	if typeExpr != nil {
		var typeName string
		if typeAtom, ok := typeExpr.(Symbol); ok {
			typeName = typeAtom.String()
		}
		if !MatchesType(expr, typeName) {
			return false
		}
	}

	// For pure matching, single expressions match all blank types
	// (sequence handling is more complex and typically needs context

	// TODO
	switch blankType {
	case BlankPattern, BlankSequencePattern, BlankNullSequencePattern:
		return true
	}

	return false
}

// testMatchList tests if two lists match
func (pm *PatternMatcher) testMatchList(patternList, exprList List) bool {
	// Use the full pattern matching logic but discard bindings
	// This properly handles sequence patterns like z___
	bindings := make(PatternBindings)
	return matchListWithBindings(patternList, exprList, bindings)
}

// PatternBindings represents variable bindings from pattern matching
type PatternBindings map[string]Expr

// MatchWithBindings performs pattern matching and captures variable bindings
// Returns (matches, bindings)
func MatchWithBindings(pattern, expr Expr) (bool, PatternBindings) {
	bindings := make(PatternBindings)
	matches := matchWithBindingsInternal(pattern, expr, bindings)
	return matches, bindings
}

// matchWithBindingsInternal implements pattern matching with binding capture
func matchWithBindingsInternal(pattern, expr Expr, bindings PatternBindings) bool {

	// Handle symbolic patterns with variable binding
	if isPattern, varName, blankExpr := IsSymbolicPattern(pattern); isPattern {
		// Test if the blank part matches
		if matchBlankWithBindings(blankExpr, expr, bindings) {
			vn := varName.String()
			// If there's a variable name, check for existing binding or create new one
			if vn != "" {
				if existingValue, exists := bindings[vn]; exists {
					// Variable already bound - check if values match
					return existingValue.Equal(expr)
				} else {
					// New binding
					bindings[vn] = expr
				}
			}
			return true
		}
		return false
	}

	// Handle direct symbolic blanks
	if isBlank, _, _ := IsSymbolicBlank(pattern); isBlank {
		return matchBlankWithBindings(pattern, expr, bindings)
	}

	// Handle different expression types
	switch p := pattern.(type) {
	case List:
		if exprList, ok := expr.(List); ok {
			return matchListWithBindings(p, exprList, bindings)
		}
		return false
	default:
		// For other types, just check equality
		return pattern.Equal(expr)
	}
}

// matchBlankWithBindings tests if a blank pattern matches an expression
func matchBlankWithBindings(blankExpr, expr Expr, bindings PatternBindings) bool {
	isBlank, _, typeExpr := IsSymbolicBlank(blankExpr)
	if !isBlank {
		return false
	}

	// If no type constraint, accept any expression
	if typeExpr == nil {
		return true
	}

	// Extract type name from type expression
	typeName, _ := ExtractSymbol(typeExpr)

	// Check type constraint
	return MatchesType(expr, typeName)
}

// matchListWithBindings tests if a list pattern matches a list expression
func matchListWithBindings(patternList, exprList List, bindings PatternBindings) bool {
	return matchListWithBindingsSequential(patternList, exprList, bindings, 0, 0)
}

// matchListWithBindingsSequential handles pattern matching with sequence patterns
func matchListWithBindingsSequential(patternList, exprList List, bindings PatternBindings, patternIdx, exprIdx int) bool {
	patternSlice := patternList.AsSlice()
	exprSlice := exprList.AsSlice()

	// If we've processed all pattern elements
	if patternIdx >= len(patternSlice) {
		// Success if we've also processed all expression elements
		return exprIdx >= len(exprSlice)
	}

	// If we've run out of expression elements but still have patterns
	if exprIdx >= len(exprSlice) {
		// Check if remaining patterns are all BlankNullSequence (which can match zero elements)
		for i := patternIdx; i < len(patternSlice); i++ {
			elem := patternSlice[i]
			if !isNullSequencePattern(elem) {
				return false
			}
			// Bind null sequence patterns to empty list
			bindNullSequencePattern(elem, bindings)
		}
		return true
	}

	patternElem := patternSlice[patternIdx]

	// Check if this is a sequence pattern
	if isSequencePattern, varName, typeName, allowZero := analyzeSequencePattern(patternElem); isSequencePattern {
		return matchSequencePatternWithBindings(patternList, exprList, bindings, patternIdx, exprIdx, varName, typeName, allowZero)
	}

	// Regular pattern - match one element
	if matchWithBindingsInternal(patternElem, exprSlice[exprIdx], bindings) {
		return matchListWithBindingsSequential(patternList, exprList, bindings, patternIdx+1, exprIdx+1)
	}

	return false
}

// analyzeSequencePattern determines if a pattern element is a sequence pattern
func analyzeSequencePattern(pattern Expr) (isSequence bool, varName string, typeName string, allowZero bool) {
	// Check for symbolic Pattern[name, BlankSequence/BlankNullSequence]
	if isPattern, nameExpr, blankExpr := IsSymbolicPattern(pattern); isPattern {
		if isBlank, blankType, typeExpr := IsSymbolicBlank(blankExpr); isBlank {
			vn, _ := ExtractSymbol(nameExpr)
			tn := ""
			if typeExpr != nil {
				tn, _ = ExtractSymbol(typeExpr)
			}

			switch blankType {
			case BlankSequencePattern:
				return true, vn, tn, false
			case BlankNullSequencePattern:
				return true, vn, tn, true
			}
		}
	}

	// Check for direct symbolic BlankSequence/BlankNullSequence
	if isBlank, blankType, typeExpr := IsSymbolicBlank(pattern); isBlank {
		tn := ""
		if typeExpr != nil {
			tn, _ = ExtractSymbol(typeExpr)
		}

		switch blankType {
		case BlankSequencePattern:
			return true, "", tn, false
		case BlankNullSequencePattern:
			return true, "", tn, true
		}
	}
	return false, "", "", false
}

// isNullSequencePattern checks if a pattern is a BlankNullSequence that can match zero elements
func isNullSequencePattern(pattern Expr) bool {
	// Check for symbolic Pattern[name, BlankNullSequence]
	if isPattern, _, blankExpr := IsSymbolicPattern(pattern); isPattern {
		if isBlank, blankType, _ := IsSymbolicBlank(blankExpr); isBlank && blankType == BlankNullSequencePattern {
			return true
		}
	}

	// Check for direct symbolic BlankNullSequence
	if isBlank, blankType, _ := IsSymbolicBlank(pattern); isBlank && blankType == BlankNullSequencePattern {
		return true
	}

	return false
}

// bindNullSequencePattern binds a null sequence pattern to an empty list
func bindNullSequencePattern(pattern Expr, bindings PatternBindings) {
	// Extract variable name if present
	varName := ""

	// Check for symbolic Pattern[name, BlankNullSequence]
	if isPattern, nameExpr, _ := IsSymbolicPattern(pattern); isPattern {
		varName, _ = ExtractSymbol(nameExpr)
	}

	// Bind to empty list if we have a variable name
	if varName != "" {
		// Create a proper List with "List" as the first element (like the old system)
		bindings[varName] = NewList("List")
	}
}

// matchSequencePatternWithBindings handles matching sequence patterns
func matchSequencePatternWithBindings(patternList, exprList List, bindings PatternBindings, patternIdx, exprIdx int, varName, typeName string, allowZero bool) bool {

	patternSlice := patternList.AsSlice()
	exprSlice := exprList.AsSlice()

	remainingPatterns := len(patternSlice) - patternIdx - 1
	remainingExprs := len(exprSlice) - exprIdx

	// Minimum elements this sequence must consume
	minConsume := 0
	if !allowZero {
		minConsume = 1
	}

	// Maximum elements this sequence can consume
	maxConsume := remainingExprs - remainingPatterns

	if maxConsume < minConsume {
		return false
	}

	// Try consuming different numbers of elements (greedy approach)
	for consume := maxConsume; consume >= minConsume; consume-- {
		// Create a copy of bindings for this attempt
		testBindings := make(PatternBindings)
		for k, v := range bindings {
			testBindings[k] = v
		}

		// Collect the elements to bind to this sequence
		var seqElements []Expr
		for i := 0; i < consume; i++ {
			if exprIdx+i < len(exprSlice) {
				seqElements = append(seqElements, exprSlice[exprIdx+i])
			}
		}

		// Check type constraints for sequence elements
		if typeName != "" {
			allMatch := true
			for _, elem := range seqElements {
				if !MatchesType(elem, typeName) {
					allMatch = false
					break
				}
			}
			if !allMatch {
				continue // Try with fewer elements
			}
		}

		// Bind the sequence variable if we have a name
		if varName != "" {
			// Create a proper List with "List" as the first element (consistent with old system)
			listElements := append([]Expr{NewSymbol("List")}, seqElements...)
			testBindings[varName] = NewListFromExprs(listElements...)
		}

		// Try to match remaining patterns
		if matchListWithBindingsSequential(patternList, exprList, testBindings, patternIdx+1, exprIdx+consume) {
			// Success - copy bindings back
			for k, v := range testBindings {
				bindings[k] = v
			}
			return true
		}
	}

	return false
}

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
