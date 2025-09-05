# Project Status

## Syntax and Parser

- Works!
- `#` for single line comments
- `{ ... }` for 'Associations' map/dict
- `:` for `Rule` or key-value pair/tuple, i.e. a:b
- `[ ... ]` for list literals
- `fn(...)` for functional calls/objects
- proper slice/indexing in go/python style  `[2:10]`

- [x]: add Rune (character) literal parsing

- TODO: Fuzz, and negative tests.  Can panic if malformed input
- TODO: Decide on `=` or `:=` for assignment, remove other.
- TODO: Remove old `=>` for RuleDelayed (in Mathematica it's `:>` but we obsoleted it)
- TODO: Decide on multiline comments `/* ... */`
- TODO: Decide on MMA's postfix operators, e.g. `/.`, `//`, `@@`, etc.

## Symbols

- No (?) string literals in code.  Everything uses "unique" handles defined in core/symbol.
- Builtin symbols generated from code/comments in 'builtin' directory.

## Numerics

- Unlimited precision integers and rational types, handling Plus,Times,Division,Minus,Power
- Switches from machine integers to unlimited automatically.

- TODO: Real numbers are only machine precision
- TODO: Testing needs major improvement
- TODO: Complex numbers

## Lists

- Many list operations complete, Take,Drop,First,Rest,etc.
- Vectors and Matrix support is primitive

- TODO: general "Level Spec" needs  implementation. 

## Strings

- Strings are "sliceable" type, so any function that works on lists will work on strings.

- [x]: add Rune as a fundamental atom.
- TODO: check "hello"[1] = 'X'
- TODO: change parser to handle rune literals 'a', 'b', 'c' as per Go standards
- TODO: add Character based functions
- TODO: add Regexp

## Sets and Associations

- Basical Association type works (a map of Expression to Expression)
- Needs some cleanup

- TODO: pure set types.. Decide if Sets are ordered or not.
- TODO: Union, Difference, Intersection

## Pattern Matching

- Excellent matching virtual machine with no backtracking, and fast 'one-step NFA' for simple matching.

- TODO: consolidate.  Currently one system for function lookup, another for generic MatchQ stuff
- TODO: add lazy matching (should be easy)

## Documentation

- TODO: overload 'go doc' to generate documentation?

## REPL

- TODO: and "In" and "Out" support
- TODO: General cleanup
- TODO: In programming mode, do not print "last result"

## Testing

- While code coverage stats appear to be solid, the testing is quite weak
- TODO: need major rethink

## Go Integration

- TODO: Reflect to automatically map a go functions to Cardinal.

## Cleanup

- TODO: two matching systems
- TODO: unclear why there is a `Function` atom in `core`
- TODO: InputForm is implemented in the interface of Expr, however objects doesn't parse themselves so unclear why they would know how to print various forms.  Perhaps move to separate function.
- TODO: unclear if the Equal method in Expr is correct or needed.


