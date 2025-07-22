package sexpr

import (
	"fmt"
	"math"
	"unicode/utf8"
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

// Arithmetic Operations

// EvaluatePlus evaluates Plus[...] expressions
func EvaluatePlus(args []Expr) Expr {
	if len(args) == 0 {
		return NewIntAtom(0) // Plus[] = 0
	}

	// Error propagation is now handled globally in wrapBuiltinFunc

	if len(args) == 1 {
		return args[0] // Plus[x] = x
	}

	// Check if all arguments are numeric
	allNumeric := true
	for _, arg := range args {
		if !isNumeric(arg) {
			allNumeric = false
			break
		}
	}

	if allNumeric {
		sum := 0.0
		for _, arg := range args {
			if val, ok := getNumericValue(arg); ok {
				sum += val
			}
		}
		return createNumericResult(sum)
	}

	// If not all numeric, return the expression unchanged
	elements := make([]Expr, len(args)+1)
	elements[0] = NewSymbolAtom("Plus")
	copy(elements[1:], args)
	return List{Elements: elements}
}

// EvaluateTimes evaluates Times[...] expressions
func EvaluateTimes(args []Expr) Expr {
	if len(args) == 0 {
		return NewIntAtom(1) // Times[] = 1
	}

	if len(args) == 1 {
		return args[0] // Times[x] = x
	}

	// Check if all arguments are numeric
	allNumeric := true
	for _, arg := range args {
		if !isNumeric(arg) {
			allNumeric = false
			break
		}
	}

	if allNumeric {
		product := 1.0
		for _, arg := range args {
			if val, ok := getNumericValue(arg); ok {
				product *= val
			}
		}
		return createNumericResult(product)
	}

	// If not all numeric, return the expression unchanged
	elements := make([]Expr, len(args)+1)
	elements[0] = NewSymbolAtom("Times")
	copy(elements[1:], args)
	return List{Elements: elements}
}

// EvaluateSubtract evaluates Subtract[a, b] expressions
func EvaluateSubtract(args []Expr) Expr {
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("Subtract expects 2 arguments, got %d", len(args)), args)
	}

	if isNumeric(args[0]) && isNumeric(args[1]) {
		val1, _ := getNumericValue(args[0])
		val2, _ := getNumericValue(args[1])
		return createNumericResult(val1 - val2)
	}

	// Return unchanged if not numeric
	return List{Elements: []Expr{NewSymbolAtom("Subtract"), args[0], args[1]}}
}

// EvaluateDivide evaluates Divide[a, b] expressions
func EvaluateDivide(args []Expr) Expr {
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("Divide expects 2 arguments, got %d", len(args)), args)
	}

	if isNumeric(args[0]) && isNumeric(args[1]) {
		val1, _ := getNumericValue(args[0])
		val2, _ := getNumericValue(args[1])

		if val2 == 0 {
			return NewErrorExpr("DivisionByZero",
				fmt.Sprintf("Division by zero: %v / %v", args[0], args[1]), args)
		}

		return createNumericResult(val1 / val2)
	}

	// Return unchanged if not numeric
	return List{Elements: []Expr{NewSymbolAtom("Divide"), args[0], args[1]}}
}

// EvaluatePower evaluates Power[a, b] expressions
func EvaluatePower(args []Expr) Expr {
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("Power expects 2 arguments, got %d", len(args)), args)
	}

	if isNumeric(args[0]) && isNumeric(args[1]) {
		base, _ := getNumericValue(args[0])
		exp, _ := getNumericValue(args[1])
		result := math.Pow(base, exp)

		// Check for invalid results (NaN, Inf)
		if math.IsNaN(result) || math.IsInf(result, 0) {
			return NewErrorExpr("MathematicalError",
				fmt.Sprintf("Invalid mathematical operation: %v ^ %v", args[0], args[1]), args)
		}

		return createNumericResult(result)
	}

	// Return unchanged if not numeric
	return List{Elements: []Expr{NewSymbolAtom("Power"), args[0], args[1]}}
}

// Comparison Operations

// EvaluateEqual evaluates Equal[a, b] expressions
func EvaluateEqual(args []Expr) Expr {
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("Equal expects 2 arguments, got %d", len(args)), args)
	}

	// Error propagation is now handled globally in wrapBuiltinFunc

	// Try numeric comparison first
	if isNumeric(args[0]) && isNumeric(args[1]) {
		val1, _ := getNumericValue(args[0])
		val2, _ := getNumericValue(args[1])
		return NewBoolAtom(val1 == val2)
	}

	// Try boolean comparison
	if isBool(args[0]) && isBool(args[1]) {
		val1, _ := getBoolValue(args[0])
		val2, _ := getBoolValue(args[1])
		return NewBoolAtom(val1 == val2)
	}

	// String comparison by representation
	return NewBoolAtom(args[0].String() == args[1].String())
}

// EvaluateUnequal evaluates Unequal[a, b] expressions
func EvaluateUnequal(args []Expr) Expr {
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("Unequal expects 2 arguments, got %d", len(args)), args)
	}

	equalResult := EvaluateEqual(args)
	if boolResult, ok := getBoolValue(equalResult); ok {
		return NewBoolAtom(!boolResult)
	}

	return NewBoolAtom(true)
}

// EvaluateLess evaluates Less[a, b] expressions
func EvaluateLess(args []Expr) Expr {
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("Less expects 2 arguments, got %d", len(args)), args)
	}

	if isNumeric(args[0]) && isNumeric(args[1]) {
		val1, _ := getNumericValue(args[0])
		val2, _ := getNumericValue(args[1])
		return NewBoolAtom(val1 < val2)
	}

	// Return unchanged if not numeric
	return List{Elements: []Expr{NewSymbolAtom("Less"), args[0], args[1]}}
}

// EvaluateGreater evaluates Greater[a, b] expressions
func EvaluateGreater(args []Expr) Expr {
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("Greater expects 2 arguments, got %d", len(args)), args)
	}

	if isNumeric(args[0]) && isNumeric(args[1]) {
		val1, _ := getNumericValue(args[0])
		val2, _ := getNumericValue(args[1])
		return NewBoolAtom(val1 > val2)
	}

	// Return unchanged if not numeric
	return List{Elements: []Expr{NewSymbolAtom("Greater"), args[0], args[1]}}
}

// EvaluateLessEqual evaluates LessEqual[a, b] expressions
func EvaluateLessEqual(args []Expr) Expr {
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("LessEqual expects 2 arguments, got %d", len(args)), args)
	}

	if isNumeric(args[0]) && isNumeric(args[1]) {
		val1, _ := getNumericValue(args[0])
		val2, _ := getNumericValue(args[1])
		return NewBoolAtom(val1 <= val2)
	}

	// Return unchanged if not numeric
	return List{Elements: []Expr{NewSymbolAtom("LessEqual"), args[0], args[1]}}
}

// EvaluateGreaterEqual evaluates GreaterEqual[a, b] expressions
func EvaluateGreaterEqual(args []Expr) Expr {
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("GreaterEqual expects 2 arguments, got %d", len(args)), args)
	}

	if isNumeric(args[0]) && isNumeric(args[1]) {
		val1, _ := getNumericValue(args[0])
		val2, _ := getNumericValue(args[1])
		return NewBoolAtom(val1 >= val2)
	}

	// Return unchanged if not numeric
	return List{Elements: []Expr{NewSymbolAtom("GreaterEqual"), args[0], args[1]}}
}

// Logical Operations

// EvaluateNot evaluates Not[x] expressions
func EvaluateNot(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("Not expects 1 argument, got %d", len(args)), args)
	}

	if isBool(args[0]) {
		val, _ := getBoolValue(args[0])
		return NewBoolAtom(!val)
	}

	// Return unchanged if not boolean
	return List{Elements: []Expr{NewSymbolAtom("Not"), args[0]}}
}

// EvaluateSameQ evaluates SameQ[a, b] expressions (identical objects)
func EvaluateSameQ(args []Expr) Expr {
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("SameQ expects 2 arguments, got %d", len(args)), args)
	}

	// SameQ requires structural equality
	return NewBoolAtom(args[0].Equal(args[1]))
}

// EvaluateUnsameQ evaluates UnsameQ[a, b] expressions
func EvaluateUnsameQ(args []Expr) Expr {
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("UnsameQ expects 2 arguments, got %d", len(args)), args)
	}

	sameQResult := EvaluateSameQ(args)
	if boolResult, ok := getBoolValue(sameQResult); ok {
		return NewBoolAtom(!boolResult)
	}

	return NewBoolAtom(true)
}

// Introspection Operations

// EvaluateHead evaluates Head[expr] expressions - returns the head/type of an expression
func EvaluateHead(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("Head expects 1 argument, got %d", len(args)), args)
	}

	// Note: Head should work on any expression, including errors
	// So we don't propagate errors here, but analyze the error's type

	expr := args[0]
	var head string

	switch ex := expr.(type) {
	case Atom:
		switch ex.AtomType {
		case IntAtom:
			head = "Integer"
		case FloatAtom:
			head = "Real"
		case StringAtom:
			head = "String"
		case SymbolAtom:
			head = "Symbol"
		default:
			head = "Unknown"
		}
	case List:
		if len(ex.Elements) == 0 {
			head = "List"
		} else {
			// For non-empty lists, the head is the first element
			// This matches Mathematica semantics where f[x,y] has head f
			return ex.Elements[0]
		}
	case *ErrorExpr:
		head = "Error"
	case ObjectExpr:
		head = ex.TypeName
	default:
		head = "Unknown"
	}

	return NewSymbolAtom(head)
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

// Utility/Predicate Functions

// EvaluateLength evaluates Length[expr] expressions - returns the length/size of expression
func EvaluateLength(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("Length expects 1 argument, got %d", len(args)), args)
	}

	expr := args[0]

	switch ex := expr.(type) {
	case List:
		// For lists, return the number of elements (excluding the head)
		if len(ex.Elements) == 0 {
			return NewIntAtom(0) // Empty list has length 0
		}
		return NewIntAtom(len(ex.Elements) - 1) // Subtract 1 for the head
	case ObjectExpr:
		// Handle Association
		if ex.TypeName == "Association" {
			if assocValue, ok := ex.Value.(AssociationValue); ok {
				return NewIntAtom(assocValue.Len())
			}
		}
		// For other ObjectExpr types, length is 0
		return NewIntAtom(0)
	default:
		// For atoms and other expressions, length is 0
		return NewIntAtom(0)
	}
}

// EvaluateListQ evaluates ListQ[expr] expressions - returns True if expression is a List
func EvaluateListQ(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("ListQ expects 1 argument, got %d", len(args)), args)
	}

	_, isList := args[0].(List)
	return NewBoolAtom(isList)
}

// EvaluateNumberQ evaluates NumberQ[expr] expressions - returns True if expression is numeric
func EvaluateNumberQ(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("NumberQ expects 1 argument, got %d", len(args)), args)
	}

	return NewBoolAtom(isNumeric(args[0]))
}

// EvaluateBooleanQ evaluates BooleanQ[expr] expressions - returns True if expression is a boolean
func EvaluateBooleanQ(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("BooleanQ expects 1 argument, got %d", len(args)), args)
	}

	return NewBoolAtom(isBool(args[0]))
}

// EvaluateIntegerQ evaluates IntegerQ[expr] expressions - returns True if expression is an integer
func EvaluateIntegerQ(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("IntegerQ expects 1 argument, got %d", len(args)), args)
	}

	if atom, ok := args[0].(Atom); ok {
		return NewBoolAtom(atom.AtomType == IntAtom)
	}
	return NewBoolAtom(false)
}

// EvaluateFloatQ evaluates FloatQ[expr] expressions - returns True if expression is a float
func EvaluateFloatQ(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("FloatQ expects 1 argument, got %d", len(args)), args)
	}

	if atom, ok := args[0].(Atom); ok {
		return NewBoolAtom(atom.AtomType == FloatAtom)
	}
	return NewBoolAtom(false)
}

// EvaluateAtomQ evaluates AtomQ[expr] expressions - returns True if expression is an atom
func EvaluateAtomQ(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("AtomQ expects 1 argument, got %d", len(args)), args)
	}

	_, isAtom := args[0].(Atom)
	return NewBoolAtom(isAtom)
}

// EvaluateSymbolQ evaluates SymbolQ[expr] expressions - returns True if expression is a symbol
func EvaluateSymbolQ(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("SymbolQ expects 1 argument, got %d", len(args)), args)
	}

	return NewBoolAtom(isSymbol(args[0]))
}

// EvaluateStringQ evaluates StringQ[expr] expressions - returns True if expression is a string
func EvaluateStringQ(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("StringQ expects 1 argument, got %d", len(args)), args)
	}

	if atom, ok := args[0].(Atom); ok {
		return NewBoolAtom(atom.AtomType == StringAtom)
	}
	return NewBoolAtom(false)
}

// EvaluateStringLength evaluates StringLength[expr] expressions - returns length if string, error otherwise
func EvaluateStringLength(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("StringLength expects 1 argument, got %d", len(args)), args)
	}

	if atom, ok := args[0].(Atom); ok && atom.AtomType == StringAtom {
		str := atom.Value.(string)
		return NewIntAtom(utf8.RuneCountInString(str))
	}

	return NewErrorExpr("ArgumentError",
		fmt.Sprintf("StringLength expects a string, got %s", args[0].Type()), args)
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

// EvaluatePart evaluates Part[expr, index] expressions - returns the element at the specified 1-based index
func EvaluatePart(args []Expr) Expr {
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("Part expects 2 arguments, got %d", len(args)), args)
	}

	expr := args[0]
	indexExpr := args[1]

	// Handle Association access first (key-based)
	if objExpr, ok := expr.(ObjectExpr); ok && objExpr.TypeName == "Association" {
		if assocValue, ok := objExpr.Value.(AssociationValue); ok {
			// For associations, use the index argument as a key
			if value, exists := assocValue.Get(indexExpr); exists {
				return value
			}
			return NewErrorExpr("PartError",
				fmt.Sprintf("Key %s not found in association", indexExpr.String()), args)
		}
	}

	// Handle List access (integer-based indexing)
	if list, ok := expr.(List); ok {
		// Extract integer index - handle both direct integers and Minus[n] expressions
		var index int
		if indexAtom, ok := indexExpr.(Atom); ok && indexAtom.AtomType == IntAtom {
			// Direct integer atom
			index = indexAtom.Value.(int)
		} else if indexList, ok := indexExpr.(List); ok && len(indexList.Elements) == 2 {
			// Check for Minus[n] pattern (negative number)
			if headAtom, ok := indexList.Elements[0].(Atom); ok &&
				headAtom.AtomType == SymbolAtom && headAtom.Value.(string) == "Minus" {
				if valueAtom, ok := indexList.Elements[1].(Atom); ok && valueAtom.AtomType == IntAtom {
					index = -valueAtom.Value.(int)
				} else {
					return NewErrorExpr("PartError",
						fmt.Sprintf("Part index must be an integer for lists, got %s", indexExpr.String()), args)
				}
			} else {
				return NewErrorExpr("PartError",
					fmt.Sprintf("Part index must be an integer for lists, got %s", indexExpr.String()), args)
			}
		} else {
			return NewErrorExpr("PartError",
				fmt.Sprintf("Part index must be an integer for lists, got %s", indexExpr.String()), args)
		}
		// For lists, return the element at the specified index (1-based)
		if len(list.Elements) <= 1 {
			return NewErrorExpr("PartError",
				fmt.Sprintf("Part: expression %s has no elements", expr.String()), args)
		}

		// Handle negative indexing: -1 is last element, -2 is second to last, etc.
		var actualIndex int
		if index < 0 {
			// Negative indexing: -1 = last, -2 = second to last, etc.
			actualIndex = len(list.Elements) + index
		} else if index > 0 {
			// Positive 1-based indexing: convert to 0-based for internal use
			actualIndex = index
		} else {
			// index == 0 is invalid in 1-based indexing
			return NewErrorExpr("PartError",
				fmt.Sprintf("Part index %d is out of bounds (indices start at 1)", index), args)
		}

		// Check bounds (remember: list.Elements[0] is the head, actual elements start at index 1)
		if actualIndex < 1 || actualIndex >= len(list.Elements) {
			return NewErrorExpr("PartError",
				fmt.Sprintf("Part index %d is out of bounds for expression with %d elements",
					index, len(list.Elements)-1), args)
		}

		return list.Elements[actualIndex]
	}

	// For atoms and other expressions, Part is not defined
	return NewErrorExpr("PartError",
		fmt.Sprintf("Part: expression %s is not a list or association", expr.String()), args)
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

// EvaluateAssociation evaluates Association[...] expressions - creates association objects
func EvaluateAssociation(args []Expr) Expr {
	assocValue := NewAssociationValue()

	// Process each argument as a Rule[key, value] expression
	for _, arg := range args {
		if ruleList, ok := arg.(List); ok && len(ruleList.Elements) == 3 {
			// Check if this is Rule[key, value]
			if headAtom, ok := ruleList.Elements[0].(Atom); ok &&
				headAtom.AtomType == SymbolAtom && headAtom.Value.(string) == "Rule" {

				key := ruleList.Elements[1]
				value := ruleList.Elements[2]
				assocValue = assocValue.Set(key, value) // Now returns new association
				continue
			}
		}

		// Invalid argument - not a Rule
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("Association expects Rule expressions, got %s", arg.String()), args)
	}

	return NewObjectExpr("Association", assocValue)
}

// EvaluateAssociationQ evaluates AssociationQ[expr] expressions - returns True if expression is an Association
func EvaluateAssociationQ(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("AssociationQ expects 1 argument, got %d", len(args)), args)
	}

	if objExpr, ok := args[0].(ObjectExpr); ok && objExpr.TypeName == "Association" {
		return NewSymbolAtom("True")
	}
	return NewSymbolAtom("False")
}

// EvaluateKeys evaluates Keys[assoc] expressions - returns the keys of an association
func EvaluateKeys(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("Keys expects 1 argument, got %d", len(args)), args)
	}

	if objExpr, ok := args[0].(ObjectExpr); ok && objExpr.TypeName == "Association" {
		if assocValue, ok := objExpr.Value.(AssociationValue); ok {
			keys := assocValue.Keys()
			// Return as List[key1, key2, ...]
			elements := []Expr{NewSymbolAtom("List")}
			elements = append(elements, keys...)
			return NewList(elements...)
		}
	}

	return NewErrorExpr("ArgumentError",
		fmt.Sprintf("Keys expects an Association, got %s", args[0].String()), args)
}

// EvaluateValues evaluates Values[assoc] expressions - returns the values of an association
func EvaluateValues(args []Expr) Expr {
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("Values expects 1 argument, got %d", len(args)), args)
	}

	if objExpr, ok := args[0].(ObjectExpr); ok && objExpr.TypeName == "Association" {
		if assocValue, ok := objExpr.Value.(AssociationValue); ok {
			values := assocValue.Values()
			// Return as List[value1, value2, ...]
			elements := []Expr{NewSymbolAtom("List")}
			elements = append(elements, values...)
			return NewList(elements...)
		}
	}

	return NewErrorExpr("ArgumentError",
		fmt.Sprintf("Values expects an Association, got %s", args[0].String()), args)
}
