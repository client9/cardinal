package core

import (
	"fmt"
	"strconv"
	"strings"
)

// Expr is the fundamental interface for all expressions in the system
type Expr interface {
	String() string
	InputForm() string
	Type() string
	Length() int64
	Equal(rhs Expr) bool
}

// AtomType represents the type of an atomic value
type AtomType int

const (
	StringAtom AtomType = iota
	IntAtom
	FloatAtom
	SymbolAtom
)

// Atom represents atomic values (strings, integers, floats, symbols)
type Atom struct {
	AtomType AtomType
	Value    interface{}
}

func (a Atom) Length() int64 {
	if a.AtomType == StringAtom {
		return int64(len(a.Value.(string)))
	}
	return 0
}

func (a Atom) String() string {
	switch a.AtomType {
	case StringAtom:
		return fmt.Sprintf("\"%s\"", a.Value.(string))
	case IntAtom:
		return strconv.Itoa(a.Value.(int))
	case FloatAtom:
		return strconv.FormatFloat(a.Value.(float64), 'f', -1, 64)
	case SymbolAtom:
		return a.Value.(string)
	default:
		return ""
	}
}

func (a Atom) InputForm() string {
	// For atoms, InputForm is the same as String()
	return a.String()
}

func (a Atom) Type() string {
	switch a.AtomType {
	case StringAtom:
		return "string"
	case IntAtom:
		return "int"
	case FloatAtom:
		return "float64"
	case SymbolAtom:
		return "symbol"
	default:
		return "unknown"
	}
}

func (a Atom) Equal(rhs Expr) bool {
	rhsAtom, ok := rhs.(Atom)
	if !ok {
		return false
	}

	// Must have same atom type
	if a.AtomType != rhsAtom.AtomType {
		return false
	}

	// Compare values based on type
	switch a.AtomType {
	case StringAtom:
		return a.Value.(string) == rhsAtom.Value.(string)
	case IntAtom:
		return a.Value.(int) == rhsAtom.Value.(int)
	case FloatAtom:
		return a.Value.(float64) == rhsAtom.Value.(float64)
	case SymbolAtom:
		return a.Value.(string) == rhsAtom.Value.(string)
	default:
		return false
	}
}

// List represents compound expressions
type List struct {
	Elements []Expr
}

// Precedence levels for InputForm formatting
type Precedence int

const (
	PrecedenceLowest Precedence = iota
	PrecedenceAssign
	PrecedenceLogicalOr
	PrecedenceLogicalAnd
	PrecedenceEquality
	PrecedenceComparison
	PrecedenceSum
	PrecedenceProduct
)

func (l List) Length() int64 {
	// really should panic
	if len(l.Elements) == 0 {
		return 0
	}
	// element[0] is the head
	return int64(len(l.Elements)) - 1
}

func (l List) String() string {
	if len(l.Elements) == 0 {
		return "List()"
	}

	// Check if this is a List literal (head is "List")
	if len(l.Elements) > 0 {
		if headAtom, ok := l.Elements[0].(Atom); ok &&
			headAtom.AtomType == SymbolAtom && headAtom.Value.(string) == "List" {
			// This is a list literal: [element1, element2, ...]
			var elements []string
			for _, elem := range l.Elements[1:] {
				elements = append(elements, elem.String())
			}
			return fmt.Sprintf("List(%s)", strings.Join(elements, ", "))
		}
	}

	// This is a function call: head(arg1, arg2, ...)
	var elements []string
	for _, elem := range l.Elements {
		elements = append(elements, elem.String())
	}
	return fmt.Sprintf("%s(%s)", l.Elements[0].String(), strings.Join(elements[1:], ", "))
}

func (l List) InputForm() string {
	return l.inputFormWithPrecedence(PrecedenceLowest)
}

func (l List) inputFormWithPrecedence(parentPrecedence Precedence) string {
	if len(l.Elements) == 0 {
		return "List()"
	}

	// Check if this is a special function that has infix/shortcut representation
	if headAtom, ok := l.Elements[0].(Atom); ok && headAtom.AtomType == SymbolAtom {
		head := headAtom.Value.(string)

		switch head {
		case "List":
			// List(...) -> [...]
			if len(l.Elements) == 1 {
				return "[]"
			}
			var elements []string
			for _, elem := range l.Elements[1:] {
				elements = append(elements, elem.InputForm())
			}
			return fmt.Sprintf("[%s]", strings.Join(elements, ", "))

		case "Association":
			// Association(Rule(a,b), Rule(c,d)) -> {a: b, c: d}
			if len(l.Elements) == 1 {
				return "{}"
			}
			var pairs []string
			for _, elem := range l.Elements[1:] {
				if ruleList, ok := elem.(List); ok && len(ruleList.Elements) == 3 {
					if headAtom, ok := ruleList.Elements[0].(Atom); ok &&
						headAtom.AtomType == SymbolAtom && headAtom.Value.(string) == "Rule" {
						key := ruleList.Elements[1].InputForm()
						value := ruleList.Elements[2].InputForm()
						pairs = append(pairs, fmt.Sprintf("%s: %s", key, value))
						continue
					}
				}
				// Fallback for non-Rule elements
				pairs = append(pairs, elem.InputForm())
			}
			return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))

		case "Rule":
			// Rule(a, b) -> a: b
			if len(l.Elements) == 3 {
				return fmt.Sprintf("%s: %s", l.Elements[1].InputForm(), l.Elements[2].InputForm())
			}

		case "Set":
			// Set(a, b) -> a = b
			if len(l.Elements) == 3 {
				return l.formatInfixWithParens("=", PrecedenceAssign, parentPrecedence)
			}

		case "SetDelayed":
			// SetDelayed(a, b) -> a := b
			if len(l.Elements) == 3 {
				return l.formatInfixWithParens(":=", PrecedenceAssign, parentPrecedence)
			}

		case "Plus":
			// Plus(a, b, ...) -> a + b + ...
			if len(l.Elements) >= 3 {
				return l.formatLeftAssociativeInfix("+", PrecedenceSum, parentPrecedence)
			}

		case "Times":
			// Times(a, b, ...) -> a * b * ...
			if len(l.Elements) >= 3 {
				return l.formatLeftAssociativeInfix("*", PrecedenceProduct, parentPrecedence)
			}

		case "Subtract":
			// Subtract(a, b) -> a - b
			if len(l.Elements) == 3 {
				return l.formatInfixWithParens("-", PrecedenceSum, parentPrecedence)
			}

		case "Divide":
			// Divide(a, b) -> a / b
			if len(l.Elements) == 3 {
				return l.formatInfixWithParens("/", PrecedenceProduct, parentPrecedence)
			}

		case "Equal":
			// Equal(a, b) -> a == b
			if len(l.Elements) == 3 {
				return l.formatInfixWithParens("==", PrecedenceEquality, parentPrecedence)
			}

		case "Unequal":
			// Unequal(a, b) -> a != b
			if len(l.Elements) == 3 {
				return l.formatInfixWithParens("!=", PrecedenceEquality, parentPrecedence)
			}

		case "SameQ":
			// SameQ(a, b) -> a === b
			if len(l.Elements) == 3 {
				return l.formatInfixWithParens("===", PrecedenceEquality, parentPrecedence)
			}

		case "UnsameQ":
			// UnsameQ(a, b) -> a =!= b
			if len(l.Elements) == 3 {
				return l.formatInfixWithParens("=!=", PrecedenceEquality, parentPrecedence)
			}

		case "Less":
			// Less(a, b) -> a < b
			if len(l.Elements) == 3 {
				return l.formatInfixWithParens("<", PrecedenceComparison, parentPrecedence)
			}

		case "Greater":
			// Greater(a, b) -> a > b
			if len(l.Elements) == 3 {
				return l.formatInfixWithParens(">", PrecedenceComparison, parentPrecedence)
			}

		case "LessEqual":
			// LessEqual(a, b) -> a <= b
			if len(l.Elements) == 3 {
				return l.formatInfixWithParens("<=", PrecedenceComparison, parentPrecedence)
			}

		case "GreaterEqual":
			// GreaterEqual(a, b) -> a >= b
			if len(l.Elements) == 3 {
				return l.formatInfixWithParens(">=", PrecedenceComparison, parentPrecedence)
			}

		case "And":
			// And(a, b, ...) -> a && b && ...
			if len(l.Elements) >= 3 {
				return l.formatLeftAssociativeInfix("&&", PrecedenceLogicalAnd, parentPrecedence)
			}

		case "Or":
			// Or(a, b, ...) -> a || b || ...
			if len(l.Elements) >= 3 {
				return l.formatLeftAssociativeInfix("||", PrecedenceLogicalOr, parentPrecedence)
			}
		}
	}

	// Default: function call format Head(arg1, arg2, ...)
	var elements []string
	for _, elem := range l.Elements[1:] {
		elements = append(elements, elem.InputForm())
	}
	return fmt.Sprintf("%s(%s)", l.Elements[0].InputForm(), strings.Join(elements, ", "))
}

// formatInfixWithParens formats a binary infix operation with parentheses if needed
func (l List) formatInfixWithParens(op string, opPrecedence, parentPrecedence Precedence) string {
	left := l.getInputFormWithPrecedence(l.Elements[1], opPrecedence)
	right := l.getInputFormWithPrecedence(l.Elements[2], opPrecedence)
	result := fmt.Sprintf("%s %s %s", left, op, right)

	if opPrecedence < parentPrecedence {
		return fmt.Sprintf("(%s)", result)
	}
	return result
}

// formatLeftAssociativeInfix formats left-associative infix operations like a + b + c
func (l List) formatLeftAssociativeInfix(op string, opPrecedence, parentPrecedence Precedence) string {
	var parts []string
	for _, elem := range l.Elements[1:] {
		parts = append(parts, l.getInputFormWithPrecedence(elem, opPrecedence+1)) // Higher precedence for right operand
	}
	result := strings.Join(parts, fmt.Sprintf(" %s ", op))

	if opPrecedence < parentPrecedence {
		return fmt.Sprintf("(%s)", result)
	}
	return result
}

// getInputFormWithPrecedence gets InputForm with precedence context for proper parenthesization
func (l List) getInputFormWithPrecedence(expr Expr, precedence Precedence) string {
	if list, ok := expr.(List); ok {
		return list.inputFormWithPrecedence(precedence)
	}
	return expr.InputForm()
}

func (l List) Type() string {
	return "list"
}

func (l List) Equal(rhs Expr) bool {
	rhsList, ok := rhs.(List)
	if !ok {
		return false
	}

	// Lists must have same number of elements
	if len(l.Elements) != len(rhsList.Elements) {
		return false
	}

	// Recursively compare each element
	for i, elem := range l.Elements {
		if !elem.Equal(rhsList.Elements[i]) {
			return false
		}
	}

	return true
}

// StackFrame represents a single frame in the evaluation stack
type StackFrame struct {
	Function   string // Function name being evaluated
	Expression string // String representation of the expression
	Location   string // Optional location information
}

// ErrorExpr represents an error that occurred during evaluation
type ErrorExpr struct {
	ErrorType  string       // "DivisionByZero", "ArgumentError", etc.
	Message    string       // Detailed error message
	Args       []Expr       // Arguments that caused the error
	StackTrace []StackFrame // Stack trace of evaluation frames
}

// Possible we expose the arguments above differently and use ObjectExpr
// for now return 0
func (e *ErrorExpr) Length() int64 {
	return 0
}

func (e *ErrorExpr) String() string {
	return fmt.Sprintf("$Failed(%s)", e.ErrorType)
}

func (e *ErrorExpr) InputForm() string {
	// For errors, InputForm is the same as String()
	return e.String()
}

func (e *ErrorExpr) Type() string {
	return "error"
}

func (e *ErrorExpr) Equal(rhs Expr) bool {
	rhsError, ok := rhs.(*ErrorExpr)
	if !ok {
		return false
	}

	// Compare error type and message
	if e.ErrorType != rhsError.ErrorType || e.Message != rhsError.Message {
		return false
	}

	// Compare argument lists
	if len(e.Args) != len(rhsError.Args) {
		return false
	}

	for i, arg := range e.Args {
		if !arg.Equal(rhsError.Args[i]) {
			return false
		}
	}

	// Note: We don't compare stack traces as they are context-dependent
	return true
}

// ObjectExpr wraps a user-defined Expr implementation with a type name
// This allows users to register custom Go types that implement Expr
// and use them with pattern matching (e.g., x_Uint64)
type ObjectExpr struct {
	TypeName string // e.g., "Uint64", "BigInt", "Matrix"
	Value    Expr   // User-defined type that implements Expr interface
}

func (o ObjectExpr) String() string {
	return o.Value.String() // Delegate to the wrapped Expr
}

func (o ObjectExpr) InputForm() string {
	// Delegate to the wrapped Expr's InputForm if it has one,
	// otherwise fall back to String()
	return o.Value.InputForm()
}

func (o ObjectExpr) Length() int64 {
	return o.Value.Length() // Delegate to wrapper Expr
}

func (o ObjectExpr) Type() string {
	return o.TypeName // Return the registered type name
}

func (o ObjectExpr) Equal(rhs Expr) bool {
	rhsObj, ok := rhs.(ObjectExpr)
	if !ok || o.TypeName != rhsObj.TypeName {
		return false
	}
	return o.Value.Equal(rhsObj.Value) // Delegate to wrapped Expr
}
