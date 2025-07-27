package core

import (
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

	head, ok := list.Elements[0].(Symbol)
	if !ok || head.String() != "Pattern" {
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

	switch typeName {
	case "Integer":
		_, ok := expr.(Integer)
		return ok
	case "Real":
		_, ok := expr.(Real)
		return ok
	case "Number":
		if _, ok := expr.(Integer); ok {
			return true
		}
		if _, ok := expr.(Real); ok {
			return true
		}
		return false
	case "String":
		_, ok := expr.(String)
		return ok
	case "Symbol":
		_, ok := expr.(Symbol)
		return ok
	case "List":
		_, ok := expr.(List)
		return ok
	case "Atom":
		// Use the IsAtom() method which handles both old and new types
		return expr.IsAtom()
	default:
		// Check for ObjectExpr with matching type
		if obj, ok := expr.(ObjectExpr); ok {
			return obj.TypeName == typeName
		}
	}

	return false
}

// IsBuiltinType checks if a type name is a built-in type
func IsBuiltinType(typeName string) bool {
	builtinTypes := []string{"Integer", "Real", "Number", "String", "Symbol", "List", "Atom"}
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

	// Check for string-based patterns
	if atom, ok := pattern.(Symbol); ok {
		name := atom.String()
		if IsPatternVariable(name) {
			info := ParsePatternInfo(name)
			return GetPatternVariableSpecificity(info)
		}
	}

	// Literal patterns are most specific
	return SpecificityLiteral
}

// GetBlankExprSpecificity calculates specificity for blank expressions
func GetBlankExprSpecificity(blankExpr Expr) PatternSpecificity {
	isBlank, _, typeExpr := IsSymbolicBlank(blankExpr)
	if !isBlank {
		return SpecificityLiteral
	}

	if typeExpr == nil {
		return SpecificityGeneral // Plain _
	}

	if typeAtom, ok := typeExpr.(Symbol); ok {
		typeName := typeAtom.String()
		return GetTypeSpecificity(typeName)
	}

	return SpecificityTyped
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
	if info.TypeName == "" {
		return SpecificityGeneral
	}

	return GetTypeSpecificity(info.TypeName)
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
	// For pure matching, we do simple element-by-element matching
	// (sequence patterns are more complex and typically need context)

	if len(patternList.Elements) != len(exprList.Elements) {
		return false
	}

	// Match each element
	for i, patternElem := range patternList.Elements {
		if !pm.testMatchInternal(patternElem, exprList.Elements[i]) {
			return false
		}
	}

	return true
}
