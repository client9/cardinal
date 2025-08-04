package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

func With(e *engine.Evaluator, c *engine.Context, vars core.Expr, body core.Expr) core.Expr {

	// we are assuming the vars are all [ x=x0, y=y0, ], i.e List( Set(x,x0), Set(y,y0) ...)
	// Replace functions expect a list of Rule
	// Copy List of Sets, and change heads to Rule
	// Call ReplaceAll in stdlib (purely mechanical change, not the ReplaceAll in Builts

	list, ok := vars.(core.List)
	if !ok {
		return core.NewErrorExpr("ArgumentError", "With expected list for first argument", []core.Expr{vars})
	}
	rules := list.Copy()
	for i := int64(1); i <= rules.Length(); i++ {
		r, ok := rules.Elements[i].(core.List)
		if !ok || r.Head() != "Set" || r.Length() != 2 {
			return core.NewErrorExpr("ArgumentError", "With expected list of set assignments", []core.Expr{vars})
		}
		r.Elements[0] = core.NewSymbol("Rule")
	}

	modified := core.ReplaceAllWithRules(body, rules)

	result := e.Evaluate(c, modified)

	return result
}
