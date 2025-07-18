# S-Expression Evaluator

A Go implementation of an s-expression parser and evaluator with Mathematica-style semantics.

## Features

- **Complete S-Expression Parser** - Supports atoms (numbers, strings, booleans, symbols) and nested lists
- **Context-Based Architecture** - Isolated evaluator instances with no global state interference
- **Attribute System** - Full support for symbol attributes (HoldAll, Flat, Orderless, OneIdentity, etc.)
- **Thread-Safe Design** - Multiple evaluators can run concurrently without interference
- **Expression Evaluator** - Evaluates expressions with support for:
  - Arithmetic operations (Plus, Times, Subtract, Divide, Power)
  - Comparison operations (Equal, Less, Greater, etc.)
  - Logical operations (And, Or, Not)
  - Control structures (If, Hold, Evaluate)
  - Variable assignment (Set, SetDelayed, Unset)
  - Built-in constants (Pi, E, True, False)
- **Interactive REPL** - Read-Eval-Print Loop for interactive computation
- **Infix Notation Support** - Standard mathematical operators (+, -, *, /, ==, <, >, etc.)

## Installation

```bash
go get github.com/client9/sexpr
```

## Usage

### As a Library

```go
package main

import (
    "fmt"
    "github.com/client9/sexpr"
)

func main() {
    // Create a REPL instance
    repl := sexpr.NewREPL()
    
    // Evaluate expressions
    result, _ := repl.EvaluateString("1 + 2 * 3")
    fmt.Println("1 + 2 * 3 =", result) // Output: 7
    
    result, _ = repl.EvaluateString("x = 5")
    fmt.Println("x =", result) // Output: 5
    
    result, _ = repl.EvaluateString("If[x > 3, \"big\", \"small\"]")
    fmt.Println("Result:", result) // Output: "big"
}
```

### Interactive REPL

Build and run the interactive REPL:

```bash
go build -o repl ./cmd/repl
./repl
```

Example session:
```
S-Expression REPL v1.0
Type 'quit' or 'exit' to exit, 'help' for help

sexpr> 1 + 2 * 3
7
sexpr> x = 10
10
sexpr> y = 20
20
sexpr> And[x > 5, y < 25]
True
sexpr> If[x > y, "x wins", "y wins"]
"y wins"
sexpr> Plus[1, 2, 3, 4, 5]
15
sexpr> Hold[1 + 2]
Hold[Plus[1, 2]]
sexpr> Pi
3.141592653589793
sexpr> quit
Goodbye!
```

### File Execution

Execute expressions from a file:

```bash
./repl -file examples.sexpr
```

### REPL Commands

- `help` - Show help information
- `quit` or `exit` - Exit the REPL
- `clear` - Clear all variable assignments
- `attributes` - Show all symbols with their attributes

## Expression Syntax

### Function Calls
```
Plus[1, 2, 3]
Times[2, x, 4]
Greater[5, 3]
```

### Infix Notation
```
1 + 2 * 3        # Arithmetic
x == 5           # Equality
x > y && y < 10  # Logical operations
```

### Variable Assignment
```
x = 5            # Immediate assignment
y := 2 * x       # Delayed assignment
```

### Control Structures
```
If[condition, then, else]
Hold[expression]
Evaluate[expression]
```

### Built-in Functions

#### Arithmetic
- `Plus[...]` - Addition
- `Times[...]` - Multiplication  
- `Subtract[a, b]` - Subtraction
- `Divide[a, b]` - Division
- `Power[base, exp]` - Exponentiation

#### Comparison
- `Equal[a, b]` - Equality test
- `Less[a, b]` - Less than
- `Greater[a, b]` - Greater than
- `LessEqual[a, b]` - Less than or equal
- `GreaterEqual[a, b]` - Greater than or equal
- `SameQ[a, b]` - Identity test

#### Logical
- `And[...]` - Logical AND
- `Or[...]` - Logical OR
- `Not[x]` - Logical NOT

#### Control
- `If[cond, then, else]` - Conditional expression
- `Hold[expr]` - Prevent evaluation
- `Evaluate[expr]` - Force evaluation

#### Assignment
- `Set[var, value]` - Immediate assignment
- `SetDelayed[var, value]` - Delayed assignment
- `Unset[var]` - Remove variable

### Attributes

The system supports Mathematica-style attributes that control evaluation:

- **HoldAll** - Prevent evaluation of all arguments
- **HoldFirst** - Prevent evaluation of first argument
- **HoldRest** - Prevent evaluation of all but first argument
- **Flat** - Flatten nested applications (associativity)
- **Orderless** - Sort arguments (commutativity)
- **OneIdentity** - f[x] → x for single arguments

Example:
```
sexpr> attributes
Symbols with attributes:
===============================
And            : {Flat, HoldAll, Orderless}
Block          : {HoldAll}
CompoundExpression: {HoldAll}
E              : {Constant, Protected}
False          : {Constant, Protected}
Hold           : {HoldAll}
If             : {HoldRest}
Module         : {HoldAll}
Or             : {Flat, HoldAll, Orderless}
Pi             : {Constant, Protected}
Plus           : {Flat, OneIdentity, Orderless}
Power          : {OneIdentity}
Set            : {HoldFirst}
SetDelayed     : {HoldAll}
Times          : {Flat, OneIdentity, Orderless}
True           : {Constant, Protected}
Unset          : {HoldFirst}
While          : {HoldAll}
```

## Examples

See `examples.sexpr` for more examples:

```bash
# Arithmetic with precedence
1 + 2 * 3                    # → 7

# Function calls
Plus[1, 2, 3, 4, 5]          # → 15

# Variables and assignments
x = 10                       # → 10
y = 20                       # → 20
z = x + y                    # → 30

# Comparisons
x > y                        # → False
Equal[x, 10]                 # → True

# Logical operations
And[True, False]             # → False
Or[False, True]              # → True
And[x > 5, y < 25]          # → True

# Conditionals
If[x > y, "x is greater", "y is greater"]  # → "y is greater"

# Mathematical constants
Pi                           # → 3.141592653589793
E                            # → 2.718281828459045

# Attribute demonstrations
Plus[1, Plus[2, 3]]          # → 6 (Flat attribute)
Plus[3, 1, 2]                # → 6 (Orderless attribute)
Plus[42]                     # → 42 (OneIdentity attribute)

# Hold expressions
Hold[1 + 2]                  # → Hold[Plus[1, 2]]
```

## Testing

Run the test suite:

```bash
go test -v
```

## License

This project is open source. See the source code for details.