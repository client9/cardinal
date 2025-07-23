package main

import (
	"fmt"
	"github.com/client9/sexpr/stdlib"
	"testing"
)

func TestR1(t *testing.T) {

	stats, err := analyzeFunctionSignature(stdlib.LessNumber)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", stats)
}
