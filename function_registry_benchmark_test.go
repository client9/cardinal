package sexpr

import (
	"testing"
)

// BenchmarkNewContext_MapLiteral benchmarks the current implementation
func BenchmarkNewContext_MapLiteral(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := NewContext()
		_ = ctx
	}
}

// BenchmarkNewContext_Legacy simulates the legacy manual setup approach
func BenchmarkNewContext_Legacy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := &Context{
			variables:        make(map[string]Expr),
			parent:           nil,
			symbolTable:      NewSymbolTable(),
			functionRegistry: NewFunctionRegistry(),
			stack:            NewEvaluationStack(1000),
		}

		// Simulate manually registering each function (old way)
		registry := ctx.GetFunctionRegistry()
		registry.RegisterPatternBuiltins(map[string]PatternFunc{
			"Plus(x___)":       func(args []Expr, ctx *Context) Expr { return EvaluatePlus(args) },
			"Times(x___)":      func(args []Expr, ctx *Context) Expr { return EvaluateTimes(args) },
			"Subtract(x_, y_)": func(args []Expr, ctx *Context) Expr { return EvaluateSubtract(args) },
			"Divide(x_, y_)":   func(args []Expr, ctx *Context) Expr { return EvaluateDivide(args) },
			"Power(x_, y_)":    func(args []Expr, ctx *Context) Expr { return EvaluatePower(args) },
		})

		_ = ctx
	}
}

// BenchmarkFunctionLookup benchmarks function lookup performance
func BenchmarkFunctionLookup(b *testing.B) {
	ctx := NewContext()
	args := []Expr{NewIntAtom(1), NewIntAtom(2)}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Test lookup of different functions using the new pattern system
		_, _ = ctx.GetFunctionRegistry().FindMatchingFunction("Plus", args)
		_, _ = ctx.GetFunctionRegistry().FindMatchingFunction("Equal", args)
		_, _ = ctx.GetFunctionRegistry().FindMatchingFunction("Not", []Expr{NewBoolAtom(true)})
		_, _ = ctx.GetFunctionRegistry().FindMatchingFunction("Power", args)
	}
}
