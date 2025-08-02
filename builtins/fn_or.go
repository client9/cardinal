package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// OrExpr evaluates logical OR with short-circuiting: Or(expr1, expr2, ...)
func OrExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	var nonFalseArgs []core.Expr

	// Evaluate arguments and collect non-False values
	for _, arg := range args {
		result := e.Evaluate(c, arg)

		// Check if it's explicitly True - short-circuit
		if symbolName, ok := core.ExtractSymbol(result); ok && symbolName == "True" {
			return core.NewSymbol("True")
		}

		// Skip False values
		if symbolName, ok := core.ExtractSymbol(result); ok && symbolName == "False" {
			continue
		}

		// Collect non-False values
		nonFalseArgs = append(nonFalseArgs, result)
	}

	// If no non-False args remain, all were False
	if len(nonFalseArgs) == 0 {
		return core.NewSymbol("False")
	}

	// If only one non-False arg, return it directly
	if len(nonFalseArgs) == 1 {
		return nonFalseArgs[0]
	}

	// Return Or expression with remaining non-False args
	elements := make([]core.Expr, len(nonFalseArgs)+1)
	elements[0] = core.NewSymbol("Or")
	copy(elements[1:], nonFalseArgs)
	return core.List{Elements: elements}
}

