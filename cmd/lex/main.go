package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/client9/sexpr/engine"
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
	lexer := engine.NewLexer(input)
	
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
		
		if token.Type == engine.EOF {
			break
		}
		
		position = actualPos + len(token.Value)
	}
}

// getTokenTypeName returns a human-readable name for the token type
func getTokenTypeName(tokenType engine.TokenType) string {
	switch tokenType {
	case engine.EOF:
		return "EOF"
	case engine.SYMBOL:
		return "SYMBOL"
	case engine.INTEGER:
		return "INTEGER"
	case engine.FLOAT:
		return "FLOAT"
	case engine.STRING:
		return "STRING"
	case engine.BOOLEAN:
		return "BOOLEAN"
	case engine.LBRACKET:
		return "LBRACKET"
	case engine.RBRACKET:
		return "RBRACKET"
	case engine.LBRACE:
		return "LBRACE"
	case engine.RBRACE:
		return "RBRACE"
	case engine.COMMA:
		return "COMMA"
	case engine.COLON:
		return "COLON"
	case engine.PLUS:
		return "PLUS"
	case engine.MINUS:
		return "MINUS"
	case engine.MULTIPLY:
		return "MULTIPLY"
	case engine.DIVIDE:
		return "DIVIDE"
	case engine.LPAREN:
		return "LPAREN"
	case engine.RPAREN:
		return "RPAREN"
	case engine.SET:
		return "SET"
	case engine.SETDELAYED:
		return "SETDELAYED"
	case engine.UNSET:
		return "UNSET"
	case engine.EQUAL:
		return "EQUAL"
	case engine.UNEQUAL:
		return "UNEQUAL"
	case engine.LESS:
		return "LESS"
	case engine.GREATER:
		return "GREATER"
	case engine.LESSEQUAL:
		return "LESSEQUAL"
	case engine.GREATEREQUAL:
		return "GREATEREQUAL"
	case engine.AND:
		return "AND"
	case engine.OR:
		return "OR"
	case engine.NOT:
		return "NOT"
	case engine.SAMEQ:
		return "SAMEQ"
	case engine.UNSAMEQ:
		return "UNSAMEQ"
	case engine.CARET:
		return "CARET"
	case engine.SEMICOLON:
		return "SEMICOLON"
	case engine.UNDERSCORE:
		return "UNDERSCORE"
	case engine.WHITESPACE:
		return "WHITESPACE"
	case engine.ILLEGAL:
		return "ILLEGAL"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", int(tokenType))
	}
}

// formatTokenValue formats the token value for display
func formatTokenValue(token engine.Token) string {
	switch token.Type {
	case engine.EOF:
		return "end of input"
	case engine.SYMBOL:
		return fmt.Sprintf("symbol: %s", token.Value)
	case engine.INTEGER:
		return fmt.Sprintf("integer: %s", token.Value)
	case engine.FLOAT:
		return fmt.Sprintf("float: %s", token.Value)
	case engine.STRING:
		return fmt.Sprintf("string: %s", token.Value)
	case engine.BOOLEAN:
		return fmt.Sprintf("boolean: %s", token.Value)
	case engine.LBRACKET:
		return "["
	case engine.RBRACKET:
		return "]"
	case engine.LBRACE:
		return "{"
	case engine.RBRACE:
		return "}"
	case engine.COMMA:
		return ","
	case engine.COLON:
		return ":"
	case engine.PLUS:
		return "+"
	case engine.MINUS:
		return "-"
	case engine.MULTIPLY:
		return "*"
	case engine.DIVIDE:
		return "/"
	case engine.LPAREN:
		return "("
	case engine.RPAREN:
		return ")"
	case engine.SET:
		return "= (assignment)"
	case engine.SETDELAYED:
		return ":= (delayed assignment)"
	case engine.UNSET:
		return "=. (unset)"
	case engine.EQUAL:
		return "== (equal)"
	case engine.UNEQUAL:
		return "!= (not equal)"
	case engine.LESS:
		return "< (less than)"
	case engine.GREATER:
		return "> (greater than)"
	case engine.LESSEQUAL:
		return "<= (less equal)"
	case engine.GREATEREQUAL:
		return ">= (greater equal)"
	case engine.AND:
		return "&& (and)"
	case engine.OR:
		return "|| (or)"
	case engine.NOT:
		return "! (not)"
	case engine.SAMEQ:
		return "=== (same)"
	case engine.UNSAMEQ:
		return "=!= (not same)"
	case engine.CARET:
		return "^ (power)"
	case engine.SEMICOLON:
		return "; (semicolon)"
	case engine.UNDERSCORE:
		return "_ (pattern)"
	case engine.WHITESPACE:
		return "whitespace"
	case engine.ILLEGAL:
		return fmt.Sprintf("ILLEGAL: %s", token.Value)
	default:
		return token.Value
	}
}
