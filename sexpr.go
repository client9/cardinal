package sexpr

import (
	"github.com/client9/sexpr/core"
)

// Re-export core types for backward compatibility
type Expr = core.Expr
type List = core.List
type StackFrame = core.StackFrame
type ErrorExpr = core.ErrorExpr
type ObjectExpr = core.ObjectExpr

// Re-export core constructor functions for backward compatibility
var (
	NewErrorExpr          = core.NewErrorExpr
	NewErrorExprWithStack = core.NewErrorExprWithStack
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
