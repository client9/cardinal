# Builtin Function Development Guidelines

## Immutability Patterns for Go-based S-Expression Systems

This document provides guidelines for writing builtin functions in our Go-based s-expression system, inspired by Mathematica. The focus is on maintaining immutability while leveraging Go's unique memory management properties.

## Core Philosophy: Structural Immutability

Our system follows **structural immutability** - expressions should never be modified in place. Instead, create new expressions that represent the result of operations.

## Go Memory Management Context

### What Go Handles Automatically

1. **Slices and Maps**: Passed by reference automatically
   - `[]Expr` slices share underlying arrays
   - Map structures share reference semantics
   - This is beneficial for our read-heavy operations

2. **Value Types**: Our core types are designed as values
   - `Atom` - struct passed by value (small, efficient)
   - `List` - struct with `[]Expr` slice (reference to slice, but struct by value)
   - This gives us copy semantics for the container, reference semantics for contents

### What We Must Manage

1. **Creating new structures** instead of modifying existing ones
2. **Sharing immutable parts** efficiently
3. **Avoiding unnecessary allocations** while maintaining immutability

## Do's and Don'ts for Builtin Functions

### ✅ DO: Always Return New Expressions

```go
// ✅ GOOD: Create new list with modified elements
func EvaluateAppend(args []Expr) Expr {
    if len(args) < 2 {
        return NewErrorExpr("ArgumentError", "Append expects at least 2 arguments", args)
    }
    
    list := args[0]
    newElements := args[1:]
    
    if l, ok := list.(List); ok {
        // Create new slice with additional capacity
        newSlice := make([]Expr, len(l.Elements), len(l.Elements)+len(newElements))
        copy(newSlice, l.Elements)
        newSlice = append(newSlice, newElements...)
        
        return List{Elements: newSlice} // New List, reuses existing elements
    }
    
    return NewErrorExpr("TypeError", "First argument must be a list", args)
}
```

### ❌ DON'T: Modify Arguments In Place

```go
// ❌ BAD: Modifying the input list
func EvaluateAppendBad(args []Expr) Expr {
    list := args[0].(List)
    list.Elements = append(list.Elements, args[1:]...) // MUTATES INPUT!
    return list
}
```

### ✅ DO: Share Immutable Substructures

```go
// ✅ GOOD: Reuse unchanged elements
func EvaluateRest(args []Expr) Expr {
    if len(args) != 1 {
        return NewErrorExpr("ArgumentError", "Rest expects 1 argument", args)
    }
    
    if list, ok := args[0].(List); ok {
        if len(list.Elements) <= 1 {
            return NewErrorExpr("PartError", "Rest: expression has no elements", args)
        }
        
        // Share the head element, create new slice for rest
        newElements := make([]Expr, len(list.Elements)-1)
        newElements[0] = list.Elements[0]          // Share head
        copy(newElements[1:], list.Elements[2:])   // Share tail elements
        
        return List{Elements: newElements}
    }
    
    return NewErrorExpr("PartError", "Rest: expression is not a list", args)
}
```

### ✅ DO: Leverage Go Slice Sharing for Read Operations

```go
// ✅ GOOD: Efficient iteration without copying
func EvaluateLength(args []Expr) Expr {
    if len(args) != 1 {
        return NewErrorExpr("ArgumentError", "Length expects 1 argument", args)
    }
    
    if list, ok := args[0].(List); ok {
        // No copying needed - just read the length
        return NewIntAtom(len(list.Elements) - 1) // Subtract 1 for head
    }
    
    return NewIntAtom(0)
}
```

### ✅ DO: Use Copy-on-Write for Large Structures

```go
// ✅ GOOD: Efficient copying only when needed
func EvaluateReplacePart(args []Expr) Expr {
    if len(args) != 3 {
        return NewErrorExpr("ArgumentError", "ReplacePart expects 3 arguments", args)
    }
    
    expr := args[0]
    index := args[1]
    newValue := args[2]
    
    if list, ok := expr.(List); ok {
        if indexAtom, ok := index.(Atom); ok && indexAtom.AtomType == IntAtom {
            idx := indexAtom.Value.(int)
            
            if idx >= 1 && idx < len(list.Elements) {
                // Copy-on-write: only allocate new slice when modifying
                newElements := make([]Expr, len(list.Elements))
                copy(newElements, list.Elements) // Shallow copy
                newElements[idx] = newValue       // Replace single element
                
                return List{Elements: newElements}
            }
        }
    }
    
    return NewErrorExpr("PartError", "Invalid index", args)
}
```

### ❌ DON'T: Create Unnecessary Deep Copies

```go
// ❌ BAD: Unnecessary deep copying
func EvaluateFirstBad(args []Expr) Expr {
    list := args[0].(List)
    
    // This creates unnecessary copies of all elements
    newList := List{Elements: make([]Expr, len(list.Elements))}
    for i, elem := range list.Elements {
        newList.Elements[i] = deepCopy(elem) // WASTEFUL!
    }
    
    return newList.Elements[1] // Only needed one element!
}

// ✅ GOOD: Direct access without copying
func EvaluateFirst(args []Expr) Expr {
    if list, ok := args[0].(List); ok {
        if len(list.Elements) <= 1 {
            return NewErrorExpr("PartError", "First: expression has no elements", args)
        }
        return list.Elements[1] // Direct return, no copying
    }
    return NewErrorExpr("PartError", "First: expression is not a list", args)
}
```

## Patterns for Different Operations

### 1. Read-Only Operations (Length, Head, Part)

- **Pattern**: Direct access, no allocation
- **Go Advantage**: Slice access is O(1), no copying needed

```go
func EvaluateHead(args []Expr) Expr {
    expr := args[0]
    switch ex := expr.(type) {
    case List:
        if len(ex.Elements) == 0 {
            return NewSymbolAtom("List")
        }
        return ex.Elements[0] // Just return reference
    case Atom:
        // Return type symbol based on atom type
        // No copying, just type analysis
    }
}
```

### 2. Structural Modifications (Append, Prepend, Insert)

- **Pattern**: Create new container, share unchanged elements
- **Go Advantage**: Slice append operations are optimized

```go
func EvaluatePrepend(args []Expr) Expr {
    list := args[0].(List)
    newElement := args[1]
    
    // Efficient prepend: allocate new slice with extra capacity
    newElements := make([]Expr, len(list.Elements)+1, len(list.Elements)*2)
    newElements[0] = list.Elements[0]  // Head
    newElements[1] = newElement         // New element
    copy(newElements[2:], list.Elements[1:]) // Original elements
    
    return List{Elements: newElements}
}
```

### 3. Transformations (Map, Select, Replace)

- **Pattern**: Process and conditionally share
- **Go Advantage**: Range iteration is efficient

```go
func EvaluateMap(args []Expr) Expr {
    fn := args[0]
    list := args[1].(List)
    
    // Pre-allocate result slice
    resultElements := make([]Expr, len(list.Elements))
    resultElements[0] = list.Elements[0] // Share head
    
    // Transform elements 1 by 1
    for i := 1; i < len(list.Elements); i++ {
        // Apply function to each element
        result := applyFunction(fn, list.Elements[i])
        resultElements[i] = result
    }
    
    return List{Elements: resultElements}
}
```

### 4. Association Operations

- **Pattern**: Leverage Go's map efficiency
- **Go Advantage**: Maps are reference types with copy-on-write semantics

```go
func EvaluateAssociationMerge(args []Expr) Expr {
    assoc1 := args[0].(ObjectExpr).Value.(AssociationValue)
    assoc2 := args[1].(ObjectExpr).Value.(AssociationValue)
    
    // Create new association - Go map will handle efficient copying
    newAssoc := NewAssociationValue()
    
    // Copy entries (Go map iteration is efficient)
    for key, value := range assoc1.data {
        newAssoc.Set(key, value) // Shares key and value references
    }
    
    for key, value := range assoc2.data {
        newAssoc.Set(key, value) // Overwrites or adds
    }
    
    return NewObjectExpr("Association", newAssoc)
}
```

## Memory Efficiency Guidelines

### 1. Pre-allocate When Size is Known

```go
// ✅ GOOD: Avoid repeated allocations
newElements := make([]Expr, expectedSize, capacity)

// ❌ BAD: Repeated allocations
var newElements []Expr
for _, elem := range input {
    newElements = append(newElements, transform(elem))
}
```

### 2. Share References to Immutable Data

```go
// ✅ GOOD: All atoms and symbols are immutable, safe to share
return List{Elements: []Expr{
    list.Elements[0], // Shared head
    NewIntAtom(42),   // New atom
    list.Elements[2], // Shared element
}}
```

### 3. Use Go's Copy Function for Slices

```go
// ✅ GOOD: Efficient bulk copying
newSlice := make([]Expr, len(oldSlice))
copy(newSlice, oldSlice)

// ❌ BAD: Element-by-element copying
newSlice := make([]Expr, len(oldSlice))
for i, elem := range oldSlice {
    newSlice[i] = elem
}
```

## Error Handling and Immutability

### Always Return New Error Expressions

```go
// ✅ GOOD: Create new error, don't modify arguments
func EvaluateDivide(args []Expr) Expr {
    if len(args) != 2 {
        return NewErrorExpr("ArgumentError",
            fmt.Sprintf("Divide expects 2 arguments, got %d", len(args)), args)
    }
    
    // args slice is not modified, safe to pass as reference
    if /* division by zero */ {
        return NewErrorExpr("DivisionByZero", "Cannot divide by zero", args)
    }
    
    // Return new result
    return NewFloatAtom(result)
}
```

## Pattern Matching and Immutability

When implementing pattern-based functions, ensure pattern variables don't create aliases to mutable state:

```go
// ✅ GOOD: Pattern matching with immutable bindings
func EvaluatePatternFunction(args []Expr, ctx *Context) Expr {
    // Pattern variables are bound to immutable expressions
    // Safe to share references in variable bindings
    tempCtx := NewChildContext(ctx) // New context, don't modify parent
    
    if matches := matchPattern(pattern, expr, tempCtx); matches {
        // Variables in tempCtx are safe references to immutable expressions
        return evaluateWithBindings(body, tempCtx)
    }
    
    return expr // Return unchanged
}
```

## Testing Immutability

When writing tests, verify that operations don't modify inputs:

```go
func TestImmutability(t *testing.T) {
    original := List{Elements: []Expr{
        NewSymbolAtom("List"),
        NewIntAtom(1),
        NewIntAtom(2),
    }}
    
    // Make a copy to verify original is unchanged
    originalCopy := List{Elements: make([]Expr, len(original.Elements))}
    copy(originalCopy.Elements, original.Elements)
    
    // Perform operation
    result := EvaluateAppend([]Expr{original, NewIntAtom(3)})
    
    // Verify original is unchanged
    if !original.Equal(originalCopy) {
        t.Error("Operation modified original list")
    }
    
    // Verify result is correct
    expected := List{Elements: []Expr{
        NewSymbolAtom("List"),
        NewIntAtom(1),
        NewIntAtom(2),
        NewIntAtom(3),
    }}
    
    if !result.Equal(expected) {
        t.Error("Operation produced incorrect result")
    }
}
```

## When to Return Input Directly vs. Copy

### ✅ Safe to Return Input Directly

When the input is immutable and no transformation is needed:

```go
// ✅ GOOD: Identity operations
func EvaluatePlus(args []Expr) Expr {
    if len(args) == 1 {
        return args[0] // Plus[x] = x, return directly
    }
    // ... handle other cases
}

// ✅ GOOD: Conditional returns
func EvaluateIf(args []Expr) Expr {
    if isTrue(condition) {
        return trueExpr  // Direct return of immutable expression
    }
    return falseExpr
}

// ✅ GOOD: Element access  
func EvaluateFirst(args []Expr) Expr {
    return list.Elements[1] // Atoms and expressions are immutable
}
```

### ❌ When to Create New Expressions

Only when the result is structurally different from the input:

```go
// ✅ GOOD: Transformation needed
func EvaluateAppend(args []Expr) Expr {
    list := args[0].(List)
    newElement := args[1]
    
    // Must create new list structure
    newElements := make([]Expr, len(list.Elements)+1)
    copy(newElements, list.Elements)
    newElements[len(list.Elements)] = newElement
    
    return List{Elements: newElements}
}
```

### 🎯 Decision Matrix

| Scenario | Action | Reason |
|----------|--------|--------|
| Identity operation (`Plus[x]` = `x`) | Return input directly | No structural change needed |
| Unchanged result (`If[True, x, y]` = `x`) | Return input directly | Input is immutable |
| Structural change (`Append[list, elem]`) | Create new expression | Result differs from input |
| Type conversion (`ToString[42]`) | Create new expression | Different type/representation |

### 🔍 Our System's Guarantees

In our s-expression system:
- **Atoms**: Always immutable - safe to return directly
- **Lists**: Structure is immutable - safe to return directly  
- **Expressions**: Designed to be immutable - safe to return directly
- **Performance**: Direct returns avoid unnecessary allocations

## Summary

1. **Never modify input arguments** - always return new expressions
2. **Share immutable references** - atoms, symbols, and unchanged sublists
3. **Return inputs directly when unchanged** - leverages immutability for performance
4. **Use Go's slice efficiency** - leverage copy(), append(), and make()
5. **Pre-allocate when possible** - avoid repeated allocations
6. **Test immutability** - verify inputs remain unchanged
7. **Leverage Go's reference semantics** - for maps and slices when reading
8. **Create minimal copies** - only copy what needs to change

This approach gives us the benefits of immutability (predictable behavior, safe concurrency, easy reasoning) while leveraging Go's efficient memory management for performance.