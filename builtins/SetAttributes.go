package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/symbol"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol SetAttributes
// @ExprAttributes HoldFirst
// TODO Cleanup -- split into two functions

// SetAttributesExpr sets attributes for a symbol: SetAttributes(symbol, attr) or SetAttributes(symbol, {attr1, attr2})
// @ExprPattern (_Symbol, _Symbol)
func SetAttributesSingle(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	sym := args[0].(core.Symbol)
	attr := engine.SymbolToAttribute(args[1].(core.Symbol))
	symbolTable := c.GetSymbolTable()
	symbolTable.SetAttributes(sym, attr)
	return symbol.Null
}

// @ExprPattern (_Symbol, List(___Symbol))
func SetAttributesList(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	sym := args[0].(core.Symbol)
	attrs := args[1].(core.List).Tail()

	var attribute engine.Attribute
	for _, a := range attrs {
		attribute |= engine.SymbolToAttribute(a)
	}

	symbolTable := c.GetSymbolTable()
	symbolTable.SetAttributes(sym, attribute)
	return symbol.Null
}
