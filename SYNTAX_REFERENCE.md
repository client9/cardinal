# S-Expression Evaluator Syntax Reference

This document describes the syntax used by our s-expression evaluator and compares it with Mathematica where relevant.

## Function Call Syntax

### Our Syntax (S-Expressions)
```
FunctionName(arg1, arg2, arg3)
Plus(1, 2, 3)
Equal(x, y)
If(True, 1, 2)
MatchQ(expr, pattern)
```

### Mathematica Syntax (for comparison)
```
FunctionName[arg1, arg2, arg3]
Plus[1, 2, 3]
Equal[x, y]
If[True, 1, 2]
MatchQ[expr, pattern]
```

**Key Difference**: We use parentheses `()` for function calls, Mathematica uses square brackets `[]`.

## Data Types

### Atoms
| Type | Our Syntax | Mathematica | Examples |
|------|------------|-------------|----------|
| Integer | `42`, `-7`, `0` | Same | `42`, `-7`, `0` |
| Float | `3.14`, `-2.5`, `1.0` | Same | `3.14`, `-2.5`, `1.0` |
| String | `"hello"`, `"world"` | Same | `"hello"`, `"world"` |
| Boolean | `True`, `False` | Same | `True`, `False` |
| Symbol | `x`, `myVar`, `Plus` | Same | `x`, `myVar`, `Plus` |

### Lists
| Our Syntax | Mathematica | Description |
|------------|-------------|-------------|
| `List(1, 2, 3)` | `{1, 2, 3}` or `List[1, 2, 3]` | List of elements |
| `List()` | `{}` or `List[]` | Empty list |

**Note**: Mathematica has shorthand `{1, 2, 3}` for lists, we always use `List(...)`.

## Pattern Syntax

### Basic Patterns
| Pattern Type | Our Syntax | Mathematica | Description |
|--------------|------------|-------------|-------------|
| Blank | `_` | `_` | Matches any expression |
| Typed Blank | `_Integer`, `_String` | `_Integer`, `_String` | Matches specific type |
| Named Pattern | `x_`, `name_Integer` | `x_`, `name_Integer` | Matches and binds variable |

### Sequence Patterns
| Pattern Type | Our Syntax | Mathematica | Description |
|--------------|------------|-------------|-------------|
| Blank Sequence | `__` | `__` | Matches 1+ expressions |
| Typed Sequence | `__Integer` | `__Integer` | Matches 1+ integers |
| Named Sequence | `x__`, `nums__Integer` | `x__`, `nums__Integer` | Matches and binds sequence |
| Null Sequence | `___` | `___` | Matches 0+ expressions |
| Typed Null Seq | `___Integer` | `___Integer` | Matches 0+ integers |
| Named Null Seq | `x___`, `opts___` | `x___`, `opts___` | Matches and binds 0+ |

### Symbolic Patterns (Advanced)
| Our Syntax | Mathematica | Description |
|------------|-------------|-------------|
| `Blank()` | `Blank[]` | Symbolic blank pattern |
| `Blank(Integer)` | `Blank[Integer]` | Symbolic typed blank |
| `BlankSequence()` | `BlankSequence[]` | Symbolic sequence pattern |
| `BlankNullSequence()` | `BlankNullSequence[]` | Symbolic null sequence |
| `Pattern(x, Blank())` | `Pattern[x, Blank[]]` | Symbolic named pattern |

## Function Definitions

### Pattern-Based Functions
```lisp
; Our syntax
f(x_) := Plus(x, 1)
g(x_Integer, y_String) := List(x, y)
factorial(0) := 1
factorial(n_Integer) := Times(n, factorial(Subtract(n, 1)))

; Mathematica equivalent
f[x_] := x + 1
g[x_Integer, y_String] := {x, y}
factorial[0] := 1
factorial[n_Integer] := n * factorial[n - 1]
```

### Function Calls with Patterns
```lisp
; Our syntax
Plus(x_, y_)        ; Matches Plus(anything, anything)
List(_, _, _)       ; Matches List with exactly 3 elements
f(x_Integer)        ; Matches f with one integer argument
MyFunc(a_, b__, c_) ; Matches MyFunc with 1+variable args

; Mathematica equivalent
Plus[x_, y_]        ; Same semantics
{_, _, _}           ; Matches list with 3 elements
f[x_Integer]        ; Same semantics
MyFunc[a_, b__, c_] ; Same semantics
```

## Special Forms and Control Structures

### Assignment
| Operation | Our Syntax | Mathematica | Description |
|-----------|------------|-------------|-------------|
| Immediate | `Set(x, value)` | `x = value` | Evaluate and assign |
| Delayed | `SetDelayed(x, expr)` | `x := expr` | Assign unevaluated |
| Unset | `Unset(x)` | `x =.` | Remove assignment |

### Control Flow
| Construct | Our Syntax | Mathematica | Description |
|-----------|------------|-------------|-------------|
| Conditional | `If(test, then, else)` | `If[test, then, else]` | Conditional evaluation |
| Loop | `While(test, body)` | `While[test, body]` | While loop |
| Sequence | `CompoundExpression(a, b, c)` | `a; b; c` | Sequential evaluation |

### Logical Operations
| Operation | Our Syntax | Mathematica | Description |
|-----------|------------|-------------|-------------|
| And | `And(a, b, c)` | `a && b && c` | Logical AND |
| Or | `Or(a, b, c)` | `a \|\| b \|\| c` | Logical OR |
| Not | `Not(x)` | `!x` | Logical NOT |

## Arithmetic Operations

### Basic Arithmetic
| Operation | Our Syntax | Mathematica | Description |
|-----------|------------|-------------|-------------|
| Addition | `Plus(a, b, c)` | `a + b + c` | Addition |
| Multiplication | `Times(a, b, c)` | `a * b * c` | Multiplication |
| Subtraction | `Subtract(a, b)` | `a - b` | Subtraction |
| Division | `Divide(a, b)` | `a / b` | Division |
| Power | `Power(a, b)` | `a ^ b` | Exponentiation |

### Comparison
| Operation | Our Syntax | Mathematica | Description |
|-----------|------------|-------------|-------------|
| Equal | `Equal(a, b)` | `a == b` | Equality test |
| Less | `Less(a, b)` | `a < b` | Less than |
| Greater | `Greater(a, b)` | `a > b` | Greater than |
| Less Equal | `LessEqual(a, b)` | `a <= b` | Less or equal |
| Greater Equal | `GreaterEqual(a, b)` | `a >= b` | Greater or equal |
| Same | `SameQ(a, b)` | `a === b` | Identity test |

## Built-in Functions

### Type Testing
| Function | Our Syntax | Mathematica | Description |
|----------|------------|-------------|-------------|
| Integer test | `IntegerQ(x)` | `IntegerQ[x]` | Test if integer |
| Number test | `NumberQ(x)` | `NumberQ[x]` | Test if number |
| String test | `StringQ(x)` | `StringQ[x]` | Test if string |
| List test | `ListQ(x)` | `ListQ[x]` | Test if list |
| Atom test | `AtomQ(x)` | `AtomQ[x]` | Test if atom |

### Structure Functions
| Function | Our Syntax | Mathematica | Description |
|----------|------------|-------------|-------------|
| Head | `Head(expr)` | `Head[expr]` | Get expression head |
| Length | `Length(expr)` | `Length[expr]` | Get expression length |
| First | `First(expr)` | `First[expr]` | Get first element |
| Last | `Last(expr)` | `Last[expr]` | Get last element |

### Pattern Matching
| Function | Our Syntax | Mathematica | Description |
|----------|------------|-------------|-------------|
| Match test | `MatchQ(expr, pattern)` | `MatchQ[expr, pattern]` | Test pattern match |

## Attributes

### Function Attributes
| Attribute | Our Syntax | Mathematica | Description |
|-----------|------------|-------------|-------------|
| Flat | `Flat` | `Flat` | Associative function |
| Orderless | `Orderless` | `Orderless` | Commutative function |
| Hold First | `HoldFirst` | `HoldFirst` | Don't evaluate 1st arg |
| Hold All | `HoldAll` | `HoldAll` | Don't evaluate any args |
| One Identity | `OneIdentity` | `OneIdentity` | f(x) simplifies to x |

### Setting Attributes
```lisp
; Our syntax
SetAttributes(Plus, List(Flat, Orderless))
SetAttributes(MatchQ, List(HoldFirst))

; Mathematica syntax  
SetAttributes[Plus, {Flat, Orderless}]
SetAttributes[MatchQ, HoldFirst]
```

## Examples

### Pattern Matching Examples
```lisp
; Our syntax
MatchQ(42, _)                    ; True
MatchQ(42, _Integer)             ; True  
MatchQ(List(1, 2, 3), List(_, _, _))  ; True
MatchQ(Plus(1, 2), Plus(_, _))   ; False (Plus(1,2) → 3)
MatchQ(f(1, 2), f(_, _))         ; True

; Mathematica equivalent
MatchQ[42, _]                    (* True *)
MatchQ[42, _Integer]             (* True *)
MatchQ[{1, 2, 3}, {_, _, _}]     (* True *)
MatchQ[1 + 2, Plus[_, _]]        (* False *)
MatchQ[f[1, 2], f[_, _]]         (* True *)
```

### Function Definition Examples
```lisp
; Our syntax - pattern-based function
double(x_) := Times(2, x)
myMax() := 0
myMax(x_) := x  
myMax(x_, y_) := If(Greater(x, y), x, y)

; Mathematica equivalent
double[x_] := 2 * x
myMax[] := 0
myMax[x_] := x
myMax[x_, y_] := If[x > y, x, y]
```

## Key Differences Summary

1. **Function Calls**: We use `()`, Mathematica uses `[]`
2. **Lists**: We use `List(...)`, Mathematica has `{...}` shorthand
3. **Operators**: We use function calls `Plus(a, b)`, Mathematica has infix `a + b`
4. **Assignment**: We use `Set(x, y)`, Mathematica uses `x = y`
5. **Patterns**: Identical syntax and semantics
6. **Attributes**: Same concepts, slightly different syntax for setting

## Notes

- Our evaluator follows Mathematica's evaluation semantics but uses Lisp-style syntax
- All pattern matching behavior is designed to match Mathematica exactly
- Function attributes work the same way as in Mathematica
- The symbolic pattern system allows introspection just like Mathematica