package main

import (
	"flag"
	"fmt"
	"os"
	//	"github.com/client9/sexpr"
)

func main() {
	// Define command line flags
	var (
		prompt = flag.String("prompt", "sexpr> ", "REPL prompt string")
		help   = flag.Bool("help", false, "Show help message")
		//file   = flag.String("file", "", "Execute expressions from file instead of interactive mode")
		cmd = flag.String("c", "", "Execute expression from command line")
		//withUint64 = flag.Bool("with-uint64", false, "Enable experimental Uint64 type system")
	)

	flag.Parse()
	args := flag.Args()
	file := ""
	if len(args) > 0 {
		file = args[0]
	}

	// Show help if requested
	if *help {
		showHelp()
		return
	}

	// Create REPL instance
	repl := NewREPL()
	repl.SetPrompt(*prompt)
	/*
		// Enable Uint64 extension if requested
		if *withUint64 {
			if err := sexpr.RegisterUint64(repl.GetEvaluator().GetContext().GetFunctionRegistry()); err != nil {
				fmt.Fprintf(os.Stderr, "Error enabling Uint64 system: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Uint64 type system enabled. Try: Uint64(42), Uint64(\"#FF\"), Plus(Uint64(10), 5)")
		}
	*/
	// if expression is entered on command line, execute it
	if *cmd != "" {
		if err := repl.ExecuteString(*cmd); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		return
	}

	// If file is specified, execute it
	if file != "" {
		if err := repl.ExecuteFile(file); err != nil {
			fmt.Fprintf(os.Stderr, "Error executing file: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Start interactive REPL
	if err := repl.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "REPL error: %v\n", err)
		os.Exit(1)
	}
}

// showHelp displays help information
func showHelp() {
	fmt.Println(`S-Expression REPL - Symbolic computation system

Usage:
  repl [flags] [file]

Flags:
  -prompt string    Set the REPL prompt (default "sexpr> ")
  -c expression     Evaluate expression and exit
  -help             Show this help message

Examples:
  repl                               # Start interactive REPL
  repl -c 'InputForm(List(1,2,3))'   # Prints [1,2,3]
  repl examples.sexpr                # Execute file and exit

For detailed usage information, start the REPL and type 'help'.`)
}
