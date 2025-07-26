package stdlib

import (
	"fmt"
	"github.com/client9/sexpr/core"
)

// Print outputs the expression and returns it unchanged
// This allows debugging intermediate values in compound statements
func Print(expr core.Expr) core.Expr {
	fmt.Println(expr.String())
	return expr
}
