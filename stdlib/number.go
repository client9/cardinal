package stdlib

import (
        "github.com/client9/sexpr/core"
)

// Number indicates something was a number before being cast to a float64
type Number float64

func ExtractNumber(e core.Expr) (Number, bool) {
	val, ok := core.GetNumericValue(e)
	return Number(val), ok
}
