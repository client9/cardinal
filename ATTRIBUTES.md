# Symbol Attributes System

This document describes the Symbol Attributes system implemented for the s-expression parser, designed to match Mathematica's attribute system.

## Overview

The attribute system allows you to associate properties with symbols that control how they behave during evaluation. This is essential for implementing a Mathematica-style symbolic computation system.

## Supported Attributes

The system supports the following standard Mathematica attributes:

- **HoldAll**: Prevents evaluation of all arguments
- **HoldFirst**: Prevents evaluation of the first argument
- **HoldRest**: Prevents evaluation of all arguments except the first
- **Flat**: Treats the function as associative (e.g., `Plus[a, Plus[b, c]]` → `Plus[a, b, c]`)
- **Orderless**: Treats the function as commutative (arguments can be reordered)
- **OneIdentity**: Function with one argument returns that argument unchanged
- **Listable**: Function automatically threads over lists
- **Constant**: Symbol represents a constant value
- **NumericFunction**: Function always returns numeric values for numeric inputs
- **Protected**: Symbol cannot be redefined
- **ReadProtected**: Symbol definition cannot be read
- **Locked**: Symbol and its attributes cannot be modified
- **Temporary**: Symbol is removed at the end of the session

## API Reference

### Core Functions

```go
// Set attributes for a symbol
SetAttributes("Plus", []Attribute{Flat, Orderless, OneIdentity})

// Clear specific attributes
ClearAttributes("Plus", []Attribute{OneIdentity})

// Clear all attributes
ClearAllAttributes("Plus")

// Get all attributes for a symbol
attrs := Attributes("Plus")

// Check if symbol has specific attribute
hasFlat := HasAttribute("Plus", Flat)

// Get all symbols with attributes
symbols := AllSymbolsWithAttributes()
```

### Utility Functions

```go
// Format attributes for display
str := AttributesToString(attrs) // Returns "{Flat, Orderless}"

// Convert attribute to string
name := Flat.String() // Returns "Flat"

// Reset global symbol table (useful for testing)
ResetGlobalSymbolTable()
```

## Thread Safety

The attribute system is fully thread-safe. All operations use read-write mutexes to ensure concurrent access is safe.

## Usage Examples

### Basic Usage

```go
// Set up mathematical operators
SetAttributes("Plus", []Attribute{Flat, Orderless, OneIdentity})
SetAttributes("Times", []Attribute{Flat, Orderless, OneIdentity})

// Check attributes
if HasAttribute("Plus", Flat) {
    fmt.Println("Plus is associative")
}

// Display all attributes
fmt.Println("Plus attributes:", AttributesToString(Attributes("Plus")))
```

### Common Mathematica Symbols

```go
// Arithmetic operators
SetAttributes("Plus", []Attribute{Flat, Orderless, OneIdentity})
SetAttributes("Times", []Attribute{Flat, Orderless, OneIdentity})

// Control structures
SetAttributes("Hold", []Attribute{HoldAll})
SetAttributes("If", []Attribute{HoldRest})

// Mathematical functions
SetAttributes("Sin", []Attribute{Listable, NumericFunction})
SetAttributes("Cos", []Attribute{Listable, NumericFunction})

// Constants
SetAttributes("Pi", []Attribute{Constant, Protected})
SetAttributes("E", []Attribute{Constant, Protected})
```

### Integration with Evaluator

```go
// Example evaluator logic
func evaluate(expr Expr) Expr {
    if list, ok := expr.(*List); ok && len(list.Elements) > 0 {
        if head, ok := list.Elements[0].(*Atom); ok && head.AtomType == SymbolAtom {
            symbol := head.Value.(string)
            
            // Check if arguments should be held
            if HasAttribute(symbol, HoldAll) {
                return expr // Don't evaluate arguments
            }
            
            // Evaluate arguments
            args := make([]Expr, len(list.Elements)-1)
            for i := 1; i < len(list.Elements); i++ {
                if HasAttribute(symbol, HoldFirst) && i == 1 {
                    args[i-1] = list.Elements[i] // Don't evaluate first arg
                } else {
                    args[i-1] = evaluate(list.Elements[i])
                }
            }
            
            // Apply transformations based on attributes
            if HasAttribute(symbol, Flat) {
                args = flatten(symbol, args)
            }
            if HasAttribute(symbol, Orderless) {
                args = sort(args)
            }
            
            return NewList(append([]Expr{head}, args...)...)
        }
    }
    return expr
}
```

## Testing

The attribute system includes comprehensive tests:

```bash
# Run all attribute tests
go test -run "TestAttribute" -v

# Run integration tests
go test -run "TestAttributeSystemIntegration" -v

# Run benchmarks
go test -bench="BenchmarkAttributeOperations" -v
```

## Implementation Details

- **Global Symbol Table**: Attributes are stored in a global, thread-safe symbol table
- **Persistent Storage**: Attributes persist across expression evaluations
- **Memory Efficient**: Only symbols with attributes consume memory
- **Atomic Operations**: All operations are atomic and thread-safe
- **Sorted Output**: Attributes are always returned in sorted order for consistency

## Performance

The attribute system is designed for high performance:

- **O(1) attribute lookup** for checking if a symbol has an attribute
- **O(log n) attribute retrieval** where n is the number of attributes for a symbol
- **Thread-safe** with minimal contention using read-write locks
- **Memory efficient** with cleanup of empty attribute maps

## Future Extensions

The system is designed to be extensible:

- **Custom Attributes**: Easy to add new attribute types
- **Attribute Validation**: Can add validation logic for attribute combinations
- **Persistence**: Can be extended to save/load attributes from disk
- **Scoping**: Can be extended to support local attribute scoping