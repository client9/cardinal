package engine

import (
	"github.com/client9/sexpr/core"
)

// PatternExecutor handles pattern matching with variable binding and Context
type PatternExecutor struct {
	matcher *core.PatternMatcher
}

// NewPatternExecutor creates a new PatternExecutor
func NewPatternExecutor() *PatternExecutor {
	return &PatternExecutor{
		matcher: core.NewPatternMatcher(),
	}
}

// MatchWithBinding performs pattern matching with variable binding in the provided Context
func (pe *PatternExecutor) MatchWithBinding(pattern, expr core.Expr, ctx *Context) bool {
	// Convert pattern to symbolic if needed
	symbolicPattern := core.ConvertToSymbolicPattern(pattern)
	return pe.matchWithBindingInternal(symbolicPattern, expr, ctx, false)
}

// matchWithBindingInternal implements the core matching logic with binding
func (pe *PatternExecutor) matchWithBindingInternal(pattern, expr core.Expr, ctx *Context, isParameter bool) bool {
	// Handle symbolic patterns first
	if isPattern, nameExpr, blankExpr := core.IsSymbolicPattern(pattern); isPattern {
		var varName string
		if name, ok := core.ExtractSymbol(nameExpr); ok {
			varName = name
		}

		// Check if the blank expression matches
		if pe.matchBlankWithBinding(blankExpr, expr, ctx) {
			// Bind the variable if it has a name
			if varName != "" {
				if err := ctx.Set(varName, expr); err != nil {
					// Pattern matching should not fail due to protection - log but continue
					// This is a temporary binding during pattern matching
					return false
				}
			}
			return true
		}
		return false
	}

	// Handle direct symbolic blanks
	if isBlank, _, _ := core.IsSymbolicBlank(pattern); isBlank {
		return pe.matchBlankWithBinding(pattern, expr, ctx)
	}

	// Handle different expression types
	switch p := pattern.(type) {
	case core.Symbol:
		varName := p.String()
		// Check if it's a pattern variable
		if core.IsPatternVariable(varName) {
			info := core.ParsePatternInfo(varName)
			if info.Type == core.BlankPattern {
				// Check type constraint if present
				if !core.MatchesType(expr, info.TypeName) {
					return false
				}

				// Bind the variable if named
				if info.VarName != "" {
					if err := ctx.Set(info.VarName, expr); err != nil {
						return false
					}
				}
				return true
			}
			// Sequence patterns need special handling in list context
			return false
		} else {
			// Regular symbol behavior depends on context
			if isParameter {
				// In parameter lists, regular symbols bind to values
				if err := ctx.Set(varName, expr); err != nil {
					return false
				}
				return true
			} else {
				// In head patterns, regular symbols require exact matches
				if exprName, ok := core.ExtractSymbol(expr); ok {
					return exprName == varName
				}
				return false
			}
		}
	case core.List:
		exprList, ok := expr.(core.List)
		if !ok {
			return false
		}

		return pe.matchListWithBinding(p, exprList, ctx)

	default:
		// All other types must match exactly
		return pattern.Equal(expr)
	}
}

// matchBlankWithBinding matches a blank expression with binding
func (pe *PatternExecutor) matchBlankWithBinding(blankExpr, expr core.Expr, ctx *Context) bool {
	isBlank, blankType, typeExpr := core.IsSymbolicBlank(blankExpr)
	if !isBlank {
		return false
	}

	// Check type constraint if present
	if typeExpr != nil {
		var typeName string
		if name, ok := core.ExtractSymbol(typeExpr); ok {
			typeName = name
		}
		if !core.MatchesType(expr, typeName) {
			return false
		}
	}

	// For now, single blank expressions always match single expressions
	// BlankSequence and BlankNullSequence handling for sequences happens in list context
	switch blankType {
	case core.BlankPattern, core.BlankSequencePattern, core.BlankNullSequencePattern:
		return true
	}

	return false
}

// matchListWithBinding matches two lists with binding
func (pe *PatternExecutor) matchListWithBinding(patternList, exprList core.List, ctx *Context) bool {
	lhs := patternList.AsSlice()
	rhs := exprList.AsSlice()

	// Simple case: same length, no sequence patterns
	if !pe.hasSequencePatterns(patternList) {
		if len(lhs) != len(rhs) {
			return false
		}

		// Match each element
		for i, patternElem := range lhs {
			// Element 0 is head (literal), elements 1+ are parameters (pattern match)
			isParameterPosition := i > 0
			if !pe.matchWithBindingInternal(patternElem, rhs[i], ctx, isParameterPosition) {
				return false
			}
		}

		return true
	}

	// Complex case: has sequence patterns - need sequence matching
	return pe.matchSequencePatternsWithBinding(lhs, rhs, ctx)
}

// hasSequencePatterns checks if a pattern list contains sequence patterns
func (pe *PatternExecutor) hasSequencePatterns(patternList core.List) bool {
	for _, elem := range patternList.AsSlice() {
		// Check for symbolic sequence patterns
		if isBlank, blankType, _ := core.IsSymbolicBlank(elem); isBlank {
			if blankType == core.BlankSequencePattern || blankType == core.BlankNullSequencePattern {
				return true
			}
		}

		// Check for symbolic Pattern[name, sequence]
		if isPattern, _, blankExpr := core.IsSymbolicPattern(elem); isPattern {
			if isBlank, blankType, _ := core.IsSymbolicBlank(blankExpr); isBlank {
				if blankType == core.BlankSequencePattern || blankType == core.BlankNullSequencePattern {
					return true
				}
			}
		}

		// Check for string-based sequence patterns
		if atom, ok := elem.(core.String); ok {
			name := atom.String()
			if core.IsPatternVariable(name) {
				info := core.ParsePatternInfo(name)
				if info.Type == core.BlankSequencePattern || info.Type == core.BlankNullSequencePattern {
					return true
				}
			}
		}
	}
	return false
}

// matchSequencePatternsWithBinding handles complex sequence pattern matching
func (pe *PatternExecutor) matchSequencePatternsWithBinding(patterns, exprs []core.Expr, ctx *Context) bool {
	patternIndex := 0
	exprIndex := 0

	// Skip head in both patterns and expressions for parameter matching
	if len(patterns) > 0 && len(exprs) > 0 {
		// Match head literally (no binding)
		if !pe.matchWithBindingInternal(patterns[0], exprs[0], ctx, false) {
			return false
		}
		patternIndex = 1
		exprIndex = 1
	}

	// Match parameters with sequence support
	for patternIndex < len(patterns) && exprIndex <= len(exprs) {
		pattern := patterns[patternIndex]

		// Check what kind of pattern this is
		patternInfo := pe.getPatternTypeInfo(pattern)

		switch patternInfo.Type {
		case core.BlankPattern:
			// Single pattern - match one expression
			if exprIndex >= len(exprs) {
				return false
			}

			// Check type constraint
			if !core.MatchesType(exprs[exprIndex], patternInfo.TypeName) {
				return false
			}

			// Bind variable if named
			if patternInfo.VarName != "" {
				if err := ctx.Set(patternInfo.VarName, exprs[exprIndex]); err != nil {
					return false
				}
			}

			exprIndex++
			patternIndex++

		case core.BlankSequencePattern:
			// Sequence pattern - match 1 or more expressions
			sequenceExprs := pe.collectSequenceExprs(patterns, patternIndex, exprs, exprIndex, len(exprs), patternInfo.TypeName)
			if len(sequenceExprs) == 0 {
				return false // BlankSequence requires at least 1 match
			}

			// Bind variable if named
			if patternInfo.VarName != "" {
				if err := ctx.Set(patternInfo.VarName, core.NewList("List", sequenceExprs...)); err != nil {
					return false
				}
			}

			exprIndex += len(sequenceExprs)
			patternIndex++

		case core.BlankNullSequencePattern:
			// Null sequence pattern - match 0 or more expressions
			sequenceExprs := pe.collectSequenceExprs(patterns, patternIndex, exprs, exprIndex, len(exprs), patternInfo.TypeName)

			// Bind variable if named (can be empty list)
			if patternInfo.VarName != "" {
				if err := ctx.Set(patternInfo.VarName, core.NewList("List", sequenceExprs...)); err != nil {
					return false
				}
			}

			exprIndex += len(sequenceExprs)
			patternIndex++

		default:
			// Regular pattern - match one expression
			if exprIndex >= len(exprs) {
				return false
			}

			if !pe.matchWithBindingInternal(pattern, exprs[exprIndex], ctx, true) {
				return false
			}

			exprIndex++
			patternIndex++
		}
	}

	// Check if we matched all patterns and expressions
	return patternIndex == len(patterns) && exprIndex == len(exprs)
}

// getPatternTypeInfo extracts pattern type information from any pattern format
func (pe *PatternExecutor) getPatternTypeInfo(pattern core.Expr) core.PatternInfo {
	// Check symbolic Pattern[name, blank]
	if isPattern, _, _ := core.IsSymbolicPattern(pattern); isPattern {
		return core.GetSymbolicPatternInfo(pattern)
	}

	// Check direct symbolic blank
	if isBlank, _, _ := core.IsSymbolicBlank(pattern); isBlank {
		return core.GetSymbolicPatternInfo(pattern)
	}

	// Check string-based pattern
	if name, ok := core.ExtractSymbol(pattern); ok {
		if core.IsPatternVariable(name) {
			return core.ParsePatternInfo(name)
		}
	}

	return core.PatternInfo{}
}

// collectSequenceExprs collects expressions for sequence patterns
func (pe *PatternExecutor) collectSequenceExprs(patterns []core.Expr, patternIndex int, exprs []core.Expr, startExpr, endExpr int, typeName string) []core.Expr {
	// Simple greedy collection for now
	// More sophisticated matching would consider remaining patterns

	var collected []core.Expr
	for i := startExpr; i < endExpr; i++ {
		// Check type constraint if present
		if typeName != "" && !core.MatchesType(exprs[i], typeName) {
			break
		}
		collected = append(collected, exprs[i])

		// Check if we need to leave expressions for remaining patterns
		remainingPatterns := len(patterns) - patternIndex - 1
		remainingExprs := endExpr - i - 1
		if remainingExprs < remainingPatterns {
			break
		}
	}

	return collected
}
