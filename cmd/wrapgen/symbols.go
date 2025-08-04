package main

import (
	"github.com/client9/sexpr/builtins"
	"github.com/client9/sexpr/stdlib"
)

type rule struct {
	Pattern  string
	Function any
}

// SymbolSpec defines a complete symbol with its attributes and functions
type SymbolSpec struct {
	Name       string   // "Plus" - the symbol name
	Attributes []string // ["Flat", "Orderless"] - symbol attributes
	Functions  []rule   // "(x__Integer)" , stdlib.PlusIntegers
	Constants  []rule   // For symbols like Pi, E that have constant values
}

// Symbol specifications organized by symbol name
var symbolSpecs = []SymbolSpec{
	{
		Name: "Boole",
		Functions: []rule{
			{"(x_)", stdlib.Boole},
		},
	},
	{
		Name: "RandomReal",
		Functions: []rule{
			{"()", builtins.RandomReal},
		},
	},
	{
		Name: "PatternSpecificity",
		Functions: []rule{
			{"(pattern_)", builtins.PatternSpecificity},
		},
	},
	{
		Name: "ShowPatterns",
		Functions: []rule{
			{"(functionName_Symbol)", builtins.ShowPatterns},
		},
	},
	{
		Name: "Blank",
	},
	{
		Name: "BlankSequence",
	},
	{
		Name: "BlankNullSequence",
	},
	{
		Name:       "Pattern",
		Attributes: []string{"HoldAll"},
	},
	{
		Name:       "Assert",
		Attributes: []string{"HoldAll"},
		Functions: []rule{
			{"(x_)", builtins.Assert},
		},
	},
	{
		Name:       "Plus",
		Attributes: []string{"Flat", "Listable", "NumericFunction", "OneIdentity", "Orderless", "Protected"},
		Functions: []rule{
			{"(x___)", stdlib.PlusExpr},
		},
	},
	{
		Name:       "Times",
		Attributes: []string{"Flat", "Orderless", "OneIdentity"},
		Functions: []rule{
			{"(x___)", stdlib.TimesExpr},
		},
	},
	{
		Name:       "Power",
		Attributes: []string{"OneIdentity"},
		Functions: []rule{
			{"(x_Integer, y_Integer)", stdlib.PowerInteger},
			{"(x_Number, y_Number)", stdlib.PowerExprs},
		},
	},
	{
		Name:       "Subtract",
		Attributes: []string{},
		Functions: []rule{
			{"(x_Integer, y_Integer)", stdlib.SubtractIntegers},
			{"(x_Number, y_Number)", stdlib.SubtractExprs},
		},
	},
	{
		Name:       "Minus",
		Attributes: []string{},
		Functions: []rule{
			{"(x_Integer)", stdlib.MinusInteger},
			{"(x_Real)", stdlib.MinusReal},
			{"(x_)", stdlib.MinusExpr},
		},
	},
	{
		Name:       "Divide",
		Attributes: []string{},
		Functions: []rule{
			{"(x_Integer, y_Integer)", stdlib.DivideIntegers},
			{"(x_Number, y_Number)", stdlib.DivideExprs},
		},
	},
	{
		Name:       "Equal",
		Attributes: []string{},
		Functions: []rule{
			{"(x_Integer, y_Integer)", stdlib.EqualInts},
			{"(x_Real, y_Real)", stdlib.EqualFloats},
			{"(x_Number, y_Number)", stdlib.EqualNumbers},
			{"(x_String, y_String)", stdlib.EqualStrings},
			{"(x_, y_)", stdlib.EqualExprs},
		},
	},
	{
		Name:       "Unequal",
		Attributes: []string{},
		Functions: []rule{
			{"(x_Integer, y_Integer)", stdlib.UnequalInts},
			{"(x_Real, y_Real)", stdlib.UnequalFloats},
			{"(x_Number, y_Number)", stdlib.UnequalNumbers},
			{"(x_String, y_String)", stdlib.UnequalStrings},
			{"(x_, y_)", stdlib.UnequalExprs},
		},
	},
	{
		Name:       "Less",
		Attributes: []string{},
		Functions: []rule{
			{"(x_Number, y_Number)", stdlib.LessNumber},
		},
	},
	{
		Name:       "Greater",
		Attributes: []string{},
		Functions: []rule{
			{"(x_Number, y_Number))", stdlib.GreaterNumber},
		},
	},
	{
		Name:       "LessEqual",
		Attributes: []string{},
		Functions: []rule{
			{"(x_Number, y_Number)", stdlib.LessEqualNumber},
		},
	},
	{
		Name:       "GreaterEqual",
		Attributes: []string{},
		Functions: []rule{
			{"(x_Number, y_Number)", stdlib.GreaterEqualNumber},
		},
	},
	{
		Name:       "SameQ",
		Attributes: []string{},
		Functions: []rule{
			{"(x_, y_)", stdlib.SameQExprs},
		},
	},
	{
		Name:       "UnsameQ",
		Attributes: []string{},
		Functions: []rule{
			{"(x_, y_)", stdlib.UnsameQExprs},
		},
	},

	// Type Predicates
	{
		Name:       "IntegerQ",
		Attributes: []string{},
		Functions: []rule{
			{"(x_)", stdlib.IntegerQExpr},
		},
	},
	{
		Name:       "FloatQ",
		Attributes: []string{},
		Functions: []rule{
			{"(x_)", stdlib.FloatQExpr},
		},
	},
	{
		Name:       "NumberQ",
		Attributes: []string{},
		Functions: []rule{
			{"(x_)", stdlib.NumberQExpr},
		},
	},
	{
		Name:       "StringQ",
		Attributes: []string{},
		Functions: []rule{
			{"(x_)", stdlib.StringQExpr},
		},
	},
	{
		Name:       "BooleanQ",
		Attributes: []string{},
		Functions: []rule{
			{"(x_)", stdlib.BooleanQExpr},
		},
	},
	{
		Name:       "SymbolQ",
		Attributes: []string{},
		Functions: []rule{
			{"(x_)", stdlib.SymbolQExpr},
		},
	},
	{
		Name:       "TrueQ",
		Attributes: []string{},
		Functions: []rule{
			{"(x_)", stdlib.TrueQExpr},
		},
	},
	{
		Name:       "ListQ",
		Attributes: []string{},
		Functions: []rule{
			{"(x_)", stdlib.ListQExpr},
		},
	},
	{
		Name:       "AtomQ",
		Attributes: []string{},
		Functions: []rule{
			{"(x_)", stdlib.AtomQExpr},
		},
	},
	{
		Name:       "Head",
		Attributes: []string{},
		Functions: []rule{
			{"(x_)", stdlib.HeadExpr},
		},
	},

	// Output Format Functions
	{
		Name:       "FullForm",
		Attributes: []string{},
		Functions: []rule{
			{"(x_)", stdlib.FullFormExpr},
		},
	},
	{
		Name:       "InputForm",
		Attributes: []string{},
		Functions: []rule{
			{"(x_)", stdlib.InputFormExpr},
		},
	},

	// List Operations
	{
		Name:       "Length",
		Attributes: []string{},
		Functions: []rule{
			{"(x_)", stdlib.LengthExpr},
		},
	},
	{
		Name:       "First",
		Attributes: []string{},
		Functions: []rule{
			{"(x_List)", stdlib.FirstExpr},
			{"(x_)", stdlib.First},
		},
	},
	{
		Name:       "Last",
		Attributes: []string{},
		Functions: []rule{
			{"(x_List)", stdlib.LastExpr},
			{"(x_)", stdlib.Last},
		},
	},
	{
		Name:       "Rest",
		Attributes: []string{},
		Functions: []rule{
			{"(x_List)", stdlib.RestExpr},
			{"(x_)", stdlib.Rest},
		},
	},
	{
		Name:       "Most",
		Attributes: []string{},
		Functions: []rule{
			{"(x_List)", stdlib.MostExpr},
			{"(x_)", stdlib.Most},
		},
	},
	{
		Name:       "Append",
		Attributes: []string{},
		Functions: []rule{
			{"(x_List, y_)", stdlib.ListAppend},
			{"(x_String, y_String)", stdlib.StringAppend},
		},
	},
	{
		Name:       "Flatten",
		Attributes: []string{},
		Functions: []rule{
			{"(x_)", stdlib.FlattenExpr},
		},
	},
	{
		Name:       "Sort",
		Attributes: []string{},
		Functions: []rule{
			{"(x_)", stdlib.Sort},
		},
	},

	// Sequence Operations
	{
		Name:       "Take",
		Attributes: []string{},
		Functions: []rule{
			{"(x_, n_Integer)", stdlib.Take},
			{"(x_, List(n_Integer, m_Integer))", stdlib.TakeRange},
		},
	},
	{
		Name:       "Drop",
		Attributes: []string{},
		Functions: []rule{
			{"(x_, n_Integer)", stdlib.Drop},
			{"(x_, List(n_Integer, m_Integer))", stdlib.DropRange},
		},
	},
	{
		Name:       "Part",
		Attributes: []string{},
		Functions: []rule{
			{"(x_, n_Integer)", stdlib.Part},
			{"(x_Association, y_)", stdlib.PartAssociation},
		},
	},
	{
		Name:       "Reverse",
		Attributes: []string{},
		Functions: []rule{
			{"(x_String)", stdlib.StringReverse},
		},
	},
	{
		Name:       "RotateLeft",
		Attributes: []string{},
		Functions: []rule{
			{"(x_, n_Integer)", stdlib.RotateLeft},
		},
	},
	{
		Name:       "RotateRight",
		Attributes: []string{},
		Functions: []rule{
			{"(x_, n_Integer)", stdlib.RotateRight},
		},
	},

	// Logical Operations
	{
		Name:       "Not",
		Attributes: []string{},
		Functions: []rule{
			{"(x_)", stdlib.NotExpr},
		},
	},
	{
		Name:       "MatchQ",
		Attributes: []string{},
		Functions: []rule{
			{"(x_, y_)", stdlib.MatchQExprs},
		},
	},

	// Control Structures and Special Forms
	{
		Name:       "If",
		Attributes: []string{"HoldRest"},
		Functions: []rule{
			{"(args___)", builtins.IfExpr},
		},
	},
	{
		Name:       "Set",
		Attributes: []string{"HoldFirst"},
		Functions: []rule{
			{"(lhs_, rhs_)", builtins.SetExpr},
		},
	},
	{
		Name:       "Unset",
		Attributes: []string{"HoldFirst"},
		Functions: []rule{
			{"(lhs_)", builtins.Unset},
		},
	},
	{
		Name:       "SetDelayed",
		Attributes: []string{"HoldAll"},
		Functions: []rule{
			{"(lhs_, rhs_)", builtins.SetDelayedExpr},
		},
	},
	{
		Name:       "Hold",
		Attributes: []string{"HoldAll"},
		Functions: []rule{
			{"(args___)", builtins.HoldExpr},
		},
	},
	{
		Name:       "Evaluate",
		Attributes: []string{},
		Functions: []rule{
			{"(expr_)", builtins.EvaluateExpr},
		},
	},
	{
		Name:       "Timing",
		Attributes: []string{"HoldAll"},
		Functions: []rule{
			{"(expr_)", builtins.Timing},
		},
	},
	{
		Name:       "CompoundExpression",
		Attributes: []string{"HoldAll"},
		Functions: []rule{
			{"(args___)", builtins.CompoundExpression},
		},
	},
	{
		Name:       "And",
		Attributes: []string{"HoldAll"},
		Functions: []rule{
			{"(args___)", builtins.AndExpr},
		},
	},
	{
		Name:       "Or",
		Attributes: []string{"HoldAll"},
		Functions: []rule{
			{"(args___)", builtins.OrExpr},
		},
	},
	{
		Name:       "Function",
		Attributes: []string{"HoldAll"},
		Functions: []rule{
			{"(args__)", builtins.Function},
		},
	},
	{
		Name:       "Block",
		Attributes: []string{"HoldAll"},
		Functions: []rule{
			{"(vars_, body_)", builtins.BlockExpr},
		},
	},
	{
		Name:       "With",
		Attributes: []string{"HoldAll"},
		Functions: []rule{
			{"(vars_, body_)", builtins.With},
		},
	},
	{
		Name:       "Module",
		Attributes: []string{"HoldAll"},
		Functions: []rule{
			{"(vars_, body_)", builtins.Module},
		},
	},
	{
		Name:       "RuleDelayed",
		Attributes: []string{"HoldRest"},
		Functions: []rule{
			{"(lhs_, rhs_)", builtins.RuleDelayed},
		},
	},
	{
		Name:       "Replace",
		Attributes: []string{},
		Functions: []rule{
			{"(expr_, rule_)", builtins.Replace},
		},
	},
	{
		Name:       "ReplaceAll",
		Attributes: []string{},
		Functions: []rule{
			{"(expr_, rule_)", builtins.ReplaceAll},
		},
	},
	{
		Name:       "Table",
		Attributes: []string{"HoldAll"},
		Functions: []rule{
			{"(expr_, iterator_)", builtins.Table},
		},
	},
	{
		Name:       "Do",
		Attributes: []string{"HoldAll"},
		Functions: []rule{
			{"(expr_, iterator_)", builtins.Do},
		},
	},

	// String Operations
	{
		Name:       "StringLength",
		Attributes: []string{},
		Functions: []rule{
			{"(x_String)", stdlib.StringLengthRunes},
		},
	},
	{
		Name:       "ByteArray",
		Attributes: []string{},
		Functions: []rule{
			{"(x_String)", stdlib.ByteArrayFromString},
		},
	},

	// Association Operations
	{
		Name:       "AssociationQ",
		Attributes: []string{},
		Functions: []rule{
			{"(x_)", stdlib.AssociationQExpr},
		},
	},
	{
		Name:       "Keys",
		Attributes: []string{},
		Functions: []rule{
			{"(x_Association)", stdlib.KeysExpr},
		},
	},
	{
		Name:       "Values",
		Attributes: []string{},
		Functions: []rule{
			{"(x_Association)", stdlib.ValuesExpr},
		},
	},
	{
		Name:       "Association",
		Attributes: []string{},
		Functions: []rule{
			{"(x___Rule)", stdlib.AssociationRules},
		},
	},

	// Output Operations
	{
		Name:       "Print",
		Attributes: []string{},
		Functions: []rule{
			{"(x_)", stdlib.Print},
		},
	},

	// Attribute manipulation functions
	{
		Name:       "SetAttributes",
		Attributes: []string{"HoldFirst"},
		Functions: []rule{
			{"(symbol_, attrs_)", builtins.SetAttributesExpr},
		},
	},
	{
		Name:       "ClearAttributes",
		Attributes: []string{"HoldFirst"},
		Functions: []rule{
			{"(symbol_, attrs_)", builtins.ClearAttributesExpr},
		},
	},
	{
		Name:       "Attributes",
		Attributes: []string{"HoldFirst"},
		Functions: []rule{
			{"(symbol_)", builtins.AttributesExpr},
		},
	},

	// Functional Programming Functions
	{
		Name:       "Map",
		Attributes: []string{},
		Functions: []rule{
			{"(f_, list_)", builtins.MapExpr},
		},
	},
	{
		Name:       "Apply",
		Attributes: []string{},
		Functions: []rule{
			{"(f_, list_)", builtins.ApplyExpr},
		},
	},
}

/*
	 TODO
	 {
			Name:       "Function",
			Attributes: []string{"HoldAll"},
			Functions: []rule{
				"(args___)": builtins.FunctionExpr,
			},
		},
*/
