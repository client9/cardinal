package stdlib

import (
	"github.com/client9/sexpr/core"
)

// PatternSpecificity is a placeholder - the actual implementation is in the main package
// where it has access to pattern parsing and specificity calculation
func PatternSpecificity(patternStr core.Expr) core.Expr {
	// This function will be called by WrapPatternSpecificity in builtin_funcs.go
	// which has access to the full context and functionality
	return core.NewErrorExpr("InternalError", "PatternSpecificity should be called via wrapper", []core.Expr{patternStr})
}

// ShowPatterns is a placeholder - the actual implementation is in the main package
// where it has access to the function registry
func ShowPatterns(functionName core.Expr) core.Expr {
	// This function will be called by WrapShowPatterns in builtin_funcs.go
	// which has access to the full context and functionality
	return core.NewErrorExpr("InternalError", "ShowPatterns should be called via wrapper", []core.Expr{functionName})
}
