package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/client9/sexpr"
)

func main() {
	if len(os.Args) > 1 {
		// Process command line argument
		input := strings.Join(os.Args[1:], " ")
		processInput(input)
	} else {
		// Interactive mode
		fmt.Println("Lexer Debug Tool")
		fmt.Println("Enter expressions to tokenize (Ctrl+C to exit):")
		fmt.Println()

		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("lex> ")
			if !scanner.Scan() {
				break
			}
			
			input := strings.TrimSpace(scanner.Text())
			if input == "" {
				continue
			}
			
			if input == "exit" || input == "quit" {
				break
			}
			
			processInput(input)
			fmt.Println()
		}
	}
}

func processInput(input string) {
	lexer := sexpr.NewLexer(input)
	
	fmt.Printf("Input: %s\n", input)
	fmt.Println("Tokens:")
	fmt.Println("-------")
	fmt.Printf("%-4s %-15s %-20s %s\n", "Pos", "Type", "Value", "Display")
	fmt.Println(strings.Repeat("-", 60))
	
	position := 0
	for {
		token := lexer.NextToken()
		
		// Calculate actual position in input
		actualPos := token.Position - 1
		if actualPos < 0 {
			actualPos = position
		}
		
		// Format token type name
		typeName := getTokenTypeName(token.Type)
		
		// Format value for display
		displayValue := formatTokenValue(token)
		
		// Print token information
		fmt.Printf("%-4d %-15s %-20s %s\n", 
			actualPos, 
			typeName, 
			fmt.Sprintf("\"%s\"", token.Value), 
			displayValue)
		
		if token.Type == sexpr.EOF {
			break
		}
		
		position = actualPos + len(token.Value)
	}
}

// getTokenTypeName returns a human-readable name for the token type
func getTokenTypeName(tokenType sexpr.TokenType) string {
	switch tokenType {
	case sexpr.EOF:
		return "EOF"
	case sexpr.SYMBOL:
		return "SYMBOL"
	case sexpr.INTEGER:
		return "INTEGER"
	case sexpr.FLOAT:
		return "FLOAT"
	case sexpr.STRING:
		return "STRING"
	case sexpr.BOOLEAN:
		return "BOOLEAN"
	case sexpr.LBRACKET:
		return "LBRACKET"
	case sexpr.RBRACKET:
		return "RBRACKET"
	case sexpr.LBRACE:
		return "LBRACE"
	case sexpr.RBRACE:
		return "RBRACE"
	case sexpr.COMMA:
		return "COMMA"
	case sexpr.COLON:
		return "COLON"
	case sexpr.PLUS:
		return "PLUS"
	case sexpr.MINUS:
		return "MINUS"
	case sexpr.MULTIPLY:
		return "MULTIPLY"
	case sexpr.DIVIDE:
		return "DIVIDE"
	case sexpr.LPAREN:
		return "LPAREN"
	case sexpr.RPAREN:
		return "RPAREN"
	case sexpr.SET:
		return "SET"
	case sexpr.SETDELAYED:
		return "SETDELAYED"
	case sexpr.UNSET:
		return "UNSET"
	case sexpr.EQUAL:
		return "EQUAL"
	case sexpr.UNEQUAL:
		return "UNEQUAL"
	case sexpr.LESS:
		return "LESS"
	case sexpr.GREATER:
		return "GREATER"
	case sexpr.LESSEQUAL:
		return "LESSEQUAL"
	case sexpr.GREATEREQUAL:
		return "GREATEREQUAL"
	case sexpr.AND:
		return "AND"
	case sexpr.OR:
		return "OR"
	case sexpr.NOT:
		return "NOT"
	case sexpr.SAMEQ:
		return "SAMEQ"
	case sexpr.UNSAMEQ:
		return "UNSAMEQ"
	case sexpr.CARET:
		return "CARET"
	case sexpr.SEMICOLON:
		return "SEMICOLON"
	case sexpr.UNDERSCORE:
		return "UNDERSCORE"
	case sexpr.WHITESPACE:
		return "WHITESPACE"
	case sexpr.ILLEGAL:
		return "ILLEGAL"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", int(tokenType))
	}
}

// formatTokenValue formats the token value for display
func formatTokenValue(token sexpr.Token) string {
	switch token.Type {
	case sexpr.EOF:
		return "end of input"
	case sexpr.SYMBOL:
		return fmt.Sprintf("symbol: %s", token.Value)
	case sexpr.INTEGER:
		return fmt.Sprintf("integer: %s", token.Value)
	case sexpr.FLOAT:
		return fmt.Sprintf("float: %s", token.Value)
	case sexpr.STRING:
		return fmt.Sprintf("string: %s", token.Value)
	case sexpr.BOOLEAN:
		return fmt.Sprintf("boolean: %s", token.Value)
	case sexpr.LBRACKET:
		return "["
	case sexpr.RBRACKET:
		return "]"
	case sexpr.LBRACE:
		return "{"
	case sexpr.RBRACE:
		return "}"
	case sexpr.COMMA:
		return ","
	case sexpr.COLON:
		return ":"
	case sexpr.PLUS:
		return "+"
	case sexpr.MINUS:
		return "-"
	case sexpr.MULTIPLY:
		return "*"
	case sexpr.DIVIDE:
		return "/"
	case sexpr.LPAREN:
		return "("
	case sexpr.RPAREN:
		return ")"
	case sexpr.SET:
		return "= (assignment)"
	case sexpr.SETDELAYED:
		return ":= (delayed assignment)"
	case sexpr.UNSET:
		return "=. (unset)"
	case sexpr.EQUAL:
		return "== (equal)"
	case sexpr.UNEQUAL:
		return "!= (not equal)"
	case sexpr.LESS:
		return "< (less than)"
	case sexpr.GREATER:
		return "> (greater than)"
	case sexpr.LESSEQUAL:
		return "<= (less equal)"
	case sexpr.GREATEREQUAL:
		return ">= (greater equal)"
	case sexpr.AND:
		return "&& (and)"
	case sexpr.OR:
		return "|| (or)"
	case sexpr.NOT:
		return "! (not)"
	case sexpr.SAMEQ:
		return "=== (same)"
	case sexpr.UNSAMEQ:
		return "=!= (not same)"
	case sexpr.CARET:
		return "^ (power)"
	case sexpr.SEMICOLON:
		return "; (semicolon)"
	case sexpr.UNDERSCORE:
		return "_ (pattern)"
	case sexpr.WHITESPACE:
		return "whitespace"
	case sexpr.ILLEGAL:
		return fmt.Sprintf("ILLEGAL: %s", token.Value)
	default:
		return token.Value
	}
}