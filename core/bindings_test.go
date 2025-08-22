package core

import (
	"fmt"
	"testing"
)

// haven't been referenced by anyone else, so it's save to modify
func TestBindingsNoCopy(t *testing.T) {
	c := NewCaptures(1)

	c = c.AddStart(0, nil, 10)
	d := c.AddStart(0, nil, 20)

	if d.Length() != c.Length() {
		t.Errorf("Expected same object")
	}

	_, a, b := c.GetSlot(0)
	if a != 10 || b != 20 {
		t.Errorf("Failed got %d %d", a, b)
	}
}

func TestBindingsWithCopy(t *testing.T) {
	fmt.Printf("-------------\n")
	c := NewCaptures(1)
	c.Inc()
	fmt.Printf("Refs is %d\n", c.refs)

	d := c.AddStart(0, nil, 10)
	d = d.AddStart(0, nil, 20)

	_, a, b := c.GetSlot(0)
	if a != 0 || b != 0 {
		t.Errorf("Original object should be all zeros")
	}
	_, a, b = d.GetSlot(0)
	if a != 10 || b != 20 {
		t.Errorf("Modified copy failed, got %d %d", a, b)
	}
}
