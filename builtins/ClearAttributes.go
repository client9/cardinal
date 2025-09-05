package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/symbol"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol ClearAttributes
// @ExprAttributes HoldFirst

// ClearAttributesExpr clears attributes from a symbol: ClearAttributes(symbol, attr) or ClearAttributes(symbol, {attr1, attr2})
//
// @ExprPattern (_Symbol, _Symbol)
func ClearAttributesSingle(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	sym := args[0].(core.Symbol)
	attrName := args[1].(core.Symbol)
	symbolTable := c.GetSymbolTable()
	symbolTable.ClearAttributes(sym, engine.SymbolToAttribute(attrName))
	return symbol.Null
}

// TODO: no error handling
//
// @ExprPattern (_Symbol, List(___Symbol)))
func ClearAttributesList(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	sym := args[0].(core.Symbol)
	attrList := args[1].(core.List)
	symbolTable := c.GetSymbolTable()
	var attributes engine.Attribute
	for _, arg := range attrList.Tail() {
		attrName := arg.(core.Symbol)
		attributes |= engine.SymbolToAttribute(attrName)
	}
	if attributes != 0 {
		symbolTable.ClearAttributes(sym, attributes)
	}

	return symbol.Null
}
