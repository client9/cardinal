package sexpr

// Special Forms - these have non-standard evaluation semantics

// evaluateSet evaluates Set[lhs, rhs] expressions (x = value)
func (e *Evaluator) evaluateSet(args []Expr, ctx *Context) Expr {
	if len(args) != 2 {
		return List{Elements: []Expr{NewSymbolAtom("Set"), args[0], args[1]}}
	}

	lhs := args[0]
	rhs := e.evaluate(args[1], ctx) // Evaluate the right-hand side

	// Handle simple symbol assignment
	if atom, ok := lhs.(Atom); ok && atom.AtomType == SymbolAtom {
		symbolName := atom.Value.(string)
		ctx.Set(symbolName, rhs)
		return rhs
	}

	// TODO: Handle more complex patterns (f[x_] := body, etc.)

	return List{Elements: []Expr{NewSymbolAtom("Set"), lhs, rhs}}
}

// evaluateSetDelayed evaluates SetDelayed[lhs, rhs] expressions (x := value)
func (e *Evaluator) evaluateSetDelayed(args []Expr, ctx *Context) Expr {
	if len(args) != 2 {
		return List{Elements: []Expr{NewSymbolAtom("SetDelayed"), args[0], args[1]}}
	}

	lhs := args[0]
	rhs := args[1] // Don't evaluate the right-hand side for SetDelayed

	// Handle simple symbol assignment
	if atom, ok := lhs.(Atom); ok && atom.AtomType == SymbolAtom {
		symbolName := atom.Value.(string)
		ctx.Set(symbolName, rhs)
		return NewSymbolAtom("Null")
	}

	// Handle function definition: f(x) := 2*x
	if list, ok := lhs.(List); ok && len(list.Elements) > 0 {
		if head, ok := list.Elements[0].(Atom); ok && head.AtomType == SymbolAtom {
			// Register the function definition in the function registry
			err := ctx.functionRegistry.RegisterUserFunction(lhs, rhs)
			if err != nil {
				// Log error for debugging
				// fmt.Printf("DEBUG: RegisterUserFunction failed: %v\n", err)
				// If registration fails, fall back to the old variable-based system for compatibility
				ctx.Set(head.Value.(string), List{Elements: []Expr{
					NewSymbolAtom("Function"),
					List{Elements: list.Elements[1:]},
					rhs,
				}})
			}

			return NewSymbolAtom("Null")
		}
	}

	// Handle more complex patterns later
	return List{Elements: []Expr{NewSymbolAtom("SetDelayed"), lhs, rhs}}
}

// evaluateUnset evaluates Unset[x] expressions (x =.)
func (e *Evaluator) evaluateUnset(args []Expr, ctx *Context) Expr {
	if len(args) != 1 {
		return List{Elements: []Expr{NewSymbolAtom("Unset"), args[0]}}
	}

	if atom, ok := args[0].(Atom); ok && atom.AtomType == SymbolAtom {
		symbolName := atom.Value.(string)
		// Remove from context by setting to itself (undefined)
		delete(ctx.variables, symbolName)
		return NewSymbolAtom("Null")
	}

	return List{Elements: []Expr{NewSymbolAtom("Unset"), args[0]}}
}

// evaluateHold evaluates Hold[...] expressions (prevents evaluation)
func (e *Evaluator) evaluateHold(args []Expr, ctx *Context) Expr {
	// Hold simply returns its arguments without evaluation
	elements := make([]Expr, len(args)+1)
	elements[0] = NewSymbolAtom("Hold")
	copy(elements[1:], args)
	return List{Elements: elements}
}

// evaluateEvaluate evaluates Evaluate[...] expressions (forces evaluation)
func (e *Evaluator) evaluateEvaluate(args []Expr, ctx *Context) Expr {
	if len(args) == 0 {
		return NewSymbolAtom("Null")
	}

	if len(args) == 1 {
		return e.evaluate(args[0], ctx)
	}

	// Multiple arguments - evaluate all and return a sequence
	evaluatedArgs := make([]Expr, len(args))
	for i, arg := range args {
		evaluatedArgs[i] = e.evaluate(arg, ctx)
	}

	// Return as a List for now (in full Mathematica this would be a Sequence)
	return List{Elements: evaluatedArgs}
}

// evaluateIf evaluates If[condition, then, else] expressions
func (e *Evaluator) evaluateIf(args []Expr, ctx *Context) Expr {
	if len(args) < 2 || len(args) > 3 {
		// Return unchanged if wrong number of arguments
		elements := make([]Expr, len(args)+1)
		elements[0] = NewSymbolAtom("If")
		copy(elements[1:], args)
		return List{Elements: elements}
	}

	// Evaluate the condition
	condition := e.evaluate(args[0], ctx)

	// Check if condition is boolean
	if isBool(condition) {
		condValue, _ := getBoolValue(condition)
		if condValue {
			// Condition is true, evaluate and return the 'then' branch
			return e.evaluate(args[1], ctx)
		} else {
			// Condition is false
			if len(args) == 3 {
				// Evaluate and return the 'else' branch
				return e.evaluate(args[2], ctx)
			} else {
				// No else branch, return Null
				return NewSymbolAtom("Null")
			}
		}
	}

	// If condition is not boolean, return unchanged
	elements := make([]Expr, len(args)+1)
	elements[0] = NewSymbolAtom("If")
	elements[1] = condition
	copy(elements[2:], args[1:])
	return List{Elements: elements}
}

// evaluateWhile evaluates While[condition, body] expressions
func (e *Evaluator) evaluateWhile(args []Expr, ctx *Context) Expr {
	if len(args) != 2 {
		// Return unchanged if wrong number of arguments
		elements := make([]Expr, len(args)+1)
		elements[0] = NewSymbolAtom("While")
		copy(elements[1:], args)
		return List{Elements: elements}
	}

	var lastResult Expr = NewSymbolAtom("Null")

	for {
		// Evaluate the condition
		condition := e.evaluate(args[0], ctx)

		// Check if condition is boolean
		if isBool(condition) {
			condValue, _ := getBoolValue(condition)
			if !condValue {
				break // Exit loop if condition is false
			}
		} else {
			// If condition is not boolean, exit loop
			break
		}

		// Evaluate the body
		lastResult = e.evaluate(args[1], ctx)
	}

	return lastResult
}

// evaluateCompoundExpression evaluates CompoundExpression[...] expressions (;)
func (e *Evaluator) evaluateCompoundExpression(args []Expr, ctx *Context) Expr {
	if len(args) == 0 {
		return NewSymbolAtom("Null")
	}

	var lastResult Expr = NewSymbolAtom("Null")

	// Evaluate all expressions in sequence, return the last result
	for _, arg := range args {
		lastResult = e.evaluate(arg, ctx)
	}

	return lastResult
}

// evaluateModule evaluates Module[vars, body] expressions (local variables)
func (e *Evaluator) evaluateModule(args []Expr, ctx *Context) Expr {
	if len(args) != 2 {
		// Return unchanged if wrong number of arguments
		elements := make([]Expr, len(args)+1)
		elements[0] = NewSymbolAtom("Module")
		copy(elements[1:], args)
		return List{Elements: elements}
	}

	varsExpr := args[0]
	body := args[1]

	// Create a new child context for local variables
	childCtx := NewChildContext(ctx)

	// Initialize local variables (simplified - assumes List[var1, var2, ...])
	if varsList, ok := varsExpr.(List); ok {
		for _, varExpr := range varsList.Elements {
			if atom, ok := varExpr.(Atom); ok && atom.AtomType == SymbolAtom {
				varName := atom.Value.(string)
				childCtx.Set(varName, NewSymbolAtom(varName)) // Initialize to symbol
			}
		}
	}

	// Evaluate the body in the child context
	return e.evaluate(body, childCtx)
}

// evaluateBlock evaluates Block[vars, body] expressions (local variables with dynamic scoping)
func (e *Evaluator) evaluateBlock(args []Expr, ctx *Context) Expr {
	if len(args) != 2 {
		// Return unchanged if wrong number of arguments
		elements := make([]Expr, len(args)+1)
		elements[0] = NewSymbolAtom("Block")
		copy(elements[1:], args)
		return List{Elements: elements}
	}

	varsExpr := args[0]
	body := args[1]

	// Save current values of variables
	savedValues := make(map[string]Expr)

	// Block variables (simplified - assumes List[var1, var2, ...])
	if varsList, ok := varsExpr.(List); ok {
		for _, varExpr := range varsList.Elements {
			if atom, ok := varExpr.(Atom); ok && atom.AtomType == SymbolAtom {
				varName := atom.Value.(string)
				if oldValue, exists := ctx.Get(varName); exists {
					savedValues[varName] = oldValue
				}
				ctx.Set(varName, NewSymbolAtom(varName)) // Initialize to symbol
			}
		}
	}

	// Evaluate the body
	result := e.evaluate(body, ctx)

	// Restore previous values
	for varName, oldValue := range savedValues {
		ctx.Set(varName, oldValue)
	}

	return result
}

// evaluateAnd evaluates And[...] expressions with short-circuit evaluation
func (e *Evaluator) evaluateAnd(args []Expr, ctx *Context) Expr {
	if len(args) == 0 {
		return NewBoolAtom(true) // And[] = True
	}

	// Short-circuit evaluation - evaluate each argument
	for _, arg := range args {
		evaluatedArg := e.evaluate(arg, ctx)

		// Propagate errors
		if IsError(evaluatedArg) {
			return evaluatedArg
		}

		if isBool(evaluatedArg) {
			if val, _ := getBoolValue(evaluatedArg); !val {
				return NewBoolAtom(false)
			}
		} else {
			// If any argument is not boolean, return symbolic form
			elements := make([]Expr, len(args)+1)
			elements[0] = NewSymbolAtom("And")
			copy(elements[1:], args)
			return List{Elements: elements}
		}
	}

	return NewBoolAtom(true)
}

// evaluateOr evaluates Or[...] expressions with short-circuit evaluation
func (e *Evaluator) evaluateOr(args []Expr, ctx *Context) Expr {
	if len(args) == 0 {
		return NewBoolAtom(false) // Or[] = False
	}

	// Short-circuit evaluation - evaluate each argument
	for _, arg := range args {
		evaluatedArg := e.evaluate(arg, ctx)

		// Propagate errors
		if IsError(evaluatedArg) {
			return evaluatedArg
		}

		if isBool(evaluatedArg) {
			if val, _ := getBoolValue(evaluatedArg); val {
				return NewBoolAtom(true)
			}
		} else {
			// If any argument is not boolean, return symbolic form
			elements := make([]Expr, len(args)+1)
			elements[0] = NewSymbolAtom("Or")
			copy(elements[1:], args)
			return List{Elements: elements}
		}
	}

	return NewBoolAtom(false)
}
