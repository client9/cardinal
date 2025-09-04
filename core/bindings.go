package core

import (
	"fmt"

	"github.com/client9/sexpr/core/symbol"
)

type CaptureIndex struct {
	start int32
	end   int32
	exprs []Expr
}

type Captures struct {
	refs     int
	captures []CaptureIndex
}

func NewCaptures(n int) *Captures {
	if n == 0 {
		return &Captures{
			refs: 1,
		}
	}
	return &Captures{
		refs:     1,
		captures: make([]CaptureIndex, n),
	}
}

func (c *Captures) Inc() *Captures {
	c.refs += 1
	return c
}

func (c *Captures) Dec() {
	c.refs -= 1
	if c.refs == 0 {
		//fmt.Printf("refs to zero, has %d entries\n", c.Length())
		c.Clear()
	}
}

// DANGER
func (c *Captures) Clear() {
	c.refs = 1
	//if c.captures != nil {
	//for i:=0;i<len(c.captures);i++ {
	//	c.captures[i].end = -1
	//}
	clear(c.captures)
	//}
}

func (c *Captures) GetSlot(n int) (expr []Expr, start int32, end int32) {
	group := c.captures[n]
	return group.exprs, group.start, group.end
}

func (c *Captures) AddStart(slot int32, exprs []Expr, start int32) *Captures {
	if c.refs == 1 {
		group := c.captures[slot]
		if group.end == -1 {
			c.captures[slot].end = start
		} else {
			c.captures[slot] = CaptureIndex{start, -1, exprs}
		}
		return c
	}

	nc := NewCaptures(c.Length())
	copy(nc.captures, c.captures)
	c.Dec()

	group := nc.captures[slot]
	if group.end == -1 {
		nc.captures[slot].end = start
	} else {
		nc.captures[slot] = CaptureIndex{start, -1, exprs}
	}
	return nc
}

func (c *Captures) Length() int {
	return len(c.captures)
}

func (c *Captures) Dump(names []string) {
	fmt.Printf("\nCaptures %d\n", len(c.captures))
	for i, cc := range c.captures {
		fmt.Printf("  capture %d: %s %d %d for %s\n", i, names[i], cc.start, cc.end, cc.exprs)
	}
}

func (c *Captures) String() string {
	return c.AsRules(nil).String()
}

func (c *Captures) AsRules(names []Symbol) Expr {

	var name Symbol

	if c.Length() == 0 {
		//fmt.Printf("Capture had nothing\n")
		return ListFrom(symbol.List)
	}
	rules := make([]Expr, 0, len(c.captures))
	for i, cap := range c.captures {
		var target Expr
		// TODO do we still have the dangling names problem?
		if cap.end == -1 {
			continue
		}
		if cap.end-cap.start == 0 {
			// nil binding, or set as nil
			continue

		} else if cap.end-cap.start == 1 {
			// binding is single element
			target = cap.exprs[cap.start]
		} else {
			// mutliple elements go in a list`
			target = ListFrom(symbol.List, cap.exprs[cap.start:cap.end]...)
		}
		if names != nil {
			name = names[i]
		} else {
			name = NewSymbol(fmt.Sprintf("$%d", i+1))
		}
		rules = append(rules, ListFrom(symbol.Rule, name, target))
	}
	return ListFrom(symbol.List, rules...)
}
