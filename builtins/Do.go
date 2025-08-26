package builtins

import (
	"fmt"

	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol Do
// @ExprAttributes HoldAll

// evaluateDo implements the Do special form for iteration without collecting results
// Do(expr, n) evaluates expr n times and returns Null
// Do(expr, core.List(i, start, end, increment)) iterates with variable binding
//
// @ExprPattern (_,_)
func Do(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	if len(args) != 2 {
		return core.NewError("ArgumentError", fmt.Sprintf("Do expects 2 arguments, got %d", len(args)))
	}

	expr := args[0] // Don't evaluate expr yet - Do has HoldAll
	spec := args[1] // Don't evaluate spec yet

	// if list, assume iterator spec
	if spec.Head() == "List" {
		return doIterator(e, c, expr, spec.(core.List))
	}

	// it's something else, evaluate it.
	val := e.Evaluate(spec)
	if core.IsError(val) {
		return val
	}

	// Check if second argument is an integer (simple replication form)
	// TODO
	// Want "GetNumericInt" if int64, return int64, if float64 return int64
	if n, ok := core.GetNumericValue(val); ok {
		return doSimple(e, c, expr, int64(n))
	}

	return core.NewError("ArgumentError", "Do second argument must be integer or core.List")
}

// evaluateDoSimple implements Do(expr, n) - evaluates expr n times and returns Null
func doSimple(e *engine.Evaluator, c *engine.Context, expr core.Expr, n int64) core.Expr {
	if n < 0 {
		return core.NewError("ArgumentError", fmt.Sprintf("Do count must be non-negative, got %d", n))
	}

	// Evaluate expr n times (side effects only)
	for i := int64(0); i < n; i++ {
		result := e.Evaluate(expr)
		if core.IsError(result) {
			return result // Return error immediately
		}
		// Discard result - Do is for side effects only
	}

	return core.NewSymbolNull()
}

// evaluateDoIterator handles Do with iterator specification core.List(i, start, end, increment)
func doIterator(e *engine.Evaluator, c *engine.Context, expr core.Expr, iterSpec core.List) core.Expr {
	// Parse iterator specification into normalized form
	variable, start, end, increment, err := parseTableIteratorSpec(e, c, iterSpec)
	if err != nil {
		return err
	}

	current := start
	const maxIterations = 10000 // Prevent infinite loops

	for iteration := 0; iteration < maxIterations; iteration++ {
		// Check if we should continue iterating
		shouldContinue := evaluateIteratorCondition(e, c, current, end, increment)
		if !shouldContinue {
			break
		}

		// Evaluate expression with current iterator value (for side effects only)
		blockResult := evaluateWithIteratorBinding(e, c, expr, variable, current)
		if core.IsError(blockResult) {
			return blockResult // Return error immediately
		}
		// Discard result - Do is for side effects only

		// Increment for next iteration
		current = evaluateIteratorIncrement(e, c, current, increment)
		if core.IsError(current) {
			return current
		}
	}

	return core.NewSymbolNull()
}

/*
// DoExpr executes an expression multiple times: Do(expr, {i, n})
func DoExpr(e *engine.Evaluator, ctx *engine.Context, expr core.Expr, iterator core.Expr) core.Expr {
	// Handle simple count: Do(expr, n)
	if intVal, ok := iterator.(core.Integer); ok {
		count := int(intVal)
		for i := 0; i < count; i++ {
			e.Evaluate(ctx, expr)
		}
		return core.NewSymbolNull()
	}

	// Handle iterator: List(i, n) or List(i, start, end)
	if iterList, ok := iterator.(core.List); ok && len(iterList.Elements) >= 2 {
		if symbolName, ok := core.ExtractSymbol(iterList.Elements[0]); ok && symbolName == "List" {
			if len(iterList.Elements) == 3 {
				// List(i, n) - iterate from 1 to n
				if varName, ok := core.ExtractSymbol(iterList.Elements[1]); ok {
					if countVal, ok := iterList.Elements[2].(core.Integer); ok {
						count := int(countVal)

						// Save original variable value
						var savedValue core.Expr
						var hadValue bool
						if oldValue, exists := ctx.Get(varName); exists {
							savedValue = oldValue
							hadValue = true
						}

						// Iterate from 1 to count
						for i := 1; i <= count; i++ {
							ctx.Set(varName, core.Integer(i))
							e.Evaluate(ctx, expr)
						}

						// Restore original value
						if hadValue {
							ctx.Set(varName, savedValue)
						} else {
							ctx.Delete(varName)
						}
					}
				}
			} else if len(iterList.Elements) == 4 {
				// List(i, start, end) - iterate from start to end
				if varName, ok := core.ExtractSymbol(iterList.Elements[1]); ok {
					if startVal, ok := iterList.Elements[2].(core.Integer); ok {
						if endVal, ok := iterList.Elements[3].(core.Integer); ok {
							start := int(startVal)
							end := int(endVal)

							// Save original variable value
							var savedValue core.Expr
							var hadValue bool
							if oldValue, exists := ctx.Get(varName); exists {
								savedValue = oldValue
								hadValue = true
							}

							// Iterate from start to end
							for i := start; i <= end; i++ {
								ctx.Set(varName, core.Integer(i))
								e.Evaluate(ctx,expr)
							}

							// Restore original value
							if hadValue {
								ctx.Set(varName, savedValue)
							} else {
								ctx.Delete(varName)
							}
						}
					}
				}
			}
		}
	}

	return core.NewSymbolNull()
}
*/
