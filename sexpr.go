package sexpr

import (
	"fmt"
	"strconv"
	"strings"
)

type Expr interface {
	String() string
	Type() string
}

type AtomType int

const (
	StringAtom AtomType = iota
	IntAtom
	FloatAtom
	BoolAtom
	SymbolAtom
)

type Atom struct {
	AtomType AtomType
	Value    interface{}
}

func (a *Atom) String() string {
	switch a.AtomType {
	case StringAtom:
		return fmt.Sprintf("\"%s\"", a.Value.(string))
	case IntAtom:
		return strconv.Itoa(a.Value.(int))
	case FloatAtom:
		return strconv.FormatFloat(a.Value.(float64), 'f', -1, 64)
	case BoolAtom:
		if a.Value.(bool) {
			return "True"
		}
		return "False"
	case SymbolAtom:
		return a.Value.(string)
	default:
		return ""
	}
}

func (a *Atom) Type() string {
	switch a.AtomType {
	case StringAtom:
		return "string"
	case IntAtom:
		return "int"
	case FloatAtom:
		return "float64"
	case BoolAtom:
		return "bool"
	case SymbolAtom:
		return "symbol"
	default:
		return "unknown"
	}
}

type List struct {
	Elements []Expr
}

func (l *List) String() string {
	if len(l.Elements) == 0 {
		return "List()"
	}
	
	// Check if this is a List literal (head is "List")
	if len(l.Elements) > 0 {
		if headAtom, ok := l.Elements[0].(*Atom); ok && 
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

func (l *List) Type() string {
	return "list"
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

func (e *ErrorExpr) String() string {
	return fmt.Sprintf("$Failed(%s)", e.ErrorType)
}

func (e *ErrorExpr) Type() string {
	return "error"
}

// NewErrorExpr creates a new error expression
func NewErrorExpr(errorType, message string, args []Expr) *ErrorExpr {
	return &ErrorExpr{
		ErrorType:  errorType,
		Message:    message,
		Args:       args,
		StackTrace: []StackFrame{},
	}
}

// NewErrorExprWithStack creates a new error expression with stack trace
func NewErrorExprWithStack(errorType, message string, args []Expr, stack []StackFrame) *ErrorExpr {
	return &ErrorExpr{
		ErrorType:  errorType,
		Message:    message,
		Args:       args,
		StackTrace: stack,
	}
}

// GetStackTrace returns a formatted string representation of the stack trace
func (e *ErrorExpr) GetStackTrace() string {
	if len(e.StackTrace) == 0 {
		return "No stack trace available"
	}
	
	var trace strings.Builder
	trace.WriteString("Stack trace:\n")
	for i, frame := range e.StackTrace {
		trace.WriteString(fmt.Sprintf("  %d. %s: %s", i+1, frame.Function, frame.Expression))
		if frame.Location != "" {
			trace.WriteString(fmt.Sprintf(" at %s", frame.Location))
		}
		trace.WriteString("\n")
	}
	return trace.String()
}

// GetDetailedMessage returns the error message with stack trace
func (e *ErrorExpr) GetDetailedMessage() string {
	return fmt.Sprintf("Error: %s\nMessage: %s\n%s", e.ErrorType, e.Message, e.GetStackTrace())
}

// IsError checks if an expression is an error
func IsError(expr Expr) bool {
	_, ok := expr.(*ErrorExpr)
	return ok
}

func NewStringAtom(value string) *Atom {
	return &Atom{AtomType: StringAtom, Value: value}
}

func NewIntAtom(value int) *Atom {
	return &Atom{AtomType: IntAtom, Value: value}
}

func NewFloatAtom(value float64) *Atom {
	return &Atom{AtomType: FloatAtom, Value: value}
}

func NewBoolAtom(value bool) *Atom {
	return &Atom{AtomType: BoolAtom, Value: value}
}

func NewSymbolAtom(value string) *Atom {
	return &Atom{AtomType: SymbolAtom, Value: value}
}

func NewList(elements ...Expr) *List {
	return &List{Elements: elements}
}
