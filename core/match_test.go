package core

import (
	"testing"
)

func TestMatchLiteralSequencesBinding(t *testing.T) {
	e := NewList("List", NewInteger(1), NewInteger(2), NewInteger(3))

	pattern1 := MatchSequence(
		Named("x",MatchHead("Integer")),
		Named("y", OneOrMore(MatchHead("Integer"))),
	)
	
	compiled1, _ := CompilePattern(pattern1)
	result := compiled1.Match(e)
	if !result.Matched {
		t.Errorf("Match of expression %s with %s Failed: %v", e, pattern1, result)
	}
	// expect: { y: 1, x:[1,2] }
}

func TestMatchSequencesLiteralBinding(t *testing.T) {

	e := NewList("List", NewInteger(1), NewInteger(2), NewInteger(3))

	pattern1 := MatchSequence(
		Named("x", OneOrMore(MatchHead("Integer"))),
		Named("y",MatchHead("Integer")),
	)
	
	compiled1, _ := CompilePattern(pattern1)
	result := compiled1.Match(e)
	if !result.Matched {
		t.Errorf("Match of expression %s with %s Failed: %v", e, pattern1, result)
	}
	// expect { x:[1,2], y: 3}

}

func TestMatchLiteralSequenceBinding(t *testing.T) {

	e := NewList("List", NewInteger(1), NewInteger(2), NewInteger(3))

	pattern1 := MatchSequence(
		MatchHead("Integer"),
		OneOrMore(MatchHead("Integer")),
	)
	
	compiled1, _ := CompilePattern(pattern1)
	result := compiled1.Match(e)
	if !result.Matched {
		t.Errorf("Match of expression %s with %s Failed: %v", e, pattern1, result)
	}
}
func TestMatchTwoSequencesNoBinding(t *testing.T) {

	e := NewList("List", NewInteger(1), NewInteger(2), NewInteger(3))

	pattern1 := MatchSequence(
		OneOrMore(MatchHead("Integer")),
		OneOrMore(MatchHead("Integer")),
	)
	
	compiled1, _ := CompilePattern(pattern1)
	result := compiled1.Match(e)
	if !result.Matched {
		t.Errorf("Match of expression %s with %s Failed: %v", e, pattern1, result)
	}
}


func TestMatchTwoSequences(t *testing.T) {

	e := NewList("List", NewInteger(1), NewInteger(2), NewInteger(3))

	pattern1 := MatchSequence(
		Named("x", OneOrMore(MatchHead("Integer"))),
		Named("y", OneOrMore(MatchHead("Integer"))),
	)
	
	compiled1, _ := CompilePattern(pattern1)
	result := compiled1.Match(e)
	if !result.Matched {
		t.Errorf("Match of expression %s with %s Failed: %v", e, pattern1, result)
	}
}


func BenchmarkMatchNoBindings(b *testing.B) {
	e := NewList("List", NewInteger(1), NewInteger(2))
	p := NewPatternMatcher()
	m := NewList("List", CreateBlankExpr(NewSymbol("Integer")), CreateBlankExpr(NewSymbol("Integer")))

	for b.Loop() {
		p.TestMatch(e, m)
	}
}
func BenchmarkMatchWithBindings(b *testing.B) {
	e := NewList("List", NewInteger(1), NewInteger(2))
	m := NewList("List", CreatePatternExpr(NewSymbol("x"), CreateBlankExpr(NewSymbol("Integer"))),
		CreatePatternExpr(NewSymbol("y"), CreateBlankExpr(NewSymbol("Integer"))))

	for b.Loop() {
		MatchWithBindings(e, m)
	}
}

// TestExampleUsage demonstrates the s-expression regex system with practical examples
func BenchmarkMatchNew(b *testing.B) {
	// Example 1: Match a list with specific structure
	// Pattern: List(1, MatchHead("String"), MatchAny())
	// Should match: [1, "hello", anything]
	pattern1 := MatchSequence(
		MatchHead("Integer"),
		MatchHead("Integer"),
	)

	expr1 := NewList("List",
		NewInteger(1),
		NewInteger(2),
	)

	compiled1, _ := CompilePattern(pattern1)
	counter := 0
	for b.Loop() {
		compiled1.Match(expr1)
		counter++
	}
}

// TestExampleUsage demonstrates the s-expression regex system with practical examples
func BenchmarkMatchNewBindings(b *testing.B) {
	// Example 1: Match a list with specific structure
	// Pattern: List(1, MatchHead("String"), MatchAny())
	// Should match: [1, "hello", anything]
	pattern1 := MatchSequence(
		Named("x", MatchHead("Integer")),
		Named("y", MatchHead("Integer")),
	)

	expr1 := NewList("List",
		NewInteger(1),
		NewInteger(2),
	)

	compiled1, _ := CompilePattern(pattern1)
	counter := 0
	for b.Loop() {
		compiled1.Match(expr1)
		counter++
	}
}
