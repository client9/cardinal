package stdlib

import (
	"fmt"
	"strings"
	
	"github.com/client9/sexpr/core"
)

// AssociationEntry represents a key-value pair in an association
type AssociationEntry struct {
	Key   core.Expr
	Value core.Expr
}

// AssociationValue implements Expr for association data structure
// Uses hash buckets for efficient lookup while preserving insertion order
type AssociationValue struct {
	buckets map[string][]AssociationEntry // Hash buckets for lookup
	order   []core.Expr                   // Preserve insertion order of keys
}

// NewAssociationValue creates a new empty association
func NewAssociationValue() AssociationValue {
	return AssociationValue{
		buckets: make(map[string][]AssociationEntry),
		order:   []core.Expr{},
	}
}

// NewAssociationFromPairs creates an association from alternating key-value pairs
func NewAssociationFromPairs(pairs ...core.Expr) (AssociationValue, error) {
	if len(pairs)%2 != 0 {
		return AssociationValue{}, fmt.Errorf("association requires even number of arguments (key-value pairs)")
	}

	assoc := NewAssociationValue()
	for i := 0; i < len(pairs); i += 2 {
		key := pairs[i]
		value := pairs[i+1]
		assoc = assoc.Set(key, value) // Now returns new association
	}
	return assoc, nil
}

// hashKey computes a hash string for any expression
func hashKey(key core.Expr) string {
	return key.String()
}

// Set adds or updates a key-value pair and returns a new AssociationValue
// This maintains immutability by creating a new association instead of modifying in place
func (a AssociationValue) Set(key core.Expr, value core.Expr) AssociationValue {
	hash := hashKey(key)

	// Create a copy of the association
	newAssoc := AssociationValue{
		buckets: make(map[string][]AssociationEntry, len(a.buckets)),
		order:   make([]core.Expr, len(a.order)),
	}

	// Copy order slice
	copy(newAssoc.order, a.order)

	// Copy buckets (shallow copy of slices, but new slice containers)
	for k, bucket := range a.buckets {
		newBucket := make([]AssociationEntry, len(bucket))
		copy(newBucket, bucket)
		newAssoc.buckets[k] = newBucket
	}

	keyExists := false

	// Check if key already exists in bucket
	if bucket, exists := newAssoc.buckets[hash]; exists {
		for i, entry := range bucket {
			if entry.Key.Equal(key) {
				// Update existing entry in the new copy
				bucket[i].Value = value
				keyExists = true
				break
			}
		}
		if !keyExists {
			// Key not found in bucket, append new entry
			newAssoc.buckets[hash] = append(bucket, AssociationEntry{Key: key, Value: value})
		}
	} else {
		// New bucket
		newAssoc.buckets[hash] = []AssociationEntry{{Key: key, Value: value}}
	}

	// Add to order if it's a new key
	if !keyExists {
		newAssoc.order = append(newAssoc.order, key)
	}

	return newAssoc
}

// Get retrieves a value by key
func (a AssociationValue) Get(key core.Expr) (core.Expr, bool) {
	hash := hashKey(key)

	if bucket, exists := a.buckets[hash]; exists {
		for _, entry := range bucket {
			if entry.Key.Equal(key) {
				return entry.Value, true
			}
		}
	}
	return nil, false
}

// Keys returns all keys in insertion order
func (a AssociationValue) Keys() []core.Expr {
	result := make([]core.Expr, 0, len(a.order))

	// Filter out duplicates while preserving order
	// Note: With immutable Set operations, duplicates in order should be rare
	seen := make(map[string]bool, len(a.order))
	for _, key := range a.order {
		hash := hashKey(key)
		if !seen[hash] {
			// Since associations are now immutable, we can trust the order
			// Only check existence in buckets, which is more efficient than Get()
			if bucket, exists := a.buckets[hash]; exists {
				// Verify the key actually exists in the bucket
				for _, entry := range bucket {
					if entry.Key.Equal(key) {
						result = append(result, key)
						seen[hash] = true
						break
					}
				}
			}
		}
	}
	return result
}

// Values returns all values in key insertion order
func (a AssociationValue) Values() []core.Expr {
	keys := a.Keys()
	result := make([]core.Expr, len(keys))

	for i, key := range keys {
		if value, exists := a.Get(key); exists {
			result[i] = value
		}
	}
	return result
}

// Length returns the number of key-value pairs
func (a AssociationValue) Len() int {
	return len(a.Keys()) // Account for any duplicates in order slice
}

// String implements Expr interface - used for FullForm representation
func (a AssociationValue) String() string {
	if a.Len() == 0 {
		return "Association()"
	}

	var parts []string
	keys := a.Keys()
	for _, key := range keys {
		if value, exists := a.Get(key); exists {
			parts = append(parts, fmt.Sprintf("Rule(%s, %s)", key.String(), value.String()))
		}
	}
	return fmt.Sprintf("Association(%s)", strings.Join(parts, ", "))
}

// InputForm implements Expr interface - used for user-friendly representation
func (a AssociationValue) InputForm() string {
	if a.Len() == 0 {
		return "{}"
	}

	var parts []string
	keys := a.Keys()
	for _, key := range keys {
		if value, exists := a.Get(key); exists {
			parts = append(parts, fmt.Sprintf("%s: %s", key.InputForm(), value.InputForm()))
		}
	}
	return fmt.Sprintf("{%s}", strings.Join(parts, ", "))
}

// Type implements Expr interface
func (a AssociationValue) Type() string {
	return "Association"
}

// Equal implements Expr interface
func (a AssociationValue) Equal(rhs core.Expr) bool {
	rhsAssoc, ok := rhs.(AssociationValue)
	if !ok {
		return false
	}

	// Must have same number of entries
	if a.Len() != rhsAssoc.Len() {
		return false
	}

	// All keys and values must match
	aKeys := a.Keys()
	for _, key := range aKeys {
		aValue, aExists := a.Get(key)
		rhsValue, rhsExists := rhsAssoc.Get(key)

		if !aExists || !rhsExists || !aValue.Equal(rhsValue) {
			return false
		}
	}

	return true
}

// NewAssociation creates an ObjectExpr wrapping an AssociationValue
func NewAssociation(pairs ...core.Expr) (core.ObjectExpr, error) {
	assocValue, err := NewAssociationFromPairs(pairs...)
	if err != nil {
		return core.ObjectExpr{}, err
	}
	return core.NewObjectExpr("Association", assocValue), nil
}

// Association functions - all work with ObjectExpr of type "Association"

// AssociationQExpr checks if an expression is an Association
func AssociationQExpr(expr core.Expr) bool {
	if objExpr, ok := expr.(core.ObjectExpr); ok && objExpr.TypeName == "Association" {
		return true
	}
	return false
}

// KeysExpr returns the keys of an association as a List
func KeysExpr(assoc core.ObjectExpr) core.Expr {
	if assoc.TypeName == "Association" {
		if assocValue, ok := assoc.Value.(AssociationValue); ok {
			keys := assocValue.Keys()
			// Return as List[key1, key2, ...]
			elements := []core.Expr{core.NewSymbolAtom("List")}
			elements = append(elements, keys...)
			return core.NewList(elements...)
		}
	}
	// This should not happen with proper pattern matching, but return error as fallback
	return core.NewErrorExpr("ArgumentError",
		fmt.Sprintf("Keys expects an Association, got %s", assoc.String()), []core.Expr{assoc})
}

// ValuesExpr returns the values of an association as a List
func ValuesExpr(assoc core.ObjectExpr) core.Expr {
	if assoc.TypeName == "Association" {
		if assocValue, ok := assoc.Value.(AssociationValue); ok {
			values := assocValue.Values()
			// Return as List[value1, value2, ...]
			elements := []core.Expr{core.NewSymbolAtom("List")}
			elements = append(elements, values...)
			return core.NewList(elements...)
		}
	}
	// This should not happen with proper pattern matching, but return error as fallback
	return core.NewErrorExpr("ArgumentError",
		fmt.Sprintf("Values expects an Association, got %s", assoc.String()), []core.Expr{assoc})
}

// AssociationRules creates an Association from a sequence of Rule expressions
func AssociationRules(rules ...core.Expr) core.Expr {
	assocValue := NewAssociationValue()

	// Process each Rule expression
	for _, rule := range rules {
		if ruleList, ok := rule.(core.List); ok && len(ruleList.Elements) == 3 {
			// Check if this is Rule[key, value]
			if headAtom, ok := ruleList.Elements[0].(core.Atom); ok &&
				headAtom.AtomType == core.SymbolAtom && headAtom.Value.(string) == "Rule" {

				key := ruleList.Elements[1]
				value := ruleList.Elements[2]
				assocValue = assocValue.Set(key, value) // Returns new association (immutable)
				continue
			}
		}

		// Invalid argument - not a Rule
		return core.NewErrorExpr("ArgumentError",
			fmt.Sprintf("Association expects Rule expressions, got %s", rule.String()), []core.Expr{rule})
	}

	return core.NewObjectExpr("Association", assocValue)
}

// PartAssociation extracts a value from an association by key
func PartAssociation(assoc core.ObjectExpr, key core.Expr) core.Expr {
	if assoc.TypeName == "Association" {
		if assocValue, ok := assoc.Value.(AssociationValue); ok {
			// For associations, use the key argument to lookup value
			if value, exists := assocValue.Get(key); exists {
				return value
			}
			return core.NewErrorExpr("PartError",
				fmt.Sprintf("Key %s not found in association", key.String()), []core.Expr{assoc, key})
		}
	}
	// This should not happen with proper pattern matching, but return error as fallback
	return core.NewErrorExpr("ArgumentError",
		fmt.Sprintf("Part expects an Association, got %s", assoc.String()), []core.Expr{assoc})
}