package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// AttributesExpr gets the attributes of a symbol
func AttributesExpr(evaluator Evaluator, symbol core.Expr) core.Expr {
	ctx := evaluator.GetContext()
	symbolTable := ctx.GetSymbolTable()

	// The argument should be a symbol
	if symbolName, ok := core.ExtractSymbol(symbol); ok {
		// Get the attributes from the symbol table
		attrs := symbolTable.Attributes(symbolName)

		// Convert attributes to a list of symbols
		attrElements := make([]core.Expr, len(attrs)+1)
		attrElements[0] = core.NewSymbol("List")

		for i, attr := range attrs {
			attrElements[i+1] = core.NewSymbol(attr.String())
		}

		return core.List{Elements: attrElements}
	}

	return core.NewErrorExpr("ArgumentError",
		"Attributes expects a symbol as argument", []core.Expr{symbol})
}

// SetAttributesExpr sets attributes for a symbol: SetAttributes(symbol, attr) or SetAttributes(symbol, {attr1, attr2})
func SetAttributesExpr(evaluator Evaluator, symbol, attrs core.Expr) core.Expr {
	ctx := evaluator.GetContext()
	symbolTable := ctx.GetSymbolTable()

	// The first argument should be a symbol
	if symbolName, ok := core.ExtractSymbol(symbol); ok {
		// Handle single attribute
		if attrName, ok := core.ExtractSymbol(attrs); ok {
			if attr, ok := parseAttribute(attrName); ok {
				symbolTable.SetAttributes(symbolName, []engine.Attribute{attr})
				return core.NewSymbol("Null")
			}
		}

		// Handle list of attributes
		if attrList, ok := attrs.(core.List); ok && len(attrList.Elements) > 0 {
			if listHead, ok := core.ExtractSymbol(attrList.Elements[0]); ok && listHead == "List" {
				var attributes []engine.Attribute
				for i := 1; i < len(attrList.Elements); i++ {
					if attrName, ok := core.ExtractSymbol(attrList.Elements[i]); ok {
						if attr, ok := parseAttribute(attrName); ok {
							attributes = append(attributes, attr)
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

// ClearAttributesExpr clears attributes from a symbol: ClearAttributes(symbol, attr) or ClearAttributes(symbol, {attr1, attr2})
func ClearAttributesExpr(evaluator Evaluator, symbol, attrs core.Expr) core.Expr {
	ctx := evaluator.GetContext()
	symbolTable := ctx.GetSymbolTable()

	// The first argument should be a symbol
	if symbolName, ok := core.ExtractSymbol(symbol); ok {
		// Handle single attribute
		if attrName, ok := core.ExtractSymbol(attrs); ok {
			if attr, ok := parseAttribute(attrName); ok {
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
						if attr, ok := parseAttribute(attrName); ok {
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

// Helper function to parse attribute names to engine.Attribute objects
func parseAttribute(name string) (engine.Attribute, bool) {
	return engine.StringToAttribute(name)
}
