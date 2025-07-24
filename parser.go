package sexpr

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
	PrecedenceLogicalOr  // ||
	PrecedenceLogicalAnd // &&
	PrecedenceEquality   // ==, !=
	PrecedenceComparison // <, >, <=, >=
	PrecedenceSum        // +, -
	PrecedenceProduct    // *, /
	PrecedencePrefix     // -x, +x
)

var precedences = map[TokenType]Precedence{
	LBRACKET:     PrecedencePrefix, // High precedence for postfix indexing
	SEMICOLON:    PrecedenceCompound,
	SET:          PrecedenceAssign,
	SETDELAYED:   PrecedenceAssign,
	UNSET:        PrecedenceAssign,
	COLON:        PrecedenceAssign,
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
	return expr, nil
}

func (p *Parser) parseExpression() Expr {
	return p.parseInfixExpression(PrecedenceLowest)
}

func (p *Parser) parseInfixExpression(precedence Precedence) Expr {
	left := p.ParseAtom()

	for p.currentToken.Type != EOF && precedence < p.currentPrecedence() {
		if p.currentToken.Type == LBRACKET {
			left = p.parseIndexOrSlice(left)
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
	case INTEGER:
		expr = p.parseInteger()
		p.nextToken()
	case FLOAT:
		expr = p.parseFloat()
		p.nextToken()
	case STRING:
		expr = p.parseString()
		p.nextToken()
	case BOOLEAN:
		expr = p.parseBoolean()
		p.nextToken()
	case LBRACKET:
		expr = p.parseListLiteral()
	case LBRACE:
		expr = p.parseAssociationLiteral()
	case MINUS:
		expr = p.parsePrefixExpression()
	case PLUS:
		expr = p.parsePrefixExpression()
	case LPAREN:
		expr = p.parseGroupedExpression()
	default:
		p.addError(fmt.Sprintf("unexpected token: %s", p.currentToken.String()))
		return nil
	}

	return expr
}

func (p *Parser) parseSymbolOrList() Expr {
	symbolToken := p.currentToken
	p.nextToken()

	if p.currentToken.Type == LPAREN {
		return p.parseList(symbolToken.Value)
	}

	return NewSymbolAtom(symbolToken.Value)
}

func (p *Parser) parseList(head string) Expr {
	p.nextToken() // consume '('

	elements := []Expr{NewSymbolAtom(head)}

	if p.currentToken.Type == RPAREN {
		p.nextToken() // consume ')'
		return NewList(elements...)
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

	return NewList(elements...)
}

func (p *Parser) parseListLiteral() Expr {
	p.nextToken() // consume '['

	// Create a List expression with "List" as the head
	elements := []Expr{NewSymbolAtom("List")}

	// Handle empty list []
	if p.currentToken.Type == RBRACKET {
		p.nextToken() // consume ']'
		return NewList(elements...)
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

	return NewList(elements...)
}

func (p *Parser) parseAssociationLiteral() Expr {
	p.nextToken() // consume '{'

	// Create rules slice for Rule expressions
	var rules []Expr

	// Handle empty association {}
	if p.currentToken.Type == RBRACE {
		p.nextToken() // consume '}'
		// Create Association function call with no arguments for empty association
		return NewList(NewSymbolAtom("Association"))
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
	elements := []Expr{NewSymbolAtom("Association")}
	elements = append(elements, rules...)
	return NewList(elements...)
}

func (p *Parser) parseInteger() Expr {
	value, err := strconv.Atoi(p.currentToken.Value)
	if err != nil {
		p.addError(fmt.Sprintf("invalid integer: %s", p.currentToken.Value))
		return nil
	}

	return NewIntAtom(value)
}

func (p *Parser) parseFloat() Expr {
	value, err := strconv.ParseFloat(p.currentToken.Value, 64)
	if err != nil {
		p.addError(fmt.Sprintf("invalid float: %s", p.currentToken.Value))
		return nil
	}

	return NewFloatAtom(value)
}

func (p *Parser) parseString() Expr {
	value := p.unescapeString(p.currentToken.Value)
	return NewStringAtom(value)
}

func (p *Parser) parseBoolean() Expr {
	value := p.currentToken.Value == "True"
	return NewBoolAtom(value)
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

func (p *Parser) peekPrecedence() Precedence {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return PrecedenceLowest
}

func (p *Parser) currentPrecedence() Precedence {
	if prec, ok := precedences[p.currentToken.Type]; ok {
		return prec
	}
	return PrecedenceLowest
}

func (p *Parser) IsInfixOperator(tokenType TokenType) bool {
	switch tokenType {
	case SEMICOLON, SET, SETDELAYED, UNSET, COLON, OR, AND, EQUAL, UNEQUAL, SAMEQ, UNSAMEQ, LESS, GREATER, LESSEQUAL, GREATEREQUAL, PLUS, MINUS, MULTIPLY, DIVIDE:
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

	p.nextToken()
	right := p.parseInfixExpression(precedence)

	return p.createInfixExpr(operator.Type, left, right)
}

func (p *Parser) parsePrefixExpression() Expr {
	operator := p.currentToken
	p.nextToken()
	right := p.parseInfixExpression(PrecedencePrefix)

	return p.createPrefixExpr(operator.Type, right)
}

func (p *Parser) createInfixExpr(operator TokenType, left, right Expr) Expr {
	switch operator {
	case SEMICOLON:
		return NewList(NewSymbolAtom("CompoundStatement"), left, right)
	case SET:
		return NewList(NewSymbolAtom("Set"), left, right)
	case SETDELAYED:
		return NewList(NewSymbolAtom("SetDelayed"), left, right)
	case UNSET:
		return NewList(NewSymbolAtom("Unset"), left)
	case COLON:
		return NewList(NewSymbolAtom("Rule"), left, right)
	case OR:
		return NewList(NewSymbolAtom("Or"), left, right)
	case AND:
		return NewList(NewSymbolAtom("And"), left, right)
	case EQUAL:
		return NewList(NewSymbolAtom("Equal"), left, right)
	case UNEQUAL:
		return NewList(NewSymbolAtom("Unequal"), left, right)
	case SAMEQ:
		return NewList(NewSymbolAtom("SameQ"), left, right)
	case UNSAMEQ:
		return NewList(NewSymbolAtom("UnsameQ"), left, right)
	case LESS:
		return NewList(NewSymbolAtom("Less"), left, right)
	case GREATER:
		return NewList(NewSymbolAtom("Greater"), left, right)
	case LESSEQUAL:
		return NewList(NewSymbolAtom("LessEqual"), left, right)
	case GREATEREQUAL:
		return NewList(NewSymbolAtom("GreaterEqual"), left, right)
	case PLUS:
		return NewList(NewSymbolAtom("Plus"), left, right)
	case MINUS:
		return NewList(NewSymbolAtom("Subtract"), left, right)
	case MULTIPLY:
		return NewList(NewSymbolAtom("Times"), left, right)
	case DIVIDE:
		return NewList(NewSymbolAtom("Divide"), left, right)
	default:
		p.addError(fmt.Sprintf("unknown infix operator: %d", operator))
		return nil
	}
}

func (p *Parser) createPrefixExpr(operator TokenType, operand Expr) Expr {
	switch operator {
	case MINUS:
		return NewList(NewSymbolAtom("Minus"), operand)
	case PLUS:
		return operand // unary plus is identity
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
		return NewList(NewSymbolAtom("Part"), expr, firstExpr)

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
			// For start:, we want everything from start onwards
			// If start is negative, use Take(expr, start) for last n elements
			// If start is positive, use Drop(expr, start-1) since Drop removes the first n elements
			// But we can't easily detect negative at parse time, so we'll use a special function
			return NewList(NewSymbolAtom("TakeFrom"), expr, startExpr)
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
				return NewList(NewSymbolAtom("Take"), expr, endExpr)
			} else {
				// [start:end] syntax - Slice operation
				return NewList(NewSymbolAtom("SliceRange"), expr, startExpr, endExpr)
			}
		}
	} else {
		p.addError("expected ':' or ']' after slice expression")
		return expr
	}
}

// parseSliceExpression parses expressions inside slice brackets, treating colons as separators
func (p *Parser) parseSliceExpression() Expr {
	// Parse a simple expression that stops at colons and brackets
	// We'll use a custom precedence that's higher than colon assignment
	return p.parseInfixExpression(PrecedenceLogicalOr)
}

func ParseString(input string) (Expr, error) {
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	return parser.Parse()
}
