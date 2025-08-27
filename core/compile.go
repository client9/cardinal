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
func (c *Compile) Compile(e Expr) Prog {
	if c.Simple(e) {
		return c.compileOneStep(e)
	}
	return c.compileNFA(e)
}

func (c *Compile) CompileList(exprs []Expr) Prog {

	if c.SimpleList(exprs) {
		return c.compileListOneStep(exprs)
	}
	return c.compileNFAList(exprs)
}

func (c *Compile) reset() {
	c.pc = 0
	c.ops = c.ops[:0]
	c.groups = c.groups[:0]
}

func (c *Compile) addLink(origin int32, next int32) {
	i := &c.ops[origin]
	i.Next = next
}

func (c *Compile) addAlt(origin int32, alt int32) {
	i := &c.ops[origin]
	i.Alt = alt
}

func (c *Compile) add(x Inst) int32 {
	x.Id = c.pc
	c.ops = append(c.ops, x)
	c.pc += 1
	return x.Id
}

func (c *Compile) compileNFA(e Expr) Prog {
	c.reset()

	c.groups = c.getGroups(e, nil)

	c.emit(e)
	c.add(Inst{
		Op: InstMatchEnd,
	})

	return Prog{
		Inst:   c.ops,
		groups: c.groups,
	}
}

func (c *Compile) compileNFAList(exprs []Expr) Prog {
	c.reset()

	c.groups = c.getGroupsList(exprs, nil)

	for _, e := range exprs {
		c.emit(e)
	}
	c.add(Inst{
		Op: InstMatchEnd,
	})

	return Prog{
		Inst:   c.ops,
		groups: c.getGroupsList(exprs, nil),
	}
}

// Hack until we can do this in the instruction
func (c *Compile) getSlot(name string) int32 {
	for i, s := range c.groups {
		if s == name {
			return int32(i)
		}
	}
	return -1
}

func (c *Compile) getGroups(e Expr, names []string) []string {
	if list, ok := e.(List); ok {
		switch list.HeadAtom() {
		case atom.Pattern:
			args := list.Tail()
			// args[0] is the binding name
			// args[1] is the pattern
			names = append(names, args[0].String())
			names = c.getGroups(args[1], names)
			return names
		case atom.PatternSequence, atom.List:
			return c.getGroupsList(list.Tail(), names)
		}
	}

	return names

}

func (c *Compile) getGroupsList(exprs []Expr, names []string) []string {
	for _, arg := range exprs {
		names = c.getGroups(arg, names)
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
		case atom.MatchStar, atom.MatchPlus, atom.MatchQuest, atom.BlankSequence, atom.BlankNullSequence, atom.Optional:
			args := list.Tail()
			return c.Simple(args[0])
		case atom.MatchHead, atom.MatchAny, atom.Blank:
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
		if c.isSequencePattern(e[i]) {
			s = false
			break
		}
	}
	if s {
		return true
	}

	s = true
	for i = 0; i < len(e); i++ {
		if c.isSequencePattern(e[i]) {
			break
		}
	}
	for ; i < len(e); i++ {
		if !c.isZeroPattern(e[i]) {
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

func (c *Compile) isZeroPattern(e Expr) bool {
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

func (c *Compile) isSequencePattern(e Expr) bool {
	if list, ok := e.(List); ok {

		switch list.HeadAtom() {
		case atom.MatchStar, atom.MatchPlus, atom.MatchQuest:
			return true

		// MMA compatible
		case atom.BlankNullSequence, atom.BlankSequence, atom.Optional:
			return true

		case atom.Pattern, atom.PatternSequence, atom.List:
			for _, a := range list.Tail() {
				if c.isSequencePattern(a) {
					return true
				}
			}
			return false
		}
	}
	return false
}

func (c *Compile) compileOneStep(e Expr) Prog {
	c.reset()
	c.groups = c.getGroups(e, nil)
	c.emitOneStep(e)
	c.add(Inst{
		Op: InstMatchEnd,
	})
	eof := c.add(Inst{
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
		onestep: true,
		groups:  c.groups,
	}
}
func (c *Compile) compileListOneStep(exprs []Expr) Prog {
	c.reset()
	c.groups = c.getGroupsList(exprs, nil)
	for _, e := range exprs {
		c.emitOneStep(e)
	}
	c.add(Inst{
		Op: InstMatchEnd,
	})
	eof := c.add(Inst{
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
		onestep: true,
		groups:  c.getGroupsList(exprs, nil),
	}
}
func (c *Compile) IsListLiteral(list List) bool {
	// TODO: if head is a pattern

	for _, a := range list.Tail() {
		if alist, ok := a.(List); ok {
			switch alist.HeadAtom() {
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
			default:
				if !c.IsListLiteral(alist) {
					return false
				}
			}
		}
	}
	return true
}

func (c *Compile) emitOneStep(e Expr) {

	list, ok := e.(List)

	// not a list, some other primitive literal
	if !ok {
		// some literal
		op := c.add(Inst{
			Op:  InstMatchLiteral,
			Val: e,
		})
		c.addLink(op, c.pc)
		c.addAlt(op, -1)
		return
	}

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
		c.emitOneStep(arg)
	case atom.BlankSequence:
		arg := ListFirstArg(e)
		if arg == nil {
			arg = ListFrom(atom.MatchAny)
		} else {
			arg = ListFrom(atom.MatchHead, arg)
		}
		c.emitOneStep(ListFrom(atom.MatchPlus, arg))
	case atom.BlankNullSequence:
		arg := ListFirstArg(e)
		if arg == nil {
			arg = ListFrom(atom.MatchAny)
		} else {
			arg = ListFrom(atom.MatchHead, arg)
		}
		c.emitOneStep(ListFrom(atom.MatchStar, arg))
	case atom.Optional:
		// TODO default value
		arg := ListFirstArg(e)
		if arg == nil {
			arg = ListFrom(atom.MatchAny)
		} else {
			arg = ListFrom(atom.MatchHead, arg)
		}
		c.emitOneStep(ListFrom(atom.MatchQuest, arg))

	case atom.Pattern:
		// Pattern("x", expression)
		args := list.Tail()
		name := args[0]
		slot := c.getSlot(name.String())
		expr := args[1]
		cstart := c.add(Inst{
			Op:  InstCaptureStart,
			Alt: slot,
		})
		c.addLink(cstart, c.pc)
		c.emitOneStep(expr)
		cend := c.add(Inst{
			Op:  InstCaptureStart,
			Alt: slot,
		})
		c.addLink(cend, c.pc)
	case atom.PatternSequence:
		args := list.Tail()
		for _, arg := range args {
			c.emitOneStep(arg)
		}

	case atom.MatchHead:
		val := list.Tail()[0]
		op := c.add(Inst{
			Op:   InstMatchHead,
			Name: val.String(),
		})
		c.addLink(op, c.pc)
		c.addAlt(op, -1)
	case atom.MatchAny:
		op := c.add(Inst{
			Op: InstMatchAny,
		})
		c.addLink(op, c.pc)
		c.addAlt(op, -1)
	case atom.MatchPlus:
		arg := ListFirstArg(e)
		// only has a single argument
		c.emitOneStep(arg)

		L1 := c.pc
		c.emitOneStep(arg)
		op := c.pc - 1
		L3 := c.pc
		c.addLink(op, L1)
		c.addAlt(op, L3)
	case atom.MatchQuest:
		arg := list.Tail()[0]
		// only has a single argument
		c.emitOneStep(arg)
		op := c.pc - 1
		L1 := c.pc
		c.addLink(op, L1)
		c.addAlt(op, L1)
	case atom.MatchStar:
		L1 := c.pc

		// only has a single argument
		c.emitOneStep(list.Tail()[0])
		op := c.pc - 1
		L3 := c.pc
		c.addLink(op, L1)
		c.addAlt(op, L3)
	default:
		if c.IsListLiteral(list) {
			op := c.add(Inst{
				Op:  InstMatchLiteral,
				Val: e,
			})
			c.addLink(op, c.pc)
			c.addAlt(op, -1)
			return
		}

		// list-like object that contains pattern primitives
		//
		// figure out next program for list
		nc := NewCompiler()
		newprog := nc.compileListOneStep(list.Tail())

		op := c.add(Inst{
			Op:   InstMatchList,
			Data: newprog,
			Name: list.Head(),
		})
		c.addLink(op, c.pc)
		c.addAlt(op, -1)
	}
}

func (c *Compile) emit(e Expr) {

	list, ok := e.(List)

	// not a list, some other primitive literal
	if !ok {
		// some literal
		op := c.add(Inst{
			Op:  InstMatchLiteral,
			Val: e,
		})
		c.addLink(op, c.pc)
		c.addAlt(op, -1)
		return
	}

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
		c.emit(arg)
	case atom.BlankSequence:
		var arg Expr
		if list.Length() == 0 {
			arg = ListFrom(atom.MatchAny)
		} else {
			arg = ListFrom(atom.MatchHead, list.Tail()[0])
		}
		c.emit(ListFrom(atom.MatchPlus, arg))
	case atom.BlankNullSequence:
		var arg Expr
		if list.Length() == 0 {
			arg = ListFrom(atom.MatchAny)
		} else {
			arg = ListFrom(atom.MatchHead, list.Tail()[0])
		}
		c.emit(ListFrom(atom.MatchStar, arg))
	case atom.Optional:
		// TODO default value
		var arg Expr
		if list.Length() == 0 {
			arg = ListFrom(atom.MatchAny)
		} else {
			arg = ListFrom(atom.MatchHead, list.Tail()[0])
		}
		c.emit(ListFrom(atom.MatchQuest, arg))

	case atom.Pattern:
		// Pattern("x", expression)
		args := list.Tail()
		name := args[0]
		slot := c.getSlot(name.String())
		expr := args[1]
		cstart := c.add(Inst{
			Op: InstCaptureStart,
			//Name: name.String(),
			Alt: slot,
		})
		c.addLink(cstart, c.pc)
		c.emit(expr)
		cend := c.add(Inst{
			Op:  InstCaptureStart,
			Alt: slot,
		})
		c.addLink(cend, c.pc)
	case atom.PatternSequence:
		args := list.Tail()
		for _, arg := range args {
			c.emit(arg)
		}
	case atom.MatchHead:
		val := list.Tail()[0]

		op := c.add(Inst{
			Op: InstMatchHead,
			// head: val.String(),
			Name: val.String(),
		})
		c.addLink(op, c.pc)
	case atom.MatchAny:
		// this has a dangling pointer
		// it will be fixed at the end
		op := c.add(Inst{
			Op: InstMatchAny,
		})
		c.addLink(op, c.pc)
	case atom.MatchPlus:
		current := c.pc
		list, _ := e.(List)
		// only has a single argument
		c.emit(list.Tail()[0])

		op := c.add(Inst{
			Op: InstSplit,
		})
		c.addLink(op, current)
		c.addAlt(op, c.pc)
	case atom.MatchQuest:
		op := c.add(Inst{
			Op: InstSplit,
		})
		L1 := c.pc
		list, _ := e.(List)
		c.emit(list.Tail()[0])
		L2 := c.pc
		c.addLink(op, L1)
		c.addAlt(op, L2)
	case atom.MatchStar:
		//L1 := c.pc
		op := c.add(Inst{
			Op: InstSplit,
		})
		L2 := c.pc
		list, _ := e.(List)
		c.emit(list.Tail()[0])
		op2 := c.add(Inst{
			Op: InstJump,
		})
		L3 := c.pc
		c.addLink(op2, op)
		c.addLink(op, L2)
		c.addAlt(op, L3)
	default:
		if c.IsListLiteral(list) {
			// has no pattern operators, match as literal
			op := c.add(Inst{
				Op:  InstMatchLiteral,
				Val: e,
			})
			c.addLink(op, c.pc)
			c.addAlt(op, -1)
			return
		}

		// Some other list that has pattern operators in it
		//
		// figure out next program for list
		nc := NewCompiler()
		newprog := nc.compileNFAList(list.Tail())

		op := c.add(Inst{
			Op:   InstMatchList,
			Data: newprog,
			Name: list.Head(),
		})
		c.addLink(op, c.pc)
		c.addAlt(op, -1)
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
func (p Prog) getSlot(name string) int32 {
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
