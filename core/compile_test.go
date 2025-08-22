package core

import (
	"fmt"
	"testing"
)

func TestCompileCheck(t *testing.T) {
	c := NewCompiler()
	//e := MustParse(" 1,2,3 ")
	e := MustParse("[ MatchAny(), MatchStar(MatchAny()), MatchStar(MatchAny()) ]")
	//e := MustParse("[MatchPlus(MatchAny(Integer)), MatchStar(MatchHead(Integer))]")
	p := c.Compile(e)
	p.Dump()
	fmt.Printf("----\n")
	p2 := c.CompileOneStep(e)
	p2.Dump()
}

func TestCompileGroups(t *testing.T) {
	c := NewCompiler()
	e := MustParse("[ [ Named(x,MatchHead(Integer)), Named(y,MatchHead(Integer)) ], Named(n,MatchHead(Integer)) ]")
	elist := e.(List)
	p := c.CompileList(elist.Tail())
	g := p.Groups()
	if len(g) != 3 {
		t.Errorf("Expected 3 groups, got %v", g)
	}
}

func TestCompileSimplePositive(t *testing.T) {
	cases := []string{
		"1",
		"\"foo\"",
		"x",
		"[1,2,3]",
		"MatchAny()",
		"MatchHead(Integer)",
		"MatchSequence(1,2,3)",
		"MatchStar(1)",
		"MatchPlus(1)",
		"MatchQuest(1)",
		"MatchStar(MatchAny())",
		"MatchPlus(MatchAny())",
		"MatchQuest(MatchAny())",
		"MatchStar(MatchHead(Integer))",
		"MatchPlus(MatchAny(Integer))",
		"MatchQuest(MatchAny(Integer))",
		"[ MatchAny() ]",
		"[ MatchAny(), MatchAny() ]",
		"[ Named(x,MatchAny()), Named(y,MatchAny()) ]",
		"[ MatchStar(MatchAny()) ]",
		"[ 1, MatchStar(MatchAny()) ]",
		"[ 1, MatchStar(MatchAny()), MatchStar(MatchAny()) ]",
		"[ MatchAny(), [ MatchAny(), MatchAny()], MatchAny() ]",

		// "[ MatchStar(MatchHead(String)), 1] "
	}

	c := NewCompiler()

	for i, tt := range cases {
		e := MustParse(tt)
		if !c.Simple(e) {
			t.Errorf("Case %d, should be simple: %s", i, tt)
		}

		if elist, ok := e.(List); ok {
			if !c.SimpleList(elist.Tail()) {
				t.Errorf("Case %d, should be simple: %s", i, tt)
			}
		}
	}
}

func TestCompileSimpleNegative(t *testing.T) {
	cases := []string{
		"[ MatchStar(1), 1 ]",
		"[ MatchStar(MatchAny()), 1 ]",
		"[ MatchStar(MatchHead(Integer)), 1 ]",
		"[ 1, MatchPlus(MatchAny()), MatchPlus(MatchAny()) ]",
		"[ MatchQuest(a),MatchQuest(a),MatchQuest(a),a,a,a ]",
	}

	c := NewCompiler()

	for i, tt := range cases {
		e := MustParse(tt)
		if c.Simple(e) {
			t.Errorf("Case %d, should be not simple : %s", i, tt)
		}
		if elist, ok := e.(List); ok {
			if c.SimpleList(elist.Tail()) {
				t.Errorf("Case %d, should be not simple: %s", i, tt)
			}
		}
	}
}
