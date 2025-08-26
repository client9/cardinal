package core

// TODO if still valid or not

// Number indicates something was a number before being cast to a float64
type Number float64

func ExtractNumber(e Expr) (Number, bool) {
	val, ok := GetNumericValue(e)
	return Number(val), ok
}
