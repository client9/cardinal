package engine

import (
	"testing"

	"github.com/client9/sexpr/core"
)

func TestFunctionRegistry(t *testing.T) {
	e := NewEvaluator()
	parentCtx := NewContext()

	// Add a custom function to the parent
	customSquare := func(args []core.Expr) core.Expr {
		if len(args) != 1 {
			return core.NewErrorExpr("ArgumentError", "Square expects 1 argument", args)
		}
		
		if val, ok := core.GetNumericValue(args[0]); ok {
			return core.NewReal(val * val)
		}

		return core.List{Elements: []core.Expr{core.NewSymbol("Square"), args[0]}}
	}

	// Register custom function using pattern system
	customSquarePattern := func(e *Evaluator, c *Context, args []core.Expr) core.Expr {
		return customSquare(args)
	}
	err := parentCtx.GetFunctionRegistry().RegisterPatternBuiltins(map[string]PatternFunc{
		"Square(x_Integer)": customSquarePattern,
	})
	if err != nil {
		t.Fatalf("Failed to register Square function: %v", err)
	}

	// Create child context
	childCtx := NewChildContext(parentCtx)

	// Child should inherit the custom function
	args := []core.Expr{core.NewInteger(7)}
	funcDef, _ := childCtx.functionRegistry.FindMatchingFunction("Square", args)
	if funcDef == nil {
		t.Error("Child context should inherit Square function")
	}

	// Test that the inherited function works
	result := funcDef.GoImpl(e, childCtx, args)
	expected := "49.0"
	if result.String() != expected {
		t.Errorf("expected %s, got %s", expected, result.String())
	}
}
