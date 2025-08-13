package core

type Symbol string

func NewSymbol(s string) Symbol { return Symbol(s) }

func NewBool(value bool) Symbol {
	if value {
		return Symbol("True")
	} else {
		return Symbol("False")
	}
}

// NewSymbolNull creates the Null symbol
func NewSymbolNull() Symbol { return Symbol("Null") }

// Symbol type implementation
func (s Symbol) String() string {
	return string(s)
}

func (s Symbol) InputForm() string {
	return s.String()
}

func (s Symbol) Head() string {
	return "Symbol"
}

func (s Symbol) Length() int64 {
	return 0
}

func (s Symbol) Equal(rhs Expr) bool {
	if other, ok := rhs.(Symbol); ok {
		return s == other
	}
	return false
}

func (s Symbol) IsAtom() bool {
	return true
}
