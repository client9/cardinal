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

// PrintLabel outputs the expression with a label and returns it unchanged
func PrintLabel(label core.Expr, expr core.Expr) core.Expr {
	if labelAtom, ok := label.(core.Atom); ok && labelAtom.AtomType == core.StringAtom {
		fmt.Printf("%s: %s\n", labelAtom.Value.(string), expr.String())
	} else {
		fmt.Printf("%s: %s\n", label.String(), expr.String())
	}
	return expr
}
