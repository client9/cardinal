package main

import (
	"reflect"
	"testing"
)

func TestParseSymbolSpecs(t *testing.T) {
	source := []byte(`
package builtins

// @ExprSymbol And
// @ExprAttributes Protected HoldAll

// @ExprPattern (___)
func AndExpr(args []interface{}) interface{} {
	return nil
}

// @ExprSymbol Plus
// @ExprAttributes Flat Orderless OneIdentity

// @ExprPattern (x__Integer)
func PlusIntegers(x []int) int {
	sum := 0
	for _, v := range x {
		sum += v
	}
	return sum
}

// @ExprPattern (x__Real)
func PlusReals(x []float64) float64 {
	sum := 0.0
	for _, v := range x {
		sum += v
	}
	return sum
}

// Regular function without pattern
func HelperFunc() {}

// @ExprSymbol Times
// @ExprAttributes Flat Orderless OneIdentity

// No function for Times symbol
`)

	symbols, err := ParseSymbolSpecs(source)
	if err != nil {
		t.Fatalf("ParseSymbolSpecs failed: %v", err)
	}

	expected := []SymbolSpec{
		{
			Name:       "And",
			Attributes: []string{"Protected", "HoldAll"},
			Functions: []Rule{
				{
					Pattern:  "(___)",
					Function: "builtins.AndExpr",
				},
			},
		},
		{
			Name:       "Plus",
			Attributes: []string{"Flat", "Orderless", "OneIdentity"},
			Functions: []Rule{
				{
					Pattern:  "(x__Integer)",
					Function: "builtins.PlusIntegers",
				},
				{
					Pattern:  "(x__Real)",
					Function: "builtins.PlusReals",
				},
			},
		},
		{
			Name:       "Times",
			Attributes: []string{"Flat", "Orderless", "OneIdentity"},
			Functions:  []Rule{},
		},
	}

	if len(symbols) != len(expected) {
		t.Fatalf("Expected %d symbols, got %d", len(expected), len(symbols))
	}

	for i, symbol := range symbols {
		if symbol.Name != expected[i].Name {
			t.Errorf("Symbol %d: expected name %s, got %s", i, expected[i].Name, symbol.Name)
		}
		if !reflect.DeepEqual(symbol.Attributes, expected[i].Attributes) {
			t.Errorf("Symbol %d (%s): expected attributes %v, got %v", i, symbol.Name, expected[i].Attributes, symbol.Attributes)
		}
		if len(symbol.Functions) != len(expected[i].Functions) {
			t.Errorf("Symbol %d (%s): expected %d functions, got %d", i, symbol.Name, len(expected[i].Functions), len(symbol.Functions))
			continue
		}
		for j, rule := range symbol.Functions {
			if rule.Pattern != expected[i].Functions[j].Pattern {
				t.Errorf("Symbol %d (%s), Rule %d: expected pattern %s, got %s", i, symbol.Name, j, expected[i].Functions[j].Pattern, rule.Pattern)
			}
			if rule.Function != expected[i].Functions[j].Function {
				t.Errorf("Symbol %d (%s), Rule %d: expected function %s, got %s", i, symbol.Name, j, expected[i].Functions[j].Function, rule.Function)
			}
		}
	}
}

func TestParseSymbolSpecsEmpty(t *testing.T) {
	source := []byte(`
package test

// Regular comment
func regularFunc() {}
`)

	symbols, err := ParseSymbolSpecs(source)
	if err != nil {
		t.Fatalf("ParseSymbolSpecs failed: %v", err)
	}

	if len(symbols) != 0 {
		t.Fatalf("Expected 0 symbols, got %d", len(symbols))
	}
}
