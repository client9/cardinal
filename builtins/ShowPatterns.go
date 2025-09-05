package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/symbol"
	"github.com/client9/cardinal/engine"

	"fmt"
)

// @ExprSymbol ShowPatterns
// ShowPatternsExpr lists all registered patterns for a function name

// @ExprPattern (_Symbol)
func ShowPatterns(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {

	funcName := args[0].(core.Symbol)

	// Get function definitions from the registry
	definitions := c.GetFunctionDefinitions(funcName)
	if definitions == nil {
		return core.NewError("ArgumentError",
			fmt.Sprintf("No patterns found for function: %s", funcName))
	}

	// Create a list of pattern information
	elements := make([]core.Expr, len(definitions)+1)
	elements[0] = symbol.List

	for i, def := range definitions {
		// Create a rule showing pattern -> specificity
		patternStr := def.Pattern.String()
		specificityStr := fmt.Sprintf("%d", def.Specificity)

		ruleElements := []core.Expr{
			core.NewString(patternStr),
			core.NewString(specificityStr),
		}

		elements[i+1] = core.ListFrom(symbol.Rule, ruleElements...)
	}

	return core.NewListFromExprs(elements...)

}
