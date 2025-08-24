package core

import (
	"fmt"
	"strconv"
	"strings"
)

type Precedence int

const (
	_ Precedence = iota
	PrecedenceLowest
	PrecedenceCompound   // ; (compound statements)
	PrecedenceAssign     // =, :=, =.
	PrecedenceRule       // : (rule shorthand)
	PrecedenceLogicalOr  // ||
	PrecedenceLogicalAnd // &&
	PrecedenceEquality   // ==, !=
	PrecedenceComparison // <, >, <=, >=
	PrecedenceSum        // +, -
	PrecedenceProduct    // *, /
	PrecedenceUnary      // unary -x, +x (lower than power)
	PrecedencePower      // ^ (right associative)
	PrecedencePostfix    // high precedence postfix operators
)

var precedences = map[TokenType]Precedence{
	LBRACKET:     PrecedencePostfix, // High precedence for postfix indexing
	LPAREN:       PrecedencePostfix, // High precedence for postfix function application
	AMPERSAND:    PrecedenceRule,    // Low precedence for Function syntax (&) to bind to larger expressions
	SEMICOLON:    PrecedenceCompound,
	SET:          PrecedenceAssign,
	SETDELAYED:   PrecedenceAssign,
	UNSET:        PrecedenceAssign,
	COLON:        PrecedenceRule,
	RULEDELAYED:  PrecedenceRule,
	OR:           PrecedenceLogicalOr,
	AND:          PrecedenceLogicalAnd,
	EQUAL:        PrecedenceEquality,
	UNEQUAL:      PrecedenceEquality,
	SAMEQ:        PrecedenceEquality,
	UNSAMEQ:      PrecedenceEquality,
	LESS:         PrecedenceComparison,
	GREATER:      PrecedenceComparison,
	LESSEQUAL:    PrecedenceComparison,
	GREATEREQUAL: PrecedenceComparison,
	PLUS:         PrecedenceSum,
	MINUS:        PrecedenceSum,
	MULTIPLY:     PrecedenceProduct,
	DIVIDE:       PrecedenceProduct,
	CARET:        PrecedencePower,
	NOT:          PrecedenceUnary,
}

type Parser struct {
	lexer        *Lexer
	currentToken Token
	peekToken    Token
	errors       []string
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, fmt.Sprintf("Parse error at position %d: %s", p.currentToken.Position, msg))
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) Parse() (Expr, error) {
	expr := p.parseExpression()
	if len(p.errors) > 0 {
		return nil, fmt.Errorf("parse errors: %s", strings.Join(p.errors, "; "))
	}
	if expr == nil {
		// triggered when input is nothing
		return NewSymbolNull(), nil
	}
	return expr, nil
}

func (p *Parser) parseExpression() Expr {
	return p.parseInfixExpression(PrecedenceLowest)
}

func (p *Parser) parseInfixExpression(precedence Precedence) Expr {
	// Special case: if we hit RPAREN, RBRACKET, RBRACE, or EOF when expecting an expression,
	// return nil to signal that no expression is available
	if p.currentToken.Type == RPAREN || p.currentToken.Type == RBRACKET || p.currentToken.Type == RBRACE || p.currentToken.Type == EOF {
		return nil
	}

	left := p.ParseAtom()

	for p.currentToken.Type != EOF && precedence < p.currentPrecedence() {
		if p.currentToken.Type == LBRACKET {
			left = p.parseIndexOrSlice(left)
		} else if p.currentToken.Type == LPAREN {
			left = p.parseFunctionApplication(left)
		} else if p.currentToken.Type == AMPERSAND {
			left = p.parseFunctionShorthand(left)
		} else if p.IsInfixOperator(p.currentToken.Type) {
			left = p.parseInfixOperation(left)
		} else {
			break
		}
	}

	return left
}

func (p *Parser) ParseAtom() Expr {
	var expr Expr

	switch p.currentToken.Type {
	case SYMBOL:
		expr = p.parseSymbolOrList()
	case UNDERSCORE:
		expr = p.parseUnderscorePattern()
	case INTEGER:
		expr = p.parseInteger()
		p.nextToken()
	case FLOAT:
		expr = p.parseFloat()
		p.nextToken()
	case STRING:
		expr = p.parseString()
		p.nextToken()
	case LBRACKET:
		expr = p.parseListLiteral()
	case LBRACE:
		expr = p.parseAssociationLiteral()
	case MINUS:
		expr = p.parsePrefixExpression()
	case PLUS:
		expr = p.parsePrefixExpression()
	case NOT:
		expr = p.parsePrefixExpression()
	case LPAREN:
		expr = p.parseGroupedExpression()
	case SEMICOLON, EOF:
		// Empty expression - return Null without consuming token
		return NewSymbolNull()
	default:
		p.addError(fmt.Sprintf("unexpected token: %s", p.currentToken.String()))
		return nil
	}

	return expr
}

func (p *Parser) parseSymbolOrList() Expr {
	symbolToken := p.currentToken
	p.nextToken()

	// Check if this is a pattern: SYMBOL + UNDERSCORE(s) + optional SYMBOL
	if p.currentToken.Type == UNDERSCORE {
		return p.parsePatternFromSymbol(symbolToken.Value)
	}

	if p.currentToken.Type == LPAREN {
		return p.parseList(symbolToken.Value)
	}

	return NewSymbol(symbolToken.Value)
}

func (p *Parser) parseList(head string) Expr {
	p.nextToken() // consume '('

	elements := []Expr{}

	if p.currentToken.Type == RPAREN {
		p.nextToken() // consume ')'
		return NewList(head, elements...)
	}

	for {
		expr := p.parseExpression()
		if expr != nil {
			elements = append(elements, expr)
		}

		if p.currentToken.Type == RPAREN {
			p.nextToken() // consume ')'
			break
		}

		if p.currentToken.Type == COMMA {
			p.nextToken() // consume ','
			continue
		}

		if p.currentToken.Type == EOF {
			p.addError("unexpected EOF, expected ')'")
			break
		}

		p.addError(fmt.Sprintf("expected ',' or ')', got %s", p.currentToken.String()))
		p.nextToken()
	}

	return NewList(head, elements...)
}

func (p *Parser) parseListLiteral() Expr {
	p.nextToken() // consume '['

	// Create a List expression with "List" as the head
	elements := []Expr{}

	// Handle empty list []
	if p.currentToken.Type == RBRACKET {
		p.nextToken() // consume ']'
		return NewList("List", elements...)
	}

	// Parse list elements
	for {
		expr := p.parseExpression()
		if expr != nil {
			elements = append(elements, expr)
		}

		// Check for closing bracket
		if p.currentToken.Type == RBRACKET {
			p.nextToken() // consume ']'
			break
		}

		// Check for comma separator
		if p.currentToken.Type == COMMA {
			p.nextToken() // consume ','

			// Handle optional trailing comma: [1,2,3,]
			if p.currentToken.Type == RBRACKET {
				p.nextToken() // consume ']'
				break
			}
			continue
		}

		// Handle EOF
		if p.currentToken.Type == EOF {
			p.addError("unexpected EOF, expected ']'")
			break
		}

		// Unexpected token
		p.addError(fmt.Sprintf("expected ',' or ']', got %s", p.currentToken.String()))
		p.nextToken()
	}

	return NewList("List", elements...)
}

func (p *Parser) parseAssociationLiteral() Expr {
	p.nextToken() // consume '{'

	// Create rules slice for Rule expressions
	var rules []Expr

	// Handle empty association {}
	if p.currentToken.Type == RBRACE {
		p.nextToken() // consume '}'
		// Create Association function call with no arguments for empty association
		return NewList("Association")
	}

	// Parse expressions (expecting Rule expressions from key:value infix parsing)
	for {
		// Parse expression (should be key:value which becomes Rule(key, value))
		expr := p.parseExpression()
		if expr != nil {
			rules = append(rules, expr)
		}

		// Check for closing brace
		if p.currentToken.Type == RBRACE {
			p.nextToken() // consume '}'
			break
		}

		// Check for comma separator
		if p.currentToken.Type == COMMA {
			p.nextToken() // consume ','

			// Handle optional trailing comma: {a: 1, b: 2,}
			if p.currentToken.Type == RBRACE {
				p.nextToken() // consume '}'
				break
			}
			continue
		}

		// Handle EOF
		if p.currentToken.Type == EOF {
			p.addError("unexpected EOF, expected '}'")
			break
		}

		// Unexpected token
		p.addError(fmt.Sprintf("expected ',' or '}', got %s", p.currentToken.String()))
		p.nextToken()
	}

	// Create Association function call with Rule expressions
	return NewList("Association", rules...)
}

func (p *Parser) parseInteger() Expr {
	value, err := strconv.ParseInt(p.currentToken.Value, 10, 64)
	if err != nil {
		p.addError(fmt.Sprintf("invalid integer: %s", p.currentToken.Value))
		return nil
	}

	return NewInteger(value)
}

func (p *Parser) parseFloat() Expr {
	value, err := strconv.ParseFloat(p.currentToken.Value, 64)
	if err != nil {
		p.addError(fmt.Sprintf("invalid float: %s", p.currentToken.Value))
		return nil
	}

	return NewReal(value)
}

func (p *Parser) parseString() Expr {
	value := p.unescapeString(p.currentToken.Value)
	return NewString(value)
}

func (p *Parser) unescapeString(s string) string {
	result := strings.Builder{}
	i := 0
	for i < len(s) {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n':
				result.WriteByte('\n')
			case 't':
				result.WriteByte('\t')
			case 'r':
				result.WriteByte('\r')
			case '\\':
				result.WriteByte('\\')
			case '"':
				result.WriteByte('"')
			default:
				result.WriteByte(s[i+1])
			}
			i += 2
		} else {
			result.WriteByte(s[i])
			i++
		}
	}
	return result.String()
}

func (p *Parser) currentPrecedence() Precedence {
	if prec, ok := precedences[p.currentToken.Type]; ok {
		return prec
	}
	return PrecedenceLowest
}

func (p *Parser) IsInfixOperator(tokenType TokenType) bool {
	switch tokenType {
	case SEMICOLON, SET, SETDELAYED, UNSET, COLON, RULEDELAYED, OR, AND, EQUAL, UNEQUAL, SAMEQ, UNSAMEQ, LESS, GREATER, LESSEQUAL, GREATEREQUAL, PLUS, MINUS, MULTIPLY, DIVIDE, CARET:
		return true
	default:
		return false
	}
}

func (p *Parser) parseInfixOperation(left Expr) Expr {
	operator := p.currentToken
	precedence := p.currentPrecedence()

	// Special case for UNSET: it's a postfix unary operator
	if operator.Type == UNSET {
		p.nextToken()
		return p.createInfixExpr(operator.Type, left, nil)
	}

	// Special case for SET with slice/part expressions: convert to slice assignment
	if operator.Type == SET && p.isSliceExpression(left) {
		p.nextToken()
		right := p.parseInfixExpression(precedence)
		return p.createSliceAssignment(left, right)
	}

	p.nextToken()

	// Power (^) is right-associative, so use precedence - 1
	if operator.Type == CARET {
		right := p.parseInfixExpression(precedence - 1)
		return p.createInfixExpr(operator.Type, left, right)
	}

	right := p.parseInfixExpression(precedence)

	// Special case for semicolon: if no right operand, use Null
	if operator.Type == SEMICOLON && right == nil {
		right = NewSymbol("Null")
	}

	return p.createInfixExpr(operator.Type, left, right)
}

func (p *Parser) parsePrefixExpression() Expr {
	operator := p.currentToken
	p.nextToken()
	right := p.parseInfixExpression(PrecedenceUnary)

	return p.createPrefixExpr(operator.Type, right)
}

func (p *Parser) createInfixExpr(operator TokenType, left, right Expr) Expr {
	switch operator {
	case SEMICOLON:
		// Flatten nested CompoundExpressions into a single flat list
		if leftList, ok := left.(List); ok {
			if leftList.Head() == "CompoundExpression" {
				// Left is already a CompoundExpression, append right to it
				elements := make([]Expr, leftList.Length()+2)
				copy(elements, leftList.AsSlice())
				elements[len(elements)-1] = right
				return NewListFromExprs(elements...)
			}
		}
		return NewList("CompoundExpression", left, right)
	case SET:
		return NewList("Set", left, right)
	case SETDELAYED:
		return NewList("SetDelayed", left, right)
	case UNSET:
		return NewList("Unset", left)
	case COLON:
		return NewList("Rule", left, right)
	case RULEDELAYED:
		return NewList("RuleDelayed", left, right)
	case OR:
		return NewList("Or", left, right)
	case AND:
		return NewList("And", left, right)
	case EQUAL:
		return NewList("Equal", left, right)
	case UNEQUAL:
		return NewList("Unequal", left, right)
	case SAMEQ:
		return NewList("SameQ", left, right)
	case UNSAMEQ:
		return NewList("UnsameQ", left, right)
	case LESS:
		return NewList("Less", left, right)
	case GREATER:
		return NewList("Greater", left, right)
	case LESSEQUAL:
		return NewList("LessEqual", left, right)
	case GREATEREQUAL:
		return NewList("GreaterEqual", left, right)
	case PLUS:
		// Flatten nested Plus expressions into a single flat list
		if leftList, ok := left.(List); ok {
			if leftList.Head() == "Plus" {
				// Left is already a Plus, append right to it
				elements := make([]Expr, leftList.Length()+2)
				copy(elements, leftList.AsSlice())
				elements[len(elements)-1] = right
				return NewListFromExprs(elements...)
			}
		}
		return NewList("Plus", left, right)
	case MINUS:
		return NewList("Subtract", left, right)
	case MULTIPLY:
		// Flatten nested Times expressions into a single flat list
		if leftList, ok := left.(List); ok {
			if leftList.Head() == "Times" {
				// Left is already a Times, append right to it
				elements := make([]Expr, leftList.Length()+2)
				copy(elements, leftList.AsSlice())
				elements[len(elements)-1] = right
				return NewListFromExprs(elements...)
			}
		}
		return NewList("Times", left, right)
	case DIVIDE:
		return NewList("Divide", left, right)
	case CARET:
		return NewList("Power", left, right)
	default:
		p.addError(fmt.Sprintf("unknown infix operator: %d", operator))
		return nil
	}
}

func (p *Parser) createPrefixExpr(operator TokenType, operand Expr) Expr {
	switch operator {
	case MINUS:
		return NewList("Minus", operand)
	case PLUS:
		return operand // unary plus is identity
	case NOT:
		return NewList("Not", operand)
	default:
		p.addError(fmt.Sprintf("unknown prefix operator: %d", operator))
		return nil
	}
}

func (p *Parser) parseGroupedExpression() Expr {
	p.nextToken() // consume '('

	expr := p.parseExpression()

	if p.currentToken.Type != RPAREN {
		p.addError("expected ')' after grouped expression")
		return nil
	}

	p.nextToken() // consume ')'
	return expr
}

// parseIndexOrSlice handles postfix [index] and [start:end] syntax
func (p *Parser) parseIndexOrSlice(expr Expr) Expr {
	p.nextToken() // consume '['

	// Check for empty brackets []
	if p.currentToken.Type == RBRACKET {
		p.addError("empty brackets are not allowed")
		return expr
	}

	// Parse the first expression (could be index, start, or just ':')
	var firstExpr Expr
	var hasFirstExpr bool

	if p.currentToken.Type == COLON {
		// [:end] syntax - no start expression
		hasFirstExpr = false
	} else {
		// Parse first expression, but stop at colon (don't treat colon as infix operator here)
		firstExpr = p.parseSliceExpression()
		hasFirstExpr = true
	}

	// Check what comes next
	if p.currentToken.Type == RBRACKET {
		// Simple index: expr[index]
		if !hasFirstExpr {
			p.addError("expected expression before ']'")
			return expr
		}
		p.nextToken() // consume ']'
		return NewList("Part", expr, firstExpr)

	} else if p.currentToken.Type == COLON {
		// Slice syntax: expr[start:end] or expr[:end] or expr[start:]
		p.nextToken() // consume ':'

		var startExpr, endExpr Expr

		if hasFirstExpr {
			startExpr = firstExpr
		}

		// Check for end expression
		if p.currentToken.Type == RBRACKET {
			// expr[start:] syntax - no end expression
			if !hasFirstExpr {
				p.addError("slice cannot be empty on both sides of ':'")
				return expr
			}
			p.nextToken() // consume ']'
			// Convert to Drop operation: Drop(expr, start-1)
			if startExpr == nil {
				return expr
			}
			return NewList("Take", expr, NewList("List", startExpr, NewInteger(-1)))
		} else {
			// Parse end expression
			endExpr = p.parseSliceExpression()
			if p.currentToken.Type != RBRACKET {
				p.addError("expected ']' after slice expression")
				return expr
			}
			p.nextToken() // consume ']'

			// Generate appropriate slice expression
			if startExpr == nil {
				// [:end] syntax - Take first n elements
				return NewList("Take", expr, endExpr)
			} else {
				// [start:end] syntax - Slice operation
				return NewList("Take", expr, NewList("List", startExpr, endExpr))
			}
		}
	} else {
		p.addError("expected ':' or ']' after slice expression")
		return expr
	}
}

// parseFunctionApplication handles postfix function application like Function(x, x+1)(5)
func (p *Parser) parseFunctionApplication(expr Expr) Expr {
	p.nextToken() // consume '('

	var args []Expr

	// Handle empty argument list: f()
	if p.currentToken.Type == RPAREN {
		p.nextToken()                 // consume ')'
		return NewListFromExprs(expr) // Just the function with no arguments
	}

	// Parse arguments
	for {
		arg := p.parseExpression()
		args = append(args, arg)

		if p.currentToken.Type == COMMA {
			p.nextToken() // consume ','
			continue
		} else if p.currentToken.Type == RPAREN {
			p.nextToken() // consume ')'
			break
		} else {
			p.addError("expected ',' or ')' in function application")
			return expr
		}
	}

	// Create function application: put the function expression as head, followed by arguments
	elements := make([]Expr, len(args)+1)
	elements[0] = expr
	copy(elements[1:], args)

	return NewListFromExprs(elements...)
}

// parseSliceExpression parses expressions inside slice brackets, treating colons as separators
func (p *Parser) parseSliceExpression() Expr {
	// Parse a simple expression that stops at colons and brackets
	// We'll use a custom precedence that's higher than colon assignment
	return p.parseInfixExpression(PrecedenceLogicalOr)
}

// isSliceExpression checks if an expression is a slice or part expression
func (p *Parser) isSliceExpression(expr Expr) bool {
	list, ok := expr.(List)
	if !ok {
		return false
	}
	headName := list.Head()
	return headName == "Part" || headName == "Take"
}

// createSliceAssignment creates the appropriate slice assignment AST node
func (p *Parser) createSliceAssignment(sliceExpr Expr, value Expr) Expr {
	list := sliceExpr.(List)
	headName := list.Head()

	switch headName {
	case "Part":
		// Part(expr, index) = value -> PartSet(expr, index, value)
		if list.Length() != 2 {
			p.addError("Part expression must have exactly 2 arguments for assignment")
			return nil
		}
		args := list.Tail()
		return NewList("PartSet", args[0], args[1], value)
	case "Take":
		// Take(expr, n) = value -> SliceSet(expr, 1, n, value)
		// Take(expr, [n, m]) = value --> SliceSet(expr, n, m, value)
		if list.Length() != 2 {
			p.addError("Take expression must have exactly 2 arguments for assignment")
			return nil
		}
		args := list.Tail()
		if _, ok := ExtractInt64(args[1]); ok {
			return NewList("SliceSet", args[0], NewInteger(1), args[1], value)
		}
		if rangelist, ok := args[1].(List); ok && rangelist.Length() == 2 {
			largs := rangelist.Tail()
			return NewList("SliceSet", args[0], largs[0], largs[1], value)
		}
	}
	p.addError(fmt.Sprintf("Unknown slice expression type: %s", headName))
	return nil
}

// parseUnderscorePattern parses anonymous patterns (_, __, ___, _Integer, __Integer, ___Integer)
func (p *Parser) parseUnderscorePattern() Expr {
	// Get underscore count from token value
	underscoreToken := p.currentToken
	underscoreCount := len(underscoreToken.Value)
	p.nextToken() // consume the underscore token

	// Check if there's a type after the underscores
	var typeName string
	if p.currentToken.Type == SYMBOL {
		typeName = p.currentToken.Value
		p.nextToken()
	}

	// Create the appropriate blank expression based on underscore count
	var blankExpr Expr
	if underscoreCount >= 3 {
		if typeName != "" {
			blankExpr = NewList("BlankNullSequence", NewSymbol(typeName))
		} else {
			blankExpr = NewList("BlankNullSequence")
		}
	} else if underscoreCount == 2 {
		if typeName != "" {
			blankExpr = NewList("BlankSequence", NewSymbol(typeName))
		} else {
			blankExpr = NewList("BlankSequence")
		}
	} else {
		if typeName != "" {
			blankExpr = NewList("Blank", NewSymbol(typeName))
		} else {
			blankExpr = NewList("Blank")
		}
	}

	// Anonymous pattern - just return the blank expression
	return blankExpr
}

// parsePatternFromSymbol parses named patterns (x_, x__, x___, x_Integer, x__Integer, x___Integer)
func (p *Parser) parsePatternFromSymbol(varName string) Expr {
	// Get underscore count from token value
	underscoreToken := p.currentToken
	underscoreCount := len(underscoreToken.Value)
	p.nextToken() // consume the underscore token

	// Check if there's a type after the underscores
	var typeName string
	if p.currentToken.Type == SYMBOL {
		typeName = p.currentToken.Value
		p.nextToken()
	}

	// Create the appropriate blank expression based on underscore count
	var blankExpr Expr
	if underscoreCount >= 3 {
		if typeName != "" {
			blankExpr = NewList("BlankNullSequence", NewSymbol(typeName))
		} else {
			blankExpr = NewList("BlankNullSequence")
		}
	} else if underscoreCount == 2 {
		if typeName != "" {
			blankExpr = NewList("BlankSequence", NewSymbol(typeName))
		} else {
			blankExpr = NewList("BlankSequence")
		}
	} else {
		if typeName != "" {
			blankExpr = NewList("Blank", NewSymbol(typeName))
		} else {
			blankExpr = NewList("Blank")
		}
	}

	// Named pattern - wrap in Pattern(varName, blankExpr)
	return NewList("Pattern", NewSymbol(varName), blankExpr)
}

// parseFunctionShorthand handles the & postfix operator: expr & -> Function(expr)
func (p *Parser) parseFunctionShorthand(expr Expr) Expr {
	p.nextToken() // consume '&'
	return NewList("Function", expr)
}

func ParseString(input string) (Expr, error) {
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	return parser.Parse()
}

func MustParse(input string) Expr {
	out, err := ParseString(input)
	if err == nil {
		return out
	}
	panic("Unable to parse: " + input)
}
