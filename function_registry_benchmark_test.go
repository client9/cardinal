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

// BenchmarkNewContext_Imperative simulates the old imperative approach
func BenchmarkNewContext_Imperative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := &Context{
			variables:   make(map[string]Expr),
			parent:      nil,
			symbolTable: NewSymbolTable(),
			builtins:    make(map[string]BuiltinFunc),
		}

		// Simulate the old imperative setup
		ctx.builtins["Plus"] = EvaluatePlus
		ctx.builtins["Times"] = EvaluateTimes
		ctx.builtins["Subtract"] = EvaluateSubtract
		ctx.builtins["Divide"] = EvaluateDivide
		ctx.builtins["Power"] = EvaluatePower
		ctx.builtins["Equal"] = EvaluateEqual
		ctx.builtins["Unequal"] = EvaluateUnequal
		ctx.builtins["Less"] = EvaluateLess
		ctx.builtins["Greater"] = EvaluateGreater
		ctx.builtins["LessEqual"] = EvaluateLessEqual
		ctx.builtins["GreaterEqual"] = EvaluateGreaterEqual
		ctx.builtins["Not"] = EvaluateNot
		ctx.builtins["SameQ"] = EvaluateSameQ
		ctx.builtins["UnsameQ"] = EvaluateUnsameQ

		_ = ctx
	}
}

// BenchmarkFunctionLookup benchmarks function lookup performance
func BenchmarkFunctionLookup(b *testing.B) {
	ctx := NewContext()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Test lookup of different functions
		_, _ = ctx.GetBuiltin("Plus")
		_, _ = ctx.GetBuiltin("Equal")
		_, _ = ctx.GetBuiltin("Not")
		_, _ = ctx.GetBuiltin("Power")
	}
}
