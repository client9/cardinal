package core

import (
	"testing"
)

func compareBindings(t *testing.T,  want, got map[string]Expr) {
 	t.Helper()
	for k,v := range want {
		if v == nil && got[k] == nil {
			continue
		}
		if v == nil && got[k] != nil {
			t.Errorf("binding %q: want nil but got %v", k, got[k])
			continue
		}
		if got[k] == nil {
			t.Errorf("binding %q: want %v, got nil", k, v)
			continue
		}
		if !v.Equal(got[k]) {
			t.Errorf("binding %q: want %v, got %v", k, v, got[k])
		}
	}

}
func TestMatchLiteralSequencesBinding(t *testing.T) {
	e := MustParse("[1,2,3]")

	pattern1 := MatchSequence(
		Named("x", MatchHead("Integer")),
		Named("y", OneOrMore(MatchHead("Integer"))),
	)

	want := map[string]Expr{
		"x": MustParse("1"),
		"y": MustParse("[2,3]"),
	}

	compiled1, _ := CompilePattern(pattern1)
	result := compiled1.Match(e)
	if !result.Matched {
		t.Errorf("Match of expression %s with %s Failed: %v", e, pattern1, result)
	}
	got := result.Bindings
	compareBindings(t, want, got)
}

func TestMatchSequenceLeadingBinding1(t *testing.T) {
	e := MustParse("[1,2,3]")

	pattern1 := MatchSequence(
		Named("x", OneOrMore(MatchHead("Integer"))),
		Named("y", MatchHead("Integer")),
	)

	want := map[string]Expr{
		"x": MustParse("[1,2]"),
		"y": MustParse("3"),
	}

	compiled1, _ := CompilePattern(pattern1)
	result := compiled1.Match(e)
	if !result.Matched {
		t.Errorf("Match of expression %s with %s Failed: %v", e, pattern1, result)
	}
	got := result.Bindings
	compareBindings(t, want, got)
}

func TestMatchSequenceLeadingBinding2(t *testing.T) {
	e := MustParse("[1,2,3]")

	pattern1 := MatchSequence(
		Named("x", ZeroOrMore(MatchHead("Integer"))),
		Named("y", MatchHead("Integer")),
	)

	want := map[string]Expr{
		"x": MustParse("[1,2]"),
		"y": MustParse("3"),
	}

	compiled1, _ := CompilePattern(pattern1)
	result := compiled1.Match(e)
	if !result.Matched {
		t.Errorf("Match of expression %s with %s Failed: %v", e, pattern1, result)
	}
	got := result.Bindings
	compareBindings(t, want, got)
}

func TestMatchTwoSequences1(t *testing.T) {
	e := MustParse("[1,2,3]")

	pattern1 := MatchSequence(
		Named("x", OneOrMore(MatchHead("Integer"))),
		Named("y", OneOrMore(MatchHead("Integer"))),
	)

	want := map[string]Expr{
		"x": MustParse("[1,2]"),
		"y": MustParse("3"),
	}

	compiled1, _ := CompilePattern(pattern1)
	result := compiled1.Match(e)
	if !result.Matched {
		t.Errorf("Match of expression %s with %s Failed: %v", e, pattern1, result)
	}
	got := result.Bindings
	compareBindings(t, want, got)
}

func TestMatchTwoSequences2(t *testing.T) {
	e := MustParse("[1,2,3]")

	pattern1 := MatchSequence(
		Named("x", ZeroOrMore(MatchHead("Integer"))),
		Named("y", ZeroOrMore(MatchHead("Integer"))),
	)

	want := map[string]Expr{
		"x": MustParse("[1,2,3]"),
		"y": nil,
	}

	compiled1, _ := CompilePattern(pattern1)
	result := compiled1.Match(e)
	if !result.Matched {
		t.Errorf("Match of expression %s with %s Failed: %v", e, pattern1, result)
	}
	got := result.Bindings
	compareBindings(t, want, got)
}
func TestMatchTwoSequences3(t *testing.T) {
	e := MustParse("[1,2,3]")

	pattern1 := MatchSequence(
		Named("x", ZeroOrMore(MatchHead("Integer"))),
		Named("y", OneOrMore(MatchHead("Integer"))),
	)

	want := map[string]Expr{
		"x": MustParse("[1,2]"),
		"y": MustParse("3"),
	}

	compiled1, _ := CompilePattern(pattern1)
	result := compiled1.Match(e)
	if !result.Matched {
		t.Errorf("Match of expression %s with %s Failed: %v", e, pattern1, result)
	}
	got := result.Bindings
	compareBindings(t, want, got)
}

func TestMatchTwoSequencesNoBinding(t *testing.T) {

	e := MustParse("[1,2,3]")

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

func BenchmarkMatchNoBindings(b *testing.B) {
	e := NewList("List", NewInteger(1), NewInteger(2))
	p := NewPatternMatcher()
	m := NewList("List", CreateBlankExpr(NewSymbol("Integer")), CreateBlankExpr(NewSymbol("Integer")))

	for b.Loop() {
		p.TestMatch(e, m)
	}
}
func BenchmarkMatchWithBindings(b *testing.B) {
	e := MustParse("[1,2]")
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

	expr1 := MustParse("[1,2]")

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

	expr1 := MustParse("[1,2]")

	compiled1, _ := CompilePattern(pattern1)
	counter := 0
	for b.Loop() {
		compiled1.Match(expr1)
		counter++
	}
}

func BenchmarkMatchSequenceZeroOrMore(b *testing.B) {
	pattern1 := MatchSequence(
		ZeroOrMore(MatchHead("Integer")),
	)

	expr1 := MustParse("[1,2,3]")

	compiled1, _ := CompilePattern(pattern1)
	counter := 0
	for b.Loop() {
		compiled1.Match(expr1)
		counter++
	}
}

func BenchmarkMatchBindingSequenceZeroOrMore(b *testing.B) {
	pattern1 := MatchSequence(
		Named("x", ZeroOrMore(MatchHead("Integer"))),
	)

	expr1 := MustParse("[1,2,3]")

	compiled1, _ := CompilePattern(pattern1)
	counter := 0
	for b.Loop() {
		compiled1.Match(expr1)
		counter++
	}
}

func BenchmarkMatchSequenceOneOrMore(b *testing.B) {
	pattern1 := MatchSequence(
		OneOrMore(MatchHead("Integer")),
	)

	expr1 := MustParse("[1,2,3]")

	compiled1, _ := CompilePattern(pattern1)
	counter := 0
	for b.Loop() {
		compiled1.Match(expr1)
		counter++
	}
}

func BenchmarkMatchBindingSequenceOneOrMore(b *testing.B) {
	pattern1 := MatchSequence(
		Named("x", OneOrMore(MatchHead("Integer"))),
	)

	expr1 := MustParse("[1,2,3]")

	compiled1, _ := CompilePattern(pattern1)
	counter := 0
	for b.Loop() {
		compiled1.Match(expr1)
		counter++
	}
}

func BenchmarkMatchLeadingSequenceZeroOrMore(b *testing.B) {
	pattern1 := MatchSequence(
		ZeroOrMore(MatchHead("Integer")),
		MatchHead("Integer"),
	)

	expr1 := MustParse("[1,2,3]")

	compiled1, _ := CompilePattern(pattern1)
	counter := 0
	for b.Loop() {
		compiled1.Match(expr1)
		counter++
	}
}

func BenchmarkMatchBindingLeadingSequenceZeroOrMore(b *testing.B) {
	pattern1 := MatchSequence(
		Named("x", ZeroOrMore(MatchHead("Integer"))),
		Named("y", MatchHead("Integer")),
	)

	expr1 := MustParse("[1,2,3]")

	compiled1, _ := CompilePattern(pattern1)
	counter := 0
	for b.Loop() {
		compiled1.Match(expr1)
		counter++
	}
}
func BenchmarkMatchLeadingSequenceOneOrMore(b *testing.B) {
	pattern1 := MatchSequence(
		OneOrMore(MatchHead("Integer")),
		MatchHead("Integer"),
	)

	expr1 := MustParse("[1,2,3]")

	compiled1, _ := CompilePattern(pattern1)
	counter := 0
	for b.Loop() {
		compiled1.Match(expr1)
		counter++
	}
}

func BenchmarkMatchBindingLeadingSequenceOneOrMore(b *testing.B) {
	pattern1 := MatchSequence(
		Named("x", OneOrMore(MatchHead("Integer"))),
		Named("y", MatchHead("Integer")),
	)

	expr1 := MustParse("[1,2,3]")

	compiled1, _ := CompilePattern(pattern1)
	counter := 0
	for b.Loop() {
		compiled1.Match(expr1)
		counter++
	}
}

func BenchmarkMatchTrailingSequenceZeroOrMore(b *testing.B) {
	pattern1 := MatchSequence(
		MatchHead("Integer"),
		ZeroOrMore(MatchHead("Integer")),
	)

	expr1 := MustParse("[1,2,3]")

	compiled1, _ := CompilePattern(pattern1)
	counter := 0
	for b.Loop() {
		compiled1.Match(expr1)
		counter++
	}
}

func BenchmarkMatchBindingTrailingSequenceZeroOrMore(b *testing.B) {
	pattern1 := MatchSequence(
		Named("x", MatchHead("Integer")),
		Named("y", ZeroOrMore(MatchHead("Integer"))),
	)

	expr1 := MustParse("[1,2,3]")

	compiled1, _ := CompilePattern(pattern1)
	counter := 0
	for b.Loop() {
		compiled1.Match(expr1)
		counter++
	}
}

func BenchmarkMatchTrailingSequenceOneOrMore(b *testing.B) {
	pattern1 := MatchSequence(
		MatchHead("Integer"),
		OneOrMore(MatchHead("Integer")),
	)

	expr1 := MustParse("[1,2,3]")

	compiled1, _ := CompilePattern(pattern1)
	counter := 0
	for b.Loop() {
		compiled1.Match(expr1)
		counter++
	}
}
func BenchmarkMatchBindingTrailingSequenceOneOrMore(b *testing.B) {
	pattern1 := MatchSequence(
		Named("x", MatchHead("Integer")),
		Named("y", OneOrMore(MatchHead("Integer"))),
	)

	expr1 := MustParse("[1,2,3]")

	compiled1, _ := CompilePattern(pattern1)
	counter := 0
	for b.Loop() {
		compiled1.Match(expr1)
		counter++
	}
}
