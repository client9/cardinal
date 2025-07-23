package sexpr

import (
	"github.com/client9/sexpr/core"
)

// Re-export core types for backward compatibility
type Expr = core.Expr
type AtomType = core.AtomType
type Atom = core.Atom
type List = core.List
type StackFrame = core.StackFrame
type ErrorExpr = core.ErrorExpr
type ObjectExpr = core.ObjectExpr

// Re-export core constants
const (
	StringAtom = core.StringAtom
	IntAtom    = core.IntAtom
	FloatAtom  = core.FloatAtom
	SymbolAtom = core.SymbolAtom
)

// Re-export core constructor functions for backward compatibility
var (
	NewErrorExpr          = core.NewErrorExpr
	NewErrorExprWithStack = core.NewErrorExprWithStack
	NewStringAtom         = core.NewStringAtom
	NewIntAtom            = core.NewIntAtom
	NewFloatAtom          = core.NewFloatAtom
	NewBoolAtom           = core.NewBoolAtom
	NewSymbolAtom         = core.NewSymbolAtom
	NewList               = core.NewList
	NewObjectExpr         = core.NewObjectExpr
)

// Re-export core helper functions for backward compatibility
var (
	ExtractInt64    = core.ExtractInt64
	ExtractFloat64  = core.ExtractFloat64
	ExtractString   = core.ExtractString
	ExtractBool     = core.ExtractBool
	CopyExprList    = core.CopyExprList
	IsError         = core.IsError
	GetNumericValue = core.GetNumericValue
	IsNumeric       = core.IsNumeric
	IsBool          = core.IsBool
	IsSymbol        = core.IsSymbol
)

// Re-export functions needed by REPL
var (
	SetupBuiltinAttributes = setupBuiltinAttributes
)
