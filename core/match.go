package core

// Pure pattern matching (no variable binding)

func IsExcept(pattern Expr) Expr {
	if pattern.Head() != "Except" {
		return nil
	}
	plist, _ := pattern.(List)
	return plist.Tail()[0]
}

func IsAlternatives(pattern Expr) []Expr {
	if pattern.Head() != "Alternatives" {
		return nil
	}
	plist, _ := pattern.(List)
	return plist.Tail()
}

// PatternMatcher provides pure pattern matching without side effects
type PatternMatcher struct{}

// NewPatternMatcher creates a new PatternMatcher
func NewPatternMatcher() *PatternMatcher {
	return &PatternMatcher{}
}

// TestMatch tests if an expression matches a pattern (pure function, no binding)
func (pm *PatternMatcher) TestMatch(pattern, expr Expr) bool {
	ok, _ := MatchWithBindings(pattern, expr)
	return ok
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

	if plist := IsAlternatives(pattern); plist != nil {
		for _, p := range plist {
			if matchWithBindingsInternal(p, expr, bindings) {
				return true
			}
		}
		return false
	}

	if pinfo := GetSymbolicPatternInfo(pattern); pinfo.Type != PatternUnknown {
		if !matchBlankWithBindings(pinfo, expr, bindings) {
			return false
		}
		if vn := pinfo.VarName; vn != "" {
			if existingValue, exists := bindings[vn]; exists {
				// Variable already bound - check if values match
				return existingValue.Equal(expr)
			}
			bindings[vn] = expr
		}
		return true
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
func matchBlankWithBindings(pinfo PatternInfo, expr Expr, bindings PatternBindings) bool {
	if pinfo.Type == PatternUnknown {
		return false
	}

	// Check type constraint
	return MatchesType(expr, pinfo.TypeName)
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

			pinfo := GetSymbolicPatternInfo(elem)
			if pinfo.Type != BlankNullSequencePattern {
				return false
			}
			if vn := pinfo.VarName; vn != "" {
				// Bind null sequence patterns to empty list
				// Create a proper List with "List" as the first element (like the old system)
				bindings[vn] = NewList("List")
			}
		}
		return true
	}

	patternElem := patternSlice[patternIdx]
	pinfo := GetSymbolicPatternInfo(patternElem)
	if pinfo.Type == BlankNullSequencePattern || pinfo.Type == BlankSequencePattern {
		// Check if this is a sequence pattern
		return matchSequencePatternWithBindings(patternList, exprList, bindings, patternIdx, exprIdx, pinfo)
	}

	// Regular pattern - match one element
	if matchWithBindingsInternal(patternElem, exprSlice[exprIdx], bindings) {
		return matchListWithBindingsSequential(patternList, exprList, bindings, patternIdx+1, exprIdx+1)
	}

	return false
}

// matchSequencePatternWithBindings handles matching sequence patterns
func matchSequencePatternWithBindings(patternList, exprList List, bindings PatternBindings, patternIdx, exprIdx int, pinfo PatternInfo) bool { //
	// , varName, typeName string, allowZero bool) bool {

	typeName := pinfo.TypeName
	varName := pinfo.VarName
	allowZero := pinfo.Type == BlankNullSequencePattern

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
