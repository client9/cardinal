package core

import (
	"log"
	"strings"
)

// Pattern types and enums
type PatternType int

const (
	BlankPattern PatternType = iota
	BlankSequencePattern
	BlankNullSequencePattern
)

// PatternInfo represents complete information about a pattern variable
type PatternInfo struct {
	Type     PatternType
	VarName  string // Variable name (empty for anonymous patterns)
	TypeName string // Type constraint (empty for no constraint)
}

// PatternSpecificity represents the specificity of a pattern for ordering
type PatternSpecificity int

const (
	SpecificityGeneral     PatternSpecificity = iota // _ (most general)
	SpecificityTyped                                 // _Integer, _String, etc.
	SpecificityBuiltinType                           // Built-in types
	SpecificityUserType                              // User-defined types
	SpecificityLiteral                               // Exact literals (most specific)
)

// CompoundSpecificity represents specificity for complex patterns
type CompoundSpecificity struct {
	HeadSpecificity PatternSpecificity
	ArgsCount       int
	ArgsSpecificity []PatternSpecificity
	TotalScore      int
}

// Pattern constructors

// CreateBlankExpr creates a symbolic Blank[] expression
func CreateBlankExpr(typeExpr Expr) Expr {
	if typeExpr == nil {
		return List{Elements: []Expr{NewSymbol("Blank")}}
	}
	return List{Elements: []Expr{NewSymbol("Blank"), typeExpr}}
}

// CreateBlankSequenceExpr creates a symbolic BlankSequence[] expression
func CreateBlankSequenceExpr(typeExpr Expr) Expr {
	if typeExpr == nil {
		return List{Elements: []Expr{NewSymbol("BlankSequence")}}
	}
	return List{Elements: []Expr{NewSymbol("BlankSequence"), typeExpr}}
}

// CreateBlankNullSequenceExpr creates a symbolic BlankNullSequence[] expression
func CreateBlankNullSequenceExpr(typeExpr Expr) Expr {
	if typeExpr == nil {
		return List{Elements: []Expr{NewSymbol("BlankNullSequence")}}
	}
	return List{Elements: []Expr{NewSymbol("BlankNullSequence"), typeExpr}}
}

// CreatePatternExpr creates a symbolic Pattern[name, blank] expression
func CreatePatternExpr(nameExpr, blankExpr Expr) Expr {
	return List{Elements: []Expr{NewSymbol("Pattern"), nameExpr, blankExpr}}
}

// Pattern analysis functions

// IsSymbolicBlank checks if an expression is a symbolic blank pattern
func IsSymbolicBlank(expr Expr) (bool, string, Expr) {
	list, ok := expr.(List)
	if !ok || len(list.Elements) < 1 {
		return false, "", nil
	}

	head, ok := list.Elements[0].(Symbol)
	if !ok {
		return false, "", nil
	}

	headName := head.String()
	blankTypes := []string{"Blank", "BlankSequence", "BlankNullSequence"}
	for _, bt := range blankTypes {
		if headName == bt {
			var typeExpr Expr
			if len(list.Elements) > 1 {
				typeExpr = list.Elements[1]
			}
			return true, headName, typeExpr
		}
	}

	return false, "", nil
}

// IsSymbolicPattern checks if an expression is a symbolic Pattern[name, blank]
func IsSymbolicPattern(expr Expr) (bool, Expr, Expr) {
	list, ok := expr.(List)
	if !ok || len(list.Elements) != 3 {
		return false, nil, nil
	}

	if headName, ok := ExtractSymbol(list.Elements[0]); !ok || headName != "Pattern" {
		return false, nil, nil
	}

	return true, list.Elements[1], list.Elements[2]
}

// GetSymbolicPatternInfo extracts pattern information from a symbolic pattern
func GetSymbolicPatternInfo(expr Expr) PatternInfo {
	// Check if it's a Pattern[name, blank] first
	if isPattern, nameExpr, blankExpr := IsSymbolicPattern(expr); isPattern {
		info := PatternInfo{}

		// Extract variable name
		if nameAtom, ok := nameExpr.(Symbol); ok {
			info.VarName = nameAtom.String()
		}

		// Analyze the blank expression
		if isBlank, blankType, typeExpr := IsSymbolicBlank(blankExpr); isBlank {
			switch blankType {
			case "Blank":
				info.Type = BlankPattern
			case "BlankSequence":
				info.Type = BlankSequencePattern
			case "BlankNullSequence":
				info.Type = BlankNullSequencePattern
			}

			// Extract type constraint
			if typeExpr != nil {
				if typeAtom, ok := typeExpr.(Symbol); ok {
					info.TypeName = typeAtom.String()
				}
			}
		}

		return info
	}

	// Check if it's a direct blank expression
	if isBlank, blankType, typeExpr := IsSymbolicBlank(expr); isBlank {
		info := PatternInfo{}

		switch blankType {
		case "Blank":
			info.Type = BlankPattern
		case "BlankSequence":
			info.Type = BlankSequencePattern
		case "BlankNullSequence":
			info.Type = BlankNullSequencePattern
		}

		// Extract type constraint
		if typeExpr != nil {
			if typeAtom, ok := typeExpr.(Symbol); ok {
				info.TypeName = typeAtom.String()
			}
		}

		return info
	}

	return PatternInfo{}
}

// Pattern parsing functions

// IsPatternVariable checks if a string represents a pattern variable
func IsPatternVariable(name string) bool {
	return strings.Contains(name, "_")
}

// ParsePatternVariable extracts variable name and type from pattern string
func ParsePatternVariable(name string) (varName string, typeName string) {
	if !strings.Contains(name, "_") {
		return "", ""
	}

	parts := strings.Split(name, "_")
	if len(parts) == 2 {
		varName = parts[0]
		typeName = parts[1]
		if varName == "" {
			varName = "" // Anonymous pattern
		}
		if typeName == "" {
			typeName = "" // No type constraint
		}
	}
	return
}

// ConvertPatternStringToSymbolic converts string-based patterns to symbolic expressions
func ConvertPatternStringToSymbolic(name string) Expr {
	if !IsPatternVariable(name) {
		return NewSymbol(name)
	}

	// Count underscores to determine pattern type
	underscoreCount := strings.Count(name, "_")
	if underscoreCount > 3 {
		return NewSymbol(name) // Invalid pattern
	}

	// Extract variable name and type
	var varName, typeName string

	// Handle different underscore patterns
	if strings.HasSuffix(name, "___") {
		// BlankNullSequence pattern
		prefix := strings.TrimSuffix(name, "___")
		parts := strings.Split(prefix, "_")
		if len(parts) == 2 {
			varName = parts[0]
			typeName = parts[1]
		} else if len(parts) == 1 && parts[0] != "" {
			varName = parts[0]
		}
		underscoreCount = 3
	} else if strings.HasSuffix(name, "__") {
		// BlankSequence pattern
		prefix := strings.TrimSuffix(name, "__")
		parts := strings.Split(prefix, "_")
		if len(parts) == 2 {
			varName = parts[0]
			typeName = parts[1]
		} else if len(parts) == 1 && parts[0] != "" {
			varName = parts[0]
		}
		underscoreCount = 2
	} else {
		// Single blank pattern
		varName, typeName = ParsePatternVariable(name)
		underscoreCount = 1
	}

	// Create type expression if present
	var typeExpr Expr
	if typeName != "" {
		typeExpr = NewSymbol(typeName)
	}

	// Create appropriate blank expression
	var blankExpr Expr
	switch underscoreCount {
	case 1:
		blankExpr = CreateBlankExpr(typeExpr)
	case 2:
		blankExpr = CreateBlankSequenceExpr(typeExpr)
	case 3:
		blankExpr = CreateBlankNullSequenceExpr(typeExpr)
	default:
		return NewSymbol(name) // Invalid pattern, return as symbol
	}

	// If there's a variable name, wrap in Pattern[name, blank]
	if varName != "" {
		return CreatePatternExpr(NewSymbol(varName), blankExpr)
	}

	// Anonymous pattern, just return the blank expression
	return blankExpr
}

// ParsePatternInfo parses a pattern variable name and returns complete pattern information
func ParsePatternInfo(name string) PatternInfo {
	if !IsPatternVariable(name) {
		return PatternInfo{}
	}

	// Convert to symbolic and extract info
	symbolic := ConvertPatternStringToSymbolic(name)
	return GetSymbolicPatternInfo(symbolic)
}

// ConvertToSymbolicPattern converts a pattern to symbolic representation if it's a string-based pattern
func ConvertToSymbolicPattern(pattern Expr) Expr {
	switch p := pattern.(type) {
	case Symbol:
		patternStr := p.String()
		// Check if it's a pattern variable
		if IsPatternVariable(patternStr) {
			return ConvertPatternStringToSymbolic(patternStr)
		}
		return pattern
	case List:
		// Check if this is already a symbolic Pattern - don't convert its elements
		if isPattern, _, _ := IsSymbolicPattern(p); isPattern {
			return pattern
		}

		// Check if this is already a symbolic Blank - don't convert its elements
		if isBlank, _, _ := IsSymbolicBlank(p); isBlank {
			return pattern
		}

		// Convert elements recursively
		newElements := make([]Expr, len(p.Elements))
		for i, elem := range p.Elements {
			newElements[i] = ConvertToSymbolicPattern(elem)
		}
		return List{Elements: newElements}
	default:
		return pattern
	}
}

// Type matching functions

// MatchesType checks if an expression matches a given type name
func MatchesType(expr Expr, typeName string) bool {
	if typeName == "" {
		return true // No type constraint
	}
	if true {
		// Number is a virtual type, and probably should be changed.
		if typeName == "Number" {
			_, ok := GetNumericValue(expr)
			return ok
		}
		return expr.Head() == typeName
	} else {
		log.Printf("Checking typename=%s with expr=%s", typeName, expr.Head())
		switch typeName {
		case "Association":
			log.Printf("typeName = %s, expr = %s", typeName, expr.Head())
			return expr.Head() == typeName
		case "Integer":

			_, ok := expr.(Integer)
			return ok

		case "Real":

			_, ok := expr.(Real)
			return ok
		case "Number":
			// Number is a virtual type.. can be Real or Integer
			_, ok := GetNumericValue(expr)
			return ok
		case "String":

			_, ok := expr.(String)
			return ok

		case "Symbol":

			_, ok := expr.(Symbol)
			return ok

		case "List":

			_, ok := expr.(List)
			return ok
		case "Rule":
			log.Printf("IN RULE")
			// Check if it's a List with "Rule" as the head
			if list, ok := expr.(List); ok && len(list.Elements) >= 1 {
				if head, ok := list.Elements[0].(Symbol); ok {
					return head.String() == "Rule"
				}
			}
			return false

		default:

			// Check for ObjectExpr with matching type
			if obj, ok := expr.(ObjectExpr); ok {
				return obj.TypeName == typeName
			}
		}

		return false
	}
}

// IsBuiltinType checks if a type name is a built-in type
// This is used to give higher weighting to rules that use builtin types
func IsBuiltinType(typeName string) bool {
	builtinTypes := []string{"Integer", "Real", "Number", "String", "Symbol", "List", "Association", "ByteArray", "Rule"}
	for _, bt := range builtinTypes {
		if bt == typeName {
			return true
		}
	}
	return false
}

// Pattern specificity functions

// GetPatternSpecificity calculates the specificity of a pattern for ordering
func GetPatternSpecificity(pattern Expr) PatternSpecificity {
	// Check if it's a symbolic pattern
	if isPattern, _, blankExpr := IsSymbolicPattern(pattern); isPattern {
		return GetBlankExprSpecificity(blankExpr)
	}

	// Check if it's a direct blank
	if isBlank, _, _ := IsSymbolicBlank(pattern); isBlank {
		return GetBlankExprSpecificity(pattern)
	}

	// Check for compound patterns (Lists like Plus(), Times(x__Integer), etc.)
	if list, ok := pattern.(List); ok {
		cs := CalculateCompoundSpecificity(list)
		return PatternSpecificity(cs.TotalScore)
	}

	// Check for string-based patterns
	if atom, ok := pattern.(Symbol); ok {
		name := atom.String()
		if IsPatternVariable(name) {
			info := ParsePatternInfo(name)
			return GetPatternVariableSpecificity(info)
		}
	}

	// Literal patterns are most specific - boost them above all pattern types
	// Use a high multiplier to ensure they're always highest
	return SpecificityLiteral * 100
}

// GetBlankExprSpecificity calculates specificity for blank expressions
func GetBlankExprSpecificity(blankExpr Expr) PatternSpecificity {
	isBlank, blankType, typeExpr := IsSymbolicBlank(blankExpr)
	if !isBlank {
		return SpecificityLiteral
	}

	// Base specificity from type constraint
	var baseSpecificity PatternSpecificity
	if typeExpr == nil {
		baseSpecificity = SpecificityGeneral
	} else if typeAtom, ok := typeExpr.(Symbol); ok {
		typeName := typeAtom.String()
		baseSpecificity = GetTypeSpecificity(typeName)
	} else {
		baseSpecificity = SpecificityTyped
	}

	// Adjust specificity based on blank type (sequence vs single)
	// Use the same multiplier approach as GetPatternVariableSpecificity
	switch blankType {
	case "BlankNullSequence":
		// ___ - most general (can match 0 or more)
		return baseSpecificity*10 + 0

	case "BlankSequence":
		// __ - less general than null sequence (must match 1 or more)
		return baseSpecificity*10 + 1

	case "Blank":
		// _ - single patterns are more specific than sequences
		return baseSpecificity*10 + 2

	default:
		// Fallback for unknown blank types - treat as single pattern
		return baseSpecificity*10 + 2
	}
}

// GetTypeSpecificity calculates specificity for type constraints
func GetTypeSpecificity(typeName string) PatternSpecificity {
	if typeName == "" {
		return SpecificityGeneral
	}

	if IsBuiltinType(typeName) {
		return SpecificityBuiltinType
	}

	return SpecificityUserType
}

// GetPatternVariableSpecificity calculates specificity for pattern variables
func GetPatternVariableSpecificity(info PatternInfo) PatternSpecificity {
	// Base specificity from type constraint
	var baseSpecificity PatternSpecificity
	if info.TypeName == "" {
		baseSpecificity = SpecificityGeneral
	} else {
		baseSpecificity = GetTypeSpecificity(info.TypeName)
	}

	// Adjust specificity based on pattern type (sequence vs single)
	// Use a large multiplier to create clear separation between pattern types
	switch info.Type {
	case BlankNullSequencePattern:
		// x___ - most general (can match 0 or more)
		// Multiply by 10 and add 0 (lowest)
		return baseSpecificity*10 + 0

	case BlankSequencePattern:
		// x__ - less general than null sequence (must match 1 or more)
		// Multiply by 10 and add 1
		return baseSpecificity*10 + 1

	case BlankPattern:
		// x_ - single patterns are more specific than sequences
		// Multiply by 10 and add 2 (highest for each type)
		return baseSpecificity*10 + 2

	default:
		// Fallback for unknown pattern types - treat as single pattern
		return baseSpecificity*10 + 2
	}
}

// CalculateCompoundSpecificity calculates specificity for compound patterns (lists)
func CalculateCompoundSpecificity(pattern List) CompoundSpecificity {
	if len(pattern.Elements) == 0 {
		return CompoundSpecificity{}
	}

	cs := CompoundSpecificity{
		ArgsCount:       len(pattern.Elements) - 1, // Exclude head
		ArgsSpecificity: make([]PatternSpecificity, 0),
	}

	// Calculate head specificity
	cs.HeadSpecificity = GetPatternSpecificity(pattern.Elements[0])

	// Calculate argument specificities
	totalArgScore := 0
	for i := 1; i < len(pattern.Elements); i++ {
		argSpec := GetPatternSpecificity(pattern.Elements[i])
		cs.ArgsSpecificity = append(cs.ArgsSpecificity, argSpec)
		totalArgScore += int(argSpec)
	}

	// Calculate total score (higher is more specific)
	cs.TotalScore = int(cs.HeadSpecificity)*1000 + cs.ArgsCount*100 + totalArgScore

	return cs
}

// Pure pattern matching (no variable binding)

// PatternMatcher provides pure pattern matching without side effects
type PatternMatcher struct{}

// NewPatternMatcher creates a new PatternMatcher
func NewPatternMatcher() *PatternMatcher {
	return &PatternMatcher{}
}

// TestMatch tests if an expression matches a pattern (pure function, no binding)
func (pm *PatternMatcher) TestMatch(pattern, expr Expr) bool {
	// Convert pattern to symbolic if needed
	symbolicPattern := ConvertToSymbolicPattern(pattern)
	return pm.testMatchInternal(symbolicPattern, expr)
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

	// Handle different expression types
	switch p := pattern.(type) {
	case Symbol:
		varName := p.String()
		// Check if it's a pattern variable (legacy)
		if IsPatternVariable(varName) {
			info := ParsePatternInfo(varName)
			if info.Type == BlankPattern {
				// Check type constraint if present
				return MatchesType(expr, info.TypeName)
			}
			// Sequence patterns don't match single expressions in pure matching
			return false
		}

		// Regular symbol - must match literally
		if exprAtom, ok := expr.(Symbol); ok {
			return exprAtom.String() == varName
		}
		return false
	case List:
		exprList, ok := expr.(List)
		if !ok {
			return false
		}

		return pm.testMatchList(p, exprList)

	default:
		// All other types must match exactly
		return pattern.Equal(expr)
	}
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
	// (sequence handling is more complex and typically needs context)
	switch blankType {
	case "Blank", "BlankSequence", "BlankNullSequence":
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
	// Convert pattern to symbolic if needed
	symbolicPattern := ConvertToSymbolicPattern(pattern)

	// Handle symbolic patterns with variable binding
	if isPattern, varName, blankExpr := IsSymbolicPattern(symbolicPattern); isPattern {
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
	if isBlank, _, _ := IsSymbolicBlank(symbolicPattern); isBlank {
		return matchBlankWithBindings(symbolicPattern, expr, bindings)
	}

	// Handle different expression types
	switch p := symbolicPattern.(type) {
	case Symbol:
		// Direct symbol comparison
		if exprSym, ok := expr.(Symbol); ok {
			return p == exprSym
		}
		return false

	case Integer:
		if exprInt, ok := expr.(Integer); ok {
			return p == exprInt
		}
		return false

	case Real:
		if exprReal, ok := expr.(Real); ok {
			return p == exprReal
		}
		return false

	case String:
		if exprStr, ok := expr.(String); ok {
			return p == exprStr
		}
		return false

	case List:
		if exprList, ok := expr.(List); ok {
			return matchListWithBindings(p, exprList, bindings)
		}
		return false

	default:
		// For other types, just check equality
		return symbolicPattern.Equal(expr)
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
	// If we've processed all pattern elements
	if patternIdx >= len(patternList.Elements) {
		// Success if we've also processed all expression elements
		return exprIdx >= len(exprList.Elements)
	}

	// If we've run out of expression elements but still have patterns
	if exprIdx >= len(exprList.Elements) {
		// Check if remaining patterns are all BlankNullSequence (which can match zero elements)
		for i := patternIdx; i < len(patternList.Elements); i++ {
			elem := patternList.Elements[i]
			if !isNullSequencePattern(elem) {
				return false
			}
			// Bind null sequence patterns to empty list
			bindNullSequencePattern(elem, bindings)
		}
		return true
	}

	patternElem := patternList.Elements[patternIdx]

	// Check if this is a sequence pattern
	if isSequencePattern, varName, typeName, allowZero := analyzeSequencePattern(patternElem); isSequencePattern {
		return matchSequencePatternWithBindings(patternList, exprList, bindings, patternIdx, exprIdx, varName, typeName, allowZero)
	}

	// Regular pattern - match one element
	if matchWithBindingsInternal(patternElem, exprList.Elements[exprIdx], bindings) {
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
			case "BlankSequence":
				return true, vn, tn, false
			case "BlankNullSequence":
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
		case "BlankSequence":
			return true, "", tn, false
		case "BlankNullSequence":
			return true, "", tn, true
		}
	}

	// Check for legacy string-based sequence patterns
	if sym, ok := pattern.(Symbol); ok {
		name := sym.String()
		if IsPatternVariable(name) {
			info := ParsePatternInfo(name)
			switch info.Type {
			case BlankSequencePattern:
				return true, info.VarName, info.TypeName, false
			case BlankNullSequencePattern:
				return true, info.VarName, info.TypeName, true
			}
		}
	}

	return false, "", "", false
}

// isNullSequencePattern checks if a pattern is a BlankNullSequence that can match zero elements
func isNullSequencePattern(pattern Expr) bool {
	// Check for symbolic Pattern[name, BlankNullSequence]
	if isPattern, _, blankExpr := IsSymbolicPattern(pattern); isPattern {
		if isBlank, blankType, _ := IsSymbolicBlank(blankExpr); isBlank && blankType == "BlankNullSequence" {
			return true
		}
	}

	// Check for direct symbolic BlankNullSequence
	if isBlank, blankType, _ := IsSymbolicBlank(pattern); isBlank && blankType == "BlankNullSequence" {
		return true
	}

	// Check for legacy string-based null sequence patterns
	if sym, ok := pattern.(Symbol); ok {
		name := sym.String()
		if IsPatternVariable(name) {
			info := ParsePatternInfo(name)
			return info.Type == BlankNullSequencePattern
		}
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
	} else if sym, ok := pattern.(Symbol); ok {
		name := sym.String()
		if IsPatternVariable(name) {
			info := ParsePatternInfo(name)
			varName = info.VarName
		}
	}

	// Bind to empty list if we have a variable name
	if varName != "" {
		// Create a proper List with "List" as the first element (like the old system)
		bindings[varName] = List{Elements: []Expr{NewSymbol("List")}}
	}
}

// matchSequencePatternWithBindings handles matching sequence patterns
func matchSequencePatternWithBindings(patternList, exprList List, bindings PatternBindings, patternIdx, exprIdx int, varName, typeName string, allowZero bool) bool {
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
		// Create a copy of bindings for this attempt
		testBindings := make(PatternBindings)
		for k, v := range bindings {
			testBindings[k] = v
		}

		// Collect the elements to bind to this sequence
		var seqElements []Expr
		for i := 0; i < consume; i++ {
			if exprIdx+i < len(exprList.Elements) {
				seqElements = append(seqElements, exprList.Elements[exprIdx+i])
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
			testBindings[varName] = List{Elements: listElements}
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
			if boundList, ok := boundValue.(List); ok && len(boundList.Elements) > 0 {
				// Check if the List head is "List" (our sequence marker)
				if headSym, ok := boundList.Elements[0].(Symbol); ok {
					return headSym.String() == "List"
				}
			}
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
		newElements := make([]Expr, 0, len(e.Elements))
		changed := false

		for i, elem := range e.Elements {
			newElem := SubstituteBindings(elem, bindings)

			// Check if this is a sequence variable substitution that needs splicing
			if i > 0 && needsSequenceSplicing(elem, newElem, bindings) {
				// This is a sequence variable - splice its elements
				if elemSym, ok := elem.(Symbol); ok {
					varName := string(elemSym)
					if boundValue, exists := bindings[varName]; exists {
						if boundList, ok := boundValue.(List); ok {
							// Check if it's an empty sequence (just "List" head)
							if len(boundList.Elements) == 1 {
								// Empty sequence - skip adding anything
								changed = true
								continue
							} else if len(boundList.Elements) > 1 {
								// Non-empty sequence - skip the "List" head and add the actual elements
								newElements = append(newElements, boundList.Elements[1:]...)
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
			return List{Elements: newElements}
		}
		return e

	default:
		// For atomic types, no substitution needed
		return e
	}
}
