# S-Expression Evaluator Architecture

## Overview

The s-expression evaluator system is designed with a context-based architecture that provides isolation between different evaluator instances while supporting efficient attribute management and evaluation.

## Core Components

### 1. Context (`Context`)

The `Context` struct serves as the evaluation environment, containing:
- **Variable bindings**: Local variable assignments (`x = 5`, `y := 2*x`)
- **Symbol table**: Attribute definitions for symbols
- **Parent context**: Hierarchical scoping for nested evaluations

```go
type Context struct {
    variables   map[string]Expr
    parent      *Context
    symbolTable *SymbolTable
}
```

### 2. Symbol Table (`SymbolTable`)

The `SymbolTable` manages symbol attributes in a thread-safe manner:
- **Attribute storage**: Maps symbols to their attributes (HoldAll, Flat, etc.)
- **Thread safety**: Uses RWMutex for concurrent access
- **Isolation**: Each context has its own symbol table instance

```go
type SymbolTable struct {
    attributes map[string]map[Attribute]bool
    mu         sync.RWMutex
}
```

### 3. Evaluator (`Evaluator`)

The `Evaluator` performs expression evaluation using context-aware attribute checking:
- **Context-based evaluation**: All operations respect the context's symbol table
- **Attribute-aware processing**: Handles Hold, Flat, Orderless, OneIdentity attributes
- **Isolation**: Each evaluator has its own context and symbol table

## Key Benefits

### 1. **Instance Isolation**

Different evaluator instances are completely isolated:

```go
eval1 := NewEvaluator()
eval2 := NewEvaluator()

// These won't interfere with each other
eval1.context.symbolTable.SetAttributes("MyFunc", []Attribute{Flat})
eval2.context.symbolTable.SetAttributes("MyFunc", []Attribute{HoldAll})
```

### 2. **Concurrent Safety**

Multiple evaluators can run safely in parallel:

```go
// Safe to use concurrently
repl1 := NewREPL()  // Has its own evaluator/context
repl2 := NewREPL()  // Completely independent

// These can run in parallel without interference
go func() { repl1.EvaluateString("x = 10") }()
go func() { repl2.EvaluateString("x = 20") }()
```

### 3. **Hierarchical Contexts**

Child contexts inherit symbol tables but maintain variable isolation:

```go
parentCtx := NewContext()
childCtx := NewChildContext(parentCtx)  // Shares symbol table

// Child sees parent's attributes
parentCtx.symbolTable.SetAttributes("Plus", []Attribute{Flat})
// childCtx can see Plus has Flat attribute

// But variables are isolated
parentCtx.Set("x", NewIntAtom(10))
childCtx.Set("x", NewIntAtom(20))  // Independent variable
```

## Migration from Global State

Previously, the system used a global symbol table which caused:
- **Interference**: Multiple instances affected each other
- **Concurrency issues**: Required global locking
- **Testing complexity**: Tests had to reset global state

The new context-based approach eliminates these issues:

### Before (Global)
```go
// Global state - problems with isolation
SetAttributes("Plus", []Attribute{Flat})
eval1 := NewEvaluator()  // Uses global table
eval2 := NewEvaluator()  // Same global table - interference!
```

### After (Context-based)
```go
// Instance isolation - no interference
eval1 := NewEvaluator()  // Has own symbol table
eval2 := NewEvaluator()  // Completely independent
```

## Usage Patterns

### 1. **REPL Usage**
```go
repl := NewREPL()  // Creates evaluator with setup attributes
result, _ := repl.EvaluateString("Plus[1, 2, 3]")  // Uses instance context
```

### 2. **Custom Evaluator**
```go
eval := NewEvaluator()
setupBuiltinAttributes(eval.context.symbolTable)  // Setup standard attributes
result := eval.Evaluate(expr)
```

### 3. **Testing Isolation**
```go
func TestFeature(t *testing.T) {
    eval := NewEvaluator()  // Clean instance
    eval.context.symbolTable.SetAttributes("TestFunc", []Attribute{Flat})
    // Test without affecting other tests
}
```

## Performance Characteristics

- **Memory**: Each context has its own symbol table (small overhead)
- **Concurrency**: No global locks, excellent parallel performance
- **Isolation**: Complete separation between instances
- **Inheritance**: Efficient sharing of symbol tables in child contexts

## Thread Safety

- **SymbolTable**: Thread-safe with RWMutex
- **Context**: Safe for single-thread use, parent sharing is safe
- **Evaluator**: Safe for single-thread use
- **Multiple instances**: Fully concurrent safe