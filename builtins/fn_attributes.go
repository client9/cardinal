package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// AttributesExpr gets the attributes of a symbol
func AttributesExpr(e *engine.Evaluator, c *engine.Context, symbol core.Expr) core.Expr {
	symbolTable := c.GetSymbolTable()

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
