package builtins

import (
	"github.com/client9/cardinal/core"
	"github.com/client9/cardinal/engine"

	"slices"
	"sort"
)

// @ExprSymbol Union

// @ExprPattern (__)
func Union(e *engine.Evaluator, c *engine.Context, args []core.Expr) core.Expr {

	// TODO: Pattern should validate all arguments are list-like
	// TODO: check all heads are the same
	n := 0
	for _, a := range args {
		n += int(a.Length())
	}

	union := make([]core.Expr, 0, n)

	for _, a := range args {
		union = append(union, a.(core.List).Tail()...)
	}

	// canoncialcompare has wrong signature
	//slices.SortFunc(union, core.CanonicalCompare)

	// Sort arguments using canonical ordering
	sort.Slice(union, func(i, j int) bool {
		return core.CanonicalCompare(union[i], union[j])
	})
	union = slices.CompactFunc(union, func(a core.Expr, b core.Expr) bool {
		return a.Equal(b)
	})
	return core.ListFrom(args[0].Head(), union...)
}
