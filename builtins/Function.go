package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/core/symbol"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Function
// @ExprAttributes HoldAll

// @ExprPattern (__)
func Function(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	switch len(args) {
	case 1:
		return FunctionPure(e, c, args[0])
	case 2:
		return FunctionNamed(e, c, args[0], args[1])
	}
	return core.NewError("ArgumentError", "Function expected 1 or 2 args")
}

func FunctionPure(e *engine.Evaluator, c *engine.Context, body core.Expr) core.Expr {
	mbody := partiallyEvaluateForFunction(e, c, body)
	return core.NewFunction(nil, mbody)
}

func varsToSymbolList(vars core.Expr) []core.Expr {
	if _, ok := core.ExtractSymbol(vars); ok {
		return []core.Expr{vars}
	}

	// convert sexpression to native slice
	if vlist, ok := vars.(core.List); ok {
		// could validate here that all are symbols
		return vlist.Tail()
	}

	return nil

}

func FunctionNamed(e *engine.Evaluator, c *engine.Context, vars, body core.Expr) core.Expr {
	vlist := varsToSymbolList(vars)
	mbody := partiallyEvaluateForFunction(e, c, body)
	return core.NewFunction(vlist, mbody)
}

// partiallyEvaluateForFunction evaluates nested Function calls but preserves slot variables
func partiallyEvaluateForFunction(e *engine.Evaluator, c *engine.Context, expr core.Expr) core.Expr {
	list, ok := expr.(core.List)
	if !ok {
		return expr
	}

	if list.HeadExpr() == symbol.Function {
		return Function(e, c, list.Tail())
	}

	largs := list.AsSlice()
	newElements := make([]core.Expr, len(largs))
	for i, elem := range largs {
		newElements[i] = partiallyEvaluateForFunction(e, c, elem)
	}
	return core.NewListFromExprs(newElements...)
}
