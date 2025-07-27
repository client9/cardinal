package main

import (
	// Import stdlib functions for reflection
	"github.com/client9/sexpr/stdlib"
)

// SymbolSpec defines a complete symbol with its attributes and functions
type SymbolSpec struct {
	Name       string         // "Plus" - the symbol name
	Attributes []string       // ["Flat", "Orderless"] - symbol attributes
	Functions  map[string]any // "(x__Integer)" -> stdlib.PlusIntegers
	Constants  map[string]any // For symbols like Pi, E that have constant values
}

// Symbol specifications organized by symbol name
var symbolSpecs = map[string]SymbolSpec{
	// Arithmetic Operations
	"Plus": {
		Name:       "Plus",
		Attributes: []string{"Flat", "Listable", "NumericFunction", "OneIdentity", "Orderless", "Protected"},
		Functions: map[string]any{
			"()":           stdlib.PlusIdentity,
			"(x__Integer)": stdlib.PlusIntegers,
			"(x__Real)":    stdlib.PlusReals,
			"(x__Number)":  stdlib.PlusNumbers,
		},
	},
	"Times": {
		Name:       "Times",
		Attributes: []string{"Flat", "Orderless", "OneIdentity"},
		Functions: map[string]any{
			"()":           stdlib.TimesIdentity,
			"(x__Integer)": stdlib.TimesIntegers,
			"(x__Real)":    stdlib.TimesReals,
			"(x__Number)":  stdlib.TimesNumbers,
		},
	},
	"Power": {
		Name:       "Power",
		Attributes: []string{"OneIdentity"},
		Functions: map[string]any{
			"(base_Real, exp_Integer)": stdlib.PowerReal,
			"(x_Number, y_Number)":     stdlib.PowerExprs,
		},
	},
	"Subtract": {
		Name:       "Subtract",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_Integer, y_Integer)": stdlib.SubtractIntegers,
			"(x_Number, y_Number)":   stdlib.SubtractExprs,
		},
	},
	"Minus": {
		Name:       "Minus",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_Integer)": stdlib.MinusInteger,
			"(x_Real)":    stdlib.MinusReal,
			"(x_)":        stdlib.MinusExpr,
		},
	},
	"Divide": {
		Name:       "Divide",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_Integer, y_Integer)": stdlib.DivideIntegers,
			"(x_Number, y_Number)":   stdlib.DivideExprs,
		},
	},

	// Comparison Operations
	"Equal": {
		Name:       "Equal",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_Integer, y_Integer)": stdlib.EqualInts,
			"(x_Real, y_Real)":       stdlib.EqualFloats,
			"(x_Number, y_Number)":   stdlib.EqualNumbers,
			"(x_String, y_String)":   stdlib.EqualStrings,
			"(x_, y_)":               stdlib.EqualExprs,
		},
	},
	"Unequal": {
		Name:       "Unequal",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_Integer, y_Integer)": stdlib.UnequalInts,
			"(x_Real, y_Real)":       stdlib.UnequalFloats,
			"(x_Number, y_Number)":   stdlib.UnequalNumbers,
			"(x_String, y_String)":   stdlib.UnequalStrings,
			"(x_, y_)":               stdlib.UnequalExprs,
		},
	},
	"Less": {
		Name:       "Less",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_Number, y_Number)": stdlib.LessNumber,
		},
	},
	"Greater": {
		Name:       "Greater",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_Number, y_Number))": stdlib.GreaterNumber,
		},
	},
	"LessEqual": {
		Name:       "LessEqual",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_Number, y_Number)": stdlib.LessEqualNumber,
		},
	},
	"GreaterEqual": {
		Name:       "GreaterEqual",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_Number, y_Number)": stdlib.GreaterEqualNumber,
		},
	},
	"SameQ": {
		Name:       "SameQ",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_, y_)": stdlib.SameQExprs,
		},
	},
	"UnsameQ": {
		Name:       "UnsameQ",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_, y_)": stdlib.UnsameQExprs,
		},
	},

	// Type Predicates
	"IntegerQ": {
		Name:       "IntegerQ",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_)": stdlib.IntegerQExpr,
		},
	},
	"FloatQ": {
		Name:       "FloatQ",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_)": stdlib.FloatQExpr,
		},
	},
	"NumberQ": {
		Name:       "NumberQ",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_)": stdlib.NumberQExpr,
		},
	},
	"StringQ": {
		Name:       "StringQ",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_)": stdlib.StringQExpr,
		},
	},
	"BooleanQ": {
		Name:       "BooleanQ",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_)": stdlib.BooleanQExpr,
		},
	},
	"SymbolQ": {
		Name:       "SymbolQ",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_)": stdlib.SymbolQExpr,
		},
	},
	"TrueQ": {
		Name:       "TrueQ",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_)": stdlib.TrueQExpr,
		},
	},
	"ListQ": {
		Name:       "ListQ",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_)": stdlib.ListQExpr,
		},
	},
	"AtomQ": {
		Name:       "AtomQ",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_)": stdlib.AtomQExpr,
		},
	},
	"Head": {
		Name:       "Head",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_)": stdlib.HeadExpr,
		},
	},

	// Output Format Functions
	"FullForm": {
		Name:       "FullForm",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_)": stdlib.FullFormExpr,
		},
	},
	"InputForm": {
		Name:       "InputForm",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_)": stdlib.InputFormExpr,
		},
	},

	// List Operations
	"Length": {
		Name:       "Length",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_)": stdlib.LengthExpr,
		},
	},
	"First": {
		Name:       "First",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_List)": stdlib.FirstExpr,
			"(x_)":     stdlib.First,
		},
	},
	"Last": {
		Name:       "Last",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_List)": stdlib.LastExpr,
			"(x_)":     stdlib.Last,
		},
	},
	"Rest": {
		Name:       "Rest",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_List)": stdlib.RestExpr,
			"(x_)":     stdlib.Rest,
		},
	},
	"Most": {
		Name:       "Most",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_List)": stdlib.MostExpr,
			"(x_)":     stdlib.Most,
		},
	},
	"Append": {
		Name:       "Append",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_List, y_)":         stdlib.ListAppend,
			"(x_String, y_String)": stdlib.StringAppend,
		},
	},

	// Sequence Operations
	"Take": {
		Name:       "Take",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_, n_Integer)":                  stdlib.Take,
			"(x_, List(n_Integer, m_Integer))": stdlib.TakeRange,
		},
	},
	"Drop": {
		Name:       "Drop",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_, n_Integer)":                  stdlib.Drop,
			"(x_, List(n_Integer, m_Integer))": stdlib.DropRange,
		},
	},
	"Part": {
		Name:       "Part",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_, n_Integer)":     stdlib.Part,
			"(x_Association, y_)": stdlib.PartAssociation,
		},
	},
	"Reverse": {
		Name:       "Reverse",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_String)": stdlib.StringReverse,
		},
	},
	"RotateLeft": {
		Name:       "RotateLeft",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_, n_Integer)": stdlib.RotateLeft,
		},
	},
	"RotateRight": {
		Name:       "RotateRight",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_, n_Integer)": stdlib.RotateRight,
		},
	},

	// Logical Operations
	"Not": {
		Name:       "Not",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_)": stdlib.NotExpr,
		},
	},
	"MatchQ": {
		Name:       "MatchQ",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_, y_)": stdlib.MatchQExprs,
		},
	},
	"Replace": {
		Name:       "Replace",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_, r_Rule)": stdlib.ReplaceExpr,
		},
	},

	// String Operations
	"StringLength": {
		Name:       "StringLength",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_String)": stdlib.StringLengthRunes,
		},
	},
	"ByteArray": {
		Name:       "ByteArray",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_String)": stdlib.ByteArrayFromString,
		},
	},

	// Association Operations
	"AssociationQ": {
		Name:       "AssociationQ",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_)": stdlib.AssociationQExpr,
		},
	},
	"Keys": {
		Name:       "Keys",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_Association)": stdlib.KeysExpr,
		},
	},
	"Values": {
		Name:       "Values",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_Association)": stdlib.ValuesExpr,
		},
	},
	"Association": {
		Name:       "Association",
		Attributes: []string{},
		Functions: map[string]any{
			"(x___Rule)": stdlib.AssociationRules,
		},
	},

	// Output Operations
	"Print": {
		Name:       "Print",
		Attributes: []string{},
		Functions: map[string]any{
			"(x_)": stdlib.Print,
		},
	},

	// Constants (symbols with values but no functions)
	"Pi": {
		Name:       "Pi",
		Attributes: []string{"Constant", "Protected"},
		Constants:  map[string]any{"Pi": 3.141592653589793},
	},
	"E": {
		Name:       "E",
		Attributes: []string{"Constant", "Protected"},
		Constants:  map[string]any{"E": 2.718281828459045},
	},
	"True": {
		Name:       "True",
		Attributes: []string{"Constant", "Protected"},
		Constants:  map[string]any{"True": "True"},
	},
	"False": {
		Name:       "False",
		Attributes: []string{"Constant", "Protected"},
		Constants:  map[string]any{"False": "False"},
	},
}
