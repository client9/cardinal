package core

import "fmt"

// RuleDelayedExpr represents a delayed rule created by RuleDelayed(pattern, rhs)
// The rhs is held unevaluated until the rule is applied, providing lexical scoping
type RuleDelayedExpr struct {
	Pattern Expr // Pattern (evaluated)
	RHS     Expr // Right-hand side (held unevaluated)
}

// String returns the string representation of the rule
func (r RuleDelayedExpr) String() string {
	return fmt.Sprintf("RuleDelayed(%s, %s)", r.Pattern.String(), r.RHS.String())
}

// InputForm returns the input form representation
func (r RuleDelayedExpr) InputForm() string {
	return r.String() // Same as String for now
}

// Head returns the head of the expression
func (r RuleDelayedExpr) Head() string {
	return "RuleDelayed"
}

// Length returns the length (always 2 for pattern and rhs)
func (r RuleDelayedExpr) Length() int64 {
	return 2
}

// Equal checks equality with another expression
func (r RuleDelayedExpr) Equal(rhs Expr) bool {
	if other, ok := rhs.(RuleDelayedExpr); ok {
		return r.Pattern.Equal(other.Pattern) && r.RHS.Equal(other.RHS)
	}
	return false
}

// IsAtom returns false since rules are composite
func (r RuleDelayedExpr) IsAtom() bool {
	return false
}

// NewRuleDelayed creates a new RuleDelayedExpr
func NewRuleDelayed(pattern, rhs Expr) RuleDelayedExpr {
	return RuleDelayedExpr{
		Pattern: pattern,
		RHS:     rhs,
	}
}
