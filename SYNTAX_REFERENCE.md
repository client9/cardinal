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
| `[1, 2, 3]` | `{1, 2, 3}` | List literal shorthand |
| `List()` or `[]` | `{}` or `List[]` | Empty list |

### Associations (Modern Map/Dict Syntax)
| Our Syntax | Mathematica | Description |
|------------|-------------|-------------|
| `{key: value}` | `<\|key -> value\|>` | Association with one pair |
| `{name: "Bob", age: 30}` | `<\|"name" -> "Bob", "age" -> 30\|>` | Multiple key-value pairs |
| `{}` | `<\|\|>` | Empty association |

**Note**: We use modern `{key: value}` syntax like JavaScript/JSON, while Mathematica uses `<|key -> value|>`.

## Multi-line Expressions

Our evaluator supports multi-line expressions both in files and interactive REPL:

### File Support
```javascript
# Comments start with #
{
  person: {
    name: "Bob",
    age: 30
  },
  location: {
    city: "New York", 
    state: "NY"
  }
}

# Multi-line function calls  
Plus(
  1,
  2,
  3
)
```

### Interactive REPL
```
sexpr> {
   ... name: "Bob",
   ... age: 30
   ... }
{name: "Bob", age: 30}

sexpr> Plus(
   ... 1,
   ... 2
   ... )
3
```

The evaluator automatically detects when expressions are complete and evaluates them.

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
| Part | `Part(expr, index)` | `expr[[index]]` | Access element by index/key |

### Association Functions
| Function | Our Syntax | Mathematica | Description |
|----------|------------|-------------|-------------|
| Association test | `AssociationQ(x)` | `AssociationQ[x]` | Test if association |
| Get keys | `Keys(assoc)` | `Keys[assoc]` | Get list of keys |
| Get values | `Values(assoc)` | `Values[assoc]` | Get list of values |
| Access value | `Part(assoc, key)` | `assoc[key]` | Get value by key |

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

### Association Examples
```javascript
// Our syntax - modern association syntax
{name: "Bob", age: 30, active: True}
Part({name: "Bob", age: 30}, name)     // "Bob"
Keys({name: "Bob", age: 30})           // List(name, age)
Values({name: "Bob", age: 30})         // List("Bob", 30)
Length({name: "Bob", age: 30})         // 2

// Nested associations
{
  person: {name: "Bob", age: 30},
  scores: [85, 92, 78]
}

// Mathematica equivalent
<|"name" -> "Bob", "age" -> 30, "active" -> True|>
assoc["name"]                          (* "Bob" *)
Keys[assoc]                            (* {"name", "age"} *)
Values[assoc]                          (* {"Bob", 30} *)
Length[assoc]                          (* 2 *)
```

## InputForm vs FullForm

Our evaluator supports two output formats:

### FullForm (String() method)
- Verbose symbolic representation: `Plus(1, 2)`, `List(1, 2, 3)`
- Always unambiguous and complete
- Used for internal representation and debugging

### InputForm (InputForm() method)
- Compact, user-friendly representation with infix operators and shortcuts
- Supports operator precedence and automatic parenthesization
- Resembles common mathematical notation

| Expression | FullForm | InputForm |
|------------|----------|-----------|
| Plus(1, 2) | `Plus(1, 2)` | `1 + 2` |
| Times(3, 4) | `Times(3, 4)` | `3 * 4` |
| List(1, 2, 3) | `List(1, 2, 3)` | `[1, 2, 3]` |
| Set(x, 5) | `Set(x, 5)` | `x = 5` |
| SetDelayed(f, Plus(x, 1)) | `SetDelayed(f, Plus(x, 1))` | `f := x + 1` |
| Equal(a, b) | `Equal(a, b)` | `a == b` |
| And(True, False) | `And(True, False)` | `True && False` |
| Association(Rule(a, b)) | `Association(Rule(a, b))` | `{a: b}` |
| Plus(1, Times(2, 3)) | `Plus(1, Times(2, 3))` | `1 + 2 * 3` |
| Times(Plus(1, 2), 3) | `Times(Plus(1, 2), 3)` | `(1 + 2) * 3` |

### Precedence and Parenthesization

InputForm automatically adds parentheses based on operator precedence:

1. **Assignment**: `=`, `:=` (lowest precedence)
2. **Logical OR**: `||`
3. **Logical AND**: `&&`
4. **Equality**: `==`, `!=`, `===`, `=!=`
5. **Comparison**: `<`, `>`, `<=`, `>=`
6. **Addition**: `+`, `-`
7. **Multiplication**: `*`, `/` (highest precedence)

Examples:
- `Plus(1, Times(2, 3))` → `1 + 2 * 3` (no parentheses needed)
- `Times(Plus(1, 2), 3)` → `(1 + 2) * 3` (parentheses added)
- `Set(x, Plus(1, 2))` → `x = 1 + 2` (assignment has lowest precedence)

## Key Differences Summary

1. **Function Calls**: We use `()`, Mathematica uses `[]`
2. **Lists**: We support both `List(...)` and `[...]` shorthand, Mathematica uses `{...}`
3. **Associations**: We use modern `{key: value}` syntax, Mathematica uses `<|key -> value|>`
4. **Operators**: We support both function calls `Plus(a, b)` and infix `a + b` (InputForm)
5. **Assignment**: We support both `Set(x, y)` and `x = y` (InputForm)
6. **Multi-line**: We support multi-line expressions with automatic completion detection
7. **Patterns**: Identical syntax and semantics
8. **Attributes**: Same concepts, slightly different syntax for setting
9. **Output Formats**: FullForm for symbolic representation, InputForm for user-friendly display

## Notes

- Our evaluator follows Mathematica's evaluation semantics but uses Lisp-style syntax
- All pattern matching behavior is designed to match Mathematica exactly
- Function attributes work the same way as in Mathematica
- The symbolic pattern system allows introspection just like Mathematica