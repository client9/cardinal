package builtins

import (
	"fmt"
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol SetAttributes
// @ExprAttributes HoldFirst
// TODO Cleanup -- split into two functions

// SetAttributesExpr sets attributes for a symbol: SetAttributes(symbol, attr) or SetAttributes(symbol, {attr1, attr2})
// @ExprPattern (_Symbol,_)
func SetAttributesExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {

	symbol := args[0]
	attrs := args[1]

	symbolTable := c.GetSymbolTable()

	// The first argument should be a symbol
	if symbolName, ok := core.ExtractSymbol(symbol); ok {
		// Handle single attribute
		if attrName, ok := core.ExtractSymbol(attrs); ok {
			if attr, ok := engine.StringToAttribute(attrName); ok {
				symbolTable.SetAttributes(symbolName, []engine.Attribute{attr})
				return core.NewSymbol("Null")
			}
			return core.NewError("Attribute", fmt.Sprintf("unknown attribute %q", attrName))
		}

		// Handle list of attributes
		if attrList, ok := attrs.(core.List); ok && attrList.Length() > 0 {
			if attrList.Head() == "List" {
				var attributes []engine.Attribute
				for _, arg := range attrList.Tail() {
					if attrName, ok := core.ExtractSymbol(arg); ok {
						if attr, ok := engine.StringToAttribute(attrName); ok {
							attributes = append(attributes, attr)
						} else {
							return core.NewError("Attribute", fmt.Sprintf("unknown attribute %q", attrName))
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

	return core.NewError("ArgumentError", "Invalid arguments to SetAttributes")
}
