package core

import (
	"fmt"
	// "log"

	"github.com/client9/cardinal/core/symbol"
)

type Thread struct {
	pc       *Inst     // "program counter"
	captures *Captures // COW, ref counted
}

func NewThread(pc *Inst, captures *Captures) Thread {
	return Thread{
		pc:       pc,
		captures: captures,
	}
}

func (t Thread) String() string {
	return fmt.Sprintf("Thread{ pc:%s, binding: %s }", t.pc.String(), t.captures.String())
}

type ThompsonVM struct {
	gen         int32
	genState    []int32
	currentList []Thread
	nextList    []Thread
}

func NewRegexp() *ThompsonVM {
	return &ThompsonVM{}
}

func (r *ThompsonVM) Match(prog Prog, expr Expr) (bool, *Captures) {
	sub := NewCaptures(len(prog.Groups()))
	return r.match(prog, expr, sub)
}

func (r *ThompsonVM) match(prog Prog, expr Expr, sub *Captures) (bool, *Captures) {
	if prog.IsOneStep() {
		return r.matchM4(prog, expr, sub)
	}
	return r.matchNfa(prog, expr, sub)
}

func (r *ThompsonVM) MatchList(prog Prog, exprs []Expr) (bool, *Captures) {
	sub := NewCaptures(len(prog.Groups()))
	return r.matchList(prog, exprs, sub)
}

func (r *ThompsonVM) matchList(prog Prog, exprs []Expr, sub *Captures) (bool, *Captures) {
	if prog.IsOneStep() {
		return r.matchSequenceM4(prog, exprs, sub)
	}
	return r.matchNfaSequence(prog, exprs, sub)
}

// the rest is internal functions

func (r *ThompsonVM) reset(n int) {

	r.gen = 0

	if cap(r.genState) < n {
		r.genState = make([]int32, n)
	} else {
		// keep length, but reset all values to 0
		clear(r.genState)
	}
	if cap(r.currentList) < n {
		r.currentList = make([]Thread, 0, n)
	} else {
		r.currentList = r.currentList[:0]
	}
	if cap(r.nextList) < n {
		r.nextList = make([]Thread, 0, n)
	} else {
		r.nextList = r.nextList[:0]
	}
}

func (r *ThompsonVM) AddThread(tlist *[]Thread, prog Prog, pc int32, exprs []Expr, j int, sub *Captures) {

	if r.genState[pc] == r.gen {
		sub.Dec()
		return
	}
	r.genState[pc] = r.gen

	i := &prog.Inst[pc]
	switch i.Op {
	case InstSplit:
		r.AddThread(tlist, prog, i.Next, exprs, j, sub.Inc())
		r.AddThread(tlist, prog, i.Alt, exprs, j, sub)
	case InstCaptureStart:
		r.AddThread(tlist, prog, i.Next, exprs, j, sub.AddStart(i.Alt, exprs, int32(j)))
	case InstJump:
		r.AddThread(tlist, prog, i.Next, exprs, j, sub)
	default:
		*tlist = append(*tlist, NewThread(i, sub))
	}
}

func (r *ThompsonVM) matchNfa(prog Prog, expr Expr, sub *Captures) (bool, *Captures) {
	r.reset(prog.Length())
	r.gen += 1
	pc := prog.First()

	r.currentList = append(r.currentList, NewThread(&prog.Inst[pc], sub))
	//r.AddThread(&r.currentList, prog, pc, nil, 0, sub)

	captureWholeInput := false
	for {
		if len(r.currentList) == 0 {
			return false, nil
		}
		r.gen += 1
		for _, c := range r.currentList {
			i := c.pc
			add := false
			switch i.Op {
			case InstMatchEnd:
				if captureWholeInput {
					c.captures.captures[0].exprs = []Expr{expr}
					c.captures.captures[0].start = 0
					c.captures.captures[0].end = 1
				}
				return true, c.captures
			case InstMatchAny:
				add = true
			case InstMatchHead:
				head := expr.Head()
				add = (i.Val == head) || (i.Val == symbol.Number && (head == symbol.Integer || head == symbol.Real))
			case InstMatchLiteral:
				add = expr.Equal(i.Val)
			case InstCaptureStart:
				add = true
				captureWholeInput = true
			case InstMatchList:
				if list, ok := expr.(List); ok && (i.Val == symbol.Null || list.Head() == i.Val) {
					// TODO
					// the new pc.prog has ID that conflict with outer ones
					// technically they could be renumbered and just call recursively
					// would need to adjust max-number of threads however.
					proglist := i.Data.(Prog)
					r2 := NewRegexp()
					r2.reset(proglist.Length())

					if ok, binding := r2.matchNfaSequence(proglist, list.Tail(), c.captures.Inc()); ok {
						c.captures.Dec()
						// ?? binding.Inc()
						add = true
						c.captures = binding
					}
				}
			} // end switch
			if add {
				r.AddThread(&r.nextList, prog, i.Next, nil, 0, c.captures)
			}
			r.currentList, r.nextList = r.nextList, r.currentList
			r.nextList = r.nextList[:0]
		}
	}
	/*
		for _, c := range r.currentList {
			switch c.pc.Op {
			case InstMatchEnd:
				if captureWholeInput {
					fmt.Printf("GOT WHOLE INPUT\n")
					c.captures.refs = 1
					c.captures.AddStart(0, []Expr{expr}, 0)
					c.captures.AddStart(0, []Expr{expr}, 1)
				}
				return true, c.captures
			default:
				fmt.Printf("GOT EXTRA OF %v\n", c.pc)
			}
		}
	*/
	return false, nil

}

func (r *ThompsonVM) matchNfaSequence(prog Prog, exprs []Expr, sub *Captures) (bool, *Captures) {
	r.reset(prog.Length())
	r.gen += 1
	r.AddThread(&r.currentList, prog, prog.First(), exprs, 0, sub)

	for j, elem := range exprs {
		if len(r.currentList) == 0 {
			return false, nil
			//break
		}
		r.gen += 1
		for _, c := range r.currentList {
			add := false
			i := c.pc
			switch i.Op {
			case InstMatchAny:
				add = true
			case InstMatchHead:
				head := elem.Head()
				add = (head == i.Val) || (i.Val == symbol.Number && (head == symbol.Integer || head == symbol.Real))
			case InstMatchLiteral:
				add = elem.Equal(i.Val)
			case InstMatchList:
				if lst, ok := elem.(List); ok && (i.Val == symbol.Null || lst.Head() == i.Val) {
					// loop detection issue
					// program ids conflict
					listprog := i.Data.(Prog)
					r2 := NewRegexp()
					r2.reset(listprog.Length())
					if ok, binding := r2.matchNfaSequence(listprog, lst.Tail(), c.captures.Inc()); ok {
						c.captures.Dec()
						//??binding.Inc()
						add = true
						c.captures = binding
					}
				}
			case InstMatchEnd:
				// this is different than RCS's regexp where this could be
				// reached after everything was matched (trailing \0)
				//
				// Here, we've reached the end before we have consumed all characters
				// nothing to do and let this thread die
				//fmt.Printf("******* In OpMatchEnd\n")
			}

			if add {
				r.AddThread(&r.nextList, prog, i.Next, exprs, j+1, c.captures)
			} else {
				c.captures.Dec()
			}
		}
		r.currentList, r.nextList = r.nextList, r.currentList
		r.nextList = r.nextList[:0]
	}

	for _, c := range r.currentList {
		switch c.pc.Op {
		case InstMatchEnd:
			//fmt.Printf("%d IN MATCH END\n", i)
			//c.captures.Dump()
			//fmt.Printf("End dump")
			return true, c.captures

			/*
				// unclear if needed
				for j := i+1; j < len(currentList); j++ {
					currentList[j].captures.Dec()
				}
				break
				//currentList = currentList[:i]
				//return true, c.captures
			*/
		}
	}
	/*
				r.currentList, r.nextList = r.nextList, r.currentList
				r.nextList = r.nextList[:0]
		}
	*/
	return false, nil
}

func (r *ThompsonVM) MatchM2(prog Prog, e Expr) (bool, *Captures) {
	sub := NewCaptures(len(prog.Groups()))
	return r.matchSequenceM2(prog, []Expr{e}, sub)
}

func (r *ThompsonVM) matchM2(prog Prog, e Expr, sub *Captures) (bool, *Captures) {
	return r.matchSequenceM2(prog, []Expr{e}, sub)
}

func (r *ThompsonVM) MatchSequenceM2(prog Prog, expr []Expr) (bool, *Captures) {
	sub := NewCaptures(len(prog.Groups()))
	return r.matchSequenceM2(prog, expr, sub)
}

func (r *ThompsonVM) matchSequenceM2(prog Prog, expr []Expr, sub *Captures) (bool, *Captures) {
	pc := prog.First()
	for j := 0; j < len(expr); {

		// Assume true
		last := true

		e := expr[j]
		i := &prog.Inst[pc]
		switch i.Op {

		// special case
		case InstMatchEnd:
			return false, nil

		// special case
		case InstSplit:
			if !last {
				pc = i.Alt
				continue
			}
			last = false
		case InstMatchAny:
			last = true
		case InstMatchHead:
			last = e.Head() == i.Val
		case InstMatchLiteral:
			last = e.Equal(i.Val)
		case InstMatchList:
			last = false
			if list, ok := e.(List); ok && (i.Val == symbol.Null || list.Head() == i.Val) {
				if ok, nsub := r.matchSequenceM2(i.Data.(Prog), list.Tail(), sub); ok {
					sub = nsub
					last = true
				}
			}
		case InstCaptureStart:
			last = false
			sub = sub.AddStart(i.Alt, expr, int32(j))
		} // end switch
		if last {
			j += 1
		}
		pc = i.Next
	}

	// why this works...
	// if the last MatchAny (or other simgle matcher) is done on the last
	// element in the sequence, and it's part of MatchStar
	// the next instruction is Split, and takes the Alt path
	//
	// if the next pattern is a MatchStar, then it's a OpSplit
	// again and takes the alt path until it hits MatchEnd
	//
	//for j := 0; j < 1; {
	sub.refs = 1
	for {
		i := &prog.Inst[pc]
		switch i.Op {
		case InstMatchEnd:
			return true, sub
		case InstSplit:
			pc = i.Alt
		// if we don't do this, then the end position is not set
		//   maybe we post process this.
		case InstCaptureStart:
			sub = sub.AddStart(i.Alt, expr, int32(len(expr)))
			pc = i.Next

		// TBD: jump is probably a fail case
		case InstJump:
			//return false, nil
			pc = i.Next
		default: // InstMatchAny, InstMatchHead, InstMatchLiteral:
			return false, nil
			//j += 1
			//pc = i.Next
		}
	}

	return false, nil
}
func (r *ThompsonVM) MatchM3(prog Prog, e Expr) (bool, *Captures) {
	sub := NewCaptures(len(prog.Groups()))
	return r.matchM3(prog, e, sub)
}
func (r *ThompsonVM) matchM3(prog Prog, expr Expr, sub *Captures) (bool, *Captures) {
	n := len(prog.Groups())
	if sub.Length() < n {
		sub = NewCaptures(n)
	} else {
		sub.Clear()
	}

	pc := prog.First()
	capture0 := false

	e := expr
	consume := false
	for {
		i := &prog.Inst[pc]
		switch i.Op {
		case InstMatchList:
			consume = false
			if list, ok := e.(List); ok && (i.Val == symbol.Null || list.Head() == i.Val) {
				if ok, nsub := r.matchSequenceM3(i.Data.(Prog), list.Tail(), sub); ok {
					sub = nsub
					consume = true
				}
			}
		case InstMatchEnd:
			if e == nil && capture0 {
				sub.AddStart(0, []Expr{expr}, 0)
				sub.AddStart(0, []Expr{expr}, 1)
			}
			return e == nil, sub
		case InstFail:
			return false, nil
		case InstMatchAny:
			consume = e != nil
		case InstMatchHead:
			consume = e != nil && e.Head() == i.Val
		case InstMatchLiteral:
			consume = e != nil && e.Equal(i.Val)
		case InstCaptureStart:
			consume = false
			capture0 = true
			//sub = sub.AddStart(i.Alt, []Expr{e}, 0)
			pc = i.Next
			continue // continue for loop
		} // end switch

		if consume {
			e = nil
			pc = i.Next
		} else {
			pc = i.Alt
		}
	}
}
func (r *ThompsonVM) MatchSequenceM3(prog Prog, expr []Expr) (bool, *Captures) {
	sub := NewCaptures(len(prog.Groups()))
	return r.matchSequenceM3(prog, expr, sub)
}
func (r *ThompsonVM) matchListM3(prog Prog, e Expr, sub *Captures) (bool, *Captures) {
	if elist, ok := e.(List); ok {
		return r.matchSequenceM3(prog, elist.Tail(), sub)
	}
	return false, nil
}
func (r *ThompsonVM) matchSequenceM3(prog Prog, expr []Expr, sub *Captures) (bool, *Captures) {
	n := len(prog.Groups())
	if sub.Length() < n {
		sub = NewCaptures(n)
	} else {
		sub.Clear()
	}

	pc := prog.First()
	for j := 0; j < len(expr); {
		var consume bool
		e := expr[j]

		i := &prog.Inst[pc]

		switch i.Op {

		// special case
		case InstMatchEnd, InstFail:
			return false, nil
		case InstMatchAny:
			consume = true
		case InstMatchHead:
			consume = e.Head() == i.Val
		case InstMatchLiteral:
			consume = e.Equal(i.Val)
		case InstMatchList:
			sub.refs = 1
			if list, ok := e.(List); ok && (i.Val == symbol.Null || list.Head() == i.Val) {
				if ok, nsub := r.matchSequenceM3(i.Data.(Prog), list.Tail(), sub); ok {
					sub = nsub
					consume = true
				}
			}
		case InstCaptureStart:
			consume = false
			sub = sub.AddStart(i.Alt, expr, int32(j))
			pc = i.Next
			continue

		} // end switch

		if consume {
			j += 1
			pc = i.Next
		} else {
			pc = i.Alt
		}
	}

	for {
		i := &prog.Inst[pc]
		switch i.Op {
		case InstMatchEnd:
			return true, sub
		case InstFail:
			return false, nil
		case InstCaptureStart:
			sub = sub.AddStart(i.Alt, expr, int32(len(expr)))
			pc = i.Next
		default:
			pc = i.Alt
		}
	}
	return false, nil
}

func (r *ThompsonVM) MatchM4(prog Prog, e Expr) (bool, *Captures) {
	sub := NewCaptures(len(prog.Groups()))
	return r.matchM4(prog, e, sub)
}
func (r *ThompsonVM) MatchSequenceM4(prog Prog, args []Expr) (bool, *Captures) {
	sub := NewCaptures(len(prog.Groups()))
	return r.matchSequenceM4(prog, args, sub)
}

func (r *ThompsonVM) matchM4(prog Prog, expr Expr, sub *Captures) (bool, *Captures) {

	n := len(prog.Groups())
	if sub.Length() < n {
		sub = NewCaptures(n)
	} else {
		sub.Clear()
	}

	captureWholeInput := false

	pc := prog.First()
	e := expr
	consume := false
	for {
		i := &prog.Inst[pc]
		switch i.Op {
		case InstMatchList:
			consume = false

			fmt.Println("list = ", e, "iVal = ", i.Val, ", iVal==Null", i.Val == symbol.Null)
			i.Data.(Prog).Dump()
			if list, ok := e.(List); ok && (i.Val == symbol.Null || list.Head() == i.Val) {
				fmt.Println("    Passed condition, tail: ", list.Tail())
				if ok, nsub := r.matchSequenceM4(i.Data.(Prog), list.Tail(), sub); ok {
					fmt.Println("   And passed again")
					sub = nsub
					consume = true
				}
			}
		case InstMatchEnd:
			if e == nil && captureWholeInput {
				sub.AddStart(0, []Expr{expr}, 0)
				sub.AddStart(0, []Expr{expr}, 1)
			}
			return e == nil, sub
		case InstFail:
			return false, nil
		case InstMatchAny:
			consume = e != nil
		case InstMatchHead:
			consume = e != nil && matchHead(e.Head(), i.Val.(Symbol))
		case InstMatchLiteral:
			consume = e != nil && e.Equal(i.Val)
		case InstCaptureStart:
			consume = false
			captureWholeInput = true
			//sub = sub.AddStart(i.Alt, []Expr{e}, 0)
			pc = i.Next
			continue // continue for loop
		} // end switch

		if consume {
			e = nil
			pc = i.Next
		} else {
			pc = i.Alt
		}
	}
}

func (r *ThompsonVM) matchListM4(prog Prog, e Expr, sub *Captures) (bool, *Captures) {
	if elist, ok := e.(List); ok {
		return r.matchSequenceM4(prog, elist.Tail(), sub)
	}
	return false, nil
}

func (r *ThompsonVM) matchSequenceM4(prog Prog, args []Expr, sub *Captures) (bool, *Captures) {
	pc := prog.First()
	var j int
	var consume bool
	for {
		var e Expr

		if j < len(args) {
			e = args[j]
		}
	Again:
		i := &prog.Inst[pc]
		switch i.Op {
		case InstMatchList:
			consume = false
			if list, ok := e.(List); ok && (i.Val == symbol.Null || list.Head() == i.Val) {
				if ok, nsub := r.matchSequenceM4(i.Data.(Prog), list.Tail(), sub); ok {
					sub = nsub
					consume = true
				}
			}
		case InstMatchEnd:
			return e == nil, sub
		case InstFail:
			return false, nil
		case InstMatchAny:
			consume = e != nil
		case InstMatchHead:
			consume = e != nil && matchHead(e.Head(), i.Val.(Symbol))
		case InstMatchLiteral:
			consume = e != nil && e.Equal(i.Val)
		case InstCaptureStart:
			consume = false
			sub.AddStart(i.Alt, args, int32(j))
			pc = i.Next
			goto Again // reprocess input
		}

		if consume {
			j += 1
			pc = i.Next
		} else {
			pc = i.Alt
		}
	}
}

func matchHead(head Expr, sym Symbol) bool {
	return sym == head || (sym == symbol.Number && (head == symbol.Integer || head == symbol.Real))
}
