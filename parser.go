package sexpr

import (
	"fmt"
	"strconv"
	"strings"
)

type Precedence int

const (
	_ Precedence = iota
	LOWEST
	ASSIGN      // =, :=, =.
	LOGICAL_OR  // ||
	LOGICAL_AND // &&
	EQUALITY    // ==, !=
	COMPARISON  // <, >, <=, >=
	SUM         // +, -
	PRODUCT     // *, /
	PREFIX      // -x, +x
)

var precedences = map[TokenType]Precedence{
	SET:          ASSIGN,
	SETDELAYED:   ASSIGN,
	UNSET:        ASSIGN,
	OR:           LOGICAL_OR,
	AND:          LOGICAL_AND,
	EQUAL:        EQUALITY,
	UNEQUAL:      EQUALITY,
	SAMEQ:        EQUALITY,
	UNSAMEQ:      EQUALITY,
	LESS:         COMPARISON,
	GREATER:      COMPARISON,
	LESSEQUAL:    COMPARISON,
	GREATEREQUAL: COMPARISON,
	PLUS:         SUM,
	MINUS:        SUM,
	MULTIPLY:     PRODUCT,
	DIVIDE:       PRODUCT,
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
	return p.parseInfixExpression(LOWEST)
}

func (p *Parser) parseInfixExpression(precedence Precedence) Expr {
	left := p.ParseAtom()

	for p.currentToken.Type != EOF && precedence < p.currentPrecedence() {
		if !p.IsInfixOperator(p.currentToken.Type) {
			break
		}

		left = p.parseInfixOperation(left)
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
	return LOWEST
}

func (p *Parser) currentPrecedence() Precedence {
	if prec, ok := precedences[p.currentToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) IsInfixOperator(tokenType TokenType) bool {
	switch tokenType {
	case SET, SETDELAYED, UNSET, OR, AND, EQUAL, UNEQUAL, SAMEQ, UNSAMEQ, LESS, GREATER, LESSEQUAL, GREATEREQUAL, PLUS, MINUS, MULTIPLY, DIVIDE:
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
	right := p.parseInfixExpression(PREFIX)

	return p.createPrefixExpr(operator.Type, right)
}

func (p *Parser) createInfixExpr(operator TokenType, left, right Expr) Expr {
	switch operator {
	case SET:
		return NewList(NewSymbolAtom("Set"), left, right)
	case SETDELAYED:
		return NewList(NewSymbolAtom("SetDelayed"), left, right)
	case UNSET:
		return NewList(NewSymbolAtom("Unset"), left)
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

func ParseString(input string) (Expr, error) {
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	return parser.Parse()
}
