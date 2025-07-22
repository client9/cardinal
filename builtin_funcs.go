package sexpr

import (
	"fmt"
	"math"
	"unicode/utf8"
)

//go:generate go run cmd/wrapgen/main.go -single builtin_wrappers.go

// Business logic functions for mathematical operations
// These are pure Go functions that are automatically wrapped

// PlusIntegers adds a sequence of integers
func PlusIntegers(args ...int64) int64 {
	sum := int64(0)
	for _, v := range args {
		sum += v
	}
	return sum
}

// TimesIntegers multiplies a sequence of integers
func TimesIntegers(args ...int64) int64 {
	if len(args) == 0 {
		return 1
	}
	product := int64(1)
	for _, v := range args {
		product *= v
	}
	return product
}

// PlusReals adds a sequence of real numbers
func PlusReals(args ...float64) float64 {
	sum := 0.0
	for _, v := range args {
		sum += v
	}
	return sum
}

// TimesReals multiplies a sequence of real numbers
func TimesReals(args ...float64) float64 {
	if len(args) == 0 {
		return 1.0
	}
	product := 1.0
	for _, v := range args {
		product *= v
	}
	return product
}

// Mixed numeric arithmetic - handles both integers and floats, returns float64

// PlusNumbers adds a sequence of mixed numeric expressions
func PlusNumbers(args ...Expr) float64 {
	sum := 0.0
	for _, arg := range args {
		if val, ok := getNumericValue(arg); ok {
			sum += val
		}
		// Skip non-numeric values - they'll be caught by pattern matching
	}
	return sum
}

// TimesNumbers multiplies a sequence of mixed numeric expressions
func TimesNumbers(args ...Expr) float64 {
	if len(args) == 0 {
		return 1.0
	}
	product := 1.0
	for _, arg := range args {
		if val, ok := getNumericValue(arg); ok {
			product *= val
		}
		// Skip non-numeric values - they'll be caught by pattern matching
	}
	return product
}

// EqualExprs checks if two expressions are equal
func EqualExprs(x, y Expr) bool {
	return x.Equal(y)
}

// PowerReal computes real base to integer exponent
func PowerReal(base float64, exp int64) float64 {
	if exp == 0 {
		return 1.0
	}
	if exp < 0 {
		return 1.0 / PowerReal(base, -exp)
	}

	result := 1.0
	for i := int64(0); i < exp; i++ {
		result *= base
	}
	return result
}

// StringLengthFunc returns the length of a string
func StringLengthFunc(s string) int64 {
	return int64(len(s))
}

// Type predicate functions - all return bool

// IntegerQExpr checks if an expression is an integer
func IntegerQExpr(expr Expr) bool {
	if atom, ok := expr.(Atom); ok {
		return atom.AtomType == IntAtom
	}
	return false
}

// FloatQExpr checks if an expression is a float
func FloatQExpr(expr Expr) bool {
	if atom, ok := expr.(Atom); ok {
		return atom.AtomType == FloatAtom
	}
	return false
}

// NumberQExpr checks if an expression is numeric (int or float)
func NumberQExpr(expr Expr) bool {
	return isNumeric(expr)
}

// StringQExpr checks if an expression is a string
func StringQExpr(expr Expr) bool {
	if atom, ok := expr.(Atom); ok {
		return atom.AtomType == StringAtom
	}
	return false
}

// BooleanQExpr checks if an expression is a boolean (True/False symbol)
func BooleanQExpr(expr Expr) bool {
	return isBool(expr)
}

// SymbolQExpr checks if an expression is a symbol
func SymbolQExpr(expr Expr) bool {
	return isSymbol(expr)
}

// ListQExpr checks if an expression is a list
func ListQExpr(expr Expr) bool {
	_, isList := expr.(List)
	return isList
}

// AtomQExpr checks if an expression is an atom
func AtomQExpr(expr Expr) bool {
	_, isAtom := expr.(Atom)
	return isAtom
}

// Output format functions - all return string

// FullFormExpr returns the full string representation of an expression
func FullFormExpr(expr Expr) string {
	// Convert pattern strings to their symbolic form before getting string representation
	if atom, ok := expr.(Atom); ok && atom.AtomType == SymbolAtom {
		symbolStr := atom.Value.(string)
		if isPatternVariable(symbolStr) {
			// Convert pattern string to symbolic form
			symbolicExpr := ConvertPatternStringToSymbolic(symbolStr)
			return symbolicExpr.String()
		}
	}

	// Return the string representation
	return expr.String()
}

// InputFormExpr returns the user-friendly InputForm representation of an expression
func InputFormExpr(expr Expr) string {
	// Convert pattern strings to their symbolic form before getting InputForm representation
	if atom, ok := expr.(Atom); ok && atom.AtomType == SymbolAtom {
		symbolStr := atom.Value.(string)
		if isPatternVariable(symbolStr) {
			// Convert pattern string to symbolic form
			symbolicExpr := ConvertPatternStringToSymbolic(symbolStr)
			return symbolicExpr.InputForm()
		}
	}

	// Return the InputForm representation
	return expr.InputForm()
}

// Length and string functions - all return int64

// LengthExpr returns the length of an expression
func LengthExpr(expr Expr) int64 {
	switch ex := expr.(type) {
	case List:
		// For lists, return the number of elements (excluding the head)
		if len(ex.Elements) == 0 {
			return 0 // Empty list has length 0
		}
		return int64(len(ex.Elements) - 1) // Subtract 1 for the head
	case ObjectExpr:
		// Handle Association
		if ex.TypeName == "Association" {
			if assocValue, ok := ex.Value.(AssociationValue); ok {
				return int64(assocValue.Len())
			}
		}
		// For other ObjectExpr types, length is 0
		return 0
	default:
		// For atoms and other expressions, length is 0
		return 0
	}
}

// StringLengthStr returns the UTF-8 rune count of a string
func StringLengthStr(s string) int64 {
	return int64(utf8.RuneCountInString(s))
}

// List access functions - all return Expr

// FirstExpr returns the first element of a list (after the head)
func FirstExpr(list List) Expr {
	// For lists, return the first element after the head
	if len(list.Elements) <= 1 {
		// Return error for empty lists
		return NewErrorExpr("PartError",
			fmt.Sprintf("First: expression %s has no elements", list.String()), []Expr{list})
	}
	return list.Elements[1] // Index 1 is first element after head (index 0)
}

// LastExpr returns the last element of a list
func LastExpr(list List) Expr {
	// For lists, return the last element
	if len(list.Elements) <= 1 {
		// Return error for empty lists
		return NewErrorExpr("PartError",
			fmt.Sprintf("Last: expression %s has no elements", list.String()), []Expr{list})
	}
	return list.Elements[len(list.Elements)-1] // Last element
}

// RestExpr returns a new list with the first element after head removed
func RestExpr(list List) Expr {
	// For lists, return a new list with the first element after head removed
	if len(list.Elements) <= 1 {
		// Return error for empty lists
		return NewErrorExpr("PartError",
			fmt.Sprintf("Rest: expression %s has no elements", list.String()), []Expr{list})
	}

	// Create new list: head + elements[2:] (skip first element after head)
	if len(list.Elements) == 2 {
		// Special case: if only head and one element, return just the head
		return List{Elements: []Expr{list.Elements[0]}}
	}

	newElements := make([]Expr, len(list.Elements)-1)
	newElements[0] = list.Elements[0]        // Keep the head
	copy(newElements[1:], list.Elements[2:]) // Copy everything after the first element
	return List{Elements: newElements}
}

// MostExpr returns a new list with the last element removed
func MostExpr(list List) Expr {
	// For lists, return a new list with the last element removed
	if len(list.Elements) <= 1 {
		// Return error for empty lists
		return NewErrorExpr("PartError",
			fmt.Sprintf("Most: expression %s has no elements", list.String()), []Expr{list})
	}

	// Create new list with all elements except the last one
	if len(list.Elements) == 2 {
		// Special case: if only head and one element, return just the head
		return List{Elements: []Expr{list.Elements[0]}}
	}

	newElements := make([]Expr, len(list.Elements)-1)
	copy(newElements, list.Elements[:len(list.Elements)-1])
	return List{Elements: newElements}
}

// Comparison operators - all return bool

// UnequalExprs checks if two expressions are not equal
func UnequalExprs(x, y Expr) bool {
	return !x.Equal(y)
}

// LessExprs checks if x < y for numeric types
func LessExprs(x, y Expr) bool {
	// Extract numeric values
	val1, ok1 := getNumericValue(x)
	val2, ok2 := getNumericValue(y)
	if ok1 && ok2 {
		return val1 < val2
	}
	return false // Fallback case - will be handled by wrapper
}

// GreaterExprs checks if x > y for numeric types
func GreaterExprs(x, y Expr) bool {
	// Extract numeric values
	val1, ok1 := getNumericValue(x)
	val2, ok2 := getNumericValue(y)
	if ok1 && ok2 {
		return val1 > val2
	}
	return false // Fallback case - will be handled by wrapper
}

// LessEqualExprs checks if x <= y for numeric types
func LessEqualExprs(x, y Expr) bool {
	// Extract numeric values
	val1, ok1 := getNumericValue(x)
	val2, ok2 := getNumericValue(y)
	if ok1 && ok2 {
		return val1 <= val2
	}
	return false // Fallback case - will be handled by wrapper
}

// GreaterEqualExprs checks if x >= y for numeric types
func GreaterEqualExprs(x, y Expr) bool {
	// Extract numeric values
	val1, ok1 := getNumericValue(x)
	val2, ok2 := getNumericValue(y)
	if ok1 && ok2 {
		return val1 >= val2
	}
	return false // Fallback case - will be handled by wrapper
}

// SameQExprs checks if two expressions are structurally equal
func SameQExprs(x, y Expr) bool {
	return x.Equal(y)
}

// UnsameQExprs checks if two expressions are not structurally equal
func UnsameQExprs(x, y Expr) bool {
	return !x.Equal(y)
}

// Association functions - all work with ObjectExpr of type "Association"

// AssociationQExpr checks if an expression is an Association
func AssociationQExpr(expr Expr) bool {
	if objExpr, ok := expr.(ObjectExpr); ok && objExpr.TypeName == "Association" {
		return true
	}
	return false
}

// KeysExpr returns the keys of an association as a List
func KeysExpr(assoc ObjectExpr) Expr {
	if assoc.TypeName == "Association" {
		if assocValue, ok := assoc.Value.(AssociationValue); ok {
			keys := assocValue.Keys()
			// Return as List[key1, key2, ...]
			elements := []Expr{NewSymbolAtom("List")}
			elements = append(elements, keys...)
			return NewList(elements...)
		}
	}
	// This should not happen with proper pattern matching, but return error as fallback
	return NewErrorExpr("ArgumentError",
		fmt.Sprintf("Keys expects an Association, got %s", assoc.String()), []Expr{assoc})
}

// ValuesExpr returns the values of an association as a List
func ValuesExpr(assoc ObjectExpr) Expr {
	if assoc.TypeName == "Association" {
		if assocValue, ok := assoc.Value.(AssociationValue); ok {
			values := assocValue.Values()
			// Return as List[value1, value2, ...]
			elements := []Expr{NewSymbolAtom("List")}
			elements = append(elements, values...)
			return NewList(elements...)
		}
	}
	// This should not happen with proper pattern matching, but return error as fallback
	return NewErrorExpr("ArgumentError",
		fmt.Sprintf("Values expects an Association, got %s", assoc.String()), []Expr{assoc})
}

// AssociationRules creates an Association from a sequence of Rule expressions
func AssociationRules(rules ...Expr) Expr {
	assocValue := NewAssociationValue()

	// Process each Rule expression
	for _, rule := range rules {
		if ruleList, ok := rule.(List); ok && len(ruleList.Elements) == 3 {
			// Check if this is Rule[key, value]
			if headAtom, ok := ruleList.Elements[0].(Atom); ok &&
				headAtom.AtomType == SymbolAtom && headAtom.Value.(string) == "Rule" {

				key := ruleList.Elements[1]
				value := ruleList.Elements[2]
				assocValue = assocValue.Set(key, value) // Returns new association (immutable)
				continue
			}
		}

		// Invalid argument - not a Rule
		return NewErrorExpr("ArgumentError",
			fmt.Sprintf("Association expects Rule expressions, got %s", rule.String()), []Expr{rule})
	}

	return NewObjectExpr("Association", assocValue)
}

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

// NotExpr performs logical negation on boolean expressions
func NotExpr(expr Expr) Expr {
	if isBool(expr) {
		val, _ := getBoolValue(expr)
		return NewBoolAtom(!val)
	}

	// Return unchanged expression if not boolean (symbolic behavior)
	return List{Elements: []Expr{NewSymbolAtom("Not"), expr}}
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

// Arithmetic functions - type-specific operations

// SubtractIntegers performs integer subtraction
func SubtractIntegers(x, y int64) int64 {
	return x - y
}

// SubtractNumbers performs mixed numeric subtraction (returns float64)
// Pattern constraint ensures both arguments are numeric
func SubtractNumbers(x, y Expr) float64 {
	val1, _ := getNumericValue(x)
	val2, _ := getNumericValue(y)
	return val1 - val2
}

// PowerNumbers performs power operation on numeric arguments
// Returns (float64, error) for clear type safety
func PowerNumbers(base, exp Expr) (float64, error) {
	if !isNumeric(base) || !isNumeric(exp) {
		return 0, fmt.Errorf("MathematicalError")
	}

	baseVal, _ := getNumericValue(base)
	expVal, _ := getNumericValue(exp)
	result := math.Pow(baseVal, expVal)

	// Check for invalid results (NaN, Inf)
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, fmt.Errorf("MathematicalError")
	}

	return result, nil
}

// DivideIntegers performs integer division on int64 arguments
// Returns (int64, error) for clear type safety using Go's integer division
func DivideIntegers(x, y int64) (int64, error) {
	if y == 0 {
		return 0, fmt.Errorf("DivisionByZero")
	}

	return x / y, nil
}

// DivideNumbers performs division on numeric arguments
// Returns (float64, error) for clear type safety
func DivideNumbers(x, y Expr) (float64, error) {
	if !isNumeric(x) || !isNumeric(y) {
		return 0, fmt.Errorf("MathematicalError")
	}

	val1, _ := getNumericValue(x)
	val2, _ := getNumericValue(y)

	if val2 == 0 {
		return 0, fmt.Errorf("DivisionByZero")
	}

	return val1 / val2, nil
}

// Identity functions for empty arithmetic operations (direct pattern mappings)

// PlusEmpty returns the additive identity (0) for empty Plus()
func PlusEmpty() Expr {
	return NewIntAtom(0)
}

// TimesEmpty returns the multiplicative identity (1) for empty Times()
func TimesEmpty() Expr {
	return NewIntAtom(1)
}

// Part functions - separated for Lists and Associations

// PartList extracts an element from a list by integer index (1-based)
func PartList(list List, index int64) Expr {
	// For lists, return the element at the specified index (1-based)
	if len(list.Elements) <= 1 {
		return NewErrorExpr("PartError",
			fmt.Sprintf("Part: expression %s has no elements", list.String()), []Expr{list})
	}

	// Handle negative indexing: -1 is last element, -2 is second to last, etc.
	var actualIndex int
	if index < 0 {
		// Negative indexing: -1 = last, -2 = second to last, etc.
		actualIndex = len(list.Elements) + int(index)
	} else if index > 0 {
		// Positive 1-based indexing: convert to 0-based for internal use
		actualIndex = int(index)
	} else {
		// index == 0 is invalid in 1-based indexing
		return NewErrorExpr("PartError",
			fmt.Sprintf("Part index %d is out of bounds (indices start at 1)", index), []Expr{list})
	}

	// Check bounds (remember: list.Elements[0] is the head, actual elements start at index 1)
	if actualIndex < 1 || actualIndex >= len(list.Elements) {
		return NewErrorExpr("PartError",
			fmt.Sprintf("Part index %d is out of bounds for expression with %d elements",
				index, len(list.Elements)-1), []Expr{list})
	}

	return list.Elements[actualIndex]
}

// PartAssociation extracts a value from an association by key
func PartAssociation(assoc ObjectExpr, key Expr) Expr {
	if assoc.TypeName == "Association" {
		if assocValue, ok := assoc.Value.(AssociationValue); ok {
			// For associations, use the key argument to lookup value
			if value, exists := assocValue.Get(key); exists {
				return value
			}
			return NewErrorExpr("PartError",
				fmt.Sprintf("Key %s not found in association", key.String()), []Expr{assoc, key})
		}
	}
	// This should not happen with proper pattern matching, but return error as fallback
	return NewErrorExpr("ArgumentError",
		fmt.Sprintf("Part expects an Association, got %s", assoc.String()), []Expr{assoc})
}

// HeadExpr returns the head/type of an expression
func HeadExpr(expr Expr) Expr {
	switch ex := expr.(type) {
	case Atom:
		switch ex.AtomType {
		case IntAtom:
			return NewSymbolAtom("Integer")
		case FloatAtom:
			return NewSymbolAtom("Real")
		case StringAtom:
			return NewSymbolAtom("String")
		case SymbolAtom:
			return NewSymbolAtom("Symbol")
		default:
			return NewSymbolAtom("Unknown")
		}
	case List:
		if len(ex.Elements) == 0 {
			return NewSymbolAtom("List")
		} else {
			// For non-empty lists, the head is the first element
			// This matches Mathematica semantics where f[x,y] has head f
			return ex.Elements[0]
		}
	case ObjectExpr:
		return NewSymbolAtom(ex.TypeName)
	default:
		return NewSymbolAtom("Unknown")
	}
	// Note: ErrorExpr is not handled here - wrapper will propagate errors
}
