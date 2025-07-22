package sexpr

import (
	"fmt"
	"strings"
)

// PatternType represents the type of pattern
type PatternType int

const (
	BlankPattern             PatternType = iota // _ - matches exactly one expression
	BlankSequencePattern                        // __ - matches one or more expressions
	BlankNullSequencePattern                    // ___ - matches zero or more expressions
)

// PatternInfo contains information about a parsed pattern
type PatternInfo struct {
	Type     PatternType
	VarName  string
	TypeName string
}

// PatternSpecificity represents how specific a pattern is (higher = more specific)
type PatternSpecificity int

const (
	SpecificityNullSequence    PatternSpecificity = 10 // x___ (least specific)
	SpecificitySequence        PatternSpecificity = 20 // x__
	SpecificityGeneral         PatternSpecificity = 30 // x_
	SpecificityBuiltinGeneral  PatternSpecificity = 40 // x_Number, x_Numeric (general builtin types)
	SpecificityBuiltinSpecific PatternSpecificity = 50 // x_Integer, x_Real, x_String (specific builtin types)
	SpecificityCustomType      PatternSpecificity = 60 // x_Color, x_Point, etc.
	SpecificityLiteral         PatternSpecificity = 70 // 42, "hello", exact values (most specific)
	
	// Bonus multipliers for compound patterns
	LiteralBonusMultiplier = 2  // Literals get double weight in compounds
	SpecificTypeBonusPoints = 5 // Extra points for specific types
	PatternComplexityBonus = 1  // Per additional argument
)

// CompoundSpecificity provides detailed scoring breakdown for compound patterns
// This allows for detailed analysis of why certain patterns have higher specificity
// than others, which is crucial for debugging pattern matching issues.
//
// The scoring algorithm works as follows:
// 1. BaseScore: Sum of individual argument specificities + complexity bonus
// 2. LiteralBonus: Additional points for literal values (multiplied by LiteralBonusMultiplier)
// 3. TypeBonus: Additional points for specific type constraints (SpecificTypeBonusPoints each)
// 4. TotalScore: BaseScore + LiteralBonus + TypeBonus
//
// Examples:
// - Plus(1, 2): BaseScore=140+2, LiteralBonus=140, TypeBonus=0, TotalScore=282
// - Plus(1, y_Integer): BaseScore=70+50+2, LiteralBonus=70, TypeBonus=5, TotalScore=197
// - Plus(x_Integer, y_Integer): BaseScore=50+50+2, LiteralBonus=0, TypeBonus=10, TotalScore=112
type CompoundSpecificity struct {
	BaseScore    int // Sum of individual argument specificities
	LiteralBonus int // Extra points for literals
	TypeBonus    int // Extra points for specific types
	TotalScore   int // Final calculated score
}

// isLiteral checks if an expression is a literal value (number, string, etc.)
func isLiteral(expr Expr) bool {
	if atom, ok := expr.(Atom); ok {
		switch atom.AtomType {
		case IntAtom, FloatAtom, StringAtom:
			return true
		case SymbolAtom:
			symbolName := atom.Value.(string)
			// Built-in constants are considered literals
			return symbolName == "True" || symbolName == "False" || symbolName == "Pi" || symbolName == "E"
		}
	}
	return false
}

// isSpecificType checks if a pattern variable has a specific builtin type constraint
func isSpecificType(expr Expr) bool {
	if atom, ok := expr.(Atom); ok && atom.AtomType == SymbolAtom {
		symbolName := atom.Value.(string)
		if isPatternVariable(symbolName) {
			info := parsePatternInfo(symbolName)
			return getTypeSpecificity(info.TypeName) == SpecificityBuiltinSpecific
		}
	}
	return false
}

// calculateCompoundSpecificity calculates specificity for compound patterns
func calculateCompoundSpecificity(pattern List) CompoundSpecificity {
	if len(pattern.Elements) <= 1 {
		// Simple pattern, use basic calculation
		baseSpec := int(getPatternSpecificityBasic(pattern))
		return CompoundSpecificity{
			BaseScore:  baseSpec,
			TotalScore: baseSpec,
		}
	}
	
	args := pattern.Elements[1:] // Skip function head
	
	var baseScore, literalBonus, typeBonus int
	
	for _, arg := range args {
		argSpec := int(getPatternSpecificityBasic(arg))
		baseScore += argSpec
		
		// Bonus scoring
		if isLiteral(arg) {
			literalBonus += argSpec * LiteralBonusMultiplier / 2 // Apply multiplier bonus
		} else if isSpecificType(arg) {
			typeBonus += SpecificTypeBonusPoints
		}
		
		// Complexity bonus for each additional argument
		baseScore += PatternComplexityBonus
	}
	
	totalScore := baseScore + literalBonus + typeBonus
	
	return CompoundSpecificity{
		BaseScore:    baseScore,
		LiteralBonus: literalBonus,
		TypeBonus:    typeBonus,
		TotalScore:   totalScore,
	}
}

// DebugPatternSpecificity returns detailed information about pattern specificity calculation
// This is useful for debugging pattern matching issues and understanding why certain patterns
// are chosen over others. It returns a human-readable breakdown of the specificity calculation.
func DebugPatternSpecificity(pattern Expr) string {
	switch pat := pattern.(type) {
	case List:
		if len(pat.Elements) > 1 {
			// Check if this is a compound pattern
			if head, ok := pat.Elements[0].(Atom); ok && head.AtomType == SymbolAtom {
				compound := calculateCompoundSpecificity(pat)
				result := fmt.Sprintf("Pattern: %s\n", pattern.String())
				result += fmt.Sprintf("Type: Compound Pattern\n")
				result += fmt.Sprintf("Arguments: %d\n", len(pat.Elements)-1)
				result += fmt.Sprintf("Base Score: %d\n", compound.BaseScore)
				result += fmt.Sprintf("Literal Bonus: %d\n", compound.LiteralBonus)
				result += fmt.Sprintf("Type Bonus: %d\n", compound.TypeBonus)
				result += fmt.Sprintf("Total Score: %d\n", compound.TotalScore)
				
				// Break down individual arguments
				result += "Argument Analysis:\n"
				for i, arg := range pat.Elements[1:] {
					argSpec := getPatternSpecificityBasic(arg)
					result += fmt.Sprintf("  Arg %d (%s): %d", i+1, arg.String(), argSpec)
					if isLiteral(arg) {
						result += " [LITERAL]"
					} else if isSpecificType(arg) {
						result += " [SPECIFIC TYPE]"
					}
					result += "\n"
				}
				return result
			}
		}
		// Simple pattern
		spec := getPatternSpecificityBasic(pattern)
		return fmt.Sprintf("Pattern: %s\nType: Simple Pattern\nScore: %d\n", pattern.String(), spec)
	default:
		spec := getPatternSpecificityBasic(pattern)
		return fmt.Sprintf("Pattern: %s\nType: Basic Pattern\nScore: %d\n", pattern.String(), spec)
	}
}

// getPatternSpecificityBasic calculates specificity for a single pattern element (non-compound)
func getPatternSpecificityBasic(pattern Expr) PatternSpecificity {
	switch pat := pattern.(type) {
	case Atom:
		if pat.AtomType == SymbolAtom {
			symbolName := pat.Value.(string)
			if isPatternVariable(symbolName) {
				patternInfo := parsePatternInfo(symbolName)
				return getPatternVariableSpecificity(patternInfo)
			}
			// Regular symbol (should not occur in patterns, but handle it)
			return SpecificityLiteral
		}
		// Literal atom (number, string, boolean)
		return SpecificityLiteral
	case List:
		// Handle symbolic Pattern expressions: Pattern(x, Blank(...))
		if len(pat.Elements) >= 3 && pat.Elements[0].String() == "Pattern" {
			// This is Pattern(variable, BlankType) - analyze the blank type
			blankExpr := pat.Elements[2]
			return getBlankExprSpecificity(blankExpr)
		}
		// For lists in basic calculation, return minimal specificity
		return SpecificityGeneral
	default:
		return SpecificityGeneral
	}
}

// matchBlankExpression matches a blank expression against an expression
func (e *Evaluator) matchBlankExpression(blankExpr Expr, expr Expr, ctx *Context) bool {
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

// MatchPattern matches a pattern against an expression and binds variables (exported for testing)
func (e *Evaluator) MatchPattern(pattern Expr, expr Expr, ctx *Context) bool {
	return e.matchPatternInternal(pattern, expr, ctx, false)
}

// matchPattern matches a pattern against an expression and binds variables
func (e *Evaluator) matchPattern(pattern Expr, expr Expr, ctx *Context) bool {
	return e.matchPatternInternal(pattern, expr, ctx, false)
}

// matchPatternInternal matches a pattern with control over parameter binding behavior
func (e *Evaluator) matchPatternInternal(pattern Expr, expr Expr, ctx *Context, isParameterList bool) bool {
	// First, check if the pattern is a symbolic pattern (Pattern[], Blank[], etc.)
	if isPattern, nameExpr, blankExpr := isSymbolicPattern(pattern); isPattern {
		// This is a Pattern[name, blank] expression
		var varName string
		if nameAtom, ok := nameExpr.(Atom); ok && nameAtom.AtomType == SymbolAtom {
			varName = nameAtom.Value.(string)
		}

		// Check if the blank expression matches
		if e.matchBlankExpression(blankExpr, expr, ctx) {
			// Bind the variable if it has a name
			if varName != "" {
				ctx.Set(varName, expr)
			}
			return true
		}
		return false
	}

	// Check if the pattern is a direct blank expression
	if isBlank, _, _ := isSymbolicBlank(pattern); isBlank {
		return e.matchBlankExpression(pattern, expr, ctx)
	}

	switch pat := pattern.(type) {
	case Atom:
		if pat.AtomType == SymbolAtom {
			symbolName := pat.Value.(string)
			// Check if this is a pattern variable (ends with _) - backward compatibility
			if isPatternVariable(symbolName) {
				// Convert to symbolic and match
				symbolicPattern := ConvertPatternStringToSymbolic(symbolName)
				return e.matchPatternInternal(symbolicPattern, expr, ctx, isParameterList)
			} else {
				// Regular symbol behavior depends on context
				if isParameterList {
					// In parameter lists, regular symbols bind to values
					ctx.Set(symbolName, expr)
					return true
				} else {
					// In head patterns, regular symbols require exact matches
					if exprAtom, ok := expr.(Atom); ok && exprAtom.AtomType == SymbolAtom {
						return exprAtom.Value.(string) == symbolName
					}
					return false
				}
			}
		} else {
			// Literal atom - exact match required
			if exprAtom, ok := expr.(Atom); ok {
				return pat.AtomType == exprAtom.AtomType && pat.Value == exprAtom.Value
			}
			return false
		}
	case List:
		// Match structured expressions like Plus(x_, y_)
		if exprList, ok := expr.(List); ok {
			// For function parameters, we need to handle sequence patterns
			if isParameterList && len(pat.Elements) > 1 {
				// First, check that the heads match exactly (head is never a parameter)
				if !e.matchPatternInternal(pat.Elements[0], exprList.Elements[0], ctx, false) {
					return false
				}
				// Then use sequence pattern matching for the arguments
				return e.matchSequencePatterns(pat.Elements[1:], exprList.Elements[1:], ctx)
			}

			// For non-parameter lists, require exact length match
			if len(pat.Elements) != len(exprList.Elements) {
				return false
			}

			for i, patElem := range pat.Elements {
				// First element is the head (requires exact match), rest are parameters
				isParam := i > 0 && isParameterList
				if !e.matchPatternInternal(patElem, exprList.Elements[i], ctx, isParam) {
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

// matchSequencePatterns handles matching patterns that may contain sequence patterns
func (e *Evaluator) matchSequencePatterns(patterns []Expr, exprs []Expr, ctx *Context) bool {
	patternIndex := 0
	exprIndex := 0

	for patternIndex < len(patterns) {
		pattern := patterns[patternIndex]

		// Check if this pattern is a sequence pattern
		if atom, ok := pattern.(Atom); ok && atom.AtomType == SymbolAtom {
			symbolName := atom.Value.(string)
			if isPatternVariable(symbolName) {
				patternInfo := parsePatternInfo(symbolName)

				switch patternInfo.Type {
				case BlankPattern:
					// Single pattern: matches exactly one expression
					if exprIndex >= len(exprs) {
						return false // No more expressions to match
					}

					// Check type constraint
					if !matchesType(exprs[exprIndex], patternInfo.TypeName) {
						return false
					}

					// Bind variable if named
					if patternInfo.VarName != "" {
						ctx.Set(patternInfo.VarName, exprs[exprIndex])
					}

					exprIndex++
					patternIndex++

				case BlankSequencePattern:
					// Sequence pattern: matches one or more expressions
					if exprIndex >= len(exprs) {
						return false // Need at least one expression
					}

					// Calculate how many expressions this pattern should consume
					remainingPatterns := len(patterns) - patternIndex - 1
					minExprsNeeded := remainingPatterns // Minimum expressions needed for remaining patterns
					maxExprsAvailable := len(exprs) - exprIndex - minExprsNeeded

					if maxExprsAvailable < 1 {
						return false // Need at least one expression for BlankSequence
					}

					// Consume expressions for this sequence pattern
					sequenceExprs := make([]Expr, 0)
					for i := 0; i < maxExprsAvailable; i++ {
						if exprIndex >= len(exprs) {
							break
						}

						// Check type constraint
						if !matchesType(exprs[exprIndex], patternInfo.TypeName) {
							break
						}

						sequenceExprs = append(sequenceExprs, exprs[exprIndex])
						exprIndex++
					}

					if len(sequenceExprs) == 0 {
						return false // BlankSequence must match at least one expression
					}

					// Bind variable if named
					if patternInfo.VarName != "" {
						ctx.Set(patternInfo.VarName, NewList(append([]Expr{NewSymbolAtom("List")}, sequenceExprs...)...))
					}

					patternIndex++

				case BlankNullSequencePattern:
					// Null sequence pattern: matches zero or more expressions
					// Calculate how many expressions this pattern should consume
					remainingPatterns := len(patterns) - patternIndex - 1
					minExprsNeeded := remainingPatterns // Minimum expressions needed for remaining patterns
					maxExprsAvailable := len(exprs) - exprIndex - minExprsNeeded

					if maxExprsAvailable < 0 {
						maxExprsAvailable = 0 // Can consume zero expressions
					}

					// Consume expressions for this sequence pattern
					sequenceExprs := make([]Expr, 0)
					for i := 0; i < maxExprsAvailable; i++ {
						if exprIndex >= len(exprs) {
							break
						}

						// Check type constraint
						if !matchesType(exprs[exprIndex], patternInfo.TypeName) {
							break
						}

						sequenceExprs = append(sequenceExprs, exprs[exprIndex])
						exprIndex++
					}

					// Bind variable if named (can be empty list)
					if patternInfo.VarName != "" {
						ctx.Set(patternInfo.VarName, NewList(append([]Expr{NewSymbolAtom("List")}, sequenceExprs...)...))
					}

					patternIndex++

				default:
					return false // Unknown pattern type
				}
			} else {
				// Regular symbol pattern
				if exprIndex >= len(exprs) {
					return false
				}

				// Regular symbol in parameter list binds to value
				ctx.Set(symbolName, exprs[exprIndex])
				exprIndex++
				patternIndex++
			}
		} else {
			// Non-symbol pattern (literal)
			if exprIndex >= len(exprs) {
				return false
			}

			if !e.matchPatternInternal(pattern, exprs[exprIndex], ctx, false) {
				return false
			}

			exprIndex++
			patternIndex++
		}
	}

	// All patterns matched, check that all expressions were consumed
	return exprIndex == len(exprs)
}

// getPatternSpecificity calculates the specificity of a single pattern
func getPatternSpecificity(pattern Expr) PatternSpecificity {
	switch pat := pattern.(type) {
	case Atom:
		return getPatternSpecificityBasic(pattern)
	case List:
		// Handle symbolic Pattern expressions: Pattern(x, Blank(...))
		if len(pat.Elements) >= 3 && pat.Elements[0].String() == "Pattern" {
			// This is Pattern(variable, BlankType) - analyze the blank type
			blankExpr := pat.Elements[2]
			return getBlankExprSpecificity(blankExpr)
		}
		
		// Compound pattern like Plus(x_, y_) - use compound specificity calculation
		if len(pat.Elements) > 1 {
			// Check if head is a function name (symbol literal)
			if head, ok := pat.Elements[0].(Atom); ok && head.AtomType == SymbolAtom {
				// This is a compound pattern - use compound specificity
				compound := calculateCompoundSpecificity(pat)
				return PatternSpecificity(compound.TotalScore)
			}
		}
		
		// Simple list pattern
		return getPatternSpecificityBasic(pattern)
	default:
		return SpecificityGeneral
	}
}

// getBlankExprSpecificity calculates specificity for symbolic Blank expressions
func getBlankExprSpecificity(blankExpr Expr) PatternSpecificity {
	switch blank := blankExpr.(type) {
	case List:
		if len(blank.Elements) == 0 {
			return SpecificityGeneral // Blank()
		}

		head := blank.Elements[0]
		switch head.String() {
		case "Blank":
			if len(blank.Elements) == 1 {
				return SpecificityGeneral // Blank() - matches anything
			} else if len(blank.Elements) == 2 {
				// Blank(Type) - check the type
				typeExpr := blank.Elements[1]
				return getTypeSpecificity(typeExpr.String())
			}
		case "BlankSequence":
			if len(blank.Elements) == 1 {
				return SpecificitySequence // BlankSequence() - matches sequence
			} else if len(blank.Elements) == 2 {
				// BlankSequence(Type) - check the type but still sequence
				return SpecificitySequence // Typed sequence is still sequence level
			}
		case "BlankNullSequence":
			return SpecificityNullSequence // BlankNullSequence() - least specific
		}
	case Atom:
		if blank.AtomType == SymbolAtom {
			symbol := blank.Value.(string)
			switch symbol {
			case "Blank":
				return SpecificityGeneral
			case "BlankSequence":
				return SpecificitySequence
			case "BlankNullSequence":
				return SpecificityNullSequence
			}
		}
	}
	return SpecificityGeneral
}

// getTypeSpecificity determines specificity based on type name
func getTypeSpecificity(typeName string) PatternSpecificity {
	// Specific builtin types (more specific)
	specificTypes := []string{"Integer", "Real", "Float", "String", "Boolean", "Bool", "Symbol", "Atom", "List"}
	for _, specific := range specificTypes {
		if typeName == specific {
			return SpecificityBuiltinSpecific
		}
	}

	// General builtin types (less specific)
	generalTypes := []string{"Number", "Numeric"}
	for _, general := range generalTypes {
		if typeName == general {
			return SpecificityBuiltinGeneral
		}
	}

	// Custom types (like Uint64, Color, Point, etc.)
	return SpecificityCustomType
}

// getPatternVariableSpecificity calculates specificity for pattern variables
func getPatternVariableSpecificity(info PatternInfo) PatternSpecificity {
	// Base specificity on pattern type
	var baseSpecificity PatternSpecificity
	switch info.Type {
	case BlankNullSequencePattern:
		baseSpecificity = SpecificityNullSequence
	case BlankSequencePattern:
		baseSpecificity = SpecificitySequence
	case BlankPattern:
		baseSpecificity = SpecificityGeneral
	default:
		baseSpecificity = SpecificityGeneral
	}

	// Increase specificity if there's a type constraint
	if info.TypeName != "" {
		return getTypeSpecificity(info.TypeName)
	}

	return baseSpecificity
}

// isBuiltinType checks if a type name is a built-in type
func isBuiltinType(typeName string) bool {
	switch typeName {
	case "Integer", "Real", "Float", "Number", "Numeric", "String", "Boolean", "Bool", "Symbol", "Atom", "List":
		return true
	default:
		return false
	}
}

// isPatternVariable checks if a symbol name represents a pattern variable
func isPatternVariable(name string) bool {
	// Pattern variables have the form: [varname][underscores][typename]
	// Examples: x_, x__, x___, _Integer, x_Integer, x__Integer, x___Integer
	// But NOT regular symbols with underscores in the middle: bool_test, my_function, etc.

	// Must start with underscore OR contain underscore followed by uppercase letter (type)
	if strings.HasPrefix(name, "_") {
		return true // _Integer, _, __, ___, etc.
	}

	// Look for pattern: letter(s) + underscore(s) + [optional type starting with uppercase]
	for i := 0; i < len(name); i++ {
		if name[i] == '_' {
			// Found underscore - check if it's followed by more underscores or uppercase letter or end of string
			remaining := name[i:]
			if remaining == "_" || remaining == "__" || remaining == "___" {
				return true // x_, x__, x___
			}
			// Check if it's underscore(s) followed by type name (starts with uppercase)
			if len(remaining) > 1 {
				// Skip consecutive underscores
				j := 1
				for j < len(remaining) && remaining[j] == '_' {
					j++
				}
				if j < len(remaining) && remaining[j] >= 'A' && remaining[j] <= 'Z' {
					return true // x_Integer, x__Integer, x___Integer
				}
			}
			// If we found underscore but it doesn't match pattern variable format, it's not a pattern variable
			return false
		}
	}

	return false
}

// parsePatternVariable parses a pattern variable name and returns the variable name and type constraint
func parsePatternVariable(name string) (varName string, typeName string) {
	if !isPatternVariable(name) {
		return "", ""
	}

	// Split by underscore
	parts := strings.Split(name, "_")

	if len(parts) == 2 {
		if parts[0] == "" && parts[1] == "" {
			// Anonymous pattern: _ -> varName="", typeName=""
			return "", ""
		} else if parts[1] == "" {
			// Simple pattern: x_ -> varName="x", typeName=""
			return parts[0], ""
		} else if parts[0] != "" {
			// Named pattern: x_Integer -> varName="x", typeName="Integer"
			return parts[0], parts[1]
		}
	}

	// Invalid pattern
	return "", ""
}

// ConvertPatternStringToSymbolic converts a pattern string (like "x_Integer") to a symbolic expression
func ConvertPatternStringToSymbolic(name string) Expr {
	if !isPatternVariable(name) {
		// Not a pattern variable, return as regular symbol
		return NewSymbolAtom(name)
	}

	// Count consecutive underscores to determine pattern type
	underscoreCount := 0
	underscoreStart := -1

	// Find the first underscore
	for i, ch := range name {
		if ch == '_' {
			underscoreStart = i
			break
		}
	}

	if underscoreStart == -1 {
		return NewSymbolAtom(name) // No underscore found, regular symbol
	}

	// Count consecutive underscores
	for i := underscoreStart; i < len(name) && name[i] == '_'; i++ {
		underscoreCount++
	}

	// Extract variable name (before underscores)
	varName := name[:underscoreStart]

	// Extract type name (after underscores)
	typeStart := underscoreStart + underscoreCount
	var typeName string
	if typeStart < len(name) {
		typeName = name[typeStart:]
	}

	// Create the appropriate Blank expression
	var blankExpr Expr
	var typeExpr Expr
	if typeName != "" {
		typeExpr = NewSymbolAtom(typeName)
	}

	switch underscoreCount {
	case 1:
		blankExpr = CreateBlankExpr(typeExpr)
	case 2:
		blankExpr = CreateBlankSequenceExpr(typeExpr)
	case 3:
		blankExpr = CreateBlankNullSequenceExpr(typeExpr)
	default:
		return NewSymbolAtom(name) // Invalid pattern, return as symbol
	}

	// If there's a variable name, wrap in Pattern[name, blank]
	if varName != "" {
		return CreatePatternExpr(NewSymbolAtom(varName), blankExpr)
	}

	// Anonymous pattern, just return the blank expression
	return blankExpr
}

// parsePatternInfo parses a pattern variable name and returns complete pattern information
func parsePatternInfo(name string) PatternInfo {
	if !isPatternVariable(name) {
		return PatternInfo{}
	}

	// Count consecutive underscores to determine pattern type
	underscoreCount := 0
	underscoreStart := -1

	// Find the first underscore
	for i, ch := range name {
		if ch == '_' {
			underscoreStart = i
			break
		}
	}

	if underscoreStart == -1 {
		return PatternInfo{} // No underscore found
	}

	// Count consecutive underscores
	for i := underscoreStart; i < len(name) && name[i] == '_'; i++ {
		underscoreCount++
	}

	// Determine pattern type based on underscore count
	var patternType PatternType
	switch underscoreCount {
	case 1:
		patternType = BlankPattern
	case 2:
		patternType = BlankSequencePattern
	case 3:
		patternType = BlankNullSequencePattern
	default:
		return PatternInfo{} // Invalid pattern
	}

	// Extract variable name (before underscores)
	varName := name[:underscoreStart]

	// Extract type name (after underscores)
	typeStart := underscoreStart + underscoreCount
	var typeName string
	if typeStart < len(name) {
		typeName = name[typeStart:]
	}

	return PatternInfo{
		Type:     patternType,
		VarName:  varName,
		TypeName: typeName,
	}
}

// isSymbolicBlank checks if an expression is a symbolic Blank[], BlankSequence[], or BlankNullSequence[]
func isSymbolicBlank(expr Expr) (bool, string, Expr) {
	if list, ok := expr.(List); ok && len(list.Elements) >= 1 {
		if head, ok := list.Elements[0].(Atom); ok && head.AtomType == SymbolAtom {
			headName := head.Value.(string)
			switch headName {
			case "Blank", "BlankSequence", "BlankNullSequence":
				var typeExpr Expr
				if len(list.Elements) > 1 {
					typeExpr = list.Elements[1]
				}
				return true, headName, typeExpr
			}
		}
	}
	return false, "", nil
}

// isSymbolicPattern checks if an expression is a Pattern[name, blank]
func isSymbolicPattern(expr Expr) (bool, Expr, Expr) {
	if list, ok := expr.(List); ok && len(list.Elements) == 3 {
		if head, ok := list.Elements[0].(Atom); ok && head.AtomType == SymbolAtom {
			if head.Value.(string) == "Pattern" {
				return true, list.Elements[1], list.Elements[2] // name, blank
			}
		}
	}
	return false, nil, nil
}

// getSymbolicPatternInfo extracts pattern information from a symbolic pattern expression
func getSymbolicPatternInfo(expr Expr) PatternInfo {
	// Check if it's a Pattern[name, blank]
	if isPattern, nameExpr, blankExpr := isSymbolicPattern(expr); isPattern {
		var varName string
		if nameAtom, ok := nameExpr.(Atom); ok && nameAtom.AtomType == SymbolAtom {
			varName = nameAtom.Value.(string)
		}

		if isBlank, blankType, typeExpr := isSymbolicBlank(blankExpr); isBlank {
			var typeName string
			if typeExpr != nil {
				if typeAtom, ok := typeExpr.(Atom); ok && typeAtom.AtomType == SymbolAtom {
					typeName = typeAtom.Value.(string)
				}
			}

			var patternType PatternType
			switch blankType {
			case "Blank":
				patternType = BlankPattern
			case "BlankSequence":
				patternType = BlankSequencePattern
			case "BlankNullSequence":
				patternType = BlankNullSequencePattern
			}

			return PatternInfo{Type: patternType, VarName: varName, TypeName: typeName}
		}
	}

	// Check if it's a direct blank expression
	if isBlank, blankType, typeExpr := isSymbolicBlank(expr); isBlank {
		var typeName string
		if typeExpr != nil {
			if typeAtom, ok := typeExpr.(Atom); ok && typeAtom.AtomType == SymbolAtom {
				typeName = typeAtom.Value.(string)
			}
		}

		var patternType PatternType
		switch blankType {
		case "Blank":
			patternType = BlankPattern
		case "BlankSequence":
			patternType = BlankSequencePattern
		case "BlankNullSequence":
			patternType = BlankNullSequencePattern
		}

		return PatternInfo{Type: patternType, VarName: "", TypeName: typeName}
	}

	// Not a pattern
	return PatternInfo{}
}

// matchesType checks if an expression matches a given type constraint
func matchesType(expr Expr, typeName string) bool {
	if typeName == "" {
		// No type constraint, matches anything
		return true
	}

	// Handle built-in types first
	switch typeName {
	case "Integer":
		if atom, ok := expr.(Atom); ok {
			return atom.AtomType == IntAtom
		}
	case "Real", "Float":
		if atom, ok := expr.(Atom); ok {
			return atom.AtomType == FloatAtom
		}
	case "Number", "Numeric":
		if atom, ok := expr.(Atom); ok {
			return atom.AtomType == IntAtom || atom.AtomType == FloatAtom
		}
	case "String":
		if atom, ok := expr.(Atom); ok {
			return atom.AtomType == StringAtom
		}
	case "Boolean", "Bool":
		return isBool(expr)
	case "Symbol":
		if atom, ok := expr.(Atom); ok {
			return atom.AtomType == SymbolAtom
		}
	case "Atom":
		_, ok := expr.(Atom)
		return ok
	case "List":
		_, ok := expr.(List)
		return ok
	default:
		// Handle ObjectExpr with custom TypeName
		if objExpr, ok := expr.(ObjectExpr); ok {
			return objExpr.TypeName == typeName
		}

		// Handle arbitrary head symbols (e.g., Color, Point, etc.)
		if list, ok := expr.(List); ok && len(list.Elements) > 0 {
			if headAtom, ok := list.Elements[0].(Atom); ok && headAtom.AtomType == SymbolAtom {
				return headAtom.Value.(string) == typeName
			}
		}
		return false
	}

	return false
}
