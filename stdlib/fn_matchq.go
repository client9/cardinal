package stdlib

import (
	"github.com/client9/sexpr/core"
)

// MatchQExprs checks if an expression matches a pattern (pure test, no variable binding)
func MatchQExprs(expr, pattern core.Expr) bool {
	// Use the pure pattern matcher from core (no Context needed for pure testing)
	matcher := core.NewPatternMatcher()
	return matcher.TestMatch(pattern, expr)
}
