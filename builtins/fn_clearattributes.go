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
		if attrList, ok := attrs.(core.List); ok && len(attrList.Elements) > 0 {
			if listHead, ok := core.ExtractSymbol(attrList.Elements[0]); ok && listHead == "List" {
				var attributes []engine.Attribute
				for i := 1; i < len(attrList.Elements); i++ {
					if attrName, ok := core.ExtractSymbol(attrList.Elements[i]); ok {
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

	return core.NewErrorExpr("ArgumentError", "Invalid arguments to ClearAttributes", []core.Expr{symbol, attrs})
}

/*
// Helper function to parse attribute names to engine.Attribute objects
func parseAttribute(name string) (engine.Attribute, bool) {
	return engine.StringToAttribute(name)
}
*/
