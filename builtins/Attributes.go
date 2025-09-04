package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/symbol"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Attributes
// @ExprAttributes HoldFirst

// AttributesExpr gets the attributes of a symbol
//
// @ExprPattern (_Symbol)
func AttributesExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {

	symbolTable := c.GetSymbolTable()

	symbolName := args[0].(core.Symbol)

	attrs := symbolTable.Attributes(symbolName)

	return core.NewList(symbol.List, engine.AttributeToSymbols(attrs)...)
}
