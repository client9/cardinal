package core

// Pattern types and enums
type PatternType int

const (
	PatternUnknown PatternType = iota
	BlankPattern
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
		return NewList("Blank")
	}
	return NewList("Blank", typeExpr)
}

// CreateBlankSequenceExpr creates a symbolic BlankSequence[] expression
func CreateBlankSequenceExpr(typeExpr Expr) Expr {
	if typeExpr == nil {
		return NewList("BlankSequence")
	}
	return NewList("BlankSequence", typeExpr)
}

// CreateBlankNullSequenceExpr creates a symbolic BlankNullSequence[] expression
func CreateBlankNullSequenceExpr(typeExpr Expr) Expr {
	if typeExpr == nil {
		return NewList("BlankNullSequence")
	}
	return NewList("BlankNullSequence", typeExpr)
}

// CreatePatternExpr creates a symbolic Pattern[name, blank] expression
func CreatePatternExpr(nameExpr, blankExpr Expr) Expr {
	return NewList("Pattern", nameExpr, blankExpr)
}

// Pattern analysis functions

// IsSymbolicBlank checks if an expression is a symbolic blank pattern
func IsSymbolicBlank(expr Expr) (bool, PatternType, Expr) {
	list, ok := expr.(List)
	if !ok {
		return false, PatternUnknown, nil
	}

	var ptype PatternType
	switch list.Head() {
	case "Blank":
		ptype = BlankPattern
	case "BlankSequence":
		ptype = BlankSequencePattern
	case "BlankNullSequence":
		ptype = BlankNullSequencePattern
	default:
		return false, PatternUnknown, nil
	}
	var typeExpr Expr
	if list.Length() > 0 {
		args := list.Tail()
		typeExpr = args[0]
	}
	return true, ptype, typeExpr

}

// IsSymbolicPattern checks if an expression is a symbolic Pattern[name, blank]
func IsSymbolicPattern(expr Expr) (bool, Expr, Expr) {
	if list, ok := expr.(List); ok && list.Length() == 2 && list.Head() == "Pattern" {
		args := list.Tail()
		return true, args[0], args[1]
	}
	return false, nil, nil
}

// GetSymbolicPatternInfo extracts pattern information from a symbolic pattern
func GetSymbolicPatternInfo(expr Expr) PatternInfo {
	info := PatternInfo{}
	blankExpr := expr
	// Check if it's a Pattern[name, blank] first
	if isPattern, nameExpr, blank := IsSymbolicPattern(expr); isPattern {
		// Extract variable name
		if nameAtom, ok := nameExpr.(Symbol); ok {
			info.VarName = nameAtom.String()
		}
		blankExpr = blank
	}

	// Check if it's a direct blank expression
	if isBlank, blankType, typeExpr := IsSymbolicBlank(blankExpr); isBlank {
		info.Type = blankType

		// Extract type constraint
		if typeExpr != nil {
			if typeAtom, ok := typeExpr.(Symbol); ok {
				info.TypeName = typeAtom.String()
			}
		}
	}

	return info
}

// Type matching functions

// MatchesType checks if an expression matches a given type name
func MatchesType(expr Expr, typeName string) bool {
	if typeName == "" {
		return true // No type constraint
	}
	// Number is a virtual type, and probably should be changed.
	if typeName == "Number" {
		_, ok := GetNumericValue(expr)
		return ok
	}
	return expr.Head() == typeName
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
	case BlankNullSequencePattern:
		// ___ - most general (can match 0 or more)
		return baseSpecificity*10 + 0

	case BlankSequencePattern:
		// __ - less general than null sequence (must match 1 or more)
		return baseSpecificity*10 + 1

	case BlankPattern:
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

	cs := CompoundSpecificity{
		ArgsCount:       int(pattern.Length()),
		ArgsSpecificity: make([]PatternSpecificity, 0),
	}

	// Calculate head specificity
	cs.HeadSpecificity = GetPatternSpecificity(pattern.HeadExpr())

	// Calculate argument specificities
	totalArgScore := 0
	for _, e := range pattern.Tail() {
		argSpec := GetPatternSpecificity(e)
		cs.ArgsSpecificity = append(cs.ArgsSpecificity, argSpec)
		totalArgScore += int(argSpec)
	}

	// Calculate total score (higher is more specific)
	cs.TotalScore = int(cs.HeadSpecificity)*1000 + cs.ArgsCount*100 + totalArgScore

	return cs
}

// patternsEqual compares two patterns for equivalence
// This ignores variable names and only compares pattern structure and types
func PatternsEqual(pattern1, pattern2 Expr) bool {
	info1 := GetSymbolicPatternInfo(pattern1)
	info2 := GetSymbolicPatternInfo(pattern2)

	// If both are patterns, compare their structure (ignoring variable names)
	if info1.Type != PatternUnknown && info2.Type != PatternUnknown {
		return info1.Type == info2.Type && info1.TypeName == info2.TypeName
	}

	// For non-patterns or when one is a pattern and one isn't, do exact comparison
	switch p1 := pattern1.(type) {
	case Integer, Real, String, Symbol:
		return pattern1.Equal(pattern2)
	case List:
		if p2, ok := pattern2.(List); ok {
			s1 := p1.AsSlice()
			s2 := p2.AsSlice()

			if len(s1) != len(s2) {
				return false
			}
			for i := range s1 {
				if !PatternsEqual(s1[i], s2[i]) {
					return false
				}
			}
			return true
		}
		return false
	default:
		return false
	}
}
