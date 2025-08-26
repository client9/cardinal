package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol With
// @ExprAttributes HoldAll

// @ExprPattern (_,_)
func With(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {

	vars := args[0]
	body := args[1]

	// we are assuming the vars are all [ x=x0, y=y0, ], i.e List( Set(x,x0), Set(y,y0) ...)
	// Replace functions expect a list of Rule
	// Copy List of Sets, and change heads to Rule
	// Call ReplaceAll in stdlib (purely mechanical change, not the ReplaceAll in Builts

	list, ok := vars.(core.List)
	if !ok {
		return core.NewError("ArgumentError", "With expected list for first argument")
	}
	rules := list.Copy()
	for _, arg := range rules.Tail() {
		r, ok := arg.(core.List)
		if !ok || r.Head() != "Set" || r.Length() != 2 {
			return core.NewError("ArgumentError", "With expected list of set assignments")
		}
		// TODO: DANGER
		r.SetHead("Rule")
		//(*r.Elements)[0] = core.SymbolFor(atom.Rule)
	}

	modified := core.ReplaceAllWithRules(body, rules)

	result := e.Evaluate(modified)

	return result
}
