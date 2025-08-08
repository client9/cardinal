package engine

import (
	"fmt"
	"sort"

	"github.com/client9/sexpr/core"
)

// PatternFunc represents a Go function that can be called with pattern-matched arguments
// The evaluator parameter allows access to the calling evaluator for recursive evaluation
type PatternFunc func(e *Evaluator, c *Context, args []core.Expr) core.Expr

type PatternRule struct {
	PatternString string
	Function      PatternFunc
}

// FunctionDef represents a single function definition with pattern and implementation
type FunctionDef struct {
	Pattern     core.Expr   // The pattern to match (e.g., Plus(x_Integer, y_Integer))
	Body        core.Expr   // The body expression for user-defined functions (nil for Go implementations)
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
func calculatePatternSpecificity(pattern core.Expr) int {
	// For compound patterns (List), use compound specificity calculation
	if list, ok := pattern.(core.List); ok {
		cs := core.CalculateCompoundSpecificity(list)
		return cs.TotalScore
	}

	// For simple patterns, use the regular specificity (cast to int)
	return int(core.GetPatternSpecificity(pattern))
}

func sortBySpec(v []FunctionDef) {
	sort.Slice(v, func(i, j int) bool {
		// Higher specificity comes first
		if v[i].Specificity != v[j].Specificity {
			return v[i].Specificity > v[j].Specificity
		}
		// Tie-breaker: use lexicographic order of pattern strings for stability
		// This ensures Integer patterns come before Number patterns when specificity is equal
		return v[i].Pattern.String() < v[j].Pattern.String()
	})
}

// RegisterPatternBuiltins registers multiple built-in functions from a map
func (r *FunctionRegistry) RegisterPatternBuiltins(patterns []PatternRule) error {
	for _, rule := range patterns {
		if err := r.registerPatternBuiltin(rule.PatternString, rule.Function); err != nil {
			return fmt.Errorf("failed to register pattern %s: %v", rule.PatternString, err)
		}
	}

	for k, v := range r.functions {
		sortBySpec(v)
		r.functions[k] = v
	}

	return nil
}

// RegisterPatternBuiltin registers a built-in function with a pattern from Go code
func (r *FunctionRegistry) registerPatternBuiltin(patternStr string, impl PatternFunc) error {
	// Parse the pattern stringp
	// 'RReal(max_Number)' -> RReal(Pattern(max, Blank(Number)))
	pattern, err := ParseString(patternStr)
	if err != nil {
		return fmt.Errorf("invalid pattern syntax: %v", err)
	}

	functionName := pattern.Head()
	specificity := calculatePatternSpecificity(pattern)
	funcDef := FunctionDef{
		Pattern:     pattern,
		Body:        nil,
		GoImpl:      impl,
		Specificity: specificity,
		IsBuiltin:   true,
	}

	definitions := r.functions[functionName]
	definitions = append(definitions, funcDef)
	r.functions[functionName] = definitions

	return nil
}

// registerFunctionDef adds or replaces a function definition
func (r *FunctionRegistry) registerFunctionDef(functionName string, newDef FunctionDef) {
	definitions := r.functions[functionName]

	// Check if we need to replace an existing equivalent pattern
	for i, existingDef := range definitions {
		if core.PatternsEqual(existingDef.Pattern, newDef.Pattern) {
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
	sortBySpec(definitions)
	r.functions[functionName] = definitions
}

// RegisterUserFunction registers a user-defined function with pattern and body
func (r *FunctionRegistry) RegisterUserFunction(pattern core.Expr, body core.Expr) error {
	functionName := pattern.Head()

	funcDef := FunctionDef{
		Pattern:     pattern,
		Body:        body,
		GoImpl:      nil,
		Specificity: calculatePatternSpecificity(pattern),
		IsBuiltin:   false,
	}

	r.registerFunctionDef(functionName, funcDef)
	return nil
}

func (r *FunctionRegistry) FindMatchingFunction2(fn core.Expr) (*FunctionDef, map[string]core.Expr) {
	fname := fn.Head()

	definitions, exists := r.functions[fname]
	if !exists {
		return nil, nil
	}

	for _, def := range definitions {
		if matches, bindings := matchesPattern2(def.Pattern, fn); matches {
			return &def, bindings
		}
	}
	return nil, nil

}

// FindMatchingFunction finds the best matching function definition for given arguments
func (r *FunctionRegistry) FindMatchingFunction(functionName string, args []core.Expr) (*FunctionDef, map[string]core.Expr) {
	definitions, exists := r.functions[functionName]
	if !exists {
		return nil, nil
	}
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
func (r *FunctionRegistry) CallFunction(callExpr core.Expr, ctx *Context, e *Evaluator) (core.Expr, bool) {
	// Extract function name and arguments from the call expression
	list, ok := callExpr.(core.List)
	if !ok {
		return nil, false
	}
	// Check for new Symbol type first
	//functionName, ok := core.ExtractSymbol(list.HeadExpr())

	if !ok {
		return nil, false
	}
	args := list.Tail()

	// Find matching function definition
	//funcDef, bindings := r.FindMatchingFunction(functionName, args)

	funcDef, bindings := r.FindMatchingFunction2(callExpr)
	if funcDef == nil {
		return nil, false
	}

	//log.Printf("BINDINGS: %v, args: %v", bindings, args)
	// If pattern matches, substitute variables in replacement and return it
	//return core.SubstituteBindings(replacement, bindings), true

	// Call the function
	if funcDef.GoImpl != nil {
		result := funcDef.GoImpl(e, ctx, args)

		// the downstream code doesn't have access to the single expression
		// so we can add it here.
		if err, ok := core.AsError(result); ok {
			err.Arg = callExpr
			return err, true
		}

		return result, true
	}

	// TODO: substitute bindings seems to do the same thing as the code below
	// using ReplaceAllWithRules
	return core.SubstituteBindings(funcDef.Body, bindings), true
	/*
		rules := make([]core.Expr, 0, len(bindings)+1)
		for varName, value := range bindings {
			rules = append(rules, core.NewList("Rule", core.NewSymbol(varName), value))
		}
		rlist := core.NewList("List", rules...)
		mbody := core.ReplaceAllWithRules(funcDef.Body, rlist)

		return mbody, true
	*/
}

// couldPatternsConflict checks if two patterns could potentially match the same arguments
// Returns true only if there's genuine ambiguity that could cause pattern matching issues
func couldPatternsConflict(pattern1, pattern2 core.Expr) bool {
	// Extract type constraints from both patterns
	types1 := extractTypeConstraints(pattern1)
	types2 := extractTypeConstraints(pattern2)

	// If patterns have no overlapping type constraints, they won't conflict
	return hasOverlappingTypes(types1, types2)
}

// extractTypeConstraints extracts type names from a pattern
func extractTypeConstraints(pattern core.Expr) []string {
	var types []string

	switch p := pattern.(type) {
	case core.List:
		// Process each element in the pattern
		for _, elem := range p.AsSlice() {
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
func extractFunctionName(pattern core.Expr) (string, error) {
	switch p := pattern.(type) {

	// unclear how core.Symbol would be triggered.
	case core.Symbol:
		// This is questionable.  Indicates some other issue
		panic("Why symbol?")
		//return string(p), nil
	case core.List:
		return p.Head(), nil
	default:
		return "", fmt.Errorf("invalid pattern type")
	}
}

// matchesPattern checks if a pattern matches the given arguments and returns variable bindings
func matchesPattern(pattern core.Expr, functionName string, args []core.Expr) (bool, map[string]core.Expr) {

	// TODO: for unknown reasons the original expression is chopped up into the
	// function name and args.  But now it needs to restored to a complete expressions
	// Since it's immutable unclear why we are copying it.

	// Create a mock function call to match against the pattern

	part := make([]core.Expr, len(args)+1)
	part[0] = core.NewSymbol(functionName)
	copy(part[1:], args)
	mockCall := core.NewListFromExprs(part...)

	// Use the new unified pattern matching system with sequence pattern support
	matches, bindings := core.MatchWithBindings(pattern, mockCall)

	if matches {
		// Convert core.PatternBindings to map[string]Expr for compatibility
		result := make(map[string]core.Expr)
		for varName, value := range bindings {
			result[varName] = value
		}
		return true, result
	}
	return false, nil
}

// matchesPattern checks if a pattern matches the given arguments and returns variable bindings
func matchesPattern2(pattern core.Expr, fn core.Expr) (bool, map[string]core.Expr) {

	// Use the new unified pattern matching system with sequence pattern support
	matches, bindings := core.MatchWithBindings(pattern, fn)

	if matches {
		// Convert core.PatternBindings to map[string]Expr for compatibility
		result := make(map[string]core.Expr)
		for varName, value := range bindings {
			result[varName] = value
		}
		return true, result
	}
	return false, nil
}
