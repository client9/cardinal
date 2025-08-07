package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"

	"fmt"
)

// PatternSpecificityExpr calculates the specificity of a pattern expression for debugging
// TODO: could move to core and directly use
func PatternSpecificity(e *engine.Evaluator, c *engine.Context, arg core.Expr) core.Expr {
	specificity := core.GetPatternSpecificity(arg)
	return core.NewInteger(int64(specificity))
}

// ShowPatternsExpr lists all registered patterns for a function name
func ShowPatterns(e *engine.Evaluator, c *engine.Context, functionName core.Expr) core.Expr {
	if funcName, ok := core.ExtractSymbol(functionName); ok {

		// Get function definitions from the registry
		definitions := c.GetFunctionDefinitions(funcName)
		if definitions == nil {
			return core.NewError("ArgumentError",
				fmt.Sprintf("No patterns found for function: %s", funcName))
		}

		// Create a list of pattern information
		elements := make([]core.Expr, len(definitions)+1)
		elements[0] = core.NewSymbol("List")

		for i, def := range definitions {
			// Create a rule showing pattern -> specificity
			patternStr := def.Pattern.String()
			specificityStr := fmt.Sprintf("%d", def.Specificity)

			ruleElements := []core.Expr{
				core.NewSymbol("Rule"),
				core.NewString(patternStr),
				core.NewString(specificityStr),
			}

			elements[i+1] = core.NewListFromExprs(ruleElements...)
		}

		return core.NewListFromExprs(elements...)
	}

	return core.NewError("ArgumentError", "ShowPatterns expects a symbol")
}
