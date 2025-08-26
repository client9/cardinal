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

	symbol := args[0]
	attrName := args[1]

	symbolTable := c.GetSymbolTable()

	symbolName, _ := core.ExtractSymbol(symbol)

	// Handle single attribute
	if attrName, ok := core.ExtractSymbol(attrName); ok {
		if attr, ok := engine.StringToAttribute(attrName); ok {
			symbolTable.ClearAttributes(symbolName, []engine.Attribute{attr})
		}
	}
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

	var attributes []engine.Attribute
	for _, arg := range attrList.Tail() {
		if attrName, ok := core.ExtractSymbol(arg); ok {
			if attr, ok := engine.StringToAttribute(attrName); ok {
				attributes = append(attributes, attr)
			}
		}
	}
	if len(attributes) > 0 {
		symbolTable.ClearAttributes(symbolName, attributes)
	}

	return core.NewSymbol("Null")
}

/*
// Helper function to parse attribute names to engine.Attribute objects
func parseAttribute(name string) (engine.Attribute, bool) {
	return engine.StringToAttribute(name)
}
*/
