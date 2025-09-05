package builtins

import (
	"fmt"

	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/symbol"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Association
// @ExprAttributes Protected
//
//

// AssociationRules createsd an Association from a sequence of Rule expressions
// @ExprPattern (___Rule)
func AssociationRules(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	assoc := core.NewAssociation()

	// Process each Rule expression
	for _, rule := range args {
		if ruleList, ok := rule.(core.List); ok && ruleList.Length() == 2 && ruleList.Head() == symbol.Rule {
			args := ruleList.Tail()
			assoc = assoc.Set(args[0], args[1]) // Returns new association (immutable)
			continue
		}

		// Invalid argument - not a Rule
		return core.NewError("ArgumentError",
			fmt.Sprintf("Association expects Rule expressions, got %s", rule.String()))
	}

	return assoc
}
