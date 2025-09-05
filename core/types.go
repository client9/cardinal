package core

import (
	"github.com/client9/cardinal/core/symbol"
)

type Expr = symbol.Expr
type Symbol = symbol.SymbolExpr

func NewSymbol(s string) Symbol {
	return symbol.NewSymbol(s)
}
