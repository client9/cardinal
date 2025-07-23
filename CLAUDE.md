# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Build
```bash
make build    # Runs wrapgen code generation, builds main package and REPL
```

The build process uses reflection-based code generation via `cmd/wrapgen/` to create wrapper functions before compilation.

### Testing
```bash
make test     # Runs build then executes all tests
go test .     # Run tests directly without build
```

### Linting
```bash
make lint     # Runs go mod tidy, gofmt, and golangci-lint
```

### Clean
```bash
make clean    # Removes generated files and binaries
```

### Run REPL
```bash
./repl                    # Interactive mode
./repl -file examples.sexpr   # Execute file
```

## Architecture Overview

This is a Go implementation of an s-expression parser and evaluator with Mathematica-inspired semantics but distinct syntax using parentheses for function calls and modern association syntax.

### Core Components

**Context-Based Architecture**: The system uses isolated `Context` instances instead of global state, providing:
- **Evaluator**: Each instance has its own `Context` containing variable bindings and symbol table
- **SymbolTable**: Thread-safe attribute management (HoldAll, Flat, Orderless, OneIdentity, etc.)
- **Instance Isolation**: Multiple evaluators can run concurrently without interference

**Key Types** (in `core/` package):
- `Expr`: Base interface for all expressions (atoms, lists, errors, objects)
- `Atom`: Numbers, strings, booleans, symbols
- `List`: Nested s-expressions
- `Context`: Evaluation environment with variable/attribute storage

**Expression Evaluation**:
- Attribute-aware processing respects symbol attributes during evaluation
- Recursion tracking with stack depth limits
- Support for immediate (`Set`) and delayed (`SetDelayed`) assignments

### Code Generation System

The `cmd/wrapgen/` tool uses reflection to automatically generate wrapper functions:
- Analyzes Go functions in `stdlib/` packages
- Generates type-safe wrappers that convert between Go types and `Expr`
- Creates `builtin_setup.go` with attribute configurations
- Run via `make build` before compilation

### Package Structure

- `core/`: Core expression types and interfaces
- `stdlib/`: Built-in function implementations (math, logic, lists, strings, etc.)
- `cmd/repl/`: Interactive REPL implementation
- `cmd/wrapgen/`: Code generation for stdlib wrappers
- Root: Parser, lexer, evaluator, pattern matching, attributes

### Attribute System

Symbols have Mathematica-style attributes affecting evaluation:
- **HoldAll/HoldFirst/HoldRest**: Prevent argument evaluation
- **Flat**: Enable associative flattening (Plus(1, Plus(2, 3)) → Plus(1, 2, 3))
- **Orderless**: Sort arguments (Plus(3, 1, 2) → Plus(1, 2, 3))
- **OneIdentity**: Single arguments unwrap (Plus(42) → 42)

Use `attributes` command in REPL to view all symbol attributes.

## Syntax Key Differences from Mathematica

**Function Calls**: Use parentheses `Plus(1, 2, 3)` instead of brackets `Plus[1, 2, 3]`

**Lists**: Support both `List(1, 2, 3)` and shorthand `[1, 2, 3]` (Mathematica uses `{1, 2, 3}`)

**Associations**: Modern `{key: value, age: 30}` syntax instead of `<|key -> value, age -> 30|>`

**Rules**: Single colon `a:b` creates `Rule(a,b)` (Mathematica uses `a->b`)

**Output Formats**:
- **FullForm**: Verbose symbolic `Plus(1, 2)` for debugging
- **InputForm**: User-friendly `1 + 2` with infix operators and precedence

**Multi-line Support**: Files and REPL support multi-line expressions with automatic completion detection

**Numeric Computation**: Distinguishes integers (`int64`) from reals (`float64`) with immediate machine-precision arithmetic, unlike Mathematica's arbitrary-precision symbolic computation

## Testing Practices

- Each test creates isolated `Evaluator` instances to avoid interference
- Tests cover parser, evaluator, pattern matching, attributes, and REPL functionality
- Use `NewEvaluator()` for clean test contexts
- Pattern matching tests verify specificity ordering and complex nested patterns