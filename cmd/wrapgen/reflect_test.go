package main

import (
	"fmt"
	"github.com/client9/sexpr/stdlib"
	"testing"
)

func TestList(t *testing.T) {
	stats, err := analyzeFunctionSignature(stdlib.ListAppend)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	fmt.Printf("%v\n", stats)
}
func TestR1(t *testing.T) {

	stats, err := analyzeFunctionSignature(stdlib.LessNumber)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", stats)
}
