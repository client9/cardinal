package core

import (
	"fmt"
	"strings"

	"github.com/client9/sexpr/core/atom"
)

type Compile struct {
	pc     int32 // program counter
	ops    []Inst
	groups []string
}

func NewCompiler() *Compile {
	return &Compile{}
}
func (c *Compile) reset() {
	c.pc = 0
	c.ops = nil
	c.groups = nil
}

func (c *Compile) AddLink(origin int32, next int32) {
	i := &c.ops[origin]
	i.Next = next
}

func (c *Compile) AddAlt(origin int32, alt int32) {
	i := &c.ops[origin]
	i.Alt = alt
}

func (c *Compile) Add(x Inst) int32 {
	x.Id = c.pc
	c.ops = append(c.ops, x)
	c.pc += 1
	return x.Id
}

func (c *Compile) Compile(e Expr) Prog {
	c.reset()

	c.groups = c.Groups(e, nil)

	c.Emit(e)
	c.Add(Inst{
		Op: InstMatchEnd,
	})

	return Prog{
		Inst:    c.ops,
		onestep: c.Simple(e),
		groups:  c.groups,
	}
}

func (c *Compile) CompileList(exprs []Expr) Prog {
	c.reset()

	c.groups = c.GroupsList(exprs, nil)

	for _, e := range exprs {
		c.Emit(e)
	}
	c.Add(Inst{
		Op: InstMatchEnd,
	})

	return Prog{
		Inst:    c.ops,
		onestep: c.SimpleList(exprs),
		groups:  c.GroupsList(exprs, nil),
	}
}

// Hack until we can do this in the instruction
func (c *Compile) GetSlot(name string) int32 {
	for i, s := range c.groups {
		if s == name {
			return int32(i)
		}
	}
	return -1
}

func (c *Compile) Groups(e Expr, names []string) []string {
	if list, ok := e.(List); ok {
		switch list.HeadAtom() {
		case atom.Pattern:
			args := list.Tail()
			// args[0] is the binding name
			// args[1] is the pattern
			names = append(names, args[0].String())
			names = c.Groups(args[1], names)
			return names
		case atom.PatternSequence, atom.List:
			return c.GroupsList(list.Tail(), names)
		}
	}

	return names

}

func (c *Compile) GroupsList(exprs []Expr, names []string) []string {
	for _, arg := range exprs {
		names = c.Groups(arg, names)
	}
	return names
}

// for a single atomic Expr
func (c *Compile) Simple(e Expr) bool {
	if list, ok := e.(List); ok {
		switch list.HeadAtom() {
		case atom.Pattern:
			args := list.Tail()
			// args[0] is the binding name
			// args[1] is the pattern
			return c.Simple(args[1])
		case atom.MatchStar, atom.MatchPlus, atom.MatchQuest:
			args := list.Tail()
			return c.Simple(args[0])
		case atom.MatchHead, atom.MatchAny:
			return true
		case atom.PatternSequence, atom.List:
			return c.SimpleList(list.Tail())
		}
	}
	return true
}

// Order of elements matters
// *, 1 --> not simple
// 1, * --> Simple
// 1, *, *, * -> simple
func (c *Compile) SimpleList(e []Expr) bool {

	// simple check... no sequence patterns, except for the last one

	if len(e) == 1 {
		return c.Simple(e[0])
	}

	s := true
	i := 0

	for i = 0; i < len(e)-1; i++ {
		if c.IsSequencePattern(e[i]) {
			s = false
			break
		}
	}
	if s {
		return true
	}

	s = true
	for i = 0; i < len(e); i++ {
		if c.IsSequencePattern(e[i]) {
			break
		}
	}
	for ; i < len(e); i++ {
		if !c.IsZeroPattern(e[i]) {
			s = false
			break
		}
	}

	if s {
		return true
	}

	return false
	// a sequence pattern is ok if it is followed by pattern
	// that has a different predicate
	//
	// _Integer*, _String   --> ok
	// _Integer*, "foo"     --> ok
	// _Integer*, _Integer  --> not ok
	// _Integer*, 1         --> not ok
	// _*, 1                --> not ok
	// do other checks

	return false
}

func (c *Compile) IsZeroPattern(e Expr) bool {
	if list, ok := e.(List); ok {
		switch list.HeadAtom() {

		case atom.MatchStar, atom.MatchQuest:
			return true

		// MMA compatible
		case atom.BlankNullSequence, atom.Optional:
			return true

		}
	}
	return false
}

func (c *Compile) IsSequencePattern(e Expr) bool {
	if list, ok := e.(List); ok {

		switch list.HeadAtom() {
		case atom.MatchStar, atom.MatchPlus, atom.MatchQuest:
			return true

		// MMA compatible
		case atom.BlankNullSequence, atom.BlankSequence, atom.Optional:
			return true

		case atom.Pattern, atom.PatternSequence, atom.List:
			for _, a := range list.Tail() {
				if c.IsSequencePattern(a) {
					return true
				}
			}
			return false
		}
	}
	return false
}

func (c *Compile) CompileOneStep(e Expr) Prog {
	c.reset()
	c.groups = c.Groups(e, nil)
	c.EmitOneStep(e)
	c.Add(Inst{
		Op: InstMatchEnd,
	})
	eof := c.Add(Inst{
		Op: InstFail,
	})

	// change all Alts that have -1 to Fail.
	for i, op := range c.ops {
		if op.Alt == -1 {
			op.Alt = eof
			c.ops[i] = op
		}
	}

	return Prog{
		Inst:    c.ops,
		onestep: c.Simple(e),
		groups:  c.groups,
	}
}
func (c *Compile) CompileListOneStep(exprs []Expr) Prog {
	c.reset()
	c.groups = c.GroupsList(exprs, nil)
	for _, e := range exprs {
		c.EmitOneStep(e)
	}
	c.Add(Inst{
		Op: InstMatchEnd,
	})
	eof := c.Add(Inst{
		Op: InstFail,
	})

	// change all Alts that have -1  to MatchFail
	for i, op := range c.ops {
		if op.Alt == -1 {
			op.Alt = eof
			c.ops[i] = op
		}
	}

	return Prog{
		Inst:    c.ops,
		onestep: c.SimpleList(exprs),
		groups:  c.GroupsList(exprs, nil),
	}
}
func (c *Compile) IsLiteral(e Expr) bool {
	list, ok := e.(List)
	if !ok {
		return true
	}

	switch list.HeadAtom() {
	case atom.Pattern, atom.PatternSequence:
		return false

	// mma primitives
	case atom.Blank, atom.BlankSequence, atom.BlankNullSequence, atom.Optional:
		return false

	// low level primitives
	case atom.MatchStar, atom.MatchPlus, atom.MatchQuest,
		atom.MatchAny, atom.MatchHead, atom.MatchLiteral:
		return false

	// TBD
	case atom.MatchOr:
		return false
	}

	for _, a := range list.Tail() {
		if !c.IsLiteral(a) {
			return false
		}
	}
	return true
}

func (c *Compile) EmitOneStep(e Expr) {

	if c.IsLiteral(e) {
		op := c.Add(Inst{
			Op:  InstMatchLiteral,
			Val: e,
		})
		c.AddLink(op, c.pc)
		c.AddAlt(op, -1)
		return
	}

	if list, ok := e.(List); ok {
		switch list.HeadAtom() {

		case atom.Blank:
			// Blank[]       --> MatchAny()
			// Blank[symbol] --> MatchHead(symbol)

			arg := ListFirstArg(e)
			if arg == nil {
				arg = ListFrom(atom.MatchAny)
			} else {
				arg = ListFrom(atom.MatchHead, arg)
			}
			c.EmitOneStep(arg)
		case atom.BlankSequence:
			arg := ListFirstArg(e)
			if arg == nil {
				arg = ListFrom(atom.MatchAny)
			} else {
				arg = ListFrom(atom.MatchHead, arg)
			}
			c.EmitOneStep(ListFrom(atom.MatchPlus, arg))
		case atom.BlankNullSequence:
			arg := ListFirstArg(e)
			if arg == nil {
				arg = ListFrom(atom.MatchAny)
			} else {
				arg = ListFrom(atom.MatchHead, arg)
			}
			c.EmitOneStep(ListFrom(atom.MatchStar, arg))
		case atom.Optional:
			// TODO default value
			arg := ListFirstArg(e)
			if arg == nil {
				arg = ListFrom(atom.MatchAny)
			} else {
				arg = ListFrom(atom.MatchHead, arg)
			}
			c.EmitOneStep(ListFrom(atom.MatchQuest, arg))

		case atom.Pattern:
			// Pattern("x", expression)
			args := list.Tail()
			name := args[0]
			slot := c.GetSlot(name.String())
			expr := args[1]
			cstart := c.Add(Inst{
				Op:  InstCaptureStart,
				Alt: slot,
			})
			c.AddLink(cstart, c.pc)
			c.EmitOneStep(expr)
			cend := c.Add(Inst{
				Op:  InstCaptureStart,
				Alt: slot,
			})
			c.AddLink(cend, c.pc)
		case atom.PatternSequence:
			args := list.Tail()
			for _, arg := range args {
				c.EmitOneStep(arg)
			}

		case atom.MatchHead:
			val := list.Tail()[0]
			op := c.Add(Inst{
				Op:   InstMatchHead,
				Name: val.String(),
			})
			c.AddLink(op, c.pc)
			c.AddAlt(op, -1)
		case atom.MatchAny:
			op := c.Add(Inst{
				Op: InstMatchAny,
			})
			c.AddLink(op, c.pc)
			c.AddAlt(op, -1)
		case atom.MatchPlus:
			arg := ListFirstArg(e)
			// only has a single argument
			c.EmitOneStep(arg)

			L1 := c.pc
			c.EmitOneStep(arg)
			op := c.pc - 1
			L3 := c.pc
			c.AddLink(op, L1)
			c.AddAlt(op, L3)
		case atom.MatchQuest:
			arg := list.Tail()[0]
			// only has a single argument
			c.EmitOneStep(arg)
			op := c.pc - 1
			L1 := c.pc
			c.AddLink(op, L1)
			c.AddAlt(op, L1)
		case atom.MatchStar:
			L1 := c.pc

			// only has a single argument
			c.EmitOneStep(list.Tail()[0])
			op := c.pc - 1
			L3 := c.pc
			c.AddLink(op, L1)
			c.AddAlt(op, L3)
		default:
			// list-like object that contains pattern primitives
			//
			// figure out next program for list
			nc := NewCompiler()
			newprog := nc.CompileListOneStep(list.Tail())

			op := c.Add(Inst{
				Op:   InstMatchList,
				Data: newprog,
				Name: list.Head(),
			})
			c.AddLink(op, c.pc)
			c.AddAlt(op, -1)
		}
	}
}

func (c *Compile) Emit(e Expr) {
	if c.IsLiteral(e) {
		// Not a pattern operator, match as literal
		op := c.Add(Inst{
			Op:  InstMatchLiteral,
			Val: e,
		})
		c.AddLink(op, c.pc)
		c.AddAlt(op, -1)
		return
	}

	list := e.(List)
	switch list.HeadAtom() {

	// MMA
	case atom.Blank:
		// Blank[]       --> MatchAny()
		// Blank[symbol] --> MatchHead(symbol)

		var arg Expr
		if list.Length() == 0 {
			arg = ListFrom(atom.MatchAny)
		} else {
			arg = ListFrom(atom.MatchHead, list.Tail()[0])
		}
		c.Emit(arg)
	case atom.BlankSequence:
		var arg Expr
		if list.Length() == 0 {
			arg = ListFrom(atom.MatchAny)
		} else {
			arg = ListFrom(atom.MatchHead, list.Tail()[0])
		}
		c.Emit(ListFrom(atom.MatchPlus, arg))
	case atom.BlankNullSequence:
		var arg Expr
		if list.Length() == 0 {
			arg = ListFrom(atom.MatchAny)
		} else {
			arg = ListFrom(atom.MatchHead, list.Tail()[0])
		}
		c.Emit(ListFrom(atom.MatchStar, arg))
	case atom.Optional:
		// TODO default value
		var arg Expr
		if list.Length() == 0 {
			arg = ListFrom(atom.MatchAny)
		} else {
			arg = ListFrom(atom.MatchHead, list.Tail()[0])
		}
		c.Emit(ListFrom(atom.MatchQuest, arg))

	case atom.Pattern:
		// Pattern("x", expression)
		args := list.Tail()
		name := args[0]
		slot := c.GetSlot(name.String())
		expr := args[1]
		cstart := c.Add(Inst{
			Op: InstCaptureStart,
			//Name: name.String(),
			Alt: slot,
		})
		c.AddLink(cstart, c.pc)
		c.Emit(expr)
		cend := c.Add(Inst{
			Op:  InstCaptureStart,
			Alt: slot,
		})
		c.AddLink(cend, c.pc)
	case atom.PatternSequence:
		args := list.Tail()
		for _, arg := range args {
			c.Emit(arg)
		}
	case atom.MatchHead:
		val := list.Tail()[0]

		op := c.Add(Inst{
			Op: InstMatchHead,
			// head: val.String(),
			Name: val.String(),
		})
		c.AddLink(op, c.pc)
	case atom.MatchAny:
		// this has a dangling pointer
		// it will be fixed at the end
		op := c.Add(Inst{
			Op: InstMatchAny,
		})
		c.AddLink(op, c.pc)
	case atom.MatchPlus:
		current := c.pc
		list, _ := e.(List)
		// only has a single argument
		c.Emit(list.Tail()[0])

		op := c.Add(Inst{
			Op: InstSplit,
		})
		c.AddLink(op, current)
		c.AddAlt(op, c.pc)
	case atom.MatchQuest:
		op := c.Add(Inst{
			Op: InstSplit,
		})
		L1 := c.pc
		list, _ := e.(List)
		c.Emit(list.Tail()[0])
		L2 := c.pc
		c.AddLink(op, L1)
		c.AddAlt(op, L2)
	case atom.MatchStar:
		//L1 := c.pc
		op := c.Add(Inst{
			Op: InstSplit,
		})
		L2 := c.pc
		list, _ := e.(List)
		c.Emit(list.Tail()[0])
		op2 := c.Add(Inst{
			Op: InstJump,
		})
		L3 := c.pc
		c.AddLink(op2, op)
		c.AddLink(op, L2)
		c.AddAlt(op, L3)
	default:
		// figure out next program for list
		nc := NewCompiler()
		newprog := nc.CompileList(list.Tail())

		op := c.Add(Inst{
			Op:   InstMatchList,
			Data: newprog,
			Name: list.Head(),
		})
		c.AddLink(op, c.pc)
		c.AddAlt(op, -1)
		return
	}
}

type InstOp uint8

const (
	InstMatchLiteral = iota
	InstMatchHead
	InstMatchAny
	InstMatchList
	InstSplit
	InstJump
	InstCaptureStart
	InstCaptureEnd
	InstMatchEnd
	InstFail
)

var instOpNames = []string{
	"InstMatchLiteral",
	"InstMatchHead",
	"InstMatchAny",
	"InstMatchList",
	"InstSplit",
	"InstJump",
	"InstCaptureStart",
	"InstCaptureEnd",
	"InstMatchEnd",
	"InstFail",
}

func (i InstOp) String() string {
	if uint(i) >= uint(len(instOpNames)) {
		return ""
	}
	return instOpNames[i]
}

type Inst struct {
	Op   InstOp
	Id   int32  // More for debugging
	Next int32  // Everyone
	Alt  int32  // Split, Capture
	Val  Expr   // use in MatchLiteral only?
	Name string // MatchHead only
	Data any    // Predictates, MatchList
}

func (i Inst) String() string {
	return fmt.Sprintf("id=%d [%s %d %d] ", i.Id, i.Op.String(), i.Next, i.Alt)
}

// a "Program" is a list of a Instructions
type Prog struct {
	Inst    []Inst
	groups  []string
	onestep bool
}

func (p Prog) IsOneStep() bool {
	return p.onestep
}

func (p Prog) Groups() []string {
	return p.groups
}

// Hack until we can do this in the instruction
func (p Prog) GetSlot(name string) int32 {
	for i, s := range p.groups {
		if s == name {
			return int32(i)
		}
	}
	return -1
}

func (p Prog) First() int32 {
	return 0
}

func (p Prog) Length() int {
	return len(p.Inst)
}

func (p Prog) String() string {
	out := "[]Inst{\n"

	for _, inst := range p.Inst {
		out += fmt.Sprintf("   %d: %s\n", inst.Id, inst)
	}

	out += "}\n"
	return out
}

func (p Prog) dump(n int) {
	indent := strings.Repeat(" ", n)
	for _, inst := range p.Inst {
		fmt.Printf("%s%d: %s\n", indent, inst.Id, inst)
		if inst.Op == InstMatchList {
			p2 := inst.Data.(Prog)
			p2.dump(n + 1)
		}
	}
}
func (p Prog) Dump() {
	p.dump(0)
}
