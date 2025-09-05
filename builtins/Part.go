package builtins

import (
	"fmt"

	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Part

// PartList extracts an element from a list by integer index (1-based)
// @ExprPattern (_, _Integer)
func PartList(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	expr := args[0]
	n, _ := core.ExtractInt64(args[1])
	return core.Part(expr, n)
}

// PartAssociation extracts a value from an association by key
// @ExprPattern (_Association, _)
func PartAssociation(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	assoc := args[0].(core.Association)
	key := args[1]
	// For associations, use the key argument to lookup value
	if value, exists := assoc.Get(key); exists {
		return value
	}
	return core.NewError("PartError",
		fmt.Sprintf("Key %s not found in association", key.String()))
}
