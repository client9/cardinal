package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Attributes
// @ExprAttributes HoldFirst

// AttributesExpr gets the attributes of a symbol
//
// @ExprPattern (_Symbol)
func AttributesExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {

	symbolTable := c.GetSymbolTable()

	symbolName := args[0].String()

	attrs := symbolTable.Attributes(symbolName)

	// Convert attributes to a list of symbols
	attrElements := make([]core.Expr, len(attrs)+1)
	attrElements[0] = core.NewSymbol("List")

	for i, attr := range attrs {
		attrElements[i+1] = core.NewSymbol(attr.String())
	}

	return core.NewListFromExprs(attrElements...)
}
