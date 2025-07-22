package sexpr

// Hybrid comparison wrappers that handle fallback behavior correctly
// These are used instead of the generated wrappers for numeric comparison functions

// WrapLessExprsHybrid wraps LessExprs with proper fallback behavior
func WrapLessExprsHybrid(args []Expr, ctx *Context) Expr {
	// Validate argument count
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			"Less expects 2 arguments", args)
	}

	arg0 := args[0]
	arg1 := args[1]

	// Check if both arguments are numeric
	if isNumeric(arg0) && isNumeric(arg1) {
		// Use the generated wrapper for numeric types
		result := LessExprs(arg0, arg1)
		return NewBoolAtom(result)
	}

	// Fall back to unchanged expression for non-numeric types
	return List{Elements: []Expr{NewSymbolAtom("Less"), arg0, arg1}}
}

// WrapGreaterExprsHybrid wraps GreaterExprs with proper fallback behavior
func WrapGreaterExprsHybrid(args []Expr, ctx *Context) Expr {
	// Validate argument count
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			"Greater expects 2 arguments", args)
	}

	arg0 := args[0]
	arg1 := args[1]

	// Check if both arguments are numeric
	if isNumeric(arg0) && isNumeric(arg1) {
		// Use the business logic function for numeric types
		result := GreaterExprs(arg0, arg1)
		return NewBoolAtom(result)
	}

	// Fall back to unchanged expression for non-numeric types
	return List{Elements: []Expr{NewSymbolAtom("Greater"), arg0, arg1}}
}

// WrapLessEqualExprsHybrid wraps LessEqualExprs with proper fallback behavior
func WrapLessEqualExprsHybrid(args []Expr, ctx *Context) Expr {
	// Validate argument count
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			"LessEqual expects 2 arguments", args)
	}

	arg0 := args[0]
	arg1 := args[1]

	// Check if both arguments are numeric
	if isNumeric(arg0) && isNumeric(arg1) {
		// Use the business logic function for numeric types
		result := LessEqualExprs(arg0, arg1)
		return NewBoolAtom(result)
	}

	// Fall back to unchanged expression for non-numeric types
	return List{Elements: []Expr{NewSymbolAtom("LessEqual"), arg0, arg1}}
}

// WrapGreaterEqualExprsHybrid wraps GreaterEqualExprs with proper fallback behavior
func WrapGreaterEqualExprsHybrid(args []Expr, ctx *Context) Expr {
	// Validate argument count
	if len(args) != 2 {
		return NewErrorExpr("ArgumentError",
			"GreaterEqual expects 2 arguments", args)
	}

	arg0 := args[0]
	arg1 := args[1]

	// Check if both arguments are numeric
	if isNumeric(arg0) && isNumeric(arg1) {
		// Use the business logic function for numeric types
		result := GreaterEqualExprs(arg0, arg1)
		return NewBoolAtom(result)
	}

	// Fall back to unchanged expression for non-numeric types
	return List{Elements: []Expr{NewSymbolAtom("GreaterEqual"), arg0, arg1}}
}
