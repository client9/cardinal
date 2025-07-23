package stdlib

// Attribute manipulation functions
// These need access to core types and Context

// AttributesExpr gets the attributes of a symbol
// func AttributesExpr(expr Expr, ctx *Context) Expr {
// 	// The argument should be a symbol
// 	if atom, ok := expr.(Atom); ok && atom.AtomType == SymbolAtom {
// 		symbolName := atom.Value.(string)

// 		// Get the attributes from the symbol table
// 		attrs := ctx.symbolTable.Attributes(symbolName)

// 		// Convert attributes to a list of symbols
// 		attrElements := make([]Expr, len(attrs)+1)
// 		attrElements[0] = NewSymbolAtom("List")

// 		for i, attr := range attrs {
// 			attrElements[i+1] = NewSymbolAtom(attr.String())
// 		}

// 		return List{Elements: attrElements}
// 	}

// 	return NewErrorExpr("ArgumentError",
// 		"Attributes expects a symbol as argument", []Expr{expr})
// }

// SetAttributesSingle sets a single attribute on a symbol
// func SetAttributesSingle(symbol Expr, attr Expr, ctx *Context) Expr {
// 	// The first argument should be a symbol
// 	if atom, ok := symbol.(Atom); ok && atom.AtomType == SymbolAtom {
// 		symbolName := atom.Value.(string)

// 		// The second argument should be an attribute symbol
// 		if attrAtom, ok := attr.(Atom); ok && attrAtom.AtomType == SymbolAtom {
// 			attrName := attrAtom.Value.(string)

// 			// Convert string to Attribute
// 			if attribute, ok := StringToAttribute(attrName); ok {
// 				// Set the attribute on the symbol
// 				ctx.symbolTable.SetAttributes(symbolName, []Attribute{attribute})
// 				return NewSymbolAtom("Null")
// 			}

// 			return NewErrorExpr("ArgumentError",
// 				fmt.Sprintf("Unknown attribute: %s", attrName), []Expr{attr})
// 		}
// 	}

// 	return NewErrorExpr("ArgumentError",
// 		"SetAttributes expects (symbol, attribute)", []Expr{symbol, attr})
// }

// SetAttributesList sets multiple attributes on a symbol
// func SetAttributesList(symbol Expr, attrList List, ctx *Context) Expr {
// 	// The first argument should be a symbol
// 	if atom, ok := symbol.(Atom); ok && atom.AtomType == SymbolAtom {
// 		symbolName := atom.Value.(string)

// 		// Process each attribute in the list (skip head at index 0)
// 		var attributes []Attribute
// 		for i := 1; i < len(attrList.Elements); i++ {
// 			attrExpr := attrList.Elements[i]

// 			if attrAtom, ok := attrExpr.(Atom); ok && attrAtom.AtomType == SymbolAtom {
// 				attrName := attrAtom.Value.(string)

// 				// Convert string to Attribute
// 				if attribute, ok := StringToAttribute(attrName); ok {
// 					attributes = append(attributes, attribute)
// 				} else {
// 					return NewErrorExpr("ArgumentError",
// 						fmt.Sprintf("Unknown attribute: %s", attrName), []Expr{attrExpr})
// 				}
// 			} else {
// 				return NewErrorExpr("ArgumentError",
// 					"Attributes list must contain symbols", []Expr{attrExpr})
// 			}
// 		}

// 		// Set all attributes on the symbol
// 		ctx.symbolTable.SetAttributes(symbolName, attributes)
// 		return NewSymbolAtom("Null")
// 	}

// 	return NewErrorExpr("ArgumentError",
// 		"SetAttributes expects (symbol, attribute list)", []Expr{symbol, attrList})
// }

// ClearAttributesSingle clears a single attribute from a symbol
// func ClearAttributesSingle(symbol Expr, attr Expr, ctx *Context) Expr {
// 	// The first argument should be a symbol
// 	if atom, ok := symbol.(Atom); ok && atom.AtomType == SymbolAtom {
// 		symbolName := atom.Value.(string)

// 		// The second argument should be an attribute symbol
// 		if attrAtom, ok := attr.(Atom); ok && attrAtom.AtomType == SymbolAtom {
// 			attrName := attrAtom.Value.(string)

// 			// Convert string to Attribute
// 			if attribute, ok := StringToAttribute(attrName); ok {
// 				// Clear the attribute from the symbol
// 				ctx.symbolTable.ClearAttributes(symbolName, []Attribute{attribute})
// 				return NewSymbolAtom("Null")
// 			}

// 			return NewErrorExpr("ArgumentError",
// 				fmt.Sprintf("Unknown attribute: %s", attrName), []Expr{attr})
// 		}
// 	}

// 	return NewErrorExpr("ArgumentError",
// 		"ClearAttributes expects (symbol, attribute)", []Expr{symbol, attr})
// }

// ClearAttributesList clears multiple attributes from a symbol
// func ClearAttributesList(symbol Expr, attrList List, ctx *Context) Expr {
// 	// The first argument should be a symbol
// 	if atom, ok := symbol.(Atom); ok && atom.AtomType == SymbolAtom {
// 		symbolName := atom.Value.(string)

// 		// Process each attribute in the list (skip head at index 0)
// 		for i := 1; i < len(attrList.Elements); i++ {
// 			attrExpr := attrList.Elements[i]

// 			if attrAtom, ok := attrExpr.(Atom); ok && attrAtom.AtomType == SymbolAtom {
// 				attrName := attrAtom.Value.(string)

// 				// Convert string to Attribute
// 				if attribute, ok := StringToAttribute(attrName); ok {
// 					// Clear this attribute from the symbol
// 					ctx.symbolTable.ClearAttributes(symbolName, []Attribute{attribute})
// 				} else {
// 					return NewErrorExpr("ArgumentError",
// 						fmt.Sprintf("Unknown attribute: %s", attrName), []Expr{attrExpr})
// 				}
// 			} else {
// 				return NewErrorExpr("ArgumentError",
// 					"Attributes list must contain symbols", []Expr{attrExpr})
// 			}
// 		}

// 		return NewSymbolAtom("Null")
// 	}

// 	return NewErrorExpr("ArgumentError",
// 		"ClearAttributes expects (symbol, attribute list)", []Expr{symbol, attrList})
// }
