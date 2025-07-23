package sexpr

import (
	"fmt"
)

//go:generate go run cmd/wrapgen/main.go -single builtin_wrappers.go
//go:generate go run cmd/wrapgen/main.go -setup builtin_setup.go

// MatchQExprs checks if an expression matches a pattern
func MatchQExprs(expr, pattern Expr, ctx *Context) bool {
	// Convert string-based pattern to symbolic if needed
	symbolicPattern := convertToSymbolicPattern(pattern)

	// Create a temporary context for pattern matching (don't pollute original context)
	tempCtx := NewChildContext(ctx)

	// Use the existing pattern matching logic
	return matchPatternForMatchQ(symbolicPattern, expr, tempCtx)
}

// WrapMatchQExprs is a clean wrapper for MatchQ that uses the business logic function
func WrapMatchQExprs(args []Expr, ctx *Context) Expr {
	// Validate argument count
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			"MatchQ expects 2 arguments", args)
	}

	// Check for errors in arguments first
	for _, arg := range args {
		if IsError(arg) {
			return arg
		}
	}

	// Call business logic function
	result := MatchQExprs(args[0], args[1], ctx)

	// Convert result back to Expr
	return NewBoolAtom(result)
}


// AttributesExpr gets the attributes of a symbol
func AttributesExpr(expr Expr, ctx *Context) Expr {
	// The argument should be a symbol
	if atom, ok := expr.(Atom); ok && atom.AtomType == SymbolAtom {
		symbolName := atom.Value.(string)

		// Get the attributes from the symbol table
		attrs := ctx.symbolTable.Attributes(symbolName)

		// Convert attributes to a list of symbols
		attrElements := make([]Expr, len(attrs)+1)
		attrElements[0] = NewSymbolAtom("List")

		for i, attr := range attrs {
			attrElements[i+1] = NewSymbolAtom(attr.String())
		}

		return List{Elements: attrElements}
	}

	return NewErrorExpr("ArgumentError",
		"Attributes expects a symbol as argument", []Expr{expr})
}

// WrapAttributesExpr is a clean wrapper for Attributes that uses the business logic function
func WrapAttributesExpr(args []Expr, ctx *Context) Expr {
	// Validate argument count
	if len(args) != 1 {
		return NewErrorExpr("ArgumentError",
			"Attributes expects 1 argument", args)
	}

	// Check for errors in arguments first
	if IsError(args[0]) {
		return args[0]
	}

	// Call business logic function
	return AttributesExpr(args[0], ctx)
}

// SetAttributesSingle sets a single attribute on a symbol
func SetAttributesSingle(symbol Expr, attr Expr, ctx *Context) Expr {
	// The first argument should be a symbol
	if atom, ok := symbol.(Atom); ok && atom.AtomType == SymbolAtom {
		symbolName := atom.Value.(string)

		// The second argument should be an attribute symbol
		if attrAtom, ok := attr.(Atom); ok && attrAtom.AtomType == SymbolAtom {
			attrName := attrAtom.Value.(string)

			// Convert string to Attribute
			if attribute, ok := StringToAttribute(attrName); ok {
				// Set the attribute on the symbol
				ctx.symbolTable.SetAttributes(symbolName, []Attribute{attribute})
				return NewSymbolAtom("Null")
			}

			return NewErrorExpr("ArgumentError",
				fmt.Sprintf("Unknown attribute: %s", attrName), []Expr{attr})
		}
	}

	return NewErrorExpr("ArgumentError",
		"SetAttributes expects (symbol, attribute)", []Expr{symbol, attr})
}

// SetAttributesList sets multiple attributes on a symbol
func SetAttributesList(symbol Expr, attrList List, ctx *Context) Expr {
	// The first argument should be a symbol
	if atom, ok := symbol.(Atom); ok && atom.AtomType == SymbolAtom {
		symbolName := atom.Value.(string)

		// Process each attribute in the list (skip head at index 0)
		var attributes []Attribute
		for i := 1; i < len(attrList.Elements); i++ {
			attrExpr := attrList.Elements[i]

			if attrAtom, ok := attrExpr.(Atom); ok && attrAtom.AtomType == SymbolAtom {
				attrName := attrAtom.Value.(string)

				// Convert string to Attribute
				if attribute, ok := StringToAttribute(attrName); ok {
					attributes = append(attributes, attribute)
				} else {
					return NewErrorExpr("ArgumentError",
						fmt.Sprintf("Unknown attribute: %s", attrName), []Expr{attrExpr})
				}
			} else {
				return NewErrorExpr("ArgumentError",
					"Attributes list must contain symbols", []Expr{attrExpr})
			}
		}

		// Set all attributes on the symbol
		ctx.symbolTable.SetAttributes(symbolName, attributes)
		return NewSymbolAtom("Null")
	}

	return NewErrorExpr("ArgumentError",
		"SetAttributes expects (symbol, attribute list)", []Expr{symbol, attrList})
}

// ClearAttributesSingle clears a single attribute from a symbol
func ClearAttributesSingle(symbol Expr, attr Expr, ctx *Context) Expr {
	// The first argument should be a symbol
	if atom, ok := symbol.(Atom); ok && atom.AtomType == SymbolAtom {
		symbolName := atom.Value.(string)

		// The second argument should be an attribute symbol
		if attrAtom, ok := attr.(Atom); ok && attrAtom.AtomType == SymbolAtom {
			attrName := attrAtom.Value.(string)

			// Convert string to Attribute
			if attribute, ok := StringToAttribute(attrName); ok {
				// Clear the attribute from the symbol
				ctx.symbolTable.ClearAttributes(symbolName, []Attribute{attribute})
				return NewSymbolAtom("Null")
			}

			return NewErrorExpr("ArgumentError",
				fmt.Sprintf("Unknown attribute: %s", attrName), []Expr{attr})
		}
	}

	return NewErrorExpr("ArgumentError",
		"ClearAttributes expects (symbol, attribute)", []Expr{symbol, attr})
}

// ClearAttributesList clears multiple attributes from a symbol
func ClearAttributesList(symbol Expr, attrList List, ctx *Context) Expr {
	// The first argument should be a symbol
	if atom, ok := symbol.(Atom); ok && atom.AtomType == SymbolAtom {
		symbolName := atom.Value.(string)

		// Process each attribute in the list (skip head at index 0)
		for i := 1; i < len(attrList.Elements); i++ {
			attrExpr := attrList.Elements[i]

			if attrAtom, ok := attrExpr.(Atom); ok && attrAtom.AtomType == SymbolAtom {
				attrName := attrAtom.Value.(string)

				// Convert string to Attribute
				if attribute, ok := StringToAttribute(attrName); ok {
					// Clear this attribute from the symbol
					ctx.symbolTable.ClearAttributes(symbolName, []Attribute{attribute})
				} else {
					return NewErrorExpr("ArgumentError",
						fmt.Sprintf("Unknown attribute: %s", attrName), []Expr{attrExpr})
				}
			} else {
				return NewErrorExpr("ArgumentError",
					"Attributes list must contain symbols", []Expr{attrExpr})
			}
		}

		return NewSymbolAtom("Null")
	}

	return NewErrorExpr("ArgumentError",
		"ClearAttributes expects (symbol, attribute list)", []Expr{symbol, attrList})
}

// Clean wrappers for SetAttributes and ClearAttributes

// WrapSetAttributesSingle is a clean wrapper for SetAttributes with single attribute
func WrapSetAttributesSingle(args []Expr, ctx *Context) Expr {
	// Validate argument count
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			"SetAttributes expects 2 arguments", args)
	}

	// Check for errors in arguments first
	for _, arg := range args {
		if IsError(arg) {
			return arg
		}
	}

	// Call business logic function
	return SetAttributesSingle(args[0], args[1], ctx)
}

// WrapSetAttributesList is a clean wrapper for SetAttributes with attribute list
func WrapSetAttributesList(args []Expr, ctx *Context) Expr {
	// Validate argument count
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			"SetAttributes expects 2 arguments", args)
	}

	// Check for errors in arguments first
	for _, arg := range args {
		if IsError(arg) {
			return arg
		}
	}

	// Extract list from second argument
	attrList, ok := args[1].(List)
	if !ok {
		return NewErrorExpr("ArgumentError",
			"Second argument must be a list", args)
	}

	// Call business logic function
	return SetAttributesList(args[0], attrList, ctx)
}

// WrapClearAttributesSingle is a clean wrapper for ClearAttributes with single attribute
func WrapClearAttributesSingle(args []Expr, ctx *Context) Expr {
	// Validate argument count
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			"ClearAttributes expects 2 arguments", args)
	}

	// Check for errors in arguments first
	for _, arg := range args {
		if IsError(arg) {
			return arg
		}
	}

	// Call business logic function
	return ClearAttributesSingle(args[0], args[1], ctx)
}

// WrapClearAttributesList is a clean wrapper for ClearAttributes with attribute list
func WrapClearAttributesList(args []Expr, ctx *Context) Expr {
	// Validate argument count
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			"ClearAttributes expects 2 arguments", args)
	}

	// Check for errors in arguments first
	for _, arg := range args {
		if IsError(arg) {
			return arg
		}
	}

	// Extract list from second argument
	attrList, ok := args[1].(List)
	if !ok {
		return NewErrorExpr("ArgumentError",
			"Second argument must be a list", args)
	}

	// Call business logic function
	return ClearAttributesList(args[0], attrList, ctx)
}

