package builtins

import (
	"github.com/client9/cardinal/core"
)

// ListQ checks if an expression is a list
func ListQ(expr core.Expr) core.Expr {
	_, ok := expr.(core.List)
	return core.NewBool(ok)
}
