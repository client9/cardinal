package sexpr

import (
	"fmt"
)

type TokenType int

const (
	EOF TokenType = iota
	SYMBOL
	INTEGER
	FLOAT
	STRING
	BOOLEAN
	LBRACKET
	RBRACKET
	LBRACE
	RBRACE
	COMMA
	COLON
	PLUS
	MINUS
	MULTIPLY
	DIVIDE
	LPAREN
	RPAREN
	SET
	SETDELAYED
	UNSET
	EQUAL
	UNEQUAL
	LESS
	GREATER
	LESSEQUAL
	GREATEREQUAL
	AND
	OR
	NOT
	SAMEQ
	UNSAMEQ
	CARET
	SEMICOLON
	UNDERSCORE // _
	WHITESPACE
	ILLEGAL
)

type Token struct {
	Type     TokenType
	Value    string
	Position int
}

func (t Token) String() string {
	switch t.Type {
	case EOF:
		return "EOF"
	case SYMBOL:
		return fmt.Sprintf("SYMBOL(%s)", t.Value)
	case INTEGER:
		return fmt.Sprintf("INTEGER(%s)", t.Value)
	case FLOAT:
		return fmt.Sprintf("FLOAT(%s)", t.Value)
	case STRING:
		return fmt.Sprintf("STRING(%s)", t.Value)
	case BOOLEAN:
		return fmt.Sprintf("BOOLEAN(%s)", t.Value)
	case LBRACKET:
		return "LBRACKET"
	case RBRACKET:
		return "RBRACKET"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case COMMA:
		return "COMMA"
	case COLON:
		return "COLON"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case MULTIPLY:
		return "MULTIPLY"
	case DIVIDE:
		return "DIVIDE"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case SET:
		return "SET"
	case SETDELAYED:
		return "SETDELAYED"
	case UNSET:
		return "UNSET"
	case EQUAL:
		return "EQUAL"
	case UNEQUAL:
		return "UNEQUAL"
	case LESS:
		return "LESS"
	case GREATER:
		return "GREATER"
	case LESSEQUAL:
		return "LESSEQUAL"
	case GREATEREQUAL:
		return "GREATEREQUAL"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case NOT:
		return "NOT"
	case SAMEQ:
		return "SAMEQ"
	case UNSAMEQ:
		return "UNSAMEQ"
	case CARET:
		return "CARET"
	case UNDERSCORE:
		return "UNDERSCORE"
	case WHITESPACE:
		return "WHITESPACE"
	case ILLEGAL:
		return fmt.Sprintf("ILLEGAL(%s)", t.Value)
	default:
		return "UNKNOWN"
	}
}

type Lexer struct {
	input    string
	position int
	ch       byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:    input,
		position: 0,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.position >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.position]
	}
	l.position++
}

func (l *Lexer) peekChar() byte {
	if l.position >= len(l.input) {
		return 0
	}
	return l.input[l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	// Skip from # to end of line
	for l.ch != '\n' && l.ch != '\r' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) readString() string {
	position := l.position
	l.readChar() // skip opening quote

	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar() // skip escape character
			if l.ch != 0 {
				l.readChar() // skip escaped character
			}
		} else {
			l.readChar()
		}
	}

	if l.ch == '"' {
		result := l.input[position : l.position-1]
		l.readChar() // skip closing quote
		return result
	}

	// Handle unclosed string - return what we have
	if l.position > len(l.input) {
		return l.input[position:]
	}
	return l.input[position:l.position]
}

func (l *Lexer) readIdentifier() string {
	position := l.position - 1
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position : l.position-1]
}

func (l *Lexer) readUnderscores() string {
	position := l.position - 1
	for l.ch == '_' {
		l.readChar()
	}
	return l.input[position : l.position-1]
}

func (l *Lexer) readNumber() (string, TokenType) {
	position := l.position - 1
	tokenType := INTEGER

	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '.' && isDigit(l.peekChar()) {
		tokenType = FLOAT
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position : l.position-1], tokenType
}

func (l *Lexer) NextToken() Token {
	var tok Token

	// Skip whitespace and comments
	for {
		l.skipWhitespace()
		if l.ch == '#' {
			l.skipComment()
			// Continue to skip any additional whitespace after comment
			continue
		}
		break
	}

	switch l.ch {
	case '[':
		tok = Token{Type: LBRACKET, Value: string(l.ch), Position: l.position - 1}
	case ']':
		tok = Token{Type: RBRACKET, Value: string(l.ch), Position: l.position - 1}
	case '{':
		tok = Token{Type: LBRACE, Value: string(l.ch), Position: l.position - 1}
	case '}':
		tok = Token{Type: RBRACE, Value: string(l.ch), Position: l.position - 1}
	case ',':
		tok = Token{Type: COMMA, Value: string(l.ch), Position: l.position - 1}
	case ';':
		tok = Token{Type: SEMICOLON, Value: string(l.ch), Position: l.position - 1}
	case '+':
		tok = Token{Type: PLUS, Value: string(l.ch), Position: l.position - 1}
	case '-':
		tok = Token{Type: MINUS, Value: string(l.ch), Position: l.position - 1}
	case '*':
		tok = Token{Type: MULTIPLY, Value: string(l.ch), Position: l.position - 1}
	case '/':
		tok = Token{Type: DIVIDE, Value: string(l.ch), Position: l.position - 1}
	case '^':
		tok = Token{Type: CARET, Value: string(l.ch), Position: l.position - 1}
	case '(':
		tok = Token{Type: LPAREN, Value: string(l.ch), Position: l.position - 1}
	case ')':
		tok = Token{Type: RPAREN, Value: string(l.ch), Position: l.position - 1}
	case '=':
		position := l.position - 1
		if l.peekChar() == '!' {
			l.readChar() // move to '!'
			if l.peekChar() == '=' {
				l.readChar() // move to final '='
				l.readChar() // move past final '='
				tok = Token{Type: UNSAMEQ, Value: "=!=", Position: position}
				return tok
			} else {
				tok = Token{Type: ILLEGAL, Value: "=!", Position: position}
				return tok
			}
		} else if l.peekChar() == '=' {
			l.readChar() // move to second '='
			if l.peekChar() == '=' {
				l.readChar() // move to third '='
				l.readChar() // move past third '='
				tok = Token{Type: SAMEQ, Value: "===", Position: position}
				return tok
			} else {
				l.readChar() // move past second '='
				tok = Token{Type: EQUAL, Value: "==", Position: position}
				return tok
			}
		} else if l.peekChar() == '.' {
			l.readChar() // move to '.'
			l.readChar() // move past '.'
			tok = Token{Type: UNSET, Value: "=.", Position: position}
			return tok
		} else {
			tok = Token{Type: SET, Value: "=", Position: position}
		}
	case ':':
		if l.peekChar() == '=' {
			tok = Token{Type: SETDELAYED, Value: ":=", Position: l.position - 1}
			l.readChar() // consume ':'
			l.readChar() // consume '='
			return tok
		} else {
			tok = Token{Type: COLON, Value: string(l.ch), Position: l.position - 1}
		}
	case '!':
		if l.peekChar() == '=' {
			tok = Token{Type: UNEQUAL, Value: "!=", Position: l.position - 1}
			l.readChar() // consume '!'
			l.readChar() // consume '='
			return tok
		} else {
			tok = Token{Type: NOT, Value: string(l.ch), Position: l.position - 1}
		}
	case '<':
		if l.peekChar() == '=' {
			tok = Token{Type: LESSEQUAL, Value: "<=", Position: l.position - 1}
			l.readChar() // consume '<'
			l.readChar() // consume '='
			return tok
		} else {
			tok = Token{Type: LESS, Value: string(l.ch), Position: l.position - 1}
		}
	case '>':
		if l.peekChar() == '=' {
			tok = Token{Type: GREATEREQUAL, Value: ">=", Position: l.position - 1}
			l.readChar() // consume '>'
			l.readChar() // consume '='
			return tok
		} else {
			tok = Token{Type: GREATER, Value: string(l.ch), Position: l.position - 1}
		}
	case '&':
		if l.peekChar() == '&' {
			tok = Token{Type: AND, Value: "&&", Position: l.position - 1}
			l.readChar() // consume first '&'
			l.readChar() // consume second '&'
			return tok
		} else {
			tok = Token{Type: ILLEGAL, Value: string(l.ch), Position: l.position - 1}
		}
	case '|':
		if l.peekChar() == '|' {
			tok = Token{Type: OR, Value: "||", Position: l.position - 1}
			l.readChar() // consume first '|'
			l.readChar() // consume second '|'
			return tok
		} else {
			tok = Token{Type: ILLEGAL, Value: string(l.ch), Position: l.position - 1}
		}
	case '"':
		tok.Type = STRING
		tok.Value = l.readString()
		tok.Position = l.position - len(tok.Value) - 2
		return tok
	case '_':
		tok.Position = l.position - 1
		underscoreValue := l.readUnderscores()
		if len(underscoreValue) > 3 {
			tok = Token{Type: ILLEGAL, Value: underscoreValue, Position: tok.Position}
		} else {
			tok = Token{Type: UNDERSCORE, Value: underscoreValue, Position: tok.Position}
		}
		return tok
	case 0:
		tok.Type = EOF
		tok.Value = ""
		tok.Position = l.position
		return tok
	default:
		if isLetter(l.ch) {
			tok.Position = l.position - 1
			tok.Value = l.readIdentifier()
			tok.Type = SYMBOL
			return tok
		} else if isDigit(l.ch) {
			tok.Position = l.position - 1
			tok.Value, tok.Type = l.readNumber()
			return tok
		} else {
			tok = Token{Type: ILLEGAL, Value: string(l.ch), Position: l.position - 1}
		}
	}

	l.readChar()
	return tok
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) Tokenize() []Token {
	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == EOF {
			break
		}
	}
	return tokens
}
