package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol ClearAttributes
// @ExprAttributes HoldFirst

// ClearAttributesExpr clears attributes from a symbol: ClearAttributes(symbol, attr) or ClearAttributes(symbol, {attr1, attr2})
//
// @ExprPattern (_Symbol, _Symbol)
func ClearAttributesSingle(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {

	symbol := args[0].(core.Symbol)
	attrName := args[1].(core.Symbol)

	symbolTable := c.GetSymbolTable()

	symbolName := symbol.String()

	symbolTable.ClearAttributes(symbolName, engine.SymbolToAttribute(attrName))
	return core.NewSymbol("Null")
}

// TODO: no error handling
//
// @ExprPattern (_Symbol, List(___Symbol)))
func ClearAttributesList(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	symbol := args[0]
	attrList := args[1].(core.List)

	symbolTable := c.GetSymbolTable()

	symbolName, _ := core.ExtractSymbol(symbol)

	var attributes engine.Attribute
	for _, arg := range attrList.Tail() {
		attrName := arg.(core.Symbol)
		attributes |= engine.SymbolToAttribute(attrName)
	}
	if attributes != 0 {
		symbolTable.ClearAttributes(symbolName, attributes)
	}

	return core.NewSymbol("Null")
}
