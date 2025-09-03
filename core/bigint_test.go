package core

import (
	"fmt"
	"math"
	"testing"
)

func TestBigIntAdd(t *testing.T) {
	fmt.Println(math.MaxInt64)
	x := newMachineInt(math.MaxInt64)
	fmt.Println(x.String())
	z := addInteger(x, x)
	fmt.Println(z.String())
	z2 := addInteger(z, z)

	fmt.Println(z2.String())
}
func TestBigIntMult(t *testing.T) {
	fmt.Println(math.MaxInt64)
	x := newMachineInt(math.MaxInt64)
	fmt.Println(x.String())
	z := timesInteger(x, x)
	fmt.Println(z.String())
	z2 := timesInteger(z, z)

	fmt.Println(z2.String())
}
