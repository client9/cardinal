package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/symbol"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol And
// @ExprAttributes HoldAll

// AndExpr evaluates logical AND with short-circuiting: And(expr1, expr2, ...)
//
// @ExprPattern (___)
func AndExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	var unevaluatedArgs []core.Expr

	// Short-circuit evaluation: stop at first false, collect non-boolean true values
	for _, arg := range args {
		result := e.Evaluate(arg)

		// Check if it's explicitly False
		if symbolName, ok := core.ExtractSymbol(result); ok && symbolName == "False" {
			return core.NewBool(false)
		}

		// Check if it's explicitly True - continue without adding to unevaluated
		if symbolName, ok := core.ExtractSymbol(result); ok && symbolName == "True" {
			continue
		}

		// For non-boolean values, collect them
		unevaluatedArgs = append(unevaluatedArgs, result)
	}

	// If no unevaluated args remain, all were True
	if len(unevaluatedArgs) == 0 {
		return core.NewBool(true)
	}

	// If only one arg remains, return it directly
	if len(unevaluatedArgs) == 1 {
		return unevaluatedArgs[0]
	}

	// Return And expression with remaining args
	elements := make([]core.Expr, len(unevaluatedArgs)+1)
	elements[0] = symbol.And
	copy(elements[1:], unevaluatedArgs)
	return core.NewListFromExprs(elements...)
}
