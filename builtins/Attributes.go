package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
	"log"
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

	log.Printf("attr symbols %v", engine.AttributeToSymbols(attrs))
	return core.ListExpr(engine.AttributeToSymbols(attrs)...)
}
