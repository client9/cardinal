package sexpr

import (
	"fmt"
	"sort"
)

// PatternFunc represents a Go function that can be called with pattern-matched arguments
type PatternFunc func(args []Expr, ctx *Context) Expr

// FunctionDef represents a single function definition with pattern and implementation
type FunctionDef struct {
	Pattern     Expr        // The pattern to match (e.g., Plus(x_Integer, y_Integer))
	Body        Expr        // The body expression for user-defined functions (nil for Go implementations)
	GoImpl      PatternFunc // Go implementation for built-in functions (nil for user-defined)
	Specificity int         // Auto-calculated pattern specificity for ordering
	IsBuiltin   bool        // Whether this definition came from system registration
}

// FunctionRegistry manages all function definitions (user-defined and built-in) with pattern-based dispatch
type FunctionRegistry struct {
	functions map[string][]FunctionDef // function name -> ordered list of patterns
}

// NewFunctionRegistry creates a new function registry
func NewFunctionRegistry() *FunctionRegistry {
	return &FunctionRegistry{
		functions: make(map[string][]FunctionDef),
	}
}

// matchBlankExpression matches a blank expression (Blank[], BlankSequence[], BlankNullSequence[]) against an expression
func matchBlankExpression(blankExpr Expr, expr Expr, ctx *Context) bool {
	if isBlank, blankType, typeExpr := isSymbolicBlank(blankExpr); isBlank {
		// Check type constraint if present
		if typeExpr != nil {
			var typeName string
			if typeAtom, ok := typeExpr.(Atom); ok && typeAtom.AtomType == SymbolAtom {
				typeName = typeAtom.Value.(string)
			}
			if !matchesType(expr, typeName) {
				return false
			}
		}

		// For now, single blank expressions always match single expressions
		// BlankSequence and BlankNullSequence handling for sequences happens elsewhere
		switch blankType {
		case "Blank":
			return true // Single expression always matches Blank[]
		case "BlankSequence":
			return true // Single expression matches BlankSequence[] (at least one)
		case "BlankNullSequence":
			return true // Single expression matches BlankNullSequence[] (zero or more)
		}
	}
	return false
}

// convertParsedPatternToSymbolic converts a parsed pattern to use symbolic pattern representation
func convertParsedPatternToSymbolic(pattern Expr) Expr {
	switch p := pattern.(type) {
	case Atom:
		if p.AtomType == SymbolAtom {
			patternStr := p.Value.(string)
			// Convert pattern strings like "x_", "_Integer", etc. to symbolic
			return ConvertPatternStringToSymbolic(patternStr)
		}
		return p
	case List:
		// Convert all elements in the list
		newElements := make([]Expr, len(p.Elements))
		for i, elem := range p.Elements {
			newElements[i] = convertParsedPatternToSymbolic(elem)
		}
		return List{Elements: newElements}
	default:
		return pattern
	}
}

// RegisterPatternBuiltin registers a built-in function with a pattern from Go code
func (r *FunctionRegistry) RegisterPatternBuiltin(patternStr string, impl PatternFunc) error {
	// Parse the pattern string
	pattern, err := ParseString(patternStr)
	if err != nil {
		return fmt.Errorf("invalid pattern syntax: %v", err)
	}

	// Convert parsed pattern to symbolic representation
	symbolicPattern := convertParsedPatternToSymbolic(pattern)

	// Debug: print what was parsed and converted
	// fmt.Printf("DEBUG: Parsed pattern '%s' -> %v\n", patternStr, pattern)
	// fmt.Printf("DEBUG: Symbolic pattern -> %v\n", symbolicPattern)

	// Extract function name from original pattern (before conversion)
	functionName, err := extractFunctionName(pattern)
	if err != nil {
		return fmt.Errorf("cannot extract function name from pattern: %v", err)
	}

	// Create function definition with symbolic pattern
	funcDef := FunctionDef{
		Pattern:     symbolicPattern,
		Body:        nil,
		GoImpl:      impl,
		Specificity: calculatePatternSpecificity(symbolicPattern),
		IsBuiltin:   true,
	}

	// Register the function definition
	r.registerFunctionDef(functionName, funcDef)
	return nil
}

// RegisterPatternBuiltins registers multiple built-in functions from a map
func (r *FunctionRegistry) RegisterPatternBuiltins(patterns map[string]PatternFunc) error {
	for patternStr, impl := range patterns {
		if err := r.RegisterPatternBuiltin(patternStr, impl); err != nil {
			return fmt.Errorf("failed to register pattern %s: %v", patternStr, err)
		}
	}
	return nil
}

// RegisterUserFunction registers a user-defined function with pattern and body
func (r *FunctionRegistry) RegisterUserFunction(pattern Expr, body Expr) error {
	// Extract function name from pattern
	functionName, err := extractFunctionName(pattern)
	if err != nil {
		return fmt.Errorf("cannot extract function name from pattern: %v", err)
	}

	// fmt.Printf("DEBUG: RegisterUserFunction: function=%s, pattern=%v, body=%v\\n", functionName, pattern, body)

	// Create function definition
	funcDef := FunctionDef{
		Pattern:     pattern,
		Body:        body,
		GoImpl:      nil,
		Specificity: calculatePatternSpecificity(pattern),
		IsBuiltin:   false,
	}

	// Register the function definition (will replace equivalent patterns)
	r.registerFunctionDef(functionName, funcDef)
	return nil
}

// FindMatchingFunction finds the best matching function definition for given arguments
func (r *FunctionRegistry) FindMatchingFunction(functionName string, args []Expr) (*FunctionDef, map[string]Expr) {
	definitions, exists := r.functions[functionName]
	if !exists {
		// fmt.Printf("DEBUG: No functions found for '%s'. Available functions: %v\\n", functionName, r.GetAllFunctionNames())
		return nil, nil
	}
	// fmt.Printf("DEBUG: Found %d definitions for function '%s'\\n", len(definitions), functionName)

	// Try each definition in order (most specific first)
	for _, def := range definitions {
		// fmt.Printf("DEBUG: Trying to match pattern %v against %s(%v)\\n", def.Pattern, functionName, args)
		if matches, bindings := matchesPattern(def.Pattern, functionName, args); matches {
			// fmt.Printf("DEBUG: Pattern MATCHED with bindings: %v\\n", bindings)
			return &def, bindings
		}
		// fmt.Printf("DEBUG: Pattern did NOT match\\n")
	}
	return nil, nil
}

// GetFunctionDefinitions returns all definitions for a function name (for debugging/introspection)
func (r *FunctionRegistry) GetFunctionDefinitions(functionName string) []FunctionDef {
	if definitions, exists := r.functions[functionName]; exists {
		// Return a copy to prevent external modification
		result := make([]FunctionDef, len(definitions))
		copy(result, definitions)
		return result
	}
	return nil
}

// GetAllFunctionNames returns all registered function names
func (r *FunctionRegistry) GetAllFunctionNames() []string {
	names := make([]string, 0, len(r.functions))
	for name := range r.functions {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// CallFunction attempts to call a function with the given expression and returns (result, found)
func (r *FunctionRegistry) CallFunction(callExpr Expr, ctx *Context) (Expr, bool) {
	// Extract function name and arguments from the call expression
	if list, ok := callExpr.(List); ok && len(list.Elements) > 0 {
		if headAtom, ok := list.Elements[0].(Atom); ok && headAtom.AtomType == SymbolAtom {
			functionName := headAtom.Value.(string)
			args := list.Elements[1:]

			// Find matching function definition
			funcDef, bindings := r.FindMatchingFunction(functionName, args)
			if funcDef == nil {
				return nil, false
			}

			// Create child context with pattern variable bindings
			funcCtx := NewChildContext(ctx)
			for varName, value := range bindings {
				funcCtx.Set(varName, value)
			}

			// Call the function
			if funcDef.GoImpl != nil {
				// Built-in function - call Go implementation
				return funcDef.GoImpl(args, funcCtx), true
			} else {
				// User-defined function - evaluate body
				// Create an evaluator to evaluate the body
				evaluator := NewEvaluatorWithContext(funcCtx)
				return evaluator.evaluate(funcDef.Body, funcCtx), true
			}
		}
	}
	return nil, false
}

// RegisterFunction is an alias for RegisterUserFunction for backward compatibility
func (r *FunctionRegistry) RegisterFunction(functionName string, pattern Expr, implementation func([]Expr, *Context) Expr) error {
	// This is a simplified version that assumes the pattern contains the function name
	// For the refactored code, we need to create a proper function definition
	funcDef := FunctionDef{
		Pattern:     pattern,
		Body:        nil,
		GoImpl:      implementation,
		Specificity: calculatePatternSpecificity(pattern),
		IsBuiltin:   false,
	}

	r.registerFunctionDef(functionName, funcDef)
	return nil
}

// registerFunctionDef adds or replaces a function definition
func (r *FunctionRegistry) registerFunctionDef(functionName string, newDef FunctionDef) {
	definitions := r.functions[functionName]

	// Check if we need to replace an existing equivalent pattern
	for i, existingDef := range definitions {
		if patternsEquivalent(existingDef.Pattern, newDef.Pattern) {
			// Replace existing definition
			definitions[i] = newDef
			r.functions[functionName] = definitions
			return
		}
	}

	// Add new definition and re-sort by specificity
	definitions = append(definitions, newDef)
	sort.Slice(definitions, func(i, j int) bool {
		// Higher specificity comes first
		return definitions[i].Specificity > definitions[j].Specificity
	})

	r.functions[functionName] = definitions
}

// extractFunctionName extracts the function name from a pattern
func extractFunctionName(pattern Expr) (string, error) {
	switch p := pattern.(type) {
	case Atom:
		if p.AtomType == SymbolAtom {
			return p.Value.(string), nil
		}
		return "", fmt.Errorf("pattern must be a symbol or function call")
	case List:
		if len(p.Elements) == 0 {
			return "", fmt.Errorf("empty list pattern")
		}
		if head, ok := p.Elements[0].(Atom); ok && head.AtomType == SymbolAtom {
			return head.Value.(string), nil
		}
		return "", fmt.Errorf("pattern head must be a symbol")
	default:
		return "", fmt.Errorf("invalid pattern type")
	}
}

// matchesPattern checks if a pattern matches the given arguments and returns variable bindings
func matchesPattern(pattern Expr, functionName string, args []Expr) (bool, map[string]Expr) {
	// Create a mock function call to match against the pattern
	mockCall := List{Elements: make([]Expr, len(args)+1)}
	mockCall.Elements[0] = NewSymbolAtom(functionName)
	copy(mockCall.Elements[1:], args)

	// fmt.Printf("DEBUG: matchesPattern: pattern=%v, mockCall=%v\\n", pattern, mockCall)

	// Use a simple pattern matching approach without creating a new evaluator
	// This avoids infinite recursion from evaluator creation
	ctx := NewContext()

	// Store the initial variable names to track what gets bound
	initialVars := make(map[string]bool)
	for varName := range ctx.variables {
		initialVars[varName] = true
	}

	// Use the direct pattern matching function
	matches := directMatchPattern(pattern, mockCall, ctx)

	if matches {
		// Extract the variable bindings that were created during pattern matching
		bindings := make(map[string]Expr)
		for varName, value := range ctx.variables {
			// Only include variables that were bound during this match
			if !initialVars[varName] {
				bindings[varName] = value
			}
		}
		return true, bindings
	}
	return false, nil
}

// directMatchPattern is a simplified pattern matching function without evaluator dependencies
func directMatchPattern(pattern Expr, expr Expr, ctx *Context) bool {
	return directMatchPatternWithContext(pattern, expr, ctx, false)
}

// directMatchPatternWithContext handles pattern matching with context about parameter vs literal positions
func directMatchPatternWithContext(pattern Expr, expr Expr, ctx *Context, isParameter bool) bool {
	// fmt.Printf("DEBUG: directMatchPattern: pattern=%v (%T), expr=%v (%T)\\n", pattern, pattern, expr, expr)

	// First, check if this is a symbolic pattern (new system)
	if isPattern, nameExpr, blankExpr := isSymbolicPattern(pattern); isPattern {
		// Handle Pattern[name, blank] - extract the name and match the blank
		varName := ""
		if nameAtom, ok := nameExpr.(Atom); ok && nameAtom.AtomType == SymbolAtom {
			varName = nameAtom.Value.(string)
		}

		// Match the blank part
		if matchBlankExpression(blankExpr, expr, ctx) {
			// Bind the variable if we have a name
			if varName != "" {
				ctx.Set(varName, expr)
			}
			return true
		}
		return false
	}

	// Check if this is a direct symbolic blank (new system)
	if isBlank, _, _ := isSymbolicBlank(pattern); isBlank {
		return matchBlankExpression(pattern, expr, ctx)
	}

	// Fall back to string-based pattern matching (legacy system)
	switch p := pattern.(type) {
	case Atom:
		if p.AtomType == SymbolAtom {
			varName := p.Value.(string)
			// fmt.Printf("DEBUG: Symbol atom case: varName='%s'\\n", varName)
			// Check if it's a pattern variable
			if isPatternVariable(varName) {
				// fmt.Printf("DEBUG: '%s' is a pattern variable\\n", varName)
				info := parsePatternInfo(varName)
				if info.Type == BlankNullSequencePattern {
					// This should not happen for single expressions - sequence patterns need special handling
					return false
				} else if info.Type == BlankSequencePattern {
					// This should not happen for single expressions - sequence patterns need special handling
					return false
				} else if info.Type == BlankPattern {
					// Check type constraint if present
					if !matchesType(expr, info.TypeName) {
						return false
					}
					// Bind the variable
					if info.VarName != "" {
						ctx.Set(info.VarName, expr)
					}
					return true
				}
			} else {
				// fmt.Printf("DEBUG: '%s' is NOT a pattern variable\\n", varName)
				// Regular symbol (not a pattern variable)
				if isParameter {
					// In parameter position, regular symbols bind to arguments
					ctx.Set(varName, expr)
					return true
				} else {
					// Not in parameter position, must match literally
					if exprAtom, ok := expr.(Atom); ok && exprAtom.AtomType == SymbolAtom {
						// fmt.Printf("DEBUG: Comparing literal symbols: pattern '%s' vs expr '%s'\\n", varName, exprAtom.Value.(string))
						result := varName == exprAtom.Value.(string)
						// fmt.Printf("DEBUG: Literal symbol comparison result: %v\\n", result)
						return result
					}
					// fmt.Printf("DEBUG: expr is not a symbol atom: %T\\n", expr)
					return false
				}
			}
		}
		// For literal atoms, they must be exactly equal
		if exprAtom, ok := expr.(Atom); ok {
			// fmt.Printf("DEBUG: Comparing literal atoms: pattern %v (%T, type=%v, value=%v) vs expr %v (%T, type=%v, value=%v)\\n",
			//	p, p, p.AtomType, p.Value, exprAtom, exprAtom, exprAtom.AtomType, exprAtom.Value)
			result := p.AtomType == exprAtom.AtomType && p.Value == exprAtom.Value
			// fmt.Printf("DEBUG: Literal atom comparison result: %v\\n", result)
			return result
		}
		return false

	case List:
		if exprList, ok := expr.(List); ok {
			// Both are lists - need to match structure and handle sequence patterns
			return matchListPatternWithContext(p, exprList, ctx, isParameter)
		}
		return false

	default:
		return false
	}
}

// matchListPattern matches a list pattern against a list expression
func matchListPattern(patternList List, exprList List, ctx *Context) bool {
	return matchListPatternWithContext(patternList, exprList, ctx, false)
}

// matchListPatternWithContext matches a list pattern against a list expression with context
func matchListPatternWithContext(patternList List, exprList List, ctx *Context, isParameter bool) bool {
	// Handle empty patterns
	if len(patternList.Elements) == 0 {
		return len(exprList.Elements) == 0
	}

	patternIdx := 0
	exprIdx := 0

	for patternIdx < len(patternList.Elements) && exprIdx < len(exprList.Elements) {
		patternElem := patternList.Elements[patternIdx]

		// Check if this pattern element is a sequence pattern (symbolic or string-based)

		// First check for symbolic Pattern[name, BlankSequence/BlankNullSequence]
		if isPattern, nameExpr, blankExpr := isSymbolicPattern(patternElem); isPattern {
			if isBlank, blankType, typeExpr := isSymbolicBlank(blankExpr); isBlank {
				varName := ""
				if nameAtom, ok := nameExpr.(Atom); ok && nameAtom.AtomType == SymbolAtom {
					varName = nameAtom.Value.(string)
				}

				typeName := ""
				if typeExpr != nil {
					if typeAtom, ok := typeExpr.(Atom); ok && typeAtom.AtomType == SymbolAtom {
						typeName = typeAtom.Value.(string)
					}
				}

				// Create PatternInfo for compatibility with existing sequence matching
				info := PatternInfo{
					VarName:  varName,
					TypeName: typeName,
				}

				if blankType == "BlankNullSequence" {
					info.Type = BlankNullSequencePattern
					return matchSequencePattern(patternList, patternIdx, exprList, exprIdx, ctx, info, true)
				} else if blankType == "BlankSequence" {
					info.Type = BlankSequencePattern
					return matchSequencePattern(patternList, patternIdx, exprList, exprIdx, ctx, info, false)
				}
			}
		}

		// Check for direct symbolic BlankSequence/BlankNullSequence
		if isBlank, blankType, typeExpr := isSymbolicBlank(patternElem); isBlank {
			typeName := ""
			if typeExpr != nil {
				if typeAtom, ok := typeExpr.(Atom); ok && typeAtom.AtomType == SymbolAtom {
					typeName = typeAtom.Value.(string)
				}
			}

			// Create PatternInfo for compatibility
			info := PatternInfo{
				VarName:  "", // Anonymous sequence
				TypeName: typeName,
			}

			if blankType == "BlankNullSequence" {
				info.Type = BlankNullSequencePattern
				return matchSequencePattern(patternList, patternIdx, exprList, exprIdx, ctx, info, true)
			} else if blankType == "BlankSequence" {
				info.Type = BlankSequencePattern
				return matchSequencePattern(patternList, patternIdx, exprList, exprIdx, ctx, info, false)
			}
		}

		// Fall back to string-based sequence pattern detection
		if atom, ok := patternElem.(Atom); ok && atom.AtomType == SymbolAtom {
			varName := atom.Value.(string)
			if isPatternVariable(varName) {
				info := parsePatternInfo(varName)

				if info.Type == BlankNullSequencePattern {
					// Handle x___ pattern - match zero or more elements
					return matchSequencePattern(patternList, patternIdx, exprList, exprIdx, ctx, info, true)
				} else if info.Type == BlankSequencePattern {
					// Handle x__ pattern - match one or more elements
					return matchSequencePattern(patternList, patternIdx, exprList, exprIdx, ctx, info, false)
				}
			}
		}

		// Regular pattern element - must match exactly one expression element
		// Element 0 is head (literal), elements 1+ are parameters (bind)
		isParameterPosition := exprIdx > 0
		if !directMatchPatternWithContext(patternElem, exprList.Elements[exprIdx], ctx, isParameterPosition) {
			return false
		}

		patternIdx++
		exprIdx++
	}

	// Check if we've consumed all elements appropriately
	if patternIdx < len(patternList.Elements) {
		// Remaining pattern elements - check if they're optional (null sequence patterns)
		for patternIdx < len(patternList.Elements) {
			patternElem := patternList.Elements[patternIdx]

			// Check for symbolic null sequence patterns
			if isPattern, nameExpr, blankExpr := isSymbolicPattern(patternElem); isPattern {
				if isBlank, blankType, _ := isSymbolicBlank(blankExpr); isBlank && blankType == "BlankNullSequence" {
					// Bind to empty list
					if nameAtom, ok := nameExpr.(Atom); ok && nameAtom.AtomType == SymbolAtom {
						varName := nameAtom.Value.(string)
						if varName != "" {
							ctx.Set(varName, List{Elements: []Expr{}})
						}
					}
					patternIdx++
					continue
				}
			}

			// Check for direct symbolic BlankNullSequence
			if isBlank, blankType, _ := isSymbolicBlank(patternElem); isBlank && blankType == "BlankNullSequence" {
				// Anonymous null sequence - just continue
				patternIdx++
				continue
			}

			// Fall back to string-based null sequence patterns
			if atom, ok := patternElem.(Atom); ok && atom.AtomType == SymbolAtom {
				varName := atom.Value.(string)
				if isPatternVariable(varName) {
					info := parsePatternInfo(varName)
					if info.Type == BlankNullSequencePattern {
						// Bind to empty list
						if info.VarName != "" {
							ctx.Set(info.VarName, List{Elements: []Expr{}})
						}
						patternIdx++
						continue
					}
				}
			}
			return false // Non-optional pattern element without matching expression
		}
	}

	return exprIdx == len(exprList.Elements) // All expression elements consumed
}

// matchSequencePattern handles sequence patterns (x__ and x___)
func matchSequencePattern(patternList List, patternIdx int, exprList List, exprIdx int, ctx *Context, info PatternInfo, allowZero bool) bool {
	// Calculate how many elements this sequence pattern should consume
	remainingPatterns := len(patternList.Elements) - patternIdx - 1
	remainingExprs := len(exprList.Elements) - exprIdx

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
		// Create a copy of context for this attempt
		testCtx := NewChildContext(ctx)

		// Collect the elements to bind to this sequence
		var seqElements []Expr
		for i := 0; i < consume; i++ {
			if exprIdx+i < len(exprList.Elements) {
				seqElements = append(seqElements, exprList.Elements[exprIdx+i])
			}
		}

		// Check type constraints for sequence elements
		if info.TypeName != "" {
			allMatch := true
			for _, elem := range seqElements {
				if !matchesType(elem, info.TypeName) {
					allMatch = false
					break
				}
			}
			if !allMatch {
				continue // Try with fewer elements - this attempt fails
			}
		}

		// Bind the sequence variable
		if info.VarName != "" {
			// fmt.Printf("DEBUG: Binding sequence var '%s' to List with %d elements: %v\n", info.VarName, len(seqElements), seqElements)
			// Create a proper List with "List" as the first element
			listElements := append([]Expr{NewSymbolAtom("List")}, seqElements...)
			testCtx.Set(info.VarName, List{Elements: listElements})
		}

		// Try to match remaining patterns
		if patternIdx+1 >= len(patternList.Elements) {
			// This was the last pattern - check if we consumed all expressions
			if exprIdx+consume == len(exprList.Elements) {
				// Copy bindings to original context
				for varName, value := range testCtx.variables {
					if _, exists := ctx.variables[varName]; !exists {
						// fmt.Printf("DEBUG: Copying variable '%s' = %v (type: %T)\n", varName, value, value)
						ctx.Set(varName, value)
					}
				}
				return true
			}
		} else {
			// Try to match remaining patterns
			remainingPattern := List{Elements: patternList.Elements[patternIdx+1:]}
			remainingExpr := List{Elements: exprList.Elements[exprIdx+consume:]}

			if matchListPatternWithContext(remainingPattern, remainingExpr, testCtx, true) {
				// Copy bindings to original context
				for varName, value := range testCtx.variables {
					if _, exists := ctx.variables[varName]; !exists {
						ctx.Set(varName, value)
					}
				}
				return true
			}
		}
	}

	return false
}

// calculatePatternSpecificity calculates the specificity score for a pattern
// Higher scores indicate more specific patterns
func calculatePatternSpecificity(pattern Expr) int {
	return int(getPatternSpecificity(pattern))
}

// patternsEquivalent checks if two patterns are structurally equivalent (ignoring variable names)
func patternsEquivalent(pattern1, pattern2 Expr) bool {
	return patternsEqual(pattern1, pattern2)
}
