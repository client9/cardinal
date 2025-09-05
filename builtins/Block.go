package builtins

import (
	"fmt"
	"log"

	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/core/symbol"
	"github.com/client9/cardinal/engine"
)

// @ExprSymbol Block
// @ExprAttributes HoldAll

// @ExprPattern (_,_)
func BlockExpr(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {

	vars := args[0]
	body := args[1]

	// Store original variable values to restore later
	savedVars, err := saveBlockVars(e, c, vars)
	if err != nil {
		log.Printf("SAVE BLOCK VARS: %s", err)
	}

	// Evaluate the body in the modified context
	result := e.Evaluate(body)

	restoreBlockVars(e, c, savedVars)

	return result
}

func restoreBlockVars(e *engine.Evaluator, ctx *engine.Context, saved map[core.Symbol]core.Expr) error {

	// Restore original values
	for varName, oldValue := range saved {
		// variable was specified, but didn't exist in context, delete it
		if oldValue == nil {
			ctx.Delete(varName)
			continue
		}
		if err := ctx.Set(varName, oldValue); err != nil {
			return err
		}
	}
	return nil
}

func saveBlockVars(e *engine.Evaluator, c *engine.Context, vars core.Expr) (map[core.Symbol]core.Expr, error) {
	// Store original variable values to restore later
	savedVars := make(map[core.Symbol]core.Expr)

	varList, ok := vars.(core.List)
	if !ok || varList.Length() == 0 {
		return nil, fmt.Errorf("first arg not a list")
	}
	// Expect List(Set(x, value), Set(y, value), ...)
	for _, arg := range varList.Tail() {
		// can be a single symbol name
		if varName, ok := arg.(core.Symbol); ok {
			// Save current value
			if oldValue, exists := c.Get(varName); exists {
				savedVars[varName] = oldValue
			} else {
				savedVars[varName] = nil
			}
			continue
		}

		setvar, ok := arg.(core.List)
		if !ok || setvar.Length() != 2 || setvar.Head() != symbol.Set {
			return nil, fmt.Errorf("variable not a symbol or assignment. %s, len=%d, head=%s", setvar.String(), setvar.Length(), setvar.String())
			// ERROR
		}
		setvarArgs := setvar.Tail()

		varName, ok := setvarArgs[0].(core.Symbol)
		if !ok {
			return nil, fmt.Errorf("Set malformed")
		}
		// Save old value
		if oldValue, exists := c.Get(varName); exists {
			savedVars[varName] = oldValue
		} else {
			// mark for delete at end
			savedVars[varName] = nil
		}
		// Set new value (evaluate the RHS) -- TODO ERROR CHECK
		newValue := e.Evaluate(setvarArgs[1])
		c.Set(varName, newValue)
	}
	return savedVars, nil
}
