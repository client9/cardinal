# Project Status

## Syntax and Parser

- Works!
- `#` for single line comments
- `{ ... }` for 'Associations' map/dict
- `:` for `Rule` or key-value pair/tuple, i.e. a:b
- `[ ... ]` for list literals
- `fn(...)` for functional calls/objects
- Double quotes for strings, single quotes for runes
- proper slice/indexing in go/python style  `[2:10]`

- [x]: add Rune (character) literal parsing

- TODO: Fuzz, and negative tests. 
- TODO: Decide on `=` or `:=` for assignment, remove other.
- TODO: Remove old `->`, `=>` for RuleDelayed (in Mathematica it's `:>` but we obsoleted it)
- TODO: Decide on multiline comments `/* ... */`, `(* ... *)`, or some other python-like things with `#`
- TODO: Decide on MMA's postfix operators, e.g. `./`, `//`, `@@`, etc.

## Symbols

- Everything uses "unique" handles defined in core/symbol.  no string literals for symbols.
- Fast pointer check for equality
- Builtin symbols generated from code/comments in 'builtin' directory.
- [x] add `SymbolName` (trivial)
- [x] add `Symbol` constructor
- [x] Fix parser to accept unicode symbols

## Numerics

- Unlimited precision integers and rational types, handling Plus,Times,Division,Minus,Power
- Switches from machine integers to unlimited automatically.

- [x] TODO: Real numbers are only machine precision
- [ ] TODO: Testing needs major improvement, especially for precision issues
- [ ] TODO: Complex numbers
- [ ] TODO: Trig functions
- [ ] TODO: E, Log
- [x] Abs- Done.
- [x] Sqrt - Done
- [ ] TODO: Internal precomputed versions of E, Pi, etc to 5,000 digits.
- [ ] Improve N function

## Lists

- Many list operations complete, Take,Drop,First,Rest,etc.
- [ ] TODO: what basics are missing?
- Vectors and Matrix support is primitive

- TODO: general "Level Spec" needs  implementation. 

## Strings

- Strings are "sliceable" type, so any function that works on lists will work on strings.

- [x]: add Rune as a fundamental atom.
- [ ]: Do we need to expose Rune as a type, or is it just an optimization
       SOmehow Python just exposes characters as strings of length 1
- [x]: change parser to handle rune literals 'a', 'b', 'c' as per Go standards
- [ ]: TODO:  ToString function
- [ ]: TODO: "hello"[1] = 'X', fails.
- [ ]: TODO: add Character based functions
- [ ]: TODO: add Regexp
- [ ]: TODO: consider adding []rune specialization instead of List of Expr/Rune.

## Sets and Associations

- Basic `Association` type works (a map of Expression to Expression)
- Needs some cleanup

- TODO: pure set types.. Decide if Sets are ordered or not.  Sets aren't really needed but
        easy to add since we have the `{... }` syntax already.  See Python Sets.
- TODO: Union, works only for true List objects, need to expand to any list-like object. See Pattern matching below for blocker.
- TODO: Difference, Intersection, Complement

## Pattern Matching

- Excellent matching virtual machine with no backtracking, and fast 'one-step NFA' for simple matching.

- TODO: Condition predicate
- TODO: consolidate.  Currently one system for function lookup, another for generic MatchQ stuff
- TODO: add lazy matching (should be easy)
- TODO: Add  pattern or one-or-more or zero-or-more "list-like" objects.  ANy list object can be expressed with `_(...)`.  Need to extend to `__(...)` and `___(...)`
- TODO: Optimzed pattern matching on above, if `_(___)`, `__(___)`, `___(___)` then only check if input a list-like.  No need to descend into list.

## Documentation

- TODO: overload 'go doc' to generate documentation?

## REPL

- TODO: Fix "In" and "Out" support
- TODO: General cleanup
- TODO: In programming mode, do not print "last result"

## System

- [ ] Contexts... Currently everything is in one context.
- [ ] Loadable packages.  Will be not like MMA
- [ ] Add atttribute for Unevaluated is Error.  If a function has no match in the registry, return error.  For numeric function, no evaluation is fine(e.g. Cos(1)). 

## Testing

- While code coverage stats appear to be solid, the testing is quite weak
- TODO: need major rethink

## Go Integration

- TODO: Use Reflect to automatically convert go values to Cardinal and back

## Cleanup

- [ ] TODO: two matching systems
- [ ] TODO: unclear why there is a `Function` atom in `core`
- [ ] TODO: InputForm is implemented in the interface of Expr, however objects doesn't parse themselves so unclear why they would know how to print various forms.  Perhaps move to separate function.
- [  ] TODO: unclear if the Equal method in Expr is correct or needed.
- [  ] TODO: Slice assignment is goofy, and uses `PartSet(expression, index, value)` as a special thing in the Parser/Evaluator.  Overload `Set` such as  `Set( Part(expresssion, index), newvalue)`, which can call the same code. 
- [ ] TODO: in core, the function "CanoncialCompare" is really "Less", consider making it a true compare that return -1,0,1 an rename to "CompareExpr" or so.
- [ ] TODO:  `a / b` is parsed as Divide(a, b) then transformed into Times(a, Power(b,-1)).  The Divide step could be skipped.
