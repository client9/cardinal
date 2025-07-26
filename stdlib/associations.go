package stdlib

import (
	"fmt"

	"github.com/client9/sexpr/core"
)

func AssociationLength(x core.Association) int64 {
	return x.Length()
}

// Association functions - all work with ObjectExpr of type "Association"

// AssociationQExpr checks if an expression is an Association
func AssociationQExpr(expr core.Expr) bool {
	_, ok := expr.(core.Association)
	return ok
}

// KeysExpr returns the keys of an association as a List
func KeysExpr(assoc core.Association) core.Expr {
	keys := assoc.Keys()
	// Return as List[key1, key2, ...]
	elements := []core.Expr{core.NewSymbolAtom("List")}
	elements = append(elements, keys...)
	return core.NewList(elements...)
}

// ValuesExpr returns the values of an association as a List
func ValuesExpr(assoc core.Association) core.Expr {
	values := assoc.Values()
	// Return as List[value1, value2, ...]
	elements := []core.Expr{core.NewSymbolAtom("List")}
	elements = append(elements, values...)
	return core.NewList(elements...)
}

// AssociationRules creates an Association from a sequence of Rule expressions
func AssociationRules(rules ...core.Expr) core.Expr {
	assoc := core.NewAssociation()

	// Process each Rule expression
	for _, rule := range rules {
		if ruleList, ok := rule.(core.List); ok && len(ruleList.Elements) == 3 {
			// Check if this is Rule[key, value]
			if headName, ok := core.ExtractSymbol(ruleList.Elements[0]); ok && headName == "Rule" {

				key := ruleList.Elements[1]
				value := ruleList.Elements[2]
				assoc = assoc.Set(key, value) // Returns new association (immutable)
				continue
			}
		}

		// Invalid argument - not a Rule
		return core.NewErrorExpr("ArgumentError",
			fmt.Sprintf("Association expects Rule expressions, got %s", rule.String()), []core.Expr{rule})
	}

	return assoc
}

// PartAssociation extracts a value from an association by key
func PartAssociation(assoc core.Association, key core.Expr) core.Expr {
	// For associations, use the key argument to lookup value
	if value, exists := assoc.Get(key); exists {
		return value
	}
	return core.NewErrorExpr("PartError",
		fmt.Sprintf("Key %s not found in association", key.String()), []core.Expr{assoc, key})
}
