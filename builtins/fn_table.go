package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"

	"fmt"
)

func Table(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	if len(args) != 2 {
		return core.NewErrorExpr("ArgumentError", "Table expects 2 arguments", args)
	}

	expr := args[0] // Don't evaluate expr yet - Table has HoldAll
	spec := args[1] // Don't evaluate spec yet

	// Check if second argument is an integer (simple replication form)
	if n, ok := core.ExtractInt64(spec); ok {
		return tableSimple(e, c, expr, n)
	}

	// Check if second argument is a core.List (iterator form)
	if iterList, ok := spec.(core.List); ok {
		return tableIterator(e, c, expr, iterList)
	}

	return core.NewErrorExpr("ArgumentError", "Table second argument must be integer or core.List", args)
}

// evaluateTableSimple implements Table(expr, n) - creates n copies of expr
func tableSimple(e *engine.Evaluator, c *engine.Context, expr core.Expr, n int64) core.Expr {
	if n < 0 {
		return core.NewErrorExpr("ArgumentError", fmt.Sprintf("Table count must be non-negative, got %d", n), []core.Expr{core.NewInteger(n)})
	}

	if n == 0 {
		return core.NewList("List")
	}

	// Create result list with proper capacity
	elements := make([]core.Expr, n)

	// Evaluate expr once for each position
	for i := 0; i < int(n); i++ {
		// Evaluate expr in current context for each iteration
		// This allows expressions with side effects to work correctly
		evaluated := e.Evaluate(c, expr)
		if core.IsError(evaluated) {
			return evaluated
		}
		elements[i] = evaluated
	}

	return core.NewList("List", elements...)
}

// evaluateTableIterator implements Table(expr, core.List(i, start, end, increment))
// Handles all iterator forms using the general case with expression-based arithmetic
func tableIterator(e *engine.Evaluator, c *engine.Context, expr core.Expr, iterSpec core.List) core.Expr {
	// Parse iterator specification into normalized form
	variable, start, end, increment, err := parseTableIteratorSpec(e, c, iterSpec)
	if err != nil {
		return err
	}

	var results []core.Expr

	current := start
	const maxIterations = 10000 // Prevent infinite loops

	for iteration := 0; iteration < maxIterations; iteration++ {
		// Check if we should continue iterating
		shouldContinue := evaluateIteratorCondition(e, c, current, end, increment)
		if !shouldContinue {
			break
		}

		// Use Block to bind iterator variable and evaluate expression
		blockResult := evaluateWithIteratorBinding(e, c, expr, variable, current)
		if core.IsError(blockResult) {
			return blockResult
		}
		results = append(results, blockResult)

		// Increment current value using expression arithmetic
		current = evaluateIteratorIncrement(e, c, current, increment)
		if core.IsError(current) {
			return current
		}
	}

	return core.NewList("List", results...)
}

// parseTableIteratorSpec parses iterator specifications and normalizes them
// core.List(i, max) → core.List(i, 1, max, 1)
// core.List(i, start, end) → core.List(i, start, end, 1)
// core.List(i, start, end, increment) → core.List(i, start, end, increment)
// IMPORTANT: Evaluates start, end, and increment expressions and validates they are numeric
func parseTableIteratorSpec(e *engine.Evaluator, c *engine.Context, iterSpec core.List) (variable string, start, end, increment core.Expr, err core.Expr) {
	if len(iterSpec.Elements) < 3 || len(iterSpec.Elements) > 5 {
		return "", nil, nil, nil, core.NewErrorExpr("ArgumentError",
			"Table iterator must be core.List(var, max), core.List(var, start, end), or core.List(var, start, end, step)", []core.Expr{iterSpec})
	}

	// Extract variable name
	if varSymbol, ok := core.ExtractSymbol(iterSpec.Elements[1]); ok {
		variable = varSymbol
	} else {
		return "", nil, nil, nil, core.NewErrorExpr("ArgumentError", "Table iterator variable must be a symbol", []core.Expr{iterSpec.Elements[1]})
	}

	// Parse and evaluate based on number of arguments
	switch len(iterSpec.Elements) {
	case 3: // core.List(i, max) → core.List(i, 1, max, 1)
		start = core.NewInteger(1)
		end = e.Evaluate(c, iterSpec.Elements[2])
		if core.IsError(end) {
			return "", nil, nil, nil, end
		}
		increment = core.NewInteger(1)

	case 4: // core.List(i, start, end) → core.List(i, start, end, 1)
		start = e.Evaluate(c, iterSpec.Elements[2])
		if core.IsError(start) {
			return "", nil, nil, nil, start
		}
		end = e.Evaluate(c, iterSpec.Elements[3])
		if core.IsError(end) {
			return "", nil, nil, nil, end
		}
		increment = core.NewInteger(1)

	case 5: // core.List(i, start, end, increment)
		start = e.Evaluate(c, iterSpec.Elements[2])
		if core.IsError(start) {
			return "", nil, nil, nil, start
		}
		end = e.Evaluate(c, iterSpec.Elements[3])
		if core.IsError(end) {
			return "", nil, nil, nil, end
		}
		increment = e.Evaluate(c, iterSpec.Elements[4])
		if core.IsError(increment) {
			return "", nil, nil, nil, increment
		}

	default:
		return "", nil, nil, nil, core.NewErrorExpr("ArgumentError", "Invalid Table iterator specification", []core.Expr{iterSpec})
	}

	// Validate that arithmetic and comparison operations can be evaluated
	// Test if Plus(start, increment) evaluates to something different (not unevaluated)
	testPlus := core.NewList("Plus", start, increment)
	plusResult := e.Evaluate(c, testPlus)
	if core.IsError(plusResult) {
		return "", nil, nil, nil, core.NewErrorExpr("ArgumentError",
			fmt.Sprintf("Table iterator arithmetic failed: %s", plusResult), []core.Expr{plusResult})
	}
	if plusResult.Equal(testPlus) && !plusResult.IsAtom() {
		return "", nil, nil, nil, core.NewErrorExpr("ArgumentError",
			fmt.Sprintf("Table iterator arithmetic unevaluated: Plus(%s, %s) - missing arithmetic definition", start, increment), []core.Expr{start, increment})
	}

	// Test if comparison operation evaluates
	testLessEqual := core.NewList("LessEqual", start, end)
	compareResult := e.Evaluate(c, testLessEqual)
	if core.IsError(compareResult) {
		return "", nil, nil, nil, core.NewErrorExpr("ArgumentError",
			fmt.Sprintf("Table iterator comparison failed: %s", compareResult), []core.Expr{compareResult})
	}
	if compareResult.Equal(testLessEqual) && !compareResult.IsAtom() {
		return "", nil, nil, nil, core.NewErrorExpr("ArgumentError",
			fmt.Sprintf("Table iterator comparison unevaluated: LessEqual(%s, %s) - missing comparison definition", start, end), []core.Expr{start, end})
	}

	return variable, start, end, increment, nil
}

// evaluateIteratorCondition determines if iteration should continue
// Uses expression-based comparison with proper handling of increment direction
func evaluateIteratorCondition(e *engine.Evaluator, c *engine.Context, current, end, increment core.Expr) bool {
	// Determine comparison operator based on increment sign
	var compSymbol string
	isNegative := isNegativeIncrement(e, c, increment)
	if isNegative {
		compSymbol = "GreaterEqual" // For negative increment, continue while current >= end
	} else {
		compSymbol = "LessEqual" // For positive increment, continue while current <= end
	}

	// Create and evaluate comparison expression
	compExpr := core.NewList(compSymbol, current, end)
	result := e.Evaluate(c, compExpr)

	// Extract boolean result
	if boolVal, ok := core.ExtractBool(result); ok {
		return boolVal
	}

	// If comparison remains unevaluated, check if it's the exact same expression
	if result.Equal(compExpr) && !result.IsAtom() {
		// Comparison is unevaluated - this indicates missing comparison definition
		// This should have been caught during validation, but stop iteration safely
		return false
	}

	// If we get here, the comparison evaluated to something other than a boolean
	// This might be valid in some mathematical contexts, so be conservative
	return false
}

// evaluateIteratorIncrement adds increment to current value using expression arithmetic
func evaluateIteratorIncrement(e *engine.Evaluator, c *engine.Context, current, increment core.Expr) core.Expr {
	plusExpr := core.NewList("Plus", current, increment)
	return e.Evaluate(c, plusExpr)
}

// isNegativeIncrement determines if increment is negative using expression evaluation
func isNegativeIncrement(e *engine.Evaluator, c *engine.Context, increment core.Expr) bool {
	// Create comparison: increment < 0
	zeroExpr := core.NewInteger(0)
	lessExpr := core.NewList("Less", increment, zeroExpr)
	result := e.Evaluate(c, lessExpr)

	if boolVal, ok := core.ExtractBool(result); ok {
		return boolVal
	}

	// Default to positive if comparison fails
	return false
}

// evaluateWithIteratorBinding uses Block to bind iterator variable and evaluate expression
func evaluateWithIteratorBinding(e *engine.Evaluator, c *engine.Context, expr core.Expr, variable string, value core.Expr) core.Expr {

	// Create Block(List(Set(variable, value)), expr)
	setExpr := core.NewList("Set", core.NewSymbol(variable), value)
	blockVars := core.NewList("List", setExpr)
	//blockArgs := []core.Expr{blockVars, expr}
	return BlockExpr(e, c, blockVars, expr)
}
