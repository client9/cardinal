package core

import (
	"fmt"
	"strings"
)

/*
// Precedence levels for InputForm formatting
type Precedence int

const (
	PrecedenceLowest Precedence = iota
	PrecedenceAssign
	PrecedenceLogicalOr
	PrecedenceLogicalAnd
	PrecedenceEquality
	PrecedenceComparison
	PrecedenceSum
	PrecedenceProduct
)
*/

// inputFormWithPrecedence formats a List with precedence-aware operator handling
func (l List) inputFormWithPrecedence(parentPrecedence Precedence) string {
	if l.Length() == 0 {
		return "List()"
	}

	// Check if this is a special function that has infix/shortcut representation
	switch l.Head() {
	case symbolList:
		// List(...) -> [...]
		if l.Length() == 0 {
			return "[]"
		}
		var elements []string
		for _, elem := range l.Tail() {
			elements = append(elements, elem.InputForm())
		}
		return fmt.Sprintf("[%s]", strings.Join(elements, ", "))

	case symbolAssociation:
		// Association(Rule(a,b), Rule(c,d)) -> {a: b, c: d}
		if l.Length() == 1 {
			return "{}"
		}
		var pairs []string
		for _, elem := range l.Tail() {
			if ruleList, ok := elem.(List); ok && ruleList.Length() == 2 && ruleList.Head() == symbolRule {
				args := ruleList.Tail()
				key := args[0].InputForm()
				value := args[1].InputForm()
				pairs = append(pairs, fmt.Sprintf("%s: %s", key, value))
				continue
			}
			// Fallback for non-Rule elements
			pairs = append(pairs, elem.InputForm())
		}
		return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))

	case symbolRule:
		// Rule(a, b) -> a: b
		if l.Length() == 2 {
			e := l.Tail()
			return fmt.Sprintf("%s: %s", e[0].InputForm(), e[1].InputForm())
		}

	case symbolRuleDelayed:
		// RuleDelayed(a, b) -> a => b
		if l.Length() == 2 {
			e := l.Tail()
			return fmt.Sprintf("%s => %s", e[0].InputForm(), e[1].InputForm())
		}

	case symbolSet:
		// Set(a, b) -> a = b
		if l.Length() == 2 {
			return l.formatInfixWithParens("=", PrecedenceAssign, parentPrecedence)
		}

	case symbolSetDelayed:
		// SetDelayed(a, b) -> a := b
		if l.Length() == 2 {
			return l.formatInfixWithParens(":=", PrecedenceAssign, parentPrecedence)
		}

	case symbolPlus:
		// Plus(a, b, ...) -> a + b + ...
		if l.Length() > 1 {
			return l.formatLeftAssociativeInfix("+", PrecedenceSum, parentPrecedence)
		}

	case symbolTimes:
		// Times(a, b, ...) -> a * b * ...
		if l.Length() > 1 {
			return l.formatLeftAssociativeInfix("*", PrecedenceProduct, parentPrecedence)
		}

	case symbolSubtract:
		// Subtract(a, b) -> a - b
		if l.Length() == 2 {
			return l.formatInfixWithParens("-", PrecedenceSum, parentPrecedence)
		}

	case symbolDivide:
		// Divide(a, b) -> a / b
		if l.Length() == 2 {
			return l.formatInfixWithParens("/", PrecedenceProduct, parentPrecedence)
		}

	case symbolEqual:
		// Equal(a, b) -> a == b
		if l.Length() == 2 {
			return l.formatInfixWithParens("==", PrecedenceEquality, parentPrecedence)
		}

	case symbolUnequal:
		// Unequal(a, b) -> a != b
		if l.Length() == 2 {
			return l.formatInfixWithParens("!=", PrecedenceEquality, parentPrecedence)
		}

	case symbolSameQ:
		// SameQ(a, b) -> a === b
		if l.Length() == 2 {
			return l.formatInfixWithParens("===", PrecedenceEquality, parentPrecedence)
		}

	case symbolUnsameQ:
		// UnsameQ(a, b) -> a =!= b
		if l.Length() == 2 {
			return l.formatInfixWithParens("=!=", PrecedenceEquality, parentPrecedence)
		}

	case symbolLess:
		// Less(a, b) -> a < b
		if l.Length() == 2 {
			return l.formatInfixWithParens("<", PrecedenceComparison, parentPrecedence)
		}

	case symbolGreater:
		// Greater(a, b) -> a > b
		if l.Length() == 2 {
			return l.formatInfixWithParens(">", PrecedenceComparison, parentPrecedence)
		}

	case symbolLessEqual:
		// LessEqual(a, b) -> a <= b
		if l.Length() == 2 {
			return l.formatInfixWithParens("<=", PrecedenceComparison, parentPrecedence)
		}

	case symbolGreaterEqual:
		// GreaterEqual(a, b) -> a >= b
		if l.Length() == 2 {
			return l.formatInfixWithParens(">=", PrecedenceComparison, parentPrecedence)
		}

	case symbolAnd:
		// And(a, b, ...) -> a && b && ...
		if l.Length() > 1 {
			return l.formatLeftAssociativeInfix("&&", PrecedenceLogicalAnd, parentPrecedence)
		}

	case symbolOr:
		// Or(a, b, ...) -> a || b || ...
		if l.Length() > 1 {
			return l.formatLeftAssociativeInfix("||", PrecedenceLogicalOr, parentPrecedence)
		}
	}

	// Default: function call format Head(arg1, arg2, ...)
	var elements []string
	for _, elem := range l.Tail() {
		elements = append(elements, elem.InputForm())
	}
	return fmt.Sprintf("%s(%s)", l.Head().String(), strings.Join(elements, ", "))
}

// formatInfixWithParens formats a binary infix operation with parentheses if needed
func (l List) formatInfixWithParens(op string, opPrecedence, parentPrecedence Precedence) string {
	args := l.Tail()
	left := l.getInputFormWithPrecedence(args[0], opPrecedence)
	right := l.getInputFormWithPrecedence(args[1], opPrecedence)
	result := fmt.Sprintf("%s %s %s", left, op, right)

	if opPrecedence < parentPrecedence {
		return fmt.Sprintf("(%s)", result)
	}
	return result
}

// formatLeftAssociativeInfix formats left-associative infix operations like a + b + c
func (l List) formatLeftAssociativeInfix(op string, opPrecedence, parentPrecedence Precedence) string {
	var parts []string
	for _, elem := range l.Tail() {
		parts = append(parts, l.getInputFormWithPrecedence(elem, opPrecedence+1)) // Higher precedence for right operand
	}
	result := strings.Join(parts, fmt.Sprintf(" %s ", op))

	if opPrecedence < parentPrecedence {
		return fmt.Sprintf("(%s)", result)
	}
	return result
}

// getInputFormWithPrecedence gets InputForm with precedence context for proper parenthesization
func (l List) getInputFormWithPrecedence(expr Expr, precedence Precedence) string {
	if list, ok := expr.(List); ok {
		return list.inputFormWithPrecedence(precedence)
	}
	return expr.InputForm()
}
