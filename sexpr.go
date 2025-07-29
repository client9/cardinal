package sexpr

import (
	"github.com/client9/sexpr/core"
)

// Re-export core types for backward compatibility
type StackFrame = core.StackFrame
type ErrorExpr = core.ErrorExpr

// Re-export core constructor functions for backward compatibility
var (
	NewErrorExpr = core.NewErrorExpr
)
