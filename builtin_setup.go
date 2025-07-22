package sexpr

import (
	"fmt"
)

// setupBuiltinAttributes sets up standard attributes for built-in functions
func setupBuiltinAttributes(symbolTable *SymbolTable) {
	// Reset attributes
	symbolTable.Reset()

	// Arithmetic operations
	symbolTable.SetAttributes("Plus", []Attribute{Flat, Listable, NumericFunction, OneIdentity, Orderless, Protected})
	symbolTable.SetAttributes("Times", []Attribute{Flat, Orderless, OneIdentity})
	symbolTable.SetAttributes("Power", []Attribute{OneIdentity})

	// Control structures
	symbolTable.SetAttributes("Hold", []Attribute{HoldAll})
	symbolTable.SetAttributes("If", []Attribute{HoldRest})
	symbolTable.SetAttributes("While", []Attribute{HoldAll})
	symbolTable.SetAttributes("CompoundExpression", []Attribute{HoldAll})
	symbolTable.SetAttributes("Module", []Attribute{HoldAll})
	symbolTable.SetAttributes("Block", []Attribute{HoldAll})

	// Assignment operations
	symbolTable.SetAttributes("Set", []Attribute{HoldFirst})
	symbolTable.SetAttributes("SetDelayed", []Attribute{HoldAll})
	symbolTable.SetAttributes("Unset", []Attribute{HoldFirst})

	// Pattern matching operations
	// symbolTable.SetAttributes("MatchQ", []Attribute{HoldFirst})

	// Attribute functions
	symbolTable.SetAttributes("Attributes", []Attribute{HoldFirst})
	symbolTable.SetAttributes("SetAttributes", []Attribute{HoldFirst})
	symbolTable.SetAttributes("ClearAttributes", []Attribute{HoldFirst})

	// Logical operations
	symbolTable.SetAttributes("And", []Attribute{Flat, HoldAll})
	symbolTable.SetAttributes("Or", []Attribute{Flat, HoldAll})

	// Constants
	symbolTable.SetAttributes("Pi", []Attribute{Constant, Protected})
	symbolTable.SetAttributes("E", []Attribute{Constant, Protected})
	symbolTable.SetAttributes("True", []Attribute{Constant, Protected})
	symbolTable.SetAttributes("False", []Attribute{Constant, Protected})

	// Output functions
	// Both FullForm and InputForm evaluate their arguments, then show different representations
	// FullForm shows the full symbolic form, InputForm shows user-friendly infix notation

	// Pattern symbols
	symbolTable.SetAttributes("Blank", []Attribute{Protected})
	symbolTable.SetAttributes("BlankSequence", []Attribute{Protected})
	symbolTable.SetAttributes("BlankNullSequence", []Attribute{Protected})
	symbolTable.SetAttributes("Pattern", []Attribute{Protected})
}

// registerDefaultBuiltins registers all built-in functions with their patterns
func registerDefaultBuiltins(registry *FunctionRegistry) {
	// Register built-in functions with pattern-based dispatch
	builtinPatterns := map[string]PatternFunc{
		// Arithmetic operations - hierarchical patterns with type-specific optimizations
		// Order matters: more specific patterns first, general fallback last

		// Empty arithmetic operations - direct function mappings for identity values
		"Plus()":  func(args []Expr, ctx *Context) Expr { return PlusEmpty() },   // Additive identity: 0
		"Times()": func(args []Expr, ctx *Context) Expr { return TimesEmpty() },  // Multiplicative identity: 1

		// Fast paths for homogeneous types
		"Plus(x__Integer)":  WrapPlusIntegers,  // Fast integer-only addition
		"Plus(x__Real)":     WrapPlusReals,     // Fast real-only addition
		"Times(x__Integer)": WrapTimesIntegers, // Fast integer-only multiplication
		"Times(x__Real)":    WrapTimesReals,    // Fast real-only multiplication

		// Mixed numeric types - using generated wrappers with Number type constraints
		"Plus(x__Number)":  WrapPlusNumbers,  // Mixed numeric addition (Integer + Real)
		"Times(x__Number)": WrapTimesNumbers, // Mixed numeric multiplication (Integer + Real)

		// Subtraction operations - more specific pattern first
		"Subtract(x_Integer, y_Integer)": WrapSubtractIntegers, // Fast integer-only subtraction
		"Subtract(x_Number, y_Number)":   WrapSubtractNumbers,  // Mixed numeric subtraction (returns float64)

		// Division operations - more specific pattern first
		"Divide(x_Integer, y_Integer)": WrapDivideIntegers, // Fast integer-only division (returns int64)
		"Divide(x_Number, y_Number)":   WrapDivideNumbers,  // Mixed numeric division with error handling
		"Power(x_Number, y_Number)":    WrapPowerNumbers,   // Mixed numeric power operations with error handling

		// Comparison operations - all using generated wrappers
		"Equal(x_, y_)":                    WrapEqualExprs,        // Generated wrapper for equality
		"Unequal(x_, y_)":                  WrapUnequalExprs,      // Generated wrapper for inequality
		"Less(x_Number, y_Number)":         WrapLessExprs,         // Generated wrapper - only for numeric args
		"Greater(x_Number, y_Number)":      WrapGreaterExprs,      // Generated wrapper - only for numeric args
		"LessEqual(x_Number, y_Number)":    WrapLessEqualExprs,    // Generated wrapper - only for numeric args
		"GreaterEqual(x_Number, y_Number)": WrapGreaterEqualExprs, // Generated wrapper - only for numeric args
		"SameQ(x_, y_)":                    WrapSameQExprs,        // Generated wrapper for same (structural equality)
		"UnsameQ(x_, y_)":                  WrapUnsameQExprs,      // Generated wrapper for unsame

		// Logical operations (Not - And/Or are special forms)
		"Not(x_)": WrapNotExpr, // Generated wrapper for logical negation

		// List operations
		"Length(x_)":              WrapLengthExpr,      // Generated wrapper for length calculation
		"First(x_List)":           WrapFirstExpr,       // Generated wrapper for first element access
		"Last(x_List)":            WrapLastExpr,        // Generated wrapper for last element access
		"Rest(x_List)":            WrapRestExpr,        // Generated wrapper for rest of list
		"Most(x_List)":            WrapMostExpr,        // Generated wrapper for most of list
		"Part(x_List, i_Integer)": WrapPartList,        // Generated wrapper for list part access
		"Part(x_Association, y_)": WrapPartAssociation, // Generated wrapper for association part access

		// Type predicates - using generated wrappers
		"IntegerQ(x_)":                  WrapIntegerQExpr,          // Generated wrapper for integer check
		"FloatQ(x_)":                    WrapFloatQExpr,            // Generated wrapper for float check
		"NumberQ(x_)":                   WrapNumberQExpr,           // Generated wrapper for number check
		"StringQ(x_)":                   WrapStringQExpr,           // Generated wrapper for string check
		"BooleanQ(x_)":                  WrapBooleanQExpr,          // Generated wrapper for boolean check
		"SymbolQ(x_)":                   WrapSymbolQExpr,           // Generated wrapper for symbol check
		"ListQ(x_)":                     WrapListQExpr,             // Generated wrapper for list check
		"AtomQ(x_)":                     WrapAtomQExpr,             // Generated wrapper for atom check
		"Head(x_)":                      WrapHeadExpr,              // Generated wrapper for head analysis
		"Attributes(x_)":                WrapAttributesExpr,        // Clean wrapper for symbol attributes
		"SetAttributes(x_, y_Symbol)":   WrapSetAttributesSingle,   // Clean wrapper for single attribute
		"SetAttributes(x_, y_List)":     WrapSetAttributesList,     // Clean wrapper for attribute list
		"ClearAttributes(x_, y_Symbol)": WrapClearAttributesSingle, // Clean wrapper for single attribute
		"ClearAttributes(x_, y_List)":   WrapClearAttributesList,   // Clean wrapper for attribute list
		"MatchQ(x_, y_)":                WrapMatchQExprs,           // Clean wrapper for pattern matching

		// String functions - using generated wrappers
		"StringLength(x_String)": WrapStringLengthStr, // Generated wrapper for string length (pattern-validated)

		// Output format functions - using generated wrappers
		"FullForm(x_)":  WrapFullFormExpr,  // Generated wrapper for full form output
		"InputForm(x_)": WrapInputFormExpr, // Generated wrapper for input form output

		// Association functions
		"Association(x___Rule)": WrapAssociationRules, // Generated wrapper for Association constructor (covers empty case)
		"AssociationQ(x_)":      WrapAssociationQExpr, // Generated wrapper for type check
		"Keys(x_Association)":   WrapKeysExpr,         // Generated wrapper for keys (Association-only)
		"Values(x_Association)": WrapValuesExpr,       // Generated wrapper for values (Association-only)
	}

	// Register all patterns
	err := registry.RegisterPatternBuiltins(builtinPatterns)
	if err != nil {
		panic(fmt.Sprintf("Failed to register built-in patterns: %v", err))
	}
}

// wrapBuiltinFunc wraps a builtin function to work with the new PatternFunc signature
func wrapBuiltinFunc(builtin func([]Expr) Expr) PatternFunc {
	return func(args []Expr, ctx *Context) Expr {
		// Check for errors in arguments and propagate them
		// Note: Stack frame addition happens in the caller (evaluatePatternFunction)
		for _, arg := range args {
			if IsError(arg) {
				return arg
			}
		}

		return builtin(args)
	}
}

// wrapBuiltinFuncNoErrorProp wraps a builtin function that should NOT propagate errors
// (e.g., Head should analyze error expressions, not propagate them)
func wrapBuiltinFuncNoErrorProp(builtin func([]Expr) Expr) PatternFunc {
	return func(args []Expr, ctx *Context) Expr {
		// No error propagation - let the builtin handle errors as data
		return builtin(args)
	}
}
