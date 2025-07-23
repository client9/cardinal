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

// EvaluateAttributes evaluates Attributes[symbol] expressions - returns the attributes of a symbol
func EvaluateAttributes(args []Expr, ctx *Context) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("Attributes expects 1 argument, got %d", len(args)), args)
	}

	// Check for errors in arguments first
	if IsError(args[0]) {
		return args[0]
	}

	// The argument should be a symbol
	if atom, ok := args[0].(Atom); ok && atom.AtomType == SymbolAtom {
		symbolName := atom.Value.(string)

		// Get the attributes from the symbol table
		attrs := ctx.symbolTable.Attributes(symbolName)

		// Convert attributes to a list of symbols
		attrElements := make([]Expr, len(attrs)+1)
		attrElements[0] = NewSymbolAtom("List")

		for i, attr := range attrs {
			attrElements[i+1] = NewSymbolAtom(attr.String())
		}

		return List{Elements: attrElements}
	}

	return NewErrorExpr("ArgumentError",
		"Attributes expects a symbol as argument", args)
}

// EvaluateMatchQ evaluates MatchQ[expr, pattern] expressions - tests if expr matches pattern
func EvaluateMatchQ(args []Expr, ctx *Context) Expr {
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("MatchQ expects 2 arguments, got %d", len(args)), args)
	}

	// Check for errors in arguments first
	for _, arg := range args {
		if IsError(arg) {
			return arg
		}
	}

	expr := args[0]
	pattern := args[1]

	// Convert string-based pattern to symbolic if needed
	symbolicPattern := convertToSymbolicPattern(pattern)

	// Create a temporary context for pattern matching (don't pollute original context)
	tempCtx := NewChildContext(ctx)

	// Use improved pattern matching logic
	matches := matchPatternForMatchQ(symbolicPattern, expr, tempCtx)
	return NewBoolAtom(matches)
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

// EvaluateFullForm evaluates FullForm[expr] expressions - returns the string representation of the expression
func EvaluateFullForm(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("FullForm expects 1 argument, got %d", len(args)), args)
	}

	expr := args[0]

	// Convert pattern strings to their symbolic form before getting string representation
	if atom, ok := expr.(Atom); ok && atom.AtomType == SymbolAtom {
		symbolStr := atom.Value.(string)
		if isPatternVariable(symbolStr) {
			// Convert pattern string to symbolic form
			symbolicExpr := ConvertPatternStringToSymbolic(symbolStr)
			return NewStringAtom(symbolicExpr.String())
		}
	}

	// Return the string representation as a string atom
	return NewStringAtom(expr.String())
}

// EvaluateInputForm evaluates InputForm[expr] expressions - returns the user-friendly InputForm representation of the expression
func EvaluateInputForm(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("InputForm expects 1 argument, got %d", len(args)), args)
	}

	expr := args[0]

	// Convert pattern strings to their symbolic form before getting InputForm representation
	if atom, ok := expr.(Atom); ok && atom.AtomType == SymbolAtom {
		symbolStr := atom.Value.(string)
		if isPatternVariable(symbolStr) {
			// Convert pattern string to symbolic form
			symbolicExpr := ConvertPatternStringToSymbolic(symbolStr)
			return NewStringAtom(symbolicExpr.InputForm())
		}
	}

	// Return the InputForm representation as a string atom
	return NewStringAtom(expr.InputForm())
}

// EvaluateFirst evaluates First[expr] expressions - returns the first element after the head, if any
func EvaluateFirst(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("First expects 1 argument, got %d", len(args)), args)
	}

	if list, ok := args[0].(List); ok {
		// For lists, return the first element after the head
		if len(list.Elements) <= 1 {
			return NewErrorExpr("PartError",
				fmt.Sprintf("First: expression %s has no elements", args[0].String()), args)
		}
		return list.Elements[1] // Index 1 is first element after head (index 0)
	}

	// For atoms, First is not defined
	return NewErrorExpr("PartError",
		fmt.Sprintf("First: expression %s is not a list", args[0].String()), args)
}

// EvaluateLast evaluates Last[expr] expressions - returns the last element that is not the head
func EvaluateLast(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("Last expects 1 argument, got %d", len(args)), args)
	}

	if list, ok := args[0].(List); ok {
		// For lists, return the last element
		if len(list.Elements) <= 1 {
			return NewErrorExpr("PartError",
				fmt.Sprintf("Last: expression %s has no elements", args[0].String()), args)
		}
		return list.Elements[len(list.Elements)-1] // Last element
	}

	// For atoms, Last is not defined
	return NewErrorExpr("PartError",
		fmt.Sprintf("Last: expression %s is not a list", args[0].String()), args)
}

// EvaluateRest evaluates Rest[expr] expressions - returns the expression with the first element (after head) removed
func EvaluateRest(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("Rest expects 1 argument, got %d", len(args)), args)
	}

	if list, ok := args[0].(List); ok {
		// For lists, return a new list with the first element after head removed
		if len(list.Elements) <= 1 {
			return NewErrorExpr("PartError",
				fmt.Sprintf("Rest: expression %s has no elements", args[0].String()), args)
		}

		// Create new list: head + elements[2:] (skip first element after head)
		if len(list.Elements) == 2 {
			// Special case: if only head and one element, return just the head
			return List{Elements: []Expr{list.Elements[0]}}
		}

		newElements := make([]Expr, len(list.Elements)-1)
		newElements[0] = list.Elements[0]        // Keep the head
		copy(newElements[1:], list.Elements[2:]) // Copy everything after the first element
		return List{Elements: newElements}
	}

	// For atoms, Rest is not defined
	return NewErrorExpr("PartError",
		fmt.Sprintf("Rest: expression %s is not a list", args[0].String()), args)
}

// EvaluateMost evaluates Most[expr] expressions - returns the expression with the last element removed
func EvaluateMost(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("Most expects 1 argument, got %d", len(args)), args)
	}

	if list, ok := args[0].(List); ok {
		// For lists, return a new list with the last element removed
		if len(list.Elements) <= 1 {
			return NewErrorExpr("PartError",
				fmt.Sprintf("Most: expression %s has no elements", args[0].String()), args)
		}

		// Create new list with all elements except the last one
		if len(list.Elements) == 2 {
			// Special case: if only head and one element, return just the head
			return List{Elements: []Expr{list.Elements[0]}}
		}

		newElements := make([]Expr, len(list.Elements)-1)
		copy(newElements, list.Elements[:len(list.Elements)-1])
		return List{Elements: newElements}
	}

	// For atoms, Most is not defined
	return NewErrorExpr("PartError",
		fmt.Sprintf("Most: expression %s is not a list", args[0].String()), args)
}

// EvaluateSetAttributes evaluates SetAttributes[symbol, attrs] expressions - sets attributes for a symbol
func EvaluateSetAttributes(args []Expr, ctx *Context) Expr {
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("SetAttributes expects 2 arguments, got %d", len(args)), args)
	}

	// Check for errors in arguments first
	for _, arg := range args {
		if IsError(arg) {
			return arg
		}
	}

	// First argument should be a symbol
	var symbolName string
	if atom, ok := args[0].(Atom); ok && atom.AtomType == SymbolAtom {
		symbolName = atom.Value.(string)
	} else {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("SetAttributes: first argument must be a symbol, got %s", args[0].String()), args)
	}

	// Second argument should be an attribute or list of attributes
	attributes := []Attribute{}

	if atom, ok := args[1].(Atom); ok && atom.AtomType == SymbolAtom {
		// Single attribute
		attrName := atom.Value.(string)
		if attr, err := parseAttribute(attrName); err != nil {
			return NewErrorExpr("AttributeError",
				fmt.Sprintf("SetAttributes: unknown attribute %s", attrName), args)
		} else {
			attributes = append(attributes, attr)
		}
	} else if list, ok := args[1].(List); ok {
		// List of attributes
		if len(list.Elements) < 1 {
			return NewErrorExpr("ArgumentError",
				"SetAttributes: attribute list cannot be empty", args)
		}

		// Skip the head element (should be "List")
		for i := 1; i < len(list.Elements); i++ {
			if atom, ok := list.Elements[i].(Atom); ok && atom.AtomType == SymbolAtom {
				attrName := atom.Value.(string)
				if attr, err := parseAttribute(attrName); err != nil {
					return NewErrorExpr("AttributeError",
						fmt.Sprintf("SetAttributes: unknown attribute %s", attrName), args)
				} else {
					attributes = append(attributes, attr)
				}
			} else {
				return NewErrorExpr("ArgumentError",
					fmt.Sprintf("SetAttributes: attribute must be a symbol, got %s", list.Elements[i].String()), args)
			}
		}
	} else {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("SetAttributes: second argument must be an attribute or list of attributes, got %s", args[1].String()), args)
	}

	// Set the attributes
	ctx.symbolTable.SetAttributes(symbolName, attributes)

	return NewSymbolAtom("Null")
}

// EvaluateClearAttributes evaluates ClearAttributes[symbol, attrs] expressions - clears specific attributes from a symbol
func EvaluateClearAttributes(args []Expr, ctx *Context) Expr {
	if len(args) == 0 || len(args) > 2 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("ClearAttributes expects 1 or 2 arguments, got %d", len(args)), args)
	}

	// Check for errors in arguments first
	for _, arg := range args {
		if IsError(arg) {
			return arg
		}
	}

	// First argument should be a symbol
	var symbolName string
	if atom, ok := args[0].(Atom); ok && atom.AtomType == SymbolAtom {
		symbolName = atom.Value.(string)
	} else {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("ClearAttributes: first argument must be a symbol, got %s", args[0].String()), args)
	}

	// If only one argument, clear all attributes
	if len(args) == 1 {
		ctx.symbolTable.ClearAllAttributes(symbolName)
		return NewSymbolAtom("Null")
	}

	// Second argument should be an attribute or list of attributes to clear
	attributesToClear := []Attribute{}

	if atom, ok := args[1].(Atom); ok && atom.AtomType == SymbolAtom {
		// Single attribute
		attrName := atom.Value.(string)
		if attr, err := parseAttribute(attrName); err != nil {
			return NewErrorExpr("AttributeError",
				fmt.Sprintf("ClearAttributes: unknown attribute %s", attrName), args)
		} else {
			attributesToClear = append(attributesToClear, attr)
		}
	} else if list, ok := args[1].(List); ok {
		// List of attributes
		if len(list.Elements) < 1 {
			return NewErrorExpr("ArgumentError",
				"ClearAttributes: attribute list cannot be empty", args)
		}

		// Skip the head element (should be "List")
		for i := 1; i < len(list.Elements); i++ {
			if atom, ok := list.Elements[i].(Atom); ok && atom.AtomType == SymbolAtom {
				attrName := atom.Value.(string)
				if attr, err := parseAttribute(attrName); err != nil {
					return NewErrorExpr("AttributeError",
						fmt.Sprintf("ClearAttributes: unknown attribute %s", attrName), args)
				} else {
					attributesToClear = append(attributesToClear, attr)
				}
			} else {
				return NewErrorExpr("ArgumentError",
					fmt.Sprintf("ClearAttributes: attribute must be a symbol, got %s", list.Elements[i].String()), args)
			}
		}
	} else {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("ClearAttributes: second argument must be an attribute or list of attributes, got %s", args[1].String()), args)
	}

	// Clear the specified attributes
	ctx.symbolTable.ClearAttributes(symbolName, attributesToClear)

	return NewSymbolAtom("Null")
}

// EvaluateSpecialClearAttributes handles the single-argument version of ClearAttributes[symbol]
func EvaluateSpecialClearAttributes(args []Expr, ctx *Context) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("ClearAttributes expects 1 argument, got %d", len(args)), args)
	}

	// Check for errors in arguments first
	if IsError(args[0]) {
		return args[0]
	}

	// First argument should be a symbol
	var symbolName string
	if atom, ok := args[0].(Atom); ok && atom.AtomType == SymbolAtom {
		symbolName = atom.Value.(string)
	} else {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("ClearAttributes: first argument must be a symbol, got %s", args[0].String()), args)
	}

	// Clear all attributes
	ctx.symbolTable.ClearAllAttributes(symbolName)
	return NewSymbolAtom("Null")
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
