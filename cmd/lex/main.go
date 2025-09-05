package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/client9/cardinal/core"
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
	lexer := core.NewLexer(input)

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
		typeName := token.String() //getTokenTypeName(token.Type)

		// Format value for display
		displayValue := formatTokenValue(token)

		// Print token information
		fmt.Printf("%-4d %-15s %-20s %s\n",
			actualPos,
			typeName,
			fmt.Sprintf("\"%s\"", token.Value),
			displayValue)

		if token.Type == core.EOF {
			break
		}

		position = actualPos + len(token.Value)
	}
}

// formatTokenValue formats the token value for display
func formatTokenValue(token core.Token) string {
	return token.String()
}
