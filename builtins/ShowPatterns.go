package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/atom"
	"github.com/client9/sexpr/engine"

	"fmt"
)

// @ExprSymbol ShowPatterns
// ShowPatternsExpr lists all registered patterns for a function name

// @ExprPattern (_Symbol)
func ShowPatterns(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {

	functionName := args[0]

	if funcName, ok := core.ExtractSymbol(functionName); ok {

		// Get function definitions from the registry
		definitions := c.GetFunctionDefinitions(funcName)
		if definitions == nil {
			return core.NewError("ArgumentError",
				fmt.Sprintf("No patterns found for function: %s", funcName))
		}

		// Create a list of pattern information
		elements := make([]core.Expr, len(definitions)+1)
		elements[0] = core.SymbolFor(atom.List)

		for i, def := range definitions {
			// Create a rule showing pattern -> specificity
			patternStr := def.Pattern.String()
			specificityStr := fmt.Sprintf("%d", def.Specificity)

			ruleElements := []core.Expr{
				core.SymbolFor(atom.Rule),
				core.NewString(patternStr),
				core.NewString(specificityStr),
			}

			elements[i+1] = core.NewListFromExprs(ruleElements...)
		}

		return core.NewListFromExprs(elements...)
	}

	return core.NewError("ArgumentError", "ShowPatterns expects a symbol")
}
