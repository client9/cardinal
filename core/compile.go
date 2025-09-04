package core

import (
	"fmt"
	"strings"
)

type Compile struct {
	pc     int32 // program counter
	ops    []Inst
	groups []Symbol
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
func (c *Compile) getSlot(name Symbol) int32 {
	for i, s := range c.groups {
		if s == name {
			return int32(i)
		}
	}
	return -1
}

func (c *Compile) getGroups(e Expr, names []Symbol) []Symbol {
	list, ok := e.(List)
	if !ok {
		return names
	}

	// Head need not be a symbol
	head := list.Head()
	sym, ok := head.(Symbol)
	if !ok {
		// it's a list-like thing, but has functional head
		// MatchAny()(1,2,3)

		// ignore the head, and go into the tail
		return c.getGroupsList(list.Tail(), names)

	}
	switch sym {
	case symbolPattern:
		args := list.Tail()
		// args[0] is the binding name
		// args[1] is the pattern
		names = append(names, args[0].(Symbol))
		names = c.getGroups(args[1], names)
		return names
	case symbolPatternSequence, symbolList:
		return c.getGroupsList(list.Tail(), names)
	}

	return names
}

func (c *Compile) getGroupsList(exprs []Expr, names []Symbol) []Symbol {
	for _, arg := range exprs {
		names = c.getGroups(arg, names)
	}
	return names
}

// for a single atomic Expr
func (c *Compile) Simple(e Expr) bool {
	list, ok := e.(List)
	if !ok {
		return true
	}

	head := list.Head()
	sym, ok := head.(Symbol)
	if !ok {
		// we'll assume this is Blank or MatchAny
		return c.SimpleList(list.Tail())
	}

	// normal list with symbol head
	switch sym {
	case symbolPattern:
		args := list.Tail()
		// args[0] is the binding name
		// args[1] is the pattern
		return c.Simple(args[1])
	case symbolMatchStar, symbolMatchPlus, symbolMatchQuest,
		symbolBlankSequence, symbolBlankNullSequence, symbolOptional:
		args := list.Tail()
		if len(args) == 0 {
			return true
		}
		return c.Simple(args[0])
	case symbolMatchHead, symbolMatchAny, symbolBlank:
		return true
	case symbolPatternSequence, symbolList:
		return c.SimpleList(list.Tail())
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
		switch list.Head() {

		case symbolMatchStar, symbolMatchQuest:
			return true

		// MMA compatible
		case symbolBlankNullSequence, symbolOptional:
			return true

		}
	}
	return false
}

func (c *Compile) isSequencePattern(e Expr) bool {
	if list, ok := e.(List); ok {

		switch list.Head() {
		case symbolMatchStar, symbolMatchPlus, symbolMatchQuest:
			return true

		// MMA compatible
		case symbolBlankNullSequence, symbolBlankSequence, symbolOptional:
			return true

		case symbolPattern, symbolPatternSequence, symbolList:
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
			switch alist.Head() {
			case symbolPattern, symbolPatternSequence:
				return false

			// mma primitives
			case symbolBlank, symbolBlankSequence, symbolBlankNullSequence, symbolOptional:
				return false

			// low level primitives
			case symbolMatchStar, symbolMatchPlus, symbolMatchQuest,
				symbolMatchAny, symbolMatchHead, symbolMatchLiteral:
				return false

			// TBD
			case symbolMatchOr:
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

	head := list.Head()
	sym, ok := head.(Symbol)

	if !ok {
		// assume it's a Blank() or MatchAny()

		nc := NewCompiler()
		newprog := nc.compileListOneStep(list.Tail())

		op := c.add(Inst{
			Op:   InstMatchList,
			Data: newprog,
			Val:  symbolNull,
		})
		c.addLink(op, c.pc)
		c.addAlt(op, -1)
		return
	}

	switch sym {

	case symbolBlank:
		// Blank[]       --> MatchAny()
		// Blank[symbol] --> MatchHead(symbol)

		arg := ListFirstArg(e)
		if arg == nil {
			arg = ListFrom(symbolMatchAny)
		} else {
			arg = ListFrom(symbolMatchHead, arg)
		}
		c.emitOneStep(arg)
	case symbolBlankSequence:
		arg := ListFirstArg(e)
		if arg == nil {
			arg = ListFrom(symbolMatchAny)
		} else {
			arg = ListFrom(symbolMatchHead, arg)
		}
		c.emitOneStep(ListFrom(symbolMatchPlus, arg))
	case symbolBlankNullSequence:
		arg := ListFirstArg(e)
		if arg == nil {
			arg = ListFrom(symbolMatchAny)
		} else {
			arg = ListFrom(symbolMatchHead, arg)
		}
		c.emitOneStep(ListFrom(symbolMatchStar, arg))
	case symbolOptional:
		// TODO default value
		arg := ListFirstArg(e)
		if arg == nil {
			arg = ListFrom(symbolMatchAny)
		} else {
			arg = ListFrom(symbolMatchHead, arg)
		}
		c.emitOneStep(ListFrom(symbolMatchQuest, arg))

	case symbolPattern:
		// Pattern("x", expression)
		args := list.Tail()
		name := args[0].(Symbol)
		slot := c.getSlot(name)
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
	case symbolPatternSequence:
		args := list.Tail()
		for _, arg := range args {
			c.emitOneStep(arg)
		}

	case symbolMatchHead:
		op := c.add(Inst{
			Op:  InstMatchHead,
			Val: list.Tail()[0],
		})
		c.addLink(op, c.pc)
		c.addAlt(op, -1)
	case symbolMatchAny:
		op := c.add(Inst{
			Op: InstMatchAny,
		})
		c.addLink(op, c.pc)
		c.addAlt(op, -1)
	case symbolMatchPlus:
		arg := ListFirstArg(e)
		// only has a single argument
		c.emitOneStep(arg)

		L1 := c.pc
		c.emitOneStep(arg)
		op := c.pc - 1
		L3 := c.pc
		c.addLink(op, L1)
		c.addAlt(op, L3)
	case symbolMatchQuest:
		arg := list.Tail()[0]
		// only has a single argument
		c.emitOneStep(arg)
		op := c.pc - 1
		L1 := c.pc
		c.addLink(op, L1)
		c.addAlt(op, L1)
	case symbolMatchStar:
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
			Val:  list.Head(),
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

	//
	head := list.Head()
	sym, ok := head.(Symbol)
	if !ok {
		// head is not symbol
		// likely _[...]  "list of any head"
		// so it better be a Blank() or MatchAny()
		fn := head.Head()
		if fn == symbolBlank || fn == symbolMatchAny {

			nc := NewCompiler()
			newprog := nc.compileNFAList(list.Tail())

			op := c.add(Inst{
				Op:   InstMatchList,
				Data: newprog,
				Val:  symbolNull,
			})
			c.addLink(op, c.pc)
			c.addAlt(op, -1)
			return

		}
		panic("Unknown pattern type")

	}

	switch sym {

	// MMA
	case symbolBlank:
		// Blank[]       --> MatchAny()
		// Blank[symbol] --> MatchHead(symbol)

		var arg Expr
		if list.Length() == 0 {
			arg = ListFrom(symbolMatchAny)
		} else {
			arg = ListFrom(symbolMatchHead, list.Tail()[0])
		}
		c.emit(arg)
	case symbolBlankSequence:
		var arg Expr
		if list.Length() == 0 {
			arg = ListFrom(symbolMatchAny)
		} else {
			arg = ListFrom(symbolMatchHead, list.Tail()[0])
		}
		c.emit(ListFrom(symbolMatchPlus, arg))
	case symbolBlankNullSequence:
		var arg Expr
		if list.Length() == 0 {
			arg = ListFrom(symbolMatchAny)
		} else {
			arg = ListFrom(symbolMatchHead, list.Tail()[0])
		}
		c.emit(ListFrom(symbolMatchStar, arg))
	case symbolOptional:
		// TODO default value
		var arg Expr
		if list.Length() == 0 {
			arg = ListFrom(symbolMatchAny)
		} else {
			arg = ListFrom(symbolMatchHead, list.Tail()[0])
		}
		c.emit(ListFrom(symbolMatchQuest, arg))

	case symbolPattern:
		// Pattern("x", expression)
		args := list.Tail()
		name := args[0]
		slot := c.getSlot(name.(Symbol))
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
	case symbolPatternSequence:
		args := list.Tail()
		for _, arg := range args {
			c.emit(arg)
		}
	case symbolMatchHead:
		op := c.add(Inst{
			Op:  InstMatchHead,
			Val: list.Tail()[0],
		})
		c.addLink(op, c.pc)
	case symbolMatchAny:
		// this has a dangling pointer
		// it will be fixed at the end
		op := c.add(Inst{
			Op: InstMatchAny,
		})
		c.addLink(op, c.pc)
	case symbolMatchPlus:
		current := c.pc
		list, _ := e.(List)
		// only has a single argument
		c.emit(list.Tail()[0])

		op := c.add(Inst{
			Op: InstSplit,
		})
		c.addLink(op, current)
		c.addAlt(op, c.pc)
	case symbolMatchQuest:
		op := c.add(Inst{
			Op: InstSplit,
		})
		L1 := c.pc
		list, _ := e.(List)
		c.emit(list.Tail()[0])
		L2 := c.pc
		c.addLink(op, L1)
		c.addAlt(op, L2)
	case symbolMatchStar:
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
			Val:  list.Head(),
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
	Id   int32 // More for debugging
	Next int32 // Everyone
	Alt  int32 // Split, Capture
	Val  Expr  // use in MatchLiteral (anything), MatchHead, MatchList (Symbol)
	Data any   // Predictates, MatchList
}

func (i Inst) String() string {
	return fmt.Sprintf("id=%d [%s %d %d] ", i.Id, i.Op.String(), i.Next, i.Alt)
}

// a "Program" is a list of a Instructions
type Prog struct {
	Inst    []Inst
	groups  []Symbol
	onestep bool
}

func (p Prog) IsZero() bool {
	return len(p.Inst) == 0
}

func (p Prog) IsOneStep() bool {
	return p.onestep
}

func (p Prog) Groups() []Symbol {
	return p.groups
}

// Hack until we can do this in the instruction
func (p Prog) getSlot(name Symbol) int32 {
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
	fmt.Printf("One Step: %v\n", p.onestep)
	indent := strings.Repeat(" ", n)
	for _, inst := range p.Inst {
		fmt.Printf("%s%d: %s %s\n", indent, inst.Id, inst, inst.Val)
		if inst.Op == InstMatchList {
			p2 := inst.Data.(Prog)
			p2.dump(n + 1)
		}
	}
}
func (p Prog) Dump() {
	p.dump(0)
}
