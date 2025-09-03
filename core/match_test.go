package core

import (
	"testing"
)

func BenchmarkSREOLDMatchNoBindings(b *testing.B) {
	e := MustParse("[1,2]")
	m := MustParse("[ _Integer, _Integer]")

	p := NewPatternMatcher()
	for b.Loop() {
		p.TestMatch(e, m)
	}
}
func BenchmarkSREOLDMatchWithBindings(b *testing.B) {
	e := MustParse("[1,2]")
	m := MustParse("[ x_Integer, y_Integer ]")
	for b.Loop() {
		MatchWithBindings(e, m)
	}
}
