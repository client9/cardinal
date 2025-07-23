package sexpr

import (
	"fmt"
	"strconv"
	"strings"
)

// Uint64Value implements Expr for 64-bit unsigned integers
type Uint64Value struct {
	value uint64
}

func (u Uint64Value) String() string {
	return fmt.Sprintf("#%X", u.value)
}

func (u Uint64Value) Length() int64 {
	return 0
}

func (u Uint64Value) InputForm() string {
	return u.String() // For Uint64Value, InputForm is the same as String()
}

func (u Uint64Value) Type() string {
	return "uint64" // Internal Go type
}

func (u Uint64Value) Equal(rhs Expr) bool {
	if rhsUint64, ok := rhs.(Uint64Value); ok {
		return u.value == rhsUint64.value
	}
	return false
}

// NewUint64Value creates a new Uint64Value
func NewUint64Value(value uint64) Uint64Value {
	return Uint64Value{value: value}
}

// GetValue returns the underlying uint64 value
func (u Uint64Value) GetValue() uint64 {
	return u.value
}

// RegisterUint64Type registers the Uint64 type constructor and operations
func RegisterUint64Type(registry *FunctionRegistry) error {
	// Constructor from hex string: Uint64("#FFFFFFFF") -> ObjectExpr{TypeName: "Uint64", Value: Uint64Value}
	err := registry.RegisterPatternBuiltin("Uint64(x_String)", func(args []Expr, ctx *Context) Expr {
		if len(args) != 1 {
			return NewErrorExpr("ArgumentError", "Uint64 expects 1 argument", args)
		}

		strAtom, ok := args[0].(Atom)
		if !ok || strAtom.AtomType != StringAtom {
			return NewErrorExpr("ArgumentError", "Uint64 constructor requires string argument", args)
		}

		hexStr := strAtom.Value.(string)
		if !strings.HasPrefix(hexStr, "#") {
			return NewErrorExpr("ArgumentError", "Uint64 hex string must start with #", args)
		}

		value, err := strconv.ParseUint(hexStr[1:], 16, 64)
		if err != nil {
			return NewErrorExpr("ParseError", fmt.Sprintf("Invalid hex string: %v", err), args)
		}

		uint64Expr := NewUint64Value(value)
		return NewObjectExpr("Uint64", uint64Expr)
	})
	if err != nil {
		return fmt.Errorf("failed to register Uint64 string constructor: %v", err)
	}

	// Constructor from integer: Uint64(42) -> ObjectExpr{TypeName: "Uint64", Value: Uint64Value}
	err = registry.RegisterPatternBuiltin("Uint64(x_Integer)", func(args []Expr, ctx *Context) Expr {
		if len(args) != 1 {
			return NewErrorExpr("ArgumentError", "Uint64 expects 1 argument", args)
		}

		intAtom, ok := args[0].(Atom)
		if !ok || intAtom.AtomType != IntAtom {
			return NewErrorExpr("ArgumentError", "Uint64 constructor requires integer argument", args)
		}

		intVal := intAtom.Value.(int)
		if intVal < 0 {
			return NewErrorExpr("ValueError", "Uint64 cannot be negative", args)
		}

		uint64Expr := NewUint64Value(uint64(intVal))
		return NewObjectExpr("Uint64", uint64Expr)
	})
	if err != nil {
		return fmt.Errorf("failed to register Uint64 integer constructor: %v", err)
	}

	// Handle error cases with catch-all pattern
	err = registry.RegisterPatternBuiltin("Uint64(x___)", func(args []Expr, ctx *Context) Expr {
		if len(args) == 0 {
			return NewErrorExpr("ArgumentError", "Uint64 expects 1 argument, got 0", args)
		}
		if len(args) > 1 {
			return NewErrorExpr("ArgumentError", fmt.Sprintf("Uint64 expects 1 argument, got %d", len(args)), args)
		}

		// Single argument but not a string
		if atom, ok := args[0].(Atom); !ok || atom.AtomType != StringAtom {
			return NewErrorExpr("ArgumentError", "Uint64 constructor requires string argument", args)
		}

		// This should never be reached due to pattern specificity, but just in case
		return NewErrorExpr("ArgumentError", "Uint64 constructor error", args)
	})
	if err != nil {
		return fmt.Errorf("failed to register Uint64 error handler: %v", err)
	}

	// Type predicate: Uint64Q(x) -> True/False
	err = registry.RegisterPatternBuiltin("Uint64Q(x_)", func(args []Expr, ctx *Context) Expr {
		if len(args) != 1 {
			return NewErrorExpr("ArgumentError", "Uint64Q expects 1 argument", args)
		}

		if objExpr, ok := args[0].(ObjectExpr); ok && objExpr.TypeName == "Uint64" {
			return NewSymbolAtom("True")
		}
		return NewSymbolAtom("False")
	})
	if err != nil {
		return fmt.Errorf("failed to register Uint64Q predicate: %v", err)
	}

	return nil
}

// RegisterUint64Operations registers arithmetic operations for Uint64
func RegisterUint64Operations(registry *FunctionRegistry) error {
	// Pure Uint64 addition: Plus(x_Uint64, y_Uint64)
	err := registry.RegisterPatternBuiltin("Plus(x_Uint64, y_Uint64)", func(args []Expr, ctx *Context) Expr {
		if len(args) != 2 {
			return NewErrorExpr("ArgumentError", "Plus expects 2 arguments", args)
		}

		x := args[0].(ObjectExpr).Value.(Uint64Value).GetValue()
		y := args[1].(ObjectExpr).Value.(Uint64Value).GetValue()

		// Check for overflow
		if x > ^uint64(0)-y {
			return NewErrorExpr("OverflowError", "Uint64 addition overflow", args)
		}

		result := NewUint64Value(x + y)
		return NewObjectExpr("Uint64", result)
	})
	if err != nil {
		return fmt.Errorf("failed to register Plus(Uint64, Uint64): %v", err)
	}

	// Mixed Integer + Uint64: Plus(x_Integer, y_Uint64)
	err = registry.RegisterPatternBuiltin("Plus(x_Integer, y_Uint64)", func(args []Expr, ctx *Context) Expr {
		if len(args) != 2 {
			return NewErrorExpr("ArgumentError", "Plus expects 2 arguments", args)
		}

		x := args[0].(Atom).Value.(int)
		y := args[1].(ObjectExpr).Value.(Uint64Value).GetValue()

		if x < 0 {
			return NewErrorExpr("ValueError", "Cannot add negative integer to Uint64", args)
		}

		ux := uint64(x)
		if ux > ^uint64(0)-y {
			return NewErrorExpr("OverflowError", "Uint64 addition overflow", args)
		}

		result := NewUint64Value(ux + y)
		return NewObjectExpr("Uint64", result)
	})
	if err != nil {
		return fmt.Errorf("failed to register Plus(Integer, Uint64): %v", err)
	}

	// Mixed Uint64 + Integer: Plus(x_Uint64, y_Integer)
	err = registry.RegisterPatternBuiltin("Plus(x_Uint64, y_Integer)", func(args []Expr, ctx *Context) Expr {
		if len(args) != 2 {
			return NewErrorExpr("ArgumentError", "Plus expects 2 arguments", args)
		}

		x := args[0].(ObjectExpr).Value.(Uint64Value).GetValue()
		y := args[1].(Atom).Value.(int)

		if y < 0 {
			return NewErrorExpr("ValueError", "Cannot add negative integer to Uint64", args)
		}

		uy := uint64(y)
		if x > ^uint64(0)-uy {
			return NewErrorExpr("OverflowError", "Uint64 addition overflow", args)
		}

		result := NewUint64Value(x + uy)
		return NewObjectExpr("Uint64", result)
	})
	if err != nil {
		return fmt.Errorf("failed to register Plus(Uint64, Integer): %v", err)
	}

	// Pure Uint64 multiplication: Times(x_Uint64, y_Uint64)
	err = registry.RegisterPatternBuiltin("Times(x_Uint64, y_Uint64)", func(args []Expr, ctx *Context) Expr {
		if len(args) != 2 {
			return NewErrorExpr("ArgumentError", "Times expects 2 arguments", args)
		}

		x := args[0].(ObjectExpr).Value.(Uint64Value).GetValue()
		y := args[1].(ObjectExpr).Value.(Uint64Value).GetValue()

		// Check for overflow
		if x != 0 && y > ^uint64(0)/x {
			return NewErrorExpr("OverflowError", "Uint64 multiplication overflow", args)
		}

		result := NewUint64Value(x * y)
		return NewObjectExpr("Uint64", result)
	})
	if err != nil {
		return fmt.Errorf("failed to register Times(Uint64, Uint64): %v", err)
	}

	return nil
}

// RegisterUint64 registers the complete Uint64 type system
func RegisterUint64(registry *FunctionRegistry) error {
	if err := RegisterUint64Type(registry); err != nil {
		return err
	}
	if err := RegisterUint64Operations(registry); err != nil {
		return err
	}
	return nil
}
