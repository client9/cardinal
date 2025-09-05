package core

import (
	"strconv"

	"github.com/client9/cardinal/core/symbol"
)

type Rune rune

func NewRune(s rune) Rune { return Rune(s) }

// String type implementation
func (s Rune) String() string {
	return strconv.QuoteRune(rune(s))
}

func (s Rune) InputForm() string {
	return strconv.QuoteRune(rune(s))
}

func (s Rune) Head() Expr {
	return symbol.Rune
}

func (s Rune) Length() int64 {
	return 0
}

func (s Rune) Equal(rhs Expr) bool {
	return s == rhs
}

func (s Rune) IsAtom() bool {
	return true
}
