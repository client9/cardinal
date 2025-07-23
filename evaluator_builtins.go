package sexpr

import (
	"fmt"
)

// Pattern and Blank constructors

// CreateBlankExpr creates a symbolic Blank[] expression
func CreateBlankExpr(typeExpr Expr) Expr {
	if typeExpr == nil {
		return List{Elements: []Expr{NewSymbolAtom("Blank")}}
	}
	return List{Elements: []Expr{NewSymbolAtom("Blank"), typeExpr}}
}

// CreateBlankSequenceExpr creates a symbolic BlankSequence[] expression
func CreateBlankSequenceExpr(typeExpr Expr) Expr {
	if typeExpr == nil {
		return List{Elements: []Expr{NewSymbolAtom("BlankSequence")}}
	}
	return List{Elements: []Expr{NewSymbolAtom("BlankSequence"), typeExpr}}
}

// CreateBlankNullSequenceExpr creates a symbolic BlankNullSequence[] expression
func CreateBlankNullSequenceExpr(typeExpr Expr) Expr {
	if typeExpr == nil {
		return List{Elements: []Expr{NewSymbolAtom("BlankNullSequence")}}
	}
	return List{Elements: []Expr{NewSymbolAtom("BlankNullSequence"), typeExpr}}
}

// CreatePatternExpr creates a symbolic Pattern[name, blank] expression
func CreatePatternExpr(nameExpr, blankExpr Expr) Expr {
	return List{Elements: []Expr{NewSymbolAtom("Pattern"), nameExpr, blankExpr}}
}

// convertToSymbolicPattern converts a pattern to symbolic representation if it's a string-based pattern
func convertToSymbolicPattern(pattern Expr) Expr {
	switch p := pattern.(type) {
	case Atom:
		if p.AtomType == SymbolAtom {
			patternStr := p.Value.(string)
			// Check if it's a pattern variable
			if isPatternVariable(patternStr) {
				return ConvertPatternStringToSymbolic(patternStr)
			}
		}
		return p
	case List:
		// Convert all elements in the list
		newElements := make([]Expr, len(p.Elements))
		for i, elem := range p.Elements {
			newElements[i] = convertToSymbolicPattern(elem)
		}
		return List{Elements: newElements}
	default:
		return pattern
	}
}

// matchPatternForMatchQ implements pattern matching specifically for MatchQ
func matchPatternForMatchQ(pattern Expr, expr Expr, ctx *Context) bool {
	return matchPatternForMatchQWithContext(pattern, expr, ctx, false)
}

// matchPatternForMatchQWithContext handles pattern matching with context for MatchQ
func matchPatternForMatchQWithContext(pattern Expr, expr Expr, ctx *Context, isParameter bool) bool {
	// Handle symbolic patterns first
	if isPattern, _, blankExpr := isSymbolicPattern(pattern); isPattern {
		// Handle Pattern[name, blank] - just match the blank part for MatchQ
		return matchBlankExpression(blankExpr, expr, ctx)
	}

	// Handle direct symbolic blanks
	if isBlank, _, _ := isSymbolicBlank(pattern); isBlank {
		return matchBlankExpression(pattern, expr, ctx)
	}

	// Handle different expression types
	switch p := pattern.(type) {
	case Atom:
		if p.AtomType == SymbolAtom {
			varName := p.Value.(string)
			// Check if it's a pattern variable
			if isPatternVariable(varName) {
				info := parsePatternInfo(varName)
				if info.Type == BlankPattern {
					// Check type constraint if present
					if !matchesType(expr, info.TypeName) {
						return false
					}
					return true // Don't bind variables in MatchQ, just test matching
				}
			} else {
				// Regular symbol - must match literally
				if exprAtom, ok := expr.(Atom); ok && exprAtom.AtomType == SymbolAtom {
					return varName == exprAtom.Value.(string)
				}
				return false
			}
		}
		// For literal atoms, they must be exactly equal
		if exprAtom, ok := expr.(Atom); ok {
			return p.AtomType == exprAtom.AtomType && p.Value == exprAtom.Value
		}
		return false

	case List:
		if exprList, ok := expr.(List); ok {
			// Both are lists - need to match structure
			return matchListForMatchQ(p, exprList, ctx)
		}
		return false

	default:
		return false
	}
}

// matchListForMatchQ matches list patterns for MatchQ
func matchListForMatchQ(patternList List, exprList List, ctx *Context) bool {
	// Handle empty patterns
	if len(patternList.Elements) == 0 {
		return len(exprList.Elements) == 0
	}

	// Check if the length matches exactly (for now, no sequence patterns in MatchQ)
	if len(patternList.Elements) != len(exprList.Elements) {
		return false
	}

	// Match each element
	for i, patternElem := range patternList.Elements {
		// Element 0 is head (literal), elements 1+ are parameters (pattern match)
		isParameterPosition := i > 0
		if !matchPatternForMatchQWithContext(patternElem, exprList.Elements[i], ctx, isParameterPosition) {
			return false
		}
	}

	return true
}

// parseAttribute parses an attribute name string into an Attribute enum value
func parseAttribute(attrName string) (Attribute, error) {
	switch attrName {
	case "HoldAll":
		return HoldAll, nil
	case "HoldFirst":
		return HoldFirst, nil
	case "HoldRest":
		return HoldRest, nil
	case "Flat":
		return Flat, nil
	case "Orderless":
		return Orderless, nil
	case "OneIdentity":
		return OneIdentity, nil
	case "Listable":
		return Listable, nil
	case "Constant":
		return Constant, nil
	case "NumericFunction":
		return NumericFunction, nil
	case "Protected":
		return Protected, nil
	case "ReadProtected":
		return ReadProtected, nil
	case "Locked":
		return Locked, nil
	case "Temporary":
		return Temporary, nil
	default:
		return HoldAll, fmt.Errorf("unknown attribute: %s", attrName)
	}
}
