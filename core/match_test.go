package core

import (
	"testing"
)

func BenchmarkSREOLDMatchNoBindings(b *testing.B) {
	e := NewList("List", NewInteger(1), NewInteger(2))
	p := NewPatternMatcher()
	m := NewList("List", CreateBlankExpr(NewSymbol("Integer")), CreateBlankExpr(NewSymbol("Integer")))

	for b.Loop() {
		p.TestMatch(e, m)
	}
}
func BenchmarkSREOLDMatchWithBindings(b *testing.B) {
	e := MustParse("[1,2]")
	m := NewList("List", CreatePatternExpr(NewSymbol("x"), CreateBlankExpr(NewSymbol("Integer"))),
		CreatePatternExpr(NewSymbol("y"), CreateBlankExpr(NewSymbol("Integer"))))

	for b.Loop() {
		MatchWithBindings(e, m)
	}
}
