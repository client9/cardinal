package sexpr

import (
	"fmt"
	"sort"

	"github.com/client9/sexpr/core"
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

// calculatePatternSpecificity calculates specificity for compound patterns
func calculatePatternSpecificity(pattern Expr) int {
	// For compound patterns (List), use compound specificity calculation
	if list, ok := pattern.(List); ok {
		cs := core.CalculateCompoundSpecificity(list)
		return cs.TotalScore
	}

	// For simple patterns, use the regular specificity (cast to int)
	return int(core.GetPatternSpecificity(pattern))
}

// RegisterPatternBuiltin registers a built-in function with a pattern from Go code
func (r *FunctionRegistry) RegisterPatternBuiltin(patternStr string, impl PatternFunc) error {
	// Parse the pattern string
	pattern, err := ParseString(patternStr)
	if err != nil {
		return fmt.Errorf("invalid pattern syntax: %v", err)
	}

	// Debug: print what was parsed and converted
	//fmt.Printf("DEBUG: Parsed pattern '%s' -> %v\n", patternStr, pattern)

	// Extract function name from original pattern (before conversion)
	functionName, err := extractFunctionName(pattern)
	if err != nil {
		return fmt.Errorf("cannot extract function name from pattern: %v", err)
	}

	// Create function definition with symbolic pattern using compound specificity
	funcDef := FunctionDef{
		Pattern:     pattern,
		Body:        nil,
		GoImpl:      impl,
		Specificity: calculatePatternSpecificity(pattern),
		IsBuiltin:   true,
	}

	// Debug the stored pattern
	// if list, ok := symbolicPattern.(List); ok && len(list.Elements) > 0 {
	//	// fmt.Printf("DEBUG: Stored pattern head type: %T\n", list.Elements[0])
	// }

	// Register the function definition
	// // fmt.Printf("DEBUG: Registering function '%s' with specificity %d\n", functionName, funcDef.Specificity)
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

	// // fmt.Printf("DEBUG: RegisterUserFunction: function=%s, pattern=%v, body=%v\\n", functionName, pattern, body)

	// Create function definition using compound specificity
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
		// fmt.Printf("DEBUG: No functions found for '%s'. Available functions: %v\n", functionName, r.GetAllFunctionNames())
		return nil, nil
	}
	// fmt.Printf("DEBUG: Found %d definitions for function '%s'\n", len(definitions), functionName)

	// Try each definition in order (most specific first)
	for _, def := range definitions {
		if matches, bindings := matchesPattern(def.Pattern, functionName, args); matches {
			return &def, bindings
		}
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
		var functionName string
		var found bool

		// Check for new Symbol type first
		if name, ok := core.ExtractSymbol(list.Elements[0]); ok {
			functionName = name
			found = true
		}

		if found {
			args := list.Elements[1:]

			// Find matching function definition
			funcDef, bindings := r.FindMatchingFunction(functionName, args)
			if funcDef == nil {
				return nil, false
			}

			// Create child context with pattern variable bindings
			funcCtx := NewChildContext(ctx)
			for varName, value := range bindings {
				funcCtx.AddScopedVar(varName) // Keep pattern variables local to this context
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
		if patternsEqual(existingDef.Pattern, newDef.Pattern) {
			// Replace existing definition
			definitions[i] = newDef
			r.functions[functionName] = definitions
			return
		}

		// Check for specificity collision with different patterns
		// Note: Disabled warnings for now as they need fine-tuning for type overlap detection
		// TODO: Re-enable after fixing type constraint extraction
		_ = couldPatternsConflict // Prevent unused function warning
		/*
			if existingDef.Specificity == newDef.Specificity && couldPatternsConflict(existingDef.Pattern, newDef.Pattern) {
				fmt.Printf("WARNING: Pattern specificity collision for function '%s'!\n", functionName)
				fmt.Printf("  Existing: %s (specificity: %d)\n", existingDef.Pattern.String(), existingDef.Specificity)
				fmt.Printf("  New:      %s (specificity: %d)\n", newDef.Pattern.String(), newDef.Specificity)
				fmt.Printf("  Order will be determined by lexicographic tie-breaker: '%s' vs '%s'\n",
					existingDef.Pattern.String(), newDef.Pattern.String())
				if existingDef.Pattern.String() < newDef.Pattern.String() {
					fmt.Printf("  Result: '%s' will match first\n", existingDef.Pattern.String())
				} else {
					fmt.Printf("  Result: '%s' will match first\n", newDef.Pattern.String())
				}
				fmt.Printf("  Consider adjusting pattern specificity to make matching order explicit.\n\n")
			}
		*/
	}

	// Add new definition and re-sort by specificity
	definitions = append(definitions, newDef)
	sort.Slice(definitions, func(i, j int) bool {
		// Higher specificity comes first
		if definitions[i].Specificity != definitions[j].Specificity {
			return definitions[i].Specificity > definitions[j].Specificity
		}
		// Tie-breaker: use lexicographic order of pattern strings for stability
		// This ensures Integer patterns come before Number patterns when specificity is equal
		return definitions[i].Pattern.String() < definitions[j].Pattern.String()
	})

	r.functions[functionName] = definitions
}

// couldPatternsConflict checks if two patterns could potentially match the same arguments
// Returns true only if there's genuine ambiguity that could cause pattern matching issues
func couldPatternsConflict(pattern1, pattern2 Expr) bool {
	// Extract type constraints from both patterns
	types1 := extractTypeConstraints(pattern1)
	types2 := extractTypeConstraints(pattern2)

	// If patterns have no overlapping type constraints, they won't conflict
	return hasOverlappingTypes(types1, types2)
}

// extractTypeConstraints extracts type names from a pattern
func extractTypeConstraints(pattern Expr) []string {
	var types []string

	switch p := pattern.(type) {
	case List:
		// Process each element in the pattern
		for _, elem := range p.Elements {
			types = append(types, extractTypeConstraints(elem)...)
		}
	default:
		// Check if this is a symbolic pattern with type constraint
		if isPattern, _, blankExpr := core.IsSymbolicPattern(pattern); isPattern {
			if isBlank, _, typeExpr := core.IsSymbolicBlank(blankExpr); isBlank && typeExpr != nil {
				if typeName, ok := core.ExtractSymbol(typeExpr); ok {
					types = append(types, typeName)
				}
			}
		} else if isBlank, _, typeExpr := core.IsSymbolicBlank(pattern); isBlank && typeExpr != nil {
			if typeName, ok := core.ExtractSymbol(typeExpr); ok {
				types = append(types, typeName)
			}
		}
	}

	return types
}

// hasOverlappingTypes checks if two sets of type constraints could overlap
func hasOverlappingTypes(types1, types2 []string) bool {
	// If either pattern has no type constraints, they could potentially conflict
	if len(types1) == 0 || len(types2) == 0 {
		return true
	}

	// Check for direct matches or subtype relationships
	for _, t1 := range types1 {
		for _, t2 := range types2 {
			if typesCouldOverlap(t1, t2) {
				return true
			}
		}
	}

	return false
}

// typesCouldOverlap checks if two specific types could overlap in pattern matching
func typesCouldOverlap(type1, type2 string) bool {
	// Same type always overlaps
	if type1 == type2 {
		return true
	}

	// Number is a supertype of Integer and Real
	if (type1 == "Number" && (type2 == "Integer" || type2 == "Real")) ||
		(type2 == "Number" && (type1 == "Integer" || type1 == "Real")) {
		return true
	}

	// Different concrete types don't overlap
	// Integer vs String, Integer vs Real, String vs Real, etc.
	return false
}

// extractFunctionName extracts the function name from a pattern
func extractFunctionName(pattern Expr) (string, error) {
	switch p := pattern.(type) {

	// unclear how core.Symbol would be triggered.
	case core.Symbol:
		// This is questionable.  Indicates some other issue
		panic("Why symbol?")
		return string(p), nil
	case List:
		return p.Type(), nil
	default:
		return "", fmt.Errorf("invalid pattern type")
	}
}

// matchesPattern checks if a pattern matches the given arguments and returns variable bindings
func matchesPattern(pattern Expr, functionName string, args []Expr) (bool, map[string]Expr) {

	// TODO: for unknown reasons the original expression is chopped up into the
	// function name and args.  But now it needs to restored to a complete express
	// Since it's immutable unclear why we are copying it.

	// Create a mock function call to match against the pattern
	mockCall := List{Elements: make([]Expr, len(args)+1)}
	mockCall.Elements[0] = core.NewSymbol(functionName)
	copy(mockCall.Elements[1:], args)

	// Use the new unified pattern matching system with sequence pattern support
	matches, bindings := core.MatchWithBindings(pattern, mockCall)

	if matches {
		// Convert core.PatternBindings to map[string]Expr for compatibility
		result := make(map[string]Expr)
		for varName, value := range bindings {
			result[varName] = value
		}
		return true, result
	}
	return false, nil
}
