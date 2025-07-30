package engine

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/client9/sexpr/core"
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

func (p *Parser) Parse() (core.Expr, error) {
	expr := p.parseExpression()
	if len(p.errors) > 0 {
		return nil, fmt.Errorf("parse errors: %s", strings.Join(p.errors, "; "))
	}
	return expr, nil
}

func (p *Parser) parseExpression() core.Expr {
	return p.parseInfixExpression(PrecedenceLowest)
}

func (p *Parser) parseInfixExpression(precedence Precedence) core.Expr {
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

func (p *Parser) ParseAtom() core.Expr {
	var expr core.Expr

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
	case NOT:
		expr = p.parsePrefixExpression()
	case LPAREN:
		expr = p.parseGroupedExpression()
	case SEMICOLON, EOF:
		// Empty expression - return Null without consuming token
		return core.NewSymbolNull()
	default:
		p.addError(fmt.Sprintf("unexpected token: %s", p.currentToken.String()))
		return nil
	}

	return expr
}

func (p *Parser) parseSymbolOrList() core.Expr {
	symbolToken := p.currentToken
	p.nextToken()

	// Check if this is a pattern: SYMBOL + UNDERSCORE(s) + optional SYMBOL
	if p.currentToken.Type == UNDERSCORE {
		return p.parsePatternFromSymbol(symbolToken.Value)
	}

	if p.currentToken.Type == LPAREN {
		return p.parseList(symbolToken.Value)
	}

	return core.NewSymbol(symbolToken.Value)
}

func (p *Parser) parseList(head string) core.Expr {
	p.nextToken() // consume '('

	elements := []core.Expr{}

	if p.currentToken.Type == RPAREN {
		p.nextToken() // consume ')'
		return core.NewList(head, elements...)
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

	return core.NewList(head, elements...)
}

func (p *Parser) parseListLiteral() core.Expr {
	p.nextToken() // consume '['

	// Create a List expression with "List" as the head
	elements := []core.Expr{}

	// Handle empty list []
	if p.currentToken.Type == RBRACKET {
		p.nextToken() // consume ']'
		return core.NewList("List", elements...)
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

	return core.NewList("List", elements...)
}

func (p *Parser) parseAssociationLiteral() core.Expr {
	p.nextToken() // consume '{'

	// Create rules slice for Rule expressions
	var rules []core.Expr

	// Handle empty association {}
	if p.currentToken.Type == RBRACE {
		p.nextToken() // consume '}'
		// Create Association function call with no arguments for empty association
		return core.NewList("Association")
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
	return core.NewList("Association", rules...)
}

func (p *Parser) parseInteger() core.Expr {
	value, err := strconv.ParseInt(p.currentToken.Value, 10, 64)
	if err != nil {
		p.addError(fmt.Sprintf("invalid integer: %s", p.currentToken.Value))
		return nil
	}

	return core.NewInteger(value)
}

func (p *Parser) parseFloat() core.Expr {
	value, err := strconv.ParseFloat(p.currentToken.Value, 64)
	if err != nil {
		p.addError(fmt.Sprintf("invalid float: %s", p.currentToken.Value))
		return nil
	}

	return core.NewReal(value)
}

func (p *Parser) parseString() core.Expr {
	value := p.unescapeString(p.currentToken.Value)
	return core.NewString(value)
}

func (p *Parser) parseBoolean() core.Expr {
	value := p.currentToken.Value == "True"
	return core.NewBool(value)
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

func (p *Parser) parseInfixOperation(left core.Expr) core.Expr {
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
	return p.createInfixExpr(operator.Type, left, right)
}

func (p *Parser) parsePrefixExpression() core.Expr {
	operator := p.currentToken
	p.nextToken()
	right := p.parseInfixExpression(PrecedenceUnary)

	return p.createPrefixExpr(operator.Type, right)
}

func (p *Parser) createInfixExpr(operator TokenType, left, right core.Expr) core.Expr {
	switch operator {
	case SEMICOLON:
		return core.NewList("CompoundStatement", left, right)
	case SET:
		return core.NewList("Set", left, right)
	case SETDELAYED:
		return core.NewList("SetDelayed", left, right)
	case UNSET:
		return core.NewList("Unset", left)
	case COLON:
		return core.NewList("Rule", left, right)
	case RULEDELAYED:
		return core.NewList("RuleDelayed", left, right)
	case OR:
		return core.NewList("Or", left, right)
	case AND:
		return core.NewList("And", left, right)
	case EQUAL:
		return core.NewList("Equal", left, right)
	case UNEQUAL:
		return core.NewList("Unequal", left, right)
	case SAMEQ:
		return core.NewList("SameQ", left, right)
	case UNSAMEQ:
		return core.NewList("UnsameQ", left, right)
	case LESS:
		return core.NewList("Less", left, right)
	case GREATER:
		return core.NewList("Greater", left, right)
	case LESSEQUAL:
		return core.NewList("LessEqual", left, right)
	case GREATEREQUAL:
		return core.NewList("GreaterEqual", left, right)
	case PLUS:
		return core.NewList("Plus", left, right)
	case MINUS:
		return core.NewList("Subtract", left, right)
	case MULTIPLY:
		return core.NewList("Times", left, right)
	case DIVIDE:
		return core.NewList("Divide", left, right)
	case CARET:
		return core.NewList("Power", left, right)
	default:
		p.addError(fmt.Sprintf("unknown infix operator: %d", operator))
		return nil
	}
}

func (p *Parser) createPrefixExpr(operator TokenType, operand core.Expr) core.Expr {
	switch operator {
	case MINUS:
		return core.NewList("Minus", operand)
	case PLUS:
		return operand // unary plus is identity
	case NOT:
		return core.NewList("Not", operand)
	default:
		p.addError(fmt.Sprintf("unknown prefix operator: %d", operator))
		return nil
	}
}

func (p *Parser) parseGroupedExpression() core.Expr {
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
func (p *Parser) parseIndexOrSlice(expr core.Expr) core.Expr {
	p.nextToken() // consume '['

	// Check for empty brackets []
	if p.currentToken.Type == RBRACKET {
		p.addError("empty brackets are not allowed")
		return expr
	}

	// Parse the first expression (could be index, start, or just ':')
	var firstExpr core.Expr
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
		return core.NewList("Part", expr, firstExpr)

	} else if p.currentToken.Type == COLON {
		// Slice syntax: expr[start:end] or expr[:end] or expr[start:]
		p.nextToken() // consume ':'

		var startExpr, endExpr core.Expr

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
			return core.NewList("TakeFrom", expr, startExpr)
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
				return core.NewList("Take", expr, endExpr)
			} else {
				// [start:end] syntax - Slice operation
				return core.NewList("SliceRange", expr, startExpr, endExpr)
			}
		}
	} else {
		p.addError("expected ':' or ']' after slice expression")
		return expr
	}
}

// parseFunctionApplication handles postfix function application like Function(x, x+1)(5)
func (p *Parser) parseFunctionApplication(expr core.Expr) core.Expr {
	p.nextToken() // consume '('

	var args []core.Expr

	// Handle empty argument list: f()
	if p.currentToken.Type == RPAREN {
		p.nextToken()                                 // consume ')'
		return core.List{Elements: []core.Expr{expr}} // Just the function with no arguments
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
	elements := make([]core.Expr, len(args)+1)
	elements[0] = expr
	copy(elements[1:], args)

	return core.List{Elements: elements}
}

// parseSliceExpression parses expressions inside slice brackets, treating colons as separators
func (p *Parser) parseSliceExpression() core.Expr {
	// Parse a simple expression that stops at colons and brackets
	// We'll use a custom precedence that's higher than colon assignment
	return p.parseInfixExpression(PrecedenceLogicalOr)
}

// isSliceExpression checks if an expression is a slice or part expression
func (p *Parser) isSliceExpression(expr core.Expr) bool {
	list, ok := expr.(core.List)
	if !ok || len(list.Elements) == 0 {
		return false
	}

	headName, ok := core.ExtractSymbol(list.Elements[0])
	if !ok {
		return false
	}
	return headName == "Part" || headName == "SliceRange" || headName == "Take" || headName == "TakeFrom"
}

// createSliceAssignment creates the appropriate slice assignment AST node
func (p *Parser) createSliceAssignment(sliceExpr core.Expr, value core.Expr) core.Expr {
	list := sliceExpr.(core.List)
	headName, ok := core.ExtractSymbol(list.Elements[0])
	if !ok {
		p.addError(fmt.Sprintf("Unknown slice expression type: %v", list.Elements[0]))
		return nil
	}

	switch headName {
	case "Part":
		// Part(expr, index) = value -> PartSet(expr, index, value)
		if len(list.Elements) != 3 {
			p.addError("Part expression must have exactly 2 arguments for assignment")
			return nil
		}
		return core.NewList("PartSet", list.Elements[1], list.Elements[2], value)

	case "SliceRange":
		// SliceRange(expr, start, end) = value -> SliceSet(expr, start, end, value)
		if len(list.Elements) != 4 {
			p.addError("SliceRange expression must have exactly 3 arguments for assignment")
			return nil
		}
		return core.NewList("SliceSet", list.Elements[1], list.Elements[2], list.Elements[3], value)

	case "Take":
		// Take(expr, n) = value -> SliceSet(expr, 1, n, value)
		if len(list.Elements) != 3 {
			p.addError("Take expression must have exactly 2 arguments for assignment")
			return nil
		}
		return core.NewList("SliceSet", list.Elements[1], core.NewInteger(1), list.Elements[2], value)

	case "TakeFrom":
		// TakeFrom(expr, start) = value -> SliceSet(expr, start, -1, value)
		// Note: -1 represents "to end" in our slice assignment semantics
		if len(list.Elements) != 3 {
			p.addError("TakeFrom expression must have exactly 2 arguments for assignment")
			return nil
		}
		return core.NewList("SliceSet", list.Elements[1], list.Elements[2], core.NewInteger(-1), value)

	default:
		p.addError(fmt.Sprintf("Unknown slice expression type: %s", headName))
		return nil
	}
}

// parseUnderscorePattern parses anonymous patterns (_, __, ___, _Integer, __Integer, ___Integer)
func (p *Parser) parseUnderscorePattern() core.Expr {
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
	var blankExpr core.Expr
	if underscoreCount >= 3 {
		if typeName != "" {
			blankExpr = core.NewList("BlankNullSequence", core.NewSymbol(typeName))
		} else {
			blankExpr = core.NewList("BlankNullSequence")
		}
	} else if underscoreCount == 2 {
		if typeName != "" {
			blankExpr = core.NewList("BlankSequence", core.NewSymbol(typeName))
		} else {
			blankExpr = core.NewList("BlankSequence")
		}
	} else {
		if typeName != "" {
			blankExpr = core.NewList("Blank", core.NewSymbol(typeName))
		} else {
			blankExpr = core.NewList("Blank")
		}
	}

	// Anonymous pattern - just return the blank expression
	return blankExpr
}

// parsePatternFromSymbol parses named patterns (x_, x__, x___, x_Integer, x__Integer, x___Integer)
func (p *Parser) parsePatternFromSymbol(varName string) core.Expr {
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
	var blankExpr core.Expr
	if underscoreCount >= 3 {
		if typeName != "" {
			blankExpr = core.NewList("BlankNullSequence", core.NewSymbol(typeName))
		} else {
			blankExpr = core.NewList("BlankNullSequence")
		}
	} else if underscoreCount == 2 {
		if typeName != "" {
			blankExpr = core.NewList("BlankSequence", core.NewSymbol(typeName))
		} else {
			blankExpr = core.NewList("BlankSequence")
		}
	} else {
		if typeName != "" {
			blankExpr = core.NewList("Blank", core.NewSymbol(typeName))
		} else {
			blankExpr = core.NewList("Blank")
		}
	}

	// Named pattern - wrap in Pattern(varName, blankExpr)
	return core.NewList("Pattern", core.NewSymbol(varName), blankExpr)
}

// parseFunctionShorthand handles the & postfix operator: expr & -> Function(expr)
func (p *Parser) parseFunctionShorthand(expr core.Expr) core.Expr {
	p.nextToken() // consume '&'
	return core.NewList("Function", expr)
}

func ParseString(input string) (core.Expr, error) {
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	return parser.Parse()
}
