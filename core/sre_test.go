package core

import (
	"fmt"
	"strings"
	"testing"
)

type tc struct {
	name    string
	expr    string
	pattern string
	binding string
	match   bool
}

var cases = []tc{
	{
		name:    "Match Literal",
		expr:    "1",
		pattern: "1",
		binding: "",
		match:   true,
	},
	{
		name:    "Match Literal",
		expr:    "[1]",
		pattern: "[1]",
		binding: "",
		match:   true,
	},
	{
		name:    "Match Literal",
		expr:    "[1,2,3]",
		pattern: "[1,2,3]",
		binding: "",
		match:   true,
	},
	{
		name:    "Match Literal",
		expr:    "[1,2,3]",
		pattern: "[1,2,4]",
		binding: "",
		match:   false,
	},
	{
		name:    "Match Literal",
		expr:    "[1,2]",
		pattern: "[1,2,3]",
		binding: "",
		match:   false,
	},
	{
		name:    "Match Literal",
		expr:    "[1,2,3]",
		pattern: "[1,2]",
		binding: "",
		match:   false,
	},
	{
		name:    "Match Literal",
		expr:    "[1,[2,[3]]]",
		pattern: "[1,[2,[3]]]",
		binding: "",
		match:   true,
	},
	{
		name:    "MatchAny in a sequence",
		expr:    "[ a,b,c ]",
		pattern: "[ MatchAny(), MatchAny(), MatchAny() ]",
		binding: "",
		match:   true,
	},
	{
		name:    "MatchAny, MatchStar2",
		expr:    "[ a ]",
		pattern: "[ MatchAny(), MatchStar(MatchAny()), MatchStar(MatchAny()) ]",
		binding: "",
		match:   true,
	},
	{
		name:    "MatchAny,sequence3,many,binding",
		expr:    "[a, b, c]",
		pattern: "[ Pattern(x, MatchAny()), Pattern(y, MatchAny()), Pattern(z, MatchAny()) ]",
		binding: "[ x:a, y:b, z:c ]",
		match:   true,
	},
	{
		name:    "MatchAny,sequence3,matchany2",
		expr:    "[a, b]",
		pattern: "[ Pattern(x, MatchAny()), Pattern(y, MatchAny()), Pattern(z, MatchAny()) ]",
		binding: "",
		match:   false,
	},
	{
		name:    "MatchAny,sequence3,matchany2,matchstar",
		expr:    "[a, b]",
		pattern: "[ Pattern(x, MatchAny()), Pattern(y, MatchAny()), Pattern(z, MatchStar(MatchAny())) ]",
		binding: "[ x:a, y:b ]",
		match:   true,
	},
	{
		name:    "MatchAny,sequence3,matchany2,matchplus",
		expr:    "[a, b]",
		pattern: "[ Pattern(x, MatchAny()), Pattern(y, MatchAny()), Pattern(z, MatchPlus(MatchAny())) ]",
		binding: "",
		match:   false,
	},
	{
		name:    "MatchAny,sequence3,matchany2,matchstar",
		expr:    "[a]",
		pattern: "[ Pattern(x, MatchAny()), Pattern(y, MatchStar(MatchAny())), Pattern(z, MatchStar(MatchAny())) ]",
		binding: "[ x:a ]",
		match:   true,
	},
	{
		name:    "MatchAny,sublist",
		expr:    "[ [ 1, 2 ], 10 ]",
		pattern: "[ [ MatchAny(), MatchAny() ], MatchAny() ]",
		binding: "",
		match:   true,
	},
	{
		name:    "MatchAny,sublist,binding",
		expr:    "[ [ 1, 2 ], 10 ]",
		pattern: "[ [ Pattern(x,MatchAny()), Pattern(y,MatchAny()) ], Pattern(n,MatchAny()) ]",
		binding: "[x:1,y:2,n:10]",
		match:   true,
	},
	{
		name:    "MatchAny,sequence,single",
		expr:    "[ a ]",
		pattern: "[ MatchAny() ]",
		binding: "",
		match:   true,
	},
	{
		name:    "MatchAny,sequence,single,binding",
		expr:    "[ a ]",
		pattern: "[ Pattern(x, MatchAny()) ]",
		binding: "[ x:a ]",
		match:   true,
	},
	{
		name:    "MatchAny in a sequence",
		expr:    "[ a,b,c ]",
		pattern: "[ MatchAny(), MatchAny(), MatchAny() ]",
		binding: "",
		match:   true,
	},
	{
		name:    "MatchQuest,single",
		expr:    "[ a ]",
		pattern: "[ MatchQuest(MatchAny()) ]",
		binding: "",
		match:   true,
	},
	{
		name:    "MatchQuest,sequence",
		expr:    "[ a,b,c ]",
		pattern: "[ MatchQuest(MatchAny()), MatchQuest(MatchAny()), MatchQuest(MatchAny()) ]",
		binding: "",
		match:   true,
	},
	{
		name:    "MatchAny,sequence2,many,binding",
		expr:    "[a, b]",
		pattern: "[ Pattern(x, MatchAny()), Pattern(y, MatchAny()) ]",
		binding: "[ x:a, y:b ]",
		match:   true,
	},
	{
		name:    "MatchAny,sequence3,many,binding",
		expr:    "[a, b, c]",
		pattern: "[ Pattern(x, MatchAny()), Pattern(y, MatchAny()), Pattern(z, MatchAny()) ]",
		binding: "[ x:a, y:b, z:c ]",
		match:   true,
	},
	{
		name:    "MatchAny,sequence,many,negative",
		expr:    "[ a,b,c ]",
		pattern: "[ MatchAny(), MatchAny() ]",
		binding: "",
		match:   false,
	},
	{
		name:    "MatchStar,sequence,many",
		expr:    "[ a,b,c ]",
		pattern: "[ MatchStar(MatchAny()) ]",
		binding: "",
		match:   true,
	},
	{
		name:    "MatchStar,sequence,many,binding",
		expr:    "[a, b, c]",
		pattern: "[ Pattern(x,MatchStar(MatchAny())) ]",
		binding: "[ x:[a, b, c] ]",
		match:   true,
	},
	{
		name:    "MatchStar,sequence,single",
		expr:    "[ a ]",
		pattern: "[ MatchStar(MatchAny()) ]",
		binding: "",
		match:   true,
	},
	{
		name:    "MatchStar,sequence,single,binding",
		expr:    "[ a ]",
		pattern: "[ Pattern(x,MatchStar(MatchAny())) ]",
		binding: "[ x:a ]",
		match:   true,
	},
	{
		name:    "MatchStar,MatchStar,binding",
		expr:    "[ a,b,c ]",
		pattern: "[ Pattern(x,MatchStar(MatchAny())), Pattern(y, MatchStar(MatchAny())) ]",
		binding: "[ x:[a,b,c] ]",
		match:   true,
	},
	{
		name:    "MatchPlus,MatchStar",
		expr:    "[ a,b,c ]",
		pattern: "[ MatchPlus(MatchAny()), MatchStar(MatchAny()) ]",
		binding: "",
		match:   true,
	},
	{
		name:    "MatchPlus,MatchStar,binding",
		expr:    "[ a,b,c ]",
		pattern: "[ Pattern(x,MatchPlus(MatchAny())), Pattern(y, MatchStar(MatchAny())) ]",
		binding: "[ x:[a,b,c] ]",
		match:   true,
	},
	{
		name:    "MatchStar,MatchPlus",
		expr:    "[ a,b,c ]",
		pattern: "[ MatchStar(MatchAny()), MatchPlus(MatchAny()) ]",
		binding: "",
		match:   true,
	},
	{
		name:    "MatchStar,MatchPlus,binding",
		expr:    "[ a,b,c ]",
		pattern: "[ Pattern(x,MatchStar(MatchAny())), Pattern(y, MatchPlus(MatchAny())) ]",
		binding: "[ x:[a,b], y: c ]",
		match:   true,
	},
	{
		name:    "MatchHead,",
		expr:    `[ 1, "a" ]`,
		pattern: "[ Pattern(x,MatchHead(Integer)), Pattern(y, MatchHead(String)) ]",
		binding: `[ x:1, y:"a"]`,
		match:   true,
	},
	{
		name:    "MatchHead,negative",
		expr:    `[ 1, 1 ]`,
		pattern: "[ Pattern(x,MatchHead(Integer)), Pattern(y, MatchHead(String)) ]",
		binding: "",
		match:   false,
	},
	{
		name:    "MatchLiteral,",
		expr:    `[ 1, "a", b ]`,
		pattern: `[ Pattern(x,1), Pattern(y, "a"), Pattern(z, b) ]`,
		binding: `[ x:1, y:"a", z:b]`,
		match:   true,
	},
	{
		name:    "MatchHead,negative",
		expr:    `[ 1, 1 ]`,
		pattern: "[ Pattern(x,2), Pattern(y, 1) ]",
		binding: "",
		match:   false,
	},
	{
		name:    "PatternSequence,",
		expr:    `[ a, b ]`,
		pattern: `[ PatternSequence(a, b) ]`,
		binding: "",
		match:   true,
	},
	{
		name:    "PatternSequence,binding",
		expr:    `[ a, b ]`,
		pattern: `[ Pattern(x, PatternSequence(a, b)) ]`,
		binding: "[ x:[a,b] ]",
		match:   true,
	},
	{
		name:    "MatchAny,sequence3,matchany2,matchstar",
		expr:    "[a]",
		pattern: "[ Pattern(x, MatchAny()), Pattern(y, MatchStar(MatchAny())), Pattern(z, MatchStar(MatchAny())) ]",
		binding: "[ x:a ]",
		match:   true,
	},
	{
		name:    "MatchExponential",
		expr:    "[a,a,a]",
		pattern: "[ MatchQuest(a),MatchQuest(a),MatchQuest(a),a,a,a ]",
		binding: "",
		match:   true,
	},
	{
		name:    "MatchExponential,binding",
		expr:    "[a,a,a]",
		pattern: "[ Pattern(x, PatternSequence(MatchQuest(a),MatchQuest(a),MatchQuest(a),a,a,a))]",
		binding: "[ x:[a,a,a] ]",
		match:   true,
	},
	{
		name:    "MatchSingle,binding,alt",
		expr:    "a",
		pattern: "Pattern(x, MatchAny())",
		binding: "[ x:a ]",
		match:   true,
	},
	{
		name:    "MatchExponential,binding,alt",
		expr:    "[a,a,a]",
		pattern: "Pattern(x, [MatchQuest(a),MatchQuest(a),MatchQuest(a),a,a,a])",
		binding: "[ x:[a,a,a] ]",
		match:   true,
	},
	{
		name:    "MatchAny,sublist,integer",
		expr:    "[ [ 1, 2 ] ]",
		pattern: "[ [ MatchHead(Integer), MatchHead(Integer) ] ]",
		binding: "",
		match:   true,
	},
	{
		name:    "MatchAny,sublist,integer,binding",
		expr:    "[ [ 1, 2 ] ]",
		pattern: "[ [ Pattern(x,MatchHead(Integer)), Pattern(y,MatchHead(Integer)) ] ]",
		binding: "[ x:1, y:2 ]",
		match:   true,
	},
	{
		name:    "MatchHead,sublist,nested",
		expr:    "[ [ 1, 2 ], 10 ]",
		pattern: "[ [ MatchHead(Integer), MatchHead(Integer) ], MatchHead(Integer) ]",
		binding: "",
		match:   true,
	},
	{
		name:    "MatchHead,sublist,nested,binding",
		expr:    "[ [ 1, 2 ], 10 ]",
		pattern: "[ [ Pattern(x,MatchHead(Integer)), Pattern(y,MatchHead(Integer)) ], Pattern(n,MatchHead(Integer)) ]",
		binding: "[x:1,y:2,n:10]",
		match:   true,
	},
	{
		name:    "RGBColor literal",
		expr:    "RGBColor(1,2,3)",
		pattern: "RGBColor(1,2,3)",
		binding: "",
		match:   true,
	},
	{
		name:    "RGBColor literal",
		expr:    "RGBColor(1,2,3)",
		pattern: "[1,2,3]",
		binding: "",
		match:   false,
	},
	{
		name:    "RGBColor pattern1",
		expr:    "RGBColor(1,2,3)",
		pattern: "RGBColor(MatchAny(),MatchAny(),MatchAny())",
		binding: "",
		match:   true,
	},
	{
		name:    "RGBColor pattern1",
		expr:    "RGBColor(1,2,3)",
		pattern: "[ MatchAny(),MatchAny(),MatchAny() ]",
		binding: "",
		match:   false,
	},
	/*
		{
			name:    "List with any head",
			expr:    "RGBColor(1,2,3)",
			pattern: "MatchAny(MatchStar(MatchAny()))",
			binding: "",
			match:   true,
		},

		{
			name:    "List with any head, 2 elements",
			expr:    "RGBColor(1,2)",
			pattern: "MatchAny(MatchAny(), MatchAny())",
			binding: "",
			match:   true,
		},

		{
			name:    "List with any head, 2 elements",
			expr:    "RGBColor(1)",
			pattern: "MatchAny(MatchAny(), MatchAny())",
			binding: "",
			match:   false,
		},
	*/
}

func TestSREM1(t *testing.T) {
	for _, tt := range cases {

		t.Run(tt.name, func(t *testing.T) {
			e := MustParse(tt.expr)
			p := MustParse(tt.pattern)
			c := NewCompiler()
			re := NewRegexp()
			prog := c.compileNFA(p)
			sub := NewCaptures(len(prog.Groups()))
			matched, bind := re.matchNfa(prog, e, sub)

			if matched != tt.match {
				t.Errorf("Expression %q with pattern %q was %v, expected %v",
					tt.expr, tt.pattern, matched, tt.match)
				return
			}
			if tt.binding == "" && (bind != nil && bind.Length() != 0) {
				// odd case that will never happen
				t.Errorf("Expected no bindings but got some: %s", bind.AsRules(prog.Groups()))
			}
			if bind != nil && bind.Length() == 0 && tt.binding != "" {
				t.Errorf("Expecting binding of %s, but got nothing", tt.binding)
			}
			// we know both are not nill
			if tt.binding != "" && bind.Length() != 0 {
				blist := MustParse(tt.binding)
				rlist := bind.AsRules(prog.Groups())
				if !blist.Equal(rlist) {
					t.Errorf("Bindings: expected %s, got %s", blist, rlist)
				}
			}

		})
	}
}

func TestSREM2(t *testing.T) {
	for _, tt := range cases {
		e := MustParse(tt.expr)
		p := MustParse(tt.pattern)
		c := NewCompiler()
		re := NewRegexp()

		if !c.Simple(p) {
			continue
		}

		t.Run(tt.name, func(t *testing.T) {
			prog := c.compileNFA(p)
			matched, bind := re.MatchM2(prog, e)

			if matched != tt.match {
				t.Errorf("Expression %q with pattern %q was %v, expected %v",
					tt.expr, tt.pattern, matched, tt.match)
				return
			}
			if tt.binding == "" && (bind != nil && bind.Length() != 0) {
				// odd case that will never happen
				t.Errorf("Expected no bindings but got some: %s", bind.AsRules(prog.Groups()))
			}
			if bind != nil && bind.Length() == 0 && tt.binding != "" {
				t.Errorf("Expecting binding of %s, but got nothing", tt.binding)
			}
			// we know both are not nill
			if tt.binding != "" && bind.Length() != 0 {
				blist := MustParse(tt.binding)
				rlist := bind.AsRules(prog.Groups())
				if !blist.Equal(rlist) {
					t.Errorf("Bindings: expected %s, got %s", blist, rlist)
				}
			}
		})
	}
}
func TestSREM3(t *testing.T) {
	for _, tt := range cases {
		e := MustParse(tt.expr)
		p := MustParse(tt.pattern)
		c := NewCompiler()
		re := NewRegexp()
		if !c.Simple(p) {
			continue
		}

		prog := c.compileOneStep(p)
		if !prog.IsOneStep() {
			t.Errorf("%s: M3 program is not one step!", tt.name)
		}
		matched, bind := re.MatchM3(prog, e)
		if matched != tt.match {
			t.Errorf("Expression %q with pattern %q was %v, expected %v",
				tt.expr, tt.pattern, matched, tt.match)
			return
		}
		if tt.binding == "" && (bind != nil && bind.Length() != 0) {
			// odd case that will never happen
			t.Errorf("Expected no bindings but got some: %s", bind.AsRules(prog.Groups()))
		}
		if bind != nil && bind.Length() == 0 && tt.binding != "" {
			t.Errorf("Expecting binding of %s, but got nothing", tt.binding)
		}
		// we know both are not nill
		if tt.binding != "" && bind.Length() != 0 {
			blist := MustParse(tt.binding)
			rlist := bind.AsRules(prog.Groups())
			if !blist.Equal(rlist) {
				t.Errorf("Bindings: expected %s, got %s", blist, rlist)
			}
		}
	}
}

func TestSREM4(t *testing.T) {
	for _, tt := range cases {
		e := MustParse(tt.expr)
		p := MustParse(tt.pattern)
		c := NewCompiler()
		re := NewRegexp()
		if !c.Simple(p) {
			continue
		}

		prog := c.compileOneStep(p)

		if !prog.IsOneStep() {
			t.Errorf("%s: M4 program is not one step !", tt.name)
		}

		matched, bind := re.MatchM4(prog, e)
		if matched != tt.match {
			t.Errorf("M4 Expression %q with pattern %q was %v, expected %v",
				tt.expr, tt.pattern, matched, tt.match)
			return
		}
		if tt.binding == "" && (bind != nil && bind.Length() != 0) {
			// odd case that will never happen
			t.Errorf("Expected no bindings but got some: %s", bind.AsRules(prog.Groups()))
		}
		if bind != nil && bind.Length() == 0 && tt.binding != "" {
			t.Errorf("Expecting binding of %s, but got nothing", tt.binding)
		}
		// we know both are not nill
		if tt.binding != "" && bind.Length() != 0 {
			blist := MustParse(tt.binding)
			rlist := bind.AsRules(prog.Groups())
			if !blist.Equal(rlist) {
				t.Errorf("Bindings: expected %s, got %v", blist, rlist)
			}
		}
	}
}

// TestExampleUsage demonstrates the s-expression regex system with practical examples
func BenchmarkSRE(b *testing.B) {

	cases := []tc{
		{
			name:    "Match1",
			expr:    "[a]",
			pattern: "[ MatchAny() ]",
			binding: "",
			match:   true,
		},
		{
			name:    "Match1,binding",
			expr:    "[a]",
			pattern: "[ Pattern(x, MatchAny()) ]",
			binding: "[ x:a ]",
			match:   true,
		},
		{
			name:    "Match3",
			expr:    "[a, b, c]",
			pattern: "[ MatchAny(), MatchAny(), MatchAny() ]",
			binding: "",
			match:   true,
		},
		{
			name:    "Match3,binding",
			expr:    "[a, b, c]",
			pattern: "[ Pattern(x, MatchAny()), Pattern(y, MatchAny()), Pattern(z, MatchAny()) ]",
			binding: "[ x:a, y:b, z:c ]",
			match:   true,
		},
		{
			name:    "MatchStar3",
			expr:    "[a, b, c]",
			pattern: "[MatchStar(MatchAny()) ]",
			binding: "",
			match:   true,
		},
		{
			name:    "MatchStar3,binding",
			expr:    "[a, b, c]",
			pattern: "[ Pattern(x,MatchStar(MatchAny())) ]",
			binding: "[ x:[a, b, c] ]",
			match:   true,
		},
		{
			name:    "MatchSublist2",
			expr:    "[ [ 1, 2 ] ]",
			pattern: "[ [ MatchHead(Integer), MatchHead(Integer) ] ]",
			binding: "",
			match:   true,
		},
		{
			name:    "MatchSublist2,binding",
			expr:    "[ [ 1, 2 ] ]",
			pattern: "[ [ Pattern(x,MatchHead(Integer)), Pattern(y,MatchHead(Integer)) ] ]",
			binding: "[ x:1, y:2 ]",
			match:   true,
		},
		{
			name:    "Match3Sublist2",
			expr:    "[ [ 1, 2 ], 10 ]",
			pattern: "[ [ MatchHead(Integer), MatchHead(Integer) ], MatchHead(Integer) ]",
			binding: "",
			match:   true,
		},
		{
			name:    "Match3Sublist2,binding",
			expr:    "[ [ 1, 2 ], 10 ]",
			pattern: "[ [ Pattern(x,MatchHead(Integer)), Pattern(y,MatchHead(Integer)) ], Pattern(n,MatchHead(Integer)) ]",
			binding: "[ x:1, y:2, n:10 ]",
			match:   true,
		},
	}

	for _, tt := range cases {
		typeNFA := true
		typeA := false
		typeB := true

		e := MustParse(tt.expr)

		p := MustParse(tt.pattern)
		list := p.(List)
		g := NewCaptures(3)

		if typeNFA {
			b.Run("M1,"+tt.name, func(b *testing.B) {
				c := NewCompiler()
				re := NewRegexp()
				prog := c.compileNFAList(list.Tail())
				for b.Loop() {
					re.matchNfaSequence(prog, e.(List).Tail(), g)
				}
			})
		}
		if typeA {
			b.Run("M2a,"+tt.name, func(b *testing.B) {
				c := NewCompiler()
				re := NewRegexp()
				prog := c.compileNFA(p)
				for b.Loop() {
					re.matchM2(prog, e, g)
				}
			})
		}
		if typeB {
			b.Run("M2b,"+tt.name, func(b *testing.B) {
				c := NewCompiler()
				re := NewRegexp()
				prog := c.compileNFAList(list.Tail())
				for b.Loop() {
					re.matchSequenceM2(prog, e.(List).Tail(), g)
				}
			})
		}
		/*
			if typeA {
				b.Run("M3a,"+tt.name, func(b *testing.B) {
					c := NewCompiler()
					re := NewRegexp()
					prog := c.compileOneStep(p)
					for b.Loop() {
						re.matchM3(prog, e, g)
					}
				})
			}
			if typeB {
				b.Run("M3b,"+tt.name, func(b *testing.B) {
					c := NewCompiler()
					re := NewRegexp()
					prog := c.compileListOneStep(list.Tail())
					for b.Loop() {
						re.matchSequenceM3(prog, e.(List).Tail(), g)
					}
				})
			}
		*/
		if typeA {
			b.Run("M4a,"+tt.name, func(b *testing.B) {
				c := NewCompiler()
				re := NewRegexp()
				prog := c.compileOneStep(p)
				for b.Loop() {
					re.matchM4(prog, e, g)
				}
			})
		}

		if typeB {
			b.Run("M4b,"+tt.name, func(b *testing.B) {
				c := NewCompiler()
				re := NewRegexp()
				prog := c.compileListOneStep(p.(List).Tail())
				for b.Loop() {
					re.matchListM4(prog, e, g)
				}
			})
		}
	}
}

func TestSREHack(t *testing.T) {
	tt := tc{
		name:    "MatchStar,sequence,single",
		expr:    "[ a ]",
		pattern: "[ MatchStar(MatchAny()) ]",
		binding: "",
		match:   true,
	}
	e := MustParse(tt.expr)
	p := MustParse(tt.pattern)
	c := NewCompiler()
	re := NewRegexp()

	prog := c.compileNFA(p)
	matched, _ := re.MatchM2(prog, e)
	if matched != tt.match {
		t.Error("Hack failed")
	}
}

func BenchmarkSREProfile(b *testing.B) {

	tt := tc{
		name:    "Match3Sublist2,binding",
		expr:    "[ [ 1, 2 ], 10 ]",
		pattern: "[ [ Pattern(x,MatchHead(Integer)), Pattern(y,MatchHead(Integer)) ], Pattern(n,MatchHead(Integer)) ]",
		binding: "[ x:1, y:2, n:10 ]",
		match:   true,
	}

	e := MustParse(tt.expr)
	elist := e.(List)
	args := elist.Tail()

	p := MustParse(tt.pattern)
	list := p.(List)
	c := NewCompiler()
	re := NewRegexp()
	groups := NewCaptures(3)
	prog := c.compileListOneStep(list.Tail())
	for b.Loop() {
		re.matchSequenceM4(prog, args, groups)
		//re.matchSequenceM3(prog, args, groups)
	}
}
func BenchmarkSRECrazy(b *testing.B) {

	for n := 1; n < 30; n++ {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {

			// make the expression, a list of n "1"s  1,1,1,1,....
			parts := make([]string, n)
			for i := 0; i < n; i++ {
				parts[i] = "1"
			}
			e := MustParse("[" + strings.Join(parts, ",") + "]")
			elist := e.(List)
			args := elist.Tail()

			// make pattern, n  1?,1?,... followed by n "1"
			// n =3 --> 1?,1?,1?,1,1,1

			parts = make([]string, 0, n*2)
			for i := 0; i < n; i++ {
				parts = append(parts, "MatchQuest(1)")
			}
			for i := 0; i < n; i++ {
				parts = append(parts, "1")
			}
			//p := MustParse("[" + strings.Join(parts, ",") +"]")
			p := MustParse("[Pattern(x,PatternSequence(" + strings.Join(parts, ",") + "))]")
			list := p.(List)

			c := NewCompiler()
			prog := c.compileNFAList(list.Tail())

			re := NewRegexp()
			for b.Loop() {
				ok, _ := re.MatchList(prog, args)
				if !ok {
					b.Errorf("Match failed")
				}
			}
		})
	}
}
