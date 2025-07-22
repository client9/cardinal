package sexpr

import (
	"fmt"
	"strings"
)

// AssociationEntry represents a key-value pair in an association
type AssociationEntry struct {
	Key   Expr
	Value Expr
}

// AssociationValue implements Expr for association data structure
// Uses hash buckets for efficient lookup while preserving insertion order
type AssociationValue struct {
	buckets map[string][]AssociationEntry // Hash buckets for lookup
	order   []Expr                        // Preserve insertion order of keys
}

// NewAssociationValue creates a new empty association
func NewAssociationValue() AssociationValue {
	return AssociationValue{
		buckets: make(map[string][]AssociationEntry),
		order:   []Expr{},
	}
}

// NewAssociationFromPairs creates an association from alternating key-value pairs
func NewAssociationFromPairs(pairs ...Expr) (AssociationValue, error) {
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
func hashKey(key Expr) string {
	return key.String()
}

// Set adds or updates a key-value pair and returns a new AssociationValue
// This maintains immutability by creating a new association instead of modifying in place
func (a AssociationValue) Set(key Expr, value Expr) AssociationValue {
	hash := hashKey(key)

	// Create a copy of the association
	newAssoc := AssociationValue{
		buckets: make(map[string][]AssociationEntry, len(a.buckets)),
		order:   make([]Expr, len(a.order)),
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
func (a AssociationValue) Get(key Expr) (Expr, bool) {
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
func (a AssociationValue) Keys() []Expr {
	result := make([]Expr, 0, len(a.order))

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
func (a AssociationValue) Values() []Expr {
	keys := a.Keys()
	result := make([]Expr, len(keys))

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

// String implements Expr interface
func (a AssociationValue) String() string {
	if a.Len() == 0 {
		return "{}"
	}

	var parts []string
	keys := a.Keys()
	for _, key := range keys {
		if value, exists := a.Get(key); exists {
			parts = append(parts, fmt.Sprintf("%s: %s", key.String(), value.String()))
		}
	}
	return fmt.Sprintf("{%s}", strings.Join(parts, ", "))
}

// InputForm implements Expr interface
func (a AssociationValue) InputForm() string {
	return a.String() // For associations, InputForm is the same as String() (already uses {key: value} format)
}

// Type implements Expr interface
func (a AssociationValue) Type() string {
	return "Association"
}

// Equal implements Expr interface
func (a AssociationValue) Equal(rhs Expr) bool {
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
func NewAssociation(pairs ...Expr) (ObjectExpr, error) {
	assocValue, err := NewAssociationFromPairs(pairs...)
	if err != nil {
		return ObjectExpr{}, err
	}
	return NewObjectExpr("Association", assocValue), nil
}
