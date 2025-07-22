# Built-in Functions Reference

This document lists all built-in functions available in the s-expression evaluator with their signatures and descriptions.

## Arithmetic Operations

### Plus(x_, y_, ...)
**Description**: Addition of numbers  
**Examples**: `Plus(1, 2, 3)` → `6`

### Times(x_, y_, ...)
**Description**: Multiplication of numbers  
**Examples**: `Times(2, 3, 4)` → `24`

### Subtract(x_, y_)
**Description**: Subtraction (x - y)  
**Examples**: `Subtract(10, 3)` → `7`

### Divide(x_, y_)
**Description**: Division (x / y)  
**Examples**: `Divide(15, 3)` → `5`

### Power(x_, y_)
**Description**: Exponentiation (x^y)  
**Examples**: `Power(2, 3)` → `8`

### Mod(x_, y_)
**Description**: Modulo operation (x % y)  
**Examples**: `Mod(10, 3)` → `1`

### Abs(x_)
**Description**: Absolute value  
**Examples**: `Abs(-5)` → `5`

### Min(x_, y_, ...)
**Description**: Minimum of values  
**Examples**: `Min(3, 1, 4)` → `1`

### Max(x_, y_, ...)
**Description**: Maximum of values  
**Examples**: `Max(3, 1, 4)` → `4`

## Comparison Operations

### Equal(x_, y_)
**Description**: Equality test (==)  
**Examples**: `Equal(5, 5)` → `True`

### Unequal(x_, y_)
**Description**: Inequality test (!=)  
**Examples**: `Unequal(5, 3)` → `True`

### Less(x_, y_)
**Description**: Less than test (<)  
**Examples**: `Less(3, 5)` → `True`

### LessEqual(x_, y_)
**Description**: Less than or equal test (<=)  
**Examples**: `LessEqual(3, 3)` → `True`

### Greater(x_, y_)
**Description**: Greater than test (>)  
**Examples**: `Greater(5, 3)` → `True`

### GreaterEqual(x_, y_)
**Description**: Greater than or equal test (>=)  
**Examples**: `GreaterEqual(5, 5)` → `True`

### SameQ(x_, y_)
**Description**: Identity test (===) - stricter than Equal  
**Examples**: `SameQ(3, 3)` → `True`

### UnsameQ(x_, y_)
**Description**: Non-identity test (!==)  
**Examples**: `UnsameQ(3, 5)` → `True`

## Logical Operations

### And(x_, y_, ...)
**Description**: Logical AND with short-circuit evaluation  
**Attributes**: Flat, Orderless, HoldAll  
**Examples**: `And(True, True)` → `True`

### Or(x_, y_, ...)
**Description**: Logical OR with short-circuit evaluation  
**Attributes**: Flat, Orderless, HoldAll  
**Examples**: `Or(False, True)` → `True`

### Not(x_)
**Description**: Logical NOT  
**Examples**: `Not(True)` → `False`

## Control Flow

### If(test_, then_, else_)
**Description**: Conditional expression  
**Attributes**: HoldRest  
**Examples**: `If(Greater(5, 3), "big", "small")` → `"big"`

### While(test_, body_)
**Description**: While loop  
**Attributes**: HoldAll  
**Examples**: `While(Less(x, 10), Set(x, Plus(x, 1)))`

### CompoundExpression(expr1_, expr2_, ...)
**Description**: Sequential evaluation, returns last result  
**Examples**: `CompoundExpression(Set(x, 5), Plus(x, 1))` → `6`

## Assignment Operations

### Set(symbol_, value_)
**Description**: Immediate assignment (=)  
**Attributes**: HoldFirst  
**Examples**: `Set(x, 5)` → `5`

### SetDelayed(symbol_, value_)
**Description**: Delayed assignment (:=)  
**Attributes**: HoldAll  
**Examples**: `SetDelayed(f(x_), Plus(x, 1))`

### Unset(symbol_)
**Description**: Remove assignment  
**Attributes**: HoldFirst  
**Examples**: `Unset(x)`

## Type Testing Functions

### IntegerQ(x_)
**Description**: Test if expression is an integer  
**Examples**: `IntegerQ(42)` → `True`

### FloatQ(x_)
**Description**: Test if expression is a float  
**Examples**: `FloatQ(3.14)` → `True`

### NumberQ(x_)
**Description**: Test if expression is a number (integer or float)  
**Examples**: `NumberQ(42)` → `True`

### StringQ(x_)
**Description**: Test if expression is a string  
**Examples**: `StringQ("hello")` → `True`

### BooleanQ(x_)
**Description**: Test if expression is a boolean  
**Examples**: `BooleanQ(True)` → `True`

### SymbolQ(x_)
**Description**: Test if expression is a symbol  
**Examples**: `SymbolQ(x)` → `True`

### ListQ(x_)
**Description**: Test if expression is a list  
**Examples**: `ListQ(List(1, 2, 3))` → `True`

### AtomQ(x_)
**Description**: Test if expression is an atom (not a list)  
**Examples**: `AtomQ(42)` → `True`

## Structure Functions

### Head(expr_)
**Description**: Get the head/type of an expression  
**Examples**: `Head(42)` → `Integer`, `Head(Plus(1, 2))` → `Integer` (after evaluation)

### Length(expr_)
**Description**: Get the length of a list or association  
**Examples**: 
- `Length(List(1, 2, 3))` → `3`
- `Length({name: "Bob", age: 30})` → `2`

### First(expr_)
**Description**: Get the first element of a list  
**Examples**: `First(List(1, 2, 3))` → `1`

### Last(expr_)
**Description**: Get the last element of a list  
**Examples**: `Last(List(1, 2, 3))` → `3`

### Rest(expr_)
**Description**: Get all elements except the first  
**Examples**: `Rest(List(1, 2, 3))` → `List(2, 3)`

### Most(expr_)
**Description**: Get all elements except the last  
**Examples**: `Most(List(1, 2, 3))` → `List(1, 2)`

### Part(expr_, index_)
**Description**: Get element by index (lists) or key (associations)
- For lists: 1-indexed access, supports negative indices  
- For associations: key-based access  
**Examples**: 
- `Part(List(1, 2, 3), 2)` → `2`
- `Part({name: "Bob", age: 30}, name)` → `"Bob"`

## Association Functions

### Association(rules_...)
**Description**: Create an association from Rule expressions  
**Examples**: `Association(Rule(name, "Bob"), Rule(age, 30))` → `{name: "Bob", age: 30}`
**Note**: Typically created using `{key: value}` syntax

### AssociationQ(expr_)
**Description**: Test if expression is an association  
**Examples**: `AssociationQ({name: "Bob"})` → `True`, `AssociationQ(List(1, 2))` → `False`

### Keys(assoc_)
**Description**: Get all keys from an association as a list  
**Examples**: `Keys({name: "Bob", age: 30})` → `List(name, age)`

### Values(assoc_)
**Description**: Get all values from an association as a list  
**Examples**: `Values({name: "Bob", age: 30})` → `List("Bob", 30)`

**Note**: Keys and values are returned in insertion order.

## Pattern Matching

### MatchQ(expr_, pattern_)
**Description**: Test if expression matches pattern (evaluates expr first)  
**Examples**: `MatchQ(42, _Integer)` → `True`, `MatchQ(Plus(1, 2), _Integer)` → `True`

## List Functions

### List(elem1_, elem2_, ...)
**Description**: Create a list  
**Examples**: `List(1, 2, 3)` → `List(1, 2, 3)`

### Append(list_, elem_)
**Description**: Add element to end of list  
**Examples**: `Append(List(1, 2), 3)` → `List(1, 2, 3)`

### Prepend(list_, elem_)
**Description**: Add element to beginning of list  
**Examples**: `Prepend(List(2, 3), 1)` → `List(1, 2, 3)`

## Mathematical Constants

### Pi
**Description**: The mathematical constant π  
**Examples**: `Pi` → `3.141592653589793`

### E
**Description**: Euler's number e  
**Examples**: `E` → `2.718281828459045`

## Mathematical Functions

### Sin(x_)
**Description**: Sine function  
**Examples**: `Sin(0)` → `0`

### Cos(x_)
**Description**: Cosine function  
**Examples**: `Cos(0)` → `1`

### Tan(x_)
**Description**: Tangent function  
**Examples**: `Tan(0)` → `0`

### Log(x_)
**Description**: Natural logarithm  
**Examples**: `Log(E)` → `1`

### Exp(x_)
**Description**: Exponential function (e^x)  
**Examples**: `Exp(1)` → `2.718281828459045`

### Sqrt(x_)
**Description**: Square root  
**Examples**: `Sqrt(4)` → `2`

## Evaluation Control

### Hold(expr_)
**Description**: Prevent evaluation of expression  
**Attributes**: HoldAll  
**Examples**: `Hold(Plus(1, 2))` → `Hold(Plus(1, 2))`

### Evaluate(expr_)
**Description**: Force evaluation of expression  
**Examples**: `Evaluate(Hold(Plus(1, 2)))` → `3`

## Symbolic Pattern Functions

### Blank()
**Description**: Create a blank pattern (_)  
**Examples**: `Blank()` → `Blank()`, `Blank(Integer)` → `Blank(Integer)`

### BlankSequence()
**Description**: Create a blank sequence pattern (__)  
**Examples**: `BlankSequence()` → `BlankSequence()`

### BlankNullSequence()
**Description**: Create a blank null sequence pattern (___)  
**Examples**: `BlankNullSequence()` → `BlankNullSequence()`

### Pattern(name_, blank_)
**Description**: Create a named pattern  
**Examples**: `Pattern(x, Blank())` → `Pattern(x, Blank())`

## Function Attributes

Functions can have attributes that modify their behavior:

- **Flat**: Associative function (e.g., Plus, Times)
- **Orderless**: Commutative function (e.g., Plus, Times)  
- **HoldFirst**: Don't evaluate first argument
- **HoldRest**: Don't evaluate arguments after the first
- **HoldAll**: Don't evaluate any arguments
- **OneIdentity**: Single-argument form simplifies (e.g., Plus(x) → x)
- **Listable**: Function can be applied to lists element-wise
- **NumericFunction**: Function expects numeric arguments
- **Constant**: Symbol represents a constant value
- **Protected**: Symbol is protected from modification
- **ReadProtected**: Symbol cannot be read
- **Locked**: Symbol definition is locked
- **Temporary**: Symbol is temporary

### SetAttributes(symbol_, attributes_)
**Description**: Set one or more attributes for a function/symbol  
**Attributes**: HoldFirst  
**Examples**: 
- `SetAttributes(f, Protected)` → `Null`
- `SetAttributes(f, List(Flat, Orderless))` → `Null`

### ClearAttributes(symbol_, attributes_)
**Description**: Clear specific attributes from a function/symbol  
**Attributes**: HoldFirst  
**Examples**: 
- `ClearAttributes(f, Protected)` → `Null`
- `ClearAttributes(f, List(Flat, Orderless))` → `Null`

### ClearAttributes(symbol_)
**Description**: Clear all attributes from a function/symbol  
**Attributes**: HoldFirst  
**Examples**: `ClearAttributes(f)` → `Null`

### Attributes(symbol_)
**Description**: Get all attributes of a function/symbol (sorted alphabetically)  
**Attributes**: HoldFirst  
**Examples**: 
- `Attributes(Plus)` → `List(Flat, Listable, NumericFunction, OneIdentity, Orderless, Protected)`
- `Attributes(newFunc)` → `List()`

## Error Handling

Functions automatically propagate errors - if any argument is an error, the error is returned without evaluation.

## Output Formats

Our evaluator provides two output formats:

### FullForm (Default)
- Complete symbolic representation
- Example: `Plus(1, 2)`, `List(1, 2, 3)`, `Set(x, 5)`

### InputForm 
- User-friendly infix notation with operator precedence
- Example: `1 + 2`, `[1, 2, 3]`, `x = 5`

| Function | FullForm | InputForm |
|----------|----------|-----------|
| Plus(1, 2) | `Plus(1, 2)` | `1 + 2` |
| Times(a, b) | `Times(a, b)` | `a * b` |
| Set(x, 5) | `Set(x, 5)` | `x = 5` |
| SetDelayed(f, Plus(x, 1)) | `SetDelayed(f, Plus(x, 1))` | `f := x + 1` |
| Equal(a, b) | `Equal(a, b)` | `a == b` |
| And(True, False) | `And(True, False)` | `True && False` |
| List(1, 2, 3) | `List(1, 2, 3)` | `[1, 2, 3]` |
| Association(Rule(a, b)) | `Association(Rule(a, b))` | `{a: b}` |

## Usage Notes

1. **Function Evaluation**: All functions evaluate their arguments unless they have Hold attributes
2. **Pattern Matching**: Use underscore patterns (_) for flexible function definitions
3. **Attributes**: Function attributes control evaluation and algebraic properties
4. **Error Propagation**: Errors automatically bubble up through function calls
5. **Symbolic Computation**: Functions work with both numeric and symbolic expressions
6. **Output Formats**: Use InputForm for readable output, FullForm for debugging

## Examples

```lisp
# Basic arithmetic
Plus(1, 2, 3)                    # → 6
Times(2, 3)                      # → 6
Power(2, 3)                      # → 8

# Comparisons and logic
Greater(5, 3)                    # → True
And(True, Greater(5, 3))         # → True

# Type testing
NumberQ(42)                      # → True
StringQ("hello")                 # → True

# Pattern matching
MatchQ(42, _Integer)             # → True
MatchQ(List(1, 2), List(_, _))   # → True

# List operations
First(List(1, 2, 3))             # → 1
Length(List(1, 2, 3))            # → 3

# Control flow
If(Greater(5, 3), "big", "small") # → "big"

# Function definition
SetDelayed(factorial(0), 1)
SetDelayed(factorial(n_Integer), Times(n, factorial(Subtract(n, 1))))
factorial(5)                     # → 120

# Associations (modern syntax)
{name: "Bob", age: 30}           # → {name: "Bob", age: 30}
Part({name: "Bob", age: 30}, name) # → "Bob"
Keys({name: "Bob", age: 30})     # → List(name, age)
Values({name: "Bob", age: 30})   # → List("Bob", 30)
Length({name: "Bob", age: 30})   # → 2

# Nested associations
{
  person: {name: "Bob", age: 30},
  scores: [85, 92, 78]
}                                # → {person: {name: "Bob", age: 30}, scores: List(85, 92, 78)}

# Multi-line expressions (both file and REPL)
Plus(
  1,
  2,
  3
)                                # → 6
```