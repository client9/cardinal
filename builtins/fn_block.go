package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
	"fmt"	
	"log"
)

func BlockExpr(e *engine.Evaluator, c *engine.Context, vars core.Expr, body core.Expr) core.Expr {

	// Store original variable values to restore later
	savedVars, err := saveBlockVars(e,c,vars)
	if err != nil {
		log.Printf("SAVE BLOCK VARS: %s", err)
	}

	// Evaluate the body in the modified context
	result := e.Evaluate(c, body)

	restoreBlockVars(e,c, savedVars)

	return result
}

func restoreBlockVars(e *engine.Evaluator, ctx *engine.Context, saved map[string]core.Expr) error {

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

func saveBlockVars(e *engine.Evaluator, c *engine.Context, vars core.Expr) (map[string]core.Expr, error) {
	// Store original variable values to restore later
	savedVars := make(map[string]core.Expr)

	varList, ok := vars.(core.List)
	if !ok || varList.Length() == 0 {
		return nil, fmt.Errorf("first arg not a list")
	}
	// Expect List(Set(x, value), Set(y, value), ...)
	for i := 1; i < len(varList.Elements); i++ {
		arg := varList.Elements[i]

		// can be a single symbol name
		if varName, ok := core.ExtractSymbol(arg); ok {
			// Save current value
			if oldValue, exists := c.Get(varName); exists {
				savedVars[varName] = oldValue
			} else {
				savedVars[varName] = nil
			}
			continue
		}

		setvar, ok := arg.(core.List)
		if !ok || setvar.Length() != 2 || setvar.Head() != "Set" {
			return nil, fmt.Errorf("variable not a symbol or assignment. %s, len=%d, head=%s", setvar.String(), setvar.Length(), setvar.Head())
			// ERROR
		}

		varName, ok := core.ExtractSymbol(setvar.Elements[1])
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
		newValue := e.Evaluate(c, setvar.Elements[2])
		c.Set(varName, newValue)
	}
	return savedVars, nil
}
