package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// ClearAttributesExpr clears attributes from a symbol: ClearAttributes(symbol, attr) or ClearAttributes(symbol, {attr1, attr2})
func ClearAttributesExpr(e *engine.Evaluator, c *engine.Context, symbol, attrs core.Expr) core.Expr {
	symbolTable := c.GetSymbolTable()

	// The first argument should be a symbol
	if symbolName, ok := core.ExtractSymbol(symbol); ok {
		// Handle single attribute
		if attrName, ok := core.ExtractSymbol(attrs); ok {
			if attr, ok := engine.StringToAttribute(attrName); ok {
				symbolTable.ClearAttributes(symbolName, []engine.Attribute{attr})
				return core.NewSymbol("Null")
			}
		}

		// Handle list of attributes
		if attrList, ok := attrs.(core.List); ok && attrList.Length() > 0 {
			if attrList.Head() == "List" {
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
					return core.NewSymbol("Null")
				}
			}
		}
	}

	return core.NewError("ArgumentError", "Invalid arguments to ClearAttributes")
}

/*
// Helper function to parse attribute names to engine.Attribute objects
func parseAttribute(name string) (engine.Attribute, bool) {
	return engine.StringToAttribute(name)
}
*/
