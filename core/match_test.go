package core

import (
	"testing"
)

func compareBindings(t *testing.T, want, got map[string]Expr) {
	t.Helper()
	for k, v := range want {
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

func TestMatchLiteralListSequence1a(t *testing.T) {
	e := MustParse("1")

	tests := []struct {
		pattern  Pattern
		expected bool
	}{
		{MatchLiteral(NewInteger(1)), true},
		//{MatchLiteral(NewList("List", MatchLiteral(NewInteger(1)))), true},
		{MatchList(MatchLiteral(NewInteger(1))), false},

		// TBD: Maybe true?
		{MatchSequence(MatchLiteral(NewInteger(1))), false},
	}

	for _, tt := range tests {
		compiled, err := CompilePattern(tt.pattern)
		if err != nil {
			t.Errorf("Compile of %v failed", tt.pattern)
			continue
		}
		result := compiled.Match(e)
		if result.Matched != tt.expected {
			t.Errorf("Match of expression %s with %v expected: %v, got %v", e, tt.pattern, tt.expected, result.Matched)
		}
	}
}

func TestMatchLiteralListSequence1b(t *testing.T) {
	e := MustParse("[1]")

	tests := []struct {
		pattern  Pattern
		expected bool
	}{
		{MatchLiteral(NewInteger(1)), false},
		//{MatchLiteral(NewList("List", MatchLiteral(NewInteger(1)))), true},
		{MatchList(MatchLiteral(NewInteger(1))), true},
		{MatchSequence(MatchLiteral(NewInteger(1))), false},
	}

	for _, tt := range tests {
		compiled, err := CompilePattern(tt.pattern)
		if err != nil {
			t.Errorf("Compile of %v failed", tt.pattern)
			continue
		}
		result := compiled.Match(e)
		if result.Matched != tt.expected {
			t.Errorf("Match of expression %s with %v expected: %v, got %v", e, tt.pattern, tt.expected, result.Matched)
		}
	}
}

func TestMatchLiteralListsBinding(t *testing.T) {
	e := MustParse("[1,2,3]")

	pattern1 := MatchList(
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

func TestMatchListLeadingBinding1(t *testing.T) {
	e := MustParse("[1,2,3]")

	pattern1 := MatchList(
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

func TestMatchListLeadingBinding2(t *testing.T) {
	e := MustParse("[1,2,3]")

	pattern1 := MatchList(
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

func TestMatchTwoListsDirect1a(t *testing.T) {
	e := MustParse("[1,2,3]")

	// inner is MatchSequence, which should match
	pattern1 := MatchList(
		MatchSequence(MatchHead("Integer"), MatchHead("Integer")),
		MatchHead("Integer"),
	)

	compiled1, _ := CompilePattern(pattern1)
	result := compiled1.Match(e)
	if !result.Matched {
		t.Errorf("Match of expression %s with %s Failed: %v", e, pattern1, result)
	}

	// innert is MatchList which should NOT match
	pattern2 := MatchList(
		MatchList(MatchHead("Integer"), MatchHead("Integer")),
		MatchHead("Integer"),
	)

	compiled2, _ := CompilePattern(pattern2)
	result = compiled2.Match(e)
	if result.Matched {
		t.Errorf("Unexpected match of expression %s with %v", e, pattern2)
	}
}

func TestMatchTwoListsDirect1b(t *testing.T) {
	e := MustParse("[[1,2],3]")

	// inner is MatchSequence, which should NOT match, first element is a list.
	pattern1 := MatchList(
		MatchSequence(MatchHead("Integer"), MatchHead("Integer")),
		MatchHead("Integer"),
	)

	// inner is MatchList which should match [[1, 2], 3]
	pattern2 := MatchList(
		MatchList(MatchHead("Integer"), MatchHead("Integer")),
		MatchHead("Integer"),
	)

	// check pattern 1 - expect NO match (MatchSequence should not match a List element)

	compiled1, _ := CompilePattern(pattern1)
	result := compiled1.Match(e)
	if result.Matched {
		t.Errorf("Pattern 1 should NOT match %s with MatchSequence, but it did: %v", e, result)
	}

	// check pattern 2 - expect match (MatchList should match a List element)

	compiled2, _ := CompilePattern(pattern2)
	result = compiled2.Match(e)
	if !result.Matched {
		t.Errorf("Pattern 2 should match %s with MatchList: %v", e, result)
	}
}

func TestMatchBindingTwoListsDirect1a(t *testing.T) {
	e1 := MustParse("[1,2,3]")

	pattern1 := MatchList(
		Named("x", MatchSequence(MatchHead("Integer"), MatchHead("Integer"))),
		Named("y", MatchHead("Integer")),
	)

	want := map[string]Expr{
		"x": MustParse("[1,2]"),
		"y": MustParse("3"),
	}

	compiled1, _ := CompilePattern(pattern1)
	result1 := compiled1.Match(e1)
	if !result1.Matched {
		t.Errorf("Pattern should match [[1,2], 3] when using inner MatchList: %v", result1)
	}
	got := result1.Bindings
	compareBindings(t, want, got)

	pattern2 := MatchList(
		Named("x", MatchList(MatchHead("Integer"), MatchHead("Integer"))),
		Named("y", MatchHead("Integer")),
	)

	// With MatchList, [1,2,3] should NOT match (first element is not a list)
	compiled2, _ := CompilePattern(pattern2)
	result2 := compiled2.Match(e1)
	if result2.Matched {
		t.Errorf("Pattern should NOT match [1,2,3] when using inner MatchList, but it did: %v", result2)
	}

}

func TestMatchBindingTwoListsDirect1b(t *testing.T) {
	e1 := MustParse("[[1,2],3]")

	want := map[string]Expr{
		"x": MustParse("[1,2]"),
		"y": MustParse("3"),
	}

	// inner is MatchList which should match [[1,2],3]
	pattern1 := MatchList(
		Named("x", MatchList(MatchHead("Integer"), MatchHead("Integer"))),
		Named("y", MatchHead("Integer")),
	)

	// inner is MatchSequence which should not match [[1,2],3]
	pattern2 := MatchList(
		Named("x", MatchSequence(MatchHead("Integer"), MatchHead("Integer"))),
		Named("y", MatchHead("Integer")),
	)

	// check pattern 1 - expect match

	compiled1, _ := CompilePattern(pattern1)
	result1 := compiled1.Match(e1)
	if !result1.Matched {
		t.Errorf("Pattern should match [[1,2], 3] when using inner MatchList: %v", result1)
	}

	got := result1.Bindings
	compareBindings(t, want, got)

	// check pattern 2 - expect no match

	compiled2, _ := CompilePattern(pattern2)
	result2 := compiled2.Match(e1)
	if result2.Matched {
		t.Errorf("Pattern should NOT match [[1,2],3] when using inner MatchSequence, but it did: %v", result2)
	}
}

func TestMatchTwoLists1(t *testing.T) {
	e := MustParse("[1,2,3]")

	pattern1 := MatchList(
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

func TestMatchTwoLists2(t *testing.T) {
	e := MustParse("[1,2,3]")

	pattern1 := MatchList(
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
func TestMatchTwoLists3(t *testing.T) {
	e := MustParse("[1,2,3]")

	pattern1 := MatchList(
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

func TestMatchTwoListsNoBinding(t *testing.T) {

	e := MustParse("[1,2,3]")

	pattern1 := MatchList(
		OneOrMore(MatchHead("Integer")),
		OneOrMore(MatchHead("Integer")),
	)

	compiled1, _ := CompilePattern(pattern1)
	result := compiled1.Match(e)
	if !result.Matched {
		t.Errorf("Match of expression %s with %s Failed: %v", e, pattern1, result)
	}
}

func TestMatchBindingSimpleOr(t *testing.T) {
	e := MustParse("1")
	pattern1 := Named("x", MatchOr(MatchLiteral(NewInteger(2)), MatchLiteral(NewInteger(1))))
	want := map[string]Expr{
		"x": MustParse("1"),
	}

	compiled1, _ := CompilePattern(pattern1)
	result := compiled1.Match(e)
	if !result.Matched {
		t.Errorf("Match of expression %s with %s Failed: %v", e, pattern1, result)
	}
	got := result.Bindings
	compareBindings(t, want, got)
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
	pattern1 := MatchList(
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
	pattern1 := MatchList(
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

func BenchmarkMatchListZeroOrMore(b *testing.B) {
	pattern1 := MatchList(
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
	pattern1 := MatchList(
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

func BenchmarkMatchListOneOrMore(b *testing.B) {
	pattern1 := MatchList(
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
	pattern1 := MatchList(
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
	pattern1 := MatchList(
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
	pattern1 := MatchList(
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
	pattern1 := MatchList(
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
	pattern1 := MatchList(
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
	pattern1 := MatchList(
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
	pattern1 := MatchList(
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
	pattern1 := MatchList(
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
	pattern1 := MatchList(
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
func BenchmarkMatchSimpleOr(b *testing.B) {
	pattern1 := MatchOr(MatchLiteral(NewInteger(2)), MatchLiteral(NewInteger(1)))

	expr1 := MustParse("1")

	compiled1, _ := CompilePattern(pattern1)
	counter := 0
	for b.Loop() {
		compiled1.Match(expr1)
		counter++
	}
}
func BenchmarkMatchBindingSimpleOr(b *testing.B) {
	pattern1 := Named("x", MatchOr(MatchLiteral(NewInteger(2)), MatchLiteral(NewInteger(1))))

	expr1 := MustParse("1")

	compiled1, _ := CompilePattern(pattern1)
	counter := 0
	for b.Loop() {
		compiled1.Match(expr1)
		counter++
	}
}

func BenchmarkMatchTwoListsSequenceDirect(b *testing.B) {
	pattern1 := MatchList(
		MatchSequence(MatchHead("Integer"), MatchHead("Integer")),
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
func BenchmarkMatchBindingTwoListsDirectSequence(b *testing.B) {
	pattern1 := MatchList(
		Named("x", MatchSequence(MatchHead("Integer"), MatchHead("Integer"))),
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

func BenchmarkMatchTwoListsOneOrMore(b *testing.B) {
	e := MustParse("[1,2,3]")
	pattern1 := MatchList(
		OneOrMore(MatchHead("Integer")),
		OneOrMore(MatchHead("Integer")),
	)
	compiled1, _ := CompilePattern(pattern1)
	counter := 0
	for b.Loop() {
		compiled1.Match(e)
		counter++
	}
}
func BenchmarkMatchBindingTwoListsOneOrMore(b *testing.B) {
	e := MustParse("[1,2,3]")
	pattern1 := MatchList(
		Named("x", OneOrMore(MatchHead("Integer"))),
		Named("y", OneOrMore(MatchHead("Integer"))),
	)
	compiled1, _ := CompilePattern(pattern1)
	counter := 0
	for b.Loop() {
		compiled1.Match(e)
		counter++
	}
}
