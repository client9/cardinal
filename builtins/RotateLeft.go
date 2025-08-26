package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol RotateLeft

// @ExprPattern (_, _Integer)
func RotateLeft(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	expr := args[0]
	n, _ := core.ExtractInt64(args[1])

	sliceable := core.AsSliceable(expr)
	if sliceable == nil {
		return core.NewError("TypeError",
			"RotateLeft requires a sliceable expression (List, String, or ByteArray)")
	}

	length := expr.Length()
	if length == 0 {
		return expr // Empty expressions remain unchanged
	}

	// Normalize n to be within [0, length)
	n = n % length
	if n < 0 {
		n += length
	}

	if n == 0 {
		return expr // No rotation needed
	}

	// Rotate left by n: take elements [n+1:] + [1:n]
	// For 1-based indexing: slice from (n+1) to end, then from 1 to n
	rightPart := sliceable.Slice(n+1, length)
	leftPart := sliceable.Slice(1, n)

	// Concatenate the parts
	return core.ConcatenateSliceable(rightPart, leftPart)
}
