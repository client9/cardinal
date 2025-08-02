package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
	"fmt"
)

// SetAttributesExpr sets attributes for a symbol: SetAttributes(symbol, attr) or SetAttributes(symbol, {attr1, attr2})
func SetAttributesExpr(e *engine.Evaluator, c *engine.Context,  symbol, attrs core.Expr) core.Expr {
	symbolTable := c.GetSymbolTable()

	// The first argument should be a symbol
	if symbolName, ok := core.ExtractSymbol(symbol); ok {
		// Handle single attribute
		if attrName, ok := core.ExtractSymbol(attrs); ok {
			if attr, ok := engine.StringToAttribute(attrName); ok {
				symbolTable.SetAttributes(symbolName, []engine.Attribute{attr})
				return core.NewSymbol("Null")
			}
			return core.NewErrorExpr("Attribute", fmt.Sprintf("unknown attribute %q", attrName), nil)
		}

		// Handle list of attributes
		if attrList, ok := attrs.(core.List); ok && len(attrList.Elements) > 0 {
			if listHead, ok := core.ExtractSymbol(attrList.Elements[0]); ok && listHead == "List" {
				var attributes []engine.Attribute
				for i := 1; i < len(attrList.Elements); i++ {
					if attrName, ok := core.ExtractSymbol(attrList.Elements[i]); ok {
						if attr, ok := engine.StringToAttribute(attrName); ok {
							attributes = append(attributes, attr)
						} else {
							return core.NewErrorExpr("Attribute",  fmt.Sprintf("unknown attribute %q", attrName), nil)
						}
					}
				}
				if len(attributes) > 0 {
					symbolTable.SetAttributes(symbolName, attributes)
					return core.NewSymbol("Null")
				}
			}
		}
	}

	return core.NewErrorExpr("ArgumentError", "Invalid arguments to SetAttributes", []core.Expr{symbol, attrs})
}

