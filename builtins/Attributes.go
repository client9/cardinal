package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/symbol"
	"github.com/client9/cardinal/engine"
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
