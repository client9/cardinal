package builtins

import (
	"github.com/client9/sexpr/core"
	"github.com/client9/sexpr/engine"
)

// @ExprSymbol RotateRight

// RotateRight rotates elements of a sliceable expression to the right by n positions
// RotateRight(expr, n) moves the last n elements to the beginning
//
// @ExprPattern (_, _Integer)
func RotateRight(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {
	expr := args[0]
	n, _ := core.ExtractInt64(args[1])
	sliceable := core.AsSliceable(expr)
	if sliceable == nil {
		return core.NewError("TypeError",
			"RotateRight requires a sliceable expression (List, String, or ByteArray)")
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

	// Rotate right by n: take elements [length-n+1:] + [1:length-n]
	// For 1-based indexing: slice from (length-n+1) to end, then from 1 to (length-n)
	rightPart := sliceable.Slice(length-n+1, length)
	leftPart := sliceable.Slice(1, length-n)

	// Concatenate the parts
	return core.ConcatenateSliceable(rightPart, leftPart)
}
